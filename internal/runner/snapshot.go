/*
copyright 2020 the Escli authors

licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at

    http://www.apache.org/licenses/license-2.0

unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license.
*/

package runner

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"

	"github.com/DevopsArtFactory/escli/internal/constants"
	snapshotSchema "github.com/DevopsArtFactory/escli/internal/schema/snapshot"
	"github.com/DevopsArtFactory/escli/internal/util"
)

type SnapshotSegment struct {
	Key          string
	StorageClass string
}

type SegmentError struct {
	Key   string
	Error error
}

func (s *SegmentError) PrintError() {
	color.Red("Key: %s, error: %s", s.Key, s.Error.Error())
}

func (r Runner) ListSnapshot(out io.Writer) error {
	repositories, err := r.Client.GetRepositories()
	if err != nil {
		return err
	}

	for repositoryID, repositorySetting := range repositories {
		if r.Flag.WithRepo == repositoryID || r.Flag.WithRepo == constants.EmptyString {
			fmt.Fprintf(out, "Repository ID : %s\n", util.StringWithColor(repositoryID))
			fmt.Fprintf(out, "Type : %s\n", repositorySetting.Type)
		}
		if r.Flag.RepoOnly {
			prettyJSON, err := json.MarshalIndent(repositorySetting.Settings, " ", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "Setting : %s\n\n", string(prettyJSON))
		} else if r.Flag.WithRepo == repositoryID || r.Flag.WithRepo == constants.EmptyString {
			snapshotsMetadata, err := r.Client.GetSnapshots(repositoryID)
			if err != nil {
				continue
			}
			fmt.Fprintf(out, "%-50s\t%s\n",
				"ID",
				"state")
			for _, snapshot := range snapshotsMetadata.Snapshots {
				fmt.Fprintf(out, "%-50s\t%s\n",
					snapshot.Snapshot,
					util.StringWithColor(snapshot.State))
				fmt.Fprintf(out, "%s\n",
					snapshot.Indices)
			}
		}
	}

	return nil
}

func (r Runner) CreateSnapshot(out io.Writer, args []string) error {
	repositoryID := args[0]
	snapshotID := args[1]
	indices := args[2]

	requestBody, err := util.JSONtoPrettyString(
		snapshotSchema.RequestBody{
			Indices: indices,
		})

	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s\n", util.YellowString(requestBody))

	if !r.Flag.Force {
		if err := util.AskContinue("Are you sure to create snapshot"); err != nil {
			return errors.New("task has benn canceled")
		}
	}

	statusCode, err := r.Client.CreateSnapshot(repositoryID, snapshotID, requestBody)
	switch statusCode {
	case 200:
		fmt.Fprintf(out, "%s\n", util.GreenString(snapshotID+" is created"))
	default:
		fmt.Fprintf(out, "%s\n", util.RedString(snapshotID+" is not created"))
	}

	return err
}

func (r Runner) DeleteSnapshot(out io.Writer, args []string) error {
	repositoryID := args[0]
	snapshotID := args[1]

	if !r.Flag.Force {
		if err := util.AskContinue("Are you sure to delete snapshot"); err != nil {
			return errors.New("task has benn canceled")
		}
	}

	statusCode, err := r.Client.DeleteSnapshot(repositoryID, snapshotID)
	switch statusCode {
	case 200:
		fmt.Fprintf(out, "%s\n", util.GreenString(snapshotID+" is deleted"))
	default:
		fmt.Fprintf(out, "%s\n", util.RedString(snapshotID+" is not deleted"))
	}

	return err
}

func (r Runner) ArchiveSnapshot(out io.Writer, args []string) error {
	repositoryID := args[0]
	snapshotID := args[1]

	var wait sync.WaitGroup
	runtime.GOMAXPROCS(4)

	repository := r.Client.GetRepository(repositoryID)

	if repository.Type != "s3" {
		return errors.New("archive is only supported s3 type repository")
	}

	if repository.Settings.Bucket == constants.EmptyString {
		return errors.New("archive is only supported normal s3 type repository. this repository doesn't have settings information")
	}

	fmt.Fprintf(out, "bucket name : %s\n", util.StringWithColor(repository.Settings.Bucket))
	fmt.Fprintf(out, "base path : %s\n", util.StringWithColor(repository.Settings.BasePath))

	snapshot := r.Client.GetSnapshot(repositoryID, snapshotID)

	var basePath string

	if repository.Settings.BasePath == constants.EmptyString {
		basePath = constants.EmptyString
	} else {
		basePath = repository.Settings.BasePath + "/"
	}

	snapshotsIndicesS3, err := r.getIndexIDFromS3(
		aws.String(repository.Settings.Bucket),
		aws.String(basePath),
	)

	if err != nil {
		return err
	}

	for _, indexName := range snapshot.Indices {
		fmt.Fprintf(out, "index name : %s\n", util.StringWithColor(indexName))
		metaData := snapshotsIndicesS3.Indices[indexName]
		prefix := basePath + "indices/" + metaData.ID + "/"

		objs, _ := r.Client.GetObjects(aws.String(repository.Settings.Bucket), aws.String(prefix), nil, nil)

		for {
			for _, item := range objs.Contents {
				fmt.Fprintf(out, "%s\n", *item.Key)
				if *item.StorageClass != "GLACIER" {
					if r.Flag.Force {
						color.Green("Change Storage Class to %s -> GLACIER", *item.StorageClass)
						wait.Add(1)
						go func(key string) {
							defer wait.Done()
							_, err := r.Client.TransitObject(aws.String(repository.Settings.Bucket), aws.String(key), "GLACIER")
							if err != nil {
								panic(err)
							}
						}(*item.Key)
					} else {
						if err := util.AskContinue("Change Storage Class to GLACIER "); err != nil {
							color.Red("Don't change storage class %s", *item.Key)
							return errors.New("task has been cancelled")
						}

						color.Green("Change Storage Class to %s -> GLACIER", *item.StorageClass)
						wait.Add(1)
						_, err := r.Client.TransitObject(aws.String(repository.Settings.Bucket), item.Key, "GLACIER")
						if err != nil {
							panic(err)
						}
					}
				}
			}
			wait.Wait()

			if *objs.IsTruncated {
				objs, _ = r.Client.GetObjects(aws.String(repository.Settings.Bucket), aws.String(prefix), nil, objs.NextContinuationToken)
			} else {
				break
			}
		}
	}
	return nil
}

func (r Runner) getIndexIDFromS3(bucket *string, prefix *string) (*snapshotSchema.SnapshotsIndicesS3, error) {
	var snapshotsIndicesS3 snapshotSchema.SnapshotsIndicesS3
	resp, err := r.Client.GetObjects(bucket, prefix, aws.String("/"), nil)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("[0-9]+")

	var snapshotMetadataKey string
	keyNumber := 0

	for _, item := range resp.Contents {
		if strings.Contains(*item.Key, "index-") {
			currentKeyNumber := util.ParseInt(re.FindString(*item.Key))
			if currentKeyNumber >= keyNumber {
				keyNumber = currentKeyNumber
				snapshotMetadataKey = *item.Key
			}
		}
	}

	downloadFileName := r.Client.DownloadFileFromS3(bucket, &snapshotMetadataKey)

	jsonFile, err := os.Open(downloadFileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	util.ConvertJSONtoMetadata(jsonFile, &snapshotsIndicesS3)

	return &snapshotsIndicesS3, nil
}

func (r Runner) RestoreSnapshot(out io.Writer, args []string) error {
	maxSemaphore := r.Flag.MaxConcurrentJob
	if maxSemaphore == 0 {
		maxSemaphore = constants.DefaultMaxConcurrentJob
	}

	repositoryID := args[0]
	snapshotID := args[1]
	indexName := args[2]
	areAllObjectsStandard := true

	repository := r.Client.GetRepository(repositoryID)

	if repository.Type == "s3" {
		fmt.Fprintf(out, "bucket name : %s\n", util.StringWithColor(repository.Settings.Bucket))
		fmt.Fprintf(out, "base path : %s\n", util.StringWithColor(repository.Settings.BasePath))

		var basePath string
		var result []SegmentError

		if repository.Settings.BasePath == constants.EmptyString {
			basePath = constants.EmptyString
		} else {
			basePath = repository.Settings.BasePath + "/"
		}

		snapshotsIndicesS3, err := r.getIndexIDFromS3(
			aws.String(repository.Settings.Bucket),
			aws.String(basePath),
		)

		if err != nil {
			return err
		}

		metaData := snapshotsIndicesS3.Indices[indexName]
		prefix := repository.Settings.BasePath + "/indices/" + metaData.ID + "/"

		segments := r.AddObjectSegments(repository.Settings.Bucket, prefix, nil)
		bar := pb.New(len(segments))
		bar.SetRefreshRate(time.Second)
		bar.SetWriter(out)

		var wg sync.WaitGroup
		semaphore := make(chan int, maxSemaphore)
		output := make(chan []SegmentError)
		input := make(chan SegmentError)
		defer close(output)

		go func(input chan SegmentError, output chan []SegmentError, wg *sync.WaitGroup, bar *pb.ProgressBar) {
			var ret []SegmentError
			for se := range input {
				ret = append(ret, se)
				bar.Add(1)
				wg.Done()
			}
			output <- ret
		}(input, output, &wg, bar)

		f := func(out io.Writer, bucket string, segment SnapshotSegment, force bool, ch chan SegmentError, sem chan int) {
			sem <- 1
			time.Sleep(1 * time.Second)
			if force {
				//color.Green("Restore Storage Class to %s -> STANDARD", segment.StorageClass)
				err := r.restoreObject(out, aws.String(bucket), aws.String(segment.Key))
				ch <- SegmentError{
					Key:   segment.Key,
					Error: err,
				}
			} else {
				reader := bufio.NewReader(os.Stdin)

				color.Blue("Change Storage Class to STANDARD [y/n]: ")

				resp, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(resp)) == "y" {
					color.Green("Change Storage Class to %s -> STANDARD", segment.StorageClass)
					err := r.restoreObject(out, aws.String(bucket), aws.String(segment.Key))
					ch <- SegmentError{
						Key:   segment.Key,
						Error: err,
					}
				} else {
					color.Red("Don't change storage class %s", segment.Key)
				}
			}
			<-sem
		}

		bar.Start()
		for _, s := range segments {
			if s.StorageClass == "GLACIER" {
				areAllObjectsStandard = false
				wg.Add(1)
				go f(out, repository.Settings.Bucket, s, r.Flag.Force, input, semaphore)
			} else {
				bar.Add(1)
			}
		}
		wg.Wait()
		close(input)

		bar.Finish()

		result = <-output

		if len(result) > 0 {
			for _, s := range result {
				if s.Error != nil {
					s.PrintError()
				}
			}
		}
	}

	if areAllObjectsStandard {
		requestBody, _ := json.Marshal(
			snapshotSchema.RequestBody{
				Indices: indexName,
			},
		)

		resp, _ := r.Client.RestoreSnapshot(string(requestBody), repositoryID, snapshotID)
		fmt.Println(resp)
	}

	return nil
}

func (r Runner) restoreObject(_ io.Writer, bucket *string, key *string) error {
	resp, err := r.Client.HeadObject(bucket, key)

	if err != nil {
		return err
	}

	if resp.Restore != nil {
		if *resp.Restore == "ongoing-request=\"true\"" {
			return fmt.Errorf("%s is ongoing-restore", util.StringWithColor(*key))
		}
		r.Client.TransitObject(bucket, key, "STANDARD")
		return nil
	}

	err = r.Client.RestoreObject(bucket, key, r.Flag.RestoreTier)

	if err != nil {
		return err
	}

	return nil
}

// AddObjectSegments gather all segments recursively
func (r *Runner) AddObjectSegments(bucket, prefix string, token *string) []SnapshotSegment {
	var segments []SnapshotSegment
	objs, _ := r.Client.GetObjects(aws.String(bucket), aws.String(prefix), nil, token)
	for _, item := range objs.Contents {
		segments = append(segments, SnapshotSegment{
			Key:          *item.Key,
			StorageClass: *item.StorageClass,
		})
	}
	if *objs.IsTruncated {
		return append(segments, r.AddObjectSegments(bucket, prefix, objs.NextContinuationToken)...)
	}

	return segments
}

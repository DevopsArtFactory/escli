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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/fatih/color"

	"github.com/DevopsArtFactory/escli/internal/constants"
	snapshotSchema "github.com/DevopsArtFactory/escli/internal/schema/snapshot"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) ListSnapshot(out io.Writer) error {
	repositories, err := r.Client.GetRepositories()
	if err != nil {
		return err
	}

	for _, repository := range repositories {
		if r.Flag.WithRepo == repository.ID || r.Flag.WithRepo == constants.EmptyString {
			fmt.Fprintf(out, "Repository ID : %s\n", util.StringWithColor(repository.ID))
		}
		if !r.Flag.RepoOnly {
			if r.Flag.WithRepo == repository.ID || r.Flag.WithRepo == constants.EmptyString {
				snapshotsMetadata, err := r.Client.GetSnapshots(repository.ID)
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

	snapshotsIndicesS3, err := r.getIndexIDFromS3(
		aws.String(repository.Settings.Bucket),
		aws.String(repository.Settings.BasePath+"/"),
	)

	if err != nil {
		return err
	}

	for _, indexName := range snapshot.Indices {
		fmt.Fprintf(out, "index name : %s\n", util.StringWithColor(indexName))
		metaData := snapshotsIndicesS3.Indices[indexName]
		prefix := repository.Settings.BasePath + "/indices/" + metaData.ID + "/"

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
			if currentKeyNumber > keyNumber {
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
	repositoryID := args[0]
	snapshotID := args[1]
	indexName := args[2]
	areAllObjectsStandard := true

	repository := r.Client.GetRepository(repositoryID)

	if repository.Type == "s3" {
		fmt.Fprintf(out, "bucket name : %s\n", util.StringWithColor(repository.Settings.Bucket))
		fmt.Fprintf(out, "base path : %s\n", util.StringWithColor(repository.Settings.BasePath))

		snapshotsIndicesS3, err := r.getIndexIDFromS3(
			aws.String(repository.Settings.Bucket),
			aws.String(repository.Settings.BasePath+"/"),
		)

		if err != nil {
			return err
		}

		metaData := snapshotsIndicesS3.Indices[indexName]
		prefix := repository.Settings.BasePath + "/indices/" + metaData.ID + "/"
		objs, _ := r.Client.GetObjects(aws.String(repository.Settings.Bucket), aws.String(prefix), nil, nil)

		for {
			for _, item := range objs.Contents {
				fmt.Println(*item.Key)
				if *item.StorageClass == "GLACIER" {
					areAllObjectsStandard = false
					if r.Flag.Force {
						color.Green("Restore Storage Class to %s -> STANDARD", *item.StorageClass)
						err := r.restoreObject(out, aws.String(repository.Settings.Bucket), item.Key)
						if err != nil {
							panic(err)
						}
					} else {
						reader := bufio.NewReader(os.Stdin)

						color.Blue("Change Storage Class to STANDARD [y/n]: ")

						resp, _ := reader.ReadString('\n')
						if strings.ToLower(strings.TrimSpace(resp)) == "y" {
							color.Green("Change Storage Class to %s -> STANDARD", *item.StorageClass)
							err := r.restoreObject(out, aws.String(repository.Settings.Bucket), item.Key)
							if err != nil {
								panic(err)
							}
						} else {
							color.Red("Don't change storage class %s", *item.Key)
						}
					}
				}
			}

			if *objs.IsTruncated {
				objs, _ = r.Client.GetObjects(aws.String(repository.Settings.Bucket), aws.String(prefix), nil, objs.NextContinuationToken)
			} else {
				break
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

func (r Runner) restoreObject(out io.Writer, bucket *string, key *string) error {
	resp, err := r.Client.HeadObject(bucket, key)

	if err != nil {
		return err
	}

	if resp.Restore != nil {
		if *resp.Restore == "ongoing-request=\"true\"" {
			fmt.Fprintf(out, "%s is ongoing-restore\n", util.StringWithColor(*key))
			return nil
		}
		r.Client.TransitObject(bucket, key, "STANDARD")
		return nil
	}

	err = r.Client.RestoreObject(bucket, key)

	if err != nil {
		return err
	}

	return nil
}

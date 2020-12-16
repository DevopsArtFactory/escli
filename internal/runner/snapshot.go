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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/fatih/color"

	"github.com/DevopsArtFactory/escli/internal/schema"
	"github.com/DevopsArtFactory/escli/internal/util"
)

func (r Runner) ListSnapshot(out io.Writer) error {
	r.Client.ListSnapshot()
	return nil
}

func (r Runner) ArchiveSnapshot(out io.Writer, args []string) error {
	var repositoryName string
	var snapshotName string
	var snapshotMetadata schema.SnapshotMetadata
	var indexMetadata map[string]interface{}

	var wait sync.WaitGroup
	runtime.GOMAXPROCS(4)

	repositoryName = args[0]
	snapshotName = args[1]

	bucketName, basePath := r.Client.GetRepositoryMetadata(repositoryName)

	fmt.Println(bucketName)
	fmt.Println(basePath)

	util.ConvertToSnapshotMetadata(
		r.Client.GetSnapshotMetadata(repositoryName, snapshotName),
		&snapshotMetadata)

	fmt.Println(snapshotMetadata.Snapshots[0].Indices)

	resp := r.getIndexIDFromS3(aws.String(bucketName), aws.String(basePath+"/"))

	indexMetadata = resp["indices"].(map[string]interface{})

	for i := range snapshotMetadata.Snapshots[0].Indices {
		indexName := snapshotMetadata.Snapshots[0].Indices[i]
		metaData := indexMetadata[indexName].(map[string]interface{})
		prefix := basePath + "/indices/" + metaData["id"].(string) + "/"

		objs := r.Client.GetObjects(aws.String(bucketName), aws.String(prefix), nil, nil)

		for {
			for _, item := range objs.Contents {
				fmt.Println(*item.Key)
				if *item.StorageClass != "GLACIER" {
					if r.Flag.Force {
						color.Green("Change Storage Class to %s -> GLACIER", *item.StorageClass)
						wait.Add(1)
						go func(key string) {
							defer wait.Done()
							_, err := r.Client.TransitObject(aws.String(bucketName), aws.String(key), "GLACIER")
							if err != nil {
								panic(err)
							}
						}(*item.Key)
					} else {
						reader := bufio.NewReader(os.Stdin)

						color.Blue("Change Storage Class to GLACIER [y/n]: ")

						resp, _ := reader.ReadString('\n')
						if strings.ToLower(strings.TrimSpace(resp)) == "y" {
							color.Green("Change Storage Class to %s -> GLACIER", *item.StorageClass)
							wait.Add(1)
							_, err := r.Client.TransitObject(aws.String(bucketName), item.Key, "GLACIER")
							if err != nil {
								panic(err)
							}
						} else {
							color.Red("Don't change storage class %s", *item.Key)
						}
					}
				}
			}
			wait.Wait()

			if *objs.IsTruncated {
				objs = r.Client.GetObjects(aws.String(bucketName), aws.String(prefix), nil, objs.NextContinuationToken)
			} else {
				break
			}
		}
	}
	return nil
}

func (r Runner) getIndexIDFromS3(bucket *string, prefix *string) map[string]interface{} {
	resp := r.Client.GetObjects(bucket, prefix, aws.String("/"), nil)

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

	// Open our jsonFile
	jsonFile, err := os.Open(downloadFileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	return result
}

func (r Runner) RestoreSnapshot(out io.Writer, args []string) error {
	var repositoryName string
	var snapshotName string
	var indexName string
	var requestBody string
	var indexMetadata map[string]interface{}
	areAllObjectsStandard := true

	repositoryName = args[0]
	snapshotName = args[1]
	indexName = args[2]

	bucketName, basePath := r.Client.GetRepositoryMetadata(repositoryName)

	resp := r.getIndexIDFromS3(aws.String(bucketName), aws.String(basePath+"/"))
	indexMetadata = resp["indices"].(map[string]interface{})

	metaData := indexMetadata[indexName].(map[string]interface{})
	prefix := basePath + "/indices/" + metaData["id"].(string) + "/"

	objs := r.Client.GetObjects(aws.String(bucketName), aws.String(prefix), nil, nil)

	for {
		for _, item := range objs.Contents {
			fmt.Println(*item.Key)
			if *item.StorageClass == "GLACIER" {
				areAllObjectsStandard = false
				if r.Flag.Force {
					color.Green("Restore Storage Class to %s -> STANDARD", *item.StorageClass)
					err := r.Client.RestoreObject(aws.String(bucketName), item.Key)
					if err != nil {
						panic(err)
					}
				} else {
					reader := bufio.NewReader(os.Stdin)

					color.Blue("Change Storage Class to STANDARD [y/n]: ")

					resp, _ := reader.ReadString('\n')
					if strings.ToLower(strings.TrimSpace(resp)) == "y" {
						color.Green("Change Storage Class to %s -> STANDARD", *item.StorageClass)
						err := r.Client.RestoreObject(aws.String(bucketName), item.Key)
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
			objs = r.Client.GetObjects(aws.String(bucketName), aws.String(prefix), nil, objs.NextContinuationToken)
		} else {
			break
		}
	}

	if areAllObjectsStandard {
		resp, _ := r.Client.RestoreSnapshot(requestBody, repositoryName, snapshotName)
		fmt.Println(resp)
	}

	return nil
}

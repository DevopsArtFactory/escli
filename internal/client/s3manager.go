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

package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func GetS3DownloaderClientFn(sess client.ConfigProvider, creds *credentials.Credentials) *s3manager.Downloader {
	if creds == nil {
		return s3manager.NewDownloader(sess)
	}
	return s3manager.NewDownloader(sess)
}

func (c Client) DownloadFileFromS3(bucket *string, key *string) (filename string) {
	t := strings.Split(*key, "/")
	downloadFileName := "/tmp/" + t[len(t)-1]

	file, _ := os.Create(downloadFileName)

	defer file.Close()

	numBytes, err := c.S3Downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return downloadFileName
}

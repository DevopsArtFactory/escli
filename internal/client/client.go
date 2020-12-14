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
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/elastic/go-elasticsearch"

	"github.com/DevopsArtFactory/escli/internal/schema"
)

type Client struct {
	ESClient     *elasticsearch.Client
	S3Client     *s3.S3
	S3Downloader *s3manager.Downloader
	Region       string
}

func NewClient(sess client.ConfigProvider, creds *credentials.Credentials, config *schema.Config) Client {
	return Client{
		ESClient:     GetESClientFn(config.ElasticSearchURL),
		S3Client:     GetS3ClientFn(sess, creds),
		S3Downloader: GetS3DownloaderClientFn(sess, creds),
	}
}

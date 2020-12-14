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

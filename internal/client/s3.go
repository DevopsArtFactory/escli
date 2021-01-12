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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetS3ClientFn(sess client.ConfigProvider, creds *credentials.Credentials) *s3.S3 {
	if creds == nil {
		return s3.New(sess)
	}
	return s3.New(sess, &aws.Config{Credentials: creds})
}

func (c Client) GetObjects(bucket *string, prefix *string, delimeter *string, continuationToken *string) (*s3.ListObjectsV2Output, error) {
	resp, err := c.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: bucket, Prefix: prefix, Delimiter: delimeter, ContinuationToken: continuationToken})

	return resp, err
}

func (c Client) TransitObject(bucket *string, key *string, storageClass string) (*s3.CopyObjectOutput, error) {
	copySource := fmt.Sprintf("%s/%s", *bucket, *key)
	resp, err := c.S3Client.CopyObject(&s3.CopyObjectInput{Bucket: bucket, CopySource: aws.String(copySource), Key: key, StorageClass: aws.String(storageClass)})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func (c Client) HeadObject(bucket *string, key *string) (*s3.HeadObjectOutput, error) {
	resp, err := c.S3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp, nil
}

func (c Client) RestoreObject(bucket *string, key *string) error {
	_, err := c.S3Client.RestoreObject(&s3.RestoreObjectInput{
		Bucket: bucket,
		Key:    key,
		RestoreRequest: &s3.RestoreRequest{
			Days: aws.Int64(10),
			GlacierJobParameters: &s3.GlacierJobParameters{
				Tier: aws.String("Standard"),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

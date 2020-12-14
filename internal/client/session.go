package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// GetAwsSession creates new session for AWS
func GetAwsSession(region string) *session.Session {
	mySession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	return mySession
}

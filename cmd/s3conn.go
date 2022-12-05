package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

func s3conn() *session.Session {

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(currentCreds.awsKey, currentCreds.awsSecretKey, ""),
		Region:      aws.String(currentConfig.awsRegion),
	}

	if currentConfig.s3Endpoint != " " {
		s3Config.Endpoint = aws.String(currentConfig.s3Endpoint)
		s3Config.S3ForcePathStyle = aws.Bool(true)
		s3Config.DisableSSL = aws.Bool(true)
	}

	s3session, s3err := session.NewSession(s3Config)

	if s3err != nil {
		log.Println("Failed to connect to S3 bucket using provided credentials")
		log.Println(s3err)
		return nil
	}

	return s3session

}

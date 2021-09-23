package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func uploadBackup(backupfilename string, backupkeyname string, bucketname string, awskey string, awssecretkey string, awsregion string) bool {

	log.Println("Uploading file", backupfilename, "to bucket", bucketname, "as", backupkeyname)

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(awskey, awssecretkey, ""),
		Region:      aws.String(awsregion),
	}

	localArchive, ferr := os.Open(backupfilename)
	if ferr != nil {
		panic(ferr)
	}

	defer localArchive.Close()

	newSession, s3err := session.NewSession(s3Config)
	if s3err != nil {
		log.Println("Failed to connect to S3 bucket using provided credentials")
		log.Println(s3err)
		return false
	}

	s3Client := s3.New(newSession)

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(backupkeyname),
		Body:   localArchive,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			log.Printf("Upload canceled due to timeout, %v\n", err)
		} else {
			log.Printf("Failed to upload object, %v\n", err)
		}
		return false
	}

	log.Printf("Successfully uploaded file to %s/%s \n", bucketname, backupkeyname)

	return true

}

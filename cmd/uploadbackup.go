package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func uploadBackup(backupfilename string, backupkeyname string, bucketname string) bool {

	log.Println("Uploading file", backupfilename, "to bucket", bucketname, "as", backupkeyname)

	localArchive, ferr := os.Open(backupfilename)
	if ferr != nil {
		panic(ferr)
	}

	defer localArchive.Close()

	s3Client := s3.New(s3conn())

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

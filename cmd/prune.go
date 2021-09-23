package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

func prune(filesList []string, filesKeep int, bucketname string, awskey string, awssecretkey string, awsregion string) {

	log.Println("Starting to prune objects...")

	if len(filesList) <= filesKeep {
		log.Println("Files list is too short, nothing to prune")
		return
	}

	pruneList := filesList[filesKeep:]

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(awskey, awssecretkey, ""),
		Region:      aws.String(awsregion),
	}

	bucket := aws.String(bucketname)

	newSession, s3err := session.NewSession(s3Config)
	if s3err != nil {
		log.Println("Failed to connect to S3 bucket using provided credentials")
		log.Println(s3err)
		return
	}

	pruneObjectsList := make([]*s3.ObjectIdentifier, 0, 1000)

	for _, objectname := range pruneList {
		obj := s3.ObjectIdentifier{
			Key: aws.String(objectname),
		}
		pruneObjectsList = append(pruneObjectsList, &obj)
	}

	s3objects := s3.New(newSession)

	input := &s3.DeleteObjectsInput{
		Bucket: bucket,
		Delete: &s3.Delete{
			Objects: pruneObjectsList,
			Quiet:   aws.Bool(false),
		},
	}

	_, err := s3objects.DeleteObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return
	}

	log.Println("Following objects were deleted from S3 bucket:")
	for _, name := range pruneList {
		fmt.Println(name)
	}

}

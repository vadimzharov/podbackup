package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"sort"
)

func listBackups(backupkeypath string, backupkeyprefix string, bucketname string, awskey string, awssecretkey string, awsregion string) []string {

	log.Println("Listing files from bucket", bucketname, "directory", backupkeypath)

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(awskey, awssecretkey, ""),
		Region:      aws.String(awsregion),
	}

	bucket := aws.String(bucketname)

	newSession, s3err := session.NewSession(s3Config)
	if s3err != nil {
		log.Println("Failed to connect to S3 bucket using provided credentials")
		log.Println(s3err)
		return nil
	}

	s3Client := s3.New(newSession)

	s3FilelistFilter := backupkeypath + backupkeyprefix

	log.Println(s3FilelistFilter)

	input := &s3.ListObjectsInput{
		Bucket: bucket,
		Prefix: aws.String(s3FilelistFilter),
	}

	objlist, err := s3Client.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				log.Println(aerr.Error())
				return nil
			}
		} else {
			log.Println(err.Error())
		}
		return nil

	}

	listFiles := objlist.Contents

	sort.Slice(listFiles, func(i, j int) bool {
		return listFiles[i].LastModified.After(*listFiles[j].LastModified)
	})

	var namesList []string

	for _, object := range listFiles {
		namesList = append(namesList, *object.Key)
	}

	return namesList

}

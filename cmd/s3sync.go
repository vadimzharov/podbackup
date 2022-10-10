package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
	"log"
)

func syncToS3() {

	localfolder := currentConfig.backupDir

	s3path := "s3://" + currentConfig.bucketName + "/" + currentConfig.bucketFolder

	log.Println("Uploading files from ", localfolder, "to bucket", s3path)

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(currentCreds.awsKey, currentCreds.awsSecretKey, ""),
		Region:      aws.String(currentConfig.awsRegion),
	}

	if currentConfig.s3Endpoint != " " {
		s3Config.Endpoint = aws.String(currentConfig.s3Endpoint)
		s3Config.S3ForcePathStyle = aws.Bool(true)
		s3Config.DisableSSL = aws.Bool(true)
	}

	newSession, s3err := session.NewSession(s3Config)
	if s3err != nil {
		log.Println("Failed to connect to S3 bucket using provided credentials")
		log.Println(s3err)
		panic(s3err)
	}

	s3err = s3sync.New(newSession, s3sync.WithParallel(currentConfig.s3SyncParallelism)).Sync(localfolder, s3path)

	if s3err != nil {
		log.Println("Failed to upload files to S3 bucket!")
		log.Println(s3err)
		panic(s3err)
	}

	log.Println("Successfully uploaded all files from ", localfolder, " to ", s3path)

}

func syncFromS3() {

	localfolder := currentConfig.restoreDir

	s3path := "s3://" + currentConfig.bucketName + "/" + currentConfig.bucketFolder

	log.Println("Downloading files from ", s3path, "to folder", localfolder)

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(currentCreds.awsKey, currentCreds.awsSecretKey, ""),
		Region:      aws.String(currentConfig.awsRegion),
	}

	if currentConfig.s3Endpoint != " " {
		s3Config.Endpoint = aws.String(currentConfig.s3Endpoint)
		s3Config.S3ForcePathStyle = aws.Bool(true)
		s3Config.DisableSSL = aws.Bool(true)
	}

	newSession, s3err := session.NewSession(s3Config)
	if s3err != nil {
		log.Println("Failed to connect to S3 bucket using provided credentials")
		log.Println(s3err)
		panic(s3err)
	}

	s3err = s3sync.New(newSession, s3sync.WithParallel(currentConfig.s3SyncParallelism)).Sync(s3path, localfolder)

	if s3err != nil {
		log.Println("Failed to download files from S3 bucket!")
		log.Println(s3err)
		panic(s3err)
	}

	log.Println("Successfully downloaded all files from ", s3path, " to ", localfolder)

}

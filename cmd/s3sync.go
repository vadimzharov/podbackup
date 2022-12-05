package cmd

import (
	"github.com/seqsense/s3sync"
	"log"
)

func syncToS3() {

	localfolder := currentConfig.backupDir

	s3path := "s3://" + currentConfig.bucketName + "/" + currentConfig.bucketFolder

	log.Println("Uploading files from ", localfolder, "to bucket", s3path)

	s3err := s3sync.New(s3conn(), s3sync.WithParallel(currentConfig.s3SyncParallelism)).Sync(localfolder, s3path)

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

	s3err := s3sync.New(s3conn(), s3sync.WithParallel(currentConfig.s3SyncParallelism)).Sync(s3path, localfolder)

	if s3err != nil {
		log.Println("Failed to download files from S3 bucket!")
		log.Println(s3err)
		panic(s3err)
	}

	log.Println("Successfully downloaded all files from ", s3path, " to ", localfolder)

}

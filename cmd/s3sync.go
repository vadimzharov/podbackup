package cmd

import (
	"log"

	"github.com/seqsense/s3sync"
)

var s3err error

func syncToS3() {

	localfolder := currentConfig.backupDir

	s3path := "s3://" + currentConfig.bucketName + "/" + currentConfig.bucketFolder

	if currentConfig.s3CopyWithDelete {
		log.Println("Uploading files from ", localfolder, "to bucket", s3path, "and deleting files in destination if they don't exist in source")
		s3err = s3sync.New(s3conn(), s3sync.WithParallel(currentConfig.s3SyncParallelism), s3sync.WithDelete()).Sync(localfolder, s3path)
	} else {
		log.Println("Uploading files from ", localfolder, "to bucket", s3path)
		s3err = s3sync.New(s3conn(), s3sync.WithParallel(currentConfig.s3SyncParallelism)).Sync(localfolder, s3path)
	}

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

	if currentConfig.s3CopyWithDelete {
		log.Println("Downloading files from ", s3path, "to folder", localfolder, "and deleting files in destination if they don't exist in source")
		s3err = s3sync.New(s3conn(), s3sync.WithParallel(currentConfig.s3SyncParallelism), s3sync.WithDelete()).Sync(s3path, localfolder)
	} else {
		log.Println("Downloading files from ", s3path, "to folder", localfolder)
		s3err = s3sync.New(s3conn(), s3sync.WithParallel(currentConfig.s3SyncParallelism)).Sync(s3path, localfolder)
	}

	if s3err != nil {
		log.Println("Failed to download files from S3 bucket!")
		log.Println(s3err)
		panic(s3err)
	}

	log.Println("Successfully downloaded all files from ", s3path, " to ", localfolder)

}

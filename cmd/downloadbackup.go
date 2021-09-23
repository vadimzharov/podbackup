package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

func downloadBackup(restorefilename string, backupkeyname string, bucketname string, awskey string, awssecretkey string, awsregion string) (restoredfile *string) {

	log.Println("Downloading file", backupkeyname, "from bucket", bucketname, "as", restorefilename)

	file, err := os.Create(restorefilename)
	if err != nil {
		log.Println("Failed to create localfile ", restorefilename, err)
		return nil
	}
	defer file.Close()

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(awskey, awssecretkey, ""),
		Region:      aws.String(awsregion),
	}

	bucket := aws.String(bucketname)
	key := aws.String(backupkeyname)

	newSession, s3err := session.NewSession(s3Config)
	if s3err != nil {
		log.Println("Failed to connect to S3 bucket using provided credentials")
		log.Println(s3err)
		return nil
	}

	s3Downloader := s3manager.NewDownloader(newSession)

	numbytes, err := s3Downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		})
	if err != nil {
		log.Println("Failed to download file", backupkeyname, "from bucket", bucketname)
		log.Println(err)
		return nil
	}

	log.Println("Downloaded file", backupkeyname, "to", file.Name(), numbytes, "bytes total")

	return &restorefilename
}

package cmd

import (
	"log"
	"os"
	"time"
)

func backupFiles() {

	var backupkeyname string

	if currentConfig.backupDir == "" {
		log.Println("DIR_TO_BACKUP variable is not set or empty, don't know what to backup. Exiting..")
		return
	}

	t := time.Now().UTC().Format("20060102150405")

	switch currentConfig.archiveType {
	case "tarzip":
		{
			currentConfig.backupLocalFile = backuptempdir + "tararchive/backup.tar"
			makeTarBackup(currentConfig.backupDir, currentConfig.backupLocalFile)
			currentConfig.backupLocalFile = backuptempdir + "backup.zip"
			backupkeyname = currentConfig.bucketFolder + currentConfig.keyPrefix + "-" + t + ".tar.zip"
			currentConfig.backupDir = backuptempdir + "tararchive/"

			makeBackup(currentConfig.backupDir, currentConfig.backupLocalFile, currentCreds.encryptpassword)
		}

	case "zip":
		{
			currentConfig.backupLocalFile = backuptempdir + "backup.zip"
			backupkeyname = currentConfig.bucketFolder + currentConfig.keyPrefix + "-" + t + ".zip"
			makeBackup(currentConfig.backupDir, currentConfig.backupLocalFile, currentCreds.encryptpassword)
		}

	default:

		log.Panic("ARCHIVE_TYPE environment variable is not correct (should be zip or tarzip")
		os.Exit(1)

	}

	//if currentConfig.useTar {
	//
	//	currentConfig.backupLocalFile = backuptempdir + "tararchive/backup.tar"
	//	makeTarBackup(currentConfig.backupDir, currentConfig.backupLocalFile)

	//	currentConfig.backupLocalFile = backuptempdir + "backup.zip"
	//	backupkeyname = currentConfig.bucketFolder + currentConfig.keyPrefix + "-" + t + ".tar.zip"
	//	currentConfig.backupDir = backuptempdir + "tararchive/"

	//	makeBackup(currentConfig.backupDir, currentConfig.backupLocalFile, currentCreds.encryptpassword)

	//} else {
	//	currentConfig.backupLocalFile = backuptempdir + "backup.zip"
	//	backupkeyname = currentConfig.bucketFolder + currentConfig.keyPrefix + "-" + t + ".zip"
	//	makeBackup(currentConfig.backupDir, currentConfig.backupLocalFile, currentCreds.encryptpassword)

	//}

	uploadBackup(currentConfig.backupLocalFile, backupkeyname, currentConfig.bucketName, currentCreds.awsKey, currentCreds.awsSecretKey, currentConfig.awsRegion)

	os.Remove(currentConfig.backupLocalFile)
}

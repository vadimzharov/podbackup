package cmd

import (
	"log"
	"os"
	"time"
)

func backupFiles() {

	if currentConfig.backupDir == "" {
		log.Println("DIR_TO_BACKUP variable is not set or empty, don't know what to backup. Exiting..")
		return
	}

	t := time.Now().UTC().Format("20060102150405")

	backupkeyname := currentConfig.bucketFolder + currentConfig.keyPrefix + "-" + t + ".zip"

	makeBackup(currentConfig.backupDir, currentConfig.backupLocalFile, currentCreds.encryptpassword)

	uploadBackup(currentConfig.backupLocalFile, backupkeyname, currentConfig.bucketName, currentCreds.awsKey, currentCreds.awsSecretKey, currentConfig.awsRegion)

	os.Remove(currentConfig.backupLocalFile)
}

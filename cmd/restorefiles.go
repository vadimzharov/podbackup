package cmd

import (
	"log"
	"os"
)

func restoreFiles(cmdargs []string) {

	var backupkeyname string

	if currentConfig.restoreDir == "" {
		log.Println("DIR_TO_RESTORE variable is not set or empty, don't know where to restore. Exiting..")
		return
	}

	if len(cmdargs) > 2 {
		backupkeyname = cmdargs[2]
	} else {
		filesList := listBackups(currentConfig.bucketFolder, currentConfig.keyPrefix, currentConfig.bucketName, currentCreds.awsKey, currentCreds.awsSecretKey, currentConfig.awsRegion)

		if filesList == nil {
			log.Println("Cannot list files in bucket", currentConfig.bucketName)

			if currentConfig.forceRestore {
				log.Fatal("Cannot list files in bucket", currentConfig.bucketName, ". Cannot continue due to FORCE_RESTORE set to True. Exiting with error")
				os.Exit(1)
			}

			return
		}
		backupkeyname = filesList[0]
	}

	downloadedFile := downloadBackup(currentConfig.backupLocalFile, backupkeyname, currentConfig.bucketName, currentCreds.awsKey, currentCreds.awsSecretKey, currentConfig.awsRegion)

	if downloadedFile == nil {
		log.Println("File could not be downloaded from S3 storage. Nothing to restore.")

		if currentConfig.forceRestore {
			log.Fatal("File", backupkeyname, "Could not be downloaded from S3 storage. Cannot continue due to FORCE_RESTORE set to True. Exiting with error...")
			os.Exit(1)
		}

		os.Exit(0)

	} else {

		restoredFiles, err := restoreBackup(currentConfig.restoreDir, *downloadedFile, currentCreds.encryptpassword)

		if restoredFiles == nil || err != nil {
			log.Println("File was downloaded, but cannot upzip it")

			if currentConfig.forceRestore {
				log.Fatal("Cannot restore files from archive. Cannot continue due to FORCE_RESTORE set to True. Exiting with error...")
			}

		}

		os.Remove(currentConfig.backupLocalFile)
	}

}

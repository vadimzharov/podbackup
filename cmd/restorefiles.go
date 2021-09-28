package cmd

import (
	"log"
	"os"
)

func restoreFiles(cmdargs []string) {

	var backupkeyname string
	var restoredFiles []string
	var resterr error

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

	currentConfig.backupLocalFile = backuptempdir + "backup.zip"

	if currentConfig.useTar {

		currentConfig.backupLocalFile = backuptempdir + "backup.tar.zip"

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

		if currentConfig.useTar {

			err := os.RemoveAll(backuptempdir + "restoredzip/")
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			err = os.Mkdir(backuptempdir+"restoredzip/", os.ModePerm)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			restoredTarFile, resttarerr := restoreBackup(backuptempdir+"restoredzip/", *downloadedFile, currentCreds.encryptpassword)

			if restoredTarFile == nil || resttarerr != nil {
				log.Println("File was downloaded, but cannot upzip it")
			}

			restoredFiles, resterr = restoreTarBackup(currentConfig.restoreDir, backuptempdir+"restoredzip/"+"backup.tar")

		} else {

			restoredFiles, resterr = restoreBackup(currentConfig.restoreDir, *downloadedFile, currentCreds.encryptpassword)
		}

		if restoredFiles == nil || resterr != nil {
			log.Println("File was downloaded, but cannot upzip it")

			if currentConfig.forceRestore {
				log.Fatal("Cannot restore files from archive. Cannot continue due to FORCE_RESTORE set to True. Exiting with error...")
			}

		}

		os.Remove(currentConfig.backupLocalFile)
	}

}

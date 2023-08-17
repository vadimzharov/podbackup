package cmd

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var currentConfig backupConfig
var currentCreds backupCreds
var isConfigValid bool
var backuptempdir string

func Main(cmdargs []string) {

	if len(cmdargs) == 1 {

		printHelp()

		os.Exit(0)

	}

	currentConfig, currentCreds, isConfigValid = getConfig()

	if !isConfigValid {
		log.Fatal("Configuration is not valid, some mandatory variables are not set. Exiting...")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	backuptempdir = "/tmp/podbackup-" + strconv.Itoa(rand.Intn(9999)) + "/"

	err := os.Mkdir(backuptempdir, os.ModePerm)
	if err != nil {
		log.Fatal("Cannot create temp directory", backuptempdir, err)
		os.Exit(1)
	}

	//backupTicker := time.NewTicker(time.Duration(currentConfig.backupInverval) * time.Minute)
	backupTicker := time.NewTicker(currentConfig.backupInverval)

	pruneTicker := time.NewTicker(currentConfig.pruneInverval)

	if len(cmdargs) > 1 {

		switch cmdargs[1] {
		case "backup-daemon":

			log.Println("Working as a daemon to make files backup, backup interval is:", currentConfig.backupInverval, ", and pruning interval is:", currentConfig.pruneInverval)

			for {
				select {

				case <-backupTicker.C:

					backupFiles()

				case <-pruneTicker.C:

					pruneCosObjects()

				}
			}

		case "backup-sql-daemon":

			log.Println("Working as a daemon to make MySQL database backup, backup interval is:", currentConfig.backupInverval, ", and pruning interval is:", currentConfig.pruneInverval)

			for {
				select {

				case <-backupTicker.C:

					backupSqlDb()

				case <-pruneTicker.C:

					pruneCosObjects()

				}
			}

		case "backup-pgsql-daemon":

			log.Println("Working as a daemon to make PostgreSQL database backup, backup interval is:", currentConfig.backupInverval, ", and pruning interval is:", currentConfig.pruneInverval)

			for {
				select {

				case <-backupTicker.C:

					backupPGSqlDb()

				case <-pruneTicker.C:

					pruneCosObjects()

				}
			}

		case "sync-to-s3":

			log.Println("Working as a daemon to sync content from ", currentConfig.backupDir, " to S3 bucket ", currentConfig.bucketName, " folder ", currentConfig.bucketFolder)

			if currentConfig.s3CopyBeforeSync {
				log.Println("S3_COPY_BEFORE_SYNC set to true, copying data from S3 before starting S3 sync process")
				syncFromS3()
			} else {
				syncToS3()
			}

			for {

				select {

				case <-backupTicker.C:
					syncToS3()

				}
			}

		case "sync-from-s3":

			log.Println("Working as a daemon to sync content from S3", currentConfig.bucketName, currentConfig.bucketFolder, "to localfolder ", currentConfig.restoreDir)

			syncFromS3()

			for {

				select {

				case <-backupTicker.C:
					syncFromS3()

				}
			}

		case "copy-to-s3":

			syncToS3()

		case "copy-from-s3":

			syncFromS3()

		case "prune":

			pruneCosObjects()

		case "backup":

			backupFiles()

		case "backup-sql":

			backupSqlDb()

		case "list":

			listCosFiles()

		case "restore":

			restoreFiles(cmdargs)

		case "restore-sql":

			restoreSqlDb(cmdargs)

		case "backup-pgsql":

			backupPGSqlDb()

		case "restore-pgsql":

			restorePGSqlDb(cmdargs)

		default:

			printHelp()

		}

	} else {

		printHelp()

	}

	err = os.RemoveAll(backuptempdir)
	if err != nil {
		log.Println("Cannot delete temp directory ", backuptempdir, err)
	}

}

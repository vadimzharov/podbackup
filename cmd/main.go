package cmd

import (
	"log"
	"os"
	"time"
)

var currentConfig backupConfig
var currentCreds backupCreds
var isConfigValid bool
var backuptempdir string

func Main(cmdargs []string) {

	currentConfig, currentCreds, isConfigValid = getConfig()

	if !isConfigValid {
		log.Fatal("Configuration is not valid, some mandatory variables are not set. Exiting...")
		os.Exit(1)
	}

	t := time.Now().UTC().Format("20060102150405")

	backuptempdir = "/tmp/podbackup-" + t + "/"

	log.Println("Creating temp directory " + backuptempdir)

	err := os.Mkdir(backuptempdir, os.ModePerm)
	if err != nil {
		log.Fatal("Cannot create temp directory", backuptempdir, err)
		os.Exit(1)
	}

	backupTicker := time.NewTicker(time.Duration(currentConfig.backupInverval) * time.Second)

	pruneTicker := time.NewTicker(time.Duration(currentConfig.pruneInverval) * time.Second)

	if len(cmdargs) > 1 {

		switch cmdargs[1] {
		case "backup-daemon":

			for {
				select {

				case <-backupTicker.C:

					backupFiles()

				case <-pruneTicker.C:

					pruneCosObjects()

				}
			}

		case "backup-sql-daemon":

			for {
				select {

				case <-backupTicker.C:

					backupSqlDb()

				case <-pruneTicker.C:

					pruneCosObjects()

				}
			}

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

		default:

			printHelp()

		}

	} else {

		printHelp()

	}

	log.Println("Deleting temp directory " + backuptempdir)

	err = os.RemoveAll(backuptempdir)
	if err != nil {
		log.Println("Cannot delete temp directory ", backuptempdir, err)
	}

}

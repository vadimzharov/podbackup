package cmd

import (
	"log"
	"time"
)

var currentConfig backupConfig
var currentCreds backupCreds
var isConfigValid bool

func Main(cmdargs []string) {

	currentConfig, currentCreds, isConfigValid = getConfig()

	if !isConfigValid {
		log.Fatal("Configuration is not valid, some mandatory variables are not set. Exiting...")
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

}

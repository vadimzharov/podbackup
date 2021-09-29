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

	currentConfig, currentCreds, isConfigValid = getConfig()

	if !isConfigValid {
		log.Fatal("Configuration is not valid, some mandatory variables are not set. Exiting...")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	backuptempdir = "/tmp/podbackup-" + strconv.Itoa(rand.Intn(99)) + "/"

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

	err = os.RemoveAll(backuptempdir)
	if err != nil {
		log.Println("Cannot delete temp directory ", backuptempdir, err)
	}

}

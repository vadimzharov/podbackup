package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func restoreSqlDb(cmdargs []string) {

	mysqlConnSleepInternval := 20

	currentSqlConfig, isSqlConfigValid := getSqlConfig()

	if !isSqlConfigValid {
		log.Fatal("Mandatory variables to make SQL backup are not defined. Cannot make backup, exiting with error....")
		os.Exit(1)
	}

	currentConfig.restoreDir = currentSqlConfig.sqlDumpDir

	restoreFiles(cmdargs)

	//Check if dump file was restored

	dumpFile := currentSqlConfig.sqlDumpDir + currentSqlConfig.sqlDumpFileName

	if _, err := os.Stat(dumpFile); err == nil {

		var mysqlclientargs []string

		mysqlclientargs = append(mysqlclientargs, fmt.Sprintf("-u%s", currentSqlConfig.sqlUser))
		mysqlclientargs = append(mysqlclientargs, fmt.Sprintf("--password=%s", currentSqlConfig.sqlPassword))
		mysqlclientargs = append(mysqlclientargs, fmt.Sprintf("-h%s", currentSqlConfig.sqlAddr))
		mysqlclientargs = append(mysqlclientargs, fmt.Sprintf("-P%s", currentSqlConfig.sqlPort))
		mysqlclientargs = append(mysqlclientargs, fmt.Sprintf("-e source %s", dumpFile))

		mysqlrestorecmd := exec.Command(currentSqlConfig.sqlClientCmdPath, mysqlclientargs...)

		log.Println("Trying to connect to MySQL DB on", currentSqlConfig.sqlAddr+":"+currentSqlConfig.sqlPort)

		for {

			mysqlconn, connerr := net.DialTimeout("tcp", net.JoinHostPort(currentSqlConfig.sqlAddr, currentSqlConfig.sqlPort), time.Duration(mysqlConnSleepInternval)*time.Second)

			if connerr != nil {

				log.Println("Cannot connect to MySQL, sleeping for", mysqlConnSleepInternval, "seconds...")

				time.Sleep(time.Duration(mysqlConnSleepInternval) * time.Second)
			}

			if mysqlconn != nil {
				mysqlconn.Close()
				log.Println("Connected to MySQL DB on", currentSqlConfig.sqlAddr+":"+currentSqlConfig.sqlPort)
				break
			}

		}

		log.Println("Executing mysql command to restore databases from the dump file...")

		var out bytes.Buffer
		var stderr bytes.Buffer

		mysqlrestorecmd.Stdout = &out
		mysqlrestorecmd.Stderr = &stderr

		mysqlrestorecmderr := mysqlrestorecmd.Run()

		if mysqlrestorecmderr != nil {

			log.Println(fmt.Sprint(mysqlrestorecmderr) + ": " + stderr.String())

		}

		log.Println("All databases were restored from the dump file")

		err := os.RemoveAll(currentSqlConfig.sqlDumpDir)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = os.Stat(currentSqlConfig.sqlRestoreMarkDir); os.IsNotExist(err) {
			err = os.Mkdir(currentSqlConfig.sqlRestoreMarkDir, os.ModePerm)
			if err != nil {
				log.Println(err)
			}

		}

		restoredMarkFilePath := currentSqlConfig.sqlRestoreMarkDir + "restored"

		log.Println("Creating file to mark restore process completed", restoredMarkFilePath)

		restoredMarkFile, err := os.Create(restoredMarkFilePath)

		if err != nil {
			log.Println("Cannot create mark file ", err)
		}

		restoredMarkFile.Close()

	} else {

		if currentConfig.forceRestore {
			log.Fatal("MYSQL dump file was not found in latest archive downloaded from S3 bucket. FORCE_RESTORE set to TRUE, but cannot restore MySQL database from archive, exiting with error...")
			os.Exit(1)
		}

		log.Println("MYSQL dump file was not found in latest archive dowloaded from S3 bucket. Skipping MYSQL restore, exiting...")
		return

	}

}

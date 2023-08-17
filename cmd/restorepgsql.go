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

func restorePGSqlDb(cmdargs []string) {

	sqlConnSleepInternval := 20

	currentSqlConfig, isSqlConfigValid := getSqlConfig("pgsql")

	if !isSqlConfigValid {
		log.Fatal("Mandatory variables to make SQL backup are not defined. Cannot make backup, exiting with error....")
		os.Exit(1)
	}

	currentConfig.restoreDir = currentSqlConfig.sqlDumpDir

	restoreFiles(cmdargs)

	//Check if dump file was restored

	dumpFile := currentSqlConfig.sqlDumpDir + currentSqlConfig.sqlDumpFileName

	if _, err := os.Stat(dumpFile); err == nil {

		var pgsqlclientargs []string

		pgsqlclientargs = append(pgsqlclientargs, fmt.Sprintf("-U%s", currentSqlConfig.sqlUser))
		pgsqlclientargs = append(pgsqlclientargs, fmt.Sprintf("-h%s", currentSqlConfig.sqlAddr))
		pgsqlclientargs = append(pgsqlclientargs, fmt.Sprintf("-p%s", currentSqlConfig.sqlPort))
		pgsqlclientargs = append(pgsqlclientargs, fmt.Sprintf("-d%s", currentSqlConfig.sqlDatabase))
		pgsqlclientargs = append(pgsqlclientargs, fmt.Sprintf("-f%s", dumpFile))

		pgsqlrestorecmd := exec.Command(currentSqlConfig.sqlClientCmdPath, pgsqlclientargs...)

		pgsqlrestorecmd.Env = append(pgsqlrestorecmd.Env, "PGPASSWORD="+currentSqlConfig.sqlPassword)

		log.Println("Trying to connect to PGSQL DB on", currentSqlConfig.sqlAddr+":"+currentSqlConfig.sqlPort)

		for {

			sqlconn, connerr := net.DialTimeout("tcp", net.JoinHostPort(currentSqlConfig.sqlAddr, currentSqlConfig.sqlPort), time.Duration(sqlConnSleepInternval)*time.Second)

			if connerr != nil {

				log.Println("Cannot connect to PGSQL, sleeping for", sqlConnSleepInternval, "seconds...")

				time.Sleep(time.Duration(sqlConnSleepInternval) * time.Second)
			}

			if sqlconn != nil {
				sqlconn.Close()
				log.Println("Connected to PGSQL DB on", currentSqlConfig.sqlAddr+":"+currentSqlConfig.sqlPort)
				break
			}

		}

		log.Println("Executing pgsql command to restore databases from the dump file...")

		var out bytes.Buffer
		var stderr bytes.Buffer

		pgsqlrestorecmd.Stdout = &out
		pgsqlrestorecmd.Stderr = &stderr

		err := pgsqlrestorecmd.Run()

		if err != nil {

			log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
			if currentConfig.forceRestore {
				log.Fatal("PGSQL restore process failed. FORCE_RESTORE set to TRUE, but cannot restore PGSQL database from the archive, exiting with error...")
				os.Exit(1)
			}

			log.Println("PGSQL dump file was not found in latest archive dowloaded from S3 bucket. Skipping PGSQL restore, exiting...")
			log.Println(fmt.Sprint(err) + ":" + stderr.String())
			return

		} else {

			log.Println("All databases were restored from the dump file, cleaning FS...")

		}

		err = os.RemoveAll(currentSqlConfig.sqlDumpDir)
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
			log.Fatal("PGSQL dump file was not found in latest archive downloaded from S3 bucket. FORCE_RESTORE set to TRUE, but cannot restore PGSQL database from archive, exiting with error...")
			os.Exit(1)
		}

		log.Println("PGSQL dump file was not found in latest archive dowloaded from S3 bucket. Skipping PGSQL restore, exiting...")
		return

	}

}

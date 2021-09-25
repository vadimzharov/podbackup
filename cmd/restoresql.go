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

	var mysqlclientargs []string

	dumpFile := currentSqlConfig.sqlDumpDir + currentSqlConfig.sqlDumpFileName

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

		fmt.Println(fmt.Sprint(mysqlrestorecmderr) + ": " + stderr.String())

	}

	log.Println("All databases were restored from the dump file")

	err := os.RemoveAll(currentSqlConfig.sqlDumpDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

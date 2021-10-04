package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func backupSqlDb() {

	currentSqlConfig, isSqlConfigValid := getSqlConfig()

	var mysqldumpargs []string

	dumpFile := currentSqlConfig.sqlDumpDir + currentSqlConfig.sqlDumpFileName

	if !isSqlConfigValid {
		log.Fatal("Mandatory variables to make SQL backup are not defined. Cannot make backup, exiting with error....")
		os.Exit(1)
	}

	err := os.RemoveAll(currentSqlConfig.sqlDumpDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = os.Mkdir(currentSqlConfig.sqlDumpDir, os.ModePerm)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	mysqldumpargs = append(mysqldumpargs, fmt.Sprintf("-u%s", currentSqlConfig.sqlUser))
	mysqldumpargs = append(mysqldumpargs, fmt.Sprintf("--password=%s", currentSqlConfig.sqlPassword))
	mysqldumpargs = append(mysqldumpargs, fmt.Sprintf("-h%s", currentSqlConfig.sqlAddr))
	mysqldumpargs = append(mysqldumpargs, fmt.Sprintf("-P%s", currentSqlConfig.sqlPort))
	mysqldumpargs = append(mysqldumpargs, fmt.Sprintf("-r%s", dumpFile))
	mysqldumpargs = append(mysqldumpargs, "--all-databases")
	mysqldumpargs = append(mysqldumpargs, "--flush-privileges")

	mysqldumpcmd := exec.Command(currentSqlConfig.sqlDumpCmdPath, mysqldumpargs...)

	log.Println("Executing mysqldump command to dump all databases...")

	var out bytes.Buffer
	var stderr bytes.Buffer

	mysqldumpcmd.Stdout = &out
	mysqldumpcmd.Stderr = &stderr
	mysqldumpcmderr := mysqldumpcmd.Run()

	if mysqldumpcmderr != nil {
		fmt.Println(fmt.Sprint(mysqldumpcmderr) + ": " + stderr.String())
		os.Exit(1)
	}

	currentConfig.backupDir = currentSqlConfig.sqlDumpDir

	backupFiles()

}

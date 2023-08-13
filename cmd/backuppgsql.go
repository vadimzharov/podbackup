package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func backupPGSqlDb() {

	currentSqlConfig, isSqlConfigValid := getSqlConfig("pgsql")

	var pgsqldumpargs []string

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

	pgsqldumpargs = append(pgsqldumpargs, fmt.Sprintf("-U%s", currentSqlConfig.sqlUser))
	pgsqldumpargs = append(pgsqldumpargs, fmt.Sprintf("-h%s", currentSqlConfig.sqlAddr))
	pgsqldumpargs = append(pgsqldumpargs, fmt.Sprintf("-p%s", currentSqlConfig.sqlPort))
	pgsqldumpargs = append(pgsqldumpargs, fmt.Sprintf(currentSqlConfig.sqlDatabase))
	pgsqldumpargs = append(pgsqldumpargs, fmt.Sprintf("--file=%s", dumpFile))

	pgsqldumpcmd := exec.Command(currentSqlConfig.sqlDumpCmdPath, pgsqldumpargs...)

	pgsqldumpcmd.Env = append(pgsqldumpcmd.Env, "PGPASSWORD="+currentSqlConfig.sqlPassword)

	log.Println("Executing pg_dump command to dump database" + currentSqlConfig.sqlDatabase)

	var out bytes.Buffer
	var stderr bytes.Buffer

	pgsqldumpcmd.Stdout = &out
	pgsqldumpcmd.Stderr = &stderr
	pgsqldumpcmderr := pgsqldumpcmd.Run()

	if pgsqldumpcmderr != nil {
		fmt.Println(fmt.Sprint(pgsqldumpcmderr) + ": " + stderr.String())
		os.Exit(1)
	}

	currentConfig.backupDir = currentSqlConfig.sqlDumpDir

	backupFiles()

}

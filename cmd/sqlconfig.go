package cmd

import (
	"log"
	"os"
)

type sqlBackupConfig struct {
	sqlUser         string
	sqlPassword     string
	sqlAddr         string
	sqlPort         string
	sqlDumpCmdPath  string
	sqlDumpDir      string
	sqlDumpFileName string
}

func setDefaultSqlConfig() sqlBackupConfig {
	return sqlBackupConfig{
		sqlUser:         "root",
		sqlAddr:         "127.0.0.1",
		sqlPort:         "3306",
		sqlDumpCmdPath:  "/usr/bin/mysqldump",
		sqlDumpDir:      "/tmp/mysqldump/",
		sqlDumpFileName: "mysql.dump",
	}
}

func getSqlConfig() (sqlBackupParams sqlBackupConfig, configvalid bool) {
	var currentSqlConfig sqlBackupConfig
	var isSqlConfgValid bool

	currentSqlConfig = setDefaultSqlConfig()

	isSqlConfgValid = true

	if sqlUserenv := os.Getenv("MYSQL_USER"); sqlUserenv != "" {
		currentSqlConfig.sqlUser = sqlUserenv
	} else {
		log.Println("MYSQL_USER environment variable is not set, using the default", currentSqlConfig.sqlUser)
	}

	if sqlPasswordenv := os.Getenv("MYSQL_PASSWORD"); sqlPasswordenv != "" {
		currentSqlConfig.sqlPassword = sqlPasswordenv
	} else {
		log.Println("MYSQL_PASSWORD environment variable is not set, this is required to perform MYSQL backup")
		isSqlConfgValid = false
	}

	if sqlHostenv := os.Getenv("MYSQL_HOST"); sqlHostenv != "" {
		currentSqlConfig.sqlAddr = sqlHostenv
	} else {
		log.Println("MYSQL_HOST environment variable is not set, using the default", currentSqlConfig.sqlAddr)
	}

	if sqlPortenv := os.Getenv("MYSQL_PORT"); sqlPortenv != "" {
		currentSqlConfig.sqlPort = sqlPortenv
	} else {
		log.Println("MYSQL_PORT environment variable is not set, using the default", currentSqlConfig.sqlPort)
	}

	return currentSqlConfig, isSqlConfgValid
}

package cmd

import (
	"log"
	"os"
)

type sqlBackupConfig struct {
	sqlUser           string
	sqlPassword       string
	sqlDatabase       string
	sqlAddr           string
	sqlPort           string
	sqlDumpCmdPath    string
	sqlDumpDir        string
	sqlDumpFileName   string
	sqlClientCmdPath  string
	sqlRestoreMarkDir string
}

func setDefaultSqlConfig() sqlBackupConfig {
	return sqlBackupConfig{
		sqlUser:           "root",
		sqlAddr:           "127.0.0.1",
		sqlPort:           "3306",
		sqlDatabase:       "postgres",
		sqlDumpCmdPath:    "/usr/bin/mysqldump",
		sqlDumpDir:        "/tmp/sqldump/",
		sqlDumpFileName:   "sql.dump",
		sqlClientCmdPath:  "/usr/bin/mysql",
		sqlRestoreMarkDir: "/tmp/mysqlrestore/",
	}
}

func getSqlConfig(sqltype string) (sqlBackupParams sqlBackupConfig, configvalid bool) {
	var currentSqlConfig sqlBackupConfig
	var isSqlConfgValid bool

	currentSqlConfig = setDefaultSqlConfig()

	isSqlConfgValid = true

	switch sqltype {

	case "mysql":

		currentSqlConfig.sqlClientCmdPath = "/usr/bin/mysql"
		currentSqlConfig.sqlDumpCmdPath = "/usr/bin/mysqldump"

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

	case "pgsql":

		currentSqlConfig.sqlClientCmdPath = "/usr/bin/psql"
		currentSqlConfig.sqlDumpCmdPath = "/usr/bin/pg_dump"
		currentSqlConfig.sqlRestoreMarkDir = "/tmp/pgsqlrestore/"

		if sqlUserenv := os.Getenv("PGSQL_USER"); sqlUserenv != "" {
			currentSqlConfig.sqlUser = sqlUserenv
		} else {
			log.Println("PGSQL_USER environment variable is not set, using the default", currentSqlConfig.sqlUser)
		}

		if sqlPasswordenv := os.Getenv("PGSQL_PASSWORD"); sqlPasswordenv != "" {
			currentSqlConfig.sqlPassword = sqlPasswordenv
		} else {
			log.Println("PGSQL_PASSWORD environment variable is not set, this is required to perform PGSQL backup")
			isSqlConfgValid = false
		}

		if sqlDatabaseenv := os.Getenv("PGSQL_DATABASE"); sqlDatabaseenv != "" {
			currentSqlConfig.sqlDatabase = sqlDatabaseenv
		} else {
			log.Println("PGSQL_DATABASE environment variable is not set, using the default", currentSqlConfig.sqlDatabase)
		}

		if sqlHostenv := os.Getenv("PGSQL_HOST"); sqlHostenv != "" {
			currentSqlConfig.sqlAddr = sqlHostenv
		} else {
			log.Println("PGSQL_HOST environment variable is not set, using the default", currentSqlConfig.sqlAddr)
		}

		if sqlPortenv := os.Getenv("PGSQL_PORT"); sqlPortenv != "" {
			currentSqlConfig.sqlPort = sqlPortenv
		} else {
			currentSqlConfig.sqlPort = "5432"
			log.Println("PGSQL_PORT environment variable is not set, using the default", currentSqlConfig.sqlPort)

		}

	}

	return currentSqlConfig, isSqlConfgValid
}

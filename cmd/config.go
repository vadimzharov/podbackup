package cmd

import (
	"log"
	"os"
	"strconv"

	//	"strings"
	"time"
)

type backupConfig struct {
	s3Endpoint        string
	bucketName        string
	awsRegion         string
	backupDir         string
	restoreDir        string
	bucketFolder      string
	keyPrefix         string
	backupLocalFile   string
	backupInverval    time.Duration
	pruneInverval     time.Duration
	filesKeep         int
	forceRestore      bool
	archiveType       string
	s3SyncParallelism int
	s3CopyBeforeSync  bool
	s3CopyWithDelete  bool
}

type backupCreds struct {
	awsKey          string
	awsSecretKey    string
	encryptpassword string
}

func setDefaultConfig() backupConfig {
	return backupConfig{
		s3Endpoint:        "",
		awsRegion:         "us-east-1",
		bucketFolder:      "podbackup",
		keyPrefix:         "podbackup",
		backupLocalFile:   "backup.zip",
		backupInverval:    3600000000000,
		pruneInverval:     7200000000000,
		filesKeep:         3,
		forceRestore:      false,
		archiveType:       "zip",
		s3SyncParallelism: 3,
		s3CopyBeforeSync:  false,
		s3CopyWithDelete:  false,
	}
}

func getConfig() (backupConfigParams backupConfig, backupCredentials backupCreds, configvalid bool) {
	var currentConfig backupConfig
	var currentCreds backupCreds
	var isConfigValid bool
	isConfigValid = true

	// Set default values
	currentConfig = setDefaultConfig()

	// Check mandatory values

	if awskeyenv := os.Getenv("AWS_KEY"); awskeyenv != "" {
		currentCreds.awsKey = awskeyenv
	} else {
		isConfigValid = false
		log.Println("AWS_KEY environment variable is not set, config is invalid")
	}

	if awssecretkeyenv := os.Getenv("AWS_SECRET_KEY"); awssecretkeyenv != "" {
		currentCreds.awsSecretKey = awssecretkeyenv
	} else {
		isConfigValid = false
		log.Println("AWS_SECRET_KEY environment variable is not set, config is invalid")
	}

	if encryptpasswordenv := os.Getenv("ENCRYPT_PASSWORD"); encryptpasswordenv != "" {
		currentCreds.encryptpassword = encryptpasswordenv
		log.Println("Using encryption...")
	} else {
		log.Println("ENCRYPT_PASSWORD environment variable is not set, encryption is not using")
		currentCreds.encryptpassword = ""
	}

	if bucketnameenv := os.Getenv("AWS_BUCKET"); bucketnameenv != "" {
		currentConfig.bucketName = bucketnameenv
	} else {
		isConfigValid = false
		log.Println("AWS_BUCKET environment variable is not set, config is invalid")
	}

	if backupdirenv := os.Getenv("DIR_TO_BACKUP"); backupdirenv != "" {
		currentConfig.backupDir = backupdirenv
	}

	if restoredirenv := os.Getenv("DIR_TO_RESTORE"); restoredirenv != "" {
		currentConfig.restoreDir = restoredirenv
	}

	// Check optional values

	if awsregionenv := os.Getenv("AWS_REGION"); awsregionenv != "" {
		currentConfig.awsRegion = awsregionenv
	} else {
		log.Println("AWS_REGION environment variable is not set, using the default", currentConfig.awsRegion)
	}

	if s3endpointenv := os.Getenv("S3_ENDPOINT"); s3endpointenv != "" {
		currentConfig.s3Endpoint = s3endpointenv
	} else {
		log.Println("S3_ENDPOINT environment variable is not set, using AWS default")
	}
	if s3syncparenv := os.Getenv("S3_SYNC_PARALLELISM"); s3syncparenv != "" {
		currentConfig.s3SyncParallelism, _ = strconv.Atoi(s3syncparenv)
	} else {
		log.Println("S3_SYNC_PARALLELISM environment variable is not set, using the default", currentConfig.s3SyncParallelism)
	}

	if backupbucketfolderenv := os.Getenv("S3_BUCKET_FOLDER"); backupbucketfolderenv != "" {
		currentConfig.bucketFolder = backupbucketfolderenv + "/"
	} else {
		log.Println("S3_BUCKET_FOLDER environment variable is not set, using the default", currentConfig.bucketFolder)
	}

	if backupkeyprefixenv := os.Getenv("S3_FILE_PREFIX"); backupkeyprefixenv != "" {
		currentConfig.keyPrefix = backupkeyprefixenv
	} else {
		log.Println("S3_FILE_PREFIX environment variable is not set, using the default", currentConfig.keyPrefix)
	}

	if backupinvervalenv := os.Getenv("BACKUP_INTERVAL"); backupinvervalenv != "" {

		currentConfig.backupInverval = parsedInterval(backupinvervalenv)

	} else {
		log.Println("BACKUP_INTERVAL environment variable is not set, using the default", currentConfig.backupInverval)

	}

	if filesKeepenv := os.Getenv("COPIES_TO_KEEP"); filesKeepenv != "" {
		currentConfig.filesKeep, _ = strconv.Atoi(filesKeepenv)
	} else {
		log.Println("COPIES_TO_KEEP environment variable is not set, using the default", currentConfig.filesKeep)
	}

	if pruneinvervalenv := os.Getenv("PRUNE_INTERVAL"); pruneinvervalenv != "" {
		currentConfig.pruneInverval = parsedInterval(pruneinvervalenv)
	} else {
		log.Println("PRUNE_INTERVAL environment variable is not set, using the default", currentConfig.pruneInverval)
	}

	if forcerestoreenv := os.Getenv("FORCE_RESTORE"); forcerestoreenv != "" {
		currentConfig.forceRestore, _ = strconv.ParseBool(forcerestoreenv)
	} else {
		log.Println("FORCE_RESTORE environment variable is not set, using the default", currentConfig.forceRestore)
	}

	if s3CopyBeforeSyncenv := os.Getenv("S3_COPY_BEFORE_SYNC"); s3CopyBeforeSyncenv != "" {
		currentConfig.s3CopyBeforeSync, _ = strconv.ParseBool(s3CopyBeforeSyncenv)
	} else {
		log.Println("S3_COPY_BEFORE_SYNC environment variable is not set, using the default", currentConfig.s3CopyBeforeSync)
	}

	if s3CopyWithDeleteenv := os.Getenv("S3_COPY_WITH_DELETE"); s3CopyWithDeleteenv != "" {
		currentConfig.s3CopyWithDelete, _ = strconv.ParseBool(s3CopyWithDeleteenv)
	} else {
		log.Println("S3_COPY_WITH_DELETE environment variable is not set, using the default", currentConfig.s3CopyWithDelete)
	}

	if archivetypeenv := os.Getenv("ARCHIVE_TYPE"); archivetypeenv != "" {
		currentConfig.archiveType = archivetypeenv
		log.Println("Using", currentConfig.archiveType, " archive type")
	} else {
		log.Println("ARCHIVE_TYPE environment variable is not set, using the default", currentConfig.archiveType)
	}

	log.Println("The following configuration parameters will be used:")
	log.Printf("%+v\n", currentConfig)

	return currentConfig, currentCreds, isConfigValid
}

func parsedInterval(cfgInterval string) time.Duration {

	var interval time.Duration
	var convinterval int

	convinterval, converr := strconv.Atoi(cfgInterval)

	if converr == nil {

		interval = time.Duration(convinterval) * time.Second

	} else {

		intinterval, _ := strconv.Atoi(cfgInterval[:len(cfgInterval)-1])

		switch cfgInterval[len(cfgInterval)-1:] {

		case "m":
			{
				interval = time.Duration(intinterval) * time.Minute
			}

		case "h":
			{
				interval = time.Duration(intinterval) * time.Hour
			}

		default:
			{
				log.Panic("Cannot convert interval variable", cfgInterval, ". Check BACKUP_INTERVAL or PRUNE_INTERVAL environment variables")
				os.Exit(1)
			}

		}

	}

	return interval

}

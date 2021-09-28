package cmd

import (
	"log"
	"os"
	"strconv"
)

type backupConfig struct {
	bucketName      string
	awsRegion       string
	backupDir       string
	restoreDir      string
	bucketFolder    string
	keyPrefix       string
	backupLocalFile string
	backupInverval  int
	pruneInverval   int
	filesKeep       int
	forceRestore    bool
	useTar          bool
}

type backupCreds struct {
	awsKey          string
	awsSecretKey    string
	encryptpassword string
}

func setDefaultConfig() backupConfig {
	return backupConfig{
		awsRegion:       "us-east-1",
		bucketFolder:    "podbackup",
		keyPrefix:       "podbackup",
		backupLocalFile: "backup.zip",
		backupInverval:  3600,
		pruneInverval:   6000,
		filesKeep:       3,
		forceRestore:    false,
		useTar:          false,
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
		currentConfig.backupInverval, _ = strconv.Atoi(backupinvervalenv)
	} else {
		log.Println("BACKUP_INTERVAL environment variable is not set, using the default", currentConfig.backupInverval)
	}

	if filesKeepenv := os.Getenv("COPIES_TO_KEEP"); filesKeepenv != "" {
		currentConfig.filesKeep, _ = strconv.Atoi(filesKeepenv)
	} else {
		log.Println("COPIES_TO_KEEP environment variable is not set, using the default", currentConfig.filesKeep)
	}

	if pruneinvervalenv := os.Getenv("PRUNE_INTERVAL"); pruneinvervalenv != "" {
		currentConfig.pruneInverval, _ = strconv.Atoi(pruneinvervalenv)
	} else {
		log.Println("PRUNE_INTERVAL environment variable is not set, using the default", currentConfig.pruneInverval)
	}

	if forcerestoreenv := os.Getenv("FORCE_RESTORE"); forcerestoreenv != "" {
		currentConfig.forceRestore, _ = strconv.ParseBool(forcerestoreenv)
	} else {
		log.Println("FORCE_RESTORE environment variable is not set, using the default", currentConfig.forceRestore)
	}

	if usetarenv := os.Getenv("USE_TAR"); usetarenv != "" {
		currentConfig.useTar, _ = strconv.ParseBool(usetarenv)
		log.Println("USE_TAR flag is set - using TAR to make archive")
	} else {
		log.Println("USE_TAR environment variable is not set, using the default", currentConfig.useTar)
	}

	log.Println("The following configuration parameters will be used:")
	log.Printf("%+v\n", currentConfig)

	return currentConfig, currentCreds, isConfigValid
}

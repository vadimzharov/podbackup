package cmd

import (
	"fmt"
	"os"
)

const (
	helpCmd = `
	Tool to backup files and subdirectories from any local directory as ZIP archive and upload it to S3 bucket.
	For tool to work set the following mandatory environment variables:
	AWS_BUCKET - S3 bucket to use.
	AWS_KEY - key to use to access to the bucket
	AWS_SECRET_KEY - secret key to use to access to the bucket
	DIR_TO_BACKUP - absolute path for directory to backup (tool will backup all files and subdirectories inside it)
	DIR_TO_RESTORE - absolute path for directory to restore into

	Optionally set the following variables:
	S3_BUCKET_FOLDER - folder where to store ZIP archive. "podbackup" by default

	S3_FILE_PREFIX - ZIP archive name prefix. "podbackup" by default. Full filename will be <prefix>-<timestamp>.zip

	ENCRYPT_PASSWORD - encrypt/decrypt ZIP archives using this password. 

	BACKUP_INTERVAL - interval in seconds to run periodical backup (if running as daemon). 3600 seconds by default.

	COPIES_TO_KEEP - number of copies to keep in S3 folder when executing pruning.

	FORCE_RESTORE - set to True if requied tool to fail (exit with code 1) if it cannot restore files from backup.
	
	ARCHIVE_TYPE - 	by default set to 'zip' - tool will create ZIP archive (and encrypt it if ENCRYPT_PASSWORD is set). 
	            	Set to 'tarzip' - to archive all files as tar archive and then zip it (encrypted if ENCRYPT_PASSWORD is set).
			Use it if you need to save original ownership and mode of the files.
			Set to 'targz' - to archive all files as tar compressed archive. File mode and ownership persist during unpacking, however encryption is not supported.
	
	Commands:
	backup			run one time backup

	backup-daemon		work as daemon and run periodical backups according to BACKUP_INTERVAL environment variable (3600 seconds by default).
				In this mode daemon will do automatic pruning and keep only # of copies based on COPIES_TO_KEEP environment variable (3 by default)

	prune			manually run pruning (delete all old archives)

	list			list files in S3 folder (based on S3_BUCKET_FOLDER environment variable)
	
	restore			download file from S3 and restore files to directory (DIR_TO_RESTORE environment variable). Most recent archive will be used. 
				To restore from another file provide archive name based on 'podbackup list' output (like podbackup/podbackup-20210802213807.zip)
	`
)

func printHelp() {

	fmt.Println(helpCmd)
	os.Exit(0)

}

package cmd

func pruneCosObjects() {

	filesList := listBackups(currentConfig.bucketFolder, currentConfig.keyPrefix, currentConfig.bucketName)

	prune(filesList, currentConfig.filesKeep, currentConfig.bucketName)
}

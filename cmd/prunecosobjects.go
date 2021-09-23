package cmd

func pruneCosObjects() {

	filesList := listBackups(currentConfig.bucketFolder, currentConfig.keyPrefix, currentConfig.bucketName, currentCreds.awsKey, currentCreds.awsSecretKey, currentConfig.awsRegion)

	prune(filesList, currentConfig.filesKeep, currentConfig.bucketName, currentCreds.awsKey, currentCreds.awsSecretKey, currentConfig.awsRegion)
}

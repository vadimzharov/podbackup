package cmd

import (
	"fmt"
	"log"
)

func listCosFiles() {

	filesList := listBackups(currentConfig.bucketFolder, currentConfig.keyPrefix, currentConfig.bucketName)

	if filesList == nil {
		log.Println("Cannot list files in bucket", currentConfig.bucketName)
	} else {
		log.Println("List of files, sorted in descending order based on time created:")

		for _, num := range filesList {
			fmt.Println(num)
		}
	}
}

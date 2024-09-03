package cmd

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Function makeTarBackup creates a tar archive from the specified <backupDirPath> directory as specified file <backupFileName>.tar

func makeTarBackup(backupDirPath string, backupFileName string) {

	// Clean up

	err := os.RemoveAll(filepath.Dir(backupFileName))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Create directory to store file with backup

	err = os.Mkdir(filepath.Dir(backupFileName), os.ModePerm)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("Making TAR archive from", backupDirPath, "to", backupFileName)

	// Create file with backup

	localArchive, err := os.Create(backupFileName)

	if err != nil {
		log.Println("Failed to create local file with archive", backupFileName, err)
		panic(err)
	}

	defer localArchive.Close()

	tarWriter := tar.NewWriter(localArchive)

	if currentConfig.archiveType == "targz" {

		gzWriter := gzip.NewWriter(localArchive)

		defer gzWriter.Close()

		tarWriter = tar.NewWriter(gzWriter)
	}

	defer tarWriter.Close()

	filepath.Walk(backupDirPath, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if currentConfig.skipPath != "" && strings.Contains(filepath.Join(file), currentConfig.skipPath) {
			println("Skipping file ", filepath.Join(file), "because it is in SKIP_PATH")
		} else {

			// Get file header
			header, err := tar.FileInfoHeader(fi, fi.Name())

			if err != nil {
				return err
			}

			// Modify header to backup relative path
			header.Name = filepath.Join(strings.TrimPrefix(file, backupDirPath))

			hdrerr := tarWriter.WriteHeader(header)

			if hdrerr != nil {
				panic(err)
			}

			if fi.IsDir() {
				return nil
			}

			f1, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			defer f1.Close()

			if _, err := io.Copy(tarWriter, f1); err != nil {
				panic(err)
			}

		}

		return nil

	})

}

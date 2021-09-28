package cmd

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func makeTarBackup(backupdirpath string, backupfilename string) {

	err := os.RemoveAll(filepath.Dir(backupfilename))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = os.Mkdir(filepath.Dir(backupfilename), os.ModePerm)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("Making TAR archive from", backupdirpath, "to", backupfilename)

	localArchive, err := os.Create(backupfilename)
	if err != nil {
		log.Println("Failed to create localfile ", backupfilename, err)
		panic(err)
	}

	tarWriter := tar.NewWriter(localArchive)

	defer localArchive.Close()

	filepath.Walk(backupdirpath, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// Get file header
		header, err := tar.FileInfoHeader(fi, fi.Name())

		if err != nil {
			return err
		}

		// Modify header to backup relative path
		header.Name = filepath.Join(strings.TrimPrefix(file, backupdirpath))

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

		return nil
	})

	tarWriter.Close()

}

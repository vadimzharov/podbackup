package cmd

import (
	"github.com/alexmullins/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func makeBackup(backupdirpath string, backupfilename string, encryptpassword string) {

	log.Println("Compressing all files from", backupdirpath, "to", backupfilename)

	encryptArchive := false

	if encryptpassword != "" {
		encryptArchive = true
		log.Println("Encrypt password is set, using encryption")
	}

	localArchive, err := os.Create(backupfilename)
	if err != nil {
		log.Println("Failed to create localfile ", backupfilename, err)
		panic(err)
	}

	defer localArchive.Close()

	zipWriter := zip.NewWriter(localArchive)

	filepath.Walk(backupdirpath, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// Get file header
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// Modify header to backup relative path
		header.Name = filepath.Join(strings.TrimPrefix(file, backupdirpath))

		// if it is directory - add slash to archive properly
		if fi.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		if fi.IsDir() {
			return nil
		}

		f1, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f1.Close()

		// Add encrypt password if need to encrypt

		if encryptArchive {

			header.SetPassword(encryptpassword)

		}

		ar_file, err := zipWriter.CreateHeader(header)

		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(ar_file, f1); err != nil {
			panic(err)
		}

		return nil
	})

	zipWriter.Close()

}

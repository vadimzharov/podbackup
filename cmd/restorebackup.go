package cmd

import (
	"fmt"
	"github.com/alexmullins/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func restoreBackup(restoredir string, backupfilename string, encryptpassword string) ([]string, error) {

	log.Println("Restoring backup from file", backupfilename, "to", restoredir, "directory")

	encryptArchive := false

	if encryptpassword != "" {
		encryptArchive = true
		log.Println("Encrypt password is set, trying to use it to decrypt files (password applied only if file is encrypted)")
	}

	var fileNames []string

	localArchive, err := zip.OpenReader(backupfilename)
	if err != nil {
		return fileNames, err
	}
	defer localArchive.Close()

	os.MkdirAll(restoredir, 0755)

	for _, f := range localArchive.File {

		if (f.IsEncrypted()) && (!encryptArchive) {
			log.Println("Error: files are encrypted but encrypt password is not set (ENCRYPT_PASSWORD variable)")
			return fileNames, nil
		}

		if (f.IsEncrypted()) && (encryptArchive) {
			f.SetPassword(encryptpassword)
		}

		// Store filename/path for returning and using later on
		fpath := filepath.Join(restoredir, f.Name)

		log.Println("Restoring", fpath)

		if !strings.HasPrefix(fpath, filepath.Clean(restoredir)+string(os.PathSeparator)) {
			return fileNames, fmt.Errorf("%s: illegal file path", fpath)
		}

		fileNames = append(fileNames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fileNames, err
		}

		rc, err := f.Open()
		if err != nil {
			log.Println("Error during upacking:")
			log.Println(err)
			return fileNames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fileNames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return fileNames, err
		}
	}

	return fileNames, nil
}

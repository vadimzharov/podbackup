package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func restoreTarBackup(restoredir string, backupfilename string) ([]string, error) {

	log.Println("Restoring backup from file", backupfilename, "to", restoredir, "directory")

	var fileNames []string

	localTarArchive, err := os.Open(backupfilename)
	if err != nil {
		return fileNames, err
	}
	defer localTarArchive.Close()

	tarReader := tar.NewReader(localTarArchive)

	if currentConfig.archiveType == "targz" {

		gzReader, err := gzip.NewReader(localTarArchive)

		if err != nil {
			return fileNames, err
		}

		defer gzReader.Close()

		tarReader = tar.NewReader(gzReader)
	}

	os.MkdirAll(restoredir, os.ModePerm)

	for {

		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error during upacking, cannot read file:")
			log.Println(err)
			return fileNames, err
		}

		fi := header.FileInfo()

		fpath := filepath.Join(restoredir, header.Name)

		fileNames = append(fileNames, fpath)

		if fi.Mode().IsDir() {

			os.MkdirAll(fpath, 0755)

			log.Println("Creating directory", fpath)

			err = os.Chown(fpath, header.Uid, header.Gid)

			if err != nil {
				log.Println("Cannot set ownership", header.Uid, ":", header.Gid, "to file", fpath)
				log.Println(err)
			}

			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())

		log.Println("Restoring", fpath)

		if !strings.HasPrefix(fpath, filepath.Clean(restoredir)+string(os.PathSeparator)) {
			return fileNames, fmt.Errorf("%s: illegal file path", fpath)
		}

		if err != nil {
			log.Println("Error during upacking file", fpath, ":")
			log.Println(err)
			return fileNames, err
		}

		_, err = io.Copy(outFile, tarReader)

		if err != nil {
			log.Println("Error during upacking:")
			log.Println(err)
			return fileNames, err
		}

		// Close the file without defer to close before next iteration of loop
		outFile.Close()

		err = os.Chown(fpath, header.Uid, header.Gid)
		if err != nil {
			log.Println("Cannot set ownership", header.Uid, ":", header.Gid, "to file", fpath)
			log.Println(err)
		}
		err = os.Chmod(fpath, fi.Mode())
		if err != nil {
			log.Println("Cannot set permissions", fi.Mode(), "to file", fpath)
			log.Println(err)
		}

	}

	return fileNames, nil
}

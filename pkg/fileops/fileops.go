package fileops

import (
	"filesynctool/models"
	"io"
	"log"
	"os"
)

func CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Println("Failed to copy file:", src)
		return err
	}

	destFile, err := os.Create(models.Destination + dest)
	if err != nil {
		log.Println("Failed to copy file:", src)
		return err
	}

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		log.Println("Failed to copy file:", src)
		return err
	}
	return err
}

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

func DeleteFile(src string) error {
	if err := os.Remove(models.Destination + src); err != nil {
		log.Println("Failed to delete file:", src)
		return err
	}
	return nil
}

func RemoveDir(src string) error {
	if err := os.RemoveAll(models.Destination + src); err != nil {
		log.Println("Failed to delete dir:", src)
		return err
	}
	return nil
}

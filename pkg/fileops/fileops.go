package fileops

import (
	"filesynctool/models"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destPath := filepath.Join(models.Destination, dest)
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func DeleteFile(src string) error {
	err := os.Remove(filepath.Join(models.Destination, src))
	return err
}

func RemoveDir(src string) error {
	err := os.RemoveAll(filepath.Join(models.Destination, src))
	return err
}

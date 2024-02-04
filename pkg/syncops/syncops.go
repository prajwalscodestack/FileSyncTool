package syncops

import (
	"filesynctool/models"
	"filesynctool/pkg/fileops"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SyncFileCreate(syncChan chan models.FileEntry) {
	for {
		file := <-syncChan
		dst := strings.Replace(file.Path, models.Source, "", 1)
		if err := os.MkdirAll(models.Destination+filepath.Dir(dst), 0755); err != nil {
			log.Println("Failed to create directory:", err)
		}
		fmt.Println("New File", file.FileInfo.Name(), file.Path)
		// Do something with the file here
		if err := fileops.CopyFile(file.Path, dst); err != nil {
			log.Println("Failed to copy file:", err)
		}
	}
}

func SyncFileUpdate(syncChan chan models.FileEntry) {
	for {
		file := <-syncChan
		fmt.Println("Update File:", file.FileInfo.Name(), file.Path)
		// Do something with the file here
		dest := strings.Replace(file.Path, models.Source, "", 1)
		if err := fileops.CopyFile(file.Path, dest); err != nil {
			log.Println("Failed to sync file:", file.FileInfo.Name(), err)
		}
	}
}

func SyncFileDelete(syncChan chan models.FileEntry) {
	for {
		deletedFile := <-syncChan
		fmt.Println("Delete File:", deletedFile.Path)
		dst := strings.Replace(deletedFile.Path, models.Source, "", 1)
		if _, err := os.Stat(models.Destination + dst); err == nil {
			if deletedFile.FileInfo.IsDir() {
				if err := fileops.RemoveDir(dst); err != nil {
					log.Println("Failed to delete directory:", err)
				}
			}
			if deletedFile.FileInfo.Mode().IsRegular() {
				if err := fileops.DeleteFile(dst); err != nil {
					log.Println("Failed to delete file:", deletedFile.Path, err)
				}
			}
		}
	}
}

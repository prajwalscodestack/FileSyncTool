package syncops

import (
	"filesynctool/models"
	"filesynctool/pkg/fileops"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func SyncFileCreate(syncChan chan models.FileEntry) {
	for {
		file := <-syncChan
		dst, err := filepath.Rel(models.Source, file.Path)
		if err != nil {
			log.Println("Failed to get relative path:", err)
		}
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
		dst, err := filepath.Rel(models.Source, file.Path)
		if err != nil {
			log.Println("Failed to get relative path:", err)
		}
		if err := fileops.CopyFile(file.Path, dst); err != nil {
			log.Println("Failed to sync file:", file.FileInfo.Name(), err)
		}
	}
}

func SyncFileDelete(syncChan chan models.FileEntry) {
	for {
		file := <-syncChan
		fmt.Println("Delete File:", file.Path)
		dst, err := filepath.Rel(models.Source, file.Path)
		if err != nil {
			log.Println("Failed to get relative path:", err)
		}
		if _, err := os.Stat(models.Destination + dst); err == nil {
			if file.FileInfo.IsDir() {
				if err := fileops.RemoveDir(dst); err != nil {
					log.Println("Failed to delete directory:", err)
				}
			}
			if file.FileInfo.Mode().IsRegular() {
				if err := fileops.DeleteFile(dst); err != nil {
					log.Println("Failed to delete file:", file.Path, err)
				}
			}
		}
	}
}

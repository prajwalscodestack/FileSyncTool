package main

import (
	"filesynctool/models"
	"filesynctool/pkg/fileops"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SyncFileCreate(syncChan chan FileEntry) {
	for {
		file := <-syncChan
		dst := strings.Replace(file.Path, models.Source, "", 1)
		if err := os.MkdirAll(models.Destination+filepath.Dir(dst), 0755); err != nil {
			log.Println("Failed to create directory:", err)
		}
		fmt.Println("New File", file.FileInfo.Name(), file.Path)
		// Do something with the file here
		sourceFile, err := os.Open(file.Path)
		if err != nil {
			log.Println("Failed to sync file:", file.FileInfo.Name())
		}

		destFile, err := os.Create(models.Destination + dst)
		if err != nil {
			log.Println("Failed to sync file:", file.FileInfo.Name())
		}

		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			log.Println("Failed to sync file:", file.FileInfo.Name())
		}
	}
}

func SyncFileUpdate(syncChan chan FileEntry) {
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

// TODO: handle delete operation with directory
func SyncFileDelete(syncChan chan FileEntry) {
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

type FileEntry struct {
	FileInfo os.FileInfo
	Path     string
}
type FileMod struct {
	FileEntry FileEntry
	ModTime   time.Time
}

func main() {
	models.SetFlags()
	// Define the directory to watch
	//validate path
	if models.Source == "" || models.Destination == "" {
		panic("source or destination path can't be empty")
	}
	syncFileCreate := make(chan FileEntry)
	syncFileUpdate := make(chan FileEntry)
	syncFileDelete := make(chan FileEntry)

	// Create a map to store file modification times
	fileModTimes := make(map[string]FileMod)

	go SyncFileCreate(syncFileCreate)
	go SyncFileUpdate(syncFileUpdate)
	go SyncFileDelete(syncFileDelete)

	// Start an infinite loop to check for changes
	//TODO: should detect whether whole directory is deleted
	for {
		entries := make([]FileEntry, 0)
		err := filepath.Walk(models.Source, func(path string, entry os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			entries = append(entries, FileEntry{
				Path:     path,
				FileInfo: entry,
			})
			return nil
		})
		if err != nil {
			fmt.Println("Error walking directory:", err)
		}
		//exsting file
		existingFiles := make(map[string]bool)

		// Check for new or modified files
		for _, fileEntry := range entries {
			existingFiles[fileEntry.Path] = true
			if fileEntry.FileInfo.Mode().IsRegular() {
				// Check if file exists in the map
				if filemod, ok := fileModTimes[fileEntry.Path]; ok {
					// Compare modification times
					if filemod.ModTime != fileEntry.FileInfo.ModTime() {
						syncFileUpdate <- fileEntry
						// Update modification time in the map
						fileModTimes[fileEntry.Path] = FileMod{
							FileEntry: fileEntry,
							ModTime:   fileEntry.FileInfo.ModTime(),
						}
					}
				} else {
					syncFileCreate <- fileEntry
					// Add new file to the map
					fileModTimes[fileEntry.Path] = FileMod{
						FileEntry: fileEntry,
						ModTime:   fileEntry.FileInfo.ModTime(),
					}
				}
			}
			if fileEntry.FileInfo.IsDir() {
				fileModTimes[fileEntry.Path] = FileMod{
					FileEntry: fileEntry,
					ModTime:   fileEntry.FileInfo.ModTime(),
				}
			}
		}
		//check for deleted files
		for path, filemod := range fileModTimes {
			if _, ok := existingFiles[path]; !ok {
				syncFileDelete <- FileEntry{
					Path:     path,
					FileInfo: filemod.FileEntry.FileInfo,
				}
				delete(fileModTimes, path)
				delete(existingFiles, path)
			}
		}
		//make entries nil
		entries = nil

		// // Sleep for a while before checking again
		time.Sleep(1 * time.Second)
	}
}

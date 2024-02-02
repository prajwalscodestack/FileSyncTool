package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var source = "./source"
var destination = "./destination"

func SyncFileCreate(syncChan chan FileEntry) {
	for {
		file := <-syncChan
		actualFilePath := strings.Replace(file.Path, "source", "", 1)
		if err := os.MkdirAll(destination+filepath.Dir(actualFilePath), 0755); err != nil {
			log.Println("Failed to create directory:", err)
		}
		fmt.Println("New File", file.FileInfo.Name(), file.Path)
		// Do something with the file here
		sourceFile, err := os.Open(file.Path)
		if err != nil {
			log.Println("Failed to sync file:", file.FileInfo.Name())
		}

		destFile, err := os.Create(destination + actualFilePath)
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
	}
}

type FileEntry struct {
	FileInfo os.FileInfo
	Path     string
}

func main() {
	// Define the directory to watch

	syncFileCreate := make(chan FileEntry)
	syncFileUpdate := make(chan FileEntry)
	// syncFileDelete := make(chan FileEntry)

	// Create a map to store file modification times
	fileModTimes := make(map[string]time.Time)

	go SyncFileCreate(syncFileCreate)
	go SyncFileUpdate(syncFileUpdate)
	// Start an infinite loop to check for changes
	for {
		entries := make([]FileEntry, 0)
		err := filepath.Walk(source, func(path string, entry os.FileInfo, err error) error {
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

		// // Check for new or modified files
		for _, fileEntry := range entries {
			if fileEntry.FileInfo.Mode().IsRegular() {
				// Check if file exists in the map
				if modTime, ok := fileModTimes[fileEntry.FileInfo.Name()]; ok {
					// Compare modification times
					if modTime != fileEntry.FileInfo.ModTime() {
						syncFileUpdate <- fileEntry
						// Update modification time in the map
						fileModTimes[fileEntry.FileInfo.Name()] = fileEntry.FileInfo.ModTime()
					}
				} else {
					syncFileCreate <- fileEntry
					// Add new file to the map
					fileModTimes[fileEntry.FileInfo.Name()] = fileEntry.FileInfo.ModTime()
				}
			}
		}
		//make entries nil
		entries = nil
		// // Sleep for a while before checking again
		time.Sleep(1 * time.Second)
	}
}

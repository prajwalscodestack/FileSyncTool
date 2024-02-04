package main

import (
	"filesynctool/models"
	"filesynctool/pkg/syncops"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func watch(syncFileCreate, syncFileUpdate, syncFileDelete chan models.FileEntry) {
	// Create a map to store file modification times
	fileModTimes := make(map[string]models.FileMod)
	// Start an infinite loop to check for changes
	for {
		entries := make([]models.FileEntry, 0)
		err := filepath.Walk(models.Source, func(path string, entry os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			entries = append(entries, models.NewFileEntry(path, entry))
			return nil
		})
		if err != nil {
			fmt.Println("Error walking directory:", err)
		}

		// Check for new or modified files
		for _, fileEntry := range entries {
			if fileEntry.FileInfo.Mode().IsRegular() {
				// Check if file exists in the map
				if filemod, ok := fileModTimes[fileEntry.Path]; ok {
					// Compare modification times
					if filemod.ModTime != fileEntry.FileInfo.ModTime() {
						syncFileUpdate <- fileEntry
						// Update modification time in the map
						fileModTimes[fileEntry.Path] = models.NewFileMod(fileEntry)
						models.NewFileMod(fileEntry)
					}
				} else {
					syncFileCreate <- fileEntry
					// Add new file to the map
					fileModTimes[fileEntry.Path] = models.NewFileMod(fileEntry)
				}
			}
			if fileEntry.FileInfo.IsDir() {
				fileModTimes[fileEntry.Path] = models.NewFileMod(fileEntry)
			}
		}
		//check for deleted files
		for path, filemod := range fileModTimes {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				syncFileDelete <- models.FileEntry{
					Path:     path,
					FileInfo: filemod.FileEntry.FileInfo,
				}
				delete(fileModTimes, path)
			}
		}
		// // Sleep for a while before checking again
		time.Sleep(1 * time.Second)
	}
}
func main() {
	models.SetFlags()
	// Define the directory to watch
	//validate path
	if models.Source == "" || models.Destination == "" {
		panic("source or destination path can't be empty")
	}
	syncFileCreate := make(chan models.FileEntry, models.SyncBuffer)
	syncFileUpdate := make(chan models.FileEntry, models.SyncBuffer)
	syncFileDelete := make(chan models.FileEntry, models.SyncBuffer)

	go syncops.SyncFileCreate(syncFileCreate)
	go syncops.SyncFileUpdate(syncFileUpdate)
	go syncops.SyncFileDelete(syncFileDelete)

	watch(syncFileCreate, syncFileUpdate, syncFileDelete)
}

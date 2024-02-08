package syncops

import (
	"filesynctool/models"
	"filesynctool/pkg/fileops"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Watcher struct {
	SyncFileCreate chan models.FileEntry
	SyncFileUpdate chan models.FileEntry
	SyncFileDelete chan models.FileEntry
	Location       string
}

func NewWatcher(location string) *Watcher {
	return &Watcher{
		SyncFileCreate: make(chan models.FileEntry),
		SyncFileUpdate: make(chan models.FileEntry),
		SyncFileDelete: make(chan models.FileEntry),
		Location:       location,
	}
}
func (w *Watcher) LauchSyncWorker() {
	go w.SyncCreate()
	go w.SyncUpdate()
	go w.SyncDelete()
}

// TODO: Should watch location provided and send the changes to broker
func (w *Watcher) Watch() {
	// Create a map to store file modification times
	fileModTimes := make(map[string]models.FileMod)
	// Start an infinite loop to check for changes
	for {
		entries := make([]models.FileEntry, 0)
		err := filepath.Walk(w.Location, func(path string, entry os.FileInfo, err error) error {
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
						w.SyncFileUpdate <- fileEntry
						// Update modification time in the map
						fileModTimes[fileEntry.Path] = models.NewFileMod(fileEntry)
						models.NewFileMod(fileEntry)
					}
				} else {
					w.SyncFileCreate <- fileEntry
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
				w.SyncFileDelete <- models.FileEntry{
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
func (w *Watcher) SyncCreate() {
	for {
		file := <-w.SyncFileCreate
		dst, err := filepath.Rel(w.Location, file.Path)
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

func (w *Watcher) SyncUpdate() {
	for {
		file := <-w.SyncFileUpdate
		fmt.Println("Update File:", file.FileInfo.Name(), file.Path)
		// Do something with the file here
		dst, err := filepath.Rel(w.Location, file.Path)
		if err != nil {
			log.Println("Failed to get relative path:", err)
		}
		if err := fileops.CopyFile(file.Path, dst); err != nil {
			log.Println("Failed to sync file:", file.FileInfo.Name(), err)
		}
	}
}

func (w *Watcher) SyncDelete() {
	for {
		file := <-w.SyncFileDelete
		fmt.Println("Delete File:", file.Path)
		dst, err := filepath.Rel(w.Location, file.Path)
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

package main

import (
	"fmt"
	"os"
	"time"
)

func SyncDestination(syncChan chan string) {
	for {
		file := <-syncChan
		fmt.Println("Syncing file:", file)
		// Do something with the file here
	}
}
func main() {
	// Define the directory to watch
	dir := "./source"
	syncChannel := make(chan string)
	// Create a map to store file modification times
	fileModTimes := make(map[string]time.Time)

	go SyncDestination(syncChannel)
	// Start an infinite loop to check for changes
	for {
		// Open the directory
		d, err := os.Open(dir)
		if err != nil {
			fmt.Println("Error opening directory:", err)
			return
		}
		defer d.Close()

		// Read the directory entries
		entries, err := d.Readdir(-1)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		// Check for new or modified files
		for _, entry := range entries {
			if entry.Mode().IsRegular() {
				// Check if file exists in the map
				if modTime, ok := fileModTimes[entry.Name()]; ok {
					// Compare modification times
					if modTime != entry.ModTime() {
						fmt.Println("File modified:", entry.Name())
						syncChannel <- entry.Name()
						// Update modification time in the map
						fileModTimes[entry.Name()] = entry.ModTime()
					}
				} else {
					// Add new file to the map
					fmt.Println("New file:", entry.Name())
					fileModTimes[entry.Name()] = entry.ModTime()
				}
			}
		}

		// Sleep for a while before checking again
		time.Sleep(1 * time.Second)
	}
}

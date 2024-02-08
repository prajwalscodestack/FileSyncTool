package main

import (
	"filesynctool/models"
	"filesynctool/pkg/syncops"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	models.Destination = "./test/destination"
	sigCh := make(chan os.Signal, 1)
	// Notify the sigCh channel for SIGINT (Ctrl+C) and SIGTERM (termination)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sourceWatcher := syncops.NewWatcher("./test/source")
	sourceWatcher.LauchSyncWorker()
	go sourceWatcher.Watch()
	Node1 := syncops.NewWatcher("./test/Node-1")
	Node1.LauchSyncWorker()
	go Node1.Watch()
	// Block until a signal is received
	fmt.Println("Application is running..")
	<-sigCh
}

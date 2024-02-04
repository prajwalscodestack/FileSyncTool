package models

import (
	"flag"
	"os"
	"time"
)

var (
	Source      = ""
	Destination = ""
	SyncBuffer  = 100
)

type FileEntry struct {
	os.FileInfo
	Path string
}
type FileMod struct {
	FileEntry FileEntry
	ModTime   time.Time
}

var src = flag.String("source", "", "considered as source location")
var dest = flag.String("destination", "", "considered as destination location")
var syncBuffer = flag.Int("syncBuffer", 100, "Specifies the buffer size used for synchronization operations.")

func SetFlags() {
	flag.Parse()
	Source = *src
	Destination = *dest
	SyncBuffer = *syncBuffer
}

func NewFileMod(fileInfo FileEntry) FileMod {
	return FileMod{
		FileEntry: fileInfo,
		ModTime:   fileInfo.FileInfo.ModTime(),
	}
}

func NewFileEntry(path string, fileInfo os.FileInfo) FileEntry {
	return FileEntry{
		FileInfo: fileInfo,
		Path:     path,
	}
}

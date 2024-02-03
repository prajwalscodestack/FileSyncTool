package models

import "flag"

var (
	Source      = ""
	Destination = ""
)

var src = flag.String("source", "", "considered as source location")
var dest = flag.String("destination", "", "considered as destination location")

func SetFlags() {
	flag.Parse()
	Source = *src
	Destination = *dest
}

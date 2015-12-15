package main

import (
	"github.com/curt-labs/heavyduty/importer"
	"log"
)

func main() {
	//TODO - those almost-universal 5th wheel parts
	err := importer.Get()
	if err != nil {
		log.Print(err)
	}
}

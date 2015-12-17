package main

import (
	"flag"
	"log"

	"github.com/curt-labs/heavyduty/importer"
	"github.com/curt-labs/heavyduty/merger"
)

var (
	skipImport = flag.Bool("skipimport", false, "Skip Import")
)

func main() {
	//TODO - those almost-universal 5th wheel parts
	var err error
	if *skipImport == false {
		err = importer.Get()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = merger.Merge()
	if err != nil {
		log.Fatal(err)
	}

}

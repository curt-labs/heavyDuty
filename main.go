package main

import (
	"log"

	"github.com/curt-labs/heavyduty/importer"
	"github.com/curt-labs/heavyduty/merger"
)

func main() {
	//TODO - those almost-universal 5th wheel parts
	var err error
	if 1 == 2 {

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

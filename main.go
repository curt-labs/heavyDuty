package main

import (
	"flag"
	"log"

	"github.com/curt-labs/heavierduty/importer"
)

var (
	skipImport = flag.Bool("skipimport", false, "Skip Import")
	skipMerge  = flag.Bool("skipmerge", false, "Skip Merge")
)

func main() {
	//TODO - those almost-universal 5th wheel parts

	if *skipImport == false {
		var err error
		vps, err := importer.GetDataStructure()
		if err != nil {
			log.Fatal(err)
		}

		vps, err = importer.MatchYears(vps)
		if err != nil {
			log.Fatal(err)
		}

		vps, err = importer.MatchMakes(vps)
		if err != nil {
			log.Fatal(err)
		}

		vps, err = importer.MatchModels(vps)
		if err != nil {
			log.Fatal(err)
		}

		vps, err = importer.MatchStyles(vps)
		if err != nil {
			log.Fatal(err)
		}

		vps, err = importer.MatchVehicles(vps)
		if err != nil {
			log.Fatal(err)
		}

		vps, err = importer.MatchVehicleParts(vps)
		if err != nil {
			log.Fatal(err)
		}

		err = importer.CreateRelatedParts(vps)
		if err != nil {
			log.Fatal(err)
		}

		err = importer.CreateStmts()
	}

}

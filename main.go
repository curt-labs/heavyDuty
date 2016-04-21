package main

import (
	"flag"
	"log"

	"github.com/curt-labs/heavierduty/deleter"
	"github.com/curt-labs/heavierduty/importer"
)

var (
	deleteRecords = flag.Bool("delete", false, "Do Delete Part Applications")
	insertRecords = flag.Bool("insert", false, "Do Insert Part Applications")
)

func main() {
	var err error
	flag.Parse()
	if *deleteRecords {
		err = deleteApplications()
		if err != nil {
			log.Fatal(err)
		}
	}

	if *insertRecords {
		err = insertApplications()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Print("END ", err)
	return
}

func deleteApplications() error {
	var err error
	ids, err := deleter.GetDataStructure()
	if err != nil {
		return err
	}
	vehiclePartsQuery := deleter.BuildDeleteVehiclePartsQuery(ids)
	relatedPartsQuery := deleter.BuildDeleteRelatedPartsQuery(ids)
	return deleter.FileOutput(vehiclePartsQuery, relatedPartsQuery)
}

func insertApplications() error {
	var err error
	vps, err := importer.GetDataStructure()
	if err != nil {
		return err
	}
	// log.Print(vps)

	vps, err = importer.MatchYears(vps)
	if err != nil {
		return err
	}

	vps, err = importer.MatchMakes(vps)
	if err != nil {
		return err
	}

	vps, err = importer.MatchModels(vps)
	if err != nil {
		return err
	}

	vps, err = importer.MatchStyles(vps)
	if err != nil {
		return err
	}

	vps, err = importer.MatchVehicles(vps)
	if err != nil {
		return err
	}

	vps, err = importer.MatchVehicleParts(vps)
	if err != nil {
		return err
	}

	err = importer.CreateRelatedParts(vps)
	if err != nil {
		return err
	}

	return importer.CreateStmts()
}

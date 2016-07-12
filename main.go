package main

import (
	"flag"
	"log"

	"github.com/curt-labs/heavierduty/csv_conversion"
	"github.com/curt-labs/heavierduty/deleter"
	"github.com/curt-labs/heavierduty/importer"
)

var (
	deleteRecords = flag.Bool("delete", false, "Do Delete Part Applications")
	insertRecords = flag.Bool("insert", false, "Do Insert Part Applications")
	checkRecords  = flag.Bool("check", false, "Do Insert Part Applications")
	path          = flag.String("path", "", "Path to csv")
)

var (
	// maps correspond to the rows of a company-provided csv
	DoubleLockMap  = map[string][]int{"Year": []int{0}, "Make": []int{1}, "Model": []int{2}, "Style": []int{3}, "PartNumber": []int{6}, "ShortDesc": []int{8}, "RelatedParts": []int{10, 11}, "Drilling": []int{21}, "Notes": []int{22, 23, 24, 25}, "InstallTime": []int{27}, "UPC": []int{28}, "List": []int{29}, "MAP": []int{30}, "Jobber": []int{31}, "Weight": []int{32}, "Length": []int{33}, "Height": []int{34}, "Width": []int{35}, "Bullets": []int{37, 38, 39, 40, 41}}
	FoldingBallMap = map[string][]int{"Year": []int{0}, "Make": []int{1}, "Model": []int{2}, "Style": []int{3}, "PartNumber": []int{6}, "ShortDesc": []int{8}, "RelatedParts": []int{10}, "Drilling": []int{20}, "Notes": []int{21, 22, 23, 24}, "InstallTime": []int{26}, "UPC": []int{27}, "List": []int{28}, "MAP": []int{29}, "Jobber": []int{30}, "Weight": []int{31}, "Length": []int{32}, "Height": []int{33}, "Width": []int{34}, "Bullets": []int{36, 37, 38, 39, 40}}
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

	if *checkRecords {
		err = checkApplications()
		if err != nil {
			log.Fatal(err)
		}
	}

	// if *insertRecords {
	// 	err = insertApplications()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	log.Print("END ", err)
	return
}

func checkApplications() error {
	var err error
	vps, err := csv_conversion.GetApplications(*path, DoubleLockMap)
	if err != nil {
		return err
	}

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

	vps, err = importer.ConfirmPartExistance(vps)
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

func deleteApplications() error {
	var err error
	ids, err := deleter.GetDataStructure(*path)
	if err != nil {
		return err
	}
	vehiclePartsQuery := deleter.BuildDeleteVehiclePartsQuery(ids)
	relatedPartsQuery := deleter.BuildDeleteRelatedPartsQuery(ids)
	relatedRelatedPartsQuery := deleter.BuildDeleteRelatedRelatedPartsQuery(ids)
	return deleter.FileOutput(vehiclePartsQuery, relatedPartsQuery, relatedRelatedPartsQuery)
}

// func insertApplications() error {
// 	var err error
// 	vps, err := importer.GetDataStructure(*path)
// 	if err != nil {
// 		return err
// 	}

// 	vps, err = importer.MatchYears(vps)
// 	if err != nil {
// 		return err
// 	}

// 	vps, err = importer.MatchMakes(vps)
// 	if err != nil {
// 		return err
// 	}

// 	vps, err = importer.MatchModels(vps)
// 	if err != nil {
// 		return err
// 	}

// 	vps, err = importer.MatchStyles(vps)
// 	if err != nil {
// 		return err
// 	}

// 	vps, err = importer.MatchVehicles(vps)
// 	if err != nil {
// 		return err
// 	}

// 	vps, err = importer.ConfirmPartExistance(vps)
// 	if err != nil {
// 		return err
// 	}

// 	vps, err = importer.MatchVehicleParts(vps)
// 	if err != nil {
// 		return err
// 	}

// 	err = importer.CreateRelatedParts(vps)
// 	if err != nil {
// 		return err
// 	}

// 	return importer.CreateStmts()
// }

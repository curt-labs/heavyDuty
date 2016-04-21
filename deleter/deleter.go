package deleter

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	path = flag.String("csv_path", "", "Path to csv")
)

// Get Unique PartIDs from Csv
func GetDataStructure() ([]int, error) {
	var ids []int
	idMap := make(map[int]int)
	var counter int
	flag.Parse()
	if *path == "" {
		*path = "Fifth Wheel Bracket 10.26.15.csv"
	}
	f, err := os.Open(*path)
	if err != nil {
		return ids, err
	}
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	rawData, err := reader.ReadAll()
	if err != nil {
		return ids, err
	}
	for i, line := range rawData {
		if i == 0 {
			continue //header
		}
		if line[0] == "" {
			break //end of file
		}

		// parse part
		partInt, err := strconv.Atoi(line[6])
		if err != nil {
			return ids, err
		}

		idMap[partInt] = partInt
		counter++
	}
	for id, _ := range idMap {
		ids = append(ids, id)
	}
	fmt.Println(counter, " rows examined")
	return ids, err
}

// Delete VehiclePart for those partIDs
func BuildDeleteVehiclePartsQuery(ids []int) string {
	var partIdStr string
	for i, id := range ids {
		if i > 0 {
			partIdStr += ","
		}
		partIdStr += strconv.Itoa(id)
	}

	query := fmt.Sprintf("delete from VehiclePart where partID in (%s)", partIdStr)
	return query
}

// Delete RelatedParts for those partIDs
func BuildDeleteRelatedPartsQuery(ids []int) string {
	var partIdStr string
	for i, id := range ids {
		if i > 0 {
			partIdStr += ","
		}
		partIdStr += strconv.Itoa(id)
	}

	query := fmt.Sprintf("delete from RelatedPart where partID in (%s)", partIdStr)
	return query
}

// Writes queries for you to execute in an Sql Client
func FileOutput(queries ...string) error {
	f, err := os.Create("delete_queries.txt")
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err = f.WriteString(query + ";\n")
		if err != nil {
			return err
		}
	}
	return nil
}

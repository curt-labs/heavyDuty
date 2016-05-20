package importer

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//Vehicle (style-based)
//The 'temp' fields indicate whether these attributes were entered in temp table
type Vehicle struct {
	ID      int
	Year    float64
	YearID  int
	Make    string
	MakeID  int
	Model   string
	ModelID int
	Style   string
	StyleID int
}

//VehiclePart links VehicleID to PartID
//The VehicleTemp field indicates whether this vehicle existed in the DB already
type VehiclePart struct {
	PartID       int
	Vehicle      Vehicle
	Drilling     string
	InstallTime  string
	RelatedParts []int
}

var (
	path           = flag.String("path", "", "Path to csv")
	yearMap        map[float64]int
	makeMap        map[string]int
	modelMap       map[string]int
	styleMap       map[string]int
	vehicleMap     map[string]int
	vehiclePartMap map[string]int
	relatePartsMap map[string]int
	partMap        map[int]int
)

// Init creates maps
func Init() {
	var err error
	//maps
	yearMap, err = getYearMap()
	if err != nil {
		log.Fatal(err)
	}
	makeMap, err = getMakeMap()
	if err != nil {
		log.Fatal(err)
	}
	modelMap, err = getModelMap()
	if err != nil {
		log.Fatal(err)
	}
	styleMap, err = getStyleMap()
	if err != nil {
		log.Fatal(err)
	}
	vehicleMap, vehiclePartMap, err = getVehiclePartMap()
	if err != nil {
		log.Fatal(err)
	}
	relatePartsMap, err = getRelatedPartsMap()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Maps obtained.")
}

// Get - reading csv and return Vehicleparts
func GetDataStructure() ([]VehiclePart, error) {
	// Init() //init maps
	var vps []VehiclePart
	var counter int
	flag.Parse()
	if *path == "" {
		*path = "Double Gooseneck Install Kits.csv"
	}
	f, err := os.Open(*path)
	if err != nil {
		return vps, err
	}
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	rawData, err := reader.ReadAll()
	if err != nil {
		return vps, err
	}
	for i, line := range rawData {
		if i == 0 {
			continue //header
		}
		if line[0] == "" {
			break //end of file
		}
		partToRelatedParts := make(map[int][]int)
		partToRelatedParts[7] = []int{10, 11, 12, 13, 14, 18, 19, 20, 21, 22, 23, 29, 30}
		partToRelatedParts[27] = []int{10, 11, 12, 13, 14, 18, 19, 20, 21, 22, 23}
		partToRelatedParts[28] = []int{10, 11, 12, 13, 14, 18, 19, 20, 21, 22, 23}
		for partIdCell, relatedPartCells := range partToRelatedParts {
			vp, err := parseVehiclePart(partIdCell, relatedPartCells, line)
			if err != nil {
				return vps, err
			}
			vps = append(vps, vp)
		}

		counter++
	}
	fmt.Println(counter, " vehicleParts examined")
	return vps, err
}

func parseVehiclePart(partIdCell int, relatedPartCells []int, line []string) (VehiclePart, error) {
	var err error
	//Base Vehicle
	vp := VehiclePart{
		Vehicle: Vehicle{
			Make:  strings.TrimSpace(line[1]),
			Model: strings.TrimSpace(line[2]),
			Style: strings.TrimSpace(line[3]),
		},
	}

	//parse year
	vp.Vehicle.Year, err = strconv.ParseFloat(line[0], 64)
	if err != nil {
		return vp, err
	}

	// parse part
	vp.PartID, err = strconv.Atoi(line[partIdCell])
	if err != nil {
		return vp, err
	}

	//drilling/install time
	vp.Drilling = "Yes"
	if strings.ToLower(line[40]) == "no" {
		vp.Drilling = "No"
	}
	vp.InstallTime = line[46]
	//related parts
	for _, j := range relatedPartCells {
		if line[j] != "" && strings.ToLower(line[j]) != "n/a" {
			partId, err := strconv.Atoi(line[j])
			if err != nil {
				return vp, err
			}
			vp.RelatedParts = append(vp.RelatedParts, partId)
		}
	}
	return vp, nil
}

func MatchYears(vps []VehiclePart) ([]VehiclePart, error) {
	var err error
	yearMap, err = getYearMap()
	if err != nil {
		log.Fatal(err)
	}
	for i, vp := range vps {
		var ok bool

		if vps[i].Vehicle.YearID, ok = yearMap[vp.Vehicle.Year]; !ok {
			createYear(vp.Vehicle.Year)
		}
	}
	return vps, err
}

func MatchMakes(vps []VehiclePart) ([]VehiclePart, error) {
	var err error
	makeMap, err = getMakeMap()
	if err != nil {
		log.Fatal(err)
	}
	for i, vp := range vps {
		var ok bool

		if vps[i].Vehicle.MakeID, ok = makeMap[strings.ToLower(vp.Vehicle.Make)]; !ok {
			createMake(vp.Vehicle.Make)
		}
	}
	return vps, err
}

func MatchModels(vps []VehiclePart) ([]VehiclePart, error) {
	var err error
	modelMap, err = getModelMap()
	if err != nil {
		log.Fatal(err)
	}
	for i, vp := range vps {
		var ok bool

		if vps[i].Vehicle.ModelID, ok = modelMap[strings.ToLower(vp.Vehicle.Model)]; !ok {
			createModel(vp.Vehicle.Model)
		}
	}
	return vps, err
}

func MatchStyles(vps []VehiclePart) ([]VehiclePart, error) {
	var err error
	styleMap, err = getStyleMap()
	if err != nil {
		log.Fatal(err)
	}
	for i, vp := range vps {
		var ok bool

		if vps[i].Vehicle.StyleID, ok = styleMap[strings.ToLower(vp.Vehicle.Style)]; !ok {
			createStyle(vp.Vehicle.Style)
		}
	}
	return vps, err
}

//Requires year make model style IDs
func MatchVehicles(vps []VehiclePart) ([]VehiclePart, error) {
	var err error
	vehicleMap, vehiclePartMap, err = getVehiclePartMap()
	if err != nil {
		log.Fatal(err)
	}
	for i, vp := range vps {
		var ok bool

		vehicleKey := strings.Join([]string{strconv.FormatFloat(vp.Vehicle.Year, 'f', 1, 64), strings.ToLower(vp.Vehicle.Make), strings.ToLower(vp.Vehicle.Model), strings.ToLower(vp.Vehicle.Style)}, "|")
		if vps[i].Vehicle.ID, ok = vehicleMap[strings.ToLower(vehicleKey)]; !ok {
			createVehicle(vp.Vehicle, vehicleKey)
		}
	}
	return vps, err
}

func MatchVehicleParts(vps []VehiclePart) ([]VehiclePart, error) {
	var err error
	if len(vehiclePartMap) == 0 {
		_, vehiclePartMap, err = getVehiclePartMap()
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, vp := range vps {
		var ok bool

		vehiclePartKey := strings.Join([]string{strconv.FormatFloat(vp.Vehicle.Year, 'f', 1, 64), strings.ToLower(vp.Vehicle.Make), strings.ToLower(vp.Vehicle.Model), strings.ToLower(vp.Vehicle.Style), strconv.Itoa(vp.PartID)}, "|")
		if _, ok = vehiclePartMap[strings.ToLower(vehiclePartKey)]; !ok {
			createVehiclePart(vp, vehiclePartKey)
		}
	}

	return vps, err
}

func ConfirmPartExistance(vps []VehiclePart) ([]VehiclePart, error) {
	var err error
	if len(partMap) == 0 {
		partMap, err = getPartMap()
		if err != nil {
			return vps, err
		}
	}
	for _, vp := range vps {
		var ok bool
		if _, ok = partMap[vp.PartID]; !ok {
			createPart(vp)
		}
	}
	return vps, nil
}

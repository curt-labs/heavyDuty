package importer

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/curt-labs/heavyduty/database"
)

//Vehicle (style-based)
//The 'temp' fields indicate whether these attributes were entered in temp table
type Vehicle struct {
	ID        int
	Year      float64
	YearID    int
	YearTemp  bool
	Make      string
	MakeID    int
	MakeTemp  bool
	Model     string
	ModelID   int
	ModelTemp bool
	Style     string
	StyleID   int
	StyleTemp bool
}

//VehiclePart links VehicleID to PartID
//The VehicleTemp field indicates whether this vehicle existed in the DB already
type VehiclePart struct {
	PartID      int
	Vehicle     Vehicle
	Drilling    string
	VehicleTemp bool
}

var (
	path           = flag.String("path", "", "Path to csv")
	yearMap        map[float64]int
	makeMap        map[string]int
	modelMap       map[string]int
	styleMap       map[string]int
	vehicleMap     map[string]int
	vehiclePartMap map[string]int
	makeToModelMap map[string]string

	//new
	yearMapNew    = make(map[float64]int)
	makeMapNew    = make(map[string]int)
	modelMapNew   = make(map[string]int)
	styleMapNew   = make(map[string]int)
	vehicleMapNew = make(map[string]int)
)

//Init creates maps
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
	vehicleMap, err = getVehicleMap()
	if err != nil {
		log.Fatal(err)
	}
	vehiclePartMap, err = getVehiclePartMap()
	if err != nil {
		log.Fatal(err)
	}
	makeToModelMap = make(map[string]string)

	err = database.CreateNewTables()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Maps obtained.")
}

// Get - main function for reading csv and databasing
func Get() error {
	Init()
	var counter int
	flag.Parse()
	f, err := os.Open(*path)
	if err != nil {
		return err
	}
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	rawData, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i, line := range rawData {
		if i == 0 {
			continue
		}
		//get make, model, stylesPart1 data from line above if this line is empty (cols 0, 2, 3)
		if line[0] == "" {
			line[0] = rawData[i-1][0]
		}
		if line[2] == "" && line[0] == rawData[i-1][0] {
			line[2] = rawData[i-1][2]
		}
		if line[0] == "" && line[2] == rawData[i-1][2] {
			line[0] = rawData[i-1][0]
		}

		//parse line to array of vehicleApplications
		vps, err := parseLine(line)
		if err != nil {
			return err
		}

		//range over vehicleApps from csv parsing
		for i := range vps {
			//create vehicleApps (with yearId, makeId, modelId, styleId)
			skip, err := vps[i].Build()
			if err != nil {
				return err
			}

			//insert each vehicle & vehiclePart into MySqlDB
			if skip {
				continue
			}
			log.Print(vps[i])
			err = vps[i].insert()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println(counter, " vehicleParts examined")
	return nil
}

// Get Make from line or line above
// Get Models from Models => array
// Get styles part1 from styles part1 => array
// Get Style from Style
// Create styles array; append style to each stylePart1 array item
// Get years from YearsFrom/To => array
// Get parts from remaining cells (except notes)
// Get Drilling notes; NOTE: Z or No Drill = No Drilling Required; Drill = Drilling Required
// Loop years, loop models, loop styles, loop parts -s> make vehiclePart => array
func parseLine(line []string) ([]VehiclePart, error) {
	var err error
	var vps []VehiclePart

	//make
	makeArray := strings.Split(line[0], "/")

	//model array
	modelsArray := strings.Split(line[2], ",")

	//style array
	stylesPart1Array := strings.Split(line[3], ",")
	var stylesArray []string
	for _, sp := range stylesPart1Array {
		var space string
		if sp != "" {
			space = " "
		}
		stylesArray = append(stylesArray, fmt.Sprintf("%s%s%s", sp, space, line[4]))
	}

	//years array -- ints only!!!
	yearStart, err := strconv.Atoi(line[5])
	if err != nil {
		return vps, err
	}
	yearEnd, err := strconv.Atoi(line[6])
	if err != nil {
		return vps, err
	}
	var yearsArray []float64
	for i := yearStart; i <= yearEnd; i++ {
		yearsArray = append(yearsArray, float64(i))
	}

	//parts array
	var partsArray []int
	for _, x := range []int{7, 8, 9, 10, 11} {
		if partID, err := strconv.Atoi(line[x]); err == nil {
			partsArray = append(partsArray, partID)
		}
		err = nil
	}

	//drill
	drillNote := line[12]
	if strings.Contains(drillNote, "No") {
		drillNote = "No drilling required"
	} else if strings.Contains(drillNote, "Drill") || strings.Contains(drillNote, "Z") {
		drillNote = "Drilling required"
	}

	//vehiclePart
	for _, model := range modelsArray {
		vehicleMake, err := obtainMake(makeArray, model)
		if err != nil {
			return vps, err
		}
		for _, year := range yearsArray {
			for _, style := range stylesArray {
				for _, part := range partsArray {
					vp := VehiclePart{
						Vehicle: Vehicle{
							Year:  year,
							Make:  strings.ToLower(strings.TrimSpace(vehicleMake)),
							Model: strings.ToLower(strings.TrimSpace(model)),
							Style: strings.ToLower(strings.TrimSpace(style)),
						},
						PartID:   part,
						Drilling: drillNote,
					}
					vps = append(vps, vp)
				}
			}
		}
	}
	fmt.Println("CSV line parsed.")
	return vps, err
}

//if more than one make on csv line, uses make-model map or user input to determine make
func obtainMake(makeArray []string, model string) (string, error) {
	//check make if more than one on the line (e.g. Chevrolet/GMC)
	if len(makeArray) == 1 {
		return strings.ToLower(makeArray[0]), nil
	}
	if len(makeArray) < 1 {
		return "", fmt.Errorf("No vehicle make in row")
	}

	var vehicleMake string
	var ok bool
	if vehicleMake, ok = makeToModelMap[strings.ToLower(model)]; ok {
		return vehicleMake, nil
	}

	//user choose
	var enterChoice string
	makeChoiceString := fmt.Sprintf("Is the '%s' a ", model)
	for i, m := range makeArray {
		makeChoiceString += m + " (" + strconv.Itoa(i+1) + ") "
	}
	makeChoiceString += "product? "
	fmt.Printf(makeChoiceString)
	if _, err := fmt.Scanf("%s", &enterChoice); err != nil {
		return vehicleMake, err
	}

	choice, err := strconv.Atoi(enterChoice)
	log.Print(choice)
	if err != nil {
		return vehicleMake, err
	}
	if choice < 1 || choice > len(makeArray) {
		return "", fmt.Errorf("Choice is not allowed.")
	}
	vehicleMake = strings.ToLower(makeArray[choice-1])
	makeToModelMap[strings.ToLower(model)] = vehicleMake //add to map
	return vehicleMake, nil
}

// Build makes vehicle from string values
// get year, make, model, style maps (to lower)
// look for year, make, model, style matches (to lower)
// if no match on some attribute-> insert that attribute into CurtData (capitalize first letter) & add to corresponding map
// look for vehicle match
// if no match-> insert into Vehicle Table
// look for VehiclePart match (should not be one - log existing ones)
// if no match -> insert vehicle part (w/ drilling)
func (vp *VehiclePart) Build() (bool, error) {
	var err error
	skip := true
	skip, err = vp.findYearID()
	if err != nil || skip == true {
		return true, err
	}
	skip, err = vp.findMakeID()
	if err != nil || skip == true {
		return true, err
	}
	skip, err = vp.findModelID()
	if err != nil || skip == true {
		return true, err
	}
	skip, err = vp.findStyleID()
	if err != nil || skip == true {
		return true, err
	}
	fmt.Println("Vehicle established.")
	return false, nil
}

func (vp *VehiclePart) findYearID() (bool, error) {
	var id int
	var err error
	var ok bool
	if id, ok = yearMap[vp.Vehicle.Year]; ok {
		vp.Vehicle.YearID = id
		return false, nil
	}
	if id, ok = yearMapNew[vp.Vehicle.Year]; ok {
		vp.Vehicle.YearID = id
		vp.Vehicle.YearTemp = true
		return false, nil
	}
	var enterYear string
	fmt.Printf("Enter year '%f' into database? y/n: ", vp.Vehicle.Year)
	if _, err := fmt.Scanf("%s", &enterYear); err != nil {
		return true, err
	}
	if strings.ToLower(enterYear) == "y" {
		//create year & put in map
		id, err = addYear(vp.Vehicle.Year)
		if err != nil {
			return true, err
		}
		vp.Vehicle.YearID = id
		vp.Vehicle.YearTemp = true
		return false, nil
	}
	//save missing vp to csv
	return true, vp.toErrFile(fmt.Sprintf("Opted not to enter year %f", vp.Vehicle.Year))
}
func (vp *VehiclePart) findMakeID() (bool, error) {
	var id int
	var err error
	var ok bool
	if id, ok = makeMap[vp.Vehicle.Make]; ok {
		vp.Vehicle.MakeID = id
		return false, nil
	}
	if id, ok = makeMapNew[vp.Vehicle.Make]; ok {
		vp.Vehicle.MakeID = id
		vp.Vehicle.MakeTemp = true
		return false, nil
	}

	var enterMake string
	fmt.Printf("Enter make '%s' into database? y/n: ", vp.Vehicle.Make)
	if _, err := fmt.Scanf("%s", &enterMake); err != nil {
		return true, err
	}
	if strings.ToLower(enterMake) == "y" {
		//create make (capitalize) & put in map
		id, err = addMake(vp.Vehicle.Make)
		if err != nil {
			return true, err
		}
		vp.Vehicle.MakeID = id
		vp.Vehicle.MakeTemp = true
		return false, nil
	}
	//choose and alter map
	vp.Vehicle.MakeID = chooseFromMap("make", vp.Vehicle.Make)
	if vp.Vehicle.MakeID == 0 {
		//save missing vp to csv
		return true, vp.toErrFile(fmt.Sprintf("Opted not to enter make %s", vp.Vehicle.Make))
	}
	return false, nil
}
func (vp *VehiclePart) findModelID() (bool, error) {
	var id int
	var err error
	var ok bool
	if id, ok = modelMap[vp.Vehicle.Model]; ok {
		vp.Vehicle.ModelID = id
		return false, nil
	}
	if id, ok = modelMapNew[vp.Vehicle.Model]; ok {
		vp.Vehicle.ModelID = id
		vp.Vehicle.ModelTemp = true
		return false, nil
	}
	var enterModel string
	fmt.Printf("Enter model '%s' into database? y/n: ", vp.Vehicle.Model)
	if _, err := fmt.Scanf("%s", &enterModel); err != nil {
		return true, err
	}
	if strings.ToLower(enterModel) == "y" {
		//create model(capitalize) & put in map
		id, err = addModel(vp.Vehicle.Model)
		if err != nil {
			return true, err
		}
		vp.Vehicle.ModelID = id
		vp.Vehicle.ModelTemp = true
		return false, err
	}
	//choose and alter map
	vp.Vehicle.ModelID = chooseFromMap("model", vp.Vehicle.Model)
	if vp.Vehicle.ModelID == 0 {
		//save missing vp to csv
		return true, vp.toErrFile(fmt.Sprintf("Opted not to enter model %s", vp.Vehicle.Model))
	}
	return false, nil
}
func (vp *VehiclePart) findStyleID() (bool, error) {
	var id int
	var err error
	var ok bool
	if id, ok = styleMap[vp.Vehicle.Style]; ok {
		vp.Vehicle.StyleID = id
		return false, nil
	}
	if id, ok = styleMapNew[vp.Vehicle.Style]; ok {
		vp.Vehicle.StyleID = id
		vp.Vehicle.StyleTemp = true
		return false, nil
	}

	// log.Print(styleMap)
	var enterStyle string
	fmt.Printf("Enter style '%s' into database? y/n: ", vp.Vehicle.Style)
	if _, err := fmt.Scanf("%s", &enterStyle); err != nil {
		return true, err
	}
	if strings.ToLower(enterStyle) == "y" {
		//create style(capitalize) & put in map
		id, err = addStyle(vp.Vehicle.Style)
		if err != nil {
			return true, err
		}
		vp.Vehicle.StyleID = id
		vp.Vehicle.StyleTemp = true
		return false, err
	}
	//choose and alter map
	vp.Vehicle.StyleID = chooseFromMap("style", vp.Vehicle.Style)
	if vp.Vehicle.StyleID == 0 {
		//save missing vp to csv
		return true, vp.toErrFile(fmt.Sprintf("Opted not to enter style %s", vp.Vehicle.Style))
	}
	return false, nil

}

//insert vehicle and vehicle part into MySql DB
//checks for existence of each first
func (vp *VehiclePart) insert() error {
	var ok bool
	var err error
	var vpID int //vehiclePartID - only used here
	vehicleMapKey := strconv.Itoa(vp.Vehicle.YearID) + "|" + strconv.Itoa(vp.Vehicle.MakeID) + "|" + strconv.Itoa(vp.Vehicle.ModelID) + "|" + strconv.Itoa(vp.Vehicle.StyleID)
	vehiclePartMapKey := strconv.Itoa(vp.Vehicle.ID) + "|" + strconv.Itoa(vp.PartID)

	if vp.Vehicle.ID, ok = vehicleMap[vehicleMapKey]; !ok || vp.Vehicle.ID == 0 {
		err = vp.addVehicle()
		if err != nil {
			return err
		}
		vp.VehicleTemp = true
	}

	if vpID, ok = vehiclePartMap[vehiclePartMapKey]; !ok || vpID == 0 {
		err = vp.add()
		if err != nil {
			return err
		}
	}
	return nil
}

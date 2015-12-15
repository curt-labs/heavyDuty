package importer

import (
	// "database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/curt-labs/heavyduty/database"
	"log"
	"os"
	"strconv"
	"strings"
)

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

type VehiclePart struct {
	PartID   int
	Vehicle  Vehicle
	Drilling string
}

var (
	path     = flag.String("path", "", "Path to csv")
	yearMap  map[float64]int
	makeMap  map[string]int
	modelMap map[string]int
	styleMap map[string]int
)

func init() {
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
	err = database.CreateNewTables()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Maps obtained.")
}

func Get() error {
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
	var vehicleParts []VehiclePart //all the vehicle parts that are about to happen
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
		vehicleParts = append(vehicleParts, vps...)

		//insert each vehicleApp into MySqlDB
		for i, _ := range vps {
			err = vps[i].ToDB()
			if err != nil {
				return err
			}

		}
	}
	log.Print(len(vehicleParts))
	//index (mongo-ize) all new parts

	return nil
}

func parseLine(line []string) ([]VehiclePart, error) {
	// Get Make from line or line above
	// Get Models from Models => array
	// Get styles part1 from styles part1 => array
	// Get Style from Style
	// Create styles array; append style to each stylePart1 array item
	// Get years from YearsFrom/To => array
	// Get parts from remaining cells (except notes)
	// Get Drilling notes; NOTE: Z or No Drill = No Drilling Required; Drill = Drilling Required
	// Loop years, loop models, loop styles, loop parts -s> make vehiclePart => array
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
		if partId, err := strconv.Atoi(line[x]); err == nil {
			partsArray = append(partsArray, partId)
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
		//check make if more than one on the line (e.g. Chevrolet/GMC)
		//TODO - this is totally trouble!!! - split makes eariler
		var vehicleMake string
		if len(makeArray) > 1 {
			var enterChoice string
			makeChoiceString := fmt.Sprintf("Is the %s a ", model)
			for i, m := range makeArray {
				makeChoiceString += m + " (" + strconv.Itoa(i+1) + ") "
			}
			makeChoiceString += "product? "
			fmt.Printf(makeChoiceString)
			if _, err := fmt.Scanf("%s", &enterChoice); err != nil {
				return vps, err
			}

			choice, err := strconv.Atoi(enterChoice)
			if err != nil {
				//TODO - reenter or quit
				return vps, err
			}
			vehicleMake = strings.ToLower(makeArray[choice-1]) //Capitalize
		} else {
			vehicleMake = strings.ToLower(makeArray[0])
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
	fmt.Println("CSV parsed.")
	return vps, err
}

func (vp *VehiclePart) ToDB() error {
	var err error
	// get year, make, model, style maps (to lower)
	// look for year, make, model, style matches (to lower)
	// if no match on some attribute-> insert that attribute into CurtData (capitalize first letter) & add to corresponding map
	// look for vehicle match
	// if no match-> insert into Vehicle Table
	// look for VehiclePart match (should not be one - log existing ones)
	// if no match -> insert vehicle part (w/ drilling)

	var ok bool
	log.Print("HERE")
	if vp.Vehicle.YearID, ok = yearMap[vp.Vehicle.Year]; !ok {
		var enterYear string
		fmt.Printf("Enter year %d into database? y/n: ", vp.Vehicle.Year)
		if _, err := fmt.Scanf("%s", &enterYear); err != nil {
			return err
		}
		if strings.ToLower(enterYear) == "y" {
			//create year & put in map

		} else {
			//save missing vp to csv
		}
	}
	if vp.Vehicle.MakeID, ok = makeMap[vp.Vehicle.Make]; !ok {
		var enterMake string
		fmt.Printf("Enter make %s into database? y/n: ", vp.Vehicle.Make)
		if _, err := fmt.Scanf("%s", &enterMake); err != nil {
			return err
		}
		if strings.ToLower(enterMake) == "y" {
			//create make (capitalize) & put in map
		} else {
			//save missing vp to csv
		}
	}
	if vp.Vehicle.ModelID, ok = modelMap[vp.Vehicle.Model]; !ok {
		var enterModel string
		fmt.Printf("Enter model %s into database? y/n: ", vp.Vehicle.Model)
		if _, err := fmt.Scanf("%s", &enterModel); err != nil {
			return err
		}
		if strings.ToLower(enterModel) == "y" {
			//create model(capitalize) & put in map
		} else {
			//save missing vp to csv
		}
	}
	if vp.Vehicle.StyleID, ok = styleMap[vp.Vehicle.Style]; !ok {
		var enterStyle string
		fmt.Printf("Enter style %s into database? y/n: ", vp.Vehicle.Style)
		if _, err := fmt.Scanf("%s", &enterStyle); err != nil {
			return err
		}
		if strings.ToLower(enterStyle) == "y" {
			//create style(capitalize) & put in map
		} else {
			//save missing vp to csv
		}
	}

	fmt.Println("VehicleApplication DB-ed.")
	return err
}

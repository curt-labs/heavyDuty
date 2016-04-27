package importer

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	newYear        []string
	newMake        []string
	newModel       []string
	newStyle       []string
	newVehicle     []Vehicle
	newVehiclePart []VehiclePart
	newPart        []int

	newYearStmt        string
	newMakeStmt        string
	newModelStmt       string
	newStyleStmt       string
	newVehicleStmt     string
	newVehiclePartStmt string
	relatedPartStmt    string
	newPartStmt        string
)

func createYear(year float64) {
	yearStr := strconv.FormatFloat(year, 'f', 1, 64)
	newYear = append(newYear, yearStr)
	id := rand.Int()
	yearMap[year] = int(id)
}

func createMake(ma string) {
	newMake = append(newMake, ma)
	id := rand.Int()
	makeMap[strings.ToLower(ma)] = int(id)
}

func createModel(mo string) {
	newModel = append(newModel, mo)
	id := rand.Int()
	modelMap[strings.ToLower(mo)] = int(id)
}

func createStyle(st string) {
	newStyle = append(newStyle, st)
	id := rand.Int()
	styleMap[strings.ToLower(st)] = int(id)
}

func createVehicle(v Vehicle, key string) {
	newVehicle = append(newVehicle, v)
	id := rand.Int()
	vehicleMap[key] = int(id)
}

func createVehiclePart(vp VehiclePart, key string) {
	newVehiclePart = append(newVehiclePart, vp)
	id := rand.Int()
	vehiclePartMap[key] = int(id)
}

func createPart(vp VehiclePart) {
	newPart = append(newPart, vp.PartID)
	id := time.Now().UnixNano()
	partMap[vp.PartID] = int(id)
}

func CreateRelatedParts(vps []VehiclePart) error {
	var str string
	relatePartsMap, err := getRelatedPartsMap()
	if err != nil {
		return err
	}
	for _, vp := range vps {
		for _, r := range vp.RelatedParts {
			key := strconv.Itoa(vp.PartID) + "|" + strconv.Itoa(r)
			if _, ok := relatePartsMap[key]; ok {
				continue
			}
			str += "(" + strconv.Itoa(vp.PartID) + "," + strconv.Itoa(r) + ", 0),"
			relatePartsMap[key] = rand.Int()
		}
	}
	str = strings.TrimRight(str, ",")
	relatedPartStmt = fmt.Sprintf("insert into RelatedPart (partID, relatedID, rTypeID) values %s;", str)
	f, err := os.Create("relatedparts.txt")
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(relatedPartStmt))
	return err
}

func CreateStmts() error {
	var year string
	for i, y := range newYear {
		if i > 0 {
			year += ","
		}
		year += "(" + y + ")"
	}
	newYearStmt = fmt.Sprintf("insert into Year (year) values %s;", year)

	var ma string
	for i, y := range newMake {
		if i > 0 {
			ma += ","
		}
		ma += "(\"" + y + "\")"
	}
	newMakeStmt = fmt.Sprintf("insert into Make (make) values %s;", ma)

	var mo string
	for i, y := range newModel {
		if i > 0 {
			mo += ","
		}
		mo += "(\"" + y + "\")"
	}
	newModelStmt = fmt.Sprintf("insert into Model (model) values %s;", mo)

	var st string
	for i, y := range newStyle {
		if i > 0 {
			st += ","
		}
		st += "(\"" + y + "\",0)"
	}
	newStyleStmt = fmt.Sprintf("insert into Style (style, aaiaID) values %s;", st)

	var vstr string
	for i, v := range newVehicle {
		if v.YearID == 0 || v.MakeID == 0 || v.ModelID == 0 || v.StyleID == 0 {
			vp := VehiclePart{
				Vehicle: v,
			}
			vp.toErrFile(fmt.Sprintf("error on new vehicle"))
			continue
		}
		if i > 0 {
			vstr += ","
		}
		vstr += `(` + strconv.Itoa(v.YearID) + `,` + strconv.Itoa(v.MakeID) + `,` + strconv.Itoa(v.ModelID) + `,` + strconv.Itoa(v.StyleID) + `,now())`

	}
	newVehicleStmt = fmt.Sprintf("insert into Vehicle (yearID, makeID, modelID, styleID, dateAdded) values %s;", vstr)

	var pstr string
	for i, v := range newPart {
		if v < 1 {
			continue
		}
		if i > 0 {
			pstr += ","
		}
		pstr += `(` + strconv.Itoa(v) + `, 800, now(), now(), 0, 0, 1)`
	}
	newPartStmt = fmt.Sprintf("insert into Part (partID, status, dateModified, dateAdded, classID, featured, brandID) values %s;", pstr)

	var vpstr string
	for i, v := range newVehiclePart {
		if v.InstallTime == "" {
			v.InstallTime = "NULL"
		}
		if v.Vehicle.ID == 0 || v.PartID == 0 {
			v.toErrFile(fmt.Sprintf("error on new vehicle part"))
			continue
		}
		if i > 0 {
			vpstr += ","
		}
		vpstr += `(` + strconv.Itoa(v.Vehicle.ID) + `,` + strconv.Itoa(v.PartID) + `,"` + v.Drilling + `",` + v.InstallTime + `)`

	}
	newVehiclePartStmt = fmt.Sprintf("insert into VehiclePart (vehicleID, partID, drilling, installTime) values %s;", vpstr)

	//OLD INNER INNER SELECT QUERIES
	// var vStr string
	// for i, v := range newVehicle {
	// 	if i > 0 {
	// 		vStr += ","
	// 	}
	// 	vStr += `((select y.yearID from Year y where y.year = "` + strconv.FormatFloat(v.Year, 'f', 1, 64) + `"),
	// 			(select ma.makeID from Make ma where ma.make = "` + v.Make + `"),
	// 			(select mo.modelID from Model mo where mo.model = "` + v.Model + `"),
	// 			(select st.styleID from Style st where st.style = "` + v.Style + `"))`
	// }
	// newVehicleStmt = fmt.Sprintf("insert into Vehicle (yearID, makeID, modelID, styleID) values %s;", vStr)

	// var vpStr string
	// for i, v := range newVehiclePart {
	// 	if i > 0 {
	// 		vpStr += ","
	// 	}
	// 	vpStr += `((select v.vehicleID from
	// 		Vehicle v
	// 		where v.yearID = (select y.yearID from Year y where y.year = "` + strconv.FormatFloat(v.Vehicle.Year, 'f', 1, 64) + `")
	// 		and v.makeID = (select ma.makeID from Make ma where ma.make = "` + v.Vehicle.Make + `")
	// 		and v.modelID = (select mo.modelID from Model mo where mo.model = "` + v.Vehicle.Model + `")
	// 		and v.styleID = (select st.styleID from Style st where st.style = "` + v.Vehicle.Style + `")
	// 		), ` + strconv.Itoa(v.PartID) + `,"` + v.Drilling + `",` + v.InstallTime + ` )`
	// }
	// newVehiclePartStmt = fmt.Sprintf("insert into VehiclePart (vehicleID, partID, drilling, installTime) values %s;", vpStr)

	f, err := os.Create("sql.txt")
	if err != nil {
		return err
	}
	var offset int64
	if len(newYear) > 0 {
		n, err := f.WriteAt([]byte(newYearStmt), offset)
		if err != nil {
			return err
		}
		offset += int64(n)
	}
	if len(newMake) > 0 {
		n, err := f.WriteAt([]byte(newMakeStmt), offset)
		if err != nil {
			return err
		}
		offset += int64(n)
	}
	if len(newModel) > 0 {
		n, err := f.WriteAt([]byte(newModelStmt), offset)
		if err != nil {
			return err
		}
		offset += int64(n)
	}
	if len(newStyle) > 0 {
		n, err := f.WriteAt([]byte(newStyleStmt), offset)
		if err != nil {
			return err
		}
		offset += int64(n)
	}
	if len(newVehicle) > 0 {
		n, err := f.WriteAt([]byte(newVehicleStmt), offset)
		if err != nil {
			return err
		}
		offset += int64(n)
	}
	if len(newPart) > 0 {
		n, err := f.WriteAt([]byte(newPartStmt), offset)
		if err != nil {
			return err
		}
		offset += int64(n)
	}
	if len(newVehiclePart) > 0 {
		n, err := f.WriteAt([]byte(newVehiclePartStmt), offset)
		if err != nil {
			return err
		}
		offset += int64(n)
	}

	return nil
}

package importer

import (
	"database/sql"
	"strconv"

	"github.com/curt-labs/heavyduty/database"
)

func addYear(year float64) (int, error) {
	var ID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ID, err
	}
	defer db.Close()

	res, err := db.Exec("insert into "+database.YEAR_TABLE+" (year) values (?)", year)
	if err != nil {
		return ID, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return ID, err
	}
	ID = int(id)

	//add to map
	yearMap[year] = ID
	return ID, err
}

func addMake(makeName string) (int, error) {
	var ID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ID, err
	}
	defer db.Close()

	res, err := db.Exec("insert into "+database.MAKE_TABLE+" (make) values (?)", capitalize(makeName))
	if err != nil {
		return ID, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return ID, err
	}
	ID = int(id)

	//add to map
	makeMap[makeName] = ID
	return ID, err
}

func addModel(model string) (int, error) {
	var ID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ID, err
	}
	defer db.Close()

	res, err := db.Exec("insert into "+database.MODEL_TABLE+" (model) values (?)", capitalize(model))
	if err != nil {
		return ID, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return ID, err
	}
	ID = int(id)

	//add to map
	modelMap[model] = ID
	return ID, err
}

func addStyle(style string) (int, error) {
	var ID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ID, err
	}
	defer db.Close()

	res, err := db.Exec("insert into "+database.STYLE_TABLE+" (style, aaiaID) values (?, 0)", capitalize(style))
	if err != nil {
		return ID, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return ID, err
	}
	ID = int(id)

	//add to map
	styleMap[style] = ID
	return ID, err
}

func (vp *VehiclePart) addVehicle() error {
	var ID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(
		"insert into "+database.VEHICLE_TABLE+" (yearID, makeID, modelID, styleID, dateAdded, yearTemp, makeTemp, modelTemp, styleTemp) values (?, ?, ?, ?, NOW(),?,?,?,?)",
		vp.Vehicle.YearID,
		vp.Vehicle.MakeID,
		vp.Vehicle.ModelID,
		vp.Vehicle.StyleID,
		vp.Vehicle.YearTemp,
		vp.Vehicle.MakeTemp,
		vp.Vehicle.ModelTemp,
		vp.Vehicle.StyleTemp,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	vp.Vehicle.ID = int(id)
	v := strconv.Itoa(vp.Vehicle.YearID) + "|" + strconv.Itoa(vp.Vehicle.MakeID) + "|" + strconv.Itoa(vp.Vehicle.ModelID) + "|" + strconv.Itoa(vp.Vehicle.StyleID)
	//add to map
	vehicleMap[v] = ID
	return err
}

func (vp *VehiclePart) add() error {
	var ID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec("insert into "+database.VEHICLE_PART_TABLE+" (vehicleID, partID, drilling, vehicleTemp) values (?,?,?,?)", vp.Vehicle.ID, vp.PartID, vp.Drilling, vp.VehicleTemp)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	vp.Vehicle.ID = int(id)

	//add to map
	v := strconv.Itoa(vp.Vehicle.ID) + "|" + strconv.Itoa(vp.PartID)
	vehiclePartMap[v] = ID
	return err
}

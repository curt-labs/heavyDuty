package importer

import (
	"database/sql"
	"strconv"

	"github.com/curt-labs/heavierduty/database"
)

func addYear(year float64) (int, error) {
	var ID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ID, err
	}
	defer db.Close()

	res, err := db.Exec("insert into "+database.YearTable+" (year) values (?)", year)
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

	res, err := db.Exec("insert into "+database.MakeTable+" (make) values (?)", capitalize(makeName))
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

	res, err := db.Exec("insert into "+database.ModelTable+" (model) values (?)", capitalize(model))
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

	res, err := db.Exec("insert into "+database.StyleTable+" (style, aaiaID) values (?, 0)", capitalize(style))
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
		"insert into "+database.VehicleTable+" (yearID, makeID, modelID, styleID, dateAdded) values (?, ?, ?, ?, NOW())",
		vp.Vehicle.YearID,
		vp.Vehicle.MakeID,
		vp.Vehicle.ModelID,
		vp.Vehicle.StyleID,
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

func (vp *VehiclePart) add() (int, error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec("insert into "+database.VehiclePartTable+" (vehicleID, partID, drilling) values (?,?,?)", vp.Vehicle.ID, vp.PartID, vp.Drilling)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	//add to map
	v := strconv.Itoa(vp.Vehicle.ID) + "|" + strconv.Itoa(vp.PartID)
	vehiclePartMap[v] = int(id)
	return int(id), err
}

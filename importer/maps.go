package importer

import (
	"database/sql"
	"fmt"

	"github.com/curt-labs/heavyduty/database"
)

func getYearMap() (map[float64]int, error) {
	zMap := make(map[float64]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select yearID, year from %s", database.YearTable))
	if err != nil {
		return zMap, err
	}
	var i *int
	var y *float64
	for rows.Next() {
		err = rows.Scan(&i, &y)
		if err != nil {
			return zMap, err
		}
		zMap[*y] = *i
	}
	return zMap, nil
}

func getMakeMap() (map[string]int, error) {
	zMap := make(map[string]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	rows, err := db.Query(fmt.Sprintf("select makeID, lower(make) from %s", database.MakeTable))
	if err != nil {
		return zMap, err
	}
	var i *int
	var m *string
	for rows.Next() {
		err = rows.Scan(&i, &m)
		if err != nil {
			return zMap, err
		}
		zMap[*m] = *i
	}
	return zMap, nil
}

func getModelMap() (map[string]int, error) {
	zMap := make(map[string]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("select modelID, lower(model) from %s", database.ModelTable))
	if err != nil {
		return zMap, err
	}
	var i *int
	var m *string
	for rows.Next() {
		err = rows.Scan(&i, &m)
		if err != nil {
			return zMap, err
		}
		zMap[*m] = *i
	}
	return zMap, nil
}

func getStyleMap() (map[string]int, error) {
	zMap := make(map[string]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select styleID, lower(style) from %s", database.StyleTable))
	if err != nil {
		return zMap, err
	}
	var i *int
	var m *string
	for rows.Next() {
		err = rows.Scan(&i, &m)
		if err != nil {
			return zMap, err
		}
		zMap[*m] = *i
	}
	return zMap, nil
}

func getVehicleMap() (map[string]int, error) {
	zMap := make(map[string]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select vehicleID, concat_ws('|', yearID, makeID, modelID, styleID) as v from %s", database.VehicleTable))
	if err != nil {
		return zMap, err
	}
	var i *int
	var m *string
	for rows.Next() {
		err = rows.Scan(&i, &m)
		if err != nil {
			return zMap, err
		}
		zMap[*m] = *i
	}
	return zMap, nil
}

func getVehiclePartMap() (map[string]int, error) {
	zMap := make(map[string]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select vPartID, concat_ws('|', vehicleID, partID) as v from %s", database.VehiclePartTable))
	if err != nil {
		return zMap, err
	}
	var i *int
	var m *string
	for rows.Next() {
		err = rows.Scan(&i, &m)
		if err != nil {
			return zMap, err
		}
		zMap[*m] = *i
	}
	return zMap, nil
}

package importer

import (
	"database/sql"
	"github.com/curt-labs/heavyduty/database"
)

func getYearMap() (map[float64]int, error) {
	zMap := make(map[float64]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()

	rows, err := db.Query("select yearID, year from Year")
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
	rows, err := db.Query("select makeID, lower(make) from Make")
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
	rows, err := db.Query("select modelID, lower(model) from Model")
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

	rows, err := db.Query("select styleID, lower(style) from Style")
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

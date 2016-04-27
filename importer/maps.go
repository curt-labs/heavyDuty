package importer

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/curt-labs/heavierduty/database"
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
		zMap[strings.TrimSpace(*m)] = *i
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
		zMap[strings.TrimSpace(*m)] = *i
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
		zMap[strings.TrimSpace(*m)] = *i
	}
	return zMap, nil
}

func getVehiclePartMap() (map[string]int, map[string]int, error) {
	vMap := make(map[string]int)
	vpMap := make(map[string]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vMap, vpMap, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf(`select v.vehicleID, y.year, ma.make, mo.model, s.style, vp.partID from
		Vehicle v
		join Year y on y.yearID = v.yearID
		join Make ma on ma.makeID = v.makeID
		join Model mo on mo.modelID = v.modelID
		join Style s on s.styleID = v.styleID
		left join VehiclePart vp on vp.vehicleID = v.vehicleID
		where ma.make in ('Chevrolet','Dodge','Ford','GMC','Nissan','Ram','Toyota')
		and y.year > 1979`))
	if err != nil {
		return vMap, vpMap, err
	}
	var i, p *int
	var y *float64
	var ma, mo, st *string
	for rows.Next() {
		err = rows.Scan(&i, &y, &ma, &mo, &st, &p)
		if err != nil {
			return vMap, vpMap, err
		}

		vehicleKey := strings.Join([]string{strconv.FormatFloat(*y, 'f', 1, 64), strings.TrimSpace(strings.ToLower(*ma)), strings.TrimSpace(strings.ToLower(*mo)), strings.TrimSpace(strings.ToLower(*st))}, "|")
		vMap[vehicleKey] = *i

		if p != nil {
			vehiclePartKey := strings.Join([]string{strconv.FormatFloat(*y, 'f', 1, 64), strings.TrimSpace(strings.ToLower(*ma)), strings.TrimSpace(strings.ToLower(*mo)), strings.TrimSpace(strings.ToLower(*st)), strconv.Itoa(*p)}, "|")
			vpMap[vehiclePartKey] = *i
		}
	}
	return vMap, vpMap, nil
}

func getRelatedPartsMap() (map[string]int, error) {
	zMap := make(map[string]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select relPartID, partID, relatedID from  %s", database.RelatedPartsTable))
	if err != nil {
		return zMap, err
	}
	var i, p, r *int

	for rows.Next() {
		err = rows.Scan(&i, &p, &r)
		if err != nil {
			return zMap, err
		}
		key := strconv.Itoa(*p) + "|" + strconv.Itoa(*r)
		zMap[key] = *i
	}
	return zMap, nil
}

func getPartMap() (map[int]int, error) {
	zMap := make(map[int]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return zMap, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select partID from %s", database.PartTable))
	if err != nil {
		return zMap, err
	}
	var i *int
	for rows.Next() {
		err = rows.Scan(&i)
		if err != nil {
			return zMap, err
		}
		zMap[*i] = *i
	}
	return zMap, nil
}

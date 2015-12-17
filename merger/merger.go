package merger

import (
	"database/sql"
	"fmt"
	// "log"
	"os"
	"strings"

	"github.com/curt-labs/heavyduty/database"
	"github.com/curt-labs/heavyduty/importer"
	_ "github.com/go-sql-driver/mysql" //you need this
)

var (
	getNewVehiclesStmt = `select vpt.partID, vpt.vehicleID, vpt.vehicleTemp, v.yearID, vt.yearID, vt.yearTemp, 
	yt.year, v.makeID, vt.makeID, vt.makeTemp, mat.make, v.modelID, vt.modelID, vt.modelTemp, mt.model, 
	v.styleID, vt.styleID, vt.styleTemp, st.style, vpt.drilling
    from VehiclePartTemp vpt
    left join VehicleTemp as vt on vt.vehicleID = vpt.vehicleID and vpt.vehicleTemp = 1
    left join Vehicle as v on v.vehicleID = vpt.vehicleID and vpt.vehicleTemp = 0
    left join YearTEMP yt on yt.yearID = vt.yearID
    left join MakeTEMP mat on mat.makeID = vt.makeID
    left join ModelTEMP mt on mt.modelID = vt.modelID
    left join StyleTEMP st on st.styleID = vt.styleID`
	insertedYears  = make(map[int]int) //tempID to new ID
	insertedMakes  = make(map[int]int)
	insertedModels = make(map[int]int)
	insertedStyles = make(map[int]int)
)

//Merge puts TEMP table data into CurtData DB
func Merge() error {
	var input string
	fmt.Printf("Merge TEMP tables into CurtData at DB_HOST: '%s' (y/n)?\n", os.Getenv("DATABASE_HOST"))
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		return err
	}
	if strings.ToLower(input) != "y" {
		fmt.Println("Exiting program.")
		return nil
	}
	vps, err := getNewVehicles()
	if err != nil {
		return err
	}
	err = portYears()
	if err != nil {
		return err
	}

	err = portMakes()
	if err != nil {
		return err
	}

	err = portModels()
	if err != nil {
		return err
	}

	err = portStyles()
	if err != nil {
		return err
	}

	//insert new vehicles and vehicleParts
	err = portVehicles(vps)
	if err != nil {
		return err
	}

	return nil
}

func portVehicles(vps []importer.VehiclePart) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	vehicleInsertStmt, err := tx.Prepare(fmt.Sprintf("insert into %s (yearID, makeID, modelID, styleID, dateAdded) values (?,?,?,?,NOW()) ", database.VehicleTable))
	if err != nil {
		return err
	}

	vehiclePartInsertStmt, err := tx.Prepare(fmt.Sprintf("insert into %s (vehicleID, partID, drilling) values(?,?,?)", database.VehiclePartTable))
	if err != nil {
		return err
	}

	for _, vp := range vps {
		// fmt.Println("called")
		//convert yearIDs if they were temp inserted
		if vp.Vehicle.YearTemp {
			vp.Vehicle.YearID = insertedYears[vp.Vehicle.YearID]
		}
		if vp.Vehicle.MakeTemp {
			vp.Vehicle.MakeID = insertedMakes[vp.Vehicle.MakeID]
		}
		if vp.Vehicle.ModelTemp {
			vp.Vehicle.ModelID = insertedModels[vp.Vehicle.ModelID]
		}
		if vp.Vehicle.StyleTemp {
			vp.Vehicle.StyleID = insertedStyles[vp.Vehicle.StyleID]
		}

		res, err := tx.Stmt(vehicleInsertStmt).Exec(vp.Vehicle.YearID, vp.Vehicle.MakeID, vp.Vehicle.ModelID, vp.Vehicle.StyleID)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		vp.Vehicle.ID = int(id)

		_, err = tx.Stmt(vehiclePartInsertStmt).Exec(vp.Vehicle.ID, vp.PartID, vp.Drilling)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// func portVehicles(vps []importer.VehiclePart) error {
// 	var valueStrings []string
// 	valueArgs := make([]interface{}, 0, len(vps))
// 	for _, vp := range vps {
// 		//convert yearIDs if they were temp inserted
// 		if vp.Vehicle.YearTemp {
// 			vp.Vehicle.YearID = insertedYears[vp.Vehicle.YearID]
// 		}
// 		if vp.Vehicle.MakeTemp {
// 			vp.Vehicle.MakeID = insertedMakes[vp.Vehicle.MakeID]
// 		}
// 		if vp.Vehicle.ModelTemp {
// 			vp.Vehicle.ModelID = insertedModels[vp.Vehicle.ModelID]
// 		}
// 		if vp.Vehicle.StyleTemp {
// 			vp.Vehicle.StyleID = insertedStyles[vp.Vehicle.StyleID]
// 		}
// 		valueStrings = append(valueStrings, "(?,?,?,?,NOW())")
// 		valueArgs = append(valueArgs, vp.Vehicle.YearID, vp.Vehicle.MakeID, vp.Vehicle.ModelID, vp.Vehicle.StyleID)
// 	}

// 	vehicleInsertStmt := fmt.Sprintf("insert into %s (yearID, makeID, modelID, styleID, dateAdded) values %s ", database.VehicleTable, strings.Join(valueStrings, ","))

// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()
// 	res, err := db.Exec(vehicleInsertStmt, valueArgs...)
// 	fmt.Println(res.LastInsertId())
// 	return err
// }

func portYears() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("select yearID, year from %s", database.YearTableTemp))
	if err != nil {
		return err
	}
	var year *float64
	var tempID *int
	for rows.Next() {
		err = rows.Scan(&tempID, &year)
		if err != nil {
			return err
		}

		res, err := db.Exec(fmt.Sprintf("insert into %s(year) values(?)", database.YearTable), *year)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		insertedYears[*tempID] = int(id) //map tempID to new id
	}
	return err
}

func portMakes() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("select makeID, make from %s", database.MakeTableTemp))
	if err != nil {
		return err
	}
	var ma *string
	var tempID *int
	for rows.Next() {
		err = rows.Scan(&tempID, &ma)
		if err != nil {
			return err
		}

		res, err := db.Exec(fmt.Sprintf("insert into %s(make) values(?)", database.MakeTable), *ma)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		insertedMakes[*tempID] = int(id) //map tempID to new id
	}
	return err
}

func portModels() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("select modelID, model from %s", database.ModelTableTemp))
	if err != nil {
		return err
	}
	var ma *string
	var tempID *int
	for rows.Next() {
		err = rows.Scan(&tempID, &ma)
		if err != nil {
			return err
		}

		res, err := db.Exec(fmt.Sprintf("insert into %s(model) values(?)", database.ModelTable), *ma)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		insertedModels[*tempID] = int(id) //map tempID to new id
	}
	return err
}

func portStyles() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("select styleID, style from %s", database.StyleTableTemp))
	if err != nil {
		return err
	}
	var ma *string
	var tempID *int
	for rows.Next() {
		err = rows.Scan(&tempID, &ma)
		if err != nil {
			return err
		}

		res, err := db.Exec(fmt.Sprintf("insert into %s(style, aaiaID) values(?, 0)", database.StyleTable), *ma)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		insertedStyles[*tempID] = int(id) //map tempID to new id
	}
	return err
}

func getNewVehicles() ([]importer.VehiclePart, error) {
	var vps []importer.VehiclePart
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vps, err
	}
	defer db.Close()
	res, err := db.Query(getNewVehiclesStmt)
	if err != nil {
		return vps, err
	}
	var vp importer.VehiclePart
	var yearID, makeID, modelID, styleID, yearTID, makeTID, modelTID, styleTID *int
	var year *float64
	var ma, mo, style, drilling *string
	var yearBool, makeBool, modelBool, styleBool *bool
	for res.Next() {
		err = res.Scan(
			&vp.PartID,
			&vp.Vehicle.ID,
			&vp.VehicleTemp,
			&yearID,
			&yearTID,
			&yearBool,
			&year,
			&makeID,
			&makeTID,
			&makeBool,
			&ma,
			&modelID,
			&modelTID,
			&modelBool,
			&mo,
			&styleID,
			&styleTID,
			&styleBool,
			&style,
			&drilling,
		)
		if err != nil {
			return vps, err
		}
		//Nulls
		if year != nil {
			vp.Vehicle.Year = *year
		}
		if ma != nil {
			vp.Vehicle.Make = *ma
		}
		if mo != nil {
			vp.Vehicle.Model = *mo
		}
		if style != nil {
			vp.Vehicle.Style = *style
		}
		if drilling != nil {
			vp.Drilling = *drilling
		}

		//Original or New ID's - one will be nil
		if yearTID != nil {
			vp.Vehicle.YearID = *yearTID
		} else {
			vp.Vehicle.YearID = *yearID
		}
		if makeTID != nil {
			vp.Vehicle.MakeID = *makeTID
		} else {
			vp.Vehicle.MakeID = *makeID
		}
		if modelTID != nil {
			vp.Vehicle.ModelID = *modelTID
		} else {
			vp.Vehicle.ModelID = *modelID
		}
		if styleTID != nil {
			vp.Vehicle.StyleID = *styleTID
		} else {
			vp.Vehicle.StyleID = *styleID
		}

		//new or existing value flags from DB
		if yearBool != nil {
			vp.Vehicle.YearTemp = *yearBool
		}
		if makeBool != nil {
			vp.Vehicle.MakeTemp = *makeBool
		}
		if modelBool != nil {
			vp.Vehicle.ModelTemp = *modelBool
		}
		if styleBool != nil {
			vp.Vehicle.StyleTemp = *styleBool
		}

		vps = append(vps, vp)
	}
	return vps, err
}

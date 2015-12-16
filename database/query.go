package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	//temp table names
	STYLE_TABLE        = "StyleTEMP"
	MAKE_TABLE         = "MakeTEMP"
	MODEL_TABLE        = "ModelTEMP"
	YEAR_TABLE         = "YearTEMP"
	VEHICLE_TABLE      = "VehicleTEMP"
	VEHICLE_PART_TABLE = "VehiclePartTEMP"
)

//create and/or truncate temp tables
func CreateNewTables() error {
	styles := `CREATE TABLE IF NOT EXISTS ` + STYLE_TABLE + ` (
		  styleID int(11) NOT NULL AUTO_INCREMENT,
		  style varchar(255) NOT NULL,
		  aaiaID int(11) NOT NULL,
		  PRIMARY KEY (styleID)
		) ENGINE=InnoDB AUTO_INCREMENT=699 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	makes := `CREATE TABLE IF NOT EXISTS ` + MAKE_TABLE + ` (
		  makeID int(11) NOT NULL AUTO_INCREMENT,
		  make varchar(255) NOT NULL,
		  PRIMARY KEY (makeID)
		) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	models := `CREATE TABLE IF NOT EXISTS ` + MODEL_TABLE + ` (
		  modelID int(11) NOT NULL AUTO_INCREMENT,
		  model varchar(255) NOT NULL,
		  PRIMARY KEY (modelID)
		) ENGINE=InnoDB AUTO_INCREMENT=783 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	years := `CREATE TABLE IF NOT EXISTS ` + YEAR_TABLE + ` (
		  yearID int(11) NOT NULL AUTO_INCREMENT,
		  year double NOT NULL,
		  PRIMARY KEY (yearID)
		) ENGINE=InnoDB AUTO_INCREMENT=289 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	vehicles := `CREATE TABLE IF NOT EXISTS ` + VEHICLE_TABLE + ` (
		  vehicleID int(11) NOT NULL AUTO_INCREMENT,
		  yearID int(11) NOT NULL,
		  makeID int(11) NOT NULL,
		  modelID int(11) NOT NULL,
		  styleID int(11) NOT NULL,
		  dateAdded timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		  PRIMARY KEY (vehicleID)
		) ENGINE=InnoDB AUTO_INCREMENT=250265 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	vehicleParts := `CREATE TABLE IF NOT EXISTS ` + VEHICLE_PART_TABLE + ` (
		  vPartID int(11) NOT NULL AUTO_INCREMENT,
		  vehicleID int(11) NOT NULL,
		  partID int(11) NOT NULL,
		  drilling varchar(100) DEFAULT NULL,
		  exposed varchar(100) DEFAULT NULL,
		  installTime int(11) DEFAULT NULL,
		  PRIMARY KEY (vPartID)
		) ENGINE=InnoDB AUTO_INCREMENT=51821 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	t1 := `TRUNCATE TABLE ` + STYLE_TABLE + `;`
	t2 := `TRUNCATE TABLE ` + MAKE_TABLE + `;`
	t3 := `TRUNCATE TABLE ` + MODEL_TABLE + `;`
	t4 := `TRUNCATE TABLE ` + YEAR_TABLE + `;`
	t5 := `TRUNCATE TABLE ` + VEHICLE_TABLE + `;`
	t6 := `TRUNCATE TABLE ` + VEHICLE_PART_TABLE + `;`

	db, err := sql.Open("mysql", ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	//create temp tables
	_, err = db.Exec(styles)
	if err != nil {
		return err
	}
	_, err = db.Exec(makes)
	if err != nil {
		return err
	}
	_, err = db.Exec(models)
	if err != nil {
		return err
	}
	_, err = db.Exec(years)
	if err != nil {
		return err
	}
	_, err = db.Exec(vehicles)
	if err != nil {
		return err
	}
	_, err = db.Exec(vehicleParts)
	if err != nil {
		return err
	}

	//truncate tables
	_, err = db.Exec(t1)
	if err != nil {
		return err
	}
	_, err = db.Exec(t2)
	if err != nil {
		return err
	}
	_, err = db.Exec(t3)
	if err != nil {
		return err
	}
	_, err = db.Exec(t4)
	if err != nil {
		return err
	}
	_, err = db.Exec(t5)
	if err != nil {
		return err
	}
	_, err = db.Exec(t6)
	if err != nil {
		return err
	}

	return nil
}

package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func CreateNewTables() error {
	styles := `CREATE TABLE IF NOT EXISTS Style (
		  styleID int(11) NOT NULL AUTO_INCREMENT,
		  style varchar(255) NOT NULL,
		  aaiaID int(11) NOT NULL,
		  PRIMARY KEY (styleID)
		) ENGINE=InnoDB AUTO_INCREMENT=699 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	makes := `CREATE TABLE IF NOT EXISTS MakeNEW (
		  makeID int(11) NOT NULL AUTO_INCREMENT,
		  make varchar(255) NOT NULL,
		  PRIMARY KEY (makeID)
		) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	models := `CREATE TABLE IF NOT EXISTS ModelNEW (
		  modelID int(11) NOT NULL AUTO_INCREMENT,
		  model varchar(255) NOT NULL,
		  PRIMARY KEY (modelID)
		) ENGINE=InnoDB AUTO_INCREMENT=783 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	years := `CREATE TABLE IF NOT EXISTS YearNEW (
		  yearID int(11) NOT NULL AUTO_INCREMENT,
		  year double NOT NULL,
		  PRIMARY KEY (yearID)
		) ENGINE=InnoDB AUTO_INCREMENT=289 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	vehicles := `CREATE TABLE IF NOT EXISTS VehicleNEW (
		  vehicleID int(11) NOT NULL AUTO_INCREMENT,
		  yearID int(11) NOT NULL,
		  makeID int(11) NOT NULL,
		  modelID int(11) NOT NULL,
		  styleID int(11) NOT NULL,
		  dateAdded timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		  PRIMARY KEY (vehicleID)
		) ENGINE=InnoDB AUTO_INCREMENT=250265 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`

	db, err := sql.Open("mysql", ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
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
	return nil
}

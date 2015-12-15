package importer

import (
	"database/sql"
	"github.com/curt-labs/heavyduty/database"
)

func addToYearMap(year float64) (int, error) {
	var yearID int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return yearID, err
	}
	defer db.Close()

	//TODO -- all this crap
	return yearID, err

}

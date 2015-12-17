package merger

import (
	"fmt"
	"os"
	"strings"
)

//Merge puts TEMP table data into CurtData DB
func Merge() error {
	var input string
	_ = fmt.Sprintf("Merge TEMP tables into CurtData at '%s' (y/n)?", os.Getenv("DATABASE_HOST"))
	_, err := fmt.Scanf("%s", &input)
	if err != nil {
		return err
	}
	if strings.ToLower(input) != "y" {
		fmt.Println("Exiting program.")
		return nil
	}

	return nil
}

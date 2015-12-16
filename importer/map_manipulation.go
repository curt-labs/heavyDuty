package importer

import (
	"fmt"
	"strconv"
	"strings"
)

//user chooses the make, model, style they want from the map
//that choice is assigned to map
func chooseFromMap(mapType string, vehicleAspect string) int {
	var m map[string]int
	switch {
	case mapType == "style":
		m = styleMap
	case mapType == "make":
		m = makeMap
	case mapType == "model":
		m = modelMap
	}
	var input string
	var output int
	var err error
	fmt.Printf("Choose style from available %s? y/n: ", mapType)
	if _, err := fmt.Scanf("%s", &input); err != nil {
		return 0
	}
	if strings.ToLower(input) != "y" {
		return 0
	}
	for key, value := range m {
		fmt.Printf("%s...(%d) \n", key, value)
	}
	fmt.Print("Choose (NUMBER) above or press (0) to cancel: ")

	if _, err := fmt.Scanf("%s", &input); err != nil {
		return 0
	}

	output, err = strconv.Atoi(input)
	if err != nil || output == 0 {
		return 0
	}
	//save to map - vehicle aspect (make model style) already ToLower
	m[vehicleAspect] = output
	return output
}

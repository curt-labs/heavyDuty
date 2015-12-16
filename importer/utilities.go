package importer

import (
	"fmt"
	"os"
	"strings"
)

const (
	ERR_FILE = "ErrFile.csv"
)

var (
	offset int64
)

//makes stuff, like styles, capitalized real good
func capitalize(s string) string {
	if s == "" {
		return "" //nil is bad
	}
	inArray := func(arr []string, str string) bool {
		for _, a := range arr {
			if strings.ToLower(str) == strings.ToLower(a) {
				return true
			}
		}
		return false
	}
	strArray := strings.Split(s, " ")
	exceptions := []string{"a", "an", "the", "or"}
	for i, st := range strArray {
		if i != 0 {
			if inArray(exceptions, st) {
				strArray[i] = strings.ToLower(st)
				continue
			}
		}
		strArray[i] = strings.ToUpper(st[:1]) + strings.ToLower(st[1:])
	}
	return strings.Join(strArray, " ")
}

//write vp with message to error file
func (vp *VehiclePart) toErrFile(msg string) error {
	f, err := findFile(ERR_FILE)
	if err != nil {
		return err
	}

	n, err := f.WriteAt([]byte(fmt.Sprintf("%d,%f,%s,%s,%s,%s\n", vp.PartID, vp.Vehicle.Year, vp.Vehicle.Make, vp.Vehicle.Model, vp.Vehicle.Style, msg)), offset)
	if err != nil {
		return err
	}
	offset += int64(n)
	return nil
}

//file file by name - open for write & append
func findFile(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		f, err = os.Create(name)
		if err != nil {
			return f, err
		}
	}
	return f, err
}

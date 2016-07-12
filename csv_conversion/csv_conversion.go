package csv_conversion

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/curt-labs/heavierduty/importer"
)

type Line struct {
	Year         string   `field:"Year"`
	Make         string   `field:"Make"`
	Model        string   `field:"Model"`
	Style        string   `field:"Style"`
	PartNumber   string   `field:"PartNumber"`
	ShortDesc    string   `field:"ShortDesc"`
	InstallTime  string   `field:"InstallTime"`
	RelatedParts []string `field:"RelatedParts"`
	Drilling     string   `field:"Drilling"`
	Notes        []string `field:"Notes"`
	UPC          string   `field:"UPC"`
	List         string   `field:"List"`
	MAP          string   `field:"MAP"`
	Jobber       string   `field:"Jobber"`
	Weight       string   `field:"Weight"`
	Length       string   `field:"Length"`
	Height       string   `field:"Height"`
	Width        string   `field:"Width"`
	Bullets      []string `field:"Bullets"`
}

// GetApplications Parses CSV according to provided map and filepath
// Returns a slice of VehicleParts
func GetApplications(path string, csvMap map[string][]int) ([]importer.VehiclePart, error) {
	var vps []importer.VehiclePart
	f, err := os.Open(path)
	if err != nil {
		return vps, err
	}
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	rawData, err := reader.ReadAll()
	if err != nil {
		return vps, err
	}
	for i, row := range rawData {
		if i == 0 {
			continue //header
		}
		if row[0] == "" {
			break //end of file
		}
		line, err := processRow(row, csvMap)
		if err != nil {
			return vps, err
		}

		vp, err := line.BuildApplications()
		if err != nil {
			return vps, err
		}
		vps = append(vps, vp)
	}
	return vps, err
}

func processRow(row []string, csvMap map[string][]int) (Line, error) {
	var line Line
	value := reflect.ValueOf(&line).Elem()
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get("field")
		datumPosition := csvMap[tag]
		fieldKind := value.Field(i).Kind()

		switch fieldKind {
		case reflect.String:
			for _, v := range datumPosition {
				if row[v] != "" {
					value.Field(i).SetString(row[v])
				}
			}

		case reflect.Slice:
			var strSlice []string
			for _, v := range datumPosition {
				if row[v] != "" {
					strSlice = append(strSlice, row[v])
				}
			}
			value.Field(i).Set(reflect.ValueOf(strSlice))
		}
	}
	return line, nil
}

func (l *Line) BuildApplications() (importer.VehiclePart, error) {
	year, err := strconv.ParseFloat(l.Year, 64)
	if err != nil {
		return importer.VehiclePart{}, err
	}

	partID, err := strconv.Atoi(l.PartNumber)
	if err != nil {
		return importer.VehiclePart{}, err
	}

	var relparts []int
	for _, rel := range l.RelatedParts {
		relInt, err := strconv.Atoi(rel)
		if err != nil {
			// sometimes is 'n/a'
			continue
		}
		relparts = append(relparts, relInt)
	}

	vp := importer.VehiclePart{
		PartID:      partID,
		Drilling:    strings.TrimSpace(l.Drilling),
		InstallTime: strings.TrimSpace(l.InstallTime),
		Vehicle: importer.Vehicle{
			Year:  year,
			Make:  strings.TrimSpace(l.Make),
			Model: strings.TrimSpace(l.Model),
			Style: strings.TrimSpace(l.Style),
		},
		RelatedParts: relparts,
	}
	return vp, nil
}

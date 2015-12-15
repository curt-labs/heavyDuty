package importer

import (
	"strings"
)

func capitalize(s string) string {
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

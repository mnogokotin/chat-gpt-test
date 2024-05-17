package csv

import (
	"fmt"
	"reflect"
	"strings"
)

func GetIdFromLine(line string, csvSeparator string) string {
	lineData := strings.Split(line, csvSeparator)
	return lineData[0]
}

func GetStringSliceFromModel(model interface{}) []string {
	var line []string
	v := reflect.ValueOf(model)
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i).Interface()
		line = append(line, fmt.Sprint(fieldValue))
	}
	return line
}

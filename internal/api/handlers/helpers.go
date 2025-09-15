package handlers

import (
	"chickenTrade/API/pkg/utils"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Getting the fields name names only by removing the omiempty tag
func GetFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	fields := []string{}

	for i:=0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.Split(field.Tag.Get("json"), ",")[0]
		fmt.Println("Added Fiedls: ", fieldToAdd)
		fields = append(fields, fieldToAdd)
	}

	return fields

}

func CheckEmptyields(value interface{}) error {
	val := reflect.ValueOf(value)
	typ := val.Type()

	for i:=0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		jsonTag := fieldType.Tag.Get("json")

		if strings.Contains(jsonTag, "omitempty") {
			continue
		}

		fmt.Println("field:", field)
		if field.Kind() == reflect.String && field.String() == "" {
			return utils.ErrorHandler(errors.New("all fields aree required"), "all fields are reuired")
		}
	}
	return nil
}
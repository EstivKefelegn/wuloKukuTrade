package utils

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model) // Get the type information of all the models
	var columns, placeholders string

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println("The tag", dbTag)

		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}

			columns += dbTag
			placeholders += "?"
		}
	}
	fmt.Printf("INSERT INTO %s (%s) VALUES (%s)\n", tableName, columns, placeholders)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

// func GetStructValues(model interface{}) []interface{} {
// 	modelValue := reflect.ValueOf(model)
// 	modelType := modelValue.Type()
// 	values := []interface{}{}

// 	for i := 0; i < modelType.NumField(); i++ {
// 		dbTag := modelType.Field(i).Tag.Get("db")
// 		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
// 		if dbTag != "" && dbTag != "id" {
// 			values = append(values, modelValue.Field(i).Interface())
// 		}
		

// 	}

// 	log.Println("Values:", values)
// 	return values
// }


func GetStructValues(model interface{}) []interface{} {
    modelValue := reflect.ValueOf(model)
    modelType := modelValue.Type()
    values := []interface{}{}

    for i := 0; i < modelType.NumField(); i++ {
        dbTag := modelType.Field(i).Tag.Get("db")
        dbTag = strings.TrimSuffix(dbTag, ",omitempty")

        if dbTag != "" && dbTag != "id" {
            values = append(values, modelValue.Field(i).Interface())
        }
    }

    log.Println("Values:", values)
    return values
}

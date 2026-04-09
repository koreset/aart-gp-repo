package utils

import (
	"api/models"
	"reflect"
	"strings"
)

func StructToMapWithNonZeroFields(s models.PricingDistribution) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	for i := 0; i < val.NumField(); i++ {
		//fieldName := typ.Field(i).Tag.Get("json")
		//fieldValue := val.Field(i).Interface()
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i).Interface()

		if fieldName != "" {
			// Check if the field value is non-zero
			if !isZeroValue(fieldValue) {
				result[fieldName] = fieldValue
			}
		}
	}

	return result
}

func isZeroValue(v interface{}) bool {
	switch value := v.(type) {
	case int:
		return value == 0
	case float64:
		return value == 0.0
	default:
		return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
	}
}

func MatchCSVTags(strs []string, s interface{}) bool {
	// Get the type of the struct
	structType := reflect.TypeOf(s)

	// Loop through the strings and check if they match CSV struct tags
	for _, str := range strs {
		// Convert the string to lowercase for case-insensitive comparison
		str = strings.ToLower(str)

		// Loop through the fields of the struct
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)

			// Get the CSV struct tag value
			csvTag := field.Tag.Get("csv")

			// If the string matches the CSV struct tag, return true
			if strings.ToLower(csvTag) == str {
				return true
			}
		}
	}

	// If no match is found, return false
	return false
}

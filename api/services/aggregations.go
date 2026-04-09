package services

import (
	"api/models"
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func GetAggregationVariables() []string {
	agg := models.AggregatedProjection{}
	return getFloat64FieldJSONTags(agg)
}

func getFloat64FieldJSONTags(obj interface{}) []string {
	var float64Fields []string
	val := reflect.ValueOf(obj)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.Kind() == reflect.Float64 {
			// Get the json tag
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				float64Fields = append(float64Fields, jsonTag)
			}
			//} else {
			//	// If no json tag, append the field name (optional fallback)
			//	float64Fields = append(float64Fields, field.Name)
			//}
		}
	}

	return float64Fields
}

func GetAggregations(runId int, productCode string, spCode string, variables []string) ([]map[string]interface{}, error) {
	// Use caching to improve performance for frequently accessed aggregations
	cacheKey := fmt.Sprintf("aggregations_%d_%s_%s_%s", runId, productCode, spCode, strings.Join(variables, "_"))
	var results []map[string]interface{}

	err := QueryWithContextAndCache(cacheKey, &results, 15*time.Minute, func(ctx context.Context) error {
		// Build a query string to get the result aggregations
		query := "SELECT projection_month, "
		for i, variable := range variables {
			query += "SUM(" + variable + ") AS " + variable
			if i < len(variables)-1 {
				query += ", "
			}
		}

		// Add conditions based on parameters
		var queryParams []interface{}
		whereClause := " FROM aggregated_projections WHERE run_id = ?"
		queryParams = append(queryParams, runId)

		if productCode != "" {
			whereClause += " AND product_code = ?"
			queryParams = append(queryParams, productCode)
		}

		if spCode != "" {
			whereClause += " AND sp_code = ?"
			queryParams = append(queryParams, spCode)
		}

		// Add group by and order by
		query += whereClause + " GROUP BY projection_month ORDER BY projection_month ASC"

		// Execute the query with context
		return DB.WithContext(ctx).Table("aggregated_projections").Raw(query, queryParams...).Scan(&results).Error
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func filterStructFields(data interface{}, variables []string) map[string]interface{} {
	result := make(map[string]interface{})

	// Use reflection to access the struct fields
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	// Iterate through the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the JSON tag
		jsonTag := fieldType.Tag.Get("json")

		// Check if the jsonTag is in the variables list
		for _, variable := range variables {
			if jsonTag == variable {
				// Add the matching field to the result map
				result[jsonTag] = field.Interface()
			}
		}
	}
	return result
}

func GetAggregationSpCodes(runId string, productCode string) []string {
	var spCodes []string
	DB.Table("aggregated_projections").Where("run_id = ? AND product_code = ?", runId, productCode).Distinct("sp_code").Pluck("sp_code", &spCodes)
	return spCodes
}

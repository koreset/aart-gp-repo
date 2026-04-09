package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// StructSchema represents the schema of a struct
type StructSchema struct {
	StructName string                 `json:"struct_name"`
	Fields     map[string]FieldSchema `json:"fields"`
}

// FieldSchema represents the schema of a field
type FieldSchema struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Tag      string `json:"tag"`
	GormTag  string `json:"gorm_tag"`
	JSONTag  string `json:"json_tag"`
	Exported bool   `json:"exported"`
}

// GetStructSchema extracts the schema from a struct type
func GetStructSchema(structType reflect.Type) StructSchema {
	schema := StructSchema{
		StructName: structType.Name(),
		Fields:     make(map[string]FieldSchema),
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		
		// Skip fields that don't have a database column
		if field.Tag.Get("gorm") == "-" {
			continue
		}

		fieldSchema := FieldSchema{
			Name:     field.Name,
			Type:     field.Type.String(),
			Tag:      string(field.Tag),
			GormTag:  field.Tag.Get("gorm"),
			JSONTag:  field.Tag.Get("json"),
			Exported: field.PkgPath == "",
		}

		schema.Fields[field.Name] = fieldSchema
	}

	return schema
}

// SaveStructSchema saves the struct schema to a file
func SaveStructSchema(schema StructSchema) error {
	// Create schemas directory if it doesn't exist
	schemasDir := filepath.Join("schemas")
	if err := os.MkdirAll(schemasDir, 0755); err != nil {
		return err
	}

	// Create filename
	filename := fmt.Sprintf("%s.json", schema.StructName)
	filePath := filepath.Join(schemasDir, filename)

	// Marshal schema to JSON
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

// LoadStructSchema loads the struct schema from a file
func LoadStructSchema(structName string) (StructSchema, error) {
	var schema StructSchema

	// Create filename
	filename := fmt.Sprintf("%s.json", structName)
	filePath := filepath.Join("schemas", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return schema, fmt.Errorf("schema file for struct %s does not exist", structName)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return schema, err
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, &schema); err != nil {
		return schema, err
	}

	return schema, nil
}

// CompareStructSchemas compares two struct schemas and returns the differences
func CompareStructSchemas(oldSchema, newSchema StructSchema) (added, removed, modified []string) {
	// Check for added and modified fields
	for fieldName, newField := range newSchema.Fields {
		if oldField, exists := oldSchema.Fields[fieldName]; !exists {
			// Field doesn't exist in old schema, so it's added
			added = append(added, fieldName)
		} else if oldField.Type != newField.Type || oldField.GormTag != newField.GormTag {
			// Field exists but has changed type or gorm tag
			modified = append(modified, fieldName)
		}
	}

	// Check for removed fields
	for fieldName := range oldSchema.Fields {
		if _, exists := newSchema.Fields[fieldName]; !exists {
			// Field exists in old schema but not in new schema, so it's removed
			removed = append(removed, fieldName)
		}
	}

	return added, removed, modified
}
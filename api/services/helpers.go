package services

import (
	"api/models"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	"gorm.io/gorm"
)

func CalculatePV(interest float64, term float64, instalment float64) float64 {
	var presentValue float64

	for i := term; i > 0; i-- {
		factor := math.Pow(1+interest, -float64(i)/12.0)

		pv := (presentValue + instalment) * factor
		presentValue = pv
	}

	return presentValue
}

func colIndexToExcelColName(index int) string {
	result := ""
	for index >= 0 {
		remainder := index % 26
		result = string(rune('A'+remainder)) + result
		index = (index / 26) - 1
	}
	return result
}

func convertUint8ToType(v []uint8) interface{} {
	strVal := string(v)

	// Try to convert to an integer
	if intVal, err := strconv.Atoi(strVal); err == nil {
		return intVal
	}

	// Try to convert to a float
	if floatVal, err := strconv.ParseFloat(strVal, 64); err == nil {
		return floatVal
	}

	// Otherwise, return as string
	return strVal
}

func exportTableToExcel(query string) ([]byte, error) {
	// Execute the custom query to get the rows
	rows, err := DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names from the rows
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create a new Excel file
	f := excelize.NewFile()
	// Use the default sheet that already exists in a new file
	sheetName := f.GetSheetName(f.GetActiveSheetIndex())

	// Use StreamWriter for large datasets for better performance and lower memory usage
	sw, err := f.NewStreamWriter(sheetName)
	if err != nil {
		return nil, err
	}

	// Write the column headers in a single call
	header := make([]interface{}, len(columns))
	for i, colName := range columns {
		header[i] = colName
	}
	if err := sw.SetRow("A1", header); err != nil {
		return nil, err
	}

	// Write rows from the database efficiently
	rowIndex := 2
	for rows.Next() {
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Build the row data with proper types
		rowData := make([]interface{}, len(columns))
		for i, value := range columnValues {
			switch v := value.(type) {
			case nil:
				rowData[i] = ""
			case []uint8:
				// Convert byte slice to appropriate type (int, float, or string)
				rowData[i] = convertUint8ToType(v)
			case int64:
				rowData[i] = v
			case float64:
				rowData[i] = v
			case bool:
				rowData[i] = v
			case string:
				rowData[i] = v
			default:
				rowData[i] = fmt.Sprintf("%v", v)
			}
		}

		axis, _ := excelize.CoordinatesToCellName(1, rowIndex)
		if err := sw.SetRow(axis, rowData); err != nil {
			return nil, err
		}

		rowIndex++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := sw.Flush(); err != nil {
		return nil, err
	}

	// Package the workbook into an XLSX (ZIP of XML parts). This step can be CPU-heavy for large sheets due to ZIP compression.
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExtractHeaders retrieves JSON tags as headers in Title Case
func extractHeaders(model interface{}) []string {
	var headers []string
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			headers = append(headers, jsonTag)
		} else {
			headers = append(headers, field.Name) // Fallback to field name if JSON tag is missing
		}
	}
	return headers
}

// ExtractValues retrieves values from the struct as an array
func extractValues(model interface{}) []interface{} {
	var values []interface{}
	v := reflect.ValueOf(model)
	for i := 0; i < v.NumField(); i++ {
		values = append(values, v.Field(i).Interface())
	}
	return values
}

func zipMultipleFiles(files map[string][]byte) ([]byte, error) {
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	// Create each file inside the zip archive
	for filename, data := range files {
		zipFile, err := zipWriter.Create(filename)
		if err != nil {
			return nil, err
		}

		// Write data to the file
		_, err = zipFile.Write(data)
		if err != nil {
			return nil, err
		}
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ===== Generic audit helpers =====

type AuditContext struct {
	Area      string
	Entity    string
	EntityID  string
	Action    string // CREATE|UPDATE|DELETE
	Route     string
	RequestID string
	ChangedBy string
}

type Diff struct {
	Field string      `json:"field"`
	From  interface{} `json:"from"`
	To    interface{} `json:"to"`
}

func computeDiff(before interface{}, after interface{}) ([]Diff, map[string]interface{}, map[string]interface{}) {
	bj, _ := json.Marshal(before)
	aj, _ := json.Marshal(after)
	var bm, am map[string]interface{}
	_ = json.Unmarshal(bj, &bm)
	_ = json.Unmarshal(aj, &am)

	diffs := []Diff{}
	// Changes in common keys
	for k, bv := range bm {
		if av, ok := am[k]; ok {
			if !reflect.DeepEqual(bv, av) {
				diffs = append(diffs, Diff{Field: k, From: bv, To: av})
			}
		}
	}
	// New keys in after
	for k, av := range am {
		if _, ok := bm[k]; !ok {
			diffs = append(diffs, Diff{Field: k, From: nil, To: av})
		}
	}
	// Keys removed in after
	for k, bv := range bm {
		if _, ok := am[k]; !ok {
			diffs = append(diffs, Diff{Field: k, From: bv, To: nil})
		}
	}
	return diffs, bm, am
}

func writeAudit(tx *gorm.DB, ctx AuditContext, prev interface{}, next interface{}) error {
	diffs, pm, nm := computeDiff(prev, next)
	prevJSON, _ := json.Marshal(pm)
	nextJSON, _ := json.Marshal(nm)
	diffJSON, _ := json.Marshal(diffs)

	rec := models.AuditLog{
		Area:       ctx.Area,
		Entity:     ctx.Entity,
		EntityID:   ctx.EntityID,
		Action:     ctx.Action,
		Route:      ctx.Route,
		RequestID:  ctx.RequestID,
		PrevValues: string(prevJSON),
		NewValues:  string(nextJSON),
		Diff:       string(diffJSON),
		ChangedBy:  ctx.ChangedBy,
		ChangedAt:  time.Now(),
	}
	return tx.Create(&rec).Error
}

// GetAuditLogs returns generic audit logs filtered by area/entity/id
func GetAuditLogs(area, entity, entityID string, limit, offset int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	q := DB.Where("area = ?", area)
	if entity != "" {
		q = q.Where("entity = ?", entity)
	}
	if entityID != "" {
		q = q.Where("entity_id = ?", entityID)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	q = q.Order("changed_at DESC")
	if err := q.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetMemberFullAuditHistory returns audit history rows for a given member id within the
// group-pricing area. When includeAssociated is true, it currently returns the same
// consolidated logs filtered only by area and entity_id, covering all entities that
// logged with this member's id as their entity_id (e.g., g_pricing_member_data_in_forces).
// Pagination is controlled by limit and offset.
func GetMemberFullAuditHistory(memberID int, includeAssociated bool, limit, offset int) ([]models.AuditLog, error) {
	// For now, we return all audit logs in the group-pricing area for this entity_id,
	// regardless of entity type. This aligns with the controller's need to fetch
	// associated history and is consistent with how writeAudit sets EntityID for
	// member-related logs.
	return GetAuditLogs("group-pricing", "", strconv.Itoa(memberID), limit, offset)
}

// GetMemberActivityHistory returns structured activity history for a member.
func GetMemberActivityHistory(memberID string, limit, offset int) ([]models.MemberActivity, error) {
	var activities []models.MemberActivity
	q := DB.Model(&models.MemberActivity{}).
		Where("member_id = ? OR member_id_number = ?", memberID, memberID).
		Order("timestamp DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

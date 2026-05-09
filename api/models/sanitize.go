package models

import (
	"math"
	"reflect"

	"gorm.io/gorm"
)

// sanitizeFloat64NaNInf walks v in-place, replacing every float64 field whose
// value is NaN or ±Inf with 0. It returns the dotted names of fields that
// were scrubbed so callers can log them.
//
// MySQL's text protocol cannot represent NaN/Inf, so a non-finite float in
// any persisted struct surfaces as `Error 1054 (42S22): Unknown column 'NaN'`
// when GORM renders the row.
func sanitizeFloat64NaNInf(v reflect.Value) []string {
	return sanitizeFloat64NaNInfWithPrefix(v, "")
}

func sanitizeFloat64NaNInfWithPrefix(v reflect.Value, prefix string) []string {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	var scrubbed []string
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldT := t.Field(i)
		if !fieldT.IsExported() {
			continue
		}
		fv := v.Field(i)
		name := fieldT.Name
		if prefix != "" {
			name = prefix + "." + name
		}
		switch fv.Kind() {
		case reflect.Float64, reflect.Float32:
			f := fv.Float()
			if math.IsNaN(f) || math.IsInf(f, 0) {
				if fv.CanSet() {
					fv.SetFloat(0)
					scrubbed = append(scrubbed, name)
				}
			}
		case reflect.Struct:
			scrubbed = append(scrubbed, sanitizeFloat64NaNInfWithPrefix(fv, name)...)
		case reflect.Pointer:
			if !fv.IsNil() && fv.Elem().Kind() == reflect.Struct {
				scrubbed = append(scrubbed, sanitizeFloat64NaNInfWithPrefix(fv.Elem(), name)...)
			}
		}
	}
	return scrubbed
}

// BeforeSave scrubs any NaN/Inf float fields to 0 before persisting and logs
// the offending field names so the upstream calc that produced the bad value
// can be traced.
func (m *MemberRatingResultSummary) BeforeSave(tx *gorm.DB) error {
	scrubbed := sanitizeFloat64NaNInf(reflect.ValueOf(m).Elem())
	if len(scrubbed) > 0 {
		tx.Logger.Warn(tx.Statement.Context,
			"MemberRatingResultSummary: sanitized non-finite float fields to 0 before save (quote_id=%d, fields=%v)",
			m.QuoteId, scrubbed)
	}
	return nil
}

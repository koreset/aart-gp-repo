package models

import (
	"strconv"
	"strings"
)

// CsvFloat is a float64 wrapper used by csvutil-decoded structs that should be
// tolerant of blank values. csvutil's default float64 decoder fails on "" with
// "cannot unmarshal \"\" into Go value of type float64"; that bails out of the
// row loop and bypasses the structured presence-check error path. With CsvFloat,
// an empty/whitespace cell is left at its zero value and the downstream
// validation layer decides whether that constitutes a blocking error.
type CsvFloat float64

func (f *CsvFloat) UnmarshalCSV(b []byte) error {
	value := strings.Trim(string(b), `"`)
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, "null") {
		return nil
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	*f = CsvFloat(v)
	return nil
}

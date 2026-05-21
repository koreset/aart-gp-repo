package services

import (
	"api/models"
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetCpiIndices returns all CPI rows ordered most-recent first. The dataset
// is small (≤ ~12 rows per year), so we return the whole table; callers
// page on the client.
func GetCpiIndices() ([]models.CpiIndex, error) {
	var rows []models.CpiIndex
	if err := DB.Order("year_index DESC, month_index DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// UpsertCpiIndex writes a single CPI row, overwriting any existing row for
// the same (year, month). On conflict we also refresh created_by so the
// column reflects the most recent author rather than getting stuck blank on
// rows that pre-date the column.
func UpsertCpiIndex(row models.CpiIndex) (models.CpiIndex, error) {
	if err := validateCpiRow(row); err != nil {
		return row, err
	}
	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "year_index"}, {Name: "month_index"}},
		DoUpdates: clause.AssignmentColumns([]string{"cpi_index", "created_by"}),
	}).Create(&row).Error; err != nil {
		return row, err
	}
	// Re-fetch by natural key so the caller gets the persisted id (which the
	// upsert path may not populate on all dialects).
	var stored models.CpiIndex
	if err := DB.Where("year_index = ? AND month_index = ?", row.YearIndex, row.MonthIndex).
		First(&stored).Error; err != nil {
		return row, err
	}
	return stored, nil
}

// BulkUpsertCpiIndices writes many rows in a single transaction, upserting
// on (year, month). Validation runs over the whole batch first so a bad row
// rejects the entire upload.
func BulkUpsertCpiIndices(rows []models.CpiIndex) (int, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	for i, r := range rows {
		if err := validateCpiRow(r); err != nil {
			return 0, fmt.Errorf("row %d: %w", i+1, err)
		}
	}
	return len(rows), DB.Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "year_index"}, {Name: "month_index"}},
			DoUpdates: clause.AssignmentColumns([]string{"cpi_index", "created_by"}),
		}).CreateInBatches(rows, 200).Error
	})
}

func validateCpiRow(r models.CpiIndex) error {
	if r.YearIndex < 1900 || r.YearIndex > 2100 {
		return fmt.Errorf("year_index %d outside accepted range 1900-2100", r.YearIndex)
	}
	if r.MonthIndex < 1 || r.MonthIndex > 12 {
		return fmt.Errorf("month_index %d outside 1-12", r.MonthIndex)
	}
	if r.CpiIndex <= 0 {
		return fmt.Errorf("cpi_index must be positive")
	}
	return nil
}

// ParseCpiUploadFile reads either a .csv or .xlsx file and returns CpiIndex
// rows. Expected headers (case/space insensitive): year_index, month_index,
// cpi_index. The file format is detected by extension.
func ParseCpiUploadFile(path, fileName string) ([]models.CpiIndex, error) {
	lower := strings.ToLower(fileName)
	switch {
	case strings.HasSuffix(lower, ".xlsx"):
		return parseCpiExcel(path)
	case strings.HasSuffix(lower, ".csv"):
		return parseCpiCSV(path)
	default:
		return nil, fmt.Errorf("unsupported file format: only .csv and .xlsx are accepted")
	}
}

func parseCpiCSV(path string) ([]models.CpiIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	delim := detectCpiDelimiter(f)
	reader := csv.NewReader(f)
	reader.Comma = delim
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}
	return rowsToCpiIndices(rows[0], rows[1:])
}

func parseCpiExcel(path string) ([]models.CpiIndex, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("workbook contains no sheets")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}
	return rowsToCpiIndices(rows[0], rows[1:])
}

func detectCpiDelimiter(f *os.File) rune {
	delimiters := []rune{',', ';', '|', '\t'}
	maxCols := 0
	best := ','
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		f.Seek(0, 0)
		return ','
	}
	line := scanner.Text()
	for _, d := range delimiters {
		r := csv.NewReader(strings.NewReader(line))
		r.Comma = d
		if cells, err := r.ReadAll(); err == nil && len(cells) > 0 {
			if len(cells[0]) > maxCols {
				maxCols = len(cells[0])
				best = d
			}
		}
	}
	f.Seek(0, 0)
	return best
}

// rowsToCpiIndices maps a (header, dataRows) pair into CpiIndex structs.
// Accepts common header variants so a file copy-pasted from a published
// stats spreadsheet doesn't need pre-cleaning.
func rowsToCpiIndices(headers []string, dataRows [][]string) ([]models.CpiIndex, error) {
	norm := make([]string, len(headers))
	for i, h := range headers {
		norm[i] = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(h, " ", "_")))
	}
	col := func(candidates ...string) int {
		for _, c := range candidates {
			for i, h := range norm {
				if h == c {
					return i
				}
			}
		}
		return -1
	}
	idxYear := col("year_index", "year")
	idxMonth := col("month_index", "month")
	idxCpi := col("cpi_index", "cpi", "value", "index")
	if idxYear < 0 || idxMonth < 0 || idxCpi < 0 {
		return nil, fmt.Errorf("missing required columns: year_index, month_index, cpi_index")
	}

	getCell := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	out := make([]models.CpiIndex, 0, len(dataRows))
	for i, row := range dataRows {
		yearStr := getCell(row, idxYear)
		monthStr := getCell(row, idxMonth)
		cpiStr := getCell(row, idxCpi)
		if yearStr == "" && monthStr == "" && cpiStr == "" {
			continue // skip blank rows
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid year_index %q", i+2, yearStr)
		}
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid month_index %q", i+2, monthStr)
		}
		cpi, err := strconv.ParseFloat(strings.ReplaceAll(cpiStr, ",", ""), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid cpi_index %q", i+2, cpiStr)
		}
		out = append(out, models.CpiIndex{
			YearIndex:  year,
			MonthIndex: month,
			CpiIndex:   cpi,
		})
	}
	return out, nil
}

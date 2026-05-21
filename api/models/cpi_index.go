package models

import "time"

// CpiIndex captures the published Consumer Price Index value for a single
// month. The (YearIndex, MonthIndex) pair is the natural key; uploads upsert
// on that pair so re-importing a file does not duplicate rows.
type CpiIndex struct {
	ID         int       `json:"id" gorm:"primary_key"`
	YearIndex  int       `json:"year_index" gorm:"index:idx_cpi_indices_year_month,unique"`
	MonthIndex int       `json:"month_index" gorm:"index:idx_cpi_indices_year_month,unique"`
	CpiIndex   float64   `json:"cpi_index"`
	Created    time.Time `json:"created" gorm:"autoCreateTime"`
	CreatedBy  string    `json:"created_by" gorm:"size:128"`
}

func (CpiIndex) TableName() string { return "cpi_indices" }

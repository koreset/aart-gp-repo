package models

type BaseMortalityBand struct {
	ID         int     `gorm:"primary_key"`
	LowerBound int     `gorm:"unique_index" json:"lower_bound"`
	UpperBound int     `json:"upper_bound"`
	Male       float64 `json:"male"`
	Female     float64 `json:"female"`
}

package models

type OlympicWinner struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Athlete      string `json:"athlete"`
	Age          int    `json:"age"`
	Country      string `json:"country"`
	CountryGroup string `json:"country_group"`
	Year         int    `json:"year"`
	Date         string `json:"date"`
	Sport        string `json:"sport"`
	Gold         int    `json:"gold"`
	Silver       int    `json:"silver"`
	Bronze       int    `json:"bronze"`
	Total        int    `json:"total"`
}

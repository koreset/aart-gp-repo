package models

type RetrenchmentRate struct {
	ID                   uint    `gorm:"primary_key" json:"id"`
	DurationInForceYears uint    `json:"duration_in_force_years" gorm:"column:duration_if_y"`
	Value                float64 `json:"value" gorm:"column:retr_rate"`
}

type DisabilityIncidenceFactors struct {
	ID                  int
	Age                 int
	SocialEconomicClass int
	OccupationalClass   int
	Gender              string
}

type ColumnNames struct {
	Field string
	Type  string
}

package models

type TransitionState struct {
	TransitionID     int    `gorm:"primary_key"`
	ProductCode      string `json:"product_code"`
	StartState       string `json:"start_state"`
	EndState         string `json:"end_state"`
	BaseTable        string `json:"base_table"`
	AdditionalTable1 string `json:"additional_table_1"`
	AdditionalTable2 string `json:"additional_table_2"`
	AdditionalTable3 string `json:"additional_table_3"`
	Absorbing        bool   `json:"absorbing"`
}

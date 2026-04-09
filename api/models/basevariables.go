package models

import "api/enums"

type BaseAssumptionVariable struct {
	ID       int                `gorm:"primary_key" json:"id"`
	Name     string             `json:"name" gorm:"unique_index:idx_ass_id"`
	DataType enums.VariableType `json:"data_type"`
	Value    float64            `json:"value"`
}

type BaseFeature struct {
	ID          int                `gorm:"primary_key" json:"id"`
	Name        string             `json:"name" `
	Description string             `json:"description"`
	DataType    enums.VariableType `json:"data_type"`
}

package models

import "api/enums"

//ModelPointVariable example
type ModelPointVariable struct {
	ID         int                `gorm:"primary_key" json:"id" example:"1"`
	Name       string             `json:"name" gorm:"unique" example:"Policy Number"`
	Code       string             `json:"code" gorm:"unique" example:"POLICY_NUMBER"`
	DataType   enums.VariableType `json:"type" example:"1"`
	Funeral    bool               `json:"funeral"`
	CreditLife bool               `json:"credit_life"`
}

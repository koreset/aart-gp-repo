package models

// MarkovState example
type MarkovState struct {
	ID    int    `json:"id" gorm:"primary_key" example:"1"`
	State string `json:"state" form:"state" example:"Healthy"`
}

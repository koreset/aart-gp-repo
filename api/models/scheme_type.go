package models

import "time"

// SchemeType represents a classification such as Executive, Management, General, etc.
// It can be referenced by GroupPricingQuote via SchemeTypeID (added via migration when needed).
type SchemeType struct {
    ID          int       `json:"id" gorm:"primary_key"`
    Name        string    `json:"name" gorm:"uniqueIndex;size:100"`
    Description string    `json:"description" gorm:"size:500"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func (SchemeType) TableName() string { return "scheme_types" }

package models

import "time"

type IFRS17Amendment struct {
	ID               int        `json:"id" gorm:"primaryKey;autoIncrement"`
	OriginalRunID    int        `json:"original_run_id"`
	AmendmentRunID   int        `json:"amendment_run_id"`
	AmendmentType    string     `json:"amendment_type"`
	Reason           string     `json:"reason"`
	Status           string     `json:"status" gorm:"default:draft"`
	CreatedBy        string     `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	ApprovedBy       string     `json:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at"`
	OriginalRunName  string     `json:"original_run_name" gorm:"-"`
	AmendmentRunName string     `json:"amendment_run_name" gorm:"-"`
}

type CreateAmendmentRequest struct {
	OriginalRunID  int    `json:"original_run_id"`
	AmendmentRunID int    `json:"amendment_run_id"`
	AmendmentType  string `json:"amendment_type"`
	Reason         string `json:"reason"`
}

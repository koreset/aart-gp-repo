package models

import "time"

type TransitionAdjustment struct {
	ID             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UploadBatch    string     `json:"upload_batch"`
	ProductCode    string     `json:"product_code"`
	IFRS17Group    string     `json:"ifrs17_group"`
	AdjustmentType string     `json:"adjustment_type"`
	Amount         float64    `json:"amount"`
	Status         string     `json:"status" gorm:"default:pending"`
	ReviewedBy     string     `json:"reviewed_by"`
	ReviewedAt     *time.Time `json:"reviewed_at"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
}

type TransitionAdjustmentReviewRequest struct {
	Notes string `json:"notes"`
}

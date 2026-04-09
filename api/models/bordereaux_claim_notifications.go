package models

import "time"

// ClaimNotificationLog tracks the cadence of claim notifications sent to reinsurers.
// NotificationType: initial | status_update | final
// Status: pending | sent | acknowledged | overdue
type ClaimNotificationLog struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
	ClaimID       int    `json:"claim_id" gorm:"default:0;index"`
	ClaimNumber   string `json:"claim_number" gorm:"default:''"`
	SchemeID      int    `json:"scheme_id" gorm:"default:0;index"`
	SchemeName    string `json:"scheme_name" gorm:"default:''"`
	ReinsurerName string `json:"reinsurer_name" gorm:"default:''"`
	ReinsurerCode string `json:"reinsurer_code" gorm:"default:''"`
	// NotificationType: initial | status_update | final
	NotificationType string `json:"notification_type" gorm:"not null"`
	// Status: pending | sent | acknowledged | overdue
	Status         string     `json:"status" gorm:"default:'pending'"`
	DueDate        string     `json:"due_date" gorm:"default:''"`
	SentAt         *time.Time `json:"sent_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	Notes          string     `json:"notes" gorm:"type:text"`
	CreatedBy      string     `json:"created_by" gorm:"default:''"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreateClaimNotificationRequest is the request body for creating a notification log entry.
type CreateClaimNotificationRequest struct {
	ClaimID          int    `json:"claim_id"`
	ClaimNumber      string `json:"claim_number" binding:"required"`
	SchemeID         int    `json:"scheme_id"`
	SchemeName       string `json:"scheme_name"`
	ReinsurerName    string `json:"reinsurer_name"`
	ReinsurerCode    string `json:"reinsurer_code"`
	NotificationType string `json:"notification_type" binding:"required"`
	DueDate          string `json:"due_date"`
	Notes            string `json:"notes"`
}

// MarkNotificationSentRequest is the request body for marking a notification as sent or acknowledged.
type MarkNotificationSentRequest struct {
	Notes string `json:"notes"`
}

// GenerateMonthEndNotificationsRequest drives auto-generation of month-end notifications.
type GenerateMonthEndNotificationsRequest struct {
	SchemeID int `json:"scheme_id" binding:"required"`
	Month    int `json:"month"`
	Year     int `json:"year"`
}

// NotificationStats aggregates counts by status for dashboard display.
type NotificationStats struct {
	Total        int `json:"total"`
	Pending      int `json:"pending"`
	Sent         int `json:"sent"`
	Acknowledged int `json:"acknowledged"`
	Overdue      int `json:"overdue"`
}

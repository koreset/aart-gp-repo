package models

import "time"

// BordereauxDeadline tracks the expected submission deadline for a scheme's monthly data.
// Status lifecycle: pending → received (auto on upload) | overdue | waived
type BordereauxDeadline struct {
	ID                int        `gorm:"primaryKey;autoIncrement" json:"id"`
	SchemeID          int        `json:"scheme_id" gorm:"not null;index"`
	SchemeName        string     `json:"scheme_name" gorm:"default:''"`
	Month             int        `json:"month" gorm:"not null"`
	Year              int        `json:"year" gorm:"not null"`
	// Type: member_submission (extensible to claims_notification etc. in future phases)
	Type              string     `json:"type" gorm:"default:'member_submission'"`
	DueDate           string     `json:"due_date" gorm:"default:''"`
	GracePeriodDays   int        `json:"grace_period_days" gorm:"default:0"`
	// Status: pending | received | overdue | waived
	Status            string     `json:"status" gorm:"default:'pending'"`
	LinkedSubmissionID *int      `json:"linked_submission_id"`
	WaivedBy          string     `json:"waived_by" gorm:"default:''"`
	WaivedAt          *time.Time `json:"waived_at"`
	WaiverReason      string     `json:"waiver_reason" gorm:"default:''"`
	CreatedBy         string     `json:"created_by" gorm:"default:''"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// CreateDeadlineRequest is the request body for creating a single deadline.
type CreateDeadlineRequest struct {
	SchemeID        int    `json:"scheme_id" binding:"required"`
	Month           int    `json:"month" binding:"required"`
	Year            int    `json:"year" binding:"required"`
	DueDate         string `json:"due_date" binding:"required"`
	GracePeriodDays int    `json:"grace_period_days"`
	Type            string `json:"type"`
}

// GenerateDeadlinesRequest is the request body for bulk-generating deadlines for all in-force schemes.
type GenerateDeadlinesRequest struct {
	Month           int `json:"month" binding:"required"`
	Year            int `json:"year" binding:"required"`
	DueDayOfMonth   int `json:"due_day_of_month"` // defaults to 15
	GracePeriodDays int `json:"grace_period_days"`
}

// UpdateDeadlineStatusRequest is the request body for patching a deadline's status.
type UpdateDeadlineStatusRequest struct {
	Status       string `json:"status" binding:"required"` // waived | pending
	WaiverReason string `json:"waiver_reason"`
}

// DeadlineStats is the summary returned by the stats endpoint.
type DeadlineStats struct {
	OverdueCount  int `json:"overdue_count"`
	PendingCount  int `json:"pending_count"`
	ReceivedCount int `json:"received_count"`
	WaivedCount   int `json:"waived_count"`
}

// GenerateDeadlinesResult summarises the outcome of a bulk deadline generation.
type GenerateDeadlinesResult struct {
	Total   int `json:"total"`
	Created int `json:"created"`
	Skipped int `json:"skipped"`
}

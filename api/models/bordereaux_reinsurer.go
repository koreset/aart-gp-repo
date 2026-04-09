package models

import "time"

// ReinsurerAcceptance tracks whether a reinsurer has accepted a submitted bordereaux
type ReinsurerAcceptance struct {
	ID                    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	GeneratedBordereauxID string `json:"generated_bordereaux_id" gorm:"index;size:255"`
	ReinsurerName         string `json:"reinsurer_name" gorm:"default:''"`
	ReinsurerCode         string `json:"reinsurer_code" gorm:"default:''"`
	// Status: pending | accepted | queried | rejected
	Status          string     `json:"status" gorm:"default:'pending'"`
	QueryDetails    string     `json:"query_details" gorm:"type:text"`
	SubmittedAmount float64    `json:"submitted_amount" gorm:"default:0"`
	AcceptedAmount  float64    `json:"accepted_amount" gorm:"default:0"`
	Variance        float64    `json:"variance" gorm:"default:0"`
	DueDate         string     `json:"due_date" gorm:"default:''"`
	ReceivedDate    *time.Time `json:"received_date"`
	Notes           string     `json:"notes" gorm:"type:text"`
	UpdatedBy       string     `json:"updated_by" gorm:"default:''"`
	CreatedBy       string     `json:"created_by" gorm:"default:''"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ReinsurerRecovery tracks claim recovery amounts received from reinsurers
type ReinsurerRecovery struct {
	ID                    int     `json:"id" gorm:"primaryKey;autoIncrement"`
	GeneratedBordereauxID string  `json:"generated_bordereaux_id" gorm:"index;size:255"`
	ClaimReference        string  `json:"claim_reference" gorm:"default:'';index"`
	ReinsurerName         string  `json:"reinsurer_name" gorm:"default:''"`
	ReinsurerCode         string  `json:"reinsurer_code" gorm:"default:''"`
	ClaimAmount           float64 `json:"claim_amount" gorm:"default:0"`
	RecoveredAmount       float64 `json:"recovered_amount" gorm:"default:0"`
	RecoveryPercentage    float64 `json:"recovery_percentage" gorm:"default:0"`
	// Status: pending | partial | full | disputed
	Status       string     `json:"status" gorm:"default:'pending'"`
	ReceivedDate *time.Time `json:"received_date"`
	Notes        string     `json:"notes" gorm:"type:text"`
	RecordedBy   string     `json:"recorded_by" gorm:"default:''"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// CreateReinsurerAcceptanceRequest is the payload for creating a new acceptance record
type CreateReinsurerAcceptanceRequest struct {
	GeneratedBordereauxID string  `json:"generated_bordereaux_id" binding:"required"`
	ReinsurerName         string  `json:"reinsurer_name" binding:"required"`
	ReinsurerCode         string  `json:"reinsurer_code"`
	SubmittedAmount       float64 `json:"submitted_amount"`
	DueDate               string  `json:"due_date"`
	Notes                 string  `json:"notes"`
}

// UpdateReinsurerAcceptanceRequest is the payload for updating an acceptance record
type UpdateReinsurerAcceptanceRequest struct {
	Status         string     `json:"status"`
	AcceptedAmount float64    `json:"accepted_amount"`
	QueryDetails   string     `json:"query_details"`
	Notes          string     `json:"notes"`
	ReceivedDate   *time.Time `json:"received_date"`
}

// CreateReinsurerRecoveryRequest is the payload for recording a recovery
type CreateReinsurerRecoveryRequest struct {
	GeneratedBordereauxID string  `json:"generated_bordereaux_id" binding:"required"`
	ClaimReference        string  `json:"claim_reference" binding:"required"`
	ReinsurerName         string  `json:"reinsurer_name" binding:"required"`
	ReinsurerCode         string  `json:"reinsurer_code"`
	ClaimAmount           float64 `json:"claim_amount"`
	RecoveredAmount       float64 `json:"recovered_amount"`
	Notes                 string  `json:"notes"`
}

// UpdateReinsurerRecoveryRequest is the payload for updating a recovery record
type UpdateReinsurerRecoveryRequest struct {
	RecoveredAmount float64    `json:"recovered_amount"`
	Status          string     `json:"status"`
	Notes           string     `json:"notes"`
	ReceivedDate    *time.Time `json:"received_date"`
}

// AcceptanceStats summarises acceptance status across all reinsurers for a bordereaux
type AcceptanceStats struct {
	Total          int     `json:"total"`
	Pending        int     `json:"pending"`
	Accepted       int     `json:"accepted"`
	Queried        int     `json:"queried"`
	Rejected       int     `json:"rejected"`
	TotalSubmitted float64 `json:"total_submitted"`
	TotalAccepted  float64 `json:"total_accepted"`
	TotalVariance  float64 `json:"total_variance"`
}

// RecoveryStats summarises recovery amounts for a bordereaux
type RecoveryStats struct {
	Total            int     `json:"total"`
	Pending          int     `json:"pending"`
	Partial          int     `json:"partial"`
	Full             int     `json:"full"`
	Disputed         int     `json:"disputed"`
	TotalClaimAmount float64 `json:"total_claim_amount"`
	TotalRecovered   float64 `json:"total_recovered"`
	TotalOutstanding float64 `json:"total_outstanding"`
}

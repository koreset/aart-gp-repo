package models

import "time"

// BAVVerificationLog records every bank account verification call made
// through services/bav — both successes and failures — so auditors can
// reconstruct the timeline for any claim. RequestPayload and ResponsePayload
// are stored as plain JSON strings via models.JSON (no datatypes.JSON
// dependency) for portability across MySQL, PostgreSQL, and SQL Server.
//
// Retention policy is still open per the Phase 5 plan; when it lands,
// enforce at read-time via an RLS-style filter or a scheduled purge job.
type BAVVerificationLog struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ClaimID           *int      `json:"claim_id" gorm:"index"`
	Provider          string    `json:"provider" gorm:"size:64;index"`
	ProviderRequestID string    `json:"provider_request_id" gorm:"size:128"`
	IdempotencyKey    string    `json:"idempotency_key" gorm:"size:128;index:idx_bav_log_key_created"`
	Status            string    `json:"status" gorm:"size:32"`
	RequestPayload    JSON      `json:"request_payload"`
	ResponsePayload   JSON      `json:"response_payload"`
	ErrorMessage      string    `json:"error_message" gorm:"size:1024"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime;index:idx_bav_log_key_created"`
}

// TableName pins the GORM table name so a future model rename won't silently
// migrate data onto a new table.
func (BAVVerificationLog) TableName() string {
	return "bav_verification_logs"
}

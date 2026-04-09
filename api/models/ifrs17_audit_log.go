package models

import "time"

// IFRS17AuditLog records every significant event in the IFRS 17 workflow —
// run lifecycle transitions, configuration changes, file uploads, etc.
type IFRS17AuditLog struct {
	ID         int       `json:"id"          gorm:"primaryKey;autoIncrement"`
	EventType  string    `json:"event_type"`  // run_created, run_deleted, run_reviewed, run_approved, run_locked, run_returned_to_draft, config_saved, config_deleted, risk_driver_saved, risk_driver_deleted, file_uploaded, file_deleted
	ObjectType string    `json:"object_type"` // csm_run, aos_config, risk_driver, finance_file, transition_file, sap_file
	ObjectID   int       `json:"object_id"`
	ObjectName string    `json:"object_name"`
	ChangedBy  string    `json:"changed_by"`
	ChangedAt  time.Time `json:"changed_at"`
	Details    string    `json:"details"` // JSON string or human-readable description
}

package models

import "time"

// SanctionsScreening records the result of a sanctions / PEP / adverse-media
// check on a payment schedule line item. One row per (schedule_item_id,
// provider) — re-running with the same provider updates the row.
//
// Statuses:
//   - "pending"      — screening initiated, awaiting result
//   - "clear"        — provider returned no hits
//   - "hit"          — provider returned at least one hit; line cannot be
//                      authorised until "manual_clear" is recorded with a
//                      reason in the audit
//   - "manual_clear" — finance has reviewed a hit and manually cleared it
//   - "skipped"      — screening was bypassed (e.g. provider unavailable);
//                      blocks authorisation the same way as "hit"
//
// Provider names mirror services/sanctions/provider.go::Name() values; "manual"
// is the bootstrap default that ships with Phase 3.
// Table name: sanctions_screenings
type SanctionsScreening struct {
	ID             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID     int        `json:"schedule_id" gorm:"index"`
	ScheduleItemID int        `json:"schedule_item_id" gorm:"index;uniqueIndex:idx_sanctions_item_provider"`
	ClaimID        int        `json:"claim_id" gorm:"index"`
	Provider       string     `json:"provider" gorm:"size:64;uniqueIndex:idx_sanctions_item_provider"`
	Status         string     `json:"status" gorm:"size:32"`
	ProviderRef    string     `json:"provider_ref" gorm:"size:128"`
	HitSummary     string     `json:"hit_summary"`
	Notes          string     `json:"notes"`
	ScreenedBy     string     `json:"screened_by"`
	ScreenedAt     *time.Time `json:"screened_at"`
	ClearedBy      string     `json:"cleared_by"`
	ClearedAt      *time.Time `json:"cleared_at"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

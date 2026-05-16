package models

import "time"

// PaymentCutoffConfig stores the tenant's payment-run cut-off schedule and the
// daily-payment safety limit. CutoffTimes is a CSV of HH:MM values in the
// configured timezone (e.g. "11:00,15:00"); the scheduler ticks every minute
// and generates a payment schedule when "now" matches a configured time and
// no PaymentCutoffRun row exists for that license + scheduled_at yet.
//
// LicenseId is included so this can later go truly multi-tenant; for now,
// the singleton install row uses an empty license. Looked up most-recently-
// updated row first.
//
// DailyPaymentLimit (ZAR) is checked at first finance authorisation. A value
// of 0 means "no limit".
// Table name: payment_cutoff_configs
type PaymentCutoffConfig struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LicenseId         string    `json:"license_id" gorm:"size:191;uniqueIndex"`
	Enabled           bool      `json:"enabled" gorm:"default:true"`
	CutoffTimes       string    `json:"cutoff_times" gorm:"size:255"`
	DailyPaymentLimit float64   `json:"daily_payment_limit"`
	Timezone          string    `json:"timezone" gorm:"size:64"`
	UpdatedBy         string    `json:"updated_by"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// PaymentCutoffRun records every auto- or manual-triggered cut-off execution.
// The unique index on (license_id, scheduled_at, trigger_type) makes the
// scheduler idempotent: if a row already exists for today's 11:00 auto cut-off,
// the loop simply skips it on the next minute-tick.
//
// Status:
//   - "ok"          — schedule generated
//   - "no_claims"   — no approved claims found at the cut-off, no schedule made
//   - "error"       — generation failed; ErrorMessage carries the reason
// Table name: payment_cutoff_runs
type PaymentCutoffRun struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LicenseId     string    `json:"license_id" gorm:"size:191;index:idx_cutoff_run_dedup,unique,priority:1"`
	ScheduledAt   time.Time `json:"scheduled_at" gorm:"index:idx_cutoff_run_dedup,unique,priority:2"`
	TriggerType   string    `json:"trigger_type" gorm:"size:16;index:idx_cutoff_run_dedup,unique,priority:3"`
	Status        string    `json:"status" gorm:"size:16"`
	ErrorMessage  string    `json:"error_message"`
	ScheduleID    *int      `json:"schedule_id"`
	ClaimsCount   int       `json:"claims_count"`
	TotalAmount   float64   `json:"total_amount"`
	TriggeredBy   string    `json:"triggered_by"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

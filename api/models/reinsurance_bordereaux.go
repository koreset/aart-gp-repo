package models

import "time"

// RIBordereauxRun represents a single reinsurance bordereaux generation run
type RIBordereauxRun struct {
	ID            int    `json:"id" gorm:"primaryKey;autoIncrement"`
	RunID         string `json:"run_id" gorm:"uniqueIndex;size:100"`
	TreatyID      int    `json:"treaty_id" gorm:"index"`
	TreatyNumber  string `json:"treaty_number" gorm:"default:''"`
	ReinsurerName string `json:"reinsurer_name" gorm:"default:''"`
	PeriodStart   string `json:"period_start" gorm:"default:''"`
	PeriodEnd     string `json:"period_end" gorm:"default:''"`
	PeriodLabel   string `json:"period_label" gorm:"default:''"`
	// Type: member_census | claims_run
	Type      string `json:"type" gorm:"default:'member_census'"`
	SchemeIDs string `json:"scheme_ids" gorm:"type:text"`
	// Status: draft | generated | submitted | acknowledged | queried | settled
	Status              string     `json:"status" gorm:"default:'draft'"`
	TotalLives          int        `json:"total_lives" gorm:"default:0"`
	TotalCededLives     int        `json:"total_ceded_lives" gorm:"default:0"`
	GrossPremium        float64    `json:"gross_premium" gorm:"default:0"`
	CededPremium        float64    `json:"ceded_premium" gorm:"default:0"`
	RetainedPremium     float64    `json:"retained_premium" gorm:"default:0"`
	GrossClaimsIncurred float64    `json:"gross_claims_incurred" gorm:"default:0"`
	CededClaimsIncurred float64    `json:"ceded_claims_incurred" gorm:"default:0"`
	FileName            string     `json:"file_name" gorm:"default:''"`
	FilePath            string     `json:"file_path" gorm:"default:''"`
	SubmittedAt         *time.Time `json:"submitted_at"`
	AcknowledgedAt      *time.Time `json:"acknowledged_at"`
	GeneratedBy         string     `json:"generated_by" gorm:"default:''"`
	SubmittedBy         string     `json:"submitted_by" gorm:"default:''"`
	// Amendment versioning: RunVersion starts at 1; amendments create a new run with ParentRunID set
	RunVersion     int    `json:"run_version" gorm:"default:1"`
	AmendmentNotes string `json:"amendment_notes" gorm:"type:text"`
	ParentRunID    *uint  `json:"parent_run_id"`
	// Submission register fields (Phase 3)
	// BPR: Bordereaux Processing Reference, auto-generated on creation (BPR-YYYYMM-NNN)
	BPR            string    `json:"bpr" gorm:"uniqueIndex;size:50;default:''"`
	ReceivedDate   string    `json:"received_date" gorm:"default:''"`
	AcknowledgedBy string    `json:"acknowledged_by" gorm:"default:''"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RIBordereauxMemberRow stores per-member cession data for a run
type RIBordereauxMemberRow struct {
	ID              int     `json:"id" gorm:"primaryKey;autoIncrement"`
	RunID           string  `json:"run_id" gorm:"index;size:100"`
	SchemeID        int     `json:"scheme_id" gorm:"index"`
	SchemeName      string  `json:"scheme_name" gorm:"default:''"`
	MemberIDNumber  string  `json:"member_id_number" gorm:"default:''"`
	MemberName      string  `json:"member_name" gorm:"default:''"`
	DateOfBirth     string  `json:"date_of_birth" gorm:"default:''"`
	Age             int     `json:"age" gorm:"default:0"`
	Gender          string  `json:"gender" gorm:"default:''"`
	EntryDate       string  `json:"entry_date" gorm:"default:''"`
	BenefitCode     string  `json:"benefit_code" gorm:"default:''"`
	BenefitName     string  `json:"benefit_name" gorm:"default:''"`
	SumAssured      float64 `json:"sum_assured" gorm:"default:0"`
	AnnualSalary    float64 `json:"annual_salary" gorm:"default:0"`
	GrossPremium    float64 `json:"gross_premium" gorm:"default:0"`
	CededPremium    float64 `json:"ceded_premium" gorm:"default:0"`
	RetainedPremium float64 `json:"retained_premium" gorm:"default:0"`
	RetentionAmount float64 `json:"retention_amount" gorm:"default:0"`
	CededAmount     float64 `json:"ceded_amount" gorm:"default:0"`
	// CessionBasis: quota_share | surplus | xl_excess
	CessionBasis string `json:"cession_basis" gorm:"default:''"`
	// ChangeType: new | amendment | deletion | nil
	ChangeType string `json:"change_type" gorm:"default:'nil'"`
	// MemberStatus: in_force | new_entrant | exit | amended
	MemberStatus string `json:"member_status" gorm:"default:'in_force'"`
	// Phase 1 additions — mandatory §3.1 fields
	EndorsementNumber string `json:"endorsement_number" gorm:"default:''"`
	// SanctionsScreeningStatus: cleared | pending | flagged
	SanctionsScreeningStatus string  `json:"sanctions_screening_status" gorm:"default:'cleared'"`
	CumulativePremiumYTD     float64 `json:"cumulative_premium_ytd" gorm:"default:0"`
	TreatySection            string  `json:"treaty_section" gorm:"default:''"`
	// PolicyType: individual | group
	PolicyType       string    `json:"policy_type" gorm:"default:'group'"`
	ExchangeRate     float64   `json:"exchange_rate" gorm:"default:1"`
	PeriodAdjustment float64   `json:"period_adjustment" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at"`
}

// RIBordereauxClaimsRow stores per-claim cession data for a claims run
type RIBordereauxClaimsRow struct {
	ID               int     `json:"id" gorm:"primaryKey;autoIncrement"`
	RunID            string  `json:"run_id" gorm:"index;size:100"`
	SchemeID         int     `json:"scheme_id" gorm:"index"`
	SchemeName       string  `json:"scheme_name" gorm:"default:''"`
	ClaimNumber      string  `json:"claim_number" gorm:"default:''"`
	MemberIDNumber   string  `json:"member_id_number" gorm:"default:''"`
	MemberName       string  `json:"member_name" gorm:"default:''"`
	DateOfEvent      string  `json:"date_of_event" gorm:"default:''"`
	DateNotified     string  `json:"date_notified" gorm:"default:''"`
	BenefitCode      string  `json:"benefit_code" gorm:"default:''"`
	GrossClaimAmount float64 `json:"gross_claim_amount" gorm:"default:0"`
	ExcessRetention  float64 `json:"excess_retention" gorm:"default:0"`
	CededClaimAmount float64 `json:"ceded_claim_amount" gorm:"default:0"`
	RecoveryReceived float64 `json:"recovery_received" gorm:"default:0"`
	ClaimStatus      string  `json:"claim_status" gorm:"default:''"`
	IsBelowRetention bool    `json:"is_below_retention" gorm:"default:false"`
	IsIBNR           bool    `json:"is_ibnr" gorm:"default:false"`
	// Phase 1 additions — mandatory §3.2 fields
	ReinsurerClaimReference string    `json:"reinsurer_claim_reference" gorm:"default:''"`
	CauseOfLoss             string    `json:"cause_of_loss" gorm:"default:''"`
	GrossPaidLosses         float64   `json:"gross_paid_losses" gorm:"default:0"`
	GrossOutstandingReserve float64   `json:"gross_outstanding_reserve" gorm:"default:0"`
	Recoveries              float64   `json:"recoveries" gorm:"default:0"`
	PriorPeriodMovement     float64   `json:"prior_period_movement" gorm:"default:0"`
	CumulativeLossYTD       float64   `json:"cumulative_loss_ytd" gorm:"default:0"`
	IBNRFlag                bool      `json:"ibnr_flag" gorm:"default:false"`
	LargeLossFlag           bool      `json:"large_loss_flag" gorm:"default:false"`
	CatastropheEventCode    string    `json:"catastrophe_event_code" gorm:"default:''"`
	CreatedAt               time.Time `json:"created_at"`
}

// LargeClaimNotice tracks individual large claim notifications to reinsurers
type LargeClaimNotice struct {
	ID                   int     `json:"id" gorm:"primaryKey;autoIncrement"`
	TreatyID             int     `json:"treaty_id" gorm:"index"`
	TreatyNumber         string  `json:"treaty_number" gorm:"default:''"`
	ClaimID              int     `json:"claim_id" gorm:"index"`
	ClaimNumber          string  `json:"claim_number" gorm:"default:''"`
	SchemeID             int     `json:"scheme_id" gorm:"index"`
	SchemeName           string  `json:"scheme_name" gorm:"default:''"`
	ReinsurerName        string  `json:"reinsurer_name" gorm:"default:''"`
	EventDate            string  `json:"event_date" gorm:"default:''"`
	NotifiedDate         string  `json:"notified_date" gorm:"default:''"`
	DueDate              string  `json:"due_date" gorm:"default:''"`
	BenefitCode          string  `json:"benefit_code" gorm:"default:''"`
	GrossClaimAmount     float64 `json:"gross_claim_amount" gorm:"default:0"`
	ExcessAmount         float64 `json:"excess_amount" gorm:"default:0"`
	EstimatedCededAmount float64 `json:"estimated_ceded_amount" gorm:"default:0"`
	// Status: pending | sent | acknowledged | late | queried
	Status         string     `json:"status" gorm:"default:'pending'"`
	LateFlag       bool       `json:"late_flag" gorm:"default:false"`
	SentAt         *time.Time `json:"sent_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	QueryDetails   string     `json:"query_details" gorm:"type:text"`
	ResponseNotes  string     `json:"response_notes" gorm:"type:text"`
	CreatedBy      string     `json:"created_by" gorm:"default:''"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type GenerateRIBordereauxRequest struct {
	TreatyID  int    `json:"treaty_id" binding:"required"`
	Type      string `json:"type"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	SchemeIDs []int  `json:"scheme_ids"`
}

type SubmitRIBordereauxRequest struct {
	RunID   string `json:"run_id" binding:"required"`
	Message string `json:"message"`
}

type AcknowledgeReceiptRequest struct {
	ReceivedDate string `json:"received_date"`
	Notes        string `json:"notes"`
}

type AmendRIBordereauxRequest struct {
	AmendmentNotes string `json:"amendment_notes"`
}

type MonitorLargeClaimsRequest struct {
	TreatyID int `json:"treaty_id" binding:"required"`
}

type UpdateLargeClaimNoticeRequest struct {
	Status         string     `json:"status"`
	QueryDetails   string     `json:"query_details"`
	ResponseNotes  string     `json:"response_notes"`
	SentAt         *time.Time `json:"sent_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
}

// RIValidationResult stores a single validation finding for a bordereaux run
type RIValidationResult struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	RunID string `json:"run_id" gorm:"index;size:100"`
	// RowType: member | claims | run  (run = run-level / header checks)
	RowType  string `json:"row_type" gorm:"default:''"`
	RowIndex int    `json:"row_index" gorm:"default:0"`
	// Level: 1 = structural, 2 = data_integrity, 3 = business_rules
	Level int `json:"level" gorm:"default:1"`
	// Severity: critical | major | minor
	Severity  string    `json:"severity" gorm:"default:''"`
	FieldName string    `json:"field_name" gorm:"default:''"`
	ErrorCode string    `json:"error_code" gorm:"default:''"`
	Message   string    `json:"message" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

// ValidationSummary is returned by the validate endpoint
type ValidationSummary struct {
	RunID    string               `json:"run_id"`
	Status   string               `json:"status"`
	Critical int                  `json:"critical"`
	Major    int                  `json:"major"`
	Minor    int                  `json:"minor"`
	Total    int                  `json:"total"`
	Results  []RIValidationResult `json:"results"`
}

// RIBordereauxKPIs holds all 7 §8.2 KPIs computed for a treaty+period window
type RIBordereauxKPIs struct {
	// KPI 1 — Submission timeliness (target ≥90%): % of runs submitted within 10 days of period end
	SubmissionTimelinessPct    float64 `json:"submission_timeliness_pct"`
	SubmissionTimelinessTarget float64 `json:"submission_timeliness_target"`
	SubmissionTotal            int     `json:"submission_total"`
	SubmissionOnTime           int     `json:"submission_on_time"`

	// KPI 2 — Processing timeliness (target ≥95%): % acknowledged within 10 days of submission
	ProcessingTimelinessPct    float64 `json:"processing_timeliness_pct"`
	ProcessingTimelinessTarget float64 `json:"processing_timeliness_target"`
	ProcessingTotal            int     `json:"processing_total"`
	ProcessingOnTime           int     `json:"processing_on_time"`

	// KPI 3 — First-time acceptance rate (target ≥85%): % of runs never requiring an amendment
	FirstTimeAcceptancePct    float64 `json:"first_time_acceptance_pct"`
	FirstTimeAcceptanceTarget float64 `json:"first_time_acceptance_target"`
	FirstTimeTotal            int     `json:"first_time_total"`
	FirstTimeAccepted         int     `json:"first_time_accepted"`

	// KPI 4 — Avg error resolution days (target ≤5): avg days from run creation to validated status
	AvgErrorResolutionDays   float64 `json:"avg_error_resolution_days"`
	AvgErrorResolutionTarget float64 `json:"avg_error_resolution_target"`
	ErrorResolutionSamples   int     `json:"error_resolution_samples"`

	// KPI 5 — Settlement timeliness (target ≥98%): % of technical accounts settled within 30 days of agreement
	SettlementTimelinessPct    float64 `json:"settlement_timeliness_pct"`
	SettlementTimelinessTarget float64 `json:"settlement_timeliness_target"`
	SettlementTotal            int     `json:"settlement_total"`
	SettlementOnTime           int     `json:"settlement_on_time"`

	// KPI 6 — Claims register completeness (target 100%): % of large loss claims with a LargeClaimNotice
	ClaimsCompletenessPct    float64 `json:"claims_completeness_pct"`
	ClaimsCompletenessTarget float64 `json:"claims_completeness_target"`
	ClaimsLargeLossTotal     int     `json:"claims_large_loss_total"`
	ClaimsWithNotice         int     `json:"claims_with_notice"`

	// KPI 7 — Open query backlog (target <5 per quarter): runs stuck in validation_failed >30 days
	OpenQueryBacklog       int `json:"open_query_backlog"`
	OpenQueryBacklogTarget int `json:"open_query_backlog_target"`

	// Filters applied
	TreatyID   int    `json:"treaty_id"`
	PeriodFrom string `json:"period_from"`
	PeriodTo   string `json:"period_to"`
	ComputedAt string `json:"computed_at"`
}

type RIBordereauxStats struct {
	TotalRuns         int     `json:"total_runs"`
	MemberCensus      int     `json:"member_census"`
	ClaimsRuns        int     `json:"claims_runs"`
	Submitted         int     `json:"submitted"`
	TotalCededPremium float64 `json:"total_ceded_premium"`
}

type LargeClaimStats struct {
	Total        int `json:"total"`
	Pending      int `json:"pending"`
	Sent         int `json:"sent"`
	Acknowledged int `json:"acknowledged"`
	Late         int `json:"late"`
	Queried      int `json:"queried"`
}

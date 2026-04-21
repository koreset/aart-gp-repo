package models

import "time"

// BordereauxFieldMapping represents a mapping from a source field to a target field in the template
type BordereauxFieldMapping struct {
	SourceField string `json:"source_field"`
	TargetField string `json:"target_field"`
	Required    bool   `json:"required"`
}

// BordereauxValidationRules captures the validation settings for a template
type BordereauxValidationRules struct {
	ValidateIDNumbers      bool `json:"validate_id_numbers"`
	ValidateBankingDetails bool `json:"validate_banking_details"`
	ValidateAmounts        bool `json:"validate_amounts"`
	ExcludeInvalid         bool `json:"exclude_invalid"`
	RequireBeneficiaries   bool `json:"require_beneficiaries"`
	ValidateDates          bool `json:"validate_dates"`
}

// BordereauxTemplate stores the configuration for generating or parsing bordereaux files
type BordereauxTemplate struct {
	ID              int                       `json:"id" gorm:"primaryKey"`
	Name            string                    `json:"name"`
	InsurerID       int                       `json:"insurer_id"`
	Type            string                    `json:"type"`   // e.g. "member"
	Status          string                    `json:"status"` // e.g. "draft"
	Format          string                    `json:"format"` // e.g. "excel"
	Description     string                    `json:"description"`
	FieldMappings   []BordereauxFieldMapping  `json:"field_mappings" gorm:"serializer:json"`
	ValidationRules BordereauxValidationRules `json:"validation_rules" gorm:"serializer:json"`
	UsageCount      int                       `json:"usage_count"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// GeneratedBordereaux tracks generated bordereaux files
type GeneratedBordereaux struct {
	ID             int                  `json:"id" gorm:"primaryKey"`
	GeneratedID    string               `json:"generated_id" gorm:"uniqueIndex;size:255"` // e.g. BRD-YYYYMMDD-XXXXXX
	TemplateID     int                  `json:"template_id"`
	Type           string               `json:"type"` // member, premium, claim
	FileName       string               `json:"file_name"`
	FilePath       string               `json:"file_path"`
	FileSize       int64                `json:"file_size"`
	Records        int                  `json:"records"`
	ProcessingTime string               `json:"processing_time"`
	SchemeName     string               `json:"scheme_name"`
	InsurerName    string               `json:"insurer_name"`
	Status         string               `json:"status"`
	Progress       int                  `json:"progress"`
	SubmissionDate *time.Time           `json:"submission_date"`
	LastUpdated    time.Time            `json:"last_updated"`
	SLAStatus      string               `json:"sla_status"`
	PeriodStart    time.Time            `json:"period_start"`
	PeriodEnd      time.Time            `json:"period_end"`
	Timeline       []BordereauxTimeline `json:"timeline" gorm:"foreignKey:GeneratedBordereauxID;references:GeneratedID"`
	ReviewedBy     string               `json:"reviewed_by" gorm:"default:''"`
	ReviewedAt     *time.Time           `json:"reviewed_at"`
	ApprovedBy     string               `json:"approved_by" gorm:"default:''"`
	ApprovedAt     *time.Time           `json:"approved_at"`
	ReturnReason   string               `json:"return_reason" gorm:"default:''"`
	CreatedBy      string               `json:"created_by"`
	CreatedAt      time.Time            `json:"created_at"`
	// RequestPayload stores the original GenerateBordereauxRequest so the row
	// can be regenerated in-place while keeping the same generated_id, timeline
	// and audit trail. Null for rows created before regenerate was added.
	RequestPayload JSON `json:"request_payload,omitempty" gorm:"type:json"`
}

// BordereauxTimeline tracks events for a generated bordereaux
type BordereauxTimeline struct {
	ID                    int       `json:"id" gorm:"primaryKey"`
	GeneratedBordereauxID string    `json:"generated_bordereaux_id" gorm:"index;size:255"`
	Date                  time.Time `json:"date"`
	Type                  string    `json:"type"` // generated, submitted, confirmed, etc.
	Title                 string    `json:"title"`
	Description           string    `json:"description"`
}

// BordereauxConfiguration represents a saved bordereaux configuration
type BordereauxConfiguration struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ConfigData  JSON       `json:"config_data" gorm:"type:json"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// BordereauxConfirmation represents the confirmation file data received from a scheme
type BordereauxConfirmation struct {
	ID                    int        `json:"id" gorm:"primaryKey"`
	GeneratedBordereauxID string     `json:"generated_bordereaux_id" gorm:"index;size:255"`
	SchemeID              int        `json:"scheme_id"`
	SchemeName            string     `json:"scheme_name"`
	FileName              string     `json:"file_name"`
	FilePath              string     `json:"file_path"`
	FileType              string     `json:"file_type"` // premium, claim, etc.
	ValuationMonth        int        `json:"valuation_month"`
	ValuationYear         int        `json:"valuation_year"`
	Status                string     `json:"status"` // pending, processed, error
	ImportedBy            string     `json:"imported_by"`
	ImportedAt            time.Time  `json:"imported_at"`
	MatchedCount          int        `json:"matched_count"`
	DiscrepancyCount      int        `json:"discrepancy_count"`
	SubmittedAmount       float64    `json:"submitted_amount"`
	ConfirmedAmount       float64    `json:"confirmed_amount"`
	Variance              float64    `json:"variance"`
	MatchScore            float64    `json:"match_score"`
	LastReconciled        *time.Time `json:"last_reconciled"`
}

// BordereauxConfirmationRecord persists the content of a confirmation file
type BordereauxConfirmationRecord struct {
	ID                       int       `json:"id" gorm:"primaryKey"`
	BordereauxConfirmationID int       `json:"bordereaux_confirmation_id" gorm:"index"`
	BordereauxID             string    `json:"bordereaux_id" gorm:"index;size:255"`
	MemberID                 string    `json:"member_id" gorm:"index"`
	MemberName               string    `json:"member_name"`
	EmployeeNumber           string    `json:"employee_number" gorm:"index"`
	IDNumber                 string    `json:"id_number" gorm:"index"`
	Amount                   float64   `json:"amount"`
	PremiumAmount            float64   `json:"premium_amount"`
	FuneralPremiumAmount     float64   `json:"funeral_premium_amount"`
	ValuationMonth           int       `json:"valuation_month"`
	ValuationYear            int       `json:"valuation_year"`
	SchemeName               string    `json:"scheme_name"`
	RawData                  JSON      `json:"raw_data" gorm:"type:json"`
	CreatedAt                time.Time `json:"created_at"`
}

// BordereauxReconciliationResult stores the outcome of matching a confirmation record with a bordereaux record
type BordereauxReconciliationResult struct {
	ID                       int       `json:"id" gorm:"primaryKey"`
	BordereauxConfirmationID int       `json:"bordereaux_confirmation_id" gorm:"index"`
	GeneratedBordereauxID    string    `json:"generated_bordereaux_id" gorm:"index;size:255"`
	RecordID                 string    `json:"record_id" gorm:"index"` // e.g. Member ID or unique ref
	MemberName               string    `json:"member_name"`
	Field                    string    `json:"field"` // the field that has a discrepancy
	ExpectedValue            string    `json:"expected_value"`
	ActualValue              string    `json:"actual_value"`
	Variance                 float64   `json:"variance"`
	Status                   string    `json:"status" gorm:"size:64"` // matched, discrepancy, missing, extra, escalated, resolved
	IsResolved               bool      `json:"is_resolved"`
	Comments                 string    `json:"comments"`
	// Escalation workflow fields — populated by EscalateDiscrepancy.
	AssignedTo  string     `json:"assigned_to" gorm:"default:''"`
	Priority    string     `json:"priority" gorm:"default:''"` // low | medium | high
	EscalatedBy string     `json:"escalated_by" gorm:"default:''"`
	EscalatedAt *time.Time `json:"escalated_at"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
}

// BordereauxConfirmationNote is a free-text annotation on a confirmation.
// Previously these were stored as synthetic BordereauxReconciliationResult rows
// with Field="_note"; this dedicated table removes the need to filter "_note"
// rows out of every consumer of the results table.
type BordereauxConfirmationNote struct {
	ID                       int       `json:"id" gorm:"primaryKey"`
	BordereauxConfirmationID int       `json:"bordereaux_confirmation_id" gorm:"index"`
	Note                     string    `json:"note" gorm:"type:text"`
	CreatedBy                string    `json:"created_by" gorm:"default:''"`
	CreatedAt                time.Time `json:"created_at"`
}

// BordereauxDashboardStats represents the statistics for the dashboard
type BordereauxDashboardStats struct {
	GeneratedThisMonth int `json:"generated_this_month"`
	PendingSubmissions int `json:"pending_submissions"`
	ReconciledThisWeek int `json:"reconciled_this_week"`
	ActiveTemplates    int `json:"active_templates"`
}

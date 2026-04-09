package models

import "time"

// PremiumResponse is the standard API response envelope for premium module endpoints.
type PremiumResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse is the standard paginated API response envelope.
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ---------------------------------------------------------------------------
// Contribution Config
// ---------------------------------------------------------------------------

// ContributionConfig defines how premiums are split between employer and employee
// for a given scheme.
type ContributionConfig struct {
	ID               int                           `json:"id" gorm:"primaryKey;autoIncrement"`
	SchemeID         int                           `json:"scheme_id" binding:"required"`
	ContributionType string                        `json:"contribution_type" binding:"required"` // "employer_only" | "split"
	EmployerPercent  float64                       `json:"employer_percent"`
	EmployeePercent  float64                       `json:"employee_percent"`
	EffectiveDate    string                        `json:"effective_date" binding:"required"`
	BenefitOverrides []BenefitContributionOverride `json:"benefit_overrides,omitempty" gorm:"-"`
	UpdatedBy        string                        `json:"updated_by"`
	UpdatedAt        time.Time                     `json:"updated_at" gorm:"autoUpdateTime"`
}

// BenefitContributionOverride allows per-benefit contribution splits that
// differ from the scheme-level default.
type BenefitContributionOverride struct {
	BenefitCode     string  `json:"benefit_code"`
	EmployerPercent float64 `json:"employer_percent"`
	EmployeePercent float64 `json:"employee_percent"`
}

// ---------------------------------------------------------------------------
// Premium Schedule
// ---------------------------------------------------------------------------

// PremiumSchedule is a summary row for a monthly premium schedule.
type PremiumSchedule struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	SchemeID           int        `json:"scheme_id"`
	SchemeName         string     `json:"scheme_name"`
	Month              int        `json:"month"`
	Year               int        `json:"year"`
	MemberCount        int        `json:"member_count"`
	GrossPremium       float64    `json:"gross_premium"`
	NetPayable         float64    `json:"net_payable"`
	Status             string     `json:"status"` // "draft" | "reviewed" | "approved" | "finalized" | "invoiced" | "void" | "cancelled"
	GeneratedDate      string     `json:"generated_date"`
	GeneratedBy        string     `json:"generated_by"`
	ReviewedBy         string     `json:"reviewed_by" gorm:"default:''"`
	ReviewedAt         *time.Time `json:"reviewed_at"`
	ApprovedBy         string     `json:"approved_by" gorm:"default:''"`
	ApprovedAt         *time.Time `json:"approved_at"`
	FinalizedBy        string     `json:"finalized_by" gorm:"default:''"`
	FinalizedAt        *time.Time `json:"finalized_at"`
	VoidedBy           string     `json:"voided_by" gorm:"default:''"`
	VoidedAt           *time.Time `json:"voided_at"`
	VoidReason         string     `json:"void_reason" gorm:"default:''"`
	LinkedSubmissionID *int       `json:"linked_submission_id"`
	IsSupplementary    bool       `json:"is_supplementary" gorm:"default:false"`
	SupplementaryNote  string     `json:"supplementary_note" gorm:"default:''"`
}

// PremiumScheduleDetail extends PremiumSchedule with full member-level rows.
type PremiumScheduleDetail struct {
	PremiumSchedule
	NewJoiners int                 `json:"new_joiners"`
	Exits      int                 `json:"exits"`
	Members    []ScheduleMemberRow `json:"members"`
}

// ScheduleMemberRow is a single member line within a premium schedule.
type ScheduleMemberRow struct {
	ID                   int     `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID           int     `json:"schedule_id"`
	MemberID             int     `json:"member_id"`
	MemberName           string  `json:"member_name"`
	MemberIDNumber       string  `json:"member_id_number"`
	Benefit              string  `json:"benefit"`
	AnnualSalary         float64 `json:"annual_salary"`
	Rate                 float64 `json:"rate"`
	FullMonthPremium     float64 `json:"full_month_premium"`
	ProRataDays          *int    `json:"pro_rata_days,omitempty"`
	ActualPremium        float64 `json:"actual_premium"`
	EmployerContribution float64 `json:"employer_contribution"`
	EmployeeContribution float64 `json:"employee_contribution"`
	IsProRata            bool    `json:"is_pro_rata"`
}

// GenerateScheduleRequest is the request body for generating a new premium schedule.
type GenerateScheduleRequest struct {
	SchemeID int `json:"scheme_id" binding:"required"`
	Month    int `json:"month" binding:"required"`
	Year     int `json:"year" binding:"required"`
}

// BulkGenerateScheduleRequest is the request body for generating schedules for all in-force schemes.
type BulkGenerateScheduleRequest struct {
	Month int `json:"month" binding:"required"`
	Year  int `json:"year" binding:"required"`
}

// BulkScheduleResult is the per-scheme outcome of a bulk schedule generation.
type BulkScheduleResult struct {
	SchemeID   int    `json:"scheme_id"`
	SchemeName string `json:"scheme_name"`
	Status     string `json:"status"` // "success" | "skipped" | "failed"
	Message    string `json:"message,omitempty"`
	ScheduleID int    `json:"schedule_id,omitempty"`
}

// BulkGenerateResult is the aggregate outcome of generating schedules for all schemes.
type BulkGenerateResult struct {
	Total     int                  `json:"total"`
	Succeeded int                  `json:"succeeded"`
	Skipped   int                  `json:"skipped"`
	Failed    int                  `json:"failed"`
	Results   []BulkScheduleResult `json:"results"`
}

// ---------------------------------------------------------------------------
// Invoice
// ---------------------------------------------------------------------------

// Invoice is a summary row for a premium invoice.
type Invoice struct {
	ID                int     `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceNumber     string  `json:"invoice_number"`
	SchemeID          int     `json:"scheme_id"`
	SchemeName        string  `json:"scheme_name"`
	ScheduleID        int     `json:"schedule_id"`
	Month             int     `json:"month"`
	Year              int     `json:"year"`
	IssueDate         string  `json:"issue_date"`
	DueDate           string  `json:"due_date"`
	GrossAmount       float64 `json:"gross_amount"`
	Commission        float64 `json:"commission"`
	NetPayable        float64 `json:"net_payable"`
	PaidAmount        float64 `json:"paid_amount"`
	Balance           float64 `json:"balance"`
	Status            string  `json:"status"` // "draft" | "sent" | "paid" | "partial" | "overdue"
	IsRetroAdjustment bool    `json:"is_retro_adjustment" gorm:"default:false"`
}

// InvoiceDetail extends Invoice with line items, adjustments, payment history
// and contact information.
type InvoiceDetail struct {
	Invoice
	LineItems      []InvoiceLineItem   `json:"line_items"`
	Adjustments    []InvoiceAdjustment `json:"adjustments"`
	PaymentHistory []Payment           `json:"payment_history"`
	ContactName    string              `json:"contact_name"`
	ContactEmail   string              `json:"contact_email"`
}

// InvoiceLineItem represents a single benefit line on an invoice.
type InvoiceLineItem struct {
	Benefit     string  `json:"benefit"`
	MemberCount int     `json:"member_count"`
	BasePremium float64 `json:"base_premium"`
	Adjustment  float64 `json:"adjustment"`
	Total       float64 `json:"total"`
}

// InvoiceAdjustment represents a credit note or debit note applied to an invoice.
type InvoiceAdjustment struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceID   int     `json:"invoice_id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // "credit" | "debit"
	CreatedBy   string  `json:"created_by"`
	CreatedAt   string  `json:"created_at"`
}

// GenerateInvoiceRequest is the optional request body for generating an invoice.
type GenerateInvoiceRequest struct {
	DueDate string `json:"due_date"` // optional override, format "YYYY-MM-DD"
}

// EmailInvoiceRequest is the request body for emailing an invoice.
type EmailInvoiceRequest struct {
	RecipientEmail string `json:"recipient_email" binding:"required"`
	RecipientName  string `json:"recipient_name"`
	Message        string `json:"message"`
}

// ---------------------------------------------------------------------------
// Payment
// ---------------------------------------------------------------------------

// Payment represents a premium payment record.
type Payment struct {
	ID            int     `json:"id" gorm:"primaryKey;autoIncrement"`
	SchemeID      int     `json:"scheme_id"`
	SchemeName    string  `json:"scheme_name"`
	InvoiceID     *int    `json:"invoice_id,omitempty"`
	InvoiceNumber string  `json:"invoice_number"`
	PaymentDate   string  `json:"payment_date"`
	Method        string  `json:"method"` // "eft" | "cheque" | "cash" | "debit_order"
	BankReference string  `json:"bank_reference"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"` // "pending" | "matched" | "unmatched" | "voided"
	Notes         string  `json:"notes"`
	RecordedBy    string  `json:"recorded_by"`
	RecordedAt    string  `json:"recorded_at"`
}

// RecordPaymentRequest is the request body for recording a new payment.
type RecordPaymentRequest struct {
	SchemeID      int     `json:"scheme_id" binding:"required"`
	InvoiceID     *int    `json:"invoice_id,omitempty"`
	PaymentDate   string  `json:"payment_date" binding:"required"`
	Method        string  `json:"method" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	BankReference string  `json:"bank_reference"`
	Notes         string  `json:"notes"`
}

// BulkImportResult summarises the outcome of a bulk CSV payment import.
type BulkImportResult struct {
	TotalRows int             `json:"total_rows"`
	Matched   int             `json:"matched"`
	Unmatched int             `json:"unmatched"`
	Errors    int             `json:"errors"`
	Rows      []BulkImportRow `json:"rows"`
}

// BulkImportRow is the per-row result of a bulk payment import.
type BulkImportRow struct {
	RowNumber      int     `json:"row_number"`
	Reference      string  `json:"reference"`
	Amount         float64 `json:"amount"`
	MatchedInvoice string  `json:"matched_invoice"`
	Status         string  `json:"status"` // "matched" | "unmatched" | "error"
	ErrorReason    string  `json:"error_reason,omitempty"`
}

// ---------------------------------------------------------------------------
// Reconciliation
// ---------------------------------------------------------------------------

// ManualMatchRequest is the request body for manually linking a payment to an invoice.
type ManualMatchRequest struct {
	PaymentID int `json:"payment_id" binding:"required"`
	InvoiceID int `json:"invoice_id" binding:"required"`
}

// AutoMatchResult summarises the outcome of an automatic reconciliation run.
type AutoMatchResult struct {
	Matched   int `json:"matched"`
	Remaining int `json:"remaining"`
}

// AdjustmentNoteRequest is the request body for creating a credit or debit note.
// InvoiceID is optional for credit/debit notes raised from the reconciliation screen;
// the controller resolves the invoice from the scheme when InvoiceID is 0.
// Type is set by the controller (credit-note vs debit-note endpoint), not the client.
type AdjustmentNoteRequest struct {
	InvoiceID int     `json:"invoice_id"`
	SchemeID  int     `json:"scheme_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
	Reason    string  `json:"reason" binding:"required"`
	Type      string  `json:"type"` // "credit" | "debit" — set by controller
	Reference string  `json:"reference"`
}

// ---------------------------------------------------------------------------
// Arrears
// ---------------------------------------------------------------------------

// ArrearsRecord is a per-scheme arrears aging entry.
type ArrearsRecord struct {
	SchemeID          int     `json:"scheme_id"`
	SchemeName        string  `json:"scheme_name"`
	Days0To30         float64 `json:"days_0_to_30"`
	Days31To60        float64 `json:"days_31_to_60"`
	Days61To90        float64 `json:"days_61_to_90"`
	DaysOver90        float64 `json:"days_over_90"`
	TotalOutstanding  float64 `json:"total_outstanding"`
	GracePeriodExpiry *string `json:"grace_period_expiry,omitempty"`
	Status            string  `json:"status"` // "current" | "in_arrears" | "suspended"
}

// ArrearsHistory records a single event in a scheme's arrears history.
type ArrearsHistory struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SchemeID    int       `json:"scheme_id"`
	EventType   string    `json:"event_type"` // "REMINDER" | "PAYMENT_PLAN" | "SUSPENDED" | "REINSTATED"
	EventDate   string    `json:"event_date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	PerformedBy string    `json:"performed_by"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// SuspendCoverRequest is the request body for suspending cover for a scheme.
type SuspendCoverRequest struct {
	EffectiveDate string `json:"effective_date" binding:"required"`
	Reason        string `json:"reason" binding:"required"`
}

// ReinstateCoverRequest is the request body for reinstating cover for a scheme.
type ReinstateCoverRequest struct {
	ReinstatementDate string  `json:"reinstatement_date" binding:"required"`
	BackPremium       float64 `json:"back_premium"`
	Notes             string  `json:"notes"`
}

// PaymentPlan records an agreed repayment arrangement for arrears.
type PaymentPlan struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SchemeID  int       `json:"scheme_id"`
	Notes     string    `json:"notes"`
	CreatedBy string    `json:"created_by"`
	Status    string    `json:"status"` // "active" | "completed" | "defaulted"
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// PaymentPlanInstalment is a single instalment within a payment plan.
type PaymentPlanInstalment struct {
	ID     int        `json:"id" gorm:"primaryKey;autoIncrement"`
	PlanID int        `json:"plan_id"`
	Date   string     `json:"date" binding:"required"`
	Amount float64    `json:"amount" binding:"required"`
	Status string     `json:"status"` // "pending" | "paid" | "overdue"
	PaidAt *time.Time `json:"paid_at,omitempty"`
}

// PaymentPlanRequest is the request body for recording a payment plan.
type PaymentPlanRequest struct {
	Instalments []PaymentPlanInstalment `json:"instalments" binding:"required"`
	Notes       string                  `json:"notes"`
}

// ArrearsHistoryItem records a single event in a scheme's arrears history.
type ArrearsHistoryItem struct {
	EventType   string   `json:"event_type"`
	EventDate   string   `json:"event_date"`
	Description string   `json:"description"`
	PerformedBy string   `json:"performed_by"`
	Amount      *float64 `json:"amount,omitempty"`
}

// SendReminderRequest is the request body for sending a payment reminder.
type SendReminderRequest struct {
	Message string `json:"message"`
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// PremiumDashboardData is the full payload returned by the dashboard endpoint.
type PremiumDashboardData struct {
	KPIs            PremiumDashboardKPIs       `json:"kpis"`
	MonthlyTrend    []MonthlyPremiumTrend      `json:"monthly_trend"`
	StatusBreakdown []PremiumStatusBreakdown   `json:"status_breakdown"`
	TopOutstanding  []SchemeOutstandingSummary `json:"top_outstanding"`
}

// PremiumDashboardKPIs are the headline figures shown on the premium dashboard.
type PremiumDashboardKPIs struct {
	DueThisMonth       float64 `json:"due_this_month"`
	Collected          float64 `json:"collected"`
	CollectionRate     float64 `json:"collection_rate"`
	Outstanding        float64 `json:"outstanding"`
	Overdue            float64 `json:"overdue"`
	OverdueSchemeCount int     `json:"overdue_scheme_count"`
}

// MonthlyPremiumTrend is a single month data point for the trend chart.
type MonthlyPremiumTrend struct {
	Month     string  `json:"month"`
	Due       float64 `json:"due"`
	Collected float64 `json:"collected"`
}

// PremiumStatusBreakdown breaks down invoice amounts by status.
type PremiumStatusBreakdown struct {
	Status string  `json:"status"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

// SchemeOutstandingSummary is a per-scheme outstanding balance summary.
type SchemeOutstandingSummary struct {
	SchemeID         int     `json:"scheme_id"`
	SchemeName       string  `json:"scheme_name"`
	AmountDue        float64 `json:"amount_due"`
	AmountPaid       float64 `json:"amount_paid"`
	Balance          float64 `json:"balance"`
	Status           string  `json:"status"`
	DaysSincePayment *int    `json:"days_since_payment,omitempty"`
}

// ---------------------------------------------------------------------------
// Statements
// ---------------------------------------------------------------------------

// EmployerStatement is the detailed ledger statement sent to an employer.
type EmployerStatement struct {
	SchemeID       int                 `json:"scheme_id"`
	SchemeName     string              `json:"scheme_name"`
	Period         string              `json:"period"`
	OpeningBalance float64             `json:"opening_balance"`
	InvoicedAmount float64             `json:"invoiced_amount"`
	Received       float64             `json:"received"`
	Adjustments    float64             `json:"adjustments"`
	ClosingBalance float64             `json:"closing_balance"`
	LineItems      []StatementLineItem `json:"line_items"`
}

// BrokerCommissionStatement is a broker's commission statement for a period.
type BrokerCommissionStatement struct {
	BrokerID    int                      `json:"broker_id"`
	BrokerName  string                   `json:"broker_name"`
	Period      string                   `json:"period"`
	TotalEarned float64                  `json:"total_earned"`
	Schemes     []BrokerSchemeCommission `json:"schemes"`
}

// BrokerSchemeCommission is the commission breakdown for a single scheme within
// a broker statement.
type BrokerSchemeCommission struct {
	SchemeID         int     `json:"scheme_id"`
	SchemeName       string  `json:"scheme_name"`
	PremiumCollected float64 `json:"premium_collected"`
	CommissionRate   float64 `json:"commission_rate"`
	CommissionEarned float64 `json:"commission_earned"`
	Status           string  `json:"status"`
}

// StatementLineItem is a single ledger entry in an employer statement.
type StatementLineItem struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Balance     float64 `json:"balance"`
}

// EmailStatementRequest is the request body for emailing a statement.
type EmailStatementRequest struct {
	Type        string `json:"type" binding:"required"` // "employer" | "broker"
	RecipientID int    `json:"recipient_id" binding:"required"`
	FromDate    string `json:"from_date" binding:"required"`
	ToDate      string `json:"to_date" binding:"required"`
}

// ---------------------------------------------------------------------------
// Schedule Lifecycle
// ---------------------------------------------------------------------------

// VoidScheduleRequest is the request body for voiding or cancelling a schedule.
type VoidScheduleRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ScheduleMemberRateRequest is the request body for updating a member row's rate.
type ScheduleMemberRateRequest struct {
	Rate            float64 `json:"rate"`
	ActualPremium   float64 `json:"actual_premium"`
	EmployerContrib float64 `json:"employer_contribution"`
	EmployeeContrib float64 `json:"employee_contribution"`
}

// AddScheduleMemberRequest is the request body for adding a member to a draft schedule.
type AddScheduleMemberRequest struct {
	MemberID int `json:"member_id" binding:"required"`
}

// ---------------------------------------------------------------------------
// Contribution summary helpers
// ---------------------------------------------------------------------------

// SchemePaymentSummary provides the latest payment position for a scheme.
type SchemePaymentSummary struct {
	SchemeID           int      `json:"scheme_id"`
	LastPaymentDate    *string  `json:"last_payment_date,omitempty"`
	LastPaymentAmount  *float64 `json:"last_payment_amount,omitempty"`
	OutstandingBalance float64  `json:"outstanding_balance"`
	Status             string   `json:"status"`
}

// MemberContribution is the contribution breakdown for a single member in a month.
type MemberContribution struct {
	MemberID             int     `json:"member_id"`
	EmployerContribution float64 `json:"employer_contribution"`
	EmployeeContribution float64 `json:"employee_contribution"`
	TotalContribution    float64 `json:"total_contribution"`
	EffectiveMonth       string  `json:"effective_month"`
}

// CollectionRateResponse provides yearly collection rate data with monthly breakdown.
type CollectionRateResponse struct {
	Year           int                     `json:"year"`
	CollectionRate float64                 `json:"collection_rate"`
	MonthlyRates   []MonthlyCollectionRate `json:"monthly_rates"`
}

// MonthlyCollectionRate is a single month's collection rate.
type MonthlyCollectionRate struct {
	Month string  `json:"month"`
	Rate  float64 `json:"rate"`
}

// ---------------------------------------------------------------------------
// Schedule Coverage Matrix
// ---------------------------------------------------------------------------

// ScheduleCoverageMatrix is the top-level response for the coverage overview.
type ScheduleCoverageMatrix struct {
	Months  []string                   `json:"months"`  // e.g. ["Mar 2026","Feb 2026",...]
	Schemes []SchemeCoverageRow        `json:"schemes"`
}

// SchemeCoverageRow is one row in the coverage matrix (one scheme).
type SchemeCoverageRow struct {
	SchemeID         int                      `json:"scheme_id"`
	SchemeName       string                   `json:"scheme_name"`
	CommencementDate string                   `json:"commencement_date"`
	Cells            []ScheduleCoverageCell    `json:"cells"`
}

// ScheduleCoverageCell represents one month-cell for a scheme.
type ScheduleCoverageCell struct {
	Month              int    `json:"month"`
	Year               int    `json:"year"`
	Exists             bool   `json:"exists"`
	ScheduleID         int    `json:"schedule_id,omitempty"`
	Status             string `json:"status,omitempty"`
	BeforeCommencement bool   `json:"before_commencement"`
}

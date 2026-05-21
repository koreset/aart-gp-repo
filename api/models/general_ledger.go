package models

import "time"

// Operational General Ledger.
//
// This file owns the persistent double-entry posting layer for the operational
// side of the business: claim payments confirmed by the bank, premium receipts
// allocated to invoices, write-offs, refunds, and manually posted journals.
//
// It is intentionally independent of the IFRS17 / CSM reporting stack — the
// CSM-driven JournalTransactions / TrialBalanceReportEntry / ChartAccountItem
// structures in csm_engine.go and ifrs17.go are read-only valuation outputs
// and share nothing with the tables defined here (no chart, no journal, no FK).

// ---------------------------------------------------------------------------
// Chart of Accounts
// ---------------------------------------------------------------------------

// GLAccount is one node in the operational chart of accounts. Hierarchies are
// expressed via ParentID; the same account may not be active in two periods
// simultaneously (we soft-deactivate by setting IsActive=false rather than
// deleting, so historical postings keep their reference).
//
// Master-data mutations follow a maker/checker pattern: a request from the
// admin UI writes the proposed change into PendingChangeJSON and flips
// ApprovalStatus to "pending_*". A second user (different from
// PendingRequestedBy) calls the approve-change endpoint, which applies the
// change to the live row.
type GLAccount struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Code               string     `json:"code" gorm:"type:varchar(50);uniqueIndex"`
	Name               string     `json:"name" gorm:"type:varchar(191)"`
	AccountType        string     `json:"account_type" gorm:"size:32;index"`   // asset | liability | equity | income | expense
	NormalBalance      string     `json:"normal_balance" gorm:"size:8"`        // debit | credit
	ParentID           *int       `json:"parent_id" gorm:"index"`
	IsActive           bool       `json:"is_active" gorm:"default:true"`
	Description        string     `json:"description"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy          string     `json:"created_by" gorm:"size:191"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy          string     `json:"updated_by" gorm:"size:191"`

	// Maker/checker workflow fields.
	ApprovalStatus     string     `json:"approval_status" gorm:"size:24;default:'active';index"` // active | pending_create | pending_update | pending_deactivate
	PendingChangeJSON  string     `json:"pending_change_json,omitempty" gorm:"type:text"`
	PendingRequestedBy string     `json:"pending_requested_by,omitempty" gorm:"size:191"`
	PendingRequestedAt *time.Time `json:"pending_requested_at,omitempty"`
}

func (GLAccount) TableName() string { return "gl_accounts" }

// ---------------------------------------------------------------------------
// Accounting Periods
// ---------------------------------------------------------------------------

// AccountingPeriod gates posting. Postings can only land in a period whose
// status is "open". Period names follow `YYYY-MM` for the monthly cadence.
//
// Closing is a two-step dual-control flow: any user with the
// gl:request_close_period permission requests the close (open → close_requested),
// and a different user with gl:close_period commits it (close_requested →
// closed). Self-approval is rejected at the service layer.
type AccountingPeriod struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name               string     `json:"name" gorm:"type:varchar(20);uniqueIndex"` // "2026-05"
	StartDate          time.Time  `json:"start_date" gorm:"index"`
	EndDate            time.Time  `json:"end_date"`
	Status             string     `json:"status" gorm:"size:24;default:'open';index"` // open | close_requested | closed
	CloseRequestedBy   string     `json:"close_requested_by,omitempty" gorm:"size:191"`
	CloseRequestedAt   *time.Time `json:"close_requested_at,omitempty"`
	ClosedAt           *time.Time `json:"closed_at"`
	ClosedBy           string     `json:"closed_by" gorm:"size:191"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AccountingPeriod) TableName() string { return "accounting_periods" }

// ---------------------------------------------------------------------------
// Journal Entry + Lines (the immutable posted ledger)
// ---------------------------------------------------------------------------

// JournalEntry is the header of one balanced (DR=CR) accounting event.
// Posted entries are never edited or physically deleted; corrections are made
// by posting a reversal (which sets IsReversed=true on this row and creates a
// new entry with source_type="reversal").
//
// Manually-authored entries flow through a maker/checker state machine:
//
//   draft → submitted → approved → posted
//                                 ↘ reversal_pending → reversal_approved → reversed
//
// System-sourced entries (claim_payment, premium_allocation, write_off, refund)
// are authorised upstream and skip the workflow — they're created with
// Status="posted" by GLPost() directly.
//
// The approver (ApprovedBy) and the submitter (SubmittedBy) must be different
// users; the same rule applies to the reversal workflow.
type JournalEntry struct {
	ID                int           `json:"id" gorm:"primaryKey;autoIncrement"`
	EntryNumber       string        `json:"entry_number" gorm:"type:varchar(32);uniqueIndex"` // JE-YYYYMM-NNNN (assigned at post time)
	PeriodID          int           `json:"period_id" gorm:"index;not null"`
	Status            string        `json:"status" gorm:"size:24;default:'posted';index"` // draft | submitted | approved | posted | reversal_pending | reversal_approved | reversed
	PostedAt          time.Time     `json:"posted_at" gorm:"index"`
	PostedBy          string        `json:"posted_by" gorm:"size:191"`
	SourceType        string        `json:"source_type" gorm:"size:32;index"` // manual|claim_payment|premium_allocation|premium_allocation_reversal|write_off|refund|reversal
	SourceID          int           `json:"source_id" gorm:"index"`           // FK-shaped int into the source table (no GORM FK)
	Description       string        `json:"description"`
	IsReversed        bool          `json:"is_reversed" gorm:"default:false"`
	ReversedByEntryID *int          `json:"reversed_by_entry_id"`
	TotalDebit        float64       `json:"total_debit"`
	TotalCredit       float64       `json:"total_credit"`
	CreatedAt         time.Time     `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy         string        `json:"created_by" gorm:"size:191"` // drafter — set on first save
	UpdatedAt         time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy         string        `json:"updated_by" gorm:"size:191"` // last person to mutate the draft

	// Maker/checker workflow timestamps. Posted{By,At} above is the final step.
	SubmittedBy       string        `json:"submitted_by,omitempty" gorm:"size:191"`
	SubmittedAt       *time.Time    `json:"submitted_at,omitempty"`
	ApprovedBy        string        `json:"approved_by,omitempty" gorm:"size:191"`
	ApprovedAt        *time.Time    `json:"approved_at,omitempty"`

	// Reversal workflow — same maker/checker shape applied to a reversal request.
	ReversalReason          string     `json:"reversal_reason,omitempty"`
	ReversalRequestedBy     string     `json:"reversal_requested_by,omitempty" gorm:"size:191"`
	ReversalRequestedAt     *time.Time `json:"reversal_requested_at,omitempty"`
	ReversalApprovedBy      string     `json:"reversal_approved_by,omitempty" gorm:"size:191"`
	ReversalApprovedAt      *time.Time `json:"reversal_approved_at,omitempty"`

	Lines []JournalLine `json:"lines" gorm:"foreignKey:EntryID;references:ID"`
}

func (JournalEntry) TableName() string { return "journal_entries" }

// JournalLine is one debit or credit line of a JournalEntry. Exactly one of
// Debit and Credit is non-zero per line; the SchemeID and CostCentre are
// optional reporting dimensions and may be left zero/empty.
type JournalLine struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	EntryID     int     `json:"entry_id" gorm:"index;not null"`
	AccountID   int     `json:"account_id" gorm:"index;not null"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Description string  `json:"description"`
	SchemeID    int     `json:"scheme_id" gorm:"index"`
	CostCentre  string  `json:"cost_centre" gorm:"size:64;index"`
	LineOrder   int     `json:"line_order"`
}

func (JournalLine) TableName() string { return "journal_lines" }

// ---------------------------------------------------------------------------
// Posting Rules (event → DR/CR account mapping)
// ---------------------------------------------------------------------------

// PostingRule maps an operational event (e.g. "claim_payment.confirmed") to
// the pair of accounts that should be debited/credited when the event fires.
// Lookup is by EventKey; if no active rule is found the posting call fails
// loudly — missing rules are configuration bugs, not warnings.
//
// Posting rules drive automatic GL postings — a rule swap can silently redirect
// every future system entry to the wrong account. Edits therefore follow the
// same maker/checker pending-change pattern as GLAccount.
type PostingRule struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	EventKey           string     `json:"event_key" gorm:"type:varchar(64);uniqueIndex"`
	DebitAccountID     int        `json:"debit_account_id" gorm:"index"`
	CreditAccountID    int        `json:"credit_account_id" gorm:"index"`
	IsActive           bool       `json:"is_active" gorm:"default:true"`
	Notes              string     `json:"notes"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy          string     `json:"created_by" gorm:"size:191"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy          string     `json:"updated_by" gorm:"size:191"`

	// Maker/checker workflow fields — same shape as GLAccount.
	ApprovalStatus     string     `json:"approval_status" gorm:"size:24;default:'active';index"` // active | pending_create | pending_update | pending_delete
	PendingChangeJSON  string     `json:"pending_change_json,omitempty" gorm:"type:text"`
	PendingRequestedBy string     `json:"pending_requested_by,omitempty" gorm:"size:191"`
	PendingRequestedAt *time.Time `json:"pending_requested_at,omitempty"`
}

func (PostingRule) TableName() string { return "posting_rules" }

// ---------------------------------------------------------------------------
// Bank account sub-ledger + statement import
// ---------------------------------------------------------------------------

// BankAccount is the operational view of a bank account — one row per real
// bank account, each pinned to a single GLAccount whose ledger is the source
// of truth for posted activity. Used to scope statement imports and the
// reconciliation workspace.
//
// Bank account number / GL account changes are a payment-fraud vector
// (silently redirect payments to an attacker-controlled destination), so all
// mutations follow the same maker/checker pending-change pattern as GLAccount.
type BankAccount struct {
	ID                 int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Code               string     `json:"code" gorm:"type:varchar(50);uniqueIndex"`
	Name               string     `json:"name" gorm:"type:varchar(191)"`
	BankName           string     `json:"bank_name"`
	AccountNumber      string     `json:"account_number"`
	GLAccountID        int        `json:"gl_account_id" gorm:"index"`
	Currency           string     `json:"currency" gorm:"size:8;default:'ZAR'"`
	IsActive           bool       `json:"is_active" gorm:"default:true"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CreatedBy          string     `json:"created_by" gorm:"size:191"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy          string     `json:"updated_by" gorm:"size:191"`

	// Maker/checker workflow fields — same shape as GLAccount.
	ApprovalStatus     string     `json:"approval_status" gorm:"size:24;default:'active';index"` // active | pending_create | pending_update | pending_deactivate
	PendingChangeJSON  string     `json:"pending_change_json,omitempty" gorm:"type:text"`
	PendingRequestedBy string     `json:"pending_requested_by,omitempty" gorm:"size:191"`
	PendingRequestedAt *time.Time `json:"pending_requested_at,omitempty"`
}

func (BankAccount) TableName() string { return "bank_accounts" }

// BankStatementLine is a single row from an imported bank statement, paired
// to a posted JournalLine once matched. Amount is signed (positive=credit
// into the bank account, negative=debit out).
//
// Match/ignore decisions made by one user (MatchedBy) require sign-off from
// a second user (ReviewedBy) before the line is considered reconciled. The
// review step is enforced at the service layer — ReviewedBy must differ from
// MatchedBy.
type BankStatementLine struct {
	ID                   int        `json:"id" gorm:"primaryKey;autoIncrement"`
	BankAccountID        int        `json:"bank_account_id" gorm:"index"`
	StatementDate        time.Time  `json:"statement_date" gorm:"index"`
	ValueDate            *time.Time `json:"value_date"`
	Description          string     `json:"description"`
	Amount               float64    `json:"amount"`
	Reference            string     `json:"reference" gorm:"index"`
	ImportBatchID        string     `json:"import_batch_id" gorm:"size:64;index"`
	ImportedBy           string     `json:"imported_by" gorm:"size:191"`
	MatchedJournalLineID *int       `json:"matched_journal_line_id" gorm:"index"`
	MatchStatus          string     `json:"match_status" gorm:"size:16;default:'unmatched';index"` // unmatched | matched | ignored
	MatchedAt            *time.Time `json:"matched_at"`
	MatchedBy            string     `json:"matched_by" gorm:"size:191"`

	// Reviewer sign-off — populated only when ReviewStatus moves out of pending.
	ReviewStatus string     `json:"review_status" gorm:"size:24;default:'not_required';index"` // not_required (unmatched) | pending_review | reviewed | rejected
	ReviewedBy   string     `json:"reviewed_by,omitempty" gorm:"size:191"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
	ReviewNotes  string     `json:"review_notes,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (BankStatementLine) TableName() string { return "bank_statement_lines" }

// ---------------------------------------------------------------------------
// Request / response DTOs (used by the controller layer)
// ---------------------------------------------------------------------------

// ManualJournalLineInput is a single line on a caller-supplied manual JE.
type ManualJournalLineInput struct {
	AccountID   int     `json:"account_id" binding:"required"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Description string  `json:"description"`
	SchemeID    int     `json:"scheme_id"`
	CostCentre  string  `json:"cost_centre"`
}

// ManualJournalRequest is the inbound payload for a manually entered JE.
type ManualJournalRequest struct {
	Description string                   `json:"description" binding:"required"`
	PeriodID    int                      `json:"period_id"` // optional; defaults to today's open period
	Lines       []ManualJournalLineInput `json:"lines" binding:"required,min=2"`
}

// ReverseJournalRequest is the inbound payload for requesting a reversal of
// a posted JE. Approval is a separate endpoint that takes no body.
type ReverseJournalRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// StatementLineReviewRequest is the inbound payload for the reviewer sign-off
// on a match or ignore action. Outcome is "reviewed" (sign-off accepted) or
// "rejected" (kick the line back to the matcher; clears match metadata).
type StatementLineReviewRequest struct {
	Outcome string `json:"outcome" binding:"required,oneof=reviewed rejected"`
	Notes   string `json:"notes"`
}

// PendingChangeApprovalRequest is the inbound payload for approving a pending
// master-data change (GLAccount, PostingRule, BankAccount). Body is optional —
// notes are appended to the audit row's Details if supplied.
type PendingChangeApprovalRequest struct {
	Notes string `json:"notes"`
}

// TrialBalanceRow is one row of a trial balance report.
type TrialBalanceRow struct {
	AccountID     int     `json:"account_id"`
	AccountCode   string  `json:"account_code"`
	AccountName   string  `json:"account_name"`
	AccountType   string  `json:"account_type"`
	NormalBalance string  `json:"normal_balance"`
	TotalDebit    float64 `json:"total_debit"`
	TotalCredit   float64 `json:"total_credit"`
	NetBalance    float64 `json:"net_balance"` // signed by normal_balance
}

// LedgerRow is one line in an account ledger drill — denormalised join of
// JournalLine + JournalEntry for display, with running balance computed
// in the service layer.
type LedgerRow struct {
	EntryID        int       `json:"entry_id"`
	EntryNumber    string    `json:"entry_number"`
	PostedAt       time.Time `json:"posted_at"`
	SourceType     string    `json:"source_type"`
	SourceID       int       `json:"source_id"`
	Description    string    `json:"description"`
	LineDescription string   `json:"line_description"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	RunningBalance float64   `json:"running_balance"`
}

// BankStatementImportRow is one row of the CSV ingest payload.
type BankStatementImportRow struct {
	StatementDate string  `json:"statement_date" binding:"required"`
	ValueDate     string  `json:"value_date"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount" binding:"required"`
	Reference     string  `json:"reference"`
}

// BankStatementImportRequest is the inbound payload for a statement upload.
type BankStatementImportRequest struct {
	BankAccountID int                      `json:"bank_account_id" binding:"required"`
	Rows          []BankStatementImportRow `json:"rows" binding:"required,min=1"`
}

// BankMatchRequest matches a statement line to a posted journal line.
type BankMatchRequest struct {
	StatementLineID int `json:"statement_line_id" binding:"required"`
	JournalLineID   int `json:"journal_line_id" binding:"required"`
}

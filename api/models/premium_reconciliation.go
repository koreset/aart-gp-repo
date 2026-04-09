package models

import "time"

// ---------------------------------------------------------------------------
// Reconciliation Run
// ---------------------------------------------------------------------------

// ReconciliationRun is an auditable batch of allocation actions.
type ReconciliationRun struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	RunDate         string     `json:"run_date"`
	RunType         string     `json:"run_type"` // "auto" | "bank_import" | "manual"
	Status          string     `json:"status"`   // "in_progress" | "completed" | "rolled_back"
	InitiatedBy     string     `json:"initiated_by"`
	CompletedAt     *time.Time `json:"completed_at"`
	TotalProcessed  int        `json:"total_processed"`
	TotalMatched    int        `json:"total_matched"`
	TotalUnmatched  int        `json:"total_unmatched"`
	TotalAllocated  float64    `json:"total_allocated"`
	MatchingRuleSet string     `json:"matching_rule_set" gorm:"default:'default'"` // rule profile name
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// ---------------------------------------------------------------------------
// Payment Allocation — the many-to-many ledger
// ---------------------------------------------------------------------------

// PaymentAllocation links a portion of a payment to an invoice.
// This is the core journal entry of the reconciliation ledger.
type PaymentAllocation struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentID       int       `json:"payment_id" gorm:"index"`
	InvoiceID       int       `json:"invoice_id" gorm:"index"`
	RunID           *int      `json:"run_id" gorm:"index"` // nullable for manual allocations outside a run
	AllocatedAmount float64   `json:"allocated_amount"`
	AllocationType  string    `json:"allocation_type"` // "payment" | "credit_note" | "debit_note" | "write_off" | "refund" | "reversal"
	Reference       string    `json:"reference"`
	Notes           string    `json:"notes"`
	AllocatedBy     string    `json:"allocated_by"`
	ReversedByID    *int      `json:"reversed_by_id"` // points to the reversal PaymentAllocation.ID
	IsReversal      bool      `json:"is_reversal" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ---------------------------------------------------------------------------
// Reconciliation Item — workspace item for each payment or invoice
// ---------------------------------------------------------------------------

// ReconciliationItem tracks each payment or invoice through the reconciliation
// workflow. This provides the suspense-account view and aging.
type ReconciliationItem struct {
	ID                int        `json:"id" gorm:"primaryKey;autoIncrement"`
	ItemType          string     `json:"item_type" gorm:"index"` // "payment" | "invoice"
	PaymentID         *int       `json:"payment_id" gorm:"index"`
	InvoiceID         *int       `json:"invoice_id" gorm:"index"`
	SchemeID          int        `json:"scheme_id" gorm:"index"`
	SchemeName        string     `json:"scheme_name"`
	OriginalAmount    float64    `json:"original_amount"`
	AllocatedAmount   float64    `json:"allocated_amount"`
	UnallocatedAmount float64    `json:"unallocated_amount"`
	Status            string     `json:"status"` // "open" | "partial" | "matched" | "written_off" | "refunded" | "suspended"
	SuspenseReason    string     `json:"suspense_reason"`
	AgeInDays         int        `json:"age_in_days" gorm:"-"`             // computed, not stored
	Priority          string     `json:"priority" gorm:"default:'normal'"` // "low" | "normal" | "high" | "critical"
	AssignedTo        string     `json:"assigned_to"`
	LastActionDate    *time.Time `json:"last_action_date"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// ---------------------------------------------------------------------------
// Matching Rule Configuration
// ---------------------------------------------------------------------------

// MatchingRule defines a rule in the multi-strategy matching engine.
type MatchingRule struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RuleSet       string    `json:"rule_set" gorm:"default:'default'"` // group rules into profiles
	Priority      int       `json:"priority"`                          // execution order (1 = first)
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Strategy      string    `json:"strategy"`                            // "exact_reference" | "scheme_amount" | "scheme_amount_tolerance" | "scheme_date_range" | "amount_only"
	ToleranceType string    `json:"tolerance_type"`                      // "absolute" | "percentage"
	ToleranceVal  float64   `json:"tolerance_value" gorm:"default:0.01"` // e.g. 0.01 for ±R0.01 or 1.0 for ±1%
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	AllowPartial  bool      `json:"allow_partial" gorm:"default:true"`       // allow partial allocation
	AllowMulti    bool      `json:"allow_multi_invoice" gorm:"default:true"` // allocate across multiple invoices
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ---------------------------------------------------------------------------
// Write-off / Refund
// ---------------------------------------------------------------------------

// WriteOffRequest is the request body for writing off a small balance.
type WriteOffRequest struct {
	ReconciliationItemID int     `json:"reconciliation_item_id" binding:"required"`
	Amount               float64 `json:"amount" binding:"required"`
	Reason               string  `json:"reason" binding:"required"`
	InvoiceID            int     `json:"invoice_id"`
}

// RefundRequest is the request body for initiating a refund of overpayment.
type RefundRequest struct {
	ReconciliationItemID int     `json:"reconciliation_item_id" binding:"required"`
	Amount               float64 `json:"amount" binding:"required"`
	Reason               string  `json:"reason" binding:"required"`
	RefundMethod         string  `json:"refund_method" binding:"required"` // "eft" | "cheque"
	BankDetails          string  `json:"bank_details"`
}

// ---------------------------------------------------------------------------
// Request / Response types
// ---------------------------------------------------------------------------

// RunAutoMatchRequest lets the caller optionally scope the run.
type RunAutoMatchRequest struct {
	SchemeID int    `json:"scheme_id"` // 0 = all schemes
	RuleSet  string `json:"rule_set"`  // "" = "default"
	DryRun   bool   `json:"dry_run"`   // preview without persisting
}

// AllocatePaymentRequest is for manual allocation of a payment to one or more invoices.
type AllocatePaymentRequest struct {
	PaymentID   int                  `json:"payment_id" binding:"required"`
	Allocations []AllocationLineItem `json:"allocations" binding:"required"`
	Notes       string               `json:"notes"`
}

// AllocationLineItem is a single allocation line within an allocation request.
type AllocationLineItem struct {
	InvoiceID int     `json:"invoice_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

// ReverseAllocationRequest reverses one or more allocations.
type ReverseAllocationRequest struct {
	AllocationIDs []int  `json:"allocation_ids" binding:"required"`
	Reason        string `json:"reason" binding:"required"`
}

// ReconciliationSummary is the response for the reconciliation dashboard.
type ReconciliationSummary struct {
	TotalUnallocatedPayments float64             `json:"total_unallocated_payments"`
	UnallocatedPaymentCount  int                 `json:"unallocated_payment_count"`
	TotalUnpaidInvoices      float64             `json:"total_unpaid_invoices"`
	UnpaidInvoiceCount       int                 `json:"unpaid_invoice_count"`
	SuspenseBalance          float64             `json:"suspense_balance"`
	SuspenseCount            int                 `json:"suspense_count"`
	AgedOver30Days           float64             `json:"aged_over_30_days"`
	AgedOver60Days           float64             `json:"aged_over_60_days"`
	AgedOver90Days           float64             `json:"aged_over_90_days"`
	RecentRuns               []ReconciliationRun `json:"recent_runs"`
}

// ReconciliationRunDetail extends a run with its allocations.
type ReconciliationRunDetail struct {
	ReconciliationRun
	Allocations []PaymentAllocation `json:"allocations"`
}

// AllocationHistory is the full allocation trail for a payment or invoice.
type AllocationHistory struct {
	EntityType  string              `json:"entity_type"` // "payment" | "invoice"
	EntityID    int                 `json:"entity_id"`
	Total       float64             `json:"total"`
	Allocated   float64             `json:"allocated"`
	Unallocated float64             `json:"unallocated"`
	Allocations []PaymentAllocation `json:"allocations"`
}

// AutoMatchPreview is a dry-run result showing what would be matched.
type AutoMatchPreview struct {
	ProposedAllocations []ProposedAllocation `json:"proposed_allocations"`
	TotalAllocated      float64              `json:"total_allocated"`
	TotalMatched        int                  `json:"total_matched"`
	TotalRemaining      int                  `json:"total_remaining"`
}

// ProposedAllocation is a single proposed match in a dry-run.
type ProposedAllocation struct {
	PaymentID     int     `json:"payment_id"`
	InvoiceID     int     `json:"invoice_id"`
	Amount        float64 `json:"amount"`
	MatchedBy     string  `json:"matched_by"` // rule name
	Confidence    string  `json:"confidence"` // "high" | "medium" | "low"
	PaymentRef    string  `json:"payment_ref"`
	InvoiceNumber string  `json:"invoice_number"`
	SchemeName    string  `json:"scheme_name"`
}

// ReassignItemRequest reassigns a reconciliation item to another user.
type ReassignItemRequest struct {
	AssignedTo string `json:"assigned_to" binding:"required"`
	Priority   string `json:"priority"`
}

package models

import "time"

// GLAuditLog records every significant event on the operational General Ledger
// and bank/cash sub-ledger — journal lifecycle transitions, period close
// requests/approvals, master-data pending changes, reconciliation reviews.
//
// Deliberately independent of the IFRS17 / CSM audit log: that lives under
// the read-only valuation stack and has different retention and reporting
// requirements. Nothing in this package imports or shares structure with the
// IFRS17 audit log beyond the obvious column shape.
type GLAuditLog struct {
	ID         int       `json:"id"          gorm:"primaryKey;autoIncrement"`
	EventType  string    `json:"event_type"  gorm:"size:64;index"`  // journal_drafted | journal_submitted | journal_approved | journal_posted | journal_reversal_requested | journal_reversal_approved | journal_reversed | period_close_requested | period_closed | gl_account_change_requested | gl_account_change_approved | posting_rule_change_requested | posting_rule_change_approved | bank_account_change_requested | bank_account_change_approved | statement_line_matched | statement_line_match_reviewed | statement_line_ignored | statement_line_ignore_reviewed
	ObjectType string    `json:"object_type" gorm:"size:32;index"`  // journal_entry | accounting_period | gl_account | posting_rule | bank_account | bank_statement_line
	ObjectID   int       `json:"object_id"   gorm:"index"`
	ObjectName string    `json:"object_name" gorm:"size:191"`       // entry number / period name / account code
	ChangedBy  string    `json:"changed_by"  gorm:"size:191;index"` // user.UserName
	ChangedAt  time.Time `json:"changed_at"  gorm:"autoCreateTime;index"`
	Details    string    `json:"details"     gorm:"type:text"` // JSON: { from, to, reason, fields_changed }
}

func (GLAuditLog) TableName() string { return "gl_audit_logs" }

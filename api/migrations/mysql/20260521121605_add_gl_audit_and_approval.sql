-- Generated 2026-05-21T12:16:05+02:00 for dialect mysql

-- Migration for: GLAccount (table: gl_accounts)

ALTER TABLE gl_accounts ADD COLUMN created_by VARCHAR(191);
ALTER TABLE gl_accounts ADD COLUMN updated_by VARCHAR(191);
ALTER TABLE gl_accounts ADD COLUMN approval_status VARCHAR(24) DEFAULT 'active';
ALTER TABLE gl_accounts ADD COLUMN pending_change_json text;
ALTER TABLE gl_accounts ADD COLUMN pending_requested_by VARCHAR(191);
ALTER TABLE gl_accounts ADD COLUMN pending_requested_at DATETIME;
CREATE INDEX idx_gl_accounts_approval_status ON gl_accounts (approval_status);

-- Migration for: AccountingPeriod (table: accounting_periods)

ALTER TABLE accounting_periods ADD COLUMN close_requested_by VARCHAR(191);
ALTER TABLE accounting_periods ADD COLUMN close_requested_at DATETIME;
CREATE INDEX idx_accounting_periods_status ON accounting_periods (status);

-- Migration for: JournalEntry (table: journal_entries)

ALTER TABLE journal_entries ADD COLUMN status VARCHAR(24) DEFAULT 'posted';
ALTER TABLE journal_entries ADD COLUMN created_by VARCHAR(191);
ALTER TABLE journal_entries ADD COLUMN updated_at DATETIME;
ALTER TABLE journal_entries ADD COLUMN updated_by VARCHAR(191);
ALTER TABLE journal_entries ADD COLUMN submitted_by VARCHAR(191);
ALTER TABLE journal_entries ADD COLUMN submitted_at DATETIME;
ALTER TABLE journal_entries ADD COLUMN approved_by VARCHAR(191);
ALTER TABLE journal_entries ADD COLUMN approved_at DATETIME;
ALTER TABLE journal_entries ADD COLUMN reversal_reason VARCHAR(255);
ALTER TABLE journal_entries ADD COLUMN reversal_requested_by VARCHAR(191);
ALTER TABLE journal_entries ADD COLUMN reversal_requested_at DATETIME;
ALTER TABLE journal_entries ADD COLUMN reversal_approved_by VARCHAR(191);
ALTER TABLE journal_entries ADD COLUMN reversal_approved_at DATETIME;
CREATE INDEX idx_journal_entries_status ON journal_entries (status);

-- Migration for: PostingRule (table: posting_rules)

ALTER TABLE posting_rules ADD COLUMN created_by VARCHAR(191);
ALTER TABLE posting_rules ADD COLUMN updated_by VARCHAR(191);
ALTER TABLE posting_rules ADD COLUMN approval_status VARCHAR(24) DEFAULT 'active';
ALTER TABLE posting_rules ADD COLUMN pending_change_json text;
ALTER TABLE posting_rules ADD COLUMN pending_requested_by VARCHAR(191);
ALTER TABLE posting_rules ADD COLUMN pending_requested_at DATETIME;
CREATE INDEX idx_posting_rules_approval_status ON posting_rules (approval_status);

-- Migration for: BankAccount (table: bank_accounts)

ALTER TABLE bank_accounts ADD COLUMN created_by VARCHAR(191);
ALTER TABLE bank_accounts ADD COLUMN updated_by VARCHAR(191);
ALTER TABLE bank_accounts ADD COLUMN approval_status VARCHAR(24) DEFAULT 'active';
ALTER TABLE bank_accounts ADD COLUMN pending_change_json text;
ALTER TABLE bank_accounts ADD COLUMN pending_requested_by VARCHAR(191);
ALTER TABLE bank_accounts ADD COLUMN pending_requested_at DATETIME;
CREATE INDEX idx_bank_accounts_approval_status ON bank_accounts (approval_status);

-- Migration for: BankStatementLine (table: bank_statement_lines)

ALTER TABLE bank_statement_lines ADD COLUMN imported_by VARCHAR(191);
ALTER TABLE bank_statement_lines ADD COLUMN review_status VARCHAR(24) DEFAULT 'not_required';
ALTER TABLE bank_statement_lines ADD COLUMN reviewed_by VARCHAR(191);
ALTER TABLE bank_statement_lines ADD COLUMN reviewed_at DATETIME;
ALTER TABLE bank_statement_lines ADD COLUMN review_notes VARCHAR(255);
CREATE INDEX idx_bank_statement_lines_match_status ON bank_statement_lines (match_status);
CREATE INDEX idx_bank_statement_lines_review_status ON bank_statement_lines (review_status);

-- Create table: gl_audit_logs
CREATE TABLE gl_audit_logs (
    id BIGINT AUTO_INCREMENT,
    event_type VARCHAR(64),
    object_type VARCHAR(32),
    object_id BIGINT,
    object_name VARCHAR(191),
    changed_by VARCHAR(191),
    changed_at DATETIME,
    details text,
    PRIMARY KEY (id)
);

CREATE INDEX idx_gl_audit_logs_event_type ON gl_audit_logs (event_type);
CREATE INDEX idx_gl_audit_logs_object_type ON gl_audit_logs (object_type);
CREATE INDEX idx_gl_audit_logs_object_id ON gl_audit_logs (object_id);
CREATE INDEX idx_gl_audit_logs_changed_by ON gl_audit_logs (changed_by);
CREATE INDEX idx_gl_audit_logs_changed_at ON gl_audit_logs (changed_at);



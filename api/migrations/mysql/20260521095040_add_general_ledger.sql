-- Generated 2026-05-21T09:50:40+02:00 for dialect mysql
-- Operational General Ledger — chart of accounts, accounting periods,
-- journal entries + lines, posting rules, bank accounts + statement lines.
-- All CREATE TABLE statements use IF NOT EXISTS so the file is safe to
-- re-apply if a previous run failed partway through.
--
-- To re-run cleanly from scratch (drops any partial state):
--   DELETE FROM migrations WHERE version = '20260521095040';
--   DROP TABLE IF EXISTS bank_statement_lines;
--   DROP TABLE IF EXISTS bank_accounts;
--   DROP TABLE IF EXISTS posting_rules;
--   DROP TABLE IF EXISTS journal_lines;
--   DROP TABLE IF EXISTS journal_entries;
--   DROP TABLE IF EXISTS accounting_periods;
--   DROP TABLE IF EXISTS gl_accounts;
-- Then restart the API.

-- Create table: gl_accounts
CREATE TABLE IF NOT EXISTS gl_accounts (
    id BIGINT AUTO_INCREMENT,
    code varchar(50),
    name varchar(191),
    account_type VARCHAR(32),
    normal_balance VARCHAR(8),
    parent_id BIGINT,
    is_active TINYINT(1) DEFAULT 1,
    description VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_gl_accounts_code ON gl_accounts (code);
CREATE INDEX idx_gl_accounts_account_type ON gl_accounts (account_type);
CREATE INDEX idx_gl_accounts_parent_id ON gl_accounts (parent_id);


-- Create table: accounting_periods
CREATE TABLE IF NOT EXISTS accounting_periods (
    id BIGINT AUTO_INCREMENT,
    name varchar(20),
    start_date DATETIME,
    end_date DATETIME,
    status VARCHAR(16) DEFAULT 'open',
    closed_at DATETIME,
    closed_by VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_accounting_periods_name ON accounting_periods (name);
CREATE INDEX idx_accounting_periods_start_date ON accounting_periods (start_date);


-- Create table: journal_entries
CREATE TABLE IF NOT EXISTS journal_entries (
    id BIGINT AUTO_INCREMENT,
    entry_number varchar(32),
    period_id BIGINT NOT NULL,
    posted_at DATETIME,
    posted_by VARCHAR(255),
    source_type VARCHAR(32),
    source_id BIGINT,
    description VARCHAR(255),
    is_reversed TINYINT(1) DEFAULT 0,
    reversed_by_entry_id BIGINT,
    total_debit DOUBLE,
    total_credit DOUBLE,
    created_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_journal_entries_entry_number ON journal_entries (entry_number);
CREATE INDEX idx_journal_entries_period_id ON journal_entries (period_id);
CREATE INDEX idx_journal_entries_posted_at ON journal_entries (posted_at);
CREATE INDEX idx_journal_entries_source_type ON journal_entries (source_type);
CREATE INDEX idx_journal_entries_source_id ON journal_entries (source_id);


-- Create table: journal_lines
CREATE TABLE IF NOT EXISTS journal_lines (
    id BIGINT AUTO_INCREMENT,
    entry_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    debit DOUBLE,
    credit DOUBLE,
    description VARCHAR(255),
    scheme_id BIGINT,
    cost_centre VARCHAR(64),
    line_order BIGINT,
    PRIMARY KEY (id)
);

CREATE INDEX idx_journal_lines_entry_id ON journal_lines (entry_id);
CREATE INDEX idx_journal_lines_account_id ON journal_lines (account_id);
CREATE INDEX idx_journal_lines_scheme_id ON journal_lines (scheme_id);
CREATE INDEX idx_journal_lines_cost_centre ON journal_lines (cost_centre);


-- Create table: posting_rules
CREATE TABLE IF NOT EXISTS posting_rules (
    id BIGINT AUTO_INCREMENT,
    event_key varchar(64),
    debit_account_id BIGINT,
    credit_account_id BIGINT,
    is_active TINYINT(1) DEFAULT 1,
    notes VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_posting_rules_event_key ON posting_rules (event_key);
CREATE INDEX idx_posting_rules_debit_account_id ON posting_rules (debit_account_id);
CREATE INDEX idx_posting_rules_credit_account_id ON posting_rules (credit_account_id);


-- Create table: bank_accounts
CREATE TABLE IF NOT EXISTS bank_accounts (
    id BIGINT AUTO_INCREMENT,
    code varchar(50),
    name varchar(191),
    bank_name VARCHAR(255),
    account_number VARCHAR(255),
    gl_account_id BIGINT,
    currency VARCHAR(8) DEFAULT 'ZAR',
    is_active TINYINT(1) DEFAULT 1,
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_bank_accounts_code ON bank_accounts (code);
CREATE INDEX idx_bank_accounts_gl_account_id ON bank_accounts (gl_account_id);


-- Create table: bank_statement_lines
CREATE TABLE IF NOT EXISTS bank_statement_lines (
    id BIGINT AUTO_INCREMENT,
    bank_account_id BIGINT,
    statement_date DATETIME,
    value_date DATETIME,
    description VARCHAR(255),
    amount DOUBLE,
    reference VARCHAR(255),
    import_batch_id VARCHAR(64),
    matched_journal_line_id BIGINT,
    match_status VARCHAR(16) DEFAULT 'unmatched',
    matched_at DATETIME,
    matched_by VARCHAR(255),
    created_at DATETIME,
    PRIMARY KEY (id)
);

CREATE INDEX idx_bank_statement_lines_bank_account_id ON bank_statement_lines (bank_account_id);
CREATE INDEX idx_bank_statement_lines_statement_date ON bank_statement_lines (statement_date);
CREATE INDEX idx_bank_statement_lines_reference ON bank_statement_lines (reference);
CREATE INDEX idx_bank_statement_lines_import_batch_id ON bank_statement_lines (import_batch_id);
CREATE INDEX idx_bank_statement_lines_matched_journal_line_id ON bank_statement_lines (matched_journal_line_id);

-- Operational General Ledger — chart of accounts, accounting periods,
-- journal entries + lines, posting rules, bank accounts + statement lines.
-- All CREATE statements use IF NOT EXISTS so the file is safe to re-apply.
--
-- To re-run cleanly from scratch (drops any partial state):
--   DELETE FROM migrations WHERE version = '20260521095040';
--   DROP TABLE IF EXISTS bank_statement_lines, bank_accounts, posting_rules,
--                        journal_lines, journal_entries, accounting_periods,
--                        gl_accounts CASCADE;
-- Then restart the API.

CREATE TABLE IF NOT EXISTS gl_accounts (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(50),
    name            VARCHAR(191),
    account_type    VARCHAR(32),
    normal_balance  VARCHAR(8),
    parent_id       BIGINT,
    is_active       BOOLEAN DEFAULT TRUE,
    description     VARCHAR(255),
    created_at      TIMESTAMP,
    updated_at      TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_gl_accounts_code ON gl_accounts (code);
CREATE INDEX IF NOT EXISTS idx_gl_accounts_account_type ON gl_accounts (account_type);
CREATE INDEX IF NOT EXISTS idx_gl_accounts_parent_id ON gl_accounts (parent_id);


CREATE TABLE IF NOT EXISTS accounting_periods (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(20),
    start_date  TIMESTAMP,
    end_date    TIMESTAMP,
    status      VARCHAR(16) DEFAULT 'open',
    closed_at   TIMESTAMP,
    closed_by   VARCHAR(255),
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_accounting_periods_name ON accounting_periods (name);
CREATE INDEX IF NOT EXISTS idx_accounting_periods_start_date ON accounting_periods (start_date);


CREATE TABLE IF NOT EXISTS journal_entries (
    id                    BIGSERIAL PRIMARY KEY,
    entry_number          VARCHAR(32),
    period_id             BIGINT NOT NULL,
    posted_at             TIMESTAMP,
    posted_by             VARCHAR(255),
    source_type           VARCHAR(32),
    source_id             BIGINT,
    description           VARCHAR(255),
    is_reversed           BOOLEAN DEFAULT FALSE,
    reversed_by_entry_id  BIGINT,
    total_debit           DOUBLE PRECISION,
    total_credit          DOUBLE PRECISION,
    created_at            TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_journal_entries_entry_number ON journal_entries (entry_number);
CREATE INDEX IF NOT EXISTS idx_journal_entries_period_id ON journal_entries (period_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_posted_at ON journal_entries (posted_at);
CREATE INDEX IF NOT EXISTS idx_journal_entries_source_type ON journal_entries (source_type);
CREATE INDEX IF NOT EXISTS idx_journal_entries_source_id ON journal_entries (source_id);


CREATE TABLE IF NOT EXISTS journal_lines (
    id           BIGSERIAL PRIMARY KEY,
    entry_id     BIGINT NOT NULL,
    account_id   BIGINT NOT NULL,
    debit        DOUBLE PRECISION,
    credit       DOUBLE PRECISION,
    description  VARCHAR(255),
    scheme_id    BIGINT,
    cost_centre  VARCHAR(64),
    line_order   BIGINT
);

CREATE INDEX IF NOT EXISTS idx_journal_lines_entry_id ON journal_lines (entry_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_account_id ON journal_lines (account_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_scheme_id ON journal_lines (scheme_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_cost_centre ON journal_lines (cost_centre);


CREATE TABLE IF NOT EXISTS posting_rules (
    id                 BIGSERIAL PRIMARY KEY,
    event_key          VARCHAR(64),
    debit_account_id   BIGINT,
    credit_account_id  BIGINT,
    is_active          BOOLEAN DEFAULT TRUE,
    notes              VARCHAR(255),
    created_at         TIMESTAMP,
    updated_at         TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_posting_rules_event_key ON posting_rules (event_key);
CREATE INDEX IF NOT EXISTS idx_posting_rules_debit_account_id ON posting_rules (debit_account_id);
CREATE INDEX IF NOT EXISTS idx_posting_rules_credit_account_id ON posting_rules (credit_account_id);


CREATE TABLE IF NOT EXISTS bank_accounts (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(50),
    name            VARCHAR(191),
    bank_name       VARCHAR(255),
    account_number  VARCHAR(255),
    gl_account_id   BIGINT,
    currency        VARCHAR(8) DEFAULT 'ZAR',
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP,
    updated_at      TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bank_accounts_code ON bank_accounts (code);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_gl_account_id ON bank_accounts (gl_account_id);


CREATE TABLE IF NOT EXISTS bank_statement_lines (
    id                       BIGSERIAL PRIMARY KEY,
    bank_account_id          BIGINT,
    statement_date           TIMESTAMP,
    value_date               TIMESTAMP,
    description              VARCHAR(255),
    amount                   DOUBLE PRECISION,
    reference                VARCHAR(255),
    import_batch_id          VARCHAR(64),
    matched_journal_line_id  BIGINT,
    match_status             VARCHAR(16) DEFAULT 'unmatched',
    matched_at               TIMESTAMP,
    matched_by               VARCHAR(255),
    created_at               TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_bank_account_id ON bank_statement_lines (bank_account_id);
CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_statement_date ON bank_statement_lines (statement_date);
CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_reference ON bank_statement_lines (reference);
CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_import_batch_id ON bank_statement_lines (import_batch_id);
CREATE INDEX IF NOT EXISTS idx_bank_statement_lines_matched_journal_line_id ON bank_statement_lines (matched_journal_line_id);

-- Operational General Ledger — chart of accounts, accounting periods,
-- journal entries + lines, posting rules, bank accounts + statement lines.
-- Each CREATE TABLE / CREATE INDEX is wrapped in an existence check so the
-- file is safe to re-apply.
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

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'gl_accounts')
BEGIN
    CREATE TABLE gl_accounts (
        id              BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_gl_accounts PRIMARY KEY,
        code            NVARCHAR(50)  NULL,
        name            NVARCHAR(191) NULL,
        account_type    NVARCHAR(32)  NULL,
        normal_balance  NVARCHAR(8)   NULL,
        parent_id       BIGINT        NULL,
        is_active       BIT           NULL CONSTRAINT df_gl_accounts_is_active DEFAULT 1,
        description     NVARCHAR(255) NULL,
        created_at      DATETIME      NULL,
        updated_at      DATETIME      NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gl_accounts_code')
    CREATE UNIQUE INDEX idx_gl_accounts_code ON gl_accounts (code);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gl_accounts_account_type')
    CREATE INDEX idx_gl_accounts_account_type ON gl_accounts (account_type);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gl_accounts_parent_id')
    CREATE INDEX idx_gl_accounts_parent_id ON gl_accounts (parent_id);


IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'accounting_periods')
BEGIN
    CREATE TABLE accounting_periods (
        id          BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_accounting_periods PRIMARY KEY,
        name        NVARCHAR(20)  NULL,
        start_date  DATETIME      NULL,
        end_date    DATETIME      NULL,
        status      NVARCHAR(16)  NULL CONSTRAINT df_accounting_periods_status DEFAULT 'open',
        closed_at   DATETIME      NULL,
        closed_by   NVARCHAR(255) NULL,
        created_at  DATETIME      NULL,
        updated_at  DATETIME      NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_accounting_periods_name')
    CREATE UNIQUE INDEX idx_accounting_periods_name ON accounting_periods (name);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_accounting_periods_start_date')
    CREATE INDEX idx_accounting_periods_start_date ON accounting_periods (start_date);


IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'journal_entries')
BEGIN
    CREATE TABLE journal_entries (
        id                    BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_journal_entries PRIMARY KEY,
        entry_number          NVARCHAR(32)  NULL,
        period_id             BIGINT        NOT NULL,
        posted_at             DATETIME      NULL,
        posted_by             NVARCHAR(255) NULL,
        source_type           NVARCHAR(32)  NULL,
        source_id             BIGINT        NULL,
        description           NVARCHAR(255) NULL,
        is_reversed           BIT           NULL CONSTRAINT df_journal_entries_is_reversed DEFAULT 0,
        reversed_by_entry_id  BIGINT        NULL,
        total_debit           FLOAT         NULL,
        total_credit          FLOAT         NULL,
        created_at            DATETIME      NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_entries_entry_number')
    CREATE UNIQUE INDEX idx_journal_entries_entry_number ON journal_entries (entry_number);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_entries_period_id')
    CREATE INDEX idx_journal_entries_period_id ON journal_entries (period_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_entries_posted_at')
    CREATE INDEX idx_journal_entries_posted_at ON journal_entries (posted_at);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_entries_source_type')
    CREATE INDEX idx_journal_entries_source_type ON journal_entries (source_type);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_entries_source_id')
    CREATE INDEX idx_journal_entries_source_id ON journal_entries (source_id);


IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'journal_lines')
BEGIN
    CREATE TABLE journal_lines (
        id           BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_journal_lines PRIMARY KEY,
        entry_id     BIGINT        NOT NULL,
        account_id   BIGINT        NOT NULL,
        debit        FLOAT         NULL,
        credit       FLOAT         NULL,
        description  NVARCHAR(255) NULL,
        scheme_id    BIGINT        NULL,
        cost_centre  NVARCHAR(64)  NULL,
        line_order   BIGINT        NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_lines_entry_id')
    CREATE INDEX idx_journal_lines_entry_id ON journal_lines (entry_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_lines_account_id')
    CREATE INDEX idx_journal_lines_account_id ON journal_lines (account_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_lines_scheme_id')
    CREATE INDEX idx_journal_lines_scheme_id ON journal_lines (scheme_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_journal_lines_cost_centre')
    CREATE INDEX idx_journal_lines_cost_centre ON journal_lines (cost_centre);


IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'posting_rules')
BEGIN
    CREATE TABLE posting_rules (
        id                 BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_posting_rules PRIMARY KEY,
        event_key          NVARCHAR(64)  NULL,
        debit_account_id   BIGINT        NULL,
        credit_account_id  BIGINT        NULL,
        is_active          BIT           NULL CONSTRAINT df_posting_rules_is_active DEFAULT 1,
        notes              NVARCHAR(255) NULL,
        created_at         DATETIME      NULL,
        updated_at         DATETIME      NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_posting_rules_event_key')
    CREATE UNIQUE INDEX idx_posting_rules_event_key ON posting_rules (event_key);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_posting_rules_debit_account_id')
    CREATE INDEX idx_posting_rules_debit_account_id ON posting_rules (debit_account_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_posting_rules_credit_account_id')
    CREATE INDEX idx_posting_rules_credit_account_id ON posting_rules (credit_account_id);


IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'bank_accounts')
BEGIN
    CREATE TABLE bank_accounts (
        id              BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_bank_accounts PRIMARY KEY,
        code            NVARCHAR(50)  NULL,
        name            NVARCHAR(191) NULL,
        bank_name       NVARCHAR(255) NULL,
        account_number  NVARCHAR(255) NULL,
        gl_account_id   BIGINT        NULL,
        currency        NVARCHAR(8)   NULL CONSTRAINT df_bank_accounts_currency DEFAULT 'ZAR',
        is_active       BIT           NULL CONSTRAINT df_bank_accounts_is_active DEFAULT 1,
        created_at      DATETIME      NULL,
        updated_at      DATETIME      NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bank_accounts_code')
    CREATE UNIQUE INDEX idx_bank_accounts_code ON bank_accounts (code);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bank_accounts_gl_account_id')
    CREATE INDEX idx_bank_accounts_gl_account_id ON bank_accounts (gl_account_id);


IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'bank_statement_lines')
BEGIN
    CREATE TABLE bank_statement_lines (
        id                       BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_bank_statement_lines PRIMARY KEY,
        bank_account_id          BIGINT        NULL,
        statement_date           DATETIME      NULL,
        value_date               DATETIME      NULL,
        description              NVARCHAR(255) NULL,
        amount                   FLOAT         NULL,
        reference                NVARCHAR(255) NULL,
        import_batch_id          NVARCHAR(64)  NULL,
        matched_journal_line_id  BIGINT        NULL,
        match_status             NVARCHAR(16)  NULL CONSTRAINT df_bank_statement_lines_match_status DEFAULT 'unmatched',
        matched_at               DATETIME      NULL,
        matched_by               NVARCHAR(255) NULL,
        created_at               DATETIME      NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bank_statement_lines_bank_account_id')
    CREATE INDEX idx_bank_statement_lines_bank_account_id ON bank_statement_lines (bank_account_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bank_statement_lines_statement_date')
    CREATE INDEX idx_bank_statement_lines_statement_date ON bank_statement_lines (statement_date);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bank_statement_lines_reference')
    CREATE INDEX idx_bank_statement_lines_reference ON bank_statement_lines (reference);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bank_statement_lines_import_batch_id')
    CREATE INDEX idx_bank_statement_lines_import_batch_id ON bank_statement_lines (import_batch_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bank_statement_lines_matched_journal_line_id')
    CREATE INDEX idx_bank_statement_lines_matched_journal_line_id ON bank_statement_lines (matched_journal_line_id);

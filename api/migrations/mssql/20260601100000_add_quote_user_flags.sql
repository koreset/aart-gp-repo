-- Quote User Flags: lets managers flag a user on the dashboard for
-- coaching or capacity attention with an internal note, and resolve the
-- flag later with a resolution note. Append-only — history matters
-- (knowing "Alice has been flagged for capacity 3× this quarter" is the
-- point). At most one open flag per (user_name, flag_reason) is enforced
-- in the service layer, not via a partial index, so MySQL is consistent.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'quote_user_flags')
CREATE TABLE quote_user_flags (
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_name NVARCHAR(255) NOT NULL,
    user_email NVARCHAR(255) NULL,
    flag_reason NVARCHAR(50) NOT NULL,
    note NVARCHAR(MAX) NOT NULL,
    opened_by NVARCHAR(255) NOT NULL,
    opened_by_name NVARCHAR(255) NOT NULL,
    opened_at DATETIME NOT NULL CONSTRAINT df_quf_opened_at DEFAULT GETDATE(),
    resolved_by NVARCHAR(255) NULL,
    resolved_by_name NVARCHAR(255) NULL,
    resolved_at DATETIME NULL,
    resolution_note NVARCHAR(MAX) NULL
);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_quf_user' AND object_id = OBJECT_ID('quote_user_flags'))
  CREATE INDEX idx_quf_user ON quote_user_flags(user_name);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_quf_opened_at' AND object_id = OBJECT_ID('quote_user_flags'))
  CREATE INDEX idx_quf_opened_at ON quote_user_flags(opened_at);

-- Filtered index lets us cheaply look up open flags for the leaderboard
-- without scanning resolved history.
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_quf_open' AND object_id = OBJECT_ID('quote_user_flags'))
  CREATE INDEX idx_quf_open ON quote_user_flags(user_name, flag_reason) WHERE resolved_at IS NULL;

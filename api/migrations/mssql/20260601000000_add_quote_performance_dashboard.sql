-- Quote Performance Dashboard: adds per-status milestone timestamps to
-- group_pricing_quotes, a dedicated audit table for quote status
-- transitions, an admin-editable SLA targets table with seed defaults,
-- and a best-effort backfill of historical data so the dashboard renders
-- meaningful numbers from day one. Idempotent on re-runs.

-- 1) Per-status milestone timestamps on group_pricing_quotes.
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'submitted_at')
  ALTER TABLE group_pricing_quotes ADD submitted_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'approved_at')
  ALTER TABLE group_pricing_quotes ADD approved_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'rejected_at')
  ALTER TABLE group_pricing_quotes ADD rejected_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'accepted_at')
  ALTER TABLE group_pricing_quotes ADD accepted_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'in_force_at')
  ALTER TABLE group_pricing_quotes ADD in_force_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'rejected_reason')
  ALTER TABLE group_pricing_quotes ADD rejected_reason NVARCHAR(500) NULL;

-- 1a) GORM's default mapping for `string` fields is NVARCHAR(MAX), which
-- can't participate in indexes. The composite index below references
-- created_by and status, so widen both to sized NVARCHARs first.
-- Idempotent: only widens when the column is currently MAX / TEXT.
-- system_type_id 35 = text, 99 = ntext.
IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('group_pricing_quotes')
             AND name = 'created_by'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE group_pricing_quotes ALTER COLUMN created_by NVARCHAR(255) NULL;

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('group_pricing_quotes')
             AND name = 'status'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE group_pricing_quotes ALTER COLUMN status NVARCHAR(50) NULL;

-- 2) Indexes for dashboard date-range and per-user aggregation queries.
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpq_submitted_at' AND object_id = OBJECT_ID('group_pricing_quotes'))
  CREATE INDEX idx_gpq_submitted_at ON group_pricing_quotes(submitted_at) WHERE submitted_at IS NOT NULL;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpq_approved_at' AND object_id = OBJECT_ID('group_pricing_quotes'))
  CREATE INDEX idx_gpq_approved_at ON group_pricing_quotes(approved_at) WHERE approved_at IS NOT NULL;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpq_accepted_at' AND object_id = OBJECT_ID('group_pricing_quotes'))
  CREATE INDEX idx_gpq_accepted_at ON group_pricing_quotes(accepted_at) WHERE accepted_at IS NOT NULL;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpq_created_by_status' AND object_id = OBJECT_ID('group_pricing_quotes'))
  CREATE INDEX idx_gpq_created_by_status ON group_pricing_quotes(created_by, status);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpq_created_by_status_submitted_at' AND object_id = OBJECT_ID('group_pricing_quotes'))
  CREATE INDEX idx_gpq_created_by_status_submitted_at ON group_pricing_quotes(created_by, status, submitted_at);

-- 3) Quote status audit table (mirrors group_scheme_status_audits).
IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'group_pricing_quote_status_audits')
CREATE TABLE group_pricing_quote_status_audits (
    id INT IDENTITY(1,1) PRIMARY KEY,
    quote_id INT NOT NULL,
    old_status NVARCHAR(50) NULL,
    new_status NVARCHAR(50) NULL,
    status_message NVARCHAR(500) NULL,
    changed_by NVARCHAR(255) NULL,
    changed_at DATETIME NOT NULL CONSTRAINT df_gpqsa_changed_at DEFAULT GETDATE(),
    duration_from_prev_secs BIGINT NOT NULL CONSTRAINT df_gpqsa_duration DEFAULT 0,
    synthetic BIT NOT NULL CONSTRAINT df_gpqsa_synthetic DEFAULT 0
);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpqsa_quote_id' AND object_id = OBJECT_ID('group_pricing_quote_status_audits'))
  CREATE INDEX idx_gpqsa_quote_id ON group_pricing_quote_status_audits(quote_id);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpqsa_changed_by' AND object_id = OBJECT_ID('group_pricing_quote_status_audits'))
  CREATE INDEX idx_gpqsa_changed_by ON group_pricing_quote_status_audits(changed_by);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpqsa_changed_at' AND object_id = OBJECT_ID('group_pricing_quote_status_audits'))
  CREATE INDEX idx_gpqsa_changed_at ON group_pricing_quote_status_audits(changed_at);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpqsa_quote_changed' AND object_id = OBJECT_ID('group_pricing_quote_status_audits'))
  CREATE INDEX idx_gpqsa_quote_changed ON group_pricing_quote_status_audits(quote_id, changed_at);

-- 4) Configurable SLA targets per (from_status, to_status, quote_type).
IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'quote_sla_targets')
CREATE TABLE quote_sla_targets (
    id INT IDENTITY(1,1) PRIMARY KEY,
    from_status NVARCHAR(50) NOT NULL,
    to_status NVARCHAR(50) NOT NULL,
    target_hours FLOAT NOT NULL,
    warning_pct_of_sla FLOAT NOT NULL CONSTRAINT df_qst_warning_pct DEFAULT 0.8,
    quote_type NVARCHAR(30) NOT NULL CONSTRAINT df_qst_quote_type DEFAULT '',
    active BIT NOT NULL CONSTRAINT df_qst_active DEFAULT 1,
    updated_by NVARCHAR(255) NULL,
    updated_at DATETIME NOT NULL CONSTRAINT df_qst_updated_at DEFAULT GETDATE()
);

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'ux_qst_pair' AND object_id = OBJECT_ID('quote_sla_targets'))
  CREATE UNIQUE INDEX ux_qst_pair ON quote_sla_targets(from_status, to_status, quote_type);

-- Seed defaults (industry-typical turnaround targets). Guard each so re-runs
-- and admin edits aren't clobbered.
IF NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'draft' AND to_status = 'submitted' AND quote_type = '')
  INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
  VALUES ('draft', 'submitted', 48, '', 'system-seed');

IF NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'submitted' AND to_status = 'approved' AND quote_type = '')
  INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
  VALUES ('submitted', 'approved', 24, '', 'system-seed');

IF NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'approved' AND to_status = 'accepted' AND quote_type = '')
  INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
  VALUES ('approved', 'accepted', 72, '', 'system-seed');

IF NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'accepted' AND to_status = 'in_force' AND quote_type = '')
  INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
  VALUES ('accepted', 'in_force', 168, '', 'system-seed');

-- 5) Backfill per-status timestamps from existing creation/modification dates
-- so historical quotes show up on the dashboard. Guarded by IS NULL so the
-- backfill is idempotent and does not overwrite a real recorded timestamp.
UPDATE group_pricing_quotes
   SET submitted_at = creation_date
 WHERE submitted_at IS NULL
   AND status IN ('submitted', 'approved', 'rejected', 'accepted', 'in_force');

UPDATE group_pricing_quotes
   SET approved_at = modification_date
 WHERE approved_at IS NULL
   AND status IN ('approved', 'accepted', 'in_force');

UPDATE group_pricing_quotes
   SET accepted_at = modification_date
 WHERE accepted_at IS NULL
   AND status IN ('accepted', 'in_force');

UPDATE group_pricing_quotes
   SET rejected_at = modification_date
 WHERE rejected_at IS NULL
   AND status = 'rejected';

UPDATE group_pricing_quotes
   SET in_force_at = commencement_date
 WHERE in_force_at IS NULL
   AND status = 'in_force'
   AND commencement_date IS NOT NULL;

-- 6) Backfill one synthetic audit row per existing quote so the funnel chart
-- and per-user counts have something to render before any new transitions
-- occur. Marked synthetic = 1 so SLA-breach queries can exclude them.
INSERT INTO group_pricing_quote_status_audits
    (quote_id, old_status, new_status, status_message, changed_by, changed_at, duration_from_prev_secs, synthetic)
SELECT q.id,
       NULL,
       q.status,
       'Backfilled from current state',
       COALESCE(NULLIF(q.modified_by, ''), q.created_by),
       COALESCE(q.modification_date, q.creation_date),
       0,
       1
  FROM group_pricing_quotes q
 WHERE q.status IS NOT NULL
   AND NOT EXISTS (
        SELECT 1 FROM group_pricing_quote_status_audits a WHERE a.quote_id = q.id
   );

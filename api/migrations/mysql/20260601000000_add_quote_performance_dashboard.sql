-- Quote Performance Dashboard: adds per-status milestone timestamps to
-- group_pricing_quotes, a dedicated audit table for quote status
-- transitions, an admin-editable SLA targets table with seed defaults,
-- and a best-effort backfill of historical data so the dashboard renders
-- meaningful numbers from day one. Idempotent on re-runs.

-- 1) Per-status milestone timestamps on group_pricing_quotes.
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'submitted_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN submitted_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'approved_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN approved_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'rejected_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN rejected_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'accepted_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN accepted_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'in_force_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN in_force_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'rejected_reason');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN rejected_reason VARCHAR(500) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 1a) GORM's default mapping for `string` fields is LONGTEXT, which MySQL
-- refuses to index without an explicit key length. The composite index
-- below references created_by and status, so widen both to sized
-- VARCHARs first (idempotent on re-runs). See the equivalent guard in
-- 20260521000000_add_group_scheme_exposures_indexes.sql.
SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'group_pricing_quotes'
                       AND column_name = 'created_by'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE group_pricing_quotes MODIFY COLUMN created_by VARCHAR(255) NULL',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @needs_alter := (SELECT COUNT(*) FROM information_schema.columns
                     WHERE table_schema = DATABASE()
                       AND table_name = 'group_pricing_quotes'
                       AND column_name = 'status'
                       AND LOWER(data_type) IN ('text', 'longtext', 'mediumtext', 'tinytext'));
SET @sql := IF(@needs_alter > 0,
  'ALTER TABLE group_pricing_quotes MODIFY COLUMN status VARCHAR(50) NULL',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 2) Indexes for dashboard date-range and per-user aggregation queries.
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND index_name = 'idx_gpq_submitted_at');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpq_submitted_at ON group_pricing_quotes (submitted_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND index_name = 'idx_gpq_approved_at');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpq_approved_at ON group_pricing_quotes (approved_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND index_name = 'idx_gpq_accepted_at');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpq_accepted_at ON group_pricing_quotes (accepted_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND index_name = 'idx_gpq_created_by_status');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpq_created_by_status ON group_pricing_quotes (created_by, status)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND index_name = 'idx_gpq_created_by_status_submitted_at');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpq_created_by_status_submitted_at ON group_pricing_quotes (created_by, status, submitted_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 3) Quote status audit table (mirrors group_scheme_status_audits).
CREATE TABLE IF NOT EXISTS group_pricing_quote_status_audits (
    id INT AUTO_INCREMENT PRIMARY KEY,
    quote_id INT NOT NULL,
    old_status VARCHAR(50) NULL,
    new_status VARCHAR(50) NULL,
    status_message VARCHAR(500) NULL,
    changed_by VARCHAR(255) NULL,
    changed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    duration_from_prev_secs BIGINT NOT NULL DEFAULT 0,
    synthetic TINYINT(1) NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quote_status_audits' AND index_name = 'idx_gpqsa_quote_id');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpqsa_quote_id ON group_pricing_quote_status_audits (quote_id)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quote_status_audits' AND index_name = 'idx_gpqsa_changed_by');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpqsa_changed_by ON group_pricing_quote_status_audits (changed_by)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quote_status_audits' AND index_name = 'idx_gpqsa_changed_at');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpqsa_changed_at ON group_pricing_quote_status_audits (changed_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quote_status_audits' AND index_name = 'idx_gpqsa_quote_changed');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpqsa_quote_changed ON group_pricing_quote_status_audits (quote_id, changed_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 4) Configurable SLA targets per (from_status, to_status, quote_type).
CREATE TABLE IF NOT EXISTS quote_sla_targets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_status VARCHAR(50) NOT NULL,
    to_status VARCHAR(50) NOT NULL,
    target_hours DOUBLE NOT NULL,
    warning_pct_of_sla DOUBLE NOT NULL DEFAULT 0.8,
    quote_type VARCHAR(30) NOT NULL DEFAULT '',
    active TINYINT(1) NOT NULL DEFAULT 1,
    updated_by VARCHAR(255) NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'quote_sla_targets' AND index_name = 'ux_qst_pair');
SET @sql := IF(@idx = 0, 'CREATE UNIQUE INDEX ux_qst_pair ON quote_sla_targets (from_status, to_status, quote_type)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Seed defaults guarded per-row so re-runs and admin edits aren't clobbered.
INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
SELECT 'draft', 'submitted', 48, '', 'system-seed'
 WHERE NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'draft' AND to_status = 'submitted' AND quote_type = '');

INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
SELECT 'submitted', 'approved', 24, '', 'system-seed'
 WHERE NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'submitted' AND to_status = 'approved' AND quote_type = '');

INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
SELECT 'approved', 'accepted', 72, '', 'system-seed'
 WHERE NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'approved' AND to_status = 'accepted' AND quote_type = '');

INSERT INTO quote_sla_targets (from_status, to_status, target_hours, quote_type, updated_by)
SELECT 'accepted', 'in_force', 168, '', 'system-seed'
 WHERE NOT EXISTS (SELECT 1 FROM quote_sla_targets WHERE from_status = 'accepted' AND to_status = 'in_force' AND quote_type = '');

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

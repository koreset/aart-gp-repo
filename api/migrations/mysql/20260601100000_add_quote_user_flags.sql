-- Quote User Flags: lets managers flag a user on the dashboard for
-- coaching or capacity attention with an internal note, and resolve the
-- flag later with a resolution note. Append-only — history matters
-- (knowing "Alice has been flagged for capacity 3× this quarter" is the
-- point). At most one open flag per (user_name, flag_reason) is enforced
-- in the service layer, not via a partial index (MySQL doesn't support
-- partial unique indexes). Idempotent on re-runs.

SET @tbl := (SELECT COUNT(*) FROM information_schema.tables
             WHERE table_schema = DATABASE() AND table_name = 'quote_user_flags');
SET @sql := IF(@tbl = 0,
'CREATE TABLE quote_user_flags (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) NULL,
    flag_reason VARCHAR(50) NOT NULL,
    note TEXT NOT NULL,
    opened_by VARCHAR(255) NOT NULL,
    opened_by_name VARCHAR(255) NOT NULL,
    opened_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_by VARCHAR(255) NULL,
    resolved_by_name VARCHAR(255) NULL,
    resolved_at DATETIME NULL,
    resolution_note TEXT NULL
)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'quote_user_flags' AND index_name = 'idx_quf_user');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_quf_user ON quote_user_flags(user_name)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'quote_user_flags' AND index_name = 'idx_quf_opened_at');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_quf_opened_at ON quote_user_flags(opened_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- MySQL has no partial indexes — a regular composite covers the same
-- queries (a small read amplification when many rows are resolved).
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'quote_user_flags' AND index_name = 'idx_quf_user_reason');
SET @sql := IF(@idx = 0,
  'CREATE INDEX idx_quf_user_reason ON quote_user_flags(user_name, flag_reason, resolved_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Treaty activation/deactivation accountability: adds activated_by /
-- activated_at / deactivated_by / deactivated_at to reinsurance_treaties so
-- the Treaty Management screen can record who promoted a treaty out of draft
-- and who later took it out of service. Activated treaties cannot be deleted
-- (kept for traceability); deactivation transitions status to 'cancelled' and
-- stamps the deactivated_* fields.
-- Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'reinsurance_treaties' AND column_name = 'activated_by');
SET @sql := IF(@col = 0, 'ALTER TABLE reinsurance_treaties ADD COLUMN activated_by VARCHAR(255) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'reinsurance_treaties' AND column_name = 'activated_at');
SET @sql := IF(@col = 0, 'ALTER TABLE reinsurance_treaties ADD COLUMN activated_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'reinsurance_treaties' AND column_name = 'deactivated_by');
SET @sql := IF(@col = 0, 'ALTER TABLE reinsurance_treaties ADD COLUMN deactivated_by VARCHAR(255) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'reinsurance_treaties' AND column_name = 'deactivated_at');
SET @sql := IF(@col = 0, 'ALTER TABLE reinsurance_treaties ADD COLUMN deactivated_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

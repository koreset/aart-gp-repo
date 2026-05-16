-- Quote close-out outcomes: NTU + Declined.
--
-- Adds per-status milestone timestamps + reason text to group_pricing_quotes
-- so an approved quote can be closed out with one of three terminal
-- states: accepted (existing), not_taken_up (broker / client went
-- elsewhere), or declined (deal fell through on the client side).
-- Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'not_taken_up_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN not_taken_up_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'not_taken_up_reason');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN not_taken_up_reason VARCHAR(500) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'declined_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN declined_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'declined_reason');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN declined_reason VARCHAR(500) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- MySQL has no partial indexes — full indexes serve the same query path
-- (NULL rows still take an index slot but the cost is tiny).
SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND index_name = 'idx_gpq_not_taken_up_at');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpq_not_taken_up_at ON group_pricing_quotes(not_taken_up_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND index_name = 'idx_gpq_declined_at');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gpq_declined_at ON group_pricing_quotes(declined_at)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

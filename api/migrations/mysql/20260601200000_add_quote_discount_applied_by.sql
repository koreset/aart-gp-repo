-- Quote discount accountability: adds discount_applied_by /
-- discount_applied_at to group_pricing_quotes so the quote list can
-- show who set the current discount. Matches the existing pattern of
-- created_by / modified_by / approved_by on the same table.
-- Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'discount_applied_by');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN discount_applied_by VARCHAR(255) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_pricing_quotes' AND column_name = 'discount_applied_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_pricing_quotes ADD COLUMN discount_applied_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

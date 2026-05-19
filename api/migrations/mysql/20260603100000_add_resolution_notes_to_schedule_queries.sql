-- Add resolution_notes column to claim_payment_schedule_queries so finance can
-- type a reply when resolving a claims follow-up or line query. Existing rows
-- get NULL (treated as empty by the model). Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'claim_payment_schedule_queries' AND column_name = 'resolution_notes');
SET @sql := IF(@col = 0, 'ALTER TABLE claim_payment_schedule_queries ADD COLUMN resolution_notes TEXT NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

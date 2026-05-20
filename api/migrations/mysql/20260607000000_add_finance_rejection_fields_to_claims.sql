-- Snapshot columns powering the finance-rejection banner on the claim
-- detail view. Populated when a payment-schedule line is rejected; cleared
-- on the next assessor re-approval. Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND column_name = 'finance_rejected_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claims ADD COLUMN finance_rejected_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND column_name = 'finance_rejected_by');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claims ADD COLUMN finance_rejected_by VARCHAR(255) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND column_name = 'finance_rejection_reason_code');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claims ADD COLUMN finance_rejection_reason_code VARCHAR(64) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND column_name = 'finance_rejection_notes');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claims ADD COLUMN finance_rejection_notes TEXT NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND column_name = 'finance_rejection_schedule_number');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claims ADD COLUMN finance_rejection_schedule_number VARCHAR(191) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Adds the pre-authorisation amount-drift snapshot + resolution audit to
-- payment schedule items. ApprovedAmountSnapshot freezes the assessor's
-- approved figure at schedule-generation time so finance can compare it
-- against the snapshotted Gross even if the underlying assessment is later
-- edited. AmountDriftResolved is the explicit acknowledgement gate that
-- prevents first authorisation while a non-zero drift remains unreviewed.
-- Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'claim_payment_schedule_items' AND column_name = 'approved_amount_snapshot');
SET @sql := IF(@col = 0, 'ALTER TABLE claim_payment_schedule_items ADD COLUMN approved_amount_snapshot DECIMAL(18,2) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'claim_payment_schedule_items' AND column_name = 'amount_drift_resolved');
SET @sql := IF(@col = 0, 'ALTER TABLE claim_payment_schedule_items ADD COLUMN amount_drift_resolved TINYINT(1) NOT NULL DEFAULT 0', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'claim_payment_schedule_items' AND column_name = 'amount_drift_resolved_by');
SET @sql := IF(@col = 0, 'ALTER TABLE claim_payment_schedule_items ADD COLUMN amount_drift_resolved_by VARCHAR(255) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'claim_payment_schedule_items' AND column_name = 'amount_drift_resolved_at');
SET @sql := IF(@col = 0, 'ALTER TABLE claim_payment_schedule_items ADD COLUMN amount_drift_resolved_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

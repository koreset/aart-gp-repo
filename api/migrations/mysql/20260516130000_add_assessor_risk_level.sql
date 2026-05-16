-- Adds assessor_risk_level to group_scheme_claim_assessments so the assessor
-- can record their own opinion alongside the system-determined fraud_risk_level.
-- The GLM trains on assessor_risk_level (human judgement), not its own past
-- output. Idempotent.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claim_assessments' AND column_name = 'assessor_risk_level');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claim_assessments ADD COLUMN assessor_risk_level VARCHAR(32) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

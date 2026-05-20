-- Adds frozen approval audit fields to claim assessments. ApprovedAmount is
-- the figure the assessor explicitly signed off on at the moment the
-- assessment outcome moved to "approved" / "partial_approval"; the schedule
-- pre-authorisation drift check uses this as the source of truth instead of
-- the editable RecommendedAmount. Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claim_assessments' AND column_name = 'approved_amount');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claim_assessments ADD COLUMN approved_amount DECIMAL(18,2) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claim_assessments' AND column_name = 'approved_at');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claim_assessments ADD COLUMN approved_at DATETIME NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claim_assessments' AND column_name = 'approved_by');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claim_assessments ADD COLUMN approved_by VARCHAR(255) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Backfill: any existing assessment whose outcome already shows the claim was
-- approved gets the recommended amount frozen as the approved amount so the
-- drift check has something to compare against on historical data.
UPDATE group_scheme_claim_assessments
SET approved_amount = recommended_amount,
    approved_at = updated_at,
    approved_by = COALESCE(created_by, 'system')
WHERE approved_amount IS NULL
  AND LOWER(COALESCE(assessment_outcome, '')) IN ('approved','partial_approval','partially_approved');

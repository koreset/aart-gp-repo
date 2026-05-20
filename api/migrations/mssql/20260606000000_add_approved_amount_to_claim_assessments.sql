-- Adds frozen approval audit fields to claim assessments. ApprovedAmount is
-- the figure the assessor explicitly signed off on at the moment the
-- assessment outcome moved to "approved" / "partial_approval"; the schedule
-- pre-authorisation drift check uses this as the source of truth instead of
-- the editable RecommendedAmount. Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claim_assessments') AND name = 'approved_amount')
  ALTER TABLE group_scheme_claim_assessments ADD approved_amount DECIMAL(18,2) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claim_assessments') AND name = 'approved_at')
  ALTER TABLE group_scheme_claim_assessments ADD approved_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claim_assessments') AND name = 'approved_by')
  ALTER TABLE group_scheme_claim_assessments ADD approved_by NVARCHAR(255) NULL;

UPDATE group_scheme_claim_assessments
SET approved_amount = recommended_amount,
    approved_at = updated_at,
    approved_by = COALESCE(created_by, 'system')
WHERE approved_amount IS NULL
  AND LOWER(ISNULL(assessment_outcome, '')) IN ('approved','partial_approval','partially_approved');

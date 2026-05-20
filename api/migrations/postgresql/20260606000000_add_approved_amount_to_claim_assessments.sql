-- Adds frozen approval audit fields to claim assessments. ApprovedAmount is
-- the figure the assessor explicitly signed off on at the moment the
-- assessment outcome moved to "approved" / "partial_approval"; the schedule
-- pre-authorisation drift check uses this as the source of truth instead of
-- the editable RecommendedAmount. Idempotent on re-runs.

ALTER TABLE group_scheme_claim_assessments
    ADD COLUMN IF NOT EXISTS approved_amount NUMERIC(18,2) NULL;

ALTER TABLE group_scheme_claim_assessments
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP NULL;

ALTER TABLE group_scheme_claim_assessments
    ADD COLUMN IF NOT EXISTS approved_by VARCHAR(255) NULL;

UPDATE group_scheme_claim_assessments
SET approved_amount = recommended_amount,
    approved_at = updated_at,
    approved_by = COALESCE(created_by, 'system')
WHERE approved_amount IS NULL
  AND LOWER(COALESCE(assessment_outcome, '')) IN ('approved','partial_approval','partially_approved');

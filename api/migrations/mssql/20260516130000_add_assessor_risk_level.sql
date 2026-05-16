-- Adds assessor_risk_level to group_scheme_claim_assessments so the assessor
-- can record their own opinion alongside the system-determined fraud_risk_level.
-- The GLM trains on assessor_risk_level (human judgement), not its own past
-- output. Idempotent.

IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('group_scheme_claim_assessments') AND name = 'assessor_risk_level')
    ALTER TABLE group_scheme_claim_assessments ADD assessor_risk_level NVARCHAR(32) NULL;

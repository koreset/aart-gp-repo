-- Adds the pre-authorisation amount-drift snapshot + resolution audit to
-- payment schedule items. ApprovedAmountSnapshot freezes the assessor's
-- approved figure at schedule-generation time so finance can compare it
-- against the snapshotted Gross even if the underlying assessment is later
-- edited. AmountDriftResolved is the explicit acknowledgement gate that
-- prevents first authorisation while a non-zero drift remains unreviewed.
-- Idempotent on re-runs.

ALTER TABLE claim_payment_schedule_items
    ADD COLUMN IF NOT EXISTS approved_amount_snapshot NUMERIC(18,2) NULL;

ALTER TABLE claim_payment_schedule_items
    ADD COLUMN IF NOT EXISTS amount_drift_resolved BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE claim_payment_schedule_items
    ADD COLUMN IF NOT EXISTS amount_drift_resolved_by VARCHAR(255) NULL;

ALTER TABLE claim_payment_schedule_items
    ADD COLUMN IF NOT EXISTS amount_drift_resolved_at TIMESTAMP NULL;

-- Adds the pre-authorisation amount-drift snapshot + resolution audit to
-- payment schedule items. ApprovedAmountSnapshot freezes the assessor's
-- approved figure at schedule-generation time so finance can compare it
-- against the snapshotted Gross even if the underlying assessment is later
-- edited. AmountDriftResolved is the explicit acknowledgement gate that
-- prevents first authorisation while a non-zero drift remains unreviewed.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('claim_payment_schedule_items') AND name = 'approved_amount_snapshot')
  ALTER TABLE claim_payment_schedule_items ADD approved_amount_snapshot DECIMAL(18,2) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('claim_payment_schedule_items') AND name = 'amount_drift_resolved')
  ALTER TABLE claim_payment_schedule_items ADD amount_drift_resolved BIT NOT NULL CONSTRAINT df_cpsi_amt_drift_resolved DEFAULT 0;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('claim_payment_schedule_items') AND name = 'amount_drift_resolved_by')
  ALTER TABLE claim_payment_schedule_items ADD amount_drift_resolved_by NVARCHAR(255) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('claim_payment_schedule_items') AND name = 'amount_drift_resolved_at')
  ALTER TABLE claim_payment_schedule_items ADD amount_drift_resolved_at DATETIME NULL;

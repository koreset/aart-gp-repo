-- Adds treaty_basis to reinsurance_treaties so the claim-aware treaty selector
-- can decide which date to compare against the treaty's effective/expiry
-- window: risk_attaching → member entry date, loss_occurring → claim date of
-- event. Default backfills every existing row to 'risk_attaching' (the SA
-- group-life norm) so callers always read a non-empty value.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('reinsurance_treaties') AND name = 'treaty_basis')
  ALTER TABLE reinsurance_treaties ADD treaty_basis NVARCHAR(32) NOT NULL CONSTRAINT DF_reinsurance_treaties_treaty_basis DEFAULT 'risk_attaching';

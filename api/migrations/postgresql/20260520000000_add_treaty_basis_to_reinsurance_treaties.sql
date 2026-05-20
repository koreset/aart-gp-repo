-- Adds treaty_basis to reinsurance_treaties so the claim-aware treaty selector
-- can decide which date to compare against the treaty's effective/expiry
-- window: risk_attaching → member entry date, loss_occurring → claim date of
-- event. Default backfills every existing row to 'risk_attaching' (the SA
-- group-life norm) so callers always read a non-empty value.
-- Idempotent on re-runs.

ALTER TABLE reinsurance_treaties
    ADD COLUMN IF NOT EXISTS treaty_basis VARCHAR(32) NOT NULL DEFAULT 'risk_attaching';

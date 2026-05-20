-- Covering indexes for the cross-claim duplicate check. The new
-- pre-authorisation suite scans group_scheme_claims by (claimant_id_number,
-- status) and (bank_account_number, status) at schedule-generation time. These
-- indexes keep the lookup constant-time even when the paid-claim history
-- grows. Idempotent on re-runs.

CREATE INDEX IF NOT EXISTS idx_gsc_claimant_id_status
    ON group_scheme_claims (claimant_id_number, status);

CREATE INDEX IF NOT EXISTS idx_gsc_bank_account_status
    ON group_scheme_claims (bank_account_number, status);

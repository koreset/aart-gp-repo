-- Snapshot columns powering the finance-rejection banner on the claim
-- detail view. Populated when a payment-schedule line is rejected; cleared
-- on the next assessor re-approval. Idempotent on re-runs.

ALTER TABLE group_scheme_claims
    ADD COLUMN IF NOT EXISTS finance_rejected_at TIMESTAMP NULL;

ALTER TABLE group_scheme_claims
    ADD COLUMN IF NOT EXISTS finance_rejected_by VARCHAR(255) NULL;

ALTER TABLE group_scheme_claims
    ADD COLUMN IF NOT EXISTS finance_rejection_reason_code VARCHAR(64) NULL;

ALTER TABLE group_scheme_claims
    ADD COLUMN IF NOT EXISTS finance_rejection_notes TEXT NULL;

ALTER TABLE group_scheme_claims
    ADD COLUMN IF NOT EXISTS finance_rejection_schedule_number VARCHAR(191) NULL;

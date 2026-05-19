-- Add resolution_notes column to claim_payment_schedule_queries so finance can
-- type a reply when resolving a claims follow-up or line query. Existing rows
-- get NULL (treated as empty by the model). Idempotent on re-runs.

ALTER TABLE claim_payment_schedule_queries
    ADD COLUMN IF NOT EXISTS resolution_notes TEXT;

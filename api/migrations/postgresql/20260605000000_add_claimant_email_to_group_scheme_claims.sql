-- Adds claimant_email so the claim payment confirmation letter can be sent
-- by email through the existing email outbox infrastructure. Phone (kept in
-- claimant_contact_number) is retained for future SMS / WhatsApp channels.
-- Idempotent on re-runs.

ALTER TABLE group_scheme_claims
    ADD COLUMN IF NOT EXISTS claimant_email VARCHAR(255) NULL;

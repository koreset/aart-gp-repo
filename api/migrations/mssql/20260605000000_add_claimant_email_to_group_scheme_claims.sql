-- Adds claimant_email so the claim payment confirmation letter can be sent
-- by email through the existing email outbox infrastructure. Phone (kept in
-- claimant_contact_number) is retained for future SMS / WhatsApp channels.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claims') AND name = 'claimant_email')
  ALTER TABLE group_scheme_claims ADD claimant_email NVARCHAR(255) NULL;

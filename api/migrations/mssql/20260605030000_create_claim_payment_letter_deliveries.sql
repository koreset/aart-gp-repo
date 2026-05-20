-- Channel-agnostic delivery log for claim payment letters. Channel is kept as
-- a free-text string so SMS / WhatsApp slot in later without a schema change.
-- When channel = 'email', outbox_id points at the row queued in email_outbox.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'claim_payment_letter_deliveries')
BEGIN
  CREATE TABLE claim_payment_letter_deliveries (
    id           BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    letter_id    BIGINT NOT NULL,
    channel      NVARCHAR(16) NOT NULL,
    recipient    NVARCHAR(255) NULL,
    status       NVARCHAR(16) NOT NULL,
    provider_ref NVARCHAR(255) NULL,
    outbox_id    BIGINT NULL,
    error        NVARCHAR(MAX) NULL,
    sent_by      NVARCHAR(255) NULL,
    sent_at      DATETIME NULL,
    created_at   DATETIME NULL
  );
  CREATE INDEX idx_letter_deliveries_letter_id ON claim_payment_letter_deliveries (letter_id);
END

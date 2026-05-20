-- Channel-agnostic delivery log for claim payment letters. Channel is kept as
-- a free-text string so SMS / WhatsApp slot in later without a schema change.
-- When channel = 'email', outbox_id points at the row queued in email_outbox.
-- Idempotent on re-runs.

CREATE TABLE IF NOT EXISTS claim_payment_letter_deliveries (
  id           BIGSERIAL PRIMARY KEY,
  letter_id    BIGINT NOT NULL,
  channel      VARCHAR(16) NOT NULL,
  recipient    VARCHAR(255) NULL,
  status       VARCHAR(16) NOT NULL,
  provider_ref VARCHAR(255) NULL,
  outbox_id    BIGINT NULL,
  error        TEXT NULL,
  sent_by      VARCHAR(255) NULL,
  sent_at      TIMESTAMP NULL,
  created_at   TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_letter_deliveries_letter_id
    ON claim_payment_letter_deliveries (letter_id);

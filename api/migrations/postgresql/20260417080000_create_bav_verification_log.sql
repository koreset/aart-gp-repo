-- Migration: create bav_verification_logs
-- Purpose: audit trail + idempotency dedup for Bank Account Verification
-- calls made through services/bav. See docs/bav-provider-abstraction-plan.md
-- (Phase 5).

CREATE TABLE IF NOT EXISTS bav_verification_logs (
    id                   SERIAL PRIMARY KEY,
    claim_id             INTEGER NULL,
    provider             VARCHAR(64) NOT NULL DEFAULT '',
    provider_request_id  VARCHAR(128) NOT NULL DEFAULT '',
    idempotency_key      VARCHAR(128) NOT NULL DEFAULT '',
    status               VARCHAR(32) NOT NULL DEFAULT '',
    request_payload      TEXT NULL,
    response_payload     TEXT NULL,
    error_message        VARCHAR(1024) NOT NULL DEFAULT '',
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bav_log_key_created ON bav_verification_logs (idempotency_key, created_at);
CREATE INDEX IF NOT EXISTS idx_bav_log_claim_id ON bav_verification_logs (claim_id);
CREATE INDEX IF NOT EXISTS idx_bav_log_provider ON bav_verification_logs (provider);

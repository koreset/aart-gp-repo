-- Migration: create bav_verification_logs
-- Purpose: audit trail + idempotency dedup for Bank Account Verification
-- calls made through services/bav. See docs/bav-provider-abstraction-plan.md
-- (Phase 5).

CREATE TABLE IF NOT EXISTS bav_verification_logs (
    id                   INT AUTO_INCREMENT PRIMARY KEY,
    claim_id             INT NULL,
    provider             VARCHAR(64) NOT NULL DEFAULT '',
    provider_request_id  VARCHAR(128) NOT NULL DEFAULT '',
    idempotency_key      VARCHAR(128) NOT NULL DEFAULT '',
    status               VARCHAR(32) NOT NULL DEFAULT '',
    request_payload      TEXT NULL,
    response_payload     TEXT NULL,
    error_message        VARCHAR(1024) NOT NULL DEFAULT '',
    created_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- idx_bav_log_key_created serves the 24h dedup lookup (key + time cutoff).
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='bav_verification_logs' AND INDEX_NAME='idx_bav_log_key_created' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE INDEX idx_bav_log_key_created ON bav_verification_logs(idempotency_key, created_at);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='bav_verification_logs' AND INDEX_NAME='idx_bav_log_claim_id' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE INDEX idx_bav_log_claim_id ON bav_verification_logs(claim_id);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='bav_verification_logs' AND INDEX_NAME='idx_bav_log_provider' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE INDEX idx_bav_log_provider ON bav_verification_logs(provider);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

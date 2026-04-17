-- Migration: create bav_verification_logs
-- Purpose: audit trail + idempotency dedup for Bank Account Verification
-- calls made through services/bav. See docs/bav-provider-abstraction-plan.md
-- (Phase 5).

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'bav_verification_logs')
BEGIN
    CREATE TABLE bav_verification_logs (
        id                  INT IDENTITY(1,1) PRIMARY KEY,
        claim_id            INT NULL,
        provider            NVARCHAR(64) NOT NULL DEFAULT '',
        provider_request_id NVARCHAR(128) NOT NULL DEFAULT '',
        idempotency_key     NVARCHAR(128) NOT NULL DEFAULT '',
        status              NVARCHAR(32) NOT NULL DEFAULT '',
        request_payload     NVARCHAR(MAX) NULL,
        response_payload    NVARCHAR(MAX) NULL,
        error_message       NVARCHAR(1024) NOT NULL DEFAULT '',
        created_at          DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME()
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bav_log_key_created' AND object_id = OBJECT_ID('bav_verification_logs'))
BEGIN
    CREATE INDEX idx_bav_log_key_created ON bav_verification_logs(idempotency_key, created_at);
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bav_log_claim_id' AND object_id = OBJECT_ID('bav_verification_logs'))
BEGIN
    CREATE INDEX idx_bav_log_claim_id ON bav_verification_logs(claim_id);
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bav_log_provider' AND object_id = OBJECT_ID('bav_verification_logs'))
BEGIN
    CREATE INDEX idx_bav_log_provider ON bav_verification_logs(provider);
END;

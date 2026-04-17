-- Migration: add reinsurer response fields to large_claim_notices
-- Purpose: P2-5 — capture reinsurer's decision on a large-claim cession.

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='large_claim_notices' AND COLUMN_NAME='response_status')
BEGIN
    ALTER TABLE large_claim_notices ADD response_status NVARCHAR(32) NOT NULL DEFAULT '';
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='large_claim_notices' AND COLUMN_NAME='responded_at')
BEGIN
    ALTER TABLE large_claim_notices ADD responded_at DATETIME2 NULL;
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='large_claim_notices' AND COLUMN_NAME='responded_by')
BEGIN
    ALTER TABLE large_claim_notices ADD responded_by NVARCHAR(255) NOT NULL DEFAULT '';
END;

IF NOT EXISTS(SELECT * FROM sys.indexes WHERE name='idx_lcn_response_status' AND object_id=OBJECT_ID('large_claim_notices'))
BEGIN
    CREATE INDEX idx_lcn_response_status ON large_claim_notices (response_status);
END;

-- Migration: add reinsurer response fields to large_claim_notices
-- Purpose: P2-5 — capture reinsurer's decision on a large-claim cession.

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='large_claim_notices' AND COLUMN_NAME='response_status' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE large_claim_notices ADD COLUMN response_status VARCHAR(32) NOT NULL DEFAULT "";'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='large_claim_notices' AND COLUMN_NAME='responded_at' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE large_claim_notices ADD COLUMN responded_at DATETIME NULL;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='large_claim_notices' AND COLUMN_NAME='responded_by' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE large_claim_notices ADD COLUMN responded_by VARCHAR(255) NOT NULL DEFAULT "";'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='large_claim_notices' AND INDEX_NAME='idx_lcn_response_status' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE INDEX idx_lcn_response_status ON large_claim_notices (response_status);'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

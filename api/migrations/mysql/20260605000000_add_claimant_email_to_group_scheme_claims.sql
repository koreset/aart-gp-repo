-- Adds claimant_email so the claim payment confirmation letter can be sent
-- by email through the existing email outbox infrastructure. Phone (kept in
-- claimant_contact_number) is retained for future SMS / WhatsApp channels.
-- Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND column_name = 'claimant_email');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claims ADD COLUMN claimant_email VARCHAR(255) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

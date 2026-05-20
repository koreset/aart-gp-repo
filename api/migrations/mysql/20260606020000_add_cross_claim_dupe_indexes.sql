-- Covering indexes for the cross-claim duplicate check. The new
-- pre-authorisation suite scans group_scheme_claims by (claimant_id_number,
-- status) and (bank_account_number, status) at schedule-generation time. These
-- indexes keep the lookup constant-time even when the paid-claim history
-- grows. Idempotent on re-runs.

-- claimant_id_number, bank_account_number, and status are all TEXT in this
-- install (GORM's default for `string` without a size tag). MySQL requires a
-- key prefix length when indexing TEXT/BLOB columns. 64 chars covers any
-- realistic ID/account number, 32 covers every status value in use.

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND index_name = 'idx_gsc_claimant_id_status');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gsc_claimant_id_status ON group_scheme_claims (claimant_id_number(64), status(32))', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND index_name = 'idx_gsc_bank_account_status');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_gsc_bank_account_status ON group_scheme_claims (bank_account_number(64), status(32))', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

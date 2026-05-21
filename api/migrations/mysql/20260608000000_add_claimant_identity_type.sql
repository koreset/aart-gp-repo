-- Adds claimant_identity_type to group_scheme_claims so the system records
-- whether the claimant's identity document is a SA ID number or a passport,
-- rather than inferring it from the digit count at BAV-call time. Backfills
-- existing rows using the same heuristic the form was using (13 digits →
-- IDNumber, otherwise Passport) so historical claims start with a sensible
-- value. Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'group_scheme_claims' AND column_name = 'claimant_identity_type');
SET @sql := IF(@col = 0, 'ALTER TABLE group_scheme_claims ADD COLUMN claimant_identity_type VARCHAR(16) NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Backfill: any row whose identity type is still empty gets it from the ID
-- format. Skips rows already populated so the migration is safe to re-run.
UPDATE group_scheme_claims
SET claimant_identity_type = CASE
    WHEN claimant_id_number REGEXP '^[0-9]{13}$' THEN 'IDNumber'
    WHEN claimant_id_number IS NOT NULL AND claimant_id_number <> '' THEN 'Passport'
    ELSE NULL
END
WHERE claimant_identity_type IS NULL OR claimant_identity_type = '';

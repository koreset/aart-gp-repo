-- Adds claimant_identity_type to group_scheme_claims so the system records
-- whether the claimant's identity document is a SA ID number or a passport,
-- rather than inferring it from the digit count at BAV-call time. Backfills
-- existing rows using the same heuristic the form was using (13 digits →
-- IDNumber, otherwise Passport) so historical claims start with a sensible
-- value. Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claims') AND name = 'claimant_identity_type')
  ALTER TABLE group_scheme_claims ADD claimant_identity_type NVARCHAR(16) NULL;

UPDATE group_scheme_claims
SET claimant_identity_type = CASE
    WHEN claimant_id_number LIKE '[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]'
         AND LEN(claimant_id_number) = 13 THEN 'IDNumber'
    WHEN claimant_id_number IS NOT NULL AND claimant_id_number <> '' THEN 'Passport'
    ELSE NULL
END
WHERE claimant_identity_type IS NULL OR claimant_identity_type = '';

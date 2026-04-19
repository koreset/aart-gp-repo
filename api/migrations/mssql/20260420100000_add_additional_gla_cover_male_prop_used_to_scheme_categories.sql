-- Migration: persist the actual male proportion used for each
-- Additional GLA Cover calc run on scheme_categories. Audit column so the
-- value that fed the rate is stored on the same row that was
-- (re)calculated, independent of any later re-derivation from members.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_gla_cover_male_prop_used')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_male_prop_used FLOAT NULL;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_male_prop_used FLOAT NULL;
END

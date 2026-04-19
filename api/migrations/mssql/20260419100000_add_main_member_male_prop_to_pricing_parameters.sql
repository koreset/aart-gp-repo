-- Migration: add main_member_male_prop to group_pricing_parameters.
-- Used to weight male vs female gla qx when computing additional-GLA-cover
-- age-band rates:
--   qx = male_prop * male_qx + (1 - male_prop) * female_qx
-- Default 0.5 preserves the straight-average behaviour for any pre-existing
-- rows. The scheme category can override this at calc time (manual entry or
-- derived from the uploaded member list).

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_parameters') AND name = 'main_member_male_prop')
BEGIN
    ALTER TABLE group_pricing_parameters ADD main_member_male_prop FLOAT NOT NULL DEFAULT 0.5;
END
ELSE
BEGIN
    ALTER TABLE group_pricing_parameters ALTER COLUMN main_member_male_prop FLOAT NOT NULL;
END

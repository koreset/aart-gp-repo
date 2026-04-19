-- Migration: add extended_family_male_prop to group_pricing_parameters.
-- Used to weight male vs female funeral qx when computing extended-family
-- age-band rates:
--   qx = male_prop * male_qx + (1 - male_prop) * female_qx
-- Default 0.5 preserves the previous straight-average behaviour for any
-- pre-existing rows.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_parameters') AND name = 'extended_family_male_prop')
BEGIN
    ALTER TABLE group_pricing_parameters ADD extended_family_male_prop FLOAT NOT NULL DEFAULT 0.5;
END
ELSE
BEGIN
    ALTER TABLE group_pricing_parameters ALTER COLUMN extended_family_male_prop FLOAT NOT NULL;
END

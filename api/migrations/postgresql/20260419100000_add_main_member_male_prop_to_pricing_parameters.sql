-- Migration: add main_member_male_prop to group_pricing_parameters.
-- Used to weight male vs female gla qx when computing additional-GLA-cover
-- age-band rates:
--   qx = male_prop * male_qx + (1 - male_prop) * female_qx
-- Default 0.5 preserves the straight-average behaviour for any pre-existing
-- rows.

ALTER TABLE group_pricing_parameters
    ADD COLUMN IF NOT EXISTS main_member_male_prop DOUBLE PRECISION NOT NULL DEFAULT 0.5;

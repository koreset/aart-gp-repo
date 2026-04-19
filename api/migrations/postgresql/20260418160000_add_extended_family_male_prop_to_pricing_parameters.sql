-- Migration: add extended_family_male_prop to group_pricing_parameters.
-- Used to weight male vs female funeral qx when computing extended-family
-- age-band rates:
--   qx = male_prop * male_qx + (1 - male_prop) * female_qx
-- Default 0.5 preserves the previous straight-average behaviour for any
-- pre-existing rows.

ALTER TABLE group_pricing_parameters
    ADD COLUMN IF NOT EXISTS extended_family_male_prop DOUBLE PRECISION NOT NULL DEFAULT 0.5;

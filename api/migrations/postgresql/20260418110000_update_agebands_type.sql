-- Migration: age-band type
-- 1) Add `type` column to group_pricing_age_bands so multiple band sets can
--    coexist (e.g. a funeral-specific set vs a GLA-specific set).
-- 2) Add extended_family_age_band_type to scheme_categories so the UI can
--    persist which band set the extended-family cover uses.

ALTER TABLE group_pricing_age_bands ADD COLUMN IF NOT EXISTS type VARCHAR(64);
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_age_band_type VARCHAR(64);

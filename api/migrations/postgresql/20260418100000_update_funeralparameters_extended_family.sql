-- Migration: extended-family funeral benefit.
-- 1) funeral_parameters: add member_income_level + extended_family_loading.
-- 2) scheme_categories: add extended-family configuration columns.

--------------------------------------------------------------------------------
-- funeral_parameters
--------------------------------------------------------------------------------

ALTER TABLE funeral_parameters ADD COLUMN IF NOT EXISTS member_income_level INTEGER NOT NULL DEFAULT 0;
ALTER TABLE funeral_parameters ADD COLUMN IF NOT EXISTS extended_family_loading DOUBLE PRECISION NOT NULL DEFAULT 0;

--------------------------------------------------------------------------------
-- scheme_categories
--------------------------------------------------------------------------------

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_benefit BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_age_band_source VARCHAR(32);
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_custom_age_bands TEXT;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_pricing_method VARCHAR(32);
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_sums_assured TEXT;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS extended_family_band_rates TEXT;

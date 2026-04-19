-- Migration: additional GLA cover (band-based, rate-only) configuration
-- columns on scheme_categories.

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_benefit BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_age_band_source VARCHAR(16);
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_age_band_type VARCHAR(64);
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_custom_age_bands TEXT;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_band_rates TEXT;

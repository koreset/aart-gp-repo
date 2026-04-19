-- Migration: mirror Additional GLA Cover config + computed rates onto
-- member_rating_result_summaries.

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS additional_gla_cover_benefit BOOLEAN DEFAULT FALSE;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS additional_gla_cover_age_band_source VARCHAR(16);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS additional_gla_cover_age_band_type VARCHAR(64);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS additional_gla_cover_band_rates TEXT;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS additional_gla_cover_male_prop_used DOUBLE PRECISION NULL;

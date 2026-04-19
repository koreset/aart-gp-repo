-- Migration: persist extended-family configuration + computed band rates on
-- member_rating_result_summaries so the Premiums Summary UI can render the
-- per-category extended-family section from a single payload.

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS extended_family_benefit BOOLEAN DEFAULT FALSE;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS extended_family_age_band_source VARCHAR(32);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS extended_family_age_band_type VARCHAR(64);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS extended_family_pricing_method VARCHAR(32);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS extended_family_band_rates TEXT;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_extended_family_monthly_premium DOUBLE PRECISION NOT NULL DEFAULT 0;

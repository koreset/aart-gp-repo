-- Migration: mirror Additional GLA Cover config + computed rates onto
-- member_rating_result_summaries so the Premium Summary screen can render
-- per-category band rates from a single payload (matching the Extended
-- Family Funeral mirror). Rate-only product -- no aggregate premium /
-- sum-assured columns needed.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_gla_cover_benefit')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD additional_gla_cover_benefit BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN additional_gla_cover_benefit BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_gla_cover_age_band_source')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD additional_gla_cover_age_band_source NVARCHAR(16);
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN additional_gla_cover_age_band_source NVARCHAR(16);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_gla_cover_age_band_type')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD additional_gla_cover_age_band_type NVARCHAR(64);
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN additional_gla_cover_age_band_type NVARCHAR(64);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_gla_cover_band_rates')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD additional_gla_cover_band_rates NVARCHAR(MAX);
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN additional_gla_cover_band_rates NVARCHAR(MAX);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_gla_cover_male_prop_used')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD additional_gla_cover_male_prop_used FLOAT NULL;
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN additional_gla_cover_male_prop_used FLOAT NULL;
END

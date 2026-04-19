-- Migration: persist extended-family configuration + computed band rates on
-- member_rating_result_summaries so the Premiums Summary UI can render the
-- per-category extended-family section from a single payload.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'extended_family_benefit')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD extended_family_benefit BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN extended_family_benefit BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'extended_family_age_band_source')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD extended_family_age_band_source NVARCHAR(32);
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN extended_family_age_band_source NVARCHAR(32);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'extended_family_age_band_type')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD extended_family_age_band_type NVARCHAR(64);
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN extended_family_age_band_type NVARCHAR(64);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'extended_family_pricing_method')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD extended_family_pricing_method NVARCHAR(32);
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN extended_family_pricing_method NVARCHAR(32);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'extended_family_band_rates')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD extended_family_band_rates NVARCHAR(MAX);
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN extended_family_band_rates NVARCHAR(MAX);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_extended_family_monthly_premium')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD total_extended_family_monthly_premium FLOAT NOT NULL DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE member_rating_result_summaries ALTER COLUMN total_extended_family_monthly_premium FLOAT NOT NULL;
END

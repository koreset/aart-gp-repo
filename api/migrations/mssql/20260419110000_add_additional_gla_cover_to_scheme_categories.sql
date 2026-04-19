-- Migration: additional GLA cover (band-based, rate-only) configuration
-- columns on scheme_categories. When age_band_source = 'standard' the user
-- picks a band type from group_pricing_age_bands (persisted in
-- additional_gla_cover_age_band_type). Otherwise a custom JSON-persisted
-- band list is used. The male proportion is always derived from the
-- uploaded member data at rate-calc time, so no manual override is stored.
-- Band rates are computed on demand and cached in
-- additional_gla_cover_band_rates.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_gla_cover_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_benefit BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_benefit BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_gla_cover_age_band_source')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_age_band_source NVARCHAR(16);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_age_band_source NVARCHAR(16);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_gla_cover_age_band_type')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_age_band_type NVARCHAR(64);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_age_band_type NVARCHAR(64);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_gla_cover_custom_age_bands')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_custom_age_bands NVARCHAR(MAX);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_custom_age_bands NVARCHAR(MAX);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'additional_gla_cover_band_rates')
BEGIN
    ALTER TABLE scheme_categories ADD additional_gla_cover_band_rates NVARCHAR(MAX);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN additional_gla_cover_band_rates NVARCHAR(MAX);
END

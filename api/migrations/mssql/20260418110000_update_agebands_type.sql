-- Migration: age-band type
-- 1) Add `type` column to group_pricing_age_bands so multiple band sets can
--    coexist (e.g. a funeral-specific set vs a GLA-specific set).
-- 2) Add extended_family_age_band_type to scheme_categories so the UI can
--    persist which band set the extended-family cover uses.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_age_bands') AND name = 'type')
BEGIN
    ALTER TABLE group_pricing_age_bands ADD type NVARCHAR(64);
END
ELSE
BEGIN
    ALTER TABLE group_pricing_age_bands ALTER COLUMN type NVARCHAR(64);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'extended_family_age_band_type')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_age_band_type NVARCHAR(64);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_age_band_type NVARCHAR(64);
END

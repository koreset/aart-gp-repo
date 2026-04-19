-- Migration: extended-family funeral benefit.
-- 1) funeral_parameters: add member_income_level + extended_family_loading.
-- 2) scheme_categories: add extended-family configuration columns.

--------------------------------------------------------------------------------
-- funeral_parameters
--------------------------------------------------------------------------------

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('funeral_parameters') AND name = 'member_income_level')
BEGIN
    ALTER TABLE funeral_parameters ADD member_income_level INT NOT NULL DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE funeral_parameters ALTER COLUMN member_income_level INT NOT NULL;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('funeral_parameters') AND name = 'extended_family_loading')
BEGIN
    ALTER TABLE funeral_parameters ADD extended_family_loading FLOAT NOT NULL DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE funeral_parameters ALTER COLUMN extended_family_loading FLOAT NOT NULL;
END

--------------------------------------------------------------------------------
-- scheme_categories
--------------------------------------------------------------------------------

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'extended_family_benefit')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_benefit BIT DEFAULT 0;
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_benefit BIT;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'extended_family_age_band_source')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_age_band_source NVARCHAR(32);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_age_band_source NVARCHAR(32);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'extended_family_custom_age_bands')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_custom_age_bands NVARCHAR(MAX);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_custom_age_bands NVARCHAR(MAX);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'extended_family_pricing_method')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_pricing_method NVARCHAR(32);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_pricing_method NVARCHAR(32);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'extended_family_sums_assured')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_sums_assured NVARCHAR(MAX);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_sums_assured NVARCHAR(MAX);
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('scheme_categories') AND name = 'extended_family_band_rates')
BEGIN
    ALTER TABLE scheme_categories ADD extended_family_band_rates NVARCHAR(MAX);
END
ELSE
BEGIN
    ALTER TABLE scheme_categories ALTER COLUMN extended_family_band_rates NVARCHAR(MAX);
END

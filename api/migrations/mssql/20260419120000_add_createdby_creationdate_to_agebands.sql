-- Migration: add audit columns (created_by, creation_date) to
-- group_pricing_age_bands so uploaded age-band rows carry the same metadata
-- as every other rate table.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_age_bands') AND name = 'creation_date')
BEGIN
    ALTER TABLE group_pricing_age_bands ADD creation_date DATETIME2 NULL;
END

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_age_bands') AND name = 'created_by')
BEGIN
    ALTER TABLE group_pricing_age_bands ADD created_by NVARCHAR(128);
END
ELSE
BEGIN
    ALTER TABLE group_pricing_age_bands ALTER COLUMN created_by NVARCHAR(128);
END

-- Migration: age-band type
-- 1) Add `type` column to group_pricing_age_bands so multiple band sets can
--    coexist (e.g. a funeral-specific set vs a GLA-specific set).
-- 2) Add extended_family_age_band_type to scheme_categories so the UI can
--    persist which band set the extended-family cover uses.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='group_pricing_age_bands' AND COLUMN_NAME='type'),
        'ALTER TABLE group_pricing_age_bands MODIFY COLUMN type VARCHAR(64)',
        'ALTER TABLE group_pricing_age_bands ADD COLUMN type VARCHAR(64)'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='extended_family_age_band_type'),
        'ALTER TABLE scheme_categories MODIFY COLUMN extended_family_age_band_type VARCHAR(64)',
        'ALTER TABLE scheme_categories ADD COLUMN extended_family_age_band_type VARCHAR(64)'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Migration: extended-family funeral benefit.
-- 1) Add member_income_level and extended_family_loading to funeral_parameters:
--    loadings vary by the main member's income level and the dependant's age.
--    Rate-table lookup now keyed by (risk_rate_code, member_income_level,
--    age_next_birthday).
-- 2) Add extended-family configuration columns to scheme_categories so the
--    benefit choice persists with the quote.

--------------------------------------------------------------------------------
-- funeral_parameters
--------------------------------------------------------------------------------

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='funeral_parameters' AND COLUMN_NAME='member_income_level'),
        'ALTER TABLE funeral_parameters MODIFY COLUMN member_income_level INT NOT NULL DEFAULT 0;',
        'ALTER TABLE funeral_parameters ADD COLUMN member_income_level INT NOT NULL DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='funeral_parameters' AND COLUMN_NAME='extended_family_loading'),
        'ALTER TABLE funeral_parameters MODIFY COLUMN extended_family_loading DOUBLE NOT NULL DEFAULT 0;',
        'ALTER TABLE funeral_parameters ADD COLUMN extended_family_loading DOUBLE NOT NULL DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

--------------------------------------------------------------------------------
-- scheme_categories (extended-family configuration)
--------------------------------------------------------------------------------

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='extended_family_benefit'),
        'ALTER TABLE scheme_categories MODIFY COLUMN extended_family_benefit TINYINT(1) DEFAULT 0;',
        'ALTER TABLE scheme_categories ADD COLUMN extended_family_benefit TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='extended_family_age_band_source'),
        'ALTER TABLE scheme_categories MODIFY COLUMN extended_family_age_band_source VARCHAR(32);',
        'ALTER TABLE scheme_categories ADD COLUMN extended_family_age_band_source VARCHAR(32);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='extended_family_custom_age_bands'),
        'ALTER TABLE scheme_categories MODIFY COLUMN extended_family_custom_age_bands TEXT;',
        'ALTER TABLE scheme_categories ADD COLUMN extended_family_custom_age_bands TEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='extended_family_pricing_method'),
        'ALTER TABLE scheme_categories MODIFY COLUMN extended_family_pricing_method VARCHAR(32);',
        'ALTER TABLE scheme_categories ADD COLUMN extended_family_pricing_method VARCHAR(32);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='extended_family_sums_assured'),
        'ALTER TABLE scheme_categories MODIFY COLUMN extended_family_sums_assured TEXT;',
        'ALTER TABLE scheme_categories ADD COLUMN extended_family_sums_assured TEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='extended_family_band_rates'),
        'ALTER TABLE scheme_categories MODIFY COLUMN extended_family_band_rates TEXT;',
        'ALTER TABLE scheme_categories ADD COLUMN extended_family_band_rates TEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

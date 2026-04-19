-- Migration: additional GLA cover (band-based, rate-only) configuration
-- columns on scheme_categories.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_gla_cover_benefit'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_gla_cover_benefit TINYINT(1) DEFAULT 0;',
        'ALTER TABLE scheme_categories ADD COLUMN additional_gla_cover_benefit TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_gla_cover_age_band_source'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_gla_cover_age_band_source VARCHAR(16);',
        'ALTER TABLE scheme_categories ADD COLUMN additional_gla_cover_age_band_source VARCHAR(16);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_gla_cover_age_band_type'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_gla_cover_age_band_type VARCHAR(64);',
        'ALTER TABLE scheme_categories ADD COLUMN additional_gla_cover_age_band_type VARCHAR(64);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_gla_cover_custom_age_bands'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_gla_cover_custom_age_bands LONGTEXT;',
        'ALTER TABLE scheme_categories ADD COLUMN additional_gla_cover_custom_age_bands LONGTEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_gla_cover_band_rates'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_gla_cover_band_rates LONGTEXT;',
        'ALTER TABLE scheme_categories ADD COLUMN additional_gla_cover_band_rates LONGTEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

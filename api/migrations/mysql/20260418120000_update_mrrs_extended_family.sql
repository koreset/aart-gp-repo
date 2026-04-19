-- Migration: persist extended-family configuration + computed band rates on
-- member_rating_result_summaries so the Premiums Summary UI can render the
-- per-category extended-family section from a single payload.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='extended_family_benefit'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN extended_family_benefit TINYINT(1) DEFAULT 0;',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN extended_family_benefit TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='extended_family_age_band_source'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN extended_family_age_band_source VARCHAR(32);',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN extended_family_age_band_source VARCHAR(32);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='extended_family_age_band_type'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN extended_family_age_band_type VARCHAR(64);',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN extended_family_age_band_type VARCHAR(64);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='extended_family_pricing_method'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN extended_family_pricing_method VARCHAR(32);',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN extended_family_pricing_method VARCHAR(32);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='extended_family_band_rates'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN extended_family_band_rates TEXT;',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN extended_family_band_rates TEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_extended_family_monthly_premium'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_extended_family_monthly_premium DOUBLE NOT NULL DEFAULT 0;',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN total_extended_family_monthly_premium DOUBLE NOT NULL DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

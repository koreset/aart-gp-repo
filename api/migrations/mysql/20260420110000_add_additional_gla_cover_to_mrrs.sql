-- Migration: mirror Additional GLA Cover config + computed rates onto
-- member_rating_result_summaries.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_gla_cover_benefit'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN additional_gla_cover_benefit TINYINT(1) DEFAULT 0;',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_gla_cover_benefit TINYINT(1) DEFAULT 0;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_gla_cover_age_band_source'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN additional_gla_cover_age_band_source VARCHAR(16);',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_gla_cover_age_band_source VARCHAR(16);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_gla_cover_age_band_type'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN additional_gla_cover_age_band_type VARCHAR(64);',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_gla_cover_age_band_type VARCHAR(64);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_gla_cover_band_rates'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN additional_gla_cover_band_rates LONGTEXT;',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_gla_cover_band_rates LONGTEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_gla_cover_male_prop_used'),
        'ALTER TABLE member_rating_result_summaries MODIFY COLUMN additional_gla_cover_male_prop_used DOUBLE NULL;',
        'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_gla_cover_male_prop_used DOUBLE NULL;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

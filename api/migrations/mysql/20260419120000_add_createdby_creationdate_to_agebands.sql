-- Migration: add audit columns (created_by, creation_date) to
-- group_pricing_age_bands so uploaded age-band rows carry the same metadata
-- as every other rate table.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='group_pricing_age_bands' AND COLUMN_NAME='creation_date'),
        'SELECT 1;',
        'ALTER TABLE group_pricing_age_bands ADD COLUMN creation_date DATETIME NULL;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='group_pricing_age_bands' AND COLUMN_NAME='created_by'),
        'ALTER TABLE group_pricing_age_bands MODIFY COLUMN created_by VARCHAR(128);',
        'ALTER TABLE group_pricing_age_bands ADD COLUMN created_by VARCHAR(128);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

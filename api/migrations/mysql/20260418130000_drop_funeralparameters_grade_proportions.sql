-- Migration: drop legacy grade proportion columns from funeral_parameters.
-- Those fields belong to the Educator Rates table, not Funeral Parameters,
-- and were confusing users viewing the funeral_parameters CSV template.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='funeral_parameters' AND COLUMN_NAME='grade0_proportion'),
        'ALTER TABLE funeral_parameters DROP COLUMN grade0_proportion;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='funeral_parameters' AND COLUMN_NAME='grade1_7_proportion'),
        'ALTER TABLE funeral_parameters DROP COLUMN grade1_7_proportion;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='funeral_parameters' AND COLUMN_NAME='grade_8_12_proportion'),
        'ALTER TABLE funeral_parameters DROP COLUMN grade_8_12_proportion;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='funeral_parameters' AND COLUMN_NAME='tertiary_proportion'),
        'ALTER TABLE funeral_parameters DROP COLUMN tertiary_proportion;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

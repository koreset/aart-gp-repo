-- Migration for struct: SchemeSizeLevel

-- Table: scheme_size_levels

-- Ensure table exists
CREATE TABLE IF NOT EXISTS scheme_size_levels (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: RiskRateCode
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_size_levels' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_size_levels MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE scheme_size_levels ADD COLUMN risk_rate_code VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MinCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_size_levels' AND COLUMN_NAME='min_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_size_levels MODIFY COLUMN min_count INT;',
    'ALTER TABLE scheme_size_levels ADD COLUMN min_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: MaxCount
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_size_levels' AND COLUMN_NAME='max_count' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_size_levels MODIFY COLUMN max_count INT;',
    'ALTER TABLE scheme_size_levels ADD COLUMN max_count INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SizeLevel
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_size_levels' AND COLUMN_NAME='size_level' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_size_levels MODIFY COLUMN size_level INT;',
    'ALTER TABLE scheme_size_levels ADD COLUMN size_level INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_size_levels' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_size_levels MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE scheme_size_levels ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_size_levels' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE scheme_size_levels MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE scheme_size_levels ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


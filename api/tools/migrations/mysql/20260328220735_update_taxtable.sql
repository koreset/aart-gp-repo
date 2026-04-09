-- Migration for struct: TaxTable

-- Table: tax_tables

-- Ensure table exists
CREATE TABLE IF NOT EXISTS tax_tables (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: RiskRateCode
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_tables' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_tables MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE tax_tables ADD COLUMN risk_rate_code VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Level
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_tables' AND COLUMN_NAME='level' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_tables MODIFY COLUMN level INT;',
    'ALTER TABLE tax_tables ADD COLUMN level INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Min
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_tables' AND COLUMN_NAME='min' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_tables MODIFY COLUMN min DOUBLE;',
    'ALTER TABLE tax_tables ADD COLUMN min DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Max
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_tables' AND COLUMN_NAME='max' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_tables MODIFY COLUMN max DOUBLE;',
    'ALTER TABLE tax_tables ADD COLUMN max DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: TaxRate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_tables' AND COLUMN_NAME='tax_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_tables MODIFY COLUMN tax_rate DOUBLE;',
    'ALTER TABLE tax_tables ADD COLUMN tax_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_tables' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_tables MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE tax_tables ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_tables' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_tables MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE tax_tables ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


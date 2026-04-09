-- Migration for struct: PremiumLoading

-- Table: premium_loadings

-- Ensure table exists
CREATE TABLE IF NOT EXISTS premium_loadings (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- Add or modify column for field: RiskRateCode
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE premium_loadings ADD COLUMN risk_rate_code VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: Channel
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='channel' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN channel VARCHAR(255);',
    'ALTER TABLE premium_loadings ADD COLUMN channel VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: SchemeSizeLevel
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='scheme_size_level' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN scheme_size_level INT;',
    'ALTER TABLE premium_loadings ADD COLUMN scheme_size_level INT;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CommissionLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='commission_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN commission_loading DOUBLE;',
    'ALTER TABLE premium_loadings ADD COLUMN commission_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ExpenseLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='expense_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN expense_loading DOUBLE;',
    'ALTER TABLE premium_loadings ADD COLUMN expense_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: AdminLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='admin_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN admin_loading DOUBLE;',
    'ALTER TABLE premium_loadings ADD COLUMN admin_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: OtherLoading
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='other_loading' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN other_loading DOUBLE;',
    'ALTER TABLE premium_loadings ADD COLUMN other_loading DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: ProfitMargin
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='profit_margin' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN profit_margin DOUBLE;',
    'ALTER TABLE premium_loadings ADD COLUMN profit_margin DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreationDate
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE premium_loadings ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Add or modify column for field: CreatedBy
-- MySQL: Add or modify column
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='premium_loadings' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE premium_loadings MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE premium_loadings ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


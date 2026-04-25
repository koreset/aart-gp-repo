-- Migration for struct: TaxRetirementTable
-- Table: tax_retirement_tables

CREATE TABLE IF NOT EXISTS tax_retirement_tables (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_retirement_tables' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_retirement_tables MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE tax_retirement_tables ADD COLUMN risk_rate_code VARCHAR(255);'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_retirement_tables' AND COLUMN_NAME='lower_bound' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_retirement_tables MODIFY COLUMN lower_bound DOUBLE;',
    'ALTER TABLE tax_retirement_tables ADD COLUMN lower_bound DOUBLE;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_retirement_tables' AND COLUMN_NAME='upper_bound' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_retirement_tables MODIFY COLUMN upper_bound DOUBLE;',
    'ALTER TABLE tax_retirement_tables ADD COLUMN upper_bound DOUBLE;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_retirement_tables' AND COLUMN_NAME='tax_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_retirement_tables MODIFY COLUMN tax_rate DOUBLE;',
    'ALTER TABLE tax_retirement_tables ADD COLUMN tax_rate DOUBLE;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_retirement_tables' AND COLUMN_NAME='cumulative_tax_relief' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_retirement_tables MODIFY COLUMN cumulative_tax_relief DOUBLE;',
    'ALTER TABLE tax_retirement_tables ADD COLUMN cumulative_tax_relief DOUBLE;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_retirement_tables' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_retirement_tables MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE tax_retirement_tables ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='tax_retirement_tables' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE tax_retirement_tables MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE tax_retirement_tables ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

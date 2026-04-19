-- Migration for struct: ReinsuranceFuneralAidsRate
-- Table: reinsurance_funeral_aids_rates

CREATE TABLE IF NOT EXISTS reinsurance_funeral_aids_rates (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_funeral_aids_rates' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_funeral_aids_rates MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN risk_rate_code VARCHAR(255);'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_funeral_aids_rates' AND COLUMN_NAME='age_next_birthday' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_funeral_aids_rates MODIFY COLUMN age_next_birthday INT;',
    'ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN age_next_birthday INT;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_funeral_aids_rates' AND COLUMN_NAME='gender' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_funeral_aids_rates MODIFY COLUMN gender VARCHAR(255);',
    'ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN gender VARCHAR(255);'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_funeral_aids_rates' AND COLUMN_NAME='fun_aids_qx' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_funeral_aids_rates MODIFY COLUMN fun_aids_qx DOUBLE;',
    'ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN fun_aids_qx DOUBLE;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_funeral_aids_rates' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_funeral_aids_rates MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN creation_date DATETIME;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_funeral_aids_rates' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_funeral_aids_rates MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN created_by VARCHAR(255);'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

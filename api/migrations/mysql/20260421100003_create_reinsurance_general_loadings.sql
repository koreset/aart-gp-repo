-- Migration for struct: ReinsuranceGeneralLoading
-- Table: reinsurance_general_loadings

CREATE TABLE IF NOT EXISTS reinsurance_general_loadings (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN risk_rate_code VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='age' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN age INT;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN age INT;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='gender' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN gender VARCHAR(255);',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN gender VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='gla_contigency_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN gla_contigency_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN gla_contigency_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ptd_contigency_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN ptd_contigency_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ptd_contigency_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ci_contigency_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN ci_contigency_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ci_contigency_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ttd_contigency_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN ttd_contigency_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ttd_contigency_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='phi_contigency_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN phi_contigency_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN phi_contigency_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='fun_contigency_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN fun_contigency_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN fun_contigency_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='continuation_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN continuation_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN continuation_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='terminal_illness_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN terminal_illness_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN terminal_illness_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ptd_accelerated_benefit_discount' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN ptd_accelerated_benefit_discount DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ptd_accelerated_benefit_discount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ci_accelerated_benefit_discount' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN ci_accelerated_benefit_discount DOUBLE;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ci_accelerated_benefit_discount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN creation_date DATETIME;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_general_loadings MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN created_by VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

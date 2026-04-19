-- Migration for struct: ReinsuranceIndustryLoading
-- Table: reinsurance_industry_loadings

CREATE TABLE IF NOT EXISTS reinsurance_industry_loadings (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN risk_rate_code VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='occupation_class' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN occupation_class INT;',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN occupation_class INT;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='gender' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN gender VARCHAR(255);',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN gender VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='gla_industry_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN gla_industry_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN gla_industry_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='ptd_industry_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN ptd_industry_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN ptd_industry_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='ci_industry_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN ci_industry_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN ci_industry_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='ttd_industry_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN ttd_industry_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN ttd_industry_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='phi_industry_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN phi_industry_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN phi_industry_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN creation_date DATETIME;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_industry_loadings' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_industry_loadings MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE reinsurance_industry_loadings ADD COLUMN created_by VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

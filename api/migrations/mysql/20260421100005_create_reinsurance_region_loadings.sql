-- Migration for struct: ReinsuranceRegionLoading
-- Table: reinsurance_region_loadings

CREATE TABLE IF NOT EXISTS reinsurance_region_loadings (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='risk_rate_code' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN risk_rate_code VARCHAR(255);',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN risk_rate_code VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='region' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN region VARCHAR(255);',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN region VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='gender' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN gender VARCHAR(255);',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN gender VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='gla_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN gla_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN gla_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='gla_aids_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN gla_aids_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN gla_aids_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='ptd_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN ptd_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN ptd_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='ci_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN ci_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN ci_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='ttd_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN ttd_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN ttd_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='phi_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN phi_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN phi_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='fun_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN fun_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN fun_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='fun_aids_region_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN fun_aids_region_loading_rate DOUBLE;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN fun_aids_region_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='creation_date' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN creation_date DATETIME;',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN creation_date DATETIME;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_region_loadings' AND COLUMN_NAME='created_by' AND TABLE_SCHEMA = DATABASE()),
    'ALTER TABLE reinsurance_region_loadings MODIFY COLUMN created_by VARCHAR(255);',
    'ALTER TABLE reinsurance_region_loadings ADD COLUMN created_by VARCHAR(255);'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

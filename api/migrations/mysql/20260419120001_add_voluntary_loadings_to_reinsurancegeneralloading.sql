-- Migration for struct: ReinsuranceGeneralLoading (voluntary loadings)
-- Table: reinsurance_general_loadings

-- gla_voluntary_loading_rate
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='gla_voluntary_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN gla_voluntary_loading_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ptd_voluntary_loading_rate
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ptd_voluntary_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ptd_voluntary_loading_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ci_voluntary_loading_rate
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ci_voluntary_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ci_voluntary_loading_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ttd_voluntary_loading_rate
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='ttd_voluntary_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN ttd_voluntary_loading_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- phi_voluntary_loading_rate
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='phi_voluntary_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN phi_voluntary_loading_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- fun_voluntary_loading_rate
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='reinsurance_general_loadings' AND COLUMN_NAME='fun_voluntary_loading_rate' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE reinsurance_general_loadings ADD COLUMN fun_voluntary_loading_rate DOUBLE;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

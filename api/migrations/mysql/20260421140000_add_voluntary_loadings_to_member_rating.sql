-- Migration: add per-benefit voluntary loading columns to member_rating_results.

CREATE TABLE IF NOT EXISTS member_rating_results (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ci_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ci_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ttd_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ttd_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN phi_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN phi_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='fun_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN fun_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN fun_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ptd_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ptd_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ptd_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ci_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ci_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ci_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ttd_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ttd_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ttd_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_phi_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_phi_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_phi_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_fun_voluntary_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_fun_voluntary_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_fun_voluntary_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


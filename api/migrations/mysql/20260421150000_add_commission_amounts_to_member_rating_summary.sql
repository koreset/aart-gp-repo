-- Migration: add per-benefit commission amount columns and scheme-wide commission totals to member_rating_result_summaries.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_gla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_gla_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_gla_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_add_acc_gla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_add_acc_gla_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_add_acc_gla_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ptd_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ptd_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ptd_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ci_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ci_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ci_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_sgla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_sgla_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_sgla_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ttd_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ttd_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ttd_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_phi_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_phi_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_phi_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_fun_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_fun_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_fun_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_educator_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_educator_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_educator_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='scheme_total_commission' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN scheme_total_commission DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN scheme_total_commission DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='scheme_total_commission_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN scheme_total_commission_rate DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN scheme_total_commission_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


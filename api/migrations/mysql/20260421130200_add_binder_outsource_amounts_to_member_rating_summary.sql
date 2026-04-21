-- Migration: add per-benefit binder and outsource fee columns to member_rating_result_summaries.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_gla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_gla_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_gla_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_gla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_gla_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_gla_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_add_acc_gla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_add_acc_gla_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_add_acc_gla_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_add_acc_gla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_add_acc_gla_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_add_acc_gla_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_add_acc_gla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_add_acc_gla_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_add_acc_gla_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_add_acc_gla_annual_outsourced_amt' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_add_acc_gla_annual_outsourced_amt DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_add_acc_gla_annual_outsourced_amt DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ptd_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ptd_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ptd_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ptd_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ptd_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ptd_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ci_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ci_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ci_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ci_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ci_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ci_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ci_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ci_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ci_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ci_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ci_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ci_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_sgla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_sgla_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_sgla_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_sgla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_sgla_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_sgla_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_sgla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_sgla_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_sgla_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_sgla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_sgla_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_sgla_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ttd_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ttd_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ttd_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ttd_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ttd_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ttd_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ttd_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ttd_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_ttd_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_ttd_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_phi_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_phi_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_phi_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_phi_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_phi_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_phi_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_phi_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_phi_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_phi_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_phi_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_phi_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_phi_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_fun_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_fun_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_fun_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_fun_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_fun_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_fun_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_fun_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_fun_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_fun_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_fun_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_fun_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_fun_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_annual_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_annual_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


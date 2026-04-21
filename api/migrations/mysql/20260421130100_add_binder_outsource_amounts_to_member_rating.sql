-- Migration: add per-benefit binder and outsource fee columns to member_rating_results.

CREATE TABLE IF NOT EXISTS member_rating_results (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN additional_accidental_gla_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN additional_accidental_gla_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_add_acc_gla_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_add_acc_gla_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_add_acc_gla_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_add_acc_gla_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_add_acc_gla_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_add_acc_gla_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ci_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ci_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ci_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ci_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ci_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ci_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ci_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ci_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ci_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ci_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_spouse_gla_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_spouse_gla_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_spouse_gla_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_spouse_gla_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_spouse_gla_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_spouse_gla_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ttd_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ttd_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ttd_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ttd_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ttd_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ttd_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ttd_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ttd_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN phi_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN phi_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN phi_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN phi_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_phi_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_phi_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_phi_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_phi_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_phi_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_phi_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_funeral_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN main_member_funeral_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN main_member_funeral_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_funeral_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN main_member_funeral_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN main_member_funeral_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_funeral_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_funeral_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_funeral_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_funeral_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_funeral_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_funeral_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='children_funeral_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN children_funeral_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN children_funeral_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='children_funeral_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN children_funeral_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN children_funeral_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependants_funeral_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN dependants_funeral_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN dependants_funeral_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependants_funeral_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN dependants_funeral_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN dependants_funeral_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_funeral_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN total_funeral_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN total_funeral_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_funeral_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN total_funeral_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN total_funeral_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_total_funeral_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_total_funeral_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_total_funeral_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_total_funeral_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_total_funeral_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_total_funeral_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN total_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN total_binder_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='total_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN total_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN total_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


-- Migration: add Final*BinderAmount + Final*OutsourcedAmount columns to
-- member_rating_result_summaries. These mirror the existing Exp*Binder /
-- Exp*Outsourced fields but are derived from the post-discount Final office
-- premium so the breakdown reconciles to FinalOfficePremium after a discount
-- is applied. Pre-discount they equal the Exp* values; on non-binder
-- distribution channels they remain 0.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_add_acc_gla_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_add_acc_gla_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_annual_binder_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_annual_outsourced_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

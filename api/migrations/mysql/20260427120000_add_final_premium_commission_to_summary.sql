-- Migration: add Discount + Final* office-premium and commission columns to
-- member_rating_result_summaries. Final*OfficePremium is persisted as the
-- post-discount pre-comm office premium (Exp*RiskPremium /
-- (1 - (SchemeTotalLoading + Discount))) plus its per-benefit commission
-- slice, so the Final* values include commission and reconcile to
-- final_total_annual_premium{,_excl_funeral}. Exp* values stay pre-commission
-- and frozen at calc time, so pre-discount Final - Exp == commission slice.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='discount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN discount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_add_acc_gla_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_office_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_annual_office_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_add_acc_gla_annual_comm_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_add_acc_gla_annual_comm_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_annual_commission_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_annual_comm_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_annual_comm_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_annual_comm_amount' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_annual_comm_amount DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_total_annual_premium_excl_funeral' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_total_annual_premium_excl_funeral DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_total_annual_premium' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_total_annual_premium DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_scheme_total_commission' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_scheme_total_commission DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_scheme_total_commission_rate' AND TABLE_SCHEMA = DATABASE()), 'SELECT 1', 'ALTER TABLE member_rating_result_summaries ADD COLUMN final_scheme_total_commission_rate DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

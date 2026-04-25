-- Migration: add the non-Exp Total*CommissionAmount mirrors plus the per-category
-- ExpTotalCommission roll-up to member_rating_result_summaries. Parallel to the
-- existing Exp*CommissionAmount columns added in 20260421150000 so the scheme-wide
-- commission can be distributed on both the Exp and Total premium bases.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_annual_premium_excl_funeral' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_annual_premium_excl_funeral DOUBLE DEFAULT 0;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_premium_excl_funeral DOUBLE DEFAULT 0;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_annual_premium_excl_funeral' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_annual_premium_excl_funeral DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_annual_premium_excl_funeral DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_commission' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_commission DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_commission DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_add_acc_gla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_add_acc_gla_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_add_acc_gla_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ci_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ci_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ci_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_sgla_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_sgla_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_sgla_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ttd_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ttd_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_phi_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_phi_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_phi_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_fun_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_fun_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_fun_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_educator_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_educator_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_educator_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_educator_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_educator_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_educator_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

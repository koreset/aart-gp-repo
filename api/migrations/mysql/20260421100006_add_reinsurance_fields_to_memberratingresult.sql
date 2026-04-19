-- Migration: add reinsurance rate/loading/premium fields to member_rating_results
-- Adds the 62 columns introduced alongside the reinsurance rate tables so that
-- GORM INSERTs from MemberRatingResult succeed against existing DBs.

CREATE TABLE IF NOT EXISTS member_rating_results (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_industry_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_industry_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_industry_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ptd_industry_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ptd_industry_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ptd_industry_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ci_industry_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ci_industry_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ci_industry_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ttd_industry_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ttd_industry_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ttd_industry_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_phi_industry_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_phi_industry_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_phi_industry_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_aids_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_aids_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_aids_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ptd_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ptd_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ptd_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ci_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ci_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ci_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ttd_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ttd_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ttd_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_phi_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_phi_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_phi_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_fun_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_fun_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_fun_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_fun_aids_region_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_fun_aids_region_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_fun_aids_region_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_contingency_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_contingency_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_contingency_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ptd_contingency_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ptd_contingency_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ptd_contingency_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ci_contingency_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ci_contingency_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ci_contingency_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ttd_contingency_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ttd_contingency_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ttd_contingency_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_phi_contingency_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_phi_contingency_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_phi_contingency_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_fun_contingency_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_fun_contingency_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_fun_contingency_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_continuation_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_continuation_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_continuation_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_terminal_illness_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_terminal_illness_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_terminal_illness_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_qx' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_qx DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_qx DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_gla_aids_qx' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_gla_aids_qx DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_gla_aids_qx DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_reins_gla_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN base_reins_gla_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN base_reins_gla_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_reins_gla_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN loaded_reins_gla_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN loaded_reins_gla_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ptd_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ptd_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ptd_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_reins_ptd_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN base_reins_ptd_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN base_reins_ptd_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_reins_ptd_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN loaded_reins_ptd_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN loaded_reins_ptd_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_ci_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_ci_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_ci_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_reins_ci_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN base_reins_ci_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN base_reins_ci_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_reins_ci_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN loaded_reins_ci_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN loaded_reins_ci_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_reins_ttd_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN base_reins_ttd_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN base_reins_ttd_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_reins_ttd_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN loaded_reins_ttd_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN loaded_reins_ttd_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_phi_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_phi_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_phi_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_reins_phi_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN base_reins_phi_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN base_reins_phi_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_reins_phi_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN loaded_reins_phi_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN loaded_reins_phi_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_spouse_gla_qx' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_spouse_gla_qx DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_spouse_gla_qx DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_spouse_gla_aids_qx' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_spouse_gla_aids_qx DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_spouse_gla_aids_qx DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='reins_spouse_gla_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN reins_spouse_gla_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN reins_spouse_gla_loading DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='base_reins_spouse_gla_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN base_reins_spouse_gla_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN base_reins_spouse_gla_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='loaded_reins_spouse_gla_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN loaded_reins_spouse_gla_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN loaded_reins_spouse_gla_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_reinsurance_base_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN main_member_reinsurance_base_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN main_member_reinsurance_base_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_reinsurance_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN main_member_reinsurance_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN main_member_reinsurance_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_reinsurance_base_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_reinsurance_base_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_reinsurance_base_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_reinsurance_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_reinsurance_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_reinsurance_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='child_reinsurance_base_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN child_reinsurance_base_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN child_reinsurance_base_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='child_reinsurance_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN child_reinsurance_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN child_reinsurance_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='parent_reinsurance_base_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN parent_reinsurance_base_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN parent_reinsurance_base_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='parent_reinsurance_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN parent_reinsurance_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN parent_reinsurance_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependant_reinsurance_base_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN dependant_reinsurance_base_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN dependant_reinsurance_base_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependant_reinsurance_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN dependant_reinsurance_rate DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN dependant_reinsurance_rate DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ci_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ci_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_gla_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_gla_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ttd_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ttd_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN phi_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN phi_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN main_member_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN main_member_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN spouse_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN spouse_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='child_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN child_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN child_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='parent_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN parent_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN parent_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependant_reinsurance_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN dependant_reinsurance_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN dependant_reinsurance_premium DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;


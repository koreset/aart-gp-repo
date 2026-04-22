-- Migration: split educator benefit tracking into GLA and PTD components.
-- The combined Educator* columns stay as the sum; the new columns let the
-- business attribute the educator premium between GLA-educator and
-- PTD-educator and expose buildup fields (premium, %salary, rate per 1000).

CREATE TABLE IF NOT EXISTS member_rating_results (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_educator_sum_assured' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_educator_sum_assured DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_educator_sum_assured DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_gla_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_gla_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_gla_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_gla_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_gla_educator_risk_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN proportion_gla_educator_risk_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_gla_educator_risk_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_gla_educator_office_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN proportion_gla_educator_office_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_gla_educator_office_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_proportion_gla_educator_risk_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_proportion_gla_educator_risk_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_proportion_gla_educator_risk_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_proportion_gla_educator_office_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_proportion_gla_educator_office_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_proportion_gla_educator_office_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_educator_risk_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN gla_educator_risk_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN gla_educator_risk_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_educator_office_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN gla_educator_office_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN gla_educator_office_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_educator_risk_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_gla_educator_risk_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_gla_educator_risk_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_educator_office_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_gla_educator_office_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_gla_educator_office_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_educator_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_ptd_educator_risk_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_ptd_educator_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_educator_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_ptd_educator_office_premium DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_ptd_educator_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_ptd_educator_risk_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN proportion_ptd_educator_risk_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_ptd_educator_risk_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_ptd_educator_office_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN proportion_ptd_educator_office_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_ptd_educator_office_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_proportion_ptd_educator_risk_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_proportion_ptd_educator_risk_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_proportion_ptd_educator_risk_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_proportion_ptd_educator_office_premium_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_proportion_ptd_educator_office_premium_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_proportion_ptd_educator_office_premium_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_educator_risk_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN ptd_educator_risk_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN ptd_educator_risk_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_educator_office_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN ptd_educator_office_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN ptd_educator_office_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_educator_risk_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_ptd_educator_risk_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_ptd_educator_risk_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_educator_office_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_ptd_educator_office_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_ptd_educator_office_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Split binder / outsource per-member columns
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_gla_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_gla_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ptd_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ptd_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ptd_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ptd_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Split binder / outsource summary totals
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_gla_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_gla_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_gla_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_gla_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_gla_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_gla_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ptd_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ptd_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_educator_binder_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_ptd_educator_binder_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_ptd_educator_binder_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_educator_outsourced_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_ptd_educator_outsourced_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_ptd_educator_outsourced_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Split commission summary totals
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_educator_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_gla_educator_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_gla_educator_commission_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_educator_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_ptd_educator_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_ptd_educator_commission_amount DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

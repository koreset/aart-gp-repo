-- Migration: add Additional Accidental GLA aggregate columns to
-- member_rating_result_summaries. The original
-- 20260417150000_add_additional_accidental_gla.sql migration covered
-- scheme_categories and member_rating_results but missed the summary
-- table — causing `Unknown column 'min_additional_accidental_gla_sum_assured'`
-- on DB.Create(&mdrs) during quote recalculation.

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='min_additional_accidental_gla_sum_assured'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN min_additional_accidental_gla_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN min_additional_accidental_gla_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='max_additional_accidental_gla_sum_assured'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN max_additional_accidental_gla_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN max_additional_accidental_gla_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='max_additional_accidental_gla_capped_sum_assured'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN max_additional_accidental_gla_capped_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN max_additional_accidental_gla_capped_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_additional_accidental_gla_sum_assured'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_additional_accidental_gla_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN total_additional_accidental_gla_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_additional_accidental_gla_capped_sum_assured'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_additional_accidental_gla_capped_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN total_additional_accidental_gla_capped_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='average_additional_accidental_gla_capped_sum_assured'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN average_additional_accidental_gla_capped_sum_assured DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN average_additional_accidental_gla_capped_sum_assured DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_additional_accidental_gla_risk_rate'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_additional_accidental_gla_risk_rate DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN total_additional_accidental_gla_risk_rate DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_additional_accidental_gla_annual_risk_premium'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_additional_accidental_gla_annual_risk_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN total_additional_accidental_gla_annual_risk_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_accidental_gla_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_additional_accidental_gla_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN proportion_additional_accidental_gla_annual_risk_premium_salary DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_additional_accidental_gla_annual_risk_premium_salary DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_additional_accidental_gla_annual_office_premium'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_additional_accidental_gla_annual_office_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN total_additional_accidental_gla_annual_office_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_accidental_gla_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_additional_accidental_gla_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN proportion_additional_accidental_gla_office_premium_salary DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_additional_accidental_gla_office_premium_salary DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_additional_accidental_gla_risk_rate'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_additional_accidental_gla_risk_rate DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_additional_accidental_gla_risk_rate DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_additional_accidental_gla_annual_risk_premium'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_additional_accidental_gla_annual_risk_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_additional_accidental_gla_annual_risk_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_additional_accidental_gla_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Column name abbreviated to fit MySQL's 64-char identifier limit. The Go
-- field uses `gorm:"column:exp_prop_additional_accidental_gla_annual_risk_premium_salary"`
-- so only the DB column is shortened — JSON/CSV serialisation keeps the full name.
SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_prop_additional_accidental_gla_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_prop_additional_accidental_gla_annual_risk_premium_salary DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_prop_additional_accidental_gla_annual_risk_premium_salary DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_additional_accidental_gla_annual_office_premium'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_additional_accidental_gla_annual_office_premium DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_additional_accidental_gla_annual_office_premium DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_additional_accidental_gla_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_additional_accidental_gla_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_proportion_additional_accidental_gla_office_premium_salary DECIMAL(15,5);',
    'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_proportion_additional_accidental_gla_office_premium_salary DECIMAL(15,5);'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

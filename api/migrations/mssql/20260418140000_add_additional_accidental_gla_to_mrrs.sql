-- Migration: add Additional Accidental GLA aggregate columns to
-- member_rating_result_summaries. The original
-- 20260417150000_add_additional_accidental_gla.sql migration covered
-- scheme_categories and member_rating_results but missed the summary
-- table — causing `Unknown column 'min_additional_accidental_gla_sum_assured'`
-- on DB.Create(&mdrs) during quote recalculation.

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'min_additional_accidental_gla_sum_assured')
BEGIN ALTER TABLE member_rating_result_summaries ADD min_additional_accidental_gla_sum_assured DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'max_additional_accidental_gla_sum_assured')
BEGIN ALTER TABLE member_rating_result_summaries ADD max_additional_accidental_gla_sum_assured DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'max_additional_accidental_gla_capped_sum_assured')
BEGIN ALTER TABLE member_rating_result_summaries ADD max_additional_accidental_gla_capped_sum_assured DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_additional_accidental_gla_sum_assured')
BEGIN ALTER TABLE member_rating_result_summaries ADD total_additional_accidental_gla_sum_assured DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_additional_accidental_gla_capped_sum_assured')
BEGIN ALTER TABLE member_rating_result_summaries ADD total_additional_accidental_gla_capped_sum_assured DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'average_additional_accidental_gla_capped_sum_assured')
BEGIN ALTER TABLE member_rating_result_summaries ADD average_additional_accidental_gla_capped_sum_assured DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_additional_accidental_gla_risk_rate')
BEGIN ALTER TABLE member_rating_result_summaries ADD total_additional_accidental_gla_risk_rate DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_additional_accidental_gla_annual_risk_premium')
BEGIN ALTER TABLE member_rating_result_summaries ADD total_additional_accidental_gla_annual_risk_premium DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'additional_accidental_gla_risk_rate_per1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries ADD additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'proportion_additional_accidental_gla_annual_risk_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries ADD proportion_additional_accidental_gla_annual_risk_premium_salary DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_additional_accidental_gla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries ADD total_additional_accidental_gla_annual_office_premium DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'additional_accidental_gla_office_rate_per1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries ADD additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'proportion_additional_accidental_gla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries ADD proportion_additional_accidental_gla_office_premium_salary DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_additional_accidental_gla_risk_rate')
BEGIN ALTER TABLE member_rating_result_summaries ADD exp_total_additional_accidental_gla_risk_rate DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_additional_accidental_gla_annual_risk_premium')
BEGIN ALTER TABLE member_rating_result_summaries ADD exp_total_additional_accidental_gla_annual_risk_premium DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_additional_accidental_gla_risk_rate_per1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries ADD exp_additional_accidental_gla_risk_rate_per1000_sa DECIMAL(15,5); END

-- Abbreviated to stay within MySQL's 64-char identifier cap so the Go
-- model can use one gorm:"column:" override across dialects.
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_prop_additional_accidental_gla_annual_risk_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries ADD exp_prop_additional_accidental_gla_annual_risk_premium_salary DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_additional_accidental_gla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries ADD exp_total_additional_accidental_gla_annual_office_premium DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_additional_accidental_gla_office_rate_per1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries ADD exp_additional_accidental_gla_office_rate_per1000_sa DECIMAL(15,5); END

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_proportion_additional_accidental_gla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries ADD exp_proportion_additional_accidental_gla_office_premium_salary DECIMAL(15,5); END

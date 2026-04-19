-- Migration: add Additional Accidental GLA aggregate columns to
-- member_rating_result_summaries. The original
-- 20260417150000_add_additional_accidental_gla.sql migration covered
-- scheme_categories and member_rating_results but missed the summary
-- table — causing `Unknown column 'min_additional_accidental_gla_sum_assured'`
-- on DB.Create(&mdrs) during quote recalculation.

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS min_additional_accidental_gla_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS max_additional_accidental_gla_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS max_additional_accidental_gla_capped_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_additional_accidental_gla_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_additional_accidental_gla_capped_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS average_additional_accidental_gla_capped_sum_assured NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_additional_accidental_gla_risk_rate NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_additional_accidental_gla_annual_risk_premium NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS additional_accidental_gla_risk_rate_per1000_sa NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS proportion_additional_accidental_gla_annual_risk_premium_salary NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_additional_accidental_gla_annual_office_premium NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS additional_accidental_gla_office_rate_per1000_sa NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS proportion_additional_accidental_gla_office_premium_salary NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_additional_accidental_gla_risk_rate NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_additional_accidental_gla_annual_risk_premium NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_additional_accidental_gla_risk_rate_per1000_sa NUMERIC(15,5);
-- Abbreviated to stay within MySQL's 64-char identifier cap so the Go
-- model can use one gorm:"column:" override across dialects.
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_prop_additional_accidental_gla_annual_risk_premium_salary NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_additional_accidental_gla_annual_office_premium NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_additional_accidental_gla_office_rate_per1000_sa NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_proportion_additional_accidental_gla_office_premium_salary NUMERIC(15,5);

-- Migration: rename four Additional Accidental GLA rate-per-1000 columns on
-- member_rating_result_summaries to match GORM's NamingStrategy output
-- (`per1000_sa`, no underscore between `per` and `1000`).

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_accidental_gla_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_accidental_gla_risk_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.additional_accidental_gla_risk_rate_per_1000_sa', 'additional_accidental_gla_risk_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_accidental_gla_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_accidental_gla_office_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.additional_accidental_gla_office_rate_per_1000_sa', 'additional_accidental_gla_office_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_additional_accidental_gla_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_additional_accidental_gla_risk_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.exp_additional_accidental_gla_risk_rate_per_1000_sa', 'exp_additional_accidental_gla_risk_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_additional_accidental_gla_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_additional_accidental_gla_office_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.exp_additional_accidental_gla_office_rate_per_1000_sa', 'exp_additional_accidental_gla_office_rate_per1000_sa', 'COLUMN';
END

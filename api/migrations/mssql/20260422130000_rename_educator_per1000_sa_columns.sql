-- Migration: rename educator GLA/PTD rate-per-1000 columns on
-- member_rating_result_summaries to match GORM's NamingStrategy output
-- (`per1000_sa`, no underscore between `per` and `1000`). The
-- 20260421170000 migration used `per_1000_sa` which GORM cannot find on
-- INSERT/UPDATE — causing "Unknown column" 1054 errors during
-- CalculateGroupPricingQuote. Follows the precedent of
-- 20260418150000_rename_additional_accidental_gla_per1000.sql.

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_educator_risk_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.gla_educator_risk_rate_per_1000_sa', 'gla_educator_risk_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_educator_office_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.gla_educator_office_rate_per_1000_sa', 'gla_educator_office_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_educator_risk_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.exp_gla_educator_risk_rate_per_1000_sa', 'exp_gla_educator_risk_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_educator_office_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.exp_gla_educator_office_rate_per_1000_sa', 'exp_gla_educator_office_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_educator_risk_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.ptd_educator_risk_rate_per_1000_sa', 'ptd_educator_risk_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_educator_office_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.ptd_educator_office_rate_per_1000_sa', 'ptd_educator_office_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_educator_risk_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_educator_risk_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.exp_ptd_educator_risk_rate_per_1000_sa', 'exp_ptd_educator_risk_rate_per1000_sa', 'COLUMN';
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_educator_office_rate_per_1000_sa')
    AND NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_educator_office_rate_per1000_sa')
BEGIN
    EXEC sp_rename 'member_rating_result_summaries.exp_ptd_educator_office_rate_per_1000_sa', 'exp_ptd_educator_office_rate_per1000_sa', 'COLUMN';
END

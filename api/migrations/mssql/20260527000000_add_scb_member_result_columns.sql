-- Add Salary Continuation Benefit (SCB) per-member and summary columns.
-- Mirrors the field structure of the conversion-on-withdrawal slices on
-- MemberRatingResult / MemberRatingResultSummary; SCB is tracked as a
-- reportable slice (NOT added to any group total).
--
-- Source-side rate tables and scheme_categories columns were added in
-- 20260519010000_add_salary_continuation_benefit_tables.sql.
-- Idempotent on re-runs.

-- ── member_rating_results ──────────────────────────────────────────────
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'scb_rate')
BEGIN
    ALTER TABLE member_rating_results ADD scb_rate FLOAT NOT NULL CONSTRAINT df_mrr_scb_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'base_scb_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_scb_rate FLOAT NOT NULL CONSTRAINT df_mrr_base_scb_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'loaded_scb_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_scb_rate FLOAT NOT NULL CONSTRAINT df_mrr_loaded_scb_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_loaded_scb_rate')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_loaded_scb_rate FLOAT NOT NULL CONSTRAINT df_mrr_exp_adj_loaded_scb_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'scb_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD scb_risk_premium FLOAT NOT NULL CONSTRAINT df_mrr_scb_risk_premium DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_scb_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD exp_adj_scb_risk_premium FLOAT NOT NULL CONSTRAINT df_mrr_exp_adj_scb_risk_premium DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'reins_scb_rate')
BEGIN
    ALTER TABLE member_rating_results ADD reins_scb_rate FLOAT NOT NULL CONSTRAINT df_mrr_reins_scb_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'base_reins_scb_rate')
BEGIN
    ALTER TABLE member_rating_results ADD base_reins_scb_rate FLOAT NOT NULL CONSTRAINT df_mrr_base_reins_scb_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'loaded_reins_scb_rate')
BEGIN
    ALTER TABLE member_rating_results ADD loaded_reins_scb_rate FLOAT NOT NULL CONSTRAINT df_mrr_loaded_reins_scb_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'reins_scb_risk_premium')
BEGIN
    ALTER TABLE member_rating_results ADD reins_scb_risk_premium FLOAT NOT NULL CONSTRAINT df_mrr_reins_scb_risk_premium DEFAULT 0;
END;

-- ── member_rating_result_summaries ─────────────────────────────────────
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_scb_annual_risk_premium')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD total_scb_annual_risk_premium FLOAT NOT NULL CONSTRAINT df_mrrs_total_scb_annual_risk_premium DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_scb_annual_risk_premium')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_scb_annual_risk_premium FLOAT NOT NULL CONSTRAINT df_mrrs_exp_adj_total_scb_annual_risk_premium DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_scb_risk_premium_salary')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD proportion_scb_risk_premium_salary FLOAT NOT NULL CONSTRAINT df_mrrs_proportion_scb_risk_premium_salary DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_proportion_scb_risk_premium_salary')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD exp_adj_proportion_scb_risk_premium_salary FLOAT NOT NULL CONSTRAINT df_mrrs_exp_adj_prop_scb_risk_prem_salary DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'scb_risk_rate_per_1000_income')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD scb_risk_rate_per_1000_income FLOAT NOT NULL CONSTRAINT df_mrrs_scb_risk_rate_per_1000_income DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_scb_risk_rate_per_1000_income')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD exp_scb_risk_rate_per_1000_income FLOAT NOT NULL CONSTRAINT df_mrrs_exp_scb_risk_rate_per_1000_income DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'final_scb_office_premium')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD final_scb_office_premium FLOAT NOT NULL CONSTRAINT df_mrrs_final_scb_office_premium DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_reins_scb_annual_risk_premium')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD total_reins_scb_annual_risk_premium FLOAT NOT NULL CONSTRAINT df_mrrs_total_reins_scb_annual_risk_premium DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'final_reins_scb_office_premium')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD final_reins_scb_office_premium FLOAT NOT NULL CONSTRAINT df_mrrs_final_reins_scb_office_premium DEFAULT 0;
END;

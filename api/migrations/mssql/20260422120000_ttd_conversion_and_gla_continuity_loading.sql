-- Migration: TTD conversion-on-withdrawal slice + GLA continuity-during-disability
-- dedicated loading. Adds per-benefit scheme_categories flags, general_loadings
-- rate columns, member_rating_results slice/loading columns, and
-- member_rating_result_summaries aggregates (4 totals + 4 proportions +
-- 4 rate-per-1000) for the new TTD slice. Follows the
-- 20260421180000_conversion_continuity_premium_slices pattern.
-- Idempotent column adds.

-- ── scheme_categories: new flags ──────────────────────────────────────
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ttd_conversion_on_withdrawal')
    ALTER TABLE scheme_categories ADD ttd_conversion_on_withdrawal BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_continuity_during_disability')
    ALTER TABLE scheme_categories ADD gla_continuity_during_disability BIT DEFAULT 0;

-- ── general_loadings: new per-slice loading rates ─────────────────────
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ttd_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD ttd_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_continuity_during_dis_loading_rate')
    ALTER TABLE general_loadings ADD gla_continuity_during_dis_loading_rate FLOAT;

-- ── member_rating_results: TTD slice + GLA continuity loading ─────────
-- Slice: TTD conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD ttd_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD ttd_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ttd_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ttd_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ttd_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD ttd_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ttd_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ttd_conv_on_wdr_office_premium FLOAT;
-- GLA continuity during disability: dedicated per-member loading (replaces reuse of continuation_loading)
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_continuity_during_dis_loading')
    ALTER TABLE member_rating_results ADD gla_continuity_during_dis_loading FLOAT;

-- ── member_rating_result_summaries: TTD slice aggregates (4+4+4) ──────
-- Slice: TTD conversion on withdrawal (denominator for rate-per-1000 = TotalTtdCappedIncome)
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ttd_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ttd_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ttd_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ttd_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ttd_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ttd_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ttd_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ttd_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ttd_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ttd_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ttd_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ttd_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ttd_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ttd_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ttd_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ttd_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ttd_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ttd_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ttd_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ttd_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ttd_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ttd_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ttd_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ttd_conv_on_wdr_office_rate_per_1000_sa FLOAT;

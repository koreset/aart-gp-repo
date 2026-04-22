-- Migration: TTD conversion-on-withdrawal slice + GLA continuity-during-disability
-- dedicated loading. Adds per-benefit scheme_categories flags, general_loadings
-- rate columns, member_rating_results slice/loading columns, and
-- member_rating_result_summaries aggregates (4 totals + 4 proportions +
-- 4 rate-per-1000) for the new TTD slice. Follows the
-- 20260421180000_conversion_continuity_premium_slices pattern.

-- ── scheme_categories: new flags ──────────────────────────────────────
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ttd_conversion_on_withdrawal BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_continuity_during_disability BOOLEAN DEFAULT FALSE;

-- ── general_loadings: new per-slice loading rates ─────────────────────
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ttd_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gla_continuity_during_dis_loading_rate NUMERIC(20,6);

-- ── member_rating_results: TTD slice + GLA continuity loading ─────────
-- Slice: TTD conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ttd_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ttd_conv_on_wdr_office_premium NUMERIC(20,6);
-- GLA continuity during disability: dedicated per-member loading (replaces reuse of continuation_loading)
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_continuity_during_dis_loading NUMERIC(20,6);

-- ── member_rating_result_summaries: TTD slice aggregates (4+4+4) ──────
-- Slice: TTD conversion on withdrawal (denominator for rate-per-1000 = TotalTtdCappedIncome)
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ttd_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ttd_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ttd_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ttd_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ttd_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ttd_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ttd_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ttd_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ttd_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ttd_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ttd_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ttd_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);

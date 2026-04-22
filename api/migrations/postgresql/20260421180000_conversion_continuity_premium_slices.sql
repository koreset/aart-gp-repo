-- Migration: conversion / continuity premium slice tracking. Adds
-- per-slice loading rates (general_loadings), per-member slice premium
-- fields (member_rating_results), and per-category summary buildup
-- (member_rating_result_summaries). Gated by new scheme_categories flags.
-- Mirrors the TaxSaver slice pattern (plan:
-- under-gla-benefit-lets-cheeky-hickey.md).

-- ── scheme_categories: new flags (8) ──────────────────────────────────
CREATE TABLE IF NOT EXISTS scheme_categories (id SERIAL PRIMARY KEY);
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS phi_conversion_on_withdrawal BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS sgla_conversion_on_withdrawal BOOLEAN DEFAULT FALSE;
ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS fun_conversion_on_withdrawal BOOLEAN DEFAULT FALSE;

-- ── general_loadings: new per-slice loading rates (12) ─────────────────
CREATE TABLE IF NOT EXISTS general_loadings (id SERIAL PRIMARY KEY);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gla_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gla_conv_on_ret_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ptd_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ci_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS phi_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS sgla_conv_on_wdr_loading_rate NUMERIC(20,6);
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS fun_conv_on_wdr_loading_rate NUMERIC(20,6);

-- ── member_rating_results: educator intermediate rates + slice columns ─
CREATE TABLE IF NOT EXISTS member_rating_results (id SERIAL PRIMARY KEY);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_gla_educator_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_gla_educator_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_ptd_educator_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_loaded_ptd_educator_rate NUMERIC(20,6);

-- Per-slice loading + 4 premium columns (risk + office + ExpAdj variants).
-- Slice 1: GLA conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_conv_on_wdr_office_premium NUMERIC(20,6);
-- Slice 9: GLA conversion on retirement
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_conv_on_ret_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_conv_on_ret_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_conv_on_ret_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_conv_on_ret_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_conv_on_ret_office_premium NUMERIC(20,6);
-- Slice 12: GLA continuity during disability (reuses ContinuationLoading)
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_cont_dur_dis_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_cont_dur_dis_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_cont_dur_dis_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_cont_dur_dis_office_premium NUMERIC(20,6);
-- Slice 2: GLA Educator conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_ed_conv_on_wdr_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_ed_conv_on_wdr_office_prem NUMERIC(20,6);
-- Slice 10: GLA Educator conversion on retirement
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_ed_conv_on_ret_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_ed_conv_on_ret_office_prem NUMERIC(20,6);
-- Slice 13: GLA Educator continuity during disability
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_ed_cont_dur_dis_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_ed_cont_dur_dis_office_prem NUMERIC(20,6);
-- Slice 4: PTD conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_conv_on_wdr_office_premium NUMERIC(20,6);
-- Slice 3: PTD Educator conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_ed_conv_on_wdr_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_ed_conv_on_wdr_office_prem NUMERIC(20,6);
-- Slice 11: PTD Educator conversion on retirement
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_ed_conv_on_ret_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_ed_conv_on_ret_office_prem NUMERIC(20,6);
-- Slice 6: CI conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ci_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ci_conv_on_wdr_office_premium NUMERIC(20,6);
-- Slice 5: PHI conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_phi_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_phi_conv_on_wdr_office_premium NUMERIC(20,6);
-- Slice 7: SGLA conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS sgla_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS sgla_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_sgla_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS sgla_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_sgla_conv_on_wdr_office_premium NUMERIC(20,6);
-- Slice 8: Funeral conversion on withdrawal
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS fun_conv_on_wdr_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS fun_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_fun_conv_on_wdr_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS fun_conv_on_wdr_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_fun_conv_on_wdr_office_premium NUMERIC(20,6);

-- ── member_rating_result_summaries: summary buildup (156) + funeral SA ─
CREATE TABLE IF NOT EXISTS member_rating_result_summaries (id SERIAL PRIMARY KEY);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_family_funeral_sum_assured NUMERIC(20,6);

-- Macro for each slice: 4 totals + 4 proportions + 4 rates-per-1000.
-- Slice 1: GLA conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 9: GLA conversion on retirement
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_conv_on_ret_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_conv_on_ret_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_conv_on_ret_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_conv_on_ret_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_conv_on_ret_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_conv_on_ret_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_conv_on_ret_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_conv_on_ret_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_conv_on_ret_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_conv_on_ret_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_conv_on_ret_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_conv_on_ret_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 12: GLA continuity during disability
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_cont_dur_dis_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_cont_dur_dis_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_cont_dur_dis_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_cont_dur_dis_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_cont_dur_dis_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_cont_dur_dis_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_cont_dur_dis_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_cont_dur_dis_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_cont_dur_dis_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_cont_dur_dis_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_cont_dur_dis_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_cont_dur_dis_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 2: GLA Educator conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_ed_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_ed_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_ed_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_ed_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_ed_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_ed_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_ed_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_ed_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_ed_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_ed_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_ed_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 10: GLA Educator conversion on retirement
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_ed_conv_on_ret_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_ed_conv_on_ret_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_ed_conv_on_ret_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_ed_conv_on_ret_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_ed_conv_on_ret_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_ed_conv_on_ret_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_ed_conv_on_ret_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_ed_conv_on_ret_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_ed_conv_on_ret_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_ed_conv_on_ret_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_ed_conv_on_ret_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 13: GLA Educator continuity during disability
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_ed_cont_dur_dis_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_ed_cont_dur_dis_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_ed_cont_dur_dis_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_ed_cont_dur_dis_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_ed_cont_dur_dis_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_gla_ed_cont_dur_dis_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_ed_cont_dur_dis_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_gla_ed_cont_dur_dis_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_ed_cont_dur_dis_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_ed_cont_dur_dis_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_ed_cont_dur_dis_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 4: PTD conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ptd_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ptd_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ptd_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ptd_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 3: PTD Educator conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_ed_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_ed_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_ed_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_ed_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ptd_ed_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ptd_ed_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ptd_ed_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ptd_ed_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_ed_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_ed_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 11: PTD Educator conversion on retirement
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_ed_conv_on_ret_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_ed_conv_on_ret_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_ed_conv_on_ret_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_ed_conv_on_ret_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ptd_ed_conv_on_ret_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ptd_ed_conv_on_ret_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ptd_ed_conv_on_ret_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ptd_ed_conv_on_ret_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_ed_conv_on_ret_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_ed_conv_on_ret_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_ed_conv_on_ret_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 5: PHI conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_phi_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_phi_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_phi_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_phi_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_phi_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_phi_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_phi_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_phi_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS phi_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS phi_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_phi_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_phi_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 6: CI conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ci_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ci_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ci_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ci_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ci_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_ci_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ci_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_ci_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ci_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ci_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ci_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ci_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 7: SGLA conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_sgla_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_sgla_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_sgla_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_sgla_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_sgla_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_sgla_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_sgla_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_sgla_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS sgla_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS sgla_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_sgla_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_sgla_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
-- Slice 8: Funeral conversion on withdrawal
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_fun_conv_on_wdr_annual_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_fun_conv_on_wdr_annual_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_fun_conv_on_wdr_ann_risk_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_fun_conv_on_wdr_ann_office_prem NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_fun_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS prop_fun_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_fun_conv_on_wdr_risk_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_prop_fun_conv_on_wdr_office_prem_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS fun_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS fun_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_fun_conv_on_wdr_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_fun_conv_on_wdr_office_rate_per_1000_sa NUMERIC(20,6);

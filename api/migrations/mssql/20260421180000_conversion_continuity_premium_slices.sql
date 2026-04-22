-- Migration: conversion / continuity premium slice tracking.
-- Generated from postgres source; idempotent column adds.

-- Migration: conversion / continuity premium slice tracking. Adds
-- per-slice loading rates (general_loadings), per-member slice premium
-- fields (member_rating_results), and per-category summary buildup
-- (member_rating_result_summaries). Gated by new scheme_categories flags.
-- Mirrors the TaxSaver slice pattern (plan:
-- under-gla-benefit-lets-cheeky-hickey.md).

-- ── scheme_categories: new flags (8) ──────────────────────────────────
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'scheme_categories')
BEGIN
    CREATE TABLE scheme_categories (id INT IDENTITY(1,1) PRIMARY KEY);
END;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_ed_conv_on_wdr')
    ALTER TABLE scheme_categories ADD gla_ed_conv_on_wdr BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_ed_conv_on_ret')
    ALTER TABLE scheme_categories ADD gla_ed_conv_on_ret BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'gla_ed_cont_dur_dis')
    ALTER TABLE scheme_categories ADD gla_ed_cont_dur_dis BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr')
    ALTER TABLE scheme_categories ADD ptd_ed_conv_on_wdr BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'ptd_ed_conv_on_ret')
    ALTER TABLE scheme_categories ADD ptd_ed_conv_on_ret BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'phi_conversion_on_withdrawal')
    ALTER TABLE scheme_categories ADD phi_conversion_on_withdrawal BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'sgla_conversion_on_withdrawal')
    ALTER TABLE scheme_categories ADD sgla_conversion_on_withdrawal BIT DEFAULT 0;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'fun_conversion_on_withdrawal')
    ALTER TABLE scheme_categories ADD fun_conversion_on_withdrawal BIT DEFAULT 0;

-- ── general_loadings: new per-slice loading rates (12) ─────────────────
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'general_loadings')
BEGIN
    CREATE TABLE general_loadings (id INT IDENTITY(1,1) PRIMARY KEY);
END;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD gla_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_conv_on_ret_loading_rate')
    ALTER TABLE general_loadings ADD gla_conv_on_ret_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_ed_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD gla_ed_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_ed_conv_on_ret_loading_rate')
    ALTER TABLE general_loadings ADD gla_ed_conv_on_ret_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'gla_ed_cont_dur_dis_loading_rate')
    ALTER TABLE general_loadings ADD gla_ed_cont_dur_dis_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ptd_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD ptd_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD ptd_ed_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ptd_ed_conv_on_ret_loading_rate')
    ALTER TABLE general_loadings ADD ptd_ed_conv_on_ret_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'ci_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD ci_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'phi_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD phi_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'sgla_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD sgla_conv_on_wdr_loading_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'general_loadings' AND COLUMN_NAME = 'fun_conv_on_wdr_loading_rate')
    ALTER TABLE general_loadings ADD fun_conv_on_wdr_loading_rate FLOAT;

-- ── member_rating_results: educator intermediate rates + slice columns ─
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_results')
BEGIN
    CREATE TABLE member_rating_results (id INT IDENTITY(1,1) PRIMARY KEY);
END;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_gla_educator_rate')
    ALTER TABLE member_rating_results ADD loaded_gla_educator_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_gla_educator_rate')
    ALTER TABLE member_rating_results ADD exp_adj_loaded_gla_educator_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'loaded_ptd_educator_rate')
    ALTER TABLE member_rating_results ADD loaded_ptd_educator_rate FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_loaded_ptd_educator_rate')
    ALTER TABLE member_rating_results ADD exp_adj_loaded_ptd_educator_rate FLOAT;

-- Per-slice loading + 4 premium columns (risk + office + ExpAdj variants).
-- Slice 1: GLA conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD gla_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD gla_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD gla_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_conv_on_wdr_office_premium FLOAT;
-- Slice 9: GLA conversion on retirement
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_conv_on_ret_loading')
    ALTER TABLE member_rating_results ADD gla_conv_on_ret_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_conv_on_ret_risk_premium')
    ALTER TABLE member_rating_results ADD gla_conv_on_ret_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_conv_on_ret_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_conv_on_ret_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_conv_on_ret_office_premium')
    ALTER TABLE member_rating_results ADD gla_conv_on_ret_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_conv_on_ret_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_conv_on_ret_office_premium FLOAT;
-- Slice 12: GLA continuity during disability (reuses ContinuationLoading)
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_cont_dur_dis_risk_premium')
    ALTER TABLE member_rating_results ADD gla_cont_dur_dis_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_cont_dur_dis_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_cont_dur_dis_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_cont_dur_dis_office_premium')
    ALTER TABLE member_rating_results ADD gla_cont_dur_dis_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_cont_dur_dis_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_gla_cont_dur_dis_office_premium FLOAT;
-- Slice 2: GLA Educator conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD gla_ed_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD gla_ed_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_ed_conv_on_wdr_risk_prem')
    ALTER TABLE member_rating_results ADD exp_adj_gla_ed_conv_on_wdr_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD gla_ed_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_ed_conv_on_wdr_office_prem')
    ALTER TABLE member_rating_results ADD exp_adj_gla_ed_conv_on_wdr_office_prem FLOAT;
-- Slice 10: GLA Educator conversion on retirement
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_conv_on_ret_loading')
    ALTER TABLE member_rating_results ADD gla_ed_conv_on_ret_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_conv_on_ret_risk_premium')
    ALTER TABLE member_rating_results ADD gla_ed_conv_on_ret_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_ed_conv_on_ret_risk_prem')
    ALTER TABLE member_rating_results ADD exp_adj_gla_ed_conv_on_ret_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_conv_on_ret_office_premium')
    ALTER TABLE member_rating_results ADD gla_ed_conv_on_ret_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_ed_conv_on_ret_office_prem')
    ALTER TABLE member_rating_results ADD exp_adj_gla_ed_conv_on_ret_office_prem FLOAT;
-- Slice 13: GLA Educator continuity during disability
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_cont_dur_dis_loading')
    ALTER TABLE member_rating_results ADD gla_ed_cont_dur_dis_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_cont_dur_dis_risk_premium')
    ALTER TABLE member_rating_results ADD gla_ed_cont_dur_dis_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_ed_cont_dur_dis_risk_prem')
    ALTER TABLE member_rating_results ADD exp_adj_gla_ed_cont_dur_dis_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'gla_ed_cont_dur_dis_office_premium')
    ALTER TABLE member_rating_results ADD gla_ed_cont_dur_dis_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_gla_ed_cont_dur_dis_office_prem')
    ALTER TABLE member_rating_results ADD exp_adj_gla_ed_cont_dur_dis_office_prem FLOAT;
-- Slice 4: PTD conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD ptd_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD ptd_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD ptd_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_conv_on_wdr_office_premium FLOAT;
-- Slice 3: PTD Educator conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD ptd_ed_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD ptd_ed_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_ed_conv_on_wdr_risk_prem')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_ed_conv_on_wdr_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD ptd_ed_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_ed_conv_on_wdr_office_prem')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_ed_conv_on_wdr_office_prem FLOAT;
-- Slice 11: PTD Educator conversion on retirement
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_ed_conv_on_ret_loading')
    ALTER TABLE member_rating_results ADD ptd_ed_conv_on_ret_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_ed_conv_on_ret_risk_premium')
    ALTER TABLE member_rating_results ADD ptd_ed_conv_on_ret_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_ed_conv_on_ret_risk_prem')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_ed_conv_on_ret_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ptd_ed_conv_on_ret_office_premium')
    ALTER TABLE member_rating_results ADD ptd_ed_conv_on_ret_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ptd_ed_conv_on_ret_office_prem')
    ALTER TABLE member_rating_results ADD exp_adj_ptd_ed_conv_on_ret_office_prem FLOAT;
-- Slice 6: CI conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD ci_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD ci_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ci_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ci_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'ci_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD ci_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_ci_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_ci_conv_on_wdr_office_premium FLOAT;
-- Slice 5: PHI conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD phi_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD phi_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_phi_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_phi_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'phi_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD phi_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_phi_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_phi_conv_on_wdr_office_premium FLOAT;
-- Slice 7: SGLA conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'sgla_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD sgla_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'sgla_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD sgla_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_sgla_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_sgla_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'sgla_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD sgla_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_sgla_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_sgla_conv_on_wdr_office_premium FLOAT;
-- Slice 8: Funeral conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'fun_conv_on_wdr_loading')
    ALTER TABLE member_rating_results ADD fun_conv_on_wdr_loading FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'fun_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD fun_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_fun_conv_on_wdr_risk_premium')
    ALTER TABLE member_rating_results ADD exp_adj_fun_conv_on_wdr_risk_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'fun_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD fun_conv_on_wdr_office_premium FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'exp_adj_fun_conv_on_wdr_office_premium')
    ALTER TABLE member_rating_results ADD exp_adj_fun_conv_on_wdr_office_premium FLOAT;

-- ── member_rating_result_summaries: summary buildup (156) + funeral SA ─
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (id INT IDENTITY(1,1) PRIMARY KEY);
END;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_family_funeral_sum_assured')
    ALTER TABLE member_rating_result_summaries ADD total_family_funeral_sum_assured FLOAT;

-- Macro for each slice: 4 totals + 4 proportions + 4 rates-per-1000.
-- Slice 1: GLA conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_conv_on_wdr_office_rate_per_1000_sa FLOAT;
-- Slice 9: GLA conversion on retirement
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_conv_on_ret_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_conv_on_ret_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_conv_on_ret_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_conv_on_ret_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_conv_on_ret_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_conv_on_ret_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_conv_on_ret_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_conv_on_ret_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_conv_on_ret_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_conv_on_ret_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_conv_on_ret_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_conv_on_ret_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_conv_on_ret_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_conv_on_ret_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_conv_on_ret_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_conv_on_ret_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_conv_on_ret_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_conv_on_ret_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_conv_on_ret_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_conv_on_ret_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_conv_on_ret_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_conv_on_ret_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_conv_on_ret_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_conv_on_ret_office_rate_per_1000_sa FLOAT;
-- Slice 12: GLA continuity during disability
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_cont_dur_dis_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_cont_dur_dis_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_cont_dur_dis_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_cont_dur_dis_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_cont_dur_dis_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_cont_dur_dis_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_cont_dur_dis_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_cont_dur_dis_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_cont_dur_dis_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_cont_dur_dis_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_cont_dur_dis_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_cont_dur_dis_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_cont_dur_dis_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_cont_dur_dis_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_cont_dur_dis_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_cont_dur_dis_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_cont_dur_dis_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_cont_dur_dis_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_cont_dur_dis_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_cont_dur_dis_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_cont_dur_dis_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_cont_dur_dis_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_cont_dur_dis_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_cont_dur_dis_office_rate_per_1000_sa FLOAT;
-- Slice 2: GLA Educator conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_ed_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_ed_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_ed_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_ed_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_ed_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_ed_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_ed_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_ed_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_ed_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_ed_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_ed_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_ed_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_ed_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_ed_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_ed_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_ed_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_ed_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_ed_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_ed_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_ed_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_ed_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_ed_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_ed_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_ed_conv_on_wdr_office_rate_per_1000_sa FLOAT;
-- Slice 10: GLA Educator conversion on retirement
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_ed_conv_on_ret_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_ed_conv_on_ret_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_ed_conv_on_ret_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_ed_conv_on_ret_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_ed_conv_on_ret_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_ed_conv_on_ret_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_ed_conv_on_ret_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_ed_conv_on_ret_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_ed_conv_on_ret_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_ed_conv_on_ret_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_ed_conv_on_ret_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_ed_conv_on_ret_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_ed_conv_on_ret_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_ed_conv_on_ret_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_ed_conv_on_ret_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_ed_conv_on_ret_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_ed_conv_on_ret_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_ed_conv_on_ret_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_ed_conv_on_ret_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_ed_conv_on_ret_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_ed_conv_on_ret_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_ed_conv_on_ret_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_ed_conv_on_ret_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_ed_conv_on_ret_office_rate_per_1000_sa FLOAT;
-- Slice 13: GLA Educator continuity during disability
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_ed_cont_dur_dis_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_ed_cont_dur_dis_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_gla_ed_cont_dur_dis_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_gla_ed_cont_dur_dis_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_ed_cont_dur_dis_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_ed_cont_dur_dis_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_gla_ed_cont_dur_dis_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_gla_ed_cont_dur_dis_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_ed_cont_dur_dis_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_ed_cont_dur_dis_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_gla_ed_cont_dur_dis_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_gla_ed_cont_dur_dis_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_ed_cont_dur_dis_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_ed_cont_dur_dis_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_gla_ed_cont_dur_dis_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_gla_ed_cont_dur_dis_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_ed_cont_dur_dis_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_ed_cont_dur_dis_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'gla_ed_cont_dur_dis_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD gla_ed_cont_dur_dis_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_ed_cont_dur_dis_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_ed_cont_dur_dis_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_gla_ed_cont_dur_dis_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_gla_ed_cont_dur_dis_office_rate_per_1000_sa FLOAT;
-- Slice 4: PTD conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ptd_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ptd_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ptd_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ptd_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ptd_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ptd_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ptd_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ptd_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_conv_on_wdr_office_rate_per_1000_sa FLOAT;
-- Slice 3: PTD Educator conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_ed_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_ed_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_ed_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_ed_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_ed_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_ed_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_ed_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_ed_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ptd_ed_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ptd_ed_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ptd_ed_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ptd_ed_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ptd_ed_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ptd_ed_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ptd_ed_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ptd_ed_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_ed_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_ed_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_ed_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_ed_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_ed_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_ed_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_ed_conv_on_wdr_office_rate_per_1000_sa FLOAT;
-- Slice 11: PTD Educator conversion on retirement
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_ed_conv_on_ret_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_ed_conv_on_ret_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ptd_ed_conv_on_ret_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ptd_ed_conv_on_ret_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_ed_conv_on_ret_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_ed_conv_on_ret_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ptd_ed_conv_on_ret_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ptd_ed_conv_on_ret_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ptd_ed_conv_on_ret_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ptd_ed_conv_on_ret_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ptd_ed_conv_on_ret_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ptd_ed_conv_on_ret_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ptd_ed_conv_on_ret_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ptd_ed_conv_on_ret_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ptd_ed_conv_on_ret_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ptd_ed_conv_on_ret_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_ed_conv_on_ret_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_ed_conv_on_ret_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ptd_ed_conv_on_ret_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ptd_ed_conv_on_ret_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_ed_conv_on_ret_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_ed_conv_on_ret_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ptd_ed_conv_on_ret_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ptd_ed_conv_on_ret_office_rate_per_1000_sa FLOAT;
-- Slice 5: PHI conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_phi_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_phi_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_phi_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_phi_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_phi_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_phi_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_phi_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_phi_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_phi_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_phi_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_phi_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_phi_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_phi_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_phi_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_phi_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_phi_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'phi_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD phi_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'phi_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD phi_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_phi_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_phi_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_phi_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_phi_conv_on_wdr_office_rate_per_1000_sa FLOAT;
-- Slice 6: CI conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ci_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ci_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_ci_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_ci_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ci_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ci_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_ci_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_ci_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ci_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ci_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_ci_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_ci_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ci_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ci_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_ci_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_ci_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ci_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ci_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'ci_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD ci_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ci_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ci_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_ci_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_ci_conv_on_wdr_office_rate_per_1000_sa FLOAT;
-- Slice 7: SGLA conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_sgla_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_sgla_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_sgla_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_sgla_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_sgla_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_sgla_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_sgla_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_sgla_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_sgla_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_sgla_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_sgla_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_sgla_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_sgla_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_sgla_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_sgla_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_sgla_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'sgla_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD sgla_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'sgla_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD sgla_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_sgla_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_sgla_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_sgla_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_sgla_conv_on_wdr_office_rate_per_1000_sa FLOAT;
-- Slice 8: Funeral conversion on withdrawal
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_fun_conv_on_wdr_annual_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD total_fun_conv_on_wdr_annual_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'total_fun_conv_on_wdr_annual_office_prem')
    ALTER TABLE member_rating_result_summaries ADD total_fun_conv_on_wdr_annual_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_fun_conv_on_wdr_ann_risk_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_fun_conv_on_wdr_ann_risk_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_total_fun_conv_on_wdr_ann_office_prem')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_fun_conv_on_wdr_ann_office_prem FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_fun_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_fun_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'prop_fun_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD prop_fun_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_fun_conv_on_wdr_risk_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_fun_conv_on_wdr_risk_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_adj_prop_fun_conv_on_wdr_office_prem_salary')
    ALTER TABLE member_rating_result_summaries ADD exp_adj_prop_fun_conv_on_wdr_office_prem_salary FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'fun_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD fun_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'fun_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD fun_conv_on_wdr_office_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_fun_conv_on_wdr_risk_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_fun_conv_on_wdr_risk_rate_per_1000_sa FLOAT;
IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_fun_conv_on_wdr_office_rate_per_1000_sa')
    ALTER TABLE member_rating_result_summaries ADD exp_fun_conv_on_wdr_office_rate_per_1000_sa FLOAT;

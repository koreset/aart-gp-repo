-- Migration: TTD conversion-on-withdrawal slice + GLA continuity-during-disability
-- dedicated loading. Generated to match the 20260421180000 mysql idempotent
-- ADD/MODIFY pattern.

-- ── scheme_categories: new flags ──────────────────────────────────────
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='ttd_conversion_on_withdrawal' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE scheme_categories MODIFY COLUMN ttd_conversion_on_withdrawal TINYINT(1) DEFAULT 0;', 'ALTER TABLE scheme_categories ADD COLUMN ttd_conversion_on_withdrawal TINYINT(1) DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='scheme_categories' AND COLUMN_NAME='gla_continuity_during_disability' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE scheme_categories MODIFY COLUMN gla_continuity_during_disability TINYINT(1) DEFAULT 0;', 'ALTER TABLE scheme_categories ADD COLUMN gla_continuity_during_disability TINYINT(1) DEFAULT 0;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ── general_loadings: new per-slice loading rates ─────────────────────
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='general_loadings' AND COLUMN_NAME='ttd_conv_on_wdr_loading_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE general_loadings MODIFY COLUMN ttd_conv_on_wdr_loading_rate DOUBLE;', 'ALTER TABLE general_loadings ADD COLUMN ttd_conv_on_wdr_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='general_loadings' AND COLUMN_NAME='gla_continuity_during_dis_loading_rate' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE general_loadings MODIFY COLUMN gla_continuity_during_dis_loading_rate DOUBLE;', 'ALTER TABLE general_loadings ADD COLUMN gla_continuity_during_dis_loading_rate DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ── member_rating_results: TTD slice + GLA continuity loading ─────────
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_conv_on_wdr_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ttd_conv_on_wdr_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ttd_conv_on_wdr_loading DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_conv_on_wdr_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ttd_conv_on_wdr_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ttd_conv_on_wdr_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_conv_on_wdr_risk_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ttd_conv_on_wdr_risk_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ttd_conv_on_wdr_risk_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_conv_on_wdr_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN ttd_conv_on_wdr_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN ttd_conv_on_wdr_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_conv_on_wdr_office_premium' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN exp_adj_ttd_conv_on_wdr_office_premium DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN exp_adj_ttd_conv_on_wdr_office_premium DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_continuity_during_dis_loading' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_results MODIFY COLUMN gla_continuity_during_dis_loading DOUBLE;', 'ALTER TABLE member_rating_results ADD COLUMN gla_continuity_during_dis_loading DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ── member_rating_result_summaries: TTD slice aggregates (4+4+4) ──────
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_conv_on_wdr_annual_risk_prem' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ttd_conv_on_wdr_annual_risk_prem DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ttd_conv_on_wdr_annual_risk_prem DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_conv_on_wdr_annual_office_prem' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN total_ttd_conv_on_wdr_annual_office_prem DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN total_ttd_conv_on_wdr_annual_office_prem DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ttd_conv_on_wdr_ann_risk_prem' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_ttd_conv_on_wdr_ann_risk_prem DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_ttd_conv_on_wdr_ann_risk_prem DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ttd_conv_on_wdr_ann_office_prem' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_total_ttd_conv_on_wdr_ann_office_prem DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_ttd_conv_on_wdr_ann_office_prem DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_ttd_conv_on_wdr_risk_prem_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN prop_ttd_conv_on_wdr_risk_prem_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_ttd_conv_on_wdr_risk_prem_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_ttd_conv_on_wdr_office_prem_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN prop_ttd_conv_on_wdr_office_prem_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_ttd_conv_on_wdr_office_prem_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_ttd_conv_on_wdr_risk_prem_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_prop_ttd_conv_on_wdr_risk_prem_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_prop_ttd_conv_on_wdr_risk_prem_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_ttd_conv_on_wdr_office_prem_salary' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_adj_prop_ttd_conv_on_wdr_office_prem_salary DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_prop_ttd_conv_on_wdr_office_prem_salary DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ttd_conv_on_wdr_risk_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN ttd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN ttd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ttd_conv_on_wdr_office_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN ttd_conv_on_wdr_office_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN ttd_conv_on_wdr_office_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ttd_conv_on_wdr_risk_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_ttd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_ttd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ttd_conv_on_wdr_office_rate_per_1000_sa' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_ttd_conv_on_wdr_office_rate_per_1000_sa DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_ttd_conv_on_wdr_office_rate_per_1000_sa DOUBLE;'));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

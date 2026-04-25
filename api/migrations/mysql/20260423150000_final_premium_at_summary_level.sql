-- Migration: collapse per-member Final* premiums and move Final premiums to
-- the summary only. On `member_rating_result_summaries`:
--   1) rename `final_adj_total_<x>_annual_office_premium` -> `final_<x>_office_premium`
--      (preserves data if old column exists)
--   2) ensure every Final*/ProportionFinal* column the model declares exists,
--      adding it with DOUBLE DEFAULT 0 when missing.
-- On `member_rating_results`, drop the old per-member Final* office-premium
-- columns — Final now lives on the summary.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);
CREATE TABLE IF NOT EXISTS member_rating_results (
    id INT AUTO_INCREMENT PRIMARY KEY
);

-- ---------------------------------------------------------------------------
-- Phase 1 — RENAMES on member_rating_result_summaries
-- Preserves existing data when the legacy Exp-adj name is present.
-- ---------------------------------------------------------------------------

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_gla_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_gla_annual_office_premium final_gla_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_tax_saver_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_tax_saver_annual_office_premium final_tax_saver_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_additional_accidental_gla_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_additional_accidental_gla_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_additional_accidental_gla_annual_office_premium final_additional_accidental_gla_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_ptd_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_ptd_annual_office_premium final_ptd_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_ci_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_ci_annual_office_premium final_ci_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_sgla_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_sgla_annual_office_premium final_sgla_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_ttd_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_ttd_annual_office_premium final_ttd_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_phi_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_phi_annual_office_premium final_phi_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_fun_annual_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_fun_annual_office_premium final_fun_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_gla_educator_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_gla_educator_office_premium final_gla_educator_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_adj_total_ptd_educator_office_premium')
    AND NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_office_premium'),
    'ALTER TABLE member_rating_result_summaries CHANGE COLUMN final_adj_total_ptd_educator_office_premium final_ptd_educator_office_premium DOUBLE;',
    'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- Phase 2 — ADD every Final*/ProportionFinal* column the model declares.
-- Idempotent: each statement is guarded by NOT EXISTS. The office_premium
-- columns are re-checked here so that when Phase-1 rename was a no-op
-- (legacy name not present) the column is still created.
-- ---------------------------------------------------------------------------

-- Base benefit: office_premium (fallback ADD for each Phase-1 rename target)
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_additional_accidental_gla_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_office_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_office_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Base benefit: risk_premium
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_additional_accidental_gla_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_risk_premium'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_risk_premium DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Base benefit: proportion office_premium_salary (AAGLA uses short gorm name)
-- proportion_final_gla_office_premium_salary is now computed inline in
-- quote_template/schema.go and is no longer persisted.
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_tax_saver_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_tax_saver_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_add_acc_gla_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_add_acc_gla_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ptd_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ci_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ci_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_sgla_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_sgla_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ttd_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ttd_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_phi_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_phi_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_fun_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_fun_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_gla_educator_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_gla_educator_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ptd_educator_office_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_educator_office_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Base benefit: proportion annual_risk_premium_salary (AAGLA uses short gorm name)
-- proportion_final_gla_annual_risk_premium_salary is now computed inline in
-- quote_template/schema.go and is no longer persisted.
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_tax_saver_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_tax_saver_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_add_acc_gla_ann_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_add_acc_gla_ann_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ptd_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ci_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ci_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_sgla_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_sgla_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ttd_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ttd_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_phi_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_phi_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_fun_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_fun_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_gla_educator_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_gla_educator_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_final_ptd_educator_annual_risk_premium_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_final_ptd_educator_annual_risk_premium_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Base benefit: office_rate_per_1000_sa
-- final_gla_office_rate_per1000_sa is now computed inline in
-- quote_template/schema.go and is no longer persisted.
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_additional_accidental_gla_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_office_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_office_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Base benefit: risk_rate_per_1000_sa
-- final_gla_risk_rate_per1000_sa is now computed inline in
-- quote_template/schema.go and is no longer persisted.
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_tax_saver_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_tax_saver_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_additional_accidental_gla_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_additional_accidental_gla_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_educator_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_educator_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_educator_risk_rate_per1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_educator_risk_rate_per1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- Phase 2b — Conversion / continuity slice columns (gorm-short names).
-- 14 (benefit × slice) pairings × 6 columns each.
-- ---------------------------------------------------------------------------

-- GLA conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- GLA conv_on_ret
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_ret_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_ret_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_conv_on_ret_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_ret_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_conv_on_ret_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_conv_on_ret_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_ret_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_conv_on_ret_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_conv_on_ret_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- GLA cont_dur_dis
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_cont_dur_dis_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_cont_dur_dis_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_cont_dur_dis_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_cont_dur_dis_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_cont_dur_dis_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_cont_dur_dis_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_cont_dur_dis_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_cont_dur_dis_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_cont_dur_dis_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- GLA Educator conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_ed_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_ed_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- GLA Educator conv_on_ret
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_ret_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_ret_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_ed_conv_on_ret_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_ret_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_ed_conv_on_ret_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_conv_on_ret_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_ret_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_conv_on_ret_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_conv_on_ret_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- GLA Educator cont_dur_dis
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_cont_dur_dis_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_cont_dur_dis_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_ed_cont_dur_dis_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_cont_dur_dis_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_gla_ed_cont_dur_dis_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_gla_ed_cont_dur_dis_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_cont_dur_dis_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_gla_ed_cont_dur_dis_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_gla_ed_cont_dur_dis_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- PTD conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ptd_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ptd_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- PTD Educator conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ptd_ed_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ptd_ed_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- PTD Educator conv_on_ret
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_ret_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_ret_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ptd_ed_conv_on_ret_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_ret_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ptd_ed_conv_on_ret_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ptd_ed_conv_on_ret_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_ret_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ptd_ed_conv_on_ret_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ptd_ed_conv_on_ret_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- PHI conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_phi_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_phi_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_phi_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_phi_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_phi_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_phi_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- TTD conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ttd_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ttd_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ttd_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ttd_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ttd_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ttd_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- CI conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ci_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ci_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_ci_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_ci_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_ci_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_ci_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- SGLA conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_sgla_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_sgla_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_sgla_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_sgla_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_sgla_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_sgla_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- FUN conv_on_wdr
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_conv_on_wdr_risk_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_risk_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_conv_on_wdr_office_prem'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_office_prem DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_fun_conv_on_wdr_risk_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_fun_conv_on_wdr_risk_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_final_fun_conv_on_wdr_office_prem_salary'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN prop_final_fun_conv_on_wdr_office_prem_salary DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_conv_on_wdr_risk_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_risk_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
SET @sql := (SELECT IF(NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='final_fun_conv_on_wdr_office_rate_per_1000_sa'),
    'ALTER TABLE member_rating_result_summaries ADD COLUMN final_fun_conv_on_wdr_office_rate_per_1000_sa DOUBLE DEFAULT 0;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- Phase 3 — DROPS on member_rating_results (Final now lives on the summary)
-- ---------------------------------------------------------------------------

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_gla_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_gla_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_tax_saver_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_tax_saver_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_additional_accidental_gla_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_additional_accidental_gla_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ptd_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_ptd_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ci_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_ci_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_spouse_gla_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_spouse_gla_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ttd_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_ttd_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_phi_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_phi_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_total_funeral_office_cost'),
    'ALTER TABLE member_rating_results DROP COLUMN final_total_funeral_office_cost;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_gla_educator_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_gla_educator_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ptd_educator_office_premium'),
    'ALTER TABLE member_rating_results DROP COLUMN final_ptd_educator_office_premium;', 'SELECT 1;'));
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

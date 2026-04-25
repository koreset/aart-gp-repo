-- Migration: drop persisted office-premium and office-rate-per-1000-SA columns
-- from member_rating_results and member_rating_result_summaries.
-- Office premium is now derived on the fly from the risk premium and the
-- scheme-level loading (expense + commission + profit) on the summary.

-- ── member_rating_results ───────────────────────────────────────────
SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='tax_saver_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN tax_saver_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_tax_saver_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_tax_saver_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_tax_saver_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_tax_saver_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='additional_accidental_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN additional_accidental_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_additional_accidental_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_additional_accidental_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_additional_accidental_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_additional_accidental_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ptd_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ptd_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_ptd_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ci_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ci_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ci_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ci_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_ci_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN spouse_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_spouse_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_spouse_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_spouse_gla_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_spouse_gla_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ttd_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ttd_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ttd_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_ttd_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN phi_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_phi_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_phi_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_phi_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_phi_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='main_member_funeral_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN main_member_funeral_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='spouse_funeral_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN spouse_funeral_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='children_funeral_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN children_funeral_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='dependants_funeral_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN dependants_funeral_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_educator_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_educator_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_gla_educator_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_gla_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_educator_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ptd_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_educator_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='final_ptd_educator_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN final_ptd_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_conv_on_ret_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_conv_on_ret_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_conv_on_ret_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_conv_on_ret_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_cont_dur_dis_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_cont_dur_dis_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_cont_dur_dis_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_cont_dur_dis_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_ed_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_ed_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_ed_conv_on_wdr_office_prem'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_ed_conv_on_wdr_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_ed_conv_on_ret_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_ed_conv_on_ret_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_ed_conv_on_ret_office_prem'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_ed_conv_on_ret_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='gla_ed_cont_dur_dis_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN gla_ed_cont_dur_dis_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_gla_ed_cont_dur_dis_office_prem'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_ed_cont_dur_dis_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ptd_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_ed_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ptd_ed_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_ed_conv_on_wdr_office_prem'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_ed_conv_on_wdr_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ptd_ed_conv_on_ret_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ptd_ed_conv_on_ret_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ptd_ed_conv_on_ret_office_prem'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_ed_conv_on_ret_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ci_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ci_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ci_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ci_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='phi_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN phi_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_phi_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_phi_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='sgla_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN sgla_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_sgla_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_sgla_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='fun_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN fun_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_fun_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_fun_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='ttd_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN ttd_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_results' AND COLUMN_NAME='exp_adj_ttd_conv_on_wdr_office_premium'),
        'ALTER TABLE member_rating_results DROP COLUMN exp_adj_ttd_conv_on_wdr_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ── member_rating_result_summaries ──────────────────────────────────
SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_gla_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_gla_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_gla_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_gla_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_gla_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_gla_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_tax_saver_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_tax_saver_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_tax_saver_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_tax_saver_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_additional_accidental_gla_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_additional_accidental_gla_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='additional_accidental_gla_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN additional_accidental_gla_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_additional_accidental_gla_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_additional_accidental_gla_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_additional_accidental_gla_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_additional_accidental_gla_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_additional_accidental_gla_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_additional_accidental_gla_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_additional_accidental_gla_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_additional_accidental_gla_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_ptd_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ptd_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ptd_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_ptd_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_ptd_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_ptd_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ci_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ci_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ci_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ci_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_ci_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ci_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ci_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_ci_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ci_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ci_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_ci_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_ci_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_sgla_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_sgla_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='sgla_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN sgla_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_sgla_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_sgla_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_sgla_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_sgla_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_sgla_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_sgla_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_sgla_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_sgla_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ttd_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ttd_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ttd_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_ttd_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ttd_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_ttd_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_ttd_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ttd_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ttd_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_ttd_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_ttd_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_phi_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_phi_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='phi_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN phi_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_phi_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_phi_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_phi_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_phi_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_phi_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_phi_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_phi_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_phi_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_fun_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_fun_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_fun_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_fun_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_fun_annual_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_fun_annual_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_proportion_fun_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_fun_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_educator_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_educator_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_gla_educator_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_gla_educator_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_proportion_gla_educator_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_proportion_gla_educator_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_educator_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_educator_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_educator_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_educator_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_educator_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_educator_office_premium'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_educator_office_premium;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='proportion_ptd_educator_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ptd_educator_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_proportion_ptd_educator_office_premium_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_proportion_ptd_educator_office_premium_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_educator_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_educator_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_educator_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_educator_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_gla_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_gla_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_conv_on_ret_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_conv_on_ret_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_conv_on_ret_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_conv_on_ret_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_gla_conv_on_ret_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_conv_on_ret_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_gla_conv_on_ret_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_conv_on_ret_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_conv_on_ret_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_conv_on_ret_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_conv_on_ret_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_conv_on_ret_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_cont_dur_dis_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_cont_dur_dis_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_cont_dur_dis_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_cont_dur_dis_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_gla_cont_dur_dis_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_cont_dur_dis_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_gla_cont_dur_dis_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_cont_dur_dis_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_cont_dur_dis_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_cont_dur_dis_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_cont_dur_dis_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_cont_dur_dis_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_ed_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_ed_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_ed_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_ed_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_gla_ed_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_ed_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_gla_ed_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_ed_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_ed_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_ed_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_ed_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_ed_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_ed_conv_on_ret_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_ed_conv_on_ret_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_ed_conv_on_ret_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_ed_conv_on_ret_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_gla_ed_conv_on_ret_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_ed_conv_on_ret_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_gla_ed_conv_on_ret_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_ed_conv_on_ret_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_ed_conv_on_ret_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_ed_conv_on_ret_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_ed_conv_on_ret_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_ed_conv_on_ret_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_gla_ed_cont_dur_dis_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_ed_cont_dur_dis_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_gla_ed_cont_dur_dis_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_ed_cont_dur_dis_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_gla_ed_cont_dur_dis_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_ed_cont_dur_dis_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_gla_ed_cont_dur_dis_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_ed_cont_dur_dis_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='gla_ed_cont_dur_dis_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN gla_ed_cont_dur_dis_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_gla_ed_cont_dur_dis_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_ed_cont_dur_dis_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_ptd_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ptd_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_ptd_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ptd_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_ed_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_ed_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_ed_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_ed_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_ptd_ed_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ptd_ed_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_ptd_ed_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ptd_ed_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_ed_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_ed_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_ed_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_ed_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ptd_ed_conv_on_ret_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_ed_conv_on_ret_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ptd_ed_conv_on_ret_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_ed_conv_on_ret_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_ptd_ed_conv_on_ret_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ptd_ed_conv_on_ret_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_ptd_ed_conv_on_ret_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ptd_ed_conv_on_ret_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ptd_ed_conv_on_ret_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_ed_conv_on_ret_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ptd_ed_conv_on_ret_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_ed_conv_on_ret_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_phi_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_phi_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_phi_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_phi_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_phi_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_phi_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_phi_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_phi_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='phi_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN phi_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_phi_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_phi_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ttd_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ttd_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ttd_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ttd_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_ttd_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ttd_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_ttd_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ttd_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ttd_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ttd_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ttd_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ttd_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_ci_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_ci_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_ci_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ci_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_ci_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ci_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_ci_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ci_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='ci_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN ci_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_ci_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ci_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_sgla_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_sgla_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_sgla_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_sgla_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_sgla_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_sgla_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_sgla_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_sgla_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='sgla_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN sgla_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_sgla_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_sgla_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='total_fun_conv_on_wdr_annual_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN total_fun_conv_on_wdr_annual_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_total_fun_conv_on_wdr_ann_office_prem'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_fun_conv_on_wdr_ann_office_prem;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='prop_fun_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN prop_fun_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_adj_prop_fun_conv_on_wdr_office_prem_salary'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_fun_conv_on_wdr_office_prem_salary;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='fun_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN fun_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_fun_conv_on_wdr_office_rate_per_1000_sa'),
        'ALTER TABLE member_rating_result_summaries DROP COLUMN exp_fun_conv_on_wdr_office_rate_per_1000_sa;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

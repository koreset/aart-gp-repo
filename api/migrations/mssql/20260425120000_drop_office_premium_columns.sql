-- Migration: drop persisted office-premium and office-rate-per-1000-SA columns
-- from member_rating_results and member_rating_result_summaries.
-- Office premium is now derived on the fly from the risk premium and the
-- scheme-level loading (expense + commission + profit) on the summary.

-- ── member_rating_results ───────────────────────────────────────────
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'tax_saver_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN tax_saver_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_tax_saver_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_tax_saver_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_tax_saver_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_tax_saver_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'additional_accidental_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN additional_accidental_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_additional_accidental_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_additional_accidental_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_additional_accidental_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_additional_accidental_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ptd_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ptd_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ptd_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_ptd_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_ptd_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ci_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ci_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ci_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ci_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_ci_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_ci_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'spouse_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN spouse_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_spouse_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_spouse_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_spouse_gla_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_spouse_gla_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ttd_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ttd_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ttd_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ttd_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_ttd_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_ttd_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'phi_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN phi_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_phi_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_phi_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_phi_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_phi_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'main_member_funeral_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN main_member_funeral_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'spouse_funeral_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN spouse_funeral_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'children_funeral_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN children_funeral_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'dependants_funeral_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN dependants_funeral_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_educator_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_educator_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_gla_educator_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_gla_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ptd_educator_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ptd_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ptd_educator_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_ptd_educator_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN final_ptd_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_conv_on_ret_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_conv_on_ret_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_conv_on_ret_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_conv_on_ret_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_cont_dur_dis_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_cont_dur_dis_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_cont_dur_dis_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_cont_dur_dis_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_ed_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_ed_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_ed_conv_on_wdr_office_prem')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_ed_conv_on_wdr_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_ed_conv_on_ret_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_ed_conv_on_ret_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_ed_conv_on_ret_office_prem')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_ed_conv_on_ret_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'gla_ed_cont_dur_dis_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN gla_ed_cont_dur_dis_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_gla_ed_cont_dur_dis_office_prem')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_gla_ed_cont_dur_dis_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ptd_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ptd_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ptd_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ptd_ed_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ptd_ed_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ptd_ed_conv_on_wdr_office_prem')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_ed_conv_on_wdr_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ptd_ed_conv_on_ret_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ptd_ed_conv_on_ret_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ptd_ed_conv_on_ret_office_prem')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ptd_ed_conv_on_ret_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ci_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ci_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ci_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ci_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'phi_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN phi_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_phi_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_phi_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'sgla_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN sgla_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_sgla_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_sgla_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'fun_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN fun_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_fun_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_fun_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'ttd_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN ttd_conv_on_wdr_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_ttd_conv_on_wdr_office_premium')
BEGIN ALTER TABLE member_rating_results DROP COLUMN exp_adj_ttd_conv_on_wdr_office_premium; END;

-- ── member_rating_result_summaries ──────────────────────────────────
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_gla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_gla_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_gla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_gla_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_gla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_gla_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_tax_saver_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_tax_saver_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_tax_saver_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_tax_saver_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_additional_accidental_gla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_additional_accidental_gla_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'additional_accidental_gla_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN additional_accidental_gla_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_additional_accidental_gla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_additional_accidental_gla_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_additional_accidental_gla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_additional_accidental_gla_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_additional_accidental_gla_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_additional_accidental_gla_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_additional_accidental_gla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_additional_accidental_gla_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ptd_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_ptd_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ptd_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_ptd_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_ptd_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_ptd_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_ptd_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ci_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ci_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ci_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ci_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_ci_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ci_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_ci_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_ci_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ci_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ci_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_ci_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_ci_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_sgla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_sgla_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'sgla_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN sgla_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_sgla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_sgla_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_sgla_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_sgla_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_sgla_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_sgla_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_sgla_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_sgla_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ttd_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ttd_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ttd_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ttd_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_ttd_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ttd_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_ttd_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_ttd_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ttd_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ttd_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_ttd_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_ttd_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_phi_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_phi_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'phi_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN phi_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_phi_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_phi_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_phi_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_phi_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_phi_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_phi_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_phi_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_phi_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_fun_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_fun_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_fun_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_fun_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_total_fun_annual_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_total_fun_annual_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_proportion_fun_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_proportion_fun_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_educator_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_gla_educator_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_gla_educator_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_gla_educator_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_proportion_gla_educator_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_proportion_gla_educator_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_educator_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_educator_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_educator_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_educator_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ptd_educator_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_ptd_educator_office_premium')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_educator_office_premium; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'proportion_ptd_educator_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN proportion_ptd_educator_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_proportion_ptd_educator_office_premium_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_proportion_ptd_educator_office_premium_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_educator_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_educator_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_educator_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_educator_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_gla_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_gla_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_gla_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_conv_on_ret_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_conv_on_ret_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_gla_conv_on_ret_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_conv_on_ret_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_gla_conv_on_ret_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_conv_on_ret_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_gla_conv_on_ret_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_conv_on_ret_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_conv_on_ret_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_conv_on_ret_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_conv_on_ret_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_conv_on_ret_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_cont_dur_dis_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_cont_dur_dis_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_gla_cont_dur_dis_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_cont_dur_dis_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_gla_cont_dur_dis_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_cont_dur_dis_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_gla_cont_dur_dis_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_cont_dur_dis_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_cont_dur_dis_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_cont_dur_dis_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_cont_dur_dis_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_cont_dur_dis_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_ed_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_ed_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_gla_ed_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_ed_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_gla_ed_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_ed_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_gla_ed_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_ed_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_ed_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_ed_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_ed_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_ed_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_ed_conv_on_ret_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_ed_conv_on_ret_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_gla_ed_conv_on_ret_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_ed_conv_on_ret_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_gla_ed_conv_on_ret_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_ed_conv_on_ret_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_gla_ed_conv_on_ret_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_ed_conv_on_ret_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_ed_conv_on_ret_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_ed_conv_on_ret_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_ed_conv_on_ret_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_ed_conv_on_ret_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_gla_ed_cont_dur_dis_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_gla_ed_cont_dur_dis_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_gla_ed_cont_dur_dis_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_gla_ed_cont_dur_dis_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_gla_ed_cont_dur_dis_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_gla_ed_cont_dur_dis_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_gla_ed_cont_dur_dis_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_gla_ed_cont_dur_dis_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_ed_cont_dur_dis_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN gla_ed_cont_dur_dis_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_gla_ed_cont_dur_dis_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_gla_ed_cont_dur_dis_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ptd_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_ptd_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_ptd_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ptd_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_ptd_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ptd_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ptd_ed_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_ed_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_ptd_ed_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_ed_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_ptd_ed_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ptd_ed_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_ptd_ed_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ptd_ed_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_ed_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_ed_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_ed_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_ed_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ptd_ed_conv_on_ret_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ptd_ed_conv_on_ret_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_ptd_ed_conv_on_ret_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ptd_ed_conv_on_ret_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_ptd_ed_conv_on_ret_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ptd_ed_conv_on_ret_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_ptd_ed_conv_on_ret_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ptd_ed_conv_on_ret_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ptd_ed_conv_on_ret_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ptd_ed_conv_on_ret_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ptd_ed_conv_on_ret_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ptd_ed_conv_on_ret_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_phi_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_phi_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_phi_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_phi_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_phi_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_phi_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_phi_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_phi_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'phi_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN phi_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_phi_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_phi_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ttd_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ttd_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_ttd_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ttd_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_ttd_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ttd_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_ttd_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ttd_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ttd_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ttd_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ttd_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ttd_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_ci_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_ci_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_ci_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_ci_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_ci_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_ci_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_ci_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_ci_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'ci_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN ci_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_ci_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_ci_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_sgla_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_sgla_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_sgla_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_sgla_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_sgla_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_sgla_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_sgla_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_sgla_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'sgla_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN sgla_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_sgla_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_sgla_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_fun_conv_on_wdr_annual_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN total_fun_conv_on_wdr_annual_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_fun_conv_on_wdr_ann_office_prem')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_total_fun_conv_on_wdr_ann_office_prem; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'prop_fun_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN prop_fun_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_prop_fun_conv_on_wdr_office_prem_salary')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_adj_prop_fun_conv_on_wdr_office_prem_salary; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'fun_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN fun_conv_on_wdr_office_rate_per_1000_sa; END;
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_fun_conv_on_wdr_office_rate_per_1000_sa')
BEGIN ALTER TABLE member_rating_result_summaries DROP COLUMN exp_fun_conv_on_wdr_office_rate_per_1000_sa; END;

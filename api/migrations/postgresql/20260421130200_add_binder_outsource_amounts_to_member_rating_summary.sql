-- Migration: add per-benefit binder and outsource fee aggregate columns to
-- member_rating_result_summaries. Mirrors the office-premium totals so the
-- summary UI and downstream reporting can break the office premium into its
-- binder and outsource slices when distribution_channel = 'binder'.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_gla_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_gla_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_add_acc_gla_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_add_acc_gla_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_add_acc_gla_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_add_acc_gla_annual_outsourced_amt NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ptd_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ptd_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ci_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ci_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ci_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ci_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_sgla_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_sgla_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_sgla_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_sgla_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ttd_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ttd_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ttd_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ttd_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_phi_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_phi_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_phi_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_phi_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_fun_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_fun_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_fun_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_fun_annual_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_annual_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_annual_outsourced_amount NUMERIC(20,6);

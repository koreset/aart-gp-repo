-- Migration: add per-benefit commission amount columns and scheme-wide
-- commission totals to member_rating_result_summaries. Commission is
-- computed on the scheme-wide total premium via the tiered
-- CommissionStructure bands and distributed proportionally per benefit.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_gla_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_add_acc_gla_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ptd_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ci_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_sgla_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_ttd_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_phi_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_fun_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_educator_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS scheme_total_commission NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS scheme_total_commission_rate NUMERIC(20,6);

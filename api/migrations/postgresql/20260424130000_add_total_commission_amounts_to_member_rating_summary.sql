-- Migration: add the non-Exp Total*CommissionAmount mirrors plus the per-category
-- ExpTotalCommission roll-up to member_rating_result_summaries. Parallel to the
-- existing Exp*CommissionAmount columns added in 20260421150000 so the scheme-wide
-- commission can be distributed on both the Exp and Total premium bases.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_annual_premium_excl_funeral NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_annual_premium_excl_funeral NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_commission NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_add_acc_gla_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ci_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_sgla_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ttd_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_phi_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_fun_annual_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_educator_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_educator_commission_amount NUMERIC(20,6);

-- Migration: add Discount + Final* office-premium and commission columns to
-- member_rating_result_summaries. Final*OfficePremium is persisted as the
-- post-discount pre-comm office premium (Exp*RiskPremium /
-- (1 - (SchemeTotalLoading + Discount))) plus its per-benefit commission
-- slice, so the Final* values include commission and reconcile to
-- final_total_annual_premium{,_excl_funeral}. Exp* values stay pre-commission
-- and frozen at calc time, so pre-discount Final - Exp == commission slice.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS discount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_add_acc_gla_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ci_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_sgla_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_tax_saver_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ttd_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_phi_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_fun_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_educator_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_educator_annual_office_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_add_acc_gla_annual_comm_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ci_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_sgla_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_tax_saver_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ttd_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_phi_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_fun_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_educator_annual_comm_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_educator_annual_comm_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_total_annual_premium_excl_funeral NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_total_annual_premium NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_scheme_total_commission NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_scheme_total_commission_rate NUMERIC(20,6) DEFAULT 0;

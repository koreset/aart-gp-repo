-- Migration: add Final*BinderAmount + Final*OutsourcedAmount columns to
-- member_rating_result_summaries. These mirror the existing Exp*Binder /
-- Exp*Outsourced fields but are derived from the post-discount Final office
-- premium so the breakdown reconciles to FinalOfficePremium after a discount
-- is applied. Pre-discount they equal the Exp* values; on non-binder
-- distribution channels they remain 0.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_add_acc_gla_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_add_acc_gla_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ci_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ci_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_sgla_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_sgla_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_tax_saver_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_tax_saver_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ttd_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ttd_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_phi_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_phi_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_fun_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_fun_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_educator_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_gla_educator_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_educator_annual_binder_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS final_ptd_educator_annual_outsourced_amount NUMERIC(20,6) DEFAULT 0;

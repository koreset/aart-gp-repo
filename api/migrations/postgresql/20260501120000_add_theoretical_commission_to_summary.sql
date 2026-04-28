-- Migration: persist book-rate (theoretical) commission slices on
-- member_rating_result_summaries. Mirrors the Exp*CommissionAmount columns but
-- uses the book scheme-total premium to derive its progressive commission rate
-- so the OutputSummary Theoretical column can show a commission-inclusive
-- office premium.

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_add_acc_gla_annual_comm_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ci_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_sgla_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_tax_saver_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ttd_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_phi_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_fun_annual_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_educator_commission_amount NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_educator_commission_amount NUMERIC(20,6) DEFAULT 0;

-- Migration: add reinsurance rate/loading/premium fields to member_rating_results
-- Adds the 62 columns introduced alongside the reinsurance rate tables so that
-- GORM INSERTs from MemberRatingResult succeed against existing DBs.

CREATE TABLE IF NOT EXISTS member_rating_results (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_industry_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ptd_industry_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ci_industry_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ttd_industry_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_phi_industry_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_aids_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ptd_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ci_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ttd_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_phi_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_fun_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_fun_aids_region_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_contingency_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ptd_contingency_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ci_contingency_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ttd_contingency_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_phi_contingency_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_fun_contingency_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_continuation_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_terminal_illness_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_qx NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_gla_aids_qx NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_reins_gla_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_reins_gla_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ptd_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_reins_ptd_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_reins_ptd_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_ci_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_reins_ci_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_reins_ci_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_reins_ttd_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_reins_ttd_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_phi_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_reins_phi_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_reins_phi_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_spouse_gla_qx NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_spouse_gla_aids_qx NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS reins_spouse_gla_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS base_reins_spouse_gla_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS loaded_reins_spouse_gla_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_reinsurance_base_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_reinsurance_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_reinsurance_base_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_reinsurance_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS child_reinsurance_base_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS child_reinsurance_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS parent_reinsurance_base_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS parent_reinsurance_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependant_reinsurance_base_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependant_reinsurance_rate NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS child_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS parent_reinsurance_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependant_reinsurance_premium NUMERIC(20,6);

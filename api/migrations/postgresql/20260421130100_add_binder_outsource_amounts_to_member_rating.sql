-- Migration: add per-benefit binder and outsource fee amount columns to
-- member_rating_results. Each pair breaks the corresponding office premium
-- into its binder-fee and outsource-fee slices when the quote is sold
-- through the binder distribution channel.

CREATE TABLE IF NOT EXISTS member_rating_results (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS additional_accidental_gla_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_add_acc_gla_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_add_acc_gla_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ci_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ci_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ci_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_gla_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_spouse_gla_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_spouse_gla_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ttd_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ttd_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ttd_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS phi_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_phi_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_phi_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_funeral_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS main_member_funeral_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_funeral_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS spouse_funeral_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS children_funeral_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS children_funeral_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependants_funeral_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS dependants_funeral_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS total_funeral_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS total_funeral_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_total_funeral_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_total_funeral_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS total_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS total_outsourced_amount NUMERIC(20,6);

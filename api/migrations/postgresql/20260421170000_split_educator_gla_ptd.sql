-- Migration: split educator benefit tracking into GLA and PTD components.
-- The combined Educator* columns stay as the sum; the new columns let the
-- business attribute the educator premium between GLA-educator and
-- PTD-educator and expose buildup fields (premium, %salary, rate per 1000).

CREATE TABLE IF NOT EXISTS member_rating_results (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_educator_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_educator_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_educator_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_educator_office_premium NUMERIC(20,6);

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_educator_sum_assured NUMERIC(20,6);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_educator_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_educator_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS proportion_gla_educator_risk_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS proportion_gla_educator_office_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_proportion_gla_educator_risk_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_proportion_gla_educator_office_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_educator_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS gla_educator_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_educator_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_gla_educator_office_rate_per_1000_sa NUMERIC(20,6);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_educator_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_educator_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_educator_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS proportion_ptd_educator_risk_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS proportion_ptd_educator_office_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_proportion_ptd_educator_risk_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_proportion_ptd_educator_office_premium_salary NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_educator_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS ptd_educator_office_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_educator_risk_rate_per_1000_sa NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_ptd_educator_office_rate_per_1000_sa NUMERIC(20,6);

-- Split binder / outsource per-member columns
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS gla_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_gla_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS ptd_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_ptd_educator_outsourced_amount NUMERIC(20,6);

-- Split binder / outsource summary totals
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_gla_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_ptd_educator_outsourced_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_educator_binder_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_educator_outsourced_amount NUMERIC(20,6);

-- Split commission summary totals
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_gla_educator_commission_amount NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_adj_total_ptd_educator_commission_amount NUMERIC(20,6);

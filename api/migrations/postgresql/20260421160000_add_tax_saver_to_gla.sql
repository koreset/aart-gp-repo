-- Migration: add optional TaxSaver rider fields. The tax-saver loading
-- itself lives on general_loadings (per age/gender/risk_rate_code), so
-- scheme_categories only carries the opt-in flag.

CREATE TABLE IF NOT EXISTS scheme_categories (
    id SERIAL PRIMARY KEY
);

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS tax_saver_benefit BOOLEAN DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS general_loadings (
    id SERIAL PRIMARY KEY
);

ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS tax_saver_loading_rate NUMERIC(20,6);

CREATE TABLE IF NOT EXISTS member_rating_results (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS tax_saver_loading NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS tax_saver_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_tax_saver_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS tax_saver_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS exp_adj_tax_saver_office_premium NUMERIC(20,6);

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS tax_saver_benefit BOOLEAN DEFAULT FALSE;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_tax_saver_annual_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS total_tax_saver_annual_office_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_tax_saver_annual_risk_premium NUMERIC(20,6);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_tax_saver_annual_office_premium NUMERIC(20,6);

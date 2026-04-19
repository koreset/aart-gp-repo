-- Migration for struct: ReinsuranceIndustryLoading
-- Table: reinsurance_industry_loadings

CREATE TABLE IF NOT EXISTS reinsurance_industry_loadings (
    id SERIAL PRIMARY KEY
);

ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS occupation_class INTEGER;
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS gender VARCHAR(255);
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS gla_industry_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS ptd_industry_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS ci_industry_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS ttd_industry_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS phi_industry_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP;
ALTER TABLE reinsurance_industry_loadings ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

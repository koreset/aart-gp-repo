-- Migration for struct: ReinsuranceRegionLoading
-- Table: reinsurance_region_loadings

CREATE TABLE IF NOT EXISTS reinsurance_region_loadings (
    id SERIAL PRIMARY KEY
);

ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS region VARCHAR(255);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS gender VARCHAR(255);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS gla_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS gla_aids_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS ptd_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS ci_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS ttd_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS phi_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS fun_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS fun_aids_region_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP;
ALTER TABLE reinsurance_region_loadings ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

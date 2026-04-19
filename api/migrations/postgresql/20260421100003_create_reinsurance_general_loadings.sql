-- Migration for struct: ReinsuranceGeneralLoading
-- Table: reinsurance_general_loadings

CREATE TABLE IF NOT EXISTS reinsurance_general_loadings (
    id SERIAL PRIMARY KEY
);

ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS age INTEGER;
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS gender VARCHAR(255);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS gla_contigency_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ptd_contigency_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ci_contigency_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ttd_contigency_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS phi_contigency_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS fun_contigency_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS continuation_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS terminal_illness_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ptd_accelerated_benefit_discount NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ci_accelerated_benefit_discount NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP;
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

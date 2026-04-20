-- Migration for struct: ReinsuranceGeneralLoading (voluntary loadings)
-- Table: reinsurance_general_loadings

ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS gla_voluntary_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ptd_voluntary_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ci_voluntary_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS ttd_voluntary_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS phi_voluntary_loading_rate NUMERIC(15,5);
ALTER TABLE reinsurance_general_loadings ADD COLUMN IF NOT EXISTS fun_voluntary_loading_rate NUMERIC(15,5);

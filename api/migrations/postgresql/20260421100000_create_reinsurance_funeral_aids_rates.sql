-- Migration for struct: ReinsuranceFuneralAidsRate
-- Table: reinsurance_funeral_aids_rates

CREATE TABLE IF NOT EXISTS reinsurance_funeral_aids_rates (
    id SERIAL PRIMARY KEY
);

ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN IF NOT EXISTS age_next_birthday INTEGER;
ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN IF NOT EXISTS gender VARCHAR(255);
ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN IF NOT EXISTS fun_aids_qx NUMERIC(15,5);
ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP;
ALTER TABLE reinsurance_funeral_aids_rates ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

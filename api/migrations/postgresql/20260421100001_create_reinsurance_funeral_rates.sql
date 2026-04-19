-- Migration for struct: ReinsuranceFuneralRate
-- Table: reinsurance_funeral_rates

CREATE TABLE IF NOT EXISTS reinsurance_funeral_rates (
    id SERIAL PRIMARY KEY
);

ALTER TABLE reinsurance_funeral_rates ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
ALTER TABLE reinsurance_funeral_rates ADD COLUMN IF NOT EXISTS age_next_birthday INTEGER;
ALTER TABLE reinsurance_funeral_rates ADD COLUMN IF NOT EXISTS gender VARCHAR(255);
ALTER TABLE reinsurance_funeral_rates ADD COLUMN IF NOT EXISTS fun_qx NUMERIC(15,5);
ALTER TABLE reinsurance_funeral_rates ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP;
ALTER TABLE reinsurance_funeral_rates ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

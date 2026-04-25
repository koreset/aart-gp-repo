-- Migration for struct: TaxRetirementTable
-- Table: tax_retirement_tables

CREATE TABLE IF NOT EXISTS tax_retirement_tables (
    id SERIAL PRIMARY KEY
);

ALTER TABLE tax_retirement_tables ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
ALTER TABLE tax_retirement_tables ADD COLUMN IF NOT EXISTS lower_bound NUMERIC(18,4);
ALTER TABLE tax_retirement_tables ADD COLUMN IF NOT EXISTS upper_bound NUMERIC(18,4);
ALTER TABLE tax_retirement_tables ADD COLUMN IF NOT EXISTS tax_rate NUMERIC(15,5);
ALTER TABLE tax_retirement_tables ADD COLUMN IF NOT EXISTS cumulative_tax_relief NUMERIC(18,4);
ALTER TABLE tax_retirement_tables ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP;
ALTER TABLE tax_retirement_tables ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

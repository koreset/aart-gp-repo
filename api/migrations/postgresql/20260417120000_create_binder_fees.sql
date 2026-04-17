-- Migration: create binder_fees reference table.
-- Uniquely keyed by (risk_rate_code, binderholder_name) so a binderholder
-- can carry different fee caps across product lines.

CREATE TABLE IF NOT EXISTS binder_fees (
    id                      SERIAL PRIMARY KEY,
    binderholder_name       VARCHAR(255),
    risk_rate_code          VARCHAR(255),
    maximum_binder_fee      DOUBLE PRECISION DEFAULT 0,
    maximum_outsource_fee   DOUBLE PRECISION DEFAULT 0,
    creation_date           TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by              VARCHAR(255)
);

-- Widen legacy columns if the table pre-existed with narrower types.
ALTER TABLE binder_fees ADD COLUMN IF NOT EXISTS binderholder_name VARCHAR(255);
ALTER TABLE binder_fees ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
ALTER TABLE binder_fees ADD COLUMN IF NOT EXISTS maximum_binder_fee DOUBLE PRECISION DEFAULT 0;
ALTER TABLE binder_fees ADD COLUMN IF NOT EXISTS maximum_outsource_fee DOUBLE PRECISION DEFAULT 0;
ALTER TABLE binder_fees ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE binder_fees ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

CREATE UNIQUE INDEX IF NOT EXISTS idx_binder_fee_rrc_holder
    ON binder_fees (risk_rate_code, binderholder_name);

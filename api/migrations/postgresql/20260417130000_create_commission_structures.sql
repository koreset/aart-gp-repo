-- Migration: create commission_structures reference table for per-channel
-- progressive sliding-scale commission. Each row is one band.
-- Uniquely keyed by (channel, lower_bound); upper_bound NULL = unbounded.
-- Direct channel is treated as flat 0% and is not stored.

CREATE TABLE IF NOT EXISTS commission_structures (
    id                   SERIAL PRIMARY KEY,
    channel              VARCHAR(20),
    lower_bound          DOUBLE PRECISION DEFAULT 0,
    upper_bound          DOUBLE PRECISION,
    maximum_commission   DOUBLE PRECISION DEFAULT 0,
    applicable_rate      DOUBLE PRECISION DEFAULT 0,
    creation_date        TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by           VARCHAR(255)
);

-- Widen legacy columns if the table pre-existed.
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS channel VARCHAR(20);
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS lower_bound DOUBLE PRECISION DEFAULT 0;
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS upper_bound DOUBLE PRECISION;
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS maximum_commission DOUBLE PRECISION DEFAULT 0;
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS applicable_rate DOUBLE PRECISION DEFAULT 0;
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);

CREATE UNIQUE INDEX IF NOT EXISTS idx_commission_channel_lower
    ON commission_structures (channel, lower_bound);

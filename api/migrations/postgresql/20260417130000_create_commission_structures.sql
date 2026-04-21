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

-- Dedupe duplicate (channel, lower_bound) rows before the unique index is
-- created. Keeps the earliest-inserted row (smallest id) for each pair.
-- No-op once the unique index exists. IS NOT DISTINCT FROM treats NULLs as
-- equal so NULL channels also dedupe.
DELETE FROM commission_structures c1
USING commission_structures c2
WHERE c1.id > c2.id
  AND c1.channel IS NOT DISTINCT FROM c2.channel
  AND c1.lower_bound IS NOT DISTINCT FROM c2.lower_bound;

CREATE UNIQUE INDEX IF NOT EXISTS idx_commission_channel_lower
    ON commission_structures (channel, lower_bound);

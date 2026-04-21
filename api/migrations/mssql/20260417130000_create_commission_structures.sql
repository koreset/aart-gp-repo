-- Migration: create commission_structures reference table for per-channel
-- progressive sliding-scale commission.
-- Uniquely keyed by (channel, lower_bound); upper_bound NULL = unbounded.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'commission_structures')
BEGIN
    CREATE TABLE commission_structures (
        id                   INT IDENTITY(1,1) PRIMARY KEY,
        channel              NVARCHAR(20),
        lower_bound          FLOAT DEFAULT 0,
        upper_bound          FLOAT,
        maximum_commission   FLOAT DEFAULT 0,
        applicable_rate      FLOAT DEFAULT 0,
        creation_date        DATETIME DEFAULT CURRENT_TIMESTAMP,
        created_by           NVARCHAR(255)
    );
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('commission_structures') AND name = 'channel')
BEGIN
    ALTER TABLE commission_structures ALTER COLUMN channel NVARCHAR(20);
END

-- Dedupe duplicate (channel, lower_bound) rows before the unique index is
-- created. Keeps the earliest-inserted row (smallest id) for each pair.
-- No-op once the unique index exists. NULL-safe so NULL channels also dedupe.
DELETE c1
FROM commission_structures c1
INNER JOIN commission_structures c2
    ON c1.id > c2.id
   AND (c1.channel = c2.channel OR (c1.channel IS NULL AND c2.channel IS NULL))
   AND (c1.lower_bound = c2.lower_bound OR (c1.lower_bound IS NULL AND c2.lower_bound IS NULL));

IF NOT EXISTS(SELECT 1 FROM sys.indexes WHERE object_id = OBJECT_ID('commission_structures') AND name = 'idx_commission_channel_lower')
BEGIN
    CREATE UNIQUE INDEX idx_commission_channel_lower
        ON commission_structures (channel, lower_bound);
END

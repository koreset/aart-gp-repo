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

IF NOT EXISTS(SELECT 1 FROM sys.indexes WHERE object_id = OBJECT_ID('commission_structures') AND name = 'idx_commission_channel_lower')
BEGIN
    CREATE UNIQUE INDEX idx_commission_channel_lower
        ON commission_structures (channel, lower_bound);
END

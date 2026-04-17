-- Migration: add holder_name to commission_structures and widen the
-- composite unique index to (channel, holder_name, lower_bound). Empty
-- holder_name ('') defines the channel default used as a fallback in pricing.

-- 1. Add the column (idempotent).
IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('commission_structures') AND name = 'holder_name')
BEGIN
    ALTER TABLE commission_structures ADD holder_name NVARCHAR(255) DEFAULT '';
END
ELSE
BEGIN
    ALTER TABLE commission_structures ALTER COLUMN holder_name NVARCHAR(255);
END

-- 2. Backfill NULLs so contents are consistent.
UPDATE commission_structures SET holder_name = '' WHERE holder_name IS NULL;

-- 3. Drop the old (channel, lower_bound) unique index if present.
IF EXISTS(SELECT 1 FROM sys.indexes WHERE object_id = OBJECT_ID('commission_structures') AND name = 'idx_commission_channel_lower')
BEGIN
    DROP INDEX idx_commission_channel_lower ON commission_structures;
END

-- 4. Create the new (channel, holder_name, lower_bound) unique index if missing.
IF NOT EXISTS(SELECT 1 FROM sys.indexes WHERE object_id = OBJECT_ID('commission_structures') AND name = 'idx_commission_channel_holder_lower')
BEGIN
    CREATE UNIQUE INDEX idx_commission_channel_holder_lower
        ON commission_structures (channel, holder_name, lower_bound);
END

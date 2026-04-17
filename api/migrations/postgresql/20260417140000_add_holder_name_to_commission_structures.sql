-- Migration: add holder_name to commission_structures so the sliding
-- scale can be scoped per broker / binder / tied-agent. Empty holder_name
-- ('') defines the channel default used as a fallback in pricing.
--
-- Widens the composite unique index to (channel, holder_name, lower_bound).

-- 1. Add the column (idempotent).
ALTER TABLE commission_structures ADD COLUMN IF NOT EXISTS holder_name VARCHAR(255) DEFAULT '';

-- Widen legacy narrower types, if any.
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='commission_structures' AND column_name='holder_name') THEN
        ALTER TABLE commission_structures ALTER COLUMN holder_name TYPE VARCHAR(255);
    END IF;
END $$;

-- 2. Backfill NULLs so contents are consistent.
UPDATE commission_structures SET holder_name = '' WHERE holder_name IS NULL;

-- 3. Drop the old (channel, lower_bound) unique index.
DROP INDEX IF EXISTS idx_commission_channel_lower;

-- 4. Create the new composite unique index (channel, holder_name, lower_bound).
CREATE UNIQUE INDEX IF NOT EXISTS idx_commission_channel_holder_lower
    ON commission_structures (channel, holder_name, lower_bound);

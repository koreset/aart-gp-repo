-- Migration: add holder_name to commission_structures and widen the
-- composite unique index to (channel, holder_name, lower_bound). Empty
-- holder_name ('') defines the channel default used as a fallback in pricing.

-- 1. Add the column if missing; ensure VARCHAR(255) so the composite
--    unique index works (MySQL rejects TEXT in a key without prefix length).
SET @sql := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'commission_structures'
              AND COLUMN_NAME = 'holder_name'
        ),
        'ALTER TABLE commission_structures MODIFY COLUMN holder_name VARCHAR(255) DEFAULT "";',
        'ALTER TABLE commission_structures ADD COLUMN holder_name VARCHAR(255) DEFAULT "";'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 2. Backfill NULLs so contents are consistent.
UPDATE commission_structures SET holder_name = '' WHERE holder_name IS NULL;

-- 3. Drop the old (channel, lower_bound) unique index if present.
SET @sql := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'commission_structures'
              AND INDEX_NAME = 'idx_commission_channel_lower'
        ),
        'DROP INDEX idx_commission_channel_lower ON commission_structures;',
        'SELECT 1;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 4. Create the new (channel, holder_name, lower_bound) unique index if missing.
SET @sql := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'commission_structures'
              AND INDEX_NAME = 'idx_commission_channel_holder_lower'
        ),
        'SELECT 1;',
        'CREATE UNIQUE INDEX idx_commission_channel_holder_lower ON commission_structures (channel, holder_name, lower_bound);'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

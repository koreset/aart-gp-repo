-- Migration: add `calculation_completed_at` to group_pricing_quotes so the
-- UI can display when the underlying calculation last produced results.
-- Nullable on purpose: pre-existing quotes have no such timestamp until
-- they are recalculated.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'group_pricing_quotes'
          AND column_name = 'calculation_completed_at'
    ) THEN
        ALTER TABLE group_pricing_quotes ADD COLUMN calculation_completed_at TIMESTAMP NULL;
    END IF;
END $$;

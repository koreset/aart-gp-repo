-- Migration: add `calculation_completed_at` to group_pricing_quotes so the
-- UI can display when the underlying calculation last produced results.
-- Nullable on purpose: pre-existing quotes have no such timestamp until
-- they are recalculated.

IF NOT EXISTS(
    SELECT * FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_NAME = 'group_pricing_quotes'
      AND COLUMN_NAME = 'calculation_completed_at'
)
BEGIN
    ALTER TABLE group_pricing_quotes ADD calculation_completed_at DATETIME NULL;
END;

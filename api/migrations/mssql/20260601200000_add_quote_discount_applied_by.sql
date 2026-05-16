-- Quote discount accountability: adds discount_applied_by /
-- discount_applied_at to group_pricing_quotes so the quote list can
-- show who set the current discount. Matches the existing pattern of
-- created_by / modified_by / approved_by on the same table.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'discount_applied_by')
  ALTER TABLE group_pricing_quotes ADD discount_applied_by NVARCHAR(255) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'discount_applied_at')
  ALTER TABLE group_pricing_quotes ADD discount_applied_at DATETIME NULL;

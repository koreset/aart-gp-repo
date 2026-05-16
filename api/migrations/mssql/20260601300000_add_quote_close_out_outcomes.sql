-- Quote close-out outcomes: NTU + Declined.
--
-- Adds per-status milestone timestamps + reason text to group_pricing_quotes
-- so an approved quote can be closed out with one of three terminal
-- states: accepted (existing), not_taken_up (broker / client went
-- elsewhere), or declined (deal fell through on the client side).
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'not_taken_up_at')
  ALTER TABLE group_pricing_quotes ADD not_taken_up_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'not_taken_up_reason')
  ALTER TABLE group_pricing_quotes ADD not_taken_up_reason NVARCHAR(500) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'declined_at')
  ALTER TABLE group_pricing_quotes ADD declined_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_pricing_quotes') AND name = 'declined_reason')
  ALTER TABLE group_pricing_quotes ADD declined_reason NVARCHAR(500) NULL;

-- Filtered indexes — same shape as the existing milestone indexes added
-- in 20260601000000_add_quote_performance_dashboard.sql. WHERE clause
-- keeps the index small (only rows that actually reached the milestone).
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpq_not_taken_up_at' AND object_id = OBJECT_ID('group_pricing_quotes'))
  CREATE INDEX idx_gpq_not_taken_up_at ON group_pricing_quotes(not_taken_up_at) WHERE not_taken_up_at IS NOT NULL;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_gpq_declined_at' AND object_id = OBJECT_ID('group_pricing_quotes'))
  CREATE INDEX idx_gpq_declined_at ON group_pricing_quotes(declined_at) WHERE declined_at IS NOT NULL;

-- Migration: add maximum allowed Free Cover Limit to the restrictions table.
-- Used as the underwriting ceiling for user-enforced quote-level FCL overrides.
-- A quote's FreeCoverLimit is honoured directly unless it exceeds this ceiling
-- by more than 20%, in which case it is clamped to the ceiling. A value of 0
-- means no ceiling is configured and user overrides pass through unclamped.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('restrictions') AND name = 'maximum_allowed_fcl')
BEGIN
    ALTER TABLE restrictions
        ADD maximum_allowed_fcl FLOAT NOT NULL CONSTRAINT df_restrictions_maximum_allowed_fcl DEFAULT 0;
END;

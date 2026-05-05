-- Migration: add per-benefit maximum cover ceilings to the restrictions table.
-- maximum_gla_cover and maximum_ptd_cover cap the member's covered sum assured
-- for the GLA and PTD benefits respectively. A value of 0 means no ceiling is
-- configured and the cap is skipped. Applied on top of the existing
-- FreeCoverLimit / AppliedFreeCoverLimit logic in group pricing.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('restrictions') AND name = 'maximum_gla_cover')
BEGIN
    ALTER TABLE restrictions
        ADD maximum_gla_cover FLOAT NOT NULL CONSTRAINT df_restrictions_maximum_gla_cover DEFAULT 0;
END;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('restrictions') AND name = 'maximum_ptd_cover')
BEGIN
    ALTER TABLE restrictions
        ADD maximum_ptd_cover FLOAT NOT NULL CONSTRAINT df_restrictions_maximum_ptd_cover DEFAULT 0;
END;

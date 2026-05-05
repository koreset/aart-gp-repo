-- Migration: add per-benefit maximum cover ceilings to the restrictions table.
-- maximum_gla_cover and maximum_ptd_cover cap the member's covered sum assured
-- for the GLA and PTD benefits respectively. A value of 0 means no ceiling is
-- configured and the cap is skipped. Applied on top of the existing
-- FreeCoverLimit / AppliedFreeCoverLimit logic in group pricing.

ALTER TABLE restrictions
    ADD COLUMN maximum_gla_cover DOUBLE NOT NULL DEFAULT 0;

ALTER TABLE restrictions
    ADD COLUMN maximum_ptd_cover DOUBLE NOT NULL DEFAULT 0;

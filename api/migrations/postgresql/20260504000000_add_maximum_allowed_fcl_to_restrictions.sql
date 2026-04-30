-- Migration: add maximum allowed Free Cover Limit to the restrictions table.
-- Used as the underwriting ceiling for user-enforced quote-level FCL overrides.
-- A quote's FreeCoverLimit is honoured directly unless it exceeds this ceiling
-- by more than 20%, in which case it is clamped to the ceiling. A value of 0
-- means no ceiling is configured and user overrides pass through unclamped.

ALTER TABLE restrictions
    ADD COLUMN IF NOT EXISTS maximum_allowed_fcl DOUBLE PRECISION NOT NULL DEFAULT 0;

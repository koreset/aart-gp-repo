-- Migration: persist the per-quote cap limits on member_rating_result_summaries
-- so bordereaux's CoveredSumsAssured (and downstream premium recomputation)
-- can reapply the same caps the pricing flow used.
-- A column value of 0 means "no limit" for that benefit.

ALTER TABLE member_rating_result_summaries
    ADD COLUMN IF NOT EXISTS maximum_gla_cover             DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS maximum_ptd_cover             DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS severe_illness_maximum_benefit DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS spouse_gla_maximum_benefit    DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ttd_maximum_monthly_benefit   DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS phi_maximum_monthly_benefit   DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_max_gla_cover           DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_max_ptd_cover           DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_max_ci_cover            DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_max_sgla_cover          DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_max_ttd_cover           DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_max_phi_cover           DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reins_max_fun_cover           DOUBLE PRECISION NOT NULL DEFAULT 0;

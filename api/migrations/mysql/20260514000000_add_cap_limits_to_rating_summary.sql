-- Migration: persist the per-quote cap limits on member_rating_result_summaries
-- so bordereaux's CoveredSumsAssured (and downstream premium recomputation)
-- can reapply the same caps the pricing flow used.
-- A column value of 0 means "no limit" for that benefit.

ALTER TABLE member_rating_result_summaries
    ADD COLUMN maximum_gla_cover              DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN maximum_ptd_cover              DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN severe_illness_maximum_benefit DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN spouse_gla_maximum_benefit     DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ttd_maximum_monthly_benefit    DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN phi_maximum_monthly_benefit    DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN reins_max_gla_cover            DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN reins_max_ptd_cover            DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN reins_max_ci_cover             DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN reins_max_sgla_cover           DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN reins_max_ttd_cover            DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN reins_max_phi_cover            DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN reins_max_fun_cover            DOUBLE NOT NULL DEFAULT 0;

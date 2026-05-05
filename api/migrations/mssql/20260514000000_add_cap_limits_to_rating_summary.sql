-- Migration: persist the per-quote cap limits on member_rating_result_summaries
-- so bordereaux's CoveredSumsAssured (and downstream premium recomputation)
-- can reapply the same caps the pricing flow used.
-- A column value of 0 means "no limit" for that benefit.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'maximum_gla_cover')
BEGIN
    ALTER TABLE member_rating_result_summaries
        ADD maximum_gla_cover              FLOAT NOT NULL CONSTRAINT df_mrrs_maximum_gla_cover              DEFAULT 0,
            maximum_ptd_cover              FLOAT NOT NULL CONSTRAINT df_mrrs_maximum_ptd_cover              DEFAULT 0,
            severe_illness_maximum_benefit FLOAT NOT NULL CONSTRAINT df_mrrs_severe_illness_maximum_benefit DEFAULT 0,
            spouse_gla_maximum_benefit     FLOAT NOT NULL CONSTRAINT df_mrrs_spouse_gla_maximum_benefit     DEFAULT 0,
            ttd_maximum_monthly_benefit    FLOAT NOT NULL CONSTRAINT df_mrrs_ttd_maximum_monthly_benefit    DEFAULT 0,
            phi_maximum_monthly_benefit    FLOAT NOT NULL CONSTRAINT df_mrrs_phi_maximum_monthly_benefit    DEFAULT 0,
            reins_max_gla_cover            FLOAT NOT NULL CONSTRAINT df_mrrs_reins_max_gla_cover            DEFAULT 0,
            reins_max_ptd_cover            FLOAT NOT NULL CONSTRAINT df_mrrs_reins_max_ptd_cover            DEFAULT 0,
            reins_max_ci_cover             FLOAT NOT NULL CONSTRAINT df_mrrs_reins_max_ci_cover             DEFAULT 0,
            reins_max_sgla_cover           FLOAT NOT NULL CONSTRAINT df_mrrs_reins_max_sgla_cover           DEFAULT 0,
            reins_max_ttd_cover            FLOAT NOT NULL CONSTRAINT df_mrrs_reins_max_ttd_cover            DEFAULT 0,
            reins_max_phi_cover            FLOAT NOT NULL CONSTRAINT df_mrrs_reins_max_phi_cover            DEFAULT 0,
            reins_max_fun_cover            FLOAT NOT NULL CONSTRAINT df_mrrs_reins_max_fun_cover            DEFAULT 0;
END;

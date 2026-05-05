-- Migration: persist the per-benefit max-cover-age limits on
-- member_rating_result_summaries so the claim flow can reapply the same
-- age guards the pricing flow used (frozen at quote-rating time).
-- A column value of 0 means "no limit" for that benefit.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'gla_max_cover_age')
BEGIN
    ALTER TABLE member_rating_result_summaries
        ADD gla_max_cover_age INT NOT NULL CONSTRAINT df_mrrs_gla_max_cover_age DEFAULT 0,
            ptd_max_cover_age INT NOT NULL CONSTRAINT df_mrrs_ptd_max_cover_age DEFAULT 0,
            ci_max_cover_age  INT NOT NULL CONSTRAINT df_mrrs_ci_max_cover_age  DEFAULT 0,
            ttd_max_cover_age INT NOT NULL CONSTRAINT df_mrrs_ttd_max_cover_age DEFAULT 0,
            phi_max_cover_age INT NOT NULL CONSTRAINT df_mrrs_phi_max_cover_age DEFAULT 0,
            fun_max_cover_age INT NOT NULL CONSTRAINT df_mrrs_fun_max_cover_age DEFAULT 0;
END;

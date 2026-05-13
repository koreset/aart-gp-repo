-- Add per-member scheme-size and educator base loading columns to
-- member_rating_results. Mirrors the source-side columns added in
-- 20260518000000_add_scheme_size_loadings.sql (scheme_size_levels) and
-- 20260518010000_add_educator_loadings.sql (general_loadings).

ALTER TABLE member_rating_results
    ADD COLUMN gla_scheme_size_loading DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ptd_scheme_size_loading DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ci_scheme_size_loading  DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ttd_scheme_size_loading DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN phi_scheme_size_loading DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN fun_scheme_size_loading DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN gla_educator_loading    DOUBLE NOT NULL DEFAULT 0,
    ADD COLUMN ptd_educator_loading    DOUBLE NOT NULL DEFAULT 0;

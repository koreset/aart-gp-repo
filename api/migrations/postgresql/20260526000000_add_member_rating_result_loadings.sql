-- Add per-member scheme-size and educator base loading columns to
-- member_rating_results. Mirrors the source-side columns added in
-- 20260518000000_add_scheme_size_loadings.sql and
-- 20260518010000_add_educator_loadings.sql.

ALTER TABLE member_rating_results
    ADD COLUMN IF NOT EXISTS gla_scheme_size_loading DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ptd_scheme_size_loading DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ci_scheme_size_loading  DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ttd_scheme_size_loading DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS phi_scheme_size_loading DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS fun_scheme_size_loading DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS gla_educator_loading    DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ptd_educator_loading    DOUBLE PRECISION NOT NULL DEFAULT 0;

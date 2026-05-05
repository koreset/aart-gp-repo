-- Migration: persist the per-benefit max-cover-age limits on
-- member_rating_result_summaries so the claim flow can reapply the same
-- age guards the pricing flow used (frozen at quote-rating time).
-- A column value of 0 means "no limit" for that benefit.

ALTER TABLE member_rating_result_summaries
    ADD COLUMN IF NOT EXISTS gla_max_cover_age INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ptd_max_cover_age INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ci_max_cover_age  INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ttd_max_cover_age INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS phi_max_cover_age INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS fun_max_cover_age INTEGER NOT NULL DEFAULT 0;

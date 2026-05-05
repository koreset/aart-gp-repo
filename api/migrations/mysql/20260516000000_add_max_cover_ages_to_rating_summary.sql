-- Migration: persist the per-benefit max-cover-age limits on
-- member_rating_result_summaries so the claim flow can reapply the same
-- age guards the pricing flow used (frozen at quote-rating time).
-- A column value of 0 means "no limit" for that benefit.

ALTER TABLE member_rating_result_summaries
    ADD COLUMN gla_max_cover_age INT NOT NULL DEFAULT 0,
    ADD COLUMN ptd_max_cover_age INT NOT NULL DEFAULT 0,
    ADD COLUMN ci_max_cover_age  INT NOT NULL DEFAULT 0,
    ADD COLUMN ttd_max_cover_age INT NOT NULL DEFAULT 0,
    ADD COLUMN phi_max_cover_age INT NOT NULL DEFAULT 0,
    ADD COLUMN fun_max_cover_age INT NOT NULL DEFAULT 0;

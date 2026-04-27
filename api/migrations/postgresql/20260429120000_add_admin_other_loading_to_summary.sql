-- Migration: add admin_loading + other_loading columns to
-- member_rating_result_summaries. These let SchemeTotalLoading() include the
-- full premium-loading sum (expense + profit + admin + other + binder +
-- outsource) at the summary level, matching the rating-phase
-- TotalPremiumLoading on MemberRatingResult.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS admin_loading NUMERIC(20,6) DEFAULT 0;
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS other_loading NUMERIC(20,6) DEFAULT 0;

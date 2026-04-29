-- Migration: add indicative_rates_count column to member_rating_result_summaries.
-- Captured from the rating loop's local indicativeRatesCount so doc-template
-- helpers can derive office-side proportion-of-salary directly from a
-- persisted Final office premium and the annual salary, instead of grossing
-- up the risk proportion at render time.

ALTER TABLE member_rating_result_summaries ADD COLUMN indicative_rates_count DOUBLE NOT NULL DEFAULT 0;

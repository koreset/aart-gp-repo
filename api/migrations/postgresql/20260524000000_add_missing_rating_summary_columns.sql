-- Backfills four columns on member_rating_result_summaries that exist on
-- the MemberRatingResultSummary Go struct but were never delivered via a
-- migration file. Fresh-install DBs picked them up via the bootstrap
-- AutoMigrate path; existing DBs past bootstrap did not, which broke the
-- slim /result-summary endpoint with "Unknown column" errors after it
-- started fetching them by name. Idempotent on re-runs.

ALTER TABLE member_rating_result_summaries
  ADD COLUMN IF NOT EXISTS additional_gla_cover_age_band_type VARCHAR(64) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS additional_gla_cover_male_prop_used DOUBLE PRECISION NULL,
  ADD COLUMN IF NOT EXISTS extended_family_age_band_type VARCHAR(64) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS extended_family_pricing_method VARCHAR(64) NOT NULL DEFAULT '';

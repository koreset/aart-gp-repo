-- Backfills four columns on member_rating_result_summaries that exist on
-- the MemberRatingResultSummary Go struct but were never delivered via a
-- migration file. Fresh-install DBs picked them up via the bootstrap
-- AutoMigrate path; existing DBs past bootstrap did not, which broke the
-- slim /result-summary endpoint with "Unknown column" errors after it
-- started fetching them by name. Idempotent on re-runs because each ALTER
-- is gated on a sys.columns existence check.

IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('member_rating_result_summaries')
                 AND name = 'additional_gla_cover_age_band_type')
  ALTER TABLE member_rating_result_summaries
    ADD additional_gla_cover_age_band_type NVARCHAR(64) NOT NULL CONSTRAINT df_mrrs_aagla_age_band_type DEFAULT '';

IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('member_rating_result_summaries')
                 AND name = 'additional_gla_cover_male_prop_used')
  ALTER TABLE member_rating_result_summaries
    ADD additional_gla_cover_male_prop_used FLOAT NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('member_rating_result_summaries')
                 AND name = 'extended_family_age_band_type')
  ALTER TABLE member_rating_result_summaries
    ADD extended_family_age_band_type NVARCHAR(64) NOT NULL CONSTRAINT df_mrrs_extfam_age_band_type DEFAULT '';

IF NOT EXISTS (SELECT 1 FROM sys.columns
               WHERE object_id = OBJECT_ID('member_rating_result_summaries')
                 AND name = 'extended_family_pricing_method')
  ALTER TABLE member_rating_result_summaries
    ADD extended_family_pricing_method NVARCHAR(64) NOT NULL CONSTRAINT df_mrrs_extfam_pricing_method DEFAULT '';

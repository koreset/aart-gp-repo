-- Migration: add admin_loading + other_loading columns to
-- member_rating_result_summaries. These let SchemeTotalLoading() include the
-- full premium-loading sum (expense + profit + admin + other + binder +
-- outsource) at the summary level, matching the rating-phase
-- TotalPremiumLoading on MemberRatingResult.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='admin_loading')
    ALTER TABLE member_rating_result_summaries ADD admin_loading FLOAT DEFAULT 0;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='other_loading')
    ALTER TABLE member_rating_result_summaries ADD other_loading FLOAT DEFAULT 0;

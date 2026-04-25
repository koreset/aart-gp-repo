-- Migration: persist binder-fee and outsource-fee rates on each member rating
-- result summary row alongside the existing scheme-level loadings.

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'binder_fee_rate')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD binder_fee_rate DECIMAL(15,5);
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'outsource_fee_rate')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD outsource_fee_rate DECIMAL(15,5);
END;

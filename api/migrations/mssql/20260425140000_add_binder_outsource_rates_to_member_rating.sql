-- Migration: persist binder-fee and outsource-fee rates on each member rating
-- result row so the per-member loading breakdown is auditable end-to-end.

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'binder_fee_rate')
BEGIN
    ALTER TABLE member_rating_results ADD binder_fee_rate DECIMAL(15,5);
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_results' AND COLUMN_NAME = 'outsource_fee_rate')
BEGIN
    ALTER TABLE member_rating_results ADD outsource_fee_rate DECIMAL(15,5);
END;

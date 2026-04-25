-- Migration: persist binder-fee and outsource-fee rates on each member rating
-- result summary row alongside the existing scheme-level loadings.

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS binder_fee_rate NUMERIC(15,5);
ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS outsource_fee_rate NUMERIC(15,5);

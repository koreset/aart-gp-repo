-- Migration: persist binder-fee and outsource-fee rates on each member rating
-- result row so the per-member loading breakdown is auditable end-to-end.

ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS binder_fee_rate NUMERIC(15,5);
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS outsource_fee_rate NUMERIC(15,5);

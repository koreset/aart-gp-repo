-- Migration: persist the actual male proportion used for each
-- Additional GLA Cover calc run on scheme_categories. Audit column so the
-- value that fed the rate is stored on the same row that was
-- (re)calculated, independent of any later re-derivation from members.

ALTER TABLE scheme_categories ADD COLUMN IF NOT EXISTS additional_gla_cover_male_prop_used DOUBLE PRECISION NULL;

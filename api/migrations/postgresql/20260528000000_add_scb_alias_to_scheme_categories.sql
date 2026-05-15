-- Add scb_alias column to scheme_categories. Mirrors the other per-benefit
-- alias columns (phi_alias, ttd_alias, etc.) added in
-- 20251204162201_update_schemecategory.sql; this row brings SCB into the
-- benefit-customisation system so the user can rename it from the UI.
-- Idempotent on re-runs.

ALTER TABLE scheme_categories
    ADD COLUMN IF NOT EXISTS scb_alias VARCHAR(255);

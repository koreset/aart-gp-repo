-- Add scb_alias column to scheme_categories. Mirrors the other per-benefit
-- alias columns (phi_alias, ttd_alias, etc.) added in
-- 20251204162201_update_schemecategory.sql; this row brings SCB into the
-- benefit-customisation system so the user can rename it from the UI.
-- Idempotent on re-runs.

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'scheme_categories' AND COLUMN_NAME = 'scb_alias')
BEGIN
    ALTER TABLE scheme_categories ADD scb_alias NVARCHAR(255);
END;

-- Add scb_alias column to scheme_categories. Mirrors the other per-benefit
-- alias columns (phi_alias, ttd_alias, etc.) added in
-- 20251204162201_update_schemecategory.sql; this row brings SCB into the
-- benefit-customisation system so the user can rename it from the UI.
-- Idempotent on re-runs.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'scheme_categories' AND column_name = 'scb_alias');
SET @sql := IF(@col = 0,
  'ALTER TABLE scheme_categories ADD COLUMN scb_alias VARCHAR(255)',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

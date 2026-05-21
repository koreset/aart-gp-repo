-- Migration: add created_by column to cpi_indices so the UI can show who
-- registered or uploaded each CPI value.

ALTER TABLE cpi_indices
    ADD COLUMN IF NOT EXISTS created_by VARCHAR(128) NOT NULL DEFAULT '';

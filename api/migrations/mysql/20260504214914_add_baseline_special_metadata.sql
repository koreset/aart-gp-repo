-- Generated 2026-05-04T21:49:14+02:00 for dialect mysql

-- Migration for: GPPermission (table: gp_permissions)

ALTER TABLE gp_permissions ADD COLUMN category VARCHAR(255);
ALTER TABLE gp_permissions ADD COLUMN tier VARCHAR(255);
ALTER TABLE gp_permissions ADD COLUMN parent_slug VARCHAR(255);
ALTER TABLE gp_permissions ADD COLUMN display_order BIGINT;


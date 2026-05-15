-- Migration: add age calculation methodology selection to the singleton
-- group_pricing_settings row, with audit pair for change tracking. Default
-- 'age_next_birthday' preserves historical behaviour (member age is rounded up
-- once the commencement date crosses the birthday in the same year). Switching
-- to 'age_last_birthday' uses the floored age:
--   ROUNDDOWN((12*(YEAR(CommenDate)-YEAR(DoB)) + (MONTH(CommenDate)-MONTH(DoB)))/12, 0)

ALTER TABLE group_pricing_settings
    ADD COLUMN age_method VARCHAR(32) NOT NULL DEFAULT 'age_next_birthday',
    ADD COLUMN age_method_updated_at DATETIME NULL,
    ADD COLUMN age_method_updated_by VARCHAR(255) NOT NULL DEFAULT '';

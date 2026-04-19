-- Migration: add main_member_male_prop to group_pricing_parameters.
-- Used to weight male vs female gla qx when computing additional-GLA-cover
-- age-band rates:
--   qx = male_prop * male_qx + (1 - male_prop) * female_qx
-- Default 0.5 preserves the straight-average behaviour for any pre-existing
-- rows.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='group_pricing_parameters' AND COLUMN_NAME='main_member_male_prop'),
        'ALTER TABLE group_pricing_parameters MODIFY COLUMN main_member_male_prop DOUBLE NOT NULL DEFAULT 0.5;',
        'ALTER TABLE group_pricing_parameters ADD COLUMN main_member_male_prop DOUBLE NOT NULL DEFAULT 0.5;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Migration: persist the actual male proportion used for each
-- Additional GLA Cover calc run on scheme_categories. Audit column so the
-- value that fed the rate is stored on the same row that was
-- (re)calculated, independent of any later re-derivation from members.

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='scheme_categories' AND COLUMN_NAME='additional_gla_cover_male_prop_used'),
        'ALTER TABLE scheme_categories MODIFY COLUMN additional_gla_cover_male_prop_used DOUBLE NULL;',
        'ALTER TABLE scheme_categories ADD COLUMN additional_gla_cover_male_prop_used DOUBLE NULL;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

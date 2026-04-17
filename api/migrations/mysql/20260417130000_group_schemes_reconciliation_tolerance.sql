-- Migration: add reconciliation_tolerance to group_schemes
-- Purpose: P2-3 — allow per-scheme configuration of the variance threshold used
-- when classifying bordereaux reconciliation lines as matched vs discrepancy.
-- See services/bordereaux.go reconciliationToleranceForScheme.

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='reconciliation_tolerance' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE group_schemes ADD COLUMN reconciliation_tolerance DOUBLE NOT NULL DEFAULT 0;'
));
PREPARE stmt FROM @s; EXECUTE stmt; DEALLOCATE PREPARE stmt;

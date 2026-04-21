-- Migration: add reconciliation_tolerance to group_schemes
-- Purpose: P2-3 — allow per-scheme configuration of the variance threshold used
-- when classifying bordereaux reconciliation lines as matched vs discrepancy.
-- See services/bordereaux.go reconciliationToleranceForScheme.

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='group_schemes' AND COLUMN_NAME='reconciliation_tolerance')
BEGIN
    ALTER TABLE group_schemes ADD reconciliation_tolerance FLOAT NOT NULL DEFAULT 0;
END;

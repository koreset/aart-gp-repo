-- Migration: add reconciliation_tolerance to group_schemes
-- Purpose: P2-3 — allow per-scheme configuration of the variance threshold used
-- when classifying bordereaux reconciliation lines as matched vs discrepancy.
-- A value of 0 (default on existing rows) means "use the codebase default
-- (0.001)" so this migration is safe to run before any operator tunes a scheme.
-- See services/bordereaux.go reconciliationToleranceForScheme.

ALTER TABLE group_schemes ADD COLUMN IF NOT EXISTS reconciliation_tolerance DOUBLE PRECISION NOT NULL DEFAULT 0;

-- Migration: create bordereaux_confirmation_notes
-- Purpose: P2-4 — dedicated table for free-text notes on a bordereaux
-- confirmation. Replaces the prior "synthetic _note row in
-- bordereaux_reconciliation_results" hack so reconciliation-results queries
-- no longer need to filter out _note entries.

CREATE TABLE IF NOT EXISTS bordereaux_confirmation_notes (
    id                          SERIAL PRIMARY KEY,
    bordereaux_confirmation_id  INTEGER NOT NULL,
    note                        TEXT NOT NULL DEFAULT '',
    created_by                  VARCHAR(255) NOT NULL DEFAULT '',
    created_at                  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bcn_confirmation_id ON bordereaux_confirmation_notes (bordereaux_confirmation_id);

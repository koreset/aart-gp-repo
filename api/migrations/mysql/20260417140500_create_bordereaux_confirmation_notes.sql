-- Migration: create bordereaux_confirmation_notes
-- Purpose: P2-4 — dedicated table for free-text notes on a bordereaux
-- confirmation. Replaces the "synthetic _note row" hack.

CREATE TABLE IF NOT EXISTS bordereaux_confirmation_notes (
    id                          INT AUTO_INCREMENT PRIMARY KEY,
    bordereaux_confirmation_id  INT NOT NULL,
    note                        TEXT NOT NULL,
    created_by                  VARCHAR(255) NOT NULL DEFAULT '',
    created_at                  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_bcn_confirmation_id (bordereaux_confirmation_id)
);

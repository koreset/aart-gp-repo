-- Migration: create bordereaux_confirmation_notes
-- Purpose: P2-4 — dedicated table for free-text notes on a bordereaux
-- confirmation. Replaces the "synthetic _note row" hack.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'bordereaux_confirmation_notes')
BEGIN
    CREATE TABLE bordereaux_confirmation_notes (
        id                          INT IDENTITY(1,1) PRIMARY KEY,
        bordereaux_confirmation_id  INT NOT NULL,
        note                        NVARCHAR(MAX) NOT NULL,
        created_by                  NVARCHAR(255) NOT NULL DEFAULT '',
        created_at                  DATETIME2 NOT NULL DEFAULT SYSUTCDATETIME()
    );
END;

IF NOT EXISTS(SELECT * FROM sys.indexes WHERE name='idx_bcn_confirmation_id' AND object_id=OBJECT_ID('bordereaux_confirmation_notes'))
BEGIN
    CREATE INDEX idx_bcn_confirmation_id ON bordereaux_confirmation_notes (bordereaux_confirmation_id);
END;

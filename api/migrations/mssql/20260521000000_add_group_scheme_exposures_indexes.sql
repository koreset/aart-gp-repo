-- Adds two composite indexes on group_scheme_exposures so the GP dashboard's
-- region and industry-by-age queries hit index range scans instead of full
-- table scans. quote_status is widened to a sized NVARCHAR first because
-- MSSQL (like MySQL) cannot index NVARCHAR(MAX) columns. Idempotent on
-- re-runs because both the column type and the index existence are checked
-- first.

-- 1) Convert quote_status to NVARCHAR(64) if it's currently NVARCHAR(MAX) /
--    VARCHAR(MAX) / TEXT / NTEXT.
IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('group_scheme_exposures')
             AND name = 'quote_status'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE group_scheme_exposures
    ALTER COLUMN quote_status NVARCHAR(64) NOT NULL;

-- 2) idx_gse_year_quote on (financial_year, quote_id)
IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_gse_year_quote'
                 AND object_id = OBJECT_ID('group_scheme_exposures'))
  CREATE INDEX idx_gse_year_quote
    ON group_scheme_exposures (financial_year, quote_id);

-- 3) idx_gse_year_status on (financial_year, quote_status)
IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_gse_year_status'
                 AND object_id = OBJECT_ID('group_scheme_exposures'))
  CREATE INDEX idx_gse_year_status
    ON group_scheme_exposures (financial_year, quote_status);

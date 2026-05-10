-- Adds composite indexes on member_rating_result_summaries and
-- member_rating_results so the GP dashboard's status/year rollups and the
-- per-quote demographic aggregations hit index range scans instead of full
-- table / full-quote scans.
--
-- String columns participating in new indexes are first converted from
-- NVARCHAR(MAX) / VARCHAR(MAX) / TEXT / NTEXT to sized NVARCHAR, since MSSQL
-- cannot index MAX-length string columns. Idempotent on re-runs because
-- both the column type and the index existence are checked first.

-- ─────────────────────────────────────────────────────────────────────────
-- 1) Convert string columns on member_rating_result_summaries to sized NVARCHAR.
-- ─────────────────────────────────────────────────────────────────────────

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('member_rating_result_summaries')
             AND name = 'if_status'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE member_rating_result_summaries
    ALTER COLUMN if_status NVARCHAR(64) NOT NULL;

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('member_rating_result_summaries')
             AND name = 'quote_type'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE member_rating_result_summaries
    ALTER COLUMN quote_type NVARCHAR(64) NOT NULL;

-- ─────────────────────────────────────────────────────────────────────────
-- 2) Convert string columns on member_rating_results to sized NVARCHAR.
-- ─────────────────────────────────────────────────────────────────────────

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('member_rating_results')
             AND name = 'category'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE member_rating_results
    ALTER COLUMN category NVARCHAR(128) NOT NULL;

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('member_rating_results')
             AND name = 'age_band'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE member_rating_results
    ALTER COLUMN age_band NVARCHAR(64) NOT NULL;

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('member_rating_results')
             AND name = 'gender'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE member_rating_results
    ALTER COLUMN gender NVARCHAR(32) NOT NULL;

-- ─────────────────────────────────────────────────────────────────────────
-- 3) Create indexes on member_rating_result_summaries.
-- ─────────────────────────────────────────────────────────────────────────

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_mrrs_status_type_creation'
                 AND object_id = OBJECT_ID('member_rating_result_summaries'))
  CREATE INDEX idx_mrrs_status_type_creation
    ON member_rating_result_summaries (if_status, quote_type, creation_date);

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_mrrs_creation_date'
                 AND object_id = OBJECT_ID('member_rating_result_summaries'))
  CREATE INDEX idx_mrrs_creation_date
    ON member_rating_result_summaries (creation_date);

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_mrrs_quote_creation'
                 AND object_id = OBJECT_ID('member_rating_result_summaries'))
  CREATE INDEX idx_mrrs_quote_creation
    ON member_rating_result_summaries (quote_id, creation_date);

-- ─────────────────────────────────────────────────────────────────────────
-- 4) Create indexes on member_rating_results.
-- ─────────────────────────────────────────────────────────────────────────

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_mrr_quote_category'
                 AND object_id = OBJECT_ID('member_rating_results'))
  CREATE INDEX idx_mrr_quote_category
    ON member_rating_results (quote_id, category);

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_mrr_quote_age_gender'
                 AND object_id = OBJECT_ID('member_rating_results'))
  CREATE INDEX idx_mrr_quote_age_gender
    ON member_rating_results (quote_id, age_band, gender);

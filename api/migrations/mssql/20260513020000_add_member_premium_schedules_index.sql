-- Add a composite index on member_premium_schedules so that
-- refreshMemberPremiumSchedules' per-quote ORDER BY (category, member_name)
-- and the per-quote DELETE both hit an index seek + ordered scan instead
-- of a full table scan + sort. MAX-length string columns are first sized
-- down to indexable NVARCHAR lengths. Idempotent on re-runs.

-- 1) Widen / size string columns participating in the new index.

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('member_premium_schedules')
             AND name = 'category'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE member_premium_schedules
    ALTER COLUMN category NVARCHAR(128) NOT NULL;

IF EXISTS (SELECT 1 FROM sys.columns
           WHERE object_id = OBJECT_ID('member_premium_schedules')
             AND name = 'member_name'
             AND (max_length = -1 OR system_type_id IN (35, 99)))
  ALTER TABLE member_premium_schedules
    ALTER COLUMN member_name NVARCHAR(255) NOT NULL;

-- 2) Create composite index.

IF NOT EXISTS (SELECT 1 FROM sys.indexes
               WHERE name = 'idx_mps_quote_category_member'
                 AND object_id = OBJECT_ID('member_premium_schedules'))
    CREATE INDEX idx_mps_quote_category_member
        ON member_premium_schedules (quote_id, category, member_name);

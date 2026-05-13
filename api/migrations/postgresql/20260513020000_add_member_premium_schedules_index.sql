-- Add a composite index on member_premium_schedules so that
-- refreshMemberPremiumSchedules' per-quote ORDER BY (category, member_name)
-- and the per-quote DELETE both hit an index seek + ordered scan instead
-- of a full table scan + sort. Idempotent on re-runs.

CREATE INDEX IF NOT EXISTS idx_mps_quote_category_member
    ON member_premium_schedules (quote_id, category, member_name);

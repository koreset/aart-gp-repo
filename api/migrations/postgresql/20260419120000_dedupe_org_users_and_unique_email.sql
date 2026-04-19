-- Dedupe org_users by email then enforce uniqueness on email.
-- Required because the refresh-from-license-server path relies on an
-- ON CONFLICT (email) upsert. Without this index the upsert silently
-- degrades to a plain INSERT and creates duplicates.
--
-- Dedup strategy: keep the row with the most useful locally-managed
-- data per email (gp_role assigned > gp_role_id assigned > role
-- assigned), tie-broken by lowest id.

-- 1. Remove rows that cannot participate in a unique-email index.
DELETE FROM org_users WHERE email IS NULL OR email = '';

-- 2. Delete duplicate rows, keeping the best per email.
WITH ranked AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY email
               ORDER BY
                   CASE WHEN gp_role IS NOT NULL AND gp_role <> '' AND gp_role <> 'None' THEN 0 ELSE 1 END,
                   CASE WHEN gp_role_id IS NOT NULL AND gp_role_id <> 0 THEN 0 ELSE 1 END,
                   CASE WHEN role IS NOT NULL AND role <> '' THEN 0 ELSE 1 END,
                   id ASC
           ) AS rn
    FROM org_users
)
DELETE FROM org_users
WHERE id IN (SELECT id FROM ranked WHERE rn > 1);

-- 3. Add the unique index. IF NOT EXISTS so re-runs are safe.
CREATE UNIQUE INDEX IF NOT EXISTS idx_org_users_email ON org_users (email);

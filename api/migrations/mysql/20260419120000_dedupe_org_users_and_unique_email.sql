-- Dedupe org_users by email then enforce uniqueness on email.
-- Required because the refresh-from-license-server path relies on an
-- ON DUPLICATE KEY upsert keyed by email. Without this index the upsert
-- silently degrades to a plain INSERT and creates duplicates.
--
-- Dedup strategy: keep the row with the most useful locally-managed
-- data per email — a row with a meaningful gp_role beats one without,
-- tie-broken by lowest id. Written as a self-join so it works on
-- MySQL 5.7 as well as 8.0+.

-- 1. Remove rows that cannot participate in a unique-email index.
DELETE FROM org_users WHERE email IS NULL OR email = '';

-- 2. For each email, delete rows that are dominated by a "better" sibling.
DELETE r1 FROM org_users r1
INNER JOIN org_users r2
    ON r1.email = r2.email
   AND r1.id <> r2.id
   AND (
        (
            r2.gp_role IS NOT NULL AND r2.gp_role <> '' AND r2.gp_role <> 'None'
            AND NOT (r1.gp_role IS NOT NULL AND r1.gp_role <> '' AND r1.gp_role <> 'None')
        )
        OR
        (
            (
                (r2.gp_role IS NOT NULL AND r2.gp_role <> '' AND r2.gp_role <> 'None')
                =
                (r1.gp_role IS NOT NULL AND r1.gp_role <> '' AND r1.gp_role <> 'None')
            )
            AND r2.id < r1.id
        )
   );

-- 3. Add the unique index. MySQL has no `CREATE UNIQUE INDEX IF NOT
--    EXISTS`, so build the right ALTER TABLE dynamically based on
--    INFORMATION_SCHEMA. The prepared SQL contains no embedded
--    semicolons (the migration runner splits on lines ending in ";").
SET @create_sql = (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'org_users'
              AND INDEX_NAME = 'idx_org_users_email'
        ),
        'SELECT 1',
        'ALTER TABLE org_users ADD UNIQUE INDEX idx_org_users_email (email)'
    )
);
PREPARE stmt FROM @create_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

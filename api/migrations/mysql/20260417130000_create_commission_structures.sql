-- Migration: create commission_structures reference table for per-channel
-- progressive sliding-scale commission.
-- Uniquely keyed by (channel, lower_bound); upper_bound NULL = unbounded.

CREATE TABLE IF NOT EXISTS commission_structures (
    id                   INT AUTO_INCREMENT PRIMARY KEY,
    channel              VARCHAR(20),
    lower_bound          DOUBLE DEFAULT 0,
    upper_bound          DOUBLE,
    maximum_commission   DOUBLE DEFAULT 0,
    applicable_rate      DOUBLE DEFAULT 0,
    creation_date        DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by           VARCHAR(255)
);

-- If the table pre-existed (e.g. via GORM AutoMigrate) ensure the key
-- columns are VARCHAR(20) rather than TEXT/LONGTEXT so the composite
-- unique index doesn't hit error 1170.
ALTER TABLE commission_structures MODIFY COLUMN channel VARCHAR(20);

-- Dedupe duplicate (channel, lower_bound) rows before the unique index is
-- created. Keeps the earliest-inserted row (smallest id) for each pair.
-- No-op once the unique index exists, so safe to run repeatedly.
-- `<=>` is MySQL's null-safe equal, so NULL channels also dedupe.
DELETE c1 FROM commission_structures c1
JOIN commission_structures c2
    ON c1.channel <=> c2.channel
   AND c1.lower_bound <=> c2.lower_bound
   AND c1.id > c2.id;

SET @sql := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'commission_structures'
              AND INDEX_NAME = 'idx_commission_channel_lower'
        ),
        'SELECT 1;',
        'CREATE UNIQUE INDEX idx_commission_channel_lower ON commission_structures (channel, lower_bound);'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Migration: create binder_fees reference table.
-- Uniquely keyed by (risk_rate_code, binderholder_name).

CREATE TABLE IF NOT EXISTS binder_fees (
    id                      INT AUTO_INCREMENT PRIMARY KEY,
    binderholder_name       VARCHAR(255),
    risk_rate_code          VARCHAR(255),
    maximum_binder_fee      DOUBLE DEFAULT 0,
    maximum_outsource_fee   DOUBLE DEFAULT 0,
    creation_date           DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by              VARCHAR(255)
);

-- If the table pre-exists (e.g. via GORM AutoMigrate), make sure the
-- key columns are VARCHAR(255) rather than TEXT/LONGTEXT — MySQL cannot
-- index TEXT without a prefix length (error 1170).
ALTER TABLE binder_fees MODIFY COLUMN binderholder_name VARCHAR(255);
ALTER TABLE binder_fees MODIFY COLUMN risk_rate_code VARCHAR(255);

SET @sql := (
    SELECT IF(
        EXISTS(
            SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS
            WHERE TABLE_SCHEMA = DATABASE()
              AND TABLE_NAME = 'binder_fees'
              AND INDEX_NAME = 'idx_binder_fee_rrc_holder'
        ),
        'SELECT 1;',
        'CREATE UNIQUE INDEX idx_binder_fee_rrc_holder ON binder_fees (risk_rate_code, binderholder_name);'
    )
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

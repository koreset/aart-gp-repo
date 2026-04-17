-- Migration: create binder_fees reference table.
-- Uniquely keyed by (risk_rate_code, binderholder_name).

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'binder_fees')
BEGIN
    CREATE TABLE binder_fees (
        id                      INT IDENTITY(1,1) PRIMARY KEY,
        binderholder_name       NVARCHAR(255),
        risk_rate_code          NVARCHAR(255),
        maximum_binder_fee      FLOAT DEFAULT 0,
        maximum_outsource_fee   FLOAT DEFAULT 0,
        creation_date           DATETIME DEFAULT CURRENT_TIMESTAMP,
        created_by              NVARCHAR(255)
    );
END

-- Guard types in case the table pre-existed with a different shape.
IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('binder_fees') AND name = 'binderholder_name')
BEGIN
    ALTER TABLE binder_fees ALTER COLUMN binderholder_name NVARCHAR(255);
END

IF EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('binder_fees') AND name = 'risk_rate_code')
BEGIN
    ALTER TABLE binder_fees ALTER COLUMN risk_rate_code NVARCHAR(255);
END

IF NOT EXISTS(SELECT 1 FROM sys.indexes WHERE object_id = OBJECT_ID('binder_fees') AND name = 'idx_binder_fee_rrc_holder')
BEGIN
    CREATE UNIQUE INDEX idx_binder_fee_rrc_holder
        ON binder_fees (risk_rate_code, binderholder_name);
END

-- Migration for struct: TableConfiguration
-- Table: table_configurations
--
-- Per-table-type "is required" flag for Group Pricing rating tables.
-- Rows are seeded at startup by services.EnsureTableConfigurations().

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'table_configurations')
BEGIN
    CREATE TABLE table_configurations (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configurations' AND COLUMN_NAME = 'table_type')
    ALTER TABLE table_configurations ADD table_type NVARCHAR(128);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configurations' AND COLUMN_NAME = 'display_name')
    ALTER TABLE table_configurations ADD display_name NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configurations' AND COLUMN_NAME = 'category')
    ALTER TABLE table_configurations ADD category NVARCHAR(64);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configurations' AND COLUMN_NAME = 'is_required')
    ALTER TABLE table_configurations ADD is_required BIT DEFAULT 1;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configurations' AND COLUMN_NAME = 'updated_by')
    ALTER TABLE table_configurations ADD updated_by NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configurations' AND COLUMN_NAME = 'updated_at')
    ALTER TABLE table_configurations ADD updated_at DATETIME;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configurations' AND COLUMN_NAME = 'created_at')
    ALTER TABLE table_configurations ADD created_at DATETIME;

IF NOT EXISTS(SELECT * FROM sys.indexes WHERE name = 'idx_table_configurations_table_type' AND object_id = OBJECT_ID('table_configurations'))
    CREATE UNIQUE INDEX idx_table_configurations_table_type ON table_configurations (table_type);

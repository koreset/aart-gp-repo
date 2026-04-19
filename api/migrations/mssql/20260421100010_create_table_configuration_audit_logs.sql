-- Migration for struct: TableConfigurationAuditLog
-- Table: table_configuration_audit_logs
--
-- Append-only audit history of every change to TableConfiguration.IsRequired.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'table_configuration_audit_logs')
BEGIN
    CREATE TABLE table_configuration_audit_logs (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configuration_audit_logs' AND COLUMN_NAME = 'table_type')
    ALTER TABLE table_configuration_audit_logs ADD table_type NVARCHAR(128);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configuration_audit_logs' AND COLUMN_NAME = 'event_type')
    ALTER TABLE table_configuration_audit_logs ADD event_type NVARCHAR(64);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configuration_audit_logs' AND COLUMN_NAME = 'old_value')
    ALTER TABLE table_configuration_audit_logs ADD old_value BIT;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configuration_audit_logs' AND COLUMN_NAME = 'new_value')
    ALTER TABLE table_configuration_audit_logs ADD new_value BIT;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configuration_audit_logs' AND COLUMN_NAME = 'changed_by')
    ALTER TABLE table_configuration_audit_logs ADD changed_by NVARCHAR(255);

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configuration_audit_logs' AND COLUMN_NAME = 'changed_at')
    ALTER TABLE table_configuration_audit_logs ADD changed_at DATETIME;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'table_configuration_audit_logs' AND COLUMN_NAME = 'details')
    ALTER TABLE table_configuration_audit_logs ADD details NVARCHAR(MAX);

IF NOT EXISTS(SELECT * FROM sys.indexes WHERE name = 'idx_table_configuration_audit_logs_table_type' AND object_id = OBJECT_ID('table_configuration_audit_logs'))
    CREATE INDEX idx_table_configuration_audit_logs_table_type ON table_configuration_audit_logs (table_type);

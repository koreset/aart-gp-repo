-- Migration for struct: TableConfigurationAuditLog
-- Table: table_configuration_audit_logs
--
-- Append-only audit history of every change to TableConfiguration.IsRequired.

CREATE TABLE IF NOT EXISTS table_configuration_audit_logs (
    id SERIAL PRIMARY KEY
);

ALTER TABLE table_configuration_audit_logs ADD COLUMN IF NOT EXISTS table_type VARCHAR(128);
ALTER TABLE table_configuration_audit_logs ADD COLUMN IF NOT EXISTS event_type VARCHAR(64);
ALTER TABLE table_configuration_audit_logs ADD COLUMN IF NOT EXISTS old_value BOOLEAN;
ALTER TABLE table_configuration_audit_logs ADD COLUMN IF NOT EXISTS new_value BOOLEAN;
ALTER TABLE table_configuration_audit_logs ADD COLUMN IF NOT EXISTS changed_by VARCHAR(255);
ALTER TABLE table_configuration_audit_logs ADD COLUMN IF NOT EXISTS changed_at TIMESTAMP;
ALTER TABLE table_configuration_audit_logs ADD COLUMN IF NOT EXISTS details TEXT;

CREATE INDEX IF NOT EXISTS idx_table_configuration_audit_logs_table_type
    ON table_configuration_audit_logs (table_type);

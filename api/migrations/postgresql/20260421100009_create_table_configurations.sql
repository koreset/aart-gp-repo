-- Migration for struct: TableConfiguration
-- Table: table_configurations
--
-- Per-table-type "is required" flag for Group Pricing rating tables.
-- Rows are seeded at startup by services.EnsureTableConfigurations() from
-- the canonical gpTableSpecs list, so no SQL seed step is required here.

CREATE TABLE IF NOT EXISTS table_configurations (
    id SERIAL PRIMARY KEY
);

ALTER TABLE table_configurations ADD COLUMN IF NOT EXISTS table_type VARCHAR(128);
ALTER TABLE table_configurations ADD COLUMN IF NOT EXISTS display_name VARCHAR(255);
ALTER TABLE table_configurations ADD COLUMN IF NOT EXISTS category VARCHAR(64);
ALTER TABLE table_configurations ADD COLUMN IF NOT EXISTS is_required BOOLEAN DEFAULT TRUE;
ALTER TABLE table_configurations ADD COLUMN IF NOT EXISTS updated_by VARCHAR(255);
ALTER TABLE table_configurations ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP;
ALTER TABLE table_configurations ADD COLUMN IF NOT EXISTS created_at TIMESTAMP;

CREATE UNIQUE INDEX IF NOT EXISTS idx_table_configurations_table_type
    ON table_configurations (table_type);

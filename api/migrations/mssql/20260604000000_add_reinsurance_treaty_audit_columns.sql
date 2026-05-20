-- Treaty activation/deactivation accountability: adds activated_by /
-- activated_at / deactivated_by / deactivated_at to reinsurance_treaties so
-- the Treaty Management screen can record who promoted a treaty out of draft
-- and who later took it out of service. Activated treaties cannot be deleted
-- (kept for traceability); deactivation transitions status to 'cancelled' and
-- stamps the deactivated_* fields.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('reinsurance_treaties') AND name = 'activated_by')
  ALTER TABLE reinsurance_treaties ADD activated_by NVARCHAR(255) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('reinsurance_treaties') AND name = 'activated_at')
  ALTER TABLE reinsurance_treaties ADD activated_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('reinsurance_treaties') AND name = 'deactivated_by')
  ALTER TABLE reinsurance_treaties ADD deactivated_by NVARCHAR(255) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('reinsurance_treaties') AND name = 'deactivated_at')
  ALTER TABLE reinsurance_treaties ADD deactivated_at DATETIME NULL;

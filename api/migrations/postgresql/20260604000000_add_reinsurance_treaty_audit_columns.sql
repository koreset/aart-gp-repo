-- Treaty activation/deactivation accountability: adds activated_by /
-- activated_at / deactivated_by / deactivated_at to reinsurance_treaties so
-- the Treaty Management screen can record who promoted a treaty out of draft
-- and who later took it out of service. Activated treaties cannot be deleted
-- (kept for traceability); deactivation transitions status to 'cancelled' and
-- stamps the deactivated_* fields.
-- Idempotent on re-runs.

ALTER TABLE reinsurance_treaties
    ADD COLUMN IF NOT EXISTS activated_by VARCHAR(255) NULL;

ALTER TABLE reinsurance_treaties
    ADD COLUMN IF NOT EXISTS activated_at TIMESTAMP NULL;

ALTER TABLE reinsurance_treaties
    ADD COLUMN IF NOT EXISTS deactivated_by VARCHAR(255) NULL;

ALTER TABLE reinsurance_treaties
    ADD COLUMN IF NOT EXISTS deactivated_at TIMESTAMP NULL;

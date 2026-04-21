-- Migration: add request_payload to generated_bordereauxes so draft rows can
-- be regenerated in place using the original request parameters.

ALTER TABLE generated_bordereauxes
    ADD COLUMN IF NOT EXISTS request_payload JSONB;

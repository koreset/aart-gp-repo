-- Migration: add request_payload to generated_bordereauxes so draft rows can
-- be regenerated in place using the original request parameters.

IF NOT EXISTS(SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('generated_bordereauxes') AND name = 'request_payload')
BEGIN
    ALTER TABLE generated_bordereauxes ADD request_payload NVARCHAR(MAX) NULL;
END

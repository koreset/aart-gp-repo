-- Add resolution_notes column to claim_payment_schedule_queries so finance can
-- type a reply when resolving a claims follow-up or line query. Existing rows
-- get NULL (treated as empty by the model). Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('claim_payment_schedule_queries') AND name = 'resolution_notes')
BEGIN
    ALTER TABLE claim_payment_schedule_queries
        ADD resolution_notes NVARCHAR(MAX) NULL;
END;

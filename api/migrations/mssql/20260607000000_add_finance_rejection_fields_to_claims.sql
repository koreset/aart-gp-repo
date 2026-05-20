-- Snapshot columns powering the finance-rejection banner on the claim
-- detail view. Populated when a payment-schedule line is rejected; cleared
-- on the next assessor re-approval. Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claims') AND name = 'finance_rejected_at')
  ALTER TABLE group_scheme_claims ADD finance_rejected_at DATETIME NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claims') AND name = 'finance_rejected_by')
  ALTER TABLE group_scheme_claims ADD finance_rejected_by NVARCHAR(255) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claims') AND name = 'finance_rejection_reason_code')
  ALTER TABLE group_scheme_claims ADD finance_rejection_reason_code NVARCHAR(64) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claims') AND name = 'finance_rejection_notes')
  ALTER TABLE group_scheme_claims ADD finance_rejection_notes NVARCHAR(MAX) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('group_scheme_claims') AND name = 'finance_rejection_schedule_number')
  ALTER TABLE group_scheme_claims ADD finance_rejection_schedule_number NVARCHAR(191) NULL;

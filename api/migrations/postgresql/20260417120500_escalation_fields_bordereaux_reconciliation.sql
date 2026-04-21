-- Migration: add escalation workflow fields to bordereaux_reconciliation_results
-- Purpose: support P1 escalation workflow — assignee, SLA due date,
-- priority, and audit timestamps so escalated discrepancies can be queued,
-- surfaced as overdue, and routed to a specific user.
-- See services/bordereaux_reconciliation.go EscalateDiscrepancy.

ALTER TABLE bordereaux_reconciliation_results ADD COLUMN IF NOT EXISTS assigned_to VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE bordereaux_reconciliation_results ADD COLUMN IF NOT EXISTS priority VARCHAR(32) NOT NULL DEFAULT '';
ALTER TABLE bordereaux_reconciliation_results ADD COLUMN IF NOT EXISTS escalated_by VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE bordereaux_reconciliation_results ADD COLUMN IF NOT EXISTS escalated_at TIMESTAMP WITH TIME ZONE NULL;
ALTER TABLE bordereaux_reconciliation_results ADD COLUMN IF NOT EXISTS due_date TIMESTAMP WITH TIME ZONE NULL;

CREATE INDEX IF NOT EXISTS idx_brr_status_due_date ON bordereaux_reconciliation_results (status, due_date);
CREATE INDEX IF NOT EXISTS idx_brr_assigned_to ON bordereaux_reconciliation_results (assigned_to);

-- Migration: add escalation workflow fields to bordereaux_reconciliation_results
-- Purpose: support P1 escalation workflow — assignee, SLA due date,
-- priority, and audit timestamps so escalated discrepancies can be queued,
-- surfaced as overdue, and routed to a specific user.
-- See services/bordereaux_reconciliation.go EscalateDiscrepancy.

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='assigned_to')
BEGIN
    ALTER TABLE bordereaux_reconciliation_results ADD assigned_to NVARCHAR(255) NOT NULL DEFAULT '';
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='priority')
BEGIN
    ALTER TABLE bordereaux_reconciliation_results ADD priority NVARCHAR(32) NOT NULL DEFAULT '';
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='escalated_by')
BEGIN
    ALTER TABLE bordereaux_reconciliation_results ADD escalated_by NVARCHAR(255) NOT NULL DEFAULT '';
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='escalated_at')
BEGIN
    ALTER TABLE bordereaux_reconciliation_results ADD escalated_at DATETIME2 NULL;
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='due_date')
BEGIN
    ALTER TABLE bordereaux_reconciliation_results ADD due_date DATETIME2 NULL;
END;

IF NOT EXISTS(SELECT * FROM sys.indexes WHERE name='idx_brr_status_due_date' AND object_id=OBJECT_ID('bordereaux_reconciliation_results'))
BEGIN
    CREATE INDEX idx_brr_status_due_date ON bordereaux_reconciliation_results (status, due_date);
END;

IF NOT EXISTS(SELECT * FROM sys.indexes WHERE name='idx_brr_assigned_to' AND object_id=OBJECT_ID('bordereaux_reconciliation_results'))
BEGIN
    CREATE INDEX idx_brr_assigned_to ON bordereaux_reconciliation_results (assigned_to);
END;

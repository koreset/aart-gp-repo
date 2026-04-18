-- Migration: add escalation workflow fields to bordereaux_reconciliation_results
-- Purpose: support P1 escalation workflow — assignee, SLA due date,
-- priority, and audit timestamps so escalated discrepancies can be queued,
-- surfaced as overdue, and routed to a specific user.
-- See services/bordereaux_reconciliation.go EscalateDiscrepancy.

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='assigned_to' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE bordereaux_reconciliation_results ADD COLUMN assigned_to VARCHAR(255) NOT NULL DEFAULT "";'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='priority' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE bordereaux_reconciliation_results ADD COLUMN priority VARCHAR(32) NOT NULL DEFAULT "";'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='escalated_by' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE bordereaux_reconciliation_results ADD COLUMN escalated_by VARCHAR(255) NOT NULL DEFAULT "";'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='escalated_at' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE bordereaux_reconciliation_results ADD COLUMN escalated_at DATETIME NULL;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND COLUMN_NAME='due_date' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'ALTER TABLE bordereaux_reconciliation_results ADD COLUMN due_date DATETIME NULL;'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- MySQL requires a key-length prefix when indexing TEXT columns. The existing
-- `status` column was created without an explicit size (GORM default → TEXT),
-- so the composite index uses a 64-char prefix. This is still selective enough
-- for the queue queries ("status = 'escalated'").
SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND INDEX_NAME='idx_brr_status_due_date' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE INDEX idx_brr_status_due_date ON bordereaux_reconciliation_results (status(64), due_date);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @s = (SELECT IF(
    EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_NAME='bordereaux_reconciliation_results' AND INDEX_NAME='idx_brr_assigned_to' AND TABLE_SCHEMA = DATABASE()),
    'SELECT 1;',
    'CREATE INDEX idx_brr_assigned_to ON bordereaux_reconciliation_results (assigned_to);'
));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

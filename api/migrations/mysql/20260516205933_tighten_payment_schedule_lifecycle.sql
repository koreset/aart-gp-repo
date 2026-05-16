-- Generated 2026-05-16T20:59:33+02:00 for dialect mysql

-- Migration for: ClaimPaymentSchedule (table: claim_payment_schedules)

ALTER TABLE claim_payment_schedules ADD COLUMN notes varchar(30);
ALTER TABLE claim_payment_schedules ADD COLUMN gross_total DOUBLE;
ALTER TABLE claim_payment_schedules ADD COLUMN deductions_total DOUBLE;
ALTER TABLE claim_payment_schedules ADD COLUMN net_total DOUBLE;
ALTER TABLE claim_payment_schedules ADD COLUMN locked_at DATETIME;
ALTER TABLE claim_payment_schedules ADD COLUMN head_of_claims_signed_by VARCHAR(255);
ALTER TABLE claim_payment_schedules ADD COLUMN head_of_claims_signed_at DATETIME;
ALTER TABLE claim_payment_schedules ADD COLUMN finance_review_started_by VARCHAR(255);
ALTER TABLE claim_payment_schedules ADD COLUMN finance_review_started_at DATETIME;
ALTER TABLE claim_payment_schedules ADD COLUMN finance_first_auth_by VARCHAR(255);
ALTER TABLE claim_payment_schedules ADD COLUMN finance_first_auth_at DATETIME;
ALTER TABLE claim_payment_schedules ADD COLUMN finance_second_auth_by VARCHAR(255);
ALTER TABLE claim_payment_schedules ADD COLUMN finance_second_auth_at DATETIME;
ALTER TABLE claim_payment_schedules ADD COLUMN submitted_to_bank_at DATETIME;
ALTER TABLE claim_payment_schedules ADD COLUMN archived_at DATETIME;
ALTER TABLE claim_payment_schedules ADD COLUMN archived_by VARCHAR(255);

-- Migration for: ClaimPaymentScheduleItem (table: claim_payment_schedule_items)

ALTER TABLE claim_payment_schedule_items ADD COLUMN gross_amount DOUBLE;
ALTER TABLE claim_payment_schedule_items ADD COLUMN premium_arrears_deduction DOUBLE;
ALTER TABLE claim_payment_schedule_items ADD COLUMN policy_loan_deduction DOUBLE;
ALTER TABLE claim_payment_schedule_items ADD COLUMN tax_withheld DOUBLE;
ALTER TABLE claim_payment_schedule_items ADD COLUMN net_payable DOUBLE;
ALTER TABLE claim_payment_schedule_items ADD COLUMN beneficiary_name VARCHAR(255);
ALTER TABLE claim_payment_schedule_items ADD COLUMN beneficiary_id_number VARCHAR(255);
ALTER TABLE claim_payment_schedule_items ADD COLUMN risk_flags BLOB;
ALTER TABLE claim_payment_schedule_items ADD COLUMN approval_reference VARCHAR(255);
ALTER TABLE claim_payment_schedule_items ADD COLUMN line_status VARCHAR(32) DEFAULT 'pending';
ALTER TABLE claim_payment_schedule_items ADD COLUMN verified_by VARCHAR(255);
ALTER TABLE claim_payment_schedule_items ADD COLUMN verified_at DATETIME;
ALTER TABLE claim_payment_schedule_items ADD COLUMN query_reason_code VARCHAR(64);
ALTER TABLE claim_payment_schedule_items ADD COLUMN query_notes VARCHAR(255);
ALTER TABLE claim_payment_schedule_items ADD COLUMN queried_by VARCHAR(255);
ALTER TABLE claim_payment_schedule_items ADD COLUMN queried_at DATETIME;

-- Create table: claim_payment_schedule_queries
CREATE TABLE claim_payment_schedule_queries (
    id BIGINT AUTO_INCREMENT,
    schedule_id BIGINT,
    schedule_item_id BIGINT,
    claim_id BIGINT,
    claim_number VARCHAR(255),
    reason_code VARCHAR(64),
    notes VARCHAR(255),
    outcome VARCHAR(32) DEFAULT 'open',
    raised_by VARCHAR(255),
    raised_at DATETIME,
    resolved_by VARCHAR(255),
    resolved_at DATETIME,
    PRIMARY KEY (id)
);

CREATE INDEX idx_claim_payment_schedule_queries_schedule_id ON claim_payment_schedule_queries (schedule_id);
CREATE INDEX idx_claim_payment_schedule_queries_schedule_item_id ON claim_payment_schedule_queries (schedule_item_id);
CREATE INDEX idx_claim_payment_schedule_queries_claim_id ON claim_payment_schedule_queries (claim_id);
CREATE INDEX idx_claim_payment_schedule_queries_reason_code ON claim_payment_schedule_queries (reason_code);


-- Create table: authority_matrix
CREATE TABLE authority_matrix (
    id BIGINT AUTO_INCREMENT,
    role VARCHAR(64),
    action VARCHAR(64),
    min_amount DOUBLE,
    max_amount DOUBLE,
    is_active TINYINT(1) DEFAULT 1,
    created_by VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE INDEX idx_authority_matrix_role ON authority_matrix (role);
CREATE INDEX idx_authority_matrix_action ON authority_matrix (action);


-- Create table: payment_schedule_audits
CREATE TABLE payment_schedule_audits (
    id BIGINT AUTO_INCREMENT,
    schedule_id BIGINT,
    from_status VARCHAR(64),
    to_status VARCHAR(64),
    actor VARCHAR(255),
    notes VARCHAR(255),
    changed_at DATETIME,
    PRIMARY KEY (id)
);

CREATE INDEX idx_payment_schedule_audits_schedule_id ON payment_schedule_audits (schedule_id);


-- Backwards-compat: existing schedules created before the new lifecycle
-- need their legacy status mapped onto the new state machine.
-- "submitted" with an ACB file generated → finance_second_authorised (or
-- submitted_to_bank — picking the more conservative state that still lets
-- proof upload close the loop).
-- "submitted" without an ACB file → claims_signed_off (waiting for finance).
-- "confirmed" stays as-is.
UPDATE claim_payment_schedules
SET status = 'submitted_to_bank'
WHERE LOWER(status) = 'submitted' AND acb_file_generated = 1;

UPDATE claim_payment_schedules
SET status = 'claims_signed_off'
WHERE LOWER(status) = 'submitted' AND (acb_file_generated IS NULL OR acb_file_generated = 0);

-- Existing line items default to "verified" since they were already
-- effectively passed through the old (no-review) flow.
UPDATE claim_payment_schedule_items
SET line_status = 'verified'
WHERE line_status IS NULL OR line_status = '';

-- Seed gross_amount / net_payable from the existing claim_amount so the
-- breakdown is visible on already-generated schedules.
UPDATE claim_payment_schedule_items
SET gross_amount = claim_amount,
    net_payable  = claim_amount
WHERE gross_amount IS NULL OR gross_amount = 0;

UPDATE claim_payment_schedules s
JOIN (
    SELECT schedule_id,
           SUM(gross_amount) AS gross,
           SUM(COALESCE(premium_arrears_deduction,0) + COALESCE(policy_loan_deduction,0) + COALESCE(tax_withheld,0)) AS ded,
           SUM(net_payable) AS net
    FROM claim_payment_schedule_items
    GROUP BY schedule_id
) t ON t.schedule_id = s.id
SET s.gross_total      = t.gross,
    s.deductions_total = t.ded,
    s.net_total        = t.net
WHERE s.gross_total IS NULL OR s.gross_total = 0;


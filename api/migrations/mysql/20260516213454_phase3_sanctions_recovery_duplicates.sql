-- Generated 2026-05-16T21:34:54+02:00 for dialect mysql
-- Phase 3 deltas only. Phase 1's ClaimPaymentScheduleItem columns are added
-- in the earlier migration 20260516205933_tighten_payment_schedule_lifecycle.

-- Migration for: ClaimPaymentScheduleItem (Phase 3 additions)

ALTER TABLE claim_payment_schedule_items ADD COLUMN reinsurance_recovery_required TINYINT(1);
ALTER TABLE claim_payment_schedule_items ADD COLUMN reinsurance_recovery_amount DOUBLE;
ALTER TABLE claim_payment_schedule_items ADD COLUMN reinsurance_recovery_raised_by VARCHAR(255);
ALTER TABLE claim_payment_schedule_items ADD COLUMN reinsurance_recovery_raised_at DATETIME;
ALTER TABLE claim_payment_schedule_items ADD COLUMN duplicate_beneficiary_flag TINYINT(1);
ALTER TABLE claim_payment_schedule_items ADD COLUMN duplicate_beneficiary_cleared TINYINT(1);

-- Create table: sanctions_screenings
CREATE TABLE sanctions_screenings (
    id BIGINT AUTO_INCREMENT,
    schedule_id BIGINT,
    schedule_item_id BIGINT,
    claim_id BIGINT,
    provider VARCHAR(64),
    status VARCHAR(32),
    provider_ref VARCHAR(128),
    hit_summary VARCHAR(255),
    notes VARCHAR(255),
    screened_by VARCHAR(255),
    screened_at DATETIME,
    cleared_by VARCHAR(255),
    cleared_at DATETIME,
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE INDEX idx_sanctions_screenings_schedule_id ON sanctions_screenings (schedule_id);
CREATE INDEX idx_sanctions_screenings_schedule_item_id ON sanctions_screenings (schedule_item_id);
CREATE UNIQUE INDEX idx_sanctions_item_provider ON sanctions_screenings (schedule_item_id, provider);
CREATE INDEX idx_sanctions_screenings_claim_id ON sanctions_screenings (claim_id);

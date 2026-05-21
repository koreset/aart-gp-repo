-- Migration: create bulk_enrollment_batches and add per-row batch tracking
-- columns to g_pricing_member_data_in_forces.
--
-- A batch is a single CSV upload pending approval. Members from the upload
-- are stored as status='draft' linked to the batch via batch_id; member-read
-- paths exclude drafts so they don't reach the live members grid, premium
-- calc, or exposure dashboards until the batch is approved.

CREATE TABLE IF NOT EXISTS bulk_enrollment_batches (
    id                    BIGINT AUTO_INCREMENT PRIMARY KEY,
    scheme_id             BIGINT NOT NULL,
    quote_id              BIGINT,
    status                VARCHAR(32) NOT NULL DEFAULT 'pending_approval',
    member_count          BIGINT,
    valid_count           BIGINT,
    blocking_count        BIGINT,
    soft_error_count      BIGINT,
    file_name             VARCHAR(512),
    file_size_bytes       BIGINT,
    file_checksum         VARCHAR(64),
    skip_duplicates       BOOLEAN,
    validation_report     LONGTEXT,
    external_id_check_run BOOLEAN,
    external_id_check_at  DATETIME(3),
    uploaded_by           VARCHAR(128),
    uploaded_at           DATETIME(3),
    approved_by           VARCHAR(128),
    approved_at           DATETIME(3),
    rejected_by           VARCHAR(128),
    rejected_at           DATETIME(3),
    rejected_reason       VARCHAR(1000),
    notes                 VARCHAR(1000)
);

CREATE INDEX idx_bulk_enrollment_batches_scheme_id ON bulk_enrollment_batches(scheme_id);
CREATE INDEX idx_bulk_enrollment_batches_status ON bulk_enrollment_batches(status);
CREATE INDEX idx_bulk_enrollment_batches_approved_at ON bulk_enrollment_batches(approved_at);
CREATE INDEX idx_bulk_enrollment_batches_rejected_at ON bulk_enrollment_batches(rejected_at);

-- row_number is a reserved word in MySQL 8 (ROW_NUMBER() window function),
-- so the column is named row_index instead.
ALTER TABLE g_pricing_member_data_in_forces
    ADD COLUMN batch_id          BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN row_index         BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN validation_status VARCHAR(16) NOT NULL DEFAULT '';

CREATE INDEX idx_g_pricing_member_data_in_forces_batch_id
    ON g_pricing_member_data_in_forces(batch_id);

-- Migration: create bulk_enrollment_batches and add per-row batch tracking
-- columns to g_pricing_member_data_in_forces.
--
-- A batch is a single CSV upload pending approval. Members from the upload
-- are stored as status='draft' linked to the batch via batch_id; member-read
-- paths exclude drafts so they don't reach the live members grid, premium
-- calc, or exposure dashboards until the batch is approved.

CREATE TABLE IF NOT EXISTS bulk_enrollment_batches (
    id                    BIGSERIAL PRIMARY KEY,
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
    validation_report     TEXT,
    external_id_check_run BOOLEAN,
    external_id_check_at  TIMESTAMP WITH TIME ZONE,
    uploaded_by           VARCHAR(128),
    uploaded_at           TIMESTAMP WITH TIME ZONE,
    approved_by           VARCHAR(128),
    approved_at           TIMESTAMP WITH TIME ZONE,
    rejected_by           VARCHAR(128),
    rejected_at           TIMESTAMP WITH TIME ZONE,
    rejected_reason       VARCHAR(1000),
    notes                 VARCHAR(1000)
);

CREATE INDEX IF NOT EXISTS idx_bulk_enrollment_batches_scheme_id ON bulk_enrollment_batches(scheme_id);
CREATE INDEX IF NOT EXISTS idx_bulk_enrollment_batches_status ON bulk_enrollment_batches(status);
CREATE INDEX IF NOT EXISTS idx_bulk_enrollment_batches_approved_at ON bulk_enrollment_batches(approved_at);
CREATE INDEX IF NOT EXISTS idx_bulk_enrollment_batches_rejected_at ON bulk_enrollment_batches(rejected_at);

-- row_number is a reserved word in MySQL 8; using row_index keeps the schema
-- consistent across dialects.
ALTER TABLE g_pricing_member_data_in_forces
    ADD COLUMN IF NOT EXISTS batch_id          BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS row_index         BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS validation_status VARCHAR(16) NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_g_pricing_member_data_in_forces_batch_id
    ON g_pricing_member_data_in_forces(batch_id);

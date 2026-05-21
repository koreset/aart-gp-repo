-- Migration: create bulk_enrollment_batches and add per-row batch tracking
-- columns to g_pricing_member_data_in_forces.
--
-- A batch is a single CSV upload pending approval. Members from the upload
-- are stored as status='draft' linked to the batch via batch_id; member-read
-- paths exclude drafts so they don't reach the live members grid, premium
-- calc, or exposure dashboards until the batch is approved.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'bulk_enrollment_batches')
BEGIN
    CREATE TABLE bulk_enrollment_batches (
        id                    BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_bulk_enrollment_batches PRIMARY KEY,
        scheme_id             BIGINT NOT NULL,
        quote_id              BIGINT NULL,
        status                NVARCHAR(32) NOT NULL CONSTRAINT df_bulk_enrollment_batches_status DEFAULT ('pending_approval'),
        member_count          BIGINT NULL,
        valid_count           BIGINT NULL,
        blocking_count        BIGINT NULL,
        soft_error_count      BIGINT NULL,
        file_name             NVARCHAR(512) NULL,
        file_size_bytes       BIGINT NULL,
        file_checksum         NVARCHAR(64) NULL,
        skip_duplicates       BIT NULL,
        validation_report     NVARCHAR(MAX) NULL,
        external_id_check_run BIT NULL,
        external_id_check_at  DATETIMEOFFSET NULL,
        uploaded_by           NVARCHAR(128) NULL,
        uploaded_at           DATETIMEOFFSET NULL,
        approved_by           NVARCHAR(128) NULL,
        approved_at           DATETIMEOFFSET NULL,
        rejected_by           NVARCHAR(128) NULL,
        rejected_at           DATETIMEOFFSET NULL,
        rejected_reason       NVARCHAR(1000) NULL,
        notes                 NVARCHAR(1000) NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bulk_enrollment_batches_scheme_id')
    CREATE INDEX idx_bulk_enrollment_batches_scheme_id ON bulk_enrollment_batches(scheme_id);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bulk_enrollment_batches_status')
    CREATE INDEX idx_bulk_enrollment_batches_status ON bulk_enrollment_batches(status);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bulk_enrollment_batches_approved_at')
    CREATE INDEX idx_bulk_enrollment_batches_approved_at ON bulk_enrollment_batches(approved_at);
IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_bulk_enrollment_batches_rejected_at')
    CREATE INDEX idx_bulk_enrollment_batches_rejected_at ON bulk_enrollment_batches(rejected_at);

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE Name = N'batch_id' AND Object_ID = Object_ID(N'g_pricing_member_data_in_forces'))
    ALTER TABLE g_pricing_member_data_in_forces ADD batch_id BIGINT NOT NULL CONSTRAINT df_gpmdif_batch_id DEFAULT 0;
-- row_number is a reserved word in MySQL 8; using row_index keeps the schema
-- consistent across dialects.
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE Name = N'row_index' AND Object_ID = Object_ID(N'g_pricing_member_data_in_forces'))
    ALTER TABLE g_pricing_member_data_in_forces ADD row_index BIGINT NOT NULL CONSTRAINT df_gpmdif_row_index DEFAULT 0;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE Name = N'validation_status' AND Object_ID = Object_ID(N'g_pricing_member_data_in_forces'))
    ALTER TABLE g_pricing_member_data_in_forces ADD validation_status NVARCHAR(16) NOT NULL CONSTRAINT df_gpmdif_validation_status DEFAULT '';

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_g_pricing_member_data_in_forces_batch_id')
    CREATE INDEX idx_g_pricing_member_data_in_forces_batch_id ON g_pricing_member_data_in_forces(batch_id);

-- Migration for struct: InsurerQuoteTemplate

-- Table: insurer_quote_templates

-- Create table if not exists
CREATE TABLE IF NOT EXISTS insurer_quote_templates (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: InsurerID
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='insurer_id') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN insurer_id INTEGER NOT NULL;
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN insurer_id SET NOT NULL;
    END IF;
END $$;

-- Create index on insurer_id
CREATE INDEX IF NOT EXISTS idx_insurer_quote_templates_insurer_id ON insurer_quote_templates(insurer_id);

-- Add or modify column for field: Version
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='version') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN version INTEGER;
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN version SET DATA TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Filename
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='filename') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN filename VARCHAR(255);
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN filename SET DATA TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: DocxBlob
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='docx_blob') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN docx_blob BYTEA;
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN docx_blob SET DATA TYPE BYTEA;
    END IF;
END $$;

-- Add or modify column for field: SizeBytes
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='size_bytes') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN size_bytes INTEGER;
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN size_bytes SET DATA TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: UploadedBy
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='uploaded_by') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN uploaded_by VARCHAR(255);
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN uploaded_by SET DATA TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: UploadedAt
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='uploaded_at') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN uploaded_at TIMESTAMP;
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN uploaded_at SET DATA TYPE TIMESTAMP;
    END IF;
END $$;

-- Add or modify column for field: IsActive
DO $$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='insurer_quote_templates' AND COLUMN_NAME='is_active') THEN
        ALTER TABLE insurer_quote_templates ADD COLUMN is_active BOOLEAN;
    ELSE
        ALTER TABLE insurer_quote_templates ALTER COLUMN is_active SET DATA TYPE BOOLEAN;
    END IF;
END $$;

-- Create index on is_active
CREATE INDEX IF NOT EXISTS idx_insurer_quote_templates_is_active ON insurer_quote_templates(is_active);

-- Migration for struct: InsurerQuoteTemplate

-- Table: insurer_quote_templates

-- Create table if not exists
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = N'insurer_quote_templates')
BEGIN
    CREATE TABLE insurer_quote_templates (
        id INT IDENTITY(1,1) PRIMARY KEY
    )
END

-- Add or modify column for field: InsurerID
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'insurer_id')
BEGIN
    ALTER TABLE insurer_quote_templates ADD insurer_id INT NOT NULL
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN insurer_id INT NOT NULL
END

-- Create index on insurer_id
IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = N'idx_insurer_id' AND object_id = OBJECT_ID(N'insurer_quote_templates'))
BEGIN
    CREATE INDEX idx_insurer_id ON insurer_quote_templates(insurer_id)
END

-- Add or modify column for field: Version
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'version')
BEGIN
    ALTER TABLE insurer_quote_templates ADD version INT
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN version INT
END

-- Add or modify column for field: Filename
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'filename')
BEGIN
    ALTER TABLE insurer_quote_templates ADD filename NVARCHAR(255)
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN filename NVARCHAR(255)
END

-- Add or modify column for field: DocxBlob
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'docx_blob')
BEGIN
    ALTER TABLE insurer_quote_templates ADD docx_blob VARBINARY(MAX)
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN docx_blob VARBINARY(MAX)
END

-- Add or modify column for field: SizeBytes
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'size_bytes')
BEGIN
    ALTER TABLE insurer_quote_templates ADD size_bytes INT
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN size_bytes INT
END

-- Add or modify column for field: UploadedBy
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'uploaded_by')
BEGIN
    ALTER TABLE insurer_quote_templates ADD uploaded_by NVARCHAR(255)
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN uploaded_by NVARCHAR(255)
END

-- Add or modify column for field: UploadedAt
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'uploaded_at')
BEGIN
    ALTER TABLE insurer_quote_templates ADD uploaded_at DATETIME2
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN uploaded_at DATETIME2
END

-- Add or modify column for field: IsActive
IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'insurer_quote_templates' AND COLUMN_NAME = N'is_active')
BEGIN
    ALTER TABLE insurer_quote_templates ADD is_active BIT
END
ELSE
BEGIN
    ALTER TABLE insurer_quote_templates ALTER COLUMN is_active BIT
END

-- Create index on is_active
IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = N'idx_is_active' AND object_id = OBJECT_ID(N'insurer_quote_templates'))
BEGIN
    CREATE INDEX idx_is_active ON insurer_quote_templates(is_active)
END

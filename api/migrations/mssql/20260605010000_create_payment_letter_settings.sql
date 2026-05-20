-- Singleton table (one row, id=1) holding letterhead, branding, and signatory
-- configuration used when generating claim payment confirmation letters.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'payment_letter_settings')
BEGIN
  CREATE TABLE payment_letter_settings (
    id                       BIGINT NOT NULL PRIMARY KEY,
    company_name             NVARCHAR(255) NULL,
    address_line1            NVARCHAR(255) NULL,
    address_line2            NVARCHAR(255) NULL,
    address_line3            NVARCHAR(255) NULL,
    city                     NVARCHAR(100) NULL,
    postal_code              NVARCHAR(50)  NULL,
    country                  NVARCHAR(100) NULL,
    phone                    NVARCHAR(100) NULL,
    email                    NVARCHAR(255) NULL,
    website                  NVARCHAR(255) NULL,
    logo                     VARBINARY(MAX) NULL,
    logo_mime_type           NVARCHAR(64)  NULL,
    signatory_name           NVARCHAR(255) NULL,
    signatory_title          NVARCHAR(255) NULL,
    signature                VARBINARY(MAX) NULL,
    signature_mime_type      NVARCHAR(64)  NULL,
    letter_intro_template    NVARCHAR(MAX) NULL,
    letter_closing_template  NVARCHAR(MAX) NULL,
    updated_at               DATETIME      NULL,
    updated_by               NVARCHAR(255) NULL
  );
END

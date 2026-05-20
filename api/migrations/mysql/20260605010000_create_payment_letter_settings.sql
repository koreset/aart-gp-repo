-- Singleton table (one row, id=1) holding letterhead, branding, and signatory
-- configuration used when generating claim payment confirmation letters.
-- Idempotent on re-runs.

CREATE TABLE IF NOT EXISTS payment_letter_settings (
  id                       BIGINT NOT NULL PRIMARY KEY,
  company_name             VARCHAR(255) NULL,
  address_line1            VARCHAR(255) NULL,
  address_line2            VARCHAR(255) NULL,
  address_line3            VARCHAR(255) NULL,
  city                     VARCHAR(100) NULL,
  postal_code              VARCHAR(50)  NULL,
  country                  VARCHAR(100) NULL,
  phone                    VARCHAR(100) NULL,
  email                    VARCHAR(255) NULL,
  website                  VARCHAR(255) NULL,
  logo                     LONGBLOB     NULL,
  logo_mime_type           VARCHAR(64)  NULL,
  signatory_name           VARCHAR(255) NULL,
  signatory_title          VARCHAR(255) NULL,
  signature                LONGBLOB     NULL,
  signature_mime_type      VARCHAR(64)  NULL,
  letter_intro_template    TEXT         NULL,
  letter_closing_template  TEXT         NULL,
  updated_at               DATETIME     NULL,
  updated_by               VARCHAR(255) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

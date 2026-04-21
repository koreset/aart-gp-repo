-- Migration: email sending system (Phase 1).
-- Creates per-license SMTP configuration, licensed email templates,
-- a persistent outbox for the worker, and a per-user signature field.

--------------------------------------------------------------------------------
-- email_settings
--------------------------------------------------------------------------------

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'email_settings')
BEGIN
    CREATE TABLE email_settings (
        id                      INT IDENTITY(1,1) PRIMARY KEY,
        license_id              NVARCHAR(191) NOT NULL,
        host                    NVARCHAR(255) NOT NULL,
        port                    INT NOT NULL DEFAULT 587,
        tls_mode                NVARCHAR(16) NOT NULL DEFAULT 'starttls',
        auth_user               NVARCHAR(255),
        auth_password_encrypted NVARCHAR(MAX),
        from_address            NVARCHAR(255) NOT NULL,
        from_name               NVARCHAR(255),
        reply_to                NVARCHAR(255),
        updated_by              NVARCHAR(255),
        updated_at              DATETIME2 NULL
    );
    CREATE UNIQUE INDEX uk_email_settings_license ON email_settings (license_id);
END

--------------------------------------------------------------------------------
-- email_templates
--------------------------------------------------------------------------------

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'email_templates')
BEGIN
    CREATE TABLE email_templates (
        id               INT IDENTITY(1,1) PRIMARY KEY,
        license_id       NVARCHAR(191) NOT NULL,
        code             NVARCHAR(128) NOT NULL,
        name             NVARCHAR(255),
        description      NVARCHAR(MAX),
        subject_template NVARCHAR(MAX),
        body_template    NVARCHAR(MAX),
        attachments_spec NVARCHAR(MAX),
        status           NVARCHAR(16) NOT NULL DEFAULT 'draft',
        updated_by       NVARCHAR(255),
        updated_at       DATETIME2 NULL
    );
    CREATE UNIQUE INDEX uk_email_templates_license_code ON email_templates (license_id, code);
END

--------------------------------------------------------------------------------
-- email_outbox
--------------------------------------------------------------------------------

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'email_outbox')
BEGIN
    CREATE TABLE email_outbox (
        id                  INT IDENTITY(1,1) PRIMARY KEY,
        license_id          NVARCHAR(191) NOT NULL,
        template_code       NVARCHAR(128),
        from_address        NVARCHAR(255),
        from_name           NVARCHAR(255),
        reply_to            NVARCHAR(255),
        to_recipients       NVARCHAR(MAX),
        cc_recipients       NVARCHAR(MAX),
        bcc_recipients      NVARCHAR(MAX),
        subject             NVARCHAR(MAX),
        body                NVARCHAR(MAX),
        attachments         NVARCHAR(MAX),
        status              NVARCHAR(16) NOT NULL DEFAULT 'pending',
        attempts            INT NOT NULL DEFAULT 0,
        max_attempts        INT NOT NULL DEFAULT 5,
        last_error          NVARCHAR(MAX),
        next_attempt_at     DATETIME2 NULL,
        scheduled_at        DATETIME2 NULL,
        sent_at             DATETIME2 NULL,
        related_object_type NVARCHAR(128),
        related_object_id   NVARCHAR(128),
        created_by          NVARCHAR(255),
        created_at          DATETIME2 NULL
    );
    CREATE INDEX idx_email_outbox_status_next ON email_outbox (status, next_attempt_at);
    CREATE INDEX idx_email_outbox_license ON email_outbox (license_id);
    CREATE INDEX idx_email_outbox_template ON email_outbox (template_code);
END

--------------------------------------------------------------------------------
-- gp_users.email_signature
--------------------------------------------------------------------------------

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('gp_users') AND name = 'email_signature')
BEGIN
    ALTER TABLE gp_users ADD email_signature NVARCHAR(MAX);
END

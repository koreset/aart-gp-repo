-- Migration: email sending system (Phase 1).
-- Creates per-license SMTP configuration, licensed email templates,
-- a persistent outbox for the worker, and a per-user signature field.

--------------------------------------------------------------------------------
-- email_settings
--------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS email_settings (
    id                       INT AUTO_INCREMENT PRIMARY KEY,
    license_id               VARCHAR(191) NOT NULL,
    host                     VARCHAR(255) NOT NULL,
    port                     INT NOT NULL DEFAULT 587,
    tls_mode                 VARCHAR(16) NOT NULL DEFAULT 'starttls',
    auth_user                VARCHAR(255),
    auth_password_encrypted  TEXT,
    from_address             VARCHAR(255) NOT NULL,
    from_name                VARCHAR(255),
    reply_to                 VARCHAR(255),
    updated_by               VARCHAR(255),
    updated_at               DATETIME NULL,
    UNIQUE KEY uk_email_settings_license (license_id)
);

--------------------------------------------------------------------------------
-- email_templates
--------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS email_templates (
    id                INT AUTO_INCREMENT PRIMARY KEY,
    license_id        VARCHAR(191) NOT NULL,
    code              VARCHAR(128) NOT NULL,
    name              VARCHAR(255),
    description       TEXT,
    subject_template  TEXT,
    body_template     MEDIUMTEXT,
    attachments_spec  TEXT,
    status            VARCHAR(16) NOT NULL DEFAULT 'draft',
    updated_by        VARCHAR(255),
    updated_at        DATETIME NULL,
    UNIQUE KEY uk_email_templates_license_code (license_id, code)
);

--------------------------------------------------------------------------------
-- email_outbox
--------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS email_outbox (
    id                   INT AUTO_INCREMENT PRIMARY KEY,
    license_id           VARCHAR(191) NOT NULL,
    template_code        VARCHAR(128),
    from_address         VARCHAR(255),
    from_name            VARCHAR(255),
    reply_to             VARCHAR(255),
    to_recipients        TEXT,
    cc_recipients        TEXT,
    bcc_recipients       TEXT,
    subject              TEXT,
    body                 MEDIUMTEXT,
    attachments          TEXT,
    status               VARCHAR(16) NOT NULL DEFAULT 'pending',
    attempts             INT NOT NULL DEFAULT 0,
    max_attempts         INT NOT NULL DEFAULT 5,
    last_error           TEXT,
    next_attempt_at      DATETIME NULL,
    scheduled_at         DATETIME NULL,
    sent_at              DATETIME NULL,
    related_object_type  VARCHAR(128),
    related_object_id    VARCHAR(128),
    created_by           VARCHAR(255),
    created_at           DATETIME NULL,
    KEY idx_email_outbox_status_next (status, next_attempt_at),
    KEY idx_email_outbox_license (license_id),
    KEY idx_email_outbox_template (template_code)
);

--------------------------------------------------------------------------------
-- gp_users.email_signature
--------------------------------------------------------------------------------

SET @sql := (
    SELECT IF(
        EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME='gp_users' AND COLUMN_NAME='email_signature'),
        'ALTER TABLE gp_users MODIFY COLUMN email_signature TEXT;',
        'ALTER TABLE gp_users ADD COLUMN email_signature TEXT;'
    )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

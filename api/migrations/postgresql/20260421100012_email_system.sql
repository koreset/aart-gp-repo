-- Migration: email sending system (Phase 1).
-- Creates per-license SMTP configuration, licensed email templates,
-- a persistent outbox for the worker, and a per-user signature field.

--------------------------------------------------------------------------------
-- email_settings
--------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS email_settings (
    id                       SERIAL PRIMARY KEY,
    license_id               VARCHAR(191) NOT NULL,
    host                     VARCHAR(255) NOT NULL,
    port                     INTEGER NOT NULL DEFAULT 587,
    tls_mode                 VARCHAR(16) NOT NULL DEFAULT 'starttls',
    auth_user                VARCHAR(255),
    auth_password_encrypted  TEXT,
    from_address             VARCHAR(255) NOT NULL,
    from_name                VARCHAR(255),
    reply_to                 VARCHAR(255),
    updated_by               VARCHAR(255),
    updated_at               TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_email_settings_license
    ON email_settings (license_id);

--------------------------------------------------------------------------------
-- email_templates
--------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS email_templates (
    id                SERIAL PRIMARY KEY,
    license_id        VARCHAR(191) NOT NULL,
    code              VARCHAR(128) NOT NULL,
    name              VARCHAR(255),
    description       TEXT,
    subject_template  TEXT,
    body_template     TEXT,
    attachments_spec  TEXT,
    status            VARCHAR(16) NOT NULL DEFAULT 'draft',
    updated_by        VARCHAR(255),
    updated_at        TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_email_templates_license_code
    ON email_templates (license_id, code);

--------------------------------------------------------------------------------
-- email_outbox
--------------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS email_outbox (
    id                   SERIAL PRIMARY KEY,
    license_id           VARCHAR(191) NOT NULL,
    template_code        VARCHAR(128),
    from_address         VARCHAR(255),
    from_name            VARCHAR(255),
    reply_to             VARCHAR(255),
    to_recipients        TEXT,
    cc_recipients        TEXT,
    bcc_recipients       TEXT,
    subject              TEXT,
    body                 TEXT,
    attachments          TEXT,
    status               VARCHAR(16) NOT NULL DEFAULT 'pending',
    attempts             INTEGER NOT NULL DEFAULT 0,
    max_attempts         INTEGER NOT NULL DEFAULT 5,
    last_error           TEXT,
    next_attempt_at      TIMESTAMP,
    scheduled_at         TIMESTAMP,
    sent_at              TIMESTAMP,
    related_object_type  VARCHAR(128),
    related_object_id    VARCHAR(128),
    created_by           VARCHAR(255),
    created_at           TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_email_outbox_status_next
    ON email_outbox (status, next_attempt_at);
CREATE INDEX IF NOT EXISTS idx_email_outbox_license
    ON email_outbox (license_id);
CREATE INDEX IF NOT EXISTS idx_email_outbox_template
    ON email_outbox (template_code);

--------------------------------------------------------------------------------
-- gp_users.email_signature
--------------------------------------------------------------------------------

ALTER TABLE gp_users ADD COLUMN IF NOT EXISTS email_signature TEXT;

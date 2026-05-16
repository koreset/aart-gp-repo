-- Generated 2026-05-16T12:23:44+02:00 for dialect mysql

-- Create table: email_settings
CREATE TABLE email_settings (
    id BIGINT AUTO_INCREMENT,
    license_id varchar(191),
    host VARCHAR(255),
    port BIGINT,
    tls_mode varchar(16),
    auth_user VARCHAR(255),
    auth_password_encrypted text,
    from_address VARCHAR(255),
    from_name VARCHAR(255),
    reply_to VARCHAR(255),
    updated_by VARCHAR(255),
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_email_settings_license_id ON email_settings (license_id);


-- Create table: email_outbox
CREATE TABLE email_outbox (
    id BIGINT AUTO_INCREMENT,
    license_id varchar(191),
    template_code varchar(128),
    from_address VARCHAR(255),
    from_name VARCHAR(255),
    reply_to VARCHAR(255),
    to_recipients text,
    cc_recipients text,
    bcc_recipients text,
    subject text,
    body mediumtext,
    attachments text,
    status varchar(16),
    attempts BIGINT,
    max_attempts BIGINT,
    last_error text,
    next_attempt_at DATETIME,
    scheduled_at DATETIME,
    sent_at DATETIME,
    related_object_type VARCHAR(255),
    related_object_id VARCHAR(255),
    created_by VARCHAR(255),
    created_at DATETIME,
    PRIMARY KEY (id)
);

CREATE INDEX idx_email_outbox_license_id ON email_outbox (license_id);
CREATE INDEX idx_email_outbox_template_code ON email_outbox (template_code);
CREATE INDEX idx_email_outbox_status ON email_outbox (status);
CREATE INDEX idx_email_outbox_next_attempt_at ON email_outbox (next_attempt_at);


-- Create table: email_templates
CREATE TABLE email_templates (
    id BIGINT AUTO_INCREMENT,
    license_id varchar(191),
    code varchar(128),
    name VARCHAR(255),
    description text,
    subject_template text,
    body_template mediumtext,
    attachments_spec text,
    status varchar(16),
    updated_by VARCHAR(255),
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_email_templates_license_code ON email_templates (license_id, code);



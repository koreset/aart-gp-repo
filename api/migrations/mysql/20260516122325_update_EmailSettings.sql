-- Generated 2026-05-16T12:23:25+02:00 for dialect mysql

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



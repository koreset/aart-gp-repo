-- Generated 2026-05-16T21:21:10+02:00 for dialect mysql

-- Create table: payment_cutoff_configs
CREATE TABLE payment_cutoff_configs (
    id BIGINT AUTO_INCREMENT,
    license_id VARCHAR(191),
    enabled TINYINT(1) DEFAULT 1,
    cutoff_times VARCHAR(255),
    daily_payment_limit DOUBLE,
    timezone VARCHAR(64),
    updated_by VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_payment_cutoff_configs_license_id ON payment_cutoff_configs (license_id);


-- Create table: payment_cutoff_runs
CREATE TABLE payment_cutoff_runs (
    id BIGINT AUTO_INCREMENT,
    license_id VARCHAR(191),
    scheduled_at DATETIME,
    trigger_type VARCHAR(16),
    status VARCHAR(16),
    error_message VARCHAR(255),
    schedule_id BIGINT,
    claims_count BIGINT,
    total_amount DOUBLE,
    triggered_by VARCHAR(255),
    created_at DATETIME,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_cutoff_run_dedup ON payment_cutoff_runs (license_id, scheduled_at, trigger_type);



-- Generated 2026-05-16T22:07:46+02:00 for dialect mysql

-- Create table: payment_tax_certificates
CREATE TABLE payment_tax_certificates (
    id BIGINT AUTO_INCREMENT,
    schedule_id BIGINT,
    schedule_item_id BIGINT,
    claim_id BIGINT,
    claim_number VARCHAR(255),
    benefit_name VARCHAR(255),
    beneficiary_name VARCHAR(255),
    beneficiary_id_number VARCHAR(255),
    tax_year BIGINT,
    gross_amount DOUBLE,
    tax_withheld DOUBLE,
    certificate_ref VARCHAR(64),
    file_name VARCHAR(255),
    storage_path VARCHAR(255),
    content_type VARCHAR(64),
    generated_by VARCHAR(255),
    generated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE INDEX idx_payment_tax_certificates_schedule_id ON payment_tax_certificates (schedule_id);
CREATE INDEX idx_payment_tax_certificates_schedule_item_id ON payment_tax_certificates (schedule_item_id);
CREATE UNIQUE INDEX idx_tax_cert_item_year ON payment_tax_certificates (schedule_item_id, tax_year);
CREATE INDEX idx_payment_tax_certificates_claim_id ON payment_tax_certificates (claim_id);
CREATE UNIQUE INDEX idx_payment_tax_certificates_certificate_ref ON payment_tax_certificates (certificate_ref);



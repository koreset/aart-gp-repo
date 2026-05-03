-- Migration: create medical_waivers rating table.
-- Keyed by risk_rate_code + gender + age_next_birthday + income_level,
-- holding the medical-waiver sum-at-risk used by the calculation engine.

CREATE TABLE IF NOT EXISTS medical_waivers (
    id                        SERIAL PRIMARY KEY,
    risk_rate_code            VARCHAR(255),
    gender                    VARCHAR(255),
    age_next_birthday         BIGINT,
    income_level              BIGINT,
    medicalwaiver_sum_at_risk DOUBLE PRECISION,
    reinsurance_medicalwaiver_sum_at_risk DOUBLE PRECISION,
    creation_date             TIMESTAMP WITH TIME ZONE,
    created_by                VARCHAR(255)
);

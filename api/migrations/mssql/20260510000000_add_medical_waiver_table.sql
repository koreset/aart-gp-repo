-- Migration: create medical_waivers rating table.
-- Keyed by risk_rate_code + gender + age_next_birthday + income_level,
-- holding the medical-waiver sum-at-risk used by the calculation engine.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'medical_waivers')
BEGIN
    CREATE TABLE medical_waivers (
        id                        BIGINT IDENTITY(1,1) NOT NULL CONSTRAINT pk_medical_waivers PRIMARY KEY,
        risk_rate_code            NVARCHAR(255) NULL,
        gender                    NVARCHAR(255) NULL,
        age_next_birthday         BIGINT NULL,
        income_level              BIGINT NULL,
        medicalwaiver_sum_at_risk FLOAT NULL,
        reinsurance_medicalwaiver_sum_at_risk FLOAT NULL,
        creation_date             DATETIMEOFFSET NULL,
        created_by                NVARCHAR(255) NULL
    );
END;

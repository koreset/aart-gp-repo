-- Migration for struct: GeneralLoading

-- Table: general_loadings

-- Ensure table exists
CREATE TABLE IF NOT EXISTS general_loadings (
    id SERIAL PRIMARY KEY
);

-- Add or modify column for field: RiskRateCode
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS risk_rate_code VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='risk_rate_code') THEN
        ALTER TABLE general_loadings ALTER COLUMN risk_rate_code TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: Age
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS age INTEGER;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='age') THEN
        ALTER TABLE general_loadings ALTER COLUMN age TYPE INTEGER;
    END IF;
END $$;

-- Add or modify column for field: Gender
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gender VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='gender') THEN
        ALTER TABLE general_loadings ALTER COLUMN gender TYPE VARCHAR(255);
    END IF;
END $$;

-- Add or modify column for field: GlaContigencyLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS gla_contigency_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='gla_contigency_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN gla_contigency_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdContigencyLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ptd_contigency_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='ptd_contigency_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN ptd_contigency_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiContigencyLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ci_contigency_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='ci_contigency_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN ci_contigency_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TtdContigencyLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ttd_contigency_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='ttd_contigency_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN ttd_contigency_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PhiContigencyLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS phi_contigency_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='phi_contigency_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN phi_contigency_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: FunContigencyLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS fun_contigency_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='fun_contigency_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN fun_contigency_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: ContinuationLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS continuation_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='continuation_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN continuation_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: TerminalIllnessLoadingRate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS terminal_illness_loading_rate NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='terminal_illness_loading_rate') THEN
        ALTER TABLE general_loadings ALTER COLUMN terminal_illness_loading_rate TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: PtdAcceleratedBenefitDiscount
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ptd_accelerated_benefit_discount NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='ptd_accelerated_benefit_discount') THEN
        ALTER TABLE general_loadings ALTER COLUMN ptd_accelerated_benefit_discount TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CiAcceleratedBenefitDiscount
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS ci_accelerated_benefit_discount NUMERIC(15,5);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='ci_accelerated_benefit_discount') THEN
        ALTER TABLE general_loadings ALTER COLUMN ci_accelerated_benefit_discount TYPE NUMERIC(15,5);
    END IF;
END $$;

-- Add or modify column for field: CreationDate
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS creation_date TIMESTAMP WITH TIME ZONE;
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='creation_date') THEN
        ALTER TABLE general_loadings ALTER COLUMN creation_date TYPE TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Add or modify column for field: CreatedBy
ALTER TABLE general_loadings ADD COLUMN IF NOT EXISTS created_by VARCHAR(255);
-- Update column type if it exists
DO $$
BEGIN
    IF EXISTS(SELECT 1 FROM information_schema.columns WHERE table_name='general_loadings' AND column_name='created_by') THEN
        ALTER TABLE general_loadings ALTER COLUMN created_by TYPE VARCHAR(255);
    END IF;
END $$;


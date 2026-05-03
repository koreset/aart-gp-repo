-- Restructure educator rating model: drop legacy grade-based columns
-- (Grade 0 / 1-7 / 8-12 / Tertiary) and switch to benefit-code +
-- income-level + sum-at-risk. member_rating_results is rebuilt every
-- quote calc, so its drops carry no data-loss risk; educator_rates
-- drops are intentional retirement of the legacy rate model.

-- member_rating_results -------------------------------------------------
ALTER TABLE member_rating_results
    DROP COLUMN IF EXISTS grade0_sum_assured,
    DROP COLUMN IF EXISTS grade17_sum_assured,
    DROP COLUMN IF EXISTS grade812_sum_assured,
    DROP COLUMN IF EXISTS tertiary_sum_assured,
    DROP COLUMN IF EXISTS grade0_risk_rate,
    DROP COLUMN IF EXISTS grade17_risk_rate,
    DROP COLUMN IF EXISTS grade812_risk_rate,
    DROP COLUMN IF EXISTS tertiary_risk_rate;

ALTER TABLE member_rating_results
    ADD COLUMN IF NOT EXISTS educator_sum_at_risk DOUBLE PRECISION NOT NULL DEFAULT 0;

-- educator_rates --------------------------------------------------------
ALTER TABLE educator_rates
    DROP COLUMN IF EXISTS average_child_age,
    DROP COLUMN IF EXISTS average_number_children,
    DROP COLUMN IF EXISTS grade0_risk_rate,
    DROP COLUMN IF EXISTS grade17_risk_rate,
    DROP COLUMN IF EXISTS grade812_risk_rate,
    DROP COLUMN IF EXISTS tertiary_risk_rate;

ALTER TABLE educator_rates
    ADD COLUMN IF NOT EXISTS educator_benefit_code            VARCHAR(255)     NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS income_level                     INTEGER          NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS educator_sum_at_risk             DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS reinsurance_educator_sum_at_risk DOUBLE PRECISION NOT NULL DEFAULT 0;

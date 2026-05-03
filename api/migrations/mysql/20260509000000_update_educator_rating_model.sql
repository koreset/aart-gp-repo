-- Restructure educator rating model: drop legacy grade-based columns
-- (Grade 0 / 1-7 / 8-12 / Tertiary) and switch to benefit-code +
-- income-level + sum-at-risk. member_rating_results is rebuilt every
-- quote calc, so its drops carry no data-loss risk; educator_rates
-- drops are intentional retirement of the legacy rate model.

-- member_rating_results -------------------------------------------------
ALTER TABLE member_rating_results
    DROP COLUMN grade0_sum_assured,
    DROP COLUMN grade17_sum_assured,
    DROP COLUMN grade812_sum_assured,
    DROP COLUMN tertiary_sum_assured,
    DROP COLUMN grade0_risk_rate,
    DROP COLUMN grade17_risk_rate,
    DROP COLUMN grade812_risk_rate,
    DROP COLUMN tertiary_risk_rate,
    ADD COLUMN educator_sum_at_risk DOUBLE NOT NULL DEFAULT 0;

-- educator_rates --------------------------------------------------------
ALTER TABLE educator_rates
    DROP COLUMN average_child_age,
    DROP COLUMN average_number_children,
    DROP COLUMN grade0_risk_rate,
    DROP COLUMN grade17_risk_rate,
    DROP COLUMN grade812_risk_rate,
    DROP COLUMN tertiary_risk_rate,
    ADD COLUMN educator_benefit_code            VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN income_level                     INT          NOT NULL DEFAULT 0,
    ADD COLUMN educator_sum_at_risk             DOUBLE       NOT NULL DEFAULT 0,
    ADD COLUMN reinsurance_educator_sum_at_risk DOUBLE       NOT NULL DEFAULT 0;

-- Restructure educator rating model: drop legacy grade-based columns
-- (Grade 0 / 1-7 / 8-12 / Tertiary) and switch to benefit-code +
-- income-level + sum-at-risk. member_rating_results is rebuilt every
-- quote calc, so its drops carry no data-loss risk; educator_rates
-- drops are intentional retirement of the legacy rate model.

-- member_rating_results -------------------------------------------------
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'grade0_sum_assured')
    ALTER TABLE member_rating_results DROP COLUMN grade0_sum_assured;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'grade17_sum_assured')
    ALTER TABLE member_rating_results DROP COLUMN grade17_sum_assured;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'grade812_sum_assured')
    ALTER TABLE member_rating_results DROP COLUMN grade812_sum_assured;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'tertiary_sum_assured')
    ALTER TABLE member_rating_results DROP COLUMN tertiary_sum_assured;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'grade0_risk_rate')
    ALTER TABLE member_rating_results DROP COLUMN grade0_risk_rate;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'grade17_risk_rate')
    ALTER TABLE member_rating_results DROP COLUMN grade17_risk_rate;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'grade812_risk_rate')
    ALTER TABLE member_rating_results DROP COLUMN grade812_risk_rate;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'tertiary_risk_rate')
    ALTER TABLE member_rating_results DROP COLUMN tertiary_risk_rate;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'educator_sum_at_risk')
BEGIN
    ALTER TABLE member_rating_results
        ADD educator_sum_at_risk FLOAT NOT NULL CONSTRAINT df_mrr_educator_sum_at_risk DEFAULT 0;
END;

-- educator_rates --------------------------------------------------------
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'average_child_age')
    ALTER TABLE educator_rates DROP COLUMN average_child_age;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'average_number_children')
    ALTER TABLE educator_rates DROP COLUMN average_number_children;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'grade0_risk_rate')
    ALTER TABLE educator_rates DROP COLUMN grade0_risk_rate;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'grade17_risk_rate')
    ALTER TABLE educator_rates DROP COLUMN grade17_risk_rate;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'grade812_risk_rate')
    ALTER TABLE educator_rates DROP COLUMN grade812_risk_rate;
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'tertiary_risk_rate')
    ALTER TABLE educator_rates DROP COLUMN tertiary_risk_rate;

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'educator_benefit_code')
BEGIN
    ALTER TABLE educator_rates
        ADD educator_benefit_code NVARCHAR(255) NOT NULL CONSTRAINT df_er_educator_benefit_code DEFAULT '';
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'income_level')
BEGIN
    ALTER TABLE educator_rates
        ADD income_level INT NOT NULL CONSTRAINT df_er_income_level DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'educator_sum_at_risk')
BEGIN
    ALTER TABLE educator_rates
        ADD educator_sum_at_risk FLOAT NOT NULL CONSTRAINT df_er_educator_sum_at_risk DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('educator_rates') AND name = 'reinsurance_educator_sum_at_risk')
BEGIN
    ALTER TABLE educator_rates
        ADD reinsurance_educator_sum_at_risk FLOAT NOT NULL CONSTRAINT df_er_reins_educator_sum_at_risk DEFAULT 0;
END;

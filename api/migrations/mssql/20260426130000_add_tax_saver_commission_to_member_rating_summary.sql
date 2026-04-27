-- Migration: add the TaxSaver per-benefit commission column to
-- member_rating_result_summaries. Parallel to the other Exp*CommissionAmount
-- columns added in 20260421150000 so the scheme-wide commission can be sliced
-- across the TaxSaver benefit, which is part of the TotalAnnualPremium formula.

IF NOT EXISTS (SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'member_rating_result_summaries')
BEGIN
    CREATE TABLE member_rating_result_summaries (
        id INT IDENTITY(1,1) PRIMARY KEY
    );
END;

IF NOT EXISTS(SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = 'member_rating_result_summaries' AND COLUMN_NAME = 'exp_total_tax_saver_annual_commission_amount')
    ALTER TABLE member_rating_result_summaries ADD exp_total_tax_saver_annual_commission_amount FLOAT;

-- Migration: add the TaxSaver per-benefit commission column to
-- member_rating_result_summaries. Parallel to the other Exp*CommissionAmount
-- columns added in 20260421150000 so the scheme-wide commission can be sliced
-- across the TaxSaver benefit, which is part of the TotalAnnualPremium formula.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id SERIAL PRIMARY KEY
);

ALTER TABLE member_rating_result_summaries ADD COLUMN IF NOT EXISTS exp_total_tax_saver_annual_commission_amount NUMERIC(20,6);

-- Migration: add the TaxSaver per-benefit commission column to
-- member_rating_result_summaries. Parallel to the other Exp*CommissionAmount
-- columns added in 20260421150000 so the scheme-wide commission can be sliced
-- across the TaxSaver benefit, which is part of the TotalAnnualPremium formula.

CREATE TABLE IF NOT EXISTS member_rating_result_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY
);

SET @s = (SELECT IF(EXISTS(SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME='member_rating_result_summaries' AND COLUMN_NAME='exp_total_tax_saver_annual_commission_amount' AND TABLE_SCHEMA = DATABASE()), 'ALTER TABLE member_rating_result_summaries MODIFY COLUMN exp_total_tax_saver_annual_commission_amount DOUBLE;', 'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_total_tax_saver_annual_commission_amount DOUBLE;'));
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

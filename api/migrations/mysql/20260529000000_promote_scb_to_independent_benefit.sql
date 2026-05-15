-- Add SCB rate aggregate columns. SCB now uses PhiCappedIncome as its income
-- base at the member level and is included in TotalAnnualPremium /
-- ExpTotalAnnualPremiumExclFuneral. These two columns expose the summed
-- LoadedScbRate / ExpAdjLoadedScbRate so the rate can be traced alongside
-- TotalScbAnnualRiskPremium / ExpAdjTotalScbAnnualRiskPremium.
-- Idempotent on re-runs: each step checks information_schema before acting.

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'total_scb_risk_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN total_scb_risk_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'exp_adj_total_scb_risk_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_scb_risk_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

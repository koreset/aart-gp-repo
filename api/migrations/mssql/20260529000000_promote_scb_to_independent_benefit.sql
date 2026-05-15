-- Add SCB rate aggregate columns. SCB now uses PhiCappedIncome as its income
-- base at the member level and is included in TotalAnnualPremium /
-- ExpTotalAnnualPremiumExclFuneral. These two columns expose the summed
-- LoadedScbRate / ExpAdjLoadedScbRate so the rate can be traced alongside
-- TotalScbAnnualRiskPremium / ExpAdjTotalScbAnnualRiskPremium.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'total_scb_risk_rate')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD total_scb_risk_rate FLOAT NOT NULL CONSTRAINT df_mrrs_total_scb_risk_rate DEFAULT 0;
END;
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_result_summaries') AND name = 'exp_adj_total_scb_risk_rate')
BEGIN
    ALTER TABLE member_rating_result_summaries ADD exp_adj_total_scb_risk_rate FLOAT NOT NULL CONSTRAINT df_mrrs_exp_adj_total_scb_risk_rate DEFAULT 0;
END;

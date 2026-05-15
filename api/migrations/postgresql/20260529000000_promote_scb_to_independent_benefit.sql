-- Add SCB rate aggregate columns. SCB now uses PhiCappedIncome as its income
-- base at the member level and is included in TotalAnnualPremium /
-- ExpTotalAnnualPremiumExclFuneral. These two columns expose the summed
-- LoadedScbRate / ExpAdjLoadedScbRate so the rate can be traced alongside
-- TotalScbAnnualRiskPremium / ExpAdjTotalScbAnnualRiskPremium.
-- Idempotent on re-runs.

ALTER TABLE member_rating_result_summaries
    ADD COLUMN IF NOT EXISTS total_scb_risk_rate         DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS exp_adj_total_scb_risk_rate DOUBLE PRECISION NOT NULL DEFAULT 0;

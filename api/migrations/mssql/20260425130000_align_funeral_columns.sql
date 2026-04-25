-- Migration: align member_rating_results funeral columns with the Go model.
-- The Go model went through a "Cost → RiskPremium / OfficePremium" rename
-- and a "dependant(s) → parent" rename, plus added FinalTotalFuneralOfficePremium.
-- The previous migrations only updated comments, not column names.

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'main_member_funeral_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'main_member_funeral_risk_premium')
    EXEC sp_rename 'member_rating_results.main_member_funeral_cost', 'main_member_funeral_risk_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'spouse_funeral_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'spouse_funeral_risk_premium')
    EXEC sp_rename 'member_rating_results.spouse_funeral_cost', 'spouse_funeral_risk_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'children_funeral_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'child_funeral_risk_premium')
    EXEC sp_rename 'member_rating_results.children_funeral_cost', 'child_funeral_risk_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'dependants_funeral_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'parent_funeral_risk_premium')
    EXEC sp_rename 'member_rating_results.dependants_funeral_cost', 'parent_funeral_risk_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'total_funeral_risk_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'total_funeral_risk_premium')
    EXEC sp_rename 'member_rating_results.total_funeral_risk_cost', 'total_funeral_risk_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_total_funeral_risk_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_total_funeral_risk_premium')
    EXEC sp_rename 'member_rating_results.exp_adj_total_funeral_risk_cost', 'exp_adj_total_funeral_risk_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'total_funeral_office_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'total_funeral_office_premium')
    EXEC sp_rename 'member_rating_results.total_funeral_office_cost', 'total_funeral_office_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_total_funeral_office_cost')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'exp_adj_total_funeral_office_premium')
    EXEC sp_rename 'member_rating_results.exp_adj_total_funeral_office_cost', 'exp_adj_total_funeral_office_premium', 'COLUMN';

IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'dependant_funeral_base_rate')
   AND NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'parent_funeral_base_rate')
    EXEC sp_rename 'member_rating_results.dependant_funeral_base_rate', 'parent_funeral_base_rate', 'COLUMN';

-- Add the new column the model introduced but no migration created.
IF NOT EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'final_total_funeral_office_premium')
    ALTER TABLE member_rating_results ADD final_total_funeral_office_premium DECIMAL(15,5);

-- Drop obsolete dependant_funeral_sum_assured (parent_funeral_sum_assured already exists).
IF EXISTS (SELECT 1 FROM sys.columns WHERE object_id = OBJECT_ID('member_rating_results') AND name = 'dependant_funeral_sum_assured')
    ALTER TABLE member_rating_results DROP COLUMN dependant_funeral_sum_assured;

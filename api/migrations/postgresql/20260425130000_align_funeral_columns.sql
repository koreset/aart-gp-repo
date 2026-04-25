-- Migration: align member_rating_results funeral columns with the Go model.
-- The Go model went through a "Cost → RiskPremium / OfficePremium" rename
-- and a "dependant(s) → parent" rename, plus added FinalTotalFuneralOfficePremium.
-- The previous migrations only updated comments, not column names.

DO $$
BEGIN
    -- main_member_funeral_cost → main_member_funeral_risk_premium
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'main_member_funeral_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'main_member_funeral_risk_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN main_member_funeral_cost TO main_member_funeral_risk_premium;
    END IF;

    -- spouse_funeral_cost → spouse_funeral_risk_premium
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'spouse_funeral_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'spouse_funeral_risk_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN spouse_funeral_cost TO spouse_funeral_risk_premium;
    END IF;

    -- children_funeral_cost → child_funeral_risk_premium  (also "children" → "child")
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'children_funeral_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'child_funeral_risk_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN children_funeral_cost TO child_funeral_risk_premium;
    END IF;

    -- dependants_funeral_cost → parent_funeral_risk_premium  (also "dependants" → "parent")
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'dependants_funeral_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'parent_funeral_risk_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN dependants_funeral_cost TO parent_funeral_risk_premium;
    END IF;

    -- total_funeral_risk_cost → total_funeral_risk_premium
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'total_funeral_risk_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'total_funeral_risk_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN total_funeral_risk_cost TO total_funeral_risk_premium;
    END IF;

    -- exp_adj_total_funeral_risk_cost → exp_adj_total_funeral_risk_premium
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'exp_adj_total_funeral_risk_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'exp_adj_total_funeral_risk_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN exp_adj_total_funeral_risk_cost TO exp_adj_total_funeral_risk_premium;
    END IF;

    -- total_funeral_office_cost → total_funeral_office_premium
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'total_funeral_office_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'total_funeral_office_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN total_funeral_office_cost TO total_funeral_office_premium;
    END IF;

    -- exp_adj_total_funeral_office_cost → exp_adj_total_funeral_office_premium
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'exp_adj_total_funeral_office_cost')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'exp_adj_total_funeral_office_premium') THEN
        ALTER TABLE member_rating_results RENAME COLUMN exp_adj_total_funeral_office_cost TO exp_adj_total_funeral_office_premium;
    END IF;

    -- dependant_funeral_base_rate → parent_funeral_base_rate
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'dependant_funeral_base_rate')
       AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'member_rating_results' AND column_name = 'parent_funeral_base_rate') THEN
        ALTER TABLE member_rating_results RENAME COLUMN dependant_funeral_base_rate TO parent_funeral_base_rate;
    END IF;
END $$;

-- New column the model added but no migration created it yet.
ALTER TABLE member_rating_results ADD COLUMN IF NOT EXISTS final_total_funeral_office_premium NUMERIC(15,5);

-- Drop the obsolete dependant_funeral_sum_assured (parent_funeral_sum_assured already exists per
-- the 20260417093659 migration).
ALTER TABLE member_rating_results DROP COLUMN IF EXISTS dependant_funeral_sum_assured;

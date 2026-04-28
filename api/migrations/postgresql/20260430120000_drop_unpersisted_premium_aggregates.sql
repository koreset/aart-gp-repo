-- Migration: drop derivable summary aggregates from member_rating_result_summaries.
--
-- total_annual_premium_excluding_funeral was unused.
-- proportion_exp_total_premium_excl_funeral_salary is now computed on the fly
-- in the renderer as exp_total_annual_premium_excl_funeral / total_annual_salary.

ALTER TABLE member_rating_result_summaries DROP COLUMN IF EXISTS total_annual_premium_excluding_funeral;
ALTER TABLE member_rating_result_summaries DROP COLUMN IF EXISTS proportion_exp_total_premium_excl_funeral_salary;

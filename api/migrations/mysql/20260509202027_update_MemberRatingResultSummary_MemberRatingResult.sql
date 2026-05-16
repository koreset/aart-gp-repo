-- Generated 2026-05-09T20:20:27+02:00 for dialect mysql

-- Migration for: MemberRatingResultSummary (table: member_rating_result_summaries)

ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_salary_gla DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_salary_ptd DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_salary_ci DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_salary_sgla DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_salary_ttd DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_salary_phi DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN total_annual_salary_fun DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN maximum_gla_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN maximum_ptd_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN severe_illness_maximum_benefit DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN spouse_gla_maximum_benefit DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN ttd_maximum_monthly_benefit DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN phi_maximum_monthly_benefit DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN gla_max_cover_age BIGINT;
ALTER TABLE member_rating_result_summaries ADD COLUMN ptd_max_cover_age BIGINT;
ALTER TABLE member_rating_result_summaries ADD COLUMN ci_max_cover_age BIGINT;
ALTER TABLE member_rating_result_summaries ADD COLUMN ttd_max_cover_age BIGINT;
ALTER TABLE member_rating_result_summaries ADD COLUMN phi_max_cover_age BIGINT;
ALTER TABLE member_rating_result_summaries ADD COLUMN fun_max_cover_age BIGINT;
ALTER TABLE member_rating_result_summaries ADD COLUMN reins_max_gla_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN reins_max_ptd_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN reins_max_ci_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN reins_max_sgla_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN reins_max_ttd_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN reins_max_phi_cover DOUBLE;
ALTER TABLE member_rating_result_summaries ADD COLUMN reins_max_fun_cover DOUBLE;

-- Migration for: MemberRatingResult (table: member_rating_results)

ALTER TABLE member_rating_results ADD COLUMN gla_scheme_size_loading DOUBLE;
ALTER TABLE member_rating_results ADD COLUMN ptd_scheme_size_loading DOUBLE;
ALTER TABLE member_rating_results ADD COLUMN ci_scheme_size_loading DOUBLE;
ALTER TABLE member_rating_results ADD COLUMN ttd_scheme_size_loading DOUBLE;
ALTER TABLE member_rating_results ADD COLUMN phi_scheme_size_loading DOUBLE;
ALTER TABLE member_rating_results ADD COLUMN fun_scheme_size_loading DOUBLE;
ALTER TABLE member_rating_results ADD COLUMN gla_educator_loading DOUBLE;
ALTER TABLE member_rating_results ADD COLUMN ptd_educator_loading DOUBLE;


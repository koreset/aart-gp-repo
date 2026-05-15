-- Add Salary Continuation Benefit (SCB) per-member and summary columns.
-- Mirrors the field structure of the conversion-on-withdrawal slices on
-- MemberRatingResult / MemberRatingResultSummary; SCB is tracked as a
-- reportable slice (NOT added to any group total).
--
-- Source-side rate tables and scheme_categories columns were added in
-- 20260519010000_add_salary_continuation_benefit_tables.sql.
-- Idempotent on re-runs: each step checks information_schema before acting.

-- ── member_rating_results ──────────────────────────────────────────────
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'scb_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN scb_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'base_scb_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN base_scb_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'loaded_scb_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN loaded_scb_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'exp_adj_loaded_scb_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN exp_adj_loaded_scb_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'scb_risk_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN scb_risk_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'exp_adj_scb_risk_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN exp_adj_scb_risk_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'reins_scb_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN reins_scb_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'base_reins_scb_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN base_reins_scb_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'loaded_reins_scb_rate');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN loaded_reins_scb_rate DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_results' AND column_name = 'reins_scb_risk_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN reins_scb_risk_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ── member_rating_result_summaries ─────────────────────────────────────
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'total_scb_annual_risk_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN total_scb_annual_risk_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'exp_adj_total_scb_annual_risk_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_total_scb_annual_risk_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'proportion_scb_risk_premium_salary');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN proportion_scb_risk_premium_salary DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'exp_adj_proportion_scb_risk_premium_salary');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_adj_proportion_scb_risk_premium_salary DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'scb_risk_rate_per_1000_income');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN scb_risk_rate_per_1000_income DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'exp_scb_risk_rate_per_1000_income');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN exp_scb_risk_rate_per_1000_income DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'final_scb_office_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN final_scb_office_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'total_reins_scb_annual_risk_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN total_reins_scb_annual_risk_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE() AND table_name = 'member_rating_result_summaries' AND column_name = 'final_reins_scb_office_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN final_reins_scb_office_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Underwriting platform — consolidated schema for Phases 1 through 7.
--
-- The application's AutoMigrate only runs on a fresh database
-- (api/services/db.go:388). Existing databases receive schema changes
-- exclusively through hand-written migration files like this one. This
-- file adds every column and table introduced by the underwriting
-- platform roadmap so all phase features start working without an
-- AutoMigrate path.
--
-- The file is idempotent: every step checks information_schema before
-- acting, so it's safe to re-run. Apply order within the file matters
-- (columns must exist before indexes that reference them; child tables
-- reference parent IDs but those PK columns exist by the time they're
-- queried).
--
-- Phase 4 note: `member_rating_results.id` is NOT added here. Adding an
-- AUTO_INCREMENT PRIMARY KEY to a table that already holds rows works
-- on small tables but locks for the duration on large ones. Schemes
-- with 10k+ members can take minutes. Add it in a dedicated follow-up
-- migration during a maintenance window, or schedule it for after the
-- next full recalc which rewrites the table from scratch.

-- ═════════════════════════════════════════════════════════════════════
-- Phase 4 — IMMEDIATE: group_pricing_quotes.rating_version
-- Without this column GORM Save() on a quote fails with Error 1054.
-- ═════════════════════════════════════════════════════════════════════

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'group_pricing_quotes'
               AND column_name = 'rating_version');
SET @sql := IF(@col = 0,
  'ALTER TABLE group_pricing_quotes ADD COLUMN rating_version INT NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 1 — member_rating_results: underwriting tier classification
-- ═════════════════════════════════════════════════════════════════════

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'underwriting_tier');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN underwriting_tier INT NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'fcl_excess_ratio');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN fcl_excess_ratio DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 1 — member_rating_result_summaries: per-category tier counts
-- ═════════════════════════════════════════════════════════════════════

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'within_free_cover_limit_count');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN within_free_cover_limit_count INT NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'short_form_underwriting_count');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN short_form_underwriting_count INT NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'full_underwriting_count');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN full_underwriting_count INT NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 4 — member_rating_results: per-benefit UW decision fields +
-- the per-member UW-adjusted office premium snapshot.
-- ═════════════════════════════════════════════════════════════════════

-- GLA
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_gla_loading');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_gla_loading DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_gla_cover_cap');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_gla_cover_cap DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_gla_declined');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_gla_declined TINYINT(1) NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- PTD
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_ptd_loading');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_ptd_loading DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_ptd_cover_cap');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_ptd_cover_cap DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_ptd_declined');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_ptd_declined TINYINT(1) NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- CI
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_ci_loading');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_ci_loading DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_ci_cover_cap');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_ci_cover_cap DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_ci_declined');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_ci_declined TINYINT(1) NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Spouse GLA
SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_sgla_loading');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_sgla_loading DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_sgla_cover_cap');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_sgla_cover_cap DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_sgla_declined');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_sgla_declined TINYINT(1) NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_results'
               AND column_name = 'uw_adjusted_annual_office_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_results ADD COLUMN uw_adjusted_annual_office_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 4 — member_rating_result_summaries: UW-adjusted aggregate
-- premium + per-benefit capped SA + re-rate provenance.
-- ═════════════════════════════════════════════════════════════════════

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'uw_adjusted_total_annual_premium');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN uw_adjusted_total_annual_premium DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'uw_adjusted_total_gla_capped_sum_assured');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN uw_adjusted_total_gla_capped_sum_assured DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'uw_adjusted_total_ptd_capped_sum_assured');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN uw_adjusted_total_ptd_capped_sum_assured DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'uw_adjusted_total_ci_capped_sum_assured');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN uw_adjusted_total_ci_capped_sum_assured DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'uw_adjusted_total_sgla_capped_sum_assured');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN uw_adjusted_total_sgla_capped_sum_assured DOUBLE NOT NULL DEFAULT 0',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'uw_re_rated_at');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN uw_re_rated_at DATETIME NULL',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.columns
             WHERE table_schema = DATABASE()
               AND table_name = 'member_rating_result_summaries'
               AND column_name = 'uw_re_rated_by');
SET @sql := IF(@col = 0,
  'ALTER TABLE member_rating_result_summaries ADD COLUMN uw_re_rated_by VARCHAR(128) NOT NULL DEFAULT ''''',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 1 — new table: underwriting_tier_configs
-- ═════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS underwriting_tier_configs (
  id              INT NOT NULL AUTO_INCREMENT,
  insurer_id      INT NOT NULL DEFAULT 0,
  tier            INT NOT NULL DEFAULT 0,
  lower_multiple  DOUBLE NOT NULL DEFAULT 0,
  upper_multiple  DOUBLE NOT NULL DEFAULT 0,
  active          TINYINT(1) NOT NULL DEFAULT 1,
  creation_date   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by      VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_utc_insurer (insurer_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 2 — case file core: cases, decisions, events, attachments
-- ═════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS underwriting_cases (
  id                          INT NOT NULL AUTO_INCREMENT,
  quote_id                    INT NOT NULL DEFAULT 0,
  member_id_number            VARCHAR(64) NOT NULL DEFAULT '',
  member_name                 VARCHAR(255) NOT NULL DEFAULT '',
  category                    VARCHAR(128) NOT NULL DEFAULT '',
  tier                        INT NOT NULL DEFAULT 0,
  fcl_excess_ratio            DOUBLE NOT NULL DEFAULT 0,
  gla_sum_assured             DOUBLE NOT NULL DEFAULT 0,
  ptd_sum_assured             DOUBLE NOT NULL DEFAULT 0,
  ci_sum_assured              DOUBLE NOT NULL DEFAULT 0,
  spouse_gla_sum_assured      DOUBLE NOT NULL DEFAULT 0,
  free_cover_limit            DOUBLE NOT NULL DEFAULT 0,
  status                      VARCHAR(32) NOT NULL DEFAULT '',
  assigned_underwriter_email  VARCHAR(128) NOT NULL DEFAULT '',
  decided_at                  DATETIME NULL,
  decided_by                  VARCHAR(128) NOT NULL DEFAULT '',
  rule_set_id                 INT NOT NULL DEFAULT 0,
  rule_set_version            INT NOT NULL DEFAULT 0,
  engine_outcome              VARCHAR(16) NOT NULL DEFAULT '',
  engine_loading              DOUBLE NOT NULL DEFAULT 0,
  engine_exclusions           TEXT,
  engine_evaluated_at         DATETIME NULL,
  creation_date               DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by                  VARCHAR(128) NOT NULL DEFAULT '',
  updated_at                  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_uwc_quote (quote_id),
  KEY idx_uwc_status (status),
  KEY idx_uwc_member_id (member_id_number),
  KEY idx_uwc_assignee (assigned_underwriter_email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS underwriting_decisions (
  id                  INT NOT NULL AUTO_INCREMENT,
  case_id             INT NOT NULL DEFAULT 0,
  benefit_type        VARCHAR(16) NOT NULL DEFAULT '',
  outcome             VARCHAR(16) NOT NULL DEFAULT '',
  loading_percent     DOUBLE NOT NULL DEFAULT 0,
  loading_flat_amount DOUBLE NOT NULL DEFAULT 0,
  exclusion_code      VARCHAR(64) NOT NULL DEFAULT '',
  exclusion_text      TEXT,
  cover_cap           DOUBLE NOT NULL DEFAULT 0,
  notes               TEXT,
  creation_date       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by          VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_uwd_case (case_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS underwriting_case_events (
  id            INT NOT NULL AUTO_INCREMENT,
  case_id       INT NOT NULL DEFAULT 0,
  event_type    VARCHAR(64) NOT NULL DEFAULT '',
  actor         VARCHAR(128) NOT NULL DEFAULT '',
  payload       TEXT,
  creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_uwce_case (case_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS underwriting_case_attachments (
  id           INT NOT NULL AUTO_INCREMENT,
  case_id      INT NOT NULL DEFAULT 0,
  kind         VARCHAR(32) NOT NULL DEFAULT '',
  file_name    VARCHAR(255) NOT NULL DEFAULT '',
  content_type VARCHAR(128) NOT NULL DEFAULT '',
  size_bytes   BIGINT NOT NULL DEFAULT 0,
  storage_path VARCHAR(512) NOT NULL DEFAULT '',
  uploaded_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  uploaded_by  VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_uwa_case (case_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 3 — rules engine: rule sets, rules, condition catalogue
-- ═════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS uw_rule_sets (
  id              INT NOT NULL AUTO_INCREMENT,
  name            VARCHAR(128) NOT NULL DEFAULT '',
  version         INT NOT NULL DEFAULT 0,
  effective_from  DATETIME NULL,
  effective_to    DATETIME NULL,
  active          TINYINT(1) NOT NULL DEFAULT 0,
  creation_date   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by      VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_uwrs_name (name),
  KEY idx_uwrs_version (version),
  KEY idx_uwrs_active (active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS uw_rules (
  id              INT NOT NULL AUTO_INCREMENT,
  rule_set_id     INT NOT NULL DEFAULT 0,
  category        VARCHAR(32) NOT NULL DEFAULT '',
  field           VARCHAR(64) NOT NULL DEFAULT '',
  op              VARCHAR(16) NOT NULL DEFAULT '',
  condition_json  TEXT,
  outcome         VARCHAR(16) NOT NULL DEFAULT '',
  loading_percent DOUBLE NOT NULL DEFAULT 0,
  exclusion_code  VARCHAR(64) NOT NULL DEFAULT '',
  priority        INT NOT NULL DEFAULT 0,
  notes           TEXT,
  creation_date   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_uwr_set (rule_set_id),
  KEY idx_uwr_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS uw_condition_codes (
  id            INT NOT NULL AUTO_INCREMENT,
  category      VARCHAR(32) NOT NULL DEFAULT '',
  code          VARCHAR(64) NOT NULL DEFAULT '',
  label         VARCHAR(255) NOT NULL DEFAULT '',
  description   TEXT,
  active        TINYINT(1) NOT NULL DEFAULT 1,
  creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY idx_uwcc_code (code),
  KEY idx_uwcc_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 4 — quote re-rate event audit
-- ═════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS quote_re_rate_events (
  id               INT NOT NULL AUTO_INCREMENT,
  quote_id         INT NOT NULL DEFAULT 0,
  triggered_by     VARCHAR(128) NOT NULL DEFAULT '',
  triggered_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  previous_premium DOUBLE NOT NULL DEFAULT 0,
  new_premium      DOUBLE NOT NULL DEFAULT 0,
  premium_delta    DOUBLE NOT NULL DEFAULT 0,
  reason           TEXT,
  case_id          INT NOT NULL DEFAULT 0,
  rating_version   INT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_qrre_quote (quote_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 5 — disclosure / attestation / consent
-- ═════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS member_disclosures (
  id                      INT NOT NULL AUTO_INCREMENT,
  case_id                 INT NOT NULL DEFAULT 0,
  height                  DOUBLE NOT NULL DEFAULT 0,
  weight                  DOUBLE NOT NULL DEFAULT 0,
  bmi                     DOUBLE NOT NULL DEFAULT 0,
  smoker                  TINYINT(1) NOT NULL DEFAULT 0,
  cigarettes_per_day      INT NOT NULL DEFAULT 0,
  alcohol_units_per_week  DOUBLE NOT NULL DEFAULT 0,
  has_hazardous_hobbies   TINYINT(1) NOT NULL DEFAULT 0,
  hazardous_hobbies       TEXT,
  occupation_risk_answers TEXT,
  disclosed_conditions    TEXT,
  additional_notes        TEXT,
  form_variant            VARCHAR(16) NOT NULL DEFAULT '',
  submitted_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  submitted_via           VARCHAR(32) NOT NULL DEFAULT '',
  submitted_by            VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_md_case (case_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS actively_at_work_attestations (
  id                INT NOT NULL AUTO_INCREMENT,
  case_id           INT NOT NULL DEFAULT 0,
  quote_id          INT NOT NULL DEFAULT 0,
  member_id_number  VARCHAR(64) NOT NULL DEFAULT '',
  member_name       VARCHAR(128) NOT NULL DEFAULT '',
  attested_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  attested_by_name  VARCHAR(255) NOT NULL DEFAULT '',
  attested_by_role  VARCHAR(128) NOT NULL DEFAULT '',
  attested_by_email VARCHAR(128) NOT NULL DEFAULT '',
  ip_address        VARCHAR(64) NOT NULL DEFAULT '',
  user_agent        VARCHAR(255) NOT NULL DEFAULT '',
  signature_hash    VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_aaw_case (case_id),
  KEY idx_aaw_quote (quote_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS consent_records (
  id               INT NOT NULL AUTO_INCREMENT,
  case_id          INT NOT NULL DEFAULT 0,
  quote_id         INT NOT NULL DEFAULT 0,
  consent_type     VARCHAR(32) NOT NULL DEFAULT '',
  granted_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  granted_by_name  VARCHAR(255) NOT NULL DEFAULT '',
  granted_by_email VARCHAR(128) NOT NULL DEFAULT '',
  ip_address       VARCHAR(64) NOT NULL DEFAULT '',
  signature_hash   VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_cr_case (case_id),
  KEY idx_cr_quote (quote_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 6 — vendor adapter audit (requests + webhooks)
-- ═════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS vendor_requests (
  id                   INT NOT NULL AUTO_INCREMENT,
  kind                 VARCHAR(16) NOT NULL DEFAULT '',
  provider             VARCHAR(64) NOT NULL DEFAULT '',
  case_id              INT NOT NULL DEFAULT 0,
  quote_id             INT NOT NULL DEFAULT 0,
  subject              VARCHAR(255) NOT NULL DEFAULT '',
  body                 TEXT,
  metadata_json        TEXT,
  request_payload_hash VARCHAR(64) NOT NULL DEFAULT '',
  external_request_id  VARCHAR(128) NOT NULL DEFAULT '',
  status               VARCHAR(32) NOT NULL DEFAULT '',
  response_json        TEXT,
  cost_cents           INT NOT NULL DEFAULT 0,
  requested_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  requested_by         VARCHAR(128) NOT NULL DEFAULT '',
  completed_at         DATETIME NULL,
  error_message        TEXT,
  PRIMARY KEY (id),
  KEY idx_vr_kind (kind),
  KEY idx_vr_case (case_id),
  KEY idx_vr_quote (quote_id),
  KEY idx_vr_status (status),
  KEY idx_vr_external (external_request_id),
  KEY idx_vr_payload (request_payload_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS vendor_webhooks (
  id                  INT NOT NULL AUTO_INCREMENT,
  provider            VARCHAR(64) NOT NULL DEFAULT '',
  kind                VARCHAR(16) NOT NULL DEFAULT '',
  external_request_id VARCHAR(128) NOT NULL DEFAULT '',
  idempotency_key     VARCHAR(128) NOT NULL DEFAULT '',
  signature_header    VARCHAR(255) NOT NULL DEFAULT '',
  body_sha256         VARCHAR(64) NOT NULL DEFAULT '',
  raw_body            TEXT,
  processed           TINYINT(1) NOT NULL DEFAULT 0,
  processed_at        DATETIME NULL,
  process_error       TEXT,
  received_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY idx_vw_idem (idempotency_key),
  KEY idx_vw_provider (provider),
  KEY idx_vw_kind (kind),
  KEY idx_vw_external (external_request_id),
  KEY idx_vw_processed (processed)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ═════════════════════════════════════════════════════════════════════
-- Phase 7 — takeover + policy admin handoff
-- ═════════════════════════════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS prior_insurer_schedules (
  id                 INT NOT NULL AUTO_INCREMENT,
  quote_id           INT NOT NULL DEFAULT 0,
  insurer_name       VARCHAR(255) NOT NULL DEFAULT '',
  certificate_number VARCHAR(128) NOT NULL DEFAULT '',
  effective_date     DATETIME NULL,
  expiry_date        DATETIME NULL,
  document_path      VARCHAR(512) NOT NULL DEFAULT '',
  uploaded_at        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  uploaded_by        VARCHAR(128) NOT NULL DEFAULT '',
  member_count       INT NOT NULL DEFAULT 0,
  in_force_count     INT NOT NULL DEFAULT 0,
  notes              TEXT,
  PRIMARY KEY (id),
  KEY idx_pis_quote (quote_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS prior_insurer_members (
  id                  INT NOT NULL AUTO_INCREMENT,
  schedule_id         INT NOT NULL DEFAULT 0,
  member_id_number    VARCHAR(64) NOT NULL DEFAULT '',
  member_name         VARCHAR(255) NOT NULL DEFAULT '',
  date_of_birth       DATETIME NULL,
  gla_sum_assured     DOUBLE NOT NULL DEFAULT 0,
  ptd_sum_assured     DOUBLE NOT NULL DEFAULT 0,
  ci_sum_assured      DOUBLE NOT NULL DEFAULT 0,
  prior_loadings      TEXT,
  prior_exclusions    TEXT,
  in_force            TINYINT(1) NOT NULL DEFAULT 0,
  matched_member_name VARCHAR(255) NOT NULL DEFAULT '',
  matched_category    VARCHAR(128) NOT NULL DEFAULT '',
  matched_case_id     INT NOT NULL DEFAULT 0,
  takeover_outcome    VARCHAR(48) NOT NULL DEFAULT '',
  PRIMARY KEY (id),
  KEY idx_pim_schedule (schedule_id),
  KEY idx_pim_member_id (member_id_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS policy_handoff_snapshots (
  id             INT NOT NULL AUTO_INCREMENT,
  quote_id       INT NOT NULL DEFAULT 0,
  scheme_id      INT NOT NULL DEFAULT 0,
  handed_off_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  handed_off_by  VARCHAR(128) NOT NULL DEFAULT '',
  reason         TEXT,
  member_count   INT NOT NULL DEFAULT 0,
  takeover_count INT NOT NULL DEFAULT 0,
  payload        LONGTEXT,
  PRIMARY KEY (id),
  KEY idx_phs_quote (quote_id),
  KEY idx_phs_scheme (scheme_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

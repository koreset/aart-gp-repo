-- Fraud-detection GLM + company-configurable risk classification rules.
--
-- Three tables:
--   * fraud_risk_models       singleton (id=1), stores the active logistic
--                             regression coefficients and training metadata.
--   * fraud_risk_rules        CRUD by an admin; matched rules override the GLM.
--   * fraud_risk_assessments  append-only audit of every fraud-check run.
-- Idempotent on re-runs.

CREATE TABLE IF NOT EXISTS fraud_risk_models (
    id              INT PRIMARY KEY,
    intercept       DOUBLE NOT NULL DEFAULT 0,
    coefficients    JSON NULL,
    trained_at      DATETIME NULL,
    trained_by      VARCHAR(255) NULL,
    sample_size     INT NOT NULL DEFAULT 0,
    positive_count  INT NOT NULL DEFAULT 0,
    auc             DOUBLE NOT NULL DEFAULT 0,
    updated_at      DATETIME NULL
);

-- Seed the singleton row so the GET endpoint always finds an entity.
INSERT INTO fraud_risk_models (id, intercept, sample_size, positive_count, auc)
SELECT 1, 0, 0, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM fraud_risk_models WHERE id = 1);

CREATE TABLE IF NOT EXISTS fraud_risk_rules (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    description  TEXT NULL,
    conditions   JSON NOT NULL,
    risk_level   VARCHAR(32) NOT NULL,
    priority     INT NOT NULL DEFAULT 50,
    enabled      TINYINT(1) NOT NULL DEFAULT 1,
    updated_by   VARCHAR(255) NULL,
    updated_at   DATETIME NULL,
    created_at   DATETIME NULL
);

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'fraud_risk_rules' AND index_name = 'idx_frr_enabled_priority');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_frr_enabled_priority ON fraud_risk_rules(enabled, priority)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS fraud_risk_assessments (
    id                INT AUTO_INCREMENT PRIMARY KEY,
    claim_id          INT NOT NULL,
    glm_score         DOUBLE NULL,
    glm_band          VARCHAR(32) NULL,
    matched_rule_id   INT NULL,
    final_risk_level  VARCHAR(32) NULL,
    features          JSON NULL,
    rationale         VARCHAR(500) NULL,
    computed_at       DATETIME NULL,
    computed_by       VARCHAR(255) NULL
);

SET @idx := (SELECT COUNT(*) FROM information_schema.statistics
             WHERE table_schema = DATABASE() AND table_name = 'fraud_risk_assessments' AND index_name = 'idx_fra_claim');
SET @sql := IF(@idx = 0, 'CREATE INDEX idx_fra_claim ON fraud_risk_assessments(claim_id)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

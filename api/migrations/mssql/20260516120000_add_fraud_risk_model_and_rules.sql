-- Fraud-detection GLM + company-configurable risk classification rules.
--
-- Three tables:
--   * fraud_risk_models       singleton (id=1), stores the active logistic
--                             regression coefficients and training metadata.
--   * fraud_risk_rules        CRUD by an admin; matched rules override the GLM.
--   * fraud_risk_assessments  append-only audit of every fraud-check run.
-- Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'fraud_risk_models')
BEGIN
    CREATE TABLE fraud_risk_models (
        id              INT PRIMARY KEY,
        intercept       FLOAT NOT NULL DEFAULT 0,
        coefficients    NVARCHAR(MAX) NULL,
        trained_at      DATETIME NULL,
        trained_by      NVARCHAR(255) NULL,
        sample_size     INT NOT NULL DEFAULT 0,
        positive_count  INT NOT NULL DEFAULT 0,
        auc             FLOAT NOT NULL DEFAULT 0,
        updated_at      DATETIME NULL
    );
END;

-- Seed the singleton row.
IF NOT EXISTS (SELECT 1 FROM fraud_risk_models WHERE id = 1)
    INSERT INTO fraud_risk_models (id, intercept, sample_size, positive_count, auc) VALUES (1, 0, 0, 0, 0);

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'fraud_risk_rules')
BEGIN
    CREATE TABLE fraud_risk_rules (
        id           INT IDENTITY(1,1) PRIMARY KEY,
        name         NVARCHAR(255) NOT NULL,
        description  NVARCHAR(MAX) NULL,
        conditions   NVARCHAR(MAX) NOT NULL,
        risk_level   NVARCHAR(32) NOT NULL,
        priority     INT NOT NULL DEFAULT 50,
        enabled      BIT NOT NULL DEFAULT 1,
        updated_by   NVARCHAR(255) NULL,
        updated_at   DATETIME NULL,
        created_at   DATETIME NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_frr_enabled_priority' AND object_id = OBJECT_ID('fraud_risk_rules'))
    CREATE INDEX idx_frr_enabled_priority ON fraud_risk_rules(enabled, priority);

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'fraud_risk_assessments')
BEGIN
    CREATE TABLE fraud_risk_assessments (
        id                INT IDENTITY(1,1) PRIMARY KEY,
        claim_id          INT NOT NULL,
        glm_score         FLOAT NULL,
        glm_band          NVARCHAR(32) NULL,
        matched_rule_id   INT NULL,
        final_risk_level  NVARCHAR(32) NULL,
        features          NVARCHAR(MAX) NULL,
        rationale         NVARCHAR(500) NULL,
        computed_at       DATETIME NULL,
        computed_by       NVARCHAR(255) NULL
    );
END;

IF NOT EXISTS (SELECT 1 FROM sys.indexes WHERE name = 'idx_fra_claim' AND object_id = OBJECT_ID('fraud_risk_assessments'))
    CREATE INDEX idx_fra_claim ON fraud_risk_assessments(claim_id);

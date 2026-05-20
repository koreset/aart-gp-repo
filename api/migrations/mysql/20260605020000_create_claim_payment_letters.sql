-- Records every generation of a claim payment confirmation letter. The bank
-- and amount fields are snapshotted so historical letters remain stable even
-- if the underlying claim is later edited. Idempotent on re-runs.

CREATE TABLE IF NOT EXISTS claim_payment_letters (
  id                   BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  claim_id             BIGINT NOT NULL,
  version              INT NOT NULL,
  format               VARCHAR(8) NOT NULL,
  filename             VARCHAR(512) NULL,
  size_bytes           BIGINT NULL,
  letter_reference     VARCHAR(64) NOT NULL,
  payment_amount       DECIMAL(18,2) NULL,
  paid_at              DATETIME NULL,
  bank_name            VARCHAR(255) NULL,
  bank_account_number  VARCHAR(255) NULL,
  account_holder_name  VARCHAR(255) NULL,
  settings_snapshot    TEXT NULL,
  generated_by         VARCHAR(255) NULL,
  generated_at         DATETIME NULL,
  INDEX idx_claim_payment_letters_claim_id (claim_id),
  UNIQUE INDEX idx_claim_payment_letters_reference (letter_reference)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

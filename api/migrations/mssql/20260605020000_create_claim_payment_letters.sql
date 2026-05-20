-- Records every generation of a claim payment confirmation letter. The bank
-- and amount fields are snapshotted so historical letters remain stable even
-- if the underlying claim is later edited. Idempotent on re-runs.

IF NOT EXISTS (SELECT 1 FROM sys.tables WHERE name = 'claim_payment_letters')
BEGIN
  CREATE TABLE claim_payment_letters (
    id                   BIGINT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    claim_id             BIGINT NOT NULL,
    version              INT NOT NULL,
    format               NVARCHAR(8) NOT NULL,
    filename             NVARCHAR(512) NULL,
    size_bytes           BIGINT NULL,
    letter_reference     NVARCHAR(64) NOT NULL,
    payment_amount       DECIMAL(18,2) NULL,
    paid_at              DATETIME NULL,
    bank_name            NVARCHAR(255) NULL,
    bank_account_number  NVARCHAR(255) NULL,
    account_holder_name  NVARCHAR(255) NULL,
    settings_snapshot    NVARCHAR(MAX) NULL,
    generated_by         NVARCHAR(255) NULL,
    generated_at         DATETIME NULL
  );
  CREATE INDEX idx_claim_payment_letters_claim_id ON claim_payment_letters (claim_id);
  CREATE UNIQUE INDEX idx_claim_payment_letters_reference ON claim_payment_letters (letter_reference);
END

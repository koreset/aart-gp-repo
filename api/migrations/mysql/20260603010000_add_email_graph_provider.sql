-- Add a delivery-provider discriminator and Microsoft Graph (Office 365)
-- credential columns to email_settings. Existing rows default to provider
-- 'smtp' so current SMTP/relay configurations keep working unchanged. The
-- Graph client secret is stored encrypted (AES-GCM, see api/services/crypto),
-- mirroring auth_password_encrypted.

ALTER TABLE email_settings
    ADD COLUMN provider                      VARCHAR(32)  NOT NULL DEFAULT 'smtp',
    ADD COLUMN graph_tenant_id               VARCHAR(191) NOT NULL DEFAULT '',
    ADD COLUMN graph_client_id               VARCHAR(191) NOT NULL DEFAULT '',
    ADD COLUMN graph_client_secret_encrypted TEXT NULL;

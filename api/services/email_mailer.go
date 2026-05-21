package services

import (
	"fmt"

	"api/config"
	"api/models"
	"api/services/crypto"
	"api/services/email"
)

// BuildMailer selects and constructs the delivery transport for a license's
// email settings, decrypting whichever secret that provider needs. It is the
// single place the outbox worker branches on Provider, keeping the email
// package transport-pure (crypto/config live in the services layer).
//
// For the Graph provider, the client id/secret fall back to the global
// config.GraphClientID/Secret when the per-license values are blank — this
// supports a single central multi-tenant Azure app where each license only
// records its tenant id.
func BuildMailer(s models.EmailSettings) (email.Mailer, error) {
	switch s.Provider {
	case models.EmailProviderGraph:
		clientID := s.GraphClientId
		clientSecret := ""
		if s.GraphClientSecretEncrypted != "" {
			dec, err := crypto.Decrypt(s.GraphClientSecretEncrypted)
			if err != nil {
				return nil, fmt.Errorf("decrypt graph client secret: %w", err)
			}
			clientSecret = dec
		}
		if clientID == "" {
			clientID = config.GraphClientID
		}
		if clientSecret == "" {
			clientSecret = config.GraphClientSecret
		}
		return email.NewGraphMailer(s, s.GraphTenantId, clientID, clientSecret), nil
	case models.EmailProviderSMTP, "":
		password, err := crypto.Decrypt(s.AuthPasswordEncrypted)
		if err != nil {
			return nil, fmt.Errorf("decrypt SMTP password: %w", err)
		}
		return email.NewSMTPMailer(s, password), nil
	default:
		return nil, fmt.Errorf("unsupported email provider %q", s.Provider)
	}
}

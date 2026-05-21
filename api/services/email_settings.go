package services

import (
	"errors"
	"fmt"
	"time"

	"api/config"
	"api/models"
	"api/services/crypto"

	"gorm.io/gorm"
)

// GetEmailSettings returns the EmailSettings row for the given license, or
// gorm.ErrRecordNotFound if none has been configured.
func GetEmailSettings(licenseId string) (models.EmailSettings, error) {
	var s models.EmailSettings
	err := DB.Where("license_id = ?", licenseId).First(&s).Error
	return s, err
}

// SaveEmailSettingsInput is the sanitised input accepted by SaveEmailSettings.
// AuthPassword and GraphClientSecret are plaintext — encrypted before write.
// When either is empty and a row already exists, the existing encrypted secret
// is preserved. Provider selects which set of fields applies.
type SaveEmailSettingsInput struct {
	Provider          string `json:"provider"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	TlsMode           string `json:"tls_mode"`
	AuthUser          string `json:"auth_user"`
	AuthPassword      string `json:"auth_password"`
	GraphTenantId     string `json:"graph_tenant_id"`
	GraphClientId     string `json:"graph_client_id"`
	GraphClientSecret string `json:"graph_client_secret"`
	FromAddress       string `json:"from_address"`
	FromName          string `json:"from_name"`
	ReplyTo           string `json:"reply_to"`
}

// SaveEmailSettings upserts the per-license configuration. Password is
// encrypted at rest; omitting AuthPassword on an update preserves the existing
// value rather than blanking it.
func SaveEmailSettings(licenseId string, in SaveEmailSettingsInput, user models.AppUser) (models.EmailSettings, error) {
	if licenseId == "" {
		return models.EmailSettings{}, errors.New("license_id is required")
	}
	if in.Provider == "" {
		in.Provider = models.EmailProviderSMTP
	}
	if in.FromAddress == "" {
		return models.EmailSettings{}, errors.New("from_address is required")
	}

	switch in.Provider {
	case models.EmailProviderSMTP:
		if in.Host == "" {
			return models.EmailSettings{}, errors.New("host is required for the smtp provider")
		}
		if in.Port == 0 {
			in.Port = 587
		}
		if in.TlsMode == "" {
			in.TlsMode = models.EmailTLSModeSTARTTLS
		}
		switch in.TlsMode {
		case models.EmailTLSModeNone, models.EmailTLSModeSTARTTLS, models.EmailTLSModeTLS:
		default:
			return models.EmailSettings{}, fmt.Errorf("unsupported tls_mode %q", in.TlsMode)
		}
	case models.EmailProviderGraph:
		if in.GraphTenantId == "" {
			return models.EmailSettings{}, errors.New("graph_tenant_id is required for the microsoft_graph provider")
		}
		if in.GraphClientId == "" && config.GraphClientID == "" {
			return models.EmailSettings{}, errors.New("graph_client_id is required for the microsoft_graph provider")
		}
	default:
		return models.EmailSettings{}, fmt.Errorf("unsupported provider %q", in.Provider)
	}

	existing, err := GetEmailSettings(licenseId)
	isUpdate := err == nil
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.EmailSettings{}, err
	}

	encryptedPassword := existing.AuthPasswordEncrypted
	if in.AuthPassword != "" {
		enc, err := crypto.Encrypt(in.AuthPassword)
		if err != nil {
			return models.EmailSettings{}, fmt.Errorf("encrypt password: %w", err)
		}
		encryptedPassword = enc
	}

	encryptedGraphSecret := existing.GraphClientSecretEncrypted
	if in.GraphClientSecret != "" {
		enc, err := crypto.Encrypt(in.GraphClientSecret)
		if err != nil {
			return models.EmailSettings{}, fmt.Errorf("encrypt graph client secret: %w", err)
		}
		encryptedGraphSecret = enc
	}

	// For Graph, ensure a usable client secret will exist after the write —
	// either supplied now, preserved from a prior save, or via the global
	// central-app fallback.
	if in.Provider == models.EmailProviderGraph &&
		encryptedGraphSecret == "" && config.GraphClientSecret == "" {
		return models.EmailSettings{}, errors.New("graph_client_secret is required for the microsoft_graph provider")
	}

	row := models.EmailSettings{
		ID:                         existing.ID,
		LicenseId:                  licenseId,
		Provider:                   in.Provider,
		Host:                       in.Host,
		Port:                       in.Port,
		TlsMode:                    in.TlsMode,
		AuthUser:                   in.AuthUser,
		AuthPasswordEncrypted:      encryptedPassword,
		GraphTenantId:              in.GraphTenantId,
		GraphClientId:              in.GraphClientId,
		GraphClientSecretEncrypted: encryptedGraphSecret,
		FromAddress:                in.FromAddress,
		FromName:                   in.FromName,
		ReplyTo:                    in.ReplyTo,
		UpdatedBy:                  user.UserName,
		UpdatedAt:                  time.Now(),
	}
	if isUpdate {
		if err := DB.Save(&row).Error; err != nil {
			return row, err
		}
	} else {
		if err := DB.Create(&row).Error; err != nil {
			return row, err
		}
	}
	return row, nil
}

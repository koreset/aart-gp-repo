package services

import (
	"errors"
	"fmt"
	"time"

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
// AuthPassword is plaintext — encrypted before write. When empty and a row
// already exists, the existing encrypted password is preserved.
type SaveEmailSettingsInput struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	TlsMode      string `json:"tls_mode"`
	AuthUser     string `json:"auth_user"`
	AuthPassword string `json:"auth_password"`
	FromAddress  string `json:"from_address"`
	FromName     string `json:"from_name"`
	ReplyTo      string `json:"reply_to"`
}

// SaveEmailSettings upserts the per-license configuration. Password is
// encrypted at rest; omitting AuthPassword on an update preserves the existing
// value rather than blanking it.
func SaveEmailSettings(licenseId string, in SaveEmailSettingsInput, user models.AppUser) (models.EmailSettings, error) {
	if licenseId == "" {
		return models.EmailSettings{}, errors.New("license_id is required")
	}
	if in.Host == "" || in.FromAddress == "" {
		return models.EmailSettings{}, errors.New("host and from_address are required")
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

	row := models.EmailSettings{
		ID:                    existing.ID,
		LicenseId:             licenseId,
		Host:                  in.Host,
		Port:                  in.Port,
		TlsMode:               in.TlsMode,
		AuthUser:              in.AuthUser,
		AuthPasswordEncrypted: encryptedPassword,
		FromAddress:           in.FromAddress,
		FromName:              in.FromName,
		ReplyTo:               in.ReplyTo,
		UpdatedBy:             user.UserName,
		UpdatedAt:             time.Now(),
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

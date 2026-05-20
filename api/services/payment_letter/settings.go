package payment_letter

import (
	"errors"
	"time"

	"api/models"

	"gorm.io/gorm"
)

// LoadSettings exposes the singleton settings row so callers (controller +
// generator) share one resolver. On first access it lazily creates an empty
// row at id=1.
func LoadSettings(db *gorm.DB) (models.PaymentLetterSetting, error) {
	return loadSettings(db)
}

// UpsertSettings updates the singleton in place. Logo / signature bytes are
// only overwritten when the caller explicitly sets them via the *Provided
// flags so the admin UI can update non-image fields without resending the
// (potentially large) blobs each time.
type UpsertSettingsInput struct {
	CompanyName    string
	AddressLine1   string
	AddressLine2   string
	AddressLine3   string
	City           string
	PostalCode     string
	Country        string
	Phone          string
	Email          string
	Website        string
	SignatoryName  string
	SignatoryTitle string
	LetterIntro    string
	LetterClosing  string
	UpdatedBy      string

	LogoProvided      bool
	Logo              []byte
	LogoMimeType      string
	SignatureProvided bool
	Signature         []byte
	SignatureMimeType string
}

// UpsertSettings applies the input to the singleton row, returning the saved
// row.
func UpsertSettings(db *gorm.DB, in UpsertSettingsInput) (models.PaymentLetterSetting, error) {
	var s models.PaymentLetterSetting
	if err := db.First(&s, 1).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return s, err
		}
		s = models.PaymentLetterSetting{ID: 1}
	}
	s.CompanyName = in.CompanyName
	s.AddressLine1 = in.AddressLine1
	s.AddressLine2 = in.AddressLine2
	s.AddressLine3 = in.AddressLine3
	s.City = in.City
	s.PostalCode = in.PostalCode
	s.Country = in.Country
	s.Phone = in.Phone
	s.Email = in.Email
	s.Website = in.Website
	s.SignatoryName = in.SignatoryName
	s.SignatoryTitle = in.SignatoryTitle
	s.LetterIntroTemplate = in.LetterIntro
	s.LetterClosingTemplate = in.LetterClosing
	s.UpdatedBy = in.UpdatedBy
	s.UpdatedAt = time.Now()
	if in.LogoProvided {
		s.Logo = in.Logo
		s.LogoMimeType = in.LogoMimeType
	}
	if in.SignatureProvided {
		s.Signature = in.Signature
		s.SignatureMimeType = in.SignatureMimeType
	}
	if s.ID == 0 {
		s.ID = 1
		if err := db.Create(&s).Error; err != nil {
			return s, err
		}
		return s, nil
	}
	if err := db.Save(&s).Error; err != nil {
		return s, err
	}
	return s, nil
}

package payment_letter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"api/models"
)

// BuildLetterDocx renders the payment confirmation letter for a single claim
// into a complete DOCX byte stream. The caller is responsible for resolving
// the LetterInput (claim, paidAt, settings, letter reference).
func BuildLetterDocx(in LetterInput) (filename string, data []byte, err error) {
	pkg := &docxPackage{
		body:          buildBodyXML(in),
		headerXML:     buildHeaderXML(in.Settings),
		footerXML:     buildFooterXML(),
		logo:          in.Settings.Logo,
		logoMIME:      in.Settings.LogoMimeType,
		signature:     in.Settings.Signature,
		signatureMIME: in.Settings.SignatureMimeType,
	}
	data, err = pkg.Build()
	if err != nil {
		return "", nil, fmt.Errorf("build docx: %w", err)
	}
	filename = buildFilename(in.Claim, in.PaidAt)
	return filename, data, nil
}

// buildFilename returns a stable, filesystem-safe filename for the letter.
// Example: PaymentConfirmation_CLM-12345_JaneDoe_2026-05-20.docx
func buildFilename(c models.GroupSchemeClaim, paidAt time.Time) string {
	parts := []string{"PaymentConfirmation"}
	if c.ClaimNumber != "" {
		parts = append(parts, sanitiseFilenamePart(c.ClaimNumber))
	}
	if c.ClaimantName != "" {
		parts = append(parts, sanitiseFilenamePart(strings.ReplaceAll(c.ClaimantName, " ", "")))
	}
	if !paidAt.IsZero() {
		parts = append(parts, paidAt.Format("2006-01-02"))
	}
	return strings.Join(parts, "_") + ".docx"
}

var filenameInvalid = regexp.MustCompile(`[/\\:*?"<>|]`)

func sanitiseFilenamePart(s string) string {
	return filenameInvalid.ReplaceAllString(s, "_")
}

// SerialiseSettingsSnapshot returns a compact JSON of the settings string
// fields actually rendered into the letter. Binary logo/signature blobs are
// excluded — the goal is a debuggable audit record, not a re-renderable copy.
func SerialiseSettingsSnapshot(s models.PaymentLetterSetting) string {
	snap := struct {
		CompanyName    string `json:"company_name,omitempty"`
		AddressLine1   string `json:"address_line1,omitempty"`
		AddressLine2   string `json:"address_line2,omitempty"`
		AddressLine3   string `json:"address_line3,omitempty"`
		City           string `json:"city,omitempty"`
		PostalCode     string `json:"postal_code,omitempty"`
		Country        string `json:"country,omitempty"`
		Phone          string `json:"phone,omitempty"`
		Email          string `json:"email,omitempty"`
		Website        string `json:"website,omitempty"`
		SignatoryName  string `json:"signatory_name,omitempty"`
		SignatoryTitle string `json:"signatory_title,omitempty"`
		HasLogo        bool   `json:"has_logo,omitempty"`
		HasSignature   bool   `json:"has_signature,omitempty"`
	}{
		CompanyName:    s.CompanyName,
		AddressLine1:   s.AddressLine1,
		AddressLine2:   s.AddressLine2,
		AddressLine3:   s.AddressLine3,
		City:           s.City,
		PostalCode:     s.PostalCode,
		Country:        s.Country,
		Phone:          s.Phone,
		Email:          s.Email,
		Website:        s.Website,
		SignatoryName:  s.SignatoryName,
		SignatoryTitle: s.SignatoryTitle,
		HasLogo:        len(s.Logo) > 0,
		HasSignature:   len(s.Signature) > 0,
	}
	b, _ := json.Marshal(snap)
	return string(b)
}

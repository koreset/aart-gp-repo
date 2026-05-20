package payment_letter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"api/models"

	"gorm.io/gorm"
)

// ErrChannelNotImplemented is returned for channels that have schemas but no
// outbound provider yet (currently WhatsApp; SMS is scaffolded but uses the
// log-only provider in dev).
var ErrChannelNotImplemented = errors.New("delivery channel not implemented yet")

// ErrRecipientMissing is returned when the user asks to send by a channel but
// the claim has no recipient on file (no email, no phone, etc.).
var ErrRecipientMissing = errors.New("no recipient on file for this channel")

// PrepareLetterAttachment re-renders the letter from its claim + current
// settings, writes it to a temporary file the email outbox worker can attach,
// and returns the file path. The temp file is owned by the OS temp dir
// lifecycle — the worker only needs to read it once before it's gone.
func PrepareLetterAttachment(db *gorm.DB, letter models.ClaimPaymentLetter, claim models.GroupSchemeClaim) (path string, err error) {
	settings, err := loadSettings(db)
	if err != nil {
		return "", err
	}
	in := LetterInput{
		Claim:     claim,
		PaidAt:    letter.PaidAt,
		Settings:  settings,
		LetterRef: letter.LetterReference,
	}
	_, data, err := BuildLetterDocx(in)
	if err != nil {
		return "", err
	}

	dir, err := os.MkdirTemp("", "payment_letter_*")
	if err != nil {
		return "", err
	}
	filename := letter.Filename
	if filename == "" {
		filename = buildFilename(claim, letter.PaidAt)
	}
	// Always stage DOCX; the email worker doesn't need the PDF for the
	// attachment since the letter content is identical and the DOCX
	// re-rendering is cheap.
	if filepath.Ext(filename) == "" {
		filename = filename + ".docx"
	} else if !strings.EqualFold(filepath.Ext(filename), ".docx") {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".docx"
	}
	full := filepath.Join(dir, filename)
	if err := os.WriteFile(full, data, 0o600); err != nil {
		return "", err
	}
	return full, nil
}

// DefaultRecipient returns the claim contact field that matches the channel.
// Returns "" when the claim has nothing on file.
func DefaultRecipient(c models.GroupSchemeClaim, channel string) string {
	switch channel {
	case models.PaymentLetterChannelEmail:
		return strings.TrimSpace(c.ClaimantEmail)
	case models.PaymentLetterChannelSMS, models.PaymentLetterChannelWhatsApp:
		return strings.TrimSpace(c.ClaimantContactNumber)
	default:
		return ""
	}
}

// RecordDelivery persists a delivery attempt. Used by the controller after
// queueing the underlying transport (e.g. email_outbox) — keeps the audit row
// consistent regardless of which channel was used.
func RecordDelivery(db *gorm.DB, d models.ClaimPaymentLetterDelivery) (models.ClaimPaymentLetterDelivery, error) {
	if d.Status == "" {
		d.Status = models.PaymentLetterDeliveryPending
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	if err := db.Create(&d).Error; err != nil {
		return d, err
	}
	return d, nil
}

// LogCommunication mirrors a delivery attempt onto the existing claim
// communication timeline so users browsing the claim see the same event.
func LogCommunication(db *gorm.DB, claimID int, method, message, createdBy string) {
	if claimID == 0 || message == "" {
		return
	}
	_ = db.Create(&models.GroupSchemeClaimCommunication{
		ClaimID:   claimID,
		Method:    method,
		Message:   message,
		CreatedBy: createdBy,
	}).Error
}

// CheckChannelImplemented returns an error for channels that are not yet
// connected to an outbound provider. WhatsApp is the headline gap; SMS is
// implemented but rejects when no provider is configured.
func CheckChannelImplemented(channel string) error {
	switch channel {
	case models.PaymentLetterChannelEmail, models.PaymentLetterChannelManual,
		models.PaymentLetterChannelSMS:
		return nil
	case models.PaymentLetterChannelWhatsApp:
		return ErrChannelNotImplemented
	default:
		return fmt.Errorf("unknown channel %q", channel)
	}
}

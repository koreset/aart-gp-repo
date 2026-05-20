package models

import "time"

// PaymentLetterSetting is a singleton table (one row, ID=1) holding the
// letterhead, branding, and signatory details used when generating claim
// payment confirmation letters. Logo and signature image are stored inline so
// the API remains the single source of truth across deployments — no external
// asset storage required.
type PaymentLetterSetting struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	CompanyName   string `json:"company_name"`
	AddressLine1  string `json:"address_line1"`
	AddressLine2  string `json:"address_line2"`
	AddressLine3  string `json:"address_line3"`
	City          string `json:"city"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Website       string `json:"website"`
	Logo          []byte `json:"-"`
	LogoMimeType  string `json:"logo_mime_type" gorm:"size:64"`
	SignatoryName  string `json:"signatory_name"`
	SignatoryTitle string `json:"signatory_title"`
	Signature      []byte `json:"-"`
	SignatureMimeType string `json:"signature_mime_type" gorm:"size:64"`
	LetterIntroTemplate   string `json:"letter_intro_template" gorm:"type:text"`
	LetterClosingTemplate string `json:"letter_closing_template" gorm:"type:text"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy string    `json:"updated_by"`
}

// ClaimPaymentLetter records each generation of a confirmation letter for a
// paid claim. One row per generation event: re-downloading the same version
// does not produce a new row, but switching format (PDF↔DOCX) or regenerating
// after a settings/data change does. PaidAt and bank fields are snapshotted so
// historical letters stay stable even if the claim is later edited.
type ClaimPaymentLetter struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ClaimID           int       `json:"claim_id" gorm:"index;not null"`
	Version           int       `json:"version"`
	Format            string    `json:"format" gorm:"type:varchar(8)"` // pdf | docx
	Filename          string    `json:"filename" gorm:"type:varchar(512)"`
	SizeBytes         int64     `json:"size_bytes"`
	LetterReference   string    `json:"letter_reference" gorm:"type:varchar(64);uniqueIndex"`
	PaymentAmount     float64   `json:"payment_amount"`
	PaidAt            time.Time `json:"paid_at"`
	BankName          string    `json:"bank_name"`
	BankAccountNumber string    `json:"bank_account_number"`
	AccountHolderName string    `json:"account_holder_name"`
	// SettingsSnapshot stores the company-name / signatory string fields used
	// when rendering so we can re-print the same letter without a logo blob.
	SettingsSnapshot string                       `json:"settings_snapshot" gorm:"type:text"`
	GeneratedBy      string                       `json:"generated_by"`
	GeneratedAt      time.Time                    `json:"generated_at" gorm:"autoCreateTime"`
	Deliveries       []ClaimPaymentLetterDelivery `json:"deliveries" gorm:"foreignKey:LetterID;references:ID"`
}

// ClaimPaymentLetterDelivery records every attempt to send a letter to the
// claimant via a particular channel. Channel is intentionally open-ended so
// SMS / WhatsApp can be added without a schema change. When the channel is
// email, OutboxID points at the email_outbox row queued for the worker.
type ClaimPaymentLetterDelivery struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	LetterID    int        `json:"letter_id" gorm:"index;not null"`
	Channel     string     `json:"channel" gorm:"type:varchar(16)"` // email | sms | whatsapp | manual
	Recipient   string     `json:"recipient"`
	Status      string     `json:"status" gorm:"type:varchar(16)"` // pending | sent | failed
	ProviderRef string     `json:"provider_ref"`
	OutboxID    *int       `json:"outbox_id"`
	Error       string     `json:"error" gorm:"type:text"`
	SentBy      string     `json:"sent_by"`
	SentAt      *time.Time `json:"sent_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// Delivery channel constants.
const (
	PaymentLetterChannelEmail    = "email"
	PaymentLetterChannelSMS      = "sms"
	PaymentLetterChannelWhatsApp = "whatsapp"
	PaymentLetterChannelManual   = "manual"
)

// Delivery status constants.
const (
	PaymentLetterDeliveryPending = "pending"
	PaymentLetterDeliverySent    = "sent"
	PaymentLetterDeliveryFailed  = "failed"
)

// Format constants.
const (
	PaymentLetterFormatDocx = "docx"
	PaymentLetterFormatPDF  = "pdf"
)

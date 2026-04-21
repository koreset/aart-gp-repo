package models

import "time"

// EmailSettings holds per-license SMTP configuration used by the email outbox
// worker. Exactly one row exists per license. AuthPasswordEncrypted stores the
// SMTP password after AES-GCM encryption (see api/services/crypto).
type EmailSettings struct {
	ID                     int       `json:"id" gorm:"primary_key"`
	LicenseId              string    `json:"license_id" gorm:"type:varchar(191);uniqueIndex"`
	Host                   string    `json:"host"`
	Port                   int       `json:"port"`
	TlsMode                string    `json:"tls_mode" gorm:"type:varchar(16)"` // starttls | tls | none
	AuthUser               string    `json:"auth_user"`
	AuthPasswordEncrypted  string    `json:"-" gorm:"type:text"`
	FromAddress            string    `json:"from_address"`
	FromName               string    `json:"from_name"`
	ReplyTo                string    `json:"reply_to"`
	UpdatedBy              string    `json:"updated_by"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// EmailTemplate is a license-scoped, code-keyed template whose subject and
// body are rendered at enqueue time against a per-send variables map.
// AttachmentsSpec is a JSON blob describing how to source any attachments
// (see api/services/email/attachments.go).
type EmailTemplate struct {
	ID              int       `json:"id" gorm:"primary_key"`
	LicenseId       string    `json:"license_id" gorm:"type:varchar(191);index:idx_email_templates_license_code,unique,priority:1"`
	Code            string    `json:"code" gorm:"type:varchar(128);index:idx_email_templates_license_code,unique,priority:2"`
	Name            string    `json:"name"`
	Description     string    `json:"description" gorm:"type:text"`
	SubjectTemplate string    `json:"subject_template" gorm:"type:text"`
	BodyTemplate    string    `json:"body_template" gorm:"type:mediumtext"`
	AttachmentsSpec string    `json:"attachments_spec" gorm:"type:text"` // JSON
	Status          string    `json:"status" gorm:"type:varchar(16)"`    // draft | active
	UpdatedBy       string    `json:"updated_by"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// EmailOutbox is a persisted queue of pending / sending / sent / failed emails.
// Subject and body are rendered at enqueue time; the worker renders nothing.
// Attachments hold JSON references resolved to bytes at send time.
type EmailOutbox struct {
	ID                int        `json:"id" gorm:"primary_key"`
	LicenseId         string     `json:"license_id" gorm:"type:varchar(191);index"`
	TemplateCode      string     `json:"template_code" gorm:"type:varchar(128);index"`
	FromAddress       string     `json:"from_address"`
	FromName          string     `json:"from_name"`
	ReplyTo           string     `json:"reply_to"`
	ToRecipients      string     `json:"to_recipients" gorm:"type:text"`  // JSON array of addresses
	CcRecipients      string     `json:"cc_recipients" gorm:"type:text"`  // JSON array
	BccRecipients     string     `json:"bcc_recipients" gorm:"type:text"` // JSON array
	Subject           string     `json:"subject" gorm:"type:text"`
	Body              string     `json:"body" gorm:"type:mediumtext"`
	Attachments       string     `json:"attachments" gorm:"type:text"` // JSON
	Status            string     `json:"status" gorm:"type:varchar(16);index"` // pending | sending | sent | failed
	Attempts          int        `json:"attempts"`
	MaxAttempts       int        `json:"max_attempts"`
	LastError         string     `json:"last_error" gorm:"type:text"`
	NextAttemptAt     time.Time  `json:"next_attempt_at" gorm:"index"`
	ScheduledAt       time.Time  `json:"scheduled_at"`
	SentAt            *time.Time `json:"sent_at"`
	RelatedObjectType string     `json:"related_object_type"`
	RelatedObjectID   string     `json:"related_object_id"`
	CreatedBy         string     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
}

// TableName keeps the table at the singular name the Phase 1 migration
// created, since GORM would otherwise pluralize EmailOutbox to email_outboxes.
func (EmailOutbox) TableName() string { return "email_outbox" }

// Status constants for EmailOutbox.
const (
	EmailOutboxPending = "pending"
	EmailOutboxSending = "sending"
	EmailOutboxSent    = "sent"
	EmailOutboxFailed  = "failed"
)

// TLS mode constants for EmailSettings.
const (
	EmailTLSModeNone     = "none"
	EmailTLSModeSTARTTLS = "starttls"
	EmailTLSModeTLS      = "tls"
)

// Template status constants.
const (
	EmailTemplateDraft  = "draft"
	EmailTemplateActive = "active"
)

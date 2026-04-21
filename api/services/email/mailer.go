// Package email builds, persists, and delivers outbound application email
// through a persistent outbox queue drained by a background worker.
package email

import "context"

// Mailer is the abstraction behind which delivery transports live. Phase 1
// ships only SMTPMailer; later phases can add SendGrid, SES, Postmark, etc.
type Mailer interface {
	Send(ctx context.Context, msg Message) error
}

// Message is the transport-neutral representation of an email ready to send.
// Subject and Body are already fully rendered; the worker hands this over to
// the Mailer verbatim.
type Message struct {
	FromAddress string
	FromName    string
	ReplyTo     string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string // HTML
	Attachments []Attachment
}

// Attachment carries a file to be attached to the outgoing email. Bytes are
// resolved from the outbox row's JSON spec at send time — the DB never stores
// attachment bytes directly.
type Attachment struct {
	Filename    string
	ContentType string // MIME type, e.g. "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	Content     []byte
}

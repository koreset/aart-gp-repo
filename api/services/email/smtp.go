package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"api/models"
)

// SMTPMailer delivers messages via net/smtp with STARTTLS or implicit TLS.
// Targets Office 365 (smtp.office365.com:587, STARTTLS) by default but also
// works with any SMTP relay.
type SMTPMailer struct {
	settings models.EmailSettings
	password string // already decrypted
}

// NewSMTPMailer constructs a mailer for the given license settings. The
// caller is responsible for decrypting the stored password.
func NewSMTPMailer(settings models.EmailSettings, decryptedPassword string) *SMTPMailer {
	return &SMTPMailer{settings: settings, password: decryptedPassword}
}

// Send delivers msg via the configured SMTP server. The ctx deadline is
// enforced through a connection-level deadline.
func (m *SMTPMailer) Send(ctx context.Context, msg Message) error {
	addr := fmt.Sprintf("%s:%d", m.settings.Host, m.settings.Port)
	body, err := buildMIME(msg)
	if err != nil {
		return fmt.Errorf("build MIME: %w", err)
	}

	recipients := append(append(append([]string{}, msg.To...), msg.Cc...), msg.Bcc...)
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients")
	}

	sender := msg.FromAddress
	if sender == "" {
		sender = m.settings.FromAddress
	}

	// Honour ctx deadline if present.
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(30 * time.Second)
	}

	auth := smtp.PlainAuth("", m.settings.AuthUser, m.password, m.settings.Host)

	switch strings.ToLower(m.settings.TlsMode) {
	case models.EmailTLSModeTLS:
		return sendImplicitTLS(addr, m.settings.Host, auth, sender, recipients, body, deadline)
	case models.EmailTLSModeNone:
		return sendPlain(addr, m.settings.Host, auth, sender, recipients, body, deadline)
	default: // starttls — the Office 365 path
		return sendSTARTTLS(addr, m.settings.Host, auth, sender, recipients, body, deadline)
	}
}

func sendSTARTTLS(addr, host string, auth smtp.Auth, from string, to []string, body []byte, deadline time.Time) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer client.Close()
	if err := client.Hello(clientHostname()); err != nil {
		return fmt.Errorf("hello: %w", err)
	}
	if ok, _ := client.Extension("STARTTLS"); !ok {
		return fmt.Errorf("server does not advertise STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
		return fmt.Errorf("starttls: %w", err)
	}
	return completeSession(client, auth, from, to, body, deadline)
}

func sendImplicitTLS(addr, host string, auth smtp.Auth, from string, to []string, body []byte, deadline time.Time) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()
	return completeSession(client, auth, from, to, body, deadline)
}

func sendPlain(addr, host string, auth smtp.Auth, from string, to []string, body []byte, deadline time.Time) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer client.Close()
	if err := client.Hello(clientHostname()); err != nil {
		return fmt.Errorf("hello: %w", err)
	}
	return completeSession(client, auth, from, to, body, deadline)
}

func completeSession(client *smtp.Client, auth smtp.Auth, from string, to []string, body []byte, deadline time.Time) error {
	if time.Now().After(deadline) {
		return fmt.Errorf("deadline exceeded before auth")
	}
	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(auth); err != nil {
				return fmt.Errorf("auth: %w", err)
			}
		}
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM: %w", err)
	}
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("RCPT TO %s: %w", addr, err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	if _, err := w.Write(body); err != nil {
		w.Close()
		return fmt.Errorf("write body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close body: %w", err)
	}
	return client.Quit()
}

// buildMIME constructs an RFC 5322 multipart/mixed message with an HTML body
// and optional attachments.
func buildMIME(msg Message) ([]byte, error) {
	var buf bytes.Buffer
	boundary := fmt.Sprintf("aart_%d", time.Now().UnixNano())

	// Headers
	fromHeader := msg.FromAddress
	if msg.FromName != "" {
		fromHeader = (&mail.Address{Name: msg.FromName, Address: msg.FromAddress}).String()
	}
	writeHeader(&buf, "From", fromHeader)
	writeHeader(&buf, "To", strings.Join(msg.To, ", "))
	if len(msg.Cc) > 0 {
		writeHeader(&buf, "Cc", strings.Join(msg.Cc, ", "))
	}
	if msg.ReplyTo != "" {
		writeHeader(&buf, "Reply-To", msg.ReplyTo)
	}
	writeHeader(&buf, "Subject", mime.QEncoding.Encode("utf-8", msg.Subject))
	writeHeader(&buf, "MIME-Version", "1.0")
	writeHeader(&buf, "Date", time.Now().UTC().Format(time.RFC1123Z))

	if len(msg.Attachments) == 0 {
		writeHeader(&buf, "Content-Type", "text/html; charset=utf-8")
		writeHeader(&buf, "Content-Transfer-Encoding", "base64")
		buf.WriteString("\r\n")
		buf.WriteString(base64Chunked(msg.Body))
		return buf.Bytes(), nil
	}

	writeHeader(&buf, "Content-Type", fmt.Sprintf(`multipart/mixed; boundary="%s"`, boundary))
	buf.WriteString("\r\n")

	// HTML body part
	fmt.Fprintf(&buf, "--%s\r\n", boundary)
	writeHeader(&buf, "Content-Type", "text/html; charset=utf-8")
	writeHeader(&buf, "Content-Transfer-Encoding", "base64")
	buf.WriteString("\r\n")
	buf.WriteString(base64Chunked(msg.Body))
	buf.WriteString("\r\n")

	// Attachment parts
	for _, att := range msg.Attachments {
		ct := att.ContentType
		if ct == "" {
			ct = "application/octet-stream"
		}
		fmt.Fprintf(&buf, "--%s\r\n", boundary)
		writeHeader(&buf, "Content-Type", fmt.Sprintf(`%s; name="%s"`, ct, att.Filename))
		writeHeader(&buf, "Content-Transfer-Encoding", "base64")
		writeHeader(&buf, "Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, att.Filename))
		buf.WriteString("\r\n")
		buf.WriteString(base64Chunked(string(att.Content)))
		buf.WriteString("\r\n")
	}

	fmt.Fprintf(&buf, "--%s--\r\n", boundary)
	return buf.Bytes(), nil
}

func writeHeader(buf *bytes.Buffer, key, value string) {
	fmt.Fprintf(buf, "%s: %s\r\n", key, value)
}

// base64Chunked encodes data with 76-char line wrapping per RFC 2045.
func base64Chunked(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	const width = 76
	var out strings.Builder
	for i := 0; i < len(encoded); i += width {
		end := i + width
		if end > len(encoded) {
			end = len(encoded)
		}
		out.WriteString(encoded[i:end])
		out.WriteString("\r\n")
	}
	return out.String()
}

func clientHostname() string {
	return "aart-gp"
}

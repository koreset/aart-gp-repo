package sms

import (
	appLog "api/log"
	"context"
	"strings"
)

// logProvider is the bootstrap implementation that ships with Phase 4. Every
// send is written to the application log instead of dispatched to an
// upstream, so end-to-end notification flows can be exercised in dev without
// an SMS account.
//
// When a real provider lands (Clickatell, Twilio, BulkSMS, etc.), it lives
// in its own file in this package and Use() picks it based on config.
type logProvider struct{}

// NewLog returns the dev-only log provider.
func NewLog() Provider { return logProvider{} }

func (logProvider) Name() string { return "log" }

func (logProvider) Send(_ context.Context, msg Message) (*Result, error) {
	if strings.TrimSpace(msg.To) == "" {
		return nil, ErrInvalidRecipient
	}
	appLog.WithField("sms_to", msg.To).
		WithField("sms_ref", msg.Reference).
		Info("SMS (log provider): " + truncate(msg.Body, 200))
	return &Result{Status: "queued", ProviderRef: "log-only"}, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// Package providers contains the concrete vendor adapters used by
// services/uwvendor. Each adapter satisfies uwvendor.Provider and is bound to
// exactly one uwvendor.Kind via its Kind() method.
//
// The Mock provider here is deliberately polymorphic: a single struct
// stamped with the desired Kind at construction time. It's enough for
// local development, integration tests, and the renderer's "Request X"
// buttons to round-trip a request → webhook → attachment cycle without
// any real third-party traffic.
package providers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"api/services/uw_vendor"
)

// MockConfig configures the polymorphic mock vendor provider.
type MockConfig struct {
	// Kind is the uwvendor.Kind this instance answers for. Required.
	Kind uwvendor.Kind
	// Async controls whether Submit returns awaiting_response (true, the
	// default) or complete (false). When async, the mock self-delivers a
	// webhook to the local server after `ResolveAfter` so the round-trip
	// is testable end-to-end.
	Async bool
	// ResolveAfter is how long until SimulateDelivery returns the
	// completion outcome. Zero defaults to 1 second.
	ResolveAfter time.Duration
	// WebhookSecret is the HMAC-SHA256 shared secret used to sign
	// outbound mock webhook payloads (so VerifyWebhook on the same
	// instance can validate them). Tests typically use a fixed string.
	WebhookSecret string
}

// Mock is a polymorphic in-memory vendor provider. One instance per kind.
type Mock struct {
	cfg MockConfig
}

// NewMock constructs a Mock provider. The caller must supply a Kind.
func NewMock(cfg MockConfig) *Mock {
	if cfg.ResolveAfter <= 0 {
		cfg.ResolveAfter = 1 * time.Second
	}
	if cfg.WebhookSecret == "" {
		cfg.WebhookSecret = "mock-secret"
	}
	return &Mock{cfg: cfg}
}

// Name returns the short identifier persisted on VendorRequest.Provider.
func (m *Mock) Name() string { return "mock" }

// Kind reports which service kind this mock handles.
func (m *Mock) Kind() uwvendor.Kind { return m.cfg.Kind }

// Submit generates a deterministic external request id and either
// resolves immediately (sync mode) or returns awaiting_response so the
// caller can simulate the webhook via SimulateDelivery.
func (m *Mock) Submit(_ context.Context, req uwvendor.SubmitRequest) (*uwvendor.SubmitResult, error) {
	externalID := "mock-" + uuid.New().String()
	if !m.cfg.Async {
		return &uwvendor.SubmitResult{
			ExternalRequestID: externalID,
			Status:            uwvendor.StatusComplete,
			ProviderName:      m.Name(),
			CostCents:         m.defaultCost(),
			RawPayload:        []byte(`{"mock":true,"resolved":"immediate"}`),
		}, nil
	}
	return &uwvendor.SubmitResult{
		ExternalRequestID: externalID,
		Status:            uwvendor.StatusAwaitingResponse,
		ProviderName:      m.Name(),
		CostCents:         m.defaultCost(),
		RawPayload:        []byte(`{"mock":true,"resolved":"async"}`),
	}, nil
}

// VerifyWebhook validates the HMAC-SHA256 signature on a webhook body.
// The signature is expected on the `X-Signature` header as a hex-encoded
// digest of the body using WebhookSecret.
func (m *Mock) VerifyWebhook(headers map[string][]string, body []byte) error {
	sig := headerValue(headers, "X-Signature")
	if sig == "" {
		return errors.New("missing X-Signature")
	}
	mac := hmac.New(sha256.New, []byte(m.cfg.WebhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}

// HandleWebhook decodes a mock webhook body into a WebhookOutcome.
//
// Body shape (kept minimal so tests can construct it inline):
//
//	{
//	  "external_request_id": "mock-...",
//	  "status": "complete" | "failed",
//	  "attachment": { "kind": "medical_report", "file_name": "panel.pdf",
//	                  "content_type": "application/pdf",
//	                  "body_base64": "..." }   // optional
//	  "message": "Panel results delivered"
//	}
type mockWebhookBody struct {
	ExternalRequestID string `json:"external_request_id"`
	Status            string `json:"status"`
	Attachment        *struct {
		Kind        string `json:"kind"`
		FileName    string `json:"file_name"`
		ContentType string `json:"content_type"`
		Body        string `json:"body"`
	} `json:"attachment"`
	Message string `json:"message"`
}

// HandleWebhook decodes the mock body and returns the canonical outcome.
func (m *Mock) HandleWebhook(_ context.Context, body []byte) (*uwvendor.WebhookOutcome, error) {
	var probe mockWebhookBody
	if err := json.Unmarshal(body, &probe); err != nil {
		return nil, fmt.Errorf("decode mock body: %w", err)
	}
	if probe.ExternalRequestID == "" {
		return nil, errors.New("external_request_id required")
	}
	status := uwvendor.Status(strings.TrimSpace(probe.Status))
	if status == "" {
		status = uwvendor.StatusComplete
	}
	outcome := &uwvendor.WebhookOutcome{
		ExternalRequestID: probe.ExternalRequestID,
		Status:            status,
		Message:           probe.Message,
	}
	if probe.Attachment != nil {
		outcome.Attachment = &uwvendor.AttachmentPayload{
			Kind:        probe.Attachment.Kind,
			FileName:    probe.Attachment.FileName,
			ContentType: probe.Attachment.ContentType,
			Body:        []byte(probe.Attachment.Body),
		}
	}
	return outcome, nil
}

// SimulateDelivery builds a signed webhook envelope for the given external
// request id. Tests + the dev "fire-mock-webhook" controller call this to
// drive the async round-trip.
func (m *Mock) SimulateDelivery(externalID string, status uwvendor.Status, attachment *uwvendor.AttachmentPayload) (headers map[string][]string, body []byte) {
	wb := mockWebhookBody{
		ExternalRequestID: externalID,
		Status:            string(status),
	}
	if attachment != nil {
		wb.Attachment = &struct {
			Kind        string `json:"kind"`
			FileName    string `json:"file_name"`
			ContentType string `json:"content_type"`
			Body        string `json:"body"`
		}{
			Kind:        attachment.Kind,
			FileName:    attachment.FileName,
			ContentType: attachment.ContentType,
			Body:        string(attachment.Body),
		}
	}
	body, _ = json.Marshal(wb)
	mac := hmac.New(sha256.New, []byte(m.cfg.WebhookSecret))
	mac.Write(body)
	sig := hex.EncodeToString(mac.Sum(nil))
	return map[string][]string{
		http.CanonicalHeaderKey("X-Signature"):    {sig},
		http.CanonicalHeaderKey("X-Actor-Email"):  {"mock@vendor"},
		http.CanonicalHeaderKey("Content-Type"):   {"application/json"},
	}, body
}

func (m *Mock) defaultCost() int {
	// Synthetic per-kind cost so the cost-cents column is exercised.
	switch m.cfg.Kind {
	case uwvendor.KindPathology:
		return 1500 // R15.00
	case uwvendor.KindGPRecords:
		return 2500
	case uwvendor.KindEsign:
		return 200
	case uwvendor.KindSMS:
		return 25
	}
	return 0
}

func headerValue(h map[string][]string, name string) string {
	canon := http.CanonicalHeaderKey(name)
	if vs, ok := h[canon]; ok && len(vs) > 0 {
		return vs[0]
	}
	for k, vs := range h {
		if strings.EqualFold(k, name) && len(vs) > 0 {
			return vs[0]
		}
	}
	return ""
}

// Compile-time assertions that Mock satisfies both interfaces. Goes
// unused at runtime but catches regressions during refactors.
var (
	_ uwvendor.Provider       = (*Mock)(nil)
	_ uwvendor.WebhookCapable = (*Mock)(nil)
)

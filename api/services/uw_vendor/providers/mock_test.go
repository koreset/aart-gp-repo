package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"api/services/uw_vendor"
)

func TestMock_Submit_Sync(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindPathology, Async: false, WebhookSecret: "s"})
	res, err := m.Submit(context.Background(), uwvendor.SubmitRequest{
		Kind: uwvendor.KindPathology, CaseID: 1, Subject: "MEMBER-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != uwvendor.StatusComplete {
		t.Errorf("sync mock should return complete, got %s", res.Status)
	}
	if !strings.HasPrefix(res.ExternalRequestID, "mock-") {
		t.Errorf("external id should be prefixed, got %s", res.ExternalRequestID)
	}
}

func TestMock_Submit_Async_ReturnsAwaiting(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindEsign, Async: true, WebhookSecret: "s"})
	res, err := m.Submit(context.Background(), uwvendor.SubmitRequest{
		Kind: uwvendor.KindEsign, Subject: "broker@example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != uwvendor.StatusAwaitingResponse {
		t.Errorf("async mock should return awaiting_response, got %s", res.Status)
	}
}

func TestMock_VerifyWebhook_HMACMismatch(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindSMS, WebhookSecret: "secret-a"})
	headers, body := m.SimulateDelivery("mock-1", uwvendor.StatusComplete, nil)

	// Swap to a different secret on the verifier side — the existing
	// signature should no longer match.
	other := NewMock(MockConfig{Kind: uwvendor.KindSMS, WebhookSecret: "secret-b"})
	if err := other.VerifyWebhook(headers, body); err == nil {
		t.Errorf("expected signature mismatch, got nil")
	}
}

func TestMock_VerifyWebhook_HMACMatch(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindSMS, WebhookSecret: "s"})
	headers, body := m.SimulateDelivery("mock-1", uwvendor.StatusComplete, nil)
	if err := m.VerifyWebhook(headers, body); err != nil {
		t.Errorf("same secret should verify, got %v", err)
	}
}

func TestMock_HandleWebhook_DecodesOutcome(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindGPRecords, WebhookSecret: "s"})
	attachment := &uwvendor.AttachmentPayload{
		Kind:        "medical_report",
		FileName:    "panel.pdf",
		ContentType: "application/pdf",
		Body:        []byte("dummy"),
	}
	_, body := m.SimulateDelivery("mock-42", uwvendor.StatusComplete, attachment)

	outcome, err := m.HandleWebhook(context.Background(), body)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if outcome.ExternalRequestID != "mock-42" {
		t.Errorf("external id mismatch: %s", outcome.ExternalRequestID)
	}
	if outcome.Status != uwvendor.StatusComplete {
		t.Errorf("status mismatch: %s", outcome.Status)
	}
	if outcome.Attachment == nil || outcome.Attachment.FileName != "panel.pdf" {
		t.Errorf("attachment missing or wrong: %+v", outcome.Attachment)
	}
}

func TestMock_HandleWebhook_RejectsMissingExternalID(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindEsign})
	_, err := m.HandleWebhook(context.Background(), []byte(`{"status":"complete"}`))
	if err == nil {
		t.Error("expected error for missing external_request_id")
	}
}

func TestMock_SimulateDelivery_SetsSignatureHeader(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindSMS, WebhookSecret: "s"})
	headers, _ := m.SimulateDelivery("mock-1", uwvendor.StatusComplete, nil)
	if headers[http.CanonicalHeaderKey("X-Signature")][0] == "" {
		t.Errorf("X-Signature header should be set")
	}
}

func TestMock_DefaultCostByKind(t *testing.T) {
	cases := map[uwvendor.Kind]int{
		uwvendor.KindPathology: 1500,
		uwvendor.KindGPRecords: 2500,
		uwvendor.KindEsign:     200,
		uwvendor.KindSMS:       25,
	}
	for kind, want := range cases {
		m := NewMock(MockConfig{Kind: kind, Async: false})
		res, _ := m.Submit(context.Background(), uwvendor.SubmitRequest{Kind: kind})
		if res.CostCents != want {
			t.Errorf("kind %s: cost %d, want %d", kind, res.CostCents, want)
		}
	}
}

func TestMock_HandleWebhook_DefaultsStatusToComplete(t *testing.T) {
	m := NewMock(MockConfig{Kind: uwvendor.KindSMS})
	body, _ := json.Marshal(map[string]any{"external_request_id": "x"})
	outcome, err := m.HandleWebhook(context.Background(), body)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Status != uwvendor.StatusComplete {
		t.Errorf("missing status should default to complete, got %s", outcome.Status)
	}
}

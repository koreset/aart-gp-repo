package providers

import (
	"api/services/bav"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const happyPathResponse = `{
  "success": true,
  "requestId": "req-abc-123",
  "service": "bank-account-verification",
  "results": {
    "identity_and_account_verified": true,
    "summary": "All checks passed",
    "verification_results": {
      "Status": "Verified",
      "accountFound": "Y",
      "accountOpen": "Y",
      "identityMatch": "Y",
      "accountTypeMatch": "Y",
      "acceptsCredits": "Y",
      "acceptsDebits": "N"
    }
  }
}`

func newTestAdapter(t *testing.T, handler http.HandlerFunc) (*VerifyNow, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	a := NewVerifyNow(VerifyNowConfig{
		APIKey:  "test-key",
		Mode:    "test",
		BaseURL: srv.URL,
		Timeout: 5 * time.Second,
	})
	return a, srv
}

func sampleRequest() bav.VerifyRequest {
	return bav.VerifyRequest{
		FirstName:         "Thandi",
		Surname:           "Nkosi",
		IdentityNumber:    "9001015009087",
		BankAccountNumber: "1234567890",
		BankBranchCode:    "250655",
		BankAccountType:   "CHEQUE",
		IdempotencyKey:    "idem-xyz",
	}
}

func TestVerifyNow_Name(t *testing.T) {
	a := NewVerifyNow(VerifyNowConfig{APIKey: "k"})
	if a.Name() != "verifynow" {
		t.Fatalf("Name = %q, want %q", a.Name(), "verifynow")
	}
}

func TestVerifyNow_Verify_HappyPath(t *testing.T) {
	var receivedHeaders http.Header
	var receivedBody verifyNowRequestBody
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header.Clone()
		if r.URL.Path != verifyNowVerifyPath {
			t.Errorf("path = %q, want %q", r.URL.Path, verifyNowVerifyPath)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &receivedBody); err != nil {
			t.Errorf("server couldn't decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, happyPathResponse)
	})

	res, err := a.Verify(context.Background(), sampleRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Status != bav.StatusComplete {
		t.Errorf("Status = %q, want complete", res.Status)
	}
	if !res.Verified {
		t.Error("Verified should be true")
	}
	if res.Provider != "verifynow" {
		t.Errorf("Provider = %q, want verifynow", res.Provider)
	}
	if res.ProviderRequestID != "req-abc-123" {
		t.Errorf("ProviderRequestID = %q, want req-abc-123", res.ProviderRequestID)
	}
	if res.AccountFound != bav.TriYes {
		t.Errorf("AccountFound = %q, want yes", res.AccountFound)
	}
	if res.AcceptsDebits != bav.TriNo {
		t.Errorf("AcceptsDebits = %q, want no", res.AcceptsDebits)
	}
	if res.ProviderStatusText != "Verified" {
		t.Errorf("ProviderStatusText = %q, want Verified", res.ProviderStatusText)
	}
	if string(res.RawPayload) == "" {
		t.Error("RawPayload should be populated")
	}

	if got := receivedHeaders.Get("x-api-key"); got != "test-key" {
		t.Errorf("x-api-key header = %q, want test-key", got)
	}
	if got := receivedHeaders.Get("Idempotency-Key"); got != "idem-xyz" {
		t.Errorf("Idempotency-Key = %q, want idem-xyz", got)
	}
	if got := receivedHeaders.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	if receivedBody.Mode != "test" {
		t.Errorf("Mode = %q, want test", receivedBody.Mode)
	}
	if receivedBody.IdentityType != "IDNumber" {
		t.Errorf("IdentityType default = %q, want IDNumber", receivedBody.IdentityType)
	}
}

func TestVerifyNow_Verify_GeneratesIdempotencyKeyWhenMissing(t *testing.T) {
	var idem string
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		idem = r.Header.Get("Idempotency-Key")
		_, _ = io.WriteString(w, happyPathResponse)
	})
	req := sampleRequest()
	req.IdempotencyKey = ""
	if _, err := a.Verify(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idem == "" {
		t.Fatal("expected Idempotency-Key header to be auto-generated")
	}
}

func TestVerifyNow_Verify_MissingAPIKey(t *testing.T) {
	a := NewVerifyNow(VerifyNowConfig{})
	_, err := a.Verify(context.Background(), sampleRequest())
	if !errors.Is(err, bav.ErrProviderNotConfigured) {
		t.Fatalf("err = %v, want ErrProviderNotConfigured", err)
	}
}

func TestVerifyNow_Verify_Unauthorized(t *testing.T) {
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"bad key"}`)
	})
	_, err := a.Verify(context.Background(), sampleRequest())
	if !errors.Is(err, bav.ErrUnauthorized) {
		t.Fatalf("err = %v, want ErrUnauthorized", err)
	}
}

func TestVerifyNow_Verify_RateLimited(t *testing.T) {
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	_, err := a.Verify(context.Background(), sampleRequest())
	if !errors.Is(err, bav.ErrRateLimited) {
		t.Fatalf("err = %v, want ErrRateLimited", err)
	}
}

func TestVerifyNow_Verify_ServerError(t *testing.T) {
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})
	_, err := a.Verify(context.Background(), sampleRequest())
	if !errors.Is(err, bav.ErrProviderUnavailable) {
		t.Fatalf("err = %v, want ErrProviderUnavailable", err)
	}
}

func TestVerifyNow_Verify_BadRequest(t *testing.T) {
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"missing field"}`)
	})
	_, err := a.Verify(context.Background(), sampleRequest())
	if !errors.Is(err, bav.ErrInvalidInput) {
		t.Fatalf("err = %v, want ErrInvalidInput", err)
	}
}

func TestVerifyNow_Verify_MalformedJSON(t *testing.T) {
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `not-json`)
	})
	_, err := a.Verify(context.Background(), sampleRequest())
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
	if !strings.Contains(err.Error(), "decode") {
		t.Fatalf("err = %v, want decode error", err)
	}
}

func TestVerifyNow_Verify_ContextCancelled(t *testing.T) {
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = io.WriteString(w, happyPathResponse)
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := a.Verify(ctx, sampleRequest())
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestVerifyNow_Verify_CustomIdentityType(t *testing.T) {
	var got verifyNowRequestBody
	a, _ := newTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &got)
		_, _ = io.WriteString(w, happyPathResponse)
	})
	req := sampleRequest()
	req.IdentityType = "Passport"
	if _, err := a.Verify(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.IdentityType != "Passport" {
		t.Errorf("IdentityType = %q, want Passport", got.IdentityType)
	}
}

func TestNewVerifyNow_Defaults(t *testing.T) {
	a := NewVerifyNow(VerifyNowConfig{APIKey: "k"})
	if a.cfg.BaseURL != verifyNowDefaultBaseURL {
		t.Errorf("BaseURL default = %q, want %q", a.cfg.BaseURL, verifyNowDefaultBaseURL)
	}
	if a.cfg.Timeout != verifyNowDefaultTimeout {
		t.Errorf("Timeout default = %v, want %v", a.cfg.Timeout, verifyNowDefaultTimeout)
	}
	if a.client.Timeout != verifyNowDefaultTimeout {
		t.Errorf("client.Timeout = %v, want %v", a.client.Timeout, verifyNowDefaultTimeout)
	}
}

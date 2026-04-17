// Package providers contains concrete adapters that implement bav.Provider.
// Each adapter is responsible for translating the canonical VerifyRequest /
// VerifyResult into a specific provider's wire format and mapping transport
// errors to the sentinel errors declared in the bav package.
package providers

import (
	"api/services/bav"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	verifyNowName           = "verifynow"
	verifyNowDefaultBaseURL = "https://www.verifynow.co.za"
	verifyNowDefaultTimeout = 45 * time.Second
	verifyNowVerifyPath     = "/api/external/bank-account-verification"
)

// VerifyNowConfig holds the credentials and endpoint configuration for the
// VerifyNow BAV adapter. BaseURL and Timeout default to production values
// when left empty. Mode is passed through to VerifyNow verbatim.
type VerifyNowConfig struct {
	APIKey  string
	Mode    string
	BaseURL string
	Timeout time.Duration
}

// VerifyNow is the VerifyNow BAV provider adapter.
type VerifyNow struct {
	cfg    VerifyNowConfig
	client *http.Client
}

// NewVerifyNow constructs a VerifyNow adapter. An empty APIKey is tolerated
// here so callers can build the adapter at startup before deciding whether to
// activate it; Verify returns ErrProviderNotConfigured if APIKey is still
// empty at call time.
func NewVerifyNow(cfg VerifyNowConfig) *VerifyNow {
	if cfg.BaseURL == "" {
		cfg.BaseURL = verifyNowDefaultBaseURL
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = verifyNowDefaultTimeout
	}
	return &VerifyNow{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// Name returns the short identifier used in configuration.
func (a *VerifyNow) Name() string { return verifyNowName }

type verifyNowRequestBody struct {
	FirstName         string `json:"firstName"`
	Surname           string `json:"surname"`
	IdentityNumber    string `json:"identityNumber"`
	IdentityType      string `json:"identityType"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankBranchCode    string `json:"bankBranchCode"`
	BankAccountType   string `json:"bankAccountType"`
	Mode              string `json:"mode"`
}

type verifyNowResponseBody struct {
	Success   bool                     `json:"success"`
	RequestID string                   `json:"requestId"`
	Service   string                   `json:"service"`
	Results   verifyNowResponseResults `json:"results"`
}

type verifyNowResponseResults struct {
	IdentityAndAccountVerified bool                           `json:"identity_and_account_verified"`
	Summary                    string                         `json:"summary"`
	VerificationResults        verifyNowVerificationResultRow `json:"verification_results"`
}

type verifyNowVerificationResultRow struct {
	Status           string `json:"Status"`
	AccountFound     string `json:"accountFound"`
	AccountOpen      string `json:"accountOpen"`
	IdentityMatch    string `json:"identityMatch"`
	AccountTypeMatch string `json:"accountTypeMatch"`
	AcceptsCredits   string `json:"acceptsCredits"`
	AcceptsDebits    string `json:"acceptsDebits"`
}

// Verify calls the VerifyNow bank-account-verification endpoint and maps the
// response onto the canonical VerifyResult.
func (a *VerifyNow) Verify(ctx context.Context, req bav.VerifyRequest) (*bav.VerifyResult, error) {
	if a.cfg.APIKey == "" {
		return nil, fmt.Errorf("%w: VERIFYNOW_API_KEY is not configured", bav.ErrProviderNotConfigured)
	}

	identityType := req.IdentityType
	if identityType == "" {
		identityType = "IDNumber"
	}

	body, err := json.Marshal(verifyNowRequestBody{
		FirstName:         req.FirstName,
		Surname:           req.Surname,
		IdentityNumber:    req.IdentityNumber,
		IdentityType:      identityType,
		BankAccountNumber: req.BankAccountNumber,
		BankBranchCode:    req.BankBranchCode,
		BankAccountType:   req.BankAccountType,
		Mode:              a.cfg.Mode,
	})
	if err != nil {
		return nil, fmt.Errorf("verifynow marshal error: %w", err)
	}

	url := a.cfg.BaseURL + verifyNowVerifyPath
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("verifynow request error: %w", err)
	}

	idempotencyKey := req.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = uuid.New().String()
	}
	httpReq.Header.Set("x-api-key", a.cfg.APIKey)
	httpReq.Header.Set("Idempotency-Key", idempotencyKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", bav.ErrProviderUnavailable, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("verifynow read body error: %w", err)
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized:
		return nil, fmt.Errorf("%w: verifynow returned 401", bav.ErrUnauthorized)
	case resp.StatusCode == http.StatusTooManyRequests:
		return nil, fmt.Errorf("%w: verifynow returned 429", bav.ErrRateLimited)
	case resp.StatusCode == http.StatusBadRequest:
		return nil, fmt.Errorf("%w: verifynow returned 400: %s", bav.ErrInvalidInput, truncate(respBody, 256))
	case resp.StatusCode >= 500:
		return nil, fmt.Errorf("%w: verifynow returned %d", bav.ErrProviderUnavailable, resp.StatusCode)
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("verifynow returned unexpected status %d: %s", resp.StatusCode, truncate(respBody, 256))
	}

	var parsed verifyNowResponseBody
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("verifynow response decode error: %w", err)
	}

	d := parsed.Results.VerificationResults
	return &bav.VerifyResult{
		Status:             bav.StatusComplete,
		Verified:           parsed.Results.IdentityAndAccountVerified,
		Summary:            parsed.Results.Summary,
		AccountFound:       bav.ParseTriState(d.AccountFound),
		AccountOpen:        bav.ParseTriState(d.AccountOpen),
		IdentityMatch:      bav.ParseTriState(d.IdentityMatch),
		AccountTypeMatch:   bav.ParseTriState(d.AccountTypeMatch),
		AcceptsCredits:     bav.ParseTriState(d.AcceptsCredits),
		AcceptsDebits:      bav.ParseTriState(d.AcceptsDebits),
		Provider:           verifyNowName,
		ProviderRequestID:  parsed.RequestID,
		ProviderStatusText: d.Status,
		RawPayload:         respBody,
	}, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}

// Poll is a no-op for VerifyNow — the provider's bank-account-verification
// call is synchronous, so any jobID is unrecognised. Returns ErrNotSupported
// so the registry surfaces an honest 501-shaped error.
func (a *VerifyNow) Poll(_ context.Context, _ string) (*bav.VerifyResult, error) {
	return nil, bav.ErrNotSupported
}

var _ bav.Provider = (*VerifyNow)(nil)

package providers

import (
	"api/services/bav"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	mockName          = "mock"
	mockAsyncResolveAfter = 3 * time.Second
)

// MockConfig configures the mock BAV provider used for local development,
// integration tests, and end-to-end smoke checks.
type MockConfig struct {
	// Async flips the mock from returning StatusComplete synchronously to
	// returning StatusPending and resolving on Poll after ResolveAfter.
	// Wired from the MOCK_BAV_ASYNC env var.
	Async bool
	// ResolveAfter is the synthetic delay before a pending job resolves.
	// Zero falls back to mockAsyncResolveAfter.
	ResolveAfter time.Duration
}

// Mock is a deterministic in-memory BAV provider. Sync mode returns a
// canned verified result on Verify. Async mode returns StatusPending with
// a generated jobID; Poll resolves to StatusComplete once ResolveAfter has
// elapsed since the original Verify call.
type Mock struct {
	cfg MockConfig

	mu   sync.Mutex
	jobs map[string]time.Time // jobID -> dispatchedAt (UTC)
}

// NewMock constructs the mock provider.
func NewMock(cfg MockConfig) *Mock {
	if cfg.ResolveAfter <= 0 {
		cfg.ResolveAfter = mockAsyncResolveAfter
	}
	return &Mock{
		cfg:  cfg,
		jobs: make(map[string]time.Time),
	}
}

// Name returns the short identifier used in configuration.
func (m *Mock) Name() string { return mockName }

// Verify returns a synthetic verified result in sync mode; in async mode it
// returns StatusPending with a fresh jobID that Poll can later resolve.
func (m *Mock) Verify(_ context.Context, _ bav.VerifyRequest) (*bav.VerifyResult, error) {
	if !m.cfg.Async {
		return m.completeResult(""), nil
	}

	jobID := uuid.New().String()
	m.mu.Lock()
	m.jobs[jobID] = time.Now().UTC()
	m.mu.Unlock()

	return &bav.VerifyResult{
		Status:             bav.StatusPending,
		Summary:            "Verification in progress (mock async)",
		Provider:           mockName,
		ProviderJobID:      jobID,
		ProviderStatusText: "Pending",
	}, nil
}

// Poll resolves a pending mock job. Returns StatusPending while the
// synthetic delay hasn't elapsed, otherwise StatusComplete with a verified
// result. An unknown jobID yields ErrInvalidInput.
func (m *Mock) Poll(_ context.Context, jobID string) (*bav.VerifyResult, error) {
	m.mu.Lock()
	dispatchedAt, ok := m.jobs[jobID]
	m.mu.Unlock()
	if !ok {
		return nil, bav.ErrInvalidInput
	}

	if time.Since(dispatchedAt) < m.cfg.ResolveAfter {
		return &bav.VerifyResult{
			Status:             bav.StatusPending,
			Summary:            "Verification still in progress (mock async)",
			Provider:           mockName,
			ProviderJobID:      jobID,
			ProviderStatusText: "Pending",
		}, nil
	}

	m.mu.Lock()
	delete(m.jobs, jobID)
	m.mu.Unlock()
	return m.completeResult(jobID), nil
}

func (m *Mock) completeResult(jobID string) *bav.VerifyResult {
	return &bav.VerifyResult{
		Status:             bav.StatusComplete,
		Verified:           true,
		Summary:            "Mock verification complete",
		AccountFound:       bav.TriYes,
		AccountOpen:        bav.TriYes,
		IdentityMatch:      bav.TriYes,
		AccountTypeMatch:   bav.TriYes,
		AcceptsCredits:     bav.TriYes,
		AcceptsDebits:      bav.TriYes,
		Provider:           mockName,
		ProviderRequestID:  "mock-" + uuid.New().String(),
		ProviderJobID:      jobID,
		ProviderStatusText: "Verified",
		RawPayload:         []byte(`{"mock":true}`),
	}
}

var _ bav.Provider = (*Mock)(nil)

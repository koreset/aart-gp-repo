package bav

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type stubProvider struct {
	name       string
	callCount  int
	pollCount  int
	result     *VerifyResult
	pollResult *VerifyResult
	err        error
	pollErr    error
	lastReq    VerifyRequest
	lastJobID  string
}

func (s *stubProvider) Name() string { return s.name }

func (s *stubProvider) Verify(_ context.Context, req VerifyRequest) (*VerifyResult, error) {
	s.callCount++
	s.lastReq = req
	return s.result, s.err
}

func (s *stubProvider) Poll(_ context.Context, jobID string) (*VerifyResult, error) {
	s.pollCount++
	s.lastJobID = jobID
	if s.pollErr != nil || s.pollResult != nil {
		return s.pollResult, s.pollErr
	}
	return nil, ErrNotSupported
}

type stubLogger struct {
	cached      *VerifyResult
	cacheErr    error
	lookupCount int
	records     []LogEntry
	lookupKey   string
	recordErr   error
}

func (l *stubLogger) LookupSuccessByKey(_ context.Context, key string, _ time.Duration) (*VerifyResult, error) {
	l.lookupCount++
	l.lookupKey = key
	return l.cached, l.cacheErr
}

func (l *stubLogger) Record(_ context.Context, entry LogEntry) error {
	l.records = append(l.records, entry)
	return l.recordErr
}

func sampleRequest() VerifyRequest {
	return VerifyRequest{
		IdentityNumber:    "9001015009087",
		BankAccountNumber: "1234567890",
		BankBranchCode:    "250655",
	}
}

func TestDeriveIdempotencyKey_StableAcrossCalls(t *testing.T) {
	r := sampleRequest()
	k1 := r.DeriveIdempotencyKey("verifynow")
	k2 := r.DeriveIdempotencyKey("verifynow")
	if k1 != k2 {
		t.Fatalf("keys differ: %q vs %q", k1, k2)
	}
	if len(k1) != 64 {
		t.Fatalf("expected sha256 hex (64 chars), got %d", len(k1))
	}
}

func TestDeriveIdempotencyKey_DiffersPerProvider(t *testing.T) {
	r := sampleRequest()
	if r.DeriveIdempotencyKey("verifynow") == r.DeriveIdempotencyKey("lexisnexis") {
		t.Fatal("expected different keys per provider")
	}
}

func TestDeriveIdempotencyKey_DiffersPerAttempt(t *testing.T) {
	r1 := sampleRequest()
	r1.Attempt = 1
	r2 := sampleRequest()
	r2.Attempt = 2
	if r1.DeriveIdempotencyKey("verifynow") == r2.DeriveIdempotencyKey("verifynow") {
		t.Fatal("expected different keys per attempt number")
	}
}

func TestDeriveIdempotencyKey_DiffersWhenBankingChanges(t *testing.T) {
	r1 := sampleRequest()
	r2 := sampleRequest()
	r2.BankAccountNumber = "9999999999"
	if r1.DeriveIdempotencyKey("verifynow") == r2.DeriveIdempotencyKey("verifynow") {
		t.Fatal("expected different keys when account number changes")
	}
}

func TestRegistry_Verify_NoProviderReturnsConfiguredError(t *testing.T) {
	reg := NewRegistry(nil)
	_, err := reg.Verify(context.Background(), sampleRequest())
	if !errors.Is(err, ErrProviderNotConfigured) {
		t.Fatalf("err = %v, want ErrProviderNotConfigured", err)
	}
}

func TestRegistry_Verify_NilRegistryReturnsConfiguredError(t *testing.T) {
	var reg *Registry
	_, err := reg.Verify(context.Background(), sampleRequest())
	if !errors.Is(err, ErrProviderNotConfigured) {
		t.Fatalf("err = %v, want ErrProviderNotConfigured", err)
	}
}

func TestRegistry_Verify_WithoutLoggerStillCallsProvider(t *testing.T) {
	p := &stubProvider{name: "test", result: &VerifyResult{Status: StatusComplete}}
	reg := NewRegistry(p)
	got, err := reg.Verify(context.Background(), sampleRequest())
	if err != nil || got == nil {
		t.Fatalf("unexpected: err=%v got=%v", err, got)
	}
	if p.callCount != 1 {
		t.Errorf("provider call count = %d, want 1", p.callCount)
	}
}

func TestRegistry_Verify_DerivesKeyWhenUnset(t *testing.T) {
	p := &stubProvider{name: "test", result: &VerifyResult{Status: StatusComplete}}
	reg := NewRegistry(p)
	req := sampleRequest()
	_, _ = reg.Verify(context.Background(), req)
	if p.lastReq.IdempotencyKey == "" {
		t.Fatal("expected derived idempotency key on request passed to adapter")
	}
	if len(p.lastReq.IdempotencyKey) != 64 {
		t.Fatalf("expected sha256 hex, got %q", p.lastReq.IdempotencyKey)
	}
}

func TestRegistry_Verify_CachedSuccessShortCircuits(t *testing.T) {
	cached := &VerifyResult{Status: StatusComplete, Verified: true, Summary: "cached"}
	p := &stubProvider{name: "test", result: &VerifyResult{Status: StatusComplete}}
	logger := &stubLogger{cached: cached}
	reg := NewRegistry(p).WithLogger(logger)

	got, err := reg.Verify(context.Background(), sampleRequest())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != cached {
		t.Error("expected cached result to be returned")
	}
	if p.callCount != 0 {
		t.Errorf("provider should not be called on cache hit; callCount=%d", p.callCount)
	}
	if len(logger.records) != 0 {
		t.Errorf("no new record should be written on cache hit; got %d", len(logger.records))
	}
}

func TestRegistry_Verify_LogsSuccess(t *testing.T) {
	result := &VerifyResult{
		Status:            StatusComplete,
		Verified:          true,
		ProviderRequestID: "req-42",
	}
	p := &stubProvider{name: "test", result: result}
	logger := &stubLogger{}
	reg := NewRegistry(p).WithLogger(logger)

	_, err := reg.Verify(context.Background(), sampleRequest())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(logger.records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logger.records))
	}
	rec := logger.records[0]
	if rec.Status != StatusComplete {
		t.Errorf("Status = %q, want complete", rec.Status)
	}
	if rec.ProviderRequestID != "req-42" {
		t.Errorf("ProviderRequestID = %q, want req-42", rec.ProviderRequestID)
	}
	if rec.Err != nil {
		t.Errorf("Err should be nil on success; got %v", rec.Err)
	}
}

func TestRegistry_Verify_LogsFailure(t *testing.T) {
	callErr := errors.New("provider down")
	p := &stubProvider{name: "test", err: callErr}
	logger := &stubLogger{}
	reg := NewRegistry(p).WithLogger(logger)

	_, err := reg.Verify(context.Background(), sampleRequest())
	if err == nil {
		t.Fatal("expected error to propagate")
	}
	if len(logger.records) != 1 {
		t.Fatalf("expected 1 record on failure, got %d", len(logger.records))
	}
	rec := logger.records[0]
	if rec.Status != StatusFailed {
		t.Errorf("Status = %q, want failed", rec.Status)
	}
	if rec.Err == nil || !strings.Contains(rec.Err.Error(), "provider down") {
		t.Errorf("Err = %v, want wrapped provider error", rec.Err)
	}
}

func TestRegistry_Verify_CacheLookupFailureFallsThrough(t *testing.T) {
	// Logger lookup errors must not break the call — fall through to the provider.
	p := &stubProvider{name: "test", result: &VerifyResult{Status: StatusComplete}}
	logger := &stubLogger{cacheErr: errors.New("db down")}
	reg := NewRegistry(p).WithLogger(logger)

	got, err := reg.Verify(context.Background(), sampleRequest())
	if err != nil || got == nil {
		t.Fatalf("unexpected: err=%v got=%v", err, got)
	}
	if p.callCount != 1 {
		t.Errorf("provider should still be called; callCount=%d", p.callCount)
	}
}

func TestRegistry_Verify_RecordErrorIsSwallowed(t *testing.T) {
	// An audit write failure must not break the user-facing call.
	p := &stubProvider{name: "test", result: &VerifyResult{Status: StatusComplete}}
	logger := &stubLogger{recordErr: errors.New("disk full")}
	reg := NewRegistry(p).WithLogger(logger)

	got, err := reg.Verify(context.Background(), sampleRequest())
	if err != nil {
		t.Fatalf("expected record failure to be swallowed; got err=%v", err)
	}
	if got == nil {
		t.Fatal("expected result to be returned even when record fails")
	}
}

func TestSetDefaultAndPackageVerify(t *testing.T) {
	p := &stubProvider{name: "test", result: &VerifyResult{Status: StatusComplete}}
	reg := NewRegistry(p)
	SetDefault(reg)
	t.Cleanup(func() { SetDefault(nil) })

	got, err := Verify(context.Background(), sampleRequest())
	if err != nil || got == nil {
		t.Fatalf("unexpected: err=%v got=%v", err, got)
	}
	if p.callCount != 1 {
		t.Errorf("expected provider call; callCount=%d", p.callCount)
	}
}

func TestPackageVerify_WithoutDefaultRegistryReturnsError(t *testing.T) {
	SetDefault(nil)
	_, err := Verify(context.Background(), sampleRequest())
	if !errors.Is(err, ErrProviderNotConfigured) {
		t.Fatalf("err = %v, want ErrProviderNotConfigured", err)
	}
}

func TestRegistry_Poll_RoutesThroughProvider(t *testing.T) {
	p := &stubProvider{
		name:       "test",
		pollResult: &VerifyResult{Status: StatusComplete, Verified: true, ProviderJobID: "job-1"},
	}
	reg := NewRegistry(p)
	res, err := reg.Poll(context.Background(), "job-1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res == nil || res.Status != StatusComplete {
		t.Fatalf("unexpected result: %+v", res)
	}
	if p.pollCount != 1 {
		t.Errorf("pollCount = %d, want 1", p.pollCount)
	}
	if p.lastJobID != "job-1" {
		t.Errorf("lastJobID = %q, want job-1", p.lastJobID)
	}
}

func TestRegistry_Poll_LogsTerminalPolls(t *testing.T) {
	result := &VerifyResult{Status: StatusComplete, Verified: true, ProviderJobID: "job-2", ProviderRequestID: "req-2"}
	p := &stubProvider{name: "test", pollResult: result}
	logger := &stubLogger{}
	reg := NewRegistry(p).WithLogger(logger)

	_, err := reg.Poll(context.Background(), "job-2")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(logger.records) != 1 {
		t.Fatalf("expected 1 log record on terminal poll, got %d", len(logger.records))
	}
	rec := logger.records[0]
	if rec.ProviderJobID != "job-2" {
		t.Errorf("ProviderJobID = %q, want job-2", rec.ProviderJobID)
	}
	if rec.Status != StatusComplete {
		t.Errorf("Status = %q, want complete", rec.Status)
	}
}

func TestRegistry_Poll_SkipsPendingPolls(t *testing.T) {
	// Intermediate pending polls must not produce audit rows — otherwise a
	// 60s verification would write ~20 near-identical rows (see
	// docs/bank-account-verification.md §9 and §12).
	result := &VerifyResult{Status: StatusPending, ProviderJobID: "job-3"}
	p := &stubProvider{name: "test", pollResult: result}
	logger := &stubLogger{}
	reg := NewRegistry(p).WithLogger(logger)

	_, err := reg.Poll(context.Background(), "job-3")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(logger.records) != 0 {
		t.Fatalf("expected 0 log records for pending poll, got %d", len(logger.records))
	}
}

func TestRegistry_Poll_LogsErrorsEvenWithoutResult(t *testing.T) {
	// Transport errors count as terminal — they mean the poll definitively
	// did not observe a pending state, and auditors need to see them.
	p := &stubProvider{name: "test", pollErr: errors.New("connection refused")}
	logger := &stubLogger{}
	reg := NewRegistry(p).WithLogger(logger)

	_, err := reg.Poll(context.Background(), "job-4")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
	if len(logger.records) != 1 {
		t.Fatalf("expected 1 log record on poll error, got %d", len(logger.records))
	}
	if logger.records[0].Status != StatusFailed {
		t.Errorf("Status = %q, want failed", logger.records[0].Status)
	}
}

func TestRegistry_Poll_NotSupportedSurfaces(t *testing.T) {
	p := &stubProvider{name: "test", pollErr: ErrNotSupported}
	reg := NewRegistry(p)
	_, err := reg.Poll(context.Background(), "any")
	if !errors.Is(err, ErrNotSupported) {
		t.Fatalf("err = %v, want ErrNotSupported", err)
	}
}

func TestRegistry_Poll_NoProviderConfigured(t *testing.T) {
	var reg *Registry
	_, err := reg.Poll(context.Background(), "job")
	if !errors.Is(err, ErrProviderNotConfigured) {
		t.Fatalf("err = %v, want ErrProviderNotConfigured", err)
	}
}

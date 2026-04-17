package providers

import (
	"api/services/bav"
	"context"
	"errors"
	"testing"
	"time"
)

func TestMock_Name(t *testing.T) {
	if NewMock(MockConfig{}).Name() != "mock" {
		t.Fatal("Name() must be 'mock' for config matching")
	}
}

func TestMock_Sync_ReturnsCompleteImmediately(t *testing.T) {
	m := NewMock(MockConfig{Async: false})
	res, err := m.Verify(context.Background(), bav.VerifyRequest{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.Status != bav.StatusComplete {
		t.Errorf("Status = %q, want complete", res.Status)
	}
	if !res.Verified {
		t.Error("sync mock should report verified")
	}
	if res.ProviderJobID != "" {
		t.Errorf("sync mock should not set ProviderJobID; got %q", res.ProviderJobID)
	}
}

func TestMock_Async_Verify_ReturnsPendingWithJobID(t *testing.T) {
	m := NewMock(MockConfig{Async: true, ResolveAfter: time.Hour})
	res, err := m.Verify(context.Background(), bav.VerifyRequest{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.Status != bav.StatusPending {
		t.Errorf("Status = %q, want pending", res.Status)
	}
	if res.ProviderJobID == "" {
		t.Error("async mock must populate ProviderJobID")
	}
	if res.Verified {
		t.Error("pending result must not be marked verified")
	}
}

func TestMock_Async_Poll_ResolvesAfterDelay(t *testing.T) {
	// Use a tiny delay so the test stays fast.
	m := NewMock(MockConfig{Async: true, ResolveAfter: 20 * time.Millisecond})

	// First Verify dispatches the job.
	dispatch, err := m.Verify(context.Background(), bav.VerifyRequest{})
	if err != nil {
		t.Fatalf("verify err: %v", err)
	}
	jobID := dispatch.ProviderJobID

	// Immediate poll should still be pending.
	early, err := m.Poll(context.Background(), jobID)
	if err != nil {
		t.Fatalf("early poll err: %v", err)
	}
	if early.Status != bav.StatusPending {
		t.Errorf("early Status = %q, want pending", early.Status)
	}

	// Wait for the synthetic delay, then poll again.
	time.Sleep(30 * time.Millisecond)
	late, err := m.Poll(context.Background(), jobID)
	if err != nil {
		t.Fatalf("late poll err: %v", err)
	}
	if late.Status != bav.StatusComplete {
		t.Errorf("late Status = %q, want complete", late.Status)
	}
	if !late.Verified {
		t.Error("resolved async mock should report verified")
	}

	// Job should be evicted from the mock's in-memory store after resolving.
	_, err = m.Poll(context.Background(), jobID)
	if !errors.Is(err, bav.ErrInvalidInput) {
		t.Errorf("poll after resolve should report unknown job; got err=%v", err)
	}
}

func TestMock_Async_Poll_UnknownJobIDIsInvalidInput(t *testing.T) {
	m := NewMock(MockConfig{Async: true})
	_, err := m.Poll(context.Background(), "does-not-exist")
	if !errors.Is(err, bav.ErrInvalidInput) {
		t.Fatalf("err = %v, want ErrInvalidInput", err)
	}
}

func TestVerifyNow_Poll_ReturnsNotSupported(t *testing.T) {
	a := NewVerifyNow(VerifyNowConfig{APIKey: "k"})
	_, err := a.Poll(context.Background(), "any-job-id")
	if !errors.Is(err, bav.ErrNotSupported) {
		t.Fatalf("err = %v, want ErrNotSupported", err)
	}
}

func TestNewRegistry_MockProvider(t *testing.T) {
	reg, err := NewRegistry(Config{Provider: "mock", MockAsync: true})
	if err != nil {
		t.Fatalf("NewRegistry err: %v", err)
	}
	if reg.Active().Name() != "mock" {
		t.Fatalf("active provider = %q, want mock", reg.Active().Name())
	}
	res, err := reg.Active().Verify(context.Background(), bav.VerifyRequest{})
	if err != nil {
		t.Fatalf("verify err: %v", err)
	}
	if res.Status != bav.StatusPending {
		t.Errorf("expected pending with MockAsync=true; got %q", res.Status)
	}
}

package uwvendor

import (
	"testing"
)

func TestIsValidKind(t *testing.T) {
	for _, k := range AllKinds {
		if !IsValidKind(k) {
			t.Errorf("AllKinds member %s rejected by IsValidKind", k)
		}
	}
	if IsValidKind("not_a_kind") {
		t.Errorf("unknown kind should not be valid")
	}
}

func TestDeriveRequestPayloadHash_Stable(t *testing.T) {
	a := DeriveRequestPayloadHash(SubmitRequest{
		Kind: KindPathology, CaseID: 1, Subject: "MEMBER-1",
		Metadata: map[string]any{"panel": "full_blood", "fasting": true},
	})
	b := DeriveRequestPayloadHash(SubmitRequest{
		Kind: KindPathology, CaseID: 1, Subject: "MEMBER-1",
		Metadata: map[string]any{"fasting": true, "panel": "full_blood"}, // different map order
	})
	if a != b {
		t.Errorf("hash should be stable across map-key order: %s vs %s", a, b)
	}
}

func TestDeriveRequestPayloadHash_DiffersOnInput(t *testing.T) {
	a := DeriveRequestPayloadHash(SubmitRequest{Kind: KindPathology, CaseID: 1, Subject: "M-1"})
	b := DeriveRequestPayloadHash(SubmitRequest{Kind: KindPathology, CaseID: 2, Subject: "M-1"})
	if a == b {
		t.Errorf("different CaseID should produce different hashes")
	}
}

func TestDeriveWebhookIdempotencyKey_IncludesBody(t *testing.T) {
	a := DeriveWebhookIdempotencyKey("mock", "mock-1", []byte("{}"))
	b := DeriveWebhookIdempotencyKey("mock", "mock-1", []byte("{\"x\":1}"))
	if a == b {
		t.Errorf("different bodies should produce different keys")
	}
}

func TestDeriveWebhookIdempotencyKey_DeterministicForSameInputs(t *testing.T) {
	a := DeriveWebhookIdempotencyKey("mock", "mock-1", []byte(`{"a":1}`))
	b := DeriveWebhookIdempotencyKey("mock", "mock-1", []byte(`{"a":1}`))
	if a != b {
		t.Errorf("same inputs should produce same key")
	}
}

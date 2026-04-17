// Package bav defines the provider-agnostic canonical model for Bank Account
// Verification. Adapters for specific providers (VerifyNow, LexisNexis, etc.)
// live under services/bav/providers and translate their wire formats to and
// from the types declared here.
package bav

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// TriState represents a three-valued result where "unknown" is distinct from
// both yes and no. BAV providers frequently return "cannot determine" for
// individual checks (e.g. accountOpen) and the UI needs to differentiate that
// from an outright "no".
type TriState string

const (
	TriUnknown TriState = "unknown"
	TriYes     TriState = "yes"
	TriNo      TriState = "no"
)

// Status describes the lifecycle state of a verification call. Synchronous
// providers return StatusComplete or StatusFailed; async providers may return
// StatusPending along with a ProviderJobID for later polling.
type Status string

const (
	StatusComplete Status = "complete"
	StatusPending  Status = "pending"
	StatusFailed   Status = "failed"
)

// VerifyRequest is the canonical input for a bank account verification call.
// Adapters are responsible for mapping these fields to their provider's wire
// format.
//
// ClaimID is optional because the current UX verifies banking details before
// a claim is saved. Attempt defaults to 1 on the server; callers increment
// it when the user explicitly retries after a failure so each retry is
// distinguishable in the audit log.
type VerifyRequest struct {
	FirstName         string
	Surname           string
	IdentityNumber    string
	IdentityType      string
	BankAccountNumber string
	BankBranchCode    string
	BankAccountType   string
	ClaimID           *int
	Attempt           int
	IdempotencyKey    string
}

// DeriveIdempotencyKey produces a stable sha256 key so identical retries
// dedup against the same log row. Includes identifying banking fields so
// the key changes when the user edits account or branch before retrying.
func (r VerifyRequest) DeriveIdempotencyKey(provider string) string {
	claimID := ""
	if r.ClaimID != nil {
		claimID = strconv.Itoa(*r.ClaimID)
	}
	attempt := r.Attempt
	if attempt <= 0 {
		attempt = 1
	}
	h := sha256.New()
	_, _ = fmt.Fprintf(
		h,
		"%s|%d|%s|%s|%s|%s",
		claimID,
		attempt,
		provider,
		r.IdentityNumber,
		r.BankAccountNumber,
		r.BankBranchCode,
	)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyResult is the canonical output of a verification call. The granular
// TriState fields mirror the BankservAfrica AVS contract that virtually every
// South African BAV provider front-ends. JSON tags define the v2 canonical
// wire shape; ProviderStatusText and RawPayload are internal-only.
type VerifyResult struct {
	Status            Status   `json:"status"`
	Verified          bool     `json:"verified"`
	Summary           string   `json:"summary"`
	AccountFound      TriState `json:"accountFound"`
	AccountOpen       TriState `json:"accountOpen"`
	IdentityMatch     TriState `json:"identityMatch"`
	AccountTypeMatch  TriState `json:"accountTypeMatch"`
	AcceptsCredits    TriState `json:"acceptsCredits"`
	AcceptsDebits     TriState `json:"acceptsDebits"`
	Provider          string   `json:"provider"`
	ProviderRequestID string   `json:"providerRequestId"`
	ProviderJobID     string   `json:"providerJobId,omitempty"`
	// ProviderStatusText is the provider's own free-form status label (e.g.
	// VerifyNow's "Verified"/"Declined"). Preserved so the legacy v1 wire
	// shape can pass it through untouched; not part of the v2 canonical shape.
	ProviderStatusText string `json:"-"`
	RawPayload         []byte `json:"-"`
}

// ParseTriState normalises the wide range of values providers use to express
// yes/no/unknown into the canonical TriState. Unrecognised or empty input
// resolves to TriUnknown so that downstream consumers always get one of the
// three defined values.
func ParseTriState(v any) TriState {
	switch t := v.(type) {
	case nil:
		return TriUnknown
	case bool:
		if t {
			return TriYes
		}
		return TriNo
	case TriState:
		return t
	case string:
		switch strings.ToLower(strings.TrimSpace(t)) {
		case "y", "yes", "true", "t", "1":
			return TriYes
		case "n", "no", "false", "f", "0":
			return TriNo
		}
	}
	return TriUnknown
}

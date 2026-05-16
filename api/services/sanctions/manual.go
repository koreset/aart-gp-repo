package sanctions

import "context"

// manualProvider is the bootstrap implementation that ships with Phase 3.
// Every call to Screen returns StatusPending — finance is expected to record
// the actual outcome (clear or hit) through the screening service's
// RecordOutcome path, not via automatic provider feedback.
//
// When a real provider lands (LexisNexis, Refinitiv, etc.), it lives in its
// own file in this package and is selected via Use() at startup based on
// configuration.
type manualProvider struct{}

// NewManual returns the manual / human-driven provider.
func NewManual() Provider { return manualProvider{} }

func (manualProvider) Name() string { return "manual" }

func (manualProvider) Screen(_ context.Context, _ Subject) (*Result, error) {
	return &Result{Status: StatusPending}, nil
}

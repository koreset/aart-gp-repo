package services

import (
	"api/log"
	"context"
)

// BordereauxProgressEvent is pushed via WebSocket during a synchronous
// GenerateBordereaux call so the UI can show phase + progress while the HTTP
// request is still in flight. See the UI subscription in App.vue / the
// progress dialog in BordereauxGenerationForm.vue.
type BordereauxProgressEvent struct {
	Type     string `json:"type"`              // member | premium | claim
	Phase    string `json:"phase"`             // start | fetching | writing | zipping | completed | failed
	Progress int    `json:"progress"`          // 0..100
	Message  string `json:"message,omitempty"` // human-readable, shown under the bar
	Scheme   string `json:"scheme,omitempty"`  // current scheme when looping per-scheme
}

// sendBordereauxProgress delivers a progress event to the user who initiated
// the generation. No-op if the WS hub is not initialised (tests, CLI) or the
// user email cannot be resolved from context.
func sendBordereauxProgress(ctx context.Context, evt BordereauxProgressEvent) {
	email, _ := ctx.Value(log.UserEmailKey).(string)
	if email == "" {
		return
	}
	hub := GetHub()
	if hub == nil {
		return
	}
	hub.SendToUser(email, WSEnvelope{
		Type:    WSBordereauxProgress,
		Payload: evt,
	})
}

// RIValidationProgressEvent reports the L1 / L2 / L3 validation pipeline
// phases for an RI bordereaux run so the UI can animate a multi-step progress
// indicator.
type RIValidationProgressEvent struct {
	RunID    string `json:"run_id"`
	Level    int    `json:"level"`             // 1..3
	Phase    string `json:"phase"`             // start | level_complete | completed | failed
	Progress int    `json:"progress"`          // 0..100
	Findings int    `json:"findings"`          // running total
	Message  string `json:"message,omitempty"`
}

// sendRIValidationProgress publishes a validation-pipeline progress event.
// Same delivery semantics as sendBordereauxProgress.
func sendRIValidationProgress(ctx context.Context, evt RIValidationProgressEvent) {
	email, _ := ctx.Value(log.UserEmailKey).(string)
	if email == "" {
		return
	}
	hub := GetHub()
	if hub == nil {
		return
	}
	hub.SendToUser(email, WSEnvelope{
		Type:    WSRIValidationProgress,
		Payload: evt,
	})
}

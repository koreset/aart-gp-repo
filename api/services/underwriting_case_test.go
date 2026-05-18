package services

import (
	"testing"

	"api/models"
)

func TestUnderwritingCaseTransitionMatrix(t *testing.T) {
	cases := []struct {
		from    models.UWCaseStatus
		to      models.UWCaseStatus
		allowed bool
	}{
		{models.UWCaseStatusPendingEvidence, models.UWCaseStatusInReview, true},
		{models.UWCaseStatusPendingEvidence, models.UWCaseStatusDeclined, true},
		{models.UWCaseStatusPendingEvidence, models.UWCaseStatusDecided, false},
		{models.UWCaseStatusInReview, models.UWCaseStatusDecided, true},
		{models.UWCaseStatusInReview, models.UWCaseStatusPendingEvidence, true},
		{models.UWCaseStatusDecided, models.UWCaseStatusInReview, false},
		{models.UWCaseStatusDecided, models.UWCaseStatusDeclined, false},
		{models.UWCaseStatusDeclined, models.UWCaseStatusInReview, false},
		{models.UWCaseStatusPostponed, models.UWCaseStatusInReview, true},
		{models.UWCaseStatusPostponed, models.UWCaseStatusPendingEvidence, true},
	}
	for _, tc := range cases {
		got := allowedUWCaseTransitions[tc.from][tc.to]
		if got != tc.allowed {
			t.Errorf("%s -> %s: got allowed=%v, want %v", tc.from, tc.to, got, tc.allowed)
		}
	}
}

func TestMemberLookupKey(t *testing.T) {
	if memberLookupKey("Jane Doe", "Standard") != "Jane Doe|Standard" {
		t.Errorf("unexpected key format")
	}
	if memberLookupKey("", "") != "|" {
		t.Errorf("empty inputs should still produce a stable key")
	}
}

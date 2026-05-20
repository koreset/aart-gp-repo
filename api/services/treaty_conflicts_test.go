package services

import (
	"errors"
	"strings"
	"testing"
)

func TestFormatConflictError(t *testing.T) {
	t.Run("empty slice returns nil", func(t *testing.T) {
		if err := FormatConflictError(nil); err != nil {
			t.Errorf("nil input: got error %v, want nil", err)
		}
		if err := FormatConflictError([]SchemeTreatyConflict{}); err != nil {
			t.Errorf("empty slice: got error %v, want nil", err)
		}
	})

	t.Run("single conflict wraps sentinel and names every part", func(t *testing.T) {
		err := FormatConflictError([]SchemeTreatyConflict{{
			SchemeID:       42,
			SchemeName:     "ABC Pension",
			LineOfBusiness: "group_life",
			TreatyType:     "proportional",
			Treaties: []ConflictParty{
				{TreatyID: 1, TreatyNumber: "T-001", TreatyName: "Old Treaty"},
			},
		}})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrSchemeTreatyConflict) {
			t.Errorf("errors.Is(err, ErrSchemeTreatyConflict) = false, want true")
		}
		msg := err.Error()
		for _, want := range []string{"42", "ABC Pension", "group_life", "proportional", "T-001"} {
			if !strings.Contains(msg, want) {
				t.Errorf("message missing %q\nfull message: %s", want, msg)
			}
		}
	})

	t.Run("multiple conflicts are joined with semicolons", func(t *testing.T) {
		err := FormatConflictError([]SchemeTreatyConflict{
			{SchemeID: 1, SchemeName: "A", LineOfBusiness: "group_life", TreatyType: "proportional",
				Treaties: []ConflictParty{{TreatyNumber: "T-001"}}},
			{SchemeID: 2, SchemeName: "B", LineOfBusiness: "group_life", TreatyType: "proportional",
				Treaties: []ConflictParty{{TreatyNumber: "T-002"}}},
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		msg := err.Error()
		if !strings.Contains(msg, "; ") {
			t.Errorf("expected semicolon separator, got: %s", msg)
		}
		if !strings.Contains(msg, "T-001") || !strings.Contains(msg, "T-002") {
			t.Errorf("expected both treaty numbers, got: %s", msg)
		}
	})

	t.Run("multiple treaties per conflict are joined with commas", func(t *testing.T) {
		err := FormatConflictError([]SchemeTreatyConflict{{
			SchemeID: 7, SchemeName: "C", LineOfBusiness: "funeral", TreatyType: "xl_risk",
			Treaties: []ConflictParty{
				{TreatyNumber: "T-100"},
				{TreatyNumber: "T-101"},
				{TreatyNumber: "T-102"},
			},
		}})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		msg := err.Error()
		if !strings.Contains(msg, "T-100, T-101, T-102") {
			t.Errorf("expected comma-joined treaty numbers, got: %s", msg)
		}
	})

	t.Run("errors.Is unwraps through the joined message", func(t *testing.T) {
		err := FormatConflictError([]SchemeTreatyConflict{{
			SchemeID: 1, SchemeName: "X", LineOfBusiness: "group_life", TreatyType: "proportional",
			Treaties: []ConflictParty{{TreatyNumber: "T-001"}},
		}})
		if !errors.Is(err, ErrSchemeTreatyConflict) {
			t.Errorf("errors.Is must unwrap through fmt.Errorf %%w")
		}
	})
}

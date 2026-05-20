package services

import (
	"testing"
	"time"

	"api/models"
)

func TestMapBenefitToLineOfBusiness(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"GLA", "group_life"},
		{"SGLA", "group_life"},
		{"gla", "group_life"},
		{"  gla  ", "group_life"},
		{"PTD", "group_disability"},
		{"CI", "group_disability"},
		{"TTD", "group_disability"},
		{"PHI", "group_health"},
		{"FUNERAL", "funeral"},
		{"SPOUSE_FUNERAL", "funeral"},
		{"CHILD_FUNERAL", "funeral"},
		{"PARENT_FUNERAL", "funeral"},
		{"EXTENDED_FUNERAL", "funeral"},
		{"unknown", ""},
		{"", ""},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := mapBenefitToLineOfBusiness(tc.input)
			if got != tc.want {
				t.Errorf("mapBenefitToLineOfBusiness(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseDateOnly(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantOK    bool
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		{"valid date", "2024-03-15", true, 2024, time.March, 15},
		{"whitespace-padded", "  2024-03-15  ", true, 2024, time.March, 15},
		{"empty", "", false, 0, 0, 0},
		{"all whitespace", "   ", false, 0, 0, 0},
		{"malformed", "not-a-date", false, 0, 0, 0},
		{"wrong layout", "15/03/2024", false, 0, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := parseDateOnly(tc.input)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if got.Year() != tc.wantYear || got.Month() != tc.wantMonth || got.Day() != tc.wantDay {
				t.Errorf("parsed = %v, want %d-%02d-%02d", got, tc.wantYear, tc.wantMonth, tc.wantDay)
			}
		})
	}
}

func TestCoversDate(t *testing.T) {
	mkTreaty := func(eff, exp string) models.ReinsuranceTreaty {
		return models.ReinsuranceTreaty{EffectiveDate: eff, ExpiryDate: exp}
	}
	mustDate := func(s string) time.Time {
		d, ok := parseDateOnly(s)
		if !ok {
			t.Fatalf("test setup: bad date %q", s)
		}
		return d
	}
	cases := []struct {
		name   string
		treaty models.ReinsuranceTreaty
		date   string
		want   bool
	}{
		{"inside the window", mkTreaty("2024-01-01", "2024-12-31"), "2024-06-15", true},
		{"equal to effective is inclusive", mkTreaty("2024-01-01", "2024-12-31"), "2024-01-01", true},
		{"equal to expiry is inclusive", mkTreaty("2024-01-01", "2024-12-31"), "2024-12-31", true},
		{"before effective", mkTreaty("2024-01-01", "2024-12-31"), "2023-12-31", false},
		{"after expiry", mkTreaty("2024-01-01", "2024-12-31"), "2025-01-01", false},
		{"empty expiry is open-ended", mkTreaty("2024-01-01", ""), "2099-01-01", true},
		{"empty effective excludes everything", mkTreaty("", "2024-12-31"), "2024-06-15", false},
		{"malformed expiry treated as open-ended", mkTreaty("2024-01-01", "garbage"), "2099-01-01", true},
		{"whitespace-only expiry is open-ended", mkTreaty("2024-01-01", "   "), "2099-01-01", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := coversDate(tc.treaty, mustDate(tc.date))
			if got != tc.want {
				t.Errorf("coversDate(%+v, %s) = %v, want %v", tc.treaty, tc.date, got, tc.want)
			}
		})
	}
}

func TestCoverageDateFor(t *testing.T) {
	member, _ := parseDateOnly("2020-01-15")
	claim, _ := parseDateOnly("2024-06-01")

	cases := []struct {
		name     string
		basis    string
		want     time.Time
		wantKind string
	}{
		{"loss_occurring uses claim date", "loss_occurring", claim, "claim"},
		{"LOSS_OCCURRING is case-insensitive", "LOSS_OCCURRING", claim, "claim"},
		{"risk_attaching uses member entry", "risk_attaching", member, "member"},
		{"empty basis defaults to member entry", "", member, "member"},
		{"unknown basis defaults to member entry", "weird_value", member, "member"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t1 := models.ReinsuranceTreaty{TreatyBasis: tc.basis}
			got := coverageDateFor(t1, member, claim)
			if !got.Equal(tc.want) {
				t.Errorf("got %v (%s), want %v", got, tc.wantKind, tc.want)
			}
		})
	}
}

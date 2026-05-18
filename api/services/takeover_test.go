package services

import (
	"testing"

	"api/models"
)

func TestDefaultOutcomeFor(t *testing.T) {
	tests := []struct {
		name string
		pm   models.PriorInsurerMember
		want string
	}{
		{
			name: "not in force → new evidence",
			pm:   models.PriorInsurerMember{InForce: false},
			want: TakeoverOutcomeNewEvidenceRequired,
		},
		{
			name: "in force, clean → continuation_no_evidence",
			pm: models.PriorInsurerMember{
				InForce: true, PriorLoadings: "{}", PriorExclusions: "[]",
			},
			want: TakeoverOutcomeContinuationNoEvidence,
		},
		{
			name: "in force with loading → continuation_with_loading",
			pm: models.PriorInsurerMember{
				InForce: true, PriorLoadings: `{"gla":25}`, PriorExclusions: "[]",
			},
			want: TakeoverOutcomeContinuationWithLoading,
		},
		{
			name: "in force with exclusion → continuation_with_loading",
			pm: models.PriorInsurerMember{
				InForce: true, PriorLoadings: "{}", PriorExclusions: `["smoker"]`,
			},
			want: TakeoverOutcomeContinuationWithLoading,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := defaultOutcomeFor(&tc.pm)
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestParseBoolTrue(t *testing.T) {
	cases := map[string]bool{
		"true":     true,
		"TRUE":     true,
		"Yes":      true,
		"y":        true,
		"1":        true,
		"in_force": true,
		"false":    false,
		"":         false,
		"maybe":    false,
	}
	for in, want := range cases {
		if got := parseBoolTrue(in); got != want {
			t.Errorf("parseBoolTrue(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestParseBenefitMap(t *testing.T) {
	got := parseBenefitMap("gla:25 | ptd:10|ci:0")
	if got["gla"] != 25 || got["ptd"] != 10 || got["ci"] != 0 {
		t.Errorf("benefit map: %v", got)
	}
	empty := parseBenefitMap("")
	if len(empty) != 0 {
		t.Errorf("empty input should produce empty map, got %v", empty)
	}
	garbage := parseBenefitMap("not_a_kv | also:bad:format")
	if len(garbage) != 1 {
		t.Errorf("garbage should produce 1 entry (also: 0), got %v", garbage)
	}
}

func TestParseExclusionList(t *testing.T) {
	got := parseExclusionList("Smoker | Diabetes |")
	if len(got) != 2 {
		t.Fatalf("want 2 exclusions, got %d (%v)", len(got), got)
	}
	if got[0] != "smoker" || got[1] != "diabetes" {
		t.Errorf("lowercase+trim failed: %v", got)
	}
}

func TestMaxLoadingFromPrior(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{`{"gla":25,"ptd":10}`, 25},
		{`{"ci":50}`, 50},
		{`{}`, 0},
		{"", 0},
		{"not json", 0},
	}
	for _, tc := range tests {
		got := maxLoadingFromPrior(tc.in)
		if got != tc.want {
			t.Errorf("maxLoadingFromPrior(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestBuildPriorMemberContext_KeysAndValues(t *testing.T) {
	ctx := buildPriorMemberContext(models.PriorInsurerMember{
		InForce:         true,
		GlaSumAssured:   500_000,
		PtdSumAssured:   200_000,
		CiSumAssured:    100_000,
		PriorLoadings:   `{"gla":25,"ptd":10}`,
		PriorExclusions: `["smoker","diabetes"]`,
	})
	if ctx["in_force"] != true {
		t.Errorf("in_force missing or wrong: %v", ctx["in_force"])
	}
	if ctx["prior_gla_sa"] != 500_000.0 {
		t.Errorf("prior_gla_sa: %v", ctx["prior_gla_sa"])
	}
	if ctx["prior_loading_gla"] != 25.0 {
		t.Errorf("prior_loading_gla: %v", ctx["prior_loading_gla"])
	}
	if ctx["prior_exclusion_smoker"] != true {
		t.Errorf("prior_exclusion_smoker: %v", ctx["prior_exclusion_smoker"])
	}
	if ctx["prior_exclusion_diabetes"] != true {
		t.Errorf("prior_exclusion_diabetes: %v", ctx["prior_exclusion_diabetes"])
	}
}

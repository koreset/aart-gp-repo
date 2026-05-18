package services

import (
	"math"
	"strings"
	"testing"
	"time"
)

func TestComputeBMI(t *testing.T) {
	tests := []struct {
		name   string
		height float64
		weight float64
		want   float64
	}{
		{"normal", 180, 75, 75 / (1.8 * 1.8)},
		{"underweight inputs zero", 0, 75, 0},
		{"weight zero", 180, 0, 0},
		{"negative ignored", -180, 75, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := computeBMI(tc.height, tc.weight)
			if math.Abs(got-tc.want) > 1e-6 {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestSignatureHash_Deterministic(t *testing.T) {
	ts := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	a := signatureHash("Jane Doe", ts, "10.0.0.1")
	b := signatureHash("Jane Doe", ts, "10.0.0.1")
	if a != b {
		t.Errorf("hash should be deterministic, got %s and %s", a, b)
	}
	if len(a) != 64 {
		t.Errorf("SHA-256 hex should be 64 chars, got %d", len(a))
	}
}

func TestSignatureHash_DifferentInputs(t *testing.T) {
	ts := time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC)
	a := signatureHash("Jane Doe", ts, "10.0.0.1")
	b := signatureHash("Jane Doe", ts, "10.0.0.2") // different IP
	if a == b {
		t.Errorf("different IP should produce different hash")
	}
}

func TestMergeDisclosureIntoContext_BMIAndSmoker(t *testing.T) {
	ctx := EvaluationContext{}
	mergeDisclosureIntoContext(ctx, DisclosureSubmission{
		Height:           180,
		Weight:           120,
		Smoker:           true,
		CigarettesPerDay: 15,
	}, 37.0)
	if ctx["bmi"] != 37.0 {
		t.Errorf("bmi not set, got %v", ctx["bmi"])
	}
	if ctx["smoker"] != true {
		t.Errorf("smoker not set, got %v", ctx["smoker"])
	}
	if ctx["cigarettes_per_day"] != 15.0 {
		t.Errorf("cigarettes not float-encoded, got %v", ctx["cigarettes_per_day"])
	}
}

func TestMergeDisclosureIntoContext_ConditionsBecomeBooleans(t *testing.T) {
	ctx := EvaluationContext{}
	mergeDisclosureIntoContext(ctx, DisclosureSubmission{
		DisclosedConditions: []string{"Diabetes_Type2", "  hypertension  ", ""},
	}, 0)
	if ctx["condition_diabetes_type2"] != true {
		t.Errorf("condition_diabetes_type2 should be true, got %v", ctx["condition_diabetes_type2"])
	}
	if ctx["condition_hypertension"] != true {
		t.Errorf("trimmed condition_hypertension should be true, got %v", ctx["condition_hypertension"])
	}
	if _, present := ctx["condition_"]; present {
		t.Errorf("empty code should be skipped")
	}
}

func TestMergeDisclosureIntoContext_OccupationAnswersFlattened(t *testing.T) {
	ctx := EvaluationContext{}
	mergeDisclosureIntoContext(ctx, DisclosureSubmission{
		OccupationRiskAnswers: map[string]any{
			"WorksAtHeight": true,
			"HoursPerWeek":  45.0,
		},
	}, 0)
	if ctx["occ_worksatheight"] != true {
		t.Errorf("occ_worksatheight not set: %v", ctx["occ_worksatheight"])
	}
	if ctx["occ_hoursperweek"] != 45.0 {
		t.Errorf("occ_hoursperweek not set: %v", ctx["occ_hoursperweek"])
	}
}

func TestMergeDisclosureIntoContext_SkipsZeroNumerics(t *testing.T) {
	ctx := EvaluationContext{}
	mergeDisclosureIntoContext(ctx, DisclosureSubmission{
		Height:              0,
		Weight:              0,
		CigarettesPerDay:    0,
		AlcoholUnitsPerWeek: 0,
	}, 0)
	for _, k := range []string{"bmi", "height_cm", "weight_kg", "cigarettes_per_day", "alcohol_units_per_week"} {
		if _, present := ctx[k]; present {
			t.Errorf("zero numeric %s should be omitted, got %v", k, ctx[k])
		}
	}
	// smoker and hazardous_hobbies are booleans and always set (even false)
	// so rules can use op=eq false. Verify they are present:
	if v, present := ctx["smoker"]; !present || v != false {
		t.Errorf("smoker should be present and false, got %v (present=%v)", v, present)
	}
}

func TestDisclosureContext_ContainsAllExpectedKeys(t *testing.T) {
	// Smoke-test: a typical full disclosure produces context keys with the
	// expected shape so rule authors can target them.
	ctx := EvaluationContext{}
	mergeDisclosureIntoContext(ctx, DisclosureSubmission{
		Height:              180,
		Weight:              90,
		Smoker:              true,
		CigarettesPerDay:    10,
		AlcoholUnitsPerWeek: 7,
		HasHazardousHobbies: true,
		DisclosedConditions: []string{"asthma"},
		OccupationRiskAnswers: map[string]any{
			"WorksAtHeight": true,
		},
	}, 27.78)
	expected := []string{
		"bmi", "height_cm", "weight_kg",
		"smoker", "cigarettes_per_day", "alcohol_units_per_week",
		"hazardous_hobbies",
		"condition_asthma",
		"occ_worksatheight",
	}
	for _, k := range expected {
		if _, ok := ctx[k]; !ok {
			t.Errorf("missing key %q in context: %v", k, ctxKeys(ctx))
		}
	}
}

func ctxKeys(ctx EvaluationContext) string {
	keys := make([]string, 0, len(ctx))
	for k := range ctx {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}

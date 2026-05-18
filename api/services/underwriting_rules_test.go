package services

import (
	"strings"
	"testing"

	"api/models"
)

func ruleEq(field string, valueA string, outcome string, loading float64) models.UWRule {
	return models.UWRule{
		Field:         field,
		Op:            UWOpEq,
		ConditionJSON: `{"value_a":` + valueA + `}`,
		Outcome:       outcome,
		LoadingPercent: loading,
	}
}

func TestEvaluate_AllOperators(t *testing.T) {
	tests := []struct {
		name     string
		rule     models.UWRule
		ctx      EvaluationContext
		matches  bool
	}{
		{"eq numeric match", ruleEq("age", "50", UWRuleOutcomeRefer, 10), EvaluationContext{"age": 50.0}, true},
		{"eq numeric mismatch", ruleEq("age", "50", UWRuleOutcomeRefer, 10), EvaluationContext{"age": 49.0}, false},
		{"eq bool match", ruleEq("smoker", "true", UWRuleOutcomeRefer, 15), EvaluationContext{"smoker": true}, true},
		{"eq string match", ruleEq("occupation", `"miner"`, UWRuleOutcomeRefer, 20), EvaluationContext{"occupation": "miner"}, true},
		{"ne mismatch (matches)", models.UWRule{Field: "smoker", Op: UWOpNe, ConditionJSON: `{"value_a":true}`, Outcome: UWRuleOutcomeAccept}, EvaluationContext{"smoker": false}, true},
		{"gt true", models.UWRule{Field: "bmi", Op: UWOpGt, ConditionJSON: `{"value_a":30}`}, EvaluationContext{"bmi": 31.0}, true},
		{"gt false at boundary", models.UWRule{Field: "bmi", Op: UWOpGt, ConditionJSON: `{"value_a":30}`}, EvaluationContext{"bmi": 30.0}, false},
		{"gte true at boundary", models.UWRule{Field: "bmi", Op: UWOpGte, ConditionJSON: `{"value_a":30}`}, EvaluationContext{"bmi": 30.0}, true},
		{"lt true", models.UWRule{Field: "age", Op: UWOpLt, ConditionJSON: `{"value_a":18}`}, EvaluationContext{"age": 17.0}, true},
		{"lte boundary", models.UWRule{Field: "age", Op: UWOpLte, ConditionJSON: `{"value_a":18}`}, EvaluationContext{"age": 18.0}, true},
		{"between true inside", models.UWRule{Field: "bmi", Op: UWOpBetween, ConditionJSON: `{"value_a":30,"value_b":35}`}, EvaluationContext{"bmi": 32.0}, true},
		{"between true upper boundary", models.UWRule{Field: "bmi", Op: UWOpBetween, ConditionJSON: `{"value_a":30,"value_b":35}`}, EvaluationContext{"bmi": 35.0}, true},
		{"between false above", models.UWRule{Field: "bmi", Op: UWOpBetween, ConditionJSON: `{"value_a":30,"value_b":35}`}, EvaluationContext{"bmi": 36.0}, false},
		{"in numeric match", models.UWRule{Field: "occupation_class", Op: UWOpIn, ConditionJSON: `{"values":[3,4]}`}, EvaluationContext{"occupation_class": 3.0}, true},
		{"in numeric mismatch", models.UWRule{Field: "occupation_class", Op: UWOpIn, ConditionJSON: `{"values":[3,4]}`}, EvaluationContext{"occupation_class": 2.0}, false},
		{"in string match", models.UWRule{Field: "condition", Op: UWOpIn, ConditionJSON: `{"values":["diabetes_type2","hypertension"]}`}, EvaluationContext{"condition": "diabetes_type2"}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := Evaluate([]models.UWRule{tc.rule}, tc.ctx)
			matched := len(out) == 1
			if matched != tc.matches {
				t.Errorf("got matched=%v, want %v (rule=%+v)", matched, tc.matches, tc.rule)
			}
		})
	}
}

func TestEvaluate_FieldMissingIsSkipped(t *testing.T) {
	rule := ruleEq("bmi", "30", UWRuleOutcomeRefer, 10)
	out := Evaluate([]models.UWRule{rule}, EvaluationContext{"age": 40.0})
	if len(out) != 0 {
		t.Errorf("missing field should produce no outcomes, got %d", len(out))
	}
}

func TestEvaluate_MalformedConditionIsSkipped(t *testing.T) {
	rule := models.UWRule{Field: "bmi", Op: UWOpEq, ConditionJSON: `not json`}
	out := Evaluate([]models.UWRule{rule}, EvaluationContext{"bmi": 30.0})
	if len(out) != 0 {
		t.Errorf("malformed rule should be skipped, got %d outcomes", len(out))
	}
}

func TestEvaluate_PriorityOrdering(t *testing.T) {
	rules := []models.UWRule{
		{ID: 1, Field: "bmi", Op: UWOpGte, ConditionJSON: `{"value_a":30}`, Outcome: UWRuleOutcomeRefer, Priority: 20},
		{ID: 2, Field: "bmi", Op: UWOpGt, ConditionJSON: `{"value_a":40}`, Outcome: UWRuleOutcomeDecline, Priority: 10},
	}
	out := Evaluate(rules, EvaluationContext{"bmi": 45.0})
	if len(out) != 2 {
		t.Fatalf("want 2 matches, got %d", len(out))
	}
	if out[0].Priority > out[1].Priority {
		t.Errorf("expected priority-asc order, got %d before %d", out[0].Priority, out[1].Priority)
	}
}

func TestSummarise_StrictestOutcomeWins(t *testing.T) {
	outcomes := []EvaluatedOutcome{
		{Outcome: UWRuleOutcomeAccept, LoadingPercent: 5},
		{Outcome: UWRuleOutcomeRefer, LoadingPercent: 25, ExclusionCode: "smoker"},
		{Outcome: UWRuleOutcomeDecline, LoadingPercent: 100},
	}
	s := summarise(outcomes)
	if s.Outcome != UWRuleOutcomeDecline {
		t.Errorf("strictest outcome should be decline, got %s", s.Outcome)
	}
	if s.MaxLoading != 100 {
		t.Errorf("max loading should be 100, got %v", s.MaxLoading)
	}
	if len(s.Exclusions) != 1 || s.Exclusions[0] != "smoker" {
		t.Errorf("exclusions should be [smoker], got %v", s.Exclusions)
	}
}

func TestSummarise_ExclusionDedup(t *testing.T) {
	outcomes := []EvaluatedOutcome{
		{Outcome: UWRuleOutcomeRefer, ExclusionCode: "smoker"},
		{Outcome: UWRuleOutcomeRefer, ExclusionCode: "smoker"},
		{Outcome: UWRuleOutcomeRefer, ExclusionCode: "obese"},
	}
	s := summarise(outcomes)
	if len(s.Exclusions) != 2 {
		t.Errorf("expected 2 unique exclusions, got %d (%v)", len(s.Exclusions), s.Exclusions)
	}
}

func TestBuildConditionFromCSV(t *testing.T) {
	cases := []struct {
		name    string
		op      string
		a, b, v string
		want    string
		errSubstr string
	}{
		{"eq numeric", UWOpEq, "30", "", "", `"value_a":30`, ""},
		{"eq bool", UWOpEq, "true", "", "", `"value_a":true`, ""},
		{"eq string", UWOpEq, "miner", "", "", `"value_a":"miner"`, ""},
		{"between", UWOpBetween, "30", "35", "", `"value_a":30,"value_b":35`, ""},
		{"between missing b", UWOpBetween, "30", "", "", "", "value_b"},
		{"in numeric pipes", UWOpIn, "", "", "3|4", `"values":[3,4]`, ""},
		{"in empty", UWOpIn, "", "", "", "", "non-empty"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildConditionFromCSV(tc.op, tc.a, tc.b, tc.v)
			if tc.errSubstr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errSubstr) {
					t.Errorf("expected error containing %q, got %v", tc.errSubstr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(got, tc.want) {
				t.Errorf("output %s does not contain %q", got, tc.want)
			}
		})
	}
}

func TestCsvValueToScalar(t *testing.T) {
	if v := csvValueToScalar("true"); v != true {
		t.Errorf("got %v", v)
	}
	if v := csvValueToScalar("42.5"); v != 42.5 {
		t.Errorf("got %v", v)
	}
	if v := csvValueToScalar("hello"); v != "hello" {
		t.Errorf("got %v", v)
	}
}

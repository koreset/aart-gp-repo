package services

import (
	"testing"

	"api/models"
)

func TestSummaryColumnsFor_KnownTableTypes(t *testing.T) {
	cases := []struct {
		tableType string
		minCols   int // sanity floor; bumping the projection is fine
	}{
		{"member_rating_results", 10},
		{"member_data", 8},
		{"member_premium_schedules", 5},
	}
	for _, tc := range cases {
		t.Run(tc.tableType, func(t *testing.T) {
			got := summaryColumnsFor(tc.tableType)
			if len(got) < tc.minCols {
				t.Errorf("summary cols for %s = %d, want at least %d (%v)",
					tc.tableType, len(got), tc.minCols, got)
			}
		})
	}
}

func TestSummaryColumnsFor_UnknownReturnsNil(t *testing.T) {
	if got := summaryColumnsFor("not_a_table"); got != nil {
		t.Errorf("unknown table should produce nil, got %v", got)
	}
}

func TestProjectStructToMap_KeepsOnlyRequestedJSONTags(t *testing.T) {
	row := models.MemberRatingResult{
		FinancialYear: 2026,
		Category:      "Provident",
		MemberName:    "A Buckle",
		Gender:        "M",
		AnnualSalary:  450_000,
		// Many other fields would be zero — projector must NOT include them.
	}
	wanted := []string{"financial_year", "category", "member_name", "gender", "annual_salary"}
	out := projectStructToMap(row, wanted)
	if len(out) != len(wanted) {
		t.Fatalf("projected map should have exactly %d keys, got %d (%v)",
			len(wanted), len(out), out)
	}
	if out["financial_year"] != 2026 {
		t.Errorf("financial_year: %v", out["financial_year"])
	}
	if out["category"] != "Provident" {
		t.Errorf("category: %v", out["category"])
	}
	if _, present := out["gla_sum_assured"]; present {
		t.Errorf("non-wanted field gla_sum_assured leaked into projected map")
	}
}

func TestProjectStructToMap_NilForNonStruct(t *testing.T) {
	if projectStructToMap("not a struct", []string{"x"}) != nil {
		t.Errorf("non-struct input should return nil")
	}
}

func TestProjectSliceToMaps_RoundTrip(t *testing.T) {
	rows := []models.MemberRatingResult{
		{FinancialYear: 2026, Category: "Cat A", MemberName: "Alice"},
		{FinancialYear: 2026, Category: "Cat B", MemberName: "Bob"},
		{FinancialYear: 2026, Category: "Cat A", MemberName: "Carol"},
	}
	wanted := []string{"financial_year", "category", "member_name"}
	maps := projectSliceToMaps(rows, wanted)
	if len(maps) != len(rows) {
		t.Fatalf("expected %d maps, got %d", len(rows), len(maps))
	}
	for i, m := range maps {
		if m["member_name"] != rows[i].MemberName {
			t.Errorf("row %d member_name: got %v want %s", i, m["member_name"], rows[i].MemberName)
		}
		if len(m) != len(wanted) {
			t.Errorf("row %d projected to %d keys, want %d", i, len(m), len(wanted))
		}
	}
}

func TestSummaryColumnsFor_AllTagsAreValid(t *testing.T) {
	// Every projection column must correspond to an actual JSON tag on
	// the matching struct — otherwise GORM Select will silently drop the
	// column and the renderer will see undefined.
	type structFor struct {
		tableType string
		validTags map[string]bool
	}
	mrr := jsonTagSet(models.MemberRatingResult{})
	md := jsonTagSet(models.GPricingMemberData{})
	mps := jsonTagSet(models.MemberPremiumSchedule{})

	cases := []structFor{
		{"member_rating_results", mrr},
		{"member_data", md},
		{"member_premium_schedules", mps},
	}
	for _, tc := range cases {
		for _, col := range summaryColumnsFor(tc.tableType) {
			if !tc.validTags[col] {
				t.Errorf("table %s: projection column %q has no matching json tag on the struct", tc.tableType, col)
			}
		}
	}
}

func jsonTagSet(v interface{}) map[string]bool {
	out := make(map[string]bool)
	for _, tag := range getJSONTags(v) {
		out[tag] = true
	}
	return out
}

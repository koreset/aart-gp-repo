package services

import (
	"api/models"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ensureTestDB initialises DB with just enough scaffolding for naming-strategy
// lookups when the test binary is run without a real database connection.
func ensureTestDB(t *testing.T) {
	t.Helper()
	if DB != nil && DB.Config != nil && DB.Config.NamingStrategy != nil {
		return
	}
	DB = &gorm.DB{Config: &gorm.Config{NamingStrategy: schema.NamingStrategy{}}}
}

func contains(cols []string, target string) bool {
	for _, c := range cols {
		if c == target {
			return true
		}
	}
	return false
}

func TestGetStructDBColumns_GPricingMemberDataExpandsBenefits(t *testing.T) {
	ensureTestDB(t)

	cols := getStructDBColumns(models.GPricingMemberData{})

	if contains(cols, "benefits") {
		t.Fatalf("expected embedded Benefits to be expanded, got bare column %q in %v", "benefits", cols)
	}

	wantBenefitCols := []string{
		"benefits_gla_enabled",
		"benefits_gla_multiple",
		"benefits_sgla_enabled",
		"benefits_sgla_multiple",
		"benefits_ptd_enabled",
		"benefits_ptd_multiple",
		"benefits_ci_enabled",
		"benefits_ci_multiple",
		"benefits_ttd_enabled",
		"benefits_ttd_multiple",
		"benefits_phi_enabled",
		"benefits_phi_multiple",
		"benefits_gff_enabled",
	}
	for _, c := range wantBenefitCols {
		if !contains(cols, c) {
			t.Errorf("missing expected column %q in %v", c, cols)
		}
	}

	for _, c := range []string{"id", "year", "scheme_name", "member_name", "quote_id"} {
		if !contains(cols, c) {
			t.Errorf("missing expected top-level column %q in %v", c, cols)
		}
	}
}

func TestGetStructDBColumns_GPricingMemberDataInForceSkipsGormDash(t *testing.T) {
	ensureTestDB(t)

	cols := getStructDBColumns(models.GPricingMemberDataInForce{})

	if contains(cols, "benefits") {
		t.Fatalf("expected embedded Benefits to be expanded, got bare column %q in %v", "benefits", cols)
	}
	if contains(cols, "scheme_category_details") {
		t.Errorf("gorm:\"-\" field SchemeCategoryDetails should be skipped, got %v", cols)
	}
	if !contains(cols, "benefits_gla_enabled") {
		t.Errorf("missing expected column %q in %v", "benefits_gla_enabled", cols)
	}
}

func TestGetStructDBColumns_MemberRatingResultUnaffected(t *testing.T) {
	ensureTestDB(t)

	cols := getStructDBColumns(models.MemberRatingResult{})

	// No embedded fields on this model — sanity check a couple of well-known
	// columns to confirm the helper still returns the flat list.
	for _, c := range []string{"financial_year", "category", "member_name", "member_count", "gender"} {
		if !contains(cols, c) {
			t.Errorf("missing expected column %q in %v", c, cols)
		}
	}
}

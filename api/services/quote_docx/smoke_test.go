package quote_docx

import (
	"archive/zip"
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"api/models"
)

// TestSmokeAssembly builds a small document end-to-end from synthetic data
// (skipping the DB-dependent GenerateQuoteDocx) and verifies the resulting
// .docx is a well-formed ZIP containing the expected parts.
func TestSmokeAssembly(t *testing.T) {
	quote := models.GroupPricingQuote{
		ID:               1,
		QuoteName:        "Q-SMOKE-001",
		SchemeName:       "Smoke Test Scheme",
		CommencementDate: mustParse("2026-05-01"),
		CreationDate:     mustParse("2026-04-14"),
		FreeCoverLimit:   500000,
		SchemeCategories: []models.SchemeCategory{
			{
				SchemeCategory:            "Staff",
				GlaBenefit:                true,
				GlaSalaryMultiple:         3,
				GlaTerminalIllnessBenefit: "Yes",
				GlaWaitingPeriod:          0,
				GlaEducatorBenefit:        "No",
				FamilyFuneralBenefit:      true,
				FamilyFuneralMainMemberFuneralSumAssured: 20000,
				FamilyFuneralSpouseFuneralSumAssured:     20000,
				FamilyFuneralChildrenFuneralSumAssured:   10000,
				FamilyFuneralMaxNumberChildren:           4,
			},
		},
		UseGlobalSalaryMultiple: true,
	}
	summaries := []models.MemberRatingResultSummary{
		{
			Category:                                  "Staff",
			MemberCount:                               50,
			TotalAnnualSalary:                         12_000_000,
			TotalSumAssured:                           36_000_000,
			TotalGlaCappedSumAssured:                  36_000_000,
			ExpTotalGlaAnnualOfficePremium:            120_000,
			ExpProportionGlaOfficePremiumSalary:       0.01,
			ExpTotalAnnualPremiumExclFuneral:          120_000,
			ProportionExpTotalPremiumExclFuneralSalary: 0.01,
			ExpTotalFunAnnualOfficePremium:            9_000,
			ExpTotalFunAnnualPremiumPerMember:         180,
			ExpTotalFunMonthlyPremiumPerMember:        15,
			TotalAnnualPremium:                        129_000,
		},
	}
	insurer := models.GroupPricingInsurerDetail{
		Name:                  "Example Insurer Ltd",
		AddressLine1:          "1 Example St",
		AddressLine2:          "Suite 200",
		City:                  "Johannesburg",
		Province:              "GP",
		PostCode:              "2000",
		Telephone:             "+27 11 000 0000",
		Email:                 "quotes@example.com",
		IntroductoryText:      "We're pleased to submit this quote.",
		GeneralProvisionsText: "Standard underwriting provisions apply.",
	}
	titles := BenefitTitles{
		GlaBenefitTitle:           "Group Life",
		SglaBenefitTitle:          "Spouse Life",
		PtdBenefitTitle:           "PTD",
		CiBenefitTitle:            "Critical Illness",
		PhiBenefitTitle:           "PHI",
		TtdBenefitTitle:           "TTD",
		FamilyFuneralBenefitTitle: "Family Funeral",
	}

	var body strings.Builder
	body.WriteString(BuildCoverAndSummarySection(quote, summaries, insurer))
	body.WriteString(BuildPremiumSummarySection(quote, summaries))
	body.WriteString(BuildPremiumBreakdownSection(quote, summaries, titles))
	body.WriteString(BuildBenefitsDefinitionsSection(quote, summaries, nil, titles))
	body.WriteString(BuildProvisionsSection(quote, insurer))
	body.WriteString(BuildAcceptanceFormSection(quote))

	fixed := fixupSectionBreaks(body.String())
	data, err := BuildFullPackage(fixed, quote.SchemeName, nil, "")
	if err != nil {
		t.Fatalf("BuildFullPackage failed: %v", err)
	}
	if len(data) < 2000 {
		t.Fatalf("output suspiciously small: %d bytes", len(data))
	}

	// Confirm valid ZIP and expected parts exist.
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("output is not a valid ZIP: %v", err)
	}
	want := map[string]bool{
		"[Content_Types].xml":          false,
		"_rels/.rels":                  false,
		"word/document.xml":            false,
		"word/styles.xml":              false,
		"word/_rels/document.xml.rels": false,
	}
	for _, f := range zr.File {
		if _, ok := want[f.Name]; ok {
			want[f.Name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("missing required part: %s", name)
		}
	}

	// Save for manual inspection.
	out := "/tmp/quote_docx_smoke.docx"
	if err := os.WriteFile(out, data, 0o644); err == nil {
		t.Logf("wrote %s (%d bytes)", out, len(data))
	}
}

func mustParse(s string) (t time.Time) {
	t, _ = time.Parse("2006-01-02", s)
	return
}

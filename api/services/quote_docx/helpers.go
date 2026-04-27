package quote_docx

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"api/models"
)

// RoundUpToTwoDecimalsAccounting rounds up to 2 decimal places and formats with space-separated thousands
func RoundUpToTwoDecimalsAccounting(num float64) string {
	// Ceiling to 2 decimals
	rounded := math.Ceil(num*100) / 100

	// Format with 2 decimal places and space-separated thousands
	// Use manual formatting (no locale) to avoid dependencies
	intPart := int64(rounded)
	fracPart := int64(math.Round((rounded - float64(intPart)) * 100))

	// Format integer part with space-separated thousands
	intStr := strconv.FormatInt(intPart, 10)
	var result strings.Builder

	// Add space separators for thousands
	for i, ch := range intStr {
		if i > 0 && (len(intStr)-i)%3 == 0 {
			result.WriteRune(' ')
		}
		result.WriteRune(ch)
	}

	// Add fractional part
	fracStr := fmt.Sprintf("%02d", fracPart)
	result.WriteString(".")
	result.WriteString(fracStr)

	return result.String()
}

// officePercentSalary derives a benefit's office premium proportion of
// salary as a printable percent string, from the persisted
// risk-proportion-of-salary field on the summary. Office =
// risk / (1 - SchemeTotalLoading()), so the same factor scales the
// proportion. Returns "0%" when the scheme loading saturates the
// denominator.
func officePercentSalary(item models.MemberRatingResultSummary, riskProportion float64) string {
	denom := 1.0 - item.SchemeTotalLoading()
	if denom <= 0 {
		return "0%"
	}
	return fmt.Sprintf("%s%%", RoundUpToTwoDecimalsAccounting(riskProportion*100/denom))
}

// FormatQuoteDate formats time.Time to "02 Jan 2006" format
func FormatQuoteDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006")
}

// CalculateQuoteTotals calculates aggregated totals across all result summaries
func CalculateQuoteTotals(summaries []models.MemberRatingResultSummary) QuoteTotals {
	totals := QuoteTotals{}
	for _, item := range summaries {
		totals.TotalLives += int(item.MemberCount)
		totals.TotalSumAssured += item.TotalSumAssured
		totals.TotalAnnualSalary += item.TotalAnnualSalary
		totals.TotalAnnualPremium += item.ExpTotalAnnualPremiumExclFuneral
	}
	return totals
}

// HasAnyNonFuneralBenefits checks if any category has non-funeral benefits
func HasAnyNonFuneralBenefits(summaries []models.MemberRatingResultSummary) bool {
	for _, item := range summaries {
		if CategoryHasNonFuneralBenefits(item) {
			return true
		}
	}
	return false
}

// CategoryHasNonFuneralBenefits checks if a single category has non-funeral benefits
func CategoryHasNonFuneralBenefits(item models.MemberRatingResultSummary) bool {
	return item.TotalGlaCappedSumAssured > 0 ||
		item.TotalPtdCappedSumAssured > 0 ||
		item.TotalCiCappedSumAssured > 0 ||
		item.TotalSglaCappedSumAssured > 0 ||
		item.TotalPhiCappedIncome > 0 ||
		item.TotalTtdCappedIncome > 0
}

// BuildInitialInfoRows builds the initial info / quote summary key-value rows
func BuildInitialInfoRows(quote models.GroupPricingQuote, totals QuoteTotals) []LabelValueRow {
	return []LabelValueRow{
		{Label: "Type of Policy:", Value: "Group Risk Assurance"},
		{Label: "Quote Number:", Value: quote.QuoteName},
		{Label: "Quote Date:", Value: FormatQuoteDate(quote.CreationDate)},
		{Label: "Scheme Name:", Value: quote.SchemeName},
		{Label: "Inception Date:", Value: FormatQuoteDate(quote.CommencementDate)},
		{Label: "Number of Lives Covered:", Value: fmt.Sprintf("%d", totals.TotalLives)},
		{Label: "Total Sum Assured:", Value: RoundUpToTwoDecimalsAccounting(totals.TotalSumAssured)},
		{Label: "Total Annual Salary:", Value: RoundUpToTwoDecimalsAccounting(totals.TotalAnnualSalary)},
		{Label: "Total Annual Premium:", Value: RoundUpToTwoDecimalsAccounting(totals.TotalAnnualPremium)},
	}
}

// BuildPremiumSummaryRows builds rows for the Premium Summary table
// Returns data rows + a totals row
func BuildPremiumSummaryRows(summaries []models.MemberRatingResultSummary) []PremiumSummaryRow {
	rows := []PremiumSummaryRow{}

	for _, item := range summaries {
		if !CategoryHasNonFuneralBenefits(item) {
			continue
		}

		percentSalary := 0.0
		if item.TotalAnnualSalary > 0 {
			percentSalary = item.ExpTotalAnnualPremiumExclFuneral / item.TotalAnnualSalary
		}
		rows = append(rows, PremiumSummaryRow{
			Category:        item.Category,
			MemberCount:     fmt.Sprintf("%.0f", item.MemberCount),
			TotalSalary:     RoundUpToTwoDecimalsAccounting(item.TotalAnnualSalary),
			TotalSumAssured: RoundUpToTwoDecimalsAccounting(item.TotalSumAssured),
			AnnualPremium:   RoundUpToTwoDecimalsAccounting(item.ExpTotalAnnualPremiumExclFuneral),
			PercentSalary:   fmt.Sprintf("%s%%", RoundUpToTwoDecimalsAccounting(percentSalary*100)),
		})
	}

	// Calculate totals row
	totalLives := 0.0
	totalSalary := 0.0
	totalSumAssured := 0.0
	totalPremium := 0.0

	for _, item := range summaries {
		totalLives += item.MemberCount
		totalSalary += item.TotalAnnualSalary
		totalSumAssured += item.TotalSumAssured
		totalPremium += item.ExpTotalAnnualPremiumExclFuneral
	}

	percentSalaryStr := "0.00%"
	if totalSalary > 0 {
		percentSalaryStr = fmt.Sprintf("%s%%", RoundUpToTwoDecimalsAccounting((totalPremium/totalSalary)*100))
	}

	rows = append(rows, PremiumSummaryRow{
		Category:        "Total",
		MemberCount:     fmt.Sprintf("%.0f", totalLives),
		TotalSalary:     RoundUpToTwoDecimalsAccounting(totalSalary),
		TotalSumAssured: RoundUpToTwoDecimalsAccounting(totalSumAssured),
		AnnualPremium:   RoundUpToTwoDecimalsAccounting(totalPremium),
		PercentSalary:   percentSalaryStr,
	})

	return rows
}

// BuildGroupFuneralRows builds rows for the Group Funeral summary table
func BuildGroupFuneralRows(summaries []models.MemberRatingResultSummary) []GroupFuneralRow {
	rows := []GroupFuneralRow{}

	for _, item := range summaries {
		rows = append(rows, GroupFuneralRow{
			Category:           item.Category,
			MemberCount:        fmt.Sprintf("%.0f", item.MemberCount),
			MonthlyPremium:     RoundUpToTwoDecimalsAccounting(item.ExpTotalFunMonthlyPremiumPerMember),
			AnnualPremium:      RoundUpToTwoDecimalsAccounting(item.ExpTotalFunAnnualPremiumPerMember),
			TotalAnnualPremium: RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalFunAnnualRiskPremium, &item)),
		})
	}

	// Totals row
	totalLives := 0.0
	totalMonthly := 0.0
	totalAnnual := 0.0
	totalAnnualOffice := 0.0

	for _, item := range summaries {
		totalLives += item.MemberCount
		totalMonthly += item.ExpTotalFunMonthlyPremiumPerMember
		totalAnnual += item.ExpTotalFunAnnualPremiumPerMember
		totalAnnualOffice += models.ComputeOfficePremium(item.ExpTotalFunAnnualRiskPremium, &item)
	}

	rows = append(rows, GroupFuneralRow{
		Category:           "Total",
		MemberCount:        fmt.Sprintf("%.0f", totalLives),
		MonthlyPremium:     RoundUpToTwoDecimalsAccounting(totalMonthly),
		AnnualPremium:      RoundUpToTwoDecimalsAccounting(totalAnnual),
		TotalAnnualPremium: RoundUpToTwoDecimalsAccounting(totalAnnualOffice),
	})

	return rows
}

// BuildPremiumBreakdownRows builds per-category benefit premium breakdown rows
func BuildPremiumBreakdownRows(item models.MemberRatingResultSummary, titles BenefitTitles) []PremiumBreakdownRow {
	return []PremiumBreakdownRow{
		{
			Benefit:         titles.GlaBenefitTitle,
			TotalSumAssured: RoundUpToTwoDecimalsAccounting(item.TotalGlaCappedSumAssured),
			AnnualPremium:   RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalGlaAnnualRiskPremium, &item)),
			PercentSalary:   officePercentSalary(item, item.ExpProportionGlaAnnualRiskPremiumSalary),
		},
		{
			Benefit:         titles.SglaBenefitTitle,
			TotalSumAssured: RoundUpToTwoDecimalsAccounting(item.TotalSglaCappedSumAssured),
			AnnualPremium:   RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalSglaAnnualRiskPremium, &item)),
			PercentSalary:   officePercentSalary(item, item.ExpProportionSglaAnnualRiskPremiumSalary),
		},
		{
			Benefit:         titles.PtdBenefitTitle,
			TotalSumAssured: RoundUpToTwoDecimalsAccounting(item.TotalPtdCappedSumAssured),
			AnnualPremium:   RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalPtdAnnualRiskPremium, &item)),
			PercentSalary:   officePercentSalary(item, item.ExpProportionPtdAnnualRiskPremiumSalary),
		},
		{
			Benefit:         titles.CiBenefitTitle,
			TotalSumAssured: RoundUpToTwoDecimalsAccounting(item.TotalCiCappedSumAssured),
			AnnualPremium:   RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalCiAnnualRiskPremium, &item)),
			PercentSalary:   officePercentSalary(item, item.ExpProportionCiAnnualRiskPremiumSalary),
		},
		{
			Benefit:         titles.PhiBenefitTitle,
			TotalSumAssured: RoundUpToTwoDecimalsAccounting(item.TotalPhiCappedIncome),
			AnnualPremium:   RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalPhiAnnualRiskPremium, &item)),
			PercentSalary:   officePercentSalary(item, item.ExpProportionPhiAnnualRiskPremiumSalary),
		},
		{
			Benefit:         titles.TtdBenefitTitle,
			TotalSumAssured: RoundUpToTwoDecimalsAccounting(item.TotalTtdCappedIncome),
			AnnualPremium:   RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalTtdAnnualRiskPremium, &item)),
			PercentSalary:   officePercentSalary(item, item.ExpProportionTtdAnnualRiskPremiumSalary),
		},
	}
}

// BuildGroupFuneralBreakdownRows builds per-category group funeral breakdown rows (key-value pairs)
func BuildGroupFuneralBreakdownRows(item models.MemberRatingResultSummary) []LabelValueRow {
	return []LabelValueRow{
		{
			Label: "Monthly Premium per Member",
			Value: RoundUpToTwoDecimalsAccounting(item.ExpTotalFunMonthlyPremiumPerMember),
		},
		{
			Label: "Annual Premium per Member",
			Value: RoundUpToTwoDecimalsAccounting(item.ExpTotalFunAnnualPremiumPerMember),
		},
		{
			Label: "Total Annual Premium",
			Value: RoundUpToTwoDecimalsAccounting(models.ComputeOfficePremium(item.ExpTotalFunAnnualRiskPremium, &item)),
		},
	}
}

// BuildCategoryCommonBenefitRows builds common category benefits key-value rows
func BuildCategoryCommonBenefitRows(cat models.SchemeCategory, quote models.GroupPricingQuote) []LabelValueRow {
	return []LabelValueRow{
		{Label: "Terminal Illness", Value: cat.GlaTerminalIllnessBenefit},
		{Label: "Free Cover Limit", Value: fmt.Sprintf("%.0f", quote.FreeCoverLimit)},
		{Label: "Gla Educator", Value: cat.GlaEducatorBenefit},
		{Label: "Ptd Educator", Value: cat.PtdEducatorBenefit},
		{Label: "Retirement Premium Waiver", Value: cat.PhiPremiumWaiver},
		{Label: "Medical Aid Premium Waiver", Value: cat.PhiMedicalAidPremiumWaiver},
	}
}

// BuildBenefitDefinitionRows builds the 7-column benefit definition rows
func BuildBenefitDefinitionRows(cat models.SchemeCategory, quote models.GroupPricingQuote, titles BenefitTitles) []BenefitDefinitionRow {
	salaryMultipleFmt := func(useGlobal bool, val float64) string {
		if useGlobal {
			return fmt.Sprintf("%.0f", val)
		}
		return "varies"
	}

	return []BenefitDefinitionRow{
		{
			Benefit:          titles.GlaBenefitTitle,
			SalaryMultiple:   salaryMultipleFmt(quote.UseGlobalSalaryMultiple, cat.GlaSalaryMultiple),
			BenefitStructure: "standalone",
			WaitingPeriod:    fmt.Sprintf("%d", cat.GlaWaitingPeriod),
			DeferredPeriod:   "n.a",
			CoverDefinition:  "n.a",
			RiskType:         "n.a",
		},
		{
			Benefit:          titles.SglaBenefitTitle,
			SalaryMultiple:   salaryMultipleFmt(quote.UseGlobalSalaryMultiple, cat.SglaSalaryMultiple),
			BenefitStructure: "rider",
			WaitingPeriod:    "0",
			DeferredPeriod:   "n.a",
			CoverDefinition:  "n.a",
			RiskType:         "n.a",
		},
		{
			Benefit:          titles.PtdBenefitTitle,
			SalaryMultiple:   salaryMultipleFmt(quote.UseGlobalSalaryMultiple, cat.PtdSalaryMultiple),
			BenefitStructure: cat.PtdBenefitType,
			WaitingPeriod:    "0",
			DeferredPeriod:   fmt.Sprintf("%d", cat.PtdDeferredPeriod),
			CoverDefinition:  cat.PtdDisabilityDefinition,
			RiskType:         cat.PtdRiskType,
		},
		{
			Benefit:          titles.CiBenefitTitle,
			SalaryMultiple:   salaryMultipleFmt(quote.UseGlobalSalaryMultiple, cat.CiCriticalIllnessSalaryMultiple),
			BenefitStructure: cat.CiBenefitStructure,
			WaitingPeriod:    "0",
			DeferredPeriod:   "0",
			CoverDefinition:  cat.CiBenefitDefinition,
			RiskType:         "n.a",
		},
		{
			Benefit:          titles.PhiBenefitTitle,
			SalaryMultiple:   fmt.Sprintf("%.2f", cat.PhiIncomeReplacementPercentage/100),
			BenefitStructure: "n.a",
			WaitingPeriod:    fmt.Sprintf("%d", cat.PhiWaitingPeriod),
			DeferredPeriod:   fmt.Sprintf("%d", cat.PhiDeferredPeriod),
			CoverDefinition:  cat.PhiDisabilityDefinition,
			RiskType:         cat.PhiRiskType,
		},
		{
			Benefit:          titles.TtdBenefitTitle,
			SalaryMultiple:   fmt.Sprintf("%.2f", cat.TtdIncomeReplacementPercentage/100),
			BenefitStructure: "n.a",
			WaitingPeriod:    fmt.Sprintf("%d", cat.TtdWaitingPeriod),
			DeferredPeriod:   fmt.Sprintf("%d", cat.TtdDeferredPeriod),
			CoverDefinition:  cat.TtdDisabilityDefinition,
			RiskType:         cat.TtdRiskType,
		},
	}
}

// BuildFuneralCoverageRows builds the group funeral coverage rows for a given scheme category
func BuildFuneralCoverageRows(cat models.SchemeCategory) []FuneralCoverageRow {
	return []FuneralCoverageRow{
		{
			Member:     "Main Member",
			SumAssured: cat.FamilyFuneralMainMemberFuneralSumAssured,
			MaxCovered: 1,
		},
		{
			Member:     "Spouse",
			SumAssured: cat.FamilyFuneralSpouseFuneralSumAssured,
			MaxCovered: 1,
		},
		{
			Member:     "Child",
			SumAssured: cat.FamilyFuneralChildrenFuneralSumAssured,
			MaxCovered: cat.FamilyFuneralMaxNumberChildren,
		},
		{
			Member:     "Parent",
			SumAssured: cat.FamilyFuneralParentFuneralSumAssured,
			MaxCovered: 1, // Safe default per requirements
		},
		{
			Member:     "Dependant",
			SumAssured: cat.FamilyFuneralAdultDependantSumAssured,
			MaxCovered: cat.FamilyFuneralMaxNumberAdultDependants,
		},
	}
}

// BuildEducatorBenefitRows builds educator benefit rows for a given scheme category
func BuildEducatorBenefitRows(benefits []interface{}, schemeCategory string) []EducatorBenefitRow {
	// Benefits interface{} represents CategoryEducatorBenefit from services
	// For now, return empty slice as educator benefits are complex to map
	return []EducatorBenefitRow{}
}

// ResolveBenefitTitles looks up each benefit_code from maps and uses alias if non-empty else name
func ResolveBenefitTitles(maps []models.GroupBenefitMapper) BenefitTitles {
	titles := BenefitTitles{
		GlaBenefitTitle:           "GLA",
		SglaBenefitTitle:          "SGLA",
		PtdBenefitTitle:           "PTD",
		CiBenefitTitle:            "CI",
		PhiBenefitTitle:           "PHI",
		TtdBenefitTitle:           "TTD",
		FamilyFuneralBenefitTitle: "Family Funeral",
	}

	// Resolve from maps (if available)
	for _, m := range maps {
		if m.BenefitAlias != "" {
			label := m.BenefitAlias
			switch m.BenefitCode {
			case "GLA":
				titles.GlaBenefitTitle = label
			case "SGLA":
				titles.SglaBenefitTitle = label
			case "PTD":
				titles.PtdBenefitTitle = label
			case "CI":
				titles.CiBenefitTitle = label
			case "PHI":
				titles.PhiBenefitTitle = label
			case "TTD":
				titles.TtdBenefitTitle = label
			case "Family Funeral":
				titles.FamilyFuneralBenefitTitle = label
			}
		} else if m.BenefitName != "" {
			label := m.BenefitName
			switch m.BenefitCode {
			case "GLA":
				titles.GlaBenefitTitle = label
			case "SGLA":
				titles.SglaBenefitTitle = label
			case "PTD":
				titles.PtdBenefitTitle = label
			case "CI":
				titles.CiBenefitTitle = label
			case "PHI":
				titles.PhiBenefitTitle = label
			case "TTD":
				titles.TtdBenefitTitle = label
			case "Family Funeral":
				titles.FamilyFuneralBenefitTitle = label
			}
		}
	}

	return titles
}

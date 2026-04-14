package quote_template

import (
	"fmt"
	"strconv"

	"api/models"
	"api/services"
	"api/services/quote_docx"
)

// Context is a map of template variables
type Context map[string]interface{}

// BuildContext constructs the complete context for template substitution
func BuildContext(quoteID string) (Context, error) {
	// Fetch quote
	quote, err := services.GetGroupPricingQuote(quoteID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	// Fetch result summaries
	summaries, err := services.GetGroupPricingQuoteResultSummary(quote.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch result summary: %w", err)
	}

	// Fetch insurer details
	insurer, err := services.GetInsurerDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch insurer details: %w", err)
	}

	// Fetch benefit maps
	benefitMaps, err := services.GetBenefitMaps()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch benefit maps: %w", err)
	}

	// Build totals
	totals := quote_docx.CalculateQuoteTotals(summaries)
	hasNonFuneral := quote_docx.HasAnyNonFuneralBenefits(summaries)

	// Build insurer context
	insurerCtx := map[string]interface{}{
		"name":                    insurer.Name,
		"address_line_1":          insurer.AddressLine1,
		"address_line_2":          insurer.AddressLine2,
		"address_line_3":          insurer.AddressLine3,
		"city":                    insurer.City,
		"province":                insurer.Province,
		"post_code":               insurer.PostCode,
		"telephone":               insurer.Telephone,
		"email":                   insurer.Email,
		"introductory_text":       insurer.IntroductoryText,
		"general_provisions_text": insurer.GeneralProvisionsText,
	}

	// Build categories
	categories := buildCategoriesContext(summaries, benefitMaps)

	// Build main context
	ctx := Context{
		"quote_name":                  quote.QuoteName,
		"scheme_name":                 quote.SchemeName,
		"quote_number":                quote.QuoteName,
		"creation_date":               quote_docx.FormatQuoteDate(quote.CreationDate),
		"commencement_date":           quote_docx.FormatQuoteDate(quote.CommencementDate),
		"industry":                    quote.Industry,
		"currency":                    quote.Currency,
		"free_cover_limit":            quote_docx.RoundUpToTwoDecimalsAccounting(quote.FreeCoverLimit),
		"total_lives":                 strconv.Itoa(totals.TotalLives),
		"total_sum_assured":           quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalSumAssured),
		"total_annual_salary":         quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualSalary),
		"total_annual_premium":        quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualPremium),
		"has_non_funeral_benefits":    hasNonFuneral,
		"insurer":                     insurerCtx,
		"categories":                  categories,
	}

	return ctx, nil
}

// buildCategoriesContext builds the categories array for the template
func buildCategoriesContext(summaries []models.MemberRatingResultSummary, benefitMaps []models.GroupBenefitMapper) []map[string]interface{} {
	var categories []map[string]interface{}

	for _, summary := range summaries {
		hasNonFuneral := quote_docx.CategoryHasNonFuneralBenefits(summary)

		// Determine which benefits are present
		hasGLA := summary.TotalGlaCappedSumAssured > 0
		hasSGLA := summary.TotalSglaCappedSumAssured > 0
		hasPTD := summary.TotalPtdCappedSumAssured > 0
		hasCI := summary.TotalCiCappedSumAssured > 0
		hasPHI := summary.TotalPhiCappedIncome > 0
		hasTTD := summary.TotalTtdCappedIncome > 0
		hasFuneral := summary.TotalFunAnnualOfficePremium > 0

		// Build benefit contexts
		glaCtx := map[string]interface{}{}
		if hasGLA {
			glaCtx = map[string]interface{}{
				"title":                     "Group Life",
				"salary_multiple":           "3", // placeholder; can be enhanced
				"waiting_period":            "0",
				"terminal_illness_benefit":  "Yes",
				"educator_benefit":          "No",
				"total_sum_assured":         quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalGlaCappedSumAssured),
				"annual_premium":            quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalGlaAnnualOfficePremium),
				"percent_salary":            formatPercent(summary.ExpProportionGlaOfficePremiumSalary),
			}
		}

		sglaCtx := map[string]interface{}{}
		if hasSGLA {
			sglaCtx = map[string]interface{}{
				"title":                   "Spouse Group Life",
				"salary_multiple":         "0",
				"waiting_period":          "0",
				"total_sum_assured":       quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalSglaCappedSumAssured),
				"annual_premium":          quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalSglaAnnualOfficePremium),
				"percent_salary":          formatPercent(summary.ExpProportionSglaOfficePremiumSalary),
			}
		}

		ptdCtx := map[string]interface{}{}
		if hasPTD {
			ptdCtx = map[string]interface{}{
				"title":                   "Permanent Total Disability",
				"waiting_period":          "90",
				"total_sum_assured":       quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalPtdCappedSumAssured),
				"annual_premium":          quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalPtdAnnualOfficePremium),
				"percent_salary":          formatPercent(summary.ExpProportionPtdOfficePremiumSalary),
			}
		}

		ciCtx := map[string]interface{}{}
		if hasCI {
			ciCtx = map[string]interface{}{
				"title":                   "Critical Illness",
				"waiting_period":          "0",
				"total_sum_assured":       quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalCiCappedSumAssured),
				"annual_premium":          quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalCiAnnualOfficePremium),
				"percent_salary":          formatPercent(summary.ExpProportionCiOfficePremiumSalary),
			}
		}

		phiCtx := map[string]interface{}{}
		if hasPHI {
			phiCtx = map[string]interface{}{
				"title":                   "Personal Accident / Income Protection",
				"waiting_period":          "30",
				"total_sum_assured":       quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalPhiCappedIncome),
				"annual_premium":          quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalPhiAnnualOfficePremium),
				"percent_salary":          formatPercent(summary.ExpProportionPhiOfficePremiumSalary),
			}
		}

		ttdCtx := map[string]interface{}{}
		if hasTTD {
			ttdCtx = map[string]interface{}{
				"title":                   "Temporary Total Disability",
				"waiting_period":          "14",
				"total_sum_assured":       quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalTtdCappedIncome),
				"annual_premium":          quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalTtdAnnualOfficePremium),
				"percent_salary":          formatPercent(summary.ExpProportionTtdOfficePremiumSalary),
			}
		}

		funeralCtx := map[string]interface{}{}
		if hasFuneral {
			funeralCtx = map[string]interface{}{
				"monthly_premium_per_member":  quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalFunMonthlyPremiumPerMember),
				"annual_premium_per_member":   quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalFunAnnualPremiumPerMember),
				"total_annual_premium":        quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalFunAnnualOfficePremium),
				"main_member_sum_assured":     "20000",  // placeholder
				"spouse_sum_assured":          "20000",  // placeholder
				"child_sum_assured":           "10000",  // placeholder
				"max_children":                "4",      // placeholder
				"parent_sum_assured":          "0",      // placeholder
				"max_dependants":              "0",      // placeholder
			}
		}

		categoryCtx := map[string]interface{}{
			"name":                        summary.Category,
			"member_count":                strconv.Itoa(int(summary.MemberCount)),
			"total_salary":                quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalAnnualSalary),
			"total_sum_assured":           quote_docx.RoundUpToTwoDecimalsAccounting(summary.TotalSumAssured),
			"annual_premium":              quote_docx.RoundUpToTwoDecimalsAccounting(summary.ExpTotalAnnualPremiumExclFuneral),
			"percent_salary":              formatPercent(summary.ProportionExpTotalPremiumExclFuneralSalary),
			"has_non_funeral_benefits":    hasNonFuneral,
			"has_gla":                     hasGLA,
			"has_sgla":                    hasSGLA,
			"has_ptd":                     hasPTD,
			"has_ci":                      hasCI,
			"has_phi":                     hasPHI,
			"has_ttd":                     hasTTD,
			"has_funeral":                 hasFuneral,
			"gla":                         glaCtx,
			"sgla":                        sglaCtx,
			"ptd":                         ptdCtx,
			"ci":                          ciCtx,
			"phi":                         phiCtx,
			"ttd":                         ttdCtx,
			"funeral":                     funeralCtx,
		}

		categories = append(categories, categoryCtx)
	}

	return categories
}

// formatPercent formats a decimal percentage (e.g. 0.01 -> "1.00%")
func formatPercent(decimal float64) string {
	return fmt.Sprintf("%.2f%%", decimal*100)
}

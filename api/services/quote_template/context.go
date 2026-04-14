package quote_template

import (
	"fmt"
	"strconv"

	"api/models"
	"api/services"
	"api/services/quote_docx"
)

// Context is the map of template variables handed to the render engine.
// Keys nested under "insurer" and each entry of "categories" are accessed
// via dot syntax in templates (e.g. {{insurer.name}}, {{gla.annual_premium}}
// inside a {{#categories}} block).
type Context map[string]interface{}

// BuildContext assembles the complete context for a given quote. It joins
// the result summaries with the quote's SchemeCategory config so per-benefit
// detail fields (salary multiple, waiting period, benefit type, etc.) come
// from the authored quote rather than hard-coded placeholders.
func BuildContext(quoteID string) (Context, error) {
	quote, err := services.GetGroupPricingQuote(quoteID)
	if err != nil {
		return nil, fmt.Errorf("fetch quote: %w", err)
	}
	summaries, err := services.GetGroupPricingQuoteResultSummary(quote.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch result summary: %w", err)
	}
	insurer, err := services.GetInsurerDetails()
	if err != nil {
		return nil, fmt.Errorf("fetch insurer: %w", err)
	}
	benefitMaps, err := services.GetBenefitMaps()
	if err != nil {
		return nil, fmt.Errorf("fetch benefit maps: %w", err)
	}

	titles := quote_docx.ResolveBenefitTitles(benefitMaps)
	totals := quote_docx.CalculateQuoteTotals(summaries)
	hasNonFuneral := quote_docx.HasAnyNonFuneralBenefits(summaries)

	// Index scheme categories by name so we can join them to summaries.
	catByName := make(map[string]models.SchemeCategory, len(quote.SchemeCategories))
	for _, c := range quote.SchemeCategories {
		catByName[c.SchemeCategory] = c
	}

	categories := make([]map[string]interface{}, 0, len(summaries))
	for _, s := range summaries {
		cat, ok := catByName[s.Category]
		if !ok {
			// Summary with no matching scheme category — populate what we can.
			cat = models.SchemeCategory{SchemeCategory: s.Category}
		}
		categories = append(categories, buildCategoryContext(s, cat, quote, titles))
	}

	return Context{
		// Quote-level
		"quote_name":                  quote.QuoteName,
		"quote_number":                quote.QuoteName,
		"scheme_name":                 quote.SchemeName,
		"creation_date":               quote_docx.FormatQuoteDate(quote.CreationDate),
		"commencement_date":           quote_docx.FormatQuoteDate(quote.CommencementDate),
		"industry":                    quote.Industry,
		"currency":                    quote.Currency,
		"free_cover_limit":            quote_docx.RoundUpToTwoDecimalsAccounting(quote.FreeCoverLimit),
		"normal_retirement_age":       strconv.Itoa(quote.NormalRetirementAge),
		"obligation_type":             quote.ObligationType,
		"quote_type":                  quote.QuoteType,
		"use_global_salary_multiple":  quote.UseGlobalSalaryMultiple,
		"total_lives":                 strconv.Itoa(totals.TotalLives),
		"total_sum_assured":           quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalSumAssured),
		"total_annual_salary":         quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualSalary),
		"total_annual_premium":        quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualPremium),
		"has_non_funeral_benefits":    hasNonFuneral,
		"insurer": map[string]interface{}{
			"name":                    insurer.Name,
			"contact_person":          insurer.ContactPerson,
			"address_line_1":          insurer.AddressLine1,
			"address_line_2":          insurer.AddressLine2,
			"address_line_3":          insurer.AddressLine3,
			"city":                    insurer.City,
			"province":                insurer.Province,
			"post_code":               insurer.PostCode,
			"country":                 insurer.Country,
			"telephone":               insurer.Telephone,
			"email":                   insurer.Email,
			"introductory_text":       insurer.IntroductoryText,
			"general_provisions_text": insurer.GeneralProvisionsText,
		},
		"categories": categories,
	}, nil
}

// buildCategoryContext returns the per-category map used inside a
// {{#categories}} block.
func buildCategoryContext(
	s models.MemberRatingResultSummary,
	cat models.SchemeCategory,
	quote models.GroupPricingQuote,
	titles quote_docx.BenefitTitles,
) map[string]interface{} {
	hasGLA := s.TotalGlaCappedSumAssured > 0
	hasSGLA := s.TotalSglaCappedSumAssured > 0
	hasPTD := s.TotalPtdCappedSumAssured > 0
	hasCI := s.TotalCiCappedSumAssured > 0
	hasPHI := s.TotalPhiCappedIncome > 0
	hasTTD := s.TotalTtdCappedIncome > 0
	hasFuneral := s.TotalFunAnnualOfficePremium > 0 || cat.FamilyFuneralBenefit

	ctx := map[string]interface{}{
		// Category summary
		"name":                     s.Category,
		"region":                   cat.Region,
		"member_count":             strconv.Itoa(int(s.MemberCount)),
		"total_salary":             quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalAnnualSalary),
		"total_sum_assured":        quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalSumAssured),
		"annual_premium":           quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalAnnualPremiumExclFuneral),
		"percent_salary":           formatPercent(s.ProportionExpTotalPremiumExclFuneralSalary),
		"free_cover_limit":         quote_docx.RoundUpToTwoDecimalsAccounting(cat.FreeCoverLimit),
		"has_non_funeral_benefits": quote_docx.CategoryHasNonFuneralBenefits(s),
		"has_gla":                  hasGLA,
		"has_sgla":                 hasSGLA,
		"has_ptd":                  hasPTD,
		"has_ci":                   hasCI,
		"has_phi":                  hasPHI,
		"has_ttd":                  hasTTD,
		"has_funeral":              hasFuneral,
		// Common / shared flags
		"retirement_premium_waiver":   orDash(cat.PhiPremiumWaiver),
		"medical_aid_premium_waiver":  orDash(cat.PhiMedicalAidPremiumWaiver),
		"gla_terminal_illness_benefit": orDash(cat.GlaTerminalIllnessBenefit),
		"gla_educator_benefit":         orDash(cat.GlaEducatorBenefit),
		"ptd_educator_benefit":         orDash(cat.PtdEducatorBenefit),
	}

	// Benefit-specific objects — populated only when the category has the
	// benefit. Missing maps resolve to empty strings in templates so users
	// can safely reference them.
	ctx["gla"] = glaContext(s, cat, quote, titles, hasGLA)
	ctx["sgla"] = sglaContext(s, cat, quote, titles, hasSGLA)
	ctx["ptd"] = ptdContext(s, cat, quote, titles, hasPTD)
	ctx["ci"] = ciContext(s, cat, quote, titles, hasCI)
	ctx["phi"] = phiContext(s, cat, titles, hasPHI)
	ctx["ttd"] = ttdContext(s, cat, titles, hasTTD)
	ctx["funeral"] = funeralContext(s, cat, titles, hasFuneral)

	return ctx
}

func glaContext(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles, has bool) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"title":                    t.GlaBenefitTitle,
		"salary_multiple":          salaryMultiple(q.UseGlobalSalaryMultiple, cat.GlaSalaryMultiple),
		"waiting_period":           strconv.Itoa(cat.GlaWaitingPeriod),
		"benefit_structure":        "standalone",
		"benefit_type":             orDash(cat.GlaBenefitType),
		"terminal_illness_benefit": orDash(cat.GlaTerminalIllnessBenefit),
		"educator_benefit":         orDash(cat.GlaEducatorBenefit),
		"total_sum_assured":        quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalGlaCappedSumAssured),
		"annual_premium":           quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalGlaAnnualOfficePremium),
		"percent_salary":           formatPercent(s.ExpProportionGlaOfficePremiumSalary),
	}
}

func sglaContext(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles, has bool) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"title":             t.SglaBenefitTitle,
		"salary_multiple":   salaryMultiple(q.UseGlobalSalaryMultiple, cat.SglaSalaryMultiple),
		"waiting_period":    "0",
		"benefit_structure": "rider",
		"max_benefit":       quote_docx.RoundUpToTwoDecimalsAccounting(cat.SglaMaxBenefit),
		"total_sum_assured": quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalSglaCappedSumAssured),
		"annual_premium":    quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalSglaAnnualOfficePremium),
		"percent_salary":    formatPercent(s.ExpProportionSglaOfficePremiumSalary),
	}
}

func ptdContext(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles, has bool) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"title":                 t.PtdBenefitTitle,
		"salary_multiple":       salaryMultiple(q.UseGlobalSalaryMultiple, cat.PtdSalaryMultiple),
		"waiting_period":        "0",
		"deferred_period":       strconv.Itoa(cat.PtdDeferredPeriod),
		"benefit_type":          orDash(cat.PtdBenefitType),
		"disability_definition": orDash(cat.PtdDisabilityDefinition),
		"risk_type":             orDash(cat.PtdRiskType),
		"educator_benefit":      orDash(cat.PtdEducatorBenefit),
		"total_sum_assured":     quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalPtdCappedSumAssured),
		"annual_premium":        quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalPtdAnnualOfficePremium),
		"percent_salary":        formatPercent(s.ExpProportionPtdOfficePremiumSalary),
	}
}

func ciContext(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles, has bool) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"title":              t.CiBenefitTitle,
		"salary_multiple":    salaryMultiple(q.UseGlobalSalaryMultiple, cat.CiCriticalIllnessSalaryMultiple),
		"waiting_period":     "0",
		"deferred_period":    "0",
		"benefit_structure":  orDash(cat.CiBenefitStructure),
		"benefit_definition": orDash(cat.CiBenefitDefinition),
		"max_benefit":        quote_docx.RoundUpToTwoDecimalsAccounting(cat.CiMaxBenefit),
		"total_sum_assured":  quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalCiCappedSumAssured),
		"annual_premium":     quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalCiAnnualOfficePremium),
		"percent_salary":     formatPercent(s.ExpProportionCiOfficePremiumSalary),
	}
}

func phiContext(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles, has bool) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"title":                         t.PhiBenefitTitle,
		"income_replacement_percentage": fmt.Sprintf("%.2f%%", cat.PhiIncomeReplacementPercentage),
		"waiting_period":                strconv.Itoa(cat.PhiWaitingPeriod),
		"deferred_period":               strconv.Itoa(cat.PhiDeferredPeriod),
		"disability_definition":         orDash(cat.PhiDisabilityDefinition),
		"risk_type":                     orDash(cat.PhiRiskType),
		"premium_waiver":                orDash(cat.PhiPremiumWaiver),
		"medical_aid_premium_waiver":    orDash(cat.PhiMedicalAidPremiumWaiver),
		"benefit_escalation":            orDash(cat.PhiBenefitEscalation),
		"total_sum_assured":             quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalPhiCappedIncome),
		"annual_premium":                quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalPhiAnnualOfficePremium),
		"percent_salary":                formatPercent(s.ExpProportionPhiOfficePremiumSalary),
	}
}

func ttdContext(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles, has bool) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"title":                         t.TtdBenefitTitle,
		"income_replacement_percentage": fmt.Sprintf("%.2f%%", cat.TtdIncomeReplacementPercentage),
		"waiting_period":                strconv.Itoa(cat.TtdWaitingPeriod),
		"deferred_period":               strconv.Itoa(cat.TtdDeferredPeriod),
		"disability_definition":         orDash(cat.TtdDisabilityDefinition),
		"risk_type":                     orDash(cat.TtdRiskType),
		"total_sum_assured":             quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalTtdCappedIncome),
		"annual_premium":                quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalTtdAnnualOfficePremium),
		"percent_salary":                formatPercent(s.ExpProportionTtdOfficePremiumSalary),
	}
}

func funeralContext(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles, has bool) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"title":                      t.FamilyFuneralBenefitTitle,
		"monthly_premium_per_member": quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunMonthlyPremiumPerMember),
		"annual_premium_per_member":  quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunAnnualPremiumPerMember),
		"total_annual_premium":       quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalFunAnnualOfficePremium),
		"main_member_sum_assured":    quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralMainMemberFuneralSumAssured),
		"spouse_sum_assured":         quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralSpouseFuneralSumAssured),
		"child_sum_assured":          quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralChildrenFuneralSumAssured),
		"max_children":               strconv.Itoa(cat.FamilyFuneralMaxNumberChildren),
		"parent_sum_assured":         quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralParentFuneralSumAssured),
		"dependant_sum_assured":      quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralAdultDependantSumAssured),
		"max_dependants":             strconv.Itoa(cat.FamilyFuneralMaxNumberAdultDependants),
	}
}

// --- formatting helpers ---

func formatPercent(decimal float64) string {
	return fmt.Sprintf("%.2f%%", decimal*100)
}

// salaryMultiple returns the category's multiple formatted for display, or
// "varies" when the quote uses per-member multiples rather than a global one
// (matching the frontend's behaviour in buildBenefitDefinitionRows).
func salaryMultiple(useGlobal bool, multiple float64) string {
	if !useGlobal {
		return "varies"
	}
	// Trim trailing zeros — "3", "3.5", "2.25"
	return fmt.Sprintf("%g", multiple)
}

// orDash returns "-" when the input is empty, otherwise the input unchanged.
func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

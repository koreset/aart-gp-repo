package quote_template

import (
	"fmt"
	"strconv"

	"api/models"
	"api/services/quote_docx"
)

// Field is one template token. A single []Field slice drives both the
// runtime Context (via fieldsToMap) and the sample template's reference
// tables (via the builders in sample.go). To add a new token, add one
// Field entry in the appropriate *Fields function below; both the render
// engine and the self-documenting sample pick it up automatically.
type Field struct {
	// Key is the token identifier. Inside a "{{#categories}}" block it is
	// referenced bare (e.g. "{{name}}"); nested objects like "insurer" or
	// "gla" reference it with dot syntax (e.g. "{{insurer.name}}").
	Key string

	// Label is the human-readable description shown in the sample's
	// reference tables.
	Label string

	// Value is the resolved runtime value. For nested objects this is
	// populated by fieldsToMap on the nested []Field. Ignored by the
	// sample builder.
	Value interface{}
}

// fieldsToMap folds a []Field into a map[string]interface{}, preserving
// insertion semantics (template engines look up by key, not order).
func fieldsToMap(fs []Field) map[string]interface{} {
	m := make(map[string]interface{}, len(fs))
	for _, f := range fs {
		m[f.Key] = f.Value
	}
	return m
}

// benefitFlags captures which benefit objects are populated for a given
// category. Pre-computing this once avoids re-deriving the flags in each
// per-benefit Fields function and keeps the "has_*" category bools in
// lock-step with the benefit sub-objects.
type benefitFlags struct {
	GLA, SGLA, PTD, CI, PHI, TTD, Funeral bool
}

func deriveBenefitFlags(s models.MemberRatingResultSummary, cat models.SchemeCategory) benefitFlags {
	return benefitFlags{
		GLA:     s.TotalGlaCappedSumAssured > 0,
		SGLA:    s.TotalSglaCappedSumAssured > 0,
		PTD:     s.TotalPtdCappedSumAssured > 0,
		CI:      s.TotalCiCappedSumAssured > 0,
		PHI:     s.TotalPhiCappedIncome > 0,
		TTD:     s.TotalTtdCappedIncome > 0,
		Funeral: s.TotalFunAnnualOfficePremium > 0 || cat.FamilyFuneralBenefit,
	}
}

// ---------------------------------------------------------------------------
// Quote-level scope
// ---------------------------------------------------------------------------

// quoteFields returns the root-scope tokens (e.g. {{scheme_name}}).
func quoteFields(
	quote models.GroupPricingQuote,
	totals quote_docx.QuoteTotals,
	hasNonFuneral bool,
) []Field {
	return []Field{
		{Key: "quote_name", Label: "Quote Name (alias)", Value: quote.QuoteName},
		{Key: "quote_number", Label: "Quote Number", Value: quote.QuoteName},
		{Key: "scheme_name", Label: "Scheme Name", Value: quote.SchemeName},
		{Key: "creation_date", Label: "Creation Date", Value: quote_docx.FormatQuoteDate(quote.CreationDate)},
		{Key: "commencement_date", Label: "Commencement Date", Value: quote_docx.FormatQuoteDate(quote.CommencementDate)},
		{Key: "industry", Label: "Industry", Value: quote.Industry},
		{Key: "currency", Label: "Currency", Value: quote.Currency},
		{Key: "free_cover_limit", Label: "Free Cover Limit", Value: quote_docx.RoundUpToTwoDecimalsAccounting(quote.FreeCoverLimit)},
		{Key: "normal_retirement_age", Label: "Normal Retirement Age", Value: strconv.Itoa(quote.NormalRetirementAge)},
		{Key: "obligation_type", Label: "Obligation Type", Value: quote.ObligationType},
		{Key: "quote_type", Label: "Quote Type", Value: quote.QuoteType},
		{Key: "use_global_salary_multiple", Label: "Use Global Salary Multiple", Value: quote.UseGlobalSalaryMultiple},
		{Key: "total_lives", Label: "Total Lives Covered", Value: strconv.Itoa(totals.TotalLives)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalSumAssured)},
		{Key: "total_annual_salary", Label: "Total Annual Salary", Value: quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualSalary)},
		{Key: "total_annual_premium", Label: "Total Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualPremium)},
		{Key: "has_non_funeral_benefits", Label: "Has Non-Funeral Benefits (flag)", Value: hasNonFuneral},
	}
}

// ---------------------------------------------------------------------------
// Insurer scope ({{insurer.*}})
// ---------------------------------------------------------------------------

func insurerFields(i models.GroupPricingInsurerDetail) []Field {
	return []Field{
		{Key: "name", Label: "Insurer Name", Value: i.Name},
		{Key: "contact_person", Label: "Contact Person", Value: i.ContactPerson},
		{Key: "address_line_1", Label: "Address Line 1", Value: i.AddressLine1},
		{Key: "address_line_2", Label: "Address Line 2", Value: i.AddressLine2},
		{Key: "address_line_3", Label: "Address Line 3", Value: i.AddressLine3},
		{Key: "city", Label: "City", Value: i.City},
		{Key: "province", Label: "Province", Value: i.Province},
		{Key: "post_code", Label: "Post Code", Value: i.PostCode},
		{Key: "country", Label: "Country", Value: i.Country},
		{Key: "telephone", Label: "Telephone", Value: i.Telephone},
		{Key: "email", Label: "Email", Value: i.Email},
		{Key: "introductory_text", Label: "Introductory Text", Value: i.IntroductoryText},
		{Key: "general_provisions_text", Label: "General Provisions Text", Value: i.GeneralProvisionsText},
	}
}

// ---------------------------------------------------------------------------
// Category scope (inside {{#categories}})
// ---------------------------------------------------------------------------

// categoryScalarFields returns the non-bool category tokens. Rendered in
// the sample as a key/value table.
func categoryScalarFields(
	s models.MemberRatingResultSummary,
	cat models.SchemeCategory,
) []Field {
	return []Field{
		{Key: "name", Label: "Category Name", Value: s.Category},
		{Key: "region", Label: "Region", Value: cat.Region},
		{Key: "member_count", Label: "Member Count", Value: strconv.Itoa(int(s.MemberCount))},
		{Key: "total_salary", Label: "Total Annual Salary", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalAnnualSalary)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalSumAssured)},
		{Key: "annual_premium", Label: "Annual Premium (excl. funeral)", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalAnnualPremiumExclFuneral)},
		{Key: "percent_salary", Label: "Premium as % of Salary", Value: formatPercent(s.ProportionExpTotalPremiumExclFuneralSalary)},
		{Key: "free_cover_limit", Label: "Free Cover Limit (category override)", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FreeCoverLimit)},
		{Key: "gla_rate_per_1000", Label: "GLA Rate per 1,000 SA", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpGlaOfficeRatePer1000SA)},
		{Key: "sgla_rate_per_1000", Label: "SGLA Rate per 1,000 SA", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpSglaOfficeRatePer1000SA)},
		{Key: "ptd_rate_per_1000", Label: "PTD Rate per 1,000 SA", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpPtdOfficeRatePer1000SA)},
		{Key: "ci_rate_per_1000", Label: "CI Rate per 1,000 SA", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpCiOfficeRatePer1000SA)},
		{Key: "retirement_premium_waiver", Label: "Retirement Premium Waiver", Value: orDash(cat.PhiPremiumWaiver)},
		{Key: "medical_aid_premium_waiver", Label: "Medical Aid Premium Waiver", Value: orDash(cat.PhiMedicalAidPremiumWaiver)},
		{Key: "gla_terminal_illness_benefit", Label: "GLA Terminal Illness Benefit", Value: orDash(cat.GlaTerminalIllnessBenefit)},
		{Key: "gla_educator_benefit", Label: "GLA Educator Benefit", Value: orDash(cat.GlaEducatorBenefit)},
		{Key: "ptd_educator_benefit", Label: "PTD Educator Benefit", Value: orDash(cat.PtdEducatorBenefit)},
	}
}

// categoryBoolFields returns the has_* flags. Rendered in the sample as
// bullet points demonstrating conditional-block syntax.
func categoryBoolFields(
	s models.MemberRatingResultSummary,
	flags benefitFlags,
) []Field {
	return []Field{
		{Key: "has_non_funeral_benefits", Label: "Category has any non-funeral benefit", Value: quote_docx.CategoryHasNonFuneralBenefits(s)},
		{Key: "has_gla", Label: "Category has GLA", Value: flags.GLA},
		{Key: "has_sgla", Label: "Category has SGLA", Value: flags.SGLA},
		{Key: "has_ptd", Label: "Category has PTD", Value: flags.PTD},
		{Key: "has_ci", Label: "Category has CI", Value: flags.CI},
		{Key: "has_phi", Label: "Category has PHI", Value: flags.PHI},
		{Key: "has_ttd", Label: "Category has TTD", Value: flags.TTD},
		{Key: "has_funeral", Label: "Category has Funeral", Value: flags.Funeral},
	}
}

// ---------------------------------------------------------------------------
// Per-benefit scopes (inside category: gla, sgla, ptd, ci, phi, ttd, funeral)
// ---------------------------------------------------------------------------
//
// Each *Fields function returns the full list of tokens the benefit
// exposes. When `has` is false (benefit not present on the category) the
// runtime path still stores an empty map — matching the prior behaviour
// that missing tokens render to empty strings rather than leaking the
// placeholder text. Callers decide whether to fold the fields or drop to
// an empty map via benefitMap().

func benefitMap(has bool, fs []Field) map[string]interface{} {
	if !has {
		return map[string]interface{}{}
	}
	return fieldsToMap(fs)
}

func glaFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.GlaBenefitTitle},
		{Key: "salary_multiple", Label: "Salary Multiple", Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.GlaSalaryMultiple)},
		{Key: "waiting_period", Label: "Waiting Period (months)", Value: strconv.Itoa(cat.GlaWaitingPeriod)},
		{Key: "benefit_structure", Label: "Benefit Structure", Value: "standalone"},
		{Key: "benefit_type", Label: "Benefit Type", Value: orDash(cat.GlaBenefitType)},
		{Key: "terminal_illness_benefit", Label: "Terminal Illness Benefit", Value: orDash(cat.GlaTerminalIllnessBenefit)},
		{Key: "educator_benefit", Label: "Educator Benefit", Value: orDash(cat.GlaEducatorBenefit)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalGlaCappedSumAssured)},
		{Key: "annual_premium", Label: "Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalGlaAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(s.ExpProportionGlaOfficePremiumSalary)},
	}
}

func sglaFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.SglaBenefitTitle},
		{Key: "salary_multiple", Label: "Salary Multiple", Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.SglaSalaryMultiple)},
		{Key: "waiting_period", Label: "Waiting Period (months)", Value: "0"},
		{Key: "benefit_structure", Label: "Benefit Structure", Value: "rider"},
		{Key: "max_benefit", Label: "Maximum Benefit", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.SglaMaxBenefit)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalSglaCappedSumAssured)},
		{Key: "annual_premium", Label: "Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalSglaAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(s.ExpProportionSglaOfficePremiumSalary)},
	}
}

func ptdFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.PtdBenefitTitle},
		{Key: "salary_multiple", Label: "Salary Multiple", Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.PtdSalaryMultiple)},
		{Key: "waiting_period", Label: "Waiting Period (months)", Value: "0"},
		{Key: "deferred_period", Label: "Deferred Period (months)", Value: strconv.Itoa(cat.PtdDeferredPeriod)},
		{Key: "benefit_type", Label: "Benefit Type", Value: orDash(cat.PtdBenefitType)},
		{Key: "disability_definition", Label: "Disability Definition", Value: orDash(cat.PtdDisabilityDefinition)},
		{Key: "risk_type", Label: "Risk Type", Value: orDash(cat.PtdRiskType)},
		{Key: "educator_benefit", Label: "Educator Benefit", Value: orDash(cat.PtdEducatorBenefit)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalPtdCappedSumAssured)},
		{Key: "annual_premium", Label: "Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalPtdAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(s.ExpProportionPtdOfficePremiumSalary)},
	}
}

func ciFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.CiBenefitTitle},
		{Key: "salary_multiple", Label: "Salary Multiple", Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.CiCriticalIllnessSalaryMultiple)},
		{Key: "waiting_period", Label: "Waiting Period (months)", Value: "0"},
		{Key: "deferred_period", Label: "Deferred Period (months)", Value: "0"},
		{Key: "benefit_structure", Label: "Benefit Structure", Value: orDash(cat.CiBenefitStructure)},
		{Key: "benefit_definition", Label: "Benefit Definition", Value: orDash(cat.CiBenefitDefinition)},
		{Key: "max_benefit", Label: "Maximum Benefit", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.CiMaxBenefit)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalCiCappedSumAssured)},
		{Key: "annual_premium", Label: "Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalCiAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(s.ExpProportionCiOfficePremiumSalary)},
	}
}

func phiFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.PhiBenefitTitle},
		{Key: "income_replacement_percentage", Label: "Income Replacement %", Value: fmt.Sprintf("%.2f%%", cat.PhiIncomeReplacementPercentage)},
		{Key: "waiting_period", Label: "Waiting Period (months)", Value: strconv.Itoa(cat.PhiWaitingPeriod)},
		{Key: "deferred_period", Label: "Deferred Period (months)", Value: strconv.Itoa(cat.PhiDeferredPeriod)},
		{Key: "disability_definition", Label: "Disability Definition", Value: orDash(cat.PhiDisabilityDefinition)},
		{Key: "risk_type", Label: "Risk Type", Value: orDash(cat.PhiRiskType)},
		{Key: "premium_waiver", Label: "Premium Waiver", Value: orDash(cat.PhiPremiumWaiver)},
		{Key: "medical_aid_premium_waiver", Label: "Medical Aid Premium Waiver", Value: orDash(cat.PhiMedicalAidPremiumWaiver)},
		{Key: "benefit_escalation", Label: "Benefit Escalation", Value: orDash(cat.PhiBenefitEscalation)},
		{Key: "total_sum_assured", Label: "Total Covered Income", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalPhiCappedIncome)},
		{Key: "annual_premium", Label: "Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalPhiAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(s.ExpProportionPhiOfficePremiumSalary)},
	}
}

func ttdFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.TtdBenefitTitle},
		{Key: "income_replacement_percentage", Label: "Income Replacement %", Value: fmt.Sprintf("%.2f%%", cat.TtdIncomeReplacementPercentage)},
		{Key: "waiting_period", Label: "Waiting Period (months)", Value: strconv.Itoa(cat.TtdWaitingPeriod)},
		{Key: "deferred_period", Label: "Deferred Period (months)", Value: strconv.Itoa(cat.TtdDeferredPeriod)},
		{Key: "disability_definition", Label: "Disability Definition", Value: orDash(cat.TtdDisabilityDefinition)},
		{Key: "risk_type", Label: "Risk Type", Value: orDash(cat.TtdRiskType)},
		{Key: "total_sum_assured", Label: "Total Covered Income", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalTtdCappedIncome)},
		{Key: "annual_premium", Label: "Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalTtdAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(s.ExpProportionTtdOfficePremiumSalary)},
	}
}

func funeralFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.FamilyFuneralBenefitTitle},
		{Key: "monthly_premium_per_member", Label: "Monthly Premium per Member", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunMonthlyPremiumPerMember)},
		{Key: "annual_premium_per_member", Label: "Annual Premium per Member", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunAnnualPremiumPerMember)},
		{Key: "total_annual_premium", Label: "Total Annual Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalFunAnnualOfficePremium)},
		{Key: "main_member_sum_assured", Label: "Main Member Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralMainMemberFuneralSumAssured)},
		{Key: "spouse_sum_assured", Label: "Spouse Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralSpouseFuneralSumAssured)},
		{Key: "child_sum_assured", Label: "Child Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralChildrenFuneralSumAssured)},
		{Key: "max_children", Label: "Max Children Covered", Value: strconv.Itoa(cat.FamilyFuneralMaxNumberChildren)},
		{Key: "parent_sum_assured", Label: "Parent Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralParentFuneralSumAssured)},
		{Key: "dependant_sum_assured", Label: "Dependant Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralAdultDependantSumAssured)},
		{Key: "max_dependants", Label: "Max Dependants Covered", Value: strconv.Itoa(cat.FamilyFuneralMaxNumberAdultDependants)},
	}
}

// ---------------------------------------------------------------------------
// BenefitSpec: a small table driving sample-template rendering of each
// benefit sub-section. The fields are resolved with a zero fixture when
// BuildSampleTemplate runs, so Label/Key pairs are what matter.
// ---------------------------------------------------------------------------

// BenefitSpec describes one nested benefit object (gla, sgla, ...) for the
// sample template. The Fields function is invoked with a zero fixture so
// its Value outputs can be discarded.
type BenefitSpec struct {
	Prefix string // "gla", "sgla", ...
	Title  string // "Group Life Assurance (GLA)"
	Fields func() []Field
}

// benefitSpecsForSample returns one spec per benefit, ordered to match
// the legacy sample layout. Each Fields closure passes zero-value inputs
// because the sample only needs Keys/Labels.
func benefitSpecsForSample() []BenefitSpec {
	var (
		zs models.MemberRatingResultSummary
		zc models.SchemeCategory
		zq models.GroupPricingQuote
		zt quote_docx.BenefitTitles
	)
	return []BenefitSpec{
		{Prefix: "gla", Title: "Group Life Assurance (GLA)", Fields: func() []Field { return glaFields(zs, zc, zq, zt) }},
		{Prefix: "sgla", Title: "Spouse Group Life (SGLA)", Fields: func() []Field { return sglaFields(zs, zc, zq, zt) }},
		{Prefix: "ptd", Title: "Permanent Total Disability (PTD)", Fields: func() []Field { return ptdFields(zs, zc, zq, zt) }},
		{Prefix: "ci", Title: "Critical Illness (CI)", Fields: func() []Field { return ciFields(zs, zc, zq, zt) }},
		{Prefix: "phi", Title: "Permanent Health Insurance (PHI)", Fields: func() []Field { return phiFields(zs, zc, zt) }},
		{Prefix: "ttd", Title: "Temporary Total Disability (TTD)", Fields: func() []Field { return ttdFields(zs, zc, zt) }},
		{Prefix: "funeral", Title: "Group Funeral", Fields: func() []Field { return funeralFields(zs, zc, zt) }},
	}
}

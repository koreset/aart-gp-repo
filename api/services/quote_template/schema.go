package quote_template

import (
	"fmt"
	"strconv"
	"strings"

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

// ---------------------------------------------------------------------------
// Benefit naming — resolves per-insurer customisations (alias title + alias
// code) so token keys/labels use the names the client has configured. When
// no customisation exists, canonical defaults mirror the base benefit
// mappers in services.getBaseBenefitMaps.
// ---------------------------------------------------------------------------

// benefitName carries the resolved code + title for one benefit. Code is the
// snake_case prefix used in token keys (e.g. "gla" or a customised "gl").
// ShortCode is an optional abbreviated variant used in a few long-key
// contexts (binder/outsource/commission amounts for Additional Accidental
// GLA, whose full code would otherwise create ungainly keys). Title is the
// human-readable name shown in labels.
type benefitName struct {
	Code      string
	ShortCode string
	Title     string
}

// Short returns ShortCode when set, else Code. Callers use this for
// abbreviated key forms.
func (b benefitName) Short() string {
	if b.ShortCode != "" {
		return b.ShortCode
	}
	return b.Code
}

// benefitNaming is the resolved set of code+title pairs for every benefit
// referenced by the category-scope token schema.
type benefitNaming struct {
	GLA, SGLA, PTD, CI, PHI, TTD, Funeral  benefitName
	AdditionalAccidentalGla, AdditionalGla benefitName
	GlaEducator, PtdEducator               benefitName
	TaxSaver, ExtendedFamily               benefitName
}

// defaultBenefitNaming returns the canonical defaults used when no DB-sourced
// customisation is available (sample template generation in tests, etc.).
// Defaults match the historical token prefixes so templates written against
// an un-customised deployment continue to resolve.
func defaultBenefitNaming() benefitNaming {
	return benefitNaming{
		GLA:                     benefitName{Code: "gla", Title: "Group Life Assurance"},
		SGLA:                    benefitName{Code: "sgla", Title: "Spouse Group Life Assurance"},
		PTD:                     benefitName{Code: "ptd", Title: "Permanent Total Disability"},
		CI:                      benefitName{Code: "ci", Title: "Critical Illness"},
		PHI:                     benefitName{Code: "phi", Title: "Personal Health Insurance"},
		TTD:                     benefitName{Code: "ttd", Title: "Temporary Total Disability"},
		Funeral:                 benefitName{Code: "fun", Title: "Group Family Funeral"},
		AdditionalAccidentalGla: benefitName{Code: "additional_accidental_gla", ShortCode: "add_acc_gla", Title: "Additional Accidental Group Life Assurance"},
		AdditionalGla:           benefitName{Code: "additional_gla", Title: "Additional Group Life Assurance"},
		GlaEducator:             benefitName{Code: "gla_educator", Title: "GLA Educator"},
		PtdEducator:             benefitName{Code: "ptd_educator", Title: "PTD Educator"},
		TaxSaver:                benefitName{Code: "tax_saver", Title: "Tax Saver"},
		ExtendedFamily:          benefitName{Code: "extended_family", Title: "Extended Family Funeral"},
	}
}

// resolveBenefitNaming layers customisations from GroupBenefitMapper rows
// onto the defaults. BenefitAliasCode overrides Code (and ShortCode) so the
// customised code flows through every token uniformly; BenefitAlias
// overrides Title. Tax saver and extended family have no mapper entries —
// they always use the defaults.
func resolveBenefitNaming(maps []models.GroupBenefitMapper) benefitNaming {
	out := defaultBenefitNaming()
	for _, m := range maps {
		target := benefitTargetFor(m.BenefitCode, &out)
		if target == nil {
			continue
		}
		if code := sanitiseCode(m.BenefitAliasCode); code != "" {
			target.Code = code
			target.ShortCode = code
		}
		if alias := strings.TrimSpace(m.BenefitAlias); alias != "" {
			target.Title = alias
		}
	}
	return out
}

func benefitTargetFor(benefitCode string, n *benefitNaming) *benefitName {
	switch benefitCode {
	case "GLA":
		return &n.GLA
	case "SGLA":
		return &n.SGLA
	case "PTD":
		return &n.PTD
	case "CI":
		return &n.CI
	case "PHI":
		return &n.PHI
	case "TTD":
		return &n.TTD
	case "GFF":
		return &n.Funeral
	case "AAGLA":
		return &n.AdditionalAccidentalGla
	case "AGLA":
		return &n.AdditionalGla
	case "GLA_EDU":
		return &n.GlaEducator
	case "PTD_EDU":
		return &n.PtdEducator
	}
	return nil
}

// sanitiseCode normalises a customised code into a valid snake_case token
// prefix: trims whitespace, lowercases, and collapses internal whitespace
// runs to underscores.
func sanitiseCode(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return ""
	}
	return strings.Join(strings.Fields(s), "_")
}

// benefitFlags captures which benefit objects are populated for a given
// category. Pre-computing this once avoids re-deriving the flags in each
// per-benefit Fields function and keeps the "has_*" category bools in
// lock-step with the benefit sub-objects.
type benefitFlags struct {
	GLA, SGLA, PTD, CI, PHI, TTD, Funeral bool

	// Additional indicators exposed only as has_* bools (no nested scope).
	AdditionalGla, AdditionalAccidentalGla bool
	GlaEducator, PtdEducator               bool
	ExtendedFamily, TaxSaver               bool
}

func deriveBenefitFlags(s models.MemberRatingResultSummary, cat models.SchemeCategory) benefitFlags {
	return benefitFlags{
		GLA:                     s.TotalGlaCappedSumAssured > 0,
		SGLA:                    s.TotalSglaCappedSumAssured > 0,
		PTD:                     s.TotalPtdCappedSumAssured > 0,
		CI:                      s.TotalCiCappedSumAssured > 0,
		PHI:                     s.TotalPhiCappedIncome > 0,
		TTD:                     s.TotalTtdCappedIncome > 0,
		Funeral:                 models.ComputeOfficePremium(s.TotalFunAnnualRiskPremium, &s) > 0 || cat.FamilyFuneralBenefit,
		AdditionalGla:           s.AdditionalGlaCoverBenefit,
		AdditionalAccidentalGla: s.TotalAdditionalAccidentalGlaCappedSumAssured > 0,
		GlaEducator:             models.ComputeOfficePremium(s.TotalGlaEducatorRiskPremium, &s) > 0 || cat.GlaEducatorBenefit != "",
		PtdEducator:             models.ComputeOfficePremium(s.TotalPtdEducatorRiskPremium, &s) > 0 || cat.PtdEducatorBenefit != "",
		ExtendedFamily:          s.ExtendedFamilyBenefit,
		TaxSaver:                s.TaxSaverBenefit || models.ComputeOfficePremium(s.ExpTotalTaxSaverAnnualRiskPremium, &s) > 0,
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
		{Key: "total_salary", Label: "Total Annual Salary", Value: quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualSalary)},
		{Key: "total_premium", Label: "Total Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualPremium)},
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
		{Key: "on_risk_letter_text", Label: "On-Risk Letter Text", Value: i.OnRiskLetterText},
	}
}

// ---------------------------------------------------------------------------
// Category scope (inside {{#categories}})
// ---------------------------------------------------------------------------

// categoryScalarFields returns the non-bool category tokens. The top of the
// slice carries truly category-scope values (name, totals, scheme commission,
// premium-waiver flags, premium-excluding-funeral aggregates); everything
// else is delegated to the per-benefit token aggregators below so each
// benefit's tokens (rate-per-1000 scalar, indicators, rating block, binder
// / outsource / commission fees, basic premium, conversion / continuity
// slices, and educator splits where applicable) live in one discoverable
// function rather than scattered across multiple multi-benefit aggregators.
func categoryScalarFields(
	s models.MemberRatingResultSummary,
	cat models.SchemeCategory,
	q models.GroupPricingQuote,
	n benefitNaming,
) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	fs := []Field{
		{Key: "name", Label: "Category Name", Value: s.Category},
		{Key: "region", Label: "Region", Value: cat.Region},
		{Key: "member_count", Label: "Member Count", Value: strconv.Itoa(int(s.MemberCount))},
		{Key: "total_salary", Label: "Total Annual Salary", Value: money(s.TotalAnnualSalary)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: money(s.TotalSumAssured)},
		{Key: "premium", Label: "Premium (excl. funeral)", Value: money(s.FinalTotalAnnualPremiumExclFuneral)},
		{Key: "percent_salary", Label: "Premium as % of Salary", Value: formatPercent(proportionOfSalary(s.FinalTotalAnnualPremiumExclFuneral, s.TotalAnnualSalary))},
		{Key: "free_cover_limit", Label: "Free Cover Limit (category override)", Value: money(cat.FreeCoverLimit)},
		{Key: "normal_retirement_age", Label: "Normal Retirement Age", Value: strconv.Itoa(q.NormalRetirementAge)},
		{Key: "retirement_premium_waiver", Label: "Retirement Premium Waiver", Value: orDash(cat.PhiPremiumWaiver)},
		{Key: "medical_aid_premium_waiver", Label: "Medical Aid Premium Waiver", Value: orDash(cat.PhiMedicalAidPremiumWaiver)},
		{Key: "total_premium_excl_funeral", Label: "Total Premium (excluding Funeral)", Value: money(s.FinalTotalAnnualPremiumExclFuneral)},
		{Key: "proportion_total_premium_excl_funeral_salary", Label: "Total Premium (excluding Funeral) as % of Salary", Value: formatPercent(proportionOfSalary(s.FinalTotalAnnualPremiumExclFuneral, s.TotalAnnualSalary))},
		{Key: "scheme_total_commission", Label: "Scheme Total Commission", Value: money(s.FinalSchemeTotalCommission)},
		{Key: "scheme_total_commission_rate", Label: "Scheme Total Commission Rate", Value: formatPercent(s.FinalSchemeTotalCommissionRate)},
	}
	fs = append(fs, glaCategoryTokens(s, cat, q, n)...)
	fs = append(fs, glaEducatorCategoryTokens(s, cat, n)...)
	fs = append(fs, additionalAccidentalGlaCategoryTokens(s, n)...)
	fs = append(fs, taxSaverCategoryTokens(s, n)...)
	fs = append(fs, ptdCategoryTokens(s, cat, q, n)...)
	fs = append(fs, ptdEducatorCategoryTokens(s, cat, n)...)
	fs = append(fs, ciCategoryTokens(s, cat, q, n)...)
	fs = append(fs, sglaCategoryTokens(s, cat, q, n)...)
	fs = append(fs, ttdCategoryTokens(s, cat, n)...)
	fs = append(fs, phiCategoryTokens(s, cat, n)...)
	fs = append(fs, funeralCategoryTokens(s, cat, n)...)
	fs = append(fs, extendedFamilyCategoryTokens(s, n)...)
	return fs
}

// sliceOfficePremium returns a conversion / continuity slice's office premium.
// Prefers the persisted Final*OfficePremium column; falls back to
// ComputeFinalOfficePremium of the slice's ExpAdj*RiskPremium when the
// persisted value is zero (legacy quote pre-dating the column, or a recalc
// whose indicative_rates_count INSERT failed because the migration hadn't
// been applied yet) so the doc still renders meaningful values.
func sliceOfficePremium(persisted, expAdjRisk float64, s *models.MemberRatingResultSummary) float64 {
	if persisted > 0 {
		return persisted
	}
	return models.ComputeFinalOfficePremium(expAdjRisk, s)
}

// ---------------------------------------------------------------------------
// Per-benefit category-token aggregators
//
// Each function below returns every category-scope token for one benefit:
// the optional rate-per-1000 / indicator scalars, the rating block, binder /
// outsource / commission fees, basic premium, and any conversion / continuity
// slices that apply. Token keys, labels, and values are unchanged from the
// previous scattered emissions — this grouping is for findability so a reader
// looking for "all GLA tokens" or "all GLA Educator tokens" finds them in
// one place.
// ---------------------------------------------------------------------------

// glaCategoryTokens emits every category-scope token for GLA: the legacy
// rate-per-1000 scalar, the terminal-illness indicator, salary multiple,
// max cover age, the rating block, binder / outsource / commission fees,
// basic premium, and the three conversion / continuity slices (conv_on_wdr,
// conv_on_ret, cont_dur_dis).
func glaCategoryTokens(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	b := n.GLA
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_rate_per_1000", b.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(s.FinalGlaAnnualOfficePremium, s.TotalGlaCappedSumAssured))},
		Field{Key: fmt.Sprintf("%s_terminal_illness_benefit", b.Code), Label: fmt.Sprintf("%s Terminal Illness Benefit", b.Title), Value: orDash(cat.GlaTerminalIllnessBenefit)},
		Field{Key: fmt.Sprintf("%s_salary_multiple", b.Code), Label: fmt.Sprintf("%s Salary Multiple", b.Title), Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.GlaSalaryMultiple)},
		Field{Key: fmt.Sprintf("%s_max_cover_age", b.Code), Label: fmt.Sprintf("%s Max Cover Age", b.Title), Value: strconv.Itoa(s.GlaMaxCoverAge)},
	)
	fs = append(fs, benefitRatingBlock(b, false,
		s.TotalGlaCappedSumAssured, s.ExpTotalGlaRiskRate, s.ExpTotalGlaAnnualRiskPremium,
		s.ExpGlaRiskRatePer1000SA, s.ExpProportionGlaAnnualRiskPremiumSalary,
		s.FinalGlaAnnualOfficePremium, officeRateFromFinalPremium(s.FinalGlaAnnualOfficePremium, s.TotalGlaCappedSumAssured), officeProportionFromFinalPremium(s.FinalGlaAnnualOfficePremium, &s))...)
	fs = append(fs, benefitBinderOutsourceBlock(b, s.FinalGlaAnnualBinderAmount, s.FinalGlaAnnualOutsourcedAmount, s.TotalGlaCappedSumAssured, &s, false)...)
	fs = append(fs, benefitCommissionBlock(b, s.FinalGlaAnnualCommissionAmount, s.TotalGlaCappedSumAssured, &s, false)...)
	fs = append(fs, benefitBasicPremiumBlock(b,
		basicPremium(s.FinalGlaAnnualOfficePremium,
			s.FinalGlaAnnualCommissionAmount,
			s.FinalGlaAnnualBinderAmount,
			s.FinalGlaAnnualOutsourcedAmount,
			s.FinalTaxSaverAnnualOfficePremium,
			s.FinalGlaConversionOnWithdrawalOfficePremium,
			s.FinalGlaConversionOnRetirementOfficePremium,
			s.FinalGlaContinuityDuringDisabilityOfficePremium),
		s.TotalGlaCappedSumAssured, &s, false)...)
	{
		p := sliceOfficePremium(s.FinalGlaConversionOnWithdrawalOfficePremium, s.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionGlaConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpGlaConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalGlaCappedSumAssured))...)
	}
	{
		p := sliceOfficePremium(s.FinalGlaConversionOnRetirementOfficePremium, s.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_ret", "Conversion on Retirement",
			s.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium, p,
			s.ExpAdjProportionGlaConversionOnRetirementRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpGlaConversionOnRetirementRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalGlaCappedSumAssured))...)
	}
	{
		p := sliceOfficePremium(s.FinalGlaContinuityDuringDisabilityOfficePremium, s.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "cont_dur_dis", "Continuity During Disability",
			s.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium, p,
			s.ExpAdjProportionGlaContinuityDuringDisabilityRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpGlaContinuityDuringDisabilityRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalGlaCappedSumAssured))...)
	}
	return fs
}

// glaEducatorCategoryTokens emits every category-scope token for GLA Educator:
// the educator-benefit indicator, the educator split block (rating + binder /
// outsource / commission fees + basic premium, all in one helper), and the
// three conversion / continuity slices (conv_on_wdr, conv_on_ret, cont_dur_dis).
// Basic premium = office − commission − binder − outsourced − conv_on_wdr −
// conv_on_ret − continuity_during_disability, clamped at zero.
func glaEducatorCategoryTokens(s models.MemberRatingResultSummary, cat models.SchemeCategory, n benefitNaming) []Field {
	b := n.GlaEducator
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_benefit", b.Code), Label: fmt.Sprintf("%s Benefit", b.Title), Value: orDash(cat.GlaEducatorBenefit)},
		Field{Key: fmt.Sprintf("%s_benefit_type", b.Code), Label: fmt.Sprintf("%s Benefit Type", b.Title), Value: orDash(cat.GlaEducatorBenefitType)},
	)
	basic := basicPremium(s.FinalGlaEducatorAnnualOfficePremium,
		s.FinalGlaEducatorAnnualCommissionAmount,
		s.FinalGlaEducatorAnnualBinderAmount,
		s.FinalGlaEducatorAnnualOutsourcedAmount,
		s.FinalGlaEducatorConversionOnWithdrawalOfficePremium,
		s.FinalGlaEducatorConversionOnRetirementOfficePremium,
		s.FinalGlaEducatorContinuityDuringDisabilityOfficePremium)
	fs = append(fs, educatorSplitBlock(b, &s, s.TotalEducatorSumAssured,
		s.ExpAdjTotalGlaEducatorRiskPremium, s.FinalGlaEducatorAnnualOfficePremium,
		s.ExpAdjProportionGlaEducatorRiskPremiumSalary, officeProportionFromFinalPremium(s.FinalGlaEducatorAnnualOfficePremium, &s),
		s.ExpGlaEducatorRiskRatePer1000SA, officeRateFromFinalPremium(s.FinalGlaEducatorAnnualOfficePremium, s.TotalEducatorSumAssured),
		s.FinalGlaEducatorAnnualBinderAmount, s.FinalGlaEducatorAnnualOutsourcedAmount,
		s.FinalGlaEducatorAnnualCommissionAmount, basic)...)
	{
		p := sliceOfficePremium(s.FinalGlaEducatorConversionOnWithdrawalOfficePremium, s.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionGlaEducatorConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpGlaEducatorConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalEducatorSumAssured))...)
	}
	{
		p := sliceOfficePremium(s.FinalGlaEducatorConversionOnRetirementOfficePremium, s.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_ret", "Conversion on Retirement",
			s.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium, p,
			s.ExpAdjProportionGlaEducatorConversionOnRetirementRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpGlaEducatorConversionOnRetirementRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalEducatorSumAssured))...)
	}
	{
		p := sliceOfficePremium(s.FinalGlaEducatorContinuityDuringDisabilityOfficePremium, s.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "cont_dur_dis", "Continuity During Disability",
			s.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium, p,
			s.ExpAdjProportionGlaEducatorContinuityDuringDisabilityRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpGlaEducatorContinuityDuringDisabilityRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalEducatorSumAssured))...)
	}
	return fs
}

// additionalAccidentalGlaCategoryTokens emits every category-scope token for
// Additional Accidental GLA: a custom-shape rating block (the risk-proportion
// key uses "prop_" mirroring the underlying column name; office-proportion
// uses "proportion_"), binder / outsource / commission fees (using the
// short code), and basic premium. There are no conversion / continuity slices
// for this benefit.
func additionalAccidentalGlaCategoryTokens(s models.MemberRatingResultSummary, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	ab := n.AdditionalAccidentalGla
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("total_%s_capped_sum_assured", ab.Code), Label: fmt.Sprintf("%s — Total Capped Sum Assured", ab.Title), Value: money(s.TotalAdditionalAccidentalGlaCappedSumAssured)},
		Field{Key: fmt.Sprintf("total_%s_risk_rate", ab.Code), Label: fmt.Sprintf("%s — Total Risk Rate", ab.Title), Value: money(s.ExpTotalAdditionalAccidentalGlaRiskRate)},
		Field{Key: fmt.Sprintf("total_%s_risk_premium", ab.Code), Label: fmt.Sprintf("%s — Total Risk Premium", ab.Title), Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium)},
		Field{Key: fmt.Sprintf("%s_risk_rate_per1000_sa", ab.Code), Label: fmt.Sprintf("%s — Risk Rate per 1,000 SA", ab.Title), Value: money(s.ExpAdditionalAccidentalGlaRiskRatePer1000SA)},
		Field{Key: fmt.Sprintf("prop_%s_risk_premium_salary", ab.Code), Label: fmt.Sprintf("%s — Risk Premium as %% of Salary", ab.Title), Value: formatPercent(s.ExpProportionAdditionalAccidentalGlaAnnualRiskPremiumSalary)},
		Field{Key: fmt.Sprintf("total_%s_office_premium", ab.Code), Label: fmt.Sprintf("%s — Total Office Premium", ab.Title), Value: money(s.FinalAdditionalAccidentalGlaAnnualOfficePremium)},
		Field{Key: fmt.Sprintf("%s_office_rate_per1000_sa", ab.Code), Label: fmt.Sprintf("%s — Office Rate per 1,000 SA", ab.Title), Value: money(officeRateFromFinalPremium(s.FinalAdditionalAccidentalGlaAnnualOfficePremium, s.TotalAdditionalAccidentalGlaCappedSumAssured))},
		Field{Key: fmt.Sprintf("proportion_%s_office_premium_salary", ab.Code), Label: fmt.Sprintf("%s — Office Premium as %% of Salary", ab.Title), Value: formatPercent(officeProportionFromFinalPremium(s.FinalAdditionalAccidentalGlaAnnualOfficePremium, &s))},
	)
	fs = append(fs, benefitBinderOutsourceBlock(ab, s.FinalAdditionalAccidentalGlaAnnualBinderAmount, s.FinalAdditionalAccidentalGlaAnnualOutsourcedAmt, s.TotalAdditionalAccidentalGlaCappedSumAssured, &s, true)...)
	fs = append(fs, benefitCommissionBlock(ab, s.FinalAdditionalAccidentalGlaAnnualCommissionAmount, s.TotalAdditionalAccidentalGlaCappedSumAssured, &s, true)...)
	fs = append(fs, benefitBasicPremiumBlock(ab,
		basicPremium(s.FinalAdditionalAccidentalGlaAnnualOfficePremium,
			s.FinalAdditionalAccidentalGlaAnnualCommissionAmount,
			s.FinalAdditionalAccidentalGlaAnnualBinderAmount,
			s.FinalAdditionalAccidentalGlaAnnualOutsourcedAmt),
		s.TotalAdditionalAccidentalGlaCappedSumAssured, &s, true)...)
	return fs
}

// taxSaverCategoryTokens emits the tax saver tokens. Tax saver SA is an
// extra cover layered on top of GLA — there is no aggregated TaxSaverSumAssured
// on MRRS, so the rate-per-1000-SA is denominated against
// TotalGlaCappedSumAssured (the parent benefit's covered amount).
// Proportion-of-salary uses the standard salary × IndicativeRatesCount
// denominator via officeProportionFromFinalPremium.
func taxSaverCategoryTokens(s models.MemberRatingResultSummary, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	tb := n.TaxSaver
	return []Field{
		{Key: fmt.Sprintf("total_%s_risk_premium", tb.Code), Label: fmt.Sprintf("%s — Total Risk Premium", tb.Title), Value: money(s.ExpTotalTaxSaverAnnualRiskPremium)},
		{Key: fmt.Sprintf("total_%s_office_premium", tb.Code), Label: fmt.Sprintf("%s — Total Office Premium", tb.Title), Value: money(s.FinalTaxSaverAnnualOfficePremium)},
		{Key: fmt.Sprintf("%s_office_rate_per1000_sa", tb.Code), Label: fmt.Sprintf("%s — Office Rate per 1,000 SA", tb.Title), Value: money(officeRateFromFinalPremium(s.FinalTaxSaverAnnualOfficePremium, s.TotalGlaCappedSumAssured))},
		{Key: fmt.Sprintf("proportion_%s_office_premium_salary", tb.Code), Label: fmt.Sprintf("%s — Office Premium as %% of Salary", tb.Title), Value: formatPercent(officeProportionFromFinalPremium(s.FinalTaxSaverAnnualOfficePremium, &s))},
	}
}

// ptdCategoryTokens emits every category-scope token for PTD: the legacy
// rate-per-1000 scalar, salary multiple, max cover age, rating block, binder
// / outsource / commission fees, basic premium, and the conv_on_wdr slice.
func ptdCategoryTokens(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	b := n.PTD
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_rate_per_1000", b.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(s.FinalPtdAnnualOfficePremium, s.TotalPtdCappedSumAssured))},
		Field{Key: fmt.Sprintf("%s_salary_multiple", b.Code), Label: fmt.Sprintf("%s Salary Multiple", b.Title), Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.PtdSalaryMultiple)},
		Field{Key: fmt.Sprintf("%s_max_cover_age", b.Code), Label: fmt.Sprintf("%s Max Cover Age", b.Title), Value: strconv.Itoa(s.PtdMaxCoverAge)},
	)
	fs = append(fs, benefitRatingBlock(b, false,
		s.TotalPtdCappedSumAssured, s.ExpTotalPtdRiskRate, s.ExpTotalPtdAnnualRiskPremium,
		s.ExpPtdRiskRatePer1000SA, s.ExpProportionPtdAnnualRiskPremiumSalary,
		s.FinalPtdAnnualOfficePremium, officeRateFromFinalPremium(s.FinalPtdAnnualOfficePremium, s.TotalPtdCappedSumAssured), officeProportionFromFinalPremium(s.FinalPtdAnnualOfficePremium, &s))...)
	fs = append(fs, benefitBinderOutsourceBlock(b, s.FinalPtdAnnualBinderAmount, s.FinalPtdAnnualOutsourcedAmount, s.TotalPtdCappedSumAssured, &s, false)...)
	fs = append(fs, benefitCommissionBlock(b, s.FinalPtdAnnualCommissionAmount, s.TotalPtdCappedSumAssured, &s, false)...)
	fs = append(fs, benefitBasicPremiumBlock(b,
		basicPremium(s.FinalPtdAnnualOfficePremium,
			s.FinalPtdAnnualCommissionAmount,
			s.FinalPtdAnnualBinderAmount,
			s.FinalPtdAnnualOutsourcedAmount,
			s.FinalPtdConversionOnWithdrawalOfficePremium),
		s.TotalPtdCappedSumAssured, &s, false)...)
	{
		p := sliceOfficePremium(s.FinalPtdConversionOnWithdrawalOfficePremium, s.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionPtdConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpPtdConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalPtdCappedSumAssured))...)
	}
	return fs
}

// ptdEducatorCategoryTokens emits every category-scope token for PTD Educator:
// the educator-benefit indicator, the educator split block, and the
// conv_on_wdr + conv_on_ret slices. PTD Educator does not have a
// continuity-during-disability slice. Basic premium = office − commission −
// binder − outsourced − conv_on_wdr − conv_on_ret, clamped at zero.
func ptdEducatorCategoryTokens(s models.MemberRatingResultSummary, cat models.SchemeCategory, n benefitNaming) []Field {
	b := n.PtdEducator
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_benefit", b.Code), Label: fmt.Sprintf("%s Benefit", b.Title), Value: orDash(cat.PtdEducatorBenefit)},
		Field{Key: fmt.Sprintf("%s_benefit_type", b.Code), Label: fmt.Sprintf("%s Benefit Type", b.Title), Value: orDash(cat.PtdEducatorBenefitType)},
	)
	basic := basicPremium(s.FinalPtdEducatorAnnualOfficePremium,
		s.FinalPtdEducatorAnnualCommissionAmount,
		s.FinalPtdEducatorAnnualBinderAmount,
		s.FinalPtdEducatorAnnualOutsourcedAmount,
		s.FinalPtdEducatorConversionOnWithdrawalOfficePremium,
		s.FinalPtdEducatorConversionOnRetirementOfficePremium)
	fs = append(fs, educatorSplitBlock(b, &s, s.TotalEducatorSumAssured,
		s.ExpAdjTotalPtdEducatorRiskPremium, s.FinalPtdEducatorAnnualOfficePremium,
		s.ExpAdjProportionPtdEducatorRiskPremiumSalary, officeProportionFromFinalPremium(s.FinalPtdEducatorAnnualOfficePremium, &s),
		s.ExpPtdEducatorRiskRatePer1000SA, officeRateFromFinalPremium(s.FinalPtdEducatorAnnualOfficePremium, s.TotalEducatorSumAssured),
		s.FinalPtdEducatorAnnualBinderAmount, s.FinalPtdEducatorAnnualOutsourcedAmount,
		s.FinalPtdEducatorAnnualCommissionAmount, basic)...)
	{
		p := sliceOfficePremium(s.FinalPtdEducatorConversionOnWithdrawalOfficePremium, s.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionPtdEducatorConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpPtdEducatorConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalEducatorSumAssured))...)
	}
	{
		p := sliceOfficePremium(s.FinalPtdEducatorConversionOnRetirementOfficePremium, s.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_ret", "Conversion on Retirement",
			s.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium, p,
			s.ExpAdjProportionPtdEducatorConversionOnRetirementRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpPtdEducatorConversionOnRetirementRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalEducatorSumAssured))...)
	}
	return fs
}

// ciCategoryTokens emits every category-scope token for CI: the legacy
// rate-per-1000 scalar, salary multiple, max cover age, rating block, binder
// / outsource / commission fees, basic premium, and the conv_on_wdr slice.
func ciCategoryTokens(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	b := n.CI
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_rate_per_1000", b.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(s.FinalCiAnnualOfficePremium, s.TotalCiCappedSumAssured))},
		Field{Key: fmt.Sprintf("%s_salary_multiple", b.Code), Label: fmt.Sprintf("%s Salary Multiple", b.Title), Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.CiCriticalIllnessSalaryMultiple)},
		Field{Key: fmt.Sprintf("%s_max_cover_age", b.Code), Label: fmt.Sprintf("%s Max Cover Age", b.Title), Value: strconv.Itoa(s.CiMaxCoverAge)},
	)
	fs = append(fs, benefitRatingBlock(b, false,
		s.TotalCiCappedSumAssured, s.ExpTotalCiRiskRate, s.ExpTotalCiAnnualRiskPremium,
		s.ExpCiRiskRatePer1000SA, s.ExpProportionCiAnnualRiskPremiumSalary,
		s.FinalCiAnnualOfficePremium, officeRateFromFinalPremium(s.FinalCiAnnualOfficePremium, s.TotalCiCappedSumAssured), officeProportionFromFinalPremium(s.FinalCiAnnualOfficePremium, &s))...)
	fs = append(fs, benefitBinderOutsourceBlock(b, s.FinalCiAnnualBinderAmount, s.FinalCiAnnualOutsourcedAmount, s.TotalCiCappedSumAssured, &s, false)...)
	fs = append(fs, benefitCommissionBlock(b, s.FinalCiAnnualCommissionAmount, s.TotalCiCappedSumAssured, &s, false)...)
	fs = append(fs, benefitBasicPremiumBlock(b,
		basicPremium(s.FinalCiAnnualOfficePremium,
			s.FinalCiAnnualCommissionAmount,
			s.FinalCiAnnualBinderAmount,
			s.FinalCiAnnualOutsourcedAmount,
			s.FinalCiConversionOnWithdrawalOfficePremium),
		s.TotalCiCappedSumAssured, &s, false)...)
	{
		p := sliceOfficePremium(s.FinalCiConversionOnWithdrawalOfficePremium, s.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionCiConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpCiConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalCiCappedSumAssured))...)
	}
	return fs
}

// sglaCategoryTokens emits every category-scope token for SGLA: the legacy
// rate-per-1000 scalar, salary multiple, rating block, binder / outsource /
// commission fees, basic premium, and the conv_on_wdr slice. SGLA has no
// dedicated max-cover-age column on the summary (it inherits GLA's).
func sglaCategoryTokens(s models.MemberRatingResultSummary, cat models.SchemeCategory, q models.GroupPricingQuote, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	b := n.SGLA
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_rate_per_1000", b.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(s.FinalSglaAnnualOfficePremium, s.TotalSglaCappedSumAssured))},
		Field{Key: fmt.Sprintf("%s_salary_multiple", b.Code), Label: fmt.Sprintf("%s Salary Multiple", b.Title), Value: salaryMultiple(q.UseGlobalSalaryMultiple, cat.SglaSalaryMultiple)},
	)
	fs = append(fs, benefitRatingBlock(b, false,
		s.TotalSglaCappedSumAssured, s.ExpTotalSglaRiskRate, s.ExpTotalSglaAnnualRiskPremium,
		s.ExpSglaRiskRatePer1000SA, s.ExpProportionSglaAnnualRiskPremiumSalary,
		s.FinalSglaAnnualOfficePremium, officeRateFromFinalPremium(s.FinalSglaAnnualOfficePremium, s.TotalSglaCappedSumAssured), officeProportionFromFinalPremium(s.FinalSglaAnnualOfficePremium, &s))...)
	fs = append(fs, benefitBinderOutsourceBlock(b, s.FinalSglaAnnualBinderAmount, s.FinalSglaAnnualOutsourcedAmount, s.TotalSglaCappedSumAssured, &s, false)...)
	fs = append(fs, benefitCommissionBlock(b, s.FinalSglaAnnualCommissionAmount, s.TotalSglaCappedSumAssured, &s, false)...)
	fs = append(fs, benefitBasicPremiumBlock(b,
		basicPremium(s.FinalSglaAnnualOfficePremium,
			s.FinalSglaAnnualCommissionAmount,
			s.FinalSglaAnnualBinderAmount,
			s.FinalSglaAnnualOutsourcedAmount,
			s.FinalSglaConversionOnWithdrawalOfficePremium),
		s.TotalSglaCappedSumAssured, &s, false)...)
	{
		p := sliceOfficePremium(s.FinalSglaConversionOnWithdrawalOfficePremium, s.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionSglaConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpSglaConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalSglaCappedSumAssured))...)
	}
	return fs
}

// ttdCategoryTokens emits every category-scope token for TTD: max cover age,
// the income-based rating block, binder / outsource / commission fees, basic
// premium, and the conv_on_wdr slice. TTD has no rate-per-1000 scalar at
// category scope, and no salary multiple (income-based benefit).
func ttdCategoryTokens(s models.MemberRatingResultSummary, _ models.SchemeCategory, n benefitNaming) []Field {
	b := n.TTD
	var fs []Field
	fs = append(fs, Field{Key: fmt.Sprintf("%s_max_cover_age", b.Code), Label: fmt.Sprintf("%s Max Cover Age", b.Title), Value: strconv.Itoa(s.TtdMaxCoverAge)})
	fs = append(fs, benefitRatingBlock(b, true,
		s.TotalTtdCappedIncome, s.ExpTotalTtdRiskRate, s.ExpTotalTtdAnnualRiskPremium,
		s.ExpTtdRiskRatePer1000SA, s.ExpProportionTtdAnnualRiskPremiumSalary,
		s.FinalTtdAnnualOfficePremium, officeRateFromFinalPremium(s.FinalTtdAnnualOfficePremium, s.TotalTtdCappedIncome), officeProportionFromFinalPremium(s.FinalTtdAnnualOfficePremium, &s))...)
	fs = append(fs, benefitBinderOutsourceBlock(b, s.FinalTtdAnnualBinderAmount, s.FinalTtdAnnualOutsourcedAmount, s.TotalTtdCappedIncome, &s, false)...)
	fs = append(fs, benefitCommissionBlock(b, s.FinalTtdAnnualCommissionAmount, s.TotalTtdCappedIncome, &s, false)...)
	fs = append(fs, benefitBasicPremiumBlock(b,
		basicPremium(s.FinalTtdAnnualOfficePremium,
			s.FinalTtdAnnualCommissionAmount,
			s.FinalTtdAnnualBinderAmount,
			s.FinalTtdAnnualOutsourcedAmount,
			s.FinalTtdConversionOnWithdrawalOfficePremium),
		s.TotalTtdCappedIncome, &s, false)...)
	{
		p := sliceOfficePremium(s.FinalTtdConversionOnWithdrawalOfficePremium, s.ExpAdjTotalTtdConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalTtdConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionTtdConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpTtdConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalTtdCappedIncome))...)
	}
	return fs
}

// phiCategoryTokens emits every category-scope token for PHI: max cover age,
// PHI-specific normal retirement age, the income-based rating block, binder /
// outsource / commission fees, basic premium, and the conv_on_wdr slice.
func phiCategoryTokens(s models.MemberRatingResultSummary, cat models.SchemeCategory, n benefitNaming) []Field {
	b := n.PHI
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_max_cover_age", b.Code), Label: fmt.Sprintf("%s Max Cover Age", b.Title), Value: strconv.Itoa(s.PhiMaxCoverAge)},
		Field{Key: fmt.Sprintf("%s_normal_retirement_age", b.Code), Label: fmt.Sprintf("%s Normal Retirement Age", b.Title), Value: strconv.Itoa(cat.PhiNormalRetirementAge)},
	)
	fs = append(fs, benefitRatingBlock(b, true,
		s.TotalPhiCappedIncome, s.ExpTotalPhiRiskRate, s.ExpTotalPhiAnnualRiskPremium,
		s.ExpPhiRiskRatePer1000SA, s.ExpProportionPhiAnnualRiskPremiumSalary,
		s.FinalPhiAnnualOfficePremium, officeRateFromFinalPremium(s.FinalPhiAnnualOfficePremium, s.TotalPhiCappedIncome), officeProportionFromFinalPremium(s.FinalPhiAnnualOfficePremium, &s))...)
	fs = append(fs, benefitBinderOutsourceBlock(b, s.FinalPhiAnnualBinderAmount, s.FinalPhiAnnualOutsourcedAmount, s.TotalPhiCappedIncome, &s, false)...)
	fs = append(fs, benefitCommissionBlock(b, s.FinalPhiAnnualCommissionAmount, s.TotalPhiCappedIncome, &s, false)...)
	fs = append(fs, benefitBasicPremiumBlock(b,
		basicPremium(s.FinalPhiAnnualOfficePremium,
			s.FinalPhiAnnualCommissionAmount,
			s.FinalPhiAnnualBinderAmount,
			s.FinalPhiAnnualOutsourcedAmount,
			s.FinalPhiConversionOnWithdrawalOfficePremium),
		s.TotalPhiCappedIncome, &s, false)...)
	{
		p := sliceOfficePremium(s.FinalPhiConversionOnWithdrawalOfficePremium, s.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionPhiConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpPhiConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalPhiCappedIncome))...)
	}
	return fs
}

// funeralCategoryTokens emits every category-scope token for the family
// funeral benefit: max cover age, the custom rating shape (no risk rate, with
// monthly / annual per-member premium splits), binder / outsource / commission
// fees, basic premium, and the conv_on_wdr slice. The "monthly" qualifier
// stays on the monthly-premium-per-member key to disambiguate it from the
// (default) annual sibling, but the "annual" qualifier is dropped elsewhere.
func funeralCategoryTokens(s models.MemberRatingResultSummary, _ models.SchemeCategory, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	b := n.Funeral
	var fs []Field
	fs = append(fs,
		Field{Key: fmt.Sprintf("%s_max_cover_age", b.Code), Label: fmt.Sprintf("%s Max Cover Age", b.Title), Value: strconv.Itoa(s.FunMaxCoverAge)},
		Field{Key: fmt.Sprintf("total_%s_risk_premium", b.Code), Label: fmt.Sprintf("%s — Total Risk Premium", b.Title), Value: money(s.ExpTotalFunAnnualRiskPremium)},
		Field{Key: fmt.Sprintf("proportion_%s_risk_premium_salary", b.Code), Label: fmt.Sprintf("%s — Risk Premium as %% of Salary", b.Title), Value: formatPercent(s.ExpProportionFunAnnualRiskPremiumSalary)},
		Field{Key: fmt.Sprintf("total_%s_office_premium", b.Code), Label: fmt.Sprintf("%s — Total Office Premium", b.Title), Value: money(s.FinalFunAnnualOfficePremium)},
		Field{Key: fmt.Sprintf("proportion_%s_office_premium_salary", b.Code), Label: fmt.Sprintf("%s — Office Premium as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(s.FinalFunAnnualOfficePremium, &s))},
		Field{Key: fmt.Sprintf("total_%s_premium_per_member", b.Code), Label: fmt.Sprintf("%s — Premium per Member", b.Title), Value: money(s.ExpTotalFunAnnualPremiumPerMember)},
		Field{Key: fmt.Sprintf("total_%s_monthly_premium_per_member", b.Code), Label: fmt.Sprintf("%s — Monthly Premium per Member", b.Title), Value: money(s.ExpTotalFunMonthlyPremiumPerMember)},
	)
	fs = append(fs, benefitBinderOutsourceBlock(b, s.FinalFunAnnualBinderAmount, s.FinalFunAnnualOutsourcedAmount, s.TotalFamilyFuneralSumAssured, &s, false)...)
	fs = append(fs, benefitCommissionBlock(b, s.FinalFunAnnualCommissionAmount, s.TotalFamilyFuneralSumAssured, &s, false)...)
	fs = append(fs, benefitBasicPremiumBlock(b,
		basicPremium(s.FinalFunAnnualOfficePremium,
			s.FinalFunAnnualCommissionAmount,
			s.FinalFunAnnualBinderAmount,
			s.FinalFunAnnualOutsourcedAmount,
			s.FinalFunConversionOnWithdrawalOfficePremium),
		s.TotalFamilyFuneralSumAssured, &s, false)...)
	{
		p := sliceOfficePremium(s.FinalFunConversionOnWithdrawalOfficePremium, s.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium, &s)
		fs = append(fs, conversionSliceBlock(b, "conv_on_wdr", "Conversion on Withdrawal",
			s.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium, p,
			s.ExpAdjProportionFunConversionOnWithdrawalRiskPremiumSalary, officeProportionFromFinalPremium(p, &s),
			s.ExpFunConversionOnWithdrawalRiskRatePer1000SA, officeRateFromFinalPremium(p, s.TotalFamilyFuneralSumAssured))...)
	}
	return fs
}

// extendedFamilyCategoryTokens emits every category-scope token for the
// extended family funeral cover: configuration scalars (age band source /
// type, pricing method), the total monthly premium across all bands, and
// the per-band detail as a nested array (`extended_family_bands`) so doc
// templates can iterate via {{#extended_family_bands}}…{{/…}} and render
// per-band sum assured + monthly premium rows.
func extendedFamilyCategoryTokens(s models.MemberRatingResultSummary, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	b := n.ExtendedFamily
	bands := make([]map[string]interface{}, 0, len(s.ExtendedFamilyBandRates))
	for _, r := range s.ExtendedFamilyBandRates {
		bands = append(bands, map[string]interface{}{
			"min_age":                strconv.Itoa(r.MinAge),
			"max_age":                strconv.Itoa(r.MaxAge),
			"sum_assured":            money(r.SumAssured),
			"monthly_premium":        money(r.MonthlyPremium),
			"office_monthly_premium": money(r.OfficeMonthlyPremium),
		})
	}
	return []Field{
		{Key: fmt.Sprintf("%s_age_band_source", b.Code), Label: fmt.Sprintf("%s Age Band Source", b.Title), Value: orDash(s.ExtendedFamilyAgeBandSource)},
		{Key: fmt.Sprintf("%s_age_band_type", b.Code), Label: fmt.Sprintf("%s Age Band Type", b.Title), Value: orDash(s.ExtendedFamilyAgeBandType)},
		{Key: fmt.Sprintf("%s_pricing_method", b.Code), Label: fmt.Sprintf("%s Pricing Method", b.Title), Value: orDash(s.ExtendedFamilyPricingMethod)},
		{Key: fmt.Sprintf("total_%s_monthly_premium", b.Code), Label: fmt.Sprintf("%s — Total Monthly Premium", b.Title), Value: money(s.TotalExtendedFamilyMonthlyPremium)},
		{Key: fmt.Sprintf("%s_bands", b.Code), Label: fmt.Sprintf("%s Bands (iterate with {{#%s_bands}}…{{/%s_bands}}; each item exposes min_age, max_age, sum_assured, monthly_premium, office_monthly_premium)", b.Title, b.Code, b.Code), Value: bands},
	}
}

// benefitRatingBlock returns the 8-field rating block for one benefit:
// capped sum assured (or income), risk rate, risk premium, risk-rate per
// 1,000 SA, risk-proportion-of-salary, office premium, office-rate per
// 1,000 SA, office-proportion-of-salary.
func benefitRatingBlock(b benefitName, incomeBased bool,
	cappedAmount, riskRate, riskPrem, riskRatePer1000, propRiskSalary,
	officePrem, officeRatePer1000, propOfficeSalary float64,
) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	cappedKey := "capped_sum_assured"
	cappedLabel := "Total Capped Sum Assured"
	if incomeBased {
		cappedKey = "capped_income"
		cappedLabel = "Total Capped Income"
	}
	return []Field{
		{Key: fmt.Sprintf("total_%s_%s", b.Code, cappedKey), Label: fmt.Sprintf("%s — %s", b.Title, cappedLabel), Value: money(cappedAmount)},
		{Key: fmt.Sprintf("total_%s_risk_rate", b.Code), Label: fmt.Sprintf("%s — Total Risk Rate", b.Title), Value: money(riskRate)},
		{Key: fmt.Sprintf("total_%s_risk_premium", b.Code), Label: fmt.Sprintf("%s — Total Risk Premium", b.Title), Value: money(riskPrem)},
		{Key: fmt.Sprintf("%s_risk_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Risk Rate per 1,000 SA", b.Title), Value: money(riskRatePer1000)},
		{Key: fmt.Sprintf("proportion_%s_risk_premium_salary", b.Code), Label: fmt.Sprintf("%s — Risk Premium as %% of Salary", b.Title), Value: formatPercent(propRiskSalary)},
		{Key: fmt.Sprintf("total_%s_office_premium", b.Code), Label: fmt.Sprintf("%s — Total Office Premium", b.Title), Value: money(officePrem)},
		{Key: fmt.Sprintf("%s_office_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Office Rate per 1,000 SA", b.Title), Value: money(officeRatePer1000)},
		{Key: fmt.Sprintf("proportion_%s_office_premium_salary", b.Code), Label: fmt.Sprintf("%s — Office Premium as %% of Salary", b.Title), Value: formatPercent(propOfficeSalary)},
	}
}

// benefitBinderOutsourceBlock returns the binder + outsourced tokens for
// one benefit. Each amount also exposes a rate-per-1000-SA and a
// proportion-of-salary derivative so doc templates can render binder /
// outsource fee structures the same way they render the parent office
// premium. When useShort is true, b.Short() is used for the key prefix
// (Add Acc GLA uses the abbreviated form to keep keys tractable). The
// outsourced token's suffix is "_amt" for Add Acc GLA (matches the
// historical attached-list naming) and "_amount" elsewhere.
func benefitBinderOutsourceBlock(b benefitName, binder, outsourced, sumAssured float64, s *models.MemberRatingResultSummary, useShort bool) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	code := b.Code
	outsourcedRoot := "outsourced"
	outsourcedSuffix := "outsourced_amount"
	if useShort {
		code = b.Short()
		// The "_amt" suffix is a quirk of the default Add Acc GLA naming
		// (ShortCode differs from Code). When a customisation is applied
		// the resolver sets ShortCode == Code, so the consistent
		// "_amount" suffix is used.
		if b.ShortCode != "" && b.ShortCode != b.Code {
			outsourcedSuffix = "outsourced_amt"
			outsourcedRoot = "outsourced"
		}
	}
	return []Field{
		{Key: fmt.Sprintf("total_%s_binder_amount", code), Label: fmt.Sprintf("%s — Total Binder Amount", b.Title), Value: money(binder)},
		{Key: fmt.Sprintf("%s_binder_rate_per1000_sa", code), Label: fmt.Sprintf("%s — Binder Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(binder, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_binder_amount_salary", code), Label: fmt.Sprintf("%s — Binder Amount as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(binder, s))},
		{Key: fmt.Sprintf("total_%s_%s", code, outsourcedSuffix), Label: fmt.Sprintf("%s — Total Outsourced Amount", b.Title), Value: money(outsourced)},
		{Key: fmt.Sprintf("%s_%s_rate_per1000_sa", code, outsourcedRoot), Label: fmt.Sprintf("%s — Outsource Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(outsourced, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_%s_amount_salary", code, outsourcedRoot), Label: fmt.Sprintf("%s — Outsource Amount as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(outsourced, s))},
	}
}

// benefitCommissionBlock returns the commission tokens for one benefit:
// the total amount, the rate-per-1000-SA, and the proportion-of-salary
// derivative — mirroring the binder / outsource shape so doc templates
// can render commission alongside the rest of the fee breakdown.
// useShort chooses the abbreviated code (Add Acc GLA).
func benefitCommissionBlock(b benefitName, amount, sumAssured float64, s *models.MemberRatingResultSummary, useShort bool) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	code := b.Code
	if useShort {
		code = b.Short()
	}
	return []Field{
		{Key: fmt.Sprintf("total_%s_commission_amount", code), Label: fmt.Sprintf("%s — Total Commission Amount", b.Title), Value: money(amount)},
		{Key: fmt.Sprintf("%s_commission_rate_per1000_sa", code), Label: fmt.Sprintf("%s — Commission Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(amount, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_commission_amount_salary", code), Label: fmt.Sprintf("%s — Commission Amount as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(amount, s))},
	}
}

// educatorSplitBlock returns the educator-cover block for one educator
// (GLA or PTD). Values come from the experience-adjusted model fields;
// the "adjusted" qualifier is omitted from keys and labels to keep tokens
// short and time-neutral. Binder, outsourced, and commission each expose
// the amount plus a rate-per-1000-SA and a proportion-of-salary derivative
// (using TotalEducatorSumAssured and the shared salary × rates count
// denominators) so doc templates can render the fee structure consistently.
func educatorSplitBlock(b benefitName, s *models.MemberRatingResultSummary, sumAssured float64,
	riskPrem, officePrem, propRiskSalary, propOfficeSalary,
	riskRatePer1000, officeRatePer1000,
	binder, outsourced, commission, basic float64,
) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	return []Field{
		{Key: fmt.Sprintf("total_%s_risk_premium", b.Code), Label: fmt.Sprintf("%s — Total Risk Premium", b.Title), Value: money(riskPrem)},
		{Key: fmt.Sprintf("total_%s_office_premium", b.Code), Label: fmt.Sprintf("%s — Total Office Premium", b.Title), Value: money(officePrem)},
		{Key: fmt.Sprintf("proportion_%s_risk_premium_salary", b.Code), Label: fmt.Sprintf("%s — Risk Premium as %% of Salary", b.Title), Value: formatPercent(propRiskSalary)},
		{Key: fmt.Sprintf("proportion_%s_office_premium_salary", b.Code), Label: fmt.Sprintf("%s — Office Premium as %% of Salary", b.Title), Value: formatPercent(propOfficeSalary)},
		{Key: fmt.Sprintf("%s_risk_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Risk Rate per 1,000 SA", b.Title), Value: money(riskRatePer1000)},
		{Key: fmt.Sprintf("%s_office_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Office Rate per 1,000 SA", b.Title), Value: money(officeRatePer1000)},
		{Key: fmt.Sprintf("total_%s_binder_amount", b.Code), Label: fmt.Sprintf("%s — Total Binder Amount", b.Title), Value: money(binder)},
		{Key: fmt.Sprintf("%s_binder_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Binder Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(binder, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_binder_amount_salary", b.Code), Label: fmt.Sprintf("%s — Binder Amount as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(binder, s))},
		{Key: fmt.Sprintf("total_%s_outsourced_amount", b.Code), Label: fmt.Sprintf("%s — Total Outsourced Amount", b.Title), Value: money(outsourced)},
		{Key: fmt.Sprintf("%s_outsourced_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Outsource Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(outsourced, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_outsourced_amount_salary", b.Code), Label: fmt.Sprintf("%s — Outsource Amount as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(outsourced, s))},
		{Key: fmt.Sprintf("total_%s_commission_amount", b.Code), Label: fmt.Sprintf("%s — Total Commission Amount", b.Title), Value: money(commission)},
		{Key: fmt.Sprintf("%s_commission_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Commission Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(commission, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_commission_amount_salary", b.Code), Label: fmt.Sprintf("%s — Commission Amount as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(commission, s))},
		{Key: fmt.Sprintf("total_%s_basic_premium", b.Code), Label: fmt.Sprintf("%s — Total Basic Premium", b.Title), Value: money(basic)},
		{Key: fmt.Sprintf("%s_basic_rate_per1000_sa", b.Code), Label: fmt.Sprintf("%s — Basic Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(basic, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_basic_premium_salary", b.Code), Label: fmt.Sprintf("%s — Basic Premium as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(basic, s))},
	}
}

// conversionSliceBlock returns the 6-field block for one (benefit × slice)
// pairing. sliceKey is the snake_case suffix appended to the benefit code
// in each token key (e.g. "conv_on_wdr"); sliceLabel is the human-readable
// descriptor used in labels (e.g. "Conversion on Withdrawal").
func conversionSliceBlock(b benefitName, sliceKey, sliceLabel string,
	riskPrem, officePrem, propRiskSalary, propOfficeSalary,
	riskRatePer1000, officeRatePer1000 float64,
) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	return []Field{
		{Key: fmt.Sprintf("total_%s_%s_risk_prem", b.Code, sliceKey), Label: fmt.Sprintf("%s — %s (Risk Premium)", b.Title, sliceLabel), Value: money(riskPrem)},
		{Key: fmt.Sprintf("total_%s_%s_office_prem", b.Code, sliceKey), Label: fmt.Sprintf("%s — %s (Office Premium)", b.Title, sliceLabel), Value: money(officePrem)},
		{Key: fmt.Sprintf("prop_%s_%s_risk_prem_salary", b.Code, sliceKey), Label: fmt.Sprintf("%s — %s Risk Premium as %% of Salary", b.Title, sliceLabel), Value: formatPercent(propRiskSalary)},
		{Key: fmt.Sprintf("prop_%s_%s_office_prem_salary", b.Code, sliceKey), Label: fmt.Sprintf("%s — %s Office Premium as %% of Salary", b.Title, sliceLabel), Value: formatPercent(propOfficeSalary)},
		{Key: fmt.Sprintf("%s_%s_risk_rate_per_1000_sa", b.Code, sliceKey), Label: fmt.Sprintf("%s — %s Risk Rate per 1,000 SA", b.Title, sliceLabel), Value: money(riskRatePer1000)},
		{Key: fmt.Sprintf("%s_%s_office_rate_per_1000_sa", b.Code, sliceKey), Label: fmt.Sprintf("%s — %s Office Rate per 1,000 SA", b.Title, sliceLabel), Value: money(officeRatePer1000)},
	}
}

// categoryBoolFields returns the has_* flags. Rendered in the sample as
// bullet points demonstrating conditional-block syntax. Per-benefit flag
// keys use the customised code where set (e.g. a customised GLA with
// code "gl" becomes {{#has_gl}}); labels use the customised title.
func categoryBoolFields(
	s models.MemberRatingResultSummary,
	flags benefitFlags,
	n benefitNaming,
) []Field {
	has := func(b benefitName, v bool) Field {
		return Field{
			Key:   fmt.Sprintf("has_%s", b.Code),
			Label: fmt.Sprintf("Category has %s", b.Title),
			Value: v,
		}
	}
	return []Field{
		{Key: "has_non_funeral_benefits", Label: "Category has any non-funeral benefit", Value: quote_docx.CategoryHasNonFuneralBenefits(s)},
		has(n.GLA, flags.GLA),
		has(n.SGLA, flags.SGLA),
		has(n.PTD, flags.PTD),
		has(n.CI, flags.CI),
		has(n.PHI, flags.PHI),
		has(n.TTD, flags.TTD),
		has(n.Funeral, flags.Funeral),
		has(n.AdditionalGla, flags.AdditionalGla),
		has(n.AdditionalAccidentalGla, flags.AdditionalAccidentalGla),
		has(n.GlaEducator, flags.GlaEducator),
		has(n.PtdEducator, flags.PtdEducator),
		has(n.ExtendedFamily, flags.ExtendedFamily),
		has(n.TaxSaver, flags.TaxSaver),
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.FinalGlaAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(officeProportionFromFinalPremium(s.FinalGlaAnnualOfficePremium, &s))},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.FinalSglaAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(officeProportionFromFinalPremium(s.FinalSglaAnnualOfficePremium, &s))},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.FinalPtdAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(officeProportionFromFinalPremium(s.FinalPtdAnnualOfficePremium, &s))},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.FinalCiAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(officeProportionFromFinalPremium(s.FinalCiAnnualOfficePremium, &s))},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.FinalPhiAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(officeProportionFromFinalPremium(s.FinalPhiAnnualOfficePremium, &s))},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.FinalTtdAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(officeProportionFromFinalPremium(s.FinalTtdAnnualOfficePremium, &s))},
	}
}

func funeralFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.FamilyFuneralBenefitTitle},
		{Key: "monthly_premium_per_member", Label: "Monthly Premium per Member", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunMonthlyPremiumPerMember)},
		{Key: "premium_per_member", Label: "Premium per Member", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunAnnualPremiumPerMember)},
		{Key: "total_premium", Label: "Total Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.FinalFunAnnualOfficePremium)},
		{Key: "main_member_sum_assured", Label: "Main Member Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralMainMemberFuneralSumAssured)},
		{Key: "spouse_sum_assured", Label: "Spouse Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralSpouseFuneralSumAssured)},
		{Key: "child_sum_assured", Label: "Child Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralChildrenFuneralSumAssured)},
		{Key: "max_children", Label: "Max Children Covered", Value: strconv.Itoa(cat.FamilyFuneralMaxNumberChildren)},
		{Key: "parent_sum_assured", Label: "Parent Sum Assured", Value: quote_docx.RoundUpToTwoDecimalsAccounting(cat.FamilyFuneralParentFuneralSumAssured)},
		{Key: "max_parents", Label: "Max Parents Covered", Value: strconv.Itoa(cat.FamilyFuneralMaxNumberParents)},
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

// ---------------------------------------------------------------------------
// Exported sample accessors — other packages (e.g. on_risk_letter_template)
// use these to build their own self-documenting samples from the shared
// schema without having to know about benefitNaming or invoke the private
// *Fields helpers directly. Each accessor resolves the insurer's
// customised benefit codes via sampleBenefitNaming() so the tokens shown
// in any derived sample match the ones the render engine will emit.
// ---------------------------------------------------------------------------

// QuoteFieldsForSample returns the root-scope quote tokens with zero-value
// fixtures.
func QuoteFieldsForSample() []Field {
	var (
		zq models.GroupPricingQuote
		zT quote_docx.QuoteTotals
	)
	return quoteFields(zq, zT, false)
}

// InsurerFieldsForSample returns the {{insurer.*}} tokens with zero-value
// fixtures.
func InsurerFieldsForSample() []Field {
	var zi models.GroupPricingInsurerDetail
	return insurerFields(zi)
}

// CategoryScalarFieldsForSample returns the non-bool category tokens using
// the resolved (DB-customised or default) benefit naming.
func CategoryScalarFieldsForSample() []Field {
	var (
		zs models.MemberRatingResultSummary
		zc models.SchemeCategory
		zq models.GroupPricingQuote
	)
	return categoryScalarFields(zs, zc, zq, sampleBenefitNaming())
}

// CategoryBoolFieldsForSample returns the has_* category flags using the
// resolved benefit naming.
func CategoryBoolFieldsForSample() []Field {
	var zs models.MemberRatingResultSummary
	return categoryBoolFields(zs, benefitFlags{}, sampleBenefitNaming())
}

// BenefitSpecsForSample returns one spec per nested benefit object, each
// prefixed by the resolved benefit code.
func BenefitSpecsForSample() []BenefitSpec {
	return benefitSpecsForSample(sampleBenefitNaming())
}

// benefitSpecsForSample returns one spec per benefit, ordered to match
// the legacy sample layout. Each Fields closure passes zero-value inputs
// because the sample only needs Keys/Labels. The Prefix on each spec
// uses the resolved benefit code (customised where set) so the sample
// document shows the same nested-scope prefixes the render engine emits.
func benefitSpecsForSample(n benefitNaming) []BenefitSpec {
	var (
		zs models.MemberRatingResultSummary
		zc models.SchemeCategory
		zq models.GroupPricingQuote
		zt quote_docx.BenefitTitles
	)
	return []BenefitSpec{
		{Prefix: n.GLA.Code, Title: n.GLA.Title, Fields: func() []Field { return glaFields(zs, zc, zq, zt) }},
		{Prefix: n.SGLA.Code, Title: n.SGLA.Title, Fields: func() []Field { return sglaFields(zs, zc, zq, zt) }},
		{Prefix: n.PTD.Code, Title: n.PTD.Title, Fields: func() []Field { return ptdFields(zs, zc, zq, zt) }},
		{Prefix: n.CI.Code, Title: n.CI.Title, Fields: func() []Field { return ciFields(zs, zc, zq, zt) }},
		{Prefix: n.PHI.Code, Title: n.PHI.Title, Fields: func() []Field { return phiFields(zs, zc, zt) }},
		{Prefix: n.TTD.Code, Title: n.TTD.Title, Fields: func() []Field { return ttdFields(zs, zc, zt) }},
		{Prefix: n.Funeral.Code, Title: n.Funeral.Title, Fields: func() []Field { return funeralFields(zs, zc, zt) }},
	}
}

// officeProportionFromFinalPremium computes office prop-of-salary directly
// from a persisted Final office premium. Mirrors the rating-time risk
// proportion formula but with the (post-discount, post-commission) Final
// office premium as numerator.
//
// Fallback: when IndicativeRatesCount is zero (legacy summaries pre-dating
// the column, or rows whose recalc-time INSERT failed because the migration
// hadn't been applied yet) the count is back-derived from any persisted
// (risk premium, risk proportion) pair: count = riskPrem / (salary * riskProp).
// GLA is preferred; PTD / CI / SGLA / PHI / TTD / FUN serve as fallbacks for
// schemes without a GLA benefit.
func officeProportionFromFinalPremium(finalOfficePrem float64, s *models.MemberRatingResultSummary) float64 {
	if s == nil || s.TotalAnnualSalary <= 0 {
		return 0
	}
	count := s.IndicativeRatesCount
	if count <= 0 {
		count = inferIndicativeRatesCount(s)
	}
	denom := s.TotalAnnualSalary * count
	if denom <= 0 {
		return 0
	}
	return finalOfficePrem / denom
}

// inferIndicativeRatesCount back-derives count from any persisted
// (risk premium, risk proportion) pair on the summary. Returns 0 only when
// every benefit's risk proportion is zero, in which case the prop-of-salary
// helper falls through to its existing zero-denom guard.
func inferIndicativeRatesCount(s *models.MemberRatingResultSummary) float64 {
	type pair struct {
		riskPrem, riskProp float64
	}
	pairs := []pair{
		{s.TotalGlaAnnualRiskPremium, s.ProportionGlaAnnualRiskPremiumSalary},
		{s.TotalPtdAnnualRiskPremium, s.ProportionPtdAnnualRiskPremiumSalary},
		{s.TotalCiAnnualRiskPremium, s.ProportionCiAnnualRiskPremiumSalary},
		{s.TotalSglaAnnualRiskPremium, s.ProportionSglaAnnualRiskPremiumSalary},
		{s.TotalPhiAnnualRiskPremium, s.ProportionPhiAnnualRiskPremiumSalary},
		{s.TotalTtdAnnualRiskPremium, s.ProportionTtdAnnualRiskPremiumSalary},
		{s.TotalFunAnnualRiskPremium, s.ProportionFunAnnualRiskPremiumSalary},
	}
	for _, p := range pairs {
		if p.riskPrem > 0 && p.riskProp > 0 {
			return p.riskPrem / (s.TotalAnnualSalary * p.riskProp)
		}
	}
	return 0
}

// officeRateFromFinalPremium computes office rate-per-1000-SA directly from
// a persisted Final office premium and the benefit's sum-assured (or capped
// income, for income-based benefits). Matches the denominator the risk-rate
// calc uses at rating time.
func officeRateFromFinalPremium(finalOfficePrem, sumAssured float64) float64 {
	if sumAssured <= 0 {
		return 0
	}
	return finalOfficePrem * 1000.0 / sumAssured
}

// basicPremium returns officePrem minus every subtracted component, clamped
// at zero so floating-point drift / mis-aligned data can't surface negatives
// in the doc-template. Components capture commission / binder / outsource /
// tax-saver / conversion-continuity slice allocations that the user wants
// stripped from a benefit's "basic" view.
func basicPremium(officePrem float64, components ...float64) float64 {
	out := officePrem
	for _, c := range components {
		out -= c
	}
	if out < 0 {
		return 0
	}
	return out
}

// benefitBasicPremiumBlock emits the three basic-premium tokens for one
// benefit: amount, rate-per-1000-SA, and proportion-of-salary. Reuses the
// office-premium rate/prop helpers so behaviour matches the office block
// exactly. useShort applies the abbreviated benefit code for Add Acc GLA.
func benefitBasicPremiumBlock(b benefitName, basic, sumAssured float64, s *models.MemberRatingResultSummary, useShort bool) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	code := b.Code
	if useShort {
		code = b.Short()
	}
	return []Field{
		{Key: fmt.Sprintf("total_%s_basic_premium", code), Label: fmt.Sprintf("%s — Total Basic Premium", b.Title), Value: money(basic)},
		{Key: fmt.Sprintf("%s_basic_rate_per1000_sa", code), Label: fmt.Sprintf("%s — Basic Rate per 1,000 SA", b.Title), Value: money(officeRateFromFinalPremium(basic, sumAssured))},
		{Key: fmt.Sprintf("proportion_%s_basic_premium_salary", code), Label: fmt.Sprintf("%s — Basic Premium as %% of Salary", b.Title), Value: formatPercent(officeProportionFromFinalPremium(basic, s))},
	}
}

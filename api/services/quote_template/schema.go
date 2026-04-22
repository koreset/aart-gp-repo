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
		Funeral:                 s.TotalFunAnnualOfficePremium > 0 || cat.FamilyFuneralBenefit,
		AdditionalGla:           s.AdditionalGlaCoverBenefit,
		AdditionalAccidentalGla: s.TotalAdditionalAccidentalGlaCappedSumAssured > 0,
		GlaEducator:             s.TotalGlaEducatorOfficePremium > 0 || cat.GlaEducatorBenefit != "",
		PtdEducator:             s.TotalPtdEducatorOfficePremium > 0 || cat.PtdEducatorBenefit != "",
		ExtendedFamily:          s.ExtendedFamilyBenefit,
		TaxSaver:                s.TaxSaverBenefit || s.ExpTotalTaxSaverAnnualOfficePremium > 0,
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
		{Key: "total_salary", Label: "Total Salary", Value: quote_docx.RoundUpToTwoDecimalsAccounting(totals.TotalAnnualSalary)},
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

// categoryScalarFields returns the non-bool category tokens. Rendered in
// the sample as a key/value table. Benefit-prefixed tokens (rate_per_1000,
// educator indicators, and all the rating/slice tokens built by the
// sub-builders below) use the customised code/title from n where the
// insurer has configured one; otherwise the canonical defaults apply.
func categoryScalarFields(
	s models.MemberRatingResultSummary,
	cat models.SchemeCategory,
	n benefitNaming,
) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	fs := []Field{
		{Key: "name", Label: "Category Name", Value: s.Category},
		{Key: "region", Label: "Region", Value: cat.Region},
		{Key: "member_count", Label: "Member Count", Value: strconv.Itoa(int(s.MemberCount))},
		{Key: "total_salary", Label: "Total Salary", Value: money(s.TotalAnnualSalary)},
		{Key: "total_sum_assured", Label: "Total Sum Assured", Value: money(s.TotalSumAssured)},
		{Key: "premium", Label: "Premium (excl. funeral)", Value: money(s.ExpTotalAnnualPremiumExclFuneral)},
		{Key: "percent_salary", Label: "Premium as % of Salary", Value: formatPercent(s.ProportionExpTotalPremiumExclFuneralSalary)},
		{Key: "free_cover_limit", Label: "Free Cover Limit (category override)", Value: money(cat.FreeCoverLimit)},
		{Key: fmt.Sprintf("%s_rate_per_1000", n.GLA.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", n.GLA.Title), Value: money(s.ExpGlaOfficeRatePer1000SA)},
		{Key: fmt.Sprintf("%s_rate_per_1000", n.SGLA.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", n.SGLA.Title), Value: money(s.ExpSglaOfficeRatePer1000SA)},
		{Key: fmt.Sprintf("%s_rate_per_1000", n.PTD.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", n.PTD.Title), Value: money(s.ExpPtdOfficeRatePer1000SA)},
		{Key: fmt.Sprintf("%s_rate_per_1000", n.CI.Code), Label: fmt.Sprintf("%s Rate per 1,000 SA", n.CI.Title), Value: money(s.ExpCiOfficeRatePer1000SA)},
		{Key: "retirement_premium_waiver", Label: "Retirement Premium Waiver", Value: orDash(cat.PhiPremiumWaiver)},
		{Key: "medical_aid_premium_waiver", Label: "Medical Aid Premium Waiver", Value: orDash(cat.PhiMedicalAidPremiumWaiver)},
		{Key: fmt.Sprintf("%s_terminal_illness_benefit", n.GLA.Code), Label: fmt.Sprintf("%s Terminal Illness Benefit", n.GLA.Title), Value: orDash(cat.GlaTerminalIllnessBenefit)},
		{Key: fmt.Sprintf("%s_benefit", n.GlaEducator.Code), Label: fmt.Sprintf("%s Benefit", n.GlaEducator.Title), Value: orDash(cat.GlaEducatorBenefit)},
		{Key: fmt.Sprintf("%s_benefit", n.PtdEducator.Code), Label: fmt.Sprintf("%s Benefit", n.PtdEducator.Title), Value: orDash(cat.PtdEducatorBenefit)},
	}
	fs = append(fs, categoryRatingSummaryFields(s, n)...)
	fs = append(fs, categoryEducatorSummaryFields(s, n)...)
	fs = append(fs, categoryConversionSliceFields(s, n)...)
	return fs
}

// categoryRatingSummaryFields exposes the full member-rating result summary
// tokens at category scope: per-benefit risk rates, risk premiums, office
// premiums, proportions of salary, rate-per-1000 figures, binder/outsource
// splits, per-benefit commission, scheme-level commission totals, and tax
// saver slices. Keys/labels use the customised benefit code + title from n
// where set, falling back to defaults. The leading "exp_" prefix from the
// underlying model fields is stripped from all token keys.
func categoryRatingSummaryFields(s models.MemberRatingResultSummary, n benefitNaming) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	var fs []Field

	// Core benefit rating blocks (8 fields each).
	fs = append(fs, benefitRatingBlock(n.GLA, false,
		s.TotalGlaCappedSumAssured, s.ExpTotalGlaRiskRate, s.ExpTotalGlaAnnualRiskPremium,
		s.ExpGlaRiskRatePer1000SA, s.ExpProportionGlaAnnualRiskPremiumSalary,
		s.ExpTotalGlaAnnualOfficePremium, s.ExpGlaOfficeRatePer1000SA, s.ExpProportionGlaOfficePremiumSalary)...)
	fs = append(fs, benefitRatingBlock(n.PTD, false,
		s.TotalPtdCappedSumAssured, s.ExpTotalPtdRiskRate, s.ExpTotalPtdAnnualRiskPremium,
		s.ExpPtdRiskRatePer1000SA, s.ExpProportionPtdAnnualRiskPremiumSalary,
		s.ExpTotalPtdAnnualOfficePremium, s.ExpPtdOfficeRatePer1000SA, s.ExpProportionPtdOfficePremiumSalary)...)
	fs = append(fs, benefitRatingBlock(n.CI, false,
		s.TotalCiCappedSumAssured, s.ExpTotalCiRiskRate, s.ExpTotalCiAnnualRiskPremium,
		s.ExpCiRiskRatePer1000SA, s.ExpProportionCiAnnualRiskPremiumSalary,
		s.ExpTotalCiAnnualOfficePremium, s.ExpCiOfficeRatePer1000SA, s.ExpProportionCiOfficePremiumSalary)...)
	fs = append(fs, benefitRatingBlock(n.SGLA, false,
		s.TotalSglaCappedSumAssured, s.ExpTotalSglaRiskRate, s.ExpTotalSglaAnnualRiskPremium,
		s.ExpSglaRiskRatePer1000SA, s.ExpProportionSglaAnnualRiskPremiumSalary,
		s.ExpTotalSglaAnnualOfficePremium, s.ExpSglaOfficeRatePer1000SA, s.ExpProportionSglaOfficePremiumSalary)...)
	fs = append(fs, benefitRatingBlock(n.TTD, true,
		s.TotalTtdCappedIncome, s.ExpTotalTtdRiskRate, s.ExpTotalTtdAnnualRiskPremium,
		s.ExpTtdRiskRatePer1000SA, s.ExpProportionTtdAnnualRiskPremiumSalary,
		s.ExpTotalTtdAnnualOfficePremium, s.ExpTtdOfficeRatePer1000SA, s.ExpProportionTtdOfficePremiumSalary)...)
	fs = append(fs, benefitRatingBlock(n.PHI, true,
		s.TotalPhiCappedIncome, s.ExpTotalPhiRiskRate, s.ExpTotalPhiAnnualRiskPremium,
		s.ExpPhiRiskRatePer1000SA, s.ExpProportionPhiAnnualRiskPremiumSalary,
		s.ExpTotalPhiAnnualOfficePremium, s.ExpPhiOfficeRatePer1000SA, s.ExpProportionPhiOfficePremiumSalary)...)

	// Funeral + aggregate. "monthly" stays on the monthly-premium-per-member
	// key to disambiguate it from the (default) annual sibling, but the
	// "annual" qualifier is dropped everywhere else.
	fs = append(fs,
		Field{Key: fmt.Sprintf("total_%s_risk_premium", n.Funeral.Code), Label: fmt.Sprintf("%s — Total Risk Premium", n.Funeral.Title), Value: money(s.ExpTotalFunAnnualRiskPremium)},
		Field{Key: fmt.Sprintf("proportion_%s_risk_premium_salary", n.Funeral.Code), Label: fmt.Sprintf("%s — Risk Premium as %% of Salary", n.Funeral.Title), Value: formatPercent(s.ExpProportionFunAnnualRiskPremiumSalary)},
		Field{Key: fmt.Sprintf("total_%s_office_premium", n.Funeral.Code), Label: fmt.Sprintf("%s — Total Office Premium", n.Funeral.Title), Value: money(s.ExpTotalFunAnnualOfficePremium)},
		Field{Key: fmt.Sprintf("proportion_%s_office_premium_salary", n.Funeral.Code), Label: fmt.Sprintf("%s — Office Premium as %% of Salary", n.Funeral.Title), Value: formatPercent(s.ExpProportionFunOfficePremiumSalary)},
		Field{Key: fmt.Sprintf("total_%s_premium_per_member", n.Funeral.Code), Label: fmt.Sprintf("%s — Premium per Member", n.Funeral.Title), Value: money(s.ExpTotalFunAnnualPremiumPerMember)},
		Field{Key: fmt.Sprintf("total_%s_monthly_premium_per_member", n.Funeral.Code), Label: fmt.Sprintf("%s — Monthly Premium per Member", n.Funeral.Title), Value: money(s.ExpTotalFunMonthlyPremiumPerMember)},
		Field{Key: "total_premium_excl_funeral", Label: "Total Premium (excluding Funeral)", Value: money(s.ExpTotalAnnualPremiumExclFuneral)},
		Field{Key: "proportion_total_premium_excl_funeral_salary", Label: "Total Premium (excluding Funeral) as % of Salary", Value: formatPercent(s.ProportionExpTotalPremiumExclFuneralSalary)},
	)

	// Additional Accidental GLA — mirrors the core rating shape, but the
	// risk-proportion key uses "prop_" (mirroring the attached list /
	// underlying column name); the office-proportion key uses "proportion_".
	ab := n.AdditionalAccidentalGla
	fs = append(fs,
		Field{Key: fmt.Sprintf("total_%s_capped_sum_assured", ab.Code), Label: fmt.Sprintf("%s — Total Capped Sum Assured", ab.Title), Value: money(s.TotalAdditionalAccidentalGlaCappedSumAssured)},
		Field{Key: fmt.Sprintf("total_%s_risk_rate", ab.Code), Label: fmt.Sprintf("%s — Total Risk Rate", ab.Title), Value: money(s.ExpTotalAdditionalAccidentalGlaRiskRate)},
		Field{Key: fmt.Sprintf("total_%s_risk_premium", ab.Code), Label: fmt.Sprintf("%s — Total Risk Premium", ab.Title), Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium)},
		Field{Key: fmt.Sprintf("%s_risk_rate_per1000_sa", ab.Code), Label: fmt.Sprintf("%s — Risk Rate per 1,000 SA", ab.Title), Value: money(s.ExpAdditionalAccidentalGlaRiskRatePer1000SA)},
		Field{Key: fmt.Sprintf("prop_%s_risk_premium_salary", ab.Code), Label: fmt.Sprintf("%s — Risk Premium as %% of Salary", ab.Title), Value: formatPercent(s.ExpProportionAdditionalAccidentalGlaAnnualRiskPremiumSalary)},
		Field{Key: fmt.Sprintf("total_%s_office_premium", ab.Code), Label: fmt.Sprintf("%s — Total Office Premium", ab.Title), Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualOfficePremium)},
		Field{Key: fmt.Sprintf("%s_office_rate_per1000_sa", ab.Code), Label: fmt.Sprintf("%s — Office Rate per 1,000 SA", ab.Title), Value: money(s.ExpAdditionalAccidentalGlaOfficeRatePer1000SA)},
		Field{Key: fmt.Sprintf("proportion_%s_office_premium_salary", ab.Code), Label: fmt.Sprintf("%s — Office Premium as %% of Salary", ab.Title), Value: formatPercent(s.ExpProportionAdditionalAccidentalGlaOfficePremiumSalary)},
	)

	// Binder & outsource amounts. Add Acc GLA uses ShortCode for these keys.
	fs = append(fs, benefitBinderOutsourceBlock(n.GLA, s.ExpTotalGlaAnnualBinderAmount, s.ExpTotalGlaAnnualOutsourcedAmount, false)...)
	fs = append(fs, benefitBinderOutsourceBlock(n.AdditionalAccidentalGla, s.ExpTotalAdditionalAccidentalGlaAnnualBinderAmount, s.ExpTotalAdditionalAccidentalGlaAnnualOutsourcedAmt, true)...)
	fs = append(fs, benefitBinderOutsourceBlock(n.PTD, s.ExpTotalPtdAnnualBinderAmount, s.ExpTotalPtdAnnualOutsourcedAmount, false)...)
	fs = append(fs, benefitBinderOutsourceBlock(n.CI, s.ExpTotalCiAnnualBinderAmount, s.ExpTotalCiAnnualOutsourcedAmount, false)...)
	fs = append(fs, benefitBinderOutsourceBlock(n.SGLA, s.ExpTotalSglaAnnualBinderAmount, s.ExpTotalSglaAnnualOutsourcedAmount, false)...)
	fs = append(fs, benefitBinderOutsourceBlock(n.TTD, s.ExpTotalTtdAnnualBinderAmount, s.ExpTotalTtdAnnualOutsourcedAmount, false)...)
	fs = append(fs, benefitBinderOutsourceBlock(n.PHI, s.ExpTotalPhiAnnualBinderAmount, s.ExpTotalPhiAnnualOutsourcedAmount, false)...)
	fs = append(fs, benefitBinderOutsourceBlock(n.Funeral, s.ExpTotalFunAnnualBinderAmount, s.ExpTotalFunAnnualOutsourcedAmount, false)...)

	// Commission amounts (per benefit + scheme totals).
	fs = append(fs,
		benefitCommissionField(n.GLA, s.ExpTotalGlaAnnualCommissionAmount, false),
		benefitCommissionField(n.AdditionalAccidentalGla, s.ExpTotalAdditionalAccidentalGlaAnnualCommissionAmount, true),
		benefitCommissionField(n.PTD, s.ExpTotalPtdAnnualCommissionAmount, false),
		benefitCommissionField(n.CI, s.ExpTotalCiAnnualCommissionAmount, false),
		benefitCommissionField(n.SGLA, s.ExpTotalSglaAnnualCommissionAmount, false),
		benefitCommissionField(n.TTD, s.ExpTotalTtdAnnualCommissionAmount, false),
		benefitCommissionField(n.PHI, s.ExpTotalPhiAnnualCommissionAmount, false),
		benefitCommissionField(n.Funeral, s.ExpTotalFunAnnualCommissionAmount, false),
		Field{Key: "scheme_total_commission", Label: "Scheme Total Commission", Value: money(s.SchemeTotalCommission)},
		Field{Key: "scheme_total_commission_rate", Label: "Scheme Total Commission Rate", Value: formatPercent(s.SchemeTotalCommissionRate)},
	)

	// Tax saver (slice of GLA office premium — no SA / rate-per-1000).
	tb := n.TaxSaver
	fs = append(fs,
		Field{Key: fmt.Sprintf("total_%s_risk_premium", tb.Code), Label: fmt.Sprintf("%s — Total Risk Premium", tb.Title), Value: money(s.ExpTotalTaxSaverAnnualRiskPremium)},
		Field{Key: fmt.Sprintf("total_%s_office_premium", tb.Code), Label: fmt.Sprintf("%s — Total Office Premium", tb.Title), Value: money(s.ExpTotalTaxSaverAnnualOfficePremium)},
	)

	return fs
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

// benefitBinderOutsourceBlock returns the binder + outsourced token pair
// for one benefit. When useShort is true, b.Short() is used for the key
// prefix (Add Acc GLA uses the abbreviated form to keep keys tractable).
// The outsourced token's suffix is "_amt" for Add Acc GLA (matches the
// historical attached-list naming) and "_amount" elsewhere.
func benefitBinderOutsourceBlock(b benefitName, binder, outsourced float64, useShort bool) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	code := b.Code
	outsourcedSuffix := "outsourced_amount"
	if useShort {
		code = b.Short()
		// The "_amt" suffix is a quirk of the default Add Acc GLA naming
		// (ShortCode differs from Code). When a customisation is applied
		// the resolver sets ShortCode == Code, so the consistent
		// "_amount" suffix is used.
		if b.ShortCode != "" && b.ShortCode != b.Code {
			outsourcedSuffix = "outsourced_amt"
		}
	}
	return []Field{
		{Key: fmt.Sprintf("total_%s_binder_amount", code), Label: fmt.Sprintf("%s — Total Binder Amount", b.Title), Value: money(binder)},
		{Key: fmt.Sprintf("total_%s_%s", code, outsourcedSuffix), Label: fmt.Sprintf("%s — Total Outsourced Amount", b.Title), Value: money(outsourced)},
	}
}

// benefitCommissionField returns the commission-amount token for one
// benefit. useShort chooses the abbreviated code (Add Acc GLA).
func benefitCommissionField(b benefitName, amount float64, useShort bool) Field {
	code := b.Code
	if useShort {
		code = b.Short()
	}
	return Field{
		Key:   fmt.Sprintf("total_%s_commission_amount", code),
		Label: fmt.Sprintf("%s — Total Commission Amount", b.Title),
		Value: quote_docx.RoundUpToTwoDecimalsAccounting(amount),
	}
}

// categoryEducatorSummaryFields exposes the GLA/PTD educator split tokens:
// risk and office premiums, proportion-of-salary, rate-per-1000, plus
// binder, outsource, and commission breakdowns for each educator cover.
// Token keys use the customised educator code (from n.GlaEducator /
// n.PtdEducator) where set, falling back to the defaults "gla_educator"
// and "ptd_educator".
func categoryEducatorSummaryFields(s models.MemberRatingResultSummary, n benefitNaming) []Field {
	var fs []Field
	fs = append(fs, educatorSplitBlock(n.GlaEducator,
		s.ExpAdjTotalGlaEducatorRiskPremium, s.ExpAdjTotalGlaEducatorOfficePremium,
		s.ExpAdjProportionGlaEducatorRiskPremiumSalary, s.ExpAdjProportionGlaEducatorOfficePremiumSalary,
		s.ExpGlaEducatorRiskRatePer1000SA, s.ExpGlaEducatorOfficeRatePer1000SA,
		s.ExpAdjTotalGlaEducatorBinderAmount, s.ExpAdjTotalGlaEducatorOutsourcedAmount,
		s.ExpAdjTotalGlaEducatorCommissionAmount)...)
	fs = append(fs, educatorSplitBlock(n.PtdEducator,
		s.ExpAdjTotalPtdEducatorRiskPremium, s.ExpAdjTotalPtdEducatorOfficePremium,
		s.ExpAdjProportionPtdEducatorRiskPremiumSalary, s.ExpAdjProportionPtdEducatorOfficePremiumSalary,
		s.ExpPtdEducatorRiskRatePer1000SA, s.ExpPtdEducatorOfficeRatePer1000SA,
		s.ExpAdjTotalPtdEducatorBinderAmount, s.ExpAdjTotalPtdEducatorOutsourcedAmount,
		s.ExpAdjTotalPtdEducatorCommissionAmount)...)
	return fs
}

// educatorSplitBlock returns the 9-field educator block for one educator
// cover (GLA or PTD). Values come from the experience-adjusted model
// fields; the "adjusted" qualifier is omitted from keys and labels to
// keep tokens short and time-neutral.
func educatorSplitBlock(b benefitName,
	riskPrem, officePrem, propRiskSalary, propOfficeSalary,
	riskRatePer1000, officeRatePer1000,
	binder, outsourced, commission float64,
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
		{Key: fmt.Sprintf("total_%s_outsourced_amount", b.Code), Label: fmt.Sprintf("%s — Total Outsourced Amount", b.Title), Value: money(outsourced)},
		{Key: fmt.Sprintf("total_%s_commission_amount", b.Code), Label: fmt.Sprintf("%s — Total Commission Amount", b.Title), Value: money(commission)},
	}
}

// categoryConversionSliceFields exposes conversion / continuity slice
// tokens. Each slice carries six variants: risk premium, office premium,
// proportion of salary (risk + office) and rate-per-1000 (risk + office).
// Keys use the underlying benefit code (possibly customised) joined with
// a slice suffix like "conv_on_wdr".
func categoryConversionSliceFields(s models.MemberRatingResultSummary, n benefitNaming) []Field {
	var fs []Field

	// GLA slices
	fs = append(fs, conversionSliceBlock(n.GLA, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalGlaConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionGlaConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionGlaConversionOnWithdrawalOfficePremiumSalary,
		s.ExpGlaConversionOnWithdrawalRiskRatePer1000SA, s.ExpGlaConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.GLA, "conv_on_ret", "Conversion on Retirement",
		s.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium, s.ExpAdjTotalGlaConversionOnRetirementAnnualOfficePremium,
		s.ExpAdjProportionGlaConversionOnRetirementRiskPremiumSalary, s.ExpAdjProportionGlaConversionOnRetirementOfficePremiumSalary,
		s.ExpGlaConversionOnRetirementRiskRatePer1000SA, s.ExpGlaConversionOnRetirementOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.GLA, "cont_dur_dis", "Continuity During Disability",
		s.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium, s.ExpAdjTotalGlaContinuityDuringDisabilityAnnualOfficePremium,
		s.ExpAdjProportionGlaContinuityDuringDisabilityRiskPremiumSalary, s.ExpAdjProportionGlaContinuityDuringDisabilityOfficePremiumSalary,
		s.ExpGlaContinuityDuringDisabilityRiskRatePer1000SA, s.ExpGlaContinuityDuringDisabilityOfficeRatePer1000SA)...)

	// GLA educator slices
	fs = append(fs, conversionSliceBlock(n.GlaEducator, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionGlaEducatorConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionGlaEducatorConversionOnWithdrawalOfficePremiumSalary,
		s.ExpGlaEducatorConversionOnWithdrawalRiskRatePer1000SA, s.ExpGlaEducatorConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.GlaEducator, "conv_on_ret", "Conversion on Retirement",
		s.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium, s.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualOfficePremium,
		s.ExpAdjProportionGlaEducatorConversionOnRetirementRiskPremiumSalary, s.ExpAdjProportionGlaEducatorConversionOnRetirementOfficePremiumSalary,
		s.ExpGlaEducatorConversionOnRetirementRiskRatePer1000SA, s.ExpGlaEducatorConversionOnRetirementOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.GlaEducator, "cont_dur_dis", "Continuity During Disability",
		s.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium, s.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualOfficePremium,
		s.ExpAdjProportionGlaEducatorContinuityDuringDisabilityRiskPremiumSalary, s.ExpAdjProportionGlaEducatorContinuityDuringDisabilityOfficePremiumSalary,
		s.ExpGlaEducatorContinuityDuringDisabilityRiskRatePer1000SA, s.ExpGlaEducatorContinuityDuringDisabilityOfficeRatePer1000SA)...)

	// PTD slices
	fs = append(fs, conversionSliceBlock(n.PTD, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalPtdConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionPtdConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionPtdConversionOnWithdrawalOfficePremiumSalary,
		s.ExpPtdConversionOnWithdrawalRiskRatePer1000SA, s.ExpPtdConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.PtdEducator, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionPtdEducatorConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionPtdEducatorConversionOnWithdrawalOfficePremiumSalary,
		s.ExpPtdEducatorConversionOnWithdrawalRiskRatePer1000SA, s.ExpPtdEducatorConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.PtdEducator, "conv_on_ret", "Conversion on Retirement",
		s.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium, s.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualOfficePremium,
		s.ExpAdjProportionPtdEducatorConversionOnRetirementRiskPremiumSalary, s.ExpAdjProportionPtdEducatorConversionOnRetirementOfficePremiumSalary,
		s.ExpPtdEducatorConversionOnRetirementRiskRatePer1000SA, s.ExpPtdEducatorConversionOnRetirementOfficeRatePer1000SA)...)

	// PHI / CI / SGLA / Funeral conversion on withdrawal
	fs = append(fs, conversionSliceBlock(n.PHI, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalPhiConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionPhiConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionPhiConversionOnWithdrawalOfficePremiumSalary,
		s.ExpPhiConversionOnWithdrawalRiskRatePer1000SA, s.ExpPhiConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.CI, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalCiConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionCiConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionCiConversionOnWithdrawalOfficePremiumSalary,
		s.ExpCiConversionOnWithdrawalRiskRatePer1000SA, s.ExpCiConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.SGLA, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalSglaConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionSglaConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionSglaConversionOnWithdrawalOfficePremiumSalary,
		s.ExpSglaConversionOnWithdrawalRiskRatePer1000SA, s.ExpSglaConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.Funeral, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalFunConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionFunConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionFunConversionOnWithdrawalOfficePremiumSalary,
		s.ExpFunConversionOnWithdrawalRiskRatePer1000SA, s.ExpFunConversionOnWithdrawalOfficeRatePer1000SA)...)
	fs = append(fs, conversionSliceBlock(n.TTD, "conv_on_wdr", "Conversion on Withdrawal",
		s.ExpAdjTotalTtdConversionOnWithdrawalAnnualRiskPremium, s.ExpAdjTotalTtdConversionOnWithdrawalAnnualOfficePremium,
		s.ExpAdjProportionTtdConversionOnWithdrawalRiskPremiumSalary, s.ExpAdjProportionTtdConversionOnWithdrawalOfficePremiumSalary,
		s.ExpTtdConversionOnWithdrawalRiskRatePer1000SA, s.ExpTtdConversionOnWithdrawalOfficeRatePer1000SA)...)

	return fs
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalGlaAnnualOfficePremium)},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalSglaAnnualOfficePremium)},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalPtdAnnualOfficePremium)},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalCiAnnualOfficePremium)},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalPhiAnnualOfficePremium)},
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
		{Key: "premium", Label: "Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalTtdAnnualOfficePremium)},
		{Key: "percent_salary", Label: "% of Salary", Value: formatPercent(s.ExpProportionTtdOfficePremiumSalary)},
	}
}

func funeralFields(s models.MemberRatingResultSummary, cat models.SchemeCategory, t quote_docx.BenefitTitles) []Field {
	return []Field{
		{Key: "title", Label: "Benefit Title", Value: t.FamilyFuneralBenefitTitle},
		{Key: "monthly_premium_per_member", Label: "Monthly Premium per Member", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunMonthlyPremiumPerMember)},
		{Key: "premium_per_member", Label: "Premium per Member", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.ExpTotalFunAnnualPremiumPerMember)},
		{Key: "total_premium", Label: "Total Premium", Value: quote_docx.RoundUpToTwoDecimalsAccounting(s.TotalFunAnnualOfficePremium)},
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
	)
	return categoryScalarFields(zs, zc, sampleBenefitNaming())
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

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
	fs := []Field{
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
	fs = append(fs, categoryRatingSummaryFields(s)...)
	fs = append(fs, categoryEducatorSummaryFields(s)...)
	fs = append(fs, categoryConversionSliceFields(s)...)
	return fs
}

// categoryRatingSummaryFields exposes the full member-rating result summary
// tokens at category scope: per-benefit risk rates, risk premiums, office
// premiums, proportions of salary, rate-per-1000 figures, binder/outsource
// splits, per-benefit commission, scheme-level commission totals, and tax
// saver slices. Keys/labels mirror the attached variable list with the
// leading "exp_" prefix stripped.
func categoryRatingSummaryFields(s models.MemberRatingResultSummary) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	return []Field{
		// GLA
		{Key: "total_gla_capped_sum_assured", Label: "total_gla_capped_sum_assured", Value: money(s.TotalGlaCappedSumAssured)},
		{Key: "total_gla_risk_rate", Label: "total_gla_risk_rate", Value: money(s.ExpTotalGlaRiskRate)},
		{Key: "total_gla_annual_risk_premium", Label: "total_gla_annual_risk_premium", Value: money(s.ExpTotalGlaAnnualRiskPremium)},
		{Key: "gla_risk_rate_per1000_sa", Label: "gla_risk_rate_per1000_sa", Value: money(s.ExpGlaRiskRatePer1000SA)},
		{Key: "proportion_gla_annual_risk_premium_salary", Label: "proportion_gla_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionGlaAnnualRiskPremiumSalary)},
		{Key: "total_gla_annual_office_premium", Label: "total_gla_annual_office_premium", Value: money(s.ExpTotalGlaAnnualOfficePremium)},
		{Key: "gla_office_rate_per1000_sa", Label: "gla_office_rate_per1000_sa", Value: money(s.ExpGlaOfficeRatePer1000SA)},
		{Key: "proportion_gla_office_premium_salary", Label: "proportion_gla_office_premium_salary", Value: formatPercent(s.ExpProportionGlaOfficePremiumSalary)},

		// PTD
		{Key: "total_ptd_capped_sum_assured", Label: "total_ptd_capped_sum_assured", Value: money(s.TotalPtdCappedSumAssured)},
		{Key: "total_ptd_risk_rate", Label: "total_ptd_risk_rate", Value: money(s.ExpTotalPtdRiskRate)},
		{Key: "total_ptd_annual_risk_premium", Label: "total_ptd_annual_risk_premium", Value: money(s.ExpTotalPtdAnnualRiskPremium)},
		{Key: "ptd_risk_rate_per1000_sa", Label: "ptd_risk_rate_per1000_sa", Value: money(s.ExpPtdRiskRatePer1000SA)},
		{Key: "proportion_ptd_annual_risk_premium_salary", Label: "proportion_ptd_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionPtdAnnualRiskPremiumSalary)},
		{Key: "total_ptd_annual_office_premium", Label: "total_ptd_annual_office_premium", Value: money(s.ExpTotalPtdAnnualOfficePremium)},
		{Key: "ptd_office_rate_per1000_sa", Label: "ptd_office_rate_per1000_sa", Value: money(s.ExpPtdOfficeRatePer1000SA)},
		{Key: "proportion_ptd_office_premium_salary", Label: "proportion_ptd_office_premium_salary", Value: formatPercent(s.ExpProportionPtdOfficePremiumSalary)},

		// CI
		{Key: "total_ci_capped_sum_assured", Label: "total_ci_capped_sum_assured", Value: money(s.TotalCiCappedSumAssured)},
		{Key: "total_ci_risk_rate", Label: "total_ci_risk_rate", Value: money(s.ExpTotalCiRiskRate)},
		{Key: "total_ci_annual_risk_premium", Label: "total_ci_annual_risk_premium", Value: money(s.ExpTotalCiAnnualRiskPremium)},
		{Key: "ci_risk_rate_per1000_sa", Label: "ci_risk_rate_per1000_sa", Value: money(s.ExpCiRiskRatePer1000SA)},
		{Key: "proportion_ci_annual_risk_premium_salary", Label: "proportion_ci_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionCiAnnualRiskPremiumSalary)},
		{Key: "total_ci_annual_office_premium", Label: "total_ci_annual_office_premium", Value: money(s.ExpTotalCiAnnualOfficePremium)},
		{Key: "ci_office_rate_per1000_sa", Label: "ci_office_rate_per1000_sa", Value: money(s.ExpCiOfficeRatePer1000SA)},
		{Key: "proportion_ci_office_premium_salary", Label: "proportion_ci_office_premium_salary", Value: formatPercent(s.ExpProportionCiOfficePremiumSalary)},

		// SGLA
		{Key: "total_sgla_capped_sum_assured", Label: "total_sgla_capped_sum_assured", Value: money(s.TotalSglaCappedSumAssured)},
		{Key: "total_sgla_risk_rate", Label: "total_sgla_risk_rate", Value: money(s.ExpTotalSglaRiskRate)},
		{Key: "total_sgla_annual_risk_premium", Label: "total_sgla_annual_risk_premium", Value: money(s.ExpTotalSglaAnnualRiskPremium)},
		{Key: "sgla_risk_rate_per1000_sa", Label: "sgla_risk_rate_per1000_sa", Value: money(s.ExpSglaRiskRatePer1000SA)},
		{Key: "proportion_sgla_annual_risk_premium_salary", Label: "proportion_sgla_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionSglaAnnualRiskPremiumSalary)},
		{Key: "total_sgla_annual_office_premium", Label: "total_sgla_annual_office_premium", Value: money(s.ExpTotalSglaAnnualOfficePremium)},
		{Key: "sgla_office_rate_per1000_sa", Label: "sgla_office_rate_per1000_sa", Value: money(s.ExpSglaOfficeRatePer1000SA)},
		{Key: "proportion_sgla_office_premium_salary", Label: "proportion_sgla_office_premium_salary", Value: formatPercent(s.ExpProportionSglaOfficePremiumSalary)},

		// TTD (income-based)
		{Key: "total_ttd_capped_income", Label: "total_ttd_capped_income", Value: money(s.TotalTtdCappedIncome)},
		{Key: "total_ttd_risk_rate", Label: "total_ttd_risk_rate", Value: money(s.ExpTotalTtdRiskRate)},
		{Key: "total_ttd_annual_risk_premium", Label: "total_ttd_annual_risk_premium", Value: money(s.ExpTotalTtdAnnualRiskPremium)},
		{Key: "ttd_risk_rate_per1000_sa", Label: "ttd_risk_rate_per1000_sa", Value: money(s.ExpTtdRiskRatePer1000SA)},
		{Key: "proportion_ttd_annual_risk_premium_salary", Label: "proportion_ttd_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionTtdAnnualRiskPremiumSalary)},
		{Key: "total_ttd_annual_office_premium", Label: "total_ttd_annual_office_premium", Value: money(s.ExpTotalTtdAnnualOfficePremium)},
		{Key: "ttd_office_rate_per1000_sa", Label: "ttd_office_rate_per1000_sa", Value: money(s.ExpTtdOfficeRatePer1000SA)},
		{Key: "proportion_ttd_office_premium_salary", Label: "proportion_ttd_office_premium_salary", Value: formatPercent(s.ExpProportionTtdOfficePremiumSalary)},

		// PHI (income-based)
		{Key: "total_phi_capped_income", Label: "total_phi_capped_income", Value: money(s.TotalPhiCappedIncome)},
		{Key: "total_phi_risk_rate", Label: "total_phi_risk_rate", Value: money(s.ExpTotalPhiRiskRate)},
		{Key: "total_phi_annual_risk_premium", Label: "total_phi_annual_risk_premium", Value: money(s.ExpTotalPhiAnnualRiskPremium)},
		{Key: "phi_risk_rate_per1000_sa", Label: "phi_risk_rate_per1000_sa", Value: money(s.ExpPhiRiskRatePer1000SA)},
		{Key: "proportion_phi_annual_risk_premium_salary", Label: "proportion_phi_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionPhiAnnualRiskPremiumSalary)},
		{Key: "total_phi_annual_office_premium", Label: "total_phi_annual_office_premium", Value: money(s.ExpTotalPhiAnnualOfficePremium)},
		{Key: "phi_office_rate_per1000_sa", Label: "phi_office_rate_per1000_sa", Value: money(s.ExpPhiOfficeRatePer1000SA)},
		{Key: "proportion_phi_office_premium_salary", Label: "proportion_phi_office_premium_salary", Value: formatPercent(s.ExpProportionPhiOfficePremiumSalary)},

		// Funeral + aggregate
		{Key: "total_fun_annual_risk_premium", Label: "total_fun_annual_risk_premium", Value: money(s.ExpTotalFunAnnualRiskPremium)},
		{Key: "proportion_fun_annual_risk_premium_salary", Label: "proportion_fun_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionFunAnnualRiskPremiumSalary)},
		{Key: "total_fun_annual_office_premium", Label: "total_fun_annual_office_premium", Value: money(s.ExpTotalFunAnnualOfficePremium)},
		{Key: "proportion_fun_office_premium_salary", Label: "proportion_fun_office_premium_salary", Value: formatPercent(s.ExpProportionFunOfficePremiumSalary)},
		{Key: "total_fun_annual_premium_per_member", Label: "total_fun_annual_premium_per_member", Value: money(s.ExpTotalFunAnnualPremiumPerMember)},
		{Key: "total_fun_monthly_premium_per_member", Label: "total_fun_monthly_premium_per_member", Value: money(s.ExpTotalFunMonthlyPremiumPerMember)},
		{Key: "total_annual_premium_excl_funeral", Label: "total_annual_premium_excl_funeral", Value: money(s.ExpTotalAnnualPremiumExclFuneral)},
		{Key: "proportion_exp_total_premium_excl_funeral_salary", Label: "proportion_exp_total_premium_excl_funeral_salary", Value: formatPercent(s.ProportionExpTotalPremiumExclFuneralSalary)},

		// Additional Accidental GLA
		{Key: "total_additional_accidental_gla_capped_sum_assured", Label: "total_additional_accidental_gla_capped_sum_assured", Value: money(s.TotalAdditionalAccidentalGlaCappedSumAssured)},
		{Key: "total_additional_accidental_gla_risk_rate", Label: "total_additional_accidental_gla_risk_rate", Value: money(s.ExpTotalAdditionalAccidentalGlaRiskRate)},
		{Key: "total_additional_accidental_gla_annual_risk_premium", Label: "total_additional_accidental_gla_annual_risk_premium", Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualRiskPremium)},
		{Key: "additional_accidental_gla_risk_rate_per1000_sa", Label: "additional_accidental_gla_risk_rate_per1000_sa", Value: money(s.ExpAdditionalAccidentalGlaRiskRatePer1000SA)},
		{Key: "prop_additional_accidental_gla_annual_risk_premium_salary", Label: "prop_additional_accidental_gla_annual_risk_premium_salary", Value: formatPercent(s.ExpProportionAdditionalAccidentalGlaAnnualRiskPremiumSalary)},
		{Key: "total_additional_accidental_gla_annual_office_premium", Label: "total_additional_accidental_gla_annual_office_premium", Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualOfficePremium)},
		{Key: "additional_accidental_gla_office_rate_per1000_sa", Label: "additional_accidental_gla_office_rate_per1000_sa", Value: money(s.ExpAdditionalAccidentalGlaOfficeRatePer1000SA)},
		{Key: "proportion_additional_accidental_gla_office_premium_salary", Label: "proportion_additional_accidental_gla_office_premium_salary", Value: formatPercent(s.ExpProportionAdditionalAccidentalGlaOfficePremiumSalary)},

		// Binder & outsource amounts (per benefit)
		{Key: "total_gla_annual_binder_amount", Label: "total_gla_annual_binder_amount", Value: money(s.ExpTotalGlaAnnualBinderAmount)},
		{Key: "total_gla_annual_outsourced_amount", Label: "total_gla_annual_outsourced_amount", Value: money(s.ExpTotalGlaAnnualOutsourcedAmount)},
		{Key: "total_add_acc_gla_annual_binder_amount", Label: "total_add_acc_gla_annual_binder_amount", Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualBinderAmount)},
		{Key: "total_add_acc_gla_annual_outsourced_amt", Label: "total_add_acc_gla_annual_outsourced_amt", Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualOutsourcedAmt)},
		{Key: "total_ptd_annual_binder_amount", Label: "total_ptd_annual_binder_amount", Value: money(s.ExpTotalPtdAnnualBinderAmount)},
		{Key: "total_ptd_annual_outsourced_amount", Label: "total_ptd_annual_outsourced_amount", Value: money(s.ExpTotalPtdAnnualOutsourcedAmount)},
		{Key: "total_ci_annual_binder_amount", Label: "total_ci_annual_binder_amount", Value: money(s.ExpTotalCiAnnualBinderAmount)},
		{Key: "total_ci_annual_outsourced_amount", Label: "total_ci_annual_outsourced_amount", Value: money(s.ExpTotalCiAnnualOutsourcedAmount)},
		{Key: "total_sgla_annual_binder_amount", Label: "total_sgla_annual_binder_amount", Value: money(s.ExpTotalSglaAnnualBinderAmount)},
		{Key: "total_sgla_annual_outsourced_amount", Label: "total_sgla_annual_outsourced_amount", Value: money(s.ExpTotalSglaAnnualOutsourcedAmount)},
		{Key: "total_ttd_annual_binder_amount", Label: "total_ttd_annual_binder_amount", Value: money(s.ExpTotalTtdAnnualBinderAmount)},
		{Key: "total_ttd_annual_outsourced_amount", Label: "total_ttd_annual_outsourced_amount", Value: money(s.ExpTotalTtdAnnualOutsourcedAmount)},
		{Key: "total_phi_annual_binder_amount", Label: "total_phi_annual_binder_amount", Value: money(s.ExpTotalPhiAnnualBinderAmount)},
		{Key: "total_phi_annual_outsourced_amount", Label: "total_phi_annual_outsourced_amount", Value: money(s.ExpTotalPhiAnnualOutsourcedAmount)},
		{Key: "total_fun_annual_binder_amount", Label: "total_fun_annual_binder_amount", Value: money(s.ExpTotalFunAnnualBinderAmount)},
		{Key: "total_fun_annual_outsourced_amount", Label: "total_fun_annual_outsourced_amount", Value: money(s.ExpTotalFunAnnualOutsourcedAmount)},

		// Commission amounts (per benefit + scheme totals)
		{Key: "total_gla_annual_commission_amount", Label: "total_gla_annual_commission_amount", Value: money(s.ExpTotalGlaAnnualCommissionAmount)},
		{Key: "total_add_acc_gla_annual_commission_amount", Label: "total_add_acc_gla_annual_commission_amount", Value: money(s.ExpTotalAdditionalAccidentalGlaAnnualCommissionAmount)},
		{Key: "total_ptd_annual_commission_amount", Label: "total_ptd_annual_commission_amount", Value: money(s.ExpTotalPtdAnnualCommissionAmount)},
		{Key: "total_ci_annual_commission_amount", Label: "total_ci_annual_commission_amount", Value: money(s.ExpTotalCiAnnualCommissionAmount)},
		{Key: "total_sgla_annual_commission_amount", Label: "total_sgla_annual_commission_amount", Value: money(s.ExpTotalSglaAnnualCommissionAmount)},
		{Key: "total_ttd_annual_commission_amount", Label: "total_ttd_annual_commission_amount", Value: money(s.ExpTotalTtdAnnualCommissionAmount)},
		{Key: "total_phi_annual_commission_amount", Label: "total_phi_annual_commission_amount", Value: money(s.ExpTotalPhiAnnualCommissionAmount)},
		{Key: "total_fun_annual_commission_amount", Label: "total_fun_annual_commission_amount", Value: money(s.ExpTotalFunAnnualCommissionAmount)},
		{Key: "scheme_total_commission", Label: "scheme_total_commission", Value: money(s.SchemeTotalCommission)},
		{Key: "scheme_total_commission_rate", Label: "scheme_total_commission_rate", Value: formatPercent(s.SchemeTotalCommissionRate)},

		// Tax saver slice (of GLA office premium)
		{Key: "total_tax_saver_annual_risk_premium", Label: "total_tax_saver_annual_risk_premium", Value: money(s.ExpTotalTaxSaverAnnualRiskPremium)},
		{Key: "total_tax_saver_annual_office_premium", Label: "total_tax_saver_annual_office_premium", Value: money(s.ExpTotalTaxSaverAnnualOfficePremium)},
	}
}

// categoryEducatorSummaryFields exposes the GLA/PTD educator split tokens:
// risk and office premiums, proportion-of-salary, rate-per-1000, plus binder,
// outsource, and commission breakdowns for each educator cover.
func categoryEducatorSummaryFields(s models.MemberRatingResultSummary) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	return []Field{
		// GLA educator
		{Key: "adj_total_gla_educator_risk_premium", Label: "adj_total_gla_educator_risk_premium", Value: money(s.ExpAdjTotalGlaEducatorRiskPremium)},
		{Key: "adj_total_gla_educator_office_premium", Label: "adj_total_gla_educator_office_premium", Value: money(s.ExpAdjTotalGlaEducatorOfficePremium)},
		{Key: "adj_proportion_gla_educator_risk_premium_salary", Label: "adj_proportion_gla_educator_risk_premium_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorRiskPremiumSalary)},
		{Key: "adj_proportion_gla_educator_office_premium_salary", Label: "adj_proportion_gla_educator_office_premium_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorOfficePremiumSalary)},
		{Key: "gla_educator_risk_rate_per1000_sa", Label: "gla_educator_risk_rate_per1000_sa", Value: money(s.ExpGlaEducatorRiskRatePer1000SA)},
		{Key: "gla_educator_office_rate_per1000_sa", Label: "gla_educator_office_rate_per1000_sa", Value: money(s.ExpGlaEducatorOfficeRatePer1000SA)},

		// PTD educator
		{Key: "adj_total_ptd_educator_risk_premium", Label: "adj_total_ptd_educator_risk_premium", Value: money(s.ExpAdjTotalPtdEducatorRiskPremium)},
		{Key: "adj_total_ptd_educator_office_premium", Label: "adj_total_ptd_educator_office_premium", Value: money(s.ExpAdjTotalPtdEducatorOfficePremium)},
		{Key: "adj_proportion_ptd_educator_risk_premium_salary", Label: "adj_proportion_ptd_educator_risk_premium_salary", Value: formatPercent(s.ExpAdjProportionPtdEducatorRiskPremiumSalary)},
		{Key: "adj_proportion_ptd_educator_office_premium_salary", Label: "adj_proportion_ptd_educator_office_premium_salary", Value: formatPercent(s.ExpAdjProportionPtdEducatorOfficePremiumSalary)},
		{Key: "ptd_educator_risk_rate_per1000_sa", Label: "ptd_educator_risk_rate_per1000_sa", Value: money(s.ExpPtdEducatorRiskRatePer1000SA)},
		{Key: "ptd_educator_office_rate_per1000_sa", Label: "ptd_educator_office_rate_per1000_sa", Value: money(s.ExpPtdEducatorOfficeRatePer1000SA)},

		// Educator binder / outsourced / commission
		{Key: "adj_total_gla_educator_binder_amount", Label: "adj_total_gla_educator_binder_amount", Value: money(s.ExpAdjTotalGlaEducatorBinderAmount)},
		{Key: "adj_total_gla_educator_outsourced_amount", Label: "adj_total_gla_educator_outsourced_amount", Value: money(s.ExpAdjTotalGlaEducatorOutsourcedAmount)},
		{Key: "adj_total_ptd_educator_binder_amount", Label: "adj_total_ptd_educator_binder_amount", Value: money(s.ExpAdjTotalPtdEducatorBinderAmount)},
		{Key: "adj_total_ptd_educator_outsourced_amount", Label: "adj_total_ptd_educator_outsourced_amount", Value: money(s.ExpAdjTotalPtdEducatorOutsourcedAmount)},
		{Key: "adj_total_gla_educator_commission_amount", Label: "adj_total_gla_educator_commission_amount", Value: money(s.ExpAdjTotalGlaEducatorCommissionAmount)},
		{Key: "adj_total_ptd_educator_commission_amount", Label: "adj_total_ptd_educator_commission_amount", Value: money(s.ExpAdjTotalPtdEducatorCommissionAmount)},
	}
}

// categoryConversionSliceFields exposes conversion / continuity slice tokens.
// Each slice carries six variants: annual risk premium, annual office premium,
// proportion of salary (risk + office) and rate-per-1000 (risk + office).
func categoryConversionSliceFields(s models.MemberRatingResultSummary) []Field {
	money := quote_docx.RoundUpToTwoDecimalsAccounting
	return []Field{
		// GLA conversion on withdrawal
		{Key: "adj_total_gla_conv_on_wdr_ann_risk_prem", Label: "adj_total_gla_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalGlaConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_gla_conv_on_wdr_ann_office_prem", Label: "adj_total_gla_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalGlaConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_gla_conv_on_wdr_risk_prem_salary", Label: "adj_prop_gla_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_gla_conv_on_wdr_office_prem_salary", Label: "adj_prop_gla_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "gla_conv_on_wdr_risk_rate_per_1000_sa", Label: "gla_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpGlaConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "gla_conv_on_wdr_office_rate_per_1000_sa", Label: "gla_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpGlaConversionOnWithdrawalOfficeRatePer1000SA)},

		// GLA conversion on retirement
		{Key: "adj_total_gla_conv_on_ret_ann_risk_prem", Label: "adj_total_gla_conv_on_ret_ann_risk_prem", Value: money(s.ExpAdjTotalGlaConversionOnRetirementAnnualRiskPremium)},
		{Key: "adj_total_gla_conv_on_ret_ann_office_prem", Label: "adj_total_gla_conv_on_ret_ann_office_prem", Value: money(s.ExpAdjTotalGlaConversionOnRetirementAnnualOfficePremium)},
		{Key: "adj_prop_gla_conv_on_ret_risk_prem_salary", Label: "adj_prop_gla_conv_on_ret_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaConversionOnRetirementRiskPremiumSalary)},
		{Key: "adj_prop_gla_conv_on_ret_office_prem_salary", Label: "adj_prop_gla_conv_on_ret_office_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaConversionOnRetirementOfficePremiumSalary)},
		{Key: "gla_conv_on_ret_risk_rate_per_1000_sa", Label: "gla_conv_on_ret_risk_rate_per_1000_sa", Value: money(s.ExpGlaConversionOnRetirementRiskRatePer1000SA)},
		{Key: "gla_conv_on_ret_office_rate_per_1000_sa", Label: "gla_conv_on_ret_office_rate_per_1000_sa", Value: money(s.ExpGlaConversionOnRetirementOfficeRatePer1000SA)},

		// GLA continuity during disability
		{Key: "adj_total_gla_cont_dur_dis_ann_risk_prem", Label: "adj_total_gla_cont_dur_dis_ann_risk_prem", Value: money(s.ExpAdjTotalGlaContinuityDuringDisabilityAnnualRiskPremium)},
		{Key: "adj_total_gla_cont_dur_dis_ann_office_prem", Label: "adj_total_gla_cont_dur_dis_ann_office_prem", Value: money(s.ExpAdjTotalGlaContinuityDuringDisabilityAnnualOfficePremium)},
		{Key: "adj_prop_gla_cont_dur_dis_risk_prem_salary", Label: "adj_prop_gla_cont_dur_dis_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaContinuityDuringDisabilityRiskPremiumSalary)},
		{Key: "adj_prop_gla_cont_dur_dis_office_prem_salary", Label: "adj_prop_gla_cont_dur_dis_office_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaContinuityDuringDisabilityOfficePremiumSalary)},
		{Key: "gla_cont_dur_dis_risk_rate_per_1000_sa", Label: "gla_cont_dur_dis_risk_rate_per_1000_sa", Value: money(s.ExpGlaContinuityDuringDisabilityRiskRatePer1000SA)},
		{Key: "gla_cont_dur_dis_office_rate_per_1000_sa", Label: "gla_cont_dur_dis_office_rate_per_1000_sa", Value: money(s.ExpGlaContinuityDuringDisabilityOfficeRatePer1000SA)},

		// GLA educator conversion on withdrawal
		{Key: "adj_total_gla_ed_conv_on_wdr_ann_risk_prem", Label: "adj_total_gla_ed_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_gla_ed_conv_on_wdr_ann_office_prem", Label: "adj_total_gla_ed_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalGlaEducatorConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_gla_ed_conv_on_wdr_risk_prem_salary", Label: "adj_prop_gla_ed_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_gla_ed_conv_on_wdr_office_prem_salary", Label: "adj_prop_gla_ed_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "gla_ed_conv_on_wdr_risk_rate_per_1000_sa", Label: "gla_ed_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpGlaEducatorConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "gla_ed_conv_on_wdr_office_rate_per_1000_sa", Label: "gla_ed_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpGlaEducatorConversionOnWithdrawalOfficeRatePer1000SA)},

		// GLA educator conversion on retirement
		{Key: "adj_total_gla_ed_conv_on_ret_ann_risk_prem", Label: "adj_total_gla_ed_conv_on_ret_ann_risk_prem", Value: money(s.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualRiskPremium)},
		{Key: "adj_total_gla_ed_conv_on_ret_ann_office_prem", Label: "adj_total_gla_ed_conv_on_ret_ann_office_prem", Value: money(s.ExpAdjTotalGlaEducatorConversionOnRetirementAnnualOfficePremium)},
		{Key: "adj_prop_gla_ed_conv_on_ret_risk_prem_salary", Label: "adj_prop_gla_ed_conv_on_ret_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorConversionOnRetirementRiskPremiumSalary)},
		{Key: "adj_prop_gla_ed_conv_on_ret_office_prem_salary", Label: "adj_prop_gla_ed_conv_on_ret_office_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorConversionOnRetirementOfficePremiumSalary)},
		{Key: "gla_ed_conv_on_ret_risk_rate_per_1000_sa", Label: "gla_ed_conv_on_ret_risk_rate_per_1000_sa", Value: money(s.ExpGlaEducatorConversionOnRetirementRiskRatePer1000SA)},
		{Key: "gla_ed_conv_on_ret_office_rate_per_1000_sa", Label: "gla_ed_conv_on_ret_office_rate_per_1000_sa", Value: money(s.ExpGlaEducatorConversionOnRetirementOfficeRatePer1000SA)},

		// GLA educator continuity during disability
		{Key: "adj_total_gla_ed_cont_dur_dis_ann_risk_prem", Label: "adj_total_gla_ed_cont_dur_dis_ann_risk_prem", Value: money(s.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualRiskPremium)},
		{Key: "adj_total_gla_ed_cont_dur_dis_ann_office_prem", Label: "adj_total_gla_ed_cont_dur_dis_ann_office_prem", Value: money(s.ExpAdjTotalGlaEducatorContinuityDuringDisabilityAnnualOfficePremium)},
		{Key: "adj_prop_gla_ed_cont_dur_dis_risk_prem_salary", Label: "adj_prop_gla_ed_cont_dur_dis_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorContinuityDuringDisabilityRiskPremiumSalary)},
		{Key: "adj_prop_gla_ed_cont_dur_dis_office_prem_salary", Label: "adj_prop_gla_ed_cont_dur_dis_office_prem_salary", Value: formatPercent(s.ExpAdjProportionGlaEducatorContinuityDuringDisabilityOfficePremiumSalary)},
		{Key: "gla_ed_cont_dur_dis_risk_rate_per_1000_sa", Label: "gla_ed_cont_dur_dis_risk_rate_per_1000_sa", Value: money(s.ExpGlaEducatorContinuityDuringDisabilityRiskRatePer1000SA)},
		{Key: "gla_ed_cont_dur_dis_office_rate_per_1000_sa", Label: "gla_ed_cont_dur_dis_office_rate_per_1000_sa", Value: money(s.ExpGlaEducatorContinuityDuringDisabilityOfficeRatePer1000SA)},

		// PTD conversion on withdrawal
		{Key: "adj_total_ptd_conv_on_wdr_ann_risk_prem", Label: "adj_total_ptd_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalPtdConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_ptd_conv_on_wdr_ann_office_prem", Label: "adj_total_ptd_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalPtdConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_ptd_conv_on_wdr_risk_prem_salary", Label: "adj_prop_ptd_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionPtdConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_ptd_conv_on_wdr_office_prem_salary", Label: "adj_prop_ptd_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionPtdConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "ptd_conv_on_wdr_risk_rate_per_1000_sa", Label: "ptd_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpPtdConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "ptd_conv_on_wdr_office_rate_per_1000_sa", Label: "ptd_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpPtdConversionOnWithdrawalOfficeRatePer1000SA)},

		// PTD educator conversion on withdrawal
		{Key: "adj_total_ptd_ed_conv_on_wdr_ann_risk_prem", Label: "adj_total_ptd_ed_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_ptd_ed_conv_on_wdr_ann_office_prem", Label: "adj_total_ptd_ed_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalPtdEducatorConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_ptd_ed_conv_on_wdr_risk_prem_salary", Label: "adj_prop_ptd_ed_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionPtdEducatorConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_ptd_ed_conv_on_wdr_office_prem_salary", Label: "adj_prop_ptd_ed_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionPtdEducatorConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "ptd_ed_conv_on_wdr_risk_rate_per_1000_sa", Label: "ptd_ed_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpPtdEducatorConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "ptd_ed_conv_on_wdr_office_rate_per_1000_sa", Label: "ptd_ed_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpPtdEducatorConversionOnWithdrawalOfficeRatePer1000SA)},

		// PTD educator conversion on retirement
		{Key: "adj_total_ptd_ed_conv_on_ret_ann_risk_prem", Label: "adj_total_ptd_ed_conv_on_ret_ann_risk_prem", Value: money(s.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualRiskPremium)},
		{Key: "adj_total_ptd_ed_conv_on_ret_ann_office_prem", Label: "adj_total_ptd_ed_conv_on_ret_ann_office_prem", Value: money(s.ExpAdjTotalPtdEducatorConversionOnRetirementAnnualOfficePremium)},
		{Key: "adj_prop_ptd_ed_conv_on_ret_risk_prem_salary", Label: "adj_prop_ptd_ed_conv_on_ret_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionPtdEducatorConversionOnRetirementRiskPremiumSalary)},
		{Key: "adj_prop_ptd_ed_conv_on_ret_office_prem_salary", Label: "adj_prop_ptd_ed_conv_on_ret_office_prem_salary", Value: formatPercent(s.ExpAdjProportionPtdEducatorConversionOnRetirementOfficePremiumSalary)},
		{Key: "ptd_ed_conv_on_ret_risk_rate_per_1000_sa", Label: "ptd_ed_conv_on_ret_risk_rate_per_1000_sa", Value: money(s.ExpPtdEducatorConversionOnRetirementRiskRatePer1000SA)},
		{Key: "ptd_ed_conv_on_ret_office_rate_per_1000_sa", Label: "ptd_ed_conv_on_ret_office_rate_per_1000_sa", Value: money(s.ExpPtdEducatorConversionOnRetirementOfficeRatePer1000SA)},

		// PHI conversion on withdrawal
		{Key: "adj_total_phi_conv_on_wdr_ann_risk_prem", Label: "adj_total_phi_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalPhiConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_phi_conv_on_wdr_ann_office_prem", Label: "adj_total_phi_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalPhiConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_phi_conv_on_wdr_risk_prem_salary", Label: "adj_prop_phi_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionPhiConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_phi_conv_on_wdr_office_prem_salary", Label: "adj_prop_phi_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionPhiConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "phi_conv_on_wdr_risk_rate_per_1000_sa", Label: "phi_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpPhiConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "phi_conv_on_wdr_office_rate_per_1000_sa", Label: "phi_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpPhiConversionOnWithdrawalOfficeRatePer1000SA)},

		// CI conversion on withdrawal
		{Key: "adj_total_ci_conv_on_wdr_ann_risk_prem", Label: "adj_total_ci_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalCiConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_ci_conv_on_wdr_ann_office_prem", Label: "adj_total_ci_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalCiConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_ci_conv_on_wdr_risk_prem_salary", Label: "adj_prop_ci_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionCiConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_ci_conv_on_wdr_office_prem_salary", Label: "adj_prop_ci_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionCiConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "ci_conv_on_wdr_risk_rate_per_1000_sa", Label: "ci_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpCiConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "ci_conv_on_wdr_office_rate_per_1000_sa", Label: "ci_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpCiConversionOnWithdrawalOfficeRatePer1000SA)},

		// SGLA conversion on withdrawal
		{Key: "adj_total_sgla_conv_on_wdr_ann_risk_prem", Label: "adj_total_sgla_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalSglaConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_sgla_conv_on_wdr_ann_office_prem", Label: "adj_total_sgla_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalSglaConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_sgla_conv_on_wdr_risk_prem_salary", Label: "adj_prop_sgla_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionSglaConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_sgla_conv_on_wdr_office_prem_salary", Label: "adj_prop_sgla_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionSglaConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "sgla_conv_on_wdr_risk_rate_per_1000_sa", Label: "sgla_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpSglaConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "sgla_conv_on_wdr_office_rate_per_1000_sa", Label: "sgla_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpSglaConversionOnWithdrawalOfficeRatePer1000SA)},

		// Funeral conversion on withdrawal
		{Key: "adj_total_fun_conv_on_wdr_ann_risk_prem", Label: "adj_total_fun_conv_on_wdr_ann_risk_prem", Value: money(s.ExpAdjTotalFunConversionOnWithdrawalAnnualRiskPremium)},
		{Key: "adj_total_fun_conv_on_wdr_ann_office_prem", Label: "adj_total_fun_conv_on_wdr_ann_office_prem", Value: money(s.ExpAdjTotalFunConversionOnWithdrawalAnnualOfficePremium)},
		{Key: "adj_prop_fun_conv_on_wdr_risk_prem_salary", Label: "adj_prop_fun_conv_on_wdr_risk_prem_salary", Value: formatPercent(s.ExpAdjProportionFunConversionOnWithdrawalRiskPremiumSalary)},
		{Key: "adj_prop_fun_conv_on_wdr_office_prem_salary", Label: "adj_prop_fun_conv_on_wdr_office_prem_salary", Value: formatPercent(s.ExpAdjProportionFunConversionOnWithdrawalOfficePremiumSalary)},
		{Key: "fun_conv_on_wdr_risk_rate_per_1000_sa", Label: "fun_conv_on_wdr_risk_rate_per_1000_sa", Value: money(s.ExpFunConversionOnWithdrawalRiskRatePer1000SA)},
		{Key: "fun_conv_on_wdr_office_rate_per_1000_sa", Label: "fun_conv_on_wdr_office_rate_per_1000_sa", Value: money(s.ExpFunConversionOnWithdrawalOfficeRatePer1000SA)},
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
		{Key: "has_additional_gla", Label: "Category has Additional GLA Cover", Value: flags.AdditionalGla},
		{Key: "has_additional_accidental_gla", Label: "Category has Additional Accidental GLA", Value: flags.AdditionalAccidentalGla},
		{Key: "has_gla_educator", Label: "Category has GLA Educator benefit", Value: flags.GlaEducator},
		{Key: "has_ptd_educator", Label: "Category has PTD Educator benefit", Value: flags.PtdEducator},
		{Key: "has_extended_family", Label: "Category has Extended Family Funeral", Value: flags.ExtendedFamily},
		{Key: "has_tax_saver", Label: "Category has Tax Saver benefit", Value: flags.TaxSaver},
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

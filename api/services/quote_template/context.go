package quote_template

import (
	"fmt"

	"api/models"
	"api/services"
	"api/services/quote_docx"
)

// Context is the map of template variables handed to the render engine.
// Keys nested under "insurer" and each entry of "categories" are accessed
// via dot syntax in templates (e.g. {{insurer.name}}, {{gla.premium}}
// inside a {{#categories}} block).
//
// Every key in this map is produced by a *Fields function in schema.go.
// schema.go is the single source of truth — add new tokens there and they
// flow automatically into both the rendered quote and the sample template.
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
	naming := resolveBenefitNaming(benefitMaps)
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
			cat = models.SchemeCategory{SchemeCategory: s.Category}
		}
		categories = append(categories, buildCategoryMap(s, cat, quote, titles, naming))
	}

	// Fold the root-scope fields into the Context map, then layer the
	// nested "insurer" object and the "categories" list on top.
	ctx := Context(fieldsToMap(quoteFields(quote, totals, hasNonFuneral)))
	ctx["insurer"] = fieldsToMap(insurerFields(insurer))
	ctx["categories"] = categories
	return ctx, nil
}

// buildCategoryMap produces the per-category map used inside {{#categories}}.
// It layers the scalar fields, the has_* bool flags, and the seven benefit
// sub-objects into a single map with the exact keys the template expects.
func buildCategoryMap(
	s models.MemberRatingResultSummary,
	cat models.SchemeCategory,
	quote models.GroupPricingQuote,
	titles quote_docx.BenefitTitles,
	naming benefitNaming,
) map[string]interface{} {
	flags := deriveBenefitFlags(s, cat)

	m := fieldsToMap(categoryScalarFields(s, cat, naming))
	for k, v := range fieldsToMap(categoryBoolFields(s, flags, naming)) {
		m[k] = v
	}
	m[naming.GLA.Code] = benefitMap(flags.GLA, glaFields(s, cat, quote, titles))
	m[naming.SGLA.Code] = benefitMap(flags.SGLA, sglaFields(s, cat, quote, titles))
	m[naming.PTD.Code] = benefitMap(flags.PTD, ptdFields(s, cat, quote, titles))
	m[naming.CI.Code] = benefitMap(flags.CI, ciFields(s, cat, quote, titles))
	m[naming.PHI.Code] = benefitMap(flags.PHI, phiFields(s, cat, titles))
	m[naming.TTD.Code] = benefitMap(flags.TTD, ttdFields(s, cat, titles))
	m[naming.Funeral.Code] = benefitMap(flags.Funeral, funeralFields(s, cat, titles))
	return m
}

// --- formatting helpers (shared by schema.go) ---

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

// proportionOfSalary returns amount / salary, guarding the zero-salary case
// so an empty fixture (or a category with no recorded salary) renders 0
// instead of NaN.
func proportionOfSalary(amount, salary float64) float64 {
	if salary <= 0 {
		return 0
	}
	return amount / salary
}

// orDash returns "-" when the input is empty, otherwise the input unchanged.
func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}


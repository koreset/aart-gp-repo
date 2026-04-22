package on_risk_letter_template

import (
	"fmt"
	"strconv"

	"api/services"
	"api/services/quote_docx"
	"api/services/quote_template"
)

// BuildContext assembles the complete context for rendering an on-risk
// letter template. The base context is produced by quote_template.BuildContext
// so every token a quote PDF template exposes (quote scalars, {{insurer.*}},
// the {{#categories}} iteration with its nested benefit sub-objects and
// the full rating/educator/conversion-slice surface) is also available to
// letter templates. Letter-specific tokens — letter metadata, cover end
// date, broker/channel info, identifiers, and the flat {{#benefit_summary}}
// cross-category rollup — are layered on top.
func BuildContext(quoteID string) (quote_template.Context, error) {
	quoteIDInt, err := strconv.Atoi(quoteID)
	if err != nil {
		return nil, fmt.Errorf("invalid quote id: %w", err)
	}

	ctx, err := quote_template.BuildContext(quoteID)
	if err != nil {
		return nil, fmt.Errorf("build base quote context: %w", err)
	}

	letterData, err := services.GetOnRiskLetterData(quoteIDInt)
	if err != nil {
		return nil, fmt.Errorf("fetch on-risk letter data: %w", err)
	}

	quote := letterData.Quote
	scheme := letterData.Scheme
	letter := letterData.Letter
	benefitSummary := letterData.BenefitSummary

	schemeContact := quote.SchemeContact
	if schemeContact == "" {
		schemeContact = scheme.ContactPerson
	}
	schemeEmail := quote.SchemeEmail
	if schemeEmail == "" {
		schemeEmail = scheme.ContactEmail
	}

	brokerName := quote.QuoteBroker.Name
	isBrokerChannel := quote.DistributionChannel == "broker"

	benefitSummaryMaps := make([]map[string]interface{}, len(benefitSummary))
	for i, line := range benefitSummary {
		benefitSummaryMaps[i] = map[string]interface{}{
			"benefit":        line.Benefit,
			"annual_premium": quote_docx.RoundUpToTwoDecimalsAccounting(line.AnnualPremium),
		}
	}

	ctx["letter_reference"] = letter.LetterReference
	ctx["letter_date"] = quote_docx.FormatQuoteDate(letter.LetterDate)
	ctx["generated_by"] = letter.GeneratedBy
	ctx["cover_end_date"] = quote_docx.FormatQuoteDate(quote.CoverEndDate)
	ctx["scheme_contact"] = schemeContact
	ctx["scheme_email"] = schemeEmail
	ctx["distribution_channel"] = string(quote.DistributionChannel)
	ctx["broker_name"] = brokerName
	ctx["is_broker_channel"] = isBrokerChannel
	ctx["member_count"] = strconv.Itoa(quote.MemberDataCount)
	ctx["quote_id"] = strconv.Itoa(quote.ID)
	ctx["quote_reference"] = quote.QuoteName
	ctx["has_benefits"] = len(benefitSummary) > 0
	ctx["benefit_summary"] = benefitSummaryMaps

	return ctx, nil
}

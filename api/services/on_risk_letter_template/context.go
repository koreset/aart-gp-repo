package on_risk_letter_template

import (
	"fmt"
	"strconv"

	"api/services"
	"api/services/quote_docx"
	"api/services/quote_template"
)

// BuildContext assembles the complete context for rendering an on-risk letter template
func BuildContext(quoteID string) (quote_template.Context, error) {
	// Convert string to int
	quoteIDInt, err := strconv.Atoi(quoteID)
	if err != nil {
		return nil, fmt.Errorf("invalid quote id: %w", err)
	}

	// Fetch on-risk letter data
	letterData, err := services.GetOnRiskLetterData(quoteIDInt)
	if err != nil {
		return nil, fmt.Errorf("fetch on-risk letter data: %w", err)
	}

	quote := letterData.Quote
	scheme := letterData.Scheme
	insurer := letterData.Insurer
	letter := letterData.Letter
	benefitSummary := letterData.BenefitSummary

	// Calculate total annual premium
	totalAnnualPremium := 0.0
	for _, benefit := range benefitSummary {
		totalAnnualPremium += benefit.AnnualPremium
	}

	// Determine scheme contact info (prefer quote fields, fall back to scheme fields)
	schemeContact := quote.SchemeContact
	if schemeContact == "" {
		schemeContact = scheme.ContactPerson
	}
	schemeEmail := quote.SchemeEmail
	if schemeEmail == "" {
		schemeEmail = scheme.ContactEmail
	}

	// Determine broker name
	brokerName := ""
	if quote.QuoteBroker.Name != "" {
		brokerName = quote.QuoteBroker.Name
	}

	// Determine distribution channel flag
	isBrokerChannel := quote.DistributionChannel == "broker"

	// Convert benefit summary to proper format for templating
	benefitSummaryMaps := make([]map[string]interface{}, len(benefitSummary))
	for i, line := range benefitSummary {
		benefitSummaryMaps[i] = map[string]interface{}{
			"benefit":          line.Benefit,
			"annual_premium":   quote_docx.RoundUpToTwoDecimalsAccounting(line.AnnualPremium),
		}
	}

	// Build insurer map with on_risk_letter_text
	insurerMap := map[string]interface{}{
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
		"on_risk_letter_text":     insurer.OnRiskLetterText,
		"introductory_text":       insurer.IntroductoryText,
		"general_provisions_text": insurer.GeneralProvisionsText,
	}

	return quote_template.Context{
		"letter_reference":      letter.LetterReference,
		"letter_date":           quote_docx.FormatQuoteDate(letter.LetterDate),
		"generated_by":          letter.GeneratedBy,
		"scheme_name":           quote.SchemeName,
		"scheme_contact":        schemeContact,
		"scheme_email":          schemeEmail,
		"commencement_date":     quote_docx.FormatQuoteDate(quote.CommencementDate),
		"cover_end_date":        quote_docx.FormatQuoteDate(quote.CoverEndDate),
		"industry":              quote.Industry,
		"obligation_type":       quote.ObligationType,
		"currency":              quote.Currency,
		"free_cover_limit":      quote_docx.RoundUpToTwoDecimalsAccounting(quote.FreeCoverLimit),
		"normal_retirement_age": strconv.Itoa(quote.NormalRetirementAge),
		"distribution_channel":  string(quote.DistributionChannel),
		"broker_name":           brokerName,
		"is_broker_channel":     isBrokerChannel,
		"member_count":          strconv.Itoa(quote.MemberDataCount),
		"quote_id":              strconv.Itoa(quote.ID),
		"quote_reference":       quote.QuoteName,
		"total_annual_premium":  quote_docx.RoundUpToTwoDecimalsAccounting(totalAnnualPremium),
		"has_benefits":          len(benefitSummary) > 0,
		"insurer":               insurerMap,
		"benefit_summary":       benefitSummaryMaps,
	}, nil
}

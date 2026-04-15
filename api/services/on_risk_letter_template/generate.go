package on_risk_letter_template

import (
	"fmt"
	"regexp"
	"time"

	"api/services"
	"api/services/quote_template"
)

// GenerateOnRiskLetterDocx is the top-level entry point for backend-templated on-risk letter generation
func GenerateOnRiskLetterDocx(quoteID string) (filename string, data []byte, err error) {
	// Fetch quote to get insurer info
	quote, err := services.GetGroupPricingQuote(quoteID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	// Fetch insurer (single global insurer)
	insurer, err := services.GetInsurerDetails()
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch insurer: %w", err)
	}

	// Check if active template exists for this insurer
	activeTemplate, err := GetActiveTemplate(insurer.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to check for active template: %w", err)
	}

	if activeTemplate == nil {
		return "", nil, fmt.Errorf("no active on-risk letter template — upload one in the insurer settings to use the backend-templated path")
	}

	// Build context
	ctx, err := BuildContext(quoteID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to build template context: %w", err)
	}

	// Render template
	renderedBytes, err := quote_template.Render(activeTemplate.DocxBlob, ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to render template: %w", err)
	}

	// Generate filename: {schemeName}_On_Risk_Letter_{YYYY-MM-DD}.docx
	sanitizedScheme := sanitizeFilename(quote.SchemeName)
	dateStr := time.Now().Format("2006-01-02")
	filename = fmt.Sprintf("%s_On_Risk_Letter_%s.docx", sanitizedScheme, dateStr)

	return filename, renderedBytes, nil
}

// sanitizeFilename removes/replaces characters invalid in Windows filenames
func sanitizeFilename(s string) string {
	invalidChars := regexp.MustCompile(`[/\\:*?"<>|]`)
	return invalidChars.ReplaceAllString(s, "_")
}

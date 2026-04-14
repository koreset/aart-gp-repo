package quote_docx

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"api/services"
)

// GenerateQuoteDocx is the top-level entry point for DOCX quote generation
func GenerateQuoteDocx(quoteID string) (filename string, data []byte, err error) {
	// 1. Fetch quote details
	quote, err := services.GetGroupPricingQuote(quoteID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	quoteIdInt := quote.ID

	// 2. Fetch result summaries
	resultSummaries, err := services.GetGroupPricingQuoteResultSummary(quoteIdInt)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch result summary: %w", err)
	}

	// 3. Fetch insurer details
	insurer, err := services.GetInsurerDetails()
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch insurer details: %w", err)
	}

	// 4. Fetch educator benefits (optional)
	var educatorBenefits []interface{}
	educatorBenefitsData, err := services.GetGroupPricingQuoteEducatorBenefits(quoteIdInt)
	if err == nil && educatorBenefitsData != nil {
		for _, eb := range educatorBenefitsData {
			educatorBenefits = append(educatorBenefits, eb)
		}
	}

	// 5. Fetch benefit maps
	benefitMaps, err := services.GetBenefitMaps()
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch benefit maps: %w", err)
	}

	// 6. Resolve benefit titles
	benefitTitles := ResolveBenefitTitles(benefitMaps)

	// 8. Build all sections
	var bodyBuf strings.Builder

	// Section 1: Cover and Summary
	section1 := BuildCoverAndSummarySection(quote, resultSummaries, insurer)
	bodyBuf.WriteString(section1)

	// Section 2: Premium Summary
	section2 := BuildPremiumSummarySection(quote, resultSummaries)
	bodyBuf.WriteString(section2)

	// Section 3: Premium Breakdown
	section3 := BuildPremiumBreakdownSection(quote, resultSummaries, benefitTitles)
	bodyBuf.WriteString(section3)

	// Section 4: Benefits and Definitions
	section4 := BuildBenefitsDefinitionsSection(quote, resultSummaries, educatorBenefits, benefitTitles)
	bodyBuf.WriteString(section4)

	// Section 5: Provisions
	section5 := BuildProvisionsSection(quote, insurer)
	bodyBuf.WriteString(section5)

	// Section 6: Acceptance Form
	section6 := BuildAcceptanceFormSection(quote)
	bodyBuf.WriteString(section6)

	// 9. Build DOCX package
	logoMIME := insurer.LogoMimeType
	if logoMIME == "" && len(insurer.Logo) > 0 {
		logoMIME = "image/png"
	}

	body := fixupSectionBreaks(bodyBuf.String())

	data, err = BuildFullPackage(body, quote.SchemeName, insurer.Logo, logoMIME)
	if err != nil {
		return "", nil, fmt.Errorf("failed to build DOCX package: %w", err)
	}

	// 10. Generate filename
	// Sanitize scheme name: replace special characters with underscores
	sanitizedScheme := sanitizeFilename(quote.SchemeName)
	dateStr := time.Now().Format("2006-01-02")
	filename = fmt.Sprintf("%s_Quotation_%s.docx", sanitizedScheme, dateStr)

	return filename, data, nil
}

// sanitizeFilename removes/replaces characters invalid in Windows filenames
func sanitizeFilename(s string) string {
	// Replace invalid characters with underscore
	invalidChars := regexp.MustCompile(`[/\\:*?"<>|]`)
	return invalidChars.ReplaceAllString(s, "_")
}

// fixupSectionBreaks reworks the concatenated section bodies so that they
// form a valid OOXML document. Section builders each append a <w:sectPr>
// directly to their content; that's incorrect for every section except the
// last one, which is the only sectPr allowed as a direct child of <w:body>.
// All other sectPr elements must live inside an (otherwise empty) paragraph's
// <w:pPr>. This function finds each <w:sectPr>...</w:sectPr>, wraps every
// one except the last, and leaves the final one alone.
func fixupSectionBreaks(body string) string {
	const open = "<w:sectPr>"
	const close = "</w:sectPr>"

	// Collect (start, end) offsets for every sectPr in order.
	type span struct{ start, end int }
	var spans []span
	i := 0
	for {
		s := strings.Index(body[i:], open)
		if s < 0 {
			break
		}
		s += i
		e := strings.Index(body[s:], close)
		if e < 0 {
			break
		}
		e += s + len(close)
		spans = append(spans, span{s, e})
		i = e
	}

	if len(spans) < 2 {
		return body // 0 or 1 sectPr — nothing to wrap
	}

	// Rebuild the body, wrapping every sectPr except the last.
	var out strings.Builder
	prev := 0
	for idx, sp := range spans {
		out.WriteString(body[prev:sp.start])
		chunk := body[sp.start:sp.end]
		if idx == len(spans)-1 {
			out.WriteString(chunk) // final sectPr stays as-is
		} else {
			out.WriteString(wrapSectPrInParagraph(chunk))
		}
		prev = sp.end
	}
	out.WriteString(body[prev:])
	return out.String()
}

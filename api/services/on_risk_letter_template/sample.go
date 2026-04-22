package on_risk_letter_template

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"strings"

	"api/services/quote_template"
)

// BuildSampleTemplate generates a sample On Risk Letter template DOCX documenting all available tokens
func BuildSampleTemplate() ([]byte, error) {
	var body strings.Builder

	// Title page
	body.WriteString(titlePage())

	// How tokens work
	body.WriteString(howTokensWork())

	// Letter-level tokens (letter metadata — unique to on-risk letters)
	body.WriteString(sectionLetterTokens())

	// Quote-level tokens (full shared surface from the quote template schema)
	body.WriteString(sectionQuoteTokens())

	// On-risk letter quote-context tokens (cover end, scheme contact, broker,
	// identifiers — sourced from the quote but not in the quote PDF token set)
	body.WriteString(sectionQuoteContextExtended())

	// Insurer tokens (driven by shared schema, incl. on_risk_letter_text)
	body.WriteString(sectionInsurerTokens())

	// Conditional flags
	body.WriteString(sectionConditionalFlags())

	// Categories iteration (inherited from the quote template surface)
	body.WriteString(sectionCategoriesIteration())

	// Nested benefits inside the categories block
	body.WriteString(sectionNestedBenefits())

	// Iteration: flat cross-category benefit summary
	body.WriteString(sectionIterationExample())

	// Worked example
	body.WriteString(sectionWorkedExample())

	// Tips
	body.WriteString(sectionTips())

	// Build DOCX
	header := buildSampleHeaderXML()
	footer := buildSampleFooterXML()

	pkg := &samplePackage{
		Body:   body.String(),
		Header: header,
		Footer: footer,
	}

	return pkg.Build()
}

type samplePackage struct {
	Body   string
	Header string
	Footer string
}

func (p *samplePackage) Build() ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	defer zw.Close()

	// Build ZIP with minimal structure
	if err := addFileToZip(zw, "[Content_Types].xml", buildSampleContentTypes()); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "_rels/.rels", buildSampleRels()); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "word/document.xml", buildSampleDocument(p.Body, p.Header, p.Footer)); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "word/_rels/document.xml.rels", buildSampleDocumentRels()); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "word/styles.xml", buildSampleStyles()); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "word/header1.xml", p.Header); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "word/_rels/header1.xml.rels", buildSampleHeaderRels()); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "word/footer1.xml", p.Footer); err != nil {
		return nil, err
	}
	if err := addFileToZip(zw, "word/_rels/footer1.xml.rels", buildSampleHeaderRels()); err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func addFileToZip(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, content)
	return err
}

// Title page — centred title, subtitle, and author hint.
func titlePage() string {
	var b strings.Builder
	b.WriteString(centeredParagraph("On Risk Letter Template Sample", true, 36))
	b.WriteString("<w:p/>")
	b.WriteString(centeredParagraph("Documentation of Available Tokens", false, 24))
	b.WriteString("<w:p/>")
	b.WriteString(centeredParagraph("Generated for Template Authors", false, 20))
	return b.String()
}

// centeredParagraph emits a properly-formed centred paragraph. Alignment
// belongs in <w:pPr>, NOT <w:rPr> — the previous helper put it in the
// wrong place AND was then double-escaped by paragraph(), which is why
// the raw XML tags showed up as literal text in Word.
func centeredParagraph(text string, bold bool, halfPoints int) string {
	var rPr strings.Builder
	rPr.WriteString("<w:rPr>")
	if bold {
		rPr.WriteString("<w:b/>")
	}
	if halfPoints > 0 {
		rPr.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/><w:szCs w:val="%d"/>`, halfPoints, halfPoints))
	}
	rPr.WriteString("</w:rPr>")
	return fmt.Sprintf(
		`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r>%s<w:t xml:space="preserve">%s</w:t></w:r></w:p>`,
		rPr.String(), xmlEscape(text),
	)
}

// How tokens work
func howTokensWork() string {
	var b strings.Builder
	b.WriteString(sectionHeading("How Tokens Work"))
	b.WriteString(paragraph("On Risk Letter templates use a template substitution engine to insert dynamic values. The following syntaxes are supported:"))
	b.WriteString(paragraph(""))

	// Simple substitution
	b.WriteString(subheading("Simple Substitution"))
	b.WriteString(bulletPara("Syntax: {{token_name}}"))
	b.WriteString(bulletPara("Example: The scheme name is {{scheme_name}}"))
	b.WriteString(bulletPara("Result: The scheme name is Acme Pension Scheme"))
	b.WriteString(paragraph(""))

	// Conditional blocks
	b.WriteString(subheading("Conditional Blocks"))
	b.WriteString(bulletPara("Syntax: {{#flag_name}}content{{/flag_name}}"))
	b.WriteString(bulletPara("Usage: Blocks are rendered only if the flag is true or non-empty"))
	b.WriteString(bulletPara("Example: {{#is_broker_channel}}Distributed via {{broker_name}}{{/is_broker_channel}}"))
	b.WriteString(paragraph(""))

	// Iteration
	b.WriteString(subheading("Iteration"))
	b.WriteString(bulletPara("Syntax: {{#array_name}}content {{field}}{{/array_name}}"))
	b.WriteString(bulletPara("Usage: Repeats content for each item in an array"))
	b.WriteString(bulletPara("Example: {{#benefit_summary}}— {{benefit}}: {{annual_premium}}{{/benefit_summary}}"))
	b.WriteString(paragraph(""))

	// Dot notation
	b.WriteString(subheading("Nested Fields"))
	b.WriteString(bulletPara("Syntax: {{object.field}}"))
	b.WriteString(bulletPara("Example: {{insurer.email}} resolves to insurer email"))
	b.WriteString(paragraph(""))

	return b.String()
}

// Section: Letter-level tokens
func sectionLetterTokens() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Letter-Level Tokens"))
	b.WriteString(paragraph("These tokens relate to the on-risk letter record itself:"))
	b.WriteString(paragraph(""))

	rows := [][]string{
		{"Token", "Description", "Example"},
		{"{{letter_reference}}", "Unique letter reference number", "ORL-7-12-1234"},
		{"{{letter_date}}", "Date the letter was generated", "14 Apr 2026"},
		{"{{generated_by}}", "Name/email of person who generated the letter", "Jome Akpoduado"},
	}
	b.WriteString(keyValueTable(rows))

	return b.String()
}

// Section: Quote-level tokens — the full shared surface from the quote
// template schema. Scalar tokens only; bool flags are documented in the
// Conditional Flags section below.
func sectionQuoteTokens() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Quote-Level Tokens"))
	b.WriteString(paragraph("These tokens resolve to fields on the associated quote. Every token a quote PDF template can reference is available here too."))
	b.WriteString(paragraph(""))

	rows := [][]string{{"Token", "Description"}}
	for _, f := range quote_template.QuoteFieldsForSample() {
		if _, isBool := f.Value.(bool); isBool {
			continue
		}
		rows = append(rows, []string{"{{" + f.Key + "}}", f.Label})
	}
	b.WriteString(keyValueTable(rows))

	return b.String()
}

// Section: Quote-context tokens unique to on-risk letters — fields sourced
// from the quote but not exposed by the quote PDF template schema.
func sectionQuoteContextExtended() string {
	var b strings.Builder
	b.WriteString(sectionHeading("On-Risk Letter Quote Context"))
	b.WriteString(paragraph("These tokens expose additional quote-adjacent fields that on-risk letters need — cover termination, scheme contact, distribution channel, and identifiers."))
	b.WriteString(paragraph(""))

	rows := [][]string{
		{"Token", "Description", "Example"},
		{"{{cover_end_date}}", "Cover end date", "30 Apr 2027"},
		{"{{scheme_contact}}", "Primary contact person for scheme", "John Doe"},
		{"{{scheme_email}}", "Contact email for scheme", "john@acme.com"},
		{"{{member_count}}", "Number of members in the quote", "120"},
		{"{{quote_id}}", "Internal quote ID", "12"},
		{"{{quote_reference}}", "Quote number / name alias", "Q-001"},
		{"{{distribution_channel}}", "Distribution channel (e.g. broker, direct)", "broker"},
		{"{{broker_name}}", "Broker name (when broker channel)", "Acme Brokers"},
	}
	b.WriteString(keyValueTable(rows))

	return b.String()
}

// Section: Insurer tokens — driven by the shared quote_template schema so
// additions (e.g. on_risk_letter_text) flow in automatically.
func sectionInsurerTokens() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Insurer Tokens"))
	b.WriteString(paragraph("Access insurer details using the insurer object with dot notation: {{insurer.field_name}}"))
	b.WriteString(paragraph(""))

	rows := [][]string{{"Token", "Description"}}
	for _, f := range quote_template.InsurerFieldsForSample() {
		rows = append(rows, []string{"{{insurer." + f.Key + "}}", f.Label})
	}
	b.WriteString(keyValueTable(rows))

	return b.String()
}

// Section: Categories iteration — identical surface to the quote PDF
// template. Templates wrap content in {{#categories}}...{{/categories}};
// the fields below are available inside the block. Benefit-prefixed
// tokens use the insurer's configured benefit code where customised.
func sectionCategoriesIteration() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Categories — Iteration and Per-Category Tokens"))
	b.WriteString(paragraph("The categories list contains one entry per scheme category in the quote. Wrap content you want repeated inside {{#categories}}...{{/categories}}. The fields below are available for each category."))
	b.WriteString(paragraph(""))

	b.WriteString(subheading("Category-level tokens"))
	scalarRows := [][]string{{"Token", "Description"}}
	for _, f := range quote_template.CategoryScalarFieldsForSample() {
		scalarRows = append(scalarRows, []string{"{{" + f.Key + "}}", f.Label})
	}
	b.WriteString(keyValueTable(scalarRows))
	b.WriteString(paragraph(""))

	b.WriteString(subheading("Category-level flags (true/false)"))
	b.WriteString(paragraph("Use these as conditional blocks to include content only when the category has that benefit."))
	for _, f := range quote_template.CategoryBoolFieldsForSample() {
		b.WriteString(bulletPara("{{#" + f.Key + "}}...{{/" + f.Key + "}} — " + f.Label))
	}
	b.WriteString(paragraph(""))

	b.WriteString(subheading("Example: one paragraph per category"))
	b.WriteString(paragraph("{{#categories}}— {{name}}: {{member_count}} lives · premium {{premium}} · {{percent_salary}} of salary{{/categories}}"))

	return b.String()
}

// Section: Nested benefits — each category exposes an object per benefit
// type. The sub-object is populated only when the category has that
// benefit. Prefixes come from the resolved benefit naming so customised
// codes (e.g. a GLA renamed to "group_life") are shown in the sample.
func sectionNestedBenefits() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Benefit Tokens — Inside the Categories Block"))
	b.WriteString(paragraph("Each category exposes an object per benefit (e.g. gla, sgla, ptd, ci, phi, ttd, funeral). Access nested fields with dot notation. Combine with the conditional flags above to render a section only when relevant."))
	b.WriteString(paragraph(""))

	for _, spec := range quote_template.BenefitSpecsForSample() {
		b.WriteString(subheading(spec.Title))
		rows := [][]string{{"Token", "Description"}}
		for _, f := range spec.Fields() {
			rows = append(rows, []string{"{{" + spec.Prefix + "." + f.Key + "}}", f.Label})
		}
		b.WriteString(keyValueTable(rows))
		b.WriteString(paragraph(""))
	}

	return b.String()
}

// Section: Conditional flags
func sectionConditionalFlags() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Conditional Flags"))
	b.WriteString(paragraph("Use these boolean flags with {{#flag}}...{{/flag}} syntax:"))
	b.WriteString(paragraph(""))

	rows := [][]string{
		{"Flag", "True When", "Example Usage"},
		{"{{#is_broker_channel}}", "Distribution channel is 'broker'", "{{#is_broker_channel}}This quote was distributed via {{broker_name}}{{/is_broker_channel}}"},
		{"{{#has_benefits}}", "Benefit summary is not empty", "{{#has_benefits}}See benefit breakdown below{{/has_benefits}}"},
		{"{{#has_non_funeral_benefits}}", "Any category has a non-funeral benefit", "{{#has_non_funeral_benefits}}Risk benefits are in force{{/has_non_funeral_benefits}}"},
		{"{{#use_global_salary_multiple}}", "Scheme uses a single global salary multiple", "{{#use_global_salary_multiple}}Uniform salary multiple applies{{/use_global_salary_multiple}}"},
	}
	b.WriteString(keyValueTable(rows))

	return b.String()
}

// Section: Iteration example
func sectionIterationExample() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Iteration: Benefit Summary"))
	b.WriteString(paragraph("The benefit_summary array contains one entry per benefit type. Each entry has:"))
	b.WriteString(paragraph(""))
	b.WriteString(bulletPara("benefit: The benefit name (e.g., 'Group Life Assurance (GLA)')"))
	b.WriteString(bulletPara("annual_premium: The annual premium for that benefit"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("Example template code:"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("{{#benefit_summary}}"))
	b.WriteString(paragraph("  — {{benefit}}: {{annual_premium}}"))
	b.WriteString(paragraph("{{/benefit_summary}}"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("Output:"))
	b.WriteString(bulletPara("  — Group Life Assurance (GLA): 120 000.00"))
	b.WriteString(bulletPara("  — Permanent Total Disability (PTD): 80 000.00"))
	b.WriteString(bulletPara("  — Critical Illness (CI): 115 000.00"))

	return b.String()
}

// Section: Worked example
func sectionWorkedExample() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Worked Example: Conditional + Iteration"))
	b.WriteString(paragraph("This example combines conditionals and iteration:"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("Template:"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("Dear {{scheme_contact}},"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("Your quotation {{quote_reference}} for {{scheme_name}} is now on risk as of {{commencement_date}}."))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("{{#has_benefits}}"))
	b.WriteString(paragraph("The agreed benefits are:"))
	b.WriteString(paragraph("{{#benefit_summary}}"))
	b.WriteString(paragraph("  • {{benefit}}: {{annual_premium}}"))
	b.WriteString(paragraph("{{/benefit_summary}}"))
	b.WriteString(paragraph("{{/has_benefits}}"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("{{#is_broker_channel}}"))
	b.WriteString(paragraph("This policy is placed through {{broker_name}}."))
	b.WriteString(paragraph("{{/is_broker_channel}}"))
	b.WriteString(paragraph(""))
	b.WriteString(paragraph("Best regards,"))
	b.WriteString(paragraph("{{insurer.name}}"))

	return b.String()
}

// Section: Tips
func sectionTips() string {
	var b strings.Builder
	b.WriteString(sectionHeading("Template Authoring Tips"))
	b.WriteString(paragraph(""))

	b.WriteString(subheading("1. Token Placement"))
	b.WriteString(bulletPara("Tokens MUST be placed inside <w:t> elements (Word text runs). Do not place them directly in paragraphs or table cells."))
	b.WriteString(bulletPara("If a token appears to be split across runs, the system will automatically stitch them together."))
	b.WriteString(paragraph(""))

	b.WriteString(subheading("2. Formatting"))
	b.WriteString(bulletPara("You can apply formatting (bold, italic, color) to the run containing a token; the token will inherit that formatting."))
	b.WriteString(bulletPara("Currency amounts are pre-formatted with space-separated thousands (e.g., '500 000.00')."))
	b.WriteString(paragraph(""))

	b.WriteString(subheading("3. Missing Tokens"))
	b.WriteString(bulletPara("If a token is referenced but missing from the context, it will be replaced with an empty string."))
	b.WriteString(bulletPara("Use conditional blocks to hide sections when data might be absent: {{#broker_name}}...{{/broker_name}}"))
	b.WriteString(paragraph(""))

	b.WriteString(subheading("4. Testing Your Template"))
	b.WriteString(bulletPara("Download the sample template from the admin UI."))
	b.WriteString(bulletPara("Edit it in Microsoft Word to add your insurer-specific text and formatting."))
	b.WriteString(bulletPara("Keep all tokens unchanged (do not retype them; copy-paste if needed)."))
	b.WriteString(bulletPara("Upload and test with a real quote in draft status."))
	b.WriteString(paragraph(""))

	b.WriteString(subheading("5. Common Patterns"))
	b.WriteString(bulletPara("Date fields are pre-formatted ('14 Apr 2026') — no additional formatting needed."))
	b.WriteString(bulletPara("Boolean flags like {{is_broker_channel}} are true/false; always use with {{#flag}}...{{/flag}} syntax."))
	b.WriteString(bulletPara("Nested objects like {{insurer.name}} use dot notation; never use curly braces inside a token."))
	b.WriteString(paragraph(""))

	return b.String()
}

// Helper: Create a paragraph containing the given text. The text is escaped
// and wrapped in <w:r><w:t> as required by OOXML — passing raw text directly
// inside <w:p> (without a run) makes the document invalid and Word refuses
// to open it.
func paragraph(content string) string {
	if content == "" {
		return "<w:p/>"
	}
	return fmt.Sprintf(`<w:p><w:r><w:t xml:space="preserve">%s</w:t></w:r></w:p>`, xmlEscape(content))
}

// Helper: Create a section heading
func sectionHeading(text string) string {
	escaped := xmlEscape(text)
	return fmt.Sprintf(`<w:p><w:pPr><w:spacing w:before="200" w:after="100"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="32"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, escaped)
}

// Helper: Create a subheading
func subheading(text string) string {
	escaped := xmlEscape(text)
	return fmt.Sprintf(`<w:p><w:pPr><w:spacing w:before="100" w:after="50"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="28"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, escaped)
}

// Helper: Create a bullet paragraph
func bulletPara(text string) string {
	escaped := xmlEscape(text)
	return fmt.Sprintf(`<w:p><w:pPr><w:pStyle w:val="ListBullet"/><w:numPr><w:ilvl w:val="0"/><w:numId w:val="1"/></w:numPr></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>`, escaped)
}

// Helper: Create a key-value table from rows
func keyValueTable(rows [][]string) string {
	var b strings.Builder
	b.WriteString(`<w:tbl><w:tblPr><w:tblW w:w="9500" w:type="dxa"/><w:tblBorders><w:top w:val="single" w:sz="8" w:color="000000"/><w:left w:val="single" w:sz="8" w:color="000000"/><w:bottom w:val="single" w:sz="8" w:color="000000"/><w:right w:val="single" w:sz="8" w:color="000000"/><w:insideH w:val="single" w:sz="8" w:color="000000"/><w:insideV w:val="single" w:sz="8" w:color="000000"/></w:tblBorders><w:tblCellMar><w:top w:w="100" w:type="dxa"/><w:left w:w="100" w:type="dxa"/><w:bottom w:w="100" w:type="dxa"/><w:right w:w="100" w:type="dxa"/></w:tblCellMar></w:tblPr>`)

	for i, row := range rows {
		b.WriteString(`<w:tr>`)
		for _, cell := range row {
			cellContent := xmlEscape(cell)
			if i == 0 {
				// Header row
				b.WriteString(fmt.Sprintf(`<w:tc><w:tcPr><w:shd w:val="clear" w:color="auto" w:fill="D3D3D3"/></w:tcPr><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>%s</w:t></w:r></w:p></w:tc>`, cellContent))
			} else {
				b.WriteString(fmt.Sprintf(`<w:tc><w:p><w:r><w:t>%s</w:t></w:r></w:p></w:tc>`, cellContent))
			}
		}
		b.WriteString(`</w:tr>`)
	}

	b.WriteString(`</w:tbl>`)
	return b.String()
}

// Helper: Escape XML special characters
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// DOCX structure builders
func buildSampleContentTypes() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/><Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/><Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/><Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/></Types>`
}

func buildSampleRels() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`
}

func buildSampleDocument(body, header, footer string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><w:body>%s<w:sectPr><w:pgSz w:w="12240" w:h="15840"/><w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440"/></w:sectPr></w:body></w:document>`, body)
}

func buildSampleDocumentRels() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/><Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/></Relationships>`
}

func buildSampleStyles() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:docDefaults><w:rPrDefault><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/><w:sz w:val="20"/></w:rPr></w:rPrDefault><w:pPrDefault/></w:docDefaults><w:style w:type="paragraph" w:default="1" w:styleId="Normal"><w:name w:val="Normal"/><w:qFormat/></w:style><w:style w:type="paragraph" w:styleId="ListBullet"><w:name w:val="List Bullet"/><w:basedOn w:val="Normal"/><w:pPr><w:numPr><w:ilvl w:val="0"/><w:numId w:val="1"/></w:numPr></w:pPr></w:style></w:styles>`
}

func buildSampleHeaderRels() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>`
}

func buildSampleHeaderXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:p><w:pPr><w:spacing w:after="100"/></w:pPr><w:r><w:t>On Risk Letter Template Sample</w:t></w:r></w:p></w:hdr>`
}

func buildSampleFooterXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:t>Page </w:t></w:r><w:fldSimple w:instr=" PAGE "><w:r><w:t>1</w:t></w:r></w:fldSimple></w:p></w:ftr>`
}

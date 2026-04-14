package quote_template

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
)

// BuildSampleTemplate produces a comprehensive .docx template that exercises
// every token the render engine supports. Admins download this as a
// self-documenting reference: each section explains a token category and
// demonstrates real usage.
//
// Every element here must be valid OOXML or Word will refuse to open the file.
// Rules worth remembering:
//   - Tokens must live inside a <w:t> element (inside a <w:r> inside a <w:p>).
//     They CANNOT appear as text directly inside <w:tr>, <w:tbl>, or <w:body>.
//   - <w:rPr> child order: rFonts, b, i, color, sz, szCs.
//   - <w:pPr> child order: pBdr, shd, spacing, jc.
//   - <w:tcPr> child order: tcW, tcBorders (wrapped), shd, tcMar, vAlign.
//   - Every <w:tc> must contain at least one <w:p>.
//   - <w:body> must end with a <w:sectPr>.
//   - <w:shd> requires w:val and w:color attributes, not just w:fill.
func BuildSampleTemplate() ([]byte, error) {
	body := buildSampleBodyXML()

	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	add := func(name, content string) error {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, content)
		return err
	}

	if err := add("[Content_Types].xml", sampleContentTypes()); err != nil {
		return nil, err
	}
	if err := add("_rels/.rels", sampleRootRels()); err != nil {
		return nil, err
	}
	if err := add("word/document.xml", sampleDocumentXML(body)); err != nil {
		return nil, err
	}
	if err := add("word/_rels/document.xml.rels", sampleDocumentRels()); err != nil {
		return nil, err
	}
	if err := add("word/styles.xml", sampleStylesXML()); err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// buildSampleBodyXML composes the template body. Organized as a reference
// document with one section per token category.
func buildSampleBodyXML() string {
	var b strings.Builder

	// ===== Cover / title =====
	b.WriteString(centeredTitle("Quote Template — Token Reference"))
	b.WriteString(centeredSubtitle("Every placeholder below will be replaced with quote data at generation time"))
	b.WriteString(bodyPara("This document is a starting point. Keep the tokens you need, delete the rest, and restyle freely. Tokens appear in double curly braces. Conditional and iteration blocks use a # to open and a / to close."))

	// ===== How tokens work =====
	b.WriteString(heading("How Tokens Work"))
	b.WriteString(bodyPara("There are three kinds of tokens:"))
	b.WriteString(bulletPara("Simple value — replaced with the corresponding data field. Example: {{scheme_name}}"))
	b.WriteString(bulletPara("Conditional block — the content between the open and close tags appears only when the condition is true. Example: {{#has_gla}}…{{/has_gla}}"))
	b.WriteString(bulletPara("Iteration block — the content is repeated once per item in a list. Example: {{#categories}}…{{/categories}}"))
	b.WriteString(bodyPara("Inside an iteration block, the per-item fields are available without any prefix (e.g. {{name}} refers to the current category's name)."))

	// ===== Quote-level tokens =====
	b.WriteString(heading("Quote-level Tokens"))
	b.WriteString(bodyPara("These tokens resolve to values at the quote level. They can be used anywhere in the document."))
	b.WriteString(keyValueTable([][2]string{
		{"Quote Number", "{{quote_number}}"},
		{"Quote Name (alias)", "{{quote_name}}"},
		{"Quote Type", "{{quote_type}}"},
		{"Obligation Type", "{{obligation_type}}"},
		{"Scheme Name", "{{scheme_name}}"},
		{"Creation Date", "{{creation_date}}"},
		{"Commencement Date", "{{commencement_date}}"},
		{"Industry", "{{industry}}"},
		{"Currency", "{{currency}}"},
		{"Free Cover Limit", "{{free_cover_limit}}"},
		{"Normal Retirement Age", "{{normal_retirement_age}}"},
		{"Total Lives Covered", "{{total_lives}}"},
		{"Total Sum Assured", "{{total_sum_assured}}"},
		{"Total Annual Salary", "{{total_annual_salary}}"},
		{"Total Annual Premium", "{{total_annual_premium}}"},
	}))

	// ===== Insurer tokens =====
	b.WriteString(heading("Insurer Tokens"))
	b.WriteString(bodyPara("The insurer object carries all the insurer details configured in the system. Access fields with dot notation."))
	b.WriteString(keyValueTable([][2]string{
		{"Insurer Name", "{{insurer.name}}"},
		{"Contact Person", "{{insurer.contact_person}}"},
		{"Address Line 1", "{{insurer.address_line_1}}"},
		{"Address Line 2", "{{insurer.address_line_2}}"},
		{"Address Line 3", "{{insurer.address_line_3}}"},
		{"City", "{{insurer.city}}"},
		{"Province", "{{insurer.province}}"},
		{"Post Code", "{{insurer.post_code}}"},
		{"Country", "{{insurer.country}}"},
		{"Telephone", "{{insurer.telephone}}"},
		{"Email", "{{insurer.email}}"},
		{"Introductory Text", "{{insurer.introductory_text}}"},
		{"General Provisions Text", "{{insurer.general_provisions_text}}"},
	}))

	// ===== Top-level conditional =====
	b.WriteString(heading("Top-level Conditional"))
	b.WriteString(bodyPara("The has_non_funeral_benefits flag is true when any category has a non-funeral benefit (GLA, SGLA, PTD, CI, PHI, or TTD)."))
	b.WriteString(bodyPara("{{#has_non_funeral_benefits}}This sentence appears only when the quote has non-funeral benefits.{{/has_non_funeral_benefits}}"))
	b.WriteString(bodyPara("You can use the same block to hide content in the opposite case, by placing the text outside the block."))

	// ===== Categories iteration =====
	b.WriteString(heading("Categories — Iteration and Per-Category Tokens"))
	b.WriteString(bodyPara("The categories list contains one entry per scheme category in the quote. Wrap content you want repeated inside {{#categories}}…{{/categories}}. The fields below are available for each category."))

	b.WriteString(subheading("Category-level tokens"))
	b.WriteString(keyValueTable([][2]string{
		{"Category Name", "{{name}}"},
		{"Region", "{{region}}"},
		{"Member Count", "{{member_count}}"},
		{"Free Cover Limit (category override)", "{{free_cover_limit}}"},
		{"Total Annual Salary", "{{total_salary}}"},
		{"Total Sum Assured", "{{total_sum_assured}}"},
		{"Annual Premium (excl. funeral)", "{{annual_premium}}"},
		{"Premium as % of Salary", "{{percent_salary}}"},
		{"Retirement Premium Waiver", "{{retirement_premium_waiver}}"},
		{"Medical Aid Premium Waiver", "{{medical_aid_premium_waiver}}"},
	}))

	b.WriteString(subheading("Category-level flags (true/false)"))
	b.WriteString(bodyPara("Use these as conditional blocks to include content only when the category has that benefit."))
	b.WriteString(bulletPara("{{#has_non_funeral_benefits}}…{{/has_non_funeral_benefits}}"))
	b.WriteString(bulletPara("{{#has_gla}}…{{/has_gla}}, {{#has_sgla}}…{{/has_sgla}}, {{#has_ptd}}…{{/has_ptd}}"))
	b.WriteString(bulletPara("{{#has_ci}}…{{/has_ci}}, {{#has_phi}}…{{/has_phi}}, {{#has_ttd}}…{{/has_ttd}}, {{#has_funeral}}…{{/has_funeral}}"))

	b.WriteString(subheading("Example: one paragraph per category"))
	b.WriteString(bodyPara("{{#categories}}— {{name}}: {{member_count}} lives · annual premium {{annual_premium}} · {{percent_salary}} of salary{{/categories}}"))

	// ===== Per-benefit tokens (inside categories) =====
	b.WriteString(heading("Benefit Tokens — Used Inside the Categories Block"))
	b.WriteString(bodyPara("Each category exposes an object per benefit. These are populated only when the category has the corresponding benefit. Combine with the conditional flags above to show a block only when relevant."))

	b.WriteString(subheading("Group Life Assurance (GLA)"))
	b.WriteString(benefitTokenTable([][2]string{
		{"Benefit Title", "{{gla.title}}"},
		{"Salary Multiple", "{{gla.salary_multiple}}"},
		{"Waiting Period (months)", "{{gla.waiting_period}}"},
		{"Benefit Structure", "{{gla.benefit_structure}}"},
		{"Benefit Type", "{{gla.benefit_type}}"},
		{"Terminal Illness Benefit", "{{gla.terminal_illness_benefit}}"},
		{"Educator Benefit", "{{gla.educator_benefit}}"},
		{"Total Sum Assured", "{{gla.total_sum_assured}}"},
		{"Annual Premium", "{{gla.annual_premium}}"},
		{"% of Salary", "{{gla.percent_salary}}"},
	}))

	b.WriteString(subheading("Spouse Group Life (SGLA)"))
	b.WriteString(benefitTokenTable([][2]string{
		{"Benefit Title", "{{sgla.title}}"},
		{"Salary Multiple", "{{sgla.salary_multiple}}"},
		{"Waiting Period (months)", "{{sgla.waiting_period}}"},
		{"Benefit Structure", "{{sgla.benefit_structure}}"},
		{"Maximum Benefit", "{{sgla.max_benefit}}"},
		{"Total Sum Assured", "{{sgla.total_sum_assured}}"},
		{"Annual Premium", "{{sgla.annual_premium}}"},
		{"% of Salary", "{{sgla.percent_salary}}"},
	}))

	b.WriteString(subheading("Permanent Total Disability (PTD)"))
	b.WriteString(benefitTokenTable([][2]string{
		{"Benefit Title", "{{ptd.title}}"},
		{"Salary Multiple", "{{ptd.salary_multiple}}"},
		{"Waiting Period (months)", "{{ptd.waiting_period}}"},
		{"Deferred Period (months)", "{{ptd.deferred_period}}"},
		{"Benefit Type", "{{ptd.benefit_type}}"},
		{"Disability Definition", "{{ptd.disability_definition}}"},
		{"Risk Type", "{{ptd.risk_type}}"},
		{"Educator Benefit", "{{ptd.educator_benefit}}"},
		{"Total Sum Assured", "{{ptd.total_sum_assured}}"},
		{"Annual Premium", "{{ptd.annual_premium}}"},
		{"% of Salary", "{{ptd.percent_salary}}"},
	}))

	b.WriteString(subheading("Critical Illness (CI)"))
	b.WriteString(benefitTokenTable([][2]string{
		{"Benefit Title", "{{ci.title}}"},
		{"Salary Multiple", "{{ci.salary_multiple}}"},
		{"Waiting Period (months)", "{{ci.waiting_period}}"},
		{"Deferred Period (months)", "{{ci.deferred_period}}"},
		{"Benefit Structure", "{{ci.benefit_structure}}"},
		{"Benefit Definition", "{{ci.benefit_definition}}"},
		{"Maximum Benefit", "{{ci.max_benefit}}"},
		{"Total Sum Assured", "{{ci.total_sum_assured}}"},
		{"Annual Premium", "{{ci.annual_premium}}"},
		{"% of Salary", "{{ci.percent_salary}}"},
	}))

	b.WriteString(subheading("Permanent Health Insurance (PHI)"))
	b.WriteString(benefitTokenTable([][2]string{
		{"Benefit Title", "{{phi.title}}"},
		{"Income Replacement %", "{{phi.income_replacement_percentage}}"},
		{"Waiting Period (months)", "{{phi.waiting_period}}"},
		{"Deferred Period (months)", "{{phi.deferred_period}}"},
		{"Disability Definition", "{{phi.disability_definition}}"},
		{"Risk Type", "{{phi.risk_type}}"},
		{"Premium Waiver", "{{phi.premium_waiver}}"},
		{"Medical Aid Premium Waiver", "{{phi.medical_aid_premium_waiver}}"},
		{"Benefit Escalation", "{{phi.benefit_escalation}}"},
		{"Total Covered Income", "{{phi.total_sum_assured}}"},
		{"Annual Premium", "{{phi.annual_premium}}"},
		{"% of Salary", "{{phi.percent_salary}}"},
	}))

	b.WriteString(subheading("Temporary Total Disability (TTD)"))
	b.WriteString(benefitTokenTable([][2]string{
		{"Benefit Title", "{{ttd.title}}"},
		{"Income Replacement %", "{{ttd.income_replacement_percentage}}"},
		{"Waiting Period (months)", "{{ttd.waiting_period}}"},
		{"Deferred Period (months)", "{{ttd.deferred_period}}"},
		{"Disability Definition", "{{ttd.disability_definition}}"},
		{"Risk Type", "{{ttd.risk_type}}"},
		{"Total Covered Income", "{{ttd.total_sum_assured}}"},
		{"Annual Premium", "{{ttd.annual_premium}}"},
		{"% of Salary", "{{ttd.percent_salary}}"},
	}))

	b.WriteString(subheading("Group Funeral"))
	b.WriteString(benefitTokenTable([][2]string{
		{"Benefit Title", "{{funeral.title}}"},
		{"Monthly Premium per Member", "{{funeral.monthly_premium_per_member}}"},
		{"Annual Premium per Member", "{{funeral.annual_premium_per_member}}"},
		{"Total Annual Premium", "{{funeral.total_annual_premium}}"},
		{"Main Member Sum Assured", "{{funeral.main_member_sum_assured}}"},
		{"Spouse Sum Assured", "{{funeral.spouse_sum_assured}}"},
		{"Child Sum Assured", "{{funeral.child_sum_assured}}"},
		{"Max Children Covered", "{{funeral.max_children}}"},
		{"Parent Sum Assured", "{{funeral.parent_sum_assured}}"},
		{"Dependant Sum Assured", "{{funeral.dependant_sum_assured}}"},
		{"Max Dependants Covered", "{{funeral.max_dependants}}"},
	}))

	// ===== Worked example =====
	b.WriteString(heading("Worked Example — Combining Everything"))
	b.WriteString(bodyPara("Here is a realistic block that loops over categories and shows per-benefit sections only where relevant:"))
	b.WriteString(exampleBlock())

	// ===== Footer note =====
	b.WriteString(heading("Tips for Template Authors"))
	b.WriteString(bulletPara("Edit this document in Word as you would any .docx — apply styles, change fonts, add your logo."))
	b.WriteString(bulletPara("Keep the curly-brace tokens intact. If a token breaks across two words (e.g. from autocorrect), the system will automatically rejoin it at render time."))
	b.WriteString(bulletPara("To test, upload your template from the insurer settings screen and generate a quote document — you'll see the replacements immediately."))
	b.WriteString(bulletPara("Missing tokens resolve to empty strings, so a template that references a token not populated for a given quote simply shows nothing rather than breaking."))

	return b.String()
}

// --- Paragraph builders ---

func centeredTitle(text string) string {
	return `<w:p><w:pPr><w:spacing w:after="120"/><w:jc w:val="center"/></w:pPr>` +
		`<w:r><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/><w:b/><w:color w:val="1E3A5F"/><w:sz w:val="44"/><w:szCs w:val="44"/></w:rPr>` +
		`<w:t xml:space="preserve">` + escape(text) + `</w:t></w:r></w:p>`
}

func centeredSubtitle(text string) string {
	return `<w:p><w:pPr><w:spacing w:after="240"/><w:jc w:val="center"/></w:pPr>` +
		`<w:r><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/><w:i/><w:color w:val="586069"/><w:sz w:val="22"/><w:szCs w:val="22"/></w:rPr>` +
		`<w:t xml:space="preserve">` + escape(text) + `</w:t></w:r></w:p>`
}

func heading(text string) string {
	return `<w:p><w:pPr><w:pBdr><w:bottom w:val="single" w:sz="8" w:space="2" w:color="1E3A5F"/></w:pBdr><w:spacing w:before="320" w:after="120"/></w:pPr>` +
		`<w:r><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/><w:b/><w:color w:val="1E3A5F"/><w:sz w:val="30"/><w:szCs w:val="30"/></w:rPr>` +
		`<w:t xml:space="preserve">` + escape(text) + `</w:t></w:r></w:p>`
}

func subheading(text string) string {
	return `<w:p><w:pPr><w:spacing w:before="200" w:after="80"/></w:pPr>` +
		`<w:r><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/><w:b/><w:color w:val="2C3E50"/><w:sz w:val="24"/><w:szCs w:val="24"/></w:rPr>` +
		`<w:t xml:space="preserve">` + escape(text) + `</w:t></w:r></w:p>`
}

func bodyPara(text string) string {
	return `<w:p><w:pPr><w:spacing w:after="120"/></w:pPr>` +
		`<w:r><w:t xml:space="preserve">` + escape(text) + `</w:t></w:r></w:p>`
}

func bulletPara(text string) string {
	// Simple hanging-indent pseudo-bullet (no numbering.xml dependency) —
	// prefixes a bullet glyph and indents the paragraph.
	return `<w:p><w:pPr><w:spacing w:after="80"/><w:ind w:left="360" w:hanging="360"/></w:pPr>` +
		`<w:r><w:t xml:space="preserve">• ` + escape(text) + `</w:t></w:r></w:p>`
}

// --- Table builders ---

// keyValueTable emits a 2-column reference table used for the token-list sections.
func keyValueTable(rows [][2]string) string {
	var b strings.Builder
	b.WriteString(`<w:tbl>`)
	b.WriteString(`<w:tblPr><w:tblW w:w="9000" w:type="dxa"/><w:tblBorders>`)
	for _, e := range []string{"top", "left", "bottom", "right", "insideH", "insideV"} {
		b.WriteString(`<w:` + e + ` w:val="single" w:sz="4" w:space="0" w:color="BFBFBF"/>`)
	}
	b.WriteString(`</w:tblBorders></w:tblPr>`)
	b.WriteString(`<w:tblGrid><w:gridCol w:w="3600"/><w:gridCol w:w="5400"/></w:tblGrid>`)

	for _, r := range rows {
		b.WriteString(`<w:tr>`)
		b.WriteString(cellWithShading(3600, "ECF0F1", r[0], true))
		b.WriteString(cellMonospace(5400, r[1]))
		b.WriteString(`</w:tr>`)
	}
	b.WriteString(`</w:tbl><w:p/>`)
	return b.String()
}

// benefitTokenTable is the same as keyValueTable but with a distinct accent
// colour, to visually separate per-benefit sections.
func benefitTokenTable(rows [][2]string) string {
	var b strings.Builder
	b.WriteString(`<w:tbl>`)
	b.WriteString(`<w:tblPr><w:tblW w:w="9000" w:type="dxa"/><w:tblBorders>`)
	for _, e := range []string{"top", "left", "bottom", "right", "insideH", "insideV"} {
		b.WriteString(`<w:` + e + ` w:val="single" w:sz="4" w:space="0" w:color="D1E4F0"/>`)
	}
	b.WriteString(`</w:tblBorders></w:tblPr>`)
	b.WriteString(`<w:tblGrid><w:gridCol w:w="3600"/><w:gridCol w:w="5400"/></w:tblGrid>`)

	for _, r := range rows {
		b.WriteString(`<w:tr>`)
		b.WriteString(cellWithShading(3600, "F1F8FF", r[0], true))
		b.WriteString(cellMonospace(5400, r[1]))
		b.WriteString(`</w:tr>`)
	}
	b.WriteString(`</w:tbl><w:p/>`)
	return b.String()
}

func cellWithShading(width int, fill, text string, bold bool) string {
	var b strings.Builder
	b.WriteString(`<w:tc><w:tcPr>`)
	b.WriteString(widthDxa(width))
	b.WriteString(tcBordersThin())
	b.WriteString(`<w:shd w:val="clear" w:color="auto" w:fill="` + fill + `"/>`)
	b.WriteString(`</w:tcPr>`)
	b.WriteString(`<w:p><w:r><w:rPr>`)
	if bold {
		b.WriteString(`<w:b/>`)
	}
	b.WriteString(`</w:rPr><w:t xml:space="preserve">` + escape(text) + `</w:t></w:r></w:p>`)
	b.WriteString(`</w:tc>`)
	return b.String()
}

func cellMonospace(width int, text string) string {
	var b strings.Builder
	b.WriteString(`<w:tc><w:tcPr>`)
	b.WriteString(widthDxa(width))
	b.WriteString(tcBordersThin())
	b.WriteString(`</w:tcPr>`)
	b.WriteString(`<w:p><w:r><w:rPr>`)
	b.WriteString(`<w:rFonts w:ascii="Consolas" w:hAnsi="Consolas"/>`)
	b.WriteString(`<w:color w:val="1E3A5F"/>`)
	b.WriteString(`</w:rPr><w:t xml:space="preserve">` + escape(text) + `</w:t></w:r></w:p>`)
	b.WriteString(`</w:tc>`)
	return b.String()
}

func widthDxa(n int) string {
	return `<w:tcW w:w="` + itoa(n) + `" w:type="dxa"/>`
}

func tcBordersThin() string {
	var b strings.Builder
	b.WriteString(`<w:tcBorders>`)
	for _, e := range []string{"top", "left", "bottom", "right"} {
		b.WriteString(`<w:` + e + ` w:val="single" w:sz="4" w:space="0" w:color="BFBFBF"/>`)
	}
	b.WriteString(`</w:tcBorders>`)
	return b.String()
}

// exampleBlock produces a short worked example combining iteration,
// conditionals, and simple tokens — rendered in a slightly accented block.
func exampleBlock() string {
	var b strings.Builder
	b.WriteString(bodyPara("For the scheme {{scheme_name}} (quote {{quote_number}}), we cover {{total_lives}} lives with a total annual premium of {{total_annual_premium}}."))
	b.WriteString(bodyPara("{{#categories}}"))
	b.WriteString(bodyPara("Category: {{name}} — {{member_count}} members, annual premium {{annual_premium}}."))
	b.WriteString(bodyPara("{{#has_gla}}• Group Life sum assured: {{gla.total_sum_assured}} · premium {{gla.annual_premium}}{{/has_gla}}"))
	b.WriteString(bodyPara("{{#has_funeral}}• Funeral cover: {{funeral.total_annual_premium}} annually across the category{{/has_funeral}}"))
	b.WriteString(bodyPara("{{/categories}}"))
	b.WriteString(bodyPara("{{#has_non_funeral_benefits}}Non-funeral benefits apply to this quote as detailed above.{{/has_non_funeral_benefits}}"))
	return b.String()
}

// --- Helpers ---

func escape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func itoa(n int) string {
	// Inline strconv.Itoa to avoid importing strconv for a single usage;
	// matches the style used by the other files in this package.
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// --- Package parts ---

func sampleDocumentXML(body string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
		`<w:body>` + body +
		`<w:sectPr>` +
		`<w:pgSz w:w="11906" w:h="16838"/>` +
		`<w:pgMar w:top="1134" w:right="1134" w:bottom="1134" w:left="1134" w:header="720" w:footer="720" w:gutter="0"/>` +
		`<w:cols w:num="1"/>` +
		`<w:docGrid w:linePitch="360"/>` +
		`</w:sectPr>` +
		`</w:body>` +
		`</w:document>`
}

func sampleContentTypes() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
		`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
		`<Default Extension="xml" ContentType="application/xml"/>` +
		`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>` +
		`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>` +
		`</Types>`
}

func sampleRootRels() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
		`</Relationships>`
}

func sampleDocumentRels() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>` +
		`</Relationships>`
}

func sampleStylesXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:docDefaults>` +
		`<w:rPrDefault><w:rPr><w:rFonts w:ascii="Calibri" w:hAnsi="Calibri"/><w:sz w:val="22"/></w:rPr></w:rPrDefault>` +
		`<w:pPrDefault/>` +
		`</w:docDefaults>` +
		`<w:style w:type="paragraph" w:default="1" w:styleId="Normal"><w:name w:val="Normal"/><w:qFormat/></w:style>` +
		`</w:styles>`
}

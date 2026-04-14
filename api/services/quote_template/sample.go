package quote_template

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
)

// BuildSampleTemplate produces a small but valid .docx template demonstrating
// every supported token type. Admins download this from the admin UI as a
// starting point.
//
// Important constraints — every element here must be valid OOXML or Word will
// refuse to open the file:
//   - Tokens must live inside a <w:t> element (i.e. inside a <w:r> inside a <w:p>).
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

// --- Body ---

func buildSampleBodyXML() string {
	var b strings.Builder

	// Centred title
	b.WriteString(`<w:p><w:pPr><w:spacing w:after="200"/><w:jc w:val="center"/></w:pPr>`)
	b.WriteString(`<w:r><w:rPr><w:b/><w:sz w:val="40"/><w:szCs w:val="40"/></w:rPr><w:t xml:space="preserve">Quote Template Sample</w:t></w:r>`)
	b.WriteString(`</w:p>`)

	// Intro paragraph
	b.WriteString(`<w:p><w:pPr><w:spacing w:after="200"/></w:pPr>`)
	b.WriteString(`<w:r><w:t xml:space="preserve">This is a starter Word template for generating quote documents. Edit the layout, branding, and wording to match your insurer&apos;s style. The placeholders below in {{double curly braces}} will be replaced with the real quote data when a document is generated.</w:t></w:r>`)
	b.WriteString(`</w:p>`)

	// Insurer details
	b.WriteString(sectionHeading("Insurer Details"))
	b.WriteString(simplePara(`Insurer: {{insurer.name}}`))
	b.WriteString(simplePara(`Address: {{insurer.address_line_1}}, {{insurer.city}}, {{insurer.post_code}}`))
	b.WriteString(simplePara(`Telephone: {{insurer.telephone}}`))
	b.WriteString(simplePara(`Email: {{insurer.email}}`))

	// Quote summary table
	b.WriteString(sectionHeading("Quote Information"))
	b.WriteString(quoteInfoTable())

	// Conditional block (paragraph-level)
	b.WriteString(sectionHeading("Conditional Block Example"))
	b.WriteString(simplePara(`The following sentence appears only when the quote includes non-funeral benefits:`))
	b.WriteString(simplePara(`{{#has_non_funeral_benefits}}This quote includes non-funeral benefits as listed in the breakdown above.{{/has_non_funeral_benefits}}`))

	// Iteration block (paragraph-level — repeated once per scheme category)
	b.WriteString(sectionHeading("Per-Category Iteration"))
	b.WriteString(simplePara(`The block between {{#categories}} and {{/categories}} repeats once per category in the quote:`))
	b.WriteString(simplePara(`{{#categories}}— {{name}}: {{member_count}} lives, annual premium {{annual_premium}} ({{percent_salary}} of salary){{/categories}}`))

	// Closing
	b.WriteString(sectionHeading("Underwriting and General Provisions"))
	b.WriteString(simplePara(`{{insurer.general_provisions_text}}`))

	return b.String()
}

func sectionHeading(text string) string {
	return `<w:p><w:pPr><w:spacing w:before="240" w:after="120"/></w:pPr>` +
		`<w:r><w:rPr><w:b/><w:color w:val="2C3E50"/><w:sz w:val="28"/><w:szCs w:val="28"/></w:rPr>` +
		`<w:t xml:space="preserve">` + text + `</w:t></w:r></w:p>`
}

func simplePara(text string) string {
	return `<w:p><w:pPr><w:spacing w:after="100"/></w:pPr>` +
		`<w:r><w:t xml:space="preserve">` + text + `</w:t></w:r></w:p>`
}

// quoteInfoTable produces a 2-column key/value table with quote summary tokens.
// All tokens live INSIDE <w:t> within properly-nested cells.
func quoteInfoTable() string {
	rows := [][2]string{
		{"Quote Number", "{{quote_number}}"},
		{"Scheme Name", "{{scheme_name}}"},
		{"Creation Date", "{{creation_date}}"},
		{"Commencement Date", "{{commencement_date}}"},
		{"Total Lives", "{{total_lives}}"},
		{"Total Sum Assured", "{{total_sum_assured}}"},
		{"Total Annual Premium", "{{total_annual_premium}}"},
	}
	var b strings.Builder
	b.WriteString(`<w:tbl>`)
	b.WriteString(`<w:tblPr><w:tblW w:w="9000" w:type="dxa"/>`)
	b.WriteString(`<w:tblBorders>`)
	for _, e := range []string{"top", "left", "bottom", "right", "insideH", "insideV"} {
		b.WriteString(`<w:` + e + ` w:val="single" w:sz="4" w:space="0" w:color="BFBFBF"/>`)
	}
	b.WriteString(`</w:tblBorders></w:tblPr>`)
	b.WriteString(`<w:tblGrid><w:gridCol w:w="3000"/><w:gridCol w:w="6000"/></w:tblGrid>`)

	for _, r := range rows {
		b.WriteString(`<w:tr>`)
		// Label cell
		b.WriteString(`<w:tc>`)
		b.WriteString(`<w:tcPr><w:tcW w:w="3000" w:type="dxa"/>`)
		b.WriteString(`<w:tcBorders>`)
		for _, e := range []string{"top", "left", "bottom", "right"} {
			b.WriteString(`<w:` + e + ` w:val="single" w:sz="4" w:space="0" w:color="BFBFBF"/>`)
		}
		b.WriteString(`</w:tcBorders>`)
		b.WriteString(`<w:shd w:val="clear" w:color="auto" w:fill="ECF0F1"/>`)
		b.WriteString(`</w:tcPr>`)
		b.WriteString(`<w:p><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">` + r[0] + `</w:t></w:r></w:p>`)
		b.WriteString(`</w:tc>`)
		// Value cell
		b.WriteString(`<w:tc>`)
		b.WriteString(`<w:tcPr><w:tcW w:w="6000" w:type="dxa"/>`)
		b.WriteString(`<w:tcBorders>`)
		for _, e := range []string{"top", "left", "bottom", "right"} {
			b.WriteString(`<w:` + e + ` w:val="single" w:sz="4" w:space="0" w:color="BFBFBF"/>`)
		}
		b.WriteString(`</w:tcBorders></w:tcPr>`)
		b.WriteString(`<w:p><w:r><w:t xml:space="preserve">` + r[1] + `</w:t></w:r></w:p>`)
		b.WriteString(`</w:tc>`)
		b.WriteString(`</w:tr>`)
	}
	b.WriteString(`</w:tbl>`)
	// Trailing empty paragraph (Word requires content between <w:tbl> and <w:sectPr>)
	b.WriteString(`<w:p/>`)
	return b.String()
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

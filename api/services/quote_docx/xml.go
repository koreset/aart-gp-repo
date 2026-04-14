package quote_docx

import (
	"fmt"
	"strings"
)

// xmlEscape escapes XML special characters
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// RunOptions configures a text run
type RunOptions struct {
	Text  string
	Bold  bool
	Italic bool
	Color string
	Size  int // half-points
}

// runXML generates <w:r> XML for a text run.
//
// Child elements of <w:rPr> must appear in the order defined by the OOXML
// schema (CT_RPr). Word for Mac validates this strictly. The order we use:
// rFonts, b, i, color, sz, szCs.
func runXML(opts RunOptions) string {
	var buf strings.Builder
	buf.WriteString("<w:r>")
	buf.WriteString("<w:rPr>")

	buf.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, Font, Font))
	if opts.Bold {
		buf.WriteString("<w:b/>")
	}
	if opts.Italic {
		buf.WriteString("<w:i/>")
	}
	if opts.Color != "" {
		buf.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, opts.Color))
	}
	if opts.Size > 0 {
		buf.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/>`, opts.Size))
		buf.WriteString(fmt.Sprintf(`<w:szCs w:val="%d"/>`, opts.Size))
	}

	buf.WriteString("</w:rPr>")
	buf.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, xmlEscape(opts.Text)))
	buf.WriteString("</w:r>")

	return buf.String()
}

// ParagraphOptions configures a paragraph
type ParagraphOptions struct {
	Alignment string // LEFT, CENTER, RIGHT, JUSTIFIED
	SpaceBefore int
	SpaceAfter  int
	Bold        bool
	Italic      bool
	Color       string
	Size        int
	Fill        string // background colour
	BorderBottom bool
	BorderColor  string
}

// paragraphXML generates <w:p> XML.
//
// Child elements of <w:pPr> must appear in schema order (CT_PPr). The
// relevant subset for this document: pBdr, shd, spacing, jc. Word for Mac
// rejects documents where these are out of order.
func paragraphXML(opts ParagraphOptions, runs []string) string {
	var buf strings.Builder
	buf.WriteString("<w:p>")
	buf.WriteString("<w:pPr>")

	// Paragraph border (must come before shd, spacing, jc)
	if opts.BorderBottom {
		borderColor := ColorPrimary
		if opts.BorderColor != "" {
			borderColor = opts.BorderColor
		}
		buf.WriteString(fmt.Sprintf(
			`<w:pBdr><w:bottom w:val="single" w:sz="8" w:color="%s" w:space="2"/></w:pBdr>`,
			borderColor,
		))
	}

	// Shading (before spacing and jc)
	if opts.Fill != "" {
		buf.WriteString(fmt.Sprintf(`<w:shd w:val="clear" w:color="auto" w:fill="%s"/>`, opts.Fill))
	}

	// Spacing (before jc)
	if opts.SpaceBefore > 0 || opts.SpaceAfter > 0 {
		buf.WriteString(`<w:spacing`)
		if opts.SpaceBefore > 0 {
			buf.WriteString(fmt.Sprintf(` w:before="%d"`, opts.SpaceBefore))
		}
		if opts.SpaceAfter > 0 {
			buf.WriteString(fmt.Sprintf(` w:after="%d"`, opts.SpaceAfter))
		}
		buf.WriteString(`/>`)
	}

	// Alignment (last of the group we emit)
	if opts.Alignment != "" {
		alignMap := map[string]string{
			"LEFT":      "left",
			"CENTER":    "center",
			"RIGHT":     "right",
			"JUSTIFIED": "both",
		}
		if align, ok := alignMap[opts.Alignment]; ok {
			buf.WriteString(fmt.Sprintf(`<w:jc w:val="%s"/>`, align))
		}
	}

	buf.WriteString("</w:pPr>")

	// Runs
	for _, run := range runs {
		buf.WriteString(run)
	}

	buf.WriteString("</w:p>")
	return buf.String()
}

// CellOptions configures a table cell
type CellOptions struct {
	Width        int
	Shading      string
	HasBorders   bool
	BorderColor  string
	VerticalAlign string
}

// cellXML generates <w:tc> XML.
//
// Child order of <w:tcPr> is schema-significant (CT_TcPr):
// tcW, gridSpan, hMerge/vMerge, tcBorders, shd, noWrap, tcMar,
// textDirection, tcFitText, vAlign, hideMark. Border edges must live
// inside a <w:tcBorders> wrapper; emitting them as direct children of
// <w:tcPr> makes Word refuse to open the document.
//
// Additionally, a <w:tc> must contain at least one <w:p> as content —
// empty cells must emit an empty paragraph, not nothing.
func cellXML(opts CellOptions, paragraphs []string) string {
	var buf strings.Builder
	buf.WriteString("<w:tc>")
	buf.WriteString("<w:tcPr>")

	// tcW
	if opts.Width > 0 {
		buf.WriteString(fmt.Sprintf(`<w:tcW w:w="%d" w:type="dxa"/>`, opts.Width))
	}

	// tcBorders — always emit a wrapper with each edge defined
	buf.WriteString(`<w:tcBorders>`)
	if opts.HasBorders {
		borderColor := "D5D8DC"
		if opts.BorderColor != "" {
			borderColor = opts.BorderColor
		}
		for _, edge := range []string{"top", "left", "bottom", "right"} {
			buf.WriteString(fmt.Sprintf(
				`<w:%s w:val="single" w:sz="4" w:space="0" w:color="%s"/>`,
				edge, borderColor,
			))
		}
	} else {
		for _, edge := range []string{"top", "left", "bottom", "right"} {
			buf.WriteString(fmt.Sprintf(`<w:%s w:val="nil"/>`, edge))
		}
	}
	buf.WriteString(`</w:tcBorders>`)

	// Shading (after tcBorders, before tcMar)
	if opts.Shading != "" {
		buf.WriteString(fmt.Sprintf(`<w:shd w:val="clear" w:color="auto" w:fill="%s"/>`, opts.Shading))
	}

	// tcMar
	buf.WriteString(`<w:tcMar>`)
	buf.WriteString(`<w:top w:w="60" w:type="dxa"/>`)
	buf.WriteString(`<w:left w:w="100" w:type="dxa"/>`)
	buf.WriteString(`<w:bottom w:w="60" w:type="dxa"/>`)
	buf.WriteString(`<w:right w:w="100" w:type="dxa"/>`)
	buf.WriteString(`</w:tcMar>`)

	// vAlign (last of the group we emit)
	if opts.VerticalAlign == "CENTER" {
		buf.WriteString(`<w:vAlign w:val="center"/>`)
	}

	buf.WriteString("</w:tcPr>")

	// Cell content — at least one paragraph is required.
	if len(paragraphs) == 0 {
		buf.WriteString(`<w:p/>`)
	} else {
		for _, para := range paragraphs {
			buf.WriteString(para)
		}
	}

	buf.WriteString("</w:tc>")
	return buf.String()
}

// RowOptions configures a table row
type RowOptions struct {
	IsHeader bool
}

// rowXML generates <w:tr> XML
func rowXML(opts RowOptions, cells []string) string {
	var buf strings.Builder
	buf.WriteString("<w:tr>")
	if opts.IsHeader {
		buf.WriteString(`<w:trPr><w:tblHeader/></w:trPr>`)
	}

	for _, cell := range cells {
		buf.WriteString(cell)
	}

	buf.WriteString("</w:tr>")
	return buf.String()
}

// TableOptions configures a table
type TableOptions struct {
	Width       int
	ColumnWidths []int
}

// tableXML generates <w:tbl> XML
func tableXML(opts TableOptions, rows []string) string {
	var buf strings.Builder
	buf.WriteString("<w:tbl>")
	buf.WriteString("<w:tblPr>")

	if opts.Width > 0 {
		buf.WriteString(fmt.Sprintf(`<w:tblW w:w="%d" w:type="dxa"/>`, opts.Width))
	}

	// Table grid (column widths)
	buf.WriteString("<w:tblGrid>")
	for _, width := range opts.ColumnWidths {
		buf.WriteString(fmt.Sprintf(`<w:gridCol w:w="%d"/>`, width))
	}
	buf.WriteString("</w:tblGrid>")

	buf.WriteString("</w:tblPr>")

	// Rows
	for _, row := range rows {
		buf.WriteString(row)
	}

	buf.WriteString("</w:tbl>")
	return buf.String()
}

// imageRunXML generates an image run with proper EMU calculation
func imageRunXML(relID string, widthPx, heightPx int, altText string) string {
	// Convert pixels to EMU (1 pt = 12700 EMU, assuming 96 DPI: 1px = 9525 EMU)
	const pxToEMU = 9525
	widthEMU := int64(widthPx * pxToEMU)
	heightEMU := int64(heightPx * pxToEMU)

	var buf strings.Builder
	buf.WriteString("<w:r>")
	buf.WriteString("<w:drawing>")
	buf.WriteString(fmt.Sprintf(
		`<wp:anchor distT="0" distB="0" distL="114300" distR="114300" simplePos="0" relativeHeight="251658240" behindDoc="0" locked="0" layoutInCell="1" allowOverlap="1">`,
	))
	buf.WriteString(`<wp:simplePos x="0" y="0"/>`)
	buf.WriteString(`<wp:positionH relativeFrom="column"><wp:align>right</wp:align></wp:positionH>`)
	buf.WriteString(`<wp:positionV relativeFrom="paragraph"><wp:posOffset>0</wp:posOffset></wp:positionV>`)
	buf.WriteString(fmt.Sprintf(`<wp:extent cx="%d" cy="%d"/>`, widthEMU, heightEMU))
	buf.WriteString(`<wp:effectExtent l="0" t="0" r="0" b="0"/>`)
	buf.WriteString(`<wp:wrapNone/>`)
	buf.WriteString(`<wp:docPr id="1" name="Image 1"/>`)
	buf.WriteString(`<wp:cNvGraphicFramePr/>`)
	buf.WriteString(`<a:graphic>`)
	buf.WriteString(`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">`)
	buf.WriteString(`<pic:pic>`)
	buf.WriteString(`<pic:nvPicPr>`)
	buf.WriteString(fmt.Sprintf(`<pic:cNvPr id="1" name="%s"/>`, altText))
	buf.WriteString(`<pic:cNvPicPr/>`)
	buf.WriteString(`</pic:nvPicPr>`)
	buf.WriteString(`<pic:blipFill>`)
	buf.WriteString(fmt.Sprintf(`<a:blip r:embed="%s"/>`, relID))
	buf.WriteString(`<a:stretch><a:fillRect/></a:stretch>`)
	buf.WriteString(`</pic:blipFill>`)
	buf.WriteString(`<pic:spPr>`)
	buf.WriteString(fmt.Sprintf(`<a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm>`, widthEMU, heightEMU))
	buf.WriteString(`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>`)
	buf.WriteString(`</pic:spPr>`)
	buf.WriteString(`</pic:pic>`)
	buf.WriteString(`</a:graphicData>`)
	buf.WriteString(`</a:graphic>`)
	buf.WriteString(`</wp:anchor>`)
	buf.WriteString(`</w:drawing>`)
	buf.WriteString(`</w:r>`)

	return buf.String()
}

// sectionPropsXML generates section properties. Child element order is
// schema-significant in OOXML: headerReference/footerReference, then pgSz,
// pgMar, cols. Header/footer rel IDs are referenced on every section so Word
// applies the shared header1.xml / footer1.xml on each page.
func sectionPropsXML(orientation string, marginT, marginB, marginL, marginR int) string {
	var buf strings.Builder
	buf.WriteString("<w:sectPr>")

	// Header/footer references — rId1 and rId2 come from document.xml.rels
	buf.WriteString(`<w:headerReference w:type="default" r:id="rId1"/>`)
	buf.WriteString(`<w:footerReference w:type="default" r:id="rId2"/>`)

	// Page size and orientation
	if orientation == "LANDSCAPE" {
		buf.WriteString(fmt.Sprintf(
			`<w:pgSz w:w="%d" w:h="%d" w:orient="landscape"/>`,
			A4Height, A4Width,
		))
	} else {
		buf.WriteString(fmt.Sprintf(
			`<w:pgSz w:w="%d" w:h="%d" w:orient="portrait"/>`,
			A4Width, A4Height,
		))
	}

	// Margins
	buf.WriteString(fmt.Sprintf(
		`<w:pgMar w:top="%d" w:right="%d" w:bottom="%d" w:left="%d" w:header="720" w:footer="720" w:gutter="0"/>`,
		marginT, marginR, marginB, marginL,
	))

	// Columns (default: 1)
	buf.WriteString(`<w:cols w:num="1"/>`)
	buf.WriteString(`<w:docGrid w:linePitch="360"/>`)

	buf.WriteString("</w:sectPr>")
	return buf.String()
}

// wrapSectPrInParagraph wraps a sectPr in an empty paragraph's pPr. This is
// the OOXML pattern for a section break that ends a section mid-document;
// the final section's sectPr must NOT be wrapped — it sits directly as a
// child of <w:body>.
func wrapSectPrInParagraph(sectPr string) string {
	return "<w:p><w:pPr>" + sectPr + "</w:pPr></w:p>"
}

// Helper: Section heading paragraph
func sectionHeadingXML(text string) string {
	runs := []string{
		runXML(RunOptions{
			Text:  text,
			Bold:  true,
			Color: ColorPrimary,
			Size:  SizeHeading,
		}),
	}
	return paragraphXML(ParagraphOptions{
		SpaceBefore:  200,
		SpaceAfter:   120,
		BorderBottom: true,
		BorderColor:  ColorPrimary,
	}, runs)
}

// Helper: Category heading paragraph
func categoryHeadingXML(text string) string {
	runs := []string{
		runXML(RunOptions{
			Text:  text,
			Bold:  true,
			Color: ColorDark,
			Size:  SizeSubheading,
		}),
	}
	return paragraphXML(ParagraphOptions{
		SpaceBefore: 160,
		SpaceAfter:  80,
	}, runs)
}

// Helper: Body text paragraph
func bodyTextXML(text string, italic bool) string {
	runs := []string{
		runXML(RunOptions{
			Text:   text,
			Italic: italic,
			Color:  ColorDark,
			Size:   SizeBody,
		}),
	}
	return paragraphXML(ParagraphOptions{
		Alignment:   "JUSTIFIED",
		SpaceAfter:  100,
	}, runs)
}

// Helper: Spacer paragraph
func spacerXML(spaceBefore int) string {
	return paragraphXML(ParagraphOptions{
		SpaceBefore: spaceBefore,
	}, []string{})
}

// Helper: Header row for tables
func headerRowXML(labels []string, widths []int) string {
	cells := []string{}
	for i, label := range labels {
		para := paragraphXML(ParagraphOptions{
			Alignment: map[int]string{0: "LEFT"}[i], // First col left, rest center
			Bold:      true,
			Color:     ColorWhite,
		}, []string{
			runXML(RunOptions{
				Text:  label,
				Bold:  true,
				Color: ColorWhite,
				Size:  SizeBody,
			}),
		})

		cells = append(cells, cellXML(CellOptions{
			Width:       widths[i],
			Shading:     ColorPrimary,
			HasBorders:  true,
			BorderColor: "D5D8DC",
		}, []string{para}))
	}

	return rowXML(RowOptions{IsHeader: true}, cells)
}

// Helper: Data row for tables
func dataRowXML(values []string, widths []int, alignments []string, bold bool, fillColor string) string {
	cells := []string{}
	for i, val := range values {
		align := alignments[i]
		if align == "" {
			if i == 0 {
				align = "LEFT"
			} else {
				align = "RIGHT"
			}
		}

		para := paragraphXML(ParagraphOptions{
			Alignment: align,
			Fill:      fillColor,
		}, []string{
			runXML(RunOptions{
				Text:  val,
				Bold:  bold && i == 0,
				Color: ColorDark,
				Size:  SizeCaption,
			}),
		})

		cells = append(cells, cellXML(CellOptions{
			Width:       widths[i],
			Shading:     fillColor,
			HasBorders:  true,
			BorderColor: "D5D8DC",
		}, []string{para}))
	}

	return rowXML(RowOptions{}, cells)
}

// Helper: Key-value table (2-column)
func keyValueTableXML(rows []LabelValueRow, contentWidth int) string {
	labelWidth := (contentWidth * 37) / 100
	valueWidth := contentWidth - labelWidth

	tableRows := []string{}

	// Header row not typically used for key-value, just data rows
	for i, row := range rows {
		fillColor := ""
		if i%2 == 1 {
			fillColor = ColorAltRow
		}

		// Label cell
		labelPara := paragraphXML(ParagraphOptions{
			Fill: ColorLightFill,
		}, []string{
			runXML(RunOptions{
				Text:  row.Label,
				Bold:  true,
				Color: ColorDark,
				Size:  SizeBody,
			}),
		})
		labelCell := cellXML(CellOptions{
			Width:       labelWidth,
			Shading:     ColorLightFill,
			HasBorders:  true,
			BorderColor: "D5D8DC",
		}, []string{labelPara})

		// Value cell
		valuePara := paragraphXML(ParagraphOptions{
			Fill: fillColor,
		}, []string{
			runXML(RunOptions{
				Text:  row.Value,
				Color: ColorSecondary,
				Size:  SizeBody,
			}),
		})
		valueCell := cellXML(CellOptions{
			Width:       valueWidth,
			Shading:     fillColor,
			HasBorders:  true,
			BorderColor: "D5D8DC",
		}, []string{valuePara})

		tableRows = append(tableRows, rowXML(RowOptions{}, []string{labelCell, valueCell}))
	}

	return tableXML(TableOptions{
		Width:        contentWidth,
		ColumnWidths: []int{labelWidth, valueWidth},
	}, tableRows)
}

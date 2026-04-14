package quote_docx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// Package represents a DOCX document package
type Package struct {
	Body       string // document.xml body content
	Header     string // header1.xml content
	Footer     string // footer1.xml content
	Logo       []byte // image bytes
	LogoMIME   string // e.g., "image/png", "image/jpeg"
}

// Build assembles the DOCX as a ZIP archive and returns the bytes
func (p *Package) Build() ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	defer zw.Close()

	// Determine logo extension
	logoExt := "png"
	if strings.Contains(p.LogoMIME, "jpeg") || strings.Contains(p.LogoMIME, "jpg") {
		logoExt = "jpg"
	}

	// 1. [Content_Types].xml
	contentTypes := buildContentTypesXML(len(p.Logo) > 0, logoExt)
	if err := addFileToZip(zw, "[Content_Types].xml", contentTypes); err != nil {
		return nil, err
	}

	// 2. _rels/.rels
	rels := buildRootRels()
	if err := addFileToZip(zw, "_rels/.rels", rels); err != nil {
		return nil, err
	}

	// 3. word/document.xml
	docXML := buildDocumentXML(p.Body, p.Header, p.Footer, len(p.Logo) > 0)
	if err := addFileToZip(zw, "word/document.xml", docXML); err != nil {
		return nil, err
	}

	// 4. word/_rels/document.xml.rels
	docRels := buildDocumentRels(len(p.Logo) > 0)
	if err := addFileToZip(zw, "word/_rels/document.xml.rels", docRels); err != nil {
		return nil, err
	}

	// 5. word/styles.xml
	styles := buildStylesXML()
	if err := addFileToZip(zw, "word/styles.xml", styles); err != nil {
		return nil, err
	}

	// 6. word/header1.xml
	if err := addFileToZip(zw, "word/header1.xml", p.Header); err != nil {
		return nil, err
	}

	// 7. word/_rels/header1.xml.rels
	headerRels := buildHeaderFooterRels(len(p.Logo) > 0)
	if err := addFileToZip(zw, "word/_rels/header1.xml.rels", headerRels); err != nil {
		return nil, err
	}

	// 8. word/footer1.xml
	if err := addFileToZip(zw, "word/footer1.xml", p.Footer); err != nil {
		return nil, err
	}

	// 9. word/_rels/footer1.xml.rels
	if err := addFileToZip(zw, "word/_rels/footer1.xml.rels", headerRels); err != nil {
		return nil, err
	}

	// 10. Logo image (if present)
	if len(p.Logo) > 0 {
		filename := fmt.Sprintf("word/media/image1.%s", logoExt)
		w, err := zw.Create(filename)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(p.Logo); err != nil {
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper: Add string content to ZIP
func addFileToZip(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, content)
	return err
}

// buildContentTypesXML generates [Content_Types].xml
func buildContentTypesXML(hasLogo bool, logoExt string) string {
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)

	// Default extensions
	buf.WriteString(`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`)
	buf.WriteString(`<Default Extension="xml" ContentType="application/xml"/>`)

	// Logo type
	if hasLogo {
		if logoExt == "jpg" || logoExt == "jpeg" {
			buf.WriteString(`<Default Extension="jpg" ContentType="image/jpeg"/>`)
		} else {
			buf.WriteString(`<Default Extension="png" ContentType="image/png"/>`)
		}
	}

	// Override for document types
	buf.WriteString(`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>`)
	buf.WriteString(`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>`)
	buf.WriteString(`<Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>`)
	buf.WriteString(`<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>`)

	buf.WriteString(`</Types>`)
	return buf.String()
}

// buildRootRels generates _rels/.rels
func buildRootRels() string {
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	buf.WriteString(`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>`)
	buf.WriteString(`</Relationships>`)
	return buf.String()
}

// buildDocumentRels generates word/_rels/document.xml.rels
func buildDocumentRels(hasLogo bool) string {
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	buf.WriteString(`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/>`)
	buf.WriteString(`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/>`)
	buf.WriteString(`<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`)

	if hasLogo {
		buf.WriteString(`<Relationship Id="rIdLogo" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/image1.png"/>`)
	}

	buf.WriteString(`</Relationships>`)
	return buf.String()
}

// buildHeaderFooterRels generates word/_rels/header1.xml.rels (or footer1.xml.rels)
func buildHeaderFooterRels(hasLogo bool) string {
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)

	if hasLogo {
		buf.WriteString(`<Relationship Id="rIdLogo" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="../media/image1.png"/>`)
	}

	buf.WriteString(`</Relationships>`)
	return buf.String()
}

// buildDocumentXML wraps body content in document structure
func buildDocumentXML(body, header, footer string, hasLogo bool) string {
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`)

	buf.WriteString(`<w:body>`)
	buf.WriteString(body)
	buf.WriteString(`</w:body>`)

	buf.WriteString(`</w:document>`)
	return buf.String()
}

// buildStylesXML generates word/styles.xml with minimal styles
func buildStylesXML() string {
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`)

	// Default Normal style
	buf.WriteString(`<w:docDefaults>`)
	buf.WriteString(`<w:rPrDefault>`)
	buf.WriteString(`<w:rPr>`)
	buf.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, Font, Font))
	buf.WriteString(`<w:sz w:val="20"/>`)
	buf.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, ColorDark))
	buf.WriteString(`</w:rPr>`)
	buf.WriteString(`</w:rPrDefault>`)
	buf.WriteString(`<w:pPrDefault/>`)
	buf.WriteString(`</w:docDefaults>`)

	// Normal style
	buf.WriteString(`<w:style w:type="paragraph" w:default="1" w:styleId="Normal">`)
	buf.WriteString(`<w:name w:val="Normal"/>`)
	buf.WriteString(`<w:qFormat/>`)
	buf.WriteString(`</w:style>`)

	buf.WriteString(`</w:styles>`)
	return buf.String()
}

// buildHeaderXML generates header1.xml with scheme name
func buildHeaderXML(schemeName string) string {
	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)

	buf.WriteString(`<w:p>`)
	buf.WriteString(`<w:pPr>`)
	buf.WriteString(`<w:pBdr>`)
	buf.WriteString(`<w:bottom w:val="single" w:sz="6" w:color="ECF0F1" w:space="1"/>`)
	buf.WriteString(`</w:pBdr>`)
	buf.WriteString(`<w:spacing w:after="100"/>`)
	buf.WriteString(`</w:pPr>`)

	buf.WriteString(`<w:r>`)
	buf.WriteString(`<w:rPr>`)
	buf.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, Font, Font))
	buf.WriteString(`<w:sz w:val="18"/>`)
	buf.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, ColorSecondary))
	buf.WriteString(`</w:rPr>`)
	buf.WriteString(fmt.Sprintf(`<w:t>%s - Quotation</w:t>`, xmlEscape(schemeName)))
	buf.WriteString(`</w:r>`)

	buf.WriteString(`</w:p>`)
	buf.WriteString(`</w:hdr>`)

	return buf.String()
}

// buildFooterXML generates footer1.xml with generated date and page number
func buildFooterXML() string {
	dateStr := time.Now().Format("02 Jan 2006")

	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`)

	buf.WriteString(`<w:p>`)
	buf.WriteString(`<w:pPr>`)
	buf.WriteString(`<w:pBdr>`)
	buf.WriteString(`<w:top w:val="single" w:sz="6" w:color="ECF0F1" w:space="1"/>`)
	buf.WriteString(`</w:pBdr>`)
	buf.WriteString(`<w:tabs>`)
	buf.WriteString(`<w:tab w:val="right" w:pos="9500"/>`)
	buf.WriteString(`</w:tabs>`)
	buf.WriteString(`</w:pPr>`)

	buf.WriteString(`<w:r>`)
	buf.WriteString(`<w:rPr>`)
	buf.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, Font, Font))
	buf.WriteString(`<w:sz w:val="18"/>`)
	buf.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, ColorSecondary))
	buf.WriteString(`</w:rPr>`)
	buf.WriteString(fmt.Sprintf(`<w:t>Generated on %s</w:t>`, xmlEscape(dateStr)))
	buf.WriteString(`</w:r>`)

	buf.WriteString(`<w:r>`)
	buf.WriteString(`<w:rPr>`)
	buf.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, Font, Font))
	buf.WriteString(`<w:sz w:val="18"/>`)
	buf.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, ColorSecondary))
	buf.WriteString(`</w:rPr>`)
	buf.WriteString(`<w:tab/>`)
	buf.WriteString(`<w:t>Page </w:t>`)
	buf.WriteString(`</w:r>`)

	// Page number field
	buf.WriteString(`<w:fldSimple w:instr=" PAGE ">`)
	buf.WriteString(`<w:r>`)
	buf.WriteString(`<w:rPr>`)
	buf.WriteString(fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, Font, Font))
	buf.WriteString(`<w:sz w:val="18"/>`)
	buf.WriteString(fmt.Sprintf(`<w:color w:val="%s"/>`, ColorSecondary))
	buf.WriteString(`</w:rPr>`)
	buf.WriteString(`<w:t>0</w:t>`)
	buf.WriteString(`</w:r>`)
	buf.WriteString(`</w:fldSimple>`)

	buf.WriteString(`</w:p>`)
	buf.WriteString(`</w:ftr>`)

	return buf.String()
}

// BuildFullPackage is a convenience function that builds header/footer and returns the full ZIP bytes
func BuildFullPackage(bodyXML string, schemeName string, logo []byte, logoMIME string) ([]byte, error) {
	header := buildHeaderXML(schemeName)
	footer := buildFooterXML()

	pkg := &Package{
		Body:     bodyXML,
		Header:   header,
		Footer:   footer,
		Logo:     logo,
		LogoMIME: logoMIME,
	}

	return pkg.Build()
}

// xmlEscapeForAttribute escapes XML attributes (already defined but repeated for clarity)
func xmlEscapeForAttribute(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

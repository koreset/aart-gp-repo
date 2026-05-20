// Package payment_letter builds claim payment confirmation letters as DOCX
// (and optionally PDF via api/services/docpdf). The package is intentionally
// self-contained — it builds a minimal OOXML document directly, rather than
// reusing the quote_docx helpers, because the letter only needs paragraphs,
// a header with logo, and a signature block. Bringing in the full table /
// section machinery from quote_docx would be more weight than the letter
// needs.
package payment_letter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	font  = "Calibri"
	color = "1F2937" // dark slate
)

// docxPackage assembles the OOXML zip with body, header (logo), footer, and
// styles. Signature image (if any) is referenced from the body via a separate
// relationship.
type docxPackage struct {
	body           string
	headerXML      string
	footerXML      string
	logo           []byte
	logoMIME       string
	signature      []byte
	signatureMIME  string
}

func (p *docxPackage) Build() ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	logoExt := mimeToExt(p.logoMIME)
	sigExt := mimeToExt(p.signatureMIME)

	hasLogo := len(p.logo) > 0
	hasSig := len(p.signature) > 0

	if err := addFile(zw, "[Content_Types].xml", contentTypesXML(hasLogo, logoExt, hasSig, sigExt)); err != nil {
		return nil, err
	}
	if err := addFile(zw, "_rels/.rels", rootRelsXML()); err != nil {
		return nil, err
	}
	if err := addFile(zw, "word/document.xml", documentXML(p.body)); err != nil {
		return nil, err
	}
	if err := addFile(zw, "word/_rels/document.xml.rels", documentRelsXML(hasLogo, hasSig, sigExt)); err != nil {
		return nil, err
	}
	if err := addFile(zw, "word/styles.xml", stylesXML()); err != nil {
		return nil, err
	}
	if err := addFile(zw, "word/header1.xml", p.headerXML); err != nil {
		return nil, err
	}
	if err := addFile(zw, "word/_rels/header1.xml.rels", headerFooterRelsXML(hasLogo, logoExt)); err != nil {
		return nil, err
	}
	if err := addFile(zw, "word/footer1.xml", p.footerXML); err != nil {
		return nil, err
	}
	if err := addFile(zw, "word/_rels/footer1.xml.rels", headerFooterRelsXML(false, "")); err != nil {
		return nil, err
	}

	if hasLogo {
		if err := addBinary(zw, fmt.Sprintf("word/media/logo.%s", logoExt), p.logo); err != nil {
			return nil, err
		}
	}
	if hasSig {
		if err := addBinary(zw, fmt.Sprintf("word/media/signature.%s", sigExt), p.signature); err != nil {
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func addFile(zw *zip.Writer, name, content string) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, content)
	return err
}

func addBinary(zw *zip.Writer, name string, data []byte) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func mimeToExt(mime string) string {
	switch {
	case strings.Contains(mime, "jpeg"), strings.Contains(mime, "jpg"):
		return "jpg"
	case strings.Contains(mime, "svg"):
		return "svg"
	default:
		return "png"
	}
}

func contentTypesXML(hasLogo bool, logoExt string, hasSig bool, sigExt string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString(`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">`)
	b.WriteString(`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>`)
	b.WriteString(`<Default Extension="xml" ContentType="application/xml"/>`)

	exts := map[string]string{}
	if hasLogo {
		exts[logoExt] = extToMime(logoExt)
	}
	if hasSig {
		exts[sigExt] = extToMime(sigExt)
	}
	for ext, mime := range exts {
		b.WriteString(fmt.Sprintf(`<Default Extension="%s" ContentType="%s"/>`, ext, mime))
	}

	b.WriteString(`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>`)
	b.WriteString(`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>`)
	b.WriteString(`<Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>`)
	b.WriteString(`<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>`)
	b.WriteString(`</Types>`)
	return b.String()
}

func extToMime(ext string) string {
	switch ext {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "svg":
		return "image/svg+xml"
	default:
		return "image/png"
	}
}

func rootRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
		`</Relationships>`
}

func documentRelsXML(hasLogo, hasSig bool, sigExt string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	b.WriteString(`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/>`)
	b.WriteString(`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/>`)
	b.WriteString(`<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`)
	if hasSig {
		b.WriteString(fmt.Sprintf(`<Relationship Id="rIdSig" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/signature.%s"/>`, sigExt))
	}
	_ = hasLogo
	b.WriteString(`</Relationships>`)
	return b.String()
}

func headerFooterRelsXML(hasLogo bool, logoExt string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	if hasLogo {
		b.WriteString(fmt.Sprintf(`<Relationship Id="rIdLogo" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="../media/logo.%s"/>`, logoExt))
	}
	b.WriteString(`</Relationships>`)
	return b.String()
}

func documentXML(body string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" ` +
		`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">` +
		`<w:body>` + body + sectionPropsXML() + `</w:body>` +
		`</w:document>`
}

func sectionPropsXML() string {
	return `<w:sectPr>` +
		`<w:headerReference w:type="default" r:id="rId1"/>` +
		`<w:footerReference w:type="default" r:id="rId2"/>` +
		`<w:pgSz w:w="12240" w:h="15840"/>` +
		`<w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/>` +
		`</w:sectPr>`
}

func stylesXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:docDefaults><w:rPrDefault><w:rPr>` +
		fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, font, font) +
		`<w:sz w:val="22"/>` +
		fmt.Sprintf(`<w:color w:val="%s"/>`, color) +
		`</w:rPr></w:rPrDefault><w:pPrDefault/></w:docDefaults>` +
		`<w:style w:type="paragraph" w:default="1" w:styleId="Normal">` +
		`<w:name w:val="Normal"/><w:qFormat/></w:style>` +
		`</w:styles>`
}

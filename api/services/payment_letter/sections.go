package payment_letter

import (
	"fmt"
	"strings"
	"time"

	"api/models"
)

// xmlEscape escapes text for embedding inside an OOXML text run.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// paragraph wraps a list of runs in <w:p> with optional alignment / spacing.
func paragraph(align string, spaceAfter int, runs ...string) string {
	var b strings.Builder
	b.WriteString(`<w:p>`)
	if align != "" || spaceAfter > 0 {
		b.WriteString(`<w:pPr>`)
		if align != "" {
			b.WriteString(fmt.Sprintf(`<w:jc w:val="%s"/>`, align))
		}
		if spaceAfter > 0 {
			b.WriteString(fmt.Sprintf(`<w:spacing w:after="%d"/>`, spaceAfter))
		}
		b.WriteString(`</w:pPr>`)
	}
	for _, r := range runs {
		b.WriteString(r)
	}
	b.WriteString(`</w:p>`)
	return b.String()
}

// textRun is a simple text run with optional bold / size.
func textRun(text string, bold bool, sizeHalfPts int) string {
	var b strings.Builder
	b.WriteString(`<w:r><w:rPr>`)
	if bold {
		b.WriteString(`<w:b/>`)
	}
	if sizeHalfPts > 0 {
		b.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/>`, sizeHalfPts))
	}
	b.WriteString(`</w:rPr>`)
	b.WriteString(fmt.Sprintf(`<w:t xml:space="preserve">%s</w:t>`, xmlEscape(text)))
	b.WriteString(`</w:r>`)
	return b.String()
}

// imageRun renders an embedded image referenced by relationship id.
// widthEMU and heightEMU are in English Metric Units (914400 per inch).
func imageRun(relID string, widthEMU, heightEMU int, altText string) string {
	return fmt.Sprintf(`<w:r><w:drawing>
<wp:inline distT="0" distB="0" distL="0" distR="0">
<wp:extent cx="%d" cy="%d"/>
<wp:docPr id="1" name="%s" descr="%s"/>
<wp:cNvGraphicFramePr><a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/></wp:cNvGraphicFramePr>
<a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">
<pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
<pic:nvPicPr><pic:cNvPr id="1" name="%s"/><pic:cNvPicPr/></pic:nvPicPr>
<pic:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>
<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr>
</pic:pic></a:graphicData></a:graphic>
</wp:inline></w:drawing></w:r>`, widthEMU, heightEMU, xmlEscape(altText), xmlEscape(altText), xmlEscape(altText), relID, widthEMU, heightEMU)
}

// LetterInput is the resolved data needed to render a single payment letter.
// It is built by GenerateAndRecord from the claim, payment proof, and settings.
type LetterInput struct {
	Claim       models.GroupSchemeClaim
	PaidAt      time.Time
	Settings    models.PaymentLetterSetting
	LetterRef   string
}

// buildHeaderXML renders the letterhead — logo on the right (if present) and
// company address block on the left.
func buildHeaderXML(s models.PaymentLetterSetting) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString(`<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" ` +
		`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`)

	if len(s.Logo) > 0 {
		b.WriteString(paragraph("right", 60, imageRun("rIdLogo", 1524000, 571500, "Company logo")))
	}
	if s.CompanyName != "" {
		b.WriteString(paragraph("", 40, textRun(s.CompanyName, true, 24)))
	}
	addr := strings.TrimSpace(strings.Join([]string{s.AddressLine1, s.AddressLine2, s.AddressLine3}, ", "))
	addr = strings.Trim(addr, ", ")
	if addr != "" {
		b.WriteString(paragraph("", 20, textRun(addr, false, 18)))
	}
	tail := strings.TrimSpace(strings.Join([]string{s.City, s.PostalCode, s.Country}, ", "))
	tail = strings.Trim(tail, ", ")
	if tail != "" {
		b.WriteString(paragraph("", 20, textRun(tail, false, 18)))
	}
	contactParts := []string{}
	if s.Phone != "" {
		contactParts = append(contactParts, "Tel: "+s.Phone)
	}
	if s.Email != "" {
		contactParts = append(contactParts, "Email: "+s.Email)
	}
	if s.Website != "" {
		contactParts = append(contactParts, s.Website)
	}
	if len(contactParts) > 0 {
		b.WriteString(paragraph("", 20, textRun(strings.Join(contactParts, "  •  "), false, 18)))
	}
	b.WriteString(`</w:hdr>`)
	return b.String()
}

// buildFooterXML renders a thin footer with the generated date.
func buildFooterXML() string {
	dateStr := time.Now().Format("02 January 2006")
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		paragraph("center", 0, textRun("Generated on "+dateStr, false, 16)) +
		`</w:ftr>`
}

// buildBodyXML composes the full letter body: salutation, payment confirmation,
// account details, member reference (if applicable), closing, signature block.
func buildBodyXML(in LetterInput) string {
	var b strings.Builder
	c := in.Claim
	s := in.Settings

	// Date line (right-aligned)
	b.WriteString(paragraph("right", 200, textRun(in.PaidAt.Format("02 January 2006"), false, 22)))

	// Reference line
	if in.LetterRef != "" {
		b.WriteString(paragraph("", 200, textRun("Ref: "+in.LetterRef, true, 20)))
	}

	// Recipient block
	b.WriteString(paragraph("", 40, textRun(c.ClaimantName, true, 22)))
	if c.ClaimantIDNumber != "" {
		b.WriteString(paragraph("", 200, textRun("ID: "+c.ClaimantIDNumber, false, 20)))
	}

	// Salutation
	b.WriteString(paragraph("", 200, textRun(fmt.Sprintf("Dear %s,", c.ClaimantName), false, 22)))

	// Subject (bold heading)
	b.WriteString(paragraph("", 200,
		textRun("RE: Confirmation of Claim Payment — Claim "+c.ClaimNumber, true, 22),
	))

	// Intro paragraph — configurable in settings, with sensible default.
	intro := strings.TrimSpace(s.LetterIntroTemplate)
	if intro == "" {
		intro = "We confirm that the above claim has been processed and the agreed payment has been disbursed as set out below."
	}
	b.WriteString(paragraph("", 200, textRun(intro, false, 22)))

	// Member reference (when claimant differs from member, e.g. funeral beneficiary)
	if c.MemberName != "" && !strings.EqualFold(strings.TrimSpace(c.MemberName), strings.TrimSpace(c.ClaimantName)) {
		memberLine := fmt.Sprintf("This payment is made in respect of the late %s (member).", c.MemberName)
		if c.RelationshipToMember != "" && !strings.EqualFold(c.RelationshipToMember, "self") {
			memberLine = fmt.Sprintf("This payment is made in respect of the late %s (member), to whom you are listed as %s.",
				c.MemberName, c.RelationshipToMember)
		}
		b.WriteString(paragraph("", 200, textRun(memberLine, false, 22)))
	}

	// Scheme / benefit context
	if c.SchemeName != "" || c.BenefitName != "" {
		ctx := "Scheme: " + c.SchemeName
		if c.BenefitName != "" {
			if ctx != "Scheme: " {
				ctx += "   |   "
			}
			ctx += "Benefit: " + c.BenefitName
		}
		b.WriteString(paragraph("", 100, textRun(ctx, false, 20)))
	}

	// Payment details block — bold key/value lines
	b.WriteString(paragraph("", 60, textRun("Payment details:", true, 22)))
	b.WriteString(paragraph("", 40, textRun(fmt.Sprintf("Amount paid:        %s", formatAmount(c.ClaimAmount)), false, 22)))
	b.WriteString(paragraph("", 40, textRun(fmt.Sprintf("Date of payment:    %s", in.PaidAt.Format("02 January 2006")), false, 22)))
	if c.BankName != "" {
		b.WriteString(paragraph("", 40, textRun("Bank:                "+c.BankName, false, 22)))
	}
	if c.BankAccountNumber != "" {
		b.WriteString(paragraph("", 40, textRun("Account number:     "+maskAccount(c.BankAccountNumber), false, 22)))
	}
	if c.AccountHolderName != "" {
		b.WriteString(paragraph("", 200, textRun("Account holder:     "+c.AccountHolderName, false, 22)))
	}

	// Account holder mismatch note
	if c.AccountHolderName != "" && c.ClaimantName != "" &&
		!strings.EqualFold(strings.TrimSpace(c.AccountHolderName), strings.TrimSpace(c.ClaimantName)) {
		mismatch := fmt.Sprintf("Please note that the payment was made to the account in the name of %s as instructed.", c.AccountHolderName)
		b.WriteString(paragraph("", 200, textRun(mismatch, false, 22)))
	}

	// Closing paragraph — configurable
	closing := strings.TrimSpace(s.LetterClosingTemplate)
	if closing == "" {
		closing = "Should you have any queries regarding this payment, please contact our claims department using the details on the letterhead above."
	}
	b.WriteString(paragraph("", 200, textRun(closing, false, 22)))

	b.WriteString(paragraph("", 400, textRun("Yours faithfully,", false, 22)))

	// Signature image (if any) followed by signatory name + title.
	if len(s.Signature) > 0 {
		b.WriteString(paragraph("", 40, imageRun("rIdSig", 1524000, 457200, "Signature")))
	}
	if s.SignatoryName != "" {
		b.WriteString(paragraph("", 40, textRun(s.SignatoryName, true, 22)))
	}
	if s.SignatoryTitle != "" {
		b.WriteString(paragraph("", 40, textRun(s.SignatoryTitle, false, 20)))
	}
	if s.CompanyName != "" {
		b.WriteString(paragraph("", 0, textRun(s.CompanyName, false, 20)))
	}

	return b.String()
}

// formatAmount renders the claim amount with thousand separators and 2dp.
// Currency symbol is omitted because schemes can be multi-currency; the
// numeric format is unambiguous in context.
func formatAmount(v float64) string {
	// Two-decimal places with comma thousands separator.
	whole := int64(v)
	frac := int64((v - float64(whole)) * 100)
	if frac < 0 {
		frac = -frac
	}
	s := fmt.Sprintf("%d", whole)
	// Insert commas
	n := len(s)
	if n > 3 {
		var b strings.Builder
		startNeg := 0
		if n > 0 && s[0] == '-' {
			b.WriteByte('-')
			s = s[1:]
			n = len(s)
			startNeg = 1
		}
		first := n % 3
		if first > 0 {
			b.WriteString(s[:first])
			if n > first {
				b.WriteByte(',')
			}
		}
		for i := first; i < n; i += 3 {
			b.WriteString(s[i : i+3])
			if i+3 < n {
				b.WriteByte(',')
			}
		}
		_ = startNeg
		s = b.String()
	}
	return fmt.Sprintf("%s.%02d", s, frac)
}

// maskAccount returns the last 4 digits prefixed with dots, e.g. ••••1234.
// Visual confirmation without disclosing the full account number in print.
func maskAccount(acc string) string {
	if len(acc) <= 4 {
		return acc
	}
	return "••••" + acc[len(acc)-4:]
}

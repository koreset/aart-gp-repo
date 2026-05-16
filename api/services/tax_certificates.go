package services

import (
	appLog "api/log"
	"api/models"
	"errors"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ──────────────────────────────────────────────
// IT3(a) tax certificates (Phase 4)
// ──────────────────────────────────────────────
//
// A paid claim against certain benefit types (lump-sum life cover, lump-sum
// disability) is a tax event for the beneficiary; SARS expects the insurer
// to issue an IT3(a) certificate. Phase 4 generates an HTML certificate
// containing the disclosure fields and stores it on disk. Production
// installations swap the renderer for a real PDF library without changing
// the calling code — the file is read back through the same id-keyed
// download endpoint regardless of format.

// taxDisclosingBenefitTokens lists case-insensitive substrings that mark a
// benefit as triggering an IT3(a). Conservative — anything not on the list
// is treated as non-disclosing. Real installs override via the
// IsTaxDisclosing hook below.
var taxDisclosingBenefitTokens = []string{
	"gla",            // group life assurance
	"life",           // any "life cover" naming
	"funeral",        // funeral payouts often disclosed
	"disability",     // lump-sum disability
	"income protect", // income protection lump sums
	"phi",            // permanent health insurance
}

// IsTaxDisclosing returns true when a paid claim should trigger certificate
// generation. The default uses substring-matching against the benefit name;
// override at startup by reassigning the variable.
var IsTaxDisclosing = func(benefitName string) bool {
	n := strings.ToLower(benefitName)
	for _, tok := range taxDisclosingBenefitTokens {
		if strings.Contains(n, tok) {
			return true
		}
	}
	return false
}

// taxYearFor returns the SARS tax year of assessment for a payment date.
// SA tax year runs 1 March → end of February; the year of assessment is
// labelled by the year it ends in.
func taxYearFor(t time.Time) int {
	y := t.Year()
	if t.Month() >= time.March {
		return y + 1
	}
	return y
}

// generateCertificateRef returns a sortable unique reference, e.g.
// "IT3A-2026-000123".
func generateCertificateRef(taxYear int) string {
	var count int64
	DB.Model(&models.PaymentTaxCertificate{}).Where("tax_year = ?", taxYear).Count(&count)
	return fmt.Sprintf("IT3A-%d-%06d", taxYear, count+1)
}

// EnsureTaxCertificateForItem upserts a certificate for the given schedule
// item when its benefit is tax-disclosing. Returns the (possibly pre-
// existing) certificate row or a zero value if no certificate is required.
// Called from the reconciliation path for every line marked "paid".
//
// Idempotent: the unique index on (schedule_item_id, tax_year) means a
// second call within the same tax year returns the same row.
func EnsureTaxCertificateForItem(scheduleID, itemID int, user models.AppUser) (models.PaymentTaxCertificate, error) {
	var item models.ClaimPaymentScheduleItem
	if err := DB.Where("id = ? AND schedule_id = ?", itemID, scheduleID).First(&item).Error; err != nil {
		return models.PaymentTaxCertificate{}, err
	}
	if !IsTaxDisclosing(item.BenefitName) {
		return models.PaymentTaxCertificate{}, nil
	}

	now := time.Now()
	taxYear := taxYearFor(now)

	var existing models.PaymentTaxCertificate
	err := DB.Where("schedule_item_id = ? AND tax_year = ?", itemID, taxYear).First(&existing).Error
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.PaymentTaxCertificate{}, err
	}

	ref := generateCertificateRef(taxYear)
	dir := filepath.Join("data", "reports", "tax_certs", fmt.Sprintf("%d", taxYear))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return models.PaymentTaxCertificate{}, fmt.Errorf("failed to create tax-cert directory: %w", err)
	}
	fileName := fmt.Sprintf("%s.html", ref)
	storagePath := filepath.Join(dir, fileName)

	beneficiary := strings.TrimSpace(item.BeneficiaryName)
	if beneficiary == "" {
		beneficiary = item.MemberName
	}
	beneficiaryID := strings.TrimSpace(item.BeneficiaryIDNumber)
	if beneficiaryID == "" {
		beneficiaryID = item.MemberIDNumber
	}

	htmlBody := renderTaxCertificateHTML(taxCertContext{
		Ref:                 ref,
		ClaimNumber:         item.ClaimNumber,
		BenefitName:         item.BenefitName,
		SchemeName:          item.SchemeName,
		BeneficiaryName:     beneficiary,
		BeneficiaryIDNumber: beneficiaryID,
		GrossAmount:         item.GrossAmount,
		TaxWithheld:         item.TaxWithheld,
		NetPayable:          item.NetPayable,
		PaidAt:              now,
		TaxYear:             taxYear,
	})
	if err := os.WriteFile(storagePath, []byte(htmlBody), 0644); err != nil {
		return models.PaymentTaxCertificate{}, fmt.Errorf("failed to write tax cert: %w", err)
	}

	row := models.PaymentTaxCertificate{
		ScheduleID:          scheduleID,
		ScheduleItemID:      itemID,
		ClaimID:             item.ClaimID,
		ClaimNumber:         item.ClaimNumber,
		BenefitName:         item.BenefitName,
		BeneficiaryName:     beneficiary,
		BeneficiaryIDNumber: beneficiaryID,
		TaxYear:             taxYear,
		GrossAmount:         item.GrossAmount,
		TaxWithheld:         item.TaxWithheld,
		CertificateRef:      ref,
		FileName:            fileName,
		StoragePath:         storagePath,
		ContentType:         "text/html",
		GeneratedBy:         user.UserName,
	}
	if err := DB.Create(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

// ListTaxCertificatesForSchedule returns all certificates generated for lines
// on the given schedule.
func ListTaxCertificatesForSchedule(scheduleID int) ([]models.PaymentTaxCertificate, error) {
	var rows []models.PaymentTaxCertificate
	err := DB.Where("schedule_id = ?", scheduleID).Order("generated_at DESC").Find(&rows).Error
	return rows, err
}

// DownloadTaxCertificate returns the raw bytes for a certificate id.
func DownloadTaxCertificate(certID int) ([]byte, string, string, error) {
	var row models.PaymentTaxCertificate
	if err := DB.First(&row, certID).Error; err != nil {
		return nil, "", "", err
	}
	data, err := os.ReadFile(row.StoragePath)
	if err != nil {
		appLog.WithField("error", err.Error()).Warn("tax cert file missing on disk")
		return nil, "", "", fmt.Errorf("file not found: %w", err)
	}
	ct := row.ContentType
	if ct == "" {
		ct = "application/octet-stream"
	}
	return data, ct, row.FileName, nil
}

// ──────────────────────────────────────────────
// HTML renderer (replaceable seam for real PDF)
// ──────────────────────────────────────────────

type taxCertContext struct {
	Ref                 string
	ClaimNumber         string
	BenefitName         string
	SchemeName          string
	BeneficiaryName     string
	BeneficiaryIDNumber string
	GrossAmount         float64
	TaxWithheld         float64
	NetPayable          float64
	PaidAt              time.Time
	TaxYear             int
}

func renderTaxCertificateHTML(c taxCertContext) string {
	esc := html.EscapeString
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>IT3(a) Tax Certificate %s</title>
<style>
body { font-family: Arial, Helvetica, sans-serif; max-width: 720px; margin: 32px auto; color: #222; }
h1 { font-size: 22px; margin: 0 0 8px; }
.muted { color: #666; font-size: 12px; }
table { width: 100%%; border-collapse: collapse; margin: 24px 0; }
th, td { padding: 8px 10px; border: 1px solid #ddd; text-align: left; font-size: 14px; }
th { background: #f5f5f5; width: 40%%; }
.amount { text-align: right; font-variant-numeric: tabular-nums; }
.footer { margin-top: 32px; font-size: 11px; color: #555; line-height: 1.4; }
</style>
</head>
<body>
<h1>IT3(a) Tax Certificate</h1>
<div class="muted">Certificate reference: <strong>%s</strong> · Year of assessment: <strong>%d</strong></div>

<table>
<tr><th>Claim number</th><td>%s</td></tr>
<tr><th>Benefit</th><td>%s</td></tr>
<tr><th>Scheme</th><td>%s</td></tr>
<tr><th>Beneficiary</th><td>%s</td></tr>
<tr><th>ID number</th><td>%s</td></tr>
<tr><th>Payment date</th><td>%s</td></tr>
<tr><th>Gross amount</th><td class="amount">R %s</td></tr>
<tr><th>Tax withheld</th><td class="amount">R %s</td></tr>
<tr><th>Net paid</th><td class="amount">R %s</td></tr>
</table>

<div class="footer">
Issued in terms of the Income Tax Act, 1962. This certificate reflects the
gross amount, employees' tax (PAYE) withheld and net amount paid for the
year of assessment indicated above. Please attach this certificate to your
SARS submission.
</div>
</body>
</html>`,
		esc(c.Ref),
		esc(c.Ref),
		c.TaxYear,
		esc(c.ClaimNumber),
		esc(c.BenefitName),
		esc(c.SchemeName),
		esc(c.BeneficiaryName),
		esc(c.BeneficiaryIDNumber),
		c.PaidAt.Format("2 January 2006"),
		formatZAR(c.GrossAmount),
		formatZAR(c.TaxWithheld),
		formatZAR(c.NetPayable),
	)
}

func formatZAR(v float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", v), ",", " ")
}

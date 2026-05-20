package payment_letter

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"api/models"
	"api/services/docpdf"

	"gorm.io/gorm"
)

// ErrClaimNotPaid is returned when a caller asks for a letter on a claim that
// has not yet been marked paid via proof-of-payment upload. Letters are only
// meaningful once the funds have actually been disbursed.
var ErrClaimNotPaid = errors.New("claim has not been paid yet")

// GenerateAndRecord resolves claim + paid_at + settings, builds the letter in
// the requested format, persists a history row, and returns the rendered
// bytes plus filename. The version field is incremented per claim so multiple
// regenerations are individually traceable.
//
// `db` is the GORM handle (api/services.DB); we accept it as a parameter to
// keep this package free of the services.* circular import.
func GenerateAndRecord(db *gorm.DB, claimID int, format string, generatedBy string) (record models.ClaimPaymentLetter, filename string, data []byte, err error) {
	if format != models.PaymentLetterFormatDocx && format != models.PaymentLetterFormatPDF {
		return record, "", nil, fmt.Errorf("unsupported format %q (expected docx or pdf)", format)
	}

	var claim models.GroupSchemeClaim
	if err := db.First(&claim, claimID).Error; err != nil {
		return record, "", nil, err
	}
	if !strings.EqualFold(claim.Status, "paid") {
		return record, "", nil, ErrClaimNotPaid
	}

	paidAt, err := resolvePaidAt(db, claimID)
	if err != nil {
		return record, "", nil, fmt.Errorf("resolve paid_at: %w", err)
	}

	settings, err := loadSettings(db)
	if err != nil {
		return record, "", nil, fmt.Errorf("load letter settings: %w", err)
	}

	// Next version per claim.
	var maxVersion int
	db.Model(&models.ClaimPaymentLetter{}).
		Where("claim_id = ?", claimID).
		Select("COALESCE(MAX(version), 0)").Scan(&maxVersion)
	version := maxVersion + 1
	letterRef := fmt.Sprintf("CPL-%d-v%d", claimID, version)

	input := LetterInput{
		Claim:     claim,
		PaidAt:    paidAt,
		Settings:  settings,
		LetterRef: letterRef,
	}

	docxName, docxData, err := BuildLetterDocx(input)
	if err != nil {
		return record, "", nil, err
	}

	switch format {
	case models.PaymentLetterFormatDocx:
		filename = docxName
		data = docxData
	case models.PaymentLetterFormatPDF:
		pdfBytes, perr := docpdf.ConvertDocxToPdf(docxData)
		if perr != nil {
			return record, "", nil, perr
		}
		filename = strings.TrimSuffix(docxName, filepath.Ext(docxName)) + ".pdf"
		data = pdfBytes
	}

	record = models.ClaimPaymentLetter{
		ClaimID:           claimID,
		Version:           version,
		Format:            format,
		Filename:          filename,
		SizeBytes:         int64(len(data)),
		LetterReference:   letterRef,
		PaymentAmount:     claim.ClaimAmount,
		PaidAt:            paidAt,
		BankName:          claim.BankName,
		BankAccountNumber: claim.BankAccountNumber,
		AccountHolderName: claim.AccountHolderName,
		SettingsSnapshot:  SerialiseSettingsSnapshot(settings),
		GeneratedBy:       generatedBy,
		GeneratedAt:       time.Now(),
	}
	if err := db.Create(&record).Error; err != nil {
		return record, "", nil, fmt.Errorf("persist letter history: %w", err)
	}
	return record, filename, data, nil
}

// resolvePaidAt finds the timestamp when the claim was actually paid. Primary
// source is the ClaimPaymentProof that confirmed the schedule the claim was
// part of; falls back to the latest status audit row whose new_status='paid'.
// Returns time.Now() as a last resort so we never block letter generation.
func resolvePaidAt(db *gorm.DB, claimID int) (time.Time, error) {
	// 1. Try via proof of payment on the schedule.
	type row struct {
		UploadedAt time.Time
	}
	var r row
	err := db.Raw(`
		SELECT p.uploaded_at AS uploaded_at
		FROM claim_payment_proofs p
		JOIN claim_payment_schedule_items i ON i.schedule_id = p.schedule_id
		WHERE i.claim_id = ?
		ORDER BY p.uploaded_at DESC
		LIMIT 1
	`, claimID).Scan(&r).Error
	if err == nil && !r.UploadedAt.IsZero() {
		return r.UploadedAt, nil
	}

	// 2. Fall back to status audit.
	var audit models.GroupSchemeClaimStatusAudit
	if err := db.Where("claim_id = ? AND new_status = ?", claimID, "paid").
		Order("changed_at DESC").
		First(&audit).Error; err == nil && !audit.ChangedAt.IsZero() {
		return audit.ChangedAt, nil
	}

	// 3. Last resort.
	return time.Now(), nil
}

// loadSettings returns the singleton settings row, creating an empty one on
// first run so the admin UI has something to render.
func loadSettings(db *gorm.DB) (models.PaymentLetterSetting, error) {
	var s models.PaymentLetterSetting
	if err := db.First(&s, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s = models.PaymentLetterSetting{ID: 1}
			if err := db.Create(&s).Error; err != nil {
				return s, err
			}
			return s, nil
		}
		return s, err
	}
	return s, nil
}

// ListHistory returns all letters generated for a given claim, newest first.
func ListHistory(db *gorm.DB, claimID int) ([]models.ClaimPaymentLetter, error) {
	var rows []models.ClaimPaymentLetter
	err := db.Preload("Deliveries").
		Where("claim_id = ?", claimID).
		Order("version DESC").
		Find(&rows).Error
	return rows, err
}

// GenerateScheduleBundle builds a ZIP containing one letter (in the requested
// format) for every claim in a confirmed payment schedule. Fails fast if the
// PDF converter is missing on the host and PDF was requested.
func GenerateScheduleBundle(db *gorm.DB, scheduleID int, format string, generatedBy string) (filename string, data []byte, err error) {
	if format != models.PaymentLetterFormatDocx && format != models.PaymentLetterFormatPDF {
		return "", nil, fmt.Errorf("unsupported format %q (expected docx or pdf)", format)
	}

	var schedule models.ClaimPaymentSchedule
	if err := db.First(&schedule, scheduleID).Error; err != nil {
		return "", nil, err
	}
	if !strings.EqualFold(schedule.Status, "confirmed") && !strings.EqualFold(schedule.Status, "archived") {
		return "", nil, fmt.Errorf("schedule %d is %s; bundle is only available once the schedule is confirmed", scheduleID, schedule.Status)
	}

	var items []models.ClaimPaymentScheduleItem
	if err := db.Where("schedule_id = ?", scheduleID).Find(&items).Error; err != nil {
		return "", nil, err
	}
	if len(items) == 0 {
		return "", nil, fmt.Errorf("schedule %d has no items", scheduleID)
	}

	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	seen := map[string]int{}
	for _, item := range items {
		_, name, bytes, gerr := GenerateAndRecord(db, item.ClaimID, format, generatedBy)
		if gerr != nil {
			if errors.Is(gerr, ErrClaimNotPaid) {
				// Skip items whose claim is not yet paid; surface as a manifest entry.
				continue
			}
			return "", nil, fmt.Errorf("generate letter for claim %d: %w", item.ClaimID, gerr)
		}
		// Make the filename unique within the zip even if two claims share one.
		entryName := name
		if n, ok := seen[name]; ok {
			ext := filepath.Ext(name)
			base := strings.TrimSuffix(name, ext)
			entryName = fmt.Sprintf("%s_%d%s", base, n+1, ext)
			seen[name] = n + 1
		} else {
			seen[name] = 1
		}
		w, ferr := zw.Create(entryName)
		if ferr != nil {
			return "", nil, ferr
		}
		if _, ferr := w.Write(bytes); ferr != nil {
			return "", nil, ferr
		}
	}
	if err := zw.Close(); err != nil {
		return "", nil, err
	}

	zipName := fmt.Sprintf("PaymentLetters_%s_%s.zip",
		sanitiseFilenamePart(schedule.ScheduleNumber),
		time.Now().Format("2006-01-02"))
	return zipName, buf.Bytes(), nil
}

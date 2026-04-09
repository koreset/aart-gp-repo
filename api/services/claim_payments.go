package services

import (
	appLog "api/log"
	"api/models"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// generateScheduleNumber creates a unique schedule reference in the form PAY-YYYYMMDD-NNNN.
func generateScheduleNumber() string {
	var count int64
	DB.Model(&models.ClaimPaymentSchedule{}).Count(&count)
	return fmt.Sprintf("PAY-%s-%04d", time.Now().Format("20060102"), count+1)
}

// CreatePaymentScheduleRequest is the inbound payload for creating a new schedule.
type CreatePaymentScheduleRequest struct {
	ClaimIDs    []int  `json:"claim_ids"`
	Description string `json:"description"`
}

// CreatePaymentSchedule builds a payment schedule from a list of approved claim IDs.
// All qualifying claims are transitioned to "submitted_for_payment".
// Returns an error if no claim IDs are supplied or if any claim is not in an
// approvable state (approved / submitted_for_payment).
func CreatePaymentSchedule(req CreatePaymentScheduleRequest, user models.AppUser) (models.ClaimPaymentSchedule, error) {
	if len(req.ClaimIDs) == 0 {
		return models.ClaimPaymentSchedule{}, errors.New("at least one claim_id is required")
	}

	// Load the requested claims
	var claims []models.GroupSchemeClaim
	if err := DB.Where("id IN ?", req.ClaimIDs).Find(&claims).Error; err != nil {
		return models.ClaimPaymentSchedule{}, err
	}

	if len(claims) == 0 {
		return models.ClaimPaymentSchedule{}, errors.New("no valid claims found for the supplied IDs")
	}

	// Validate that every claim is in an approved state
	var items []models.ClaimPaymentScheduleItem
	var totalAmount float64
	for _, c := range claims {
		lower := strings.ToLower(c.Status)
		if lower != "approved" && lower != "submitted_for_payment" {
			return models.ClaimPaymentSchedule{}, fmt.Errorf("claim %s has status '%s' — only approved claims can be added to a payment schedule", c.ClaimNumber, c.Status)
		}
		totalAmount += c.ClaimAmount
		items = append(items, models.ClaimPaymentScheduleItem{
			ClaimID:           c.ID,
			ClaimNumber:       c.ClaimNumber,
			MemberName:        c.MemberName,
			MemberIDNumber:    c.MemberIDNumber,
			BenefitName:       c.BenefitName,
			SchemeName:        c.SchemeName,
			SchemeID:          c.SchemeId,
			ClaimAmount:       c.ClaimAmount,
			BankName:          c.BankName,
			BankBranchCode:    c.BankBranchCode,
			BankAccountNumber: c.BankAccountNumber,
			BankAccountType:   c.BankAccountType,
			AccountHolderName: c.AccountHolderName,
		})
	}

	schedule := models.ClaimPaymentSchedule{
		ScheduleNumber: generateScheduleNumber(),
		Description:    req.Description,
		Status:         "submitted",
		TotalAmount:    totalAmount,
		ClaimsCount:    len(items),
		CreatedBy:      user.UserName,
	}

	// Persist schedule + items in a transaction
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&schedule).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].ScheduleID = schedule.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// Transition each claim to "submitted_for_payment"
		for _, c := range claims {
			if strings.ToLower(c.Status) == "approved" {
				audit := models.GroupSchemeClaimStatusAudit{
					ClaimID:       c.ID,
					OldStatus:     c.Status,
					NewStatus:     "submitted_for_payment",
					StatusMessage: fmt.Sprintf("Included in payment schedule %s", schedule.ScheduleNumber),
					ChangedBy:     user.UserName,
					ChangedAt:     time.Now(),
				}
				if err := tx.Model(&models.GroupSchemeClaim{}).Where("id = ?", c.ID).Update("status", "submitted_for_payment").Error; err != nil {
					return err
				}
				if err := tx.Create(&audit).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return models.ClaimPaymentSchedule{}, err
	}

	return GetPaymentSchedule(schedule.ID)
}

// GetPaymentSchedules returns all payment schedules (newest first), including item counts.
func GetPaymentSchedules() ([]models.ClaimPaymentSchedule, error) {
	var schedules []models.ClaimPaymentSchedule
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Preload("Items").Preload("ProofOfPayments").
			Order("created_at DESC").Find(&schedules).Error
	})
	return schedules, err
}

// GetPaymentSchedule returns a single schedule by ID with all related data.
func GetPaymentSchedule(id int) (models.ClaimPaymentSchedule, error) {
	var schedule models.ClaimPaymentSchedule
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Preload("Items").Preload("ProofOfPayments").First(&schedule, id).Error
	})
	return schedule, err
}

// ExportPaymentScheduleCSV streams the payment schedule as a CSV file.
// Returns the CSV bytes and a suggested filename.
func ExportPaymentScheduleCSV(scheduleID int, user models.AppUser) ([]byte, string, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return nil, "", err
	}

	// Record export metadata
	now := time.Now()
	if dbErr := DB.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
		"exported_at": &now,
		"exported_by": user.UserName,
	}).Error; dbErr != nil {
		appLog.WithField("error", dbErr.Error()).Warn("Failed to record export timestamp on schedule")
	}

	pr, pw := io.Pipe()
	errCh := make(chan error, 1)
	go func() {
		defer pw.Close()
		w := csv.NewWriter(pw)
		// Header
		_ = w.Write([]string{
			"Schedule Number", "Claim Number", "Member Name", "ID Number",
			"Scheme", "Benefit", "Claim Amount",
		})
		for _, item := range schedule.Items {
			_ = w.Write([]string{
				schedule.ScheduleNumber,
				item.ClaimNumber,
				item.MemberName,
				item.MemberIDNumber,
				item.SchemeName,
				item.BenefitName,
				fmt.Sprintf("%.2f", item.ClaimAmount),
			})
		}
		w.Flush()
		errCh <- w.Error()
	}()

	data, readErr := io.ReadAll(pr)
	writeErr := <-errCh
	if writeErr != nil {
		return nil, "", writeErr
	}
	if readErr != nil {
		return nil, "", readErr
	}

	filename := fmt.Sprintf("payment_schedule_%s.csv", schedule.ScheduleNumber)
	return data, filename, nil
}

// validatePaymentCSV checks that a CSV file uploaded as proof of payment is
// structurally consistent with the payment schedule it is being attached to.
// It requires:
//   - "Claim Number" and "Claim Amount" columns (case-insensitive header match)
//   - Every claim number in the schedule is present in the CSV
//   - Each matched claim's amount matches the schedule item amount (±0.01 tolerance)
func validatePaymentCSV(f multipart.File, schedule models.ClaimPaymentSchedule) error {
	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	headers, err := r.Read()
	if err != nil {
		return fmt.Errorf("could not read CSV headers: %w", err)
	}

	claimNumIdx, claimAmtIdx := -1, -1
	for i, h := range headers {
		switch strings.ToLower(strings.TrimSpace(h)) {
		case "claim number", "claim_number", "claimnumber":
			claimNumIdx = i
		case "claim amount", "claim_amount", "claimamount":
			claimAmtIdx = i
		}
	}
	if claimNumIdx == -1 {
		return errors.New("CSV validation failed: required column 'Claim Number' not found")
	}
	if claimAmtIdx == -1 {
		return errors.New("CSV validation failed: required column 'Claim Amount' not found")
	}

	// Build a map of claim number → expected amount from the schedule
	expected := make(map[string]float64, len(schedule.Items))
	for _, item := range schedule.Items {
		expected[item.ClaimNumber] = item.ClaimAmount
	}

	found := make(map[string]bool)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("CSV parse error: %w", err)
		}
		if claimNumIdx >= len(row) || claimAmtIdx >= len(row) {
			continue
		}
		claimNum := strings.TrimSpace(row[claimNumIdx])
		if claimNum == "" {
			continue
		}
		expectedAmt, inSchedule := expected[claimNum]
		if !inSchedule {
			continue // extra rows are allowed; we only validate what's in the schedule
		}
		var csvAmt float64
		if _, scanErr := fmt.Sscanf(strings.ReplaceAll(row[claimAmtIdx], ",", ""), "%f", &csvAmt); scanErr != nil {
			return fmt.Errorf("CSV validation failed: could not parse amount for claim %s", claimNum)
		}
		diff := csvAmt - expectedAmt
		if diff < 0 {
			diff = -diff
		}
		if diff > 0.01 {
			return fmt.Errorf("CSV validation failed: claim %s — expected amount %.2f but CSV contains %.2f", claimNum, expectedAmt, csvAmt)
		}
		found[claimNum] = true
	}

	// Every claim in the schedule must appear in the CSV
	for claimNum := range expected {
		if !found[claimNum] {
			return fmt.Errorf("CSV validation failed: claim %s is in the payment schedule but not found in the uploaded CSV", claimNum)
		}
	}
	return nil
}

// UploadPaymentProof saves a proof-of-payment file for a schedule, then confirms
// the schedule and moves all its claims to "paid".
// If the uploaded file is a CSV, its contents are validated against the schedule
// before the schedule is confirmed.
func UploadPaymentProof(scheduleID int, fileHeader *multipart.FileHeader, notes string, user models.AppUser) (models.ClaimPaymentProof, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return models.ClaimPaymentProof{}, err
	}
	if schedule.Status == "confirmed" {
		return models.ClaimPaymentProof{}, errors.New("this schedule is already confirmed — payment has been processed")
	}

	// Validate CSV contents before accepting the upload
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	contentType := strings.ToLower(fileHeader.Header.Get("Content-Type"))
	isCSV := ext == ".csv" || strings.Contains(contentType, "csv") || strings.Contains(contentType, "text/plain")
	if isCSV {
		csvFile, openErr := fileHeader.Open()
		if openErr != nil {
			return models.ClaimPaymentProof{}, fmt.Errorf("failed to open uploaded file for validation: %w", openErr)
		}
		validationErr := validatePaymentCSV(csvFile, schedule)
		csvFile.Close()
		if validationErr != nil {
			return models.ClaimPaymentProof{}, validationErr
		}
	}

	// Persist the file to disk
	uploadDir := filepath.Join("tmp", "uploads", "payment_proofs", fmt.Sprintf("schedule_%d", scheduleID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return models.ClaimPaymentProof{}, fmt.Errorf("failed to create upload directory: %w", err)
	}
	safeName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
	storagePath := filepath.Join(uploadDir, safeName)

	src, err := fileHeader.Open()
	if err != nil {
		return models.ClaimPaymentProof{}, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(storagePath)
	if err != nil {
		return models.ClaimPaymentProof{}, fmt.Errorf("failed to save file: %w", err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return models.ClaimPaymentProof{}, fmt.Errorf("failed to write file: %w", err)
	}

	proof := models.ClaimPaymentProof{
		ScheduleID:  scheduleID,
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		SizeBytes:   fileHeader.Size,
		StoragePath: storagePath,
		Notes:       notes,
		UploadedBy:  user.UserName,
	}

	// Confirm the schedule and mark all claims as paid in a transaction
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&proof).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Update("status", "confirmed").Error; err != nil {
			return err
		}
		for _, item := range schedule.Items {
			// Fetch current status to record the audit accurately
			var current models.GroupSchemeClaim
			if err := tx.Select("id, status").First(&current, item.ClaimID).Error; err != nil {
				continue
			}
			if strings.ToLower(current.Status) == "paid" {
				continue // already paid, skip
			}
			audit := models.GroupSchemeClaimStatusAudit{
				ClaimID:       item.ClaimID,
				OldStatus:     current.Status,
				NewStatus:     "paid",
				StatusMessage: fmt.Sprintf("Payment confirmed via schedule %s", schedule.ScheduleNumber),
				ChangedBy:     user.UserName,
				ChangedAt:     time.Now(),
			}
			if err := tx.Model(&models.GroupSchemeClaim{}).Where("id = ?", item.ClaimID).Update("status", "paid").Error; err != nil {
				return err
			}
			if err := tx.Create(&audit).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if txErr != nil {
		return models.ClaimPaymentProof{}, txErr
	}

	return proof, nil
}

// GetPaymentProofs returns all proof-of-payment records for a given schedule.
func GetPaymentProofs(scheduleID int) ([]models.ClaimPaymentProof, error) {
	var proofs []models.ClaimPaymentProof
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("schedule_id = ?", scheduleID).Order("uploaded_at DESC").Find(&proofs).Error
	})
	return proofs, err
}

// DownloadPaymentProof returns the raw file bytes and content type for a proof record.
func DownloadPaymentProof(proofID int) ([]byte, string, string, error) {
	var proof models.ClaimPaymentProof
	if err := DB.First(&proof, proofID).Error; err != nil {
		return nil, "", "", err
	}
	data, err := os.ReadFile(proof.StoragePath)
	if err != nil {
		return nil, "", "", fmt.Errorf("file not found on disk: %w", err)
	}
	contentType := proof.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return data, contentType, proof.FileName, nil
}

// ──────────────────────────────────────────────
// Bank Profile CRUD
// ──────────────────────────────────────────────

// CreateBankProfile creates a new ACB bank profile.
func CreateBankProfile(req models.CreateBankProfileRequest, user models.AppUser) (models.ACBBankProfile, error) {
	profile := models.ACBBankProfile{
		ProfileName:       req.ProfileName,
		BankName:          req.BankName,
		UserCode:          req.UserCode,
		UserBranchCode:    req.UserBranchCode,
		UserAccountNumber: req.UserAccountNumber,
		UserAccountType:   req.UserAccountType,
		BankTypeCode:      req.BankTypeCode,
		ServiceType:       req.ServiceType,
		GenerationNumber:  0,
		IsActive:          true,
		CreatedBy:         user.UserName,
	}
	if err := DB.Create(&profile).Error; err != nil {
		return models.ACBBankProfile{}, err
	}
	return profile, nil
}

// GetBankProfiles returns all bank profiles.
func GetBankProfiles() ([]models.ACBBankProfile, error) {
	var profiles []models.ACBBankProfile
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Order("created_at DESC").Find(&profiles).Error
	})
	return profiles, err
}

// GetBankProfile returns a single bank profile by ID.
func GetBankProfile(id int) (models.ACBBankProfile, error) {
	var profile models.ACBBankProfile
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.First(&profile, id).Error
	})
	return profile, err
}

// UpdateBankProfile updates a bank profile.
func UpdateBankProfile(id int, req models.UpdateBankProfileRequest) (models.ACBBankProfile, error) {
	var profile models.ACBBankProfile
	if err := DB.First(&profile, id).Error; err != nil {
		return models.ACBBankProfile{}, err
	}
	updates := make(map[string]interface{})
	if req.ProfileName != nil {
		updates["profile_name"] = *req.ProfileName
	}
	if req.BankName != nil {
		updates["bank_name"] = *req.BankName
	}
	if req.UserCode != nil {
		updates["user_code"] = *req.UserCode
	}
	if req.UserBranchCode != nil {
		updates["user_branch_code"] = *req.UserBranchCode
	}
	if req.UserAccountNumber != nil {
		updates["user_account_number"] = *req.UserAccountNumber
	}
	if req.UserAccountType != nil {
		updates["user_account_type"] = *req.UserAccountType
	}
	if req.BankTypeCode != nil {
		updates["bank_type_code"] = *req.BankTypeCode
	}
	if req.ServiceType != nil {
		updates["service_type"] = *req.ServiceType
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if len(updates) > 0 {
		if err := DB.Model(&profile).Updates(updates).Error; err != nil {
			return models.ACBBankProfile{}, err
		}
	}
	return GetBankProfile(id)
}

// DeleteBankProfile deletes a bank profile by ID.
func DeleteBankProfile(id int) error {
	return DB.Delete(&models.ACBBankProfile{}, id).Error
}

// ──────────────────────────────────────────────
// ACB File Generation
// ──────────────────────────────────────────────

// GenerateACBFile generates an ACB file for a payment schedule.
func GenerateACBFile(scheduleID int, req models.GenerateACBRequest, user models.AppUser) (models.ACBFileRecord, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return models.ACBFileRecord{}, err
	}

	// Validate all items have banking details
	var missing []string
	for _, item := range schedule.Items {
		if item.BankAccountNumber == "" || item.BankBranchCode == "" {
			missing = append(missing, item.ClaimNumber)
		}
	}
	if len(missing) > 0 {
		return models.ACBFileRecord{}, fmt.Errorf("claims missing banking details: %s", strings.Join(missing, ", "))
	}

	profile, err := GetBankProfile(req.BankProfileID)
	if err != nil {
		return models.ACBFileRecord{}, fmt.Errorf("bank profile not found: %w", err)
	}

	// Increment generation number
	profile.GenerationNumber++
	DB.Model(&profile).Update("generation_number", profile.GenerationNumber)

	actionDate, err := time.Parse("2006-01-02", req.ActionDate)
	if err != nil {
		return models.ACBFileRecord{}, fmt.Errorf("invalid action_date format (expected YYYY-MM-DD): %w", err)
	}

	content, hashTotal, err := GenerateACBFileContent(profile, schedule.Items, actionDate)
	if err != nil {
		return models.ACBFileRecord{}, err
	}

	// Write to disk
	outputDir := filepath.Join("data", "reports", "acb")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return models.ACBFileRecord{}, fmt.Errorf("failed to create ACB output directory: %w", err)
	}
	fileName := fmt.Sprintf("%s_%s_%04d.txt", schedule.ScheduleNumber, actionDate.Format("20060102"), profile.GenerationNumber)
	filePath := filepath.Join(outputDir, fileName)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return models.ACBFileRecord{}, fmt.Errorf("failed to write ACB file: %w", err)
	}

	// Create audit record
	fileRecord := models.ACBFileRecord{
		ScheduleID:       scheduleID,
		BankProfileID:    profile.ID,
		FileName:         fileName,
		FilePath:         filePath,
		ActionDate:       actionDate.Format("060102"),
		TransactionCount: len(schedule.Items),
		TotalAmount:      schedule.TotalAmount,
		HashTotal:        hashTotal,
		GenerationNumber: profile.GenerationNumber,
		Status:           "generated",
		GeneratedBy:      user.UserName,
	}
	if err := DB.Create(&fileRecord).Error; err != nil {
		return models.ACBFileRecord{}, err
	}

	// Update schedule ACB tracking
	now := time.Now()
	DB.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
		"acb_file_generated": true,
		"acb_generated_at":   &now,
		"acb_generated_by":   user.UserName,
		"bank_profile_id":    &profile.ID,
	})

	return fileRecord, nil
}

// DownloadACBFile reads an ACB file from disk.
func DownloadACBFile(acbFileID int) ([]byte, string, error) {
	var record models.ACBFileRecord
	if err := DB.First(&record, acbFileID).Error; err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(record.FilePath)
	if err != nil {
		return nil, "", fmt.Errorf("ACB file not found on disk: %w", err)
	}
	return data, record.FileName, nil
}

// GetACBFileRecords returns all ACB file records for a schedule.
func GetACBFileRecords(scheduleID int) ([]models.ACBFileRecord, error) {
	var records []models.ACBFileRecord
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Preload("BankProfile").Where("schedule_id = ?", scheduleID).Order("generated_at DESC").Find(&records).Error
	})
	return records, err
}

// ──────────────────────────────────────────────
// Bank Response Processing & Reconciliation
// ──────────────────────────────────────────────

// ProcessBankResponse parses a bank response file and reconciles against schedule items.
func ProcessBankResponse(acbFileID int, fileHeader *multipart.FileHeader, user models.AppUser) (models.ACBReconciliationSummary, error) {
	var fileRecord models.ACBFileRecord
	if err := DB.Preload("Schedule.Items").First(&fileRecord, acbFileID).Error; err != nil {
		return models.ACBReconciliationSummary{}, fmt.Errorf("ACB file record not found: %w", err)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return models.ACBReconciliationSummary{}, fmt.Errorf("failed to open response file: %w", err)
	}
	defer src.Close()

	rawData, err := io.ReadAll(src)
	if err != nil {
		return models.ACBReconciliationSummary{}, fmt.Errorf("failed to read response file: %w", err)
	}

	// Build lookup map from schedule items: key = accountNumber|amountCents
	type itemRef struct {
		Item  models.ClaimPaymentScheduleItem
		Found bool
	}
	itemMap := make(map[string]*itemRef)
	for _, item := range fileRecord.Schedule.Items {
		key := fmt.Sprintf("%s|%d", strings.TrimLeft(item.BankAccountNumber, "0"), amountToCents(item.ClaimAmount))
		itemMap[key] = &itemRef{Item: item}
	}

	// Parse response — detect format by extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var rows []acbResponseRow

	if ext == ".csv" {
		rows, err = parseCSVResponse(rawData)
	} else {
		rows, err = parseACBResponse(rawData)
	}
	if err != nil {
		return models.ACBReconciliationSummary{}, err
	}

	// Reconcile
	var results []models.ACBReconciliationResult
	summary := models.ACBReconciliationSummary{}

	for _, row := range rows {
		key := fmt.Sprintf("%s|%d", strings.TrimLeft(row.AccountNumber, "0"), row.AmountCents)
		ref, matched := itemMap[key]

		result := models.ACBReconciliationResult{
			ACBFileID:     acbFileID,
			AccountNumber: row.AccountNumber,
			Amount:        float64(row.AmountCents) / 100.0,
			BankReference: row.Reference,
			ResponseCode:  row.Status,
		}

		if !matched {
			result.Status = "unmatched"
			summary.Unmatched++
		} else {
			ref.Found = true
			result.ScheduleItemID = ref.Item.ID
			result.ClaimID = ref.Item.ClaimID
			result.ClaimNumber = ref.Item.ClaimNumber

			isPaid := strings.EqualFold(row.Status, "paid") || strings.EqualFold(row.Status, "accepted") ||
				strings.EqualFold(row.Status, "success") || strings.EqualFold(row.Status, "0") || row.Status == ""
			if isPaid {
				result.Status = "paid"
				summary.Paid++
				summary.TotalPaid += result.Amount
			} else {
				result.Status = "failed"
				result.FailureReason = row.Reason
				summary.Failed++
				summary.TotalFailed += result.Amount
			}
		}
		results = append(results, result)
	}
	summary.TotalTransactions = len(rows)

	// Persist results and update claim statuses in a transaction
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		if len(results) > 0 {
			if err := tx.Create(&results).Error; err != nil {
				return err
			}
		}
		for _, r := range results {
			if r.ClaimID == 0 {
				continue
			}
			var newStatus string
			var msg string
			switch r.Status {
			case "paid":
				newStatus = "paid"
				msg = fmt.Sprintf("Payment confirmed via ACB reconciliation (ref: %s)", r.BankReference)
			case "failed":
				newStatus = "payment_failed"
				msg = fmt.Sprintf("ACB payment failed: %s", r.FailureReason)
			default:
				continue
			}
			var current models.GroupSchemeClaim
			if err := tx.Select("id, status").First(&current, r.ClaimID).Error; err != nil {
				continue
			}
			audit := models.GroupSchemeClaimStatusAudit{
				ClaimID:       r.ClaimID,
				OldStatus:     current.Status,
				NewStatus:     newStatus,
				StatusMessage: msg,
				ChangedBy:     user.UserName,
				ChangedAt:     time.Now(),
			}
			tx.Model(&models.GroupSchemeClaim{}).Where("id = ?", r.ClaimID).Update("status", newStatus)
			tx.Create(&audit)
		}
		now := time.Now()
		tx.Model(&models.ACBFileRecord{}).Where("id = ?", acbFileID).Updates(map[string]interface{}{
			"status":        "reconciled",
			"reconciled_at": &now,
			"reconciled_by": user.UserName,
		})
		return nil
	})
	if txErr != nil {
		return models.ACBReconciliationSummary{}, txErr
	}

	go NotifyClaimPaymentSummary(summary, user)
	return summary, nil
}

// acbResponseRow represents a single parsed row from a bank response file.
type acbResponseRow struct {
	AccountNumber string
	AmountCents   int64
	Status        string
	Reason        string
	Reference     string
}

// parseCSVResponse parses a CSV bank response file.
func parseCSVResponse(data []byte) ([]acbResponseRow, error) {
	r := csv.NewReader(strings.NewReader(string(data)))
	r.TrimLeadingSpace = true
	headers, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read CSV headers: %w", err)
	}

	// Find column indices (case-insensitive)
	acctIdx, amtIdx, statusIdx, reasonIdx, refIdx := -1, -1, -1, -1, -1
	for i, h := range headers {
		h = strings.ToLower(strings.TrimSpace(h))
		switch {
		case h == "account_number" || h == "account" || h == "accountnumber":
			acctIdx = i
		case h == "amount":
			amtIdx = i
		case h == "status" || h == "result":
			statusIdx = i
		case h == "reason" || h == "failure_reason":
			reasonIdx = i
		case h == "reference" || h == "bank_reference":
			refIdx = i
		}
	}
	if acctIdx == -1 || amtIdx == -1 {
		return nil, errors.New("CSV response must have 'account_number'/'account' and 'amount' columns")
	}

	var rows []acbResponseRow
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if acctIdx >= len(row) || amtIdx >= len(row) {
			continue
		}

		var amt float64
		amtStr := strings.ReplaceAll(strings.TrimSpace(row[amtIdx]), ",", "")
		fmt.Sscanf(amtStr, "%f", &amt)

		entry := acbResponseRow{
			AccountNumber: strings.TrimSpace(row[acctIdx]),
			AmountCents:   amountToCents(amt),
		}
		if statusIdx >= 0 && statusIdx < len(row) {
			entry.Status = strings.TrimSpace(row[statusIdx])
		}
		if reasonIdx >= 0 && reasonIdx < len(row) {
			entry.Reason = strings.TrimSpace(row[reasonIdx])
		}
		if refIdx >= 0 && refIdx < len(row) {
			entry.Reference = strings.TrimSpace(row[refIdx])
		}
		rows = append(rows, entry)
	}

	return rows, nil
}

// parseACBResponse parses an ACB-format bank response file (fixed-width).
func parseACBResponse(data []byte) ([]acbResponseRow, error) {
	var rows []acbResponseRow
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		if len(line) < 38 {
			continue
		}
		recType := line[0:2]
		if recType != "50" {
			continue
		}
		// Extract fields from fixed positions
		acctNum := strings.TrimLeft(line[8:19], "0 ")
		amtStr := strings.TrimLeft(line[26:38], "0 ")
		var amtCents int64
		fmt.Sscanf(amtStr, "%d", &amtCents)

		var status, reason, ref string
		// If the response has extra fields beyond standard 107 chars
		if len(line) > 107 {
			respCode := strings.TrimSpace(line[107:min(111, len(line))])
			if respCode != "" && respCode != "0" && respCode != "00" {
				status = "failed"
				reason = fmt.Sprintf("response code: %s", respCode)
			}
		}
		if status == "" {
			status = "paid"
		}
		if len(line) > 77 {
			ref = strings.TrimSpace(line[47:min(77, len(line))])
		}

		rows = append(rows, acbResponseRow{
			AccountNumber: acctNum,
			AmountCents:   amtCents,
			Status:        status,
			Reason:        reason,
			Reference:     ref,
		})
	}

	return rows, nil
}

// GetReconciliationResults returns all reconciliation results for an ACB file.
func GetACBReconciliationResults(acbFileID int) ([]models.ACBReconciliationResult, error) {
	var results []models.ACBReconciliationResult
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("acb_file_id = ?", acbFileID).Order("id ASC").Find(&results).Error
	})
	return results, err
}

// GetReconciliationSummary returns the reconciliation summary for a schedule.
func GetACBReconciliationSummary(scheduleID int) (models.ACBReconciliationSummary, error) {
	// Get all ACB files for this schedule
	var fileRecords []models.ACBFileRecord
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("schedule_id = ?", scheduleID).Find(&fileRecords).Error
	})
	if err != nil {
		return models.ACBReconciliationSummary{}, err
	}

	var fileIDs []int
	for _, f := range fileRecords {
		fileIDs = append(fileIDs, f.ID)
	}
	if len(fileIDs) == 0 {
		return models.ACBReconciliationSummary{}, nil
	}

	var results []models.ACBReconciliationResult
	err = DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Where("acb_file_id IN ?", fileIDs).Find(&results).Error
	})
	if err != nil {
		return models.ACBReconciliationSummary{}, err
	}

	summary := models.ACBReconciliationSummary{TotalTransactions: len(results)}
	for _, r := range results {
		switch r.Status {
		case "paid":
			summary.Paid++
			summary.TotalPaid += r.Amount
		case "failed":
			summary.Failed++
			summary.TotalFailed += r.Amount
		default:
			summary.Unmatched++
		}
	}
	return summary, nil
}

// ──────────────────────────────────────────────
// Retry Failed Payments
// ──────────────────────────────────────────────

// RetryFailedPayments generates a new ACB file for failed items from a previous reconciliation.
func RetryFailedPayments(acbFileID int, req models.RetryFailedRequest, user models.AppUser) (models.ACBFileRecord, error) {
	// Get failed reconciliation results
	var failedResults []models.ACBReconciliationResult
	query := DB.Where("acb_file_id = ? AND status = ?", acbFileID, "failed")
	if len(req.ItemIDs) > 0 {
		query = query.Where("schedule_item_id IN ?", req.ItemIDs)
	}
	if err := query.Find(&failedResults).Error; err != nil {
		return models.ACBFileRecord{}, err
	}
	if len(failedResults) == 0 {
		return models.ACBFileRecord{}, errors.New("no failed items found to retry")
	}

	// Get original ACB file record for schedule and profile info
	var origFile models.ACBFileRecord
	if err := DB.First(&origFile, acbFileID).Error; err != nil {
		return models.ACBFileRecord{}, err
	}

	// Collect schedule item IDs
	var itemIDs []int
	for _, r := range failedResults {
		itemIDs = append(itemIDs, r.ScheduleItemID)
	}

	// Load schedule items
	var items []models.ClaimPaymentScheduleItem
	if err := DB.Where("id IN ?", itemIDs).Find(&items).Error; err != nil {
		return models.ACBFileRecord{}, err
	}

	// Reset claims from payment_failed → submitted_for_payment
	for _, item := range items {
		var current models.GroupSchemeClaim
		if err := DB.Select("id, status").First(&current, item.ClaimID).Error; err != nil {
			continue
		}
		if current.Status == "payment_failed" {
			audit := models.GroupSchemeClaimStatusAudit{
				ClaimID:       item.ClaimID,
				OldStatus:     "payment_failed",
				NewStatus:     "submitted_for_payment",
				StatusMessage: "Retrying ACB payment",
				ChangedBy:     user.UserName,
				ChangedAt:     time.Now(),
			}
			DB.Model(&models.GroupSchemeClaim{}).Where("id = ?", item.ClaimID).Update("status", "submitted_for_payment")
			DB.Create(&audit)
		}
	}

	// Generate new ACB file
	profile, err := GetBankProfile(origFile.BankProfileID)
	if err != nil {
		return models.ACBFileRecord{}, err
	}
	profile.GenerationNumber++
	DB.Model(&profile).Update("generation_number", profile.GenerationNumber)

	actionDate := time.Now().AddDate(0, 0, 1) // next business day
	content, hashTotal, err := GenerateACBFileContent(profile, items, actionDate)
	if err != nil {
		return models.ACBFileRecord{}, err
	}

	// Write to disk
	outputDir := filepath.Join("data", "reports", "acb")
	os.MkdirAll(outputDir, 0755)

	schedule, _ := GetPaymentSchedule(origFile.ScheduleID)
	fileName := fmt.Sprintf("%s_%s_%04d_retry.txt", schedule.ScheduleNumber, actionDate.Format("20060102"), profile.GenerationNumber)
	filePath := filepath.Join(outputDir, fileName)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return models.ACBFileRecord{}, fmt.Errorf("failed to write retry ACB file: %w", err)
	}

	var totalAmount float64
	for _, item := range items {
		totalAmount += item.ClaimAmount
	}

	fileRecord := models.ACBFileRecord{
		ScheduleID:       origFile.ScheduleID,
		BankProfileID:    profile.ID,
		FileName:         fileName,
		FilePath:         filePath,
		ActionDate:       actionDate.Format("060102"),
		TransactionCount: len(items),
		TotalAmount:      totalAmount,
		HashTotal:        hashTotal,
		GenerationNumber: profile.GenerationNumber,
		Status:           "generated",
		IsRetry:          true,
		GeneratedBy:      user.UserName,
	}
	if err := DB.Create(&fileRecord).Error; err != nil {
		return models.ACBFileRecord{}, err
	}

	return fileRecord, nil
}

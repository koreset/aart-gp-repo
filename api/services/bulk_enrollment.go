package services

import (
	"api/models"
	"api/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RowValidationResult is a single per-row finding from validateBatch. The
// approver UI renders these next to each member row in the batch review.
type RowValidationResult struct {
	Row            int    `json:"row"`
	Field          string `json:"field"`
	Severity       string `json:"severity"` // "blocking" | "soft"
	Message        string `json:"message"`
	MemberIdNumber string `json:"member_id_number,omitempty"`
	MemberName     string `json:"member_name,omitempty"`
}

// BulkEnrollmentBatchFileMeta describes the source file the uploader picked.
// Captured on the batch row for audit and de-dupe.
type BulkEnrollmentBatchFileMeta struct {
	FileName      string
	FileSizeBytes int64
	FileChecksum  string
}

// CreateBulkEnrollmentBatch persists an upload as a pending_approval batch.
//
// Members are inserted regardless of validation outcome so the approver can
// inspect every row in the review screen. Rows carry status='draft' and a
// per-row ValidationStatus (valid | soft_error | blocking_error) until the
// batch is approved (drafts flip to 'active') or rejected (drafts deleted).
//
// The whole operation runs in a single transaction; audit row written via
// the shared writeAudit helper.
func CreateBulkEnrollmentBatch(
	schemeID int,
	members []models.GPricingMemberDataInForce,
	skipDuplicates bool,
	file BulkEnrollmentBatchFileMeta,
	user models.AppUser,
) (models.BulkEnrollmentBatch, []RowValidationResult, error) {
	var batch models.BulkEnrollmentBatch
	if len(members) == 0 {
		return batch, nil, fmt.Errorf("no members provided")
	}

	// The scheme's active quote provides the scheme categories and the
	// UseGlobalSalaryMultiple flag that drive validation rules. Existing
	// services use scheme_quote_status='in_effect' (NOT 'in_force') to
	// identify the active quote — see services.AddMembersToSchemeBulk.
	var quote models.GroupPricingQuote
	if err := DB.Preload("SchemeCategories").
		Where("scheme_id = ? AND scheme_quote_status = ?", schemeID, models.StatusInEffect).
		Order("commencement_date desc, id desc").
		First(&quote).Error; err != nil {
		return batch, nil, fmt.Errorf("no active quote for scheme %d: %w", schemeID, err)
	}

	categoryByName := make(map[string]models.SchemeCategory, len(quote.SchemeCategories))
	for _, c := range quote.SchemeCategories {
		categoryByName[strings.TrimSpace(c.SchemeCategory)] = c
	}

	var occupations []models.OccupationClass
	if err := DB.Find(&occupations).Error; err != nil {
		return batch, nil, fmt.Errorf("failed to load occupations: %w", err)
	}
	occupationLookup := make(map[string]struct{}, len(occupations)*2)
	for _, o := range occupations {
		if o.RiskRateCode != "" {
			occupationLookup[strings.ToLower(strings.TrimSpace(o.RiskRateCode))] = struct{}{}
		}
		if o.Category != "" {
			occupationLookup[strings.ToLower(strings.TrimSpace(o.Category))] = struct{}{}
		}
	}

	var existingIDs []string
	if err := DB.Model(&models.GPricingMemberDataInForce{}).
		Where("scheme_id = ? AND status <> ?", schemeID, models.BulkEnrollmentMemberDraftStatus).
		Pluck("member_id_number", &existingIDs).Error; err != nil {
		return batch, nil, fmt.Errorf("failed to load existing members: %w", err)
	}
	existingIDSet := make(map[string]struct{}, len(existingIDs))
	for _, id := range existingIDs {
		existingIDSet[id] = struct{}{}
	}

	results := validateBatch(members, &quote, categoryByName, occupationLookup, existingIDSet, skipDuplicates)

	// Counts feed the batch summary and the Pending Approvals badge.
	var validCount, blockingCount, softCount int
	rowSeverity := make(map[int]string, len(members))
	for _, r := range results {
		switch r.Severity {
		case "blocking":
			if rowSeverity[r.Row] != "blocking" {
				if rowSeverity[r.Row] == "soft" {
					softCount-- // re-classify
				}
				blockingCount++
				rowSeverity[r.Row] = "blocking"
			}
		case "soft":
			if rowSeverity[r.Row] == "" {
				softCount++
				rowSeverity[r.Row] = "soft"
			}
		}
	}
	for i := range members {
		if _, hit := rowSeverity[i+1]; !hit {
			validCount++
		}
	}

	reportJSON, _ := json.Marshal(results)

	batch = models.BulkEnrollmentBatch{
		SchemeID:         schemeID,
		QuoteID:          quote.ID,
		Status:           models.BulkEnrollmentBatchPendingApproval,
		MemberCount:      len(members),
		ValidCount:       validCount,
		BlockingCount:    blockingCount,
		SoftErrorCount:   softCount,
		FileName:         file.FileName,
		FileSizeBytes:    file.FileSizeBytes,
		FileChecksum:     file.FileChecksum,
		SkipDuplicates:   skipDuplicates,
		ValidationReport: string(reportJSON),
		UploadedBy:       user.UserName,
		UploadedAt:       time.Now(),
	}

	if err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&batch).Error; err != nil {
			return err
		}

		prepared := make([]models.GPricingMemberDataInForce, 0, len(members))
		for i, m := range members {
			rowNum := i + 1
			sev := rowSeverity[rowNum]

			m.SchemeId = schemeID
			m.QuoteId = quote.ID
			m.BatchID = batch.ID
			m.RowIndex = rowNum
			m.Status = models.BulkEnrollmentMemberDraftStatus
			m.CreatedBy = user.UserName
			m.CreationDate = time.Now()
			if m.Year == 0 {
				m.Year = time.Now().Year()
			}
			switch sev {
			case "blocking":
				m.ValidationStatus = models.BulkEnrollmentRowBlockingError
			case "soft":
				m.ValidationStatus = models.BulkEnrollmentRowSoftError
			default:
				m.ValidationStatus = models.BulkEnrollmentRowValid
			}
			prepared = append(prepared, m)
		}

		if err := tx.CreateInBatches(&prepared, 100).Error; err != nil {
			return err
		}

		return writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "bulk_enrollment_batches",
			EntityID:  strconv.Itoa(batch.ID),
			Action:    "CREATE",
			ChangedBy: user.UserName,
		}, struct{}{}, batch)
	}); err != nil {
		return batch, results, err
	}

	return batch, results, nil
}

// validateBatch runs the per-row checks the approval gate depends on. All
// failures are blocking — there's nothing soft in the spec — but the result
// type carries Severity to leave room for soft warnings in future.
func validateBatch(
	members []models.GPricingMemberDataInForce,
	quote *models.GroupPricingQuote,
	categoryByName map[string]models.SchemeCategory,
	occupationLookup map[string]struct{},
	existingIDSet map[string]struct{},
	skipDuplicates bool,
) []RowValidationResult {
	var results []RowValidationResult
	seenInBatch := make(map[string]int, len(members))

	for i, m := range members {
		row := i + 1
		add := func(field, msg string) {
			results = append(results, RowValidationResult{
				Row:            row,
				Field:          field,
				Severity:       "blocking",
				Message:        msg,
				MemberIdNumber: m.MemberIdNumber,
				MemberName:     m.MemberName,
			})
		}

		idNumber := strings.TrimSpace(m.MemberIdNumber)
		idType := strings.TrimSpace(m.MemberIdType)

		if idNumber == "" {
			add("member_id_number", "member_id_number is required")
		}
		if idType == "" {
			add("member_id_type", "member_id_type is required")
		}
		if strings.TrimSpace(m.Gender) == "" {
			add("gender", "gender is required")
		}
		if m.DateOfBirth.IsZero() {
			add("date_of_birth", "date_of_birth is required")
		}
		if m.AnnualSalary <= 0 {
			add("annual_salary", "annual_salary must be greater than zero")
		}

		if idNumber != "" && isRSAIDType(idType) {
			if !utils.IsValidRSAID(idNumber) {
				add("member_id_number", fmt.Sprintf("invalid RSA ID '%s' (failed local Luhn check)", idNumber))
			}
		}

		schemeCat := strings.TrimSpace(m.SchemeCategory)
		var cat models.SchemeCategory
		var catOK bool
		if schemeCat == "" {
			add("scheme_category", "scheme_category is required")
		} else {
			cat, catOK = categoryByName[schemeCat]
			if !catOK {
				add("scheme_category",
					fmt.Sprintf("scheme_category '%s' is not configured for this scheme's in-force quote", schemeCat))
			}
		}

		if occ := strings.TrimSpace(m.Occupation); occ != "" {
			if _, ok := occupationLookup[strings.ToLower(occ)]; !ok {
				add("occupation", fmt.Sprintf("occupation '%s' does not match any occupation class", occ))
			}
		} else {
			add("occupation", "occupation is required")
		}

		if catOK && !quote.UseGlobalSalaryMultiple {
			if cat.GlaBenefit && m.Benefits.GlaMultiple <= 0 {
				add("benefits.gla_multiple", "gla salary multiple is required when Use Global Salary Multiple is off")
			}
			if cat.SglaBenefit && m.Benefits.SglaMultiple <= 0 {
				add("benefits.sgla_multiple", "sgla salary multiple is required when Use Global Salary Multiple is off")
			}
			if cat.PtdBenefit && m.Benefits.PtdMultiple <= 0 {
				add("benefits.ptd_multiple", "ptd salary multiple is required when Use Global Salary Multiple is off")
			}
			if cat.CiBenefit && m.Benefits.CiMultiple <= 0 {
				add("benefits.ci_multiple", "ci salary multiple is required when Use Global Salary Multiple is off")
			}
			if cat.TtdBenefit && m.Benefits.TtdMultiple <= 0 {
				add("benefits.ttd_multiple", "ttd income replacement is required when Use Global Salary Multiple is off")
			}
			if cat.PhiBenefit && m.Benefits.PhiMultiple <= 0 {
				add("benefits.phi_multiple", "phi income replacement is required when Use Global Salary Multiple is off")
			}
		}

		if idNumber != "" {
			if _, dup := existingIDSet[idNumber]; dup {
				if !skipDuplicates {
					add("member_id_number", fmt.Sprintf("member with ID Number '%s' already exists for this scheme", idNumber))
				}
			}
			if prev, dup := seenInBatch[idNumber]; dup {
				add("member_id_number", fmt.Sprintf("duplicate member_id_number within the same file (also on row %d)", prev))
			} else {
				seenInBatch[idNumber] = row
			}
		}
	}

	return results
}

// isRSAIDType returns true when the member's id_type indicates an RSA ID number.
// Accepted forms mirror the existing single-row enrollment service.
func isRSAIDType(idType string) bool {
	upper := strings.ToUpper(strings.TrimSpace(idType))
	upper = strings.ReplaceAll(upper, " ", "_")
	return upper == "RSA_ID" || upper == "ID" || upper == "RSA_ISD"
}

// RunExternalRSAIDCheckOnBatch calls the external CheckID bulk API against
// every RSA ID in the batch and merges the results into each member's
// ValidationStatus and the batch's ValidationReport JSON.
func RunExternalRSAIDCheckOnBatch(batchID int, user models.AppUser) (models.BulkEnrollmentBatch, []RowValidationResult, error) {
	var batch models.BulkEnrollmentBatch
	if err := DB.First(&batch, batchID).Error; err != nil {
		return batch, nil, fmt.Errorf("batch not found: %w", err)
	}
	if batch.Status != models.BulkEnrollmentBatchPendingApproval {
		return batch, nil, fmt.Errorf("batch is not pending approval")
	}

	var rows []models.GPricingMemberDataInForce
	if err := DB.Where("batch_id = ? AND status = ?", batchID, models.BulkEnrollmentMemberDraftStatus).
		Find(&rows).Error; err != nil {
		return batch, nil, fmt.Errorf("failed to load batch members: %w", err)
	}

	var ids []string
	rowsByID := make(map[string][]int, len(rows))
	for _, r := range rows {
		if isRSAIDType(r.MemberIdType) && strings.TrimSpace(r.MemberIdNumber) != "" {
			ids = append(ids, r.MemberIdNumber)
			rowsByID[r.MemberIdNumber] = append(rowsByID[r.MemberIdNumber], r.RowIndex)
		}
	}

	var existing []RowValidationResult
	if batch.ValidationReport != "" {
		_ = json.Unmarshal([]byte(batch.ValidationReport), &existing)
	}

	if len(ids) == 0 {
		now := time.Now()
		batch.ExternalIDCheckRun = true
		batch.ExternalIDCheckAt = &now
		if err := DB.Save(&batch).Error; err != nil {
			return batch, existing, err
		}
		return batch, existing, nil
	}

	apiResults, err := utils.ValidateRSAIDsBulk(ids)
	if err != nil {
		return batch, existing, fmt.Errorf("ID validation service error: %w", err)
	}

	// Drop any prior external-check findings; preserve everything else so
	// re-running the check doesn't duplicate rows in the report.
	filtered := existing[:0]
	for _, r := range existing {
		if !strings.Contains(r.Message, "external CheckID") {
			filtered = append(filtered, r)
		}
	}
	updatedResults := append([]RowValidationResult{}, filtered...)

	now := time.Now()
	if err := DB.Transaction(func(tx *gorm.DB) error {
		for id, valid := range apiResults {
			rowNums := rowsByID[id]
			for _, rn := range rowNums {
				if !valid {
					updatedResults = append(updatedResults, RowValidationResult{
						Row:            rn,
						Field:          "member_id_number",
						Severity:       "blocking",
						Message:        fmt.Sprintf("RSA ID '%s' failed external CheckID validation", id),
						MemberIdNumber: id,
					})
					if err := tx.Model(&models.GPricingMemberDataInForce{}).
						Where("batch_id = ? AND row_index = ?", batchID, rn).
						Update("validation_status", models.BulkEnrollmentRowBlockingError).Error; err != nil {
						return err
					}
				}
			}
		}

		reportJSON, _ := json.Marshal(updatedResults)
		blocking := 0
		bySeverity := make(map[int]string)
		for _, r := range updatedResults {
			if r.Severity == "blocking" && bySeverity[r.Row] != "blocking" {
				blocking++
				bySeverity[r.Row] = "blocking"
			}
		}

		batch.ValidationReport = string(reportJSON)
		batch.BlockingCount = blocking
		batch.ValidCount = batch.MemberCount - blocking - batch.SoftErrorCount
		batch.ExternalIDCheckRun = true
		batch.ExternalIDCheckAt = &now

		if err := tx.Save(&batch).Error; err != nil {
			return err
		}

		return writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "bulk_enrollment_batches",
			EntityID:  strconv.Itoa(batchID),
			Action:    "EXTERNAL_ID_CHECK",
			ChangedBy: user.UserName,
		}, struct{}{}, batch)
	}); err != nil {
		return batch, updatedResults, err
	}

	return batch, updatedResults, nil
}

// ApproveBulkEnrollmentBatch flips all draft rows in the batch to status='active'
// and transitions the batch to 'approved'. Blocked unless BlockingCount == 0.
func ApproveBulkEnrollmentBatch(batchID int, user models.AppUser) (models.BulkEnrollmentBatch, error) {
	var batch models.BulkEnrollmentBatch
	if err := DB.First(&batch, batchID).Error; err != nil {
		return batch, fmt.Errorf("batch not found: %w", err)
	}
	if batch.Status != models.BulkEnrollmentBatchPendingApproval {
		return batch, fmt.Errorf("batch is not pending approval (current status: %s)", batch.Status)
	}
	if batch.BlockingCount > 0 {
		return batch, fmt.Errorf("batch has %d blocking errors and cannot be approved", batch.BlockingCount)
	}

	before := batch
	now := time.Now()

	if err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.GPricingMemberDataInForce{}).
			Where("batch_id = ? AND status = ?", batchID, models.BulkEnrollmentMemberDraftStatus).
			Updates(map[string]interface{}{
				"status":            "Active",
				"validation_status": models.BulkEnrollmentRowValid,
			}).Error; err != nil {
			return err
		}

		batch.Status = models.BulkEnrollmentBatchApproved
		batch.ApprovedBy = user.UserName
		batch.ApprovedAt = &now
		if err := tx.Save(&batch).Error; err != nil {
			return err
		}

		return writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "bulk_enrollment_batches",
			EntityID:  strconv.Itoa(batchID),
			Action:    "APPROVE",
			ChangedBy: user.UserName,
		}, before, batch)
	}); err != nil {
		return batch, err
	}

	broadcastCacheInvalidation("member_enrollment")
	return batch, nil
}

// RejectBulkEnrollmentBatch hard-deletes the draft member rows for the batch
// and marks the batch as rejected with the provided reason on the audit trail.
func RejectBulkEnrollmentBatch(batchID int, reason string, user models.AppUser) (models.BulkEnrollmentBatch, error) {
	var batch models.BulkEnrollmentBatch
	if err := DB.First(&batch, batchID).Error; err != nil {
		return batch, fmt.Errorf("batch not found: %w", err)
	}
	if batch.Status != models.BulkEnrollmentBatchPendingApproval {
		return batch, fmt.Errorf("batch is not pending approval (current status: %s)", batch.Status)
	}

	before := batch
	now := time.Now()

	if err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("batch_id = ? AND status = ?", batchID, models.BulkEnrollmentMemberDraftStatus).
			Delete(&models.GPricingMemberDataInForce{}).Error; err != nil {
			return err
		}

		batch.Status = models.BulkEnrollmentBatchRejected
		batch.RejectedBy = user.UserName
		batch.RejectedAt = &now
		batch.RejectedReason = strings.TrimSpace(reason)
		if err := tx.Save(&batch).Error; err != nil {
			return err
		}

		return writeAudit(tx, AuditContext{
			Area:      "group-pricing",
			Entity:    "bulk_enrollment_batches",
			EntityID:  strconv.Itoa(batchID),
			Action:    "REJECT",
			ChangedBy: user.UserName,
		}, before, batch)
	}); err != nil {
		return batch, err
	}

	return batch, nil
}

// ListBulkEnrollmentBatches returns batches for a scheme, optionally filtered
// by status. Ordered with the newest upload first so the approvals tab can
// surface fresh work at the top.
func ListBulkEnrollmentBatches(schemeID int, status string) ([]models.BulkEnrollmentBatch, error) {
	q := DB.Where("scheme_id = ?", schemeID).Order("uploaded_at DESC")
	if strings.TrimSpace(status) != "" {
		q = q.Where("status = ?", status)
	}
	var batches []models.BulkEnrollmentBatch
	if err := q.Find(&batches).Error; err != nil {
		return nil, err
	}
	return batches, nil
}

// GetBulkEnrollmentBatch returns a single batch together with all member rows
// linked to it. The frontend renders these in the batch review screen with
// per-row validation chips driven by ValidationStatus + ValidationReport.
func GetBulkEnrollmentBatch(batchID int) (models.BulkEnrollmentBatch, []models.GPricingMemberDataInForce, []RowValidationResult, error) {
	var batch models.BulkEnrollmentBatch
	if err := DB.First(&batch, batchID).Error; err != nil {
		return batch, nil, nil, fmt.Errorf("batch not found: %w", err)
	}

	var members []models.GPricingMemberDataInForce
	if err := DB.Where("batch_id = ?", batchID).Order("row_index ASC").Find(&members).Error; err != nil {
		return batch, nil, nil, err
	}

	var report []RowValidationResult
	if batch.ValidationReport != "" {
		_ = json.Unmarshal([]byte(batch.ValidationReport), &report)
	}

	return batch, members, report, nil
}

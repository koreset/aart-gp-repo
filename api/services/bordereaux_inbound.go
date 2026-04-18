package services

import (
	appLog "api/log"
	"api/models"
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// logNotifyFailure reports a background notification failure so silent drops are
// visible in server logs. Called from fire-and-forget goroutines that wrap the
// bordereaux inbound notification hooks.
func logNotifyFailure(event string, submissionID int, err error) {
	if err == nil {
		return
	}
	appLog.WithFields(map[string]interface{}{
		"event":         event,
		"submission_id": submissionID,
		"error":         err.Error(),
	}).Error("submission notification delivery failed")
}

// CreateEmployerSubmission creates a new inbound employer submission record with status=pending_receipt.
func CreateEmployerSubmission(req models.CreateSubmissionRequest, user models.AppUser) (models.EmployerSubmission, error) {
	var scheme models.GroupScheme
	schemeName := ""
	if err := DB.First(&scheme, req.SchemeID).Error; err == nil {
		schemeName = scheme.Name
	}

	sub := models.EmployerSubmission{
		SchemeID:           req.SchemeID,
		SchemeName:         schemeName,
		Month:              req.Month,
		Year:               req.Year,
		DueDate:            req.DueDate,
		SubmittedBy:        req.SubmittedBy,
		Notes:              req.Notes,
		Status:             "pending_receipt",
		IsRetro:            req.IsRetro,
		RetroEffectiveDate: req.RetroEffectiveDate,
		CreatedBy:          user.UserName,
	}
	if err := DB.Create(&sub).Error; err != nil {
		return sub, fmt.Errorf("failed to create submission: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "employer_submissions",
		EntityID:  strconv.Itoa(sub.ID),
		Action:    "CREATE",
		ChangedBy: user.UserName,
	}, struct{}{}, sub)
	return sub, nil
}

// UploadEmployerSubmission saves the uploaded file, parses records, validates them,
// bulk-inserts EmployerSubmissionRecord rows, and advances the submission to "received".
func UploadEmployerSubmission(id int, file multipart.File, header *multipart.FileHeader, user models.AppUser) (models.EmployerSubmission, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, id).Error; err != nil {
		return sub, fmt.Errorf("submission not found: %w", err)
	}

	uploadDir := "data/bordereaux/inbound"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return sub, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// UUID prefix avoids collisions under concurrent uploads (clock resolution
	// isn't reliably nanosecond-precise on virtualised hosts) and removes a
	// predictable-filename enumeration vector.
	safeFilename := fmt.Sprintf("%s_%s", newFileToken(), header.Filename)
	filePath := filepath.Join(uploadDir, safeFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return sub, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	fileSize, err := io.Copy(dst, file)
	if err != nil {
		return sub, fmt.Errorf("failed to save file: %w", err)
	}

	records, err := parseSubmissionFile(filePath, header.Filename)
	if err != nil {
		return sub, fmt.Errorf("failed to parse submission file: %w", err)
	}

	validCount := 0
	invalidCount := 0
	for i := range records {
		records[i].SubmissionID = sub.ID
		if records[i].ValidationStatus == "valid" {
			validCount++
		} else {
			invalidCount++
		}
	}

	// Delete any prior records for this submission before re-inserting
	DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ?", sub.ID).Delete(&models.EmployerSubmissionRecord{})

	if len(records) > 0 {
		if err := DB.CreateInBatches(records, 200).Error; err != nil {
			return sub, fmt.Errorf("failed to save submission records: %w", err)
		}
	}

	before := sub
	now := time.Now()
	sub.FileName = header.Filename
	sub.FilePath = filePath
	sub.FileSize = fileSize
	sub.RecordCount = len(records)
	sub.ValidCount = validCount
	sub.InvalidCount = invalidCount
	sub.ReceivedDate = &now
	sub.Status = "received"

	if err := DB.Save(&sub).Error; err != nil {
		return sub, fmt.Errorf("failed to update submission: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "employer_submissions",
		EntityID:  strconv.Itoa(sub.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, sub)

	// Auto-link matching deadline (non-fatal if none configured)
	LinkDeadlineToSubmission(sub.SchemeID, sub.Month, sub.Year, sub.ID)

	return sub, nil
}

// parseSubmissionFile reads a CSV or Excel file and returns a slice of EmployerSubmissionRecord.
// Column matching is flexible (handles common naming variations).
func parseSubmissionFile(filePath, fileName string) ([]models.EmployerSubmissionRecord, error) {
	lower := strings.ToLower(fileName)
	if strings.HasSuffix(lower, ".xlsx") {
		return parseExcelSubmission(filePath)
	}
	if strings.HasSuffix(lower, ".csv") {
		return parseCSVSubmission(filePath)
	}
	return nil, fmt.Errorf("unsupported file format: only .csv and .xlsx are accepted")
}

func detectDelimiter(f *os.File) rune {
	delimiters := []rune{',', ';', '|', '\t'}
	maxCols := 0
	bestDelim := ','

	// Read first line
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		f.Seek(0, 0)
		return ','
	}
	line := scanner.Text()

	for _, d := range delimiters {
		r := csv.NewReader(strings.NewReader(line))
		r.Comma = d
		if rows, err := r.ReadAll(); err == nil && len(rows) > 0 {
			if len(rows[0]) > maxCols {
				maxCols = len(rows[0])
				bestDelim = d
			}
		}
	}

	f.Seek(0, 0)
	return bestDelim
}

func parseCSVSubmission(path string) ([]models.EmployerSubmissionRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	delim := detectDelimiter(f)
	reader := csv.NewReader(f)
	reader.Comma = delim
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}
	return rowsToSubmissionRecords(rows[0], rows[1:])
}

func parseExcelSubmission(path string) ([]models.EmployerSubmissionRecord, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}
	return rowsToSubmissionRecords(rows[0], rows[1:])
}

// rowsToSubmissionRecords maps raw header+data rows to EmployerSubmissionRecord structs.
func rowsToSubmissionRecords(headers []string, dataRows [][]string) ([]models.EmployerSubmissionRecord, error) {
	// Normalise header names
	norm := make([]string, len(headers))
	for i, h := range headers {
		norm[i] = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(h, " ", "_")))
	}

	colIdx := func(candidates ...string) int {
		for _, c := range candidates {
			for i, h := range norm {
				if h == c {
					return i
				}
			}
		}
		return -1
	}

	idxName := colIdx("member_name", "name", "full_name", "fullname")
	idxEmpNo := colIdx("employee_no", "employee_number", "emp_no", "emp_number", "staff_no")
	idxID := colIdx("id_number", "id_no", "sa_id", "id", "member_id_number")
	idxDOB := colIdx("dob", "date_of_birth", "birth_date")
	idxGender := colIdx("gender", "sex")
	idxSalary := colIdx("salary", "annual_salary", "gross_salary")
	idxBenefit := colIdx("benefit", "benefit_code", "benefit_type")
	idxPremium := colIdx("premium", "premium_amount", "monthly_premium")
	idxEntry := colIdx("entry_date", "commencement_date", "start_date")
	idxExit := colIdx("exit_date", "termination_date", "end_date")
	idxIDType := colIdx("id_type", "member_id_type")

	getCellStr := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}
	getCellFloat := func(row []string, idx int) float64 {
		s := getCellStr(row, idx)
		if s == "" {
			return 0
		}
		v, _ := strconv.ParseFloat(strings.ReplaceAll(s, ",", ""), 64)
		return v
	}

	var records []models.EmployerSubmissionRecord
	for rowNum, row := range dataRows {
		if len(row) == 0 {
			continue
		}

		// Serialize raw row for audit
		rawMap := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(row) {
				rawMap[h] = row[i]
			}
		}
		rawJSON, _ := json.Marshal(rawMap)

		idNum := getCellStr(row, idxID)
		idType := getCellStr(row, idxIDType)

		// Classify ID once; auto-fill type if not supplied in the CSV column.
		idValid := true
		idInvalidReason := ""
		if idNum != "" {
			detectedType, valid, reason := ClassifyMemberID(idNum)
			if idType == "" {
				idType = detectedType
			}
			idValid = valid
			idInvalidReason = reason
		}

		rec := models.EmployerSubmissionRecord{
			RowNumber:        rowNum + 2, // 1-based data row (header is row 1)
			MemberName:       getCellStr(row, idxName),
			EmployeeNumber:   getCellStr(row, idxEmpNo),
			IDNumber:         idNum,
			IDType:           idType,
			DateOfBirth:      getCellStr(row, idxDOB),
			Gender:           getCellStr(row, idxGender),
			Salary:           getCellFloat(row, idxSalary),
			BenefitCode:      getCellStr(row, idxBenefit),
			PremiumAmount:    getCellFloat(row, idxPremium),
			EntryDate:        getCellStr(row, idxEntry),
			ExitDate:         getCellStr(row, idxExit),
			ValidationStatus: "valid",
			RawData:          string(rawJSON),
		}

		// Validation
		if rec.MemberName == "" {
			rec.ValidationStatus = "missing_data"
			rec.ExclusionReason = "member name is required"
		} else if !idValid {
			rec.ValidationStatus = "id_invalid"
			rec.ExclusionReason = idInvalidReason
		}
		if rec.ValidationStatus == "valid" && rec.PremiumAmount < 0 {
			rec.ValidationStatus = "amount_invalid"
			rec.ExclusionReason = "premium amount cannot be negative"
		}

		records = append(records, rec)
	}
	return records, nil
}

// ReviewEmployerSubmission advances a submission from "received" to "under_review".
func ReviewEmployerSubmission(id int, notes string, user models.AppUser) (models.EmployerSubmission, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, id).Error; err != nil {
		return sub, fmt.Errorf("submission not found: %w", err)
	}
	if err := ValidateEmployerSubmissionTransition(sub.Status, StatusSubmissionUnderReview); err != nil {
		return sub, err
	}
	before := sub

	now := time.Now()
	sub.Status = StatusSubmissionUnderReview
	sub.ReviewedBy = user.UserName
	sub.ReviewedAt = &now
	if notes != "" {
		sub.Notes = notes
	}
	if err := DB.Save(&sub).Error; err != nil {
		return sub, fmt.Errorf("failed to update submission: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "employer_submissions",
		EntityID:  strconv.Itoa(sub.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, sub)
	go func() {
		logNotifyFailure("submission_reviewed", sub.ID, NotifySubmissionReviewed(sub, user))
	}()
	return sub, nil
}

// RaiseSubmissionQuery raises a query on an "under_review" submission.
func RaiseSubmissionQuery(id int, queryNotes string, user models.AppUser) (models.EmployerSubmission, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, id).Error; err != nil {
		return sub, fmt.Errorf("submission not found: %w", err)
	}
	if err := ValidateEmployerSubmissionTransition(sub.Status, StatusSubmissionQueriesRaised); err != nil {
		return sub, err
	}
	before := sub

	sub.Status = StatusSubmissionQueriesRaised
	sub.QueryNotes = queryNotes
	if err := DB.Save(&sub).Error; err != nil {
		return sub, fmt.Errorf("failed to update submission: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "employer_submissions",
		EntityID:  strconv.Itoa(sub.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, sub)
	go func() {
		logNotifyFailure("submission_query_raised", sub.ID, NotifySubmissionQueryRaised(sub, user))
	}()
	return sub, nil
}

// AcceptEmployerSubmission accepts a submission that is "under_review" or "queries_raised".
func AcceptEmployerSubmission(id int, notes string, user models.AppUser) (models.EmployerSubmission, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, id).Error; err != nil {
		return sub, fmt.Errorf("submission not found: %w", err)
	}
	if err := ValidateEmployerSubmissionTransition(sub.Status, StatusSubmissionAccepted); err != nil {
		return sub, err
	}
	before := sub

	now := time.Now()
	sub.Status = StatusSubmissionAccepted
	sub.AcceptedBy = user.UserName
	sub.AcceptedAt = &now
	if notes != "" {
		if sub.Notes != "" {
			sub.Notes = sub.Notes + "\n" + notes
		} else {
			sub.Notes = notes
		}
	}
	if err := DB.Save(&sub).Error; err != nil {
		return sub, fmt.Errorf("failed to update submission: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "employer_submissions",
		EntityID:  strconv.Itoa(sub.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, sub)

	// Auto-snapshot the register diff at the moment of acceptance so historical records are preserved
	// even after subsequent member register updates (e.g. February sync overwriting January state).
	if liveResult, diffErr := computeRegisterDiffLive(sub.ID); diffErr == nil {
		_ = saveRegisterDiffSnapshot(liveResult, sub.ID, user)
	}

	go func() {
		logNotifyFailure("submission_accepted", sub.ID, NotifySubmissionAccepted(sub, user))
	}()
	return sub, nil
}

// RejectEmployerSubmission rejects a submission that is received, under_review, or queries_raised.
func RejectEmployerSubmission(id int, reason string, user models.AppUser) (models.EmployerSubmission, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, id).Error; err != nil {
		return sub, fmt.Errorf("submission not found: %w", err)
	}
	if err := ValidateEmployerSubmissionTransition(sub.Status, StatusSubmissionRejected); err != nil {
		return sub, err
	}
	before := sub

	now := time.Now()
	sub.Status = StatusSubmissionRejected
	sub.RejectedBy = user.UserName
	sub.RejectedAt = &now
	sub.RejectionReason = reason
	if err := DB.Save(&sub).Error; err != nil {
		return sub, fmt.Errorf("failed to update submission: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "employer_submissions",
		EntityID:  strconv.Itoa(sub.ID),
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, sub)
	go func() {
		logNotifyFailure("submission_rejected", sub.ID, NotifySubmissionRejected(sub, user, reason))
	}()
	return sub, nil
}

// GetEmployerSubmissions returns submissions filtered by optional scheme, month, year, status.
func GetEmployerSubmissions(schemeID, month, year int, status string) ([]models.EmployerSubmission, error) {
	var subs []models.EmployerSubmission
	q := DB.Model(&models.EmployerSubmission{})
	if schemeID > 0 {
		q = q.Where("scheme_id = ?", schemeID)
	}
	if month > 0 {
		q = q.Where("month = ?", month)
	}
	if year > 0 {
		q = q.Where("year = ?", year)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at DESC").Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

// GetEmployerSubmission fetches a single submission with its records.
// For submissions with fewer than 1 000 records all rows are returned.
// For larger submissions only invalid rows are returned (so the UI can
// always display validation failures regardless of submission size).
func GetEmployerSubmission(id int) (models.EmployerSubmission, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, id).Error; err != nil {
		return sub, fmt.Errorf("submission not found: %w", err)
	}
	// Always initialise to a non-nil slice so the JSON field is never omitted.
	sub.Records = []models.EmployerSubmissionRecord{}
	if sub.RecordCount > 0 {
		var q *gorm.DB
		if sub.RecordCount < 1000 {
			// Small submission — load everything
			q = DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ?", id)
		} else {
			// Large submission — load only invalid rows so failures are always visible
			q = DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ? AND validation_status != 'valid'", id)
		}
		if err := q.Find(&sub.Records).Error; err != nil {
			fmt.Printf("GetEmployerSubmission: failed to load records for submission %d: %v\n", id, err)
		}
	}
	return sub, nil
}

// GetSubmissionRecords returns all records for a submission ordered by row number.
func GetSubmissionRecords(submissionID int) ([]models.EmployerSubmissionRecord, error) {
	var records []models.EmployerSubmissionRecord
	if err := DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ?", submissionID).Order("row_number").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// GenerateScheduleFromSubmission generates a premium schedule from an accepted submission,
// then cross-links both records.
func GenerateScheduleFromSubmission(submissionID int, user models.AppUser) (models.PremiumSchedule, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, submissionID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("submission not found: %w", err)
	}
	if sub.Status != "accepted" {
		return models.PremiumSchedule{}, fmt.Errorf("submission must be 'accepted' before generating a schedule (current: %s)", sub.Status)
	}

	// If already linked, return the existing schedule
	if sub.LinkedPremiumScheduleID != nil {
		var existing models.PremiumSchedule
		if err := DB.First(&existing, *sub.LinkedPremiumScheduleID).Error; err == nil {
			return existing, nil
		}
	}

	// Dispatch to retro path when flagged
	if sub.IsRetro {
		return generateRetroSupplementarySchedule(sub, user)
	}

	schedule, err := GenerateMonthlySchedule(sub.SchemeID, sub.Month, sub.Year, user)
	if err != nil {
		return schedule, fmt.Errorf("schedule generation failed: %w", err)
	}

	// Cross-link
	schedule.LinkedSubmissionID = &submissionID
	if err := DB.Save(&schedule).Error; err != nil {
		return schedule, fmt.Errorf("failed to link schedule to submission: %w", err)
	}

	sub.LinkedPremiumScheduleID = &schedule.ID
	if err := DB.Save(&sub).Error; err != nil {
		return schedule, fmt.Errorf("failed to link submission to schedule: %w", err)
	}

	return schedule, nil
}

// ComputeSubmissionDelta compares an accepted submission's records against the prior month's
// accepted submission for the same scheme, classifying each member row as:
// new | amendment | ceased | continuing
func ComputeSubmissionDelta(submissionID int, user models.AppUser) (models.DeltaSummary, error) {
	var sub models.EmployerSubmission
	if err := DB.First(&sub, submissionID).Error; err != nil {
		return models.DeltaSummary{}, fmt.Errorf("submission not found: %w", err)
	}
	if sub.Status != "accepted" {
		return models.DeltaSummary{}, fmt.Errorf("submission must be 'accepted' to compute delta (current: %s)", sub.Status)
	}

	// Load current submission records
	var currentRecords []models.EmployerSubmissionRecord
	if err := DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ?", submissionID).Find(&currentRecords).Error; err != nil {
		return models.DeltaSummary{}, fmt.Errorf("failed to load current records: %w", err)
	}

	// Find prior accepted submission (same scheme, prior month)
	priorMonth := sub.Month - 1
	priorYear := sub.Year
	if priorMonth == 0 {
		priorMonth = 12
		priorYear--
	}

	var priorSub models.EmployerSubmission
	priorFound := DB.Model(&models.EmployerSubmission{}).Where("scheme_id = ? AND month = ? AND year = ? AND status = ?",
		sub.SchemeID, priorMonth, priorYear, "accepted").
		Order("id DESC").First(&priorSub).Error == nil

	priorMap := make(map[string]models.EmployerSubmissionRecord)
	if priorFound {
		var priorRecords []models.EmployerSubmissionRecord
		DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ?", priorSub.ID).Find(&priorRecords)
		for _, r := range priorRecords {
			key := memberKey(r)
			priorMap[key] = r
		}
	}

	// Build current key set
	currentMap := make(map[string]models.EmployerSubmissionRecord)
	for _, r := range currentRecords {
		key := memberKey(r)
		currentMap[key] = r
	}

	summary := models.DeltaSummary{
		SubmissionID: submissionID,
	}
	if priorFound {
		summary.PriorSubmissionID = priorSub.ID
	}

	var deltaRecords []models.SubmissionDeltaRecord

	today := time.Now()

	// Classify current records (new | amendment | continuing | ceased)
	for key, cur := range currentMap {
		dr := models.SubmissionDeltaRecord{
			SubmissionID:      submissionID,
			PriorSubmissionID: summary.PriorSubmissionID,
			MemberKey:         key,
			MemberName:        cur.MemberName,
			RowNumber:         cur.RowNumber,
		}
		// Exit records (explicit exit_date or zero premium) count as ceased
		if isExitRecord(cur, today) {
			dr.ChangeType = "ceased"
			summary.Ceased++
			deltaRecords = append(deltaRecords, dr)
			continue
		}
		if prior, exists := priorMap[key]; exists {
			changed := computeChangedFields(prior, cur)
			if len(changed) == 0 {
				dr.ChangeType = "continuing"
				summary.Continuing++
			} else {
				changedJSON, _ := json.Marshal(changed)
				dr.ChangeType = "amendment"
				dr.ChangedFields = string(changedJSON)
				summary.Amendment++
			}
		} else {
			dr.ChangeType = "new"
			summary.New++
		}
		deltaRecords = append(deltaRecords, dr)
	}

	// Classify ceased (in prior but not in current)
	for key, prior := range priorMap {
		if _, exists := currentMap[key]; !exists {
			deltaRecords = append(deltaRecords, models.SubmissionDeltaRecord{
				SubmissionID:      submissionID,
				PriorSubmissionID: summary.PriorSubmissionID,
				MemberKey:         key,
				MemberName:        prior.MemberName,
				ChangeType:        "ceased",
				RowNumber:         prior.RowNumber,
			})
			summary.Ceased++
		}
	}

	summary.Total = summary.New + summary.Amendment + summary.Ceased + summary.Continuing

	// Atomically replace delta records
	DB.Where("submission_id = ?", submissionID).Delete(&models.SubmissionDeltaRecord{})
	if len(deltaRecords) > 0 {
		if err := DB.CreateInBatches(deltaRecords, 200).Error; err != nil {
			return summary, fmt.Errorf("failed to save delta records: %w", err)
		}
	}

	return summary, nil
}

// memberKey returns a stable identifier for a submission record, preferring employee_number.
func memberKey(r models.EmployerSubmissionRecord) string {
	if r.EmployeeNumber != "" {
		return "emp:" + r.EmployeeNumber
	}
	if r.IDNumber != "" {
		return "id:" + r.IDNumber
	}
	return "name:" + strings.ToLower(strings.TrimSpace(r.MemberName))
}

// computeChangedFields returns a map of field names to [old, new] pairs where values differ.
func computeChangedFields(prior, cur models.EmployerSubmissionRecord) map[string][2]string {
	changed := make(map[string][2]string)
	if fmt.Sprintf("%.2f", prior.Salary) != fmt.Sprintf("%.2f", cur.Salary) {
		changed["salary"] = [2]string{fmt.Sprintf("%.2f", prior.Salary), fmt.Sprintf("%.2f", cur.Salary)}
	}
	if fmt.Sprintf("%.2f", prior.PremiumAmount) != fmt.Sprintf("%.2f", cur.PremiumAmount) {
		changed["premium_amount"] = [2]string{fmt.Sprintf("%.2f", prior.PremiumAmount), fmt.Sprintf("%.2f", cur.PremiumAmount)}
	}
	if prior.BenefitCode != cur.BenefitCode {
		changed["benefit_code"] = [2]string{prior.BenefitCode, cur.BenefitCode}
	}
	if prior.DateOfBirth != cur.DateOfBirth {
		changed["date_of_birth"] = [2]string{prior.DateOfBirth, cur.DateOfBirth}
	}
	if prior.EntryDate != cur.EntryDate {
		changed["entry_date"] = [2]string{prior.EntryDate, cur.EntryDate}
	}
	if prior.ExitDate != cur.ExitDate {
		changed["exit_date"] = [2]string{prior.ExitDate, cur.ExitDate}
	}
	return changed
}

// GetSubmissionDeltaRecords returns all delta records for a submission.
func GetSubmissionDeltaRecords(submissionID int) ([]models.SubmissionDeltaRecord, error) {
	var records []models.SubmissionDeltaRecord
	if err := DB.Where("submission_id = ?", submissionID).
		Order("change_type, member_name").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// SyncSubmissionToMemberRegister applies accepted submission delta records to the live member
// register (GPricingMemberDataInForce). New members are enrolled, ceased members are
// deactivated, and amended members have their changed fields updated.
// The submission must be in "accepted" status and have a computed delta.
func SyncSubmissionToMemberRegister(submissionID int, user models.AppUser) (models.MemberSyncResult, error) {
	var result models.MemberSyncResult

	var sub models.EmployerSubmission
	if err := DB.First(&sub, submissionID).Error; err != nil {
		return result, fmt.Errorf("submission not found: %w", err)
	}
	if sub.Status != "accepted" {
		return result, fmt.Errorf("submission must be 'accepted' to sync members (current: %s)", sub.Status)
	}

	// Load delta records — sync requires delta to have been computed
	var deltaRecords []models.SubmissionDeltaRecord
	if err := DB.Where("submission_id = ?", submissionID).Find(&deltaRecords).Error; err != nil {
		return result, fmt.Errorf("failed to load delta records: %w", err)
	}
	if len(deltaRecords) == 0 {
		return result, fmt.Errorf("no delta records found — run 'Compute Delta' before syncing")
	}

	// Build a lookup map from member_key → submission record (for new/amendment cases)
	var submissionRecords []models.EmployerSubmissionRecord
	DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ? AND validation_status = 'valid'", submissionID).Find(&submissionRecords)
	recByKey := make(map[string]models.EmployerSubmissionRecord, len(submissionRecords))
	for _, r := range submissionRecords {
		recByKey[memberKey(r)] = r
	}

	for _, dr := range deltaRecords {
		switch dr.ChangeType {
		case "new":
			rec, ok := recByKey[dr.MemberKey]
			if !ok {
				result.Skipped++
				result.Errors = append(result.Errors, fmt.Sprintf("new member '%s': submission record not found", dr.MemberName))
				continue
			}
			member := submissionRecordToMemberInForce(rec, sub.SchemeID)
			_, err := AddMemberToScheme(member, user)
			if err != nil {
				result.Skipped++
				result.Errors = append(result.Errors, fmt.Sprintf("add '%s': %s", dr.MemberName, err.Error()))
			} else {
				result.Added++
			}

		case "ceased":
			existing, err := findMemberInScheme(sub.SchemeID, dr.MemberKey)
			if err != nil {
				result.Skipped++
				result.Errors = append(result.Errors, fmt.Sprintf("ceased '%s': member not found in register", dr.MemberName))
				continue
			}
			if existing.Status == "Inactive" {
				result.Skipped++
				continue
			}
			if err := DeactivateSchemeMember(strconv.Itoa(sub.SchemeID), existing, user); err != nil {
				result.Skipped++
				result.Errors = append(result.Errors, fmt.Sprintf("deactivate '%s': %s", dr.MemberName, err.Error()))
			} else {
				result.Deactivated++
			}

		case "amendment":
			existing, err := findMemberInScheme(sub.SchemeID, dr.MemberKey)
			if err != nil {
				result.Skipped++
				result.Errors = append(result.Errors, fmt.Sprintf("amendment '%s': member not found in register", dr.MemberName))
				continue
			}
			updated := applyDeltaChanges(existing, dr.ChangedFields)
			if _, err := UpdateMemberInForce(strconv.Itoa(existing.ID), updated, user); err != nil {
				result.Skipped++
				result.Errors = append(result.Errors, fmt.Sprintf("update '%s': %s", dr.MemberName, err.Error()))
			} else {
				result.Updated++
			}

		case "continuing":
			result.Skipped++
		}
	}

	// Persist sync timestamp and summary on the submission
	now := time.Now()
	summaryJSON, _ := json.Marshal(result)
	DB.Model(&sub).Updates(map[string]interface{}{
		"member_synced_at":    now,
		"member_sync_summary": string(summaryJSON),
	})

	return result, nil
}

// submissionRecordToMemberInForce maps an EmployerSubmissionRecord to a GPricingMemberDataInForce
// ready for AddMemberToScheme. Benefits are resolved from the scheme's active quote category
// inside AddMemberToScheme itself — we only populate what the submission provides.
func submissionRecordToMemberInForce(rec models.EmployerSubmissionRecord, schemeID int) models.GPricingMemberDataInForce {
	m := models.GPricingMemberDataInForce{
		SchemeId:       schemeID,
		MemberName:     rec.MemberName,
		MemberIdNumber: rec.IDNumber,
		EmployeeNumber: rec.EmployeeNumber,
		AnnualSalary:   rec.Salary,
	}
	if rec.DateOfBirth != "" {
		if dob, err := parseFlexDate(rec.DateOfBirth); err == nil {
			m.DateOfBirth = dob
		}
	}
	if rec.EntryDate != "" {
		if ed, err := parseFlexDate(rec.EntryDate); err == nil {
			m.EntryDate = ed
		}
	} else {
		m.EntryDate = time.Now()
	}
	if rec.ExitDate != "" {
		if ex, err := parseFlexDate(rec.ExitDate); err == nil {
			m.ExitDate = &ex
		}
	}
	return m
}

// findMemberInScheme looks up a member within a specific scheme using the member key
// (which encodes employee number, ID number, or name).
func findMemberInScheme(schemeID int, key string) (models.GPricingMemberDataInForce, error) {
	var m models.GPricingMemberDataInForce
	q := DB.Where("scheme_id = ?", schemeID)

	if strings.HasPrefix(key, "emp:") {
		empNo := strings.TrimPrefix(key, "emp:")
		if err := q.Where("employee_number = ?", empNo).
			Order("CASE WHEN status = 'Active' THEN 0 ELSE 1 END").
			First(&m).Error; err == nil {
			return m, nil
		}
	}
	if strings.HasPrefix(key, "id:") {
		idNum := strings.TrimPrefix(key, "id:")
		if err := q.Where("member_id_number = ?", idNum).
			Order("CASE WHEN status = 'Active' THEN 0 ELSE 1 END").
			First(&m).Error; err == nil {
			return m, nil
		}
	}
	if strings.HasPrefix(key, "name:") {
		name := strings.TrimPrefix(key, "name:")
		if err := q.Where("LOWER(member_name) = ?", name).
			Order("CASE WHEN status = 'Active' THEN 0 ELSE 1 END").
			First(&m).Error; err == nil {
			return m, nil
		}
	}
	return m, gorm.ErrRecordNotFound
}

// applyDeltaChanges applies only the fields listed in the changedFields JSON to the member,
// returning the updated struct for UpdateMemberInForce.
func applyDeltaChanges(m models.GPricingMemberDataInForce, changedFieldsJSON string) models.GPricingMemberDataInForce {
	if changedFieldsJSON == "" {
		return m
	}
	var fields map[string][2]string
	if err := json.Unmarshal([]byte(changedFieldsJSON), &fields); err != nil {
		return m
	}
	for field, vals := range fields {
		newVal := vals[1]
		switch field {
		case "salary":
			if v, err := strconv.ParseFloat(newVal, 64); err == nil {
				m.AnnualSalary = v
			}
		case "gender":
			m.Gender = newVal
		case "benefit_code":
			// benefit_code changes are informational; actual benefit flags are
			// managed via the scheme category, not the string code
		case "date_of_birth":
			if t, err := parseFlexDate(newVal); err == nil {
				m.DateOfBirth = t
			}
		case "entry_date":
			if t, err := parseFlexDate(newVal); err == nil {
				m.EntryDate = t
			}
		case "exit_date":
			if t, err := parseFlexDate(newVal); err == nil {
				m.ExitDate = &t
			}
		case "premium_amount":
			// premium is computed — not stored directly on the member record
		}
	}
	return m
}

// parseFlexDate attempts to parse a date string in several common formats.
func parseFlexDate(s string) (time.Time, error) {
	formats := []string{"2006-01-02", "02/01/2006", "01/02/2006", "2006/01/02", "02-01-2006"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised date format: %s", s)
}

// generateRetroSupplementarySchedule creates a supplementary PremiumSchedule representing
// the catch-up premium for a retro submission, plus a retro adjustment invoice.
// It is called by GenerateScheduleFromSubmission when sub.IsRetro == true.
func generateRetroSupplementarySchedule(sub models.EmployerSubmission, user models.AppUser) (models.PremiumSchedule, error) {
	if sub.RetroEffectiveDate == "" {
		return models.PremiumSchedule{}, fmt.Errorf("retro_effective_date is required for retro submissions")
	}

	retroDate, err := time.Parse("2006-01-02", sub.RetroEffectiveDate)
	if err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("invalid retro_effective_date format (expected YYYY-MM-DD): %w", err)
	}

	// Count full catch-up months: from retro start up to (but not including) submission month
	retroStart := time.Date(retroDate.Year(), retroDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	submissionStart := time.Date(sub.Year, time.Month(sub.Month), 1, 0, 0, 0, 0, time.UTC)

	catchUpMonths := 0
	for t := retroStart; t.Before(submissionStart); t = t.AddDate(0, 1, 0) {
		catchUpMonths++
	}
	if catchUpMonths <= 0 {
		return models.PremiumSchedule{}, fmt.Errorf("retro_effective_date must be before the submission period (%d/%d)", sub.Month, sub.Year)
	}

	var records []models.EmployerSubmissionRecord
	DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ? AND validation_status = 'valid'", sub.ID).Find(&records)
	monthlyPremium := 0.0
	for _, r := range records {
		monthlyPremium += r.PremiumAmount
	}
	totalCatchUp := monthlyPremium * float64(catchUpMonths)

	now := time.Now()
	note := fmt.Sprintf("Retro adjustment: %d month(s) from %s (submitted %02d/%d)", catchUpMonths, sub.RetroEffectiveDate, sub.Month, sub.Year)

	schedule := models.PremiumSchedule{
		SchemeID:           sub.SchemeID,
		SchemeName:         sub.SchemeName,
		Month:              sub.Month,
		Year:               sub.Year,
		MemberCount:        len(records),
		GrossPremium:       totalCatchUp,
		NetPayable:         totalCatchUp,
		Status:             "finalized",
		GeneratedDate:      now.Format("2006-01-02"),
		GeneratedBy:        user.UserName,
		FinalizedBy:        user.UserName,
		FinalizedAt:        &now,
		IsSupplementary:    true,
		SupplementaryNote:  note,
		LinkedSubmissionID: &sub.ID,
	}

	if err := DB.Create(&schedule).Error; err != nil {
		return schedule, fmt.Errorf("failed to create supplementary schedule: %w", err)
	}

	// Generate a retro adjustment invoice immediately (skip the normal finalize step)
	invoice := models.Invoice{
		InvoiceNumber:     fmt.Sprintf("INV-RETRO-%d-%02d-%d-%d", sub.SchemeID, sub.Month, sub.Year, schedule.ID),
		SchemeID:          sub.SchemeID,
		SchemeName:        sub.SchemeName,
		ScheduleID:        schedule.ID,
		Month:             sub.Month,
		Year:              sub.Year,
		IssueDate:         now.Format("2006-01-02"),
		DueDate:           now.AddDate(0, 0, 30).Format("2006-01-02"),
		GrossAmount:       totalCatchUp,
		NetPayable:        totalCatchUp,
		Balance:           totalCatchUp,
		Status:            "sent",
		IsRetroAdjustment: true,
	}
	DB.Create(&invoice)

	// Cross-link submission → schedule
	sub.LinkedPremiumScheduleID = &schedule.ID
	DB.Save(&sub)

	return schedule, nil
}

// isExitRecord returns true if a submission record represents a member exit.
// A record is an exit when it carries an explicit exit_date on or before today,
// OR when the premium amount is zero (absent/empty in the source CSV).
func isExitRecord(rec models.EmployerSubmissionRecord, today time.Time) bool {
	if rec.PremiumAmount == 0 {
		return true
	}
	if rec.ExitDate != "" {
		if t, err := parseFlexDate(rec.ExitDate); err == nil && !t.After(today) {
			return true
		}
	}
	return false
}

// computeRegisterDiffLive does the live comparison against GPricingMemberDataInForce.
// It does NOT persist anything — use saveRegisterDiffSnapshot to persist the result.
func computeRegisterDiffLive(submissionID int) (models.RegisterDiffResult, error) {
	result := models.RegisterDiffResult{SubmissionID: submissionID}

	var sub models.EmployerSubmission
	if err := DB.First(&sub, submissionID).Error; err != nil {
		return result, fmt.Errorf("submission not found: %w", err)
	}

	var records []models.EmployerSubmissionRecord
	if err := DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ? AND validation_status = 'valid'", submissionID).Find(&records).Error; err != nil {
		return result, fmt.Errorf("failed to load submission records: %w", err)
	}

	today := time.Now()

	for _, rec := range records {
		key := memberKey(rec)
		dm := models.RegisterDiffMember{
			MemberKey:      key,
			MemberName:     rec.MemberName,
			EmployeeNumber: rec.EmployeeNumber,
			IDNumber:       rec.IDNumber,
			RowNumber:      rec.RowNumber,
		}

		// Exit: explicit exit_date on or before today, or zero premium (treated as implicit exit)
		if isExitRecord(rec, today) {
			dm.ExitDate = rec.ExitDate
			result.Exits = append(result.Exits, dm)
			continue
		}

		// Try to find member in the live register
		existing, err := findMemberInScheme(sub.SchemeID, key)
		if err != nil {
			// Not found → new joiner
			result.NewJoiners = append(result.NewJoiners, dm)
			continue
		}

		// Found → check for field changes
		changed := computeChangedFieldsFromRecord(existing, rec)
		if len(changed) > 0 {
			dm.ChangedFields = changed
			result.Amendments = append(result.Amendments, dm)
		} else {
			result.Continuing++
		}
	}

	return result, nil
}

// saveRegisterDiffSnapshot persists (or overwrites) the RegisterDiffSnapshot for a submission.
func saveRegisterDiffSnapshot(result models.RegisterDiffResult, submissionID int, user models.AppUser) error {
	diffJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to serialise diff: %w", err)
	}

	now := time.Now()
	snap := models.RegisterDiffSnapshot{
		SubmissionID:    submissionID,
		SnapshotAt:      now,
		SnapshotBy:      user.UserName,
		NewJoinersCount: len(result.NewJoiners),
		ExitsCount:      len(result.Exits),
		AmendmentsCount: len(result.Amendments),
		ContinuingCount: result.Continuing,
		DiffJSON:        string(diffJSON),
	}

	// Upsert: delete existing snapshot then insert fresh (uniqueIndex on SubmissionID)
	DB.Where("submission_id = ?", submissionID).Delete(&models.RegisterDiffSnapshot{})
	return DB.Create(&snap).Error
}

// ComputeRegisterDiff returns the persisted RegisterDiffSnapshot when one exists (is_snapshot: true),
// ensuring historical accuracy even after subsequent member register updates.
// For submissions that have not yet been snapshotted (e.g. still under review), it computes live.
func ComputeRegisterDiff(submissionID int) (models.RegisterDiffResult, error) {
	// Check for an existing persisted snapshot first.
	var snap models.RegisterDiffSnapshot
	if err := DB.Where("submission_id = ?", submissionID).First(&snap).Error; err == nil {
		// Snapshot exists — deserialise and return it.
		var result models.RegisterDiffResult
		if jsonErr := json.Unmarshal([]byte(snap.DiffJSON), &result); jsonErr != nil {
			return models.RegisterDiffResult{SubmissionID: submissionID},
				fmt.Errorf("failed to deserialise snapshot: %w", jsonErr)
		}
		result.IsSnapshot = true
		result.SnapshotAt = &snap.SnapshotAt
		result.SnapshotBy = snap.SnapshotBy
		return result, nil
	}

	// No snapshot — compute against the live register.
	return computeRegisterDiffLive(submissionID)
}

// RefreshRegisterDiffSnapshot forces a fresh live computation and overwrites any existing snapshot.
// Useful for backfilling historical submissions or correcting a snapshot taken at the wrong moment.
func RefreshRegisterDiffSnapshot(submissionID int, user models.AppUser) (models.RegisterDiffResult, error) {
	result, err := computeRegisterDiffLive(submissionID)
	if err != nil {
		return result, err
	}
	if snapErr := saveRegisterDiffSnapshot(result, submissionID, user); snapErr != nil {
		return result, fmt.Errorf("diff computed but snapshot save failed: %w", snapErr)
	}
	result.IsSnapshot = true
	now := time.Now()
	result.SnapshotAt = &now
	result.SnapshotBy = user.UserName
	return result, nil
}

// computeChangedFieldsFromRecord compares a live register member against a submission record
// and returns a map of field names to [old, new] pairs for fields that differ.
func computeChangedFieldsFromRecord(existing models.GPricingMemberDataInForce, rec models.EmployerSubmissionRecord) map[string][2]string {
	changed := make(map[string][2]string)

	// Salary (only compare when submission provides a non-zero value)
	if rec.Salary > 0 && fmt.Sprintf("%.2f", existing.AnnualSalary) != fmt.Sprintf("%.2f", rec.Salary) {
		changed["salary"] = [2]string{fmt.Sprintf("%.2f", existing.AnnualSalary), fmt.Sprintf("%.2f", rec.Salary)}
	}

	// Gender (only compare when submission provides a value)
	if rec.Gender != "" && !strings.EqualFold(existing.Gender, rec.Gender) {
		changed["gender"] = [2]string{existing.Gender, rec.Gender}
	}

	// Date of birth (only compare when submission and register both have values)
	if rec.DateOfBirth != "" && !existing.DateOfBirth.IsZero() {
		if recDOB, err := parseFlexDate(rec.DateOfBirth); err == nil {
			existingStr := existing.DateOfBirth.Format("2006-01-02")
			recStr := recDOB.Format("2006-01-02")
			if existingStr != recStr {
				changed["date_of_birth"] = [2]string{existingStr, recStr}
			}
		}
	}

	// Benefit code → compared against scheme category
	if rec.BenefitCode != "" && rec.BenefitCode != existing.SchemeCategory {
		changed["benefit_code"] = [2]string{existing.SchemeCategory, rec.BenefitCode}
	}

	return changed
}

// ApplySubmissionExits finds all valid submission records that represent an exit —
// either an explicit exit_date ≤ today or a zero/absent premium amount —
// and deactivates the corresponding members in the live register.
func ApplySubmissionExits(submissionID int, user models.AppUser) (models.MemberSyncResult, error) {
	var result models.MemberSyncResult

	var sub models.EmployerSubmission
	if err := DB.First(&sub, submissionID).Error; err != nil {
		return result, fmt.Errorf("submission not found: %w", err)
	}
	if sub.Status != "accepted" {
		return result, fmt.Errorf("submission must be 'accepted' to apply exits (current: %s)", sub.Status)
	}

	var records []models.EmployerSubmissionRecord
	DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ? AND validation_status = 'valid'", submissionID).Find(&records)

	today := time.Now()
	for _, rec := range records {
		if !isExitRecord(rec, today) {
			continue
		}

		key := memberKey(rec)
		existing, err := findMemberInScheme(sub.SchemeID, key)
		if err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("exit '%s': member not found in register", rec.MemberName))
			continue
		}
		if existing.Status == "Inactive" {
			result.Skipped++
			continue
		}
		if err := DeactivateSchemeMember(strconv.Itoa(sub.SchemeID), existing, user); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("deactivate '%s': %s", rec.MemberName, err.Error()))
		} else {
			result.Deactivated++
		}
	}

	now := time.Now()
	summaryJSON, _ := json.Marshal(result)
	DB.Model(&sub).Updates(map[string]interface{}{
		"exits_synced_at":    now,
		"exits_sync_summary": string(summaryJSON),
	})

	return result, nil
}

// ApplySubmissionAmendments finds valid submission records where the member exists in the
// register and fields have changed, and applies those field updates.
func ApplySubmissionAmendments(submissionID int, user models.AppUser) (models.MemberSyncResult, error) {
	var result models.MemberSyncResult

	var sub models.EmployerSubmission
	if err := DB.First(&sub, submissionID).Error; err != nil {
		return result, fmt.Errorf("submission not found: %w", err)
	}
	if sub.Status != "accepted" {
		return result, fmt.Errorf("submission must be 'accepted' to apply amendments (current: %s)", sub.Status)
	}

	var records []models.EmployerSubmissionRecord
	DB.Model(&models.EmployerSubmissionRecord{}).Where("submission_id = ? AND validation_status = 'valid'", submissionID).Find(&records)

	today := time.Now()
	for _, rec := range records {
		// Skip exit records (handled by ApplySubmissionExits)
		if isExitRecord(rec, today) {
			continue
		}

		key := memberKey(rec)
		existing, err := findMemberInScheme(sub.SchemeID, key)
		if err != nil {
			// Not in register → new joiner, handled separately
			continue
		}

		changed := computeChangedFieldsFromRecord(existing, rec)
		if len(changed) == 0 {
			result.Skipped++
			continue
		}

		changedJSON, _ := json.Marshal(changed)
		updated := applyDeltaChanges(existing, string(changedJSON))
		if _, err := UpdateMemberInForce(strconv.Itoa(existing.ID), updated, user); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("update '%s': %s", rec.MemberName, err.Error()))
		} else {
			result.Updated++
		}
	}

	now := time.Now()
	summaryJSON, _ := json.Marshal(result)
	DB.Model(&sub).Updates(map[string]interface{}{
		"amendments_synced_at":    now,
		"amendments_sync_summary": string(summaryJSON),
	})

	return result, nil
}

// GetNewJoinerDetails returns all staged NewJoinerDetail records for a submission.
func GetNewJoinerDetails(submissionID int) ([]models.NewJoinerDetail, error) {
	var details []models.NewJoinerDetail
	if err := DB.Where("submission_id = ?", submissionID).Order("id").Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

// UploadNewJoinerDetails parses a detailed CSV/XLSX file containing full member data for
// new joiners, validates rows, and stages them in the new_joiner_details table.
func UploadNewJoinerDetails(submissionID int, file multipart.File, header *multipart.FileHeader, user models.AppUser) ([]models.NewJoinerDetail, error) {
	// Write to a temp file so the parser can re-read it
	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("nj_%d_%d%s", submissionID, time.Now().UnixNano(), filepath.Ext(header.Filename)))
	dst, err := os.Create(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		os.Remove(tmpPath)
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	dst.Close()
	defer os.Remove(tmpPath)

	details, err := parseNewJoinerFile(tmpPath, header.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse new joiner file: %w", err)
	}

	for i := range details {
		details[i].SubmissionID = submissionID
	}

	// Replace any prior staged records for this submission
	DB.Model(&models.NewJoinerDetail{}).Where("submission_id = ?", submissionID).Delete(&models.NewJoinerDetail{})

	if len(details) > 0 {
		if err := DB.CreateInBatches(details, 200).Error; err != nil {
			return nil, fmt.Errorf("failed to save new joiner details: %w", err)
		}
	}

	return details, nil
}

// SyncNewJoiners enrols staged NewJoinerDetail records into the live member register.
func SyncNewJoiners(submissionID int, user models.AppUser) (models.MemberSyncResult, error) {
	var result models.MemberSyncResult

	var sub models.EmployerSubmission
	if err := DB.First(&sub, submissionID).Error; err != nil {
		return result, fmt.Errorf("submission not found: %w", err)
	}
	if sub.Status != "accepted" {
		return result, fmt.Errorf("submission must be 'accepted' to sync new joiners (current: %s)", sub.Status)
	}

	var details []models.NewJoinerDetail
	if err := DB.Model(&models.NewJoinerDetail{}).Where("submission_id = ? AND synced_at IS NULL", submissionID).Find(&details).Error; err != nil {
		return result, fmt.Errorf("failed to load staged new joiner details: %w", err)
	}
	if len(details) == 0 {
		return result, fmt.Errorf("no staged new joiner records found — upload a new joiner detail file first")
	}

	now := time.Now()
	for i, d := range details {
		if d.ValidationStatus != "valid" {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("skip '%s': %s", d.MemberName, d.ExclusionReason))
			continue
		}
		member := newJoinerDetailToMemberInForce(d, sub.SchemeID)
		if _, err := AddMemberToScheme(member, user); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("add '%s': %s", d.MemberName, err.Error()))
		} else {
			result.Added++
			details[i].SyncedAt = &now
			DB.Model(&models.NewJoinerDetail{}).Where("id = ?", d.ID).Update("synced_at", now)
		}
	}

	summaryJSON, _ := json.Marshal(result)
	DB.Model(&sub).Updates(map[string]interface{}{
		"new_joiners_synced_at":    now,
		"new_joiners_sync_summary": string(summaryJSON),
	})

	return result, nil
}

// parseNewJoinerFile reads a CSV or XLSX new joiner detail file and returns staged records.
func parseNewJoinerFile(filePath, fileName string) ([]models.NewJoinerDetail, error) {
	lower := strings.ToLower(fileName)
	if strings.HasSuffix(lower, ".xlsx") {
		return parseExcelNewJoiners(filePath)
	}
	if strings.HasSuffix(lower, ".csv") {
		return parseCSVNewJoiners(filePath)
	}
	return nil, fmt.Errorf("unsupported file format: only .csv and .xlsx are accepted")
}

func parseCSVNewJoiners(path string) ([]models.NewJoinerDetail, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	delim := detectDelimiter(f)
	reader := csv.NewReader(f)
	reader.Comma = delim
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}
	return rowsToNewJoinerDetails(rows[0], rows[1:])
}

func parseExcelNewJoiners(path string) ([]models.NewJoinerDetail, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}
	return rowsToNewJoinerDetails(rows[0], rows[1:])
}

// rowsToNewJoinerDetails maps raw header+data rows to NewJoinerDetail structs.
func rowsToNewJoinerDetails(headers []string, dataRows [][]string) ([]models.NewJoinerDetail, error) {
	norm := make([]string, len(headers))
	for i, h := range headers {
		norm[i] = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(h, " ", "_")))
	}

	colIdx := func(candidates ...string) int {
		for _, c := range candidates {
			for i, h := range norm {
				if h == c {
					return i
				}
			}
		}
		return -1
	}

	idxName := colIdx("member_name", "name", "full_name", "fullname")
	idxEmpNo := colIdx("employee_no", "employee_number", "emp_no", "emp_number", "staff_no")
	idxID := colIdx("id_number", "id_no", "sa_id", "id")
	idxIDType := colIdx("id_type", "member_id_type")
	idxGender := colIdx("gender", "sex")
	idxDOB := colIdx("dob", "date_of_birth", "birth_date")
	idxSalary := colIdx("salary", "annual_salary", "gross_salary")
	idxCategory := colIdx("scheme_category", "category")
	idxAddr1 := colIdx("address_line1", "address", "address_line_1", "address1")
	idxAddr2 := colIdx("address_line2", "address_line_2", "address2")
	idxCity := colIdx("city")
	idxProvince := colIdx("province", "state")
	idxPostal := colIdx("postal_code", "zip_code", "post_code")
	idxPhone := colIdx("phone", "phone_number", "contact_number", "mobile")
	idxEmail := colIdx("email", "email_address")
	idxOccupation := colIdx("occupation", "job_title")
	idxOccClass := colIdx("occupational_class", "occ_class")
	idxEntry := colIdx("entry_date", "commencement_date", "start_date")

	getCellStr := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}
	getCellFloat := func(row []string, idx int) float64 {
		s := getCellStr(row, idx)
		if s == "" {
			return 0
		}
		v, _ := strconv.ParseFloat(strings.ReplaceAll(s, ",", ""), 64)
		return v
	}

	var details []models.NewJoinerDetail
	for _, row := range dataRows {
		if len(row) == 0 {
			continue
		}
		name := getCellStr(row, idxName)
		if name == "" {
			continue
		}
		empNo := getCellStr(row, idxEmpNo)
		idNum := getCellStr(row, idxID)

		// Compute member key (mirrors memberKey for EmployerSubmissionRecord)
		key := ""
		if empNo != "" {
			key = "emp:" + empNo
		} else if idNum != "" {
			key = "id:" + idNum
		} else {
			key = "name:" + strings.ToLower(strings.TrimSpace(name))
		}

		validationStatus := "valid"
		exclusionReason := ""
		memberIdType := getCellStr(row, idxIDType)
		if idNum != "" {
			detectedType, valid, reason := ClassifyMemberID(idNum)
			if !valid {
				validationStatus = "id_invalid"
				exclusionReason = reason
			}
			if memberIdType == "" {
				memberIdType = detectedType
			}
		}

		details = append(details, models.NewJoinerDetail{
			MemberKey:         key,
			MemberName:        name,
			EmployeeNumber:    empNo,
			IDNumber:          idNum,
			MemberIdType:      memberIdType,
			Gender:            getCellStr(row, idxGender),
			DateOfBirth:       getCellStr(row, idxDOB),
			AnnualSalary:      getCellFloat(row, idxSalary),
			SchemeCategory:    getCellStr(row, idxCategory),
			AddressLine1:      getCellStr(row, idxAddr1),
			AddressLine2:      getCellStr(row, idxAddr2),
			City:              getCellStr(row, idxCity),
			Province:          getCellStr(row, idxProvince),
			PostalCode:        getCellStr(row, idxPostal),
			PhoneNumber:       getCellStr(row, idxPhone),
			Email:             getCellStr(row, idxEmail),
			Occupation:        getCellStr(row, idxOccupation),
			OccupationalClass: getCellStr(row, idxOccClass),
			EntryDate:         getCellStr(row, idxEntry),
			ValidationStatus:  validationStatus,
			ExclusionReason:   exclusionReason,
		})
	}
	return details, nil
}

// newJoinerDetailToMemberInForce maps a staged NewJoinerDetail to a GPricingMemberDataInForce
// ready for AddMemberToScheme.
func newJoinerDetailToMemberInForce(d models.NewJoinerDetail, schemeID int) models.GPricingMemberDataInForce {
	m := models.GPricingMemberDataInForce{
		SchemeId:          schemeID,
		MemberName:        d.MemberName,
		MemberIdNumber:    d.IDNumber,
		MemberIdType:      d.MemberIdType,
		Gender:            d.Gender,
		EmployeeNumber:    d.EmployeeNumber,
		AnnualSalary:      d.AnnualSalary,
		SchemeCategory:    d.SchemeCategory,
		AddressLine1:      d.AddressLine1,
		AddressLine2:      d.AddressLine2,
		City:              d.City,
		Province:          d.Province,
		PostalCode:        d.PostalCode,
		PhoneNumber:       d.PhoneNumber,
		Email:             d.Email,
		Occupation:        d.Occupation,
		OccupationalClass: d.OccupationalClass,
	}
	if d.DateOfBirth != "" {
		if t, err := parseFlexDate(d.DateOfBirth); err == nil {
			m.DateOfBirth = t
		}
	}
	if d.EntryDate != "" {
		if t, err := parseFlexDate(d.EntryDate); err == nil {
			m.EntryDate = t
		}
	} else {
		m.EntryDate = time.Now()
	}
	return m
}

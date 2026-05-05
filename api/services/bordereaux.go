package services

import (
	appLog "api/log"
	"api/models"
	"api/utils"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"

	"gorm.io/gorm"
)

// newFileToken returns a short, collision-free token suitable for embedding in
// a generated bordereaux filename. Replaces earlier time.Now().Unix()/UnixNano()
// suffixes which would collide under concurrent writes.
func newFileToken() string {
	// First 12 hex chars of a v4 UUID (~48 bits) — plenty of uniqueness for
	// filenames without producing an unwieldy path.
	return strings.ReplaceAll(uuid.NewString(), "-", "")[:12]
}

// ImportSchemeConfirmationsResult represents the result of the import and reconciliation process
type ImportSchemeConfirmationsResult struct {
	Success             bool                            `json:"success"`
	Message             string                          `json:"message"`
	Error               string                          `json:"error,omitempty"`
	ImportedCount       int                             `json:"imported_count"`
	AutoProcessResults  *ReconciliationSummary          `json:"auto_process_results,omitempty"`
	ConfirmationRecords []models.BordereauxConfirmation `json:"confirmation_records,omitempty"`
}

type ReconciliationSummary struct {
	Matched       int `json:"matched"`
	Discrepancies int `json:"discrepancies"`
	Pending       int `json:"pending"`
}

// ImportSchemeConfirmations handles the import of confirmation files and optional reconciliation
func ImportSchemeConfirmations(ctx context.Context, fileType string, autoProcess bool, files []*multipart.FileHeader, importedBy string) (*ImportSchemeConfirmationsResult, error) {
	logger := appLog.WithContext(ctx)
	result := &ImportSchemeConfirmationsResult{
		Success: true,
		Message: "Files imported successfully",
	}

	var confirmations []models.BordereauxConfirmation

	// Pre-fetch all schemes to avoid N+1 queries during import
	var allSchemes []models.GroupScheme
	if err := DB.Find(&allSchemes).Error; err != nil {
		logger.WithError(err).Warn("Failed to pre-fetch schemes")
	}
	schemeMap := make(map[string]int)
	for _, s := range allSchemes {
		schemeMap[strings.ToLower(s.Name)] = s.ID
	}

	for _, fileHeader := range files {
		// Save file to a temporary location or storage
		uploadDir := "tmp/uploads/confirmations"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %w", err)
		}

		filePath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename))

		src, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("failed to copy file content: %w", err)
		}

		// Read file content
		records, err := readConfirmationFile(filePath, fileType)
		if err != nil {
			return nil, fmt.Errorf("failed to read confirmation file: %w", err)
		}

		if len(records) == 0 {
			continue
		}

		// Group records by Bordereaux ID
		// Map key is Bordereaux ID (or "unknown" if not found)
		groupedRecords := make(map[string][]map[string]any)

		for _, r := range records {
			bdxID := ""
			for k, v := range r {
				kl := strings.ToLower(k)
				if (strings.Contains(kl, "bordereaux") || strings.Contains(kl, "bdx")) && strings.Contains(kl, "id") {
					bdxID = fmt.Sprint(v)
					break
				}
			}
			if bdxID == "" {
				bdxID = "unknown"
			}
			groupedRecords[bdxID] = append(groupedRecords[bdxID], r)
		}

		// Process each group of records as a separate confirmation
		for bdxID, group := range groupedRecords {
			// Extract metadata (Month, Year) from the first record in the group
			extractedMonth := 0
			extractedYear := 0
			// Extract metadata (Scheme Name) from the first record in the group
			schemeName := ""
			currentSchemeID := 0
			if len(group) > 0 {
				first := group[0]
				for k, v := range first {
					kl := strings.ToLower(k)
					if strings.Contains(kl, "month") {
						mStr := fmt.Sprint(v)
						m, _ := strconv.Atoi(mStr)
						if m >= 1 && m <= 12 {
							extractedMonth = m
						} else {
							months := map[string]int{"january": 1, "february": 2, "march": 3, "april": 4, "may": 5, "june": 6, "july": 7, "august": 8, "september": 9, "october": 10, "november": 11, "december": 12}
							if val, ok := months[strings.ToLower(mStr)]; ok {
								extractedMonth = val
							}
						}
					}
					if strings.Contains(kl, "year") {
						extractedYear, _ = strconv.Atoi(fmt.Sprint(v))
					}
					if strings.Contains(kl, "scheme") {
						schemeName = fmt.Sprint(v)
						currentSchemeID = schemeMap[strings.ToLower(schemeName)]
					}
				}
			}

			// Find matching generated bordereaux
			var targetBordereaux models.GeneratedBordereaux
			if bdxID != "unknown" {
				DB.Where("generated_id = ?", bdxID).First(&targetBordereaux)
			}

			// If still not found by ID, try heuristic search by scheme and period
			if targetBordereaux.ID == 0 && currentSchemeID > 0 {
				query := DB.Where("scheme_name = (SELECT name FROM group_schemes WHERE id = ?) AND status IN ('submitted', 'generated')", currentSchemeID)
				if extractedMonth > 0 && extractedYear > 0 {
					periodDate := time.Date(extractedYear, time.Month(extractedMonth), 1, 0, 0, 0, 0, time.UTC)
					query = query.Where("period_start = ?", periodDate)
				}
				if err := query.Order("submission_date desc, created_at desc").First(&targetBordereaux).Error; err != nil {
					logger.Warnf("Could not find a matching bordereaux for scheme %d and bdxID %s to link with confirmation", currentSchemeID, bdxID)
				}
			}

			confirmation := models.BordereauxConfirmation{
				GeneratedBordereauxID: targetBordereaux.GeneratedID,
				SchemeID:              currentSchemeID,
				//SchemeName:            resolveSchemeNames([]int{schemeID}),
				SchemeName:     schemeName,
				FileName:       fileHeader.Filename,
				FilePath:       filePath,
				FileType:       fileType,
				ValuationMonth: extractedMonth,
				ValuationYear:  extractedYear,
				Status:         "pending",
				ImportedBy:     importedBy,
				ImportedAt:     time.Now(),
			}

			if confirmErr := DB.Create(&confirmation).Error; confirmErr != nil {
				logger.WithError(confirmErr).Error("failed to create confirmation record")
				continue
			}

			// get the associated template for the target bordereaux
			template, tmplErr := getTemplateForBordereaux(targetBordereaux)
			if tmplErr != nil {
				logger.WithError(tmplErr).Error("failed to retrieve template for bordereaux")
				continue
			} else {
				fmt.Println(template.FieldMappings)
			}

			// Persist confirmation records for this group
			var confirmationRecords []models.BordereauxConfirmationRecord
			for _, r := range group {
				confRec := models.BordereauxConfirmationRecord{
					BordereauxConfirmationID: confirmation.ID,
					BordereauxID:             bdxID,
					ValuationMonth:           extractedMonth,
					ValuationYear:            extractedYear,
					CreatedAt:                time.Now(),
				}

				// Map common fields
				for k, v := range r {
					kl := strings.ToLower(k)
					vStr := fmt.Sprint(v)
					if strings.Contains(kl, "member") && strings.Contains(kl, "id") {
						confRec.MemberID = vStr
					} else if strings.Contains(kl, "member_name") {
						if confRec.MemberName == "" || strings.Contains(kl, "full") {
							confRec.MemberName = vStr
						}
					} else if strings.Contains(kl, "id") && (strings.Contains(kl, "number") || strings.Contains(kl, "national")) {
						confRec.IDNumber = vStr
					} else if strings.Contains(kl, "amount") || strings.Contains(kl, "total_annual_premium") {
						amt, _ := strconv.ParseFloat(vStr, 64)
						if amt > 0 {
							confRec.Amount = amt
						}
					} else if strings.Contains(kl, "scheme") {
						confRec.SchemeName = vStr
					} else if strings.Contains(kl, "employee_number") || strings.Contains(kl, "employee") {
						confRec.EmployeeNumber = vStr
					}
				}

				rawJSON, _ := json.Marshal(r)
				confRec.RawData = models.JSON(rawJSON)
				confirmationRecords = append(confirmationRecords, confRec)
			}

			// Persist confirmation records and (optionally) reconcile atomically so
			// a mid-batch failure cannot leave orphaned confirmation rows with no
			// reconciliation results.
			var importSummary *ReconciliationSummary
			txErr := DB.Transaction(func(tx *gorm.DB) error {
				if len(confirmationRecords) > 0 {
					if err := tx.CreateInBatches(confirmationRecords, 500).Error; err != nil {
						return fmt.Errorf("failed to persist confirmation records: %w", err)
					}
				}
				if autoProcess {
					summary, err := ReconcileConfirmation(ctx, tx, &confirmation)
					if err != nil {
						return err
					}
					importSummary = summary
				}
				return nil
			})
			if txErr != nil {
				logger.WithError(txErr).Errorf("Confirmation import transaction failed for confirmation %d", confirmation.ID)
				confirmation.Status = "error"
				DB.Save(&confirmation)
			} else if autoProcess && importSummary != nil {
				if result.AutoProcessResults == nil {
					result.AutoProcessResults = importSummary
				} else {
					result.AutoProcessResults.Matched += importSummary.Matched
					result.AutoProcessResults.Discrepancies += importSummary.Discrepancies
					result.AutoProcessResults.Pending += importSummary.Pending
				}
			}

			confirmations = append(confirmations, confirmation)
			result.ImportedCount++
		}
	}

	result.ConfirmationRecords = confirmations
	return result, nil
}

// defaultReconciliationTolerance is the per-record variance threshold used when
// a scheme has not configured its own value. Set conservatively so legacy rows
// continue to match at the old precision.
const defaultReconciliationTolerance = 0.001

// reconciliationToleranceForScheme returns the scheme-specific tolerance from
// group_schemes.reconciliation_tolerance, falling back to the codebase default
// on zero / missing rows. Uses the supplied db handle so callers inside a
// transaction see the in-flight state.
func reconciliationToleranceForScheme(db *gorm.DB, schemeID int) float64 {
	if schemeID <= 0 {
		return defaultReconciliationTolerance
	}
	var tolerance float64
	if err := db.Model(&models.GroupScheme{}).
		Where("id = ?", schemeID).
		Select("reconciliation_tolerance").
		Row().Scan(&tolerance); err != nil {
		return defaultReconciliationTolerance
	}
	if tolerance <= 0 {
		return defaultReconciliationTolerance
	}
	return tolerance
}

func getTemplateForBordereaux(bordereaux models.GeneratedBordereaux) (models.BordereauxTemplate, error) {
	var template models.BordereauxTemplate

	if err := DB.Where("id = ?", bordereaux.TemplateID).First(&template).Error; err != nil {
		return template, fmt.Errorf("failed to retrieve template for bordereaux: %w", err)
	}

	return template, nil
}

// ReconcileConfirmation performs the matching between confirmation records and bordereaux.
// The db parameter allows callers to pass either the package-level DB or an active
// transaction so confirmation persistence and reconciliation results commit together.
func ReconcileConfirmation(ctx context.Context, db *gorm.DB, confirmation *models.BordereauxConfirmation) (*ReconciliationSummary, error) {
	if db == nil {
		db = DB
	}
	// 1. Fetch confirmation records from DB (persisted earlier)
	var records []models.BordereauxConfirmationRecord
	if err := db.Where("bordereaux_confirmation_id = ?", confirmation.ID).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch confirmation records: %w", err)
	}

	summary := &ReconciliationSummary{}

	// 2. Fetch the original bordereaux
	var bordereaux models.GeneratedBordereaux
	if err := db.Where("generated_id = ?", confirmation.GeneratedBordereauxID).First(&bordereaux).Error; err != nil {
		// If we don't have a linked bordereaux, we can't do full reconciliation
		return summary, nil
	}

	// Per-scheme tolerance overrides the codebase default when configured.
	tolerance := reconciliationToleranceForScheme(db, confirmation.SchemeID)

	// 3. Fetch original bordereaux data
	var submittedAmount float64
	var totalExpected int

	// Maps for different types
	premiumExpectedMap := make(map[string]models.PremiumBordereauxData)
	claimExpectedMap := make(map[string]models.GroupSchemeClaim)
	memberExpectedMap := make(map[string]models.MemberBordereauxData)

	if bordereaux.Type == "premium" {
		var expectedRecords []models.PremiumBordereauxData
		db.Where("bordereaux_id = ?", bordereaux.GeneratedID).Find(&expectedRecords)
		totalExpected = len(expectedRecords)
		for _, r := range expectedRecords {
			submittedAmount += r.TotalAnnualPremium
			key := r.EmployeeNumber
			if key == "" {
				key = r.MemberIdNumber
			}
			if key == "" {
				key = r.MemberName
			}
			premiumExpectedMap[key] = r
		}
	} else if bordereaux.Type == "claim" {
		var expectedRecords []models.GroupSchemeClaim
		// Using BordereauxID if it's stored in GroupSchemeClaim, or filtering by scheme and period
		// Looking at GroupSchemeClaim model, it doesn't have BordereauxID but has SchemeId.
		// Usually bordereaux generation for claims would mark them or we filter by period.
		db.Where("scheme_id = ? AND creation_date BETWEEN ? AND ?", bordereaux.ID, bordereaux.PeriodStart, bordereaux.PeriodEnd).Find(&expectedRecords)
		totalExpected = len(expectedRecords)
		for _, r := range expectedRecords {
			submittedAmount += r.ClaimAmount
			key := r.MemberIDNumber
			if key == "" {
				key = r.MemberName
			}
			claimExpectedMap[key] = r
		}
	} else if bordereaux.Type == "member" {
		// Member bordereaux are census/roster data — there is no per-member amount to
		// reconcile, so variance is identity-based (matched / missing / extra) against
		// the snapshot persisted at generation time.
		var expectedRecords []models.MemberBordereauxData
		db.Where("bordereaux_id = ?", bordereaux.GeneratedID).Find(&expectedRecords)
		totalExpected = len(expectedRecords)
		for _, r := range expectedRecords {
			key := r.EmployeeNumber
			if key == "" {
				key = r.MemberIdNumber
			}
			if key == "" {
				key = r.MemberName
			}
			memberExpectedMap[key] = r
		}
		submittedAmount = 0 // member bordereaux carry no reconcilable amount
	}

	confirmation.SubmittedAmount = submittedAmount
	confirmedAmount := 0.0
	matchedCount := 0
	discrepancyCount := 0

	var reconciliationResults []models.BordereauxReconciliationResult

	// 4. Match and compare
	for _, rec := range records {
		confirmedAmount += rec.Amount
		key := rec.EmployeeNumber
		if key == "" {
			key = rec.MemberName
		}

		if bordereaux.Type == "premium" {
			if expected, ok := premiumExpectedMap[key]; ok {
				//matchedCount++
				variance := rec.Amount - expected.TotalAnnualPremium
				//status := "matched"
				var status string
				if math.Abs(variance) > tolerance {
					discrepancyCount++
					status = "discrepancy"
				} else {
					matchedCount++
					status = "matched"
				}
				reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
					BordereauxConfirmationID: confirmation.ID,
					GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
					RecordID:                 key,
					MemberName:               rec.MemberName,
					Field:                    "Amount",
					ExpectedValue:            fmt.Sprintf("%.2f", expected.TotalAnnualPremium),
					ActualValue:              fmt.Sprintf("%.2f", rec.Amount),
					Variance:                 variance,
					Status:                   status,
					CreatedAt:                time.Now(),
				})
				delete(premiumExpectedMap, key)
			} else {
				discrepancyCount++
				reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
					BordereauxConfirmationID: confirmation.ID,
					GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
					RecordID:                 key,
					MemberName:               rec.MemberName,
					Status:                   "extra",
					Comments:                 "Record not found in submitted bordereaux",
					CreatedAt:                time.Now(),
				})
			}
		} else if bordereaux.Type == "claim" {
			if expected, ok := claimExpectedMap[key]; ok {
				matchedCount++
				variance := rec.Amount - expected.ClaimAmount
				status := "matched"
				if math.Abs(variance) > tolerance {
					discrepancyCount++
					status = "discrepancy"
				}
				reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
					BordereauxConfirmationID: confirmation.ID,
					GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
					RecordID:                 key,
					MemberName:               rec.MemberName,
					Field:                    "Amount",
					ExpectedValue:            fmt.Sprintf("%.2f", expected.ClaimAmount),
					ActualValue:              fmt.Sprintf("%.2f", rec.Amount),
					Variance:                 variance,
					Status:                   status,
					CreatedAt:                time.Now(),
				})
				delete(claimExpectedMap, key)
			} else {
				discrepancyCount++
				reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
					BordereauxConfirmationID: confirmation.ID,
					GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
					RecordID:                 key,
					MemberName:               rec.MemberName,
					Status:                   "extra",
					Comments:                 "Record not found in submitted bordereaux",
					CreatedAt:                time.Now(),
				})
			}
		} else if bordereaux.Type == "member" {
			// Census bordereaux: match incoming rows against the snapshot by identity
			// (employee number → ID number → name). No amount comparison.
			memberKey := rec.EmployeeNumber
			if memberKey == "" {
				memberKey = rec.IDNumber
			}
			if memberKey == "" {
				memberKey = rec.MemberName
			}
			if _, ok := memberExpectedMap[memberKey]; ok {
				matchedCount++
				reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
					BordereauxConfirmationID: confirmation.ID,
					GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
					RecordID:                 memberKey,
					MemberName:               rec.MemberName,
					Status:                   "matched",
					CreatedAt:                time.Now(),
				})
				delete(memberExpectedMap, memberKey)
			} else {
				discrepancyCount++
				reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
					BordereauxConfirmationID: confirmation.ID,
					GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
					RecordID:                 memberKey,
					MemberName:               rec.MemberName,
					Status:                   "extra",
					Comments:                 "Member not found in submitted bordereaux",
					CreatedAt:                time.Now(),
				})
			}
		}
	}

	// 5. Check for missing records
	if bordereaux.Type == "premium" {
		for key, expected := range premiumExpectedMap {
			discrepancyCount++
			reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
				BordereauxConfirmationID: confirmation.ID,
				GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
				RecordID:                 key,
				MemberName:               expected.MemberName,
				Status:                   "missing",
				Comments:                 "Record missing from confirmation file",
				CreatedAt:                time.Now(),
			})
		}
	} else if bordereaux.Type == "claim" {
		for key, expected := range claimExpectedMap {
			discrepancyCount++
			reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
				BordereauxConfirmationID: confirmation.ID,
				GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
				RecordID:                 key,
				MemberName:               expected.MemberName,
				Status:                   "missing",
				Comments:                 "Record missing from confirmation file",
				CreatedAt:                time.Now(),
			})
		}
	} else if bordereaux.Type == "member" {
		for key, expected := range memberExpectedMap {
			discrepancyCount++
			reconciliationResults = append(reconciliationResults, models.BordereauxReconciliationResult{
				BordereauxConfirmationID: confirmation.ID,
				GeneratedBordereauxID:    confirmation.GeneratedBordereauxID,
				RecordID:                 key,
				MemberName:               expected.MemberName,
				Status:                   "missing",
				Comments:                 "Member missing from confirmation file",
				CreatedAt:                time.Now(),
			})
		}
	}

	// Batch insert reconciliation results
	if len(reconciliationResults) > 0 {
		if err := db.CreateInBatches(reconciliationResults, 500).Error; err != nil {
			return nil, fmt.Errorf("failed to persist reconciliation results: %w", err)
		}
	}

	// Update confirmation stats
	now := time.Now()
	confirmation.LastReconciled = &now
	confirmation.ConfirmedAmount = confirmedAmount
	confirmation.Variance = confirmedAmount - submittedAmount
	confirmation.MatchedCount = matchedCount
	confirmation.DiscrepancyCount = discrepancyCount

	if totalExpected > 0 {
		confirmation.MatchScore = (float64(matchedCount) / float64(totalExpected)) * 100
	} else if len(records) == 0 {
		confirmation.MatchScore = 100
	} else {
		confirmation.MatchScore = 0
	}

	// Determine status for the table
	if confirmation.MatchScore >= 100 && discrepancyCount == 0 {
		confirmation.Status = "matched"
	} else if discrepancyCount > 0 {
		confirmation.Status = "discrepancy"
	} else {
		confirmation.Status = "pending"
	}

	// Save the updated confirmation

	if err := db.Save(confirmation).Error; err != nil {
		return nil, fmt.Errorf("failed to update confirmation: %w", err)
	}

	// Update the original bordereaux status and progress to confirmed
	if bordereaux.GeneratedID != "" {
		updates := map[string]interface{}{
			"status":       "confirmed",
			"progress":     100,
			"last_updated": time.Now(),
		}
		if err := db.Model(&bordereaux).Updates(updates).Error; err != nil {
			appLog.WithFields(map[string]interface{}{
				"generated_id": bordereaux.GeneratedID,
				"error":        err.Error(),
			}).Error("Failed to update bordereaux status to confirmed")
		}
	}

	// Add timeline entry
	AddBordereauxTimelineEntry(confirmation.GeneratedBordereauxID, models.BordereauxTimeline{
		Date:        time.Now(),
		Type:        "reconciled",
		Title:       "Reconciliation Completed",
		Description: fmt.Sprintf("Reconciliation completed with %d matches and %d discrepancies. Match score: %.2f%%", matchedCount, discrepancyCount, confirmation.MatchScore),
	})

	summary.Matched = matchedCount
	summary.Discrepancies = discrepancyCount
	summary.Pending = totalExpected - matchedCount - discrepancyCount
	if summary.Pending < 0 {
		summary.Pending = 0
	}

	return summary, nil
}

func GetReconciliationStats(schemeID int) (*ReconciliationSummary, error) {
	var summary ReconciliationSummary

	var confirmations []models.BordereauxConfirmation
	if err := DB.Where("scheme_id = ?", schemeID).Find(&confirmations).Error; err != nil {
		return nil, err
	}

	for _, c := range confirmations {
		summary.Matched += c.MatchedCount
		summary.Discrepancies += c.DiscrepancyCount
	}

	// For pending, it's a bit more complex.
	// If we define pending as records in our system that haven't been confirmed yet...
	// but the ReconciliationSummary seems to treat 'Pending' as unmatched records from the file.
	// Looking at the frontend code:
	// reconciliationStats.value.pending -= (matched + discrepancies)
	// This implies pending starts as some total and is reduced.

	// Let's count how many 'missing' records we have across all confirmations for this scheme.
	var missingCount int64
	DB.Model(&models.BordereauxReconciliationResult{}).
		Joins("JOIN bordereaux_confirmations ON bordereaux_confirmations.id = bordereaux_reconciliation_results.bordereaux_confirmation_id").
		Where("bordereaux_confirmations.scheme_id = ? AND bordereaux_reconciliation_results.status = ?", schemeID, "missing").
		Count(&missingCount)

	summary.Pending = int(missingCount)

	return &summary, nil
}

func GetBordereauxConfirmations() ([]map[string]any, error) {
	var confirmations []models.BordereauxConfirmation
	if err := DB.Order("imported_at desc").Find(&confirmations).Error; err != nil {
		return nil, err
	}

	var results []map[string]any
	for _, c := range confirmations {
		item := map[string]any{
			"id":                      c.ID,
			"generated_bordereaux_id": c.GeneratedBordereauxID,
			"scheme_id":               c.SchemeID,
			"scheme_name":             c.SchemeName,
			"file_name":               c.FileName,
			"status":                  c.Status,
			"matched_count":           c.MatchedCount,
			"discrepancy_count":       c.DiscrepancyCount,
			"submitted_amount":        c.SubmittedAmount,
			"confirmed_amount":        c.ConfirmedAmount,
			"variance":                c.Variance,
			"match_score":             c.MatchScore,
			"last_reconciled":         c.LastReconciled,
			"imported_at":             c.ImportedAt,
			"imported_by":             c.ImportedBy,
		}

		// Fetch linked bordereaux info for type
		var bordereaux models.GeneratedBordereaux
		if err := DB.Where("generated_id = ?", c.GeneratedBordereauxID).First(&bordereaux).Error; err == nil {
			item["type"] = bordereaux.Type
			if item["scheme_name"] == "" {
				item["scheme_name"] = bordereaux.SchemeName
			}
		}

		results = append(results, item)
	}

	return results, nil
}

func GetReconciliationResults(confirmationID int) ([]models.BordereauxReconciliationResult, error) {
	var results []models.BordereauxReconciliationResult
	err := DB.Where("bordereaux_confirmation_id = ?", confirmationID).Find(&results).Error
	return results, err
}

func GetUnmatchedReconciliationResults(confirmationID int) ([]models.BordereauxReconciliationResult, error) {
	var results []models.BordereauxReconciliationResult
	err := DB.Where("bordereaux_confirmation_id = ? AND status != ?", confirmationID, "matched").Find(&results).Error
	return results, err
}

func GetBordereauxConfirmation(id int) (map[string]any, error) {
	var confirmation models.BordereauxConfirmation
	if err := DB.First(&confirmation, id).Error; err != nil {
		return nil, err
	}

	results, err := GetReconciliationResults(id)
	if err != nil {
		return nil, err
	}

	item := map[string]any{
		"id":                      confirmation.ID,
		"generated_bordereaux_id": confirmation.GeneratedBordereauxID,
		"scheme_id":               confirmation.SchemeID,
		"scheme_name":             confirmation.SchemeName,
		"file_name":               confirmation.FileName,
		"file_path":               confirmation.FilePath,
		"file_type":               confirmation.FileType,
		"status":                  confirmation.Status,
		"matched_count":           confirmation.MatchedCount,
		"discrepancy_count":       confirmation.DiscrepancyCount,
		"submitted_amount":        confirmation.SubmittedAmount,
		"confirmed_amount":        confirmation.ConfirmedAmount,
		"variance":                confirmation.Variance,
		"match_score":             confirmation.MatchScore,
		"last_reconciled":         confirmation.LastReconciled,
		"imported_at":             confirmation.ImportedAt,
		"imported_by":             confirmation.ImportedBy,
		"reconciliation_results":  results,
	}

	// Fetch linked bordereaux info for type
	var bordereaux models.GeneratedBordereaux
	if err := DB.Where("generated_id = ?", confirmation.GeneratedBordereauxID).First(&bordereaux).Error; err == nil {
		item["type"] = bordereaux.Type
		if item["scheme_name"] == "" {
			item["scheme_name"] = bordereaux.SchemeName
		}
	}

	return item, nil
}

func DeleteBordereauxConfirmation(id int) error {
	var confirmation models.BordereauxConfirmation
	if err := DB.First(&confirmation, id).Error; err != nil {
		return err
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		// 1. Delete reconciliation results
		if err := tx.Where("bordereaux_confirmation_id = ?", id).Delete(&models.BordereauxReconciliationResult{}).Error; err != nil {
			return err
		}

		// 2. Delete confirmation records
		if err := tx.Where("bordereaux_confirmation_id = ?", id).Delete(&models.BordereauxConfirmationRecord{}).Error; err != nil {
			return err
		}

		// 3. Delete the confirmation itself
		if err := tx.Delete(&confirmation).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 4. Delete physical file if exists
	if confirmation.FilePath != "" {
		if _, err := os.Stat(confirmation.FilePath); err == nil {
			os.Remove(confirmation.FilePath)
		}
	}

	return nil
}

// ReconcilePendingConfirmations finds all confirmations with 'pending' status and runs reconciliation on them
func ReconcilePendingConfirmations(ctx context.Context) (int, error) {
	var pendingConfirmations []models.BordereauxConfirmation
	if err := DB.Where("status = ?", "pending").Find(&pendingConfirmations).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch pending confirmations: %w", err)
	}

	processedCount := 0
	for i := range pendingConfirmations {
		txErr := DB.Transaction(func(tx *gorm.DB) error {
			_, err := ReconcileConfirmation(ctx, tx, &pendingConfirmations[i])
			return err
		})
		if txErr != nil {
			appLog.WithContext(ctx).WithError(txErr).Errorf("Batch reconciliation failed for confirmation %d", pendingConfirmations[i].ID)
			pendingConfirmations[i].Status = "error"
			DB.Save(&pendingConfirmations[i])
			continue
		}
		processedCount++
	}

	return processedCount, nil
}

func readConfirmationFile(path string, fileType string) ([]map[string]any, error) {
	lowerPath := strings.ToLower(path)

	// Handle gzipped files
	if strings.HasSuffix(lowerPath, ".gz") {
		gzFile, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open gz file: %w", err)
		}
		defer gzFile.Close()

		zr, err := gzip.NewReader(gzFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer zr.Close()

		// Create a temporary file for the decompressed content
		// We need a path for excelize and for extension checking
		tmpName := strings.TrimSuffix(path, ".gz")
		// Ensure it has an extension if it didn't
		if !strings.Contains(filepath.Base(tmpName), ".") {
			tmpName += "." + fileType
		}

		tmpFile, err := os.CreateTemp("", "bdx_confirm_*_"+filepath.Base(tmpName))
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file for decompression: %w", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		if _, err := io.Copy(tmpFile, zr); err != nil {
			return nil, fmt.Errorf("failed to decompress gz file: %w", err)
		}
		// Close to ensure all data is written and available for reading
		tmpFile.Close()

		return readConfirmationFile(tmpFile.Name(), fileType)
	}

	if strings.HasSuffix(lowerPath, ".xlsx") {
		return readExcelConfirmation(path)
	}
	if strings.HasSuffix(lowerPath, ".csv") {
		return readCSVConfirmation(path)
	}
	return nil, fmt.Errorf("unsupported file format")
}

func readCSVConfirmation(path string) ([]map[string]any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Some CSVs might have different delimiters
	// For now, assume comma. We could add logic to detect it if needed.

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, nil
	}

	headers := rows[0]
	var records []map[string]any
	for _, row := range rows[1:] {
		record := make(map[string]any)
		for i, cell := range row {
			if i < len(headers) {
				record[headers[i]] = cell
			}
		}
		records = append(records, record)
	}
	return records, nil
}

func readExcelConfirmation(path string) ([]map[string]any, error) {
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

	headers := rows[0]
	var records []map[string]any
	for _, row := range rows[1:] {
		record := make(map[string]any)
		for i, cell := range row {
			if i < len(headers) {
				record[headers[i]] = cell
			}
		}
		records = append(records, record)
	}
	return records, nil
}

// GetBordereauxFieldsByType returns available source fields for a given bordereaux type
// The response is a list of objects: {"display_name": "Field Name", "field_name": "field_name"}
func GetBordereauxFieldsByType(bType string) ([]map[string]string, error) {
	t := strings.ToLower(strings.TrimSpace(bType))
	var out []map[string]string
	// helper to add
	add := func(display, name string) {
		out = append(out, map[string]string{"display_name": display, "field_name": name})
	}

	switch t {
	case "member", "members":
		add("Member Name", "member_name")
		add("Employee Number", "employee_number")
		add("Annual Salary", "annual_salary")
		add("Entry Date", "entry_date")
		add("Exit Date", "exit_date")
		add("GLA Salary Multiple", "gla_salary_multiple")
		add("PTD Salary Multiple", "ptd_salary_multiple")
		add("CI Salary Multiple", "ci_salary_multiple")
		add("SGLA Salary Multiple", "sgla_salary_multiple")
		add("GLA Covered Sum Assured", "gla_covered_sum_assured")
		add("PTD Covered Sum Assured", "ptd_capped_sum_assured")
		add("CI Covered Sum Assured", "ci_capped_sum_assured")
		add("SGLA Covered Sum Assured", "spouse_gla_capped_sum_assured")
		add("TTD Covered Sum Assured", "ttd_capped_income")
		add("PHI Covered Sum Assured", "phi_capped_income")
		add("GLA Ceded Sum Assured", "gla_ceded_sum_assured")
		add("PTD Ceded Sum Assured", "ptd_ceded_sum_assured")
		add("CI Ceded Sum Assured", "ci_ceded_sum_assured")
		add("SGLA Ceded Sum Assured", "spouse_gla_ceded_sum_assured")
		add("TTD Ceded Income", "ttd_ceded_income")
		add("PHI Ceded Income", "phi_ceded_income")
		add("Member Funeral Sum Assured", "member_funeral_sum_assured")
		add("Spouse Funeral Sum Assured", "spouse_funeral_sum_assured")
		add("Child Funeral Sum Assured", "child_funeral_sum_assured")
		add("Dependant Funeral Sum Assured", "dependant_funeral_sum_assured")
		add("Parent Funeral Sum Assured", "parent_funeral_sum_assured")
		add("Member Funeral Ceded Sum Assured", "member_funeral_ceded_sum_assured")
		add("Spouse Funeral Ceded Sum Assured", "spouse_funeral_ceded_sum_assured")
		add("Child Funeral Ceded Sum Assured", "child_funeral_ceded_sum_assured")
		add("Dependant Funeral Ceded Sum Assured", "dependant_funeral_ceded_sum_assured")
		add("Parent Funeral Ceded Sum Assured", "parent_funeral_ceded_sum_assured")

		// extend with additional member fields as needed
	// Premium / claim / reinsurance field catalogues temporarily disabled — only
	// member is a valid use case right now. Uncomment to re-enable.
	// case "premium", "premiums", "reinsurance_premium", "reinsurance_premiums":
	// 	add("Scheme Name", "scheme_name")
	// 	add("Member Name", "member_name")
	// 	add("Employee Number", "employee_number")
	// 	add("Member ID Number", "member_id_number")
	// 	add("Category", "category")
	// 	add("Annual Salary", "annual_salary")
	// 	add("Valuation Month", "valuation_month")
	// 	add("Age at Month", "age_at_month")
	// 	add("GLA Annual Premium", "gla_annual_premium")
	// 	add("SGLA Annual Premium", "sgla_annual_premium")
	// 	add("PTD Annual Premium", "ptd_annual_premium")
	// 	add("CI Annual Premium", "ci_annual_premium")
	// 	add("TTD Annual Premium", "ttd_annual_premium")
	// 	add("PHI Annual Premium", "phi_annual_premium")
	// 	add("Total Annual Premium Excl Funeral", "total_annual_premium_excl_funeral")
	// 	add("Total Annual Funeral Premium", "total_annual_funeral_premium")
	// 	add("Total Annual Premium", "total_annual_premium")
	// 	//add("Main Member Funeral Annual Premium", "main_member_funeral_annual_premium")
	// 	//add("Spouse Funeral Annual Premium", "spouse_funeral_annual_premium")
	// 	//add("Children Funeral Annual Premium", "children_funeral_annual_premium")
	// 	//add("Dependants Funeral Annual Premium", "dependants_funeral_annual_premium")
	// 	//add("Total Annual Premium Payable", "total_annual_premium_payable")
	// case "claim", "claims", "reinsurance_claim", "reinsurance_claims":
	// 	add("Claim Number", "claim_number")
	// 	add("Employee Number", "employee_number")
	// 	add("Member Name", "member_name")
	// 	add("Member ID Number", "member_id_number")
	// 	add("Scheme ID", "scheme_id")
	// 	add("Scheme Name", "scheme_name")
	// 	add("Benefit Name", "benefit_name")
	// 	add("Member Type", "member_type")
	// 	add("Date Of Event", "date_of_event")
	// 	add("Date Notified", "date_notified")
	// 	add("Status", "status")
	// 	add("Claim Amount", "claim_amount")
	// 	add("Date Registered", "date_registered")
	// 	add("Claimant Name", "claimant_name")
	// 	add("Claimant ID Number", "claimant_id_number")
	// 	add("Relationship To Member", "relationship_to_member")
	default:
		return nil, fmt.Errorf("unsupported bordereaux type: %s", bType)
	}
	return out, nil
}

func resolveSchemeNames(schemeIDs []int) string {
	if len(schemeIDs) == 0 {
		return ""
	}
	var names []string
	if err := DB.Model(&models.GroupScheme{}).Where("id IN ?", schemeIDs).Pluck("name", &names).Error; err != nil {
		return ""
	}
	return strings.Join(names, ", ")
}

func resolveInsurerName(insurerID int) string {
	if insurerID == 0 {
		return ""
	}
	var insurer models.GroupPricingInsurerDetail
	if err := DB.First(&insurer, insurerID).Error; err != nil {
		return ""
	}
	return insurer.Name
}

// GenerateBordereauxRequest matches the UI payload
type GenerateBordereauxRequest struct {
	Type                   string `json:"type"`
	SchemeIDs              []int  `json:"scheme_ids"`
	PeriodType             string `json:"period_type"`
	StartDate              string `json:"start_date"`
	EndDate                string `json:"end_date"`
	Month                  int    `json:"month"`
	Year                   int    `json:"year"`
	TemplateID             int    `json:"template_id"`
	OutputFormat           string `json:"output_format"`
	IncludeTerminated      bool   `json:"include_terminated"`
	IncludeDependants      bool   `json:"include_dependants"`
	IncludeBeneficiaries   bool   `json:"include_beneficiaries"`
	ValidateIDNumbers      bool   `json:"validate_id_numbers"`
	ValidateBankingDetails bool   `json:"validate_banking_details"`
	ExcludeInvalid         bool   `json:"exclude_invalid"`
	ReferenceNumber        string `json:"reference_number"`
	Notes                  string `json:"notes"`
	Category               string `json:"category"` // optional scheme category filter
	GeneratePerScheme      bool   `json:"generate_per_scheme"`
	// ExistingGeneratedID is set only by the regenerate service to replay a
	// prior request against the same generated_id. Never deserialized from
	// client input (tag is "-").
	ExistingGeneratedID string `json:"-"`
}

type BordereauxBatchSubmitRequest struct {
	BordereauxIDs  []int     `json:"bordereaux_ids"`
	DeliveryMethod string    `json:"delivery_method"`
	Message        string    `json:"message"`
	SubmissionDate time.Time `json:"submission_date"`
}

type BordereauxReportMeta struct {
	BordereauxID   int                  `json:"bordereaux_id"`
	GeneratedID    string               `json:"generated_id"`
	FilePath       string               `json:"file_path"`
	FileName       string               `json:"file_name"`
	FileExtension  string               `json:"file_extension"`
	DownloadURL    string               `json:"download_url"`
	Format         string               `json:"format"`
	Records        int                  `json:"records"`
	FileSize       int64                `json:"file_size"`
	ProcessingTime string               `json:"processing_time"`
	SchemeName     string               `json:"scheme_name"`
	InsurerName    string               `json:"insurer_name"`
	Status         string               `json:"status"`
	Progress       int                  `json:"progress"`
	SubmissionDate *time.Time           `json:"submission_date"`
	LastUpdated    time.Time            `json:"last_updated"`
	SLAStatus      string               `json:"sla_status"`
	PeriodStart    time.Time            `json:"period_start"`
	PeriodEnd      time.Time            `json:"period_end"`
	Timeline       []BordereauxTimeline `json:"timeline"`
}

type BordereauxTimeline struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

// GenerateBordereaux is the generic entrypoint that dispatches to the type-specific generators
func GenerateBordereaux(ctx context.Context, req GenerateBordereauxRequest) (any, error) {
	startTime := time.Now()
	// Normalize type
	t := strings.ToLower(strings.TrimSpace(req.Type))

	sendBordereauxProgress(ctx, BordereauxProgressEvent{
		Type:     t,
		Phase:    "start",
		Progress: 5,
		Message:  fmt.Sprintf("Preparing %s bordereaux generation", t),
	})

	if req.GeneratePerScheme && len(req.SchemeIDs) > 0 {
		var results []BordereauxReportMeta
		totalSchemes := len(req.SchemeIDs)
		for idx, id := range req.SchemeIDs {
			singleReq := req
			singleReq.SchemeIDs = []int{id}
			singleReq.GeneratePerScheme = false // prevent recursion if any

			var res BordereauxReportMeta
			var err error

			// Per-scheme progress: 5% → 85% linearly across the scheme loop,
			// leaving 85..100% for zipping and finalisation.
			perSchemeProgress := 5 + int(float64(idx)/float64(totalSchemes)*80)
			sendBordereauxProgress(ctx, BordereauxProgressEvent{
				Type:     t,
				Phase:    "generating",
				Progress: perSchemeProgress,
				Message:  fmt.Sprintf("Generating for scheme %d of %d", idx+1, totalSchemes),
				Scheme:   fmt.Sprintf("%d", id),
			})

			switch t {
			case "", "member":
				singleReq.Type = "member"
				res, err = GenerateMemberBordereaux(ctx, singleReq)
			// Premium / claim bordereaux generation temporarily disabled — only
			// member is a valid use case right now. Restore these cases when the
			// corresponding generators are re-enabled.
			// case "premium", "premiums":
			// 	singleReq.Type = "premium"
			// 	res, err = GeneratePremiumBordereaux(ctx, singleReq)
			// case "claim", "claims":
			// 	singleReq.Type = "claim"
			// 	res, err = GenerateClaimBordereaux(ctx, singleReq)
			default:
				sendBordereauxProgress(ctx, BordereauxProgressEvent{
					Type:     t,
					Phase:    "failed",
					Progress: 100,
					Message:  fmt.Sprintf("Unsupported bordereaux type: %s", req.Type),
				})
				return nil, fmt.Errorf("unsupported bordereaux type: %s", req.Type)
			}

			if err != nil {
				// Should we continue or fail fast? Failing fast for now.
				sendBordereauxProgress(ctx, BordereauxProgressEvent{
					Type:     t,
					Phase:    "failed",
					Progress: 100,
					Message:  fmt.Sprintf("Failed for scheme %d: %v", id, err),
				})
				return nil, fmt.Errorf("failed for scheme %d: %w", id, err)
			}
			results = append(results, res)
		}
		if len(results) == 1 {
			sendBordereauxProgress(ctx, BordereauxProgressEvent{
				Type: t, Phase: "completed", Progress: 100,
				Message: "Generation complete",
			})
			return results[0], nil
		}

		sendBordereauxProgress(ctx, BordereauxProgressEvent{
			Type: t, Phase: "zipping", Progress: 90,
			Message: fmt.Sprintf("Zipping %d generated files", len(results)),
		})

		// Zip results if multiple
		zipFileName := fmt.Sprintf("bordereaux_%s_%s.zip", t, newFileToken())
		reportDir := filepath.Join("data", "reports")
		zipPath := filepath.Join(reportDir, zipFileName)

		if err := createZipArchive(zipPath, results); err != nil {
			sendBordereauxProgress(ctx, BordereauxProgressEvent{
				Type: t, Phase: "failed", Progress: 100,
				Message: fmt.Sprintf("Zip failed: %v", err),
			})
			return nil, fmt.Errorf("failed to create zip archive: %w", err)
		}

		var fileSize int64
		if info, err := os.Stat(zipPath); err == nil {
			fileSize = info.Size()
		}

		generatedID, _ := generateUniqueBordereauxID(DB.WithContext(ctx))
		userEmail, _ := ctx.Value(appLog.UserEmailKey).(string)
		schemeNames := resolveSchemeNames(req.SchemeIDs)
		insurerName := ""
		if len(results) > 0 {
			insurerName = results[0].InsurerName
		}

		zipMeta := BordereauxReportMeta{
			BordereauxID:   req.TemplateID,
			GeneratedID:    generatedID,
			FilePath:       zipPath,
			FileName:       zipFileName,
			FileExtension:  ".zip",
			DownloadURL:    fmt.Sprintf("/group-pricing/bordereaux/download/%s", zipFileName),
			Format:         "zip",
			Records:        len(results),
			FileSize:       fileSize,
			ProcessingTime: time.Since(startTime).String(),
			SchemeName:     schemeNames,
			InsurerName:    insurerName,
			Status:         "generated",
			Progress:       33,
			LastUpdated:    time.Now(),
			SLAStatus:      "on_time",
			PeriodStart:    results[0].PeriodStart,
			PeriodEnd:      results[0].PeriodEnd,
			Timeline: []BordereauxTimeline{
				{
					Date:        time.Now(),
					Type:        "generated",
					Title:       "Generated",
					Description: "Bordereaux generated successfully",
				},
			},
		}

		now := time.Now()
		_ = DB.WithContext(ctx).Create(&models.GeneratedBordereaux{
			GeneratedID:    generatedID,
			TemplateID:     req.TemplateID,
			Type:           t,
			FileName:       zipFileName,
			FilePath:       zipPath,
			FileSize:       fileSize,
			Records:        len(results),
			ProcessingTime: zipMeta.ProcessingTime,
			SchemeName:     schemeNames,
			InsurerName:    insurerName,
			Status:         "generated",
			Progress:       33,
			LastUpdated:    time.Now(),
			SLAStatus:      "on_time",
			PeriodStart:    zipMeta.PeriodStart,
			PeriodEnd:      zipMeta.PeriodEnd,
			SubmissionDate: &now,
			CreatedBy:      userEmail,
			CreatedAt:      time.Now(),
			Timeline: []models.BordereauxTimeline{
				{
					Date:        time.Now(),
					Type:        "generated",
					Title:       "Generated",
					Description: "Bordereaux generated successfully",
				},
			},
		})

		sendBordereauxProgress(ctx, BordereauxProgressEvent{
			Type: t, Phase: "completed", Progress: 100,
			Message: fmt.Sprintf("Generation complete (%d schemes)", len(results)),
		})
		return zipMeta, nil
	}

	if !req.GeneratePerScheme && len(req.SchemeIDs) > 0 {
		sendBordereauxProgress(ctx, BordereauxProgressEvent{
			Type: t, Phase: "generating", Progress: 30,
			Message: "Generating combined bordereaux",
		})
		var out any
		var err error
		switch t {
		case "", "member":
			req.Type = "member"
			out, err = GenerateMemberBordereaux(ctx, req)
		// Premium / claim bordereaux generation temporarily disabled — only
		// member is a valid use case right now.
		// case "premium", "premiums":
		// 	req.Type = "premium"
		// 	out, err = GeneratePremiumBordereaux(ctx, req)
		// case "claim", "claims":
		// 	req.Type = "claim"
		// 	out, err = GenerateClaimBordereaux(ctx, req)
		default:
			sendBordereauxProgress(ctx, BordereauxProgressEvent{
				Type: t, Phase: "failed", Progress: 100,
				Message: fmt.Sprintf("Unsupported bordereaux type: %s", req.Type),
			})
			return BordereauxReportMeta{}, fmt.Errorf("unsupported bordereaux type: %s", req.Type)
		}
		if err != nil {
			sendBordereauxProgress(ctx, BordereauxProgressEvent{
				Type: t, Phase: "failed", Progress: 100,
				Message: err.Error(),
			})
			return out, err
		}
		sendBordereauxProgress(ctx, BordereauxProgressEvent{
			Type: t, Phase: "completed", Progress: 100,
			Message: "Generation complete",
		})
		return out, nil
	}
	sendBordereauxProgress(ctx, BordereauxProgressEvent{
		Type: t, Phase: "failed", Progress: 100,
		Message: fmt.Sprintf("Unsupported bordereaux type: %s", req.Type),
	})
	return nil, fmt.Errorf("unsupported bordereaux type: %s", req.Type)

}

func createZipArchive(zipPath string, results []BordereauxReportMeta) error {
	newZipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, res := range results {
		fileToZip, err := os.Open(res.FilePath)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = res.FileName
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, fileToZip); err != nil {
			return err
		}
		fileToZip.Close() // close immediately to avoid too many open files if list is long
	}
	return nil
}

// ensureNoConflictingMemberBordereaux rejects creation when any of the
// requested schemes already has an active (non-cancelled, non-failed) member
// bordereaux covering the same reporting month. Regenerate callers bypass this
// check because they reuse the existing row.
func ensureNoConflictingMemberBordereaux(ctx context.Context, req GenerateBordereauxRequest, start time.Time) error {
	if len(req.SchemeIDs) == 0 || start.IsZero() {
		return nil
	}

	monthStart := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	var existing []models.GeneratedBordereaux
	err := DB.WithContext(ctx).
		Where("type = ? AND status NOT IN ? AND period_start >= ? AND period_start < ?",
			"member",
			[]string{StatusGeneratedCancelled, StatusGeneratedFailed},
			monthStart, monthEnd).
		Find(&existing).Error
	if err != nil {
		// Don't block generation on a transient query failure; surface it instead.
		return fmt.Errorf("failed to check for existing bordereaux: %w", err)
	}
	if len(existing) == 0 {
		return nil
	}

	requested := make(map[int]struct{}, len(req.SchemeIDs))
	for _, id := range req.SchemeIDs {
		requested[id] = struct{}{}
	}

	for _, row := range existing {
		// Prefer the exact scheme set from the original request payload; fall
		// back to the human-readable scheme_name string for pre-payload rows.
		var existingReq GenerateBordereauxRequest
		if len(row.RequestPayload) > 0 {
			_ = json.Unmarshal(row.RequestPayload, &existingReq)
		}
		for _, sid := range existingReq.SchemeIDs {
			if _, clash := requested[sid]; clash {
				return fmt.Errorf(
					"scheme %q already has a member bordereaux for %s %d (existing %s, status %s) — cancel or regenerate it instead",
					resolveSchemeNames([]int{sid}), start.Month(), start.Year(), row.GeneratedID, row.Status,
				)
			}
		}
		// Legacy row without payload: best-effort name-level check.
		if len(existingReq.SchemeIDs) == 0 && row.SchemeName != "" {
			wanted := resolveSchemeNames(req.SchemeIDs)
			if wanted != "" && wanted == row.SchemeName {
				return fmt.Errorf(
					"%s already has a member bordereaux for %s %d (existing %s, status %s) — cancel or regenerate it instead",
					row.SchemeName, start.Month(), start.Year(), row.GeneratedID, row.Status,
				)
			}
		}
	}
	return nil
}

// resolvePeriod returns the start and end dates for the request
func resolvePeriod(req GenerateBordereauxRequest) (time.Time, time.Time, error) {
	switch strings.ToLower(req.PeriodType) {
	case "monthly":
		if req.Month == 0 || req.Year == 0 {
			return time.Time{}, time.Time{}, errors.New("month and year are required for monthly period")
		}
		start := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local)
		end := start.AddDate(0, 1, -1)
		return start, end, nil
	case "custom":
		if req.StartDate == "" || req.EndDate == "" {
			return time.Time{}, time.Time{}, errors.New("start_date and end_date required for custom period")
		}
		s, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		e, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return s, e, nil
	default:
		return time.Time{}, time.Time{}, errors.New("unsupported period_type")
	}
}

func generateUniqueBordereauxID(tx *gorm.DB) (string, error) {
	const maxAttempts = 5
	datePart := time.Now().Format("20060102")
	for i := 0; i < maxAttempts; i++ {
		suffix := fmt.Sprintf("%06d", (time.Now().UnixNano()+int64(i))%1000000)
		candidate := fmt.Sprintf("BRD-%s-%s", datePart, suffix)
		var count int64
		if err := tx.Model(&models.GeneratedBordereaux{}).Where("generated_id = ?", candidate).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
		time.Sleep(1 * time.Millisecond)
	}
	return "", errors.New("could not generate unique bordereaux ID after several attempts")
}

// GenerateMemberBordereaux orchestrates generation for type=="member"
func GenerateMemberBordereaux(ctx context.Context, req GenerateBordereauxRequest) (BordereauxReportMeta, error) {
	startTime := time.Now()

	if strings.ToLower(req.Type) != "member" {
		return BordereauxReportMeta{}, errors.New("unsupported bordereaux type")
	}

	// Resolve period
	start, end, err := resolvePeriod(req)
	if err != nil {
		return BordereauxReportMeta{}, err
	}

	// Load template
	tpl, err := GetBordereauxTemplateByID(req.TemplateID)
	if err != nil {
		return BordereauxReportMeta{}, err
	}

	if tpl.Type != "member" {
		return BordereauxReportMeta{}, fmt.Errorf("template type %s not supported for member generation", tpl.Type)
	}

	// Uniqueness: a scheme may only have one active member bordereaux per
	// period. Regenerate reuses the existing row, so skip this check then.
	if req.ExistingGeneratedID == "" {
		if err := ensureNoConflictingMemberBordereaux(ctx, req, start); err != nil {
			return BordereauxReportMeta{}, err
		}
	}

	// Fire-and-forget: increment usage asynchronously
	go func() { _ = IncrementTemplateUsage(tpl.ID) }()

	// Resolve latest quote IDs per scheme up to period end
	quoteIDs, err := resolveLatestQuotesForSchemes(req.SchemeIDs, end)
	if err != nil {
		return BordereauxReportMeta{}, err
	}
	if len(quoteIDs) == 0 {
		return BordereauxReportMeta{}, errors.New("no quotes found for selected schemes/period")
	}

	// -------------------------------------------------------------------------
	// Parallel fetch: in-effect quote IDs + (unblocked until quoteIDs resolved)
	// -------------------------------------------------------------------------
	var (
		inEffectQuoteIDs []int
		inEffectErr      error
	)

	inEffectQuoteIDs, inEffectErr = fetchInEffectQuoteIDs(ctx, quoteIDs)
	if inEffectErr != nil {
		return BordereauxReportMeta{}, inEffectErr
	}
	if len(inEffectQuoteIDs) == 0 {
		return BordereauxReportMeta{}, errors.New("no in-effect quotes found for generation period")
	}

	// -------------------------------------------------------------------------
	// Parallel fetch: inforce rows + rating summaries
	// -------------------------------------------------------------------------
	type fetchResult[T any] struct {
		data []T
		err  error
	}

	rowsCh := make(chan fetchResult[models.GPricingMemberDataInForce], 1)
	ratingCh := make(chan fetchResult[models.MemberRatingResultSummary], 1)

	go func() {
		var rows []models.GPricingMemberDataInForce
		err := DB.WithContext(ctx).
			Where("quote_id IN ?", inEffectQuoteIDs).
			Find(&rows).Error
		rowsCh <- fetchResult[models.GPricingMemberDataInForce]{data: rows, err: err}
	}()

	go func() {
		var summaries []models.MemberRatingResultSummary
		err := DB.WithContext(ctx).
			Where("quote_id IN ?", inEffectQuoteIDs).
			Find(&summaries).Error
		ratingCh <- fetchResult[models.MemberRatingResultSummary]{data: summaries, err: err}
	}()

	rowsResult := <-rowsCh
	if rowsResult.err != nil {
		return BordereauxReportMeta{}, rowsResult.err
	}
	if len(rowsResult.data) == 0 {
		return BordereauxReportMeta{}, errors.New("no inforce data found for resolved quotes")
	}

	ratingResult := <-ratingCh
	if ratingResult.err != nil {
		return BordereauxReportMeta{}, ratingResult.err
	}
	if len(ratingResult.data) == 0 {
		return BordereauxReportMeta{}, errors.New("no rating summaries found for in-effect quotes")
	}

	rows := rowsResult.data
	ratingSummaries := ratingResult.data

	// -------------------------------------------------------------------------
	// Build lookup maps
	// -------------------------------------------------------------------------

	// rate lookup: (QuoteID, Category) -> MemberRatingResultSummary
	type rateKey struct {
		QuoteID  int
		Category string
	}
	rates := make(map[rateKey]models.MemberRatingResultSummary, len(ratingSummaries))
	for _, rs := range ratingSummaries {
		rates[rateKey{QuoteID: rs.QuoteId, Category: rs.Category}] = rs
	}

	// -------------------------------------------------------------------------
	// Batch-fetch all distinct scheme categories needed (eliminates N+1)
	// -------------------------------------------------------------------------
	type categoryKey struct {
		QuoteID  int
		Category string
	}

	// Collect unique (quoteID-> []category) pairs
	quoteIDCategories := make(map[int][]string, len(rows))
	seen := make(map[categoryKey]struct{}, len(rows))
	for _, r := range rows {
		k := categoryKey{QuoteID: r.QuoteId, Category: r.SchemeCategory}
		if _, exists := seen[k]; !exists {
			seen[k] = struct{}{}
			quoteIDCategories[r.QuoteId] = append(quoteIDCategories[r.QuoteId], r.SchemeCategory)
		}
	}

	categoryCache, err := batchFetchSchemeCategories(ctx, quoteIDCategories)
	if err != nil {
		return BordereauxReportMeta{}, fmt.Errorf("failed to fetch scheme categories: %w", err)
	}

	// -------------------------------------------------------------------------
	// Batch-fetch all distinct scheme names needed (eliminates N+1)
	// -------------------------------------------------------------------------
	schemeIDSet := make(map[int]struct{}, len(rows))
	for _, r := range rows {
		if r.SchemeId > 0 {
			schemeIDSet[r.SchemeId] = struct{}{}
		}
	}
	schemeNames, err := batchFetchSchemeNames(ctx, schemeIDSet)
	if err != nil {
		return BordereauxReportMeta{}, fmt.Errorf("failed to fetch scheme names: %w", err)
	}

	// -------------------------------------------------------------------------
	// Build Bordereaux records
	// -------------------------------------------------------------------------
	periodStr := fmt.Sprintf("%s %d", time.Month(req.Month), req.Year)
	monthName := monthIntToName(req.Month)

	out := make([]models.Bordereaux, 0, len(rows))
	memberData := make([]models.MemberBordereauxData, 0, len(rows))
	var generatedID string
	if req.ExistingGeneratedID != "" {
		generatedID = req.ExistingGeneratedID
	} else {
		generatedID, _ = generateUniqueBordereauxID(DB.WithContext(ctx))
	}

	for _, r := range rows {
		rk := rateKey{QuoteID: r.QuoteId, Category: r.SchemeCategory}
		rate, ok := rates[rk]
		if !ok {
			// No rate found for this (quote, category) — skip
			continue
		}

		categoryMap, ok := categoryCache[r.QuoteId]
		if !ok {
			continue
		}
		category, ok := categoryMap[r.SchemeCategory]
		if !ok {
			continue
		}

		glaCappedSA := math.Max(r.AnnualSalary*r.Benefits.GlaMultiple, rate.FreeCoverLimit)
		glaRiskPremium := glaCappedSA * rate.ExpGlaRiskRatePer1000SA / 1000
		glaAnnualPremium := utils.FloatPrecision(
			models.ComputeOfficePremium(glaRiskPremium, &rate),
			AccountingPrecision,
		)

		b := models.Bordereaux{
			SchemeId:       r.SchemeId,
			MemberName:     r.MemberName,
			EmployeeNumber: r.EmployeeNumber,
			QuoteId:        r.QuoteId,
			MemberIdNumber: r.MemberIdNumber,
			Category:       r.SchemeCategory,
			Month:          time.Month(req.Month),
			Year:           req.Year,
			EntryDate:      r.EntryDate,
			ExitDate: func() time.Time {
				if r.ExitDate == nil {
					return time.Time{}
				}
				return *r.ExitDate
			}(),
			Period: periodStr,
			Gender: r.Gender,
			DateOfBirth: func() *time.Time {
				if r.DateOfBirth.IsZero() {
					return nil
				}
				t := r.DateOfBirth
				return &t
			}(),
			AnnualSalary:                r.AnnualSalary,
			GlaMultiple:                 r.Benefits.GlaMultiple,
			PtdMultiple:                 r.Benefits.PtdMultiple,
			CiMultiple:                  r.Benefits.CiMultiple,
			SglaMultiple:                r.Benefits.SglaMultiple,
			PhiReplacementMultiple:      r.Benefits.PhiMultiple,
			TtdReplacementMultiple:      r.Benefits.TtdMultiple,
			MainMemberFuneralSumAssured: category.FamilyFuneralMainMemberFuneralSumAssured,
			SpouseFuneralSumAssured:     category.FamilyFuneralSpouseFuneralSumAssured,
			ChildFuneralSumAssured:      category.FamilyFuneralChildrenFuneralSumAssured,
			ParentFuneralSumAssured:     category.FamilyFuneralParentFuneralSumAssured,
			DependantFuneralSumAssured:  category.FamilyFuneralAdultDependantSumAssured,
			GlaAnnualPremium:            glaAnnualPremium,
			GlaCoveredSumAssured:        math.Min(r.AnnualSalary*r.Benefits.GlaMultiple, category.FreeCoverLimit),
			GlaRetainedSumAssured:       math.Min(r.AnnualSalary*r.Benefits.GlaMultiple, category.FreeCoverLimit),
			GlaCededSumAssured:          math.Min(r.AnnualSalary*r.Benefits.GlaMultiple, category.FreeCoverLimit),
			PtdCoveredSumAssured:        math.Min(r.AnnualSalary*r.Benefits.PtdMultiple, category.FreeCoverLimit),
			PtdRetainedSumAssured:       math.Min(r.AnnualSalary*r.Benefits.PtdMultiple, category.FreeCoverLimit),
			PtdCededSumAssured:          math.Min(r.AnnualSalary*r.Benefits.PtdMultiple, category.FreeCoverLimit),
			CiCoveredSumAssured:         math.Min(r.AnnualSalary*r.Benefits.CiMultiple, category.FreeCoverLimit),
			CiRetainedSumAssured:        math.Min(r.AnnualSalary*r.Benefits.CiMultiple, category.FreeCoverLimit),
			CiCededSumAssured:           math.Min(r.AnnualSalary*r.Benefits.CiMultiple, category.FreeCoverLimit),
			SglaCoveredSumAssured:       math.Min(r.AnnualSalary*r.Benefits.SglaMultiple, category.FreeCoverLimit),
			SglaRetainedSumAssured:      math.Min(r.AnnualSalary*r.Benefits.SglaMultiple, category.FreeCoverLimit),
			SglaCededSumAssured:         math.Min(r.AnnualSalary*r.Benefits.SglaMultiple, category.FreeCoverLimit),
			PhiMonthlyBenefit:           math.Min(r.AnnualSalary*r.Benefits.PhiMultiple, category.FreeCoverLimit),
			PhiRetainedMonthlyBenefit:   math.Min(r.AnnualSalary*r.Benefits.PhiMultiple, category.FreeCoverLimit),
			PhiCededMonthlyBenefit:      math.Min(r.AnnualSalary*r.Benefits.PhiMultiple, category.FreeCoverLimit),
			TtdMonthlyBenefit:           math.Min(r.AnnualSalary*r.Benefits.TtdMultiple, category.FreeCoverLimit),
			TtdRetainedMonthlyBenefit:   math.Min(r.AnnualSalary*r.Benefits.TtdMultiple, category.FreeCoverLimit),
			TtdCededMonthlyBenefit:      math.Min(r.AnnualSalary*r.Benefits.TtdMultiple, category.FreeCoverLimit),
		}

		md := models.MemberBordereauxData{
			BordereauxID:   generatedID, // set after ID generation below
			SchemeName:     schemeNames[r.SchemeId],
			MemberName:     b.MemberName,
			EmployeeNumber: b.EmployeeNumber,
			MemberIdNumber: b.MemberIdNumber,
			Category:       b.Category,
			Month:          monthName,
			Year:           req.Year,
			Period:         fmt.Sprintf("%s %d", monthName, req.Year),
			Gender:         b.Gender,
			DateOfBirth: func() time.Time {
				if b.DateOfBirth == nil {
					return time.Time{}
				}
				return *b.DateOfBirth
			}(),
			AnnualSalary:           b.AnnualSalary,
			GlaMultiple:            b.GlaMultiple,
			GlaCoveredSumAssured:   b.GlaCoveredSumAssured,
			GlaRetainedSumAssured:  b.GlaRetainedSumAssured,
			GlaCededSumAssured:     b.GlaCededSumAssured,
			PtdMultiple:            b.PtdMultiple,
			PtdCoveredSumAssured:   b.PtdCoveredSumAssured,
			PtdRetainedSumAssured:  b.PtdRetainedSumAssured,
			PtdCededSumAssured:     b.PtdCededSumAssured,
			CiMultiple:             b.CiMultiple,
			CiCoveredSumAssured:    b.CiCoveredSumAssured,
			CiRetainedSumAssured:   b.CiRetainedSumAssured,
			CiCededSumAssured:      b.CiCededSumAssured,
			SglaMultiple:           b.SglaMultiple,
			SglaCoveredSumAssured:  b.SglaCoveredSumAssured,
			SglaRetainedSumAssured: b.SglaRetainedSumAssured,
			SglaCededSumAssured:    b.SglaCededSumAssured,
			PhiMultiple:            b.PhiReplacementMultiple,
			PhiCoveredIncome:       b.PhiMonthlyBenefit,
			PhiRetainedIncome:      b.PhiRetainedMonthlyBenefit,
			PhiCededIncome:         b.PhiCededMonthlyBenefit,
			TtdMultiple:            b.TtdReplacementMultiple,
			TtdCoveredIncome:       b.TtdMonthlyBenefit,
			TtdRetainedIncome:      b.TtdRetainedMonthlyBenefit,
			TtdCededIncome:         b.TtdCededMonthlyBenefit,
			MmFuneralSumAssured:    b.MainMemberFuneralSumAssured,
			SpFuneralSumAssured:    b.SpouseFuneralSumAssured,
			ChFuneralSumAssured:    b.ChildFuneralSumAssured,
			ParFuneralSumAssured:   b.ParentFuneralSumAssured,
			DepFuneralSumAssured:   b.DependantFuneralSumAssured,
		}

		out = append(out, b)
		memberData = append(memberData, md)
	}

	if len(out) == 0 {
		return BordereauxReportMeta{}, errors.New("no bordereaux records could be generated after applying filters")
	}

	// -------------------------------------------------------------------------
	// Determine output format and write file
	// -------------------------------------------------------------------------
	format := strings.ToLower(req.OutputFormat)
	if format == "" {
		format = strings.ToLower(tpl.Format)
	}
	if format == "" {
		format = "excel"
	}

	reportDir := filepath.Join("data", "reports")
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		return BordereauxReportMeta{}, fmt.Errorf("failed to create report directory: %w", err)
	}

	ext := ".xlsx"
	if format == "csv" {
		ext = ".csv"
	}
	fileName := fmt.Sprintf("bordereaux_member_%d%02d_%s%s", start.Year(), int(start.Month()), newFileToken(), ext)
	outPath := filepath.Join(reportDir, fileName)
	downloadURL := fmt.Sprintf("/group-pricing/bordereaux/download/%s", fileName)

	switch format {
	case "excel", "xlsx":
		if err := writeExcel(outPath, tpl, out, req); err != nil {
			return BordereauxReportMeta{}, err
		}
	case "csv":
		if err := writeCSV(outPath, tpl, out, req); err != nil {
			return BordereauxReportMeta{}, err
		}
	default:
		return BordereauxReportMeta{}, errors.New("unsupported output format")
	}

	var fileSize int64
	if info, err := os.Stat(outPath); err == nil {
		fileSize = info.Size()
	}

	if len(memberData) > 0 {
		appLog.WithField("count", len(memberData)).Info("Persisting member bordereaux data")
		if err := DB.WithContext(ctx).CreateInBatches(memberData, 500).Error; err != nil {
			appLog.WithField("error", err.Error()).Error("Failed to persist member bordereaux data")
			return BordereauxReportMeta{}, fmt.Errorf("failed to persist member bordereaux data: %w", err)
		}
	}

	userEmail, _ := ctx.Value(appLog.UserEmailKey).(string)
	schemeName := resolveSchemeNames(req.SchemeIDs)
	insurerName := resolveInsurerName(tpl.InsurerID)

	if req.ExistingGeneratedID != "" {
		// Regenerate path: update the existing row in place so generated_id,
		// timeline, audit and created_at/created_by survive. Reviewer / approver
		// fields are cleared because the file has been refreshed.
		var existing models.GeneratedBordereaux
		_ = DB.WithContext(ctx).Where("generated_id = ?", req.ExistingGeneratedID).First(&existing).Error
		oldFilePath := existing.FilePath

		updates := map[string]any{
			"template_id":     req.TemplateID,
			"file_name":       fileName,
			"file_path":       outPath,
			"file_size":       fileSize,
			"records":         len(out),
			"processing_time": time.Since(startTime).String(),
			"scheme_name":     schemeName,
			"insurer_name":    insurerName,
			"status":          StatusGeneratedGenerated,
			"progress":        33,
			"last_updated":    time.Now(),
			"sla_status":      "on_time",
			"period_start":    start,
			"period_end":      end,
			"submission_date": &start,
			"reviewed_by":     "",
			"reviewed_at":     gorm.Expr("NULL"),
			"approved_by":     "",
			"approved_at":     gorm.Expr("NULL"),
			"return_reason":   "",
		}
		_ = DB.WithContext(ctx).Model(&models.GeneratedBordereaux{}).
			Where("generated_id = ?", req.ExistingGeneratedID).
			Updates(updates).Error

		_ = DB.WithContext(ctx).Create(&models.BordereauxTimeline{
			GeneratedBordereauxID: req.ExistingGeneratedID,
			Date:                  time.Now(),
			Type:                  "regenerated",
			Title:                 "Regenerated",
			Description:           "Bordereaux regenerated successfully",
		})

		if oldFilePath != "" && oldFilePath != outPath {
			if err := os.Remove(oldFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
				appLog.WithField("error", err.Error()).Warn("Failed to delete old bordereaux file")
			}
		}
	} else {
		// Create path: persist the original request so a future regenerate can
		// replay it without asking the user to re-enter parameters.
		payload, _ := json.Marshal(req)
		_ = DB.WithContext(ctx).Create(&models.GeneratedBordereaux{
			GeneratedID:    generatedID,
			TemplateID:     req.TemplateID,
			Type:           "member",
			FileName:       fileName,
			FilePath:       outPath,
			FileSize:       fileSize,
			Records:        len(out),
			ProcessingTime: time.Since(startTime).String(),
			SchemeName:     schemeName,
			InsurerName:    insurerName,
			Status:         StatusGeneratedGenerated,
			Progress:       33,
			LastUpdated:    time.Now(),
			SLAStatus:      "on_time",
			PeriodStart:    start,
			PeriodEnd:      end,
			CreatedBy:      userEmail,
			CreatedAt:      time.Now(),
			SubmissionDate: &start,
			RequestPayload: models.JSON(payload),
			Timeline: []models.BordereauxTimeline{
				{
					Date:        time.Now(),
					Type:        "generated",
					Title:       "Generated",
					Description: "Bordereaux generated successfully",
				},
			},
		})
	}

	return BordereauxReportMeta{
		BordereauxID:   req.TemplateID,
		GeneratedID:    generatedID,
		FilePath:       outPath,
		FileName:       fileName,
		FileExtension:  filepath.Ext(fileName),
		DownloadURL:    downloadURL,
		Format:         format,
		Records:        len(out),
		FileSize:       fileSize,
		ProcessingTime: time.Since(startTime).String(),
		SchemeName:     schemeName,
		InsurerName:    insurerName,
		Status:         "generated",
		Progress:       33,
		LastUpdated:    time.Now(),
		SLAStatus:      "on_time",
		PeriodStart:    start,
		PeriodEnd:      end,
		Timeline: []BordereauxTimeline{
			{
				Date:        time.Now(),
				Type:        "generated",
				Title:       "Generated",
				Description: "Bordereaux generated successfully",
			},
		},
	}, nil
}

func writeExcel(path string, tpl models.BordereauxTemplate, rows []models.Bordereaux, req GenerateBordereauxRequest) error {
	f := excelize.NewFile()
	sheet := "Bordereaux"
	f.SetSheetName("Sheet1", sheet)
	// headers
	for i, m := range tpl.FieldMappings {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, m.TargetField)
	}
	// rows
	for rIdx, b := range rows {
		values, _ := mapFields(tpl, b, req)
		for cIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
			_ = f.SetCellValue(sheet, cell, v)
		}
	}
	return f.SaveAs(path)
}

func writeCSV(path string, tpl models.BordereauxTemplate, rows []models.Bordereaux, req GenerateBordereauxRequest) error {
	// simple CSV writer
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// header
	header := make([]string, len(tpl.FieldMappings))
	for i, m := range tpl.FieldMappings {
		header[i] = m.TargetField
	}
	file.WriteString(strings.Join(header, ",") + "\n")
	// rows
	for _, b := range rows {
		values, _ := mapFields(tpl, b, req)
		strVals := make([]string, len(values))
		for i, v := range values {
			strVals[i] = fmt.Sprint(v)
		}
		file.WriteString(strings.Join(strVals, ",") + "\n")
	}
	return nil
}

func maskID(id string) string {
	id = strings.TrimSpace(id)
	if len(id) <= 4 {
		return "****"
	}
	return id[:len(id)-4] + "****"
}

// formatDate renders date values consistently as YYYY-MM-DD strings so Excel /
// CSV output is readable. Non-date values pass through unchanged.
func formatDate(v any) any {
	switch d := v.(type) {
	case time.Time:
		if d.IsZero() {
			return ""
		}
		return d.Format("2006-01-02")
	case *time.Time:
		if d == nil || d.IsZero() {
			return ""
		}
		return d.Format("2006-01-02")
	default:
		return v
	}
}

// mapFields evaluates template field mappings against a Bordereaux row
func mapFields(tpl models.BordereauxTemplate, b models.Bordereaux, req GenerateBordereauxRequest) ([]any, []string) {
	vals := make([]any, 0, len(tpl.FieldMappings))
	var errs []string
	for _, m := range tpl.FieldMappings {
		v, ok := resolveSourceField(m.SourceField, b)
		if !ok && m.Required {
			errs = append(errs, fmt.Sprintf("missing required field: %s", m.SourceField))
		}
		// mask member_id or member_id_number
		if strings.EqualFold(m.SourceField, "member.member_id") || strings.EqualFold(m.SourceField, "member_id_number") || strings.EqualFold(m.SourceField, "member.id_number") {
			if s, ok := v.(string); ok && s != "" {
				v = maskID(s)
			}
		}
		// run lightweight validations for specific fields
		if req.ValidateIDNumbers && strings.Contains(strings.ToLower(m.TargetField), "id") {
			s := fmt.Sprint(v)
			if s != "" && !ValidateSouthAfricanID(s) {
				errs = append(errs, "invalid SA ID")
			}
		}
		vals = append(vals, v)
	}
	return vals, errs
}

// resolveSourceField maps a source field name to a value on models.Bordereaux.
// Covers every field advertised by GetBordereauxFieldsByType("member").
func resolveSourceField(path string, b models.Bordereaux) (any, bool) {
	switch strings.ToLower(path) {
	case "id":
		return 0, true // placeholder id; extend as needed
	case "member.member_name", "member_name":
		return b.MemberName, true
	case "employee_number":
		return b.EmployeeNumber, true
	case "member.member_id", "member.id_number", "member_id_number":
		return b.MemberIdNumber, true
	case "member.annual_salary", "annual_salary":
		return b.AnnualSalary, true
	case "entry_date":
		return formatDate(b.EntryDate), true
	case "exit_date":
		return formatDate(b.ExitDate), true
	case "gla_salary_multiple":
		return b.GlaMultiple, true
	case "sgla_salary_multiple":
		return b.SglaMultiple, true
	case "ptd_salary_multiple":
		return b.PtdMultiple, true
	case "ci_salary_multiple":
		return b.CiMultiple, true
	case "gla.sum_assured", "gla_covered_sum_assured":
		return b.GlaCoveredSumAssured, true
	case "ptd_capped_sum_assured":
		return b.PtdCoveredSumAssured, true
	case "ci_capped_sum_assured":
		return b.CiCoveredSumAssured, true
	case "spouse_gla_capped_sum_assured":
		return b.SglaCoveredSumAssured, true
	case "ttd_capped_income":
		return b.TtdMonthlyBenefit, true
	case "phi_capped_income":
		return b.PhiMonthlyBenefit, true
	case "gla_ceded_sum_assured":
		return b.GlaCededSumAssured, true
	case "ptd_ceded_sum_assured":
		return b.PtdCededSumAssured, true
	case "ci_ceded_sum_assured":
		return b.CiCededSumAssured, true
	case "spouse_gla_ceded_sum_assured":
		return b.SglaCededSumAssured, true
	case "ttd_ceded_income":
		return b.TtdCededMonthlyBenefit, true
	case "phi_ceded_income":
		return b.PhiCededMonthlyBenefit, true
	case "member_funeral_sum_assured":
		return b.MainMemberFuneralSumAssured, true
	case "spouse_funeral_sum_assured":
		return b.SpouseFuneralSumAssured, true
	case "child_funeral_sum_assured":
		return b.ChildFuneralSumAssured, true
	case "dependant_funeral_sum_assured":
		return b.DependantFuneralSumAssured, true
	case "parent_funeral_sum_assured":
		return b.ParentFuneralSumAssured, true
	case "member_funeral_ceded_sum_assured":
		return b.MainMemberCededSumAssured, true
	case "spouse_funeral_ceded_sum_assured":
		return b.SpouseCededSumAssured, true
	case "child_funeral_ceded_sum_assured":
		return b.ChildCededSumAssured, true
	case "dependant_funeral_ceded_sum_assured":
		return b.DependantCededSumAssured, true
	case "parent_funeral_ceded_sum_assured":
		return b.ParentCededSumAssured, true
	case "gla.rate_per_1000", "loaded_gla_risk_rate":
		return b.LoadedGlaRiskRate, true
	default:
		return "", false
	}
}

// ValidateSouthAfricanID performs a Luhn-like checksum validation for a 13-digit ID
func ValidateSouthAfricanID(id string) bool {
	id = strings.TrimSpace(id)
	if len(id) != 13 {
		return false
	}
	for _, c := range id {
		if c < '0' || c > '9' {
			return false
		}
	}
	sum := 0
	alt := false
	for i := len(id) - 1; i >= 0; i-- {
		d := int(id[i] - '0')
		if alt {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		alt = !alt
	}
	return sum%10 == 0
}

// isValidPassportNumber checks that a passport number is alphanumeric and between 5 and 20 characters.
func isValidPassportNumber(id string) bool {
	id = strings.TrimSpace(id)
	if len(id) < 5 || len(id) > 20 {
		return false
	}
	for _, c := range id {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

// ClassifyMemberID detects whether an ID string is an RSA ID (13 all-digit) or a passport number
// and validates it accordingly. Returns (idType, isValid, failReason).
func ClassifyMemberID(id string) (string, bool, string) {
	trimmed := strings.TrimSpace(id)
	// 13-character all-digit string → RSA ID
	allDigits := len(trimmed) > 0
	for _, c := range trimmed {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if len(trimmed) == 13 && allDigits {
		if ValidateSouthAfricanID(trimmed) {
			return "rsa_id", true, ""
		}
		return "rsa_id", false, "SA ID number failed Luhn checksum"
	}
	// Otherwise treat as passport
	if isValidPassportNumber(trimmed) {
		return "passport", true, ""
	}
	return "passport", false, "Passport number format is invalid (expected 5–20 alphanumeric characters)"
}

// resolveLatestQuotesForSchemes returns the most recent quote id for each scheme effective on/before atDate
func resolveLatestQuotesForSchemes(schemeIDs []int, atDate time.Time) ([]int, error) {
	if len(schemeIDs) == 0 {
		return nil, nil
	}
	var quotes []models.GroupPricingQuote
	if err := DB.Where("scheme_id IN ? and scheme_quote_status = ? ", schemeIDs, models.StatusInEffect).Order("commencement_date desc, id desc").Find(&quotes).Error; err != nil {
		return nil, err
	}
	// pick the first per scheme
	selected := map[int]int{}
	for _, q := range quotes {
		if !q.CommencementDate.IsZero() && q.CommencementDate.After(atDate) {
			continue
		}
		if _, ok := selected[q.SchemeID]; !ok {
			selected[q.SchemeID] = q.ID
		}
	}
	out := make([]int, 0, len(selected))
	for _, id := range selected {
		out = append(out, id)
	}
	return out, nil
}

// ================= Premiums Bordereaux =================
// Premium bordereaux generation is temporarily disabled — only member is a
// valid use case right now. The entire block below is commented out so it can
// be restored in one go.
/*
// GeneratePremiumBordereaux creates a Premiums bordereaux using MemberPremiumSchedule data
func GeneratePremiumBordereaux(ctx context.Context, req GenerateBordereauxRequest) (BordereauxReportMeta, error) {
	startTime := time.Now()
	if strings.ToLower(req.Type) != "premium" {
		return BordereauxReportMeta{}, errors.New("unsupported bordereaux type")
	}

	start, end, err := resolvePeriod(req)
	if err != nil {
		return BordereauxReportMeta{}, err
	}

	// Load template and verify type
	tpl, err := GetBordereauxTemplateByID(req.TemplateID)
	if err != nil {
		return BordereauxReportMeta{}, err
	}

	// Increment usage
	_ = IncrementTemplateUsage(tpl.ID)

	if tpl.Type != "premium" && tpl.Type != "premiums" {
		return BordereauxReportMeta{}, fmt.Errorf("template type %s not supported for premium generation", tpl.Type)
	}

	// Resolve latest quotes for schemes
	quoteIDs, err := resolveLatestQuotesForSchemes(req.SchemeIDs, time.Now())
	if err != nil {
		return BordereauxReportMeta{}, err
	}
	if len(quoteIDs) == 0 {
		return BordereauxReportMeta{}, errors.New("no quotes found for selected schemes/period")
	}

	// Fetch rows from MemberPremiumSchedule
	//var rows []models.MemberPremiumSchedule
	//q := DB.WithContext(ctx).Where("quote_id IN ?", quoteIDs)
	//if req.Category != "" {
	//	q = q.Where("category = ?", req.Category)
	//}

	//Fetch rows from GPricingMemberMemberDataInforce
	var rows []models.GPricingMemberDataInForce
	q := DB.WithContext(ctx).Where("scheme_id IN ?", req.SchemeIDs)
	if req.Category != "" {
		q = q.Where("scheme_category = ?", req.Category)
	}

	//Fetch MemberRatingResultSummary for the quote in effect
	var mrss models.MemberRatingResultSummary
	r := DB.WithContext(ctx).Where("quote_id IN ?", quoteIDs)
	if err := r.Find(&mrss).Error; err != nil {
		return BordereauxReportMeta{}, err
	}

	if err := q.Find(&rows).Error; err != nil {
		return BordereauxReportMeta{}, err
	}

	// Optional filter: terminated
	//if !req.IncludeTerminated {
	//	filtered := make([]models.MemberPremiumSchedule, 0, len(rows))
	//	for _, r := range rows {
	//		if !r.ExitDate.IsZero() && r.ExitDate.Before(start) {
	//			continue
	//		}
	//		filtered = append(filtered, r)
	//	}
	//	rows = filtered
	//}

	if !req.IncludeTerminated {
		filtered := make([]models.GPricingMemberDataInForce, 0, len(rows))
		for _, r := range rows {
			if r.Status == "INACTIVE" {
				continue
			}
			filtered = append(filtered, r)
		}
		rows = filtered
	}

	// Write output
	format := strings.ToLower(req.OutputFormat)
	if format == "" {
		format = strings.ToLower(tpl.Format)
	}
	if format == "" {
		format = "excel"
	}
	reportDir := filepath.Join("data", "reports")
	_ = os.MkdirAll(reportDir, 0o755)
	generatedID, _ := generateUniqueBordereauxID(DB.WithContext(ctx))
	//fileName := fmt.Sprintf("bordereaux_premium_%d%02d_%d.xlsx", start.Year(), int(start.Month()), time.Now().Unix())
	fileName := fmt.Sprintf("%s.xlsx", generatedID)
	if format == "csv" {
		fileName = strings.TrimSuffix(fileName, ".xlsx") + ".csv"
	}
	outPath := filepath.Join(reportDir, fileName)
	downloadURL := fmt.Sprintf("/group-pricing/bordereaux/download/%s", fileName)

	switch format {
	case "excel", "xlsx":
		if err := writeExcelPremium(outPath, tpl, rows, mrss, req, generatedID); err != nil {
			return BordereauxReportMeta{}, err
		}
	case "csv":
		if err := writeCSVPremium(outPath, tpl, rows, mrss, req, generatedID); err != nil {
			return BordereauxReportMeta{}, err
		}
	default:
		return BordereauxReportMeta{}, errors.New("unsupported output format")
	}

	// Persist premium data fields
	var premiumData []models.PremiumBordereauxData
	for _, b := range rows {
		// Calculate age
		valuationDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1)
		age := valuationDate.Year() - b.DateOfBirth.Year()
		if valuationDate.YearDay() < b.DateOfBirth.YearDay() {
			age--
		}

		convertToMonth := func(month int) string {
			months := []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
			if month >= 1 && month <= 12 {
				return months[month]
			}
			return ""
		}

		data := models.PremiumBordereauxData{
			BordereauxID:                       generatedID,
			SchemeName:                         b.SchemeName,
			MemberName:                         b.MemberName,
			MemberIdNumber:                     b.MemberIdNumber,
			EmployeeNumber:                     b.EmployeeNumber,
			Category:                           b.SchemeCategory,
			Month:                              convertToMonth(req.Month),
			Year:                               req.Year,
			Period:                             fmt.Sprintf("%s %d", convertToMonth(req.Month), req.Year),
			Age:                                age,
			AnnualSalary:                       b.AnnualSalary,
			GlaAnnualPremium:                   getPremiumValue("gla_annual_premium", b, mrss),
			PtdAnnualPremium:                   getPremiumValue("ptd_annual_premium", b, mrss),
			CiAnnualPremium:                    getPremiumValue("ci_annual_premium", b, mrss),
			PhiAnnualPremium:                   getPremiumValue("phi_annual_premium", b, mrss),
			TotalAnnualFuneralPremium:          getPremiumValue("total_annual_funeral_premium", b, mrss),
			TotalAnnualPremiumExcludingFuneral: getPremiumValue("total_annual_premium_excl_funeral", b, mrss),
			TotalAnnualPremium:                 getPremiumValue("total_annual_premium", b, mrss),
		}
		premiumData = append(premiumData, data)
	}

	if len(premiumData) > 0 {
		appLog.WithField("count", len(premiumData)).Info("Persisting premium bordereaux data")
		if err := DB.WithContext(ctx).CreateInBatches(premiumData, 500).Error; err != nil {
			appLog.WithField("error", err.Error()).Error("Failed to persist premium bordereaux data")
			return BordereauxReportMeta{}, fmt.Errorf("failed to persist premium bordereaux data: %w", err)
		}
	}

	var fileSize int64
	if info, err := os.Stat(outPath); err == nil {
		fileSize = info.Size()
	}

	userEmail, _ := ctx.Value(appLog.UserEmailKey).(string)
	schemeName := resolveSchemeNames(req.SchemeIDs)
	insurerName := resolveInsurerName(tpl.InsurerID)
	now := time.Now()
	_ = DB.WithContext(ctx).Create(&models.GeneratedBordereaux{
		GeneratedID:    generatedID,
		TemplateID:     req.TemplateID,
		Type:           "premium",
		FileName:       fileName,
		FilePath:       outPath,
		FileSize:       fileSize,
		Records:        len(rows),
		ProcessingTime: time.Since(startTime).String(),
		SchemeName:     schemeName,
		InsurerName:    insurerName,
		Status:         "generated",
		Progress:       33,
		LastUpdated:    time.Now(),
		SLAStatus:      "on_time",
		SubmissionDate: &now,
		PeriodStart:    start,
		PeriodEnd:      end,
		CreatedBy:      userEmail,
		CreatedAt:      time.Now(),
		Timeline: []models.BordereauxTimeline{
			{
				Date:        time.Now(),
				Type:        "generated",
				Title:       "Generated",
				Description: "Bordereaux generated successfully",
			},
		},
	})

	return BordereauxReportMeta{
		BordereauxID:   req.TemplateID,
		GeneratedID:    generatedID,
		FilePath:       outPath,
		FileName:       fileName,
		FileExtension:  filepath.Ext(fileName),
		DownloadURL:    downloadURL,
		Format:         format,
		Records:        len(rows),
		FileSize:       fileSize,
		ProcessingTime: time.Since(startTime).String(),
		SchemeName:     schemeName,
		InsurerName:    insurerName,
		Status:         "generated",
		Progress:       33,
		LastUpdated:    time.Now(),
		SLAStatus:      "on_time",
		PeriodStart:    start,
		PeriodEnd:      end,
		Timeline: []BordereauxTimeline{
			{
				Date:        time.Now(),
				Type:        "generated",
				Title:       "Generated",
				Description: "Bordereaux generated successfully",
			},
		},
	}, nil
}

func writeExcelPremium(path string, tpl models.BordereauxTemplate, rows []models.GPricingMemberDataInForce, mrss models.MemberRatingResultSummary, req GenerateBordereauxRequest, bordereauxId string) error {
	f := excelize.NewFile()
	sheet := "Bordereaux"
	f.SetSheetName("Sheet1", sheet)

	// Add fixed columns that are not in the template
	_ = f.SetCellValue(sheet, "A1", "Bordereaux ID")
	_ = f.SetCellValue(sheet, "B1", "Scheme Name")
	_ = f.SetCellValue(sheet, "C1", "Valuation Month")
	_ = f.SetCellValue(sheet, "D1", "Valuation Year")
	_ = f.SetCellValue(sheet, "E1", "Period Type")

	for i, m := range tpl.FieldMappings {
		cell, _ := excelize.CoordinatesToCellName(i+6, 1)
		_ = f.SetCellValue(sheet, cell, m.TargetField)
	}
	for rIdx, b := range rows {
		values, _ := mapFieldsPremium(tpl, b, mrss, req, bordereauxId)
		for cIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
			_ = f.SetCellValue(sheet, cell, v)
		}
	}
	return f.SaveAs(path)
}

func writeCSVPremium(path string, tpl models.BordereauxTemplate, rows []models.GPricingMemberDataInForce, mrss models.MemberRatingResultSummary, req GenerateBordereauxRequest, bordereauxId string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	header := make([]string, 0, len(tpl.FieldMappings)+5)
	header = append(header, "Bordereaux ID", "Scheme Name", "Valuation Month", "Valuation Year", "Period Type")
	for _, m := range tpl.FieldMappings {
		header = append(header, m.TargetField)
	}
	file.WriteString(strings.Join(header, ",") + "\n")
	for _, b := range rows {
		values, _ := mapFieldsPremium(tpl, b, mrss, req, bordereauxId)
		strVals := make([]string, len(values))
		for i, v := range values {
			strVals[i] = fmt.Sprint(v)
		}
		file.WriteString(strings.Join(strVals, ",") + "\n")
	}
	return nil
}

func mapFieldsPremium(tpl models.BordereauxTemplate, b models.GPricingMemberDataInForce, mrss models.MemberRatingResultSummary, req GenerateBordereauxRequest, bordereauxId string) ([]any, []string) {
	vals := make([]any, 0, len(tpl.FieldMappings)+5)

	convertToMonth := func(month int) string {
		switch month {
		case 1:
			return "January"
		case 2:
			return "February"
		case 3:
			return "March"
		case 4:
			return "April"
		case 5:
			return "May"
		case 6:
			return "June"
		case 7:
			return "July"
		case 8:
			return "August"
		case 9:
			return "September"
		case 10:
			return "October"
		case 11:
			return "November"
		case 12:
			return "December"
		default:
			return ""
		}
	}

	vals = append(vals, bordereauxId)
	vals = append(vals, b.SchemeName)
	vals = append(vals, convertToMonth(req.Month))
	vals = append(vals, req.Year)
	vals = append(vals, req.PeriodType)

	var errs []string
	for _, m := range tpl.FieldMappings {
		v, ok := resolvePremiumSourceField(m.SourceField, b, mrss)
		if !ok && m.Required {
			errs = append(errs, fmt.Sprintf("missing required field: %s", m.SourceField))
		}

		// mask member_id or member_id_number
		if strings.EqualFold(m.SourceField, "member.member_id_number") || strings.EqualFold(m.SourceField, "member_id_number") {
			if s, ok := v.(string); ok && s != "" {
				v = maskID(s)
			}
		}

		// optional numeric validation for amounts
		// handle special cases for valuation_month and age_at_month if needed
		if m.SourceField == "valuation_month" {
			// we need to convert the number to a month, e.g. 12 -> December
			v = convertToMonth(req.Month)
		}

		if m.SourceField == "age_at_month" {
			// calculate age at the end of the month
			valuationDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1)
			age := valuationDate.Year() - b.DateOfBirth.Year()
			if valuationDate.YearDay() < b.DateOfBirth.YearDay() {
				age--
			}
			v = age
		}

		vals = append(vals, v)
	}
	return vals, errs
}

func getPremiumValue(field string, b models.GPricingMemberDataInForce, mrss models.MemberRatingResultSummary) float64 {
	val, ok := resolvePremiumSourceField(field, b, mrss)
	if !ok {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func resolvePremiumSourceField(path string, b models.GPricingMemberDataInForce, mrss models.MemberRatingResultSummary) (any, bool) {
	switch strings.ToLower(path) {
	case "member.member_name", "member_name":
		return b.MemberName, true
	case "member.scheme_name", "scheme_name":
		return b.SchemeName, true
	case "member.member_id_number", "member_id_number":
		return b.MemberIdNumber, true
	case "employee_number":
		return b.EmployeeNumber, true
	case "category":
		return b.SchemeCategory, true
	case "annual_salary":
		return b.AnnualSalary, true
	case "gla.annual_premium", "gla_annual_premium":
		return utils.FloatPrecision(monthlyOfficeFromRiskProportion(b.AnnualSalary, mrss.ExpProportionGlaAnnualRiskPremiumSalary, &mrss), AccountingPrecision), true
	case "sgla.annual_premium", "sgla_annual_premium":
		return utils.FloatPrecision(monthlyOfficeFromRiskProportion(b.AnnualSalary, mrss.ExpProportionSglaAnnualRiskPremiumSalary, &mrss), AccountingPrecision), true
	case "ptd.annual_premium", "ptd_annual_premium":
		return utils.FloatPrecision(monthlyOfficeFromRiskProportion(b.AnnualSalary, mrss.ExpProportionPtdAnnualRiskPremiumSalary, &mrss), AccountingPrecision), true
	case "ci.annual_premium", "ci_annual_premium":
		return utils.FloatPrecision(monthlyOfficeFromRiskProportion(b.AnnualSalary, mrss.ExpProportionCiAnnualRiskPremiumSalary, &mrss), AccountingPrecision), true
	case "ttd.annual_premium", "ttd_annual_premium":
		return utils.FloatPrecision(monthlyOfficeFromRiskProportion(b.AnnualSalary, mrss.ExpProportionTtdAnnualRiskPremiumSalary, &mrss), AccountingPrecision), true
	case "phi.annual_premium", "phi_annual_premium":
		return utils.FloatPrecision(monthlyOfficeFromRiskProportion(b.AnnualSalary, mrss.ExpProportionPhiAnnualRiskPremiumSalary, &mrss), AccountingPrecision), true
	case "total_annual_premium_excl_funeral", "premium.total_annual_excl_funeral":
		return utils.FloatPrecision(b.AnnualSalary*mrss.ProportionExpTotalPremiumExclFuneralSalary/12.0, AccountingPrecision), true
	case "total_annual_funeral_premium", "total_funeral_premium":
		return utils.FloatPrecision(mrss.ExpTotalFunAnnualPremiumPerMember/12.0, AccountingPrecision), true
	case "total_annual_premium", "total_premium":
		return utils.FloatPrecision((b.AnnualSalary*mrss.ProportionExpTotalPremiumExclFuneralSalary+mrss.ExpTotalFunAnnualPremiumPerMember)/12.0, AccountingPrecision), true
	//case "funeral.spouse_annual_premium", "spouse_funeral_annual_premium":
	//	return b.SpouseFuneralAnnualPremium, true
	//case "funeral.children_annual_premium", "children_funeral_annual_premium":
	//	return b.ChildrenFuneralAnnualPremium, true
	//case "funeral.dependants_annual_premium", "dependants_funeral_annual_premium":
	//	return b.DependantsFuneralAnnualPremium, true

	default:
		return "", false
	}
}
*/
// End of disabled premium bordereaux block.

// ================= Claims Bordereaux =================
// Claim bordereaux generation is temporarily disabled — only member is a valid
// use case right now. The entire block below is commented out so it can be
// restored in one go.
/*
// GenerateClaimBordereaux creates a Claims bordereaux from GroupSchemeClaim data
func GenerateClaimBordereaux(ctx context.Context, req GenerateBordereauxRequest) (BordereauxReportMeta, error) {
	startTime := time.Now()
	if strings.ToLower(req.Type) != "claim" {
		return BordereauxReportMeta{}, errors.New("unsupported bordereaux type")
	}

	start, end, err := resolvePeriod(req)
	if err != nil {
		return BordereauxReportMeta{}, err
	}

	// Load template and verify type
	tpl, err := GetBordereauxTemplateByID(req.TemplateID)
	if err != nil {
		return BordereauxReportMeta{}, err
	}

	// Increment usage
	_ = IncrementTemplateUsage(tpl.ID)

	if tpl.Type != "claim" && tpl.Type != "claims" {
		return BordereauxReportMeta{}, fmt.Errorf("template type %s not supported for claim generation", tpl.Type)
	}

	// Fetch claims by scheme ids then period filter
	var rows []models.GroupSchemeClaim
	q := DB.WithContext(ctx)
	if len(req.SchemeIDs) > 0 {
		q = q.Where("scheme_id IN ?", req.SchemeIDs)
	}
	if err := q.Find(&rows).Error; err != nil {
		return BordereauxReportMeta{}, err
	}

	// Filter by period using DateRegistered if available else DateOfEvent
	filtered := make([]models.GroupSchemeClaim, 0, len(rows))
	for _, r := range rows {
		t := firstParsableDate(r.DateRegistered, r.DateOfEvent)
		if t.IsZero() {
			continue
		}
		if (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end)) {
			// Optionally validate ID numbers
			if req.ValidateIDNumbers {
				if r.MemberIDNumber != "" && !ValidateSouthAfricanID(r.MemberIDNumber) {
					if req.ExcludeInvalid {
						continue
					}
				}
			}
			filtered = append(filtered, r)
		}
	}

	format := strings.ToLower(req.OutputFormat)
	if format == "" {
		format = strings.ToLower(tpl.Format)
	}
	if format == "" {
		format = "excel"
	}
	reportDir := filepath.Join("data", "reports")
	_ = os.MkdirAll(reportDir, 0o755)
	fileName := fmt.Sprintf("bordereaux_claim_%d%02d_%s.xlsx", start.Year(), int(start.Month()), newFileToken())
	if format == "csv" {
		fileName = strings.TrimSuffix(fileName, ".xlsx") + ".csv"
	}
	outPath := filepath.Join(reportDir, fileName)
	downloadURL := fmt.Sprintf("/group-pricing/bordereaux/download/%s", fileName)

	// Back-join member employee numbers — GroupSchemeClaim has no EmployeeNumber,
	// so templates that request employee_number need a lookup keyed by ID number.
	empNoByID := buildClaimEmployeeNumberLookup(ctx, req.SchemeIDs, filtered)

	switch format {
	case "excel", "xlsx":
		if err := writeExcelClaim(outPath, tpl, filtered, req, empNoByID); err != nil {
			return BordereauxReportMeta{}, err
		}
	case "csv":
		if err := writeCSVClaim(outPath, tpl, filtered, req, empNoByID); err != nil {
			return BordereauxReportMeta{}, err
		}
	default:
		return BordereauxReportMeta{}, errors.New("unsupported output format")
	}

	var fileSize int64
	if info, err := os.Stat(outPath); err == nil {
		fileSize = info.Size()
	}

	generatedID, _ := generateUniqueBordereauxID(DB.WithContext(ctx))

	// Persist claim data fields
	var claimData []models.ClaimBordereauxData
	for _, b := range filtered {
		convertToMonth := func(month int) string {
			months := []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
			if month >= 1 && month <= 12 {
				return months[month]
			}
			return ""
		}

		eventDate := firstParsableDate(b.DateOfEvent)
		var eventDatePtr *time.Time
		if !eventDate.IsZero() {
			eventDatePtr = &eventDate
		}

		data := models.ClaimBordereauxData{
			BordereauxID:   generatedID,
			SchemeName:     b.SchemeName,
			MemberName:     b.MemberName,
			MemberIdNumber: b.MemberIDNumber,
			ClaimNumber:    b.ClaimNumber,
			Category:       b.MemberType, // In GroupSchemeClaim, Category might be MemberType
			Month:          convertToMonth(req.Month),
			Year:           req.Year,
			Period:         fmt.Sprintf("%s %d", convertToMonth(req.Month), req.Year),
			EventDate:      eventDatePtr,
			ClaimAmount:    b.ClaimAmount,
			ClaimType:      b.BenefitAlias,
			Status:         b.Status,
		}
		claimData = append(claimData, data)
	}

	if len(claimData) > 0 {
		appLog.WithField("count", len(claimData)).Info("Persisting claim bordereaux data")
		if err := DB.WithContext(ctx).CreateInBatches(claimData, 500).Error; err != nil {
			appLog.WithField("error", err.Error()).Error("Failed to persist claim bordereaux data")
			return BordereauxReportMeta{}, fmt.Errorf("failed to persist claim bordereaux data: %w", err)
		}
	}

	userEmail, _ := ctx.Value(appLog.UserEmailKey).(string)
	schemeName := resolveSchemeNames(req.SchemeIDs)
	insurerName := resolveInsurerName(tpl.InsurerID)
	now := time.Now()

	_ = DB.WithContext(ctx).Create(&models.GeneratedBordereaux{
		GeneratedID:    generatedID,
		TemplateID:     req.TemplateID,
		Type:           "claim",
		FileName:       fileName,
		FilePath:       outPath,
		FileSize:       fileSize,
		Records:        len(filtered),
		ProcessingTime: time.Since(startTime).String(),
		SchemeName:     schemeName,
		InsurerName:    insurerName,
		Status:         "generated",
		Progress:       33,
		LastUpdated:    time.Now(),
		SLAStatus:      "on_time",
		SubmissionDate: &now,
		PeriodStart:    start,
		PeriodEnd:      end,
		CreatedBy:      userEmail,
		CreatedAt:      time.Now(),
		Timeline: []models.BordereauxTimeline{
			{
				Date:        time.Now(),
				Type:        "generated",
				Title:       "Generated",
				Description: "Bordereaux generated successfully",
			},
		},
	})

	return BordereauxReportMeta{
		BordereauxID:   req.TemplateID,
		GeneratedID:    generatedID,
		FilePath:       outPath,
		FileName:       fileName,
		FileExtension:  filepath.Ext(fileName),
		DownloadURL:    downloadURL,
		Format:         format,
		Records:        len(filtered),
		FileSize:       fileSize,
		ProcessingTime: time.Since(startTime).String(),
		SchemeName:     schemeName,
		InsurerName:    insurerName,
		Status:         "generated",
		Progress:       33,
		LastUpdated:    time.Now(),
		SLAStatus:      "on_time",
		PeriodStart:    start,
		PeriodEnd:      end,
		Timeline: []BordereauxTimeline{
			{
				Date:        time.Now(),
				Type:        "generated",
				Title:       "Generated",
				Description: "Bordereaux generated successfully",
			},
		},
	}, nil
}

func writeExcelClaim(path string, tpl models.BordereauxTemplate, rows []models.GroupSchemeClaim, req GenerateBordereauxRequest, empNoByID map[string]string) error {
	f := excelize.NewFile()
	sheet := "Bordereaux"
	f.SetSheetName("Sheet1", sheet)
	for i, m := range tpl.FieldMappings {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, m.TargetField)
	}
	for rIdx, b := range rows {
		values, _ := mapFieldsClaim(tpl, b, req, empNoByID)
		for cIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
			_ = f.SetCellValue(sheet, cell, v)
		}
	}
	return f.SaveAs(path)
}

func writeCSVClaim(path string, tpl models.BordereauxTemplate, rows []models.GroupSchemeClaim, req GenerateBordereauxRequest, empNoByID map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	header := make([]string, len(tpl.FieldMappings))
	for i, m := range tpl.FieldMappings {
		header[i] = m.TargetField
	}
	file.WriteString(strings.Join(header, ",") + "\n")
	for _, b := range rows {
		values, _ := mapFieldsClaim(tpl, b, req, empNoByID)
		strVals := make([]string, len(values))
		for i, v := range values {
			strVals[i] = fmt.Sprint(v)
		}
		file.WriteString(strings.Join(strVals, ",") + "\n")
	}
	return nil
}

func mapFieldsClaim(tpl models.BordereauxTemplate, b models.GroupSchemeClaim, req GenerateBordereauxRequest, empNoByID map[string]string) ([]any, []string) {
	vals := make([]any, 0, len(tpl.FieldMappings))
	var errs []string
	for _, m := range tpl.FieldMappings {
		v, ok := resolveClaimSourceField(m.SourceField, b, empNoByID)
		if !ok && m.Required {
			errs = append(errs, fmt.Sprintf("missing required field: %s", m.SourceField))
		}

		// mask member_id or member_id_number
		if strings.EqualFold(m.SourceField, "member.id_number") || strings.EqualFold(m.SourceField, "member_id_number") {
			if s, ok := v.(string); ok && s != "" {
				v = maskID(s)
			}
		}

		if req.ValidateIDNumbers && strings.Contains(strings.ToLower(m.TargetField), "id") {
			s := fmt.Sprint(v)
			if s != "" && !ValidateSouthAfricanID(s) {
				errs = append(errs, "invalid SA ID")
			}
		}
		vals = append(vals, v)
	}
	return vals, errs
}

// buildClaimEmployeeNumberLookup returns a map of member_id_number ->
// employee_number, sourced from GPricingMemberDataInForce. Used to populate the
// employee_number column on claims bordereaux, since GroupSchemeClaim does not
// carry an employee number directly.
func buildClaimEmployeeNumberLookup(ctx context.Context, schemeIDs []int, claims []models.GroupSchemeClaim) map[string]string {
	out := map[string]string{}
	if len(claims) == 0 {
		return out
	}

	ids := make([]string, 0, len(claims))
	seen := map[string]struct{}{}
	for _, c := range claims {
		if c.MemberIDNumber == "" {
			continue
		}
		if _, ok := seen[c.MemberIDNumber]; ok {
			continue
		}
		seen[c.MemberIDNumber] = struct{}{}
		ids = append(ids, c.MemberIDNumber)
	}
	if len(ids) == 0 {
		return out
	}

	type row struct {
		MemberIdNumber string
		EmployeeNumber string
	}
	var rows []row
	q := DB.WithContext(ctx).
		Model(&models.GPricingMemberDataInForce{}).
		Select("member_id_number, employee_number").
		Where("member_id_number IN ?", ids).
		Order("creation_date DESC")
	if len(schemeIDs) > 0 {
		q = q.Where("scheme_id IN ?", schemeIDs)
	}
	if err := q.Find(&rows).Error; err != nil {
		return out
	}
	// First (latest) row per member wins.
	for _, r := range rows {
		if _, ok := out[r.MemberIdNumber]; ok {
			continue
		}
		if r.EmployeeNumber == "" {
			continue
		}
		out[r.MemberIdNumber] = r.EmployeeNumber
	}
	return out
}

func resolveClaimSourceField(path string, b models.GroupSchemeClaim, empNoByID map[string]string) (any, bool) {
	switch strings.ToLower(path) {
	case "claim.claim_number", "claim_number":
		return b.ClaimNumber, true
	case "employee_number":
		// GroupSchemeClaim has no employee number; join via member ID.
		return empNoByID[b.MemberIDNumber], true
	case "member.name", "member_name":
		return b.MemberName, true
	case "member.id_number", "member_id_number":
		return b.MemberIDNumber, true
	case "scheme.id", "scheme_id":
		return b.SchemeId, true
	case "scheme.name", "scheme_name":
		return b.SchemeName, true
	case "benefit.type", "benefit_type":
		return b.BenefitAlias, true
	case "benefit_name":
		return b.BenefitName, true
	case "member.type", "member_type":
		return b.MemberType, true
	case "claim.date_of_event", "date_of_event":
		return b.DateOfEvent, true
	case "claim.date_notified", "date_notified":
		return b.DateNotified, true
	case "claim.status", "status":
		return b.Status, true
	case "claim.amount", "claim_amount":
		return b.ClaimAmount, true
	case "claim.date_registered", "date_registered":
		return b.DateRegistered, true
	case "claimant.name", "claimant_name":
		return b.ClaimantName, true
	case "claimant.id_number", "claimant_id_number":
		return b.ClaimantIDNumber, true
	case "claimant.relationship", "relationship_to_member":
		return b.RelationshipToMember, true
	default:
		return "", false
	}
}
*/
// End of disabled claim bordereaux block.

// firstParsableDate tries multiple formats and returns the first successfully parsed date
func firstParsableDate(values ...string) time.Time {
	layouts := []string{"2006-01-02", "2006/01/02", "02-01-2006", time.RFC3339}
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, v); err == nil {
				return t
			}
		}
	}
	return time.Time{}
}

func GetMemberBordereaux(schemeIDs []int, req GenerateBordereauxRequest, template models.BordereauxTemplate) ([]models.GPricingMemberDataInForce, error) {
	var rows []models.GPricingMemberDataInForce
	var err error

	sourceFields := extractSourceFields(template)
	selectedVariables := strings.Join(sourceFields, ", ")

	if !req.IncludeBeneficiaries {
		if req.IncludeTerminated {
			err = DB.WithContext(ctx).
				Table("g_pricing_member_data_in_forces"). // Explicitly specify the table name
				Select(selectedVariables). // Use .Select() to specify columns
				Where("scheme_id IN ?", schemeIDs). // Example condition
				Find(&rows).
				Error
		}
		if !req.IncludeTerminated {
			err = DB.WithContext(ctx).
				Table("g_pricing_member_data_in_forces"). // Explicitly specify the table name
				Select(selectedVariables). // Use .Select() to specify columns
				Where("scheme_id IN ? and status=?", schemeIDs, models.StatusActive). // Example condition
				Find(&rows).
				Error
		}

	}
	if req.IncludeBeneficiaries {
		if req.IncludeTerminated {
			err = DB.WithContext(ctx).
				Table("g_pricing_member_data_in_forces").
				Select(selectedVariables).
				Where("scheme_id IN ?", schemeIDs).
				// The string inside Preload() must match the field name in the struct (i.e., 'Beneficiaries')
				Preload("Beneficiaries").
				Find(&rows).
				Error
		}
		if !req.IncludeTerminated {
			err = DB.WithContext(ctx).
				Table("g_pricing_member_data_in_forces").
				Select(selectedVariables).
				Where("scheme_id IN ? and status=?", schemeIDs, models.StatusActive).
				// The string inside Preload() must match the field name in the struct (i.e., 'Beneficiaries')
				Preload("Beneficiaries").
				Find(&rows).
				Error
		}

	}

	if err != nil {
		return nil, fmt.Errorf("failed to read inforce data: %w", err)
	}
	return rows, nil

}

func extractSourceFields(template models.BordereauxTemplate) []string {
	var sourceFields []string

	// Iterate through the array of field mappings
	for _, mapping := range template.FieldMappings {
		// Append the SourceField value to our result slice
		sourceFields = append(sourceFields, mapping.SourceField)
	}

	return sourceFields
}

// GetAllGeneratedBordereaux returns all generated bordereaux records ordered by creation date descending
func GetAllGeneratedBordereaux() ([]models.GeneratedBordereaux, error) {
	var list []models.GeneratedBordereaux
	if err := DB.Preload("Timeline").Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetBordereauxDashboardStats() (models.BordereauxDashboardStats, error) {
	var stats models.BordereauxDashboardStats
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Generated This Month
	var generatedCount int64
	if err := DB.Model(&models.GeneratedBordereaux{}).Where("created_at >= ?", startOfMonth).Count(&generatedCount).Error; err != nil {
		return stats, err
	}
	stats.GeneratedThisMonth = int(generatedCount)

	// Pending Submissions
	var pendingCount int64
	if err := DB.Model(&models.GeneratedBordereaux{}).Where("status IN ?", []string{"pending", "submitted", "Pending", "Submitted"}).Count(&pendingCount).Error; err != nil {
		return stats, err
	}
	stats.PendingSubmissions = int(pendingCount)

	// Reconciled This Week
	// Assuming reconciled means status is 'processed' or similar in BordereauxConfirmation
	// and we want to count those whose last_reconciled is in the current week.
	// Find start of week (Sunday or Monday, let's use Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -weekday+1)

	var reconciledCount int64
	if err := DB.Model(&models.BordereauxConfirmation{}).
		Where("status = ? AND last_reconciled >= ?", "processed", startOfWeek).
		Count(&reconciledCount).Error; err != nil {
		return stats, err
	}
	stats.ReconciledThisWeek = int(reconciledCount)

	// Active Templates
	var activeTemplatesCount int64
	if err := DB.Model(&models.BordereauxTemplate{}).Where("status = ?", "active").Count(&activeTemplatesCount).Error; err != nil {
		return stats, err
	}
	stats.ActiveTemplates = int(activeTemplatesCount)

	return stats, nil
}

// GetGeneratedBordereauxByGeneratedID retrieves a single generated bordereaux record by its generated ID
func GetGeneratedBordereauxByGeneratedID(generatedID string) (models.GeneratedBordereaux, error) {
	var record models.GeneratedBordereaux
	if err := DB.Preload("Timeline").Where("generated_id = ?", generatedID).First(&record).Error; err != nil {
		return record, err
	}
	return record, nil
}

// ErrBordereauxNotFound is returned when a generated bordereaux row matching
// a filename cannot be located. Callers should map this to HTTP 404.
var ErrBordereauxNotFound = errors.New("bordereaux not found")

// ErrBordereauxNotAuthorized is returned when the requesting user is not the
// creator, reviewer or approver of the bordereaux. Callers map this to 403.
var ErrBordereauxNotAuthorized = errors.New("not authorized to access this bordereaux")

// AuthorizeBordereauxDownload looks up the generated bordereaux by file name,
// confirms the user may access it (creator, reviewer, or approver), writes a
// download audit entry, and returns the record. 404 on missing row prevents
// filename enumeration against the reports directory.
func AuthorizeBordereauxDownload(fileName string, user models.AppUser) (models.GeneratedBordereaux, error) {
	var gen models.GeneratedBordereaux
	if err := DB.Where("file_name = ?", fileName).First(&gen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gen, ErrBordereauxNotFound
		}
		return gen, err
	}
	if gen.CreatedBy != user.UserEmail &&
		gen.ReviewedBy != user.UserName &&
		gen.ApprovedBy != user.UserName {
		return gen, ErrBordereauxNotAuthorized
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "generated_bordereaux",
		EntityID:  gen.GeneratedID,
		Action:    "DOWNLOAD",
		ChangedBy: user.UserName,
	}, gen, gen)
	return gen, nil
}

// AddBordereauxTimelineEntry adds a new timeline entry to a generated bordereaux
func AddBordereauxTimelineEntry(generatedID string, entry models.BordereauxTimeline) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var record models.GeneratedBordereaux
		if err := tx.Where("generated_id = ?", generatedID).First(&record).Error; err != nil {
			return err
		}

		entry.GeneratedBordereauxID = generatedID
		if entry.Date.IsZero() {
			entry.Date = time.Now()
		}
		if err := tx.Create(&entry).Error; err != nil {
			return err
		}

		// Update parent record
		updates := map[string]interface{}{
			"last_updated": time.Now(),
		}
		if entry.Type != "" {
			updates["status"] = strings.Title(entry.Type)
		}
		// Submitting to the employer is the final step — the cycle is
		// considered complete once a bordereaux has been submitted.
		if entry.Type == "submitted" || entry.Type == "confirmed" || entry.Type == "reconciled" {
			if entry.Type == "submitted" {
				now := time.Now()
				updates["submission_date"] = &now
			}
			updates["progress"] = 100
		} else if entry.Type == "generated" {
			updates["progress"] = 33
		}

		return tx.Model(&record).Updates(updates).Error
	})
}

// GetBordereauxConfigurations returns all saved configurations
func GetBordereauxConfigurations() ([]models.BordereauxConfiguration, error) {
	var configs []models.BordereauxConfiguration
	err := DB.Order("created_at desc").Find(&configs).Error
	return configs, err
}

// GetBordereauxConfiguration returns a single configuration by ID
func GetBordereauxConfiguration(id string) (models.BordereauxConfiguration, error) {
	var config models.BordereauxConfiguration
	err := DB.First(&config, id).Error
	return config, err
}

// SaveBordereauxConfiguration saves a new configuration
func SaveBordereauxConfiguration(config models.BordereauxConfiguration) (models.BordereauxConfiguration, error) {
	err := DB.Create(&config).Error
	return config, err
}

// UpdateBordereauxConfiguration updates an existing configuration
func UpdateBordereauxConfiguration(id string, config models.BordereauxConfiguration) (models.BordereauxConfiguration, error) {
	var existing models.BordereauxConfiguration
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}

	existing.Name = config.Name
	existing.Description = config.Description
	existing.ConfigData = config.ConfigData

	err := DB.Save(&existing).Error
	return existing, err
}

// DeleteBordereauxConfiguration deletes a configuration
func DeleteBordereauxConfiguration(id string) error {
	return DB.Delete(&models.BordereauxConfiguration{}, id).Error
}

// UpdateConfigurationUsage updates the last_used_at timestamp
func UpdateConfigurationUsage(id string) error {
	return DB.Model(&models.BordereauxConfiguration{}).Where("id = ?", id).Update("last_used_at", time.Now()).Error
}

// GetGeneratedBordereauxData retrieves the data records for a generated bordereaux
func GetGeneratedBordereauxData(generatedID string) (any, error) {
	var record models.GeneratedBordereaux
	if err := DB.Where("generated_id = ?", generatedID).First(&record).Error; err != nil {
		return nil, err
	}

	ids := []string{generatedID}

	// If it's a zip file, it means there are multiple bordereaux files inside.
	// We need to use their filenames (which are their respective IDs) to retrieve the data.
	if strings.HasSuffix(strings.ToLower(record.FileName), ".zip") {
		r, err := zip.OpenReader(record.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open zip file: %w", err)
		}
		defer r.Close()

		ids = []string{}
		for _, f := range r.File {
			// Extract ID from filename (e.g., BRD-20240101-123456.xlsx -> BRD-20240101-123456)
			id := strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
			ids = append(ids, id)
		}
	}

	switch strings.ToLower(record.Type) {
	case "premium", "premiums":
		var data []models.PremiumBordereauxData
		if err := DB.Where("bordereaux_id IN ?", ids).Find(&data).Error; err != nil {
			return nil, err
		}
		return data, nil
	case "member":
		var data []models.MemberBordereauxData
		if err := DB.Where("bordereaux_id IN ?", ids).Find(&data).Error; err != nil {
			return nil, err
		}
		return data, nil
	case "claim":
		var data []models.ClaimBordereauxData
		if err := DB.Where("bordereaux_id IN ?", ids).Find(&data).Error; err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("unsupported bordereaux type: %s", record.Type)
	}
}

// SubmitBordereauxBatch processes a batch submission of bordereaux files
func SubmitBordereauxBatch(ctx context.Context, req BordereauxBatchSubmitRequest, user models.AppUser) error {
	logger := appLog.WithContext(ctx)

	for _, id := range req.BordereauxIDs {
		err := DB.Transaction(func(tx *gorm.DB) error {
			var bordereaux models.GeneratedBordereaux
			if err := tx.First(&bordereaux, id).Error; err != nil {
				return fmt.Errorf("bordereaux %d not found: %w", id, err)
			}

			if err := ValidateGeneratedBordereauxTransition(bordereaux.Status, StatusGeneratedSubmitted); err != nil {
				return fmt.Errorf("bordereaux %d: %w", id, err)
			}

			// Retrieve GroupScheme to get contact email
			// GeneratedBordereaux has SchemeName, but we need the ID to be sure or search by name
			// Let's check if we can find the scheme by name
			var scheme models.GroupScheme
			if err := tx.Where("name = ?", bordereaux.SchemeName).First(&scheme).Error; err != nil {
				logger.WithField("scheme_name", bordereaux.SchemeName).Warn("Scheme not found for bordereaux submission")
				// If scheme not found, we might still want to proceed or fail.
				// The requirement says "retrieve the GeneratedBordereaux which will contain the generatedId as well as the schemes associated with the bordereaux"
				// GeneratedBordereaux model has SchemeName.
			}

			if req.DeliveryMethod == "email" {
				email := scheme.ContactEmail
				if email == "" {
					logger.WithField("bordereaux_id", id).Warn("No contact email for scheme, using mock")
					email = "mock-recipient@example.com"
				}

				// Mock sending email
				logger.WithFields(map[string]interface{}{
					"bordereaux_id": id,
					"email":         email,
					"file":          bordereaux.FileName,
					"message":       req.Message,
				}).Info("Mock sending bordereaux via email")
			}

			// Add timeline entry
			timelineEntry := models.BordereauxTimeline{
				GeneratedBordereauxID: bordereaux.GeneratedID,
				Date:                  req.SubmissionDate,
				Type:                  "submitted",
				Title:                 "Bordereaux Submitted",
				Description:           fmt.Sprintf("Method: %s. Message: %s", req.DeliveryMethod, req.Message),
			}
			if err := tx.Create(&timelineEntry).Error; err != nil {
				return fmt.Errorf("failed to create timeline entry: %w", err)
			}

			before := bordereaux
			// Update bordereaux status and submission date. Submission is the
			// final step in the cycle, so progress reaches 100%.
			updates := map[string]interface{}{
				"status":          "submitted",
				"submission_date": req.SubmissionDate,
				"last_updated":    time.Now(),
				"progress":        100,
			}
			if err := tx.Model(&bordereaux).Updates(updates).Error; err != nil {
				return fmt.Errorf("failed to update bordereaux status: %w", err)
			}
			// Reload so the audit "after" reflects the persisted state.
			var after models.GeneratedBordereaux
			if err := tx.First(&after, bordereaux.ID).Error; err == nil {
				_ = writeAudit(tx, AuditContext{
					Area:      "group-pricing",
					Entity:    "generated_bordereaux",
					EntityID:  bordereaux.GeneratedID,
					Action:    "UPDATE",
					ChangedBy: user.UserName,
				}, before, after)
			}
			return nil
		})

		if err != nil {
			logger.WithError(err).Errorf("Failed to process bordereaux %d", id)
			// Continue with other bordereaux in the batch even if one fails?
			// The requirement doesn't specify. Usually batch should be atomic or report partial success.
			// Let's fail for now to be safe.
			return err
		}

		// Fire notification after successful submission
		var submitted models.GeneratedBordereaux
		if DB.First(&submitted, id).Error == nil {
			go NotifyBordereauxSubmitted(submitted, user)
		}
	}

	return nil
}

// ReviewGeneratedBordereaux marks a generated bordereaux as under review (draft/generated → reviewed)
func ReviewGeneratedBordereaux(generatedID string, notes string, user models.AppUser) (models.GeneratedBordereaux, error) {
	var brd models.GeneratedBordereaux
	if err := DB.Where("generated_id = ?", generatedID).First(&brd).Error; err != nil {
		return brd, fmt.Errorf("bordereaux not found: %w", err)
	}
	if err := ValidateGeneratedBordereauxTransition(brd.Status, StatusGeneratedReviewed); err != nil {
		return brd, err
	}
	before := brd
	now := time.Now()
	brd.Status = StatusGeneratedReviewed
	brd.ReviewedBy = user.UserName
	brd.ReviewedAt = &now
	brd.LastUpdated = now
	if notes != "" {
		brd.ReturnReason = ""
	}
	if err := DB.Save(&brd).Error; err != nil {
		return brd, fmt.Errorf("failed to save review: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "generated_bordereaux",
		EntityID:  brd.GeneratedID,
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, brd)
	timeline := models.BordereauxTimeline{
		GeneratedBordereauxID: brd.GeneratedID,
		Date:                  now,
		Type:                  "reviewed",
		Title:                 "Bordereaux Reviewed",
		Description:           fmt.Sprintf("Reviewed by %s. %s", user.UserName, notes),
	}
	DB.Create(&timeline)
	go NotifyBordereauxReviewed(brd, user)
	return brd, nil
}

// ApproveGeneratedBordereaux marks a reviewed bordereaux as approved (reviewed → approved)
func ApproveGeneratedBordereaux(generatedID string, notes string, user models.AppUser) (models.GeneratedBordereaux, error) {
	var brd models.GeneratedBordereaux
	if err := DB.Where("generated_id = ?", generatedID).First(&brd).Error; err != nil {
		return brd, fmt.Errorf("bordereaux not found: %w", err)
	}
	if err := ValidateGeneratedBordereauxTransition(brd.Status, StatusGeneratedApproved); err != nil {
		return brd, err
	}
	before := brd
	now := time.Now()
	brd.Status = StatusGeneratedApproved
	brd.ApprovedBy = user.UserName
	brd.ApprovedAt = &now
	brd.LastUpdated = now
	if err := DB.Save(&brd).Error; err != nil {
		return brd, fmt.Errorf("failed to save approval: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "generated_bordereaux",
		EntityID:  brd.GeneratedID,
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, brd)
	timeline := models.BordereauxTimeline{
		GeneratedBordereauxID: brd.GeneratedID,
		Date:                  now,
		Type:                  "approved",
		Title:                 "Bordereaux Approved",
		Description:           fmt.Sprintf("Approved by %s. %s", user.UserName, notes),
	}
	DB.Create(&timeline)
	go NotifyBordereauxApproved(brd, user)
	return brd, nil
}

// RegenerateGeneratedBordereaux re-runs generation for a draft bordereaux
// in-place: same generated_id, refreshed file, status returns to "generated".
// Requires the original request to have been persisted on the row.
func RegenerateGeneratedBordereaux(ctx context.Context, generatedID string, user models.AppUser) (models.GeneratedBordereaux, error) {
	var brd models.GeneratedBordereaux
	if err := DB.Where("generated_id = ?", generatedID).First(&brd).Error; err != nil {
		return brd, fmt.Errorf("bordereaux not found: %w", err)
	}
	if err := ValidateGeneratedBordereauxTransition(brd.Status, StatusGeneratedGenerated); err != nil {
		return brd, err
	}
	if len(brd.RequestPayload) == 0 {
		return brd, fmt.Errorf("regenerate unavailable: original request not stored for this bordereaux")
	}

	var req GenerateBordereauxRequest
	if err := json.Unmarshal(brd.RequestPayload, &req); err != nil {
		return brd, fmt.Errorf("stored request payload is invalid: %w", err)
	}
	req.ExistingGeneratedID = generatedID
	if req.Type == "" {
		req.Type = brd.Type
	}

	before := brd
	if _, err := GenerateMemberBordereaux(ctx, req); err != nil {
		return brd, fmt.Errorf("regeneration failed: %w", err)
	}

	// Reload the updated row for the caller and audit trail.
	if err := DB.Where("generated_id = ?", generatedID).First(&brd).Error; err != nil {
		return brd, err
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "generated_bordereaux",
		EntityID:  brd.GeneratedID,
		Action:    "REGENERATE",
		ChangedBy: user.UserName,
	}, before, brd)
	return brd, nil
}

// ReturnOutboundToDraft returns a reviewed or approved bordereaux back to draft for rework
func ReturnOutboundToDraft(generatedID string, reason string, user models.AppUser) (models.GeneratedBordereaux, error) {
	var brd models.GeneratedBordereaux
	if err := DB.Where("generated_id = ?", generatedID).First(&brd).Error; err != nil {
		return brd, fmt.Errorf("bordereaux not found: %w", err)
	}
	if err := ValidateGeneratedBordereauxTransition(brd.Status, StatusGeneratedDraft); err != nil {
		return brd, err
	}
	before := brd
	now := time.Now()
	brd.Status = StatusGeneratedDraft
	brd.ReturnReason = reason
	brd.LastUpdated = now
	if err := DB.Save(&brd).Error; err != nil {
		return brd, fmt.Errorf("failed to return to draft: %w", err)
	}
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "generated_bordereaux",
		EntityID:  brd.GeneratedID,
		Action:    "UPDATE",
		ChangedBy: user.UserName,
	}, before, brd)
	timeline := models.BordereauxTimeline{
		GeneratedBordereauxID: brd.GeneratedID,
		Date:                  now,
		Type:                  "returned",
		Title:                 "Returned to Draft",
		Description:           fmt.Sprintf("Returned by %s. Reason: %s", user.UserName, reason),
	}
	DB.Create(&timeline)
	return brd, nil
}

func GetSchemeCategory(id int, category string) (models.SchemeCategory, error) {
	var cat models.SchemeCategory
	cacheKey := "scheme_categories" + "_" + strconv.Itoa(id) + "_" + category
	cached, found := GroupPricingCache.Get(cacheKey)

	if found {
		result := cached.(models.SchemeCategory)
		return result, nil
	}
	err := DB.Where("quote_id = ? and scheme_category=?", id, category).Find(&cat).Error

	if err != nil {
		fmt.Println(err)
		GroupPricingCache.Set(cacheKey, 0, 1)
		return cat, fmt.Errorf("category not found: %w", err)
	}
	GroupPricingCache.Set(cacheKey, cat, 1)

	return cat, nil
}

// fetchInEffectQuoteIDs filters quoteIDs to only those with in-effect status.
// fetchInEffectQuoteIDs filters quoteIDs to only those with in-effect status.
func fetchInEffectQuoteIDs(ctx context.Context, quoteIDs []int) ([]int, error) {
	var ids []int
	err := DB.WithContext(ctx).
		Model(&models.GroupPricingQuote{}).
		Where("id IN ? AND scheme_quote_status = ?", quoteIDs, "in_effect").
		Pluck("id", &ids).Error
	return ids, err
}

// batchFetchSchemeCategories uses int keys
func batchFetchSchemeCategories(
	ctx context.Context,
	quoteIDCategories map[int][]string, // quoteID -> []category
) (map[int]map[string]models.SchemeCategory, error) {

	if len(quoteIDCategories) == 0 {
		return map[int]map[string]models.SchemeCategory{}, nil
	}

	// result: quoteID -> category -> SchemeCategory
	cache := make(map[int]map[string]models.SchemeCategory, len(quoteIDCategories))

	for quoteID, categories := range quoteIDCategories {
		var batch []models.SchemeCategory
		if err := DB.WithContext(ctx).
			Where("quote_id = ? AND scheme_category IN ?", quoteID, categories).
			Find(&batch).Error; err != nil {
			return nil, fmt.Errorf("batchFetchSchemeCategories quoteID=%d: %w", quoteID, err)
		}

		inner := make(map[string]models.SchemeCategory, len(batch))
		for _, sc := range batch {
			inner[sc.SchemeCategory] = sc
		}
		cache[quoteID] = inner
	}

	return cache, nil
}

// batchFetchSchemeNames uses int keys to match SchemeId field type
func batchFetchSchemeNames(ctx context.Context, ids map[int]struct{}) (map[int]string, error) {
	if len(ids) == 0 {
		return map[int]string{}, nil
	}

	idSlice := make([]int, 0, len(ids))
	for id := range ids {
		idSlice = append(idSlice, id)
	}

	var schemes []models.GroupScheme
	if err := DB.WithContext(ctx).
		Select("id, name").
		Where("id IN ?", idSlice).
		Find(&schemes).Error; err != nil {
		return nil, err
	}

	names := make(map[int]string, len(schemes))
	for _, s := range schemes {
		names[int(s.ID)] = s.Name // cast s.ID if it's uint
	}
	return names, nil
}

// monthlyOfficeFromRiskProportion derives a member's monthly office premium
// share from the scheme-level proportion of annual risk premium to salary.
// Office = risk / (1 - SchemeTotalLoading()), then divided by 12 for monthly.
func monthlyOfficeFromRiskProportion(annualSalary, riskProportion float64, s *models.MemberRatingResultSummary) float64 {
	if s == nil {
		return 0
	}
	denom := 1.0 - s.SchemeTotalLoading()
	if denom <= 0 {
		return 0
	}
	return annualSalary * riskProportion / denom / 12.0
}

// monthIntToName converts a 1-based month integer to its English name.
func monthIntToName(month int) string {
	if month < 1 || month > 12 {
		return ""
	}
	return time.Month(month).String()
}

// CoveredSumsAssured calculates the covered, retained, and ceded sums assured
// for all benefits of a member in force, based on their category and quote.
//
// Parameters:
//   - r        : member in force pricing data (salary, benefit multiples, etc.)
//   - category : scheme category containing free cover limits and retention limits
//   - quote    : quote containing reinsurance / cession configuration
//   - mrss     : pre-resolved restriction maxes and reinsurance cover caps that
//     were applied at pricing time. Each cap of 0 means "no limit"; non-zero
//     values clamp the covered sum after the FCL clamp so this function lands
//     on the same covered figures the pricing flow stored.
//
// Returns:
//   - SumsAssuredResult with all benefit sums assured populated
func CoveredSumsAssured(
	r models.GPricingMemberDataInForce,
	category models.SchemeCategory,
	quote models.GroupPricingQuote,
	reinsuranceTreaty models.ReinsuranceTreaty,
	mrss models.MemberRatingResultSummary,
) models.SumsAssuredResult {

	// --- Helper: compute covered, retained, ceded for a given raw sum assured ---
	// covered  = min(rawSA, freeCoverLimit)
	// retained = min(covered, retentionLimit)
	// ceded    = covered - retained
	var glaCovered, glaRetained, glaCeded float64
	var ptdCovered, ptdRetained, ptdCeded float64
	var ciCovered, ciRetained, ciCeded float64
	var sglaCovered, sglaRetained, sglaCeded float64
	var phiCovered, phiRetained, phiCeded float64
	var ttdCovered, ttdRetained, ttdCeded float64
	var mmFuneralSACovered, mmFuneralRetained, mmFuneralCeded float64
	var spFuneralSACovered, spFuneralRetained, spFuneralCeded float64
	var chFuneralSACovered, chFuneralRetained, chFuneralCeded float64
	var parFuneralSACovered, parFuneralRetained, parFuneralCeded float64
	var depFuneralSACovered, depFuneralRetained, depFuneralCeded float64

	fcl := category.FreeCoverLimit

	// --- GLA ---
	if category.GlaBenefit {
		glaRaw := r.AnnualSalary * r.Benefits.GlaMultiple
		glaRaw = applyMaxCoverCap(applyMaxCoverCap(glaRaw, mrss.MaximumGlaCover), mrss.ReinsMaxGlaCover)
		glaCovered, glaRetained, glaCeded = SplitLumpSumSA(glaRaw, fcl, reinsuranceTreaty)
	}

	// --- PTD ---
	if category.PtdBenefit {
		ptdRaw := r.AnnualSalary * r.Benefits.PtdMultiple
		ptdRaw = applyMaxCoverCap(applyMaxCoverCap(ptdRaw, mrss.MaximumPtdCover), mrss.ReinsMaxPtdCover)
		ptdCovered, ptdRetained, ptdCeded = SplitLumpSumSA(ptdRaw, fcl, reinsuranceTreaty)
	}

	// --- CI ---
	if category.CiBenefit {
		ciRaw := r.AnnualSalary * r.Benefits.CiMultiple
		ciRaw = applyMaxCoverCap(applyMaxCoverCap(ciRaw, mrss.SevereIllnessMaximumBenefit), mrss.ReinsMaxCiCover)
		ciCovered, ciRetained, ciCeded = SplitLumpSumSA(ciRaw, fcl, reinsuranceTreaty)
	}

	// --- SGLA ---
	if category.SglaBenefit {
		sglaRaw := r.AnnualSalary * r.Benefits.SglaMultiple
		sglaRaw = applyMaxCoverCap(applyMaxCoverCap(sglaRaw, mrss.SpouseGlaMaximumBenefit), mrss.ReinsMaxSglaCover)
		sglaCovered, sglaRetained, sglaCeded = SplitLumpSumSA(sglaRaw, fcl, reinsuranceTreaty)
	}

	// --- PHI (monthly benefit) ---
	if category.PhiBenefit {
		phiRaw := r.AnnualSalary * r.Benefits.PhiMultiple
		phiRaw = applyMaxCoverCap(applyMaxCoverCap(phiRaw, mrss.PhiMaximumMonthlyBenefit), mrss.ReinsMaxPhiCover)
		phiCovered, phiRetained, phiCeded = SplitIncome(phiRaw, fcl, reinsuranceTreaty)
	}

	// --- TTD (monthly benefit) ---
	if category.TtdBenefit {
		ttdRaw := r.AnnualSalary * r.Benefits.TtdMultiple
		ttdRaw = applyMaxCoverCap(applyMaxCoverCap(ttdRaw, mrss.TtdMaximumMonthlyBenefit), mrss.ReinsMaxTtdCover)
		ttdCovered, ttdRetained, ttdCeded = SplitIncome(ttdRaw, fcl, reinsuranceTreaty)
	}

	// --- Funeral Sum Assured
	if category.FamilyFuneralBenefit {
		mmFuneralSARaw := applyMaxCoverCap(category.FamilyFuneralMainMemberFuneralSumAssured, mrss.ReinsMaxFunCover)
		mmFuneralSACovered, mmFuneralRetained, mmFuneralCeded = SplitFuneralSA(mmFuneralSARaw, reinsuranceTreaty)

		spFuneralSARaw := category.FamilyFuneralMainMemberFuneralSumAssured
		spFuneralSACovered, spFuneralRetained, spFuneralCeded = SplitFuneralSA(spFuneralSARaw, reinsuranceTreaty)

		chFuneralSARaw := category.FamilyFuneralMainMemberFuneralSumAssured
		chFuneralSACovered, chFuneralRetained, chFuneralCeded = SplitFuneralSA(chFuneralSARaw, reinsuranceTreaty)

		parFuneralSARaw := category.FamilyFuneralMainMemberFuneralSumAssured
		parFuneralSACovered, parFuneralRetained, parFuneralCeded = SplitFuneralSA(parFuneralSARaw, reinsuranceTreaty)

		depFuneralSARaw := category.FamilyFuneralMainMemberFuneralSumAssured
		depFuneralSACovered, depFuneralRetained, depFuneralCeded = SplitFuneralSA(depFuneralSARaw, reinsuranceTreaty)

	}

	var sumAssuredResult models.SumsAssuredResult

	sumAssuredResult.GlaCoveredSumAssured = glaCovered
	sumAssuredResult.GlaRetainedSumAssured = glaRetained
	sumAssuredResult.GlaCededSumAssured = glaCeded

	sumAssuredResult.SglaCoveredSumAssured = sglaCovered
	sumAssuredResult.SglaRetainedSumAssured = sglaRetained
	sumAssuredResult.SglaCededSumAssured = sglaCeded

	sumAssuredResult.PtdCoveredSumAssured = ptdCovered
	sumAssuredResult.PtdRetainedSumAssured = ptdRetained
	sumAssuredResult.PtdCededSumAssured = ptdCeded

	sumAssuredResult.CiCoveredSumAssured = ciCovered
	sumAssuredResult.CiRetainedSumAssured = ciRetained
	sumAssuredResult.CiCededSumAssured = ciCeded

	sumAssuredResult.PhiMonthlyBenefit = phiCovered
	sumAssuredResult.PhiRetainedMonthlyBenefit = phiRetained
	sumAssuredResult.PhiCededMonthlyBenefit = phiCeded

	sumAssuredResult.TtdMonthlyBenefit = ttdCovered
	sumAssuredResult.TtdRetainedMonthlyBenefit = ttdRetained
	sumAssuredResult.TtdCededMonthlyBenefit = ttdCeded

	sumAssuredResult.MmFuneralSumAssured = mmFuneralSACovered
	sumAssuredResult.MmRetainedFuneralSumAssured = mmFuneralRetained
	sumAssuredResult.MmCededFuneralSumAssured = mmFuneralCeded

	sumAssuredResult.SpFuneralSumAssured = spFuneralSACovered
	sumAssuredResult.SpRetainedFuneralSumAssured = spFuneralRetained
	sumAssuredResult.SpCededFuneralSumAssured = spFuneralCeded

	sumAssuredResult.ChFuneralSumAssured = chFuneralSACovered
	sumAssuredResult.ChRetainedFuneralSumAssured = chFuneralRetained
	sumAssuredResult.ChCededFuneralSumAssured = chFuneralCeded

	sumAssuredResult.DepFuneralSumAssured = depFuneralSACovered
	sumAssuredResult.DepRetainedFuneralSumAssured = depFuneralRetained
	sumAssuredResult.DepCededFuneralSumAssured = depFuneralCeded

	sumAssuredResult.ParFuneralSumAssured = parFuneralSACovered
	sumAssuredResult.ParRetainedFuneralSumAssured = parFuneralRetained
	sumAssuredResult.ParCededFuneralSumAssured = parFuneralCeded

	return sumAssuredResult
}

func SplitLumpSumSA(rawSA, freeCoverLimit float64, reinsuranceTreaty models.ReinsuranceTreaty) (covered, retained, ceded float64) {
	covered = math.Min(rawSA, freeCoverLimit)
	retained = 0
	ceded = 0

	var glalevel1value, glalevel2value, glalevel3value float64

	if reinsuranceTreaty.TreatyCode != "" {
		glalevel1value = math.Max(math.Min(covered, reinsuranceTreaty.Level1Upperbound-reinsuranceTreaty.Level1Lowerbound), 0)
		glalevel2value = math.Max(math.Min(covered-glalevel1value, reinsuranceTreaty.Level2Upperbound-reinsuranceTreaty.Level2Lowerbound), 0)
		glalevel3value = math.Max(math.Min(covered-glalevel1value-glalevel2value, reinsuranceTreaty.Level3Upperbound-reinsuranceTreaty.Level3Lowerbound), 0)
		ceded = glalevel1value*reinsuranceTreaty.Level1CededProportion/100 + glalevel2value*reinsuranceTreaty.Level2CededProportion/100 + glalevel3value*reinsuranceTreaty.Level3CededProportion/100

	}
	retained = math.Max(covered-ceded, 0)
	return covered, retained, ceded
}

func SplitIncome(rawSA, incomeLimit float64, reinsuranceTreaty models.ReinsuranceTreaty) (covered, retained, ceded float64) {
	covered = math.Min(rawSA, incomeLimit)
	retained = 0
	ceded = 0

	var glalevel1value, glalevel2value, glalevel3value float64

	if reinsuranceTreaty.TreatyCode != "" {
		glalevel1value = math.Max(math.Min(covered, reinsuranceTreaty.Level1Upperbound-reinsuranceTreaty.Level1Lowerbound), 0)
		glalevel2value = math.Max(math.Min(covered-glalevel1value, reinsuranceTreaty.Level2Upperbound-reinsuranceTreaty.Level2Lowerbound), 0)
		glalevel3value = math.Max(math.Min(covered-glalevel1value-glalevel2value, reinsuranceTreaty.Level3Upperbound-reinsuranceTreaty.Level3Lowerbound), 0)
		ceded = glalevel1value*reinsuranceTreaty.Level1CededProportion/100 + glalevel2value*reinsuranceTreaty.Level2CededProportion/100 + glalevel3value*reinsuranceTreaty.Level3CededProportion/100
	}
	retained = math.Max(covered-ceded, 0)
	return covered, retained, ceded
}

func SplitFuneralSA(rawSA float64, reinsuranceTreaty models.ReinsuranceTreaty) (covered, retained, ceded float64) {
	covered = rawSA
	retained = 0
	ceded = 0
	var glalevel1value, glalevel2value, glalevel3value float64

	if reinsuranceTreaty.TreatyCode != "" {
		glalevel1value = math.Max(math.Min(covered, reinsuranceTreaty.Level1Upperbound-reinsuranceTreaty.Level1Lowerbound), 0)
		glalevel2value = math.Max(math.Min(covered-glalevel1value, reinsuranceTreaty.Level2Upperbound-reinsuranceTreaty.Level2Lowerbound), 0)
		glalevel3value = math.Max(math.Min(covered-glalevel1value-glalevel2value, reinsuranceTreaty.Level3Upperbound-reinsuranceTreaty.Level3Lowerbound), 0)
		ceded = glalevel1value*reinsuranceTreaty.Level1CededProportion/100 + glalevel2value*reinsuranceTreaty.Level2CededProportion/100 + glalevel3value*reinsuranceTreaty.Level3CededProportion/100
	}
	retained = math.Max(covered-ceded, 0)
	return covered, retained, ceded
}

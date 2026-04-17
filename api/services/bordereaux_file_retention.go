package services

import (
	appLog "api/log"
	"api/config"
	"api/models"
	"os"
	"time"
)

// SweepExpiredBordereauxFiles removes on-disk artefacts whose retention window
// has passed. DB rows are preserved (metadata + reconciled content stays);
// only the file is deleted and FilePath/FileName are cleared so the UI knows
// the artefact is gone. Three sweeps run in sequence:
//
//   - GeneratedBordereaux: any row older than retention (output files are
//     regenerable from the underlying DB data).
//   - BordereauxConfirmation: older than retention AND in a terminal status.
//   - EmployerSubmission: older than retention AND in a terminal status.
//
// Terminal statuses are intentionally narrow for user-uploaded files so we
// don't destroy evidence that's still under active review.
func SweepExpiredBordereauxFiles() {
	days := config.BordereauxFileRetentionDays
	if days <= 0 {
		return
	}
	cutoff := time.Now().AddDate(0, 0, -days)

	deletedReports, freedReports := sweepGeneratedBordereauxFiles(cutoff)
	deletedConfirmations, freedConfirmations := sweepConfirmationFiles(cutoff)
	deletedSubmissions, freedSubmissions := sweepSubmissionFiles(cutoff)

	total := deletedReports + deletedConfirmations + deletedSubmissions
	if total == 0 {
		return
	}
	appLog.WithFields(map[string]interface{}{
		"retention_days":         days,
		"cutoff":                 cutoff.Format(time.RFC3339),
		"deleted_reports":        deletedReports,
		"deleted_confirmations":  deletedConfirmations,
		"deleted_submissions":    deletedSubmissions,
		"bytes_freed":            freedReports + freedConfirmations + freedSubmissions,
	}).Info("Bordereaux file retention sweep complete")
}

// StartBordereauxFileRetentionSweeper runs the retention sweep immediately at
// boot and then every 24 hours. Mirrors the ticker idiom used by the other
// bordereaux sweepers.
func StartBordereauxFileRetentionSweeper() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			SweepExpiredBordereauxFiles()
			<-ticker.C
		}
	}()
}

func sweepGeneratedBordereauxFiles(cutoff time.Time) (deleted int, bytesFreed int64) {
	var rows []models.GeneratedBordereaux
	if err := DB.Where("created_at < ? AND file_path != ''", cutoff).Find(&rows).Error; err != nil {
		appLog.Warn("Retention sweep failed to load generated bordereaux: " + err.Error())
		return
	}
	for _, r := range rows {
		freed := removeIfPresent(r.FilePath)
		if err := DB.Model(&models.GeneratedBordereaux{}).
			Where("id = ?", r.ID).
			Updates(map[string]interface{}{"file_path": "", "file_name": ""}).Error; err != nil {
			appLog.WithFields(map[string]interface{}{
				"generated_id": r.GeneratedID,
				"error":        err.Error(),
			}).Warn("Failed to clear file path on generated bordereaux after retention delete")
			continue
		}
		deleted++
		bytesFreed += freed
	}
	return
}

func sweepConfirmationFiles(cutoff time.Time) (deleted int, bytesFreed int64) {
	var rows []models.BordereauxConfirmation
	terminal := []string{"matched", "discrepancy", "error"}
	if err := DB.Where("imported_at < ? AND file_path != '' AND status IN ?", cutoff, terminal).
		Find(&rows).Error; err != nil {
		appLog.Warn("Retention sweep failed to load confirmations: " + err.Error())
		return
	}
	for _, r := range rows {
		freed := removeIfPresent(r.FilePath)
		if err := DB.Model(&models.BordereauxConfirmation{}).
			Where("id = ?", r.ID).
			Updates(map[string]interface{}{"file_path": "", "file_name": ""}).Error; err != nil {
			appLog.WithFields(map[string]interface{}{
				"confirmation_id": r.ID,
				"error":           err.Error(),
			}).Warn("Failed to clear file path on confirmation after retention delete")
			continue
		}
		deleted++
		bytesFreed += freed
	}
	return
}

func sweepSubmissionFiles(cutoff time.Time) (deleted int, bytesFreed int64) {
	var rows []models.EmployerSubmission
	terminal := []string{"accepted", "rejected"}
	if err := DB.Where("created_at < ? AND file_path != '' AND status IN ?", cutoff, terminal).
		Find(&rows).Error; err != nil {
		appLog.Warn("Retention sweep failed to load submissions: " + err.Error())
		return
	}
	for _, r := range rows {
		freed := removeIfPresent(r.FilePath)
		if err := DB.Model(&models.EmployerSubmission{}).
			Where("id = ?", r.ID).
			Updates(map[string]interface{}{"file_path": "", "file_name": ""}).Error; err != nil {
			appLog.WithFields(map[string]interface{}{
				"submission_id": r.ID,
				"error":         err.Error(),
			}).Warn("Failed to clear file path on submission after retention delete")
			continue
		}
		deleted++
		bytesFreed += freed
	}
	return
}

// removeIfPresent deletes a file if it still exists on disk and returns the
// byte size it freed. Missing files are not an error — the DB column can still
// be cleared (the file may have been removed out-of-band).
func removeIfPresent(path string) int64 {
	if path == "" {
		return 0
	}
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	size := info.Size()
	if err := os.Remove(path); err != nil {
		appLog.WithFields(map[string]interface{}{
			"path":  path,
			"error": err.Error(),
		}).Warn("Retention sweep failed to remove file")
		return 0
	}
	return size
}

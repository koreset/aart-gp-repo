package services

import (
	"api/models"
	"fmt"
	"time"
)

// GetBordereauxDeadlines returns deadlines filtered by optional scheme, month, year, status.
func GetBordereauxDeadlines(schemeID, month, year int, status string) ([]models.BordereauxDeadline, error) {
	var deadlines []models.BordereauxDeadline
	q := DB.Model(&models.BordereauxDeadline{})
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
	if err := q.Order("due_date ASC").Find(&deadlines).Error; err != nil {
		return nil, err
	}
	// Refresh overdue status in-memory for any pending deadlines past grace expiry
	today := time.Now().Format("2006-01-02")
	for i := range deadlines {
		if deadlines[i].Status == "pending" && deadlines[i].DueDate != "" {
			graceExpiry := graceExpiryDate(deadlines[i].DueDate, deadlines[i].GracePeriodDays)
			if graceExpiry < today {
				deadlines[i].Status = "overdue"
				DB.Model(&deadlines[i]).Update("status", "overdue")
			}
		}
	}
	return deadlines, nil
}

// CreateBordereauxDeadline creates a single deadline.
func CreateBordereauxDeadline(req models.CreateDeadlineRequest, user models.AppUser) (models.BordereauxDeadline, error) {
	var scheme models.GroupScheme
	schemeName := ""
	if err := DB.First(&scheme, req.SchemeID).Error; err == nil {
		schemeName = scheme.Name
	}

	dlType := req.Type
	if dlType == "" {
		dlType = "member_submission"
	}

	dl := models.BordereauxDeadline{
		SchemeID:        req.SchemeID,
		SchemeName:      schemeName,
		Month:           req.Month,
		Year:            req.Year,
		Type:            dlType,
		DueDate:         req.DueDate,
		GracePeriodDays: req.GracePeriodDays,
		Status:          "pending",
		CreatedBy:       user.UserName,
	}
	if err := DB.Create(&dl).Error; err != nil {
		return dl, fmt.Errorf("failed to create deadline: %w", err)
	}
	return dl, nil
}

// GenerateDeadlinesForAllSchemes auto-creates deadlines for every in-force scheme
// for the given month/year. Schemes that already have a deadline for the period are skipped.
func GenerateDeadlinesForAllSchemes(req models.GenerateDeadlinesRequest, user models.AppUser) (models.GenerateDeadlinesResult, error) {
	dueDayOfMonth := req.DueDayOfMonth
	if dueDayOfMonth <= 0 || dueDayOfMonth > 28 {
		dueDayOfMonth = 15
	}

	dueDate := fmt.Sprintf("%04d-%02d-%02d", req.Year, req.Month, dueDayOfMonth)

	var schemes []models.GroupScheme
	if err := DB.Where("in_force = ?", true).Find(&schemes).Error; err != nil {
		return models.GenerateDeadlinesResult{}, fmt.Errorf("failed to load in-force schemes: %w", err)
	}

	result := models.GenerateDeadlinesResult{Total: len(schemes)}

	for _, scheme := range schemes {
		// Check if deadline already exists for this scheme/month/year
		var existing models.BordereauxDeadline
		err := DB.Where("scheme_id = ? AND month = ? AND year = ? AND type = ?",
			scheme.ID, req.Month, req.Year, "member_submission").First(&existing).Error
		if err == nil {
			result.Skipped++
			continue
		}

		dl := models.BordereauxDeadline{
			SchemeID:        scheme.ID,
			SchemeName:      scheme.Name,
			Month:           req.Month,
			Year:            req.Year,
			Type:            "member_submission",
			DueDate:         dueDate,
			GracePeriodDays: req.GracePeriodDays,
			Status:          "pending",
			CreatedBy:       user.UserName,
		}
		if err := DB.Create(&dl).Error; err == nil {
			result.Created++
		}
	}

	return result, nil
}

// UpdateDeadlineStatus patches the status of a deadline (waive or reopen to pending).
func UpdateDeadlineStatus(id int, req models.UpdateDeadlineStatusRequest, user models.AppUser) (models.BordereauxDeadline, error) {
	var dl models.BordereauxDeadline
	if err := DB.First(&dl, id).Error; err != nil {
		return dl, fmt.Errorf("deadline not found: %w", err)
	}

	switch req.Status {
	case "waived":
		now := time.Now()
		dl.Status = "waived"
		dl.WaivedBy = user.UserName
		dl.WaivedAt = &now
		dl.WaiverReason = req.WaiverReason
	case "pending":
		dl.Status = "pending"
		dl.WaivedBy = ""
		dl.WaivedAt = nil
		dl.WaiverReason = ""
	default:
		return dl, fmt.Errorf("invalid status '%s': must be 'waived' or 'pending'", req.Status)
	}

	if err := DB.Save(&dl).Error; err != nil {
		return dl, fmt.Errorf("failed to update deadline: %w", err)
	}
	return dl, nil
}

// LinkDeadlineToSubmission marks the matching deadline as received and links it to the submission.
// Called automatically when an employer submission file is uploaded successfully.
func LinkDeadlineToSubmission(schemeID, month, year, submissionID int) {
	var dl models.BordereauxDeadline
	err := DB.Where("scheme_id = ? AND month = ? AND year = ? AND type = ?",
		schemeID, month, year, "member_submission").First(&dl).Error
	if err != nil {
		return // no deadline configured — not an error
	}
	dl.Status = "received"
	dl.LinkedSubmissionID = &submissionID
	DB.Save(&dl)
}

// GetDeadlineStats returns aggregate counts across all deadlines, refreshing overdue statuses first.
func GetDeadlineStats() (models.DeadlineStats, error) {
	// Refresh overdue: set pending deadlines past grace expiry to overdue
	today := time.Now().Format("2006-01-02")
	DB.Model(&models.BordereauxDeadline{}).
		Where("status = 'pending' AND due_date != '' AND date(due_date, '+' || grace_period_days || ' days') < ?", today).
		Update("status", "overdue")

	var stats models.DeadlineStats
	type countRow struct {
		Status string
		Count  int
	}
	var rows []countRow
	DB.Model(&models.BordereauxDeadline{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&rows)
	for _, r := range rows {
		switch r.Status {
		case "pending":
			stats.PendingCount = r.Count
		case "overdue":
			stats.OverdueCount = r.Count
		case "received":
			stats.ReceivedCount = r.Count
		case "waived":
			stats.WaivedCount = r.Count
		}
	}
	return stats, nil
}

// graceExpiryDate returns the ISO date string for due_date + graceDays.
func graceExpiryDate(dueDate string, graceDays int) string {
	t, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return dueDate
	}
	return t.AddDate(0, 0, graceDays).Format("2006-01-02")
}

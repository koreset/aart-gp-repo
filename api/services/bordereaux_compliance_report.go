package services

import (
	"api/models"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// GenerateComplianceReport assembles an xlsx workbook of bordereaux-compliance
// signals an auditor would ask for: the summary numbers, open discrepancies,
// overdue deadlines, current escalations, and pending large-claim notices.
// Returns the file plus a suggested download filename so the controller can
// stream it with a sensible Content-Disposition.
func GenerateComplianceReport(from, to time.Time) (*excelize.File, string, error) {
	if from.IsZero() {
		from = time.Now().AddDate(0, 0, -30)
	}
	if to.IsZero() {
		to = time.Now()
	}
	f := excelize.NewFile()

	if err := writeComplianceSummarySheet(f, from, to); err != nil {
		return nil, "", err
	}
	if err := writeComplianceDiscrepancySheet(f, from, to); err != nil {
		return nil, "", err
	}
	if err := writeComplianceDeadlinesSheet(f); err != nil {
		return nil, "", err
	}
	if err := writeComplianceEscalationsSheet(f); err != nil {
		return nil, "", err
	}
	if err := writeComplianceLargeClaimsSheet(f, from, to); err != nil {
		return nil, "", err
	}

	// Remove the default "Sheet1" added by excelize.NewFile — we don't use it.
	if idx, _ := f.GetSheetIndex("Sheet1"); idx >= 0 {
		_ = f.DeleteSheet("Sheet1")
	}

	filename := fmt.Sprintf("bordereaux_compliance_%s_%s.xlsx",
		from.Format("20060102"), to.Format("20060102"))
	return f, filename, nil
}

func writeComplianceSummarySheet(f *excelize.File, from, to time.Time) error {
	name := "Summary"
	idx, err := f.NewSheet(name)
	if err != nil {
		return err
	}
	f.SetActiveSheet(idx)

	_ = f.SetCellValue(name, "A1", "Bordereaux Compliance Report")
	_ = f.SetCellValue(name, "A2", fmt.Sprintf("Period: %s to %s",
		from.Format("2006-01-02"), to.Format("2006-01-02")))
	_ = f.SetCellValue(name, "A3", fmt.Sprintf("Generated: %s",
		time.Now().Format("2006-01-02 15:04")))

	var bordCount, confCount, matchedCount, discrepancyCount int64
	DB.Model(&models.GeneratedBordereaux{}).
		Where("created_at BETWEEN ? AND ?", from, to).
		Count(&bordCount)
	DB.Model(&models.BordereauxConfirmation{}).
		Where("imported_at BETWEEN ? AND ?", from, to).
		Count(&confCount)
	DB.Model(&models.BordereauxReconciliationResult{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", from, to, "matched").
		Count(&matchedCount)
	DB.Model(&models.BordereauxReconciliationResult{}).
		Where("created_at BETWEEN ? AND ? AND status IN ?", from, to,
			[]string{"discrepancy", "missing", "extra"}).
		Count(&discrepancyCount)

	var overdueDeadlines int64
	DB.Model(&models.BordereauxDeadline{}).
		Where("status = ?", "overdue").
		Count(&overdueDeadlines)

	var openEscalations, overdueEscalations int64
	DB.Model(&models.BordereauxReconciliationResult{}).
		Where("status = ?", "escalated").
		Count(&openEscalations)
	DB.Model(&models.BordereauxReconciliationResult{}).
		Where("status = ? AND due_date IS NOT NULL AND due_date < ?", "escalated", time.Now()).
		Count(&overdueEscalations)

	var pendingNotices, lateNotices int64
	DB.Model(&models.LargeClaimNotice{}).
		Where("status IN ?", []string{"pending", "sent", "queried"}).
		Count(&pendingNotices)
	DB.Model(&models.LargeClaimNotice{}).
		Where("late_flag = ?", true).
		Count(&lateNotices)

	rows := [][]interface{}{
		{"Metric", "Value"},
		{"Bordereaux generated", bordCount},
		{"Confirmations imported", confCount},
		{"Reconciled lines: matched", matchedCount},
		{"Reconciled lines: open discrepancies", discrepancyCount},
		{"Deadlines overdue (all time)", overdueDeadlines},
		{"Escalations open", openEscalations},
		{"Escalations past SLA", overdueEscalations},
		{"Large-claim notices awaiting response", pendingNotices},
		{"Large-claim notices flagged late", lateNotices},
	}
	for i, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+5)
		_ = f.SetSheetRow(name, cell, &row)
	}
	_ = f.SetColWidth(name, "A", "A", 48)
	_ = f.SetColWidth(name, "B", "B", 18)
	return nil
}

func writeComplianceDiscrepancySheet(f *excelize.File, from, to time.Time) error {
	name := "Open Discrepancies"
	if _, err := f.NewSheet(name); err != nil {
		return err
	}
	headers := []interface{}{
		"ID", "Confirmation ID", "Bordereaux ID", "Record ID", "Member Name",
		"Field", "Expected", "Actual", "Variance", "Status", "Created At",
	}
	_ = f.SetSheetRow(name, "A1", &headers)

	var rows []models.BordereauxReconciliationResult
	if err := DB.Where("created_at BETWEEN ? AND ? AND status IN ? AND is_resolved = ?",
		from, to, []string{"discrepancy", "missing", "extra"}, false).
		Order("created_at DESC").
		Limit(5000).
		Find(&rows).Error; err != nil {
		return err
	}
	for i, r := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		row := []interface{}{
			r.ID, r.BordereauxConfirmationID, r.GeneratedBordereauxID,
			r.RecordID, r.MemberName, r.Field, r.ExpectedValue, r.ActualValue,
			r.Variance, r.Status, r.CreatedAt.Format("2006-01-02 15:04"),
		}
		_ = f.SetSheetRow(name, cell, &row)
	}
	return nil
}

func writeComplianceDeadlinesSheet(f *excelize.File) error {
	name := "Overdue Deadlines"
	if _, err := f.NewSheet(name); err != nil {
		return err
	}
	headers := []interface{}{
		"ID", "Scheme", "Type", "Due Date", "Grace Days", "Status",
		"Linked Submission", "Created At",
	}
	_ = f.SetSheetRow(name, "A1", &headers)

	var rows []models.BordereauxDeadline
	if err := DB.Where("status = ?", "overdue").
		Order("due_date ASC").
		Limit(5000).
		Find(&rows).Error; err != nil {
		return err
	}
	for i, r := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		submission := ""
		if r.LinkedSubmissionID != nil {
			submission = fmt.Sprintf("%d", *r.LinkedSubmissionID)
		}
		row := []interface{}{
			r.ID, r.SchemeName, r.Type, r.DueDate, r.GracePeriodDays,
			r.Status, submission, r.CreatedAt.Format("2006-01-02 15:04"),
		}
		_ = f.SetSheetRow(name, cell, &row)
	}
	return nil
}

func writeComplianceEscalationsSheet(f *excelize.File) error {
	name := "Escalations"
	if _, err := f.NewSheet(name); err != nil {
		return err
	}
	headers := []interface{}{
		"ID", "Record ID", "Member Name", "Assigned To", "Priority",
		"Escalated By", "Escalated At", "Due Date", "Comments",
	}
	_ = f.SetSheetRow(name, "A1", &headers)

	var rows []models.BordereauxReconciliationResult
	if err := DB.Where("status = ?", "escalated").
		Order("due_date ASC").
		Limit(5000).
		Find(&rows).Error; err != nil {
		return err
	}
	for i, r := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		due := ""
		if r.DueDate != nil {
			due = r.DueDate.Format("2006-01-02 15:04")
		}
		escalated := ""
		if r.EscalatedAt != nil {
			escalated = r.EscalatedAt.Format("2006-01-02 15:04")
		}
		row := []interface{}{
			r.ID, r.RecordID, r.MemberName, r.AssignedTo, r.Priority,
			r.EscalatedBy, escalated, due, r.Comments,
		}
		_ = f.SetSheetRow(name, cell, &row)
	}
	return nil
}

func writeComplianceLargeClaimsSheet(f *excelize.File, from, to time.Time) error {
	name := "Large Claim Notices"
	if _, err := f.NewSheet(name); err != nil {
		return err
	}
	headers := []interface{}{
		"ID", "Treaty", "Claim Number", "Scheme", "Reinsurer",
		"Gross Claim", "Estimated Ceded", "Status", "Response Status",
		"Notified Date", "Due Date", "Late Flag",
	}
	_ = f.SetSheetRow(name, "A1", &headers)

	var rows []models.LargeClaimNotice
	if err := DB.Where("created_at BETWEEN ? AND ?", from, to).
		Order("created_at DESC").
		Limit(5000).
		Find(&rows).Error; err != nil {
		return err
	}
	for i, r := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		row := []interface{}{
			r.ID, r.TreatyNumber, r.ClaimNumber, r.SchemeName, r.ReinsurerName,
			r.GrossClaimAmount, r.EstimatedCededAmount, r.Status, r.ResponseStatus,
			r.NotifiedDate, r.DueDate, r.LateFlag,
		}
		_ = f.SetSheetRow(name, cell, &row)
	}
	return nil
}

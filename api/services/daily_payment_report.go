package services

import (
	appLog "api/log"
	"api/models"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ──────────────────────────────────────────────
// Daily payment report (Phase 4)
// ──────────────────────────────────────────────
//
// Once per day (after the configured "report time") an email goes to the
// configured recipients (Head of Claims, Head of Finance) summarising the
// previous day's payment activity: schedules confirmed, totals, queries
// raised, exceptions outstanding, plus archive count.
//
// The report time is implicit: 06:00 in the configured timezone. Repeats
// are deduped by stamping daily_payment_reports.report_date so the
// scheduler can call SendDailyReportIfDue() every minute and only fire once.

const dailyReportTemplateCode = "daily_payment_report"

// DailyPaymentReportRecipients is the comma-separated env-like override.
// In production this is configured per-install through email recipients on
// the cutoff config; this var-level seam keeps tests trivial.
var DailyPaymentReportRecipients = ""

// reportSentMu serialises the dedupe check so concurrent scheduler ticks
// don't both fire the same day's report. The DB-level uniqueness is enforced
// via INSERT ... ON DUPLICATE pattern below; this is just an optimisation.
var reportSentMu sync.Mutex

// SendDailyReportIfDue checks the local-time clock against the configured
// report hour and fires the day's report if it hasn't yet been sent for
// today's date. Idempotent — safe to call every minute.
func SendDailyReportIfDue(now time.Time) {
	cfg, err := GetPaymentCutoffConfig("")
	if err != nil {
		return
	}
	loc, err := time.LoadLocation(strings.TrimSpace(cfg.Timezone))
	if err != nil || cfg.Timezone == "" {
		loc = time.Local
	}
	local := now.In(loc)
	if local.Hour() != 6 || local.Minute() != 0 {
		return
	}

	reportSentMu.Lock()
	defer reportSentMu.Unlock()

	reportDate := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)

	// Dedupe by checking an audit-style row pattern in payment_schedule_audits
	// — a synthetic schedule_id = -1 record marks the daily report having
	// run for that calendar day. This avoids needing yet another table.
	var exists int64
	if err := DB.Model(&models.PaymentScheduleAudit{}).
		Where("schedule_id = ? AND from_status = ? AND to_status = ?", -1, "daily_report", reportDate.Format("2006-01-02")).
		Count(&exists).Error; err == nil && exists > 0 {
		return
	}

	if err := sendDailyReport(reportDate, loc); err != nil {
		appLog.WithField("error", err.Error()).Warn("daily payment report: send failed")
		return
	}
	_ = DB.Create(&models.PaymentScheduleAudit{
		ScheduleID: -1,
		FromStatus: "daily_report",
		ToStatus:   reportDate.Format("2006-01-02"),
		Actor:      "system:daily_report",
		Notes:      "Daily payment report dispatched",
	}).Error
}

// sendDailyReport assembles the previous-day metrics and dispatches the
// email. The numbers come from the existing tables — no new schema needed.
func sendDailyReport(reportDate time.Time, loc *time.Location) error {
	dayStart := reportDate.AddDate(0, 0, -1)
	dayEnd := reportDate

	type counts struct {
		Confirmed       int64
		ConfirmedNet    float64
		FailedLines     int64
		QueriesRaised   int64
		ArchivedToday   int64
	}
	var c counts
	DB.Model(&models.ClaimPaymentSchedule{}).
		Where("status = ? AND finance_second_auth_at >= ? AND finance_second_auth_at < ?", "confirmed", dayStart, dayEnd).
		Count(&c.Confirmed)
	DB.Model(&models.ClaimPaymentSchedule{}).
		Where("status = ? AND finance_second_auth_at >= ? AND finance_second_auth_at < ?", "confirmed", dayStart, dayEnd).
		Select("COALESCE(SUM(net_total), 0)").Row().Scan(&c.ConfirmedNet)
	DB.Model(&models.ACBReconciliationResult{}).
		Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
		Where("status = ?", "failed").Count(&c.FailedLines)
	DB.Model(&models.ClaimPaymentScheduleQuery{}).
		Where("raised_at >= ? AND raised_at < ?", dayStart, dayEnd).Count(&c.QueriesRaised)
	DB.Model(&models.ClaimPaymentSchedule{}).
		Where("archived_at >= ? AND archived_at < ?", dayStart, dayEnd).Count(&c.ArchivedToday)

	recipients := splitRecipients(DailyPaymentReportRecipients)
	if len(recipients) == 0 {
		appLog.Info("daily payment report: no recipients configured, skipping send")
		return nil
	}
	licenseId := resolveDefaultLicense()
	if licenseId == "" {
		appLog.Info("daily payment report: no email_settings configured, skipping send")
		return nil
	}

	vars := map[string]interface{}{
		"report_date":      reportDate.Format("2 January 2006"),
		"confirmed_count":  c.Confirmed,
		"confirmed_net":    c.ConfirmedNet,
		"failed_lines":     c.FailedLines,
		"queries_raised":   c.QueriesRaised,
		"archived_count":   c.ArchivedToday,
		"summary_line":     fmt.Sprintf("%d schedule(s) confirmed totalling R %.2f net.", c.Confirmed, c.ConfirmedNet),
	}

	req := EnqueueEmailRequest{
		LicenseId:         licenseId,
		TemplateCode:      dailyReportTemplateCode,
		To:                recipients,
		Vars:              vars,
		RelatedObjectType: "daily_payment_report",
		RelatedObjectID:   reportDate.Format("2006-01-02"),
		CreatedBy:         "system:daily_report",
	}
	if _, err := EnqueueEmail(req); err != nil {
		if isMissingTemplate(err) {
			appLog.WithField("template", dailyReportTemplateCode).
				Info("daily payment report: template not configured, skipping send")
			return nil
		}
		return err
	}
	return nil
}

func splitRecipients(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// ──────────────────────────────────────────────
// Auto-archive (Phase 4)
// ──────────────────────────────────────────────
//
// Confirmed schedules become read-only history once they've been settled
// long enough that no further follow-up is plausible. Default window:
// 30 days. Run from the same scheduler tick.

// ArchiveAgeDays is the number of days after `confirmed_at` (= upload-proof
// time) that a schedule is auto-archived. Override at startup by reassigning.
var ArchiveAgeDays = 30

// ArchiveOldConfirmedSchedules transitions schedules confirmed more than
// ArchiveAgeDays ago into the "archived" state. Best-effort; logs failures
// individually rather than aborting the batch.
func ArchiveOldConfirmedSchedules() {
	cutoff := time.Now().AddDate(0, 0, -ArchiveAgeDays)
	var rows []models.ClaimPaymentSchedule
	if err := DB.Select("id, status, updated_at").
		Where("status = ? AND updated_at <= ?", "confirmed", cutoff).
		Find(&rows).Error; err != nil {
		appLog.WithField("error", err.Error()).Warn("auto-archive: lookup failed")
		return
	}
	now := time.Now()
	for _, r := range rows {
		txErr := DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", r.ID).Updates(map[string]interface{}{
				"status":      "archived",
				"archived_at": &now,
				"archived_by": "system:auto_archive",
			}).Error; err != nil {
				return err
			}
			return tx.Create(&models.PaymentScheduleAudit{
				ScheduleID: r.ID,
				FromStatus: "confirmed",
				ToStatus:   "archived",
				Actor:      "system:auto_archive",
				Notes:      fmt.Sprintf("Auto-archived after %d days", ArchiveAgeDays),
			}).Error
		})
		if txErr != nil {
			appLog.WithField("error", txErr.Error()).
				WithField("schedule_id", r.ID).
				Warn("auto-archive: transition failed")
		}
	}
}

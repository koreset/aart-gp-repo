package services

import (
	appLog "api/log"
	"api/models"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ──────────────────────────────────────────────
// Payment cut-off scheduler
// ──────────────────────────────────────────────
//
// The scheduler ticks every minute, reads the active PaymentCutoffConfig for
// each license, and runs a cut-off when "now" matches a configured HH:MM
// (within a 60-second window) provided a dedupe row in payment_cutoff_runs
// doesn't already exist for that license + scheduled_at + trigger="auto".
//
// A cut-off run = call CreatePaymentSchedule with the IDs of every claim
// currently in "approved" status. If no claims are approved, a no-claims
// PaymentCutoffRun row is still recorded so dashboards can see the cut-off
// fired but produced no schedule (this is normal between batches).

var (
	cutoffSchedulerMu      sync.Mutex
	cutoffSchedulerRunning bool
)

// systemUser is used as the actor for auto-generated cut-off runs. The
// resulting audit trail rows are distinguishable from human-triggered ones
// by the "system:cutoff_scheduler" username.
var systemUser = models.AppUser{UserName: "system:cutoff_scheduler"}

// StartCutoffScheduler boots the background scheduler. Idempotent —
// calling it multiple times is a no-op after the first start.
func StartCutoffScheduler() {
	cutoffSchedulerMu.Lock()
	defer cutoffSchedulerMu.Unlock()

	if cutoffSchedulerRunning {
		appLog.Info("Payment cut-off scheduler is already running")
		return
	}
	cutoffSchedulerRunning = true

	go func() {
		appLog.Info("Payment cut-off scheduler started")
		// Align the first tick to the top of the next minute so we don't drift.
		now := time.Now()
		nextMin := now.Truncate(time.Minute).Add(time.Minute)
		time.Sleep(time.Until(nextMin))

		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		processCutoffTick(time.Now())
		for t := range ticker.C {
			processCutoffTick(t)
		}
	}()
}

// processCutoffTick is the per-minute loop body. Safe to call manually in
// tests.
func processCutoffTick(now time.Time) {
	var configs []models.PaymentCutoffConfig
	if err := DB.Where("enabled = ?", true).Find(&configs).Error; err != nil {
		appLog.WithField("error", err.Error()).Warn("cut-off scheduler: failed to load configs")
		return
	}
	for _, cfg := range configs {
		dispatchCutoff(cfg, now)
	}

	// Phase 4 hooks. The scheduler is already the only minute-ticker in the
	// process, so daily-report dispatch and auto-archive piggyback here
	// rather than spinning a second goroutine. Both functions are self-
	// dedup'ing — safe to call every tick.
	SendDailyReportIfDue(now)
	if now.Minute() == 5 { // run once an hour, away from cut-off-times
		ArchiveOldConfirmedSchedules()
	}
}

// dispatchCutoff looks at the configured HH:MM times for this config and
// triggers a run for any that match the current minute (in the config's
// timezone) and aren't already recorded in payment_cutoff_runs.
func dispatchCutoff(cfg models.PaymentCutoffConfig, now time.Time) {
	loc, err := time.LoadLocation(strings.TrimSpace(cfg.Timezone))
	if err != nil || cfg.Timezone == "" {
		loc = time.Local
	}
	local := now.In(loc)

	for _, hhmm := range parseCutoffTimes(cfg.CutoffTimes) {
		scheduledAt := combineDateTime(local, hhmm)
		if scheduledAt.Minute() != local.Minute() || scheduledAt.Hour() != local.Hour() {
			continue
		}

		var existing int64
		if err := DB.Model(&models.PaymentCutoffRun{}).
			Where("license_id = ? AND scheduled_at = ? AND trigger_type = ?", cfg.LicenseId, scheduledAt, "auto").
			Count(&existing).Error; err != nil {
			appLog.WithField("error", err.Error()).Warn("cut-off scheduler: dedupe check failed")
			continue
		}
		if existing > 0 {
			continue
		}

		appLog.WithField("scheduled_at", scheduledAt.Format(time.RFC3339)).Info("cut-off scheduler: firing auto cut-off")
		if _, err := RunCutoff(cfg.LicenseId, scheduledAt, "auto", systemUser); err != nil {
			appLog.WithField("error", err.Error()).Warn("cut-off scheduler: run failed")
		}
	}
}

// parseCutoffTimes splits a CSV like "11:00,15:00" into time-of-day structs.
// Malformed entries are skipped silently — settings UI is responsible for
// validation.
type timeOfDay struct{ Hour, Minute int }

func parseCutoffTimes(raw string) []timeOfDay {
	out := []timeOfDay{}
	for _, part := range strings.Split(raw, ",") {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		var h, m int
		if _, err := fmt.Sscanf(p, "%d:%d", &h, &m); err != nil {
			continue
		}
		if h < 0 || h > 23 || m < 0 || m > 59 {
			continue
		}
		out = append(out, timeOfDay{Hour: h, Minute: m})
	}
	return out
}

func combineDateTime(day time.Time, t timeOfDay) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), t.Hour, t.Minute, 0, 0, day.Location())
}

// RunCutoff executes a payment-schedule cut-off run for the given license at
// the given logical scheduledAt. triggerType is "auto" (scheduler) or
// "manual" (user-initiated from UI / API).
//
// The PaymentCutoffRun row is always written: either with status="ok" + a
// new schedule, status="no_claims" when nothing was approved, or
// status="error" with the message. The unique index on
// (license_id, scheduled_at, trigger_type) means concurrent runs collide
// at insert time — the second caller sees a duplicate-key error and bails.
func RunCutoff(licenseId string, scheduledAt time.Time, triggerType string, user models.AppUser) (models.PaymentCutoffRun, error) {
	if triggerType != "auto" && triggerType != "manual" {
		return models.PaymentCutoffRun{}, fmt.Errorf("invalid trigger_type %q", triggerType)
	}

	// Claim the run slot up-front so concurrent triggers don't double-fire.
	run := models.PaymentCutoffRun{
		LicenseId:   licenseId,
		ScheduledAt: scheduledAt,
		TriggerType: triggerType,
		Status:      "running",
		TriggeredBy: user.UserName,
	}
	if err := DB.Create(&run).Error; err != nil {
		return run, fmt.Errorf("could not claim cut-off run slot: %w", err)
	}

	// Find all currently-approved claim IDs.
	var claimIDs []int
	if err := DB.Model(&models.GroupSchemeClaim{}).
		Where("LOWER(status) = ?", "approved").
		Pluck("id", &claimIDs).Error; err != nil {
		_ = updateRun(run.ID, "error", err.Error(), nil, 0, 0)
		return run, err
	}

	if len(claimIDs) == 0 {
		_ = updateRun(run.ID, "no_claims", "", nil, 0, 0)
		run.Status = "no_claims"
		return run, nil
	}

	description := fmt.Sprintf("Cut-off %s (%s)", scheduledAt.Format("2006-01-02 15:04"), triggerType)
	schedule, err := CreatePaymentSchedule(CreatePaymentScheduleRequest{
		ClaimIDs:    claimIDs,
		Description: description,
	}, user)
	if err != nil {
		_ = updateRun(run.ID, "error", err.Error(), nil, 0, 0)
		return run, err
	}

	_ = updateRun(run.ID, "ok", "", &schedule.ID, schedule.ClaimsCount, schedule.GrossTotal)
	run.Status = "ok"
	run.ScheduleID = &schedule.ID
	run.ClaimsCount = schedule.ClaimsCount
	run.TotalAmount = schedule.GrossTotal
	return run, nil
}

func updateRun(id int, status, errMsg string, scheduleID *int, claimsCount int, total float64) error {
	return DB.Model(&models.PaymentCutoffRun{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"error_message": errMsg,
		"schedule_id":   scheduleID,
		"claims_count":  claimsCount,
		"total_amount":  total,
	}).Error
}

// ──────────────────────────────────────────────
// Settings CRUD
// ──────────────────────────────────────────────

// GetPaymentCutoffConfig returns the config for a license, falling back to a
// zero-value (Enabled=false) record when none is configured yet. Callers
// should treat the returned ID==0 as "no row yet — Save will Create".
func GetPaymentCutoffConfig(licenseId string) (models.PaymentCutoffConfig, error) {
	var cfg models.PaymentCutoffConfig
	err := DB.Where("license_id = ?", licenseId).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.PaymentCutoffConfig{LicenseId: licenseId, Timezone: "Africa/Johannesburg"}, nil
	}
	return cfg, err
}

// SavePaymentCutoffConfig upserts the config. Validates that cutoff_times
// parses to at least one HH:MM if Enabled.
func SavePaymentCutoffConfig(licenseId string, patch models.PaymentCutoffConfig, user models.AppUser) (models.PaymentCutoffConfig, error) {
	patch.LicenseId = licenseId
	patch.UpdatedBy = user.UserName

	if patch.Enabled && len(parseCutoffTimes(patch.CutoffTimes)) == 0 {
		return patch, errors.New("enabled config must declare at least one valid HH:MM cut-off time")
	}
	if patch.DailyPaymentLimit < 0 {
		return patch, errors.New("daily_payment_limit cannot be negative (use 0 for 'no limit')")
	}

	var existing models.PaymentCutoffConfig
	err := DB.Where("license_id = ?", licenseId).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := DB.Create(&patch).Error; err != nil {
			return patch, err
		}
		return patch, nil
	}
	if err != nil {
		return patch, err
	}
	updates := map[string]interface{}{
		"enabled":             patch.Enabled,
		"cutoff_times":        patch.CutoffTimes,
		"daily_payment_limit": patch.DailyPaymentLimit,
		"timezone":            patch.Timezone,
		"updated_by":          user.UserName,
	}
	if err := DB.Model(&existing).Updates(updates).Error; err != nil {
		return patch, err
	}
	return GetPaymentCutoffConfig(licenseId)
}

// ListRecentCutoffRuns returns the most recent runs for a license, newest
// first. Limit defaults to 30.
func ListRecentCutoffRuns(licenseId string, limit int) ([]models.PaymentCutoffRun, error) {
	if limit <= 0 || limit > 200 {
		limit = 30
	}
	var rows []models.PaymentCutoffRun
	err := DB.Where("license_id = ?", licenseId).Order("created_at DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

// NextCutoff returns the next scheduled cut-off for a license (or zero time
// if none configured / disabled). Used by the schedules list UI to show
// "Next cut-off: today 15:00".
func NextCutoff(licenseId string, now time.Time) (time.Time, bool) {
	cfg, err := GetPaymentCutoffConfig(licenseId)
	if err != nil || !cfg.Enabled {
		return time.Time{}, false
	}
	loc, err := time.LoadLocation(strings.TrimSpace(cfg.Timezone))
	if err != nil || cfg.Timezone == "" {
		loc = time.Local
	}
	local := now.In(loc)
	times := parseCutoffTimes(cfg.CutoffTimes)
	if len(times) == 0 {
		return time.Time{}, false
	}

	var earliest time.Time
	found := false
	for _, t := range times {
		candidate := combineDateTime(local, t)
		if !candidate.After(local) {
			continue
		}
		if !found || candidate.Before(earliest) {
			earliest = candidate
			found = true
		}
	}
	if !found {
		// All today's slots passed — first slot of tomorrow.
		tomorrow := local.AddDate(0, 0, 1)
		earliest = combineDateTime(tomorrow, times[0])
		found = true
	}
	return earliest, found
}

// ──────────────────────────────────────────────
// Daily payment limit guard
// ──────────────────────────────────────────────

// CheckDailyPaymentLimit returns an error when authorising this schedule's
// NetTotal would push today's already-authorised total above the configured
// daily limit. limit=0 means "no limit" and always passes.
//
// "Today's already-authorised total" is the sum of NetTotal across schedules
// for which first authorisation happened today. The proposed schedule is
// expected to be in `finance_in_review` status (i.e. not yet first-authorised),
// so it's not double-counted.
func CheckDailyPaymentLimit(licenseId string, proposed models.ClaimPaymentSchedule) error {
	cfg, err := GetPaymentCutoffConfig(licenseId)
	if err != nil {
		return nil
	}
	if cfg.DailyPaymentLimit <= 0 {
		return nil
	}

	loc, err := time.LoadLocation(strings.TrimSpace(cfg.Timezone))
	if err != nil || cfg.Timezone == "" {
		loc = time.Local
	}
	startOfDay := time.Now().In(loc).Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var todayTotal float64
	err = DB.Model(&models.ClaimPaymentSchedule{}).
		Where("finance_first_auth_at >= ? AND finance_first_auth_at < ?", startOfDay, endOfDay).
		Where("status IN ?", []string{"finance_first_authorised", "finance_second_authorised", "submitted_to_bank", "confirmed"}).
		Select("COALESCE(SUM(net_total), 0)").
		Row().Scan(&todayTotal)
	if err != nil {
		return nil
	}

	if todayTotal+proposed.NetTotal > cfg.DailyPaymentLimit {
		return fmt.Errorf("daily payment limit of %.2f exceeded: already authorised %.2f today, this schedule adds %.2f", cfg.DailyPaymentLimit, todayTotal, proposed.NetTotal)
	}
	return nil
}

package services

import (
	"api/models"
	"api/utils"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// GenerateMonthlySchedule creates a premium schedule for a scheme for a specific month/year.
// It reads the ContributionConfig from DB, applies employer/employee splits, and calculates
// pro-rata premiums for mid-month joiners.
func GenerateMonthlySchedule(schemeID int, month, year int, user models.AppUser) (models.PremiumSchedule, error) {
	var scheme models.GroupScheme
	if err := DB.First(&scheme, schemeID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("scheme not found: %w", err)
	}

	// Prevent schedule creation before scheme commencement date
	if !scheme.CommencementDate.IsZero() {
		commYear := scheme.CommencementDate.Year()
		commMonth := int(scheme.CommencementDate.Month())
		if year < commYear || (year == commYear && month < commMonth) {
			return models.PremiumSchedule{}, fmt.Errorf("cannot generate a schedule before the scheme commencement date (%s)", scheme.CommencementDate.Format("Jan 2006"))
		}
	}

	// Check if schedule already exists
	//var existing models.PremiumSchedule
	//if err := DB.Where("scheme_id = ? AND month = ? AND year = ?", schemeID, month, year).First(&existing).Error; err == nil {
	//	return existing, errors.New("schedule already exists for this period")
	//}

	// read associated quote that is in effect
	var quote models.GroupPricingQuote
	if err := DB.Where("scheme_id = ? AND scheme_quote_status = ?", schemeID, "in_effect").Find(&quote).Error; err != nil {
		return models.PremiumSchedule{}, errors.New("the scheme is currently not linked to any quote in effect")
	}

	// read categories
	var categories []models.SchemeCategory
	if err := DB.Where("quote_id = ?", quote.ID).Find(&categories).Error; err != nil {
		return models.PremiumSchedule{}, errors.New("missing scheme categories")
	}

	// read member rating result summary
	var memberRatingResultSummaries []models.MemberRatingResultSummary
	if err := DB.Where("quote_id = ?", quote.ID).Find(&memberRatingResultSummaries).Error; err != nil {
		return models.PremiumSchedule{}, errors.New("missing rating result summary")
	}

	// read reinsurance treaties
	var treatyLinks []models.TreatySchemeLink
	if err := DB.Where("scheme_id = ?", schemeID).Find(&treatyLinks).Error; err != nil {
		errors.New("no associated treaties")
	}

	// after loading treatyLinks, extract TreatyIDs
	linkedTreatyIds := extractTreatyIDs(treatyLinks)

	// read reinsurance treaties
	var reinsuranceTreaties []models.ReinsuranceTreaty
	var reinsuranceTreaty models.ReinsuranceTreaty
	if len(linkedTreatyIds) == 0 {
		reinsuranceTreaty = models.ReinsuranceTreaty{}
	} else {
		err := DB.Where("id IN ?", linkedTreatyIds).Find(&reinsuranceTreaties).Error
		if err != nil {
			reinsuranceTreaty = models.ReinsuranceTreaty{}
		} else {
			reinsuranceTreaty = reinsuranceTreaties[0]
		}
	}

	// Load contribution config; fall back to 100% employer if not configured
	var cfg models.ContributionConfig
	employerPct := 100.0
	employeePct := 0.0
	if err := DB.Where("scheme_id = ?", schemeID).Order("updated_at DESC").First(&cfg).Error; err == nil {
		if cfg.ContributionType == "split" {
			employerPct = cfg.EmployerPercent
			employeePct = cfg.EmployeePercent
		}
	}

	// Determine days in the schedule month for pro-rata
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	firstOfNext := firstOfMonth.AddDate(0, 1, 0)
	daysInMonth := int(firstOfNext.Sub(firstOfMonth).Hours() / 24)

	// Fetch premium rate from member rating result summary

	// Fetch all active members for this scheme
	var members []models.GPricingMemberDataInForce
	if err := DB.Where("scheme_id = ?", schemeID).Find(&members).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("failed to fetch members: %w", err)
	}

	// Pre-fetch all MemberPremiumSchedule rows for this scheme in one query,
	// then build a map of member_name → latest record to avoid N per-member SELECTs.
	var allMPS []models.MemberPremiumSchedule
	if err := DB.Where("scheme_id = ?", schemeID).Order("creation_date DESC").Find(&allMPS).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("failed to fetch premium schedules: %w", err)
	}
	mpsMap := make(map[string]models.MemberPremiumSchedule, len(allMPS))
	for _, mps := range allMPS {
		if _, exists := mpsMap[mps.MemberName]; !exists {
			mpsMap[mps.MemberName] = mps
		}
	}

	schedule := models.PremiumSchedule{
		SchemeID:      scheme.ID,
		SchemeName:    scheme.Name,
		Month:         month,
		Year:          year,
		Status:        "draft",
		GeneratedDate: time.Now().Format("2006-01-02"),
		GeneratedBy:   user.UserName,
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&schedule).Error; err != nil {
			return err
		}

		var grossTotal float64
		rows := make([]models.ScheduleMemberRow, 0, len(members))

		for _, m := range members {
			mps, ok := mpsMap[m.MemberName]
			if !ok {
				continue
			}

			var daysActive = daysInMonth

			fullMonthPremium := mps.TotalAnnualPremiumPayable / 12.0
			actualPremium := fullMonthPremium
			isProRata := false
			var proRataDays *int

			// Pro-rata: member joined during this month
			entryYear := m.EntryDate.Year()
			entryMonth := int(m.EntryDate.Month())
			if entryYear == year && entryMonth == month && m.EntryDate.Day() > 1 {
				daysRemaining := daysInMonth - m.EntryDate.Day() + 1
				actualPremium = fullMonthPremium * float64(daysRemaining) / float64(daysInMonth)
				isProRata = true
				proRataDays = &daysRemaining
			}

			// Pro-rata: member exited during this month
			//if !m.ExitDate.IsZero() {
			//	exitYear := m.ExitDate.Year()
			//	exitMonth := int(m.ExitDate.Month())
			//	if exitYear == year && exitMonth == month {
			//		daysActive := m.ExitDate.Day()
			//		actualPremium = fullMonthPremium * float64(daysActive) / float64(daysInMonth)
			//		isProRata = true
			//		proRataDays = &daysActive
			//	}
			//}

			if m.ExitDate != nil {
				monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				monthEnd := monthStart.AddDate(0, 1, -1) // last day of month
				daysInMonth := monthEnd.Day()

				// Normalize dates
				var effectiveStart time.Time
				if m.EntryDate.Before(monthStart) {
					effectiveStart = monthStart
				}
				if m.EntryDate.After(monthStart) {
					effectiveStart = m.EntryDate
				}

				var exitDate time.Time
				if m.ExitDate != nil {
					exitDate = *m.ExitDate
				}

				// Step 1: determine active interval within the month
				startActive := monthStart
				if !effectiveStart.IsZero() && effectiveStart.After(monthStart) {
					startActive = effectiveStart
				}

				// If there is an exit date, limit to that; otherwise assume still active through month end
				endActive := monthEnd
				if !exitDate.IsZero() && exitDate.Before(monthEnd) {
					endActive = exitDate
				}

				// Step 2: compute overlap duration
				if startActive.After(endActive) {
					// no overlap; no days active
					actualPremium = 0
					isProRata = false
					proRataDays = nil
				} else {
					// overlap exists
					daysActive = int(endActive.Sub(startActive).Hours()/24 + 1) // inclusive days
					// guard: ensure daysActive is within 0..daysInMonth
					if daysActive < 0 {
						daysActive = 0
					}
					if int(daysActive) > daysInMonth {
						daysActive = daysInMonth
					}

					// Step 3: decide premium
					if int(daysActive) == daysInMonth {
						// full month
						actualPremium = fullMonthPremium
						isProRata = false
						proRataDays = nil
					} else {
						// pro-rated
						actualPremium = fullMonthPremium * float64(daysActive) / float64(daysInMonth)
						isProRata = true
						d := daysActive
						proRataDays = &d
					}
				}
			}

			category, err := MapMemberToCategory(m, categories)
			if err != nil {
				return nil
			}
			memberratingResultSummary, err := MapMemberToMemberRatingResultSummary(m, memberRatingResultSummaries)
			coveredSA := CoveredSumsAssured(m, category, quote, reinsuranceTreaty, memberratingResultSummary)
			premium := PremiumComputation(m, coveredSA, memberratingResultSummary)

			//prorata function
			//isProRata, proRataDays = IsProRata(m, year, month, daysInMonth)
			fullMonthPremium = premium.TotalOfficePremium / 12.0
			actualPremium = fullMonthPremium * float64(daysActive) / float64(daysInMonth)

			rows = append(rows, models.ScheduleMemberRow{
				ScheduleID:           schedule.ID,
				MemberID:             m.ID,
				MemberName:           m.MemberName,
				MemberIDNumber:       m.MemberIdNumber,
				AnnualSalary:         m.AnnualSalary,
				FullMonthPremium:     fullMonthPremium,
				ActualPremium:        actualPremium,
				ProRataDays:          proRataDays,
				IsProRata:            isProRata,
				EmployerContribution: actualPremium * employerPct / 100.0,
				EmployeeContribution: actualPremium * employeePct / 100.0,
			})

			grossTotal += actualPremium
		}

		// Batch insert all member rows in one shot instead of one INSERT per member.
		if len(rows) > 0 {
			if err := tx.CreateInBatches(rows, 500).Error; err != nil {
				return err
			}
		}

		schedule.MemberCount = len(members)
		schedule.GrossPremium = grossTotal
		schedule.NetPayable = grossTotal
		return tx.Save(&schedule).Error
	})

	return schedule, err
}

// GenerateAllSchedules generates a monthly premium schedule for every in-force scheme.
// Schemes that already have a schedule for the period are counted as skipped; all
// others are attempted independently so one failure does not block the rest.
func GenerateAllSchedules(month, year int, user models.AppUser) (models.BulkGenerateResult, error) {
	var schemes []models.GroupScheme
	if err := DB.Where("in_force = ?", true).Find(&schemes).Error; err != nil {
		return models.BulkGenerateResult{}, fmt.Errorf("failed to fetch schemes: %w", err)
	}

	result := models.BulkGenerateResult{
		Total:   len(schemes),
		Results: make([]models.BulkScheduleResult, 0, len(schemes)),
	}

	for _, s := range schemes {
		row := models.BulkScheduleResult{
			SchemeID:   s.ID,
			SchemeName: s.Name,
		}

		// Skip schemes whose commencement date is after the requested period
		if !s.CommencementDate.IsZero() {
			commYear := s.CommencementDate.Year()
			commMonth := int(s.CommencementDate.Month())
			if year < commYear || (year == commYear && month < commMonth) {
				row.Status = "skipped"
				row.Message = fmt.Sprintf("Period is before scheme commencement date (%s)", s.CommencementDate.Format("Jan 2006"))
				result.Skipped++
				result.Results = append(result.Results, row)
				continue
			}
		}

		// Check for existing schedule first to distinguish skip from failure
		var existing models.PremiumSchedule
		if err := DB.Where("scheme_id = ? AND month = ? AND year = ?", s.ID, month, year).
			First(&existing).Error; err == nil {
			row.Status = "skipped"
			row.Message = "Schedule already exists for this period"
			row.ScheduleID = existing.ID
			result.Skipped++
			result.Results = append(result.Results, row)
			continue
		}

		schedule, err := GenerateMonthlySchedule(s.ID, month, year, user)
		if err != nil {
			row.Status = "failed"
			row.Message = err.Error()
			result.Failed++
		} else {
			row.Status = "success"
			row.ScheduleID = schedule.ID
			result.Succeeded++
		}
		result.Results = append(result.Results, row)
	}

	return result, nil
}

// GenerateInvoice creates an invoice from a finalized premium schedule.
// dueDate is optional — when empty the due date defaults to the last day of
// the invoice period month (e.g. Jan 2026 → 2026-01-31).
func GenerateInvoice(scheduleID int, dueDate string, user models.AppUser) (models.Invoice, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.Invoice{}, err
	}

	if schedule.Status != "finalized" {
		return models.Invoice{}, errors.New("schedule must be finalized to generate invoice")
	}

	// Check if invoice already exists for this schedule
	var existing models.Invoice
	if err := DB.Where("schedule_id = ?", scheduleID).First(&existing).Error; err == nil {
		return existing, nil
	}

	// Determine due date: use provided value, or default to end of invoice period month
	resolvedDueDate := dueDate
	if resolvedDueDate == "" {
		// Last day of the invoice period month
		lastDay := time.Date(schedule.Year, time.Month(schedule.Month)+1, 0, 0, 0, 0, 0, time.UTC)
		resolvedDueDate = lastDay.Format("2006-01-02")
	} else {
		// Validate the provided date format
		if _, err := time.Parse("2006-01-02", resolvedDueDate); err != nil {
			return models.Invoice{}, errors.New("invalid due_date format, expected YYYY-MM-DD")
		}
	}

	invoice := models.Invoice{
		InvoiceNumber: fmt.Sprintf("INV-%d-%02d-%d", schedule.SchemeID, schedule.Month, schedule.Year),
		SchemeID:      schedule.SchemeID,
		SchemeName:    schedule.SchemeName,
		ScheduleID:    schedule.ID,
		Month:         schedule.Month,
		Year:          schedule.Year,
		IssueDate:     time.Now().Format("2006-01-02"),
		DueDate:       resolvedDueDate,
		GrossAmount:   schedule.GrossPremium,
		NetPayable:    schedule.NetPayable,
		Balance:       schedule.NetPayable,
		Status:        "sent",
	}

	if err := DB.Create(&invoice).Error; err != nil {
		return models.Invoice{}, err
	}

	// Transition schedule to "invoiced"
	schedule.Status = "invoiced"
	DB.Save(&schedule)

	// Create reconciliation item for the new invoice
	EnsureReconciliationItemForInvoice(invoice)

	return invoice, nil
}

// RecordPayment records a payment and attempts to match it with an invoice.
func RecordPayment(req models.RecordPaymentRequest, user models.AppUser) (models.Payment, error) {
	payment := models.Payment{
		SchemeID:      req.SchemeID,
		InvoiceID:     req.InvoiceID,
		PaymentDate:   req.PaymentDate,
		Method:        req.Method,
		Amount:        req.Amount,
		BankReference: req.BankReference,
		Notes:         req.Notes,
		Status:        "pending",
		RecordedBy:    user.UserName,
		RecordedAt:    time.Now().Format(time.RFC3339),
	}

	var scheme models.GroupScheme
	if err := DB.First(&scheme, req.SchemeID).Error; err == nil {
		payment.SchemeName = scheme.Name
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		if payment.InvoiceID != nil {
			var invoice models.Invoice
			if err := tx.First(&invoice, *payment.InvoiceID).Error; err == nil {
				invoice.PaidAmount += payment.Amount
				invoice.Balance = invoice.NetPayable - invoice.PaidAmount
				// Use AutoCloseTolerance so sub-rand rounding residues close
				// out as 'paid' instead of lingering as 'partial'.
				if invoice.Balance <= AutoCloseTolerance {
					invoice.Status = "paid"
					invoice.Balance = 0
				} else {
					invoice.Status = "partial"
				}
				if err := tx.Save(&invoice).Error; err != nil {
					return err
				}
				payment.Status = "matched"
				payment.InvoiceNumber = invoice.InvoiceNumber
				if err := tx.Save(&payment).Error; err != nil {
					return err
				}

				// Sync v2 reconciliation ledger: ensure reconciliation items
				// exist for both sides, record the PaymentAllocation, and
				// recalculate both items from the ledger. Without this, the
				// invoice's ReconciliationItem would keep showing the full
				// original amount as unallocated.
				if _, err := getOrCreateInvoiceItemTx(tx, invoice); err != nil {
					return err
				}
				if _, err := getOrCreatePaymentItemTx(tx, payment); err != nil {
					return err
				}
				alloc := models.PaymentAllocation{
					PaymentID:       payment.ID,
					InvoiceID:       *payment.InvoiceID,
					AllocatedAmount: payment.Amount,
					AllocationType:  "payment",
					Reference:       "Direct match on payment recording",
					AllocatedBy:     user.UserName,
				}
				if err := tx.Create(&alloc).Error; err != nil {
					return err
				}
				if err := recalcReconciliationItem(tx, "payment", payment.ID); err != nil {
					return err
				}
				if err := recalcReconciliationItem(tx, "invoice", *payment.InvoiceID); err != nil {
					return err
				}
			}
		} else {
			payment.Status = "unmatched"
			if err := tx.Save(&payment).Error; err != nil {
				return err
			}
			// Unmatched payment: create the reconciliation item so it shows
			// up in the Unallocated Payments workspace for auto-match.
			if _, err := getOrCreatePaymentItemTx(tx, payment); err != nil {
				return err
			}
		}
		return nil
	})

	return payment, err
}

// GetArrearsStatus calculates the current arrears position for a scheme.
func GetArrearsStatus(schemeID int) (models.ArrearsRecord, error) {
	var invoices []models.Invoice
	if err := DB.Where("scheme_id = ? AND balance > 0", schemeID).Find(&invoices).Error; err != nil {
		return models.ArrearsRecord{}, err
	}

	var scheme models.GroupScheme
	if err := DB.First(&scheme, schemeID).Error; err != nil {
		return models.ArrearsRecord{}, err
	}

	record := models.ArrearsRecord{
		SchemeID:   scheme.ID,
		SchemeName: scheme.Name,
		Status:     "current",
	}

	now := time.Now()
	for _, inv := range invoices {
		dueDate, err := time.Parse("2006-01-02", inv.DueDate)
		if err != nil {
			continue
		}
		if now.Before(dueDate) {
			continue // Not yet due
		}
		daysOld := int(now.Sub(dueDate).Hours() / 24)

		if daysOld <= 30 {
			record.Days0To30 += inv.Balance
		} else if daysOld <= 60 {
			record.Days31To60 += inv.Balance
		} else if daysOld <= 90 {
			record.Days61To90 += inv.Balance
		} else {
			record.DaysOver90 += inv.Balance
		}
		record.TotalOutstanding += inv.Balance
	}

	if record.TotalOutstanding > 0 {
		record.Status = "in_arrears"
		if record.DaysOver90 > 0 {
			record.Status = "suspended"
		}
	}

	return record, nil
}

// ---------------------------------------------------------------------------
// Contribution Config
// ---------------------------------------------------------------------------

// GetContributionConfig returns the contribution config for a scheme, or a default if not found.
func GetContributionConfig(schemeID int) (models.ContributionConfig, error) {
	var cfg models.ContributionConfig
	err := DB.Where("scheme_id = ?", schemeID).Order("updated_at DESC").First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.ContributionConfig{
			SchemeID:         schemeID,
			ContributionType: "employer_only",
			EmployerPercent:  100,
			EmployeePercent:  0,
		}, nil
	}
	return cfg, err
}

// SaveContributionConfig upserts the contribution config for a scheme.
func SaveContributionConfig(cfg models.ContributionConfig, user models.AppUser) (models.ContributionConfig, error) {
	cfg.UpdatedBy = user.UserName

	var existing models.ContributionConfig
	err := DB.Where("scheme_id = ?", cfg.SchemeID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := DB.Create(&cfg).Error; err != nil {
			return models.ContributionConfig{}, err
		}
		return cfg, nil
	}
	if err != nil {
		return models.ContributionConfig{}, err
	}

	existing.ContributionType = cfg.ContributionType
	existing.EmployerPercent = cfg.EmployerPercent
	existing.EmployeePercent = cfg.EmployeePercent
	existing.EffectiveDate = cfg.EffectiveDate
	existing.UpdatedBy = cfg.UpdatedBy
	if err := DB.Save(&existing).Error; err != nil {
		return models.ContributionConfig{}, err
	}
	return existing, nil
}

// ---------------------------------------------------------------------------
// Schedules
// ---------------------------------------------------------------------------

// GetPremiumSchedules returns a filtered list of premium schedules.
func GetPremiumSchedules(schemeID, month, year int, status string) ([]models.PremiumSchedule, error) {
	query := DB.Model(&models.PremiumSchedule{})
	if schemeID > 0 {
		query = query.Where("scheme_id = ?", schemeID)
	}
	if month > 0 {
		query = query.Where("month = ?", month)
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var schedules []models.PremiumSchedule
	err := query.Order("year DESC, month DESC").Find(&schedules).Error
	return schedules, err
}

// GetScheduleDetail returns a schedule with its member rows and joiner/exit counts.
func GetScheduleDetail(scheduleID int) (models.PremiumScheduleDetail, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumScheduleDetail{}, fmt.Errorf("schedule not found: %w", err)
	}

	var rows []models.ScheduleMemberRow
	if err := DB.Where("schedule_id = ?", scheduleID).Find(&rows).Error; err != nil {
		return models.PremiumScheduleDetail{}, err
	}

	newJoiners := 0
	exits := 0
	for _, r := range rows {
		if r.IsProRata {
			// Heuristic: if ActualPremium < FullMonthPremium and positive, it's a new joiner
			// A more precise approach would store a reason field; this is sufficient for now.
			newJoiners++
		}
	}

	return models.PremiumScheduleDetail{
		PremiumSchedule: schedule,
		NewJoiners:      newJoiners,
		Exits:           exits,
		Members:         rows,
	}, nil
}

// FinalizeSchedule sets a schedule's status to "finalized". Requires "approved" status.
func FinalizeSchedule(scheduleID int, user models.AppUser) (models.PremiumSchedule, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status == "finalized" {
		return schedule, nil
	}
	if schedule.Status != "approved" {
		return models.PremiumSchedule{}, errors.New("schedule must be approved before finalizing")
	}
	now := time.Now()
	schedule.Status = "finalized"
	schedule.FinalizedBy = user.UserName
	schedule.FinalizedAt = &now
	if err := DB.Save(&schedule).Error; err != nil {
		return models.PremiumSchedule{}, err
	}
	go NotifyScheduleFinalized(schedule, user)
	return schedule, nil
}

// ReviewSchedule transitions a draft schedule to "reviewed".
func ReviewSchedule(scheduleID int, user models.AppUser) (models.PremiumSchedule, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "draft" {
		return models.PremiumSchedule{}, errors.New("only draft schedules can be submitted for review")
	}
	now := time.Now()
	schedule.Status = "reviewed"
	schedule.ReviewedBy = user.UserName
	schedule.ReviewedAt = &now
	if err := DB.Save(&schedule).Error; err != nil {
		return models.PremiumSchedule{}, err
	}
	go NotifyScheduleReviewed(schedule, user)
	return schedule, nil
}

// ApproveSchedule transitions a reviewed schedule to "approved".
func ApproveSchedule(scheduleID int, user models.AppUser) (models.PremiumSchedule, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "reviewed" {
		return models.PremiumSchedule{}, errors.New("only reviewed schedules can be approved")
	}
	now := time.Now()
	schedule.Status = "approved"
	schedule.ApprovedBy = user.UserName
	schedule.ApprovedAt = &now
	if err := DB.Save(&schedule).Error; err != nil {
		return models.PremiumSchedule{}, err
	}
	go NotifyScheduleApproved(schedule, user)
	return schedule, nil
}

// ReturnScheduleToDraft returns a reviewed or approved schedule back to draft.
func ReturnScheduleToDraft(scheduleID int, user models.AppUser) (models.PremiumSchedule, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "reviewed" && schedule.Status != "approved" {
		return models.PremiumSchedule{}, errors.New("only reviewed or approved schedules can be returned to draft")
	}
	schedule.Status = "draft"
	schedule.ReviewedBy = ""
	schedule.ReviewedAt = nil
	schedule.ApprovedBy = ""
	schedule.ApprovedAt = nil
	if err := DB.Save(&schedule).Error; err != nil {
		return models.PremiumSchedule{}, err
	}
	return schedule, nil
}

// VoidSchedule voids a finalized or invoiced schedule with a mandatory reason.
func VoidSchedule(scheduleID int, reason string, user models.AppUser) (models.PremiumSchedule, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "finalized" && schedule.Status != "invoiced" {
		return models.PremiumSchedule{}, errors.New("only finalized or invoiced schedules can be voided")
	}
	now := time.Now()
	schedule.Status = "void"
	schedule.VoidedBy = user.UserName
	schedule.VoidedAt = &now
	schedule.VoidReason = reason
	if err := DB.Save(&schedule).Error; err != nil {
		return models.PremiumSchedule{}, err
	}
	go NotifyScheduleVoided(schedule, user)
	return schedule, nil
}

// CancelSchedule cancels a draft, reviewed, or approved schedule.
func CancelSchedule(scheduleID int, reason string, user models.AppUser) (models.PremiumSchedule, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "draft" && schedule.Status != "reviewed" && schedule.Status != "approved" {
		return models.PremiumSchedule{}, errors.New("only draft, reviewed, or approved schedules can be cancelled")
	}
	now := time.Now()
	schedule.Status = "cancelled"
	schedule.VoidedBy = user.UserName
	schedule.VoidedAt = &now
	schedule.VoidReason = reason
	if err := DB.Save(&schedule).Error; err != nil {
		return models.PremiumSchedule{}, err
	}
	go NotifyScheduleCancelled(schedule, user)
	return schedule, nil
}

// recalculateScheduleTotals recomputes GrossPremium, NetPayable, and MemberCount from rows.
func recalculateScheduleTotals(tx *gorm.DB, schedule *models.PremiumSchedule) error {
	var rows []models.ScheduleMemberRow
	if err := tx.Where("schedule_id = ?", schedule.ID).Find(&rows).Error; err != nil {
		return err
	}
	var gross float64
	for _, r := range rows {
		gross += r.ActualPremium
	}
	schedule.GrossPremium = gross
	schedule.NetPayable = gross
	schedule.MemberCount = len(rows)
	return tx.Save(schedule).Error
}

// RegenerateSchedule re-runs the full member row calculation for a draft schedule,
// replacing all existing rows with a fresh set drawn from the current member data.
// This allows picking up mid-cycle changes (new joiners, exits, rate updates) from
// external systems without leaving the schedule.
func RegenerateSchedule(scheduleID int, user models.AppUser) (models.PremiumSchedule, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "draft" {
		return models.PremiumSchedule{}, errors.New("only draft schedules can be regenerated")
	}

	// Load contribution config; fall back to 100% employer if not configured
	employerPct := 100.0
	employeePct := 0.0
	var cfg models.ContributionConfig
	if err := DB.Where("scheme_id = ?", schedule.SchemeID).Order("updated_at DESC").First(&cfg).Error; err == nil {
		if cfg.ContributionType == "split" {
			employerPct = cfg.EmployerPercent
			employeePct = cfg.EmployeePercent
		}
	}

	//// Days in the schedule month for pro-rata calculation
	//firstOfMonth := time.Date(schedule.Year, time.Month(schedule.Month), 1, 0, 0, 0, 0, time.UTC)
	//firstOfNext := firstOfMonth.AddDate(0, 1, 0)
	//daysInMonth := int(firstOfNext.Sub(firstOfMonth).Hours() / 24)

	var scheme models.GroupScheme
	if err := DB.Where("id = ?", schedule.SchemeID).Find(&scheme).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("scheme not found: %w", err)
	}
	// Assuming 'commencementDate' is a time.Time variable holding the start of the scheme
	commencementDate := time.Date(scheme.CommencementDate.Year(), scheme.CommencementDate.Month(), scheme.CommencementDate.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate the first day of the target month and the first day of the next month
	firstOfMonth := time.Date(schedule.Year, time.Month(schedule.Month), 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	// Determine the effective start date for the calculation.
	// If the scheme commenced *after* the first of the target month, we use the commencement date.
	// Otherwise, we use the first of the month.
	effectiveStart := firstOfMonth
	if commencementDate.After(firstOfMonth) {
		effectiveStart = commencementDate
	}

	// If the effective start is on or after the first day of the next month,
	// it means the scheme started *after* the entire target month, so there are 0 days.
	var daysInMonth int
	if !effectiveStart.Before(firstOfNextMonth) {
		daysInMonth = 0
	} else {
		// Calculate the duration between the effective start and the end of the month.
		// Using Truncate to get a whole number of days.
		duration := firstOfNextMonth.Sub(effectiveStart)
		daysInMonth = int(duration.Truncate(24*time.Hour).Hours() / 24)
	}

	// Fetch all active members for this scheme
	var members []models.GPricingMemberDataInForce
	if err := DB.Where("scheme_id = ?", schedule.SchemeID).Find(&members).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("failed to fetch members: %w", err)
	}

	// Pre-fetch MemberPremiumSchedule rows; build map keyed by member name
	var allMPS []models.MemberPremiumSchedule
	if err := DB.Where("scheme_id = ?", schedule.SchemeID).Order("creation_date DESC").Find(&allMPS).Error; err != nil {
		return models.PremiumSchedule{}, fmt.Errorf("failed to fetch premium schedules: %w", err)
	}
	mpsMap := make(map[string]models.MemberPremiumSchedule, len(allMPS))
	for _, mps := range allMPS {
		if _, exists := mpsMap[mps.MemberName]; !exists {
			mpsMap[mps.MemberName] = mps
		}
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		// Drop all existing rows for this schedule
		if err := tx.Where("schedule_id = ?", scheduleID).Delete(&models.ScheduleMemberRow{}).Error; err != nil {
			return err
		}

		var grossTotal float64

		rows := make([]models.ScheduleMemberRow, 0, len(members))

		for _, m := range members {
			mps, ok := mpsMap[m.MemberName]
			if !ok {
				continue
			}

			fullMonthPremium := mps.TotalAnnualPremiumPayable / 12.0

			isProRata := false
			var proRataDays *int
			var actualPremium float64

			if daysInMonth > 0 {

				actualPremium = fullMonthPremium

				// Pro-rata: member joined during this month
				entryYear := m.EntryDate.Year()
				entryMonth := int(m.EntryDate.Month())
				if entryYear == schedule.Year && entryMonth == schedule.Month && m.EntryDate.Day() > 1 {
					daysRemaining := daysInMonth - m.EntryDate.Day() + 1
					actualPremium = fullMonthPremium * float64(daysRemaining) / float64(daysInMonth)
					isProRata = true
					proRataDays = &daysRemaining
				}

				//Pro-rata: member exited during this month
				//if !m.ExitDate.IsZero() {
				//	exitYear := m.ExitDate.Year()
				//	exitMonth := int(m.ExitDate.Month())
				//	if exitYear == schedule.Year && exitMonth == schedule.Month {
				//		daysActive := m.ExitDate.Day()
				//		actualPremium = fullMonthPremium * float64(daysActive) / float64(daysInMonth)
				//		isProRata = true
				//		proRataDays = &daysActive
				//	}
				//}
				//
				if m.ExitDate != nil {
					monthStart := time.Date(schedule.Year, time.Month(schedule.Month), 1, 0, 0, 0, 0, time.UTC)
					monthEnd := monthStart.AddDate(0, 1, -1) // last day of month
					daysInMonth := monthEnd.Day()

					// Normalize dates
					var effectiveStart time.Time
					if m.EntryDate.Before(monthStart) {
						effectiveStart = monthStart
					}
					if m.EntryDate.After(monthStart) {
						effectiveStart = m.EntryDate
					}

					var exitDate time.Time
					if m.ExitDate != nil {
						exitDate = *m.ExitDate
					}

					// Step 1: determine active interval within the month
					startActive := monthStart
					if !effectiveStart.IsZero() && effectiveStart.After(monthStart) {
						startActive = effectiveStart
					}

					// If there is an exit date, limit to that; otherwise assume still active through month end
					endActive := monthEnd
					if !exitDate.IsZero() && exitDate.Before(monthEnd) {
						endActive = exitDate
					}

					// Step 2: compute overlap duration
					if startActive.After(endActive) {
						// no overlap; no days active
						actualPremium = 0
						isProRata = false
						proRataDays = nil
					} else {
						// overlap exists
						daysActive := endActive.Sub(startActive).Hours()/24 + 1 // inclusive days
						// guard: ensure daysActive is within 0..daysInMonth
						if daysActive < 0 {
							daysActive = 0
						}
						if int(daysActive) > daysInMonth {
							daysActive = float64(daysInMonth)
						}

						// Step 3: decide premium
						if int(daysActive) == daysInMonth {
							// full month
							actualPremium = fullMonthPremium
							isProRata = false
							proRataDays = nil
						} else {
							// pro-rated
							actualPremium = fullMonthPremium * float64(daysActive) / float64(daysInMonth)
							isProRata = true
							d := int(daysActive)
							proRataDays = &d
						}
					}
				}
			}

			rows = append(rows, models.ScheduleMemberRow{
				ScheduleID:           schedule.ID,
				MemberID:             m.ID,
				MemberName:           m.MemberName,
				MemberIDNumber:       m.MemberIdNumber,
				AnnualSalary:         m.AnnualSalary,
				FullMonthPremium:     fullMonthPremium,
				ActualPremium:        actualPremium,
				ProRataDays:          proRataDays,
				IsProRata:            isProRata,
				EmployerContribution: actualPremium * employerPct / 100.0,
				EmployeeContribution: actualPremium * employeePct / 100.0,
			})
			grossTotal += actualPremium
		}

		if len(rows) > 0 {
			if err := tx.CreateInBatches(rows, 500).Error; err != nil {
				return err
			}
		}

		schedule.MemberCount = len(rows)
		schedule.GrossPremium = grossTotal
		schedule.NetPayable = grossTotal
		schedule.GeneratedBy = user.UserName
		return tx.Save(&schedule).Error
	})

	return schedule, err
}

// RemoveScheduleMemberRow removes a member row from a draft schedule.
func RemoveScheduleMemberRow(scheduleID int, rowID int, user models.AppUser) error {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "draft" {
		return errors.New("members can only be removed from draft schedules")
	}

	var row models.ScheduleMemberRow
	if err := DB.First(&row, rowID).Error; err != nil {
		return fmt.Errorf("row not found: %w", err)
	}
	if row.ScheduleID != scheduleID {
		return errors.New("row does not belong to this schedule")
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&row).Error; err != nil {
			return err
		}
		return recalculateScheduleTotals(tx, &schedule)
	})
}

// UpdateScheduleMemberRow updates rate fields on a member row in a draft schedule.
func UpdateScheduleMemberRow(scheduleID int, rowID int, req models.ScheduleMemberRateRequest, user models.AppUser) (models.ScheduleMemberRow, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return models.ScheduleMemberRow{}, fmt.Errorf("schedule not found: %w", err)
	}
	if schedule.Status != "draft" {
		return models.ScheduleMemberRow{}, errors.New("member rows can only be updated in draft schedules")
	}

	var row models.ScheduleMemberRow
	if err := DB.First(&row, rowID).Error; err != nil {
		return models.ScheduleMemberRow{}, fmt.Errorf("row not found: %w", err)
	}
	if row.ScheduleID != scheduleID {
		return models.ScheduleMemberRow{}, errors.New("row does not belong to this schedule")
	}

	row.Rate = req.Rate
	row.ActualPremium = req.ActualPremium
	row.EmployerContribution = req.EmployerContrib
	row.EmployeeContribution = req.EmployeeContrib

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&row).Error; err != nil {
			return err
		}
		return recalculateScheduleTotals(tx, &schedule)
	})

	return row, err
}

// ExportScheduleCSV returns a CSV of the member rows for a schedule.
func ExportScheduleCSV(scheduleID int) ([]byte, error) {
	var schedule models.PremiumSchedule
	if err := DB.First(&schedule, scheduleID).Error; err != nil {
		return nil, fmt.Errorf("schedule not found: %w", err)
	}

	var rows []models.ScheduleMemberRow
	if err := DB.Where("schedule_id = ?", scheduleID).Find(&rows).Error; err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{
		"Member ID", "Member Name", "ID Number", "Benefit",
		"Annual Salary", "Full Month Premium", "Pro Rata Days",
		"Actual Premium", "Employer Contribution", "Employee Contribution",
	})
	for _, r := range rows {
		proRataDays := ""
		if r.ProRataDays != nil {
			proRataDays = strconv.Itoa(*r.ProRataDays)
		}
		_ = w.Write([]string{
			strconv.Itoa(r.MemberID),
			r.MemberName,
			r.MemberIDNumber,
			r.Benefit,
			fmt.Sprintf("%.2f", r.AnnualSalary),
			fmt.Sprintf("%.2f", r.FullMonthPremium),
			proRataDays,
			fmt.Sprintf("%.2f", r.ActualPremium),
			fmt.Sprintf("%.2f", r.EmployerContribution),
			fmt.Sprintf("%.2f", r.EmployeeContribution),
		})
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}

// ---------------------------------------------------------------------------
// Invoices
// ---------------------------------------------------------------------------

// GetInvoices returns a filtered list of invoices.
func GetInvoices(schemeID, month, year int, status string) ([]models.Invoice, error) {
	query := DB.Model(&models.Invoice{})
	if schemeID > 0 {
		query = query.Where("scheme_id = ?", schemeID)
	}
	if month > 0 {
		query = query.Where("month = ?", month)
	}
	if year > 0 {
		query = query.Where("year = ?", year)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var invoices []models.Invoice
	err := query.Order("issue_date DESC").Find(&invoices).Error
	return invoices, err
}

// GetInvoiceDetail returns an invoice with line items, adjustments, and payment history.
func GetInvoiceDetail(invoiceID int) (models.InvoiceDetail, error) {
	var invoice models.Invoice
	if err := DB.First(&invoice, invoiceID).Error; err != nil {
		return models.InvoiceDetail{}, fmt.Errorf("invoice not found: %w", err)
	}

	// Build line items by aggregating schedule rows by benefit
	var rows []models.ScheduleMemberRow
	if invoice.ScheduleID > 0 {
		DB.Where("schedule_id = ?", invoice.ScheduleID).Find(&rows)
	}
	lineMap := map[string]*models.InvoiceLineItem{}
	for _, r := range rows {
		benefit := r.Benefit
		if benefit == "" {
			benefit = "General"
		}
		if _, ok := lineMap[benefit]; !ok {
			lineMap[benefit] = &models.InvoiceLineItem{Benefit: benefit}
		}
		lineMap[benefit].MemberCount++
		lineMap[benefit].BasePremium += r.ActualPremium
		lineMap[benefit].Total += r.ActualPremium
	}
	lineItems := make([]models.InvoiceLineItem, 0, len(lineMap))
	for _, li := range lineMap {
		lineItems = append(lineItems, *li)
	}

	var adjustments []models.InvoiceAdjustment
	DB.Where("invoice_id = ?", invoiceID).Find(&adjustments)

	var payments []models.Payment
	DB.Where("invoice_id = ?", invoiceID).Find(&payments)

	// Fetch contact info from scheme
	var scheme models.GroupScheme
	contactName := ""
	contactEmail := ""
	if err := DB.First(&scheme, invoice.SchemeID).Error; err == nil {
		contactName = scheme.ContactPerson
		contactEmail = scheme.ContactEmail
	}

	return models.InvoiceDetail{
		Invoice:        invoice,
		LineItems:      lineItems,
		Adjustments:    adjustments,
		PaymentHistory: payments,
		ContactName:    contactName,
		ContactEmail:   contactEmail,
	}, nil
}

// MarkInvoiceSent marks an invoice status as "sent".
func MarkInvoiceSent(invoiceID int, user models.AppUser) (models.Invoice, error) {
	var invoice models.Invoice
	if err := DB.First(&invoice, invoiceID).Error; err != nil {
		return models.Invoice{}, fmt.Errorf("invoice not found: %w", err)
	}
	invoice.Status = "sent"
	if err := DB.Save(&invoice).Error; err != nil {
		return models.Invoice{}, err
	}
	return invoice, nil
}

// CreateAdjustmentNote creates a credit or debit note on an invoice.
func CreateAdjustmentNote(req models.AdjustmentNoteRequest, user models.AppUser) (models.InvoiceAdjustment, error) {
	var invoice models.Invoice
	if err := DB.First(&invoice, req.InvoiceID).Error; err != nil {
		return models.InvoiceAdjustment{}, fmt.Errorf("invoice not found: %w", err)
	}

	adj := models.InvoiceAdjustment{
		InvoiceID:   req.InvoiceID,
		Description: req.Reason,
		Amount:      req.Amount,
		Type:        req.Type,
		CreatedBy:   user.UserName,
		CreatedAt:   time.Now().Format("2006-01-02"),
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&adj).Error; err != nil {
			return err
		}
		// Apply adjustment to invoice balance
		if req.Type == "credit" {
			invoice.Balance -= req.Amount
			invoice.PaidAmount += req.Amount
		} else {
			invoice.Balance += req.Amount
		}
		if invoice.Balance <= AutoCloseTolerance {
			invoice.Status = "paid"
			invoice.Balance = 0
		}
		return tx.Save(&invoice).Error
	})

	return adj, err
}

// ---------------------------------------------------------------------------
// Payments
// ---------------------------------------------------------------------------

// GetPayments returns a filtered list of payments.
func GetPayments(schemeID int, method, status string) ([]models.Payment, error) {
	query := DB.Model(&models.Payment{})
	if schemeID > 0 {
		query = query.Where("scheme_id = ?", schemeID)
	}
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var payments []models.Payment
	err := query.Order("payment_date DESC").Find(&payments).Error
	return payments, err
}

// VoidPayment voids a payment and reverses any invoice balance update.
func VoidPayment(paymentID int, user models.AppUser) error {
	var payment models.Payment
	if err := DB.First(&payment, paymentID).Error; err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}
	if payment.Status == "voided" {
		return errors.New("payment is already voided")
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		// Reverse invoice balance if matched
		if payment.InvoiceID != nil && payment.Status == "matched" {
			var invoice models.Invoice
			if err := tx.First(&invoice, *payment.InvoiceID).Error; err == nil {
				invoice.PaidAmount -= payment.Amount
				invoice.Balance = invoice.NetPayable - invoice.PaidAmount
				if invoice.Balance > 0 && invoice.Status == "paid" {
					invoice.Status = "partial"
				}
				if invoice.PaidAmount <= 0 {
					invoice.Status = "sent"
				}
				if err := tx.Save(&invoice).Error; err != nil {
					return err
				}
			}
		}
		payment.Status = "voided"
		return tx.Save(&payment).Error
	})
}

// BulkImportPayments parses CSV rows and records payments, matching by invoice number or bank reference.
func BulkImportPayments(rows [][]string, user models.AppUser) (models.BulkImportResult, error) {
	result := models.BulkImportResult{
		TotalRows: len(rows),
		Rows:      make([]models.BulkImportRow, 0, len(rows)),
	}

	for i, row := range rows {
		rowNum := i + 1
		if len(row) < 2 {
			result.Errors++
			result.Rows = append(result.Rows, models.BulkImportRow{
				RowNumber:   rowNum,
				Status:      "error",
				ErrorReason: "insufficient columns (need at least Reference, Amount)",
			})
			continue
		}

		reference := row[0]
		amount, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			result.Errors++
			result.Rows = append(result.Rows, models.BulkImportRow{
				RowNumber:   rowNum,
				Reference:   reference,
				Status:      "error",
				ErrorReason: "invalid amount: " + row[1],
			})
			continue
		}

		bankRef := ""
		if len(row) >= 3 {
			bankRef = row[2]
		}
		invoiceNum := ""
		if len(row) >= 4 {
			invoiceNum = row[3]
		}

		// Try to match to an invoice
		var invoice models.Invoice
		matchedInvoice := ""
		var invoiceID *int
		if invoiceNum != "" {
			if err := DB.Where("invoice_number = ?", invoiceNum).First(&invoice).Error; err == nil {
				matchedInvoice = invoice.InvoiceNumber
				id := invoice.ID
				invoiceID = &id
			}
		}
		if matchedInvoice == "" && bankRef != "" {
			if err := DB.Where("bank_reference = ? AND status = ?", bankRef, "sent").First(&invoice).Error; err == nil {
				matchedInvoice = invoice.InvoiceNumber
				id := invoice.ID
				invoiceID = &id
			}
		}

		schemeID := 0
		if invoiceID != nil {
			schemeID = invoice.SchemeID
		}

		req := models.RecordPaymentRequest{
			SchemeID:      schemeID,
			InvoiceID:     invoiceID,
			PaymentDate:   time.Now().Format("2006-01-02"),
			Method:        "eft",
			Amount:        amount,
			BankReference: bankRef,
			Notes:         "Bulk import row " + strconv.Itoa(rowNum),
		}

		if schemeID == 0 {
			// Still record unmatched payment without scheme
			payment := models.Payment{
				PaymentDate:   req.PaymentDate,
				Method:        req.Method,
				Amount:        amount,
				BankReference: bankRef,
				Status:        "unmatched",
				RecordedBy:    user.UserName,
				RecordedAt:    time.Now().Format(time.RFC3339),
			}
			DB.Create(&payment)
			result.Unmatched++
			result.Rows = append(result.Rows, models.BulkImportRow{
				RowNumber: rowNum,
				Reference: reference,
				Amount:    amount,
				Status:    "unmatched",
			})
			continue
		}

		if _, err := RecordPayment(req, user); err != nil {
			result.Errors++
			result.Rows = append(result.Rows, models.BulkImportRow{
				RowNumber:   rowNum,
				Reference:   reference,
				Amount:      amount,
				Status:      "error",
				ErrorReason: err.Error(),
			})
			continue
		}

		if matchedInvoice != "" {
			result.Matched++
			result.Rows = append(result.Rows, models.BulkImportRow{
				RowNumber:      rowNum,
				Reference:      reference,
				Amount:         amount,
				MatchedInvoice: matchedInvoice,
				Status:         "matched",
			})
		} else {
			result.Unmatched++
			result.Rows = append(result.Rows, models.BulkImportRow{
				RowNumber: rowNum,
				Reference: reference,
				Amount:    amount,
				Status:    "unmatched",
			})
		}
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Reconciliation
// ---------------------------------------------------------------------------

// AutoMatchPayments attempts to match all unmatched payments to outstanding invoices.
// Strategy: (1) match by bank_reference → invoice_number, (2) fall back to scheme + amount.
func AutoMatchPayments(user models.AppUser) (models.AutoMatchResult, error) {
	var unmatched []models.Payment
	if err := DB.Where("status = ?", "unmatched").Find(&unmatched).Error; err != nil {
		return models.AutoMatchResult{}, err
	}

	matchedCount := 0
	for _, p := range unmatched {
		var invoice models.Invoice
		var found bool

		// Strategy 1: match by bank reference → invoice number
		if p.BankReference != "" {
			err := DB.Where("invoice_number = ? AND status IN ?",
				p.BankReference, []string{"sent", "partial"}).
				First(&invoice).Error
			if err == nil {
				found = true
			}
		}

		// Strategy 2: match by scheme + amount (within tolerance)
		if !found {
			err := DB.Where("scheme_id = ? AND balance >= ? AND balance <= ? AND status IN ?",
				p.SchemeID, p.Amount-0.01, p.Amount+0.01, []string{"sent", "partial"}).
				Order("issue_date ASC").First(&invoice).Error
			if err != nil {
				continue
			}
		}

		// Link the existing payment to the invoice in a transaction
		err := DB.Transaction(func(tx *gorm.DB) error {
			invoice.PaidAmount += p.Amount
			invoice.Balance = invoice.NetPayable - invoice.PaidAmount
			if invoice.Balance <= AutoCloseTolerance {
				invoice.Status = "paid"
				invoice.Balance = 0
			} else {
				invoice.Status = "partial"
			}
			if err := tx.Save(&invoice).Error; err != nil {
				return err
			}

			p.Status = "matched"
			p.InvoiceID = &invoice.ID
			p.InvoiceNumber = invoice.InvoiceNumber
			return tx.Save(&p).Error
		})
		if err == nil {
			matchedCount++
		}
	}

	var remaining int64
	DB.Model(&models.Payment{}).Where("status = ?", "unmatched").Count(&remaining)

	return models.AutoMatchResult{
		Matched:   matchedCount,
		Remaining: int(remaining),
	}, nil
}

// ManualMatchPayment links a payment to a specific invoice.
func ManualMatchPayment(req models.ManualMatchRequest, user models.AppUser) (models.Payment, error) {
	var payment models.Payment
	if err := DB.First(&payment, req.PaymentID).Error; err != nil {
		return models.Payment{}, fmt.Errorf("payment not found: %w", err)
	}

	if payment.Status != "unmatched" {
		return models.Payment{}, fmt.Errorf("payment is already %s and cannot be matched", payment.Status)
	}

	var invoice models.Invoice
	if err := DB.First(&invoice, req.InvoiceID).Error; err != nil {
		return models.Payment{}, fmt.Errorf("invoice not found: %w", err)
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		invoice.PaidAmount += payment.Amount
		invoice.Balance = invoice.NetPayable - invoice.PaidAmount
		if invoice.Balance <= AutoCloseTolerance {
			invoice.Status = "paid"
			invoice.Balance = 0
		} else {
			invoice.Status = "partial"
		}
		if err := tx.Save(&invoice).Error; err != nil {
			return err
		}

		payment.Status = "matched"
		payment.InvoiceID = &invoice.ID
		payment.InvoiceNumber = invoice.InvoiceNumber
		return tx.Save(&payment).Error
	})
	return payment, err
}

// ---------------------------------------------------------------------------
// Arrears
// ---------------------------------------------------------------------------

// GetAllArrearsAging computes aging for all in-force schemes.
func GetAllArrearsAging(status string) ([]models.ArrearsRecord, error) {
	var schemes []models.GroupScheme
	query := DB.Model(&models.GroupScheme{}).Where("in_force = ?", true)
	if err := query.Find(&schemes).Error; err != nil {
		return nil, err
	}

	records := make([]models.ArrearsRecord, 0, len(schemes))
	for _, s := range schemes {
		rec, err := GetArrearsStatus(s.ID)
		if err != nil {
			continue
		}
		// Filter by status if provided
		if status != "" && rec.Status != status {
			continue
		}
		records = append(records, rec)
	}
	return records, nil
}

// SendArrearsReminder logs a REMINDER event in arrears history.
func SendArrearsReminder(schemeID int, req models.SendReminderRequest, user models.AppUser) error {
	event := models.ArrearsHistory{
		SchemeID:    schemeID,
		EventType:   "REMINDER",
		EventDate:   time.Now().Format("2006-01-02"),
		Description: req.Message,
		PerformedBy: user.UserName,
	}
	return DB.Create(&event).Error
}

// RecordPaymentPlan creates a payment plan with instalments for a scheme.
func RecordPaymentPlan(schemeID int, req models.PaymentPlanRequest, user models.AppUser) (models.PaymentPlan, error) {
	plan := models.PaymentPlan{
		SchemeID:  schemeID,
		Notes:     req.Notes,
		CreatedBy: user.UserName,
		Status:    "active",
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&plan).Error; err != nil {
			return err
		}
		for _, inst := range req.Instalments {
			instalment := models.PaymentPlanInstalment{
				PlanID: plan.ID,
				Date:   inst.Date,
				Amount: inst.Amount,
				Status: "pending",
			}
			if err := tx.Create(&instalment).Error; err != nil {
				return err
			}
		}

		// Log to arrears history
		event := models.ArrearsHistory{
			SchemeID:    schemeID,
			EventType:   "PAYMENT_PLAN",
			EventDate:   time.Now().Format("2006-01-02"),
			Description: fmt.Sprintf("Payment plan created with %d instalments", len(req.Instalments)),
			PerformedBy: user.UserName,
		}
		return tx.Create(&event).Error
	})

	return plan, err
}

// SuspendSchemeCover updates the scheme status to suspended and logs the event.
func SuspendSchemeCover(schemeID int, req models.SuspendCoverRequest, user models.AppUser) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.GroupScheme{}).Where("id = ?", schemeID).
			Update("status", "SUSPENDED").Error; err != nil {
			return err
		}
		event := models.ArrearsHistory{
			SchemeID:    schemeID,
			EventType:   "SUSPENDED",
			EventDate:   req.EffectiveDate,
			Description: req.Reason,
			PerformedBy: user.UserName,
		}
		return tx.Create(&event).Error
	})
	if err != nil {
		return err
	}
	var scheme models.GroupScheme
	if DB.First(&scheme, schemeID).Error == nil {
		go NotifySchemeSuspended(scheme, user)
	}
	return nil
}

// ReinstateSchemeCover reinstates cover, optionally raising a back-premium invoice.
func ReinstateSchemeCover(schemeID int, req models.ReinstateCoverRequest, user models.AppUser) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.GroupScheme{}).Where("id = ?", schemeID).
			Update("status", "ACTIVE").Error; err != nil {
			return err
		}

		// Create back-premium invoice if requested
		if req.BackPremium > 0 {
			var scheme models.GroupScheme
			tx.First(&scheme, schemeID)
			backInvoice := models.Invoice{
				InvoiceNumber: fmt.Sprintf("BACK-%d-%s", schemeID, time.Now().Format("20060102")),
				SchemeID:      schemeID,
				SchemeName:    scheme.Name,
				IssueDate:     req.ReinstatementDate,
				DueDate:       time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
				GrossAmount:   req.BackPremium,
				NetPayable:    req.BackPremium,
				Balance:       req.BackPremium,
				Status:        "sent",
			}
			if err := tx.Create(&backInvoice).Error; err != nil {
				return err
			}
		}

		event := models.ArrearsHistory{
			SchemeID:    schemeID,
			EventType:   "REINSTATED",
			EventDate:   req.ReinstatementDate,
			Description: req.Notes,
			Amount:      req.BackPremium,
			PerformedBy: user.UserName,
		}
		return tx.Create(&event).Error
	})
	if err != nil {
		return err
	}
	var scheme models.GroupScheme
	if DB.First(&scheme, schemeID).Error == nil {
		go NotifySchemeReinstated(scheme, user)
	}
	return nil
}

// GetArrearsHistory returns the arrears event history for a scheme.
func GetArrearsHistory(schemeID int) ([]models.ArrearsHistory, error) {
	var history []models.ArrearsHistory
	err := DB.Where("scheme_id = ?", schemeID).Order("created_at DESC").Find(&history).Error
	return history, err
}

// GetPaymentPlans returns all payment plans for a scheme, each with its
// ordered instalments. Plans are ordered newest-first.
func GetPaymentPlans(schemeID int) ([]models.PaymentPlanWithInstalments, error) {
	var plans []models.PaymentPlan
	if err := DB.Where("scheme_id = ?", schemeID).
		Order("created_at DESC").
		Find(&plans).Error; err != nil {
		return nil, err
	}

	result := make([]models.PaymentPlanWithInstalments, 0, len(plans))
	for _, p := range plans {
		var instalments []models.PaymentPlanInstalment
		if err := DB.Where("plan_id = ?", p.ID).
			Order("date ASC").
			Find(&instalments).Error; err != nil {
			return nil, err
		}
		result = append(result, models.PaymentPlanWithInstalments{
			PaymentPlan: p,
			Instalments: instalments,
		})
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Statements
// ---------------------------------------------------------------------------

// GetEmployerStatement aggregates invoices, payments and adjustments into a ledger.
func GetEmployerStatement(schemeID int, from, to string) (models.EmployerStatement, error) {
	var scheme models.GroupScheme
	if err := DB.First(&scheme, schemeID).Error; err != nil {
		return models.EmployerStatement{}, fmt.Errorf("scheme not found: %w", err)
	}

	// Opening balance: sum of balances on invoices before the from date
	var openingBalance float64
	DB.Model(&models.Invoice{}).
		Where("scheme_id = ? AND issue_date < ?", schemeID, from).
		Select("COALESCE(SUM(balance), 0)").Scan(&openingBalance)

	// Invoices in period
	var invoices []models.Invoice
	DB.Where("scheme_id = ? AND issue_date >= ? AND issue_date <= ?", schemeID, from, to).
		Order("issue_date ASC").Find(&invoices)

	// Payments in period
	var payments []models.Payment
	DB.Where("scheme_id = ? AND payment_date >= ? AND payment_date <= ? and status !=?", schemeID, from, to, "voided").
		Order("payment_date ASC").Find(&payments)

	// Adjustments in period
	var adjustments []models.InvoiceAdjustment
	DB.Joins("JOIN invoices ON invoices.id = invoice_adjustments.invoice_id").
		Where("invoices.scheme_id = ? AND invoice_adjustments.created_at >= ? AND invoice_adjustments.created_at <= ?",
			schemeID, from, to).Find(&adjustments)

	lineItems := make([]models.StatementLineItem, 0)
	runningBalance := openingBalance
	var invoicedAmount, received, adjustmentsTotal float64

	for _, inv := range invoices {
		runningBalance += inv.NetPayable
		invoicedAmount += inv.NetPayable
		lineItems = append(lineItems, models.StatementLineItem{
			Date:        inv.IssueDate,
			Description: fmt.Sprintf("Invoice %s", inv.InvoiceNumber),
			Debit:       inv.NetPayable,
			Balance:     runningBalance,
		})
	}
	for _, p := range payments {
		runningBalance -= p.Amount
		received += p.Amount
		lineItems = append(lineItems, models.StatementLineItem{
			Date:        p.PaymentDate,
			Description: fmt.Sprintf("Payment (%s) %s", p.Method, p.BankReference),
			Credit:      p.Amount,
			Balance:     runningBalance,
		})
	}
	for _, adj := range adjustments {
		if adj.Type == "credit" {
			runningBalance -= adj.Amount
			adjustmentsTotal -= adj.Amount
		} else {
			runningBalance += adj.Amount
			adjustmentsTotal += adj.Amount
		}
		adjCredit := 0.0
		adjDebit := 0.0
		if adj.Type == "credit" {
			adjCredit = adj.Amount
		} else {
			adjDebit = adj.Amount
		}
		lineItems = append(lineItems, models.StatementLineItem{
			Date:        adj.CreatedAt,
			Description: fmt.Sprintf("Adjustment: %s", adj.Description),
			Credit:      adjCredit,
			Debit:       adjDebit,
			Balance:     runningBalance,
		})
	}

	return models.EmployerStatement{
		SchemeID:       schemeID,
		SchemeName:     scheme.Name,
		Period:         from + " to " + to,
		OpeningBalance: openingBalance,
		InvoicedAmount: invoicedAmount,
		Received:       received,
		Adjustments:    adjustmentsTotal,
		ClosingBalance: runningBalance,
		LineItems:      lineItems,
	}, nil
}

// GetBrokerCommissionStatement aggregates commissions for all schemes under a broker.
func GetBrokerCommissionStatement(brokerID int, from, to string) (models.BrokerCommissionStatement, error) {
	var broker models.Broker
	if err := DB.First(&broker, brokerID).Error; err != nil {
		return models.BrokerCommissionStatement{}, fmt.Errorf("broker not found: %w", err)
	}
	var commRate float64

	var schemes []models.GroupScheme
	DB.Where("broker_id = ? AND in_force = ?", brokerID, true).Find(&schemes)

	commissions := make([]models.BrokerSchemeCommission, 0, len(schemes))
	totalEarned := 0.0

	for _, s := range schemes {
		var collected float64
		DB.Model(&models.Payment{}).
			Where("scheme_id = ? AND payment_date >= ? AND payment_date <= ? AND status = ?",
				s.ID, from, to, "matched").
			Select("COALESCE(SUM(amount), 0)").Scan(&collected)

		if s.AnnualPremium > 0 {
			commRate = s.Commission / s.AnnualPremium
		}

		earned := collected * commRate
		totalEarned += earned

		commissions = append(commissions, models.BrokerSchemeCommission{
			SchemeID:         s.ID,
			SchemeName:       s.Name,
			PremiumCollected: collected,
			CommissionRate:   commRate,
			CommissionEarned: earned,
			Status:           string(s.Status),
		})
	}

	return models.BrokerCommissionStatement{
		BrokerID:    brokerID,
		BrokerName:  broker.Name,
		Period:      from + " to " + to,
		TotalEarned: totalEarned,
		Schemes:     commissions,
	}, nil
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// GetPremiumDashboard returns aggregate KPIs, trends, and outstanding schemes.
func GetPremiumDashboard(year int) (models.PremiumDashboardData, error) {
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()
	if year == 0 {
		year = currentYear
	}

	// KPIs
	var dueThisMonth, collected float64
	DB.Model(&models.Invoice{}).
		Where("month = ? AND year = ?", currentMonth, currentYear).
		Select("COALESCE(SUM(net_payable), 0)").Scan(&dueThisMonth)

	DB.Model(&models.Invoice{}).
		Where("month = ? AND year = ?", currentMonth, currentYear).
		Select("COALESCE(SUM(paid_amount), 0)").Scan(&collected)

	collectionRate := 0.0
	if dueThisMonth > 0 {
		collectionRate = collected / dueThisMonth * 100
	}

	var outstanding float64
	DB.Model(&models.Invoice{}).
		Where("balance > 0").
		Select("COALESCE(SUM(balance), 0)").Scan(&outstanding)

	today := now.Format("2006-01-02")
	var overdue float64
	DB.Model(&models.Invoice{}).
		Where("balance > 0 AND due_date < ?", today).
		Select("COALESCE(SUM(balance), 0)").Scan(&overdue)

	var overdueCount int64
	DB.Model(&models.Invoice{}).
		Where("balance > 0 AND due_date < ?", today).
		Group("scheme_id").Count(&overdueCount)

	kpis := models.PremiumDashboardKPIs{
		DueThisMonth:       dueThisMonth,
		Collected:          collected,
		CollectionRate:     collectionRate,
		Outstanding:        outstanding,
		Overdue:            overdue,
		OverdueSchemeCount: int(overdueCount),
	}

	// 12-month trend
	trend := make([]models.MonthlyPremiumTrend, 0, 12)
	for i := 11; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		m := int(t.Month())
		y := t.Year()
		var due, coll float64
		DB.Model(&models.Invoice{}).Where("month = ? AND year = ?", m, y).
			Select("COALESCE(SUM(net_payable), 0)").Scan(&due)
		DB.Model(&models.Invoice{}).Where("month = ? AND year = ?", m, y).
			Select("COALESCE(SUM(paid_amount), 0)").Scan(&coll)
		trend = append(trend, models.MonthlyPremiumTrend{
			Month:     t.Format("Jan 2006"),
			Due:       due,
			Collected: coll,
		})
	}

	// Status breakdown
	statusList := []string{"draft", "sent", "partial", "paid", "overdue"}
	breakdown := make([]models.PremiumStatusBreakdown, 0, len(statusList))
	for _, s := range statusList {
		var amount float64
		var count int64
		DB.Model(&models.Invoice{}).Where("status = ?", s).
			Select("COALESCE(SUM(net_payable), 0)").Scan(&amount)
		DB.Model(&models.Invoice{}).Where("status = ?", s).Count(&count)
		breakdown = append(breakdown, models.PremiumStatusBreakdown{
			Status: s,
			Amount: amount,
			Count:  int(count),
		})
	}

	// Top 10 outstanding schemes
	type schemeBalance struct {
		SchemeID   int
		SchemeName string
		Due        float64
		Paid       float64
		Balance    float64
	}
	var topSchemes []schemeBalance
	DB.Model(&models.Invoice{}).
		Where("balance > 0").
		Select("scheme_id, scheme_name, SUM(net_payable) as due, SUM(paid_amount) as paid, SUM(balance) as balance").
		Group("scheme_id, scheme_name").
		Order("balance DESC").
		Limit(10).Scan(&topSchemes)

	topOutstanding := make([]models.SchemeOutstandingSummary, 0, len(topSchemes))
	for _, ts := range topSchemes {
		topOutstanding = append(topOutstanding, models.SchemeOutstandingSummary{
			SchemeID:   ts.SchemeID,
			SchemeName: ts.SchemeName,
			AmountDue:  ts.Due,
			AmountPaid: ts.Paid,
			Balance:    ts.Balance,
			Status:     "in_arrears",
		})
	}

	return models.PremiumDashboardData{
		KPIs:            kpis,
		MonthlyTrend:    trend,
		StatusBreakdown: breakdown,
		TopOutstanding:  topOutstanding,
	}, nil
}

// GetCollectionRate returns monthly collection rates for a year.
func GetCollectionRate(year int) (models.CollectionRateResponse, error) {
	if year == 0 {
		year = time.Now().Year()
	}

	monthlyRates := make([]models.MonthlyCollectionRate, 0, 12)
	totalDue := 0.0
	totalCollected := 0.0

	for m := 1; m <= 12; m++ {
		var due, collected float64
		DB.Model(&models.Invoice{}).Where("month = ? AND year = ?", m, year).
			Select("COALESCE(SUM(net_payable), 0)").Scan(&due)
		DB.Model(&models.Invoice{}).Where("month = ? AND year = ?", m, year).
			Select("COALESCE(SUM(paid_amount), 0)").Scan(&collected)

		rate := 0.0
		if due > 0 {
			rate = collected / due * 100
		}

		t := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		monthlyRates = append(monthlyRates, models.MonthlyCollectionRate{
			Month: t.Format("Jan 2006"),
			Rate:  rate,
		})
		totalDue += due
		totalCollected += collected
	}

	overallRate := 0.0
	if totalDue > 0 {
		overallRate = totalCollected / totalDue * 100
	}

	return models.CollectionRateResponse{
		Year:           year,
		CollectionRate: overallRate,
		MonthlyRates:   monthlyRates,
	}, nil
}

// IsProRata calculates whether a member is on a pro-rata basis for a given period.
// It supports any billing interval (monthly, quarterly, etc.) by accepting
// the period start and end dates explicitly.
//
// Parameters:
//   - m            : the member in force
//   - year         : the year of the billing period
//   - month        : the month of the billing period (1–12)
//   - daysInMonth  : number of days in the billing month (allows overriding for quarterly etc.)
//
// Returns:
//   - isProRata    : true if the member did not cover the full period
//   - proRataDays  : pointer to the number of active days (nil if full period)
func IsProRata(m models.GPricingMemberDataInForce, year, month, daysInMonth int) (isProRata bool, proRataDays *int) {
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := time.Date(year, time.Month(month), daysInMonth, 0, 0, 0, 0, time.UTC)

	// --- Determine effective start within the period ---
	effectiveStart := monthStart
	if m.EntryDate.After(monthStart) {
		effectiveStart = m.EntryDate
	}

	// --- Determine effective end within the period ---
	effectiveEnd := monthEnd
	if m.ExitDate != nil && m.ExitDate.Before(monthEnd) {
		effectiveEnd = *m.ExitDate
	}

	// --- No overlap: member not active in this period ---
	if effectiveStart.After(effectiveEnd) {
		return false, nil
	}

	// --- Calculate inclusive days active ---
	daysActive := int(effectiveEnd.Sub(effectiveStart).Hours()/24) + 1

	// --- Guard bounds ---
	if daysActive < 0 {
		daysActive = 0
	}
	if daysActive > daysInMonth {
		daysActive = daysInMonth
	}

	// --- Full period: not pro-rata ---
	if daysActive == daysInMonth {
		return false, nil
	}

	// --- Partial period: pro-rata ---
	return true, &daysActive
}

func PremiumComputation(member models.GPricingMemberDataInForce, coveredSA models.SumsAssuredResult, memberratingresultsummary models.MemberRatingResultSummary) models.PremiumComputation {

	var premiums models.PremiumComputation

	premiums.GlaRiskPremium = utils.FloatPrecision(coveredSA.GlaCoveredSumAssured*memberratingresultsummary.ExpGlaRiskRatePer1000SA/1000.0, AccountingPrecision)
	premiums.GlaOfficePremium = utils.FloatPrecision(models.ComputeOfficePremium(premiums.GlaRiskPremium, &memberratingresultsummary), AccountingPrecision)
	premiums.PtdRiskPremium = utils.FloatPrecision(coveredSA.PtdCoveredSumAssured*memberratingresultsummary.ExpPtdRiskRatePer1000SA/1000.0, AccountingPrecision)
	premiums.PtdOfficePremium = utils.FloatPrecision(models.ComputeOfficePremium(premiums.PtdRiskPremium, &memberratingresultsummary), AccountingPrecision)
	premiums.CiRiskPremium = utils.FloatPrecision(coveredSA.CiCoveredSumAssured*memberratingresultsummary.ExpCiRiskRatePer1000SA/1000.0, AccountingPrecision)
	premiums.CiOfficePremium = utils.FloatPrecision(models.ComputeOfficePremium(premiums.CiRiskPremium, &memberratingresultsummary), AccountingPrecision)
	premiums.SglaRiskPremium = utils.FloatPrecision(coveredSA.SglaCoveredSumAssured*memberratingresultsummary.ExpSglaRiskRatePer1000SA/1000.0, AccountingPrecision)
	premiums.SglaOfficePremium = utils.FloatPrecision(models.ComputeOfficePremium(premiums.SglaRiskPremium, &memberratingresultsummary), AccountingPrecision)
	premiums.PhiRiskPremium = utils.FloatPrecision(member.AnnualSalary*memberratingresultsummary.ExpProportionPhiAnnualRiskPremiumSalary, AccountingPrecision)
	premiums.PhiOfficePremium = utils.FloatPrecision(models.ComputeOfficePremium(premiums.PhiRiskPremium, &memberratingresultsummary), AccountingPrecision)
	premiums.TtdRiskPremium = utils.FloatPrecision(member.AnnualSalary*memberratingresultsummary.ExpProportionTtdAnnualRiskPremiumSalary, AccountingPrecision)
	premiums.TtdOfficePremium = utils.FloatPrecision(models.ComputeOfficePremium(premiums.TtdRiskPremium, &memberratingresultsummary), AccountingPrecision)
	premiums.FuneralRiskPremium = utils.FloatPrecision(memberratingresultsummary.ExpTotalFunAnnualRiskPremium, AccountingPrecision)
	premiums.FuneralOfficePremium = utils.FloatPrecision(models.ComputeOfficePremium(memberratingresultsummary.ExpTotalFunAnnualRiskPremium, &memberratingresultsummary), AccountingPrecision)

	premiums.TotalRiskPremiumExclFun = utils.FloatPrecision(premiums.GlaRiskPremium+premiums.PtdRiskPremium+premiums.CiRiskPremium+premiums.SglaRiskPremium+premiums.PhiRiskPremium+premiums.TtdRiskPremium, AccountingPrecision)

	premiums.TotalRiskFuneralPremium = utils.FloatPrecision(memberratingresultsummary.ExpTotalFunAnnualPremiumPerMember, AccountingPrecision)
	premiums.TotalOfficeFuneralPremium = utils.FloatPrecision(memberratingresultsummary.ExpTotalFunAnnualPremiumPerMember, AccountingPrecision)

	premiums.TotalOfficePremiumExclFun = utils.FloatPrecision(premiums.GlaOfficePremium+premiums.PtdOfficePremium+premiums.CiOfficePremium+premiums.SglaOfficePremium+premiums.PhiOfficePremium+premiums.TtdOfficePremium, AccountingPrecision)
	premiums.TotalRiskPremium = utils.FloatPrecision(premiums.TotalRiskPremiumExclFun+premiums.TotalRiskFuneralPremium, AccountingPrecision)
	premiums.TotalOfficePremium = utils.FloatPrecision(premiums.TotalOfficePremiumExclFun+premiums.TotalOfficeFuneralPremium, AccountingPrecision)

	return premiums
}

// GetScheduleCoverageMatrix returns a matrix of all in-force schemes × last 12 months
// showing whether a premium schedule exists for each month.
func GetScheduleCoverageMatrix() (models.ScheduleCoverageMatrix, error) {
	now := time.Now()

	// Build the 12-month window (most recent first)
	type monthYear struct {
		Month int
		Year  int
	}
	window := make([]monthYear, 12)
	for i := 0; i < 12; i++ {
		t := now.AddDate(0, -i, 0)
		window[i] = monthYear{Month: int(t.Month()), Year: t.Year()}
	}

	monthLabels := make([]string, 12)
	for i, my := range window {
		monthLabels[i] = time.Date(my.Year, time.Month(my.Month), 1, 0, 0, 0, 0, time.UTC).Format("Jan 2006")
	}

	// Fetch all in-force schemes
	var schemes []models.GroupScheme
	if err := DB.Where("status = ?", models.StatusInForce).Order("name ASC").Find(&schemes).Error; err != nil {
		return models.ScheduleCoverageMatrix{}, err
	}

	if len(schemes) == 0 {
		return models.ScheduleCoverageMatrix{Months: monthLabels, Schemes: []models.SchemeCoverageRow{}}, nil
	}

	// Collect all scheme IDs
	schemeIDs := make([]int, len(schemes))
	for i, s := range schemes {
		schemeIDs[i] = s.ID
	}

	// Fetch schedules within the window in a single query
	oldest := window[11]
	var schedules []models.PremiumSchedule
	if err := DB.Where("scheme_id IN ? AND status NOT IN ?", schemeIDs, []string{"void", "cancelled"}).
		Where("(year > ? OR (year = ? AND month >= ?))", oldest.Year, oldest.Year, oldest.Month).
		Where("(year < ? OR (year = ? AND month <= ?))", window[0].Year, window[0].Year, window[0].Month).
		Find(&schedules).Error; err != nil {
		return models.ScheduleCoverageMatrix{}, err
	}

	// Index schedules by scheme_id → month/year
	type schedKey struct {
		SchemeID int
		Month    int
		Year     int
	}
	schedMap := make(map[schedKey]models.PremiumSchedule, len(schedules))
	for _, s := range schedules {
		schedMap[schedKey{SchemeID: s.SchemeID, Month: s.Month, Year: s.Year}] = s
	}

	// Build the matrix
	rows := make([]models.SchemeCoverageRow, len(schemes))
	for i, scheme := range schemes {
		cells := make([]models.ScheduleCoverageCell, 12)
		// Determine commencement month boundary (year/month of commencement date)
		commYear := 0
		commMonth := 0
		if !scheme.CommencementDate.IsZero() {
			commYear = scheme.CommencementDate.Year()
			commMonth = int(scheme.CommencementDate.Month())
		}
		for j, my := range window {
			cell := models.ScheduleCoverageCell{Month: my.Month, Year: my.Year}
			// Mark cells that fall before the scheme's commencement date
			if commYear > 0 && (my.Year < commYear || (my.Year == commYear && my.Month < commMonth)) {
				cell.BeforeCommencement = true
			}
			if s, ok := schedMap[schedKey{SchemeID: scheme.ID, Month: my.Month, Year: my.Year}]; ok {
				cell.Exists = true
				cell.ScheduleID = s.ID
				cell.Status = s.Status
			}
			cells[j] = cell
		}
		commDate := ""
		if !scheme.CommencementDate.IsZero() {
			commDate = scheme.CommencementDate.Format("2006-01-02")
		}
		rows[i] = models.SchemeCoverageRow{
			SchemeID:         scheme.ID,
			SchemeName:       scheme.Name,
			CommencementDate: commDate,
			Cells:            cells,
		}
	}

	return models.ScheduleCoverageMatrix{Months: monthLabels, Schemes: rows}, nil
}
func MapMemberToCategory(member models.GPricingMemberDataInForce, categories []models.SchemeCategory) (models.SchemeCategory, error) {
	for _, c := range categories {
		// choose your matching rule; here we match on category
		if c.SchemeCategory == member.SchemeCategory {
			return c, nil
		}
	}
	return models.SchemeCategory{}, nil
}

func MapMemberToMemberRatingResultSummary(member models.GPricingMemberDataInForce, memberRatingResultSummaries []models.MemberRatingResultSummary) (models.MemberRatingResultSummary, error) {
	for _, result := range memberRatingResultSummaries {
		// choose your matching rule; here we match on category
		if result.Category == member.SchemeCategory {
			return result, nil
		}
	}
	return models.MemberRatingResultSummary{}, nil
}

// after loading treatyLinks, extract TreatyIDs
func extractTreatyIDs(treatyLinks []models.TreatySchemeLink) []int {
	ids := make([]int, 0, len(treatyLinks))
	for _, tl := range treatyLinks {
		ids = append(ids, tl.TreatyID)
	}
	return ids
}

// treaties: list of available treaties
// benefits: list of benefits to assign
// Returns the updated slice of benefits and an optional error if a benefit can't be mapped.
func MapBenefitsToTreaties(treaties []models.ReinsuranceTreaty, benefits []models.BenefitTreatyMap) ([]models.BenefitTreatyMap, error) {
	// Build lookup: TypeKey -> TreatyID
	lookup := make(map[string]int, len(treaties))
	for _, tr := range treaties {
		lookup[tr.TreatyCode] = tr.ID
	}

	// Assign
	for i := range benefits {
		b := &benefits[i]
		tid, ok := lookup[b.Name]
		if !ok {
			continue
		}
		b.TreatyID = tid
	}

	return benefits, nil
}

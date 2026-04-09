package services

import (
	appLog "api/log"
	"api/models"
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Auto-Match Engine
// ---------------------------------------------------------------------------

// RunAutoMatch executes the multi-strategy matching engine. It creates a
// ReconciliationRun, applies each active MatchingRule in priority order,
// and allocates payments to invoices. Supports dry-run mode for previews.
func RunAutoMatch(req models.RunAutoMatchRequest, user models.AppUser) (interface{}, error) {
	// Load active rules
	rules, err := getActiveRules(req.RuleSet)
	if err != nil {
		return nil, fmt.Errorf("failed to load matching rules: %w", err)
	}

	// Load unallocated payment items
	paymentItems, err := getOpenReconciliationItems("payment", req.SchemeID)
	if err != nil {
		return nil, fmt.Errorf("failed to load payment items: %w", err)
	}

	// Load unpaid invoice items
	invoiceItems, err := getOpenReconciliationItems("invoice", req.SchemeID)
	if err != nil {
		return nil, fmt.Errorf("failed to load invoice items: %w", err)
	}

	if len(paymentItems) == 0 || len(invoiceItems) == 0 {
		if req.DryRun {
			return models.AutoMatchPreview{}, nil
		}
		return models.ReconciliationRun{
			RunDate: time.Now().Format("2006-01-02"),
			RunType: "auto",
			Status:  "completed",
			Notes:   "No items to match",
		}, nil
	}

	// Load full payment and invoice records for matching
	payments, err := loadPaymentsForItems(paymentItems)
	if err != nil {
		return nil, err
	}
	invoices, err := loadInvoicesForItems(invoiceItems)
	if err != nil {
		return nil, err
	}

	// Build lookup maps
	paymentItemMap := make(map[int]*models.ReconciliationItem) // paymentID → item
	for i := range paymentItems {
		if paymentItems[i].PaymentID != nil {
			paymentItemMap[*paymentItems[i].PaymentID] = &paymentItems[i]
		}
	}
	invoiceItemMap := make(map[int]*models.ReconciliationItem) // invoiceID → item
	for i := range invoiceItems {
		if invoiceItems[i].InvoiceID != nil {
			invoiceItemMap[*invoiceItems[i].InvoiceID] = &invoiceItems[i]
		}
	}

	// Execute rules in priority order, collecting proposed allocations
	var proposed []models.ProposedAllocation
	matchedPayments := make(map[int]float64) // paymentID → amount already proposed
	matchedInvoices := make(map[int]float64) // invoiceID → amount already proposed

	for _, rule := range rules {
		matches := applyRule(rule, payments, invoices, paymentItemMap, invoiceItemMap, matchedPayments, matchedInvoices)
		proposed = append(proposed, matches...)
	}

	// Dry-run: return preview without persisting
	if req.DryRun {
		var totalAlloc float64
		for _, p := range proposed {
			totalAlloc += p.Amount
		}
		remaining := 0
		for _, pi := range paymentItems {
			if pi.PaymentID != nil {
				if alloc, ok := matchedPayments[*pi.PaymentID]; !ok || alloc < pi.UnallocatedAmount {
					remaining++
				}
			}
		}
		return models.AutoMatchPreview{
			ProposedAllocations: proposed,
			TotalAllocated:      totalAlloc,
			TotalMatched:        len(proposed),
			TotalRemaining:      remaining,
		}, nil
	}

	// Persist: create run and allocations in a transaction
	run := models.ReconciliationRun{
		RunDate:         time.Now().Format("2006-01-02"),
		RunType:         "auto",
		Status:          "in_progress",
		InitiatedBy:     user.UserName,
		TotalProcessed:  len(paymentItems),
		MatchingRuleSet: req.RuleSet,
	}
	if run.MatchingRuleSet == "" {
		run.MatchingRuleSet = "default"
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&run).Error; err != nil {
			return err
		}

		matchedCount := 0
		var totalAllocated float64

		for _, prop := range proposed {
			alloc := models.PaymentAllocation{
				PaymentID:       prop.PaymentID,
				InvoiceID:       prop.InvoiceID,
				RunID:           &run.ID,
				AllocatedAmount: prop.Amount,
				AllocationType:  "payment",
				Reference:       prop.MatchedBy,
				AllocatedBy:     user.UserName,
			}
			if err := tx.Create(&alloc).Error; err != nil {
				return err
			}

			// Update payment's reconciliation item
			if pi, ok := paymentItemMap[prop.PaymentID]; ok {
				pi.AllocatedAmount += prop.Amount
				pi.UnallocatedAmount = pi.OriginalAmount - pi.AllocatedAmount
				if pi.UnallocatedAmount <= 0.005 {
					pi.Status = "matched"
					pi.UnallocatedAmount = 0
				} else {
					pi.Status = "partial"
				}
				now := time.Now()
				pi.LastActionDate = &now
				if err := tx.Save(pi).Error; err != nil {
					return err
				}
			}

			// Update invoice's reconciliation item
			if ii, ok := invoiceItemMap[prop.InvoiceID]; ok {
				ii.AllocatedAmount += prop.Amount
				ii.UnallocatedAmount = ii.OriginalAmount - ii.AllocatedAmount
				if ii.UnallocatedAmount <= 0.005 {
					ii.Status = "matched"
					ii.UnallocatedAmount = 0
				} else {
					ii.Status = "partial"
				}
				now := time.Now()
				ii.LastActionDate = &now
				if err := tx.Save(ii).Error; err != nil {
					return err
				}
			}

			// Update the underlying Invoice record
			if err := applyAllocationToInvoice(tx, prop.InvoiceID, prop.Amount); err != nil {
				return err
			}

			// Update the underlying Payment record
			if err := updatePaymentStatus(tx, prop.PaymentID, paymentItemMap); err != nil {
				return err
			}

			matchedCount++
			totalAllocated += prop.Amount
		}

		// Finalize run
		now := time.Now()
		run.Status = "completed"
		run.CompletedAt = &now
		run.TotalMatched = matchedCount
		run.TotalAllocated = totalAllocated

		var remaining int64
		tx.Model(&models.ReconciliationItem{}).
			Where("item_type = ? AND status IN ?", "payment", []string{"open", "partial"}).
			Count(&remaining)
		run.TotalUnmatched = int(remaining)

		return tx.Save(&run).Error
	})

	if err != nil {
		return nil, err
	}

	return run, nil
}

// ---------------------------------------------------------------------------
// Manual Allocation
// ---------------------------------------------------------------------------

// AllocatePayment manually allocates a payment to one or more invoices.
func AllocatePayment(req models.AllocatePaymentRequest, user models.AppUser) ([]models.PaymentAllocation, error) {
	var payment models.Payment
	if err := DB.First(&payment, req.PaymentID).Error; err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// Verify that total allocation doesn't exceed available
	paymentItem, err := getOrCreatePaymentItem(payment)
	if err != nil {
		return nil, err
	}

	var totalRequested float64
	for _, a := range req.Allocations {
		totalRequested += a.Amount
	}
	if totalRequested > paymentItem.UnallocatedAmount+0.005 {
		return nil, fmt.Errorf("allocation total %.2f exceeds unallocated balance %.2f", totalRequested, paymentItem.UnallocatedAmount)
	}

	var allocations []models.PaymentAllocation

	err = DB.Transaction(func(tx *gorm.DB) error {
		for _, line := range req.Allocations {
			var invoice models.Invoice
			if err := tx.First(&invoice, line.InvoiceID).Error; err != nil {
				return fmt.Errorf("invoice %d not found: %w", line.InvoiceID, err)
			}

			if line.Amount <= 0 {
				return fmt.Errorf("allocation amount must be positive")
			}

			alloc := models.PaymentAllocation{
				PaymentID:       req.PaymentID,
				InvoiceID:       line.InvoiceID,
				AllocatedAmount: line.Amount,
				AllocationType:  "payment",
				Notes:           req.Notes,
				AllocatedBy:     user.UserName,
			}
			if err := tx.Create(&alloc).Error; err != nil {
				return err
			}
			allocations = append(allocations, alloc)

			// Update the reconciliation items
			if err := applyAllocationToInvoice(tx, line.InvoiceID, line.Amount); err != nil {
				return err
			}

			invoiceItem, err := getOrCreateInvoiceItemTx(tx, invoice)
			if err != nil {
				return err
			}
			invoiceItem.AllocatedAmount += line.Amount
			invoiceItem.UnallocatedAmount = invoiceItem.OriginalAmount - invoiceItem.AllocatedAmount
			if invoiceItem.UnallocatedAmount <= 0.005 {
				invoiceItem.Status = "matched"
				invoiceItem.UnallocatedAmount = 0
			} else {
				invoiceItem.Status = "partial"
			}
			now := time.Now()
			invoiceItem.LastActionDate = &now
			if err := tx.Save(&invoiceItem).Error; err != nil {
				return err
			}
		}

		// Update payment reconciliation item
		paymentItem.AllocatedAmount += totalRequested
		paymentItem.UnallocatedAmount = paymentItem.OriginalAmount - paymentItem.AllocatedAmount
		if paymentItem.UnallocatedAmount <= 0.005 {
			paymentItem.Status = "matched"
			paymentItem.UnallocatedAmount = 0
		} else {
			paymentItem.Status = "partial"
		}
		now := time.Now()
		paymentItem.LastActionDate = &now
		if err := tx.Save(&paymentItem).Error; err != nil {
			return err
		}

		// Update underlying payment status
		return updatePaymentStatus(tx, req.PaymentID, nil)
	})

	return allocations, err
}

// ---------------------------------------------------------------------------
// Allocation Reversal
// ---------------------------------------------------------------------------

// ReverseAllocations reverses one or more allocations, creating counter-entries.
func ReverseAllocations(req models.ReverseAllocationRequest, user models.AppUser) ([]models.PaymentAllocation, error) {
	var reversals []models.PaymentAllocation

	err := DB.Transaction(func(tx *gorm.DB) error {
		for _, allocID := range req.AllocationIDs {
			var original models.PaymentAllocation
			if err := tx.First(&original, allocID).Error; err != nil {
				return fmt.Errorf("allocation %d not found: %w", allocID, err)
			}
			if original.IsReversal {
				return fmt.Errorf("allocation %d is itself a reversal and cannot be reversed", allocID)
			}
			if original.ReversedByID != nil {
				return fmt.Errorf("allocation %d has already been reversed", allocID)
			}

			reversal := models.PaymentAllocation{
				PaymentID:       original.PaymentID,
				InvoiceID:       original.InvoiceID,
				RunID:           original.RunID,
				AllocatedAmount: -original.AllocatedAmount,
				AllocationType:  "reversal",
				Reference:       fmt.Sprintf("Reversal of allocation #%d", allocID),
				Notes:           req.Reason,
				AllocatedBy:     user.UserName,
				IsReversal:      true,
			}
			if err := tx.Create(&reversal).Error; err != nil {
				return err
			}

			// Link original to its reversal
			original.ReversedByID = &reversal.ID
			if err := tx.Save(&original).Error; err != nil {
				return err
			}

			// Reverse the invoice impact
			if err := applyAllocationToInvoice(tx, original.InvoiceID, -original.AllocatedAmount); err != nil {
				return err
			}

			// Update reconciliation items for both sides
			if err := recalcReconciliationItem(tx, "payment", original.PaymentID); err != nil {
				return err
			}
			if err := recalcReconciliationItem(tx, "invoice", original.InvoiceID); err != nil {
				return err
			}

			reversals = append(reversals, reversal)
		}

		return nil
	})

	return reversals, err
}

// ---------------------------------------------------------------------------
// Write-Off
// ---------------------------------------------------------------------------

// WriteOffBalance writes off a small remaining balance on an invoice.
func WriteOffBalance(req models.WriteOffRequest, user models.AppUser) (models.PaymentAllocation, error) {
	var item models.ReconciliationItem
	if err := DB.First(&item, req.ReconciliationItemID).Error; err != nil {
		return models.PaymentAllocation{}, fmt.Errorf("reconciliation item not found: %w", err)
	}

	if req.Amount > item.UnallocatedAmount+0.005 {
		return models.PaymentAllocation{}, fmt.Errorf("write-off amount %.2f exceeds unallocated balance %.2f", req.Amount, item.UnallocatedAmount)
	}

	var alloc models.PaymentAllocation
	err := DB.Transaction(func(tx *gorm.DB) error {
		invoiceID := req.InvoiceID
		if invoiceID == 0 && item.InvoiceID != nil {
			invoiceID = *item.InvoiceID
		}
		paymentID := 0
		if item.PaymentID != nil {
			paymentID = *item.PaymentID
		}

		alloc = models.PaymentAllocation{
			PaymentID:       paymentID,
			InvoiceID:       invoiceID,
			AllocatedAmount: req.Amount,
			AllocationType:  "write_off",
			Reference:       "Write-off",
			Notes:           req.Reason,
			AllocatedBy:     user.UserName,
		}
		if err := tx.Create(&alloc).Error; err != nil {
			return err
		}

		// Update reconciliation item
		item.AllocatedAmount += req.Amount
		item.UnallocatedAmount = item.OriginalAmount - item.AllocatedAmount
		if item.UnallocatedAmount <= 0.005 {
			item.Status = "written_off"
			item.UnallocatedAmount = 0
		}
		now := time.Now()
		item.LastActionDate = &now
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		// Update invoice if applicable
		if invoiceID > 0 {
			return applyAllocationToInvoice(tx, invoiceID, req.Amount)
		}
		return nil
	})

	return alloc, err
}

// ---------------------------------------------------------------------------
// Refund
// ---------------------------------------------------------------------------

// RefundOverpayment records a refund for an overpaid amount.
func RefundOverpayment(req models.RefundRequest, user models.AppUser) (models.PaymentAllocation, error) {
	var item models.ReconciliationItem
	if err := DB.First(&item, req.ReconciliationItemID).Error; err != nil {
		return models.PaymentAllocation{}, fmt.Errorf("reconciliation item not found: %w", err)
	}

	if item.ItemType != "payment" {
		return models.PaymentAllocation{}, errors.New("refunds can only be issued against payment items")
	}

	if req.Amount > item.UnallocatedAmount+0.005 {
		return models.PaymentAllocation{}, fmt.Errorf("refund amount %.2f exceeds unallocated balance %.2f", req.Amount, item.UnallocatedAmount)
	}

	var alloc models.PaymentAllocation
	err := DB.Transaction(func(tx *gorm.DB) error {
		paymentID := 0
		if item.PaymentID != nil {
			paymentID = *item.PaymentID
		}

		alloc = models.PaymentAllocation{
			PaymentID:       paymentID,
			AllocatedAmount: req.Amount,
			AllocationType:  "refund",
			Reference:       fmt.Sprintf("Refund via %s", req.RefundMethod),
			Notes:           req.Reason,
			AllocatedBy:     user.UserName,
		}
		if err := tx.Create(&alloc).Error; err != nil {
			return err
		}

		item.AllocatedAmount += req.Amount
		item.UnallocatedAmount = item.OriginalAmount - item.AllocatedAmount
		if item.UnallocatedAmount <= 0.005 {
			item.Status = "refunded"
			item.UnallocatedAmount = 0
		}
		now := time.Now()
		item.LastActionDate = &now
		return tx.Save(&item).Error
	})

	return alloc, err
}

// ---------------------------------------------------------------------------
// Query Functions
// ---------------------------------------------------------------------------

// GetReconciliationSummary returns the dashboard summary for reconciliation.
func GetReconciliationSummary() (models.ReconciliationSummary, error) {
	var summary models.ReconciliationSummary

	// Unallocated payments
	DB.Model(&models.ReconciliationItem{}).
		Where("item_type = ? AND status IN ?", "payment", []string{"open", "partial"}).
		Select("COALESCE(SUM(unallocated_amount), 0)").
		Scan(&summary.TotalUnallocatedPayments)
	DB.Model(&models.ReconciliationItem{}).
		Where("item_type = ? AND status IN ?", "payment", []string{"open", "partial"}).
		Count(new(int64)).Scan(&summary.UnallocatedPaymentCount)

	// Unpaid invoices
	DB.Model(&models.ReconciliationItem{}).
		Where("item_type = ? AND status IN ?", "invoice", []string{"open", "partial"}).
		Select("COALESCE(SUM(unallocated_amount), 0)").
		Scan(&summary.TotalUnpaidInvoices)
	DB.Model(&models.ReconciliationItem{}).
		Where("item_type = ? AND status IN ?", "invoice", []string{"open", "partial"}).
		Count(new(int64)).Scan(&summary.UnpaidInvoiceCount)

	// Suspense items
	DB.Model(&models.ReconciliationItem{}).
		Where("status = ?", "suspended").
		Select("COALESCE(SUM(unallocated_amount), 0)").
		Scan(&summary.SuspenseBalance)
	var suspenseCount int64
	DB.Model(&models.ReconciliationItem{}).
		Where("status = ?", "suspended").
		Count(&suspenseCount)
	summary.SuspenseCount = int(suspenseCount)

	// Aging
	now := time.Now()
	d30 := now.AddDate(0, 0, -30).Format("2006-01-02")
	d60 := now.AddDate(0, 0, -60).Format("2006-01-02")
	d90 := now.AddDate(0, 0, -90).Format("2006-01-02")

	DB.Model(&models.ReconciliationItem{}).
		Where("status IN ? AND created_at < ?", []string{"open", "partial"}, d30).
		Select("COALESCE(SUM(unallocated_amount), 0)").
		Scan(&summary.AgedOver30Days)
	DB.Model(&models.ReconciliationItem{}).
		Where("status IN ? AND created_at < ?", []string{"open", "partial"}, d60).
		Select("COALESCE(SUM(unallocated_amount), 0)").
		Scan(&summary.AgedOver60Days)
	DB.Model(&models.ReconciliationItem{}).
		Where("status IN ? AND created_at < ?", []string{"open", "partial"}, d90).
		Select("COALESCE(SUM(unallocated_amount), 0)").
		Scan(&summary.AgedOver90Days)

	// Recent runs
	DB.Model(&models.ReconciliationRun{}).
		Order("created_at DESC").
		Limit(10).
		Find(&summary.RecentRuns)

	return summary, nil
}

// GetReconciliationItems returns filtered reconciliation items.
func GetReconciliationItems(itemType, status string, schemeID int) ([]models.ReconciliationItem, error) {
	query := DB.Model(&models.ReconciliationItem{})
	if itemType != "" {
		query = query.Where("item_type = ?", itemType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if schemeID > 0 {
		query = query.Where("scheme_id = ?", schemeID)
	}

	var items []models.ReconciliationItem
	err := query.Order("created_at DESC").Find(&items).Error

	// Compute age for each item
	now := time.Now()
	for i := range items {
		items[i].AgeInDays = int(now.Sub(items[i].CreatedAt).Hours() / 24)
	}

	return items, err
}

// GetAllocationHistory returns the full allocation trail for a payment or invoice.
func GetAllocationHistory(entityType string, entityID int) (models.AllocationHistory, error) {
	history := models.AllocationHistory{
		EntityType: entityType,
		EntityID:   entityID,
	}

	query := DB.Model(&models.PaymentAllocation{})
	if entityType == "payment" {
		query = query.Where("payment_id = ?", entityID)
	} else {
		query = query.Where("invoice_id = ?", entityID)
	}

	if err := query.Order("created_at ASC").Find(&history.Allocations).Error; err != nil {
		return history, err
	}

	for _, a := range history.Allocations {
		history.Allocated += a.AllocatedAmount
	}

	// Get original total from the reconciliation item
	var item models.ReconciliationItem
	if entityType == "payment" {
		if err := DB.Where("payment_id = ?", entityID).First(&item).Error; err == nil {
			history.Total = item.OriginalAmount
		}
	} else {
		if err := DB.Where("invoice_id = ?", entityID).First(&item).Error; err == nil {
			history.Total = item.OriginalAmount
		}
	}
	history.Unallocated = history.Total - history.Allocated

	return history, nil
}

// GetReconciliationRunDetail returns a run with its allocations.
func GetReconciliationRunDetail(runID int) (models.ReconciliationRunDetail, error) {
	var detail models.ReconciliationRunDetail
	if err := DB.First(&detail.ReconciliationRun, runID).Error; err != nil {
		return detail, fmt.Errorf("run not found: %w", err)
	}

	DB.Where("run_id = ?", runID).Order("created_at ASC").Find(&detail.Allocations)
	return detail, nil
}

// GetReconciliationRuns returns a paginated list of runs.
func GetReconciliationRuns(page, pageSize int) ([]models.ReconciliationRun, int64, error) {
	var runs []models.ReconciliationRun
	var total int64

	DB.Model(&models.ReconciliationRun{}).Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	err := DB.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&runs).Error

	return runs, total, err
}

// RollbackRun reverses all allocations in a run.
func RollbackRun(runID int, user models.AppUser) error {
	var run models.ReconciliationRun
	if err := DB.First(&run, runID).Error; err != nil {
		return fmt.Errorf("run not found: %w", err)
	}
	if run.Status == "rolled_back" {
		return errors.New("run has already been rolled back")
	}

	var allocations []models.PaymentAllocation
	DB.Where("run_id = ? AND is_reversal = ?", runID, false).Find(&allocations)

	ids := make([]int, 0, len(allocations))
	for _, a := range allocations {
		if a.ReversedByID == nil {
			ids = append(ids, a.ID)
		}
	}

	if len(ids) == 0 {
		run.Status = "rolled_back"
		return DB.Save(&run).Error
	}

	_, err := ReverseAllocations(models.ReverseAllocationRequest{
		AllocationIDs: ids,
		Reason:        fmt.Sprintf("Rollback of run #%d", runID),
	}, user)
	if err != nil {
		return err
	}

	run.Status = "rolled_back"
	return DB.Save(&run).Error
}

// ReassignReconciliationItem updates the assignment and priority of a reconciliation item.
func ReassignReconciliationItem(itemID int, req models.ReassignItemRequest) error {
	updates := map[string]interface{}{
		"assigned_to": req.AssignedTo,
	}
	if req.Priority != "" {
		updates["priority"] = req.Priority
	}
	return DB.Model(&models.ReconciliationItem{}).Where("id = ?", itemID).Updates(updates).Error
}

// SuspendReconciliationItem moves an item to suspended status.
func SuspendReconciliationItem(itemID int, reason string) error {
	return DB.Model(&models.ReconciliationItem{}).Where("id = ?", itemID).
		Updates(map[string]interface{}{
			"status":          "suspended",
			"suspense_reason": reason,
		}).Error
}

// ---------------------------------------------------------------------------
// Matching Rule Management
// ---------------------------------------------------------------------------

// GetMatchingRules returns all matching rules for a rule set.
func GetMatchingRules(ruleSet string) ([]models.MatchingRule, error) {
	if ruleSet == "" {
		ruleSet = "default"
	}
	var rules []models.MatchingRule
	err := DB.Where("rule_set = ?", ruleSet).Order("priority ASC").Find(&rules).Error
	return rules, err
}

// SaveMatchingRule creates or updates a matching rule.
func SaveMatchingRule(rule models.MatchingRule) (models.MatchingRule, error) {
	if rule.RuleSet == "" {
		rule.RuleSet = "default"
	}
	if rule.ID > 0 {
		err := DB.Save(&rule).Error
		return rule, err
	}
	err := DB.Create(&rule).Error
	return rule, err
}

// DeleteMatchingRule deletes a matching rule.
func DeleteMatchingRule(ruleID int) error {
	return DB.Delete(&models.MatchingRule{}, ruleID).Error
}

// SeedDefaultMatchingRules creates the default matching rules if none exist.
func SeedDefaultMatchingRules() {
	var count int64
	DB.Model(&models.MatchingRule{}).Where("rule_set = ?", "default").Count(&count)
	if count > 0 {
		return
	}

	defaults := []models.MatchingRule{
		{
			RuleSet:      "default",
			Priority:     1,
			Name:         "Exact Reference Match",
			Description:  "Match payment bank reference to invoice number",
			Strategy:     "exact_reference",
			IsActive:     true,
			AllowPartial: true,
			AllowMulti:   false,
		},
		{
			RuleSet:       "default",
			Priority:      2,
			Name:          "Scheme + Exact Amount",
			Description:   "Match by scheme and exact amount within tolerance",
			Strategy:      "scheme_amount",
			ToleranceType: "absolute",
			ToleranceVal:  0.01,
			IsActive:      true,
			AllowPartial:  true,
			AllowMulti:    true,
		},
		{
			RuleSet:       "default",
			Priority:      3,
			Name:          "Scheme + Amount Tolerance",
			Description:   "Match by scheme with configurable percentage tolerance",
			Strategy:      "scheme_amount_tolerance",
			ToleranceType: "percentage",
			ToleranceVal:  1.0,
			IsActive:      true,
			AllowPartial:  true,
			AllowMulti:    true,
		},
	}

	for _, r := range defaults {
		DB.Create(&r)
	}
}

// ---------------------------------------------------------------------------
// Backfill — sync existing invoices & payments into reconciliation items
// ---------------------------------------------------------------------------

// BackfillReconciliationItems creates ReconciliationItem entries for all existing
// invoices and payments that don't already have one. This bridges pre-v2 data
// into the new ledger-based reconciliation system. Safe to call repeatedly.
func BackfillReconciliationItems() {
	// Backfill invoices
	var invoices []models.Invoice
	if err := DB.Find(&invoices).Error; err != nil {
		appLog.WithField("error", err.Error()).Error("Backfill: failed to load invoices")
		return
	}

	invoiceCount := 0
	for _, inv := range invoices {
		var existing models.ReconciliationItem
		err := DB.Where("invoice_id = ?", inv.ID).First(&existing).Error
		if err == nil {
			continue // already exists
		}

		item := models.ReconciliationItem{
			ItemType:          "invoice",
			InvoiceID:         &inv.ID,
			SchemeID:          inv.SchemeID,
			SchemeName:        inv.SchemeName,
			OriginalAmount:    inv.NetPayable,
			AllocatedAmount:   inv.PaidAmount,
			UnallocatedAmount: inv.Balance,
			Status:            "open",
		}
		if inv.Status == "paid" {
			item.Status = "matched"
		} else if inv.PaidAmount > 0 {
			item.Status = "partial"
		}

		if err := DB.Create(&item).Error; err == nil {
			invoiceCount++
		}
	}

	// Backfill payments
	var payments []models.Payment
	if err := DB.Find(&payments).Error; err != nil {
		appLog.WithField("error", err.Error()).Error("Backfill: failed to load payments")
		return
	}

	paymentCount := 0
	for _, p := range payments {
		if p.Status == "voided" {
			continue
		}

		var existing models.ReconciliationItem
		err := DB.Where("payment_id = ?", p.ID).First(&existing).Error
		if err == nil {
			continue // already exists
		}

		item := models.ReconciliationItem{
			ItemType:          "payment",
			PaymentID:         &p.ID,
			SchemeID:          p.SchemeID,
			SchemeName:        p.SchemeName,
			OriginalAmount:    p.Amount,
			AllocatedAmount:   0,
			UnallocatedAmount: p.Amount,
			Status:            "open",
		}

		if p.Status == "matched" {
			item.AllocatedAmount = p.Amount
			item.UnallocatedAmount = 0
			item.Status = "matched"
		}

		if err := DB.Create(&item).Error; err == nil {
			paymentCount++
		}
	}

	// Also backfill allocation records for payments that were already matched
	// via the legacy 1:1 system (payment.invoice_id is set)
	for _, p := range payments {
		if p.Status != "matched" || p.InvoiceID == nil {
			continue
		}

		// Check if allocation already exists
		var allocCount int64
		DB.Model(&models.PaymentAllocation{}).
			Where("payment_id = ? AND invoice_id = ?", p.ID, *p.InvoiceID).
			Count(&allocCount)
		if allocCount > 0 {
			continue
		}

		alloc := models.PaymentAllocation{
			PaymentID:       p.ID,
			InvoiceID:       *p.InvoiceID,
			AllocatedAmount: p.Amount,
			AllocationType:  "payment",
			Reference:       "Backfilled from legacy match",
			AllocatedBy:     "system",
		}
		DB.Create(&alloc)
	}

	if invoiceCount > 0 || paymentCount > 0 {
		appLog.Infof("Reconciliation backfill: created %d invoice items, %d payment items", invoiceCount, paymentCount)
	}
}

// ---------------------------------------------------------------------------
// Reconciliation Item Lifecycle  —  called from existing payment/invoice code
// ---------------------------------------------------------------------------

// EnsureReconciliationItemForPayment creates or updates a reconciliation item
// when a payment is recorded. Called from the existing RecordPayment flow.
func EnsureReconciliationItemForPayment(payment models.Payment) error {
	var item models.ReconciliationItem
	err := DB.Where("payment_id = ?", payment.ID).First(&item).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		item = models.ReconciliationItem{
			ItemType:          "payment",
			PaymentID:         &payment.ID,
			SchemeID:          payment.SchemeID,
			SchemeName:        payment.SchemeName,
			OriginalAmount:    payment.Amount,
			AllocatedAmount:   0,
			UnallocatedAmount: payment.Amount,
			Status:            "open",
		}
		if payment.Status == "matched" {
			item.AllocatedAmount = payment.Amount
			item.UnallocatedAmount = 0
			item.Status = "matched"
		}
		return DB.Create(&item).Error
	}

	return nil // already exists
}

// EnsureReconciliationItemForInvoice creates or updates a reconciliation item
// when an invoice is generated. Called from the existing GenerateInvoice flow.
func EnsureReconciliationItemForInvoice(invoice models.Invoice) error {
	var item models.ReconciliationItem
	err := DB.Where("invoice_id = ?", invoice.ID).First(&item).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		item = models.ReconciliationItem{
			ItemType:          "invoice",
			InvoiceID:         &invoice.ID,
			SchemeID:          invoice.SchemeID,
			SchemeName:        invoice.SchemeName,
			OriginalAmount:    invoice.NetPayable,
			AllocatedAmount:   invoice.PaidAmount,
			UnallocatedAmount: invoice.Balance,
			Status:            "open",
		}
		if invoice.Status == "paid" {
			item.Status = "matched"
		} else if invoice.PaidAmount > 0 {
			item.Status = "partial"
		}
		return DB.Create(&item).Error
	}

	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func getActiveRules(ruleSet string) ([]models.MatchingRule, error) {
	if ruleSet == "" {
		ruleSet = "default"
	}
	var rules []models.MatchingRule
	err := DB.Where("rule_set = ? AND is_active = ?", ruleSet, true).
		Order("priority ASC").Find(&rules).Error
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		// Seed defaults and retry
		SeedDefaultMatchingRules()
		err = DB.Where("rule_set = ? AND is_active = ?", ruleSet, true).
			Order("priority ASC").Find(&rules).Error
	}
	return rules, err
}

func getOpenReconciliationItems(itemType string, schemeID int) ([]models.ReconciliationItem, error) {
	query := DB.Where("item_type = ? AND status IN ?", itemType, []string{"open", "partial"})
	if schemeID > 0 {
		query = query.Where("scheme_id = ?", schemeID)
	}
	var items []models.ReconciliationItem
	err := query.Find(&items).Error
	return items, err
}

func loadPaymentsForItems(items []models.ReconciliationItem) ([]models.Payment, error) {
	ids := make([]int, 0, len(items))
	for _, i := range items {
		if i.PaymentID != nil {
			ids = append(ids, *i.PaymentID)
		}
	}
	if len(ids) == 0 {
		return nil, nil
	}
	var payments []models.Payment
	err := DB.Where("id IN ?", ids).Find(&payments).Error
	return payments, err
}

func loadInvoicesForItems(items []models.ReconciliationItem) ([]models.Invoice, error) {
	ids := make([]int, 0, len(items))
	for _, i := range items {
		if i.InvoiceID != nil {
			ids = append(ids, *i.InvoiceID)
		}
	}
	if len(ids) == 0 {
		return nil, nil
	}
	var invoices []models.Invoice
	err := DB.Where("id IN ? AND status IN ?", ids, []string{"sent", "partial"}).Find(&invoices).Error
	return invoices, err
}

func applyRule(
	rule models.MatchingRule,
	payments []models.Payment,
	invoices []models.Invoice,
	paymentItemMap map[int]*models.ReconciliationItem,
	invoiceItemMap map[int]*models.ReconciliationItem,
	matchedPayments map[int]float64,
	matchedInvoices map[int]float64,
) []models.ProposedAllocation {
	var results []models.ProposedAllocation

	switch rule.Strategy {
	case "exact_reference":
		results = matchByExactReference(rule, payments, invoices, paymentItemMap, invoiceItemMap, matchedPayments, matchedInvoices)
	case "scheme_amount":
		results = matchBySchemeAmount(rule, payments, invoices, paymentItemMap, invoiceItemMap, matchedPayments, matchedInvoices)
	case "scheme_amount_tolerance":
		results = matchBySchemeAmountTolerance(rule, payments, invoices, paymentItemMap, invoiceItemMap, matchedPayments, matchedInvoices)
	}

	return results
}

func matchByExactReference(
	rule models.MatchingRule,
	payments []models.Payment,
	invoices []models.Invoice,
	paymentItemMap map[int]*models.ReconciliationItem,
	invoiceItemMap map[int]*models.ReconciliationItem,
	matchedPayments map[int]float64,
	matchedInvoices map[int]float64,
) []models.ProposedAllocation {
	var results []models.ProposedAllocation

	// Build invoice lookup by invoice_number
	invByNumber := make(map[string]models.Invoice)
	for _, inv := range invoices {
		invByNumber[inv.InvoiceNumber] = inv
	}

	for _, p := range payments {
		if p.BankReference == "" {
			continue
		}
		pi := paymentItemMap[p.ID]
		if pi == nil {
			continue
		}
		paymentRemaining := pi.UnallocatedAmount - matchedPayments[p.ID]
		if paymentRemaining <= 0.005 {
			continue
		}

		inv, found := invByNumber[p.BankReference]
		if !found {
			continue
		}

		invItem, exists := invoiceItemMap[inv.ID]
		if !exists {
			continue
		}
		invoiceRemaining := invItem.UnallocatedAmount - matchedInvoices[inv.ID]
		if invoiceRemaining <= 0.005 {
			continue
		}

		allocAmt := math.Min(paymentRemaining, invoiceRemaining)

		results = append(results, models.ProposedAllocation{
			PaymentID:     p.ID,
			InvoiceID:     inv.ID,
			Amount:        allocAmt,
			MatchedBy:     rule.Name,
			Confidence:    "high",
			PaymentRef:    p.BankReference,
			InvoiceNumber: inv.InvoiceNumber,
			SchemeName:    inv.SchemeName,
		})

		matchedPayments[p.ID] += allocAmt
		matchedInvoices[inv.ID] += allocAmt
	}

	return results
}

func matchBySchemeAmount(
	rule models.MatchingRule,
	payments []models.Payment,
	invoices []models.Invoice,
	paymentItemMap map[int]*models.ReconciliationItem,
	invoiceItemMap map[int]*models.ReconciliationItem,
	matchedPayments map[int]float64,
	matchedInvoices map[int]float64,
) []models.ProposedAllocation {
	var results []models.ProposedAllocation

	// Group invoices by scheme
	invByScheme := make(map[int][]models.Invoice)
	for _, inv := range invoices {
		invByScheme[inv.SchemeID] = append(invByScheme[inv.SchemeID], inv)
	}

	for _, p := range payments {
		pi := paymentItemMap[p.ID]
		if pi == nil {
			continue
		}
		paymentRemaining := pi.UnallocatedAmount - matchedPayments[p.ID]
		if paymentRemaining <= 0.005 {
			continue
		}

		schemeInvoices := invByScheme[p.SchemeID]
		if len(schemeInvoices) == 0 {
			continue
		}

		tolerance := rule.ToleranceVal

		if rule.AllowMulti {
			// Try to allocate across multiple invoices (oldest first)
			results = allocateAcrossInvoices(rule, p, schemeInvoices, invoiceItemMap, matchedPayments, matchedInvoices, paymentRemaining, tolerance, results)
		} else {
			// Single invoice exact match
			for _, inv := range schemeInvoices {
				invItem, exists := invoiceItemMap[inv.ID]
				if !exists {
					continue
				}
				invoiceRemaining := invItem.UnallocatedAmount - matchedInvoices[inv.ID]
				if invoiceRemaining <= 0.005 {
					continue
				}

				if math.Abs(paymentRemaining-invoiceRemaining) <= tolerance {
					allocAmt := math.Min(paymentRemaining, invoiceRemaining)
					results = append(results, models.ProposedAllocation{
						PaymentID:     p.ID,
						InvoiceID:     inv.ID,
						Amount:        allocAmt,
						MatchedBy:     rule.Name,
						Confidence:    "high",
						PaymentRef:    p.BankReference,
						InvoiceNumber: inv.InvoiceNumber,
						SchemeName:    inv.SchemeName,
					})
					matchedPayments[p.ID] += allocAmt
					matchedInvoices[inv.ID] += allocAmt
					break
				}
			}
		}
	}

	return results
}

func matchBySchemeAmountTolerance(
	rule models.MatchingRule,
	payments []models.Payment,
	invoices []models.Invoice,
	paymentItemMap map[int]*models.ReconciliationItem,
	invoiceItemMap map[int]*models.ReconciliationItem,
	matchedPayments map[int]float64,
	matchedInvoices map[int]float64,
) []models.ProposedAllocation {
	var results []models.ProposedAllocation

	invByScheme := make(map[int][]models.Invoice)
	for _, inv := range invoices {
		invByScheme[inv.SchemeID] = append(invByScheme[inv.SchemeID], inv)
	}

	for _, p := range payments {
		pi := paymentItemMap[p.ID]
		if pi == nil {
			continue
		}
		paymentRemaining := pi.UnallocatedAmount - matchedPayments[p.ID]
		if paymentRemaining <= 0.005 {
			continue
		}

		schemeInvoices := invByScheme[p.SchemeID]
		tolerance := rule.ToleranceVal
		if rule.ToleranceType == "percentage" {
			tolerance = paymentRemaining * (rule.ToleranceVal / 100.0)
		}

		for _, inv := range schemeInvoices {
			invItem, exists := invoiceItemMap[inv.ID]
			if !exists {
				continue
			}
			invoiceRemaining := invItem.UnallocatedAmount - matchedInvoices[inv.ID]
			if invoiceRemaining <= 0.005 {
				continue
			}

			if math.Abs(paymentRemaining-invoiceRemaining) <= tolerance {
				allocAmt := math.Min(paymentRemaining, invoiceRemaining)
				results = append(results, models.ProposedAllocation{
					PaymentID:     p.ID,
					InvoiceID:     inv.ID,
					Amount:        allocAmt,
					MatchedBy:     rule.Name,
					Confidence:    "medium",
					PaymentRef:    p.BankReference,
					InvoiceNumber: inv.InvoiceNumber,
					SchemeName:    inv.SchemeName,
				})
				matchedPayments[p.ID] += allocAmt
				matchedInvoices[inv.ID] += allocAmt
				break
			}
		}
	}

	return results
}

func allocateAcrossInvoices(
	rule models.MatchingRule,
	payment models.Payment,
	invoices []models.Invoice,
	invoiceItemMap map[int]*models.ReconciliationItem,
	matchedPayments map[int]float64,
	matchedInvoices map[int]float64,
	paymentRemaining float64,
	tolerance float64,
	results []models.ProposedAllocation,
) []models.ProposedAllocation {
	remaining := paymentRemaining

	for _, inv := range invoices {
		if remaining <= 0.005 {
			break
		}
		invItem, exists := invoiceItemMap[inv.ID]
		if !exists {
			continue
		}
		invoiceRemaining := invItem.UnallocatedAmount - matchedInvoices[inv.ID]
		if invoiceRemaining <= 0.005 {
			continue
		}

		allocAmt := math.Min(remaining, invoiceRemaining)
		if allocAmt < 0.005 {
			continue
		}

		// For multi-invoice, allow partial allocation of last invoice
		results = append(results, models.ProposedAllocation{
			PaymentID:     payment.ID,
			InvoiceID:     inv.ID,
			Amount:        allocAmt,
			MatchedBy:     rule.Name,
			Confidence:    "medium",
			PaymentRef:    payment.BankReference,
			InvoiceNumber: inv.InvoiceNumber,
			SchemeName:    inv.SchemeName,
		})

		matchedPayments[payment.ID] += allocAmt
		matchedInvoices[inv.ID] += allocAmt
		remaining -= allocAmt
	}

	// If small remainder is left within tolerance, write it off against last allocated invoice
	if remaining > 0 && remaining <= tolerance && len(results) > 0 {
		last := &results[len(results)-1]
		last.Amount += remaining
		matchedPayments[payment.ID] += remaining
		matchedInvoices[last.InvoiceID] += remaining
	}

	return results
}

func applyAllocationToInvoice(tx *gorm.DB, invoiceID int, amount float64) error {
	var invoice models.Invoice
	if err := tx.First(&invoice, invoiceID).Error; err != nil {
		return err
	}
	invoice.PaidAmount += amount
	invoice.Balance = invoice.NetPayable - invoice.PaidAmount

	if invoice.Balance <= 0.005 {
		invoice.Status = "paid"
		invoice.Balance = 0
	} else if invoice.PaidAmount > 0 {
		invoice.Status = "partial"
	} else {
		invoice.Status = "sent"
	}

	return tx.Save(&invoice).Error
}

func updatePaymentStatus(tx *gorm.DB, paymentID int, itemMap map[int]*models.ReconciliationItem) error {
	var payment models.Payment
	if err := tx.First(&payment, paymentID).Error; err != nil {
		return err
	}

	// Check whether fully allocated via reconciliation item
	if itemMap != nil {
		if pi, ok := itemMap[paymentID]; ok {
			if pi.UnallocatedAmount <= 0.005 {
				payment.Status = "matched"
			} else {
				payment.Status = "unmatched"
			}
			return tx.Save(&payment).Error
		}
	}

	// Fallback: check from DB
	var item models.ReconciliationItem
	if err := tx.Where("payment_id = ?", paymentID).First(&item).Error; err == nil {
		if item.UnallocatedAmount <= 0.005 {
			payment.Status = "matched"
		} else {
			payment.Status = "unmatched"
		}
	}

	return tx.Save(&payment).Error
}

func getOrCreatePaymentItem(payment models.Payment) (*models.ReconciliationItem, error) {
	var item models.ReconciliationItem
	err := DB.Where("payment_id = ?", payment.ID).First(&item).Error
	if err == nil {
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = models.ReconciliationItem{
		ItemType:          "payment",
		PaymentID:         &payment.ID,
		SchemeID:          payment.SchemeID,
		SchemeName:        payment.SchemeName,
		OriginalAmount:    payment.Amount,
		AllocatedAmount:   0,
		UnallocatedAmount: payment.Amount,
		Status:            "open",
	}
	err = DB.Create(&item).Error
	return &item, err
}

func getOrCreateInvoiceItemTx(tx *gorm.DB, invoice models.Invoice) (*models.ReconciliationItem, error) {
	var item models.ReconciliationItem
	err := tx.Where("invoice_id = ?", invoice.ID).First(&item).Error
	if err == nil {
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = models.ReconciliationItem{
		ItemType:          "invoice",
		InvoiceID:         &invoice.ID,
		SchemeID:          invoice.SchemeID,
		SchemeName:        invoice.SchemeName,
		OriginalAmount:    invoice.NetPayable,
		AllocatedAmount:   invoice.PaidAmount,
		UnallocatedAmount: invoice.Balance,
		Status:            "open",
	}
	if invoice.PaidAmount > 0 {
		item.Status = "partial"
	}
	err = tx.Create(&item).Error
	return &item, err
}

// recalcReconciliationItem recalculates the allocated/unallocated amounts
// from the allocation ledger for a given entity.
func recalcReconciliationItem(tx *gorm.DB, entityType string, entityID int) error {
	var item models.ReconciliationItem
	var err error

	if entityType == "payment" {
		err = tx.Where("payment_id = ?", entityID).First(&item).Error
	} else {
		err = tx.Where("invoice_id = ?", entityID).First(&item).Error
	}
	if err != nil {
		return err
	}

	// Sum all allocations (including reversals which are negative)
	var totalAllocated float64
	col := "payment_id"
	if entityType == "invoice" {
		col = "invoice_id"
	}
	tx.Model(&models.PaymentAllocation{}).
		Where(col+" = ?", entityID).
		Select("COALESCE(SUM(allocated_amount), 0)").
		Scan(&totalAllocated)

	item.AllocatedAmount = totalAllocated
	item.UnallocatedAmount = item.OriginalAmount - totalAllocated
	if item.UnallocatedAmount <= 0.005 {
		item.Status = "matched"
		item.UnallocatedAmount = 0
	} else if totalAllocated > 0 {
		item.Status = "partial"
	} else {
		item.Status = "open"
	}
	now := time.Now()
	item.LastActionDate = &now

	return tx.Save(&item).Error
}

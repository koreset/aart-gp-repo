package services

import (
	"api/models"
	"fmt"
	"math"
	"strings"
	"time"
)

// quarterLabel converts a period_end date string to a quarter label e.g. "Q1"
func quarterLabel(periodEnd string) string {
	if len(periodEnd) < 7 {
		return "Q1"
	}
	month := 0
	fmt.Sscanf(periodEnd[5:7], "%d", &month)
	switch {
	case month <= 3:
		return "Q1"
	case month <= 6:
		return "Q2"
	case month <= 9:
		return "Q3"
	default:
		return "Q4"
	}
}

// generateAccountNumber produces a unique account number e.g. TA-2025-Q1-MUNICH
func generateAccountNumber(treaty models.ReinsuranceTreaty, periodEnd string) string {
	year := ""
	if len(periodEnd) >= 4 {
		year = periodEnd[:4]
	}
	q := quarterLabel(periodEnd)
	riCode := strings.ToUpper(treaty.ReinsurerCode)
	if len(riCode) > 6 {
		riCode = riCode[:6]
	}
	if riCode == "" {
		riCode = strings.ToUpper(strings.ReplaceAll(treaty.ReinsurerName, " ", ""))
		if len(riCode) > 6 {
			riCode = riCode[:6]
		}
	}
	base := fmt.Sprintf("TA-%s-%s-%s", year, q, riCode)
	candidate := base
	for i := 2; i < 20; i++ {
		var count int64
		DB.Model(&models.TechnicalAccount{}).Where("account_number = ?", candidate).Count(&count)
		if count == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return fmt.Sprintf("%s-%d", base, time.Now().UnixNano())
}

// GenerateTechnicalAccount creates a new technical account by aggregating RI bordereaux runs for the period
func GenerateTechnicalAccount(req models.GenerateTechnicalAccountRequest, user models.AppUser) (models.TechnicalAccount, error) {
	treaty, err := GetTreatyByID(req.TreatyID)
	if err != nil {
		return models.TechnicalAccount{}, fmt.Errorf("treaty not found: %w", err)
	}

	// Sum ceded premium from member census runs in period
	var memberRuns []models.RIBordereauxRun
	DB.Where("treaty_id = ? AND type = 'member_census' AND period_start >= ? AND period_end <= ?",
		treaty.ID, req.PeriodStart, req.PeriodEnd).Find(&memberRuns)

	var cededPremiumEarned float64
	for _, r := range memberRuns {
		cededPremiumEarned += r.CededPremium
	}

	// Sum ceded claims from claims runs in period
	var claimsRuns []models.RIBordereauxRun
	DB.Where("treaty_id = ? AND type = 'claims_run' AND period_start >= ? AND period_end <= ?",
		treaty.ID, req.PeriodStart, req.PeriodEnd).Find(&claimsRuns)

	var cededClaimsPaid float64
	for _, r := range claimsRuns {
		cededClaimsPaid += r.CededClaimsIncurred
	}

	riCommission := cededPremiumEarned * (treaty.ReinsuranceCommissionRate / 100)
	netCededPremium := cededPremiumEarned - riCommission
	totalCededClaims := cededClaimsPaid

	// Profit commission: if loss ratio < 60% apply profit commission
	profitCommission := 0.0
	if netCededPremium > 0 && treaty.ProfitCommissionRate > 0 {
		lossRatio := totalCededClaims / netCededPremium
		if lossRatio < 0.60 {
			profitCommission = (netCededPremium - totalCededClaims) * (treaty.ProfitCommissionRate / 100)
		}
	}

	// Positive net balance = cedant owes reinsurer; negative = reinsurer owes cedant
	netBalance := netCededPremium - totalCededClaims - profitCommission

	accountNumber := generateAccountNumber(treaty, req.PeriodEnd)

	account := models.TechnicalAccount{
		AccountNumber:         accountNumber,
		TreatyID:              treaty.ID,
		TreatyNumber:          treaty.TreatyNumber,
		ReinsurerName:         treaty.ReinsurerName,
		PeriodStart:           req.PeriodStart,
		PeriodEnd:             req.PeriodEnd,
		CededPremiumEarned:    cededPremiumEarned,
		ReinsuranceCommission: riCommission,
		NetCededPremium:       netCededPremium,
		CededClaimsPaid:       cededClaimsPaid,
		CededIBNR:             0,
		TotalCededClaims:      totalCededClaims,
		ProfitCommission:      profitCommission,
		NetBalance:            netBalance,
		Status:                "draft",
		Notes:                 req.Notes,
		CreatedBy:             user.UserName,
	}

	if err := DB.Create(&account).Error; err != nil {
		return account, fmt.Errorf("failed to create technical account: %w", err)
	}
	return account, nil
}

// GetTechnicalAccounts returns accounts filtered by treaty and status
func GetTechnicalAccounts(treatyID int, status string) ([]models.TechnicalAccount, error) {
	var accounts []models.TechnicalAccount
	q := DB.Order("created_at DESC")
	if treatyID > 0 {
		q = q.Where("treaty_id = ?", treatyID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetTechnicalAccountByID returns a single account
func GetTechnicalAccountByID(id int) (models.TechnicalAccount, error) {
	var account models.TechnicalAccount
	if err := DB.First(&account, id).Error; err != nil {
		return account, fmt.Errorf("technical account %d not found: %w", id, err)
	}
	return account, nil
}

// UpdateTechnicalAccount updates status, IBNR and notes; advances status timestamps
func UpdateTechnicalAccount(id int, req models.UpdateTechnicalAccountRequest, user models.AppUser) (models.TechnicalAccount, error) {
	account, err := GetTechnicalAccountByID(id)
	if err != nil {
		return account, err
	}

	now := time.Now()
	statusChanged := req.Status != "" && req.Status != account.Status
	if statusChanged {
		account.Status = req.Status
		switch req.Status {
		case "issued":
			account.IssuedAt = &now
		case "agreed":
			account.AgreedAt = &now
		case "settled":
			account.SettledAt = &now
		}
	}
	if req.DisputeNotes != "" {
		account.DisputeNotes = req.DisputeNotes
	}
	if req.Notes != "" {
		account.Notes = req.Notes
	}
	if req.CededIBNR != 0 {
		account.CededIBNR = req.CededIBNR
		account.TotalCededClaims = account.CededClaimsPaid + account.CededIBNR
		account.NetBalance = account.NetCededPremium - account.TotalCededClaims - account.ProfitCommission
	}
	account.UpdatedBy = user.UserName
	account.UpdatedAt = now

	if err := DB.Save(&account).Error; err != nil {
		return account, fmt.Errorf("failed to update technical account: %w", err)
	}
	if statusChanged {
		go NotifySettlementUpdated(account, user, req.Status)
	}
	return account, nil
}

// RecordSettlementPayment records a payment against a technical account; marks settled when fully paid
func RecordSettlementPayment(req models.RecordSettlementPaymentRequest, user models.AppUser) (models.SettlementPayment, error) {
	account, err := GetTechnicalAccountByID(req.TechnicalAccountID)
	if err != nil {
		return models.SettlementPayment{}, err
	}

	payment := models.SettlementPayment{
		TechnicalAccountID: req.TechnicalAccountID,
		PaymentDate:        req.PaymentDate,
		Amount:             req.Amount,
		Direction:          req.Direction,
		PaymentMethod:      req.PaymentMethod,
		Reference:          req.Reference,
		Notes:              req.Notes,
		RecordedBy:         user.UserName,
	}
	if err := DB.Create(&payment).Error; err != nil {
		return payment, fmt.Errorf("failed to record payment: %w", err)
	}

	// Sum all payments and mark settled if >= net balance
	var payments []models.SettlementPayment
	DB.Where("technical_account_id = ?", req.TechnicalAccountID).Find(&payments)
	var totalPaid float64
	for _, p := range payments {
		totalPaid += p.Amount
	}
	if totalPaid >= math.Abs(account.NetBalance) && account.Status != "settled" {
		now := time.Now()
		account.Status = "settled"
		account.SettledAt = &now
		account.UpdatedAt = now
		DB.Save(&account)
	}

	return payment, nil
}

// GetSettlementPayments returns payments for a given account
func GetSettlementPayments(accountID int) ([]models.SettlementPayment, error) {
	var payments []models.SettlementPayment
	q := DB.Order("created_at DESC")
	if accountID > 0 {
		q = q.Where("technical_account_id = ?", accountID)
	}
	if err := q.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// EscalateSettlementDispute escalates a disputed technical account to the next stage
func EscalateSettlementDispute(id int, req models.EscalateDisputeRequest, user models.AppUser) (models.TechnicalAccount, error) {
	validStages := map[string]bool{"stage1": true, "stage2": true, "stage3": true}
	if !validStages[req.Stage] {
		return models.TechnicalAccount{}, fmt.Errorf("stage must be one of: stage1, stage2, stage3")
	}
	account, err := GetTechnicalAccountByID(id)
	if err != nil {
		return account, err
	}
	if account.Status != "disputed" {
		return account, fmt.Errorf("account must be in disputed status to escalate (current: %s)", account.Status)
	}
	now := time.Now()
	account.DisputeStatus = req.Stage
	account.DisputeEscalatedAt = &now
	if req.Notes != "" {
		account.DisputeNotes = req.Notes
	}
	account.UpdatedBy = user.UserName
	account.UpdatedAt = now
	if err := DB.Save(&account).Error; err != nil {
		return account, fmt.Errorf("failed to escalate dispute: %w", err)
	}
	return account, nil
}

// ResolveSettlementDispute marks a dispute as resolved on the technical account
func ResolveSettlementDispute(id int, req models.ResolveDisputeRequest, user models.AppUser) (models.TechnicalAccount, error) {
	account, err := GetTechnicalAccountByID(id)
	if err != nil {
		return account, err
	}
	if account.Status != "disputed" {
		return account, fmt.Errorf("account must be in disputed status to resolve (current: %s)", account.Status)
	}
	now := time.Now()
	account.Status = "agreed"
	account.DisputeStatus = "resolved"
	account.AgreedAt = &now
	if req.Notes != "" {
		account.DisputeNotes = req.Notes
	}
	account.UpdatedBy = user.UserName
	account.UpdatedAt = now
	if err := DB.Save(&account).Error; err != nil {
		return account, fmt.Errorf("failed to resolve dispute: %w", err)
	}
	go NotifySettlementDisputeResolved(account, user)
	return account, nil
}

// GetSettlementStats returns a summary of technical account statuses and outstanding balances
func GetSettlementStats(treatyID int) (models.SettlementStats, error) {
	var accounts []models.TechnicalAccount
	q := DB.Model(&models.TechnicalAccount{})
	if treatyID > 0 {
		q = q.Where("treaty_id = ?", treatyID)
	}
	if err := q.Find(&accounts).Error; err != nil {
		return models.SettlementStats{}, err
	}
	stats := models.SettlementStats{Total: len(accounts)}
	for _, a := range accounts {
		switch a.Status {
		case "draft":
			stats.Draft++
		case "issued":
			stats.Issued++
		case "agreed":
			stats.Agreed++
		case "settled":
			stats.Settled++
		case "disputed":
			stats.Disputed++
		}
		if a.NetBalance > 0 {
			stats.NetOwedToRI += a.NetBalance
		} else {
			stats.NetOwedByCedant += math.Abs(a.NetBalance)
		}
	}
	return stats, nil
}

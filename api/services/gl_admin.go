package services

import (
	"api/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Admin/CRUD surface for the operational General Ledger.
//
// The double-entry posting path lives in gl_posting.go; this file owns
// chart-of-accounts, posting-rule, bank-account, period, and bank-rec
// management.
//
// Master-data mutations (GLAccount, PostingRule, BankAccount) and period
// close all follow the maker/checker pattern: one user requests the change,
// a different user approves it. Self-approval is rejected at the service
// layer (ErrSelfApproval). Every transition writes a GLAuditLog row inside
// the same transaction.

// ---------------------------------------------------------------------------
// Chart of accounts
// ---------------------------------------------------------------------------

func GLListAccounts() ([]models.GLAccount, error) {
	var accounts []models.GLAccount
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		return d.Order("code ASC").Find(&accounts).Error
	})
	return accounts, err
}

func GLGetAccount(id int) (models.GLAccount, error) {
	var account models.GLAccount
	err := DB.First(&account, id).Error
	return account, err
}

// GLRequestCreateAccount stages a new GL account creation for approval. The
// row is persisted immediately with ApprovalStatus="pending_create" but
// IsActive=false so it cannot be used for postings until approved.
func GLRequestCreateAccount(req models.GLAccount, user models.AppUser) (models.GLAccount, error) {
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		return req, errors.New("code and name are required")
	}
	if req.AccountType == "" || req.NormalBalance == "" {
		return req, errors.New("account_type and normal_balance are required")
	}
	now := time.Now()
	req.IsActive = false
	req.ApprovalStatus = "pending_create"
	req.CreatedBy = user.UserName
	req.UpdatedBy = user.UserName
	req.PendingRequestedBy = user.UserName
	req.PendingRequestedAt = &now

	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&req).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "gl_account_change_requested", "gl_account", req.Code, req.ID, user, map[string]any{
			"action": "create",
			"name":   req.Name,
		})
	})
	return req, err
}

// GLRequestUpdateAccount stores a proposed update as PendingChangeJSON
// without touching the live fields. The live row continues to function
// normally until an approver applies the change.
func GLRequestUpdateAccount(id int, patch models.GLAccount, user models.AppUser) (models.GLAccount, error) {
	var existing models.GLAccount
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus != "active" && existing.ApprovalStatus != "" {
		return existing, fmt.Errorf("account already has a pending change (%s)", existing.ApprovalStatus)
	}
	proposed := map[string]any{
		"name":           patch.Name,
		"account_type":   patch.AccountType,
		"normal_balance": patch.NormalBalance,
		"parent_id":      patch.ParentID,
		"is_active":      patch.IsActive,
		"description":    patch.Description,
	}
	body, err := json.Marshal(proposed)
	if err != nil {
		return existing, err
	}
	now := time.Now()
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Updates(map[string]any{
			"approval_status":      "pending_update",
			"pending_change_json":  string(body),
			"pending_requested_by": user.UserName,
			"pending_requested_at": &now,
			"updated_by":           user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "gl_account_change_requested", "gl_account", existing.Code, existing.ID, user, map[string]any{
			"action":   "update",
			"proposed": proposed,
		})
	})
	if err != nil {
		return existing, err
	}
	return GLGetAccount(id)
}

// GLRequestDeactivateAccount stages a deactivation for approval. The live
// IsActive flag is NOT touched until approval — postings to the account
// continue to work during the pending window.
func GLRequestDeactivateAccount(id int, user models.AppUser) (models.GLAccount, error) {
	var existing models.GLAccount
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus != "active" && existing.ApprovalStatus != "" {
		return existing, fmt.Errorf("account already has a pending change (%s)", existing.ApprovalStatus)
	}
	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Updates(map[string]any{
			"approval_status":      "pending_deactivate",
			"pending_requested_by": user.UserName,
			"pending_requested_at": &now,
			"updated_by":           user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "gl_account_change_requested", "gl_account", existing.Code, existing.ID, user, map[string]any{
			"action": "deactivate",
		})
	})
	if err != nil {
		return existing, err
	}
	return GLGetAccount(id)
}

// GLApproveAccountChange applies whatever pending change is on the row and
// flips ApprovalStatus back to "active". Enforces SoD — approver must differ
// from PendingRequestedBy.
func GLApproveAccountChange(id int, notes string, user models.AppUser) (models.GLAccount, error) {
	var existing models.GLAccount
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus == "active" || existing.ApprovalStatus == "" {
		return existing, fmt.Errorf("account has no pending change to approve")
	}
	if existing.PendingRequestedBy == user.UserName {
		return existing, ErrSelfApproval
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"approval_status":      "active",
			"pending_change_json":  "",
			"pending_requested_by": "",
			"pending_requested_at": nil,
			"updated_by":           user.UserName,
		}
		switch existing.ApprovalStatus {
		case "pending_create":
			updates["is_active"] = true
		case "pending_deactivate":
			updates["is_active"] = false
		case "pending_update":
			var proposed map[string]any
			if err := json.Unmarshal([]byte(existing.PendingChangeJSON), &proposed); err != nil {
				return fmt.Errorf("decode pending change: %w", err)
			}
			for k, v := range proposed {
				updates[k] = v
			}
		}
		if err := tx.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "gl_account_change_approved", "gl_account", existing.Code, existing.ID, user, map[string]any{
			"prior_status": existing.ApprovalStatus,
			"requested_by": existing.PendingRequestedBy,
			"notes":        notes,
		})
	})
	if err != nil {
		return existing, err
	}
	return GLGetAccount(id)
}

// ---------------------------------------------------------------------------
// Accounting periods
// ---------------------------------------------------------------------------

func GLListPeriods() ([]models.AccountingPeriod, error) {
	var periods []models.AccountingPeriod
	err := DB.Order("start_date DESC").Find(&periods).Error
	return periods, err
}

// GLCreatePeriod opens a new monthly period. The name is normalised to
// YYYY-MM via the start date so callers can't accidentally drift it.
func GLCreatePeriod(req models.AccountingPeriod, user models.AppUser) (models.AccountingPeriod, error) {
	if req.StartDate.IsZero() {
		return req, errors.New("start_date is required")
	}
	monthStart := time.Date(req.StartDate.Year(), req.StartDate.Month(), 1, 0, 0, 0, 0, req.StartDate.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)
	period := models.AccountingPeriod{
		Name:      monthStart.Format("2006-01"),
		StartDate: monthStart,
		EndDate:   monthEnd,
		Status:    "open",
	}
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&period).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "period_opened", "accounting_period", period.Name, period.ID, user, nil)
	})
	return period, err
}

// GLRequestClosePeriod stamps the period as close_requested. A different user
// then commits the close via GLClosePeriod.
func GLRequestClosePeriod(id int, user models.AppUser) (models.AccountingPeriod, error) {
	var period models.AccountingPeriod
	if err := DB.First(&period, id).Error; err != nil {
		return period, err
	}
	if period.Status != "open" {
		return period, fmt.Errorf("period is %s — cannot request close", period.Status)
	}
	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&period).Updates(map[string]any{
			"status":              "close_requested",
			"close_requested_by":  user.UserName,
			"close_requested_at":  &now,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "period_close_requested", "accounting_period", period.Name, period.ID, user, nil)
	})
	if err != nil {
		return period, err
	}
	return period, DB.First(&period, id).Error
}

// GLClosePeriod commits a close that was previously requested. Subsequent
// postings whose natural date falls inside the period will be rejected by
// writeJournalEntry. Enforces SoD — the closer must differ from the
// requester.
func GLClosePeriod(id int, user models.AppUser) (models.AccountingPeriod, error) {
	var period models.AccountingPeriod
	if err := DB.First(&period, id).Error; err != nil {
		return period, err
	}
	if period.Status == "closed" {
		return period, errors.New("period is already closed")
	}
	if period.Status != "close_requested" {
		return period, fmt.Errorf("period must be close_requested before it can be closed (currently %s)", period.Status)
	}
	if period.CloseRequestedBy == user.UserName {
		return period, ErrSelfApproval
	}
	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&period).Updates(map[string]any{
			"status":    "closed",
			"closed_at": &now,
			"closed_by": user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "period_closed", "accounting_period", period.Name, period.ID, user, map[string]any{
			"requested_by": period.CloseRequestedBy,
		})
	})
	if err != nil {
		return period, err
	}
	return period, DB.First(&period, id).Error
}

// ---------------------------------------------------------------------------
// Posting rules
// ---------------------------------------------------------------------------

func GLListPostingRules() ([]models.PostingRule, error) {
	var rules []models.PostingRule
	err := DB.Order("event_key ASC").Find(&rules).Error
	return rules, err
}

// GLRequestCreatePostingRule stages a new rule for approval. Created with
// IsActive=false so GLPost will not pick it up before the approver releases it.
func GLRequestCreatePostingRule(rule models.PostingRule, user models.AppUser) (models.PostingRule, error) {
	if strings.TrimSpace(rule.EventKey) == "" {
		return rule, errors.New("event_key is required")
	}
	if rule.DebitAccountID == 0 || rule.CreditAccountID == 0 {
		return rule, errors.New("debit and credit accounts are required")
	}
	now := time.Now()
	rule.IsActive = false
	rule.ApprovalStatus = "pending_create"
	rule.CreatedBy = user.UserName
	rule.UpdatedBy = user.UserName
	rule.PendingRequestedBy = user.UserName
	rule.PendingRequestedAt = &now
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&rule).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "posting_rule_change_requested", "posting_rule", rule.EventKey, rule.ID, user, map[string]any{
			"action":           "create",
			"debit_account":    rule.DebitAccountID,
			"credit_account":   rule.CreditAccountID,
		})
	})
	return rule, err
}

// GLRequestUpdatePostingRule stages a proposed edit. The live row keeps
// working until approved.
func GLRequestUpdatePostingRule(id int, patch models.PostingRule, user models.AppUser) (models.PostingRule, error) {
	var existing models.PostingRule
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus != "active" && existing.ApprovalStatus != "" {
		return existing, fmt.Errorf("rule already has a pending change (%s)", existing.ApprovalStatus)
	}
	proposed := map[string]any{
		"debit_account_id":  patch.DebitAccountID,
		"credit_account_id": patch.CreditAccountID,
		"is_active":         patch.IsActive,
		"notes":             patch.Notes,
	}
	body, err := json.Marshal(proposed)
	if err != nil {
		return existing, err
	}
	now := time.Now()
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Updates(map[string]any{
			"approval_status":      "pending_update",
			"pending_change_json":  string(body),
			"pending_requested_by": user.UserName,
			"pending_requested_at": &now,
			"updated_by":           user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "posting_rule_change_requested", "posting_rule", existing.EventKey, existing.ID, user, map[string]any{
			"action":   "update",
			"proposed": proposed,
		})
	})
	if err != nil {
		return existing, err
	}
	var refreshed models.PostingRule
	return refreshed, DB.First(&refreshed, id).Error
}

// GLRequestDeletePostingRule stages a deletion (we tombstone via
// ApprovalStatus rather than removing the row outright, so the audit trail
// keeps its FK target).
func GLRequestDeletePostingRule(id int, user models.AppUser) (models.PostingRule, error) {
	var existing models.PostingRule
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus != "active" && existing.ApprovalStatus != "" {
		return existing, fmt.Errorf("rule already has a pending change (%s)", existing.ApprovalStatus)
	}
	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Updates(map[string]any{
			"approval_status":      "pending_delete",
			"pending_requested_by": user.UserName,
			"pending_requested_at": &now,
			"updated_by":           user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "posting_rule_change_requested", "posting_rule", existing.EventKey, existing.ID, user, map[string]any{
			"action": "delete",
		})
	})
	if err != nil {
		return existing, err
	}
	var refreshed models.PostingRule
	return refreshed, DB.First(&refreshed, id).Error
}

// GLApprovePostingRuleChange applies the pending change. For pending_delete
// the row is deactivated rather than physically removed.
func GLApprovePostingRuleChange(id int, notes string, user models.AppUser) (models.PostingRule, error) {
	var existing models.PostingRule
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus == "active" || existing.ApprovalStatus == "" {
		return existing, fmt.Errorf("rule has no pending change to approve")
	}
	if existing.PendingRequestedBy == user.UserName {
		return existing, ErrSelfApproval
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"approval_status":      "active",
			"pending_change_json":  "",
			"pending_requested_by": "",
			"pending_requested_at": nil,
			"updated_by":           user.UserName,
		}
		switch existing.ApprovalStatus {
		case "pending_create":
			updates["is_active"] = true
		case "pending_delete":
			updates["is_active"] = false
		case "pending_update":
			var proposed map[string]any
			if err := json.Unmarshal([]byte(existing.PendingChangeJSON), &proposed); err != nil {
				return fmt.Errorf("decode pending change: %w", err)
			}
			for k, v := range proposed {
				updates[k] = v
			}
		}
		if err := tx.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "posting_rule_change_approved", "posting_rule", existing.EventKey, existing.ID, user, map[string]any{
			"prior_status": existing.ApprovalStatus,
			"requested_by": existing.PendingRequestedBy,
			"notes":        notes,
		})
	})
	if err != nil {
		return existing, err
	}
	var refreshed models.PostingRule
	return refreshed, DB.First(&refreshed, id).Error
}

// ---------------------------------------------------------------------------
// Journal entry queries
// ---------------------------------------------------------------------------

// GLListJournalsOptions scopes the journals list view.
type GLListJournalsOptions struct {
	PeriodID   int
	SourceType string
	Status     string
	AccountID  int
	From       time.Time
	To         time.Time
	Limit      int
}

func GLListJournals(opts GLListJournalsOptions) ([]models.JournalEntry, error) {
	q := DB.Model(&models.JournalEntry{})
	if opts.PeriodID > 0 {
		q = q.Where("period_id = ?", opts.PeriodID)
	}
	if opts.SourceType != "" {
		q = q.Where("source_type = ?", opts.SourceType)
	}
	if opts.Status != "" {
		q = q.Where("status = ?", opts.Status)
	}
	if !opts.From.IsZero() {
		q = q.Where("posted_at >= ?", opts.From)
	}
	if !opts.To.IsZero() {
		q = q.Where("posted_at <= ?", opts.To)
	}
	if opts.AccountID > 0 {
		// Filter via the join — return entries that touch the account.
		q = q.Where("id IN (?)", DB.Model(&models.JournalLine{}).Select("entry_id").Where("account_id = ?", opts.AccountID))
	}
	limit := opts.Limit
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	var rows []models.JournalEntry
	err := q.Order("created_at DESC, id DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

func GLGetJournalEntry(id int) (models.JournalEntry, error) {
	var entry models.JournalEntry
	err := DB.Preload("Lines").First(&entry, id).Error
	return entry, err
}

// ---------------------------------------------------------------------------
// Bank account + statement reconciliation
// ---------------------------------------------------------------------------

func GLListBankAccounts() ([]models.BankAccount, error) {
	var rows []models.BankAccount
	err := DB.Order("code ASC").Find(&rows).Error
	return rows, err
}

// GLRequestCreateBankAccount stages a new bank account for approval. Created
// inactive so it cannot be selected for statement imports until released.
func GLRequestCreateBankAccount(req models.BankAccount, user models.AppUser) (models.BankAccount, error) {
	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
		return req, errors.New("code and name are required")
	}
	if req.GLAccountID == 0 {
		return req, errors.New("gl_account_id is required")
	}
	now := time.Now()
	req.IsActive = false
	req.ApprovalStatus = "pending_create"
	req.CreatedBy = user.UserName
	req.UpdatedBy = user.UserName
	req.PendingRequestedBy = user.UserName
	req.PendingRequestedAt = &now
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&req).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "bank_account_change_requested", "bank_account", req.Code, req.ID, user, map[string]any{
			"action":         "create",
			"name":           req.Name,
			"bank_name":      req.BankName,
			"account_number": req.AccountNumber,
		})
	})
	return req, err
}

// GLRequestUpdateBankAccount stages a proposed edit. The live row keeps its
// existing values (including account_number) until an approver releases the
// change — this is what stops a single insider from silently rerouting
// payments.
func GLRequestUpdateBankAccount(id int, patch models.BankAccount, user models.AppUser) (models.BankAccount, error) {
	var existing models.BankAccount
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus != "active" && existing.ApprovalStatus != "" {
		return existing, fmt.Errorf("bank account already has a pending change (%s)", existing.ApprovalStatus)
	}
	proposed := map[string]any{
		"name":           patch.Name,
		"bank_name":      patch.BankName,
		"account_number": patch.AccountNumber,
		"gl_account_id":  patch.GLAccountID,
		"currency":       patch.Currency,
		"is_active":      patch.IsActive,
	}
	body, err := json.Marshal(proposed)
	if err != nil {
		return existing, err
	}
	now := time.Now()
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Updates(map[string]any{
			"approval_status":      "pending_update",
			"pending_change_json":  string(body),
			"pending_requested_by": user.UserName,
			"pending_requested_at": &now,
			"updated_by":           user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "bank_account_change_requested", "bank_account", existing.Code, existing.ID, user, map[string]any{
			"action":   "update",
			"proposed": proposed,
		})
	})
	if err != nil {
		return existing, err
	}
	var refreshed models.BankAccount
	return refreshed, DB.First(&refreshed, id).Error
}

// GLRequestDeactivateBankAccount stages a deactivation. Live IsActive stays
// true during the pending window.
func GLRequestDeactivateBankAccount(id int, user models.AppUser) (models.BankAccount, error) {
	var existing models.BankAccount
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus != "active" && existing.ApprovalStatus != "" {
		return existing, fmt.Errorf("bank account already has a pending change (%s)", existing.ApprovalStatus)
	}
	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existing).Updates(map[string]any{
			"approval_status":      "pending_deactivate",
			"pending_requested_by": user.UserName,
			"pending_requested_at": &now,
			"updated_by":           user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "bank_account_change_requested", "bank_account", existing.Code, existing.ID, user, map[string]any{
			"action": "deactivate",
		})
	})
	if err != nil {
		return existing, err
	}
	var refreshed models.BankAccount
	return refreshed, DB.First(&refreshed, id).Error
}

// GLApproveBankAccountChange applies the pending change. Enforces SoD.
func GLApproveBankAccountChange(id int, notes string, user models.AppUser) (models.BankAccount, error) {
	var existing models.BankAccount
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	if existing.ApprovalStatus == "active" || existing.ApprovalStatus == "" {
		return existing, fmt.Errorf("bank account has no pending change to approve")
	}
	if existing.PendingRequestedBy == user.UserName {
		return existing, ErrSelfApproval
	}

	err := DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"approval_status":      "active",
			"pending_change_json":  "",
			"pending_requested_by": "",
			"pending_requested_at": nil,
			"updated_by":           user.UserName,
		}
		switch existing.ApprovalStatus {
		case "pending_create":
			updates["is_active"] = true
		case "pending_deactivate":
			updates["is_active"] = false
		case "pending_update":
			var proposed map[string]any
			if err := json.Unmarshal([]byte(existing.PendingChangeJSON), &proposed); err != nil {
				return fmt.Errorf("decode pending change: %w", err)
			}
			for k, v := range proposed {
				updates[k] = v
			}
		}
		if err := tx.Model(&existing).Updates(updates).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "bank_account_change_approved", "bank_account", existing.Code, existing.ID, user, map[string]any{
			"prior_status": existing.ApprovalStatus,
			"requested_by": existing.PendingRequestedBy,
			"notes":        notes,
		})
	})
	if err != nil {
		return existing, err
	}
	var refreshed models.BankAccount
	return refreshed, DB.First(&refreshed, id).Error
}

// GLImportBankStatement stores a batch of statement lines under a generated
// batch ID. Date strings are accepted in ISO-8601 (YYYY-MM-DD) form to keep
// the CSV ingest path uncluttered. The importer is captured for audit but
// statement imports do not themselves require approval — they're raw data
// ingest; the maker/checker step happens at match/review time.
func GLImportBankStatement(req models.BankStatementImportRequest, user models.AppUser) (string, int, error) {
	if req.BankAccountID == 0 {
		return "", 0, errors.New("bank_account_id is required")
	}
	var bank models.BankAccount
	if err := DB.First(&bank, req.BankAccountID).Error; err != nil {
		return "", 0, fmt.Errorf("bank account not found: %w", err)
	}
	if !bank.IsActive {
		return "", 0, errors.New("bank account is inactive")
	}
	batchID := fmt.Sprintf("BATCH-%s", time.Now().Format("20060102-150405"))
	rows := make([]models.BankStatementLine, 0, len(req.Rows))
	for _, in := range req.Rows {
		stmtDate, err := time.Parse("2006-01-02", in.StatementDate)
		if err != nil {
			return "", 0, fmt.Errorf("invalid statement_date %q: %w", in.StatementDate, err)
		}
		var valueDate *time.Time
		if in.ValueDate != "" {
			vd, err := time.Parse("2006-01-02", in.ValueDate)
			if err != nil {
				return "", 0, fmt.Errorf("invalid value_date %q: %w", in.ValueDate, err)
			}
			valueDate = &vd
		}
		rows = append(rows, models.BankStatementLine{
			BankAccountID: req.BankAccountID,
			StatementDate: stmtDate,
			ValueDate:     valueDate,
			Description:   in.Description,
			Amount:        in.Amount,
			Reference:     in.Reference,
			ImportBatchID: batchID,
			ImportedBy:    user.UserName,
			MatchStatus:   "unmatched",
			ReviewStatus:  "not_required",
		})
	}
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(rows, 500).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "bank_statement_imported", "bank_account", bank.Code, bank.ID, user, map[string]any{
			"batch_id": batchID,
			"rows":     len(rows),
		})
	})
	if err != nil {
		return "", 0, err
	}
	return batchID, len(rows), nil
}

// GLListStatementLines returns every imported statement line for a bank
// account, optionally restricted by match status.
func GLListStatementLines(bankAccountID int, status string) ([]models.BankStatementLine, error) {
	q := DB.Where("bank_account_id = ?", bankAccountID)
	if status != "" {
		q = q.Where("match_status = ?", status)
	}
	var rows []models.BankStatementLine
	err := q.Order("statement_date DESC, id DESC").Find(&rows).Error
	return rows, err
}

// GLMatchStatementLine pairs a bank statement line to a posted journal line
// on the bank account's GL account. Both sides are checked to confirm they
// refer to the same GL account.
//
// The match goes into ReviewStatus="pending_review" — a different user must
// sign off via GLReviewStatementLine before the line is considered
// reconciled.
func GLMatchStatementLine(req models.BankMatchRequest, user models.AppUser) error {
	if req.StatementLineID == 0 || req.JournalLineID == 0 {
		return errors.New("statement_line_id and journal_line_id are required")
	}
	var stmt models.BankStatementLine
	if err := DB.First(&stmt, req.StatementLineID).Error; err != nil {
		return fmt.Errorf("statement line not found: %w", err)
	}
	if stmt.MatchStatus == "matched" || stmt.MatchStatus == "ignored" {
		return fmt.Errorf("statement line already %s", stmt.MatchStatus)
	}
	var bank models.BankAccount
	if err := DB.First(&bank, stmt.BankAccountID).Error; err != nil {
		return err
	}
	var line models.JournalLine
	if err := DB.First(&line, req.JournalLineID).Error; err != nil {
		return fmt.Errorf("journal line not found: %w", err)
	}
	if line.AccountID != bank.GLAccountID {
		return fmt.Errorf("journal line account %d does not match bank account's GL account %d", line.AccountID, bank.GLAccountID)
	}
	now := time.Now()
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&stmt).Updates(map[string]any{
			"matched_journal_line_id": &line.ID,
			"match_status":            "matched",
			"matched_at":              &now,
			"matched_by":              user.UserName,
			"review_status":           "pending_review",
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "statement_line_matched", "bank_statement_line", stmt.Reference, stmt.ID, user, map[string]any{
			"journal_line_id": line.ID,
			"amount":          stmt.Amount,
		})
	})
}

// GLIgnoreStatementLine marks a statement line as deliberately not posted
// (e.g. bank fee waived, duplicate row in upload). Like matches, ignores
// require reviewer sign-off.
func GLIgnoreStatementLine(statementLineID int, user models.AppUser) error {
	var stmt models.BankStatementLine
	if err := DB.First(&stmt, statementLineID).Error; err != nil {
		return fmt.Errorf("statement line not found: %w", err)
	}
	if stmt.MatchStatus == "matched" || stmt.MatchStatus == "ignored" {
		return fmt.Errorf("statement line already %s", stmt.MatchStatus)
	}
	now := time.Now()
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&stmt).Updates(map[string]any{
			"match_status":  "ignored",
			"matched_at":    &now,
			"matched_by":    user.UserName,
			"review_status": "pending_review",
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "statement_line_ignored", "bank_statement_line", stmt.Reference, stmt.ID, user, map[string]any{
			"amount": stmt.Amount,
		})
	})
}

// GLReviewStatementLine is the second-pair-of-eyes sign-off on a previous
// match or ignore. Reviewer must differ from MatchedBy. Outcome:
//   - "reviewed":  ReviewStatus → reviewed (line is reconciled)
//   - "rejected":  match metadata is cleared and the line returns to
//                  unmatched / not_required for re-work
func GLReviewStatementLine(statementLineID int, req models.StatementLineReviewRequest, user models.AppUser) (models.BankStatementLine, error) {
	var stmt models.BankStatementLine
	if err := DB.First(&stmt, statementLineID).Error; err != nil {
		return stmt, fmt.Errorf("statement line not found: %w", err)
	}
	if stmt.ReviewStatus != "pending_review" {
		return stmt, fmt.Errorf("statement line is not pending review (review_status=%s)", stmt.ReviewStatus)
	}
	if stmt.MatchedBy == user.UserName {
		return stmt, ErrSelfApproval
	}
	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		var updates map[string]any
		var eventType string
		if req.Outcome == "rejected" {
			updates = map[string]any{
				"review_status":           "not_required",
				"reviewed_by":             user.UserName,
				"reviewed_at":             &now,
				"review_notes":            req.Notes,
				"match_status":            "unmatched",
				"matched_journal_line_id": nil,
				"matched_at":              nil,
				"matched_by":              "",
			}
			eventType = "statement_line_review_rejected"
		} else {
			updates = map[string]any{
				"review_status": "reviewed",
				"reviewed_by":   user.UserName,
				"reviewed_at":   &now,
				"review_notes":  req.Notes,
			}
			eventType = "statement_line_match_reviewed"
		}
		if err := tx.Model(&stmt).Updates(updates).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, eventType, "bank_statement_line", stmt.Reference, stmt.ID, user, map[string]any{
			"matched_by":    stmt.MatchedBy,
			"prior_status":  stmt.MatchStatus,
			"notes":         req.Notes,
		})
	})
	if err != nil {
		return stmt, err
	}
	return stmt, DB.First(&stmt, statementLineID).Error
}

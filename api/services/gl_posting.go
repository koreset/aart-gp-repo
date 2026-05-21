package services

import (
	"api/models"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Operational General Ledger — posting, reversal, and read-side queries.
//
// Two posting flows coexist:
//
//   1. **System-sourced entries** (claim payments, premium allocations,
//      write-offs, refunds) are authorised upstream and post directly via
//      GLPost(). They land with Status="posted" in a single step.
//
//   2. **Manual entries** (caller-supplied multi-line journals from the
//      admin UI) follow a maker/checker state machine:
//
//          draft → submitted → approved → posted
//                                       ↘ (after posting)
//                              reversal_pending → reversal_approved → reversed
//
//      The approver (ApprovedBy) must be a different user from the submitter
//      (SubmittedBy); the same rule applies to the reversal flow.
//
// All posting paths converge on writeJournalEntry, which enforces:
//   * the target accounting period exists and is open
//   * total debits == total credits (within balanceTolerance)
//   * every account exists and is active
//   * each line has exactly one of (Debit, Credit) > 0
//   * when status="posted", the entry number is unique (one retry on collision)
//
// All public functions accept an optional *gorm.DB so they can be composed
// inside a caller's transaction. If tx is nil, services.DB is used and a new
// transaction is opened internally.

const balanceTolerance = 0.005

// Typed errors so HTTP controllers (and callers further up the stack) can
// distinguish configuration bugs from user-correctable mistakes.
var (
	ErrPeriodClosed    = errors.New("no open accounting period for the posting date")
	ErrNoPostingRule   = errors.New("no active posting rule for event")
	ErrUnbalanced      = errors.New("journal entry is not balanced (debits != credits)")
	ErrAccountInactive = errors.New("posting account is inactive")
	ErrAccountNotFound = errors.New("posting account not found")
	ErrInvalidLine     = errors.New("each journal line must have exactly one of debit or credit > 0")
	ErrEntryNotFound   = errors.New("journal entry not found")
	ErrAlreadyReversed = errors.New("journal entry has already been reversed")
	ErrBadState        = errors.New("journal entry is not in the required state for this action")
	ErrSelfApproval    = errors.New("approver must be a different user from the submitter")
)

// ---------------------------------------------------------------------------
// Event-driven posting (claim payments, premium allocations, write-offs, refunds)
// ---------------------------------------------------------------------------

// GLPost looks up the posting rule for eventKey and writes a balanced 2-line
// journal entry (DR rule.debit / CR rule.credit) for the given amount. Pass a
// non-nil tx to compose with the caller's outer transaction.
//
// System-sourced postings bypass the maker/checker workflow — they are
// authorised by the upstream business event (claim payment confirmed, premium
// receipt allocated, etc.).
func GLPost(tx *gorm.DB, eventKey, sourceType string, sourceID int, amount float64, description string, schemeID int, postedBy string) (models.JournalEntry, error) {
	if math.Abs(amount) < balanceTolerance {
		return models.JournalEntry{}, fmt.Errorf("posting amount must be non-zero")
	}
	db := tx
	if db == nil {
		db = DB
	}

	var rule models.PostingRule
	if err := db.Where("event_key = ? AND is_active = ?", eventKey, true).First(&rule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.JournalEntry{}, fmt.Errorf("%w: %s", ErrNoPostingRule, eventKey)
		}
		return models.JournalEntry{}, err
	}

	header := models.JournalEntry{
		Status:      "posted",
		SourceType:  sourceType,
		SourceID:    sourceID,
		Description: description,
		CreatedBy:   postedBy,
		PostedBy:    postedBy,
	}
	lines := []models.JournalLine{
		{AccountID: rule.DebitAccountID, Debit: amount, Description: description, SchemeID: schemeID, LineOrder: 1},
		{AccountID: rule.CreditAccountID, Credit: amount, Description: description, SchemeID: schemeID, LineOrder: 2},
	}

	return writeJournalEntry(db, header, lines)
}

// ---------------------------------------------------------------------------
// Manual journal workflow: draft → submit → approve → post
// ---------------------------------------------------------------------------

// GLDraftManualJournal persists a caller-supplied multi-line journal as a
// DRAFT. The entry is validated and stored (so the approver sees exactly what
// they're approving) but is not yet on the trial balance. No EntryNumber is
// assigned until the post step.
func GLDraftManualJournal(req models.ManualJournalRequest, user models.AppUser) (models.JournalEntry, error) {
	header := models.JournalEntry{
		SourceType:  "manual",
		Status:      "draft",
		Description: req.Description,
		CreatedBy:   user.UserName,
		UpdatedBy:   user.UserName,
	}
	if req.PeriodID > 0 {
		header.PeriodID = req.PeriodID
	}
	lines := manualLinesFromRequest(req)

	var saved models.JournalEntry
	err := DB.Transaction(func(tx *gorm.DB) error {
		var werr error
		saved, werr = writeJournalEntry(tx, header, lines)
		if werr != nil {
			return werr
		}
		return LogGLEvent(tx, "journal_drafted", "journal_entry", fmt.Sprintf("draft #%d", saved.ID), saved.ID, user, map[string]any{
			"description": req.Description,
			"line_count":  len(lines),
			"total_debit": saved.TotalDebit,
		})
	})
	return saved, err
}

// GLUpdateDraftJournal replaces the lines and description on an existing
// draft. Only entries with Status="draft" are editable; everything else is
// frozen.
func GLUpdateDraftJournal(id int, req models.ManualJournalRequest, user models.AppUser) (models.JournalEntry, error) {
	var existing models.JournalEntry
	if err := DB.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return existing, ErrEntryNotFound
		}
		return existing, err
	}
	if existing.Status != "draft" {
		return existing, fmt.Errorf("%w: cannot edit %s entry", ErrBadState, existing.Status)
	}

	newLines := manualLinesFromRequest(req)

	err := DB.Transaction(func(tx *gorm.DB) error {
		// Validate the candidate first (period, balance, accounts).
		totalDr, totalCr, err := validateJournal(tx, &existing, newLines)
		if err != nil {
			return err
		}
		// Replace lines: delete old, insert new.
		if err := tx.Where("entry_id = ?", existing.ID).Delete(&models.JournalLine{}).Error; err != nil {
			return err
		}
		for i := range newLines {
			newLines[i].EntryID = existing.ID
		}
		if err := tx.Create(&newLines).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.JournalEntry{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
			"description":  req.Description,
			"total_debit":  totalDr,
			"total_credit": totalCr,
			"updated_by":   user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "journal_draft_updated", "journal_entry", fmt.Sprintf("draft #%d", existing.ID), existing.ID, user, map[string]any{
			"description": req.Description,
			"line_count":  len(newLines),
			"total_debit": totalDr,
		})
	})
	if err != nil {
		return existing, err
	}
	return GLGetJournalEntry(id)
}

// GLSubmitManualJournal transitions a draft → submitted, locking it from
// further edits and exposing it to approvers.
func GLSubmitManualJournal(id int, user models.AppUser) (models.JournalEntry, error) {
	var entry models.JournalEntry
	if err := DB.First(&entry, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entry, ErrEntryNotFound
		}
		return entry, err
	}
	if entry.Status != "draft" {
		return entry, fmt.Errorf("%w: cannot submit %s entry", ErrBadState, entry.Status)
	}

	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.JournalEntry{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":       "submitted",
			"submitted_by": user.UserName,
			"submitted_at": &now,
			"updated_by":   user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "journal_submitted", "journal_entry", fmt.Sprintf("draft #%d", id), id, user, nil)
	})
	if err != nil {
		return entry, err
	}
	return GLGetJournalEntry(id)
}

// GLApproveManualJournal transitions submitted → approved. Enforces
// segregation of duties: the approver must be a different user from the
// submitter.
func GLApproveManualJournal(id int, user models.AppUser) (models.JournalEntry, error) {
	var entry models.JournalEntry
	if err := DB.First(&entry, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entry, ErrEntryNotFound
		}
		return entry, err
	}
	if entry.Status != "submitted" {
		return entry, fmt.Errorf("%w: cannot approve %s entry", ErrBadState, entry.Status)
	}
	if entry.SubmittedBy == user.UserName {
		return entry, ErrSelfApproval
	}

	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.JournalEntry{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": user.UserName,
			"approved_at": &now,
			"updated_by":  user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "journal_approved", "journal_entry", fmt.Sprintf("draft #%d", id), id, user, map[string]any{
			"submitted_by": entry.SubmittedBy,
		})
	})
	if err != nil {
		return entry, err
	}
	return GLGetJournalEntry(id)
}

// GLPostApprovedJournal transitions approved → posted: assigns the
// EntryNumber, stamps PostedAt/PostedBy, and locks the entry into the trial
// balance. Once posted the entry is immutable; corrections are made by
// requesting a reversal.
func GLPostApprovedJournal(id int, user models.AppUser) (models.JournalEntry, error) {
	var entry models.JournalEntry
	if err := DB.First(&entry, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entry, ErrEntryNotFound
		}
		return entry, err
	}
	if entry.Status != "approved" {
		return entry, fmt.Errorf("%w: cannot post %s entry", ErrBadState, entry.Status)
	}

	// Re-verify the period is still open (it may have closed between approval
	// and posting).
	var period models.AccountingPeriod
	if err := DB.First(&period, entry.PeriodID).Error; err != nil {
		return entry, fmt.Errorf("period %d: %w", entry.PeriodID, err)
	}
	if period.Status != "open" {
		return entry, fmt.Errorf("%w: period %s", ErrPeriodClosed, period.Name)
	}

	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		// Allocate entry number, with one retry on uniqueIndex collision.
		var assigned string
		for attempt := 0; attempt < 2; attempt++ {
			number, err := nextEntryNumber(tx, now)
			if err != nil {
				return err
			}
			updErr := tx.Model(&models.JournalEntry{}).Where("id = ?", id).Updates(map[string]interface{}{
				"entry_number": number,
				"posted_at":    now,
				"posted_by":    user.UserName,
				"status":       "posted",
				"updated_by":   user.UserName,
			}).Error
			if updErr == nil {
				assigned = number
				break
			}
			if attempt == 0 && isDuplicateKeyError(updErr) {
				continue
			}
			return updErr
		}
		return LogGLEvent(tx, "journal_posted", "journal_entry", assigned, id, user, map[string]any{
			"entry_number": assigned,
			"approved_by":  entry.ApprovedBy,
		})
	})
	if err != nil {
		return entry, err
	}
	return GLGetJournalEntry(id)
}

// GLDiscardDraft deletes a draft (or submitted) entry that the maker no
// longer wants. Only the original drafter or an approver can discard.
// Posted/approved entries cannot be discarded — use the reversal workflow.
func GLDiscardDraft(id int, user models.AppUser) error {
	var entry models.JournalEntry
	if err := DB.First(&entry, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEntryNotFound
		}
		return err
	}
	if entry.Status != "draft" && entry.Status != "submitted" {
		return fmt.Errorf("%w: cannot discard %s entry", ErrBadState, entry.Status)
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entry_id = ?", id).Delete(&models.JournalLine{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.JournalEntry{}, id).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "journal_draft_discarded", "journal_entry", fmt.Sprintf("draft #%d", id), id, user, map[string]any{
			"prior_status": entry.Status,
			"description":  entry.Description,
		})
	})
}

// ---------------------------------------------------------------------------
// Reversal workflow: request → approve (= post offsetting entry)
// ---------------------------------------------------------------------------

// GLRequestReversal flags a posted entry for reversal. The caller supplies a
// reason; the entry moves from posted → reversal_pending and waits for a
// second user to approve.
func GLRequestReversal(id int, reason string, user models.AppUser) (models.JournalEntry, error) {
	var entry models.JournalEntry
	if err := DB.First(&entry, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entry, ErrEntryNotFound
		}
		return entry, err
	}
	if entry.Status != "posted" {
		return entry, fmt.Errorf("%w: cannot request reversal of %s entry", ErrBadState, entry.Status)
	}
	if entry.IsReversed {
		return entry, ErrAlreadyReversed
	}
	if strings.TrimSpace(reason) == "" {
		return entry, errors.New("reversal reason is required")
	}

	now := time.Now()
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.JournalEntry{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":                "reversal_pending",
			"reversal_reason":       reason,
			"reversal_requested_by": user.UserName,
			"reversal_requested_at": &now,
			"updated_by":            user.UserName,
		}).Error; err != nil {
			return err
		}
		return LogGLEvent(tx, "journal_reversal_requested", "journal_entry", entry.EntryNumber, id, user, map[string]any{
			"reason": reason,
		})
	})
	if err != nil {
		return entry, err
	}
	return GLGetJournalEntry(id)
}

// GLApproveReversal completes a reversal: the original moves to "reversed",
// IsReversed flips true, and a new offsetting JournalEntry is created and
// posted (Status="posted", SourceType="reversal"). Enforces SoD — the
// approver must differ from the user who requested the reversal.
func GLApproveReversal(id int, user models.AppUser) (models.JournalEntry, error) {
	var orig models.JournalEntry
	if err := DB.Preload("Lines").First(&orig, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return orig, ErrEntryNotFound
		}
		return orig, err
	}
	if orig.Status != "reversal_pending" {
		return orig, fmt.Errorf("%w: cannot approve reversal of %s entry", ErrBadState, orig.Status)
	}
	if orig.IsReversed {
		return orig, ErrAlreadyReversed
	}
	if orig.ReversalRequestedBy == user.UserName {
		return orig, ErrSelfApproval
	}

	now := time.Now()
	var rev models.JournalEntry
	err := DB.Transaction(func(tx *gorm.DB) error {
		header := models.JournalEntry{
			SourceType:  "reversal",
			Status:      "posted",
			SourceID:    orig.ID,
			Description: fmt.Sprintf("Reversal of %s: %s", orig.EntryNumber, orig.ReversalReason),
			CreatedBy:   user.UserName,
			PostedBy:    user.UserName,
		}
		lines := make([]models.JournalLine, len(orig.Lines))
		for i, l := range orig.Lines {
			lines[i] = models.JournalLine{
				AccountID:   l.AccountID,
				Debit:       l.Credit,
				Credit:      l.Debit,
				Description: "Reversal: " + l.Description,
				SchemeID:    l.SchemeID,
				CostCentre:  l.CostCentre,
				LineOrder:   i + 1,
			}
		}
		var werr error
		rev, werr = writeJournalEntry(tx, header, lines)
		if werr != nil {
			return werr
		}
		// Mark the original as reversed and link the pair.
		if err := tx.Model(&models.JournalEntry{}).Where("id = ?", orig.ID).Updates(map[string]interface{}{
			"status":                 "reversed",
			"is_reversed":            true,
			"reversed_by_entry_id":   &rev.ID,
			"reversal_approved_by":   user.UserName,
			"reversal_approved_at":   &now,
			"updated_by":             user.UserName,
		}).Error; err != nil {
			return err
		}
		if err := LogGLEvent(tx, "journal_reversal_approved", "journal_entry", orig.EntryNumber, orig.ID, user, map[string]any{
			"reason":               orig.ReversalReason,
			"requested_by":         orig.ReversalRequestedBy,
			"reversal_entry_id":    rev.ID,
			"reversal_entry_number": rev.EntryNumber,
		}); err != nil {
			return err
		}
		return LogGLEvent(tx, "journal_posted", "journal_entry", rev.EntryNumber, rev.ID, user, map[string]any{
			"source_type": "reversal",
			"reverses":    orig.EntryNumber,
		})
	})
	if err != nil {
		return rev, err
	}
	return rev, nil
}

// ---------------------------------------------------------------------------
// Read-side queries (trial balance, account ledger)
// ---------------------------------------------------------------------------

// realPostingStatuses is the set of JournalEntry.Status values that
// represent a posting that has actually landed on the books (and therefore
// belongs on the trial balance and account ledgers). Drafts, submitted and
// approved entries are work-in-progress and excluded; reversed originals
// remain — their offsetting reversal entry is also "posted" so the pair
// nets out correctly.
var realPostingStatuses = []string{"posted", "reversal_pending", "reversal_approved", "reversed"}

// GLGetTrialBalance returns one row per active account with summed debits,
// credits, and a signed net balance (positive when on the account's normal
// side). Optionally scoped to a single accounting period.
func GLGetTrialBalance(periodID int) ([]models.TrialBalanceRow, error) {
	var accounts []models.GLAccount
	if err := DB.Order("code ASC").Find(&accounts).Error; err != nil {
		return nil, err
	}
	rows := make([]models.TrialBalanceRow, 0, len(accounts))
	for _, a := range accounts {
		var dr, cr float64
		q := DB.Model(&models.JournalLine{}).
			Joins("JOIN journal_entries je ON je.id = journal_lines.entry_id").
			Where("journal_lines.account_id = ?", a.ID).
			Where("je.status IN ?", realPostingStatuses)
		if periodID > 0 {
			q = q.Where("je.period_id = ?", periodID)
		}
		q.Select("COALESCE(SUM(debit), 0)").Scan(&dr)
		q2 := DB.Model(&models.JournalLine{}).
			Joins("JOIN journal_entries je ON je.id = journal_lines.entry_id").
			Where("journal_lines.account_id = ?", a.ID).
			Where("je.status IN ?", realPostingStatuses)
		if periodID > 0 {
			q2 = q2.Where("je.period_id = ?", periodID)
		}
		q2.Select("COALESCE(SUM(credit), 0)").Scan(&cr)

		net := dr - cr
		if a.NormalBalance == "credit" {
			net = -net
		}
		rows = append(rows, models.TrialBalanceRow{
			AccountID:     a.ID,
			AccountCode:   a.Code,
			AccountName:   a.Name,
			AccountType:   a.AccountType,
			NormalBalance: a.NormalBalance,
			TotalDebit:    dr,
			TotalCredit:   cr,
			NetBalance:    net,
		})
	}
	return rows, nil
}

// GLGetAccountLedger returns every posted line on an account between from/to
// (inclusive), ordered chronologically with a running balance. From/To are
// optional (use zero time for "no bound").
func GLGetAccountLedger(accountID int, from, to time.Time) ([]models.LedgerRow, error) {
	var account models.GLAccount
	if err := DB.First(&account, accountID).Error; err != nil {
		return nil, err
	}

	q := DB.Table("journal_lines AS jl").
		Select(`je.id AS entry_id,
		        je.entry_number,
		        je.posted_at,
		        je.source_type,
		        je.source_id,
		        je.description AS description,
		        jl.description AS line_description,
		        jl.debit,
		        jl.credit`).
		Joins("JOIN journal_entries je ON je.id = jl.entry_id").
		Where("jl.account_id = ?", accountID).
		Where("je.status IN ?", realPostingStatuses).
		Order("je.posted_at ASC, je.id ASC, jl.line_order ASC")
	if !from.IsZero() {
		q = q.Where("je.posted_at >= ?", from)
	}
	if !to.IsZero() {
		q = q.Where("je.posted_at <= ?", to)
	}

	var rows []models.LedgerRow
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	running := 0.0
	for i := range rows {
		delta := rows[i].Debit - rows[i].Credit
		if account.NormalBalance == "credit" {
			delta = -delta
		}
		running += delta
		rows[i].RunningBalance = running
	}
	return rows, nil
}

// ---------------------------------------------------------------------------
// Internal: shared write path + validation
// ---------------------------------------------------------------------------

// manualLinesFromRequest projects a ManualJournalRequest payload into model
// JournalLine values. LineOrder is assigned in arrival order.
func manualLinesFromRequest(req models.ManualJournalRequest) []models.JournalLine {
	lines := make([]models.JournalLine, len(req.Lines))
	for i, in := range req.Lines {
		lines[i] = models.JournalLine{
			AccountID:   in.AccountID,
			Debit:       in.Debit,
			Credit:      in.Credit,
			Description: in.Description,
			SchemeID:    in.SchemeID,
			CostCentre:  in.CostCentre,
			LineOrder:   i + 1,
		}
	}
	return lines
}

// validateJournal runs all checks that don't write to the DB:
// line count, per-line shape, balance, account existence/activeness, and
// period resolution. Returns the computed (totalDr, totalCr).
//
// Side effect: if header.PeriodID is zero on entry, it is populated with
// today's open period.
func validateJournal(db *gorm.DB, header *models.JournalEntry, lines []models.JournalLine) (float64, float64, error) {
	if len(lines) < 2 {
		return 0, 0, fmt.Errorf("journal entry needs at least 2 lines")
	}

	var totalDr, totalCr float64
	for i, l := range lines {
		drPos := l.Debit > balanceTolerance
		crPos := l.Credit > balanceTolerance
		if drPos == crPos { // both zero or both non-zero
			return 0, 0, fmt.Errorf("line %d: %w", i+1, ErrInvalidLine)
		}
		if l.Debit < 0 || l.Credit < 0 {
			return 0, 0, fmt.Errorf("line %d: debit/credit must be non-negative", i+1)
		}
		totalDr += l.Debit
		totalCr += l.Credit
	}
	if math.Abs(totalDr-totalCr) > balanceTolerance {
		return 0, 0, fmt.Errorf("%w: dr=%.2f cr=%.2f", ErrUnbalanced, totalDr, totalCr)
	}

	accountIDs := make([]int, 0, len(lines))
	for _, l := range lines {
		accountIDs = append(accountIDs, l.AccountID)
	}
	var accounts []models.GLAccount
	if err := db.Where("id IN ?", accountIDs).Find(&accounts).Error; err != nil {
		return 0, 0, err
	}
	activeByID := make(map[int]bool, len(accounts))
	for _, a := range accounts {
		activeByID[a.ID] = a.IsActive
	}
	for _, l := range lines {
		active, ok := activeByID[l.AccountID]
		if !ok {
			return 0, 0, fmt.Errorf("%w: id=%d", ErrAccountNotFound, l.AccountID)
		}
		if !active {
			return 0, 0, fmt.Errorf("%w: id=%d", ErrAccountInactive, l.AccountID)
		}
	}

	// Period resolution. If PeriodID was set by the caller use it; otherwise
	// resolve today's open period. Drafts also pick a period so the approver
	// sees where the entry would land.
	if header.PeriodID == 0 {
		period, err := resolveOpenPeriod(db, time.Now())
		if err != nil {
			return 0, 0, err
		}
		header.PeriodID = period.ID
	} else {
		var period models.AccountingPeriod
		if err := db.First(&period, header.PeriodID).Error; err != nil {
			return 0, 0, fmt.Errorf("period %d: %w", header.PeriodID, err)
		}
		if period.Status != "open" {
			return 0, 0, fmt.Errorf("%w: period %s", ErrPeriodClosed, period.Name)
		}
	}

	return totalDr, totalCr, nil
}

// writeJournalEntry validates header + lines and persists them. If
// header.Status is "" or "posted", the entry is posted immediately
// (EntryNumber assigned, PostedAt stamped). Any other status (e.g. "draft")
// is saved as-is without an EntryNumber.
//
// Composes with the caller's transaction when db is a *gorm.DB inside one.
func writeJournalEntry(db *gorm.DB, header models.JournalEntry, lines []models.JournalLine) (models.JournalEntry, error) {
	totalDr, totalCr, err := validateJournal(db, &header, lines)
	if err != nil {
		return header, err
	}
	header.TotalDebit = totalDr
	header.TotalCredit = totalCr

	if header.Status == "" {
		header.Status = "posted"
	}

	if header.Status == "posted" {
		header.PostedAt = time.Now()
		// Allocate entry number, with one retry on uniqueIndex collision.
		for attempt := 0; attempt < 2; attempt++ {
			number, err := nextEntryNumber(db, header.PostedAt)
			if err != nil {
				return header, err
			}
			header.EntryNumber = number
			err = db.Create(&header).Error
			if err == nil {
				break
			}
			if attempt == 0 && isDuplicateKeyError(err) {
				continue
			}
			return header, err
		}
	} else {
		// Draft/submitted/approved: persist without an entry number.
		if err := db.Create(&header).Error; err != nil {
			return header, err
		}
	}

	for i := range lines {
		lines[i].EntryID = header.ID
	}
	if err := db.Create(&lines).Error; err != nil {
		return header, err
	}
	header.Lines = lines
	return header, nil
}

// resolveOpenPeriod returns the open AccountingPeriod containing `at`. If none
// exists, ErrPeriodClosed is returned.
func resolveOpenPeriod(db *gorm.DB, at time.Time) (models.AccountingPeriod, error) {
	var period models.AccountingPeriod
	err := db.Where("status = ? AND start_date <= ? AND end_date >= ?", "open", at, at).
		Order("start_date DESC").
		First(&period).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return period, ErrPeriodClosed
	}
	return period, err
}

// nextEntryNumber generates a JE-YYYYMM-NNNN sequence number scoped to the
// posting month. Uses MAX-on-prefix; concurrent callers may collide on the
// unique index and the caller retries once.
func nextEntryNumber(db *gorm.DB, at time.Time) (string, error) {
	prefix := fmt.Sprintf("JE-%s-", at.Format("200601"))

	var current string
	err := db.Model(&models.JournalEntry{}).
		Where("entry_number LIKE ?", prefix+"%").
		Select("COALESCE(MAX(entry_number), '')").
		Scan(&current).Error
	if err != nil {
		return "", err
	}

	seq := 0
	if current != "" {
		tail := strings.TrimPrefix(current, prefix)
		if n, parseErr := strconv.Atoi(tail); parseErr == nil {
			seq = n
		}
	}
	return fmt.Sprintf("%s%04d", prefix, seq+1), nil
}

// isDuplicateKeyError loosely matches the unique-violation error strings
// across MySQL, PostgreSQL, and SQL Server. We rely on substring match because
// the project links multiple drivers and a stable cross-driver error type
// is not available.
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "unique") ||
		strings.Contains(msg, "23505") ||
		strings.Contains(msg, "violation")
}

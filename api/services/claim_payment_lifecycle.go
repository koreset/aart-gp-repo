package services

import (
	"api/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// This file owns the state-transition surface for a ClaimPaymentSchedule.
// The state machine is:
//
//   draft → claims_signed_off → finance_in_review → finance_first_authorised
//   → finance_second_authorised → submitted_to_bank → confirmed → archived
//
// Each transition is gated on:
//   1. The current status of the schedule.
//   2. The actor's permission slug (enforced by the controller / route).
//   3. The authority matrix (role + monetary threshold for the schedule total).
//
// Per-line transitions (verify / query / reject) are also implemented here.
// Queried and rejected lines are detached from the schedule and the underlying
// claim is returned to "approved" so it can be picked up on the next cut-off.

// buildApprovalReference snapshots who approved the claim, when, and at what
// authority level into a single string column. Pulled from the most recent
// status audit row whose new_status is "approved" (or the latest assessment
// if no such audit row exists).
func buildApprovalReference(claim models.GroupSchemeClaim, db *gorm.DB) string {
	if db == nil {
		db = DB
	}
	var audit models.GroupSchemeClaimStatusAudit
	err := db.Where("claim_id = ? AND new_status = ?", claim.ID, "approved").
		Order("changed_at DESC").
		First(&audit).Error
	if err == nil && audit.ChangedBy != "" {
		return fmt.Sprintf("%s @ %s", audit.ChangedBy, audit.ChangedAt.Format("2006-01-02 15:04"))
	}

	var assessment models.GroupSchemeClaimAssessment
	if err := db.Where("claim_id = ?", claim.ID).Order("creation_date DESC").First(&assessment).Error; err == nil {
		return fmt.Sprintf("%s @ %s", assessment.AssessorName, assessment.CreationDate.Format("2006-01-02 15:04"))
	}
	return ""
}

// QueryRequest is the inbound payload for a finance query or rejection.
type QueryRequest struct {
	ReasonCode string `json:"reason_code"`
	Notes      string `json:"notes"`
}

// recordAudit writes a PaymentScheduleAudit row inside the given tx.
func recordAudit(tx *gorm.DB, scheduleID int, from, to, actor, notes string) error {
	return tx.Create(&models.PaymentScheduleAudit{
		ScheduleID: scheduleID,
		FromStatus: from,
		ToStatus:   to,
		Actor:      actor,
		Notes:      notes,
	}).Error
}

// SignOffByHeadOfClaims transitions a schedule from draft → claims_signed_off.
// All line items must be in a non-terminal state (no queried/rejected leftovers
// from a previous review round); the authority matrix is checked against the
// schedule's NetTotal.
func SignOffByHeadOfClaims(scheduleID int, user models.AppUser) (models.ClaimPaymentSchedule, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return schedule, err
	}
	if schedule.Status != "draft" {
		return schedule, fmt.Errorf("schedule is not in draft status (current: %s)", schedule.Status)
	}
	for _, item := range schedule.Items {
		switch item.LineStatus {
		case "", "pending", "verified":
			// ok
		default:
			return schedule, fmt.Errorf("schedule has a %s line item — resolve before signing off", item.LineStatus)
		}
	}
	if err := RequireAuthority(user, AuthActionSignOffSchedule, schedule.NetTotal); err != nil {
		return schedule, err
	}

	now := time.Now()
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
			"status":                   "claims_signed_off",
			"head_of_claims_signed_by": user.UserName,
			"head_of_claims_signed_at": &now,
		}).Error; err != nil {
			return err
		}
		return recordAudit(tx, scheduleID, schedule.Status, "claims_signed_off", user.UserName, "Head of Claims sign-off")
	})
	if txErr != nil {
		return schedule, txErr
	}
	return GetPaymentSchedule(scheduleID)
}

// FinanceStartReview transitions a schedule from claims_signed_off →
// finance_in_review. This is what makes per-line verify/query/reject actions
// available to the finance team.
func FinanceStartReview(scheduleID int, user models.AppUser) (models.ClaimPaymentSchedule, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return schedule, err
	}
	if schedule.Status != "claims_signed_off" {
		return schedule, fmt.Errorf("schedule must be claims_signed_off to start finance review (current: %s)", schedule.Status)
	}
	if err := RequireAuthority(user, AuthActionFinanceReview, schedule.NetTotal); err != nil {
		return schedule, err
	}

	now := time.Now()
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
			"status":                    "finance_in_review",
			"finance_review_started_by": user.UserName,
			"finance_review_started_at": &now,
		}).Error; err != nil {
			return err
		}
		return recordAudit(tx, scheduleID, schedule.Status, "finance_in_review", user.UserName, "Finance review started")
	})
	if txErr != nil {
		return schedule, txErr
	}
	// Phase 3: re-run within-schedule duplicate detection now that the line
	// set may have shifted since creation (queries / rejections).
	_, _ = FlagDuplicateBeneficiariesForSchedule(scheduleID)
	return GetPaymentSchedule(scheduleID)
}

// VerifyLineItem marks a single line as "verified" by finance.
func VerifyLineItem(scheduleID, itemID int, user models.AppUser) (models.ClaimPaymentScheduleItem, error) {
	var item models.ClaimPaymentScheduleItem
	if err := DB.Where("id = ? AND schedule_id = ?", itemID, scheduleID).First(&item).Error; err != nil {
		return item, err
	}
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return item, err
	}
	if schedule.Status != "finance_in_review" {
		return item, fmt.Errorf("schedule is not in finance review (current: %s)", schedule.Status)
	}
	now := time.Now()
	if err := DB.Model(&item).Updates(map[string]interface{}{
		"line_status": "verified",
		"verified_by": user.UserName,
		"verified_at": &now,
	}).Error; err != nil {
		return item, err
	}
	item.LineStatus = "verified"
	item.VerifiedBy = user.UserName
	item.VerifiedAt = &now
	return item, nil
}

// QueryLineItem raises a finance query against a line. The line is removed
// from the schedule and the underlying claim is returned to "approved" so it
// can be picked up on the next cut-off after claims has resolved the issue.
// The schedule's totals are recomputed.
// ClaimStatusFinanceRejected is the terminal status applied to a claim when
// finance rejects its payment-schedule line. The claim is held out of the
// auto cut-off scheduler until an assessor explicitly acknowledges via the
// AcknowledgeFinanceRejection endpoint.
const ClaimStatusFinanceRejected = "finance_rejected"

// ClaimStatusOmbudClaim is applied to a claim that was previously declined
// and has been re-opened following an Ombudsman engagement. Used to surface
// these claims on the Regular Income Claims view so they get re-assessed.
const ClaimStatusOmbudClaim = "ombud_claim"

func QueryLineItem(scheduleID, itemID int, req QueryRequest, user models.AppUser) error {
	return removeLineItem(scheduleID, itemID, "queried", req, user)
}

// RejectLineItem rejects a line outright. Unlike a query (which sends the
// claim back to "approved" for the next cut-off), rejection marks the
// underlying claim as "finance_rejected" and snapshots the reason onto the
// claim so the claim-side banner can render without joining audits. The
// claim is held until an assessor acknowledges.
func RejectLineItem(scheduleID, itemID int, req QueryRequest, user models.AppUser) error {
	return removeLineItem(scheduleID, itemID, "rejected", req, user)
}

// removeLineItem is the shared body for query/reject. The line row is kept
// (LineStatus updated; ScheduleID retained for historical lookup) and the
// schedule totals are recomputed to exclude it. The underlying claim is
// returned to "approved" for a query, or moved to "finance_rejected" for a
// rejection.
func removeLineItem(scheduleID, itemID int, outcome string, req QueryRequest, user models.AppUser) error {
	if strings.TrimSpace(req.ReasonCode) == "" {
		return errors.New("reason_code is required")
	}

	var item models.ClaimPaymentScheduleItem
	if err := DB.Where("id = ? AND schedule_id = ?", itemID, scheduleID).First(&item).Error; err != nil {
		return err
	}
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return err
	}
	if schedule.Status != "finance_in_review" {
		return fmt.Errorf("schedule must be in finance review (current: %s)", schedule.Status)
	}

	now := time.Now()
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&item).Updates(map[string]interface{}{
			"line_status":       outcome,
			"query_reason_code": req.ReasonCode,
			"query_notes":       req.Notes,
			"queried_by":        user.UserName,
			"queried_at":        &now,
		}).Error; err != nil {
			return err
		}

		// Recompute schedule totals over the remaining "live" lines.
		var liveItems []models.ClaimPaymentScheduleItem
		if err := tx.Where("schedule_id = ? AND line_status IN ?", scheduleID, []string{"pending", "verified"}).Find(&liveItems).Error; err != nil {
			return err
		}
		var gross, ded, net float64
		for _, li := range liveItems {
			gross += li.GrossAmount
			ded += li.PremiumArrearsDeduction + li.PolicyLoanDeduction + li.TaxWithheld
			net += li.NetPayable
		}
		if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
			"total_amount":     gross,
			"gross_total":      gross,
			"deductions_total": ded,
			"net_total":        net,
			"claims_count":     len(liveItems),
		}).Error; err != nil {
			return err
		}

		// Transition the underlying claim depending on the line outcome.
		// Query → "approved" (claim eligible for next cut-off, current behaviour).
		// Reject → "finance_rejected" (held until an assessor acknowledges; the
		// reason is snapshotted onto the claim so the banner can render from
		// the claim row alone).
		var claim models.GroupSchemeClaim
		if err := tx.Select("id, status, claim_number").First(&claim, item.ClaimID).Error; err == nil {
			newStatus := "approved"
			statusMessage := fmt.Sprintf("Returned from payment schedule %s — %s (%s)", schedule.ScheduleNumber, outcome, req.ReasonCode)
			updates := map[string]interface{}{"status": newStatus}

			if outcome == "rejected" {
				newStatus = ClaimStatusFinanceRejected
				statusMessage = fmt.Sprintf("Finance rejected via schedule %s — %s: %s", schedule.ScheduleNumber, req.ReasonCode, req.Notes)
				rejectedAt := now
				updates = map[string]interface{}{
					"status":                            newStatus,
					"finance_rejected_at":               &rejectedAt,
					"finance_rejected_by":               user.UserName,
					"finance_rejection_reason_code":     req.ReasonCode,
					"finance_rejection_notes":           req.Notes,
					"finance_rejection_schedule_number": schedule.ScheduleNumber,
				}
			}

			if err := tx.Create(&models.GroupSchemeClaimStatusAudit{
				ClaimID:       claim.ID,
				OldStatus:     claim.Status,
				NewStatus:     newStatus,
				StatusMessage: statusMessage,
				ChangedBy:     user.UserName,
				ChangedAt:     now,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.GroupSchemeClaim{}).Where("id = ?", claim.ID).Updates(updates).Error; err != nil {
				return err
			}
		}

		// Log the query for reason-code analytics.
		if err := tx.Create(&models.ClaimPaymentScheduleQuery{
			ScheduleID:     scheduleID,
			ScheduleItemID: itemID,
			ClaimID:        item.ClaimID,
			ClaimNumber:    item.ClaimNumber,
			ReasonCode:     req.ReasonCode,
			Notes:          req.Notes,
			Outcome:        outcome,
			RaisedBy:       user.UserName,
		}).Error; err != nil {
			return err
		}

		return recordAudit(tx, scheduleID, schedule.Status, schedule.Status, user.UserName, fmt.Sprintf("Line %s (%s): %s", item.ClaimNumber, outcome, req.ReasonCode))
	})
}

// AcknowledgeFinanceRejection moves a claim out of finance_rejected and back
// to "draft" — separation of duties means the capturer (claims:lodge) edits
// the fields that finance flagged (banking, amount, etc.) and then submits
// the claim for assessment again. The assessor then re-runs assessment and
// approval. This keeps each role bounded to its own workflow stage.
//
// The finance_rejection_* snapshot columns are intentionally retained so the
// capturer can see exactly what finance flagged while editing. They're
// cleared the next time the claim is re-approved (handled in
// UpdateClaimAssessment).
//
// Also marks any open ClaimPaymentScheduleQuery rows tied to this rejection
// as "resolved" with the actor / timestamp, so the analytics view stays
// honest.
func AcknowledgeFinanceRejection(claimID int, user models.AppUser) (models.GroupSchemeClaim, error) {
	var claim models.GroupSchemeClaim
	if err := DB.First(&claim, claimID).Error; err != nil {
		return claim, err
	}
	if !strings.EqualFold(claim.Status, ClaimStatusFinanceRejected) {
		return claim, fmt.Errorf("claim %s is not finance-rejected (current status: %s)", claim.ClaimNumber, claim.Status)
	}

	now := time.Now()
	if err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.GroupSchemeClaim{}).
			Where("id = ?", claim.ID).
			Update("status", "draft").Error; err != nil {
			return err
		}
		if err := tx.Create(&models.GroupSchemeClaimStatusAudit{
			ClaimID:       claim.ID,
			OldStatus:     ClaimStatusFinanceRejected,
			NewStatus:     "draft",
			StatusMessage: "Finance rejection acknowledged — claim returned to capturer for amendment and re-submission for assessment",
			ChangedBy:     user.UserName,
			ChangedAt:     now,
		}).Error; err != nil {
			return err
		}
		// Best-effort: close any open queries for this claim raised against
		// the rejection. We resolve every open "rejected" outcome row so the
		// analytics view shows the loop closed; if no rows exist the update
		// is a no-op.
		if err := tx.Model(&models.ClaimPaymentScheduleQuery{}).
			Where("claim_id = ? AND outcome = ?", claim.ID, "rejected").
			Where("(resolved_at IS NULL OR resolved_at = ?)", time.Time{}).
			Updates(map[string]interface{}{
				"resolved_by": user.UserName,
				"resolved_at": &now,
			}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return claim, err
	}

	// Reload so the controller returns the updated row.
	if err := DB.First(&claim, claimID).Error; err != nil {
		return claim, err
	}
	return claim, nil
}

// ResolveScheduleQuery marks an open query as resolved, stamping the resolver
// and storing their reply. Either team can resolve a query they did not raise
// — finance closes out a line query, and the same handler is reused by finance
// to respond to a claims follow-up.
func ResolveScheduleQuery(queryID int, response string, user models.AppUser) (models.ClaimPaymentScheduleQuery, error) {
	var query models.ClaimPaymentScheduleQuery
	response = strings.TrimSpace(response)
	if response == "" {
		return query, errors.New("response is required")
	}
	if err := DB.First(&query, queryID).Error; err != nil {
		return query, err
	}
	if query.Outcome == "resolved" || query.Outcome == "cancelled" {
		return query, fmt.Errorf("query is already %s", query.Outcome)
	}
	now := time.Now()
	if err := DB.Model(&query).Updates(map[string]interface{}{
		"outcome":          "resolved",
		"resolution_notes": response,
		"resolved_by":      user.UserName,
		"resolved_at":      &now,
	}).Error; err != nil {
		return query, err
	}
	return query, DB.First(&query, queryID).Error
}

// GetScheduleQueries returns all finance queries raised against a schedule.
func GetScheduleQueries(scheduleID int) ([]models.ClaimPaymentScheduleQuery, error) {
	var rows []models.ClaimPaymentScheduleQuery
	err := DB.Where("schedule_id = ?", scheduleID).Order("raised_at DESC").Find(&rows).Error
	return rows, err
}

// FinanceFirstAuthorise transitions finance_in_review → finance_first_authorised.
// Requires all line items to be in a terminal state (no `pending` lines), and
// the actor must meet the authority matrix for the schedule's NetTotal.
func FinanceFirstAuthorise(scheduleID int, user models.AppUser) (models.ClaimPaymentSchedule, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return schedule, err
	}
	if schedule.Status != "finance_in_review" {
		return schedule, fmt.Errorf("schedule must be finance_in_review (current: %s)", schedule.Status)
	}
	for _, item := range schedule.Items {
		if item.LineStatus == "" || item.LineStatus == "pending" {
			return schedule, fmt.Errorf("line %s is still pending — verify or query before authorising", item.ClaimNumber)
		}
	}
	if err := RequireAuthority(user, AuthActionAuthoriseFirst, schedule.NetTotal); err != nil {
		return schedule, err
	}
	// Funds-available / daily payment limit (Phase 2). Empty license string
	// matches the singleton install row.
	if err := CheckDailyPaymentLimit("", schedule); err != nil {
		return schedule, err
	}
	// Phase 3 guards. Each returns a clear error pointing at the offending
	// lines so finance knows exactly what to resolve.
	if blockers, err := blockingSanctionsItems(scheduleID); err == nil && len(blockers) > 0 {
		return schedule, fmt.Errorf("sanctions / PEP screening outstanding for: %s", strings.Join(blockers, ", "))
	}
	if dupes, err := outstandingDuplicateBeneficiaries(scheduleID); err == nil && len(dupes) > 0 {
		return schedule, fmt.Errorf("duplicate beneficiary across this schedule on: %s — review and clear each flag", strings.Join(dupes, ", "))
	}
	if drifts, err := outstandingAmountDrifts(scheduleID); err == nil && len(drifts) > 0 {
		return schedule, fmt.Errorf("amount drift outstanding on: %s — acknowledge or query each line before authorising", strings.Join(drifts, ", "))
	}
	if missing, err := outstandingReinsuranceRecoveries(scheduleID); err == nil && len(missing) > 0 {
		return schedule, fmt.Errorf("reinsurance recovery not yet raised for: %s — record the recovery before authorising", strings.Join(missing, ", "))
	}

	now := time.Now()
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
			"status":                  "finance_first_authorised",
			"finance_first_auth_by":   user.UserName,
			"finance_first_auth_at":   &now,
		}).Error; err != nil {
			return err
		}
		return recordAudit(tx, scheduleID, schedule.Status, "finance_first_authorised", user.UserName, "First finance authorisation")
	})
	if txErr != nil {
		return schedule, txErr
	}
	return GetPaymentSchedule(scheduleID)
}

// FinanceSecondAuthorise transitions finance_first_authorised →
// finance_second_authorised. Enforces the four-eyes rule: the second
// authoriser must be a different user from the first.
func FinanceSecondAuthorise(scheduleID int, user models.AppUser) (models.ClaimPaymentSchedule, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return schedule, err
	}
	if schedule.Status != "finance_first_authorised" {
		return schedule, fmt.Errorf("schedule must be finance_first_authorised (current: %s)", schedule.Status)
	}
	if schedule.FinanceFirstAuthBy != "" && schedule.FinanceFirstAuthBy == user.UserName {
		return schedule, errors.New("second authoriser must be a different user from the first authoriser")
	}
	if err := RequireAuthority(user, AuthActionAuthoriseSecond, schedule.NetTotal); err != nil {
		return schedule, err
	}

	now := time.Now()
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
			"status":                  "finance_second_authorised",
			"finance_second_auth_by":  user.UserName,
			"finance_second_auth_at":  &now,
		}).Error; err != nil {
			return err
		}
		return recordAudit(tx, scheduleID, schedule.Status, "finance_second_authorised", user.UserName, "Second finance authorisation")
	})
	if txErr != nil {
		return schedule, txErr
	}
	return GetPaymentSchedule(scheduleID)
}

// ArchiveSchedule moves a confirmed schedule to "archived". Archived schedules
// remain readable but are hidden from default list views.
func ArchiveSchedule(scheduleID int, user models.AppUser) (models.ClaimPaymentSchedule, error) {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return schedule, err
	}
	if schedule.Status != "confirmed" {
		return schedule, fmt.Errorf("only confirmed schedules can be archived (current: %s)", schedule.Status)
	}
	if err := RequireAuthority(user, AuthActionArchive, schedule.NetTotal); err != nil {
		return schedule, err
	}

	now := time.Now()
	txErr := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ClaimPaymentSchedule{}).Where("id = ?", scheduleID).Updates(map[string]interface{}{
			"status":       "archived",
			"archived_by":  user.UserName,
			"archived_at":  &now,
		}).Error; err != nil {
			return err
		}
		return recordAudit(tx, scheduleID, schedule.Status, "archived", user.UserName, "Schedule archived")
	})
	if txErr != nil {
		return schedule, txErr
	}
	return GetPaymentSchedule(scheduleID)
}

// GetScheduleAuditTrail returns the full state-transition history of a schedule.
func GetScheduleAuditTrail(scheduleID int) ([]models.PaymentScheduleAudit, error) {
	var rows []models.PaymentScheduleAudit
	err := DB.Where("schedule_id = ?", scheduleID).Order("changed_at ASC").Find(&rows).Error
	return rows, err
}

// ──────────────────────────────────────────────
// Authority Matrix CRUD
// ──────────────────────────────────────────────

// ListAuthorityMatrix returns all authority matrix rows.
func ListAuthorityMatrix() ([]models.AuthorityMatrix, error) {
	var rows []models.AuthorityMatrix
	err := DB.Order("role ASC, action ASC, min_amount ASC").Find(&rows).Error
	return rows, err
}

// CreateAuthorityMatrixRow inserts a new matrix row.
func CreateAuthorityMatrixRow(row models.AuthorityMatrix, user models.AppUser) (models.AuthorityMatrix, error) {
	if strings.TrimSpace(row.Role) == "" || strings.TrimSpace(row.Action) == "" {
		return row, errors.New("role and action are required")
	}
	row.CreatedBy = user.UserName
	if err := DB.Create(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

// UpdateAuthorityMatrixRow patches an existing matrix row.
func UpdateAuthorityMatrixRow(id int, patch models.AuthorityMatrix) (models.AuthorityMatrix, error) {
	var row models.AuthorityMatrix
	if err := DB.First(&row, id).Error; err != nil {
		return row, err
	}
	updates := map[string]interface{}{
		"role":       patch.Role,
		"action":     patch.Action,
		"min_amount": patch.MinAmount,
		"max_amount": patch.MaxAmount,
		"is_active":  patch.IsActive,
	}
	if err := DB.Model(&row).Updates(updates).Error; err != nil {
		return row, err
	}
	return row, DB.First(&row, id).Error
}

// DeleteAuthorityMatrixRow removes a matrix row.
func DeleteAuthorityMatrixRow(id int) error {
	return DB.Delete(&models.AuthorityMatrix{}, id).Error
}

// ──────────────────────────────────────────────
// Reinsurance recovery (Phase 3)
// ──────────────────────────────────────────────

// SetReinsuranceRecoveryRequest is the inbound payload for either flagging
// a line as needing recovery or recording that the recovery has been raised.
type SetReinsuranceRecoveryRequest struct {
	Required bool    `json:"required"`
	Amount   float64 `json:"amount"`
}

// SetReinsuranceRecovery flags a line as requiring a reinsurance recovery
// (or unflags it). Called by claims before finance review starts.
func SetReinsuranceRecovery(scheduleID, itemID int, req SetReinsuranceRecoveryRequest, user models.AppUser) (models.ClaimPaymentScheduleItem, error) {
	var item models.ClaimPaymentScheduleItem
	if err := DB.Where("id = ? AND schedule_id = ?", itemID, scheduleID).First(&item).Error; err != nil {
		return item, err
	}
	updates := map[string]interface{}{
		"reinsurance_recovery_required": req.Required,
		"reinsurance_recovery_amount":   req.Amount,
	}
	if !req.Required {
		// Unflagging also clears any prior "raised" record.
		updates["reinsurance_recovery_raised_by"] = ""
		updates["reinsurance_recovery_raised_at"] = nil
	}
	if err := DB.Model(&item).Updates(updates).Error; err != nil {
		return item, err
	}
	return getScheduleItem(itemID)
}

// ConfirmReinsuranceRecoveryRaised records that finance has confirmed the
// recovery has been raised with the reinsurer. Must be called before first
// authorisation when reinsurance_recovery_required is true.
func ConfirmReinsuranceRecoveryRaised(scheduleID, itemID int, user models.AppUser) (models.ClaimPaymentScheduleItem, error) {
	var item models.ClaimPaymentScheduleItem
	if err := DB.Where("id = ? AND schedule_id = ?", itemID, scheduleID).First(&item).Error; err != nil {
		return item, err
	}
	if !item.ReinsuranceRecoveryRequired {
		return item, errors.New("this line is not flagged for reinsurance recovery")
	}
	now := time.Now()
	if err := DB.Model(&item).Updates(map[string]interface{}{
		"reinsurance_recovery_raised_by": user.UserName,
		"reinsurance_recovery_raised_at": &now,
	}).Error; err != nil {
		return item, err
	}
	return getScheduleItem(itemID)
}

func getScheduleItem(id int) (models.ClaimPaymentScheduleItem, error) {
	var item models.ClaimPaymentScheduleItem
	err := DB.First(&item, id).Error
	return item, err
}

// DiscardDraftSchedule hard-deletes a draft schedule and returns each of its
// line-item claims to "approved" so the next cut-off picks them up. Only drafts
// can be discarded — once a schedule has been signed off, finance owns it.
func DiscardDraftSchedule(scheduleID int, user models.AppUser) error {
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return err
	}
	if schedule.Status != "draft" {
		return fmt.Errorf("only draft schedules can be discarded (current: %s)", schedule.Status)
	}

	now := time.Now()
	return DB.Transaction(func(tx *gorm.DB) error {
		var items []models.ClaimPaymentScheduleItem
		if err := tx.Where("schedule_id = ?", scheduleID).Find(&items).Error; err != nil {
			return err
		}

		// Return each underlying claim to "approved" with an audit row.
		for _, item := range items {
			var claim models.GroupSchemeClaim
			if err := tx.Select("id, status, claim_number").First(&claim, item.ClaimID).Error; err != nil {
				continue
			}
			if err := tx.Create(&models.GroupSchemeClaimStatusAudit{
				ClaimID:       claim.ID,
				OldStatus:     claim.Status,
				NewStatus:     "approved",
				StatusMessage: fmt.Sprintf("Draft payment schedule %s discarded", schedule.ScheduleNumber),
				ChangedBy:     user.UserName,
				ChangedAt:     now,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.GroupSchemeClaim{}).Where("id = ?", claim.ID).Update("status", "approved").Error; err != nil {
				return err
			}
		}

		// Record the discard in the schedule audit trail before deleting rows.
		if err := recordAudit(tx, scheduleID, schedule.Status, "discarded", user.UserName, "Draft discarded by claims user"); err != nil {
			return err
		}

		// Wipe related rows. Drafts shouldn't have queries or ACB records, but
		// these deletes are safe no-ops if there aren't any.
		if err := tx.Where("schedule_id = ?", scheduleID).Delete(&models.ClaimPaymentScheduleQuery{}).Error; err != nil {
			return err
		}
		if err := tx.Where("schedule_id = ?", scheduleID).Delete(&models.ClaimPaymentScheduleItem{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.ClaimPaymentSchedule{}, scheduleID).Error; err != nil {
			return err
		}
		return nil
	})
}

// RaiseClaimsFollowup records a claims-side follow-up note on a submitted
// payment schedule. Reuses the queries table so the existing Queries panel
// surfaces it to finance, who can resolve it with the same controls used for
// finance-raised line queries. The follow-up is scoped to the schedule (no
// line item), so ScheduleItemID is left at zero.
func RaiseClaimsFollowup(scheduleID int, notes string, user models.AppUser) (models.ClaimPaymentScheduleQuery, error) {
	var row models.ClaimPaymentScheduleQuery
	notes = strings.TrimSpace(notes)
	if notes == "" {
		return row, errors.New("notes are required")
	}
	schedule, err := GetPaymentSchedule(scheduleID)
	if err != nil {
		return row, err
	}
	if schedule.Status == "draft" {
		return row, errors.New("follow-ups can only be raised after the schedule has been signed off")
	}
	row = models.ClaimPaymentScheduleQuery{
		ScheduleID: scheduleID,
		ReasonCode: "claims_followup",
		Notes:      notes,
		Outcome:    "open",
		RaisedBy:   user.UserName,
	}
	if err := DB.Create(&row).Error; err != nil {
		return row, err
	}
	return row, nil
}

// outstandingReinsuranceRecoveries returns the claim numbers of lines that
// require a recovery to be raised but where finance hasn't yet confirmed it.
func outstandingReinsuranceRecoveries(scheduleID int) ([]string, error) {
	var rows []struct{ ClaimNumber string }
	err := DB.Model(&models.ClaimPaymentScheduleItem{}).
		Select("claim_number").
		Where("schedule_id = ? AND reinsurance_recovery_required = ?", scheduleID, true).
		Where("reinsurance_recovery_raised_at IS NULL").
		Where("line_status IN ?", []string{"pending", "verified"}).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.ClaimNumber)
	}
	return out, nil
}

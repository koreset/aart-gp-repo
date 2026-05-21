package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api/models"
	"api/services"
	"api/services/bav"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// reverifyMinInterval guards against accidental double-clicks racking up
// provider costs. Same window the UI surfaces as a disabled-state tooltip.
const reverifyMinInterval = time.Hour

// GetLineBankVerificationStatus handles
// GET /group-pricing/claims/payment-schedules/:schedule_id/items/:item_id/bank-verification
//
// Returns the digested BAV status for the schedule line — verified vs stale,
// last attempt timestamp, the next attempt number the UI should send when the
// user clicks Re-verify.
func GetLineBankVerificationStatus(c *gin.Context) {
	scheduleID, itemID, ok := parseScheduleItemParams(c)
	if !ok {
		return
	}
	item, err := loadScheduleItem(scheduleID, itemID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	bankingChanged := services.UnmarshalRiskFlags(item.RiskFlags).BankingChange30d
	status, err := services.LatestBankVerification(services.DB, item.ClaimID, bankingChanged)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, status)
}

// ReverifyLineBankAccount handles
// POST /group-pricing/claims/payment-schedules/:schedule_id/items/:item_id/bank-verification/reverify
//
// Fires a fresh BAV provider call against the schedule item's snapshotted
// banking fields, incrementing Attempt so the registry's idempotency cache
// is bypassed. Returns 429 when the previous call was made within the
// reverifyMinInterval window — the audit log is the rate-limit oracle.
func ReverifyLineBankAccount(c *gin.Context) {
	scheduleID, itemID, ok := parseScheduleItemParams(c)
	if !ok {
		return
	}
	item, err := loadScheduleItem(scheduleID, itemID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	// Server-side cost guard.
	withinWindow, lastAt, err := services.LastReverifyWithinWindow(services.DB, item.ClaimID, reverifyMinInterval)
	if err != nil {
		InternalError(c, err)
		return
	}
	if withinWindow && lastAt != nil {
		c.JSON(http.StatusTooManyRequests, models.PremiumResponse{
			Success: false,
			Message: fmt.Sprintf("a bank verification call was made for this claim at %s — wait before re-trying to avoid duplicate provider charges", lastAt.Format(time.RFC3339)),
		})
		return
	}

	// Pull the live claim to source the name fields the BAV provider needs.
	var claim models.GroupSchemeClaim
	if err := services.DB.First(&claim, item.ClaimID).Error; err != nil {
		NotFound(c, "claim not found")
		return
	}

	// Compute next attempt from the log so the idempotency key changes.
	var priorCount int64
	_ = services.DB.Model(&models.BAVVerificationLog{}).Where("claim_id = ?", item.ClaimID).Count(&priorCount).Error

	firstName, surname := splitName(claim.ClaimantName)
	req := bav.VerifyRequest{
		FirstName:         firstName,
		Surname:           surname,
		IdentityNumber:    claim.ClaimantIDNumber,
		IdentityType:      resolveIdentityType(claim),
		BankAccountNumber: item.BankAccountNumber,
		BankBranchCode:    item.BankBranchCode,
		BankAccountType:   item.BankAccountType,
		ClaimID:           &claim.ID,
		Attempt:           int(priorCount) + 1,
	}

	result, err := bav.Verify(c.Request.Context(), req)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, result)
}

// acknowledgeAmountDriftRequest is the (currently empty) body envelope for
// the Acknowledge handler. Defined as a struct so a future "reason" field
// slots in without breaking the route shape.
type acknowledgeAmountDriftRequest struct {
	// Reserved for an optional note finance might attach to the
	// acknowledgement (Phase 6+). Empty for now.
	Note string `json:"note"`
}

// AcknowledgeAmountDrift handles
// POST /group-pricing/claims/payment-schedules/:schedule_id/items/:item_id/amount-drift/acknowledge
//
// Stamps the schedule line's amount_drift_resolved trio so the warning icon
// clears in the UI. Required before first authorisation when drift is
// non-zero (the FinanceFirstAuthorise gate enforces this).
func AcknowledgeAmountDrift(c *gin.Context) {
	scheduleID, itemID, ok := parseScheduleItemParams(c)
	if !ok {
		return
	}
	var req acknowledgeAmountDriftRequest
	_ = c.ShouldBindJSON(&req)

	item, err := loadScheduleItem(scheduleID, itemID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	delta, drifted := services.ComputeAmountDrift(item)
	if !drifted {
		BadRequestMsg(c, "no amount drift on this line — nothing to acknowledge")
		return
	}

	user := c.MustGet("user").(models.AppUser)
	now := time.Now()
	updates := map[string]interface{}{
		"amount_drift_resolved":    true,
		"amount_drift_resolved_by": user.UserName,
		"amount_drift_resolved_at": &now,
	}
	if err := services.DB.Model(&models.ClaimPaymentScheduleItem{}).
		Where("id = ?", itemID).
		Updates(updates).Error; err != nil {
		InternalError(c, err)
		return
	}
	OK(c, gin.H{
		"item_id":      itemID,
		"drift_amount": delta,
		"resolved_by":  user.UserName,
		"resolved_at":  now,
	})
}

// parseScheduleItemParams reads :schedule_id and :item_id off the request,
// writing a 400 on failure. Returns ok=false when the caller should abort.
func parseScheduleItemParams(c *gin.Context) (scheduleID, itemID int, ok bool) {
	var err error
	scheduleID, err = strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return 0, 0, false
	}
	itemID, err = strconv.Atoi(c.Param("item_id"))
	if err != nil {
		BadRequestMsg(c, "invalid item_id")
		return 0, 0, false
	}
	return scheduleID, itemID, true
}

// loadScheduleItem fetches the line by id, ensuring it belongs to the
// schedule referenced by the URL. Prevents the wrong-schedule URL crafting
// edge case where an item id from schedule A is requested under schedule B.
func loadScheduleItem(scheduleID, itemID int) (models.ClaimPaymentScheduleItem, error) {
	var item models.ClaimPaymentScheduleItem
	err := services.DB.First(&item, itemID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, fmt.Errorf("schedule item not found")
		}
		return item, err
	}
	if item.ScheduleID != scheduleID {
		return item, fmt.Errorf("schedule item does not belong to schedule %d", scheduleID)
	}
	return item, nil
}

// resolveIdentityType prefers the value the user explicitly captured on the
// claim (ClaimantIdentityType) and falls back to the 13-digit heuristic for
// legacy rows where the column is still empty.
func resolveIdentityType(claim models.GroupSchemeClaim) string {
	if t := strings.TrimSpace(claim.ClaimantIdentityType); t != "" {
		return t
	}
	return detectIdentityType(claim.ClaimantIDNumber)
}

// detectIdentityType mirrors the heuristic the claim registration form used
// before the explicit identity_type column was added: 13 digits → South
// African ID number, otherwise → Passport. Kept as the fallback path for
// legacy claims where ClaimantIdentityType is empty.
func detectIdentityType(idNumber string) string {
	trimmed := strings.TrimSpace(idNumber)
	if len(trimmed) == 13 {
		allDigits := true
		for _, r := range trimmed {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return "IDNumber"
		}
	}
	return "Passport"
}

// splitName is a best-effort splitter for the claimant's "first name surname"
// string the BAV provider expects as two fields. Returns the whole string as
// FirstName when no whitespace is present.
func splitName(full string) (first, surname string) {
	for i, r := range full {
		if r == ' ' || r == '\t' {
			return full[:i], full[i+1:]
		}
	}
	return full, ""
}

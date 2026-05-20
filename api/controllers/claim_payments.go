package controllers

import (
	"api/models"
	"api/services"
	"api/services/bav"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePaymentSchedule handles POST /group-pricing/claims/payment-schedules
// Body: { "claim_ids": [1,2,3], "description": "..." }
func CreatePaymentSchedule(c *gin.Context) {
	var req services.CreatePaymentScheduleRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}

	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.CreatePaymentSchedule(req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	Created(c, schedule)
}

// GetPaymentSchedules handles GET /group-pricing/claims/payment-schedules.
// Pass ?include_archived=1 to also surface archived schedules.
// Pass ?scope=claims to limit to schedules the current user created,
// or ?scope=finance to exclude drafts (the finance hub list).
func GetPaymentSchedules(c *gin.Context) {
	includeArchived := c.Query("include_archived") == "1" || c.Query("include_archived") == "true"
	opts := services.GetPaymentSchedulesOptions{
		IncludeArchived: includeArchived,
	}
	switch c.Query("scope") {
	case string(services.ScopeClaims):
		opts.Scope = services.ScopeClaims
		if user, ok := c.MustGet("user").(models.AppUser); ok {
			opts.CreatedBy = user.UserName
		}
	case string(services.ScopeFinance):
		opts.Scope = services.ScopeFinance
	}
	schedules, err := services.GetPaymentSchedules(opts)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, schedules)
}

// GetPaymentSchedule handles GET /group-pricing/claims/payment-schedules/:schedule_id
func GetPaymentSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}

	schedule, err := services.GetPaymentSchedule(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "payment schedule not found")
			return
		}
		InternalError(c, err)
		return
	}
	OK(c, schedule)
}

// UpdatePaymentScheduleNotes handles PATCH /group-pricing/claims/payment-schedules/:schedule_id/notes
// Body: { "notes": "free text ≤30 chars" }
func UpdatePaymentScheduleNotes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}
	var req models.UpdatePaymentScheduleNotesRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	schedule, err := services.UpdatePaymentScheduleNotes(id, req.Notes)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "payment schedule not found")
			return
		}
		InternalError(c, err)
		return
	}
	OK(c, schedule)
}

// ExportPaymentScheduleCSV handles GET /group-pricing/claims/payment-schedules/:schedule_id/export
// Returns a CSV file download of the payment schedule.
func ExportPaymentScheduleCSV(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}

	user := c.MustGet("user").(models.AppUser)
	data, filename, err := services.ExportPaymentScheduleCSV(id, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "payment schedule not found")
			return
		}
		InternalError(c, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", data)
}

// UploadPaymentProof handles POST /group-pricing/claims/payment-schedules/:schedule_id/proof
// Expects multipart/form-data with a "file" field and optional "notes" text field.
func UploadPaymentProof(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		BadRequestMsg(c, "a file must be uploaded in the 'file' field")
		return
	}

	notes := c.PostForm("notes")
	user := c.MustGet("user").(models.AppUser)

	proof, err := services.UploadPaymentProof(id, fileHeader, notes, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "payment schedule not found")
			return
		}
		BadRequest(c, err)
		return
	}
	Created(c, proof)
}

// GetPaymentProofs handles GET /group-pricing/claims/payment-schedules/:schedule_id/proof
func GetPaymentProofs(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}

	proofs, err := services.GetPaymentProofs(id)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, proofs)
}

// DownloadPaymentProof handles GET /group-pricing/claims/payment-schedules/proof/:proof_id/download
func DownloadPaymentProof(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("proof_id"))
	if err != nil {
		BadRequestMsg(c, "invalid proof_id")
		return
	}

	data, contentType, filename, err := services.DownloadPaymentProof(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "proof of payment not found")
			return
		}
		InternalError(c, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// ──────────────────────────────────────────────
// Bank Profile handlers
// ──────────────────────────────────────────────

// CreateBankProfile handles POST /group-pricing/claims/bank-profiles
func CreateBankProfile(c *gin.Context) {
	var req models.CreateBankProfileRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	profile, err := services.CreateBankProfile(req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	Created(c, profile)
}

// GetBankProfiles handles GET /group-pricing/claims/bank-profiles
func GetBankProfiles(c *gin.Context) {
	profiles, err := services.GetBankProfiles()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, profiles)
}

// GetBankProfile handles GET /group-pricing/claims/bank-profiles/:profile_id
func GetBankProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("profile_id"))
	if err != nil {
		BadRequestMsg(c, "invalid profile_id")
		return
	}
	profile, err := services.GetBankProfile(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "bank profile not found")
			return
		}
		InternalError(c, err)
		return
	}
	OK(c, profile)
}

// UpdateBankProfile handles PATCH /group-pricing/claims/bank-profiles/:profile_id
func UpdateBankProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("profile_id"))
	if err != nil {
		BadRequestMsg(c, "invalid profile_id")
		return
	}
	var req models.UpdateBankProfileRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	profile, err := services.UpdateBankProfile(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "bank profile not found")
			return
		}
		BadRequest(c, err)
		return
	}
	OK(c, profile)
}

// DeleteBankProfile handles DELETE /group-pricing/claims/bank-profiles/:profile_id
func DeleteBankProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("profile_id"))
	if err != nil {
		BadRequestMsg(c, "invalid profile_id")
		return
	}
	if err := services.DeleteBankProfile(id); err != nil {
		InternalError(c, err)
		return
	}
	OK(c, gin.H{"message": "bank profile deleted"})
}

// ──────────────────────────────────────────────
// ACB File handlers
// ──────────────────────────────────────────────

// GenerateACBFile handles POST /group-pricing/claims/payment-schedules/:schedule_id/acb
func GenerateACBFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}
	var req models.GenerateACBRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	fileRecord, err := services.GenerateACBFile(id, req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	Created(c, fileRecord)
}

// GetACBFileRecords handles GET /group-pricing/claims/payment-schedules/:schedule_id/acb-files
func GetACBFileRecords(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}
	records, err := services.GetACBFileRecords(id)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, records)
}

// DownloadACBFile handles GET /group-pricing/claims/acb-files/:acb_file_id/download
func DownloadACBFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("acb_file_id"))
	if err != nil {
		BadRequestMsg(c, "invalid acb_file_id")
		return
	}
	data, filename, err := services.DownloadACBFile(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "ACB file not found")
			return
		}
		InternalError(c, err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/plain", data)
}

// ProcessBankResponse handles POST /group-pricing/claims/acb-files/:acb_file_id/reconcile
func ProcessBankResponse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("acb_file_id"))
	if err != nil {
		BadRequestMsg(c, "invalid acb_file_id")
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		BadRequestMsg(c, "a file must be uploaded in the 'file' field")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	summary, err := services.ProcessBankResponse(id, fileHeader, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, summary)
}

// GetReconciliationResults handles GET /group-pricing/claims/acb-files/:acb_file_id/reconciliation
func GetACBReconciliationResults(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("acb_file_id"))
	if err != nil {
		BadRequestMsg(c, "invalid acb_file_id")
		return
	}
	results, err := services.GetACBReconciliationResults(id)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, results)
}

// GetReconciliationSummary handles GET /group-pricing/claims/payment-schedules/:schedule_id/reconciliation-summary
func GetACBReconciliationSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}
	summary, err := services.GetACBReconciliationSummary(id)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, summary)
}

// RetryFailedPayments handles POST /group-pricing/claims/acb-files/:acb_file_id/retry
func RetryFailedPayments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("acb_file_id"))
	if err != nil {
		BadRequestMsg(c, "invalid acb_file_id")
		return
	}
	var req models.RetryFailedRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	fileRecord, err := services.RetryFailedPayments(id, req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	Created(c, fileRecord)
}

// verifyBankAccountRequest is the frontend-facing request for bank account
// verification. ClaimID and Attempt are optional: when present the audit log
// can tie the call to a specific claim draft and retry counter; when absent
// the registry derives a stable idempotency key from the banking fields.
type verifyBankAccountRequest struct {
	FirstName         string `json:"first_name" binding:"required"`
	Surname           string `json:"surname"`
	IdentityNumber    string `json:"identity_number" binding:"required"`
	IdentityType      string `json:"identity_type"`
	BankAccountNumber string `json:"bank_account_number" binding:"required"`
	BankBranchCode    string `json:"bank_branch_code" binding:"required"`
	BankAccountType   string `json:"bank_account_type" binding:"required"`
	ClaimID           *int   `json:"claim_id"`
	Attempt           int    `json:"attempt"`
}

// toBavRequest maps the HTTP-layer DTO onto the canonical bav.VerifyRequest.
func (r verifyBankAccountRequest) toBavRequest() bav.VerifyRequest {
	return bav.VerifyRequest{
		FirstName:         r.FirstName,
		Surname:           r.Surname,
		IdentityNumber:    r.IdentityNumber,
		IdentityType:      r.IdentityType,
		BankAccountNumber: r.BankAccountNumber,
		BankBranchCode:    r.BankBranchCode,
		BankAccountType:   r.BankAccountType,
		ClaimID:           r.ClaimID,
		Attempt:           r.Attempt,
	}
}

// legacyVerifyResponse preserves the v1 wire shape the frontend consumes.
// Phase 4 replaces this with the canonical bav.VerifyResult shape under /v2.
type legacyVerifyResponse struct {
	Success   bool                        `json:"success"`
	RequestID string                      `json:"requestId"`
	Service   string                      `json:"service"`
	Results   legacyVerifyResponseResults `json:"results"`
}

type legacyVerifyResponseResults struct {
	IdentityAndAccountVerified bool                           `json:"identity_and_account_verified"`
	Summary                    string                         `json:"summary"`
	VerificationResults        legacyVerifyVerificationResult `json:"verification_results"`
}

type legacyVerifyVerificationResult struct {
	Status           string `json:"Status"`
	AccountFound     string `json:"accountFound"`
	AccountOpen      string `json:"accountOpen"`
	IdentityMatch    string `json:"identityMatch"`
	AccountTypeMatch string `json:"accountTypeMatch"`
	AcceptsCredits   string `json:"acceptsCredits"`
	AcceptsDebits    string `json:"acceptsDebits"`
}

// triToLegacy maps the canonical TriState onto the "Yes"/"No"/"Unknown"
// strings the v1 wire shape uses. The legacy frontend at
// ClaimRegistrationForm.vue branches on literal "Yes" so these spellings
// must be preserved until the v2 cutover completes.
func triToLegacy(t bav.TriState) string {
	switch t {
	case bav.TriYes:
		return "Yes"
	case bav.TriNo:
		return "No"
	default:
		return "Unknown"
	}
}

// VerifyBankAccount handles POST /group-pricing/claims/verify-bank-account
func VerifyBankAccount(c *gin.Context) {
	var req verifyBankAccountRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}

	result, err := bav.Verify(c.Request.Context(), req.toBavRequest())
	if err != nil {
		InternalError(c, err)
		return
	}

	OK(c, legacyVerifyResponse{
		Success:   result.Status == bav.StatusComplete,
		RequestID: result.ProviderRequestID,
		Service:   "bank-account-verification",
		Results: legacyVerifyResponseResults{
			IdentityAndAccountVerified: result.Verified,
			Summary:                    result.Summary,
			VerificationResults: legacyVerifyVerificationResult{
				Status:           result.ProviderStatusText,
				AccountFound:     triToLegacy(result.AccountFound),
				AccountOpen:      triToLegacy(result.AccountOpen),
				IdentityMatch:    triToLegacy(result.IdentityMatch),
				AccountTypeMatch: triToLegacy(result.AccountTypeMatch),
				AcceptsCredits:   triToLegacy(result.AcceptsCredits),
				AcceptsDebits:    triToLegacy(result.AcceptsDebits),
			},
		},
	})
}

// VerifyBankAccountV2Status handles
// POST /v2/group-pricing/claims/verify-bank-account/status/:job_id and resolves
// the current state of a pending async verification. Sync providers (VerifyNow)
// return 501 via ErrNotSupported.
func VerifyBankAccountV2Status(c *gin.Context) {
	jobID := c.Param("job_id")
	if jobID == "" {
		BadRequestMsg(c, "job_id is required")
		return
	}
	result, err := bav.Poll(c.Request.Context(), jobID)
	if err != nil {
		if errors.Is(err, bav.ErrNotSupported) {
			c.JSON(http.StatusNotImplemented, models.PremiumResponse{
				Success: false,
				Message: "active BAV provider does not support async polling",
			})
			return
		}
		if errors.Is(err, bav.ErrInvalidInput) {
			NotFound(c, "unknown job_id")
			return
		}
		InternalError(c, err)
		return
	}
	OK(c, result)
}

// BAVWebhook handles POST /bav/webhook/:provider. Scaffolded for Phase 6 but
// inert — HMAC verification and per-provider dispatch land in Phase 7 when
// a real async provider is wired in. Responding 501 makes accidental traffic
// loud rather than silently accepted.
func BAVWebhook(c *gin.Context) {
	provider := c.Param("provider")
	c.JSON(http.StatusNotImplemented, models.PremiumResponse{
		Success: false,
		Message: "BAV webhook is not configured for provider: " + provider,
	})
}

// VerifyBankAccountV2 handles POST /v2/group-pricing/claims/verify-bank-account
// and emits the canonical bav.VerifyResult shape. The legacy v1 endpoint will
// be removed one release after v2 ships.
func VerifyBankAccountV2(c *gin.Context) {
	var req verifyBankAccountRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}

	result, err := bav.Verify(c.Request.Context(), req.toBavRequest())
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, result)
}

// ──────────────────────────────────────────────
// Payment Schedule Lifecycle (Phase 1)
// ──────────────────────────────────────────────

// parseScheduleID is a small helper used by every lifecycle endpoint.
func parseScheduleID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return 0, false
	}
	return id, true
}

func parseItemID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		BadRequestMsg(c, "invalid item_id")
		return 0, false
	}
	return id, true
}

// SignOffPaymentSchedule handles POST /claims/payment-schedules/:schedule_id/signoff
func SignOffPaymentSchedule(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.SignOffByHeadOfClaims(id, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, schedule)
}

// StartFinanceReview handles POST /claims/payment-schedules/:schedule_id/finance/start-review
func StartFinanceReview(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.FinanceStartReview(id, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, schedule)
}

// VerifyScheduleLineItem handles POST /claims/payment-schedules/:schedule_id/items/:item_id/verify
func VerifyScheduleLineItem(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	item, err := services.VerifyLineItem(sid, iid, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, item)
}

// QueryScheduleLineItem handles POST /claims/payment-schedules/:schedule_id/items/:item_id/query
// Body: { "reason_code": "...", "notes": "..." }
func QueryScheduleLineItem(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	var req services.QueryRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.QueryLineItem(sid, iid, req, user); err != nil {
		BadRequest(c, err)
		return
	}
	schedule, err := services.GetPaymentSchedule(sid)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, schedule)
}

// RejectScheduleLineItem handles POST /claims/payment-schedules/:schedule_id/items/:item_id/reject
// Body: { "reason_code": "...", "notes": "..." }
func RejectScheduleLineItem(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	var req services.QueryRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.RejectLineItem(sid, iid, req, user); err != nil {
		BadRequest(c, err)
		return
	}
	schedule, err := services.GetPaymentSchedule(sid)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, schedule)
}

// AcknowledgeFinanceRejection handles
// POST /group-pricing/claims/:claim_id/finance-rejection/acknowledge
//
// Moves a finance-rejected claim back to under_assessment so the assessor
// can edit and re-approve. Refuses when the claim isn't currently in
// finance_rejected status — guards against double-clicks or stale UI.
func AcknowledgeFinanceRejection(c *gin.Context) {
	claimID, err := strconv.Atoi(c.Param("claim_id"))
	if err != nil {
		BadRequestMsg(c, "invalid claim_id")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	claim, err := services.AcknowledgeFinanceRejection(claimID, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "claim not found")
			return
		}
		BadRequest(c, err)
		return
	}
	OK(c, claim)
}

// FirstAuthorisePaymentSchedule handles POST /claims/payment-schedules/:schedule_id/finance/authorise-first
func FirstAuthorisePaymentSchedule(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.FinanceFirstAuthorise(id, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, schedule)
}

// SecondAuthorisePaymentSchedule handles POST /claims/payment-schedules/:schedule_id/finance/authorise-second
func SecondAuthorisePaymentSchedule(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.FinanceSecondAuthorise(id, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, schedule)
}

// ArchivePaymentSchedule handles POST /claims/payment-schedules/:schedule_id/archive
func ArchivePaymentSchedule(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	schedule, err := services.ArchiveSchedule(id, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, schedule)
}

// DiscardPaymentSchedule handles DELETE /claims/payment-schedules/:schedule_id.
// Only drafts can be discarded — once signed off, finance owns the schedule.
// All line-item claims are returned to "approved" so the next cut-off picks
// them up again.
func DiscardPaymentSchedule(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.DiscardDraftSchedule(id, user); err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, gin.H{"discarded": true, "schedule_id": id})
}

// PostScheduleFollowup handles POST /claims/payment-schedules/:schedule_id/followups
// Body: { "notes": "..." }
// Records a claims-side follow-up note on a submitted schedule. The note shows
// up in the existing Queries panel for finance to resolve.
func PostScheduleFollowup(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.RaiseClaimsFollowup(id, req.Notes, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, row)
}

// ResolveScheduleQuery handles POST /claims/payment-schedules/queries/:query_id/resolve
// Body: { "response": "..." }
// Finance uses this to respond to a claims follow-up or close out a line query.
func ResolveScheduleQuery(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("query_id"))
	if err != nil {
		BadRequestMsg(c, "invalid query_id")
		return
	}
	var req struct {
		Response string `json:"response"`
	}
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.ResolveScheduleQuery(id, req.Response, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, row)
}

// GetScheduleQueries handles GET /claims/payment-schedules/:schedule_id/queries
func GetScheduleQueries(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	rows, err := services.GetScheduleQueries(id)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// GetScheduleAuditTrail handles GET /claims/payment-schedules/:schedule_id/audit
func GetScheduleAuditTrail(c *gin.Context) {
	id, ok := parseScheduleID(c)
	if !ok {
		return
	}
	rows, err := services.GetScheduleAuditTrail(id)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// ──────────────────────────────────────────────
// Authority Matrix CRUD
// ──────────────────────────────────────────────

// ListAuthorityMatrix handles GET /claims/authority-matrix
func ListAuthorityMatrix(c *gin.Context) {
	rows, err := services.ListAuthorityMatrix()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// CreateAuthorityMatrixRow handles POST /claims/authority-matrix
func CreateAuthorityMatrixRow(c *gin.Context) {
	var row models.AuthorityMatrix
	if err := c.BindJSON(&row); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	created, err := services.CreateAuthorityMatrixRow(row, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	Created(c, created)
}

// UpdateAuthorityMatrixRow handles PATCH /claims/authority-matrix/:row_id
func UpdateAuthorityMatrixRow(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("row_id"))
	if err != nil {
		BadRequestMsg(c, "invalid row_id")
		return
	}
	var patch models.AuthorityMatrix
	if err := c.BindJSON(&patch); err != nil {
		BadRequest(c, err)
		return
	}
	row, err := services.UpdateAuthorityMatrixRow(id, patch)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, row)
}

// DeleteAuthorityMatrixRow handles DELETE /claims/authority-matrix/:row_id
func DeleteAuthorityMatrixRow(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("row_id"))
	if err != nil {
		BadRequestMsg(c, "invalid row_id")
		return
	}
	if err := services.DeleteAuthorityMatrixRow(id); err != nil {
		InternalError(c, err)
		return
	}
	OK(c, gin.H{"deleted": true})
}

// ──────────────────────────────────────────────
// Payment cut-off settings & manual run (Phase 2)
// ──────────────────────────────────────────────

// resolveLicense returns the license id this request operates on. Phase 2 is
// single-tenant so the empty-string license-id is the singleton install row.
// Reads X-License-Id when present for future-proofing.
func resolveLicense(c *gin.Context) string {
	return strings.TrimSpace(c.GetHeader("X-License-Id"))
}

// GetPaymentCutoffConfig handles GET /claims/payment-cutoff/config
func GetPaymentCutoffConfig(c *gin.Context) {
	cfg, err := services.GetPaymentCutoffConfig(resolveLicense(c))
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, cfg)
}

// SavePaymentCutoffConfig handles PUT /claims/payment-cutoff/config
func SavePaymentCutoffConfig(c *gin.Context) {
	var patch models.PaymentCutoffConfig
	if err := c.BindJSON(&patch); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	cfg, err := services.SavePaymentCutoffConfig(resolveLicense(c), patch, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, cfg)
}

// RunPaymentCutoffNow handles POST /claims/payment-cutoff/run — a manual
// "generate schedule for the latest cut-off" trigger. Records the run under
// trigger_type="manual".
func RunPaymentCutoffNow(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	run, err := services.RunCutoff(resolveLicense(c), time.Now(), "manual", user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, run)
}

// ListPaymentCutoffRuns handles GET /claims/payment-cutoff/runs
func ListPaymentCutoffRuns(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	rows, err := services.ListRecentCutoffRuns(resolveLicense(c), limit)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// GetNextPaymentCutoff handles GET /claims/payment-cutoff/next
func GetNextPaymentCutoff(c *gin.Context) {
	next, found := services.NextCutoff(resolveLicense(c), time.Now())
	if !found {
		OK(c, gin.H{"configured": false})
		return
	}
	OK(c, gin.H{"configured": true, "scheduled_at": next})
}

// ──────────────────────────────────────────────
// Sanctions / PEP screening (Phase 3)
// ──────────────────────────────────────────────

// ScreenScheduleLineItem handles POST /claims/payment-schedules/:schedule_id/items/:item_id/screen
func ScreenScheduleLineItem(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.ScreenLineItem(sid, iid, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, row)
}

// RecordSanctionsOutcome handles POST /claims/payment-schedules/:schedule_id/items/:item_id/sanctions-outcome
// Body: { "status": "clear" | "hit" | "manual_clear", "notes": "..." }
func RecordSanctionsOutcome(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	var req services.RecordSanctionsOutcomeRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.RecordSanctionsOutcome(sid, iid, req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, row)
}

// ListScheduleSanctionsScreenings handles GET /claims/payment-schedules/:schedule_id/sanctions
func ListScheduleSanctionsScreenings(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	rows, err := services.ListSanctionsScreeningsForSchedule(sid)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// ──────────────────────────────────────────────
// Reinsurance recovery (Phase 3)
// ──────────────────────────────────────────────

// SetReinsuranceRecovery handles PUT /claims/payment-schedules/:schedule_id/items/:item_id/reinsurance
// Body: { "required": true|false, "amount": 0 }
func SetReinsuranceRecovery(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	var req services.SetReinsuranceRecoveryRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.SetReinsuranceRecovery(sid, iid, req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, row)
}

// ConfirmReinsuranceRecoveryRaised handles POST /claims/payment-schedules/:schedule_id/items/:item_id/reinsurance/raised
func ConfirmReinsuranceRecoveryRaised(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.ConfirmReinsuranceRecoveryRaised(sid, iid, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, row)
}

// ──────────────────────────────────────────────
// Duplicate beneficiary clear (Phase 3)
// ──────────────────────────────────────────────

// ClearDuplicateBeneficiary handles POST /claims/payment-schedules/:schedule_id/items/:item_id/duplicate/clear
func ClearDuplicateBeneficiary(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	iid, ok := parseItemID(c)
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.ClearDuplicateBeneficiary(sid, iid, user); err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, gin.H{"cleared": true})
}

// ──────────────────────────────────────────────
// Payment exceptions + tax certificates (Phase 4)
// ──────────────────────────────────────────────

// ListPaymentExceptions handles GET /claims/payment-exceptions
// Query params: status (failed|unmatched|""), include_resolved=1
func ListPaymentExceptions(c *gin.Context) {
	req := services.ListPaymentExceptionsRequest{
		Status:          c.Query("status"),
		IncludeResolved: c.Query("include_resolved") == "1" || c.Query("include_resolved") == "true",
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		req.Limit = limit
	}
	rows, err := services.ListPaymentExceptions(req)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// GetPaymentExceptionsSummary handles GET /claims/payment-exceptions/summary
func GetPaymentExceptionsSummary(c *gin.Context) {
	summary, err := services.GetPaymentExceptionsSummary()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, summary)
}

// ExportPaymentExceptions handles GET /claims/payment-exceptions/export
// Honours the same status / include_resolved filters as ListPaymentExceptions
// so the CSV matches the rows currently shown on screen.
func ExportPaymentExceptions(c *gin.Context) {
	req := services.ListPaymentExceptionsRequest{
		Status:          c.Query("status"),
		IncludeResolved: c.Query("include_resolved") == "1" || c.Query("include_resolved") == "true",
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		req.Limit = limit
	}

	data, filename, err := services.ExportPaymentExceptionsCSV(req)
	if err != nil {
		InternalError(c, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", data)
}

// ListScheduleTaxCertificates handles GET /claims/payment-schedules/:schedule_id/tax-certificates
func ListScheduleTaxCertificates(c *gin.Context) {
	sid, ok := parseScheduleID(c)
	if !ok {
		return
	}
	rows, err := services.ListTaxCertificatesForSchedule(sid)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// DownloadTaxCertificate handles GET /claims/tax-certificates/:cert_id/download
func DownloadTaxCertificate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("cert_id"))
	if err != nil {
		BadRequestMsg(c, "invalid cert_id")
		return
	}
	data, contentType, filename, err := services.DownloadTaxCertificate(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "tax certificate not found")
			return
		}
		InternalError(c, err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

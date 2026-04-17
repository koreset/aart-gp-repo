package controllers

import (
	"api/models"
	"api/services"
	"api/services/bav"
	"errors"
	"net/http"
	"strconv"

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

// GetPaymentSchedules handles GET /group-pricing/claims/payment-schedules
func GetPaymentSchedules(c *gin.Context) {
	schedules, err := services.GetPaymentSchedules()
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

package controllers

import (
	"api/models"
	"api/services"
	"api/utils"
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

// verifyBankAccountRequest is the frontend-facing request for bank account verification.
type verifyBankAccountRequest struct {
	FirstName         string `json:"first_name" binding:"required"`
	Surname           string `json:"surname"`
	IdentityNumber    string `json:"identity_number" binding:"required"`
	BankAccountNumber string `json:"bank_account_number" binding:"required"`
	BankBranchCode    string `json:"bank_branch_code" binding:"required"`
	BankAccountType   string `json:"bank_account_type" binding:"required"`
}

// VerifyBankAccount handles POST /group-pricing/claims/verify-bank-account
func VerifyBankAccount(c *gin.Context) {
	var req verifyBankAccountRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}

	result, err := utils.VerifyBankAccount(utils.VerifyBankAccountRequest{
		FirstName:         req.FirstName,
		Surname:           req.Surname,
		IdentityNumber:    req.IdentityNumber,
		BankAccountNumber: req.BankAccountNumber,
		BankBranchCode:    req.BankBranchCode,
		BankAccountType:   req.BankAccountType,
	})
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, result)
}

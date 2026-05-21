package controllers

import (
	"api/models"
	"api/services"

	"github.com/gin-gonic/gin"
)

// Bank account sub-ledger and statement reconciliation. Mounted under
// /bank-accounts; depends on the same auth middleware as the /gl group.
//
// Bank account master data follows the maker/checker pattern (request →
// approve by different user). Statement imports themselves are raw data
// ingest and do not require approval — the maker/checker step applies at
// match/ignore review time.

func ListBankAccounts(c *gin.Context) {
	rows, err := services.GLListBankAccounts()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func RequestCreateBankAccount(c *gin.Context) {
	var req models.BankAccount
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestCreateBankAccount(req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	Created(c, row)
}

func RequestUpdateBankAccount(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.BankAccount
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestUpdateBankAccount(id, req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func RequestDeactivateBankAccount(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestDeactivateBankAccount(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func ApproveBankAccountChange(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.PendingChangeApprovalRequest
	_ = c.BindJSON(&req)
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLApproveBankAccountChange(id, req.Notes, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func ImportBankStatement(c *gin.Context) {
	var req models.BankStatementImportRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	batchID, count, err := services.GLImportBankStatement(req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	Created(c, gin.H{"batch_id": batchID, "imported": count})
}

func ListStatementLines(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	status := c.Query("status")
	rows, err := services.GLListStatementLines(id, status)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func MatchStatementLine(c *gin.Context) {
	var req models.BankMatchRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.GLMatchStatementLine(req, user); err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, gin.H{"status": "matched", "review_status": "pending_review"})
}

func IgnoreStatementLine(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.GLIgnoreStatementLine(id, user); err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, gin.H{"status": "ignored", "review_status": "pending_review"})
}

func ReviewStatementLine(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.StatementLineReviewRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLReviewStatementLine(id, req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

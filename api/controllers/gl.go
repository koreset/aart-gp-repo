package controllers

import (
	"api/models"
	"api/services"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HTTP surface for the operational General Ledger. All routes are mounted
// under /gl in routes.go and assume the GetActiveUser() middleware has run.
//
// Master-data mutations (chart of accounts, posting rules, periods, bank
// accounts) and manual journals all follow a maker/checker pattern — the
// request endpoint stages the change, and a different user approves it
// via a separate endpoint. Self-approval is rejected at the service layer
// with services.ErrSelfApproval (returned to the client as HTTP 409).

// respondServiceError maps service-layer typed errors to appropriate HTTP
// status codes. ErrSelfApproval is a 409 (conflict — the request is valid
// but the actor is not allowed to perform it themselves).
func respondServiceError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrSelfApproval) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, services.ErrBadState) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, services.ErrEntryNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	BadRequest(c, err)
}

// ---------------------------------------------------------------------------
// Chart of accounts (maker/checker)
// ---------------------------------------------------------------------------

func ListGLAccounts(c *gin.Context) {
	rows, err := services.GLListAccounts()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func GetGLAccount(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	row, err := services.GLGetAccount(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "GL account not found")
			return
		}
		InternalError(c, err)
		return
	}
	OK(c, row)
}

func RequestCreateGLAccount(c *gin.Context) {
	var req models.GLAccount
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestCreateAccount(req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	Created(c, row)
}

func RequestUpdateGLAccount(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.GLAccount
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestUpdateAccount(id, req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func RequestDeactivateGLAccount(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestDeactivateAccount(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func ApproveGLAccountChange(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.PendingChangeApprovalRequest
	_ = c.BindJSON(&req) // body optional
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLApproveAccountChange(id, req.Notes, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

// ---------------------------------------------------------------------------
// Accounting periods (two-step close)
// ---------------------------------------------------------------------------

func ListAccountingPeriods(c *gin.Context) {
	rows, err := services.GLListPeriods()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func CreateAccountingPeriod(c *gin.Context) {
	var req models.AccountingPeriod
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLCreatePeriod(req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	Created(c, row)
}

func RequestClosePeriod(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestClosePeriod(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func ClosePeriod(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLClosePeriod(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

// ---------------------------------------------------------------------------
// Posting rules (maker/checker)
// ---------------------------------------------------------------------------

func ListPostingRules(c *gin.Context) {
	rows, err := services.GLListPostingRules()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func RequestCreatePostingRule(c *gin.Context) {
	var req models.PostingRule
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestCreatePostingRule(req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	Created(c, row)
}

func RequestUpdatePostingRule(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.PostingRule
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestUpdatePostingRule(id, req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func RequestDeletePostingRule(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestDeletePostingRule(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func ApprovePostingRuleChange(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.PendingChangeApprovalRequest
	_ = c.BindJSON(&req)
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLApprovePostingRuleChange(id, req.Notes, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

// ---------------------------------------------------------------------------
// Journal entries: list, detail, draft/submit/approve/post, request/approve reversal
// ---------------------------------------------------------------------------

func ListJournalEntries(c *gin.Context) {
	opts := services.GLListJournalsOptions{
		SourceType: c.Query("source_type"),
		Status:     c.Query("status"),
	}
	if v := c.Query("period_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.PeriodID = n
		}
	}
	if v := c.Query("account_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.AccountID = n
		}
	}
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts.From = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts.To = t.Add(24*time.Hour - time.Second)
		}
	}
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.Limit = n
		}
	}
	rows, err := services.GLListJournals(opts)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func GetJournalEntry(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	row, err := services.GLGetJournalEntry(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "journal entry not found")
			return
		}
		InternalError(c, err)
		return
	}
	OK(c, row)
}

func DraftManualJournal(c *gin.Context) {
	var req models.ManualJournalRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLDraftManualJournal(req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	Created(c, row)
}

func UpdateDraftJournal(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.ManualJournalRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLUpdateDraftJournal(id, req, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func SubmitManualJournal(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLSubmitManualJournal(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func ApproveManualJournal(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLApproveManualJournal(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func PostApprovedJournal(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLPostApprovedJournal(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func DiscardDraftJournal(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.GLDiscardDraft(id, user); err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, gin.H{"status": "discarded"})
}

func RequestReverseJournal(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req models.ReverseJournalRequest
	if err := c.BindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	row, err := services.GLRequestReversal(id, req.Reason, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	OK(c, row)
}

func ApproveReverseJournal(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	user := c.MustGet("user").(models.AppUser)
	rev, err := services.GLApproveReversal(id, user)
	if err != nil {
		respondServiceError(c, err)
		return
	}
	Created(c, rev)
}

// ---------------------------------------------------------------------------
// Reports
// ---------------------------------------------------------------------------

func GetTrialBalance(c *gin.Context) {
	periodID := 0
	if v := c.Query("period_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			periodID = n
		}
	}
	rows, err := services.GLGetTrialBalance(periodID)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func GetAccountLedger(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var from, to time.Time
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			to = t.Add(24*time.Hour - time.Second)
		}
	}
	rows, err := services.GLGetAccountLedger(id, from, to)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "GL account not found")
			return
		}
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// ---------------------------------------------------------------------------
// Audit log
// ---------------------------------------------------------------------------

// ListGLAuditLogUsers returns the set of user names the audit-log filter
// dropdown should offer (union of names that have appeared in the log and
// known AppUser names). Useful even before the log has any entries.
func ListGLAuditLogUsers(c *gin.Context) {
	rows, err := services.GLListAuditLogUsers()
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

func ListGLAuditLog(c *gin.Context) {
	opts := services.GLListAuditLogOptions{
		EventType:  c.Query("event_type"),
		ObjectType: c.Query("object_type"),
		ChangedBy:  c.Query("changed_by"),
	}
	if v := c.Query("object_id"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.ObjectID = n
		}
	}
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts.From = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			opts.To = t.Add(24*time.Hour - time.Second)
		}
	}
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			opts.Limit = n
		}
	}
	rows, err := services.GLListAuditLog(opts)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseIDParam(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil || id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

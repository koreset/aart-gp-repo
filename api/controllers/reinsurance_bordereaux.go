package controllers

import (
	"api/log"
	"api/models"
	"api/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GenerateRIMemberBordereaux handles POST /group-pricing/reinsurance/bordereaux/member
func GenerateRIMemberBordereaux(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.GenerateRIBordereauxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Type = "member_census"
	run, err := services.GenerateRIMemberBordereaux(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

// GenerateRIClaimsBordereaux handles POST /group-pricing/reinsurance/bordereaux/claims
func GenerateRIClaimsBordereaux(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.GenerateRIBordereauxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Type = "claims_run"
	run, err := services.GenerateRIClaimsBordereaux(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

// GetRIBordereauxRuns handles GET /group-pricing/reinsurance/bordereaux
func GetRIBordereauxRuns(c *gin.Context) {
	treatyIDStr := c.Query("treaty_id")
	runType := c.Query("type")
	status := c.Query("status")
	treatyID := 0
	if treatyIDStr != "" {
		treatyID, _ = strconv.Atoi(treatyIDStr)
	}
	runs, err := services.GetRIBordereauxRuns(treatyID, runType, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": runs})
}

// GetRIBordereauxRunByID handles GET /group-pricing/reinsurance/bordereaux/:run_id
func GetRIBordereauxRunByID(c *gin.Context) {
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}
	run, err := services.GetRIBordereauxRunByID(runID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

// DiffRIBordereauxRun handles GET /group-pricing/reinsurance/bordereaux/:run_id/diff
// Optional query param `against` selects the run to compare against; when
// omitted, the run's ParentRunID is used so "what changed in this amendment?"
// is the natural default.
func DiffRIBordereauxRun(c *gin.Context) {
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}
	against := c.Query("against")
	diff, err := services.DiffRIBordereauxRuns(runID, against)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": diff})
}

// GetRIBordereauxMemberRows handles GET /group-pricing/reinsurance/bordereaux/:run_id/members
func GetRIBordereauxMemberRows(c *gin.Context) {
	runID := c.Param("run_id")
	rows, err := services.GetRIBordereauxMemberRows(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

// GetRIBordereauxClaimsRows handles GET /group-pricing/reinsurance/bordereaux/:run_id/claims
func GetRIBordereauxClaimsRows(c *gin.Context) {
	runID := c.Param("run_id")
	rows, err := services.GetRIBordereauxClaimsRows(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

// SubmitRIBordereaux handles POST /group-pricing/reinsurance/bordereaux/submit
func SubmitRIBordereaux(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.SubmitRIBordereauxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	run, err := services.SubmitRIBordereaux(req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

// AcknowledgeRIBordereaux handles POST /group-pricing/reinsurance/bordereaux/:run_id/acknowledge
func AcknowledgeRIBordereaux(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}
	run, err := services.AcknowledgeRIBordereaux(runID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

// GetRIBordereauxStats handles GET /group-pricing/reinsurance/bordereaux/stats
func GetRIBordereauxStats(c *gin.Context) {
	treatyIDStr := c.Query("treaty_id")
	treatyID := 0
	if treatyIDStr != "" {
		treatyID, _ = strconv.Atoi(treatyIDStr)
	}
	stats, err := services.GetRIBordereauxStats(treatyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

// MonitorLargeClaims handles POST /group-pricing/reinsurance/bordereaux/large-claims/monitor
func MonitorLargeClaims(c *gin.Context) {
	var req models.MonitorLargeClaimsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	count, err := services.MonitorLargeClaims(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"created": count}})
}

// GetLargeClaimNotices handles GET /group-pricing/reinsurance/bordereaux/large-claims
func GetLargeClaimNotices(c *gin.Context) {
	treatyIDStr := c.Query("treaty_id")
	schemeIDStr := c.Query("scheme_id")
	status := c.Query("status")
	treatyID := 0
	schemeID := 0
	if treatyIDStr != "" {
		treatyID, _ = strconv.Atoi(treatyIDStr)
	}
	if schemeIDStr != "" {
		schemeID, _ = strconv.Atoi(schemeIDStr)
	}
	notices, err := services.GetLargeClaimNotices(treatyID, schemeID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": notices})
}

// UpdateLargeClaimNotice handles PATCH /group-pricing/reinsurance/bordereaux/large-claims/:id
func UpdateLargeClaimNotice(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.UpdateLargeClaimNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notice, err := services.UpdateLargeClaimNotice(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": notice})
}

// AcceptLargeClaimNotice handles POST /reinsurance/bordereaux/large-claims/:id/accept
func AcceptLargeClaimNotice(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req services.LargeClaimResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notice, err := services.AcceptLargeClaimNotice(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": notice})
}

// RejectLargeClaimNotice handles POST /reinsurance/bordereaux/large-claims/:id/reject
func RejectLargeClaimNotice(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req services.LargeClaimResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Reason) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required when rejecting a notice"})
		return
	}
	notice, err := services.RejectLargeClaimNotice(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": notice})
}

// QueryLargeClaimNotice handles POST /reinsurance/bordereaux/large-claims/:id/query
func QueryLargeClaimNotice(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req services.LargeClaimResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.QueryDetails) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query_details is required"})
		return
	}
	notice, err := services.QueryLargeClaimNotice(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": notice})
}

// GetRIBordereauxKPIs handles GET /group-pricing/reinsurance/bordereaux/kpis
func GetRIBordereauxKPIs(c *gin.Context) {
	treatyID := 0
	if s := c.Query("treaty_id"); s != "" {
		treatyID, _ = strconv.Atoi(s)
	}
	periodFrom := c.Query("period_from")
	periodTo := c.Query("period_to")
	kpis, err := services.GetRIBordereauxKPIs(treatyID, periodFrom, periodTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": kpis})
}

// AcknowledgeRIBordereauxReceipt handles POST /group-pricing/reinsurance/bordereaux/:run_id/acknowledge-receipt
func AcknowledgeRIBordereauxReceipt(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}
	var req models.AcknowledgeReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	run, err := services.AcknowledgeRIBordereauxReceipt(runID, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

// AmendRIBordereaux handles POST /group-pricing/reinsurance/bordereaux/:run_id/amend
func AmendRIBordereaux(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}
	var req models.AmendRIBordereauxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	run, err := services.AmendRIBordereaux(runID, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": run})
}

// ValidateRIBordereaux handles POST /group-pricing/reinsurance/bordereaux/:run_id/validate
func ValidateRIBordereaux(c *gin.Context) {
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	ctx := log.ContextWithUserInfo(c.Request.Context(), user.UserEmail, user.UserName)
	summary, err := services.ValidateRIBordereaux(ctx, runID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

// GetRIValidationResults handles GET /group-pricing/reinsurance/bordereaux/:run_id/validation-results
func GetRIValidationResults(c *gin.Context) {
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}
	summary, err := services.GetRIValidationResults(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

// GetLargeClaimStats handles GET /group-pricing/reinsurance/bordereaux/large-claims/stats
func GetLargeClaimStats(c *gin.Context) {
	treatyIDStr := c.Query("treaty_id")
	treatyID := 0
	if treatyIDStr != "" {
		treatyID, _ = strconv.Atoi(treatyIDStr)
	}
	stats, err := services.GetLargeClaimStats(treatyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

// GetCatastropheClaimsRows handles GET /group-pricing/reinsurance/bordereaux/cat-events
func GetCatastropheClaimsRows(c *gin.Context) {
	catEventCode := c.Query("cat_event_code")
	treatyIDStr := c.Query("treaty_id")
	periodFrom := c.Query("period_from")
	periodTo := c.Query("period_to")
	treatyID := 0
	if treatyIDStr != "" {
		treatyID, _ = strconv.Atoi(treatyIDStr)
	}
	rows, err := services.GetCatastropheClaimsRows(catEventCode, treatyID, periodFrom, periodTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

// GetRunOffTreaties handles GET /group-pricing/reinsurance/bordereaux/run-off-treaties
func GetRunOffTreaties(c *gin.Context) {
	status := c.Query("status")
	treaties, err := services.GetRunOffTreaties(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": treaties})
}

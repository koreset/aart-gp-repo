package controllers

import (
	"api/log"
	"api/models"
	"api/services"
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateBordereaux handles POST /group-pricing/bordereaux/generate
func GenerateBordereaux(c *gin.Context) {
	// request context with logging/meta
	requestID, _ := c.Get("requestID")
	var ctx context.Context = context.Background()
	if v, ok := requestID.(string); ok {
		ctx = context.WithValue(ctx, log.RequestIDKey, v)
	}
	user := c.MustGet("user").(models.AppUser)
	ctx = log.ContextWithUserInfo(ctx, user.UserEmail, user.UserName)
	logger := log.WithContext(ctx)

	var req services.GenerateBordereauxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delegate to service (supports member, premium, claim)
	meta, err := services.GenerateBordereaux(ctx, req)
	if err != nil {
		logger.WithError(err).Error("GenerateBordereaux failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta)
}

// GetBordereauxFields handles GET /group-pricing/bordereaux/fields/:type
func GetBordereauxFields(c *gin.Context) {
	bType := c.Param("type")
	fields, err := services.GetBordereauxFieldsByType(bType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fields)
}

// DownloadBordereaux serves a generated bordereaux file from data/reports.
// Access is limited to the record's creator, reviewer, or approver.
func DownloadBordereaux(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	reportDir := filepath.Join("data", "reports")
	filePath := filepath.Join(reportDir, fileName)

	// Security: Ensure path is within reportDir
	absReportDir, _ := filepath.Abs(reportDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absReportDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	if _, err := services.AuthorizeBordereauxDownload(fileName, user); err != nil {
		switch {
		case errors.Is(err, services.ErrBordereauxNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		case errors.Is(err, services.ErrBordereauxNotAuthorized):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.File(absFilePath)
}

// GetGeneratedBordereauxList handles GET /group-pricing/bordereaux/generated
func GetGeneratedBordereauxList(c *gin.Context) {
	list, err := services.GetAllGeneratedBordereaux()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetGeneratedBordereauxByID handles GET /group-pricing/bordereaux/generated/:id
func GetGeneratedBordereauxByID(c *gin.Context) {
	generatedID := c.Param("id")
	if generatedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generated_id is required"})
		return
	}

	record, err := services.GetGeneratedBordereauxByGeneratedID(generatedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// GetGeneratedBordereauxData handles GET /group-pricing/bordereaux/generated/:id/data
func GetGeneratedBordereauxData(c *gin.Context) {
	generatedID := c.Param("id")
	if generatedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generated_id is required"})
		return
	}

	data, err := services.GetGeneratedBordereauxData(generatedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// AddBordereauxTimelineEntry handles POST /group-pricing/bordereaux/generated/:id/timeline
func AddBordereauxTimelineEntry(c *gin.Context) {
	generatedID := c.Param("id")
	if generatedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generated_id is required"})
		return
	}

	var entry models.BordereauxTimeline
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AddBordereauxTimelineEntry(generatedID, entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "timeline entry added"})
}

// GetBordereauxConfigurations handles GET /group-pricing/bordereaux/configurations
func GetBordereauxConfigurations(c *gin.Context) {
	configs, err := services.GetBordereauxConfigurations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

// GetBordereauxConfiguration handles GET /group-pricing/bordereaux/configurations/:id
func GetBordereauxConfiguration(c *gin.Context) {
	configId := c.Param("id")
	config, err := services.GetBordereauxConfiguration(configId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}
	c.JSON(http.StatusOK, config)
}

// SaveBordereauxConfiguration handles POST /group-pricing/bordereaux/configurations
func SaveBordereauxConfiguration(c *gin.Context) {
	var config models.BordereauxConfiguration
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedConfig, err := services.SaveBordereauxConfiguration(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, savedConfig)
}

// UpdateBordereauxConfiguration handles PUT /group-pricing/bordereaux/configurations/:id
func UpdateBordereauxConfiguration(c *gin.Context) {
	configId := c.Param("id")
	var config models.BordereauxConfiguration
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedConfig, err := services.UpdateBordereauxConfiguration(configId, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedConfig)
}

// DeleteBordereauxConfiguration handles DELETE /group-pricing/bordereaux/configurations/:id
func DeleteBordereauxConfiguration(c *gin.Context) {
	configId := c.Param("id")
	if err := services.DeleteBordereauxConfiguration(configId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "configuration deleted"})
}

// UpdateConfigurationUsage handles PATCH /group-pricing/bordereaux/configurations/:id/usage
func UpdateConfigurationUsage(c *gin.Context) {
	configId := c.Param("id")
	if err := services.UpdateConfigurationUsage(configId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "usage updated"})
}

// SubmitBordereauxBatch handles POST /group-pricing/bordereaux/batch-submit
func SubmitBordereauxBatch(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	var ctx context.Context = context.Background()
	if v, ok := requestID.(string); ok {
		ctx = context.WithValue(ctx, log.RequestIDKey, v)
	}
	user := c.MustGet("user").(models.AppUser)
	ctx = log.ContextWithUserInfo(ctx, user.UserEmail, user.UserName)
	logger := log.WithContext(ctx)

	var req services.BordereauxBatchSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SubmitBordereauxBatch(ctx, req, user); err != nil {
		logger.WithError(err).Error("SubmitBordereauxBatch failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bordereaux batch submission processed successfully"})
}

// ImportSchemeConfirmations handles POST /group-pricing/bordereaux/confirmations/import
func ImportSchemeConfirmations(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	var ctx context.Context = context.Background()
	if v, ok := requestID.(string); ok {
		ctx = context.WithValue(ctx, log.RequestIDKey, v)
	}
	user := c.MustGet("user").(models.AppUser)
	ctx = log.ContextWithUserInfo(ctx, user.UserEmail, user.UserName)
	logger := log.WithContext(ctx)

	// Parse form data
	//schemeIDStr := c.PostForm("scheme_id")
	//schemeID := 0
	//if schemeIDStr != "" {
	//	schemeID, _ = strconv.Atoi(schemeIDStr)
	//}

	fileType := c.PostForm("file_type")
	autoProcessStr := c.PostForm("auto_process")
	autoProcess := autoProcessStr == "true"

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	files := form.File["confirmation_files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	result, err := services.ImportSchemeConfirmations(ctx, fileType, autoProcess, files, user.UserEmail)
	if err != nil {
		logger.WithError(err).Error("ImportSchemeConfirmations failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetReconciliationStats handles GET /group-pricing/bordereaux/reconciliation/stats
func GetReconciliationStats(c *gin.Context) {
	schemeIDStr := c.Query("scheme_id")
	schemeID, _ := strconv.Atoi(schemeIDStr)
	if schemeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheme_id is required"})
		return
	}

	stats, err := services.GetReconciliationStats(schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetBordereauxDashboardStats handles GET /group-pricing/bordereaux/dashboard/stats
func GetBordereauxDashboardStats(c *gin.Context) {
	stats, err := services.GetBordereauxDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetBordereauxComplianceReport handles GET /group-pricing/bordereaux/compliance-report
// Query params: from=YYYY-MM-DD, to=YYYY-MM-DD. Streams an xlsx workbook with
// Summary, Open Discrepancies, Overdue Deadlines, Escalations, and Large Claim
// Notices sheets.
func GetBordereauxComplianceReport(c *gin.Context) {
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
	file, filename, err := services.GenerateComplianceReport(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// GetBordereauxAnalytics handles GET /group-pricing/bordereaux/analytics
// Query params: period=last_7_days|last_30_days|last_quarter|last_year|ytd,
// from=YYYY-MM-DD, to=YYYY-MM-DD, scheme_id=<int>.
func GetBordereauxAnalytics(c *gin.Context) {
	filter := services.BordereauxAnalyticsFilters{
		Period: c.Query("period"),
	}
	if v := c.Query("scheme_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			filter.SchemeID = id
		}
	}
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filter.From = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			// extend 'to' to end-of-day so the range is inclusive
			filter.To = t.Add(24*time.Hour - time.Second)
		}
	}
	resp, err := services.GetBordereauxAnalytics(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetBordereauxConfirmations handles GET /group-pricing/bordereaux/confirmations
func GetBordereauxConfirmations(c *gin.Context) {

	confirmations, err := services.GetBordereauxConfirmations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, confirmations)
}

// GetBordereauxConfirmation handles GET /group-pricing/bordereaux/confirmations/:id
func GetBordereauxConfirmation(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid confirmation id"})
		return
	}

	confirmation, err := services.GetBordereauxConfirmation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, confirmation)
}

// GetReconciliationResults handles GET /group-pricing/bordereaux/confirmations/:id/results
func GetReconciliationResults(c *gin.Context) {
	confirmationIDStr := c.Param("id")
	confirmationID, _ := strconv.Atoi(confirmationIDStr)
	if confirmationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid confirmation id"})
		return
	}

	results, err := services.GetReconciliationResults(confirmationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// GetUnmatchedReconciliationResults handles GET /group-pricing/bordereaux/confirmations/:id/unmatched
func GetUnmatchedReconciliationResults(c *gin.Context) {
	confirmationIDStr := c.Param("id")
	confirmationID, _ := strconv.Atoi(confirmationIDStr)
	if confirmationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid confirmation id"})
		return
	}

	results, err := services.GetUnmatchedReconciliationResults(confirmationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

// DeleteBordereauxConfirmation handles DELETE /group-pricing/bordereaux/confirmations/:id
func DeleteBordereauxConfirmation(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid confirmation id"})
		return
	}

	if err := services.DeleteBordereauxConfirmation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bordereaux confirmation and associated records deleted successfully"})
}

// ReconcilePendingConfirmations handles POST /group-pricing/bordereaux/confirmations/reconcile-pending
func ReconcilePendingConfirmations(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	var ctx context.Context = context.Background()
	if v, ok := requestID.(string); ok {
		ctx = context.WithValue(ctx, log.RequestIDKey, v)
	}
	user := c.MustGet("user").(models.AppUser)
	ctx = log.ContextWithUserInfo(ctx, user.UserEmail, user.UserName)
	logger := log.WithContext(ctx)

	count, err := services.ReconcilePendingConfirmations(ctx)
	if err != nil {
		logger.WithError(err).Error("ReconcilePendingConfirmations failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"message":         "Pending confirmations processed",
		"processed_count": count,
	})
}

// ReviewGeneratedBordereaux handles POST /group-pricing/bordereaux/generated/:id/review
func ReviewGeneratedBordereaux(c *gin.Context) {
	generatedID := c.Param("id")
	if generatedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generated_id is required"})
		return
	}
	user := c.MustGet("user").(models.AppUser)

	var body struct {
		Notes string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&body)

	brd, err := services.ReviewGeneratedBordereaux(generatedID, body.Notes, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": brd})
}

// ApproveGeneratedBordereaux handles POST /group-pricing/bordereaux/generated/:id/approve
func ApproveGeneratedBordereaux(c *gin.Context) {
	generatedID := c.Param("id")
	if generatedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generated_id is required"})
		return
	}
	user := c.MustGet("user").(models.AppUser)

	var body struct {
		Notes string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&body)

	brd, err := services.ApproveGeneratedBordereaux(generatedID, body.Notes, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": brd})
}

// ReturnOutboundToDraft handles POST /group-pricing/bordereaux/generated/:id/return-to-draft
func ReturnOutboundToDraft(c *gin.Context) {
	generatedID := c.Param("id")
	if generatedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "generated_id is required"})
		return
	}
	user := c.MustGet("user").(models.AppUser)

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}

	brd, err := services.ReturnOutboundToDraft(generatedID, body.Reason, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": brd})
}

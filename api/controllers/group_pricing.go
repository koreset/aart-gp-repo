package controllers

import (
	"api/log"
	"api/models"
	"api/services"
	"api/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// buildAttachmentViewerURL returns a relative URL that can be used by a document viewer/iframe
// to access an attachment for inline viewing or download. We return a relative path to avoid
// assumptions about external hostnames or reverse proxies.
func buildAttachmentViewerURL(attachmentID int) string {
	return fmt.Sprintf("/group-pricing/claims/attachments/%d/download", attachmentID)
}

// populateAttachmentViewerURLs sets the computed viewer_url on each attachment in-place.
func populateAttachmentViewerURLs(atts []models.GroupSchemeClaimAttachment) {
	for i := range atts {
		atts[i].ViewerURL = buildAttachmentViewerURL(atts[i].ID)
	}
}

func GenerateGroupPricingQuote(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)
	logger.Info("Processing GenerateGroupPricingQuote request")

	var groupQuote models.GroupPricingQuote

	// Bind JSON body to GroupPricingQuote struct
	err := c.BindJSON(&groupQuote)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to bind JSON to GroupPricingQuote")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// the groupQuote.CommencementDate is  GMT time, so we need to convert it to local time
	sastLocation, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		panic(err)
	}

	// 2. Use the .In() method to convert the time
	//sastTime := 	groupQuote.CommencementDate.In(sastLocation)

	sastTime := groupQuote.CommencementDate.In(sastLocation)
	fmt.Println(sastTime)
	groupQuote.CommencementDate = sastTime

	user := c.MustGet("user").(models.AppUser)
	logger.WithFields(map[string]interface{}{
		"user_email": user.UserEmail,
		"user_name":  user.UserName,
	}).Debug("User retrieved from context")

	logger.Info("Calling GenerateGroupPricingQuote service")
	err = services.GenerateGroupPricingQuote(groupQuote, user)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to generate group pricing quote")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Group pricing quote generated successfully")
	c.JSON(http.StatusCreated, nil)
}

func CalculateGroupPricingQuote(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	quoteId := c.Param("id")
	basis := c.Param("basis")
	credibility := utils.StringToFloat(c.Param("credibility"))

	logger.WithFields(map[string]interface{}{
		"quote_id": quoteId,
		"basis":    basis,
	}).Info("Processing CalculateGroupPricingQuote request")

	user := c.MustGet("user").(models.AppUser)
	logger.WithFields(map[string]interface{}{
		"user_email": user.UserEmail,
		"user_name":  user.UserName,
	}).Debug("User retrieved from context")

	quoteIDInt, _ := strconv.Atoi(quoteId)

	// Enqueue the calculation job instead of running it synchronously.
	// The background worker will pick it up and report progress via WebSocket.
	jobID, err := services.EnqueueCalculationJob(quoteIDInt, basis, credibility, user)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to enqueue calculation job")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	logger.WithField("job_id", jobID).Info("Calculation job enqueued")
	c.JSON(http.StatusAccepted, gin.H{"success": true, "jobId": jobID})
}

// GetCalculationJobStatus returns the current status of a calculation job.
func GetCalculationJobStatus(c *gin.Context) {
	jobID, err := strconv.Atoi(c.Param("jobId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid job ID"})
		return
	}

	job, err := services.GetCalculationJob(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// GetCustomTirTableStatus checks whether a quote that uses the custom tiered income
// replacement feature has had its custom table uploaded.
func GetCustomTirTableStatus(c *gin.Context) {
	quoteId := c.Param("id")

	result, err := services.CheckCustomTirTableStatus(quoteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func UpdateGroupPricingQuote(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	var quote models.GroupPricingQuote
	err := c.BindJSON(&quote)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	quoteId := quote.ID
	status := quote.Status

	logger.WithFields(map[string]interface{}{
		"quote_id": quoteId,
		"status":   status,
	}).Info("Processing UpdateGroupPricingQuote request")

	user := c.MustGet("user").(models.AppUser)
	logger.WithFields(map[string]interface{}{
		"user_email": user.UserEmail,
		"user_name":  user.UserName,
	}).Debug("User retrieved from context")

	logger.Info("Calling UpdateGroupPricingQuote service")
	err = services.UpdateGroupPricingQuote(quote, user)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to update group pricing quote")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Group pricing quote updated successfully")
	c.JSON(http.StatusOK, nil)
}

func ApproveGroupPricingQuote(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	quoteId := c.Param("id")

	logger.WithField("quote_id", quoteId).Info("Processing ApproveGroupPricingQuote request")

	user := c.MustGet("user").(models.AppUser)
	logger.WithFields(map[string]interface{}{
		"user_email": user.UserEmail,
		"user_name":  user.UserName,
	}).Debug("User retrieved from context")

	logger.Info("Calling ApproveGroupPricingQuote service")
	err := services.ApproveGroupPricingQuote(quoteId, user)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to approve group pricing quote")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Group pricing quote approved successfully")
	c.JSON(http.StatusOK, nil)
}

func AcceptGroupPricingQuote(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	quoteId := c.Param("id")

	// Parse request body to get commencement_date and term
	var requestBody struct {
		CommencementDate string `json:"commencement_date"`
		Term             string `json:"term"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.WithField("error", err.Error()).Error("Failed to parse request body")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	logger.WithFields(map[string]interface{}{
		"quote_id":          quoteId,
		"commencement_date": requestBody.CommencementDate,
		"term":              requestBody.Term,
	}).Info("Processing AcceptGroupPricingQuote request")

	user := c.MustGet("user").(models.AppUser)
	logger.WithFields(map[string]interface{}{
		"user_email": user.UserEmail,
		"user_name":  user.UserName,
	}).Debug("User retrieved from context")

	logger.Info("Calling AcceptGroupPricingQuote service")
	err := services.AcceptGroupPricingQuote(quoteId, requestBody.CommencementDate, requestBody.Term, user)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to accept group pricing quote")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("Group pricing quote accepted successfully")
	c.JSON(http.StatusOK, nil)
}

// CreateOnRiskLetter creates (or re-issues) an On Risk letter record for a quote.
func CreateOnRiskLetter(c *gin.Context) {
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid quote id")
		return
	}

	user := c.MustGet("user").(models.AppUser)

	letter, err := services.CreateOnRiskLetter(quoteId, user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, letter)
}

// GetOnRiskLetterData returns all data needed to render the On Risk letter document.
func GetOnRiskLetterData(c *gin.Context) {
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid quote id")
		return
	}

	data, err := services.GetOnRiskLetterData(quoteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetGroupPricingQuotes(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	filter := c.Param("filter")

	logger.WithField("filter", filter).Info("Processing GetGroupPricingQuotes request")

	logger.Info("Calling GetGroupPricingQuotes service")
	quotes, err := services.GetGroupPricingQuotes(filter)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get group pricing quotes")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("quote_count", len(quotes)).Info("Successfully retrieved group pricing quotes")
	c.JSON(http.StatusOK, quotes)
}

func GetGroupPricingQuote(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	id := c.Param("id")

	logger.WithField("quote_id", id).Info("Processing GetGroupPricingQuote request")

	logger.Info("Calling GetGroupPricingQuote service")
	quote, err := services.GetGroupPricingQuote(id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get group pricing quote")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithFields(map[string]interface{}{
		"quote_id":    quote.ID,
		"scheme_name": quote.SchemeName,
	}).Info("Successfully retrieved group pricing quote")
	c.JSON(http.StatusOK, quote)
}

// GetGroupPricingQuoteBySchemeName returns a quote given a scheme name by
// resolving the scheme's associated quote name and fetching the quote.
func GetGroupPricingQuoteBySchemeName(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	schemeName := c.Param("schemeName")
	logger.WithField("scheme_name", schemeName).Info("Processing GetGroupPricingQuoteBySchemeName request")

	quote, err := services.GetGroupPricingQuoteBySchemeName(schemeName)
	if err != nil {
		// Distinguish not found from other errors when possible
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithField("scheme_name", schemeName).Warn("Scheme or associated quote not found")
			c.JSON(http.StatusNotFound, fmt.Sprintf("no quote found for scheme name: %s", schemeName))
			return
		}
		logger.WithField("error", err.Error()).Error("Failed to get group pricing quote by scheme name")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, quote)
}

// GetQuotesForScheme returns all quotes associated to a scheme id
func GetQuotesForScheme(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	schemeId := c.Param("id")
	logger.WithField("scheme_id", schemeId).Info("Processing GetQuotesForScheme request")

	logger.Info("Calling GetGroupPricingQuotesBySchemeID service")
	quotes, err := services.GetGroupPricingQuotesBySchemeID(schemeId)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get quotes for scheme")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("quote_count", len(quotes)).Info("Successfully retrieved quotes for scheme")
	c.JSON(http.StatusOK, quotes)
}

func DeleteGroupPricingQuote(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	id := c.Param("id")

	logger.WithField("quote_id", id).Info("Processing DeleteGroupPricingQuote request")

	logger.Info("Calling DeleteGroupPricingQuote service")
	user := c.MustGet("user").(models.AppUser)
	err := services.DeleteGroupPricingQuote(id, user)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to delete group pricing quote")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("quote_id", id).Info("Successfully deleted group pricing quote")
	c.JSON(http.StatusOK, nil)
}

func UploadGPRateTables(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	tableTypeValues := form.Value["table_type"]
	if len(tableTypeValues) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_type is required"})
		return
	}
	tableType := tableTypeValues[0]

	fileValues := form.File["file"]
	if len(fileValues) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	file := fileValues[0]

	riskRateCodeValues := form.Value["risk_rate_code"]
	if len(riskRateCodeValues) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "risk_rate_code is required"})
		return
	}
	riskRateCode := riskRateCodeValues[0]

	// Optional scheme_name for custom tiered income replacement uploads
	var schemeName string
	if schemeNameValues, ok := form.Value["scheme_name"]; ok && len(schemeNameValues) > 0 {
		schemeName = schemeNameValues[0]
	}
	if tableType == "Custom Tiered Income Replacement" && schemeName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheme_name is required for Custom Tiered Income Replacement uploads"})
		return
	}

	fmt.Println(file.Filename)

	err = services.SaveGPTables(file, tableType, riskRateCode, user, schemeName)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully uploaded " + tableType})
}

func UploadQuoteTables(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form: " + err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)

	tableTypeValues := form.Value["table_type"]
	if len(tableTypeValues) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_type is required"})
		return
	}
	tableType := tableTypeValues[0]

	quoteIdValues := form.Value["quote_id"]
	if len(quoteIdValues) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quote_id is required"})
		return
	}
	quoteId, err := strconv.Atoi(quoteIdValues[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote_id: " + err.Error()})
		return
	}

	fileValues := form.File["file"]
	if len(fileValues) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	file := fileValues[0]

	err, count := services.SaveQuoteTables(file, tableType, quoteId, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"details": strings.Split(err.Error(), "; "),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count, "message": "Successfully uploaded " + tableType})
}

func DeleteQuoteTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	quoteId, _ := strconv.Atoi(c.Param("quote_id"))

	err := services.DeleteQuoteTableData(tableType, quoteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func GetGPTableMetadata(c *gin.Context) {
	metadata, err := services.GetGPTableMetaData()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, metadata)
}

// GetQuoteMemberGenderSplit returns the male / female counts of the uploaded
// member data for a quote. The Data Management screen renders this under the
// Member Data row once the upload completes.
func GetQuoteMemberGenderSplit(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid quote id"})
		return
	}
	split, err := services.GetQuoteMemberGenderSplit(quoteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": split})
}

// GetGroupPricingAgeBands returns the configured standard age bands used by the
// extended-family funeral benefit UI to partition the funeral rate table.
func GetGroupPricingAgeBands(c *gin.Context) {
	bands, err := services.GetGroupPricingAgeBands(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": bands})
}

func GetGPTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	results := services.GetGPTableData(tableType)

	c.JSON(http.StatusOK, results)
}

func GetGPTableYears(c *gin.Context) {
	tableType := c.Param("table_type")
	years := services.GetGPTableYears(tableType)

	c.JSON(http.StatusOK, years)
}

func GetGPTableRiskCodes(c *gin.Context) {
	tableType := c.Param("table_type")
	riskCodes := services.GetGPTableRiskCodes(tableType)

	c.JSON(http.StatusOK, riskCodes)
}

func DeleteGPTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	riskCode := c.Param("risk_code")
	err := services.DeleteGPTableData(tableType, riskCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

// RebuildGPTableStats recomputes row counts for all tracked rate tables and
// writes them to gp_table_stats. Run once after the initial deploy, or to
// repair drift caused by out-of-band database changes.
func RebuildGPTableStats(c *gin.Context) {
	if err := services.RebuildGPTableStats(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ListTableConfigurations returns every per-table-type "is required" row.
// Powers the Required-column rendering in the Tables UI when callers want
// the raw configuration without the populated-status payload.
func ListTableConfigurations(c *gin.Context) {
	rows, err := services.GetTableConfigurations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

// UpdateTableConfiguration toggles IsRequired for one table_type and writes
// an audit row capturing the active user, old value, new value, and an
// optional free-text note from the request body.
func UpdateTableConfiguration(c *gin.Context) {
	tableType := c.Param("table_type")
	if tableType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "table_type is required"})
		return
	}

	var body struct {
		IsRequired bool   `json:"is_required"`
		Note       string `json:"note"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	cfg, err := services.SetTableRequired(tableType, body.IsRequired, user, body.Note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cfg})
}

// GetTableConfigurationAuditLog returns the audit history for one table_type
// in reverse-chronological order. Powers the Info dialog history view.
func GetTableConfigurationAuditLog(c *gin.Context) {
	tableType := c.Param("table_type")
	if tableType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "table_type is required"})
		return
	}
	rows, err := services.GetTableConfigurationAudit(tableType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

func CreateBroker(c *gin.Context) {
	var broker models.Broker
	var appUser models.AppUser

	err := c.BindJSON(&broker)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	appUser = c.MustGet("user").(models.AppUser)

	err = services.CreateBroker(broker, appUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, err.Error())
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func GetBrokers(c *gin.Context) {
	brokers, err := services.GetBrokers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, brokers)
}

func GetBroker(c *gin.Context) {
	id := c.Param("id")
	broker, err := services.GetBroker(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, broker)
}

func EditBroker(c *gin.Context) {
	id := c.Param("id")
	var input models.Broker
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := services.EditBroker(id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "beneficiary not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func DeleteBroker(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteBroker(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// BinderFee CRUD operations

func CreateBinderFee(c *gin.Context) {
	var fee models.BinderFee
	if err := c.BindJSON(&fee); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	appUser := c.MustGet("user").(models.AppUser)

	if err := services.CreateBinderFee(fee, appUser); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "A binder fee for this binderholder and risk rate code already exists."})
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func GetBinderFees(c *gin.Context) {
	fees, err := services.GetBinderFees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, fees)
}

func GetBinderFee(c *gin.Context) {
	id := c.Param("id")
	fee, err := services.GetBinderFee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, fee)
}

func EditBinderFee(c *gin.Context) {
	id := c.Param("id")
	var input models.BinderFee
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := services.EditBinderFee(id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "binder fee not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func DeleteBinderFee(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteBinderFee(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

// CommissionStructure CRUD operations

func CreateCommissionBand(c *gin.Context) {
	var band models.CommissionStructure
	if err := c.BindJSON(&band); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appUser := c.MustGet("user").(models.AppUser)

	if err := services.CreateCommissionBand(band, appUser); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "A band with this lower bound already exists for this channel."})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func GetCommissionBands(c *gin.Context) {
	channel := c.Query("channel")
	// Distinguish between "holder_name param absent" and "holder_name="
	// (the latter means filter to default rows only).
	holderName, holderProvided := c.GetQuery("holder_name")
	allHolders := c.Query("all") == "1"
	bands, err := services.GetCommissionBands(channel, holderName, holderProvided, allHolders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, bands)
}

func GetCommissionBand(c *gin.Context) {
	id := c.Param("id")
	band, err := services.GetCommissionBand(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, band)
}

func EditCommissionBand(c *gin.Context) {
	id := c.Param("id")
	var input models.CommissionStructure
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := services.EditCommissionBand(id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "commission band not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func DeleteCommissionBand(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteCommissionBand(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

// Reinsurer CRUD operations

func CreateReinsurer(c *gin.Context) {
	var reinsurer models.Reinsurer
	var appUser models.AppUser

	err := c.BindJSON(&reinsurer)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	appUser = c.MustGet("user").(models.AppUser)

	err = services.CreateReinsurer(reinsurer, appUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "Reinsurer with this code already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func GetReinsurers(c *gin.Context) {
	reinsurers, err := services.GetReinsurers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reinsurers)
}

func GetReinsurer(c *gin.Context) {
	id := c.Param("id")
	reinsurer, err := services.GetReinsurer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reinsurer not found"})
		return
	}

	c.JSON(http.StatusOK, reinsurer)
}

func EditReinsurer(c *gin.Context) {
	id := c.Param("id")
	var input models.Reinsurer
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := services.EditReinsurer(id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reinsurer not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func DeleteReinsurer(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteReinsurer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// DeactivateReinsurer handles POST /group-pricing/reinsurers/:id/deactivate
func DeactivateReinsurer(c *gin.Context) {
	id := c.Param("id")
	user := c.MustGet("user").(models.AppUser)

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := services.DeactivateReinsurer(id, req.Reason, user)
	if err != nil {
		// Check if it's a validation error (active treaties)
		if strings.Contains(err.Error(), "cannot deactivate") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reinsurer deactivated successfully"})
}

func CreateGroupScheme(c *gin.Context) {
	var groupScheme models.GroupScheme
	var appUser models.AppUser

	err := c.BindJSON(&groupScheme)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	appUser = c.MustGet("user").(models.AppUser)

	err = services.CreateGroupScheme(groupScheme, appUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func CheckGroupSchemeName(c *gin.Context) {
	name := c.Param("name")
	result := services.CheckGroupSchemeName(name)
	c.JSON(http.StatusOK, gin.H{"name_exists": result})
}

func GetAllGroupSchemes(c *gin.Context) {
	groupSchemes, err := services.GetAllGroupSchemes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, groupSchemes)
}

func GetGroupSchemesInForce(c *gin.Context) {
	groupSchemes, err := services.GetGroupSchemesInforce()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, groupSchemes)
}

func GetGroupScheme(c *gin.Context) {
	id := c.Param("id")
	groupScheme, err := services.GetGroupScheme(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, groupScheme)
}

func GetSchemeCategories(c *gin.Context) {
	quote_id := c.Param("id")
	schemeCategories, err := services.GetGroupSchemeCategories(quote_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schemeCategories)
}

func AddMemberToScheme(c *gin.Context) {
	var memberData models.GPricingMemberDataInForce
	// Read the raw body so we can normalize allowed date formats (YYYY-MM-DD and YYYY/MM/DD)
	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Try default unmarshal first
	if uErr := json.Unmarshal(raw, &memberData); uErr != nil {
		// Fallback: normalize YYYY/MM/DD -> YYYY-MM-DD and try again
		normalized := normalizeSlashDates(raw)
		if uErr2 := json.Unmarshal(normalized, &memberData); uErr2 != nil {
			c.JSON(http.StatusBadRequest, uErr2.Error())
			return
		}
		// Restore body for any downstream middleware/handlers (not strictly needed here)
		c.Request.Body = io.NopCloser(strings.NewReader(string(normalized)))
	} else {
		// Restore body with the original payload
		c.Request.Body = io.NopCloser(strings.NewReader(string(raw)))
	}

	fmt.Println(memberData)

	user := c.MustGet("user").(models.AppUser)

	memberData, err = services.AddMemberToScheme(memberData, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, memberData)

}

// normalizeSlashDates converts any date strings in the JSON payload matching
// the pattern YYYY/MM/DD to YYYY-MM-DD to support both accepted formats.
func normalizeSlashDates(data []byte) []byte {
	// Regex to match full date tokens like 2025/12/21 not part of larger strings
	re := regexp.MustCompile(`\b(\d{4})/(\d{2})/(\d{2})\b`)
	return re.ReplaceAll(data, []byte("$1-$2-$3"))
}

func GetSchemeMembers(c *gin.Context) {
	schemeId := c.Param("id")
	members, err := services.GetSchemeMembers(schemeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, members)
}

// GetMemberInForce returns a single GPricingMemberDataInForce by member_id (primary key)
func GetMemberInForce(c *gin.Context) {
	memberId := c.Param("member_id")
	member, err := services.GetMemberInForceByID(memberId)
	if err != nil {
		// Map not found to 404, others to 500
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, member)
}

// GetMemberInForceByIdNumber returns a single GPricingMemberDataInForce by the member's IdNumber
func GetMemberInForceByIdNumber(c *gin.Context) {
	idNumber := c.Param("id_number")
	if idNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_number is required"})
		return
	}

	member, err := services.GetMemberInForceByIdNumber(idNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, member)
}

func SearchSchemeMembers(c *gin.Context) {
	schemeId := c.Param("id")
	quoteId := c.Param("quote_id")
	query := c.Query("query")

	if query == "" {
		c.JSON(http.StatusBadRequest, "Query parameter 'query' is required")
		return
	}

	members, err := services.SearchSchemeMembers(schemeId, quoteId, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, members)
}

// GetMemberBenefitSummaryInForce returns benefit summary for an in-force member
func GetMemberBenefitSummaryInForce(c *gin.Context) {
	memberIdStr := c.Param("member_id")
	if memberIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member_id is required"})
		return
	}
	memberID, err := strconv.Atoi(memberIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member_id must be an integer"})
		return
	}

	dto, err := services.GetMemberBenefitSummaryInForce(memberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto)
}

// GetMemberBenefitSummaryQuote returns benefit summary for a member within a quote context
func GetMemberBenefitSummaryQuote(c *gin.Context) {
	quoteIdStr := c.Param("id")
	memberIdStr := c.Param("member_id")
	if quoteIdStr == "" || memberIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quote id and member id are required"})
		return
	}
	quoteID, err := strconv.Atoi(quoteIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quote_id must be an integer"})
		return
	}
	memberID, err := strconv.Atoi(memberIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member_id must be an integer"})
		return
	}

	dto, err := services.GetMemberBenefitSummaryQuote(quoteID, memberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found in quote"})
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto)
}

func UpdateGroupSchemeStatus(c *gin.Context) {
	id := c.Param("id")
	var statusUpdate models.SchemeStatusUpdate

	var user models.AppUser
	user = c.MustGet("user").(models.AppUser)

	err := c.BindJSON(&statusUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	err = services.UpdateGroupSchemeStatus(id, statusUpdate, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, nil)
}

func DeleteGroupScheme(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteGroupScheme(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func DeactivateSchemeMember(c *gin.Context) {
	schemeId := c.Param("id")
	var member models.GPricingMemberDataInForce
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(models.AppUser)
	if err := services.DeactivateSchemeMember(schemeId, member, user); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

// UpdateMemberInForce updates a member's details and records an audit trail
func UpdateMemberInForce(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	//memberID, err := strconv.Atoi(memberIDStr)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
	//	return
	//}

	var input models.GPricingMemberDataInForce
	if inputErr := c.BindJSON(&input); inputErr != nil {
		fmt.Println(inputErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": inputErr.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	updated, err := services.UpdateMemberInForce(memberIDStr, input, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// GetMemberHistory returns generic audit history rows for a member
func GetMemberHistory(c *gin.Context) {
	memberID := c.Param("member_id")
	// Optional pagination
	limit := 0
	offset := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}

	// If 'structured' or 'activity' is requested, use the new MemberActivity table
	format := strings.ToLower(c.Query("format"))
	if format == "structured" || format == "activity" {
		rows, err := services.GetMemberActivityHistory(memberID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rows)
		return
	}

	include := strings.ToLower(strings.TrimSpace(c.Query("include")))
	if include == "associated" || include == "all" {
		mid, err := strconv.Atoi(memberID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
			return
		}
		rows, err := services.GetMemberFullAuditHistory(mid, true, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, rows)
		return
	}

	rows, err := services.GetAuditLogs("group-pricing", "g_pricing_member_data_in_forces", memberID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func GetGroupPricingParameterBases(c *gin.Context) {
	bases, err := services.GetGroupPricingParameterBases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, bases)
}

func GetGroupPricingIndustries(c *gin.Context) {
	bases, err := services.GetGroupPricingIndustries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, bases)
}

func GetGroupPricingQuoteTableData(c *gin.Context) {
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	tableType := c.Param("table_type")

	// Optional pagination via query parameters: ?offset=...&limit=...
	offset := 0
	limit := 0
	if offStr := c.Query("offset"); offStr != "" {
		if v, err := strconv.Atoi(offStr); err == nil && v >= 0 {
			offset = v
		}
	}
	if limStr := c.Query("limit"); limStr != "" {
		if v, err := strconv.Atoi(limStr); err == nil && v >= 0 {
			limit = v
		}
	}

	resultData, err := services.GetGroupPricingQuoteTableData(quoteId, tableType, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resultData)
}

func GetInForceTableData(c *gin.Context) {
	tableType := c.Param("table_type")
	schemeId := c.Param("id")
	results, err := services.GetInForceTableData(schemeId, tableType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, results)
}

func GetGroupPricingQuoteResultSummary(c *gin.Context) {
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	results, err := services.GetGroupPricingQuoteResultSummary(quoteId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, results)
}

func GetGroupPricingQuoteEducatorBenefits(c *gin.Context) {
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	results, err := services.GetGroupPricingQuoteEducatorBenefits(quoteId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, results)
}

func SaveInsurerDetails(c *gin.Context) {
	// expecting form data
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(models.AppUser)

	err = services.SaveInsurerDetails(form, user)

	c.JSON(http.StatusCreated, nil)
}

func GetInsurerDetails(c *gin.Context) {
	insurer, err := services.GetInsurerDetails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, insurer)
}

func GetGroupPricingDashboardData(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	dataSource := c.DefaultQuery("data_source", "inforce")
	benefit := c.DefaultQuery("benefit", "All")
	data, err := services.GetGroupPricingDashboardData(year, dataSource, benefit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetSchemePerformance returns per-scheme performance rows for the in-force
// Performance & Risk dashboard tab.
func GetSchemePerformance(c *gin.Context) {
	resp, err := services.GetSchemePerformanceRows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetInForceRiskProfile returns concentration KPIs, loss-ratio distribution,
// industry/region heatmap, frequency-severity scatter, and deteriorating
// schemes for the in-force Performance & Risk dashboard tab.
func GetInForceRiskProfile(c *gin.Context) {
	resp, err := services.GetInForceRiskProfile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetLossRatioTrend returns rolling-12-month ALR per scheme + portfolio over
// the trailing 24 months. Used by the Loss Ratio Trend chart.
func GetLossRatioTrend(c *gin.Context) {
	resp, err := services.GetLossRatioTrend()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetGroupSchemeClaimsDashboard returns aggregated analytics for claims dashboard
func GetGroupSchemeClaimsDashboard(c *gin.Context) {
	// Query params: scheme_id, period (ytd|last_12_months|last_3_months|custom), benefit, from, to
	var filters services.ClaimsDashboardFilters

	if sid := c.Query("scheme_id"); sid != "" {
		if id, err := strconv.Atoi(sid); err == nil {
			filters.SchemeID = &id
		} else {
			c.JSON(http.StatusBadRequest, "invalid scheme_id")
			return
		}
	}
	filters.Benefit = strings.TrimSpace(c.Query("benefit_type"))

	period := strings.ToLower(strings.TrimSpace(c.Query("period")))
	now := time.Now()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	// Helper to parse date YYYY-MM-DD
	parseDate := func(s string) (*time.Time, error) {
		if s == "" {
			return nil, nil
		}
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}

	switch period {
	case "ytd", "year_to_date", "year-to-date":
		filters.From = &startOfYear
		t := now
		filters.To = &t
	case "last_12_months", "last-12-months", "l12m":
		from := now.AddDate(-1, 0, 0)
		filters.From = &from
		t := now
		filters.To = &t
	case "last_3_months", "l3m":
		from := now.AddDate(0, -3, 0)
		filters.From = &from
		t := now
		filters.To = &t
	case "last_30_days", "last-30-days", "l30d", "30d":
		from := now.AddDate(0, 0, -30)
		filters.From = &from
		t := now
		filters.To = &t
	case "custom":
		from, err1 := parseDate(c.Query("from"))
		to, err2 := parseDate(c.Query("to"))
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, "invalid from/to date; expected YYYY-MM-DD")
			return
		}
		filters.From, filters.To = from, to
	case "":
		// default: last 12 months to match the UI screenshot
		from := now.AddDate(-1, 0, 0)
		filters.From = &from
		t := now
		filters.To = &t
	default:
		c.JSON(http.StatusBadRequest, "invalid period")
		return
	}

	// Optional limit for top claims table
	if limStr := c.Query("limit"); limStr != "" {
		if lim, err := strconv.Atoi(limStr); err == nil {
			filters.Limit = lim
		} else {
			c.JSON(http.StatusBadRequest, "invalid limit")
			return
		}
	}

	data, err := services.GetGroupSchemeClaimsDashboardData(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetGroupSchemeExposureData(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	benefit := c.Param("benefit")
	if benefit == "" {
		c.JSON(http.StatusBadRequest, "benefit is required")
		return
	}

	dataSource := c.DefaultQuery("data_source", "all")

	data, err := services.GetGroupSchemeExposureData(year, benefit, dataSource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func RebuildExposureData(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid year"})
		return
	}
	processed, err := services.RebuildExposureDataForYear(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"quotes_processed": processed}})
}

func GetExposureTrend(c *gin.Context) {
	benefit := c.DefaultQuery("benefit", "All")
	dataSource := c.DefaultQuery("data_source", "all")

	rows, err := services.GetExposureTimeSeries(benefit, dataSource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rows})
}

func GetFinancialYearInfo(c *gin.Context) {
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid year"})
		return
	}
	info := services.GetFinancialYearInfo(year)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": info})
}

func RegisterClaim(c *gin.Context) {
	var claim models.ClaimRegistration
	var documents []models.DocumentUpload

	// Parse multipart form
	err := c.Request.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get form values and bind to struct
	form, err := c.MultipartForm()

	// Basic form fields
	claim.MemberIDNumber = c.PostForm("member_id_number")
	claim.MemberName = c.PostForm("member_name")
	claim.SchemeName = c.PostForm("scheme_name")
	claim.MemberType = c.PostForm("member_type")
	claim.DateOfEvent = c.PostForm("date_of_event")
	claim.DateNotified = c.PostForm("date_notified")
	claim.Priority = c.PostForm("priority")
	claim.Status = c.PostForm("status")
	claim.DateRegistered = c.PostForm("date_registered")
	claim.ClaimNumber = c.PostForm("claim_number")
	claim.Description = c.PostForm("description")
	claim.BenefitName = c.PostForm("benefit_name")
	claim.BenefitCode = c.PostForm("benefit_code")
	claim.BenefitAlias = c.PostForm("benefit_alias")

	// Claimant information
	claim.ClaimantName = c.PostForm("claimant_name")
	claim.ClaimantIDNumber = c.PostForm("claimant_id_number")
	claim.RelationshipToMember = c.PostForm("relationship_to_member")
	claim.ClaimantContactNumber = c.PostForm("claimant_contact_number")

	// Parse numeric fields
	if schemeID := c.PostForm("scheme_id"); schemeID != "" {
		if id, err := strconv.Atoi(schemeID); err == nil {
			claim.SchemeID = id
		}
	}

	if claimAmount := c.PostForm("claim_amount"); claimAmount != "" {
		if amount, err := strconv.ParseFloat(claimAmount, 64); err == nil {
			claim.ClaimAmount = amount
		}
	}

	// Parse benefit type (if sent as JSON string)
	if benefitTypeStr := c.PostForm("benefit_type"); benefitTypeStr != "" {
		var benefitType models.BenefitType
		if err := json.Unmarshal([]byte(benefitTypeStr), &benefitType); err == nil {
			claim.BenefitType = benefitType
		}
	}

	// Parse missing required documents array
	if missingDocs := c.PostForm("missing_required_documents"); missingDocs != "" {
		var docs []string
		if err := json.Unmarshal([]byte(missingDocs), &docs); err == nil {
			claim.MissingRequiredDocuments = docs
		}
	}

	// Handle file uploads
	if form.File != nil {
		for fieldName, files := range form.File {
			// Skip non-file fields
			if fieldName == "supporting_documents" {
				continue
			}

			for _, file := range files {
				// Extract document type from field name or form data
				documentType := c.PostForm(fmt.Sprintf("%s_document_type", fieldName))
				documentName := c.PostForm(fmt.Sprintf("%s_document_name", fieldName))

				document := models.DocumentUpload{
					File:         file,
					DocumentType: documentType,
					DocumentName: documentName,
					FileName:     file.Filename,
					FileSize:     file.Size,
				}
				documents = append(documents, document)
			}
		}
	}

	// Validate required fields
	if claim.MemberIDNumber == "" || claim.ClaimAmount <= 0 || claim.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Save files to storage
	savedDocuments := make([]models.DocumentUpload, 0, len(documents))
	for _, doc := range documents {
		// Generate unique filename
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%d_%s_%s",
			timestamp,
			doc.DocumentType,
			filepath.Base(doc.File.Filename))

		// Save file
		savePath := filepath.Join("uploads/claims", claim.ClaimNumber, filename)
		if err := c.SaveUploadedFile(doc.File, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Update document with saved path
		doc.FileName = filename
		savedDocuments = append(savedDocuments, doc)
	}

	// Save claim to database
	if err := saveClaimToDatabase(claim, savedDocuments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save claim"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message":            "Claim registered successfully",
		"claim_number":       claim.ClaimNumber,
		"documents_uploaded": len(savedDocuments),
		"missing_documents":  len(claim.MissingRequiredDocuments),
	})
}

// Helper function to save claim to database
func saveClaimToDatabase(claim models.ClaimRegistration, documents []models.DocumentUpload) error {
	// Implement your database save logic here
	// This could use GORM, raw SQL, or your preferred ORM
	return nil
}

func GroupSchemeSubmitClaim(c *gin.Context) {
	// Parse multipart form (32 MB max)
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse multipart form: "+err.Error())
		return
	}

	form, _ := c.MultipartForm()
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, "User context missing")
		return
	}
	appUser := user.(models.AppUser)

	// If it's a multipart request with claim_data, handle it as a single claim with files
	if claimDataJSON := c.PostForm("claim_data"); claimDataJSON != "" {
		var claim models.GroupSchemeClaim
		if err := json.Unmarshal([]byte(claimDataJSON), &claim); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid claim_data JSON: "+err.Error())
			return
		}

		// Also handle document_metadata if provided
		if documentMetadataJSON := c.PostForm("document_metadata"); documentMetadataJSON != "" {
			var metadata []models.SupportingDocument
			if err := json.Unmarshal([]byte(documentMetadataJSON), &metadata); err != nil {
				c.JSON(http.StatusBadRequest, "Invalid document_metadata JSON: "+err.Error())
				return
			}
			claim.SupportingDocuments = metadata
		}

		// Use the existing service function that handles claim + files
		// Note: The service uses claim.SupportingDocuments to map files by index
		if err := services.GroupSchemeSubmitClaimWithFiles(claim, form.File, appUser); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusCreated, nil)
		return
	}

	// Fallback to JSON body if not multipart or claim_data not in form
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to read request body: "+err.Error())
		return
	}

	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, "Request body is empty")
		return
	}

	// Support both single object and array payloads for JSON
	if body[0] == '[' {
		var claims []models.GroupSchemeClaim
		if err := json.Unmarshal(body, &claims); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid JSON array: "+err.Error())
			return
		}
		created, err := services.GroupSchemeSubmitClaimsBatch(claims, appUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusCreated, created)
		return
	}

	var claim models.GroupSchemeClaim
	if err := json.Unmarshal(body, &claim); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	if err := services.GroupSchemeSubmitClaim(claim, appUser); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func GetGroupSchemeClaims(c *gin.Context) {
	claims, err := services.GetGroupSchemeClaims()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	// Populate viewer URLs for attachments in each claim
	for ci := range claims {
		populateAttachmentViewerURLs(claims[ci].Attachments)
	}
	c.JSON(http.StatusOK, claims)
}

func GetUpdatedClaimAmount(c *gin.Context) {
	var payload struct {
		MemberIDNumber string `json:"member_id_number"`
		SchemeID       int    `json:"scheme_id"`
		ClaimType      string `json:"member_type"`
		BenefitType    string `json:"benefit_type"`
		BenefitCode    string `json:"benefit_code"`
		BenefitAlias   string `json:"benefit_alias"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount := services.GetAcceleratedApprovedClaims(payload.MemberIDNumber, payload.SchemeID, payload.ClaimType, payload.BenefitCode)
	c.JSON(http.StatusOK, gin.H{"updated_claim_amount": amount})
}

// GetGroupSchemeClaim returns a single claim by its ID
func GetGroupSchemeClaim(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}

	claim, svcErr := services.GetGroupSchemeClaimByID(claimID)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	// Populate viewer URLs for attachments
	populateAttachmentViewerURLs(claim.Attachments)
	c.JSON(http.StatusOK, claim)
}

// ListGroupSchemeClaimAttachments returns all attachments metadata for a claim
func ListGroupSchemeClaimAttachments(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}

	atts, svcErr := services.GetClaimAttachments(claimID)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	populateAttachmentViewerURLs(atts)
	c.JSON(http.StatusOK, atts)
}

// DownloadGroupSchemeClaimAttachment streams an attachment file for inline viewing/downloading
func DownloadGroupSchemeClaimAttachment(c *gin.Context) {
	attIDStr := c.Param("attachment_id")
	attID, err := strconv.Atoi(attIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid attachment_id")
		return
	}

	att, svcErr := services.GetAttachmentByID(attID)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "attachment not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}

	// Security: ensure the file path is under our uploads base
	baseDirAbs, _ := filepath.Abs(filepath.Join("tmp", "uploads", "group_claims"))
	fileAbs, _ := filepath.Abs(att.StoragePath)
	if len(fileAbs) < len(baseDirAbs) || fileAbs[:len(baseDirAbs)] != baseDirAbs {
		c.JSON(http.StatusForbidden, "invalid file path")
		return
	}

	// Check file exists
	fi, statErr := os.Stat(fileAbs)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			c.JSON(http.StatusNotFound, "file not found on disk")
			return
		}
		c.JSON(http.StatusInternalServerError, statErr.Error())
		return
	}

	// Determine content type
	contentType := att.ContentType
	if contentType == "" {
		// Try to infer from extension
		if ext := filepath.Ext(att.FileName); ext != "" {
			if mt := mime.TypeByExtension(ext); mt != "" {
				contentType = mt
			}
		}
		if contentType == "" {
			// Fallback by sniffing
			if f, err := os.Open(fileAbs); err == nil {
				defer f.Close()
				buf := make([]byte, 512)
				n, _ := f.Read(buf)
				contentType = http.DetectContentType(buf[:n])
			}
		}
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	// Suggest inline for common viewable types
	disposition := "attachment"
	if contentType == "application/pdf" || len(contentType) >= 6 && contentType[:6] == "image/" {
		disposition = "inline"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fi.Size()))
	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, att.FileName))
	c.File(fileAbs)
}

// UploadGroupSchemeClaimAttachments handles uploading one or more files as supporting documents
// for a group scheme claim. Expects multipart/form-data. Files can be sent under any field names;
// all files in the multipart body will be processed.
func UploadGroupSchemeClaimAttachments(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}

	form, formErr := c.MultipartForm()
	if formErr != nil || form == nil {
		c.JSON(http.StatusBadRequest, "multipart/form-data with at least one file is required")
		return
	}

	user := c.MustGet("user").(models.AppUser)
	created, svcErr := services.AppendClaimAttachments(claimID, form.File, form.Value, user)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}

	populateAttachmentViewerURLs(created)
	c.JSON(http.StatusCreated, created)
}

// UpdateGroupSchemeClaim updates a claim by its ID
func UpdateGroupSchemeClaim(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}

	// Try to bind multipart/form-data similar to GroupSchemeSubmitClaim
	if form, formErr := c.MultipartForm(); formErr == nil && form != nil {
		var payload models.GroupSchemeClaim
		// Expect JSON blob under claim_data
		vals := form.Value["claim_data"]
		if len(vals) == 0 || vals[0] == "" {
			c.JSON(http.StatusBadRequest, "claim_data is required in multipart form")
			return
		}
		if err := json.Unmarshal([]byte(vals[0]), &payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		// Collect any uploaded files
		user := c.MustGet("user").(models.AppUser)
		updated, svcErr := services.UpdateGroupSchemeClaimWithFiles(claimID, payload, form.File, user)
		if svcErr != nil {
			if errors.Is(svcErr, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, "claim not found")
				return
			}
			c.JSON(http.StatusInternalServerError, svcErr.Error())
			return
		}
		// Ensure any returned attachments include viewer_url
		populateAttachmentViewerURLs(updated.Attachments)
		c.JSON(http.StatusOK, updated)
		return
	}

	// Fallback: accept JSON body for backward compatibility
	var payload models.GroupSchemeClaim
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(models.AppUser)
	updated, svcErr := services.UpdateGroupSchemeClaim(claimID, payload, user)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	populateAttachmentViewerURLs(updated.Attachments)
	c.JSON(http.StatusOK, updated)
}

// CreateGroupSchemeClaimAssessment handles creating a new assessment for a claim
func CreateGroupSchemeClaimAssessment(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}

	var payload models.GroupSchemeClaimAssessment
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// Ensure URL param prevails
	payload.ClaimID = claimID

	user := c.MustGet("user").(models.AppUser)
	assessment, svcErr := services.CreateClaimAssessment(payload, user)

	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	c.JSON(http.StatusCreated, assessment)
}

// UpdateGroupSchemeClaimAssessment updates an assessment by ID
func UpdateGroupSchemeClaimAssessment(c *gin.Context) {
	assessmentIDStr := c.Param("assessment_id")
	assessmentID, err := strconv.Atoi(assessmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid assessment_id")
		return
	}

	var payload models.GroupSchemeClaimAssessment
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(models.AppUser)
	updated, svcErr := services.UpdateClaimAssessment(assessmentID, payload, user)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "assessment not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	c.JSON(http.StatusOK, updated)
}

// GetGroupSchemeClaimAssessments fetches assessments for a claim
func GetGroupSchemeClaimAssessments(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}
	assessments, svcErr := services.GetClaimAssessmentsByClaim(claimID)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	c.JSON(http.StatusOK, assessments)
}

// CreateGroupSchemeClaimCommunication handles creating a new communication for a claim
func CreateGroupSchemeClaimCommunication(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}

	var payload models.GroupSchemeClaimCommunication
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// Ensure URL param prevails
	payload.ClaimID = claimID

	user := c.MustGet("user").(models.AppUser)
	comm, svcErr := services.CreateClaimCommunication(payload, user)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	c.JSON(http.StatusCreated, comm)
}

// CreateGroupSchemeClaimDecline handles creating a new decline record for a claim
func CreateGroupSchemeClaimDecline(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}

	var payload models.GroupSchemeClaimDecline
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// Ensure URL param prevails
	payload.ClaimID = claimID

	user := c.MustGet("user").(models.AppUser)
	rec, svcErr := services.CreateClaimDecline(payload, user)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	c.JSON(http.StatusCreated, rec)
}

// GetGroupSchemeClaimDeclines fetches decline records for a claim
func GetGroupSchemeClaimDeclines(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}
	declines, svcErr := services.GetClaimDeclinesByClaim(claimID)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	c.JSON(http.StatusOK, declines)
}

// GetGroupSchemeClaimCommunications fetches communications for a claim
func GetGroupSchemeClaimCommunications(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}
	comms, svcErr := services.GetClaimCommunicationsByClaim(claimID)
	if svcErr != nil {
		if errors.Is(svcErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, svcErr.Error())
		return
	}
	c.JSON(http.StatusOK, comms)
}

func GetSchemeMemberRating(c *gin.Context) {
	schemeId := c.Param("scheme_id")
	memberId := c.Param("member_id")
	quoteId := c.Param("quote_id")
	rating, err := services.GetSchemeMemberRating(schemeId, quoteId, memberId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rating)
}

func GetBenefitMaps(c *gin.Context) {
	maps, err := services.GetBenefitMaps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, maps)
}

func GetBenefitMapsByScheme(c *gin.Context) {
	schemeId := c.Param("scheme_id")
	maps, err := services.GetBenefitMapsByScheme(schemeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, maps)
}

func GetBenefitMapsBySchemeCategory(c *gin.Context) {
	schemeId := c.Param("scheme_id")
	categoryId := c.Param("category_id")
	maps, err := services.GetBenefitMapsBySchemeCategory(schemeId, categoryId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, maps)
}

func SaveBenefitMaps(c *gin.Context) {
	var maps []models.GroupBenefitMapper
	err := c.BindJSON(&maps)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.SaveBenefitMaps(maps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func GetBenefitDefinitions(c *gin.Context) {
	definitions, err := services.GetBenefitDefinitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, definitions)
}

func GetGPPermissions(c *gin.Context) {
	permissions, err := services.GetGPPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func GetGPUserRoles(c *gin.Context) {
	roles, err := services.GetGPUserRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func CreateGPUserRole(c *gin.Context) {
	var role models.GPUserRole
	err := c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	role, err = services.CreateGPUserRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, role)
}

func DeleteGPUserRole(c *gin.Context) {
	roleId := c.Param("role_id")
	err := services.DeleteGPUserRole(roleId)

	if err != nil && err.Error() == "role is in use" {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetRolePermissions(c *gin.Context) {
	roleId := c.Param("role_id")
	permissions, err := services.GetRolePermissions(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func AssignRoleToUser(c *gin.Context) {
	var userRole models.OrgUser
	err := c.BindJSON(&userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.AssignRoleToUser(userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func RemoveRoleFromUser(c *gin.Context) {
	var userRole models.OrgUser
	err := c.BindJSON(&userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.RemoveRoleFromUser(userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetRoleForUser(c *gin.Context) {
	licenseId := c.Param("license_id")
	role, err := services.GetRoleForUserLicense(licenseId)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, role)
}

func GetGroupPricingIndustriesForQuotes(c *gin.Context) {
	industries, err := services.GetGroupPricingIndustriesForQuotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, industries)
}

func GetBenefitEscalationOptions(c *gin.Context) {
	data, err := services.GetBenefitEscalationsOptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetTTDDisabilityDefinitions(c *gin.Context) {
	riskRateCode := c.Param("risk_rate_code")
	data, err := services.GetTTDDisabilityDefinitions(riskRateCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetPtdDisabilityDefinitions(c *gin.Context) {
	riskRateCode := c.Param("risk_rate_code")
	data, err := services.GetPTDDisabilityDefinitions(riskRateCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetEducatorBenefitTypes(c *gin.Context) {
	riskRateCode := c.Param("risk_rate_code")
	data, err := services.GetEducatorBenefitTypes(riskRateCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetPhiDisabilityDefinitions(c *gin.Context) {
	riskRateCode := c.Param("risk_rate_code")
	data, err := services.GetPhiDisabilityDefinitions(riskRateCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func GetDistinctWaitingPeriods(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	tableName := c.Param("table_type")

	logger.WithField("table_type", tableName).Info("Processing GetDistinctWaitingPeriods request")

	logger.Info("Calling GetDistinctWaitingPeriods service")
	riskRateCode := c.Param("risk_rate_code")
	waitingPeriods, err := services.GetDistinctWaitingPeriods(tableName, riskRateCode)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get distinct waiting periods")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("count", len(waitingPeriods)).Info("Successfully retrieved distinct waiting periods")
	c.JSON(http.StatusOK, waitingPeriods)
}

func GetDistinctGlaBenefitTypes(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	riskRateCode := c.Param("risk_rate_code")

	logger.WithField("risk_rate_code", riskRateCode).Info("Processing GetDistinctGlaBenefitTypes request")

	benefitTypes, err := services.GetDistinctGlaBenefitTypes(riskRateCode)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get distinct GLA benefit types")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("count", len(benefitTypes)).Info("Successfully retrieved distinct GLA benefit types")
	c.JSON(http.StatusOK, benefitTypes)
}

func GetDistinctDeferredPeriods(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	tableName := c.Param("table_type")

	logger.WithField("table_type", tableName).Info("Processing GetDistinctDeferredPeriods request")

	logger.Info("Calling GetDistinctDeferredPeriods service")
	riskRateCode := c.Param("risk_rate_code")
	deferredPeriods, err := services.GetDistinctDeferredPeriods(tableName, riskRateCode)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get distinct deferred periods")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("count", len(deferredPeriods)).Info("Successfully retrieved distinct deferred periods")
	c.JSON(http.StatusOK, deferredPeriods)
}

// GetDistinctNormalRetirementAges handles the request to retrieve distinct normal retirement ages from phi_rates
func GetDistinctNormalRetirementAges(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	logger.Info("Processing GetDistinctNormalRetirementAges request")

	logger.Info("Calling GetDistinctNormalRetirementAges service")
	normalRetirementAges, err := services.GetDistinctNormalRetirementAges()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get distinct normal retirement ages")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("count", len(normalRetirementAges)).Info("Successfully retrieved distinct normal retirement ages")
	c.JSON(http.StatusOK, normalRetirementAges)
}

// GetHistoricalCredibilityData handles the request to retrieve historical credibility data
func GetHistoricalCredibilityData(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	logger.Info("Processing GetHistoricalCredibilityData request")

	logger.Info("Calling GetHistoricalCredibilityData service")
	data, err := services.GetHistoricalCredibilityData()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get historical credibility data")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logger.WithField("count", len(data)).Info("Successfully retrieved historical credibility data")
	c.JSON(http.StatusOK, data)
}

// GetDistinctRiskTypes handles the request to retrieve distinct risk types from phi_rates and ttd_rates tables
func GetDistinctRiskTypes(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Get user info if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)

	logger.Info("Processing GetDistinctRiskTypes request")

	logger.Info("Calling GetDistinctRiskTypes service")
	riskTypes, err := services.GetDistinctRiskTypes()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get distinct risk types")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Log the count of risk types from each table
	logger.WithFields(map[string]interface{}{
		"phi_rates_count": len(riskTypes.PhiRates),
		"ttd_rates_count": len(riskTypes.TtdRates),
	}).Info("Successfully retrieved distinct risk types")

	c.JSON(http.StatusOK, riskTypes)
}

// UpdateGroupSchemeCoverEndDate handles updating the cover end date for a group scheme
func UpdateGroupSchemeCoverEndDate(c *gin.Context) {
	type request struct {
		SchemeID     int    `json:"scheme_id" binding:"required"`
		CoverEndDate string `json:"cover_end_date" binding:"required"`
	}
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	parsedDate, err := time.Parse("2006-01-02", req.CoverEndDate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}
	var user models.AppUser
	user = c.MustGet("user").(models.AppUser)
	err = services.UpdateGroupSchemeCoverEndDate(req.SchemeID, parsedDate, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cover end date updated successfully"})
}

// UpdateGroupPricingQuoteMemberStats handles updating member stats for a GroupPricingQuote
func UpdateGroupPricingQuoteMemberStats(c *gin.Context) {
	// Get request ID from context if available
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	// Bind JSON payload
	var input []models.MemberIndicativeDataSet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Add user info to context if available
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}
	logger := log.WithContext(ctx)

	// Validate quote id
	//if input.QuoteID == 0 {
	//	c.JSON(http.StatusBadRequest, "quote_id is required")
	//	return
	//}

	user := c.MustGet("user").(models.AppUser)
	//logger.WithFields(map[string]interface{}{
	//	"quote_id": input.QuoteID,
	//}).Info("Processing UpdateGroupPricingQuoteMemberStats request")

	updated, err := services.UpdateGroupPricingQuoteMemberStats(input, user)
	if err != nil {
		// Return 404 if not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "quote not found")
			return
		}
		logger.WithField("error", err.Error()).Error("Failed to update member stats")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, updated)
}

// UpdateGroupPricingQuoteIndicativeFlag toggles the member_indicative_data flag for a quote by ID
func UpdateGroupPricingQuoteIndicativeFlag(c *gin.Context) {
	// URL param
	idStr := c.Param("id")
	quoteID, err := strconv.Atoi(idStr)
	if err != nil || quoteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quote id"})
		return
	}

	// Payload can send string or boolean
	type payload struct {
		EnabledAny interface{} `json:"indicative_data_enabled"`
	}
	var body payload
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}

	// Normalize to bool
	var enabled bool
	switch v := body.EnabledAny.(type) {
	case bool:
		enabled = v
	case string:
		if v == "true" || v == "TRUE" || v == "True" || v == "1" {
			enabled = true
		} else if v == "false" || v == "FALSE" || v == "False" || v == "0" {
			enabled = false
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "indicative_data_enabled must be 'true' or 'false'"})
			return
		}
	case float64: // JSON numbers decode to float64
		enabled = v != 0
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported type for indicative_data_enabled"})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	updated, err := services.UpdateGroupPricingQuoteIndicativeFlag(quoteID, enabled, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func DeleteIndicativeData(c *gin.Context) {
	idStr := c.Param("id")
	quoteID, err := strconv.Atoi(idStr)
	if err != nil || quoteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quote id"})
		return
	}
	user := c.MustGet("user").(models.AppUser)

	services.DeleteIndicativeData(quoteID, user)

}

func GetTableDataCsvExport(c *gin.Context) {
	tableType := c.Param("table_type")
	quoteId, _ := strconv.Atoi(c.Param("quote_id"))

	excelData, err := services.GetQuoteTableDataExcel(quoteId, tableType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to zip files: %v", err)})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="export.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)

}

// GetGroupSchemeStatusAudit returns the audit trail for a given group scheme
func GetGroupSchemeStatusAudit(c *gin.Context) {
	schemeId := c.Param("id")
	audits, err := services.GetGroupSchemeStatusAudit(schemeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, audits)
}

// GetGroupPricingAuditLogs returns generic audit logs for group-pricing
func GetGroupPricingAuditLogs(c *gin.Context) {
	entity := c.Query("entity")
	entityID := c.Query("entity_id")
	limit := 50
	offset := 0
	if limStr := c.Query("limit"); limStr != "" {
		if v, err := strconv.Atoi(limStr); err == nil && v >= 0 {
			limit = v
		}
	}
	if offStr := c.Query("offset"); offStr != "" {
		if v, err := strconv.Atoi(offStr); err == nil && v >= 0 {
			offset = v
		}
	}
	rows, err := services.GetAuditLogs("group-pricing", entity, entityID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// GetMembersPaginated returns paginated member data in force with optional filters
// Query params: page, pageSize, search, schemeId, status
func GetMembersPaginated(c *gin.Context) {
	// Parse pagination
	page := 1
	pageSize := 20
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			page = p
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// Filters
	search := c.Query("search")
	schemeId := c.Query("schemeId")
	status := c.Query("status")

	items, total, err := services.GetMembersPaginated(page, pageSize, search, schemeId, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compute total pages without floating point
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	hasMore := false
	if totalPages > 0 && page > 0 {
		hasMore = page < totalPages
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       items,
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": totalPages,
		"hasMore":    hasMore,
	})
}

// -----------------------------
// Beneficiaries - Controllers
// -----------------------------

// List beneficiaries for a member
func GetMemberBeneficiaries(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}

	list, err := services.GetBeneficiariesByMemberID(memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Get a single beneficiary by ID for a member
func GetMemberBeneficiary(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	idStr := c.Param("id")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary id"})
		return
	}
	b, err := services.GetBeneficiaryByID(memberID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "beneficiary not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

// Create a beneficiary for a member
func CreateMemberBeneficiary(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}

	var input models.Beneficiary
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.MemberID = memberID

	user := c.MustGet("user").(models.AppUser)
	created, err := services.CreateBeneficiary(input, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// Update a beneficiary
func UpdateMemberBeneficiary(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	idStr := c.Param("id")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary id"})
		return
	}

	var input models.Beneficiary
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	updated, err := services.UpdateBeneficiary(memberID, id, input, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "beneficiary not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete a beneficiary
func DeleteMemberBeneficiary(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	idStr := c.Param("id")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member_id"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary id"})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	if err := services.DeleteBeneficiary(memberID, id, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateBenefitDocumentType creates a new benefit document type
func CreateBenefitDocumentType(c *gin.Context) {
	var input models.BenefitDocumentType
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := services.CreateBenefitDocumentType(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetBenefitDocumentTypes returns all benefit document types or filtered by benefit_code
func GetBenefitDocumentTypes(c *gin.Context) {
	benefitCode := c.Query("benefit_code")
	var docTypes []models.BenefitDocumentType
	var err error

	if benefitCode != "" {
		docTypes, err = services.GetBenefitDocumentTypesByBenefitCode(benefitCode)
	} else {
		docTypes, err = services.GetBenefitDocumentTypes()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docTypes)
}

// UpdateBenefitDocumentType updates an existing benefit document type
func UpdateBenefitDocumentType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.BenefitDocumentType
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := services.UpdateBenefitDocumentType(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteBenefitDocumentType deletes a benefit document type
func DeleteBenefitDocumentType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := services.DeleteBenefitDocumentType(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetRequiredDocumentsForClaim returns required document types for a specific claim
func GetRequiredDocumentsForClaim(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claim_id"})
		return
	}

	docTypes, err := services.GetBenefitDocumentTypesByClaimID(claimID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "claim not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docTypes)
}

// CRUD for SchemeCategoryMaster

func CreateSchemeCategoryMaster(c *gin.Context) {
	var input models.SchemeCategoryMaster
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	created, err := services.CreateSchemeCategoryMaster(input, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func GetSchemeCategoryMasters(c *gin.Context) {
	insurerIdStr := c.Query("insurer_id")
	var insurerId int
	if insurerIdStr != "" {
		insurerId, _ = strconv.Atoi(insurerIdStr)
	}

	categories, err := services.GetSchemeCategoryMasters(insurerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func GetSchemeCategoryMasterByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	category, err := services.GetSchemeCategoryMasterByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func UpdateSchemeCategoryMaster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.SchemeCategoryMaster
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	updated, err := services.UpdateSchemeCategoryMaster(id, input, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func DeleteSchemeCategoryMaster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := services.DeleteSchemeCategoryMaster(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetRegionsForRiskCode returns distinct region names from RegionLoading
// for the given risk_rate_code query parameter.
func GetRegionsForRiskCode(c *gin.Context) {
	riskRateCode := c.Query("risk_rate_code")
	regions, err := services.GetDistinctRegionsForRiskCode(riskRateCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": regions})
}

// GetDiscountAuthority returns the maximum discount percentage the current user may apply
// for the given risk_rate_code, based on their GP role in the DiscountAuthority table.
func GetDiscountAuthority(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	riskRateCode := c.Param("risk_rate_code")

	da, err := services.GetDiscountAuthorityForUser(user.UserEmail, riskRateCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No discount authority found for user role and risk rate code"})
		return
	}
	c.JSON(http.StatusOK, da)
}

// ApplyDiscountToQuote applies a user-entered discount percentage to all MemberRatingResult
// rows for the quote, then recalculates TotalLoading and all office premiums.

func ApplyDiscountToQuote(c *gin.Context) {
	quoteId := c.Param("id")
	discountPct := utils.StringToFloat(c.Param("discount"))
	user := c.MustGet("user").(models.AppUser)

	// Load the quote to get risk_rate_code for authority check
	var quote models.GroupPricingQuote
	if err := services.DB.Where("id = ?", quoteId).First(&quote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Quote not found"})
		return
	}

	// Validate discount against the user's authority
	da, err := services.GetDiscountAuthorityForUser(user.UserEmail, quote.RiskRateCode)
	// da.MaxDiscount is stored as a decimal (e.g. 0.05); discountPct is in % (e.g. 5)
	if err != nil || discountPct > da.MaxDiscount*100 {
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Discount of %.2f%% exceeds your maximum allowed discount authority", discountPct)})
		return
	}

	if err := services.ApplyDiscountToQuote(quoteId, discountPct, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

// GetGroupPricingSettings returns the singleton GroupPricingSetting row, which
// today holds the global discount calculation method. The row is created with
// defaults on first read so the endpoint always returns a usable payload.
func GetGroupPricingSettings(c *gin.Context) {
	var s models.GroupPricingSetting
	if err := services.DB.First(&s, 1).Error; err != nil {
		s = models.GroupPricingSetting{
			ID:                   1,
			DiscountMethod:       models.DiscountMethodLoadingAdjustment,
			FCLMethod:            models.FCLMethodPercentile,
			FCLOverrideTolerance: services.FCLOverrideToleranceDefault,
			RiskAlrCeilingPct:    100,
			RiskAlrDeltaPp:       20,
		}
		services.DB.Create(&s)
	}
	c.JSON(http.StatusOK, s)
}

// UpdateGroupPricingSettings writes the singleton GroupPricingSetting row.
// The discount method only takes effect on the next ApplyDiscountToQuote /
// recompute call — existing quotes are not retroactively recomputed.
func UpdateGroupPricingSettings(c *gin.Context) {
	// Pointer-typed fields so we can distinguish "field omitted" from
	// "explicitly set to zero" — relevant for FCLOverrideTolerance, where 0
	// is a meaningful (strict) value.
	var payload struct {
		DiscountMethod       *string  `json:"discount_method"`
		FCLMethod            *string  `json:"fcl_method"`
		FCLOverrideTolerance *float64 `json:"fcl_override_tolerance"`
		RiskAlrCeilingPct    *float64 `json:"risk_alr_ceiling_pct"`
		RiskAlrDeltaPp       *float64 `json:"risk_alr_delta_pp"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if payload.DiscountMethod != nil &&
		*payload.DiscountMethod != models.DiscountMethodLoadingAdjustment &&
		*payload.DiscountMethod != models.DiscountMethodProrata {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid discount_method"})
		return
	}
	if payload.FCLMethod != nil &&
		*payload.FCLMethod != models.FCLMethodPercentile &&
		*payload.FCLMethod != models.FCLMethodOutlier {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid fcl_method"})
		return
	}
	if payload.FCLOverrideTolerance != nil &&
		(*payload.FCLOverrideTolerance < 0 || *payload.FCLOverrideTolerance > 1) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "fcl_override_tolerance must be between 0 and 1"})
		return
	}
	if payload.RiskAlrCeilingPct != nil &&
		(*payload.RiskAlrCeilingPct < 0 || *payload.RiskAlrCeilingPct > 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "risk_alr_ceiling_pct must be between 0 and 1000"})
		return
	}
	if payload.RiskAlrDeltaPp != nil &&
		(*payload.RiskAlrDeltaPp < 0 || *payload.RiskAlrDeltaPp > 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "risk_alr_delta_pp must be between 0 and 1000"})
		return
	}
	user := c.MustGet("user").(models.AppUser)

	var s models.GroupPricingSetting
	if err := services.DB.First(&s, 1).Error; err != nil {
		s = models.GroupPricingSetting{
			ID:                   1,
			DiscountMethod:       models.DiscountMethodLoadingAdjustment,
			FCLMethod:            models.FCLMethodPercentile,
			FCLOverrideTolerance: services.FCLOverrideToleranceDefault,
			RiskAlrCeilingPct:    100,
			RiskAlrDeltaPp:       20,
		}
	}
	now := time.Now()
	if payload.DiscountMethod != nil {
		s.DiscountMethod = *payload.DiscountMethod
		s.DiscountMethodUpdatedAt = &now
		s.DiscountMethodUpdatedBy = user.UserEmail
	}
	if payload.FCLMethod != nil {
		s.FCLMethod = *payload.FCLMethod
		s.FCLMethodUpdatedAt = &now
		s.FCLMethodUpdatedBy = user.UserEmail
	}
	if payload.FCLOverrideTolerance != nil {
		s.FCLOverrideTolerance = *payload.FCLOverrideTolerance
	}
	if payload.RiskAlrCeilingPct != nil || payload.RiskAlrDeltaPp != nil {
		if payload.RiskAlrCeilingPct != nil {
			s.RiskAlrCeilingPct = *payload.RiskAlrCeilingPct
		}
		if payload.RiskAlrDeltaPp != nil {
			s.RiskAlrDeltaPp = *payload.RiskAlrDeltaPp
		}
		s.RiskThresholdsUpdatedAt = &now
		s.RiskThresholdsUpdatedBy = user.UserEmail
	}
	s.UpdatedBy = user.UserEmail
	if err := services.DB.Save(&s).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

// GetQuoteWinProbability returns the win probability score for a single quote.
func GetQuoteWinProbability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid quote id"})
		return
	}
	wp, err := services.GetQuoteWinProbabilityScore(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": wp})
}

// TrainWinProbabilityModelHandler triggers a model retraining run.
func TrainWinProbabilityModelHandler(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	go func() {
		if err := services.TrainWinProbabilityModel(user.UserEmail); err != nil {
			log.WithField("error", err.Error()).Error("Win probability training failed")
		}
	}()
	c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "Training job started"})
}

// GetWinProbabilityModelInfo returns metadata about the latest trained model.
func GetWinProbabilityModelInfo(c *gin.Context) {
	info, err := services.GetWinProbabilityModelInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": info})
}

// CheckCustomTieredTableExists checks whether a custom tiered income replacement
// table exists for the given scheme name and risk rate code.
func CheckCustomTieredTableExists(c *gin.Context) {
	schemeName := c.Query("scheme_name")
	riskRateCode := c.Query("risk_rate_code")
	if schemeName == "" || riskRateCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheme_name and risk_rate_code are required"})
		return
	}
	exists := services.HasCustomTieredIncomeReplacementTable(schemeName, riskRateCode)
	c.JSON(http.StatusOK, gin.H{"success": true, "exists": exists})
}

// RequestCustomTieredTable sends notifications to system admins that a custom
// tiered income replacement table is needed for a scheme.
func RequestCustomTieredTable(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req struct {
		SchemeName   string `json:"scheme_name" binding:"required"`
		SchemeID     int    `json:"scheme_id" binding:"required"`
		RiskRateCode string `json:"risk_rate_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go services.NotifyCustomTieredTableRequested(req.SchemeName, req.SchemeID, req.RiskRateCode, user)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Custom tiered income replacement table request sent to system administrators"})
}

// validExperienceRateOverrideBenefit returns true when the supplied benefit
// code is one of the seven supported per-benefit override slots.
func validExperienceRateOverrideBenefit(b string) bool {
	switch b {
	case models.ExperienceRateOverrideBenefitGla,
		models.ExperienceRateOverrideBenefitAagla,
		models.ExperienceRateOverrideBenefitSgla,
		models.ExperienceRateOverrideBenefitPtd,
		models.ExperienceRateOverrideBenefitTtd,
		models.ExperienceRateOverrideBenefitPhi,
		models.ExperienceRateOverrideBenefitCi,
		models.ExperienceRateOverrideBenefitFun:
		return true
	}
	return false
}

// quoteOverridesLocked returns true when the quote's status forbids edits
// to its experience-rate overrides. Approved, accepted, and in-force quotes
// represent settled commercial commitments; mutating overrides on them
// would silently change the priced premium without a fresh approval cycle.
func quoteOverridesLocked(quoteID int) (bool, models.Status, error) {
	var quote models.GroupPricingQuote
	if err := services.DB.Select("status").
		Where("id = ?", quoteID).
		First(&quote).Error; err != nil {
		return false, "", err
	}
	switch quote.Status {
	case models.StatusApproved, models.StatusAccepted, models.StatusInForce:
		return true, quote.Status, nil
	}
	return false, quote.Status, nil
}

// GetExperienceRateOverrides returns every per-(category, benefit) override
// row stored against a quote. An empty array is returned when the quote has
// no overrides yet — the frontend treats this as "show empty state".
func GetExperienceRateOverrides(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil || quoteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quote id"})
		return
	}
	rows := []models.GroupPricingExperienceRateOverride{}
	if err := services.DB.Where("quote_id = ?", quoteID).
		Order("scheme_category, benefit").
		Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// SaveExperienceRateOverrides upserts an array of (quote, category, benefit)
// override rows. The endpoint replaces the entire set for the quote (delete +
// insert) so the frontend can submit the in-memory list verbatim.
func SaveExperienceRateOverrides(c *gin.Context) {
	var rows []models.GroupPricingExperienceRateOverride
	if err := c.ShouldBindJSON(&rows); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one override row is required"})
		return
	}
	quoteID := rows[0].QuoteId
	if quoteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quote_id is required on every row"})
		return
	}
	if locked, status, err := quoteOverridesLocked(quoteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if locked {
		c.JSON(http.StatusConflict, gin.H{
			"error": "experience-rate overrides cannot be edited while the quote is " + string(status),
		})
		return
	}
	for i := range rows {
		if rows[i].QuoteId != quoteID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "all rows must share the same quote_id"})
			return
		}
		if rows[i].SchemeCategory == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "scheme_category is required"})
			return
		}
		if !validExperienceRateOverrideBenefit(rows[i].Benefit) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid benefit: " + rows[i].Benefit})
			return
		}
		if rows[i].Mode != models.ExperienceRateOverrideModeTheoretical &&
			rows[i].Mode != models.ExperienceRateOverrideModeExperienceRated {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mode: " + rows[i].Mode})
			return
		}
		if rows[i].Mode == models.ExperienceRateOverrideModeTheoretical {
			rows[i].OverrideRate = 0
		}
		if rows[i].OverrideRate < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "override_rate must be >= 0"})
			return
		}
		if rows[i].Credibility < 0 || rows[i].Credibility > 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "credibility must be between 0 and 1"})
			return
		}
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.SaveExperienceRateOverrides(quoteID, rows, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	persisted := []models.GroupPricingExperienceRateOverride{}
	services.DB.Where("quote_id = ?", quoteID).
		Order("scheme_category, benefit").
		Find(&persisted)
	c.JSON(http.StatusOK, persisted)
}

// UpdateExperienceOverrideCredibility persists the manually-entered
// credibility (0-1) the actuary supplies alongside experience-rate overrides.
// Locked quotes (approved / accepted / in_force) reject the change.
func UpdateExperienceOverrideCredibility(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil || quoteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quote id"})
		return
	}
	var payload struct {
		Credibility float64 `json:"credibility"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Credibility < 0 || payload.Credibility > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credibility must be between 0 and 1"})
		return
	}
	if locked, status, err := quoteOverridesLocked(quoteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if locked {
		c.JSON(http.StatusConflict, gin.H{
			"error": "credibility cannot be changed while the quote is " + string(status),
		})
		return
	}
	if err := services.DB.Model(&models.GroupPricingQuote{}).
		Where("id = ?", quoteID).
		Update("experience_override_credibility", payload.Credibility).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":                         true,
		"quote_id":                        quoteID,
		"experience_override_credibility": payload.Credibility,
	})
}

// DeleteExperienceRateOverrides removes every override row for a quote,
// reverting the quote to "no override entries yet". Used by the Delete-All
// button in the UI.
func DeleteExperienceRateOverrides(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil || quoteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quote id"})
		return
	}
	if locked, status, err := quoteOverridesLocked(quoteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if locked {
		c.JSON(http.StatusConflict, gin.H{
			"error": "experience-rate overrides cannot be deleted while the quote is " + string(status),
		})
		return
	}
	if err := services.DB.Where("quote_id = ?", quoteID).
		Delete(&models.GroupPricingExperienceRateOverride{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SaveAdditionalGlaCoverSmoothedRates persists smoothed OfficeRate/1000 values
// (and per-cell smoothing factors) for one scheme category's Additional GLA
// Cover bands. Body shape:
//
//	{
//	  "category": "Main",
//	  "rows": [
//	    {
//	      "min_age": 18, "max_age": 24,
//	      "smoothed_office_rate_per1000_male": 2.31,
//	      "smoothing_factor": 0.95
//	    }, ...
//	  ]
//	}
func SaveAdditionalGlaCoverSmoothedRates(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("id"))
	if err != nil || quoteID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quote id"})
		return
	}

	var body struct {
		Category string                                  `json:"category"`
		Rows     []services.AdditionalGlaSmoothedRateRow `json:"rows"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(body.Category) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}

	// Prefer the user's display name for the audit; fall back to email so the
	// banner still has a useful identifier when the name isn't populated.
	updatedBy := ""
	if u, ok := c.Get("user"); ok {
		if user, ok := u.(models.AppUser); ok {
			if strings.TrimSpace(user.UserName) != "" {
				updatedBy = user.UserName
			} else {
				updatedBy = user.UserEmail
			}
		}
	}

	result, err := services.SaveAdditionalGlaCoverSmoothedRates(quoteID, body.Category, body.Rows, updatedBy)
	if err != nil {
		if errors.Is(err, services.ErrAglaSmoothingLocked) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"category":                           body.Category,
		"additional_gla_cover_band_rates":    result.BandRates,
		"additional_gla_smoothed_updated_at": result.UpdatedAt,
		"additional_gla_smoothed_updated_by": result.UpdatedBy,
	})
}

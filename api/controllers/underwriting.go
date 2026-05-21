package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"api/models"
	"api/services"
)

// buildUnderwritingAttachmentViewerURL returns a relative URL the renderer
// uses to stream an attachment for inline viewing or download.
func buildUnderwritingAttachmentViewerURL(attachmentID int) string {
	return fmt.Sprintf("/group-pricing/underwriting/attachments/%d/download", attachmentID)
}

func populateUnderwritingViewerURLs(atts []models.UnderwritingCaseAttachment) {
	for i := range atts {
		atts[i].ViewerURL = buildUnderwritingAttachmentViewerURL(atts[i].ID)
	}
}

// ListUnderwritingCaseQuoteSummaries returns the quote-grouped queue used
// as the top-level "Underwriting" view. One row per quote that has at
// least one case, with case counts broken down by status. Filters
// (status / tier / assignee / quote_id) apply to the inner row set
// before aggregation.
func ListUnderwritingCaseQuoteSummaries(c *gin.Context) {
	filter := services.UnderwritingCaseFilter{
		Status:                   models.UWCaseStatus(c.Query("status")),
		AssignedUnderwriterEmail: c.Query("assignee"),
	}
	if v := c.Query("quote_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid quote_id")
			return
		}
		filter.QuoteID = &id
	}
	if v := c.Query("tier"); v != "" {
		tier, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid tier")
			return
		}
		filter.Tier = &tier
	}
	rows, err := services.ListUnderwritingCaseQuoteSummaries(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// ListUnderwritingCases returns underwriting cases filtered by optional
// quote_id, status, tier and assignee query parameters.
func ListUnderwritingCases(c *gin.Context) {
	filter := services.UnderwritingCaseFilter{
		Status:                   models.UWCaseStatus(c.Query("status")),
		AssignedUnderwriterEmail: c.Query("assignee"),
	}
	if v := c.Query("quote_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid quote_id")
			return
		}
		filter.QuoteID = &id
	}
	if v := c.Query("tier"); v != "" {
		tier, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, "invalid tier")
			return
		}
		filter.Tier = &tier
	}
	cases, err := services.ListUnderwritingCases(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, cases)
}

// GetUnderwritingCase returns one case with decisions, events, attachments.
func GetUnderwritingCase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	uw, err := services.GetUnderwritingCase(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "case not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	populateUnderwritingViewerURLs(uw.Attachments)
	c.JSON(http.StatusOK, uw)
}

// AssignUnderwritingCase takes a JSON body {"assignee_email": "..."} and sets
// the case assignee. Pass an empty string to unassign.
func AssignUnderwritingCase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	var payload struct {
		AssigneeEmail string `json:"assignee_email"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	uw, err := services.AssignUnderwritingCase(id, payload.AssigneeEmail, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "case not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, uw)
}

// TransitionUnderwritingCase moves a case to a new status. Body:
// {"status": "in_review|decided|postponed|declined|pending_evidence", "note": "..."}.
func TransitionUnderwritingCase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	var payload struct {
		Status string `json:"status" binding:"required"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	uw, err := services.TransitionUnderwritingCase(id, models.UWCaseStatus(payload.Status), user, payload.Note)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "case not found")
			return
		}
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, uw)
}

// CreateUnderwritingDecision appends a per-benefit decision to a case.
func CreateUnderwritingDecision(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	var payload models.UnderwritingDecision
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	decision, err := services.CreateUnderwritingDecision(id, payload, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "case not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, decision)
}

// UploadUnderwritingCaseAttachments accepts multipart/form-data with one or
// more files plus per-file `kind` form values (in upload order).
func UploadUnderwritingCaseAttachments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	form, formErr := c.MultipartForm()
	if formErr != nil || form == nil {
		c.JSON(http.StatusBadRequest, "multipart/form-data with at least one file is required")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	created, err := services.AppendCaseAttachments(id, form.File, form.Value, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "case not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	populateUnderwritingViewerURLs(created)
	c.JSON(http.StatusCreated, created)
}

// DownloadUnderwritingCaseAttachment streams an attachment by ID.
func DownloadUnderwritingCaseAttachment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("attachment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid attachment_id")
		return
	}
	var att models.UnderwritingCaseAttachment
	if err := services.DB.Where("id = ?", id).First(&att).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "attachment not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := os.Stat(att.StoragePath); err != nil {
		c.JSON(http.StatusNotFound, "file missing on disk")
		return
	}
	c.FileAttachment(att.StoragePath, att.FileName)
}

// RerateQuoteFromUWDecisions forces a Phase-4 re-rate for a quote. Useful
// when a previous re-rate failed or when decisions exist on cases that were
// recreated. Returns the resulting QuoteReRateEvent.
func RerateQuoteFromUWDecisions(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid quote_id")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	var payload struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&payload)
	if payload.Reason == "" {
		payload.Reason = "Manual re-rate"
	}
	event, err := services.ApplyDecisionsAndReRate(quoteID, user, payload.Reason, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, event)
}

// ListQuoteReRateEvents returns the re-rate history for a quote, newest first.
func ListQuoteReRateEvents(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid quote_id")
		return
	}
	var events []models.QuoteReRateEvent
	if err := services.DB.Where("quote_id = ?", quoteID).Order("triggered_at DESC").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, events)
}

// RecreateUnderwritingCasesForQuote forces case creation for a quote (useful
// for backfilling pre-Phase-2 quotes once UnderwritingTier is populated).
// Idempotent.
func RecreateUnderwritingCasesForQuote(c *gin.Context) {
	quoteID, err := strconv.Atoi(c.Param("quote_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid quote_id")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	result, err := services.CreateCasesForQuote(quoteID, user.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

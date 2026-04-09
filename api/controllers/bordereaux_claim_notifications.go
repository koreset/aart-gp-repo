package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateClaimNotification godoc
// POST /group-pricing/bordereaux/claim-notifications
func CreateClaimNotification(c *gin.Context) {
	var req models.CreateClaimNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	n, err := services.CreateClaimNotification(req, user)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, n)
}

// GetClaimNotifications godoc
// GET /group-pricing/bordereaux/claim-notifications?claim_id=&scheme_id=&reinsurer_code=&notification_type=&status=&claim_number=&page=&page_size=
func GetClaimNotifications(c *gin.Context) {
	claimID, _ := strconv.Atoi(c.Query("claim_id"))
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	records, total, err := services.GetClaimNotifications(
		claimID, schemeID, page, pageSize,
		c.Query("reinsurer_code"),
		c.Query("notification_type"),
		c.Query("status"),
		c.Query("claim_number"),
	)
	if err != nil {
		InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:  true,
		Data:     records,
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
	})
}

// MarkNotificationSent godoc
// POST /group-pricing/bordereaux/claim-notifications/:id/sent
func MarkNotificationSent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		BadRequestMsg(c, "invalid id")
		return
	}
	var req models.MarkNotificationSentRequest
	_ = c.ShouldBindJSON(&req)
	user := c.MustGet("user").(models.AppUser)
	n, svcErr := services.MarkNotificationSent(id, req.Notes, user)
	if svcErr != nil {
		BadRequest(c, svcErr)
		return
	}
	OK(c, n)
}

// MarkNotificationAcknowledged godoc
// POST /group-pricing/bordereaux/claim-notifications/:id/acknowledged
func MarkNotificationAcknowledged(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		BadRequestMsg(c, "invalid id")
		return
	}
	var req models.MarkNotificationSentRequest
	_ = c.ShouldBindJSON(&req)
	user := c.MustGet("user").(models.AppUser)
	n, svcErr := services.MarkNotificationAcknowledged(id, req.Notes, user)
	if svcErr != nil {
		BadRequest(c, svcErr)
		return
	}
	OK(c, n)
}

// GenerateMonthEndNotifications godoc
// POST /group-pricing/bordereaux/claim-notifications/generate-month-end
func GenerateMonthEndNotifications(c *gin.Context) {
	var req models.GenerateMonthEndNotificationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	user := c.MustGet("user").(models.AppUser)
	created, err := services.GenerateMonthEndNotifications(req.SchemeID, req.Month, req.Year, user)
	if err != nil {
		InternalError(c, err)
		return
	}
	// Return count alongside the data for the frontend
	c.JSON(http.StatusOK, models.PremiumResponse{
		Success: true,
		Data:    gin.H{"records": created, "count": len(created)},
	})
}

// DeleteClaimNotification godoc
// DELETE /group-pricing/bordereaux/claim-notifications/:id
func DeleteClaimNotification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		BadRequestMsg(c, "invalid id")
		return
	}
	if err := services.DeleteClaimNotification(id); err != nil {
		BadRequest(c, err)
		return
	}
	OKMsg(c, "notification deleted")
}

// GetNotificationStats godoc
// GET /group-pricing/bordereaux/claim-notifications/stats?scheme_id=
func GetNotificationStats(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	stats, err := services.GetNotificationStats(schemeID)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, stats)
}

// GetClaimsByScheme godoc
// GET /group-pricing/bordereaux/claim-notifications/claims-by-scheme/:scheme_id
func GetClaimsByScheme(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("scheme_id"))
	if err != nil || schemeID == 0 {
		BadRequestMsg(c, "invalid scheme_id")
		return
	}
	items, svcErr := services.GetClaimsByScheme(schemeID)
	if svcErr != nil {
		InternalError(c, svcErr)
		return
	}
	OK(c, items)
}

// ExportNotificationsCSV godoc
// GET /group-pricing/bordereaux/claim-notifications/export?scheme_id=&status=&notification_type=
func ExportNotificationsCSV(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	csvBytes, err := services.ExportNotificationsCSV(
		schemeID,
		c.Query("status"),
		c.Query("notification_type"),
	)
	if err != nil {
		InternalError(c, err)
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=claim_notifications.csv")
	c.Data(http.StatusOK, "text/csv", csvBytes)
}

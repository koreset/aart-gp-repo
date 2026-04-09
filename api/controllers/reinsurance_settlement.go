package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GenerateTechnicalAccount handles POST /group-pricing/reinsurance/settlement
func GenerateTechnicalAccount(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.GenerateTechnicalAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := services.GenerateTechnicalAccount(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": account})
}

// GetTechnicalAccounts handles GET /group-pricing/reinsurance/settlement
func GetTechnicalAccounts(c *gin.Context) {
	treatyIDStr := c.Query("treaty_id")
	status := c.Query("status")
	treatyID := 0
	if treatyIDStr != "" {
		treatyID, _ = strconv.Atoi(treatyIDStr)
	}
	accounts, err := services.GetTechnicalAccounts(treatyID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": accounts})
}

// GetTechnicalAccountByID handles GET /group-pricing/reinsurance/settlement/:id
func GetTechnicalAccountByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	account, err := services.GetTechnicalAccountByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": account})
}

// UpdateTechnicalAccount handles PATCH /group-pricing/reinsurance/settlement/:id
func UpdateTechnicalAccount(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.UpdateTechnicalAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := services.UpdateTechnicalAccount(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": account})
}

// RecordSettlementPayment handles POST /group-pricing/reinsurance/settlement/payments
func RecordSettlementPayment(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.RecordSettlementPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment, err := services.RecordSettlementPayment(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": payment})
}

// GetSettlementPayments handles GET /group-pricing/reinsurance/settlement/payments
func GetSettlementPayments(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	accountID := 0
	if accountIDStr != "" {
		accountID, _ = strconv.Atoi(accountIDStr)
	}
	payments, err := services.GetSettlementPayments(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": payments})
}

// GetSettlementStats handles GET /group-pricing/reinsurance/settlement/stats
func GetSettlementStats(c *gin.Context) {
	treatyIDStr := c.Query("treaty_id")
	treatyID := 0
	if treatyIDStr != "" {
		treatyID, _ = strconv.Atoi(treatyIDStr)
	}
	stats, err := services.GetSettlementStats(treatyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

// EscalateSettlementDispute handles POST /group-pricing/reinsurance/settlement/:id/escalate-dispute
func EscalateSettlementDispute(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.EscalateDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := services.EscalateSettlementDispute(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": account})
}

// ResolveSettlementDispute handles POST /group-pricing/reinsurance/settlement/:id/resolve-dispute
func ResolveSettlementDispute(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.ResolveDisputeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := services.ResolveSettlementDispute(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": account})
}

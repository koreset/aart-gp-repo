package controllers

import (
	"api/models"
	"api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateReinsurerAcceptance handles POST /group-pricing/bordereaux/reinsurer/acceptances
func CreateReinsurerAcceptance(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.CreateReinsurerAcceptanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	acc, err := services.CreateReinsurerAcceptance(req, user)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, acc)
}

// GetReinsurerAcceptances handles GET /group-pricing/bordereaux/reinsurer/acceptances
func GetReinsurerAcceptances(c *gin.Context) {
	generatedID := c.Query("generated_id")
	status := c.Query("status")
	records, err := services.GetReinsurerAcceptances(generatedID, status)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, records)
}

// UpdateReinsurerAcceptance handles PATCH /group-pricing/bordereaux/reinsurer/acceptances/:id
func UpdateReinsurerAcceptance(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid id")
		return
	}
	var req models.UpdateReinsurerAcceptanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	acc, err := services.UpdateReinsurerAcceptance(id, req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, acc)
}

// GetAcceptanceStats handles GET /group-pricing/bordereaux/reinsurer/acceptances/stats
func GetAcceptanceStats(c *gin.Context) {
	generatedID := c.Query("generated_id")
	if generatedID == "" {
		BadRequestMsg(c, "generated_id is required")
		return
	}
	stats, err := services.GetAcceptanceStats(generatedID)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, stats)
}

// CreateReinsurerRecovery handles POST /group-pricing/bordereaux/reinsurer/recoveries
func CreateReinsurerRecovery(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.CreateReinsurerRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	rec, err := services.CreateReinsurerRecovery(req, user)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rec)
}

// GetReinsurerRecoveries handles GET /group-pricing/bordereaux/reinsurer/recoveries
func GetReinsurerRecoveries(c *gin.Context) {
	generatedID := c.Query("generated_id")
	claimRef := c.Query("claim_ref")
	status := c.Query("status")
	records, err := services.GetReinsurerRecoveries(generatedID, claimRef, status)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, records)
}

// UpdateReinsurerRecovery handles PATCH /group-pricing/bordereaux/reinsurer/recoveries/:id
func UpdateReinsurerRecovery(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		BadRequestMsg(c, "invalid id")
		return
	}
	var req models.UpdateReinsurerRecoveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	rec, err := services.UpdateReinsurerRecovery(id, req, user)
	if err != nil {
		BadRequest(c, err)
		return
	}
	OK(c, rec)
}

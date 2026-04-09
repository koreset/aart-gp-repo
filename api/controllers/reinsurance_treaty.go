package controllers

import (
	"api/models"
	"api/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTreaty handles POST /group-pricing/reinsurance/treaties
func CreateTreaty(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	var req models.CreateTreatyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	treaty, err := services.CreateTreaty(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": treaty})
}

// GetTreaties handles GET /group-pricing/reinsurance/treaties
func GetTreaties(c *gin.Context) {
	status := c.Query("status")
	treatyType := c.Query("type")
	reinsurer := c.Query("reinsurer")
	treaties, err := services.GetTreaties(status, treatyType, reinsurer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": treaties})
}

// GetTreatyByID handles GET /group-pricing/reinsurance/treaties/:id
func GetTreatyByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	treaty, err := services.GetTreatyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": treaty})
}

// UpdateTreaty handles PUT /group-pricing/reinsurance/treaties/:id
func UpdateTreaty(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req models.UpdateTreatyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	treaty, err := services.UpdateTreaty(id, req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": treaty})
}

// DeleteTreaty handles DELETE /group-pricing/reinsurance/treaties/:id
func DeleteTreaty(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := services.DeleteTreaty(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "treaty deleted"})
}

// LinkSchemeToTreaty handles POST /group-pricing/reinsurance/treaties/:id/schemes
func LinkSchemeToTreaty(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	treatyID, err := strconv.Atoi(c.Param("id"))
	if err != nil || treatyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid treaty id"})
		return
	}
	var req models.LinkSchemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link, err := services.LinkSchemeToTreaty(treatyID, req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": link})
}

// BulkLinkSchemesToTreaty handles POST /group-pricing/reinsurance/treaties/:id/schemes/bulk
func BulkLinkSchemesToTreaty(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	treatyID, err := strconv.Atoi(c.Param("id"))
	if err != nil || treatyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid treaty id"})
		return
	}
	var req models.BulkLinkSchemesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := services.BulkLinkSchemesToTreaty(treatyID, req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("%d scheme(s) linked", created), "data": created})
}

// GetTreatySchemeLinks handles GET /group-pricing/reinsurance/treaties/:id/schemes
func GetTreatySchemeLinks(c *gin.Context) {
	treatyID, err := strconv.Atoi(c.Param("id"))
	if err != nil || treatyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid treaty id"})
		return
	}
	links, err := services.GetTreatySchemeLinks(treatyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": links})
}

// RemoveSchemeTreatyLink handles DELETE /group-pricing/reinsurance/treaties/scheme-links/:link_id
func RemoveSchemeTreatyLink(c *gin.Context) {
	linkID, err := strconv.Atoi(c.Param("link_id"))
	if err != nil || linkID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
		return
	}
	if err := services.RemoveSchemeTreatyLink(linkID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "scheme link removed"})
}

// BulkRemoveSchemeLinks handles DELETE /group-pricing/reinsurance/treaties/:id/schemes/bulk
func BulkRemoveSchemeLinks(c *gin.Context) {
	var req struct {
		LinkIDs []int `json:"link_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link_ids is required"})
		return
	}
	removed, err := services.BulkRemoveSchemeLinks(req.LinkIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "removed": removed})
}

// GetActiveTreatiesForScheme handles GET /group-pricing/reinsurance/treaties/active/scheme/:scheme_id
func GetActiveTreatiesForScheme(c *gin.Context) {
	schemeID, err := strconv.Atoi(c.Param("scheme_id"))
	if err != nil || schemeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheme id"})
		return
	}
	treaties, err := services.GetActiveTreatiesForScheme(schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": treaties})
}

// GetTreatyStats handles GET /group-pricing/reinsurance/treaties/stats
func GetTreatyStats(c *gin.Context) {
	stats, err := services.GetTreatyStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

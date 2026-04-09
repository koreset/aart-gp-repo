package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBordereauxDeadlines handles GET /bordereaux/deadlines
// Query params: scheme_id, month, year, status
func GetBordereauxDeadlines(c *gin.Context) {
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	status := c.Query("status")

	deadlines, err := services.GetBordereauxDeadlines(schemeID, month, year, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": deadlines})
}

// CreateBordereauxDeadline handles POST /bordereaux/deadlines
func CreateBordereauxDeadline(c *gin.Context) {
	var req models.CreateDeadlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	dl, err := services.CreateBordereauxDeadline(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": dl})
}

// GenerateBordereauxDeadlines handles POST /bordereaux/deadlines/generate
func GenerateBordereauxDeadlines(c *gin.Context) {
	var req models.GenerateDeadlinesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	result, err := services.GenerateDeadlinesForAllSchemes(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// UpdateDeadlineStatus handles PATCH /bordereaux/deadlines/:id/status
func UpdateDeadlineStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.UpdateDeadlineStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	dl, err := services.UpdateDeadlineStatus(id, req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": dl})
}

// GetDeadlineStats handles GET /bordereaux/deadlines/stats
func GetDeadlineStats(c *gin.Context) {
	stats, err := services.GetDeadlineStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

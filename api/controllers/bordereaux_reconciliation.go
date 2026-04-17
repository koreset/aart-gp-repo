package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ResolveDiscrepancy handles POST /group-pricing/bordereaux/reconciliation/results/:id/resolve
func ResolveDiscrepancy(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	var req services.ResolveDiscrepancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	result, svcErr := services.ResolveDiscrepancy(id, req, user)
	if svcErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": svcErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// EscalateDiscrepancy handles POST /group-pricing/bordereaux/reconciliation/results/:id/escalate
func EscalateDiscrepancy(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	var req services.EscalateDiscrepancyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	result, svcErr := services.EscalateDiscrepancy(id, req, user)
	if svcErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": svcErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// ListEscalations handles GET /group-pricing/bordereaux/reconciliation/escalations
// Query params: assigned_to, priority, overdue_only=true.
func ListEscalations(c *gin.Context) {
	filter := services.ListEscalationsRequest{
		AssignedTo:  c.Query("assigned_to"),
		Priority:    c.Query("priority"),
		OverdueOnly: c.Query("overdue_only") == "true",
	}
	results, err := services.ListEscalations(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

// AddDiscrepancyComment handles POST /group-pricing/bordereaux/reconciliation/results/:id/comment
func AddDiscrepancyComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	var req services.AddDiscrepancyCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	result, svcErr := services.AddDiscrepancyComment(id, req, user)
	if svcErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": svcErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// ConfirmReconciliation handles POST /group-pricing/bordereaux/confirmations/:id/confirm
func ConfirmReconciliation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if svcErr := services.ConfirmReconciliation(id, user); svcErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": svcErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "reconciliation confirmed"})
}

// ReprocessReconciliation handles POST /group-pricing/bordereaux/confirmations/:id/reprocess
func ReprocessReconciliation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	summary, svcErr := services.ReprocessReconciliation(id, user)
	if svcErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": svcErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

// AddReconciliationNote handles POST /group-pricing/bordereaux/confirmations/:id/note
func AddReconciliationNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	var req struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	confirmation, svcErr := services.AddReconciliationNote(id, req.Note, user)
	if svcErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": svcErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": confirmation})
}

// GetReconciliationNotes handles GET /group-pricing/bordereaux/confirmations/:id/notes
func GetReconciliationNotes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}
	notes, svcErr := services.GetReconciliationNotes(id)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": svcErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": notes})
}

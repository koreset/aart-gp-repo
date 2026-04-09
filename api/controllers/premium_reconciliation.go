package controllers

import (
	"api/models"
	"api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RunAutoMatchV2 executes the multi-strategy matching engine.
func RunAutoMatchV2(c *gin.Context) {
	var req models.RunAutoMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body (match everything)
		req = models.RunAutoMatchRequest{}
	}

	user := c.MustGet("user").(models.AppUser)
	result, err := services.RunAutoMatch(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: result})
}

// AllocatePaymentV2 manually allocates a payment to one or more invoices.
func AllocatePaymentV2(c *gin.Context) {
	var req models.AllocatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	allocations, err := services.AllocatePayment(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: allocations})
}

// ReverseAllocations reverses one or more payment allocations.
func ReverseAllocations(c *gin.Context) {
	var req models.ReverseAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	reversals, err := services.ReverseAllocations(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: reversals})
}

// WriteOffBalance writes off a small remaining balance.
func WriteOffBalance(c *gin.Context) {
	var req models.WriteOffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	alloc, err := services.WriteOffBalance(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: alloc})
}

// RefundOverpayment records a refund for an overpaid amount.
func RefundOverpayment(c *gin.Context) {
	var req models.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	user := c.MustGet("user").(models.AppUser)
	alloc, err := services.RefundOverpayment(req, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: alloc})
}

// GetReconciliationSummaryV2 returns the reconciliation dashboard.
func GetReconciliationSummaryV2(c *gin.Context) {
	summary, err := services.GetReconciliationSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: summary})
}

// GetReconciliationItems returns filtered reconciliation items.
func GetReconciliationItems(c *gin.Context) {
	itemType := c.Query("type")
	status := c.Query("status")
	schemeID, _ := strconv.Atoi(c.Query("scheme_id"))

	items, err := services.GetReconciliationItems(itemType, status, schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: items})
}

// GetAllocationHistory returns the full allocation trail for an entity.
func GetAllocationHistory(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityID, _ := strconv.Atoi(c.Param("entity_id"))

	if entityType != "payment" && entityType != "invoice" {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: "entity_type must be 'payment' or 'invoice'"})
		return
	}

	history, err := services.GetAllocationHistory(entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: history})
}

// GetReconciliationRunsV2 returns a paginated list of reconciliation runs.
func GetReconciliationRunsV2(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	runs, total, err := services.GetReconciliationRuns(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       runs,
		Total:      int(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// GetReconciliationRunDetailV2 returns a run with its allocations.
func GetReconciliationRunDetailV2(c *gin.Context) {
	runID, _ := strconv.Atoi(c.Param("run_id"))

	detail, err := services.GetReconciliationRunDetail(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: detail})
}

// RollbackRunV2 reverses all allocations in a run.
func RollbackRunV2(c *gin.Context) {
	runID, _ := strconv.Atoi(c.Param("run_id"))
	user := c.MustGet("user").(models.AppUser)

	if err := services.RollbackRun(runID, user); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Run rolled back successfully"})
}

// ReassignReconciliationItemV2 reassigns a reconciliation item.
func ReassignReconciliationItemV2(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("item_id"))

	var req models.ReassignItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	if err := services.ReassignReconciliationItem(itemID, req); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Item reassigned"})
}

// SuspendReconciliationItemV2 moves an item to suspense.
func SuspendReconciliationItemV2(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("item_id"))

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	if err := services.SuspendReconciliationItem(itemID, body.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Item moved to suspense"})
}

// GetMatchingRulesV2 returns matching rules for a rule set.
func GetMatchingRulesV2(c *gin.Context) {
	ruleSet := c.DefaultQuery("rule_set", "default")

	rules, err := services.GetMatchingRules(ruleSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: rules})
}

// SaveMatchingRuleV2 creates or updates a matching rule.
func SaveMatchingRuleV2(c *gin.Context) {
	var rule models.MatchingRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	saved, err := services.SaveMatchingRule(rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Data: saved})
}

// DeleteMatchingRuleV2 deletes a matching rule.
func DeleteMatchingRuleV2(c *gin.Context) {
	ruleID, _ := strconv.Atoi(c.Param("rule_id"))

	if err := services.DeleteMatchingRule(ruleID); err != nil {
		c.JSON(http.StatusInternalServerError, models.PremiumResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PremiumResponse{Success: true, Message: "Rule deleted"})
}

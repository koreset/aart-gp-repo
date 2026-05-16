package controllers

import (
	"api/models"
	"api/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RunFraudCheck scores a claim with the GLM + rule engine and returns the
// final risk level along with the rationale.
func RunFraudCheck(c *gin.Context) {
	claimIDStr := c.Param("claim_id")
	claimID, err := strconv.Atoi(claimIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid claim_id")
		return
	}
	user := c.MustGet("user").(models.AppUser)
	result, err := services.ScoreClaim(claimID, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "claim not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetFraudRiskModel returns the singleton GLM model row.
func GetFraudRiskModel(c *gin.Context) {
	model, err := services.LoadFraudRiskModel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, model)
}

// RefitFraudRiskModel runs a GLM refit from labelled historical claims.
func RefitFraudRiskModel(c *gin.Context) {
	user := c.MustGet("user").(models.AppUser)
	res, err := services.RefitFraudRiskModel(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetFraudFeatureCatalogue returns the feature vocabulary the rule editor and
// admin coefficient table render against.
func GetFraudFeatureCatalogue(c *gin.Context) {
	c.JSON(http.StatusOK, services.FraudFeatureCatalogue)
}

// ListFraudRiskRules lists every rule (enabled or not) for the admin UI.
func ListFraudRiskRules(c *gin.Context) {
	rules, err := services.ListFraudRiskRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rules)
}

// CreateFraudRiskRule creates a new rule.
func CreateFraudRiskRule(c *gin.Context) {
	var payload models.FraudRiskRule
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	created, err := services.CreateFraudRiskRule(payload, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateFraudRiskRule updates an existing rule.
func UpdateFraudRiskRule(c *gin.Context) {
	idStr := c.Param("rule_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_id")
		return
	}
	var payload models.FraudRiskRule
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	updated, err := services.UpdateFraudRiskRule(id, payload, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "rule not found")
			return
		}
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteFraudRiskRule removes a rule.
func DeleteFraudRiskRule(c *gin.Context) {
	idStr := c.Param("rule_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid rule_id")
		return
	}
	if err := services.DeleteFraudRiskRule(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "rule not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

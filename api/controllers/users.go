package controllers

import (
	"api/log"
	"api/models"
	"api/services"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetOrgUsers retrieves users for an organization
func GetOrgUsers(c *gin.Context) {
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
	logger.Info("Processing GetOrgUsers request")

	var organisation models.Organisation
	err := c.BindJSON(&organisation)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to bind JSON for organisation")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.WithField("organisation", organisation.Name).Info("Finding users for organisation")
	response := services.FindOrgUsers(organisation)

	if response == nil {
		logger.Warn("No users found for organisation")
		c.JSON(http.StatusOK, []models.OrgUser{})
		return
	}

	logger.WithField("user_count", len(response)).Info("Successfully retrieved organisation users")
	c.JSON(http.StatusOK, response)
}

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

// RefreshOrgUsers force-refreshes the local org_users cache for an
// organisation by re-fetching from the license server. Intended for
// admin use when new licenses have been provisioned upstream and need to
// appear locally without waiting for a cache miss.
func RefreshOrgUsers(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}

	logger := log.WithContext(ctx)
	logger.Info("Processing RefreshOrgUsers request")

	var organisation models.Organisation
	if err := c.BindJSON(&organisation); err != nil {
		logger.WithField("error", err.Error()).Error("Failed to bind JSON for organisation")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.WithField("organisation", organisation.Name).Info("Refreshing users for organisation")
	users, err := services.RefreshOrgUsers(organisation)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to refresh organisation users")
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	logger.WithField("user_count", len(users)).Info("Successfully refreshed organisation users")
	c.JSON(http.StatusOK, users)
}

package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListEmailOutbox returns paginated outbox rows for the caller's license.
// Query params: status, page, page_size.
func ListEmailOutbox(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	if licenseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-License-Id header is required"})
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	result, err := services.ListEmailOutbox(services.ListEmailOutboxInput{
		LicenseId: licenseId,
		Status:    c.Query("status"),
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetEmailOutboxItem returns one outbox row by id.
func GetEmailOutboxItem(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	row, err := services.GetEmailOutboxByID(licenseId, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "outbox row not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, row)
}

// RetryEmailOutbox resets the row to pending so the worker picks it up.
func RetryEmailOutbox(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	row, err := services.RetryEmailOutbox(licenseId, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "outbox row not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, row)
}

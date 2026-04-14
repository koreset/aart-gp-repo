package controllers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"api/log"
	"api/services"
	"api/services/quote_docx"
	"api/services/quote_template"
	"github.com/gin-gonic/gin"
)

// GenerateGroupPricingQuoteDocx generates a DOCX quotation document for a group pricing quote
func GenerateGroupPricingQuoteDocx(c *gin.Context) {
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

	id := c.Param("id")

	logger.WithField("quote_id", id).Info("Processing GenerateGroupPricingQuoteDocx request")

	// Try to fetch quote and check for active template
	quote, err := services.GetGroupPricingQuote(id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to fetch quote")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch quote"})
		return
	}

	var filename string
	var data []byte

	// Check if insurer has an active template
	// For now, assume insurer ID is always 1 (single insurer setup)
	// In multi-tenant, this would come from quote.InsurerID or similar
	insurer, _ := services.GetInsurerDetails()
	activeTemplate, _ := quote_template.GetActiveTemplate(insurer.ID)

	if activeTemplate != nil {
		// Use template path
		logger.WithField("quote_id", id).Info("Using custom insurer template")

		// Build context
		ctx, err := quote_template.BuildContext(id)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Failed to build template context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build context"})
			return
		}

		// Render template
		data, err = quote_template.Render(activeTemplate.DocxBlob, ctx)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Failed to render template")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to render template"})
			return
		}

		// Generate filename
		sanitizedScheme := sanitizeFilename(quote.SchemeName)
		dateStr := time.Now().Format("2006-01-02")
		filename = fmt.Sprintf("%s_Quotation_%s.docx", sanitizedScheme, dateStr)
	} else {
		// Use default from-scratch generation
		logger.WithField("quote_id", id).Info("Using default quote document generation")

		var err error
		filename, data, err = quote_docx.GenerateQuoteDocx(id)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Failed to generate DOCX quotation")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	logger.WithFields(map[string]interface{}{
		"quote_id": id,
		"filename": filename,
		"size":     len(data),
	}).Info("Successfully generated DOCX quotation")

	// Set response headers for DOCX download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", data)
}

// sanitizeFilename removes/replaces characters invalid in Windows filenames
func sanitizeFilename(s string) string {
	// Replace invalid characters with underscore
	invalidChars := regexp.MustCompile(`[/\\:*?"<>|]`)
	return invalidChars.ReplaceAllString(s, "_")
}

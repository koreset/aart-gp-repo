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
	ctx := requestContext(c)
	logger := log.WithContext(ctx)

	id := c.Param("id")
	logger.WithField("quote_id", id).Info("Processing GenerateGroupPricingQuoteDocx request")

	filename, data, err := buildQuoteDocxBytes(ctx, id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to build DOCX quotation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.WithFields(map[string]interface{}{
		"quote_id": id,
		"filename": filename,
		"size":     len(data),
	}).Info("Successfully generated DOCX quotation")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", data)
}

// buildQuoteDocxBytes resolves the active insurer template (or falls back to the
// from-scratch generator) and returns the rendered .docx bytes plus a download
// filename. Shared by the .docx and .pdf controllers.
func buildQuoteDocxBytes(_ context.Context, id string) (string, []byte, error) {
	quote, err := services.GetGroupPricingQuote(id)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	insurer, _ := services.GetInsurerDetails()
	activeTemplate, _ := quote_template.GetActiveTemplate(insurer.ID)

	if activeTemplate != nil {
		tplCtx, err := quote_template.BuildContext(id)
		if err != nil {
			return "", nil, fmt.Errorf("failed to build template context: %w", err)
		}
		data, err := quote_template.Render(activeTemplate.DocxBlob, tplCtx)
		if err != nil {
			return "", nil, fmt.Errorf("failed to render template: %w", err)
		}
		filename := fmt.Sprintf("%s_Quotation_%s.docx", sanitizeFilename(quote.SchemeName), time.Now().Format("2006-01-02"))
		return filename, data, nil
	}

	return quote_docx.GenerateQuoteDocx(id)
}

// requestContext builds a logging context from the gin request, copying the
// request id and user info attached by upstream middleware.
func requestContext(c *gin.Context) context.Context {
	ctx := context.Background()
	if requestID, ok := c.Get("requestID"); ok {
		ctx = context.WithValue(ctx, log.RequestIDKey, requestID.(string))
	}
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}
	return ctx
}

// sanitizeFilename removes/replaces characters invalid in Windows filenames
func sanitizeFilename(s string) string {
	invalidChars := regexp.MustCompile(`[/\\:*?"<>|]`)
	return invalidChars.ReplaceAllString(s, "_")
}

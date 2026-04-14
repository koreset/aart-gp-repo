package controllers

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"api/log"
	"api/services/quote_template"
	"github.com/gin-gonic/gin"
)

// UploadInsurerQuoteTemplate handles upload of a custom quote template for an insurer
func UploadInsurerQuoteTemplate(c *gin.Context) {
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

	insurerID := c.Param("id")
	insurerIDInt, err := strconv.Atoi(insurerID)
	if err != nil {
		logger.WithField("error", "invalid insurer id").Error("Invalid insurer ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid insurer id"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get uploaded file")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	// Check file extension
	if file.Filename[len(file.Filename)-5:] != ".docx" {
		logger.WithField("filename", file.Filename).Error("Invalid file type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .docx files are supported"})
		return
	}

	// Read file
	src, err := file.Open()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to open uploaded file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer src.Close()

	// Read bytes
	fileBytes := make([]byte, file.Size)
	_, err = src.Read(fileBytes)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to read file bytes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// Validate it's a valid ZIP (DOCX is a ZIP)
	_, err = zip.NewReader(bytes.NewReader(fileBytes), file.Size)
	if err != nil {
		logger.WithField("error", err.Error()).Error("File is not a valid DOCX")
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is not a valid DOCX document"})
		return
	}

	// Save template
	uploadedBy := ""
	if emailExists {
		uploadedBy = userEmail.(string)
	}

	template, err := quote_template.SaveTemplate(insurerIDInt, file.Filename, fileBytes, uploadedBy)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to save template")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save template: %v", err)})
		return
	}

	logger.WithFields(map[string]interface{}{
		"insurer_id": insurerIDInt,
		"version":    template.Version,
		"filename":   template.Filename,
		"size":       template.SizeBytes,
	}).Info("Successfully uploaded quote template")

	c.JSON(http.StatusOK, gin.H{
		"id":          template.ID,
		"insurer_id":  template.InsurerID,
		"version":     template.Version,
		"filename":    template.Filename,
		"size_bytes":  template.SizeBytes,
		"uploaded_by": template.UploadedBy,
		"uploaded_at": template.UploadedAt,
		"is_active":   template.IsActive,
	})
}

// GetActiveInsurerQuoteTemplate returns the active template for an insurer
func GetActiveInsurerQuoteTemplate(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	logger := log.WithContext(ctx)

	insurerID := c.Param("id")
	insurerIDInt, err := strconv.Atoi(insurerID)
	if err != nil {
		logger.WithField("error", "invalid insurer id").Error("Invalid insurer ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid insurer id"})
		return
	}

	template, err := quote_template.GetActiveTemplate(insurerIDInt)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get active template")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch template"})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active template found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          template.ID,
		"insurer_id":  template.InsurerID,
		"version":     template.Version,
		"filename":    template.Filename,
		"size_bytes":  template.SizeBytes,
		"uploaded_by": template.UploadedBy,
		"uploaded_at": template.UploadedAt,
		"is_active":   template.IsActive,
	})
}

// DownloadInsurerQuoteTemplate downloads a specific template by ID
func DownloadInsurerQuoteTemplate(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	logger := log.WithContext(ctx)

	templateID := c.Param("templateId")
	templateIDInt, err := strconv.Atoi(templateID)
	if err != nil {
		logger.WithField("error", "invalid template id").Error("Invalid template ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	template, err := quote_template.GetTemplate(templateIDInt)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get template")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch template"})
		return
	}

	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, template.Filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", template.DocxBlob)
}

// ListInsurerQuoteTemplates returns all template versions for an insurer
func ListInsurerQuoteTemplates(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	logger := log.WithContext(ctx)

	insurerID := c.Param("id")
	insurerIDInt, err := strconv.Atoi(insurerID)
	if err != nil {
		logger.WithField("error", "invalid insurer id").Error("Invalid insurer ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid insurer id"})
		return
	}

	templates, err := quote_template.ListTemplates(insurerIDInt)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to list templates")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// ActivateInsurerQuoteTemplate activates a specific template version
func ActivateInsurerQuoteTemplate(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	logger := log.WithContext(ctx)

	insurerID := c.Param("id")
	insurerIDInt, err := strconv.Atoi(insurerID)
	if err != nil {
		logger.WithField("error", "invalid insurer id").Error("Invalid insurer ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid insurer id"})
		return
	}

	templateID := c.Param("templateId")
	templateIDInt, err := strconv.Atoi(templateID)
	if err != nil {
		logger.WithField("error", "invalid template id").Error("Invalid template ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template id"})
		return
	}

	err = quote_template.ActivateVersion(insurerIDInt, templateIDInt)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to activate template")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to activate template: %v", err)})
		return
	}

	logger.WithFields(map[string]interface{}{
		"insurer_id":  insurerIDInt,
		"template_id": templateIDInt,
	}).Info("Successfully activated template version")

	c.JSON(http.StatusOK, gin.H{"message": "template activated"})
}

// DownloadSampleQuoteTemplate returns a sample template for admins to download and modify
func DownloadSampleQuoteTemplate(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	logger := log.WithContext(ctx)

	sampleData, err := quote_template.BuildSampleTemplate()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to build sample template")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build sample template"})
		return
	}

	logger.Info("Downloaded sample quote template")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", `attachment; filename="sample_quote_template.docx"`)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", sampleData)
}

package controllers

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"api/log"
	"api/services/on_risk_letter_template"
	"github.com/gin-gonic/gin"
)

// UploadInsurerOnRiskLetterTemplate handles upload of a custom on-risk letter template for an insurer
func UploadInsurerOnRiskLetterTemplate(c *gin.Context) {
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

	insurerID := c.Param("id")
	insurerIDInt, err := strconv.Atoi(insurerID)
	if err != nil {
		logger.WithField("error", "invalid insurer id").Error("Invalid insurer ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid insurer id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to get uploaded file")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	if len(file.Filename) < 5 || file.Filename[len(file.Filename)-5:] != ".docx" {
		logger.WithField("filename", file.Filename).Error("Invalid file type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .docx files are supported"})
		return
	}

	src, err := file.Open()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to open uploaded file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer src.Close()

	fileBytes := make([]byte, file.Size)
	_, err = src.Read(fileBytes)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to read file bytes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	_, err = zip.NewReader(bytes.NewReader(fileBytes), file.Size)
	if err != nil {
		logger.WithField("error", err.Error()).Error("File is not a valid DOCX")
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is not a valid DOCX document"})
		return
	}

	uploadedBy := ""
	if emailExists {
		uploadedBy = userEmail.(string)
	}

	template, err := on_risk_letter_template.SaveTemplate(insurerIDInt, file.Filename, fileBytes, uploadedBy)
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
	}).Info("Successfully uploaded on-risk letter template")

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

// GetActiveInsurerOnRiskLetterTemplate returns the active template for an insurer
func GetActiveInsurerOnRiskLetterTemplate(c *gin.Context) {
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

	template, err := on_risk_letter_template.GetActiveTemplate(insurerIDInt)
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

// DownloadInsurerOnRiskLetterTemplate downloads a specific template by ID
func DownloadInsurerOnRiskLetterTemplate(c *gin.Context) {
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

	template, err := on_risk_letter_template.GetTemplate(templateIDInt)
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

// ListInsurerOnRiskLetterTemplates returns all template versions for an insurer
func ListInsurerOnRiskLetterTemplates(c *gin.Context) {
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

	templates, err := on_risk_letter_template.ListTemplates(insurerIDInt)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to list templates")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// ActivateInsurerOnRiskLetterTemplate activates a specific template version
func ActivateInsurerOnRiskLetterTemplate(c *gin.Context) {
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

	err = on_risk_letter_template.ActivateVersion(insurerIDInt, templateIDInt)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to activate template")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to activate template: %v", err)})
		return
	}

	logger.WithFields(map[string]interface{}{
		"insurer_id":  insurerIDInt,
		"template_id": templateIDInt,
	}).Info("Successfully activated on-risk letter template version")

	c.JSON(http.StatusOK, gin.H{"message": "template activated"})
}

// DownloadSampleOnRiskLetterTemplate returns a sample template for admins to download and modify
func DownloadSampleOnRiskLetterTemplate(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	logger := log.WithContext(ctx)

	sampleData, err := on_risk_letter_template.BuildSampleTemplate()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to build sample template")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build sample template"})
		return
	}

	logger.Info("Downloaded sample on-risk letter template")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", `attachment; filename="sample_on_risk_letter_template.docx"`)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", sampleData)
}

// GenerateOnRiskLetterDocxTemplated renders an on-risk letter using the backend template
func GenerateOnRiskLetterDocxTemplated(c *gin.Context) {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}

	logger := log.WithContext(ctx)

	quoteID := c.Param("id")

	filename, data, err := on_risk_letter_template.GenerateOnRiskLetterDocx(quoteID)
	if err != nil {
		if err.Error() == "no active on-risk letter template — upload one in the insurer settings to use the backend-templated path" {
			logger.WithField("quote_id", quoteID).Warn("No active template for on-risk letter generation")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		logger.WithField("error", err.Error()).Error("Failed to generate on-risk letter document")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate document: %v", err)})
		return
	}

	logger.WithFields(map[string]interface{}{
		"quote_id": quoteID,
		"filename": filename,
		"size":     len(data),
	}).Info("Generated templated on-risk letter document")

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", data)
}

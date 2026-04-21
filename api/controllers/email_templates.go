package controllers

import (
	"errors"
	"net/http"
	"strings"

	"api/models"
	"api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListEmailTemplates lists all templates for the caller's license.
func ListEmailTemplates(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	if licenseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-License-Id header is required"})
		return
	}
	list, err := services.ListEmailTemplates(licenseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetEmailTemplate returns one template by code.
func GetEmailTemplate(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	code := c.Param("code")
	tpl, err := services.GetEmailTemplate(licenseId, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tpl)
}

// CreateEmailTemplate creates a new template.
func CreateEmailTemplate(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	user := c.MustGet("user").(models.AppUser)
	var in services.SaveEmailTemplateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tpl, err := services.CreateEmailTemplate(licenseId, in, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tpl)
}

// UpdateEmailTemplate updates a template in place.
func UpdateEmailTemplate(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	code := c.Param("code")
	user := c.MustGet("user").(models.AppUser)
	var in services.SaveEmailTemplateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tpl, err := services.UpdateEmailTemplate(licenseId, code, in, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tpl)
}

// DeleteEmailTemplate removes a template.
func DeleteEmailTemplate(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	code := c.Param("code")
	if err := services.DeleteEmailTemplate(licenseId, code); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// PreviewEmailTemplateRequest holds the preview body: sample variables and,
// optionally, an inline template override so admins can preview edits before
// saving them.
type PreviewEmailTemplateRequest struct {
	Vars            map[string]interface{} `json:"vars"`
	SubjectTemplate string                 `json:"subject_template"`
	BodyTemplate    string                 `json:"body_template"`
}

// PreviewEmailTemplate renders the named template with supplied vars without
// persisting or sending anything. If the request includes subject_template or
// body_template overrides they take precedence (for preview-while-editing).
func PreviewEmailTemplate(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	code := c.Param("code")

	var req PreviewEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tpl, err := services.GetEmailTemplate(licenseId, code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Fall back to a stub with the inline templates when the row doesn't
	// exist yet — lets an admin preview before the first save.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tpl = models.EmailTemplate{Code: code, LicenseId: licenseId}
	}
	if req.SubjectTemplate != "" {
		tpl.SubjectTemplate = req.SubjectTemplate
	}
	if req.BodyTemplate != "" {
		tpl.BodyTemplate = req.BodyTemplate
	}

	out, err := services.PreviewEmailTemplate(tpl, req.Vars)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

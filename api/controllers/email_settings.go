package controllers

import (
	"errors"
	"net/http"
	"strings"

	"api/models"
	"api/services"
	"api/services/email"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// emailSettingsResponse is the shape returned to clients. The encrypted
// password is never surfaced; the client gets a has_password boolean.
type emailSettingsResponse struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	TlsMode      string `json:"tls_mode"`
	AuthUser     string `json:"auth_user"`
	HasPassword  bool   `json:"has_password"`
	FromAddress  string `json:"from_address"`
	FromName     string `json:"from_name"`
	ReplyTo      string `json:"reply_to"`
	UpdatedBy    string `json:"updated_by"`
}

func toEmailSettingsResponse(s models.EmailSettings) emailSettingsResponse {
	return emailSettingsResponse{
		Host:        s.Host,
		Port:        s.Port,
		TlsMode:     s.TlsMode,
		AuthUser:    s.AuthUser,
		HasPassword: s.AuthPasswordEncrypted != "",
		FromAddress: s.FromAddress,
		FromName:    s.FromName,
		ReplyTo:     s.ReplyTo,
		UpdatedBy:   s.UpdatedBy,
	}
}

// GetEmailSettings returns the SMTP configuration for the caller's license.
// Responds 404 when no settings row has been created yet.
func GetEmailSettings(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	if licenseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-License-Id header is required"})
		return
	}
	s, err := services.GetEmailSettings(licenseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "email settings not configured"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toEmailSettingsResponse(s))
}

// SaveEmailSettings upserts the per-license SMTP configuration.
func SaveEmailSettings(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	if licenseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-License-Id header is required"})
		return
	}
	user := c.MustGet("user").(models.AppUser)

	var in services.SaveEmailSettingsInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	saved, err := services.SaveEmailSettings(licenseId, in, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toEmailSettingsResponse(saved))
}

// SendTestEmail enqueues a system test email to the calling user. Requires a
// template with code "system_test" to exist and be active.
func SendTestEmail(c *gin.Context) {
	licenseId := strings.TrimSpace(c.GetHeader("X-License-Id"))
	if licenseId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-License-Id header is required"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if user.UserEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current user has no email on file"})
		return
	}
	row, err := services.EnqueueEmail(services.EnqueueEmailRequest{
		LicenseId:    licenseId,
		TemplateCode: "system_test",
		To:           []string{user.UserEmail},
		Vars: map[string]interface{}{
			"user_name": user.UserName,
			"user_email": user.UserEmail,
		},
		Attachments:       []email.AttachmentSpec{},
		RelatedObjectType: "email_settings",
		RelatedObjectID:   licenseId,
		CreatedBy:         user.UserName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"outbox_id": row.ID})
}

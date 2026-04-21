package services

import (
	"encoding/json"
	"fmt"
	"time"

	appLog "api/log"
	"api/models"
	"api/services/email"
)

// EnqueueEmailRequest is the input to EnqueueEmail. Call sites populate it
// with the addressees, the template code, and the variables to render the
// template against. The worker's first poll after enqueue picks it up.
type EnqueueEmailRequest struct {
	LicenseId         string
	TemplateCode      string
	To                []string
	Cc                []string
	Bcc               []string
	Vars              map[string]interface{}
	Attachments       []email.AttachmentSpec
	// Optional overrides; default to values on EmailSettings / user's GPUser.
	FromAddress       string
	FromName          string
	ReplyTo           string
	// Optional traceability.
	RelatedObjectType string
	RelatedObjectID   string
	CreatedBy         string
	// ScheduledAt defaults to now.
	ScheduledAt       time.Time
}

// EnqueueEmail renders the named template and writes a pending row to the
// outbox. Returns the created outbox row.
func EnqueueEmail(req EnqueueEmailRequest) (models.EmailOutbox, error) {
	logger := appLog.WithField("action", "EnqueueEmail")

	if req.LicenseId == "" {
		return models.EmailOutbox{}, fmt.Errorf("license_id is required")
	}
	if req.TemplateCode == "" {
		return models.EmailOutbox{}, fmt.Errorf("template_code is required")
	}
	if len(req.To) == 0 {
		return models.EmailOutbox{}, fmt.Errorf("at least one To recipient is required")
	}

	var tpl models.EmailTemplate
	if err := DB.Where("license_id = ? AND code = ? AND status = ?", req.LicenseId, req.TemplateCode, models.EmailTemplateActive).
		First(&tpl).Error; err != nil {
		return models.EmailOutbox{}, fmt.Errorf("active template %q not found: %w", req.TemplateCode, err)
	}

	rendered, err := email.Render(tpl, req.Vars)
	if err != nil {
		return models.EmailOutbox{}, err
	}

	toJSON, _ := json.Marshal(req.To)
	ccJSON, _ := json.Marshal(req.Cc)
	bccJSON, _ := json.Marshal(req.Bcc)
	var attachmentsJSON []byte
	if len(req.Attachments) > 0 {
		attachmentsJSON, err = json.Marshal(req.Attachments)
		if err != nil {
			return models.EmailOutbox{}, fmt.Errorf("marshal attachments: %w", err)
		}
	}

	scheduled := req.ScheduledAt
	if scheduled.IsZero() {
		scheduled = time.Now()
	}

	row := models.EmailOutbox{
		LicenseId:         req.LicenseId,
		TemplateCode:      req.TemplateCode,
		FromAddress:       req.FromAddress,
		FromName:          req.FromName,
		ReplyTo:           req.ReplyTo,
		ToRecipients:      string(toJSON),
		CcRecipients:      string(ccJSON),
		BccRecipients:     string(bccJSON),
		Subject:           rendered.Subject,
		Body:              rendered.Body,
		Attachments:       string(attachmentsJSON),
		Status:            models.EmailOutboxPending,
		Attempts:          0,
		MaxAttempts:       5,
		NextAttemptAt:     scheduled,
		ScheduledAt:       scheduled,
		RelatedObjectType: req.RelatedObjectType,
		RelatedObjectID:   req.RelatedObjectID,
		CreatedBy:         req.CreatedBy,
		CreatedAt:         time.Now(),
	}
	if err := DB.Create(&row).Error; err != nil {
		return models.EmailOutbox{}, err
	}
	logger.WithField("outbox_id", row.ID).Info("Email enqueued")
	return row, nil
}

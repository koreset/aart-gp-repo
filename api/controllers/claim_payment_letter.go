package controllers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"api/models"
	"api/services"
	"api/services/docpdf"
	"api/services/email"
	"api/services/payment_letter"

	"github.com/gin-gonic/gin"
)

// ──────────────────────────────────────────────
// Letter generation (per claim)
// ──────────────────────────────────────────────

// DownloadClaimPaymentLetterDocx handles GET /group-pricing/claims/:id/payment-letter.docx
func DownloadClaimPaymentLetterDocx(c *gin.Context) {
	downloadClaimPaymentLetter(c, models.PaymentLetterFormatDocx)
}

// DownloadClaimPaymentLetterPdf handles GET /group-pricing/claims/:id/payment-letter.pdf
func DownloadClaimPaymentLetterPdf(c *gin.Context) {
	downloadClaimPaymentLetter(c, models.PaymentLetterFormatPDF)
}

func downloadClaimPaymentLetter(c *gin.Context, format string) {
	claimID, err := strconv.Atoi(c.Param("claim_id"))
	if err != nil {
		BadRequestMsg(c, "invalid claim id")
		return
	}
	user := c.MustGet("user").(models.AppUser)

	_, filename, data, err := payment_letter.GenerateAndRecord(services.DB, claimID, format, user.UserName)
	if err != nil {
		switch {
		case errors.Is(err, payment_letter.ErrClaimNotPaid):
			c.JSON(http.StatusConflict, gin.H{"error": "letter is only available after the claim is paid"})
		case errors.Is(err, docpdf.ErrConverterNotFound):
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "PDF conversion is unavailable on this server. Install LibreOffice (https://www.libreoffice.org) or Microsoft Word so the .docx can be rendered to PDF.",
			})
		default:
			InternalError(c, err)
		}
		return
	}

	contentType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	if format == models.PaymentLetterFormatPDF {
		contentType = "application/pdf"
	}
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, contentType, data)
}

// ListClaimPaymentLetters handles GET /group-pricing/claims/:id/payment-letter/history
func ListClaimPaymentLetters(c *gin.Context) {
	claimID, err := strconv.Atoi(c.Param("claim_id"))
	if err != nil {
		BadRequestMsg(c, "invalid claim id")
		return
	}
	rows, err := payment_letter.ListHistory(services.DB, claimID)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, rows)
}

// ──────────────────────────────────────────────
// Send / dispatch
// ──────────────────────────────────────────────

type sendLetterRequest struct {
	LetterID  int    `json:"letter_id"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
}

// SendClaimPaymentLetter handles POST /group-pricing/claims/:id/payment-letter/send
func SendClaimPaymentLetter(c *gin.Context) {
	claimID, err := strconv.Atoi(c.Param("claim_id"))
	if err != nil {
		BadRequestMsg(c, "invalid claim id")
		return
	}
	var req sendLetterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err)
		return
	}
	if req.LetterID == 0 {
		BadRequestMsg(c, "letter_id is required")
		return
	}
	if req.Channel == "" {
		BadRequestMsg(c, "channel is required")
		return
	}
	if err := payment_letter.CheckChannelImplemented(req.Channel); err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)

	var letter models.ClaimPaymentLetter
	if err := services.DB.First(&letter, req.LetterID).Error; err != nil {
		NotFound(c, "letter not found")
		return
	}
	if letter.ClaimID != claimID {
		BadRequestMsg(c, "letter does not belong to this claim")
		return
	}
	var claim models.GroupSchemeClaim
	if err := services.DB.First(&claim, claimID).Error; err != nil {
		NotFound(c, "claim not found")
		return
	}

	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		recipient = payment_letter.DefaultRecipient(claim, req.Channel)
	}
	if recipient == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": fmt.Sprintf("no %s recipient on file for this claim — please update the claim record first", req.Channel),
		})
		return
	}

	switch req.Channel {
	case models.PaymentLetterChannelEmail:
		path, err := payment_letter.PrepareLetterAttachment(services.DB, letter, claim)
		if err != nil {
			InternalError(c, err)
			return
		}
		outbox, err := services.EnqueueEmail(services.EnqueueEmailRequest{
			LicenseId:    resolveDefaultLicenseForLetters(),
			TemplateCode: "claim_payment_letter",
			To:           []string{recipient},
			Vars: map[string]interface{}{
				"claimant_name":  claim.ClaimantName,
				"claim_number":   claim.ClaimNumber,
				"payment_amount": letter.PaymentAmount,
				"paid_at":        letter.PaidAt.Format("02 January 2006"),
				"scheme_name":    claim.SchemeName,
				"benefit_name":   claim.BenefitName,
				"letter_ref":     letter.LetterReference,
			},
			Attachments: []email.AttachmentSpec{{
				Kind:     "file",
				Path:     path,
				Filename: filepath.Base(path),
			}},
			RelatedObjectType: "claim_payment_letter",
			RelatedObjectID:   strconv.Itoa(letter.ID),
			CreatedBy:         user.UserName,
		})
		if err != nil {
			// Persist the failure attempt so the user can see it in history.
			_, _ = payment_letter.RecordDelivery(services.DB, models.ClaimPaymentLetterDelivery{
				LetterID:  letter.ID,
				Channel:   req.Channel,
				Recipient: recipient,
				Status:    models.PaymentLetterDeliveryFailed,
				Error:     err.Error(),
				SentBy:    user.UserName,
			})
			BadRequest(c, err)
			return
		}
		oid := outbox.ID
		delivery, derr := payment_letter.RecordDelivery(services.DB, models.ClaimPaymentLetterDelivery{
			LetterID:  letter.ID,
			Channel:   req.Channel,
			Recipient: recipient,
			Status:    models.PaymentLetterDeliveryPending,
			OutboxID:  &oid,
			SentBy:    user.UserName,
		})
		if derr != nil {
			InternalError(c, derr)
			return
		}
		payment_letter.LogCommunication(services.DB, claim.ID, "Email",
			fmt.Sprintf("Payment confirmation letter %s emailed to %s", letter.LetterReference, recipient),
			user.UserName)
		Created(c, delivery)

	case models.PaymentLetterChannelManual:
		delivery, derr := payment_letter.RecordDelivery(services.DB, models.ClaimPaymentLetterDelivery{
			LetterID:  letter.ID,
			Channel:   req.Channel,
			Recipient: recipient,
			Status:    models.PaymentLetterDeliverySent,
			SentBy:    user.UserName,
		})
		if derr != nil {
			InternalError(c, derr)
			return
		}
		payment_letter.LogCommunication(services.DB, claim.ID, "Letter",
			fmt.Sprintf("Payment confirmation letter %s marked as delivered (%s)", letter.LetterReference, recipient),
			user.UserName)
		Created(c, delivery)

	case models.PaymentLetterChannelSMS:
		// SMS path: scaffold exists but the only registered provider is the
		// log-only dev one. Record the attempt as pending so future
		// real-provider rollout can pick it up.
		delivery, derr := payment_letter.RecordDelivery(services.DB, models.ClaimPaymentLetterDelivery{
			LetterID:  letter.ID,
			Channel:   req.Channel,
			Recipient: recipient,
			Status:    models.PaymentLetterDeliveryPending,
			Error:     "SMS provider not configured; letter URL pending future implementation",
			SentBy:    user.UserName,
		})
		if derr != nil {
			InternalError(c, derr)
			return
		}
		c.JSON(http.StatusAccepted, models.PremiumResponse{Success: true, Data: delivery})

	default:
		c.JSON(http.StatusNotImplemented, gin.H{"error": "channel not supported"})
	}
}

// resolveDefaultLicenseForLetters mirrors the single-tenant license resolver
// in claimant_notifications.go: returns the first EmailSettings row's
// license_id. Empty if email isn't configured (handled by EnqueueEmail).
func resolveDefaultLicenseForLetters() string {
	type r struct{ LicenseId string }
	var row r
	_ = services.DB.Raw(`SELECT license_id FROM email_settings ORDER BY id ASC LIMIT 1`).Scan(&row).Error
	return row.LicenseId
}

// ──────────────────────────────────────────────
// Schedule bulk bundle
// ──────────────────────────────────────────────

// DownloadScheduleLetterBundle handles
// GET /group-pricing/claims/payment-schedules/:schedule_id/payment-letters.zip?format=docx|pdf
func DownloadScheduleLetterBundle(c *gin.Context) {
	scheduleID, err := strconv.Atoi(c.Param("schedule_id"))
	if err != nil {
		BadRequestMsg(c, "invalid schedule_id")
		return
	}
	format := c.DefaultQuery("format", models.PaymentLetterFormatDocx)
	user := c.MustGet("user").(models.AppUser)

	filename, data, err := payment_letter.GenerateScheduleBundle(services.DB, scheduleID, format, user.UserName)
	if err != nil {
		if errors.Is(err, docpdf.ErrConverterNotFound) {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "PDF conversion is unavailable on this server. Install LibreOffice or Microsoft Word so the .docx can be rendered to PDF.",
			})
			return
		}
		BadRequest(c, err)
		return
	}
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/zip", data)
}

// ──────────────────────────────────────────────
// Settings (singleton letterhead/signatory)
// ──────────────────────────────────────────────

// GetPaymentLetterSettings handles GET /group-pricing/settings/payment-letter
func GetPaymentLetterSettings(c *gin.Context) {
	s, err := payment_letter.LoadSettings(services.DB)
	if err != nil {
		InternalError(c, err)
		return
	}
	// Indicate presence of binary fields without sending the bytes back.
	OK(c, gin.H{
		"id":                      s.ID,
		"company_name":            s.CompanyName,
		"address_line1":           s.AddressLine1,
		"address_line2":           s.AddressLine2,
		"address_line3":           s.AddressLine3,
		"city":                    s.City,
		"postal_code":             s.PostalCode,
		"country":                 s.Country,
		"phone":                   s.Phone,
		"email":                   s.Email,
		"website":                 s.Website,
		"signatory_name":          s.SignatoryName,
		"signatory_title":         s.SignatoryTitle,
		"letter_intro_template":   s.LetterIntroTemplate,
		"letter_closing_template": s.LetterClosingTemplate,
		"has_logo":                len(s.Logo) > 0,
		"logo_mime_type":          s.LogoMimeType,
		"has_signature":           len(s.Signature) > 0,
		"signature_mime_type":     s.SignatureMimeType,
		"updated_at":              s.UpdatedAt,
		"updated_by":              s.UpdatedBy,
	})
}

const maxLogoBytes = 2 * 1024 * 1024

var allowedLogoMimes = map[string]bool{
	"image/png":     true,
	"image/jpeg":    true,
	"image/svg+xml": true,
}

// UpdatePaymentLetterSettings handles PUT /group-pricing/settings/payment-letter
// Accepts multipart/form-data so logo and signature images can be uploaded
// alongside the text fields. Non-image fields are read from form values.
func UpdatePaymentLetterSettings(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(maxLogoBytes * 2); err != nil {
		BadRequest(c, fmt.Errorf("parse form: %w", err))
		return
	}
	user := c.MustGet("user").(models.AppUser)

	in := payment_letter.UpsertSettingsInput{
		CompanyName:    c.PostForm("company_name"),
		AddressLine1:   c.PostForm("address_line1"),
		AddressLine2:   c.PostForm("address_line2"),
		AddressLine3:   c.PostForm("address_line3"),
		City:           c.PostForm("city"),
		PostalCode:     c.PostForm("postal_code"),
		Country:        c.PostForm("country"),
		Phone:          c.PostForm("phone"),
		Email:          c.PostForm("email"),
		Website:        c.PostForm("website"),
		SignatoryName:  c.PostForm("signatory_name"),
		SignatoryTitle: c.PostForm("signatory_title"),
		LetterIntro:    c.PostForm("letter_intro_template"),
		LetterClosing:  c.PostForm("letter_closing_template"),
		UpdatedBy:      user.UserName,
	}

	if fh, err := c.FormFile("logo"); err == nil {
		bytes, mime, ferr := readImageUploadFile(fh)
		if ferr != nil {
			BadRequest(c, ferr)
			return
		}
		in.LogoProvided = true
		in.Logo = bytes
		in.LogoMimeType = mime
	}
	if fh, err := c.FormFile("signature"); err == nil {
		bytes, mime, ferr := readImageUploadFile(fh)
		if ferr != nil {
			BadRequest(c, ferr)
			return
		}
		in.SignatureProvided = true
		in.Signature = bytes
		in.SignatureMimeType = mime
	}

	saved, err := payment_letter.UpsertSettings(services.DB, in)
	if err != nil {
		InternalError(c, err)
		return
	}
	OK(c, gin.H{
		"id":                      saved.ID,
		"updated_at":              saved.UpdatedAt,
		"updated_by":              saved.UpdatedBy,
		"has_logo":                len(saved.Logo) > 0,
		"logo_mime_type":          saved.LogoMimeType,
		"has_signature":           len(saved.Signature) > 0,
		"signature_mime_type":     saved.SignatureMimeType,
	})
}

// readImageUploadFile validates the uploaded image (MIME + size cap), reads
// the bytes, and returns them along with the resolved MIME type.
func readImageUploadFile(fh *multipart.FileHeader) ([]byte, string, error) {
	if fh.Size > maxLogoBytes {
		return nil, "", fmt.Errorf("uploaded image exceeds the 2 MB size limit")
	}
	mime := fh.Header.Get("Content-Type")
	if mime == "" {
		// Fall back to extension-based guessing
		switch strings.ToLower(filepath.Ext(fh.Filename)) {
		case ".png":
			mime = "image/png"
		case ".jpg", ".jpeg":
			mime = "image/jpeg"
		case ".svg":
			mime = "image/svg+xml"
		}
	}
	if !allowedLogoMimes[mime] {
		return nil, "", fmt.Errorf("unsupported image type %q (allowed: png, jpeg, svg)", mime)
	}
	f, err := fh.Open()
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, "", err
	}
	return data, mime, nil
}

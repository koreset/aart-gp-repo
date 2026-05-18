package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"api/models"
	"api/services"
	uwvendor "api/services/uw_vendor"
	uwproviders "api/services/uw_vendor/providers"
)

// SubmitVendorRequest dispatches a vendor request via the registry. The
// kind is taken from the URL (`/.../request/:kind`) so the renderer's
// "Request pathology" etc. buttons each have a stable endpoint.
//
// Body:
//
//	{
//	  "case_id": 123,
//	  "quote_id": 456,
//	  "subject": "+27821234567" | "broker@example.com" | "MEMBER-789",
//	  "body": "...",
//	  "metadata": { ... }
//	}
func SubmitVendorRequest(c *gin.Context) {
	kind := uwvendor.Kind(c.Param("kind"))
	if !uwvendor.IsValidKind(kind) {
		c.JSON(http.StatusBadRequest, "unknown vendor kind")
		return
	}
	var payload struct {
		CaseID   int            `json:"case_id"`
		QuoteID  int            `json:"quote_id"`
		Subject  string         `json:"subject"`
		Body     string         `json:"body"`
		Metadata map[string]any `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := c.MustGet("user").(models.AppUser)
	req := uwvendor.SubmitRequest{
		Kind:     kind,
		CaseID:   payload.CaseID,
		QuoteID:  payload.QuoteID,
		Subject:  payload.Subject,
		Body:     payload.Body,
		Metadata: payload.Metadata,
	}
	row, err := services.VendorRegistry().Submit(c.Request.Context(), req, user.UserEmail)
	if err != nil {
		// Surface the persisted row even on error so the UI can show what
		// was attempted; the status will be `failed` with an error message.
		if errors.Is(err, uwvendor.ErrProviderNotConfigured) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error(), "request": row})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "request": row})
		return
	}
	c.JSON(http.StatusCreated, row)
}

// ListVendorRequestsForCase returns every VendorRequest associated with a
// case, newest first. Includes failed and cancelled requests so the case
// timeline is complete.
func ListVendorRequestsForCase(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("case_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid case_id")
		return
	}
	var rows []models.VendorRequest
	if err := services.DB.Where("case_id = ?", caseID).Order("requested_at DESC").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, rows)
}

// IngestVendorWebhook is the gateway every vendor calls back into. Path
// shape: POST /vendor-webhooks/:kind/:provider. The Registry verifies the
// signature, persists the raw body, dispatches to the provider's
// HandleWebhook, and applies the outcome to the originating VendorRequest.
//
// Status semantics:
//   - 200 OK with `{"processed": true}` when the webhook was applied.
//   - 200 OK with `{"replayed": true}` when the same webhook arrived
//     twice — idempotency dedup short-circuited.
//   - 401 Unauthorized when the signature header is missing or invalid.
//   - 422 Unprocessable Entity when the body decoded but referenced an
//     unknown external request id.
//   - 500 for anything else.
func IngestVendorWebhook(c *gin.Context) {
	kind := uwvendor.Kind(c.Param("kind"))
	provider := c.Param("provider")
	if !uwvendor.IsValidKind(kind) || provider == "" {
		c.JSON(http.StatusBadRequest, "invalid kind/provider")
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	row, ingestErr := services.VendorRegistry().IngestWebhook(c.Request.Context(), kind, provider, c.Request.Header, body)
	switch {
	case errors.Is(ingestErr, uwvendor.ErrAlreadyProcessed):
		c.JSON(http.StatusOK, gin.H{"replayed": true, "webhook": row})
	case errors.Is(ingestErr, uwvendor.ErrInvalidSignature):
		c.JSON(http.StatusUnauthorized, gin.H{"error": ingestErr.Error(), "webhook": row})
	case errors.Is(ingestErr, uwvendor.ErrUnknownExternalID):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": ingestErr.Error(), "webhook": row})
	case errors.Is(ingestErr, uwvendor.ErrProviderNotConfigured):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": ingestErr.Error()})
	case ingestErr != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ingestErr.Error(), "webhook": row})
	default:
		c.JSON(http.StatusOK, gin.H{"processed": true, "webhook": row})
	}
}

// FireMockVendorWebhook is a dev / test helper: given the ID of a
// VendorRequest awaiting response, it builds the signed mock webhook
// envelope and pushes it through the ingest pipeline. Exposed only when
// the active provider for the kind is the bundled Mock — the registry
// type-asserts here.
func FireMockVendorWebhook(c *gin.Context) {
	requestID, err := strconv.Atoi(c.Param("request_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request_id")
		return
	}
	var req models.VendorRequest
	if err := services.DB.Where("id = ?", requestID).First(&req).Error; err != nil {
		c.JSON(http.StatusNotFound, "request not found")
		return
	}
	provider := services.VendorRegistry().Active(uwvendor.Kind(req.Kind))
	mock, ok := provider.(*uwproviders.Mock)
	if !ok {
		c.JSON(http.StatusBadRequest, "active provider is not the mock; nothing to simulate")
		return
	}
	var payload struct {
		Status   string `json:"status"`
		Filename string `json:"filename"`
		BodyText string `json:"body_text"`
		AttachKind string `json:"attach_kind"`
	}
	_ = c.ShouldBindJSON(&payload)
	if payload.Status == "" {
		payload.Status = string(uwvendor.StatusComplete)
	}
	var att *uwvendor.AttachmentPayload
	if payload.Filename != "" || payload.BodyText != "" {
		att = &uwvendor.AttachmentPayload{
			Kind:        firstNonEmptyStr(payload.AttachKind, models.UWAttachmentKindMedicalReport),
			FileName:    firstNonEmptyStr(payload.Filename, "mock-result.txt"),
			ContentType: "text/plain",
			Body:        []byte(firstNonEmptyStr(payload.BodyText, "Mock delivery body")),
		}
	}
	headers, body := mock.SimulateDelivery(req.ExternalRequestID, uwvendor.Status(payload.Status), att)
	row, ingestErr := services.VendorRegistry().IngestWebhook(c.Request.Context(), uwvendor.Kind(req.Kind), req.Provider, headers, body)
	if ingestErr != nil && !errors.Is(ingestErr, uwvendor.ErrAlreadyProcessed) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ingestErr.Error(), "webhook": row})
		return
	}
	c.JSON(http.StatusOK, gin.H{"webhook": row, "replayed": errors.Is(ingestErr, uwvendor.ErrAlreadyProcessed)})
}

func firstNonEmptyStr(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}

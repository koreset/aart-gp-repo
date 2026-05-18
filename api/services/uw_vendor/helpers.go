package uwvendor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"

	"api/models"
)

// vendorAttachmentDir is the on-disk root for files delivered by vendor
// webhooks. Mirrors the case-attachment layout from underwriting_case.go
// so the existing download controller can stream them.
const vendorAttachmentDir = "tmp/uploads/underwriting_cases"

// firstHeader returns the first value for an HTTP header in case-insensitive
// fashion. Used to pull signature / event-type headers out of incoming
// webhook requests.
func firstHeader(h map[string][]string, name string) string {
	if h == nil {
		return ""
	}
	canon := http.CanonicalHeaderKey(name)
	if vs, ok := h[canon]; ok && len(vs) > 0 {
		return vs[0]
	}
	// Fall back to a case-insensitive scan for non-canonical headers.
	for k, vs := range h {
		if strings.EqualFold(k, name) && len(vs) > 0 {
			return vs[0]
		}
	}
	return ""
}

// actorFromHeaders returns the actor email a vendor webhook claims to be
// from (rare but some providers include it). Falls back to "vendor_webhook"
// so audit events always have a non-empty actor.
func actorFromHeaders(h map[string][]string) string {
	if v := firstHeader(h, "X-Actor-Email"); v != "" {
		return v
	}
	return "vendor_webhook"
}

// sha256Hex returns the hex-encoded SHA-256 of body.
func sha256Hex(body []byte) string {
	h := sha256.Sum256(body)
	return hex.EncodeToString(h[:])
}

// sniffExternalID inspects a webhook body for a top-level
// `external_request_id` or `request_id` field. Best-effort: an empty
// result is fine because the webhook is still persisted for replay.
func sniffExternalID(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var probe map[string]any
	if err := json.Unmarshal(body, &probe); err != nil {
		return ""
	}
	for _, key := range []string{"external_request_id", "request_id", "externalRequestId", "requestId"} {
		if v, ok := probe[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// persistAttachment writes the vendor-delivered file to disk and inserts
// a UnderwritingCaseAttachment row pointing at it. Re-uses the existing
// case-attachment layout so the download controller streams it
// transparently.
//
// We intentionally don't go through services.AppendCaseAttachments
// because that helper expects a *multipart.FileHeader; here we already
// have raw bytes from the webhook payload.
func persistAttachment(db *gorm.DB, caseID int, att *AttachmentPayload, actor string) error {
	if att == nil || len(att.Body) == 0 {
		return nil
	}
	baseDir := filepath.Join(vendorAttachmentDir, fmt.Sprintf("case_%d", caseID))
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return err
	}
	name := att.FileName
	if name == "" {
		name = fmt.Sprintf("vendor-%d.bin", time.Now().UnixNano())
	}
	destPath := filepath.Join(baseDir, filepath.Base(name))
	if err := os.WriteFile(destPath, att.Body, 0o644); err != nil {
		return err
	}
	row := models.UnderwritingCaseAttachment{
		CaseID:      caseID,
		Kind:        firstNonEmpty(att.Kind, models.UWAttachmentKindMedicalReport),
		FileName:    name,
		ContentType: att.ContentType,
		SizeBytes:   int64(len(att.Body)),
		StoragePath: destPath,
		UploadedBy:  actor,
	}
	return db.Create(&row).Error
}

func firstNonEmpty(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}

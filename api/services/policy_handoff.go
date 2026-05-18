package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"api/models"
)

// PolicyHandoffPayload is the canonical JSON document persisted on
// PolicyHandoffSnapshot.Payload. It snapshots everything compliance needs
// to reproduce the in-force handoff: the quote and scheme as they stood
// at promotion time, the in-force member list, every underwriting case
// decision, prior-insurer schedule references (if takeover applied), and
// an index of attached evidence.
//
// Phase 6 vendor requests and Phase 5 disclosures are referenced by ID
// only; the raw bodies live in their own tables / on-disk locations and
// the snapshot links them so an auditor can pull them on demand.
type PolicyHandoffPayload struct {
	SchemaVersion int                          `json:"schema_version"`
	HandedOffAt   time.Time                    `json:"handed_off_at"`
	HandedOffBy   string                       `json:"handed_off_by"`
	Reason        string                       `json:"reason"`
	Quote         models.GroupPricingQuote     `json:"quote"`
	Scheme        models.GroupScheme           `json:"scheme"`
	Categories    []models.SchemeCategory      `json:"categories"`
	Members       []models.GPricingMemberDataInForce `json:"members"`
	Cases         []HandoffCaseEntry           `json:"cases"`
	PriorSchedule *models.PriorInsurerSchedule `json:"prior_schedule,omitempty"`
	Stats         *models.GroupRiskQuoteStats  `json:"stats,omitempty"`
}

// HandoffCaseEntry is a denormalised case-and-decisions row for the
// payload. We embed decisions inline so reading the snapshot doesn't
// require joining four tables.
type HandoffCaseEntry struct {
	Case        models.UnderwritingCase        `json:"case"`
	Decisions   []models.UnderwritingDecision  `json:"decisions"`
	Attachments []HandoffAttachmentEntry       `json:"attachments,omitempty"`
}

// HandoffAttachmentEntry is a thin pointer at an attachment row. The body
// is not embedded — auditors fetch it via the existing download endpoint
// using `attachment_id`.
type HandoffAttachmentEntry struct {
	AttachmentID int    `json:"attachment_id"`
	Kind         string `json:"kind"`
	FileName     string `json:"file_name"`
	ContentType  string `json:"content_type"`
	SizeBytes    int64  `json:"size_bytes"`
	UploadedAt   string `json:"uploaded_at"`
}

// PolicyHandoffSchemaVersion is bumped whenever the PolicyHandoffPayload
// structure changes. Auditors consuming older snapshots can branch on it.
const PolicyHandoffSchemaVersion = 1

// BuildHandoffSnapshot reads the in-force tables and assembles the
// canonical payload. Pure read — no side effects, safe to call from a
// renderer for preview.
func BuildHandoffSnapshot(quoteID, schemeID int, reason, actor string) (*PolicyHandoffPayload, error) {
	if quoteID <= 0 {
		return nil, errors.New("invalid quote_id")
	}
	payload := &PolicyHandoffPayload{
		SchemaVersion: PolicyHandoffSchemaVersion,
		HandedOffAt:   time.Now(),
		HandedOffBy:   actor,
		Reason:        reason,
	}
	if err := DB.Preload("SchemeCategories").Where("id = ?", quoteID).First(&payload.Quote).Error; err != nil {
		return nil, fmt.Errorf("load quote: %w", err)
	}
	if schemeID == 0 {
		schemeID = payload.Quote.SchemeID
	}
	if err := DB.Where("id = ?", schemeID).First(&payload.Scheme).Error; err != nil {
		return nil, fmt.Errorf("load scheme: %w", err)
	}
	if err := DB.Where("quote_id = ?", quoteID).Find(&payload.Categories).Error; err != nil {
		return nil, fmt.Errorf("load categories: %w", err)
	}
	if err := DB.Where("quote_id = ?", quoteID).Find(&payload.Members).Error; err != nil {
		return nil, fmt.Errorf("load in-force members: %w", err)
	}

	var cases []models.UnderwritingCase
	if err := DB.Preload("Decisions").Preload("Attachments").Where("quote_id = ?", quoteID).Find(&cases).Error; err != nil {
		return nil, fmt.Errorf("load cases: %w", err)
	}
	for _, c := range cases {
		entry := HandoffCaseEntry{Case: c, Decisions: c.Decisions}
		for _, a := range c.Attachments {
			entry.Attachments = append(entry.Attachments, HandoffAttachmentEntry{
				AttachmentID: a.ID,
				Kind:         a.Kind,
				FileName:     a.FileName,
				ContentType:  a.ContentType,
				SizeBytes:    a.SizeBytes,
				UploadedAt:   a.UploadedAt.Format(time.RFC3339),
			})
		}
		payload.Cases = append(payload.Cases, entry)
	}

	// Optional: most recent prior-insurer schedule, if any.
	if prior, err := GetPriorInsurerScheduleForQuote(quoteID); err == nil && prior != nil {
		payload.PriorSchedule = prior
	}

	// Optional: quote stats (annual premium, expected claims, etc.) for
	// reconciliation against the finance module.
	var stats models.GroupRiskQuoteStats
	if err := DB.Where("quote_id = ?", quoteID).First(&stats).Error; err == nil {
		payload.Stats = &stats
	}

	return payload, nil
}

// RecordPolicyHandoffSnapshot assembles the canonical payload and writes
// a PolicyHandoffSnapshot row. Idempotency: if a snapshot already exists
// for (QuoteID, SchemeID), returns the existing row unchanged. Callers in
// the scheduler invoke this from the same transaction that flips the
// scheme to InForce — concurrent ticks therefore serialise on the row
// lock.
//
// Pass tx to use the current transaction (recommended when called from
// the scheduler); pass nil to use the package-level DB.
func RecordPolicyHandoffSnapshot(tx *gorm.DB, quoteID, schemeID int, actor, reason string) (*models.PolicyHandoffSnapshot, error) {
	db := tx
	if db == nil {
		db = DB
	}
	var existing models.PolicyHandoffSnapshot
	err := db.Where("quote_id = ? AND scheme_id = ?", quoteID, schemeID).First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("lookup existing snapshot: %w", err)
	}

	payload, err := BuildHandoffSnapshot(quoteID, schemeID, reason, actor)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	takeoverCount := 0
	if payload.PriorSchedule != nil {
		takeoverCount = payload.PriorSchedule.InForceCount
	}
	row := models.PolicyHandoffSnapshot{
		QuoteID:       quoteID,
		SchemeID:      schemeID,
		HandedOffBy:   actor,
		Reason:        reason,
		MemberCount:   len(payload.Members),
		TakeoverCount: takeoverCount,
		Payload:       string(body),
	}
	if err := db.Create(&row).Error; err != nil {
		return nil, fmt.Errorf("persist snapshot: %w", err)
	}
	return &row, nil
}

// ListPolicyHandoffSnapshotsForScheme returns every snapshot recorded for
// a scheme, newest first. A scheme renewed multiple times will have
// multiple snapshots — one per InForce transition.
func ListPolicyHandoffSnapshotsForScheme(schemeID int) ([]models.PolicyHandoffSnapshot, error) {
	var rows []models.PolicyHandoffSnapshot
	if err := DB.Where("scheme_id = ?", schemeID).Order("handed_off_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetPolicyHandoffSnapshot returns one snapshot row by ID.
func GetPolicyHandoffSnapshot(id int) (*models.PolicyHandoffSnapshot, error) {
	var row models.PolicyHandoffSnapshot
	if err := DB.Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

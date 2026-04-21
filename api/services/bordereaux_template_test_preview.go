package services

import (
	"api/models"
	"fmt"
	"strings"
)

// TestBordereauxTemplateRequest drives the "preview" endpoint. SampleSize
// defaults to 5; SchemeID is optional and filters the sample when provided.
type TestBordereauxTemplateRequest struct {
	SampleSize int `json:"sample_size"`
	SchemeID   int `json:"scheme_id"`
}

// TemplateFieldMappingCheck reports whether each configured source_field exists
// in the canonical field catalogue for the template's type.
type TemplateFieldMappingCheck struct {
	SourceField string `json:"source_field"`
	TargetField string `json:"target_field"`
	Required    bool   `json:"required"`
	Known       bool   `json:"known"`
}

// TestBordereauxTemplateResult is what the caller gets back.
type TestBordereauxTemplateResult struct {
	TemplateID     int                         `json:"template_id"`
	TemplateName   string                      `json:"template_name"`
	Type           string                      `json:"type"`
	SampleSize     int                         `json:"sample_size"`
	SampleSource   string                      `json:"sample_source"`
	MappingChecks  []TemplateFieldMappingCheck `json:"mapping_checks"`
	UnknownFields  []string                    `json:"unknown_fields"`
	MissingInData  []string                    `json:"missing_in_data"`
	PreviewRows    []map[string]interface{}    `json:"preview_rows"`
}

// TestBordereauxTemplate runs the template's field mappings against a small
// sample of live snapshot data (MemberBordereauxData / PremiumBordereauxData /
// GroupSchemeClaim depending on type) and returns what the mapped output would
// look like. This is the "Test Template" button in the UI — it does not
// generate a file, it just lets operators validate mappings before they
// kick off a full run.
func TestBordereauxTemplate(templateID int, req TestBordereauxTemplateRequest) (TestBordereauxTemplateResult, error) {
	out := TestBordereauxTemplateResult{}
	tmpl, err := GetBordereauxTemplateByID(templateID)
	if err != nil {
		return out, err
	}
	out.TemplateID = tmpl.ID
	out.TemplateName = tmpl.Name
	out.Type = tmpl.Type

	size := req.SampleSize
	if size <= 0 || size > 50 {
		size = 5
	}

	// 1. Validate source_field names against the canonical field catalogue.
	fields, err := GetBordereauxFieldsByType(tmpl.Type)
	if err != nil {
		return out, fmt.Errorf("fields for type %q: %w", tmpl.Type, err)
	}
	known := make(map[string]bool, len(fields))
	for _, f := range fields {
		known[f["field_name"]] = true
	}
	for _, m := range tmpl.FieldMappings {
		check := TemplateFieldMappingCheck{
			SourceField: m.SourceField,
			TargetField: m.TargetField,
			Required:    m.Required,
			Known:       known[m.SourceField],
		}
		out.MappingChecks = append(out.MappingChecks, check)
		if !check.Known && m.SourceField != "" {
			out.UnknownFields = append(out.UnknownFields, m.SourceField)
		}
	}

	// 2. Pull sample rows from the snapshot table matching this type.
	table, err := previewTableForType(tmpl.Type)
	if err != nil {
		return out, err
	}
	out.SampleSource = table

	q := DB.Table(table)
	if req.SchemeID > 0 {
		if tmpl.Type == "claim" {
			q = q.Where("scheme_id = ?", req.SchemeID)
		} else if tmpl.Type == "member" || tmpl.Type == "premium" {
			// MemberBordereauxData / PremiumBordereauxData have scheme_name, not scheme_id.
			var scheme models.GroupScheme
			if err := DB.Select("name").First(&scheme, req.SchemeID).Error; err == nil && scheme.Name != "" {
				q = q.Where("scheme_name = ?", scheme.Name)
			}
		}
	}
	var rows []map[string]interface{}
	if err := q.Order("id DESC").Limit(size).Find(&rows).Error; err != nil {
		return out, fmt.Errorf("load sample rows: %w", err)
	}
	out.SampleSize = len(rows)

	// 3. Apply mappings. Track which declared source_fields produced no value
	// in any row so the UI can flag silent misses separately from unknown ones.
	seen := make(map[string]bool)
	for _, row := range rows {
		mapped := make(map[string]interface{}, len(tmpl.FieldMappings))
		for _, m := range tmpl.FieldMappings {
			if m.TargetField == "" {
				continue
			}
			if v, ok := row[m.SourceField]; ok {
				mapped[m.TargetField] = v
				if v != nil {
					seen[m.SourceField] = true
				}
			} else {
				mapped[m.TargetField] = nil
			}
		}
		out.PreviewRows = append(out.PreviewRows, mapped)
	}
	if len(rows) > 0 {
		for _, m := range tmpl.FieldMappings {
			if m.SourceField == "" || !known[m.SourceField] {
				continue
			}
			if !seen[m.SourceField] {
				out.MissingInData = append(out.MissingInData, m.SourceField)
			}
		}
	}
	return out, nil
}

// previewTableForType returns the table name used for the template preview.
// Deliberately narrow — only the three live bordereaux types are supported.
func previewTableForType(t string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "member", "members":
		return "member_bordereaux_data", nil
	// Premium / claim template preview temporarily disabled — only member is a
	// valid use case right now.
	// case "premium", "premiums":
	// 	return "premium_bordereaux_data", nil
	// case "claim", "claims":
	// 	return "group_scheme_claims", nil
	default:
		return "", fmt.Errorf("template type %q has no preview source configured", t)
	}
}

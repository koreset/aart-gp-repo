package email

import (
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"
)

// AttachmentSpec is one entry in an outbox row's `attachments` JSON column.
// Phase 1 supports a "file" kind that points at a local filesystem path; other
// kinds (e.g. "generated_bordereaux") can be added by extending resolveOne.
//
//	[{"kind":"file","path":"/var/aart/bordereaux/123.xlsx","filename":"May2026.xlsx"}]
type AttachmentSpec struct {
	Kind        string `json:"kind"`
	Path        string `json:"path,omitempty"`
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	// Reserved for future "generated_bordereaux", "on_risk_letter", etc.
	RefID string `json:"ref_id,omitempty"`
}

// ResolveAttachments parses the outbox JSON string and reads each referenced
// source into a ready-to-attach []Attachment. Called by the worker at send
// time so that the DB never stores attachment bytes.
func ResolveAttachments(jsonSpec string) ([]Attachment, error) {
	if jsonSpec == "" || jsonSpec == "null" {
		return nil, nil
	}
	var specs []AttachmentSpec
	if err := json.Unmarshal([]byte(jsonSpec), &specs); err != nil {
		return nil, fmt.Errorf("parse attachments spec: %w", err)
	}
	out := make([]Attachment, 0, len(specs))
	for _, s := range specs {
		att, err := resolveOne(s)
		if err != nil {
			return nil, err
		}
		out = append(out, att)
	}
	return out, nil
}

func resolveOne(s AttachmentSpec) (Attachment, error) {
	switch s.Kind {
	case "file", "":
		if s.Path == "" {
			return Attachment{}, fmt.Errorf("file attachment requires path")
		}
		bytes, err := os.ReadFile(s.Path)
		if err != nil {
			return Attachment{}, fmt.Errorf("read %s: %w", s.Path, err)
		}
		filename := s.Filename
		if filename == "" {
			filename = filepath.Base(s.Path)
		}
		ct := s.ContentType
		if ct == "" {
			ct = mime.TypeByExtension(filepath.Ext(s.Path))
			if ct == "" {
				ct = "application/octet-stream"
			}
		}
		return Attachment{Filename: filename, ContentType: ct, Content: bytes}, nil
	default:
		return Attachment{}, fmt.Errorf("unsupported attachment kind: %q", s.Kind)
	}
}

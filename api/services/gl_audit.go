package services

import (
	"api/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Audit logger for the operational General Ledger and bank sub-ledger.
//
// Independent of the IFRS17 audit log (services.LogIFRS17Event). The two stacks
// do not share a table, a helper, or any imports — see api/models/gl_audit_log.go.
//
// All write paths are transactional: callers pass the *gorm.DB they are
// already mutating in so the audit row commits or rolls back atomically with
// the change it describes. A nil tx falls back to services.DB.

// LogGLEvent writes one GLAuditLog row inside the given transaction. Errors
// are returned so callers can roll back — unlike the fire-and-forget IFRS17
// logger, an audit miss here would invalidate the dual-control guarantee
// and cannot be silently swallowed.
//
// details may be any JSON-serialisable value (or nil); it's marshalled before
// persisting. Pass a struct or map describing the change ({from, to, reason,
// fields_changed, ...}); plain strings work too.
func LogGLEvent(tx *gorm.DB, eventType, objectType, objectName string, objectID int, user models.AppUser, details any) error {
	db := tx
	if db == nil {
		db = DB
	}

	var detailsJSON string
	if details != nil {
		switch v := details.(type) {
		case string:
			detailsJSON = v
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("LogGLEvent: marshal details: %w", err)
			}
			detailsJSON = string(b)
		}
	}

	entry := models.GLAuditLog{
		EventType:  eventType,
		ObjectType: objectType,
		ObjectID:   objectID,
		ObjectName: objectName,
		ChangedBy:  user.UserName,
		ChangedAt:  time.Now(),
		Details:    detailsJSON,
	}
	return db.Create(&entry).Error
}

// GLListAuditLogOptions scopes the audit log query.
type GLListAuditLogOptions struct {
	EventType  string
	ObjectType string
	ObjectID   int
	ChangedBy  string
	From       time.Time
	To         time.Time
	Limit      int
}

// GLListAuditLogUsers returns the set of user names that the audit log's
// "Changed by" filter dropdown should offer. It is the union of:
//   - distinct changed_by values that have ever appeared in gl_audit_logs
//   - all known AppUser.UserName values
//
// The union means the dropdown is useful even before any audit rows exist
// (you can pick a user who hasn't acted yet), and stays self-maintaining
// as new users come online.
func GLListAuditLogUsers() ([]string, error) {
	pool := map[string]struct{}{}
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		var seen []string
		if err := d.Model(&models.GLAuditLog{}).
			Distinct("changed_by").
			Where("changed_by IS NOT NULL AND changed_by <> ''").
			Pluck("changed_by", &seen).Error; err != nil {
			return err
		}
		for _, s := range seen {
			pool[s] = struct{}{}
		}
		var users []string
		if err := d.Model(&models.AppUser{}).
			Where("user_name IS NOT NULL AND user_name <> ''").
			Pluck("user_name", &users).Error; err != nil {
			return err
		}
		for _, s := range users {
			pool[s] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(pool))
	for s := range pool {
		out = append(out, s)
	}
	// Case-insensitive sort so the dropdown reads naturally regardless of
	// how the upstream system cased the names.
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if strings.ToLower(out[j]) < strings.ToLower(out[i]) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out, nil
}

// GLListAuditLog returns audit rows newest-first with the given filters
// applied. All filters are optional — zero values are skipped.
func GLListAuditLog(opts GLListAuditLogOptions) ([]models.GLAuditLog, error) {
	var rows []models.GLAuditLog
	err := DBReadWithResilience(context.Background(), func(d *gorm.DB) error {
		q := d.Model(&models.GLAuditLog{})
		if opts.EventType != "" {
			q = q.Where("event_type = ?", opts.EventType)
		}
		if opts.ObjectType != "" {
			q = q.Where("object_type = ?", opts.ObjectType)
		}
		if opts.ObjectID > 0 {
			q = q.Where("object_id = ?", opts.ObjectID)
		}
		if opts.ChangedBy != "" {
			q = q.Where("changed_by = ?", opts.ChangedBy)
		}
		if !opts.From.IsZero() {
			q = q.Where("changed_at >= ?", opts.From)
		}
		if !opts.To.IsZero() {
			q = q.Where("changed_at <= ?", opts.To)
		}
		limit := opts.Limit
		if limit <= 0 || limit > 1000 {
			limit = 200
		}
		return q.Order("changed_at DESC, id DESC").Limit(limit).Find(&rows).Error
	})
	return rows, err
}

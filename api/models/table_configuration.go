package models

import "time"

// TableConfiguration holds the global "is required" flag for one rating
// table type used in Group Pricing. One row per canonical table_type
// (matching the statName slot in services.gpTableSpecs, e.g. "regionLoading").
//
// When IsRequired is false the rating service must skip the DB read for that
// table; downstream variables resolve to zero via the existing graceful
// degradation path (see services/group_pricing.go around line 6884). The
// Tables UI shows "Not required" instead of "Empty" for these tables.
//
// Default IsRequired = true preserves existing behaviour for any table that
// has not yet been explicitly configured.
type TableConfiguration struct {
	ID          int       `json:"id"           gorm:"primaryKey;autoIncrement"`
	TableType   string    `json:"table_type"   gorm:"uniqueIndex;size:128"`
	DisplayName string    `json:"display_name" gorm:"size:255"`
	Category    string    `json:"category"     gorm:"size:64"`
	IsRequired  bool      `json:"is_required"  gorm:"default:true"`
	UpdatedBy   string    `json:"updated_by"   gorm:"size:255"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"   gorm:"autoCreateTime"`
}

// TableConfigurationAuditLog records every change to TableConfiguration.IsRequired.
// One row per change, append-only. Modelled on IFRS17AuditLog.
type TableConfigurationAuditLog struct {
	ID        int       `json:"id"         gorm:"primaryKey;autoIncrement"`
	TableType string    `json:"table_type" gorm:"size:128;index"`
	EventType string    `json:"event_type" gorm:"size:64"` // "config_changed" | "config_seeded"
	OldValue  bool      `json:"old_value"`
	NewValue  bool      `json:"new_value"`
	ChangedBy string    `json:"changed_by" gorm:"size:255"`
	ChangedAt time.Time `json:"changed_at"`
	Details   string    `json:"details"`
}

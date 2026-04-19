package services

import (
	appLog "api/log"
	"api/models"
	"sync"
	"time"
)

// tableRequiredCache stores the latest IsRequired value per table_type
// keyed by canonical statName (e.g. "regionLoading"). Loaded lazily on the
// first IsTableRequired() call and invalidated on every SetTableRequired().
//
// The cache is a sync.Map so it is safe for concurrent reads from quote
// calculation paths without holding a global lock.
var (
	tableRequiredCache   sync.Map // map[string]bool
	tableRequiredLoaded  bool
	tableRequiredLoadMtx sync.Mutex
)

// EnsureTableConfigurations is called once at startup. It seeds one row per
// entry in gpTableSpecs, defaulting IsRequired = true for any row that does
// not already exist. Existing rows are left untouched so user choices survive
// restarts. Mirrors the pattern used by EnsureGPTableStats.
func EnsureTableConfigurations() {
	if err := DB.AutoMigrate(&models.TableConfiguration{}, &models.TableConfigurationAuditLog{}); err != nil {
		appLog.WithField("error", err.Error()).Error("Failed to auto-migrate table_configurations / table_configuration_audit_logs")
		return
	}

	var existing []models.TableConfiguration
	DB.Find(&existing)
	have := make(map[string]struct{}, len(existing))
	for _, c := range existing {
		have[c.TableType] = struct{}{}
	}

	now := time.Now()
	for _, spec := range gpTableSpecs {
		if _, ok := have[spec.statName]; ok {
			continue
		}
		row := models.TableConfiguration{
			TableType:   spec.statName,
			DisplayName: spec.displayName,
			Category:    spec.category,
			IsRequired:  true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := DB.Create(&row).Error; err != nil {
			appLog.WithField("table_type", spec.statName).WithField("error", err.Error()).
				Warn("Failed to seed missing table_configurations row")
			continue
		}
		// Record the seed in the audit trail so the history starts cleanly.
		_ = DB.Create(&models.TableConfigurationAuditLog{
			TableType: spec.statName,
			EventType: "config_seeded",
			OldValue:  false,
			NewValue:  true,
			ChangedBy: "system",
			ChangedAt: now,
			Details:   "initial seed (default required)",
		}).Error
	}

	// Force a fresh cache load on next IsTableRequired() call.
	tableRequiredLoadMtx.Lock()
	tableRequiredLoaded = false
	tableRequiredLoadMtx.Unlock()
}

// loadTableRequiredCache populates tableRequiredCache from the DB once.
// Subsequent calls are no-ops until the cache is invalidated.
func loadTableRequiredCache() {
	tableRequiredLoadMtx.Lock()
	defer tableRequiredLoadMtx.Unlock()
	if tableRequiredLoaded {
		return
	}
	var rows []models.TableConfiguration
	if err := DB.Find(&rows).Error; err != nil {
		appLog.WithField("error", err.Error()).Warn("loadTableRequiredCache: failed to read table_configurations; defaulting to required=true")
		tableRequiredLoaded = true
		return
	}
	// Reset the map by ranging and deleting (sync.Map has no Clear).
	tableRequiredCache.Range(func(k, _ any) bool {
		tableRequiredCache.Delete(k)
		return true
	})
	for _, r := range rows {
		tableRequiredCache.Store(r.TableType, r.IsRequired)
	}
	tableRequiredLoaded = true
}

// IsTableRequired returns true if the named rating table is marked as
// required (i.e. the rating service should load it from the database and
// downstream variables should consume its values). Returns true (the safe
// default) if the configuration row is missing or any error occurs.
//
// tableType is the canonical statName from gpTableSpecs, e.g. "regionLoading"
// or "reinsuranceRegionLoading".
func IsTableRequired(tableType string) bool {
	loadTableRequiredCache()
	v, ok := tableRequiredCache.Load(tableType)
	if !ok {
		return true
	}
	required, _ := v.(bool)
	return required
}

// SetTableRequired flips the IsRequired flag for the named table and writes
// an audit row capturing the old/new value plus the active user. Returns the
// updated configuration row.
func SetTableRequired(tableType string, required bool, user models.AppUser, note string) (models.TableConfiguration, error) {
	var cfg models.TableConfiguration
	if err := DB.Where("table_type = ?", tableType).First(&cfg).Error; err != nil {
		return cfg, err
	}

	old := cfg.IsRequired
	cfg.IsRequired = required
	cfg.UpdatedBy = user.UserEmail
	cfg.UpdatedAt = time.Now()
	if err := DB.Save(&cfg).Error; err != nil {
		return cfg, err
	}

	_ = DB.Create(&models.TableConfigurationAuditLog{
		TableType: tableType,
		EventType: "config_changed",
		OldValue:  old,
		NewValue:  required,
		ChangedBy: user.UserEmail,
		ChangedAt: cfg.UpdatedAt,
		Details:   note,
	}).Error

	tableRequiredCache.Store(tableType, required)
	return cfg, nil
}

// GetTableConfigurations returns every table_configuration row, ordered by
// category then display_name so the UI render order is stable.
func GetTableConfigurations() ([]models.TableConfiguration, error) {
	var rows []models.TableConfiguration
	err := DB.Order("category asc, display_name asc").Find(&rows).Error
	return rows, err
}

// GetTableConfigurationByType fetches a single configuration row.
func GetTableConfigurationByType(tableType string) (models.TableConfiguration, error) {
	var cfg models.TableConfiguration
	err := DB.Where("table_type = ?", tableType).First(&cfg).Error
	return cfg, err
}

// GetTableConfigurationAudit returns the audit history for one table_type
// in reverse chronological order.
func GetTableConfigurationAudit(tableType string) ([]models.TableConfigurationAuditLog, error) {
	var rows []models.TableConfigurationAuditLog
	err := DB.Where("table_type = ?", tableType).Order("changed_at desc, id desc").Find(&rows).Error
	return rows, err
}

package services

import (
	"reflect"
	"strings"
	"sync"
)

// DefaultQuoteTableLimit is applied when GetGroupPricingQuoteTableData is
// called with limit == 0. Caps the response so a 5 000-member scheme
// doesn't ship a multi-megabyte JSON envelope back to the renderer.
// AG Grid's infinite-row model paginates above this.
const DefaultQuoteTableLimit = 500

// QuoteTableFieldSetSummary returns only the renderer-displayed columns
// for the requested table type.
// QuoteTableFieldSetFull returns the full struct (legacy behaviour).
const (
	QuoteTableFieldSetSummary = "summary"
	QuoteTableFieldSetFull    = "full"
)

// summaryColumnsFor returns the DB column projection for a given table
// type in summary mode. An empty result means "no projection available;
// fall back to full row". Column names are SQL column names — GORM's
// .Select() takes them verbatim.
//
// Each entry includes the renderer-displayed columns plus a few
// nice-to-haves (sum-assured, underwriting tier) that are cheap to ship
// and unlock future-proofing. liveColumnsFor() intersects this list with
// the actual schema at query time so a non-migrated deployment falls
// back to "what actually exists" rather than failing with a SQL error.
//
// Keep this list in sync with the visible columns in
// app/src/renderer/screens/group_pricing/QuoteResults.vue.
func summaryColumnsFor(tableType string) []string {
	switch tableType {
	case "member_rating_results":
		// Demographics + per-benefit capped SA + per-benefit risk
		// premium for every primary benefit. Funeral is aggregated
		// because the per-relationship SAs (main/spouse/child/parent)
		// vary per member; the aggregate risk premium is the operative
		// number on the rating row.
		//
		// `id` is intentionally omitted — Phase 4 added it as an
		// auto-increment primary key, but deployments that haven't run
		// the migration yet will fail the SELECT. The grid doesn't
		// display it anyway. liveColumnsFor() silently drops any
		// remaining un-migrated column.
		return []string{
			// Demographics
			"financial_year",
			"category",
			"member_name",
			"member_count",
			"gender",
			"date_of_birth",
			"is_original_member",
			"entry_date",
			"age_next_birthday",
			"annual_salary",
			// Underwriting tier (Phase 1)
			"underwriting_tier",
			"fcl_excess_ratio",
			"exceeds_free_cover_limit_indicator",
			// GLA
			"gla_sum_assured",
			"gla_capped_sum_assured",
			"gla_risk_premium",
			// PTD
			"ptd_capped_sum_assured",
			"ptd_risk_premium",
			// CI
			"ci_capped_sum_assured",
			"ci_risk_premium",
			// Spouse GLA
			"spouse_gla_capped_sum_assured",
			"spouse_gla_risk_premium",
			// TTD — income-based, not SA
			"ttd_capped_income",
			"ttd_risk_premium",
			// PHI — income-based, not SA
			"phi_capped_income",
			"phi_risk_premium",
			// Funeral — aggregate risk premium across all relationships
			"total_funeral_risk_premium",
		}
	case "member_data":
		return []string{
			"id",
			"member_name",
			"member_id_number",
			"scheme_category",
			"gender",
			"date_of_birth",
			"annual_salary",
			"occupation",
			"entry_date",
			"status",
		}
	case "member_premium_schedules":
		// MemberPremiumSchedule has no `id` primary key column — composite
		// (quote_id, scheme_id, member_name).
		return []string{
			"member_name",
			"category",
			"gender",
			"entry_date",
			"is_original_member",
			"total_annual_premium_payable",
			"gla_annual_premium",
			"ptd_annual_premium",
			"ci_annual_premium",
			"phi_annual_premium",
		}
	}
	return nil
}

// liveColumnCache caches `table_name -> set<column_name>` so we don't pay
// for an information_schema lookup on every quote-table request.
// Invalidated only by process restart — schema changes that need to be
// picked up mid-process can call ResetLiveColumnCache.
var (
	liveColumnCacheMu sync.RWMutex
	liveColumnCache   = map[string]map[string]bool{}
)

// LiveColumnsFor returns the set of columns that currently exist on the
// given table according to the live database schema. Used to intersect
// the projection so SELECT requests don't reference columns that haven't
// been migrated yet.
//
// Returns an empty map (not nil) if the introspection fails — callers
// should treat that as "fall back to full row".
func LiveColumnsFor(tableName string) map[string]bool {
	liveColumnCacheMu.RLock()
	if cols, ok := liveColumnCache[tableName]; ok {
		liveColumnCacheMu.RUnlock()
		return cols
	}
	liveColumnCacheMu.RUnlock()

	liveColumnCacheMu.Lock()
	defer liveColumnCacheMu.Unlock()
	if cols, ok := liveColumnCache[tableName]; ok {
		return cols // re-check after acquiring the write lock
	}
	cols := make(map[string]bool)
	if DB != nil {
		if types, err := DB.Migrator().ColumnTypes(tableName); err == nil {
			for _, t := range types {
				cols[t.Name()] = true
			}
		}
	}
	liveColumnCache[tableName] = cols
	return cols
}

// ResetLiveColumnCache clears the cached column lookup. Call after a
// schema migration runs in-process so the next request picks up the new
// columns without a restart.
func ResetLiveColumnCache() {
	liveColumnCacheMu.Lock()
	defer liveColumnCacheMu.Unlock()
	liveColumnCache = map[string]map[string]bool{}
}

// filterToLiveColumns intersects the requested projection with what
// actually exists in the database. Returns the surviving subset in the
// caller's order. When introspection fails (unknown table, no DB) the
// request is passed through unchanged.
func filterToLiveColumns(tableName string, requested []string) []string {
	live := LiveColumnsFor(tableName)
	if len(live) == 0 {
		return requested
	}
	out := make([]string, 0, len(requested))
	for _, c := range requested {
		if live[c] {
			out = append(out, c)
		}
	}
	return out
}

// projectStructToMap copies the json-tagged fields named in wanted from v
// into a fresh map. The map keys are the json tag names (the same strings
// the renderer's `json_tags` array uses), so the grid sees the same wire
// shape it gets in full-row mode — just with fewer keys.
//
// Cost: one reflect.NumField loop per row. ~373 lookups × 500 rows is
// roughly an order of magnitude cheaper than json.Marshal+Unmarshal of
// the full struct.
func projectStructToMap(v interface{}, wanted []string) map[string]interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	rt := rv.Type()
	want := make(map[string]struct{}, len(wanted))
	for _, t := range wanted {
		want[t] = struct{}{}
	}
	out := make(map[string]interface{}, len(wanted))
	for i := 0; i < rt.NumField(); i++ {
		tag := strings.SplitN(rt.Field(i).Tag.Get("json"), ",", 2)[0]
		if tag == "" || tag == "-" {
			continue
		}
		if _, ok := want[tag]; !ok {
			continue
		}
		out[tag] = rv.Field(i).Interface()
	}
	return out
}

// projectSliceToMaps applies projectStructToMap across a slice of structs
// addressed via reflection (rather than typed per call site). Callers
// pass a slice value; the helper iterates and produces the per-row map.
func projectSliceToMaps(slice interface{}, wanted []string) []map[string]interface{} {
	rv := reflect.ValueOf(slice)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		return nil
	}
	out := make([]map[string]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out[i] = projectStructToMap(rv.Index(i).Interface(), wanted)
	}
	return out
}

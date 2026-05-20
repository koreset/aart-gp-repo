package services

import (
	"api/models"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ErrSchemeTreatyConflict is returned by the link APIs when adding the
// requested scheme(s) would violate the "one active treaty per
// (scheme, line_of_business, treaty_type)" rule. The wrapped message names
// every conflicting treaty so the admin knows what to unlink first.
var ErrSchemeTreatyConflict = errors.New("scheme already linked to an active treaty for this line of business and treaty type")

// SchemeTreatyConflict describes one (scheme, LoB, treaty_type) tuple that
// has ≥2 active treaty links. Used by the conflict-report endpoint AND as
// the payload behind pre-link rejection messages.
type SchemeTreatyConflict struct {
	SchemeID       int             `json:"scheme_id"`
	SchemeName     string          `json:"scheme_name"`
	LineOfBusiness string          `json:"line_of_business"`
	TreatyType     string          `json:"treaty_type"`
	Treaties       []ConflictParty `json:"treaties"`
}

// ConflictParty identifies one treaty involved in a conflict.
type ConflictParty struct {
	TreatyID     int    `json:"treaty_id"`
	TreatyNumber string `json:"treaty_number"`
	TreatyName   string `json:"treaty_name"`
}

// FindConflictsForNewLinks checks whether adding schemeIDs to treatyID
// would create a violation. Returns the list of conflicts; an empty slice
// means safe to proceed. The caller wraps non-empty results in
// ErrSchemeTreatyConflict via FormatConflictError.
//
// Excludes treatyID itself from the "already linked" check so re-running
// the request (idempotency) doesn't false-positive on links that already
// belong to this treaty. Linking schemes to a draft treaty is allowed —
// only ACTIVE treaties produce conflicts.
func FindConflictsForNewLinks(db *gorm.DB, treatyID int, schemeIDs []int) ([]SchemeTreatyConflict, error) {
	if db == nil {
		db = DB
	}
	if len(schemeIDs) == 0 {
		return nil, nil
	}

	var target models.ReinsuranceTreaty
	if err := db.First(&target, treatyID).Error; err != nil {
		return nil, fmt.Errorf("failed to load target treaty: %w", err)
	}
	if !strings.EqualFold(target.Status, "active") {
		return nil, nil
	}

	type row struct {
		SchemeID     int
		SchemeName   string
		TreatyID     int
		TreatyNumber string
		TreatyName   string
	}
	var rows []row
	err := db.Table("treaty_scheme_links AS tsl").
		Select("tsl.scheme_id AS scheme_id, tsl.scheme_name AS scheme_name, rt.id AS treaty_id, rt.treaty_number AS treaty_number, rt.treaty_name AS treaty_name").
		Joins("JOIN reinsurance_treaties AS rt ON rt.id = tsl.treaty_id").
		Where("tsl.scheme_id IN ?", schemeIDs).
		Where("tsl.treaty_id <> ?", treatyID).
		Where("rt.status = ?", "active").
		Where("rt.line_of_business = ?", target.LineOfBusiness).
		Where("rt.treaty_type = ?", target.TreatyType).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	byScheme := make(map[int]*SchemeTreatyConflict, len(rows))
	for _, r := range rows {
		c, ok := byScheme[r.SchemeID]
		if !ok {
			c = &SchemeTreatyConflict{
				SchemeID:       r.SchemeID,
				SchemeName:     r.SchemeName,
				LineOfBusiness: target.LineOfBusiness,
				TreatyType:     target.TreatyType,
			}
			byScheme[r.SchemeID] = c
		}
		c.Treaties = append(c.Treaties, ConflictParty{
			TreatyID:     r.TreatyID,
			TreatyNumber: r.TreatyNumber,
			TreatyName:   r.TreatyName,
		})
	}
	out := make([]SchemeTreatyConflict, 0, len(byScheme))
	for _, c := range byScheme {
		out = append(out, *c)
	}
	return out, nil
}

// ListAllSchemeTreatyConflicts scans the whole tenant for existing
// violations. Drives the conflict banner + report on the Treaty Management
// screen so admins can clean up legacy data.
//
// Implementation note: uses a self-join (rather than a tuple-IN subquery)
// to keep the SQL portable across MySQL, PostgreSQL, and MSSQL — tuple-IN
// is patchily supported on older SQL Server editions.
func ListAllSchemeTreatyConflicts(db *gorm.DB) ([]SchemeTreatyConflict, error) {
	if db == nil {
		db = DB
	}

	type row struct {
		SchemeID       int
		SchemeName     string
		LineOfBusiness string
		TreatyType     string
		TreatyID       int
		TreatyNumber   string
		TreatyName     string
	}
	var hits []row
	err := db.Raw(`
        SELECT tsl.scheme_id        AS scheme_id,
               tsl.scheme_name      AS scheme_name,
               rt.line_of_business  AS line_of_business,
               rt.treaty_type       AS treaty_type,
               rt.id                AS treaty_id,
               rt.treaty_number     AS treaty_number,
               rt.treaty_name       AS treaty_name
        FROM treaty_scheme_links tsl
        JOIN reinsurance_treaties rt ON rt.id = tsl.treaty_id
        WHERE rt.status = 'active'
          AND EXISTS (
              SELECT 1
              FROM treaty_scheme_links tsl2
              JOIN reinsurance_treaties rt2 ON rt2.id = tsl2.treaty_id
              WHERE rt2.status = 'active'
                AND tsl2.scheme_id = tsl.scheme_id
                AND rt2.line_of_business = rt.line_of_business
                AND rt2.treaty_type = rt.treaty_type
                AND rt2.id <> rt.id
          )
        ORDER BY tsl.scheme_id, rt.line_of_business, rt.treaty_type, rt.id
    `).Scan(&hits).Error
	if err != nil {
		return nil, err
	}

	type key struct {
		schemeID   int
		lob        string
		treatyType string
	}
	grouped := make(map[key]*SchemeTreatyConflict)
	order := make([]key, 0)
	for _, h := range hits {
		k := key{h.SchemeID, h.LineOfBusiness, h.TreatyType}
		c, ok := grouped[k]
		if !ok {
			c = &SchemeTreatyConflict{
				SchemeID:       h.SchemeID,
				SchemeName:     h.SchemeName,
				LineOfBusiness: h.LineOfBusiness,
				TreatyType:     h.TreatyType,
			}
			grouped[k] = c
			order = append(order, k)
		}
		c.Treaties = append(c.Treaties, ConflictParty{
			TreatyID:     h.TreatyID,
			TreatyNumber: h.TreatyNumber,
			TreatyName:   h.TreatyName,
		})
	}
	out := make([]SchemeTreatyConflict, 0, len(order))
	for _, k := range order {
		out = append(out, *grouped[k])
	}
	return out, nil
}

// FormatConflictError turns a non-empty conflict slice into the wrapped
// error the link APIs return. Stable single-line format so the snackbar
// renders cleanly. Returns nil if the slice is empty so callers can use
// `if err := FormatConflictError(c); err != nil` directly.
func FormatConflictError(conflicts []SchemeTreatyConflict) error {
	if len(conflicts) == 0 {
		return nil
	}
	parts := make([]string, 0, len(conflicts))
	for _, c := range conflicts {
		treaties := make([]string, 0, len(c.Treaties))
		for _, p := range c.Treaties {
			treaties = append(treaties, p.TreatyNumber)
		}
		parts = append(parts, fmt.Sprintf(
			"scheme %d (%s) already on %s/%s via %s",
			c.SchemeID, c.SchemeName, c.LineOfBusiness, c.TreatyType,
			strings.Join(treaties, ", "),
		))
	}
	return fmt.Errorf("%w: %s", ErrSchemeTreatyConflict, strings.Join(parts, "; "))
}

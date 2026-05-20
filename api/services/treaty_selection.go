package services

import (
	"api/models"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Treaty-selection errors. Distinct values so callers can react differently
// (silent vs warn vs alert) — Chunk C uses these to drive log levels.
var (
	ErrTreatyNoneLinked     = errors.New("scheme has no treaty linked for this benefit and treaty type")
	ErrTreatyNoneApplicable = errors.New("scheme has treaties for this benefit and treaty type but none cover the relevant date")
	ErrTreatyAmbiguous      = errors.New("multiple treaties match — data integrity issue, awaiting uniqueness guard")
)

// SelectApplicableTreaty picks the single treaty that should apply to a given
// claim line, based on the scheme link, benefit-to-LoB mapping, treaty type,
// treaty basis, and the relevant coverage date (member entry for
// risk-attaching, claim date-of-event for loss-occurring).
//
// Returns one of the sentinel errors above when there is no clean answer.
// Callers MUST handle these — log-level mapping is the caller's choice.
func SelectApplicableTreaty(
	db *gorm.DB,
	schemeID int,
	benefitCode string,
	treatyType string,
	memberEntryDate time.Time,
	claimDateOfEvent string,
) (*models.ReinsuranceTreaty, error) {
	if db == nil {
		db = DB
	}

	lob := mapBenefitToLineOfBusiness(benefitCode)
	if lob == "" {
		return nil, ErrTreatyNoneLinked
	}

	var candidates []models.ReinsuranceTreaty
	err := db.
		Joins("JOIN treaty_scheme_links ON treaty_scheme_links.treaty_id = reinsurance_treaties.id").
		Where("treaty_scheme_links.scheme_id = ?", schemeID).
		Where("reinsurance_treaties.status = ?", "active").
		Where("reinsurance_treaties.line_of_business = ?", lob).
		Where("reinsurance_treaties.treaty_type = ?", treatyType).
		Find(&candidates).Error
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, ErrTreatyNoneLinked
	}

	claimEvent, _ := parseDateOnly(claimDateOfEvent)

	matches := make([]models.ReinsuranceTreaty, 0, len(candidates))
	for _, t := range candidates {
		coverageDate := coverageDateFor(t, memberEntryDate, claimEvent)
		if coverageDate.IsZero() {
			continue
		}
		if coversDate(t, coverageDate) {
			matches = append(matches, t)
		}
	}

	switch len(matches) {
	case 0:
		return nil, ErrTreatyNoneApplicable
	case 1:
		return &matches[0], nil
	default:
		return nil, ErrTreatyAmbiguous
	}
}

// coverageDateFor returns the date that the treaty's basis says we should
// compare against the treaty's effective/expiry window.
func coverageDateFor(t models.ReinsuranceTreaty, memberEntry, claimEvent time.Time) time.Time {
	if strings.EqualFold(t.TreatyBasis, "loss_occurring") {
		return claimEvent
	}
	return memberEntry
}

// coversDate returns true when d lies within the treaty's effective/expiry
// window. Empty ExpiryDate is treated as open-ended; malformed ExpiryDate
// also falls through to open-ended so a typo in admin data doesn't quietly
// exclude every otherwise-active claim.
func coversDate(t models.ReinsuranceTreaty, d time.Time) bool {
	effective, ok := parseDateOnly(t.EffectiveDate)
	if !ok {
		return false
	}
	if d.Before(effective) {
		return false
	}
	if strings.TrimSpace(t.ExpiryDate) == "" {
		return true
	}
	expiry, ok := parseDateOnly(t.ExpiryDate)
	if !ok {
		return true
	}
	return !d.After(expiry)
}

// mapBenefitToLineOfBusiness translates a claim benefit code into the
// treaty's LineOfBusiness value. Returns "" for benefits that aren't
// reinsurable (unknown codes). Kept as a switch — the mapping is stable and
// small; a lookup table is unwarranted for a dozen codes.
func mapBenefitToLineOfBusiness(benefitCode string) string {
	switch strings.ToUpper(strings.TrimSpace(benefitCode)) {
	case "GLA", "SGLA":
		return "group_life"
	case "PTD", "CI", "TTD":
		return "group_disability"
	case "PHI":
		return "group_health"
	case "FUNERAL", "SPOUSE_FUNERAL", "CHILD_FUNERAL",
		"PARENT_FUNERAL", "EXTENDED_FUNERAL":
		return "funeral"
	default:
		return ""
	}
}

// parseDateOnly is a thin wrapper around time.Parse for the "2006-01-02"
// layout used on treaty/claim date columns. Returns (zero, false) for empty
// or malformed input so callers can branch cleanly.
func parseDateOnly(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

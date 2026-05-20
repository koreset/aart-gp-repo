package services

import (
	"api/models"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ReinsuranceRecoveryComputation is the per-claim output of auto-flagging.
// Required and Amount are denormalised onto ClaimPaymentScheduleItem;
// SkipReason is logged (not stored) so finance can fall back to the manual
// SetReinsuranceRecovery flow when auto-detection can't reach a safe answer.
type ReinsuranceRecoveryComputation struct {
	Required   bool
	Amount     float64
	SkipReason string
}

// ComputeReinsuranceRecovery picks the treaty that applies to a claim using
// SelectApplicableTreaty (basis-aware: risk-attaching matches the member's
// EntryDate, loss-occurring matches the claim's DateOfEvent) and runs the
// cession calc. Hardcodes treaty_type="proportional" because XL/stop-loss
// treaties are event/threshold driven and tracked through a separate flow.
//
// Best-effort: every failure path returns a zeroed result so schedule
// creation never fails on this. Selector errors are mapped to log levels:
//   - ErrTreatyNoneLinked → silent (genuinely unreinsured for this benefit)
//   - ErrTreatyNoneApplicable → info log (treaty exists but date doesn't fit)
//   - ErrTreatyAmbiguous → warn log (Chunk B guard should prevent this)
//   - other DB error → warn log
//
// Member entry date is resolved via the canonical lookup against
// GPricingMemberDataInForce; if the member isn't on file, the zero time
// flows into the selector and resolves to ErrTreatyNoneApplicable for
// risk-attaching treaties.
func ComputeReinsuranceRecovery(claim models.GroupSchemeClaim, db *gorm.DB) ReinsuranceRecoveryComputation {
	if db == nil {
		db = DB
	}

	memberEntry := lookupMemberEntryDate(db, claim.SchemeId, claim.MemberIDNumber, claim.ID)

	treaty, err := SelectApplicableTreaty(
		db, claim.SchemeId, claim.BenefitCode, "proportional",
		memberEntry, claim.DateOfEvent,
	)
	if err != nil {
		return recoveryFromSelectorError(err, claim.SchemeId, claim.ID)
	}

	required, amount := applyTreatyCession(claim.ClaimAmount, *treaty)
	return ReinsuranceRecoveryComputation{Required: required, Amount: amount}
}

// lookupMemberEntryDate resolves the member's EntryDate for the
// risk-attaching coverage check. A missing member is not fatal — the
// selector handles a zero date by falling through to ErrTreatyNoneApplicable.
// Ordering by year DESC picks the most recent in-force snapshot deterministically.
func lookupMemberEntryDate(db *gorm.DB, schemeID int, memberIDNumber string, claimID int) time.Time {
	if memberIDNumber == "" {
		return time.Time{}
	}
	var member models.GPricingMemberDataInForce
	if err := db.
		Where("scheme_id = ? AND member_id_number = ?", schemeID, memberIDNumber).
		Order("year DESC").
		First(&member).Error; err != nil {
		log.Debug().Err(err).
			Int("scheme_id", schemeID).
			Int("claim_id", claimID).
			Str("member_id_number", memberIDNumber).
			Msg("reinsurance auto-flag: member entry date not found")
		return time.Time{}
	}
	return member.EntryDate
}

// recoveryFromSelectorError maps each treaty-selection sentinel to the
// appropriate skip result + log level. Pure on the result side (the log
// call is a side-effect), so the mapping is unit-testable.
func recoveryFromSelectorError(err error, schemeID, claimID int) ReinsuranceRecoveryComputation {
	switch {
	case errors.Is(err, ErrTreatyNoneLinked):
		// Genuinely unreinsured for this benefit + treaty type. Silent —
		// most schemes ARE reinsured, but a handful aren't, and logging
		// every one would drown out real signal.
		return ReinsuranceRecoveryComputation{SkipReason: "no treaty linked"}
	case errors.Is(err, ErrTreatyNoneApplicable):
		log.Info().
			Int("scheme_id", schemeID).
			Int("claim_id", claimID).
			Msg("reinsurance auto-flag: scheme has a matching treaty but none cover the relevant date — finance should review treaty coverage")
		return ReinsuranceRecoveryComputation{SkipReason: "no applicable treaty for relevant date"}
	case errors.Is(err, ErrTreatyAmbiguous):
		log.Warn().
			Int("scheme_id", schemeID).
			Int("claim_id", claimID).
			Msg("reinsurance auto-flag: multiple treaties match — Chunk B uniqueness guard should have prevented this")
		return ReinsuranceRecoveryComputation{SkipReason: "ambiguous treaty selection"}
	default:
		log.Warn().Err(err).
			Int("scheme_id", schemeID).
			Int("claim_id", claimID).
			Msg("reinsurance auto-flag: treaty selection failed")
		return ReinsuranceRecoveryComputation{SkipReason: "treaty selection failed"}
	}
}

// applyTreatyCession is the pure decision: given a gross claim amount and a
// treaty, what does the schedule item recover? Separated from the orchestrator
// so the cession matrix (treaty types, retention thresholds, tier boundaries)
// is unit-testable without DB scaffolding.
func applyTreatyCession(claimAmount float64, treaty models.ReinsuranceTreaty) (required bool, amount float64) {
	ceded, _, belowRetention := CalculateClaimCession(claimAmount, treaty)
	if belowRetention || ceded <= 0 {
		return false, 0
	}
	return true, ceded
}

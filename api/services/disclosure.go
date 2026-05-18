package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"api/models"
)

// Disclosure / consent / attestation form variants. Mirrored in the
// renderer — keep both sides in sync.
const (
	DisclosureVariantShort = "short"
	DisclosureVariantLong  = "long"

	DisclosureViaBroker     = "broker"
	DisclosureViaMemberSelf = "member_self"
	DisclosureViaUWer       = "underwriter"

	ConsentTypeMedicalInfo     = "medical_info"
	ConsentTypePathology       = "pathology"
	ConsentTypeGPRecords       = "gp_records"
	ConsentTypeActivelyAtWork  = "actively_at_work"
)

// DisclosureSubmission is the controller-facing payload. The service folds
// it into a MemberDisclosure row, computes BMI, encodes the JSON columns,
// and re-evaluates the rules engine.
type DisclosureSubmission struct {
	Height                float64        `json:"height"`
	Weight                float64        `json:"weight"`
	Smoker                bool           `json:"smoker"`
	CigarettesPerDay      int            `json:"cigarettes_per_day"`
	AlcoholUnitsPerWeek   float64        `json:"alcohol_units_per_week"`
	HasHazardousHobbies   bool           `json:"has_hazardous_hobbies"`
	HazardousHobbies      string         `json:"hazardous_hobbies"`
	OccupationRiskAnswers map[string]any `json:"occupation_risk_answers"`
	DisclosedConditions   []string       `json:"disclosed_conditions"`
	AdditionalNotes       string         `json:"additional_notes"`
	FormVariant           string         `json:"form_variant"`
	SubmittedVia          string         `json:"submitted_via"`
}

// SubmitDisclosure persists a disclosure on the case, snapshots the
// engine's outcome onto the case, and records audit events. Submitting a
// second disclosure replaces the previous engine snapshot but does NOT
// delete the prior disclosure row — every submission is retained.
//
// Caller must ensure a ConsentRecord of type `medical_info` exists on the
// case; this is enforced here.
func SubmitDisclosure(caseID int, payload DisclosureSubmission, user models.AppUser) (*models.MemberDisclosure, error) {
	var c models.UnderwritingCase
	if err := DB.Where("id = ?", caseID).First(&c).Error; err != nil {
		return nil, err
	}
	if err := requireConsent(caseID, ConsentTypeMedicalInfo); err != nil {
		return nil, err
	}

	bmi := computeBMI(payload.Height, payload.Weight)
	occRiskJSON, _ := json.Marshal(payload.OccupationRiskAnswers)
	condJSON, _ := json.Marshal(payload.DisclosedConditions)

	if payload.FormVariant == "" {
		payload.FormVariant = DisclosureVariantShort
	}
	if payload.SubmittedVia == "" {
		payload.SubmittedVia = DisclosureViaUWer
	}

	disc := models.MemberDisclosure{
		CaseID:                caseID,
		Height:                payload.Height,
		Weight:                payload.Weight,
		BMI:                   bmi,
		Smoker:                payload.Smoker,
		CigarettesPerDay:      payload.CigarettesPerDay,
		AlcoholUnitsPerWeek:   payload.AlcoholUnitsPerWeek,
		HasHazardousHobbies:   payload.HasHazardousHobbies,
		HazardousHobbies:      payload.HazardousHobbies,
		OccupationRiskAnswers: string(occRiskJSON),
		DisclosedConditions:   string(condJSON),
		AdditionalNotes:       payload.AdditionalNotes,
		FormVariant:           payload.FormVariant,
		SubmittedVia:          payload.SubmittedVia,
		SubmittedBy:           user.UserEmail,
	}
	if err := DB.Create(&disc).Error; err != nil {
		return nil, fmt.Errorf("save disclosure: %w", err)
	}

	// Engine evaluation. Snapshot the outcome onto the case so the case-list
	// can sort by engine recommendation and the case-detail UI surfaces it
	// without re-running the engine every read.
	ctx := BuildRuleContextForCase(c, nil)
	mergeDisclosureIntoContext(ctx, payload, bmi)
	summary, err := EvaluateAgainstRuleSet(c.RuleSetID, ctx)
	if err == nil {
		exclusionsJSON, _ := json.Marshal(summary.Exclusions)
		now := time.Now()
		c.EngineOutcome = summary.Outcome
		c.EngineLoading = summary.MaxLoading
		c.EngineExclusions = string(exclusionsJSON)
		c.EngineEvaluatedAt = &now
		if err := DB.Save(&c).Error; err != nil {
			return nil, fmt.Errorf("save case engine snapshot: %w", err)
		}
	}

	recordCaseEvent(caseID, "disclosure_submitted", user.UserEmail, map[string]any{
		"form_variant":         payload.FormVariant,
		"submitted_via":        payload.SubmittedVia,
		"bmi":                  bmi,
		"disclosed_conditions": payload.DisclosedConditions,
		"engine_outcome":       c.EngineOutcome,
		"engine_loading":       c.EngineLoading,
	})

	// Auto-advance the case from pending_evidence to in_review on first
	// disclosure — the underwriter now has something to assess. Failures
	// here are non-fatal; the disclosure is already persisted.
	if c.Status == models.UWCaseStatusPendingEvidence {
		_, _ = TransitionUnderwritingCase(caseID, models.UWCaseStatusInReview, user, "disclosure submitted")
	}

	return &disc, nil
}

// GetLatestDisclosure returns the most recent MemberDisclosure for the case,
// or nil if none exists.
func GetLatestDisclosure(caseID int) (*models.MemberDisclosure, error) {
	var disc models.MemberDisclosure
	err := DB.Where("case_id = ?", caseID).Order("submitted_at DESC").First(&disc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &disc, nil
}

// ListDisclosures returns all disclosures on a case, newest first.
func ListDisclosures(caseID int) ([]models.MemberDisclosure, error) {
	var rows []models.MemberDisclosure
	if err := DB.Where("case_id = ?", caseID).Order("submitted_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ─── Actively-at-work attestation ──────────────────────────────────────────

// AttestationSubmission is the controller-facing payload.
type AttestationSubmission struct {
	QuoteID        int    `json:"quote_id"`
	CaseID         int    `json:"case_id"`
	MemberIdNumber string `json:"member_id_number"`
	MemberName     string `json:"member_name"`
	AttestedByName string `json:"attested_by_name"`
	AttestedByRole string `json:"attested_by_role"`
}

// SubmitAttestation persists a digital actively-at-work attestation. The
// signature hash is SHA-256 of name + timestamp + IP — Phase 6 will
// replace this with a vendor-signed payload.
func SubmitAttestation(payload AttestationSubmission, user models.AppUser, ipAddress, userAgent string) (*models.ActivelyAtWorkAttestation, error) {
	if payload.AttestedByName == "" {
		return nil, errors.New("attested_by_name required")
	}
	if payload.QuoteID == 0 && payload.CaseID == 0 {
		return nil, errors.New("quote_id or case_id required")
	}
	if payload.QuoteID == 0 && payload.CaseID > 0 {
		var c models.UnderwritingCase
		if err := DB.Where("id = ?", payload.CaseID).First(&c).Error; err == nil {
			payload.QuoteID = c.QuoteID
			if payload.MemberIdNumber == "" {
				payload.MemberIdNumber = c.MemberIdNumber
			}
			if payload.MemberName == "" {
				payload.MemberName = c.MemberName
			}
		}
	}
	now := time.Now()
	att := models.ActivelyAtWorkAttestation{
		CaseID:          payload.CaseID,
		QuoteID:         payload.QuoteID,
		MemberIdNumber:  payload.MemberIdNumber,
		MemberName:      payload.MemberName,
		AttestedAt:      now,
		AttestedByName:  payload.AttestedByName,
		AttestedByRole:  payload.AttestedByRole,
		AttestedByEmail: user.UserEmail,
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
	}
	att.SignatureHash = signatureHash(att.AttestedByName, att.AttestedAt, att.IPAddress)
	if err := DB.Create(&att).Error; err != nil {
		return nil, err
	}
	if payload.CaseID > 0 {
		recordCaseEvent(payload.CaseID, "attestation_signed", user.UserEmail, map[string]any{
			"attested_by_name": payload.AttestedByName,
			"attested_by_role": payload.AttestedByRole,
		})
	}
	return &att, nil
}

// ListAttestations returns attestations for a case (or quote when caseID=0).
func ListAttestations(caseID, quoteID int) ([]models.ActivelyAtWorkAttestation, error) {
	q := DB.Model(&models.ActivelyAtWorkAttestation{})
	if caseID > 0 {
		q = q.Where("case_id = ?", caseID)
	} else if quoteID > 0 {
		q = q.Where("quote_id = ?", quoteID)
	} else {
		return nil, errors.New("case_id or quote_id required")
	}
	var rows []models.ActivelyAtWorkAttestation
	if err := q.Order("attested_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ─── Consent ───────────────────────────────────────────────────────────────

// ConsentSubmission is the controller-facing payload.
type ConsentSubmission struct {
	CaseID         int    `json:"case_id"`
	QuoteID        int    `json:"quote_id"`
	ConsentType    string `json:"consent_type"`
	GrantedByName  string `json:"granted_by_name"`
	GrantedByEmail string `json:"granted_by_email"`
}

// SubmitConsent persists a ConsentRecord. The signature hash is computed
// the same way as attestations; Phase 6 will swap in a vendor-signed
// payload via the e-sign abstraction.
func SubmitConsent(payload ConsentSubmission, user models.AppUser, ipAddress string) (*models.ConsentRecord, error) {
	if payload.ConsentType == "" {
		return nil, errors.New("consent_type required")
	}
	if payload.GrantedByName == "" {
		return nil, errors.New("granted_by_name required")
	}
	now := time.Now()
	if payload.GrantedByEmail == "" {
		payload.GrantedByEmail = user.UserEmail
	}
	rec := models.ConsentRecord{
		CaseID:         payload.CaseID,
		QuoteID:        payload.QuoteID,
		ConsentType:    payload.ConsentType,
		GrantedAt:      now,
		GrantedByName:  payload.GrantedByName,
		GrantedByEmail: payload.GrantedByEmail,
		IPAddress:      ipAddress,
		SignatureHash:  signatureHash(payload.GrantedByName, now, ipAddress),
	}
	if err := DB.Create(&rec).Error; err != nil {
		return nil, err
	}
	if payload.CaseID > 0 {
		recordCaseEvent(payload.CaseID, "consent_recorded", user.UserEmail, map[string]any{
			"consent_type":     payload.ConsentType,
			"granted_by_name":  payload.GrantedByName,
			"granted_by_email": payload.GrantedByEmail,
		})
	}
	return &rec, nil
}

// ListConsents returns consents for a case, newest first.
func ListConsents(caseID int) ([]models.ConsentRecord, error) {
	var rows []models.ConsentRecord
	if err := DB.Where("case_id = ?", caseID).Order("granted_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// requireConsent enforces that a non-revoked ConsentRecord of the given
// type exists on the case.
func requireConsent(caseID int, consentType string) error {
	var count int64
	if err := DB.Model(&models.ConsentRecord{}).
		Where("case_id = ? AND consent_type = ?", caseID, consentType).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("missing consent for %s; record consent before submitting", consentType)
	}
	return nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

// computeBMI returns weight (kg) / height (m)^2, or 0 when inputs are
// invalid. Height is supplied in centimetres.
func computeBMI(heightCm, weightKg float64) float64 {
	if heightCm <= 0 || weightKg <= 0 {
		return 0
	}
	heightM := heightCm / 100.0
	return weightKg / (heightM * heightM)
}

// signatureHash is a deterministic hex SHA-256 of name + timestamp + IP.
// Sufficient for "we have evidence the user typed this on this date from
// this IP"; Phase 6 swaps in a real e-signature.
func signatureHash(name string, ts time.Time, ip string) string {
	h := sha256.Sum256([]byte(name + "|" + ts.Format(time.RFC3339Nano) + "|" + ip))
	return hex.EncodeToString(h[:])
}

// mergeDisclosureIntoContext folds disclosure fields into an existing rule
// evaluation context. The keys mirror the field names rule authors are
// expected to reference: `bmi`, `smoker`, `cigarettes_per_day`,
// `alcohol_units_per_week`, `hazardous_hobbies` (bool), plus one boolean
// per disclosed condition code (e.g. `condition_diabetes_type2 = true`).
// Numeric fields are stored as float64 so the engine's comparison operators
// work consistently.
func mergeDisclosureIntoContext(ctx EvaluationContext, payload DisclosureSubmission, bmi float64) {
	if bmi > 0 {
		ctx["bmi"] = bmi
	}
	if payload.Height > 0 {
		ctx["height_cm"] = payload.Height
	}
	if payload.Weight > 0 {
		ctx["weight_kg"] = payload.Weight
	}
	ctx["smoker"] = payload.Smoker
	if payload.CigarettesPerDay > 0 {
		ctx["cigarettes_per_day"] = float64(payload.CigarettesPerDay)
	}
	if payload.AlcoholUnitsPerWeek > 0 {
		ctx["alcohol_units_per_week"] = payload.AlcoholUnitsPerWeek
	}
	ctx["hazardous_hobbies"] = payload.HasHazardousHobbies
	for _, code := range payload.DisclosedConditions {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		ctx["condition_"+strings.ToLower(code)] = true
	}
	// Free-form occupation risk answers are flattened into the context as
	// "occ_<key>" so rule authors can reference them directly.
	for k, v := range payload.OccupationRiskAnswers {
		ctx["occ_"+strings.ToLower(k)] = v
	}
}

package services

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"

	"api/models"
)

// Takeover outcome strings. Mirrored in the renderer; keep both sides in
// step when adding new ones. Outcomes correspond to UWRule rows authored
// under `category = "takeover"` — the engine returns the strictest match.
const (
	TakeoverOutcomeContinuationNoEvidence  = "continuation_no_evidence"
	TakeoverOutcomeContinuationWithLoading = "continuation_with_loading"
	TakeoverOutcomeNewEvidenceRequired     = "new_evidence_required"
	TakeoverOutcomeUnmatched               = "unmatched"
)

// priorInsurerScheduleDir is the on-disk root for uploaded prior-insurer
// schedule documents (the actual schedule PDFs / Excels / CSVs the broker
// hands over). Mirrors the Phase-2 attachment layout.
const priorInsurerScheduleDir = "tmp/uploads/prior_insurer_schedules"

// PriorInsurerCSVRow is the canonical schedule row shape. Broker uploads
// must conform to this header; per-insurer adapters can be added later
// without touching the matcher.
type PriorInsurerCSVRow struct {
	MemberIDNumber string  `csv:"member_id_number"`
	MemberName     string  `csv:"member_name"`
	DateOfBirth    string  `csv:"date_of_birth"`
	GlaSumAssured  float64 `csv:"gla_sum_assured"`
	PtdSumAssured  float64 `csv:"ptd_sum_assured"`
	CiSumAssured   float64 `csv:"ci_sum_assured"`
	PriorLoadings  string  `csv:"prior_loadings"`  // pipe-separated `benefit:percent` (e.g. "gla:25|ptd:10")
	PriorExclusions string `csv:"prior_exclusions"` // pipe-separated exclusion codes (e.g. "diabetes|smoker")
	InForce        string  `csv:"in_force"`         // "true" / "false" / "yes" / "no"
}

// TakeoverMatchSummary is returned to the renderer after MatchPriorMembersToCensus
// runs. Drives the upload-preview UI.
type TakeoverMatchSummary struct {
	ScheduleID                int `json:"schedule_id"`
	Total                     int `json:"total"`
	InForce                   int `json:"in_force"`
	MatchedByID               int `json:"matched_by_id"`
	MatchedByNameAndDOB       int `json:"matched_by_name_and_dob"`
	Unmatched                 int `json:"unmatched"`
	ContinuationNoEvidence    int `json:"continuation_no_evidence"`
	ContinuationWithLoading   int `json:"continuation_with_loading"`
	NewEvidenceRequired       int `json:"new_evidence_required"`
}

// ImportPriorInsurerScheduleCSV reads the canonical CSV, persists the
// `PriorInsurerSchedule` envelope and N `PriorInsurerMember` rows. The
// uploaded CSV body is also stored on disk so the original document is
// retrievable for audit.
//
// quoteID is the GroupPricingQuote the schedule belongs to. The schedule
// header (insurer name, certificate number, dates) is supplied as form
// values rather than embedded in the CSV because brokers receive those
// separately from the member list.
func ImportPriorInsurerScheduleCSV(quoteID int, header PriorInsurerScheduleHeader, file multipart.File, fileName string, user models.AppUser) (*models.PriorInsurerSchedule, error) {
	if quoteID <= 0 {
		return nil, errors.New("invalid quote_id")
	}
	if file == nil {
		return nil, errors.New("file required")
	}

	// Buffer the body so we can both parse it and stash it on disk.
	bodyBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read upload: %w", err)
	}

	schedule := models.PriorInsurerSchedule{
		QuoteID:           quoteID,
		InsurerName:       strings.TrimSpace(header.InsurerName),
		CertificateNumber: strings.TrimSpace(header.CertificateNumber),
		EffectiveDate:     header.EffectiveDate,
		ExpiryDate:        header.ExpiryDate,
		Notes:             strings.TrimSpace(header.Notes),
		UploadedBy:        user.UserEmail,
	}
	if err := DB.Create(&schedule).Error; err != nil {
		return nil, fmt.Errorf("persist schedule: %w", err)
	}

	// Stash the raw CSV alongside the schedule row for audit.
	dir := filepath.Join(priorInsurerScheduleDir, fmt.Sprintf("schedule_%d", schedule.ID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create schedule dir: %w", err)
	}
	if fileName == "" {
		fileName = "schedule.csv"
	}
	destPath := filepath.Join(dir, filepath.Base(fileName))
	if err := os.WriteFile(destPath, bodyBytes, 0o644); err != nil {
		return nil, fmt.Errorf("save schedule file: %w", err)
	}
	schedule.DocumentPath = destPath

	// Parse the CSV into member rows.
	reader := csv.NewReader(bufio.NewReader(strings.NewReader(string(bodyBytes))))
	reader.FieldsPerRecord = -1
	headerRow, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}
	cols := indexHeader(headerRow)
	required := []string{"member_id_number", "member_name", "in_force"}
	for _, name := range required {
		if _, ok := cols[name]; !ok {
			return nil, fmt.Errorf("missing required column %q", name)
		}
	}

	memberCount := 0
	inForceCount := 0
	rowNum := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			return nil, fmt.Errorf("CSV row %d: %w", rowNum, err)
		}
		inForce := parseBoolTrue(field(row, cols, "in_force"))
		dob := parseDateOptional(field(row, cols, "date_of_birth"))
		loadings := parseBenefitMap(field(row, cols, "prior_loadings"))
		exclusions := parseExclusionList(field(row, cols, "prior_exclusions"))
		loadJSON, _ := json.Marshal(loadings)
		excJSON, _ := json.Marshal(exclusions)
		pm := models.PriorInsurerMember{
			ScheduleID:      schedule.ID,
			MemberIdNumber:  strings.TrimSpace(field(row, cols, "member_id_number")),
			MemberName:      strings.TrimSpace(field(row, cols, "member_name")),
			DateOfBirth:     dob,
			GlaSumAssured:   atofOrZero(field(row, cols, "gla_sum_assured")),
			PtdSumAssured:   atofOrZero(field(row, cols, "ptd_sum_assured")),
			CiSumAssured:    atofOrZero(field(row, cols, "ci_sum_assured")),
			PriorLoadings:   string(loadJSON),
			PriorExclusions: string(excJSON),
			InForce:         inForce,
			TakeoverOutcome: TakeoverOutcomeUnmatched, // overwritten by MatchPriorMembersToCensus
		}
		if err := DB.Create(&pm).Error; err != nil {
			return nil, fmt.Errorf("persist prior member on row %d: %w", rowNum, err)
		}
		memberCount++
		if inForce {
			inForceCount++
		}
	}

	schedule.MemberCount = memberCount
	schedule.InForceCount = inForceCount
	if err := DB.Save(&schedule).Error; err != nil {
		return nil, fmt.Errorf("update schedule counts: %w", err)
	}
	return &schedule, nil
}

// PriorInsurerScheduleHeader is the metadata that arrives alongside the CSV
// body (form values, not CSV columns).
type PriorInsurerScheduleHeader struct {
	InsurerName       string
	CertificateNumber string
	EffectiveDate     *time.Time
	ExpiryDate        *time.Time
	Notes             string
}

// MatchPriorMembersToCensus pairs each PriorInsurerMember against the
// quote's current census (`GPricingMemberData` for the quote). Matching
// order:
//  1. Exact MemberIdNumber match (trimmed, case-insensitive).
//  2. Name + DateOfBirth fuzzy match (lowercased name, exact DOB).
//
// Match outcome columns on PriorInsurerMember are written back. A summary
// of counts is returned so the renderer can show a preview before the
// underwriter commits to "Apply terms".
func MatchPriorMembersToCensus(scheduleID int) (*TakeoverMatchSummary, error) {
	var schedule models.PriorInsurerSchedule
	if err := DB.Where("id = ?", scheduleID).First(&schedule).Error; err != nil {
		return nil, err
	}
	var priorMembers []models.PriorInsurerMember
	if err := DB.Where("schedule_id = ?", scheduleID).Find(&priorMembers).Error; err != nil {
		return nil, err
	}
	var census []models.GPricingMemberData
	if err := DB.Where("quote_id = ?", schedule.QuoteID).Find(&census).Error; err != nil {
		return nil, err
	}

	idIndex := make(map[string]models.GPricingMemberData, len(census))
	nameDobIndex := make(map[string]models.GPricingMemberData, len(census))
	for _, m := range census {
		if k := strings.ToLower(strings.TrimSpace(m.MemberIdNumber)); k != "" {
			idIndex[k] = m
		}
		nameDobIndex[memberNameDOBKey(m.MemberName, m.DateOfBirth)] = m
	}

	// Index cases by member name + category so we can populate
	// MatchedCaseID when a UnderwritingCase exists for the matched member.
	var cases []models.UnderwritingCase
	if err := DB.Where("quote_id = ?", schedule.QuoteID).Find(&cases).Error; err != nil {
		return nil, err
	}
	caseIndex := make(map[string]models.UnderwritingCase, len(cases))
	for _, c := range cases {
		caseIndex[strings.ToLower(c.MemberName)+"|"+strings.ToLower(c.Category)] = c
	}

	summary := &TakeoverMatchSummary{ScheduleID: scheduleID, Total: len(priorMembers)}
	for i := range priorMembers {
		pm := &priorMembers[i]
		var matched *models.GPricingMemberData
		// 1. Match by id-number.
		if k := strings.ToLower(strings.TrimSpace(pm.MemberIdNumber)); k != "" {
			if hit, ok := idIndex[k]; ok {
				matched = &hit
				summary.MatchedByID++
			}
		}
		// 2. Fallback: name + DOB.
		if matched == nil && pm.DateOfBirth != nil {
			if hit, ok := nameDobIndex[memberNameDOBKey(pm.MemberName, *pm.DateOfBirth)]; ok {
				matched = &hit
				summary.MatchedByNameAndDOB++
			}
		}
		if pm.InForce {
			summary.InForce++
		}

		if matched == nil {
			pm.MatchedMemberName = ""
			pm.MatchedCategory = ""
			pm.MatchedCaseID = 0
			pm.TakeoverOutcome = TakeoverOutcomeUnmatched
			summary.Unmatched++
		} else {
			pm.MatchedMemberName = matched.MemberName
			pm.MatchedCategory = matched.SchemeCategory
			caseKey := strings.ToLower(matched.MemberName) + "|" + strings.ToLower(matched.SchemeCategory)
			if c, ok := caseIndex[caseKey]; ok {
				pm.MatchedCaseID = c.ID
			}
			pm.TakeoverOutcome = defaultOutcomeFor(pm)
			switch pm.TakeoverOutcome {
			case TakeoverOutcomeContinuationNoEvidence:
				summary.ContinuationNoEvidence++
			case TakeoverOutcomeContinuationWithLoading:
				summary.ContinuationWithLoading++
			case TakeoverOutcomeNewEvidenceRequired:
				summary.NewEvidenceRequired++
			}
		}
		if err := DB.Save(pm).Error; err != nil {
			return nil, fmt.Errorf("save prior member %d: %w", pm.ID, err)
		}
	}
	return summary, nil
}

// defaultOutcomeFor is the fallback classifier used when no takeover rules
// match a prior member. Logic:
//   - InForce + no loadings + no exclusions → continuation_no_evidence
//   - InForce + has loadings or exclusions  → continuation_with_loading
//   - Not InForce                           → new_evidence_required
//
// ClassifyPriorMembersAgainstRules will overwrite this when a richer rule
// set is configured.
func defaultOutcomeFor(pm *models.PriorInsurerMember) string {
	if !pm.InForce {
		return TakeoverOutcomeNewEvidenceRequired
	}
	if pm.PriorLoadings != "" && pm.PriorLoadings != "{}" && pm.PriorLoadings != "null" {
		return TakeoverOutcomeContinuationWithLoading
	}
	if pm.PriorExclusions != "" && pm.PriorExclusions != "[]" && pm.PriorExclusions != "null" {
		return TakeoverOutcomeContinuationWithLoading
	}
	return TakeoverOutcomeContinuationNoEvidence
}

// ClassifyPriorMembersAgainstRules re-runs the rules engine (Phase 3)
// against every matched prior member in the schedule, replacing the
// default outcome with the engine's verdict when the engine returns a
// stricter one. Rules authored under `category = "takeover"` drive this.
func ClassifyPriorMembersAgainstRules(scheduleID int) error {
	var schedule models.PriorInsurerSchedule
	if err := DB.Where("id = ?", scheduleID).First(&schedule).Error; err != nil {
		return err
	}
	var priorMembers []models.PriorInsurerMember
	if err := DB.Where("schedule_id = ? AND takeover_outcome <> ?", scheduleID, TakeoverOutcomeUnmatched).Find(&priorMembers).Error; err != nil {
		return err
	}
	for i := range priorMembers {
		pm := &priorMembers[i]
		ctx := buildPriorMemberContext(*pm)
		summary, err := EvaluateAgainstActiveRuleSet(ctx)
		if err != nil {
			continue // engine failure is non-fatal; the default outcome stays.
		}
		if summary.Outcome == "" {
			continue
		}
		// Engine's outcome strings are opaque labels; treat anything
		// starting with `continuation_` or equal to a known takeover
		// outcome as authoritative. Otherwise fall back to default.
		switch summary.Outcome {
		case TakeoverOutcomeContinuationNoEvidence,
			TakeoverOutcomeContinuationWithLoading,
			TakeoverOutcomeNewEvidenceRequired:
			pm.TakeoverOutcome = summary.Outcome
		}
		if err := DB.Save(pm).Error; err != nil {
			return err
		}
	}
	return nil
}

// buildPriorMemberContext exposes prior-member fields to the rules engine
// under stable keys: `in_force`, `prior_gla_sa`, `prior_ptd_sa`,
// `prior_ci_sa`, `prior_loading_<benefit>` per benefit, and
// `prior_exclusion_<code> = true` per exclusion code.
func buildPriorMemberContext(pm models.PriorInsurerMember) EvaluationContext {
	ctx := EvaluationContext{
		"in_force":     pm.InForce,
		"prior_gla_sa": pm.GlaSumAssured,
		"prior_ptd_sa": pm.PtdSumAssured,
		"prior_ci_sa":  pm.CiSumAssured,
	}
	if pm.PriorLoadings != "" {
		var loadings map[string]float64
		if err := json.Unmarshal([]byte(pm.PriorLoadings), &loadings); err == nil {
			for k, v := range loadings {
				ctx["prior_loading_"+strings.ToLower(k)] = v
			}
		}
	}
	if pm.PriorExclusions != "" {
		var exclusions []string
		if err := json.Unmarshal([]byte(pm.PriorExclusions), &exclusions); err == nil {
			for _, code := range exclusions {
				code = strings.TrimSpace(code)
				if code != "" {
					ctx["prior_exclusion_"+strings.ToLower(code)] = true
				}
			}
		}
	}
	return ctx
}

// ApplyTakeoverTermsToCases pushes the matched prior member's outcome onto
// the corresponding UnderwritingCase as an engine snapshot. Suggest-don't-
// decide: the underwriter still commits the human decision via the existing
// form. Returns the number of cases touched.
func ApplyTakeoverTermsToCases(scheduleID int, user models.AppUser) (int, error) {
	var priorMembers []models.PriorInsurerMember
	if err := DB.Where("schedule_id = ? AND matched_case_id > 0", scheduleID).Find(&priorMembers).Error; err != nil {
		return 0, err
	}
	now := time.Now()
	touched := 0
	for _, pm := range priorMembers {
		var c models.UnderwritingCase
		if err := DB.Where("id = ?", pm.MatchedCaseID).First(&c).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return touched, err
		}
		c.EngineOutcome = pm.TakeoverOutcome
		c.EngineEvaluatedAt = &now
		// Surface prior loadings as a hint on EngineLoading: the max
		// numeric value in the prior_loadings JSON map.
		c.EngineLoading = maxLoadingFromPrior(pm.PriorLoadings)
		if err := DB.Save(&c).Error; err != nil {
			return touched, fmt.Errorf("save case %d: %w", c.ID, err)
		}
		recordCaseEvent(c.ID, "takeover_terms_applied", user.UserEmail, map[string]any{
			"schedule_id":     scheduleID,
			"takeover_outcome": pm.TakeoverOutcome,
			"prior_loadings":  json.RawMessage(pm.PriorLoadings),
			"prior_exclusions": json.RawMessage(pm.PriorExclusions),
		})
		touched++
	}
	return touched, nil
}

func maxLoadingFromPrior(raw string) float64 {
	if raw == "" {
		return 0
	}
	var loadings map[string]float64
	if err := json.Unmarshal([]byte(raw), &loadings); err != nil {
		return 0
	}
	var best float64
	for _, v := range loadings {
		if v > best {
			best = v
		}
	}
	return best
}

// GetPriorInsurerScheduleForQuote returns the most recent schedule (if any)
// for the given quote, with members preloaded. Used by the renderer to
// drive the takeover banner.
func GetPriorInsurerScheduleForQuote(quoteID int) (*models.PriorInsurerSchedule, error) {
	var schedule models.PriorInsurerSchedule
	err := DB.Preload("Members").
		Where("quote_id = ?", quoteID).
		Order("uploaded_at DESC").
		First(&schedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// ───── Helpers ─────────────────────────────────────────────────────────────

func memberNameDOBKey(name string, dob time.Time) string {
	return strings.ToLower(strings.TrimSpace(name)) + "|" + dob.Format("2006-01-02")
}

func parseBoolTrue(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "t", "yes", "y", "1", "in_force":
		return true
	}
	return false
}

func parseDateOptional(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02", "2006/01/02", "02/01/2006", "02-01-2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}

// parseBenefitMap accepts `gla:25|ptd:10` and returns {"gla": 25, "ptd": 10}.
func parseBenefitMap(s string) map[string]float64 {
	out := make(map[string]float64)
	for _, part := range strings.Split(s, "|") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := atofOrZero(strings.TrimSpace(kv[1]))
		if key != "" {
			out[key] = val
		}
	}
	return out
}

func parseExclusionList(s string) []string {
	out := make([]string, 0)
	for _, part := range strings.Split(s, "|") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, strings.ToLower(part))
		}
	}
	return out
}

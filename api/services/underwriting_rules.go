package services

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"api/models"
)

// Supported operators in UWRule.Op. The set is deliberately closed; adding a
// new one means a new const here + a new branch in evalRule.
const (
	UWOpEq      = "eq"
	UWOpNe      = "ne"
	UWOpGt      = "gt"
	UWOpGte     = "gte"
	UWOpLt      = "lt"
	UWOpLte     = "lte"
	UWOpBetween = "between"
	UWOpIn      = "in"
)

// UWRuleOutcome mirrors UWDecisionOutcome but is what the engine produces
// before a human commits to a UnderwritingDecision. `refer` means "needs
// human review" — there's no equivalent enum value on the decision side, by
// design: a decision is always a committed accept / postpone / decline.
const (
	UWRuleOutcomeAccept  = "accept"
	UWRuleOutcomeRefer   = "refer"
	UWRuleOutcomeDecline = "decline"
)

// EvaluationContext is the per-member input to the engine. Keys are stable
// field names (`bmi`, `age`, `smoker`, `occupation_class`, condition codes,
// etc.). Values are plain Go scalars or slices — the engine compares them
// via reflection. Callers populate this from MemberDisclosure (Phase 5) or
// per-member rating result fields.
type EvaluationContext map[string]any

// EvaluatedOutcome is one match. The engine returns every matching rule;
// the caller aggregates how it likes.
type EvaluatedOutcome struct {
	RuleID         int     `json:"rule_id"`
	Category       string  `json:"category"`
	Field          string  `json:"field"`
	Outcome        string  `json:"outcome"`
	LoadingPercent float64 `json:"loading_percent"`
	ExclusionCode  string  `json:"exclusion_code"`
	Priority       int     `json:"priority"`
	Notes          string  `json:"notes"`
}

// EvaluationSummary is the aggregated view: strictest outcome (decline >
// refer > accept), max loading across matching rules, and the union of
// exclusion codes. Convenient for UI display and case-decision pre-population.
type EvaluationSummary struct {
	Outcome        string             `json:"outcome"`
	MaxLoading     float64            `json:"max_loading"`
	Exclusions     []string           `json:"exclusions"`
	Outcomes       []EvaluatedOutcome `json:"outcomes"`
	RuleSetID      int                `json:"rule_set_id"`
	RuleSetVersion int                `json:"rule_set_version"`
}

// uwCondition is the canonical decoded form of UWRule.ConditionJSON. We
// accept either {value_a, value_b} (for eq/ne/gt/.../between) or {values}
// (for `in`). Unknown extra keys are ignored.
type uwCondition struct {
	ValueA json.RawMessage   `json:"value_a"`
	ValueB json.RawMessage   `json:"value_b"`
	Values []json.RawMessage `json:"values"`
}

// Evaluate is pure — no DB, no side effects. Iterates rules in (priority asc,
// id asc) order; returns every match. The caller decides aggregation.
func Evaluate(rules []models.UWRule, ctx EvaluationContext) []EvaluatedOutcome {
	sorted := make([]models.UWRule, len(rules))
	copy(sorted, rules)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Priority != sorted[j].Priority {
			return sorted[i].Priority < sorted[j].Priority
		}
		return sorted[i].ID < sorted[j].ID
	})

	out := make([]EvaluatedOutcome, 0, len(sorted))
	for _, r := range sorted {
		val, present := ctx[r.Field]
		if !present {
			continue
		}
		ok, err := evalRule(r, val)
		if err != nil || !ok {
			continue
		}
		out = append(out, EvaluatedOutcome{
			RuleID:         r.ID,
			Category:       r.Category,
			Field:          r.Field,
			Outcome:        r.Outcome,
			LoadingPercent: r.LoadingPercent,
			ExclusionCode:  r.ExclusionCode,
			Priority:       r.Priority,
			Notes:          r.Notes,
		})
	}
	return out
}

// EvaluateAgainstActiveRuleSet loads the currently-active UWRuleSet (the
// row with Active=true; latest Version wins on tie) and evaluates the
// context against it. Returns (summary, ruleSetID, ruleSetVersion, err).
// Returns an empty summary and zero IDs when no active set exists — the
// caller treats that as "no rules configured yet".
func EvaluateAgainstActiveRuleSet(ctx EvaluationContext) (EvaluationSummary, error) {
	rs, rules, err := loadActiveRuleSet()
	if err != nil {
		return EvaluationSummary{}, err
	}
	if rs == nil {
		return EvaluationSummary{}, nil
	}
	outcomes := Evaluate(rules, ctx)
	summary := summarise(outcomes)
	summary.RuleSetID = rs.ID
	summary.RuleSetVersion = rs.Version
	return summary, nil
}

// EvaluateAgainstRuleSet loads a specific rule set by ID (typically the one
// snapshotted on a UnderwritingCase) and evaluates the context against it.
func EvaluateAgainstRuleSet(ruleSetID int, ctx EvaluationContext) (EvaluationSummary, error) {
	if ruleSetID <= 0 {
		return EvaluateAgainstActiveRuleSet(ctx)
	}
	var rs models.UWRuleSet
	if err := DB.Where("id = ?", ruleSetID).First(&rs).Error; err != nil {
		return EvaluationSummary{}, err
	}
	var rules []models.UWRule
	if err := DB.Where("rule_set_id = ?", rs.ID).Find(&rules).Error; err != nil {
		return EvaluationSummary{}, err
	}
	outcomes := Evaluate(rules, ctx)
	summary := summarise(outcomes)
	summary.RuleSetID = rs.ID
	summary.RuleSetVersion = rs.Version
	return summary, nil
}

func loadActiveRuleSet() (*models.UWRuleSet, []models.UWRule, error) {
	var rs models.UWRuleSet
	err := DB.Where("active = ?", true).Order("version DESC").First(&rs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	var rules []models.UWRule
	if err := DB.Where("rule_set_id = ?", rs.ID).Find(&rules).Error; err != nil {
		return nil, nil, err
	}
	return &rs, rules, nil
}

func summarise(outcomes []EvaluatedOutcome) EvaluationSummary {
	s := EvaluationSummary{Outcomes: outcomes, Outcome: UWRuleOutcomeAccept}
	seenExclusion := make(map[string]bool)
	for _, o := range outcomes {
		if o.LoadingPercent > s.MaxLoading {
			s.MaxLoading = o.LoadingPercent
		}
		if o.ExclusionCode != "" && !seenExclusion[o.ExclusionCode] {
			seenExclusion[o.ExclusionCode] = true
			s.Exclusions = append(s.Exclusions, o.ExclusionCode)
		}
		s.Outcome = strictestOutcome(s.Outcome, o.Outcome)
	}
	return s
}

func strictestOutcome(current, next string) string {
	rank := map[string]int{
		UWRuleOutcomeAccept:  0,
		UWRuleOutcomeRefer:   1,
		UWRuleOutcomeDecline: 2,
	}
	if rank[next] > rank[current] {
		return next
	}
	return current
}

// evalRule tests one rule against one field value. Errors signal "rule is
// malformed" — those rules are skipped in Evaluate rather than crashing.
func evalRule(r models.UWRule, raw any) (bool, error) {
	if r.ConditionJSON == "" {
		return false, errors.New("empty condition")
	}
	var cond uwCondition
	if err := json.Unmarshal([]byte(r.ConditionJSON), &cond); err != nil {
		return false, fmt.Errorf("decode condition: %w", err)
	}
	switch r.Op {
	case UWOpEq, UWOpNe:
		want, err := decodeScalar(cond.ValueA)
		if err != nil {
			return false, err
		}
		equal := scalarsEqual(raw, want)
		if r.Op == UWOpNe {
			return !equal, nil
		}
		return equal, nil
	case UWOpGt, UWOpGte, UWOpLt, UWOpLte:
		want, err := jsonToFloat(cond.ValueA)
		if err != nil {
			return false, err
		}
		got, err := anyToFloat(raw)
		if err != nil {
			return false, err
		}
		switch r.Op {
		case UWOpGt:
			return got > want, nil
		case UWOpGte:
			return got >= want, nil
		case UWOpLt:
			return got < want, nil
		case UWOpLte:
			return got <= want, nil
		}
	case UWOpBetween:
		lo, err := jsonToFloat(cond.ValueA)
		if err != nil {
			return false, err
		}
		hi, err := jsonToFloat(cond.ValueB)
		if err != nil {
			return false, err
		}
		got, err := anyToFloat(raw)
		if err != nil {
			return false, err
		}
		return got >= lo && got <= hi, nil
	case UWOpIn:
		for _, v := range cond.Values {
			candidate, err := decodeScalar(v)
			if err != nil {
				continue
			}
			if scalarsEqual(raw, candidate) {
				return true, nil
			}
		}
		return false, nil
	}
	return false, fmt.Errorf("unknown op %q", r.Op)
}

// decodeScalar reads a json.RawMessage into a Go scalar. We support
// strings, numbers and booleans — enough for the closed operator set.
func decodeScalar(raw json.RawMessage) (any, error) {
	if len(raw) == 0 {
		return nil, errors.New("missing value")
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s, nil
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return b, nil
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("unsupported value %s", string(raw))
}

func jsonToFloat(raw json.RawMessage) (float64, error) {
	if len(raw) == 0 {
		return 0, errors.New("missing numeric value")
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return f, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strconv.ParseFloat(s, 64)
	}
	return 0, fmt.Errorf("not numeric: %s", string(raw))
}

func anyToFloat(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case json.Number:
		return x.Float64()
	case string:
		return strconv.ParseFloat(x, 64)
	}
	return 0, fmt.Errorf("not numeric: %T", v)
}

func scalarsEqual(a, b any) bool {
	switch av := a.(type) {
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case bool:
		bv, ok := b.(bool)
		return ok && av == bv
	}
	// Numeric path: compare as float64.
	af, errA := anyToFloat(a)
	bf, errB := anyToFloat(b)
	if errA != nil || errB != nil {
		return false
	}
	return af == bf
}

// BuildRuleContextForCase assembles the EvaluationContext for a case dry-run.
// Per-member rating-derived fields (sa, fcl, ratio, tier) are always present.
// Underwriter-supplied hypothetical values via query string override / add
// fields: `?bmi=32&smoker=true&occupation_class=3`. Phase 5 will replace the
// query-string override with persisted MemberDisclosure rows.
func BuildRuleContextForCase(c models.UnderwritingCase, query map[string][]string) EvaluationContext {
	ctx := EvaluationContext{
		"gla_sum_assured":        c.GlaSumAssured,
		"ptd_sum_assured":        c.PtdSumAssured,
		"ci_sum_assured":         c.CiSumAssured,
		"spouse_gla_sum_assured": c.SpouseGlaSumAssured,
		"free_cover_limit":       c.FreeCoverLimit,
		"fcl_excess_ratio":       c.FCLExcessRatio,
		"tier":                   float64(c.Tier),
	}
	for k, vals := range query {
		if len(vals) == 0 {
			continue
		}
		ctx[k] = csvValueToScalar(vals[0])
	}
	return ctx
}

// ----- CRUD helpers used by controllers -----------------------------------

// ListUWRuleSets returns rule sets ordered by version descending so the most
// recent set surfaces first in the admin UI.
func ListUWRuleSets() ([]models.UWRuleSet, error) {
	var sets []models.UWRuleSet
	if err := DB.Order("active DESC, version DESC, id DESC").Find(&sets).Error; err != nil {
		return nil, err
	}
	return sets, nil
}

// GetUWRuleSet returns a rule set with its rules preloaded.
func GetUWRuleSet(id int) (*models.UWRuleSet, error) {
	var rs models.UWRuleSet
	if err := DB.Preload("Rules").Where("id = ?", id).First(&rs).Error; err != nil {
		return nil, err
	}
	return &rs, nil
}

// CreateUWRuleSet persists a new (empty) rule set.
func CreateUWRuleSet(rs models.UWRuleSet) (*models.UWRuleSet, error) {
	if err := DB.Create(&rs).Error; err != nil {
		return nil, err
	}
	return &rs, nil
}

// ActivateUWRuleSet sets `active=true` on the given set and `active=false`
// on every other set. Only one rule set can be active at a time.
func ActivateUWRuleSet(id int) (*models.UWRuleSet, error) {
	return nil, DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.UWRuleSet{}).Where("id <> ?", id).Update("active", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.UWRuleSet{}).Where("id = ?", id).Update("active", true).Error; err != nil {
			return err
		}
		return nil
	})
}

// CreateUWRule appends a rule to a rule set.
func CreateUWRule(ruleSetID int, r models.UWRule) (*models.UWRule, error) {
	r.RuleSetID = ruleSetID
	if r.ConditionJSON == "" {
		body, err := buildConditionJSON(r)
		if err != nil {
			return nil, err
		}
		r.ConditionJSON = body
	}
	if err := DB.Create(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateUWRule replaces a rule's mutable fields.
func UpdateUWRule(id int, patch models.UWRule) (*models.UWRule, error) {
	var existing models.UWRule
	if err := DB.Where("id = ?", id).First(&existing).Error; err != nil {
		return nil, err
	}
	existing.Category = patch.Category
	existing.Field = patch.Field
	existing.Op = patch.Op
	existing.Outcome = patch.Outcome
	existing.LoadingPercent = patch.LoadingPercent
	existing.ExclusionCode = patch.ExclusionCode
	existing.Priority = patch.Priority
	existing.Notes = patch.Notes
	if patch.ConditionJSON != "" {
		existing.ConditionJSON = patch.ConditionJSON
	} else {
		body, err := buildConditionJSON(patch)
		if err == nil && body != "" {
			existing.ConditionJSON = body
		}
	}
	if err := DB.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

// DeleteUWRule removes a rule by ID.
func DeleteUWRule(id int) error {
	return DB.Where("id = ?", id).Delete(&models.UWRule{}).Error
}

// ExportUWRuleSetCSV returns the canonical CSV bytes for the given rule
// set — the same shape ImportUWRulesCSV consumes, so an export → import
// round-trip produces an equivalent set.
//
// Column order matches the importer header: rule_set_name, version,
// category, field, op, value_a, value_b, values, outcome,
// loading_percent, exclusion_code, priority, notes.
func ExportUWRuleSetCSV(ruleSetID int) ([]byte, string, error) {
	var rs models.UWRuleSet
	if err := DB.Preload("Rules").Where("id = ?", ruleSetID).First(&rs).Error; err != nil {
		return nil, "", err
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	header := []string{
		"rule_set_name", "version", "category", "field", "op",
		"value_a", "value_b", "values", "outcome",
		"loading_percent", "exclusion_code", "priority", "notes",
	}
	if err := w.Write(header); err != nil {
		return nil, "", err
	}
	for _, r := range rs.Rules {
		va, vb, vals := splitConditionForCSV(r.ConditionJSON)
		row := []string{
			rs.Name,
			strconv.Itoa(rs.Version),
			r.Category,
			r.Field,
			r.Op,
			va,
			vb,
			vals,
			r.Outcome,
			strconv.FormatFloat(r.LoadingPercent, 'f', -1, 64),
			r.ExclusionCode,
			strconv.Itoa(r.Priority),
			r.Notes,
		}
		if err := w.Write(row); err != nil {
			return nil, "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, "", err
	}
	fileName := fmt.Sprintf("%s_v%d.csv", sanitiseForFilename(rs.Name), rs.Version)
	return buf.Bytes(), fileName, nil
}

// splitConditionForCSV reverses buildConditionFromCSV. Returns the CSV
// columns (value_a, value_b, values) for the given ConditionJSON.
// Unknown / malformed JSON falls through as empty columns — the export
// is still well-formed CSV, just missing the condition payload.
func splitConditionForCSV(condJSON string) (valueA, valueB, values string) {
	if condJSON == "" {
		return
	}
	var c uwCondition
	if err := json.Unmarshal([]byte(condJSON), &c); err != nil {
		return
	}
	scalar := func(raw json.RawMessage) string {
		if len(raw) == 0 {
			return ""
		}
		v, err := decodeScalar(raw)
		if err != nil {
			return ""
		}
		switch t := v.(type) {
		case string:
			return t
		case bool:
			if t {
				return "true"
			}
			return "false"
		case float64:
			return strconv.FormatFloat(t, 'f', -1, 64)
		}
		return ""
	}
	valueA = scalar(c.ValueA)
	valueB = scalar(c.ValueB)
	if len(c.Values) > 0 {
		parts := make([]string, 0, len(c.Values))
		for _, raw := range c.Values {
			parts = append(parts, scalar(raw))
		}
		values = strings.Join(parts, "|")
	}
	return
}

// sanitiseForFilename replaces characters disallowed in filenames so the
// HTTP Content-Disposition header is portable.
func sanitiseForFilename(s string) string {
	if s == "" {
		return "rule_set"
	}
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '_', r == '-':
			out = append(out, r)
		case r == ' ':
			out = append(out, '_')
		}
	}
	if len(out) == 0 {
		return "rule_set"
	}
	return string(out)
}

// DuplicateUWRuleSet creates a new rule set with the same name as the
// source, version bumped to max(version) + 1, and every rule copied
// across. The new set is created inactive — the caller chooses when to
// activate.
func DuplicateUWRuleSet(sourceID int, actor string) (*models.UWRuleSet, error) {
	var source models.UWRuleSet
	if err := DB.Preload("Rules").Where("id = ?", sourceID).First(&source).Error; err != nil {
		return nil, err
	}
	var maxVersion int
	if err := DB.Model(&models.UWRuleSet{}).
		Where("name = ?", source.Name).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error; err != nil {
		return nil, fmt.Errorf("lookup max version: %w", err)
	}
	nextVersion := maxVersion + 1

	var newSet models.UWRuleSet
	err := DB.Transaction(func(tx *gorm.DB) error {
		newSet = models.UWRuleSet{
			Name:      source.Name,
			Version:   nextVersion,
			Active:    false,
			CreatedBy: actor,
		}
		if err := tx.Create(&newSet).Error; err != nil {
			return fmt.Errorf("create duplicate: %w", err)
		}
		for _, r := range source.Rules {
			copied := models.UWRule{
				RuleSetID:      newSet.ID,
				Category:       r.Category,
				Field:          r.Field,
				Op:             r.Op,
				ConditionJSON:  r.ConditionJSON,
				Outcome:        r.Outcome,
				LoadingPercent: r.LoadingPercent,
				ExclusionCode:  r.ExclusionCode,
				Priority:       r.Priority,
				Notes:          r.Notes,
			}
			if err := tx.Create(&copied).Error; err != nil {
				return fmt.Errorf("copy rule %d: %w", r.ID, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Preload rules on the returned object so the renderer can show them
	// immediately without a second fetch.
	if err := DB.Preload("Rules").Where("id = ?", newSet.ID).First(&newSet).Error; err != nil {
		return nil, err
	}
	return &newSet, nil
}

// DeleteUWRuleSet removes a rule set and all its rules — but only when
// it is safe to do so. Refuses when:
//   - the set is currently active (would orphan in-flight cases),
//   - any UnderwritingCase has snapshotted this set's id (historical
//     reproducibility — never delete a set a case still references).
func DeleteUWRuleSet(id int) error {
	var rs models.UWRuleSet
	if err := DB.Where("id = ?", id).First(&rs).Error; err != nil {
		return err
	}
	if rs.Active {
		return fmt.Errorf("cannot delete an active rule set — activate a different set first")
	}
	var caseCount int64
	if err := DB.Model(&models.UnderwritingCase{}).Where("rule_set_id = ?", id).Count(&caseCount).Error; err != nil {
		return fmt.Errorf("count referencing cases: %w", err)
	}
	if caseCount > 0 {
		return fmt.Errorf("cannot delete: %d underwriting case(s) snapshotted this set; preserve history", caseCount)
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("rule_set_id = ?", id).Delete(&models.UWRule{}).Error; err != nil {
			return fmt.Errorf("delete rules: %w", err)
		}
		if err := tx.Where("id = ?", id).Delete(&models.UWRuleSet{}).Error; err != nil {
			return fmt.Errorf("delete rule set: %w", err)
		}
		return nil
	})
}

// SeedStarterRuleSet idempotently creates a fixture rule set with a
// handful of common underwriting rules so first-time admins have
// something to iterate from. Re-runs return the existing set unchanged.
//
// Rules covered:
//   - build: BMI 30–35 → refer +25% (obese_class1)
//   - build: BMI 35–40 → refer +50% (obese_class2)
//   - build: BMI > 40 → decline (morbid_obesity)
//   - lifestyle: smoker = true → refer +15% (smoker)
//   - occupation: occupation_class in [3, 4] → refer +20% (hazardous_occ)
func SeedStarterRuleSet(actor string) (*models.UWRuleSet, error) {
	const starterName = "Standard underwriting"
	const starterVersion = 1

	var existing models.UWRuleSet
	err := DB.Preload("Rules").Where("name = ? AND version = ?", starterName, starterVersion).First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("lookup existing starter: %w", err)
	}

	type seed struct {
		Category       string
		Field          string
		Op             string
		Condition      string
		Outcome        string
		Loading        float64
		ExclusionCode  string
		Priority       int
		Notes          string
	}
	starter := []seed{
		{"build", "bmi", UWOpBetween, `{"value_a":30,"value_b":35}`, UWRuleOutcomeRefer, 25, "obese_class1", 10, "BMI 30–35: refer with +25%"},
		{"build", "bmi", UWOpBetween, `{"value_a":35,"value_b":40}`, UWRuleOutcomeRefer, 50, "obese_class2", 20, "BMI 35–40: refer with +50%"},
		{"build", "bmi", UWOpGt, `{"value_a":40}`, UWRuleOutcomeDecline, 0, "morbid_obesity", 30, "BMI > 40: decline"},
		{"lifestyle", "smoker", UWOpEq, `{"value_a":true}`, UWRuleOutcomeRefer, 15, "smoker", 40, "Disclosed smoker: +15%"},
		{"occupation", "occupation_class", UWOpIn, `{"values":[3,4]}`, UWRuleOutcomeRefer, 20, "hazardous_occ", 50, "Hazardous occupation classes 3-4: refer +20%"},
	}

	var created models.UWRuleSet
	err = DB.Transaction(func(tx *gorm.DB) error {
		created = models.UWRuleSet{
			Name:      starterName,
			Version:   starterVersion,
			Active:    false,
			CreatedBy: actor,
		}
		if err := tx.Create(&created).Error; err != nil {
			return fmt.Errorf("create starter set: %w", err)
		}
		for _, s := range starter {
			rule := models.UWRule{
				RuleSetID:      created.ID,
				Category:       s.Category,
				Field:          s.Field,
				Op:             s.Op,
				ConditionJSON:  s.Condition,
				Outcome:        s.Outcome,
				LoadingPercent: s.Loading,
				ExclusionCode:  s.ExclusionCode,
				Priority:       s.Priority,
				Notes:          s.Notes,
			}
			if err := tx.Create(&rule).Error; err != nil {
				return fmt.Errorf("create starter rule %s: %w", s.ExclusionCode, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := DB.Preload("Rules").Where("id = ?", created.ID).First(&created).Error; err != nil {
		return nil, err
	}
	return &created, nil
}

// ----- CSV bulk-import ---------------------------------------------------

// ImportUWRulesCSV reads a CSV with the following columns (header required):
//
//	rule_set_name, version, category, field, op, value_a, value_b, values,
//	outcome, loading_percent, exclusion_code, priority, notes
//
// `values` is pipe-separated for the `in` operator (e.g. `3|4`). A new rule
// set is created if (name, version) doesn't exist; rules are appended to the
// matching set. Returns the number of rules imported.
func ImportUWRulesCSV(file multipart.File, createdBy string) (int, error) {
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("read header: %w", err)
	}
	cols := indexHeader(header)
	required := []string{"rule_set_name", "category", "field", "op", "outcome"}
	for _, name := range required {
		if _, ok := cols[name]; !ok {
			return 0, fmt.Errorf("missing required column %q", name)
		}
	}

	count := 0
	setCache := map[string]int{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, fmt.Errorf("read row: %w", err)
		}
		setName := strings.TrimSpace(field(row, cols, "rule_set_name"))
		if setName == "" {
			continue
		}
		version := atoiOrZero(field(row, cols, "version"))
		key := fmt.Sprintf("%s|%d", setName, version)
		rsID, ok := setCache[key]
		if !ok {
			var rs models.UWRuleSet
			err := DB.Where("name = ? AND version = ?", setName, version).First(&rs).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rs = models.UWRuleSet{Name: setName, Version: version, CreatedBy: createdBy}
				if err := DB.Create(&rs).Error; err != nil {
					return count, fmt.Errorf("create rule set %q: %w", setName, err)
				}
			} else if err != nil {
				return count, fmt.Errorf("lookup rule set %q: %w", setName, err)
			}
			rsID = rs.ID
			setCache[key] = rsID
		}

		condBody, err := buildConditionFromCSV(
			field(row, cols, "op"),
			field(row, cols, "value_a"),
			field(row, cols, "value_b"),
			field(row, cols, "values"),
		)
		if err != nil {
			return count, fmt.Errorf("row %d: %w", count+2, err)
		}
		rule := models.UWRule{
			RuleSetID:      rsID,
			Category:       strings.TrimSpace(field(row, cols, "category")),
			Field:          strings.TrimSpace(field(row, cols, "field")),
			Op:             strings.TrimSpace(field(row, cols, "op")),
			ConditionJSON:  condBody,
			Outcome:        strings.TrimSpace(field(row, cols, "outcome")),
			LoadingPercent: atofOrZero(field(row, cols, "loading_percent")),
			ExclusionCode:  strings.TrimSpace(field(row, cols, "exclusion_code")),
			Priority:       atoiOrZero(field(row, cols, "priority")),
			Notes:          strings.TrimSpace(field(row, cols, "notes")),
		}
		if err := DB.Create(&rule).Error; err != nil {
			return count, fmt.Errorf("create rule on row %d: %w", count+2, err)
		}
		count++
	}
	return count, nil
}

func indexHeader(row []string) map[string]int {
	out := make(map[string]int, len(row))
	for i, h := range row {
		out[strings.ToLower(strings.TrimSpace(h))] = i
	}
	return out
}

func field(row []string, cols map[string]int, name string) string {
	idx, ok := cols[name]
	if !ok || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func atoiOrZero(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

func atofOrZero(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

// buildConditionFromCSV assembles the canonical ConditionJSON from the CSV
// columns. For `in`, `values` is pipe-separated; everything else uses
// value_a (and optionally value_b for between).
func buildConditionFromCSV(op, valueA, valueB, valuesRaw string) (string, error) {
	cond := map[string]any{}
	switch op {
	case UWOpIn:
		parts := strings.Split(valuesRaw, "|")
		out := make([]any, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			out = append(out, csvValueToScalar(p))
		}
		if len(out) == 0 {
			return "", errors.New("`in` requires non-empty values column")
		}
		cond["values"] = out
	case UWOpBetween:
		if valueA == "" || valueB == "" {
			return "", errors.New("`between` requires value_a and value_b")
		}
		cond["value_a"] = csvValueToScalar(strings.TrimSpace(valueA))
		cond["value_b"] = csvValueToScalar(strings.TrimSpace(valueB))
	default:
		if valueA == "" {
			return "", fmt.Errorf("op %q requires value_a", op)
		}
		cond["value_a"] = csvValueToScalar(strings.TrimSpace(valueA))
	}
	body, err := json.Marshal(cond)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// csvValueToScalar promotes a CSV string into the most specific Go scalar:
// booleans first, then numbers, then string.
func csvValueToScalar(s string) any {
	if strings.EqualFold(s, "true") {
		return true
	}
	if strings.EqualFold(s, "false") {
		return false
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

// buildConditionJSON is the API-side equivalent of buildConditionFromCSV: it
// reads the rule's already-set fields (op + Notes/etc. don't carry values
// directly, so callers POST ConditionJSON pre-assembled). This helper is the
// fall-through used when ConditionJSON is empty on create/update.
func buildConditionJSON(_ models.UWRule) (string, error) {
	return "", nil
}

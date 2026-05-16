package services

import (
	"api/models"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"gonum.org/v1/gonum/optimize"
	"gorm.io/gorm"
)

// Risk-level bands used everywhere (DB, frontend, rules).
const (
	RiskLevelLow      = "low"
	RiskLevelMedium   = "medium"
	RiskLevelHigh     = "high"
	RiskLevelCritical = "critical"
)

// validRiskLevels gates rule input from admins.
var validRiskLevels = map[string]bool{
	RiskLevelLow:      true,
	RiskLevelMedium:   true,
	RiskLevelHigh:     true,
	RiskLevelCritical: true,
}

// FraudFeatureKind classifies a feature so the rule editor can surface the right
// input control (numeric stepper vs string dropdown).
type FraudFeatureKind string

const (
	FeatureKindNumeric FraudFeatureKind = "numeric"
	FeatureKindString  FraudFeatureKind = "string"
	FeatureKindBool    FraudFeatureKind = "bool"
)

// FraudFeatureSpec describes a single feature used by both the GLM and the
// rule engine. Categorical features live twice: a raw string form (for rules)
// and one-hot 0/1 forms (for the GLM).
type FraudFeatureSpec struct {
	Name        string           `json:"name"`
	Kind        FraudFeatureKind `json:"kind"`
	Description string           `json:"description"`
	UsedByGLM   bool             `json:"used_by_glm"`
	UsedByRules bool             `json:"used_by_rules"`
	Choices     []string         `json:"choices,omitempty"`
}

// FraudFeatureCatalogue is the single source of truth for the feature
// vocabulary. Service extraction and the rule editor read from this list.
var FraudFeatureCatalogue = []FraudFeatureSpec{
	{Name: "claim_amount", Kind: FeatureKindNumeric, Description: "Claim amount (currency)", UsedByRules: true},
	{Name: "claim_amount_log", Kind: FeatureKindNumeric, Description: "log1p(claim_amount)", UsedByGLM: true},
	{Name: "days_event_to_notification", Kind: FeatureKindNumeric, Description: "Days from event to notification", UsedByGLM: true, UsedByRules: true},
	{Name: "days_notification_to_registration", Kind: FeatureKindNumeric, Description: "Days from notification to registration", UsedByGLM: true, UsedByRules: true},
	{Name: "attachment_count", Kind: FeatureKindNumeric, Description: "Number of attached documents", UsedByGLM: true, UsedByRules: true},
	{Name: "member_type", Kind: FeatureKindString, Description: "member|spouse|child|other", UsedByRules: true, Choices: []string{"member", "spouse", "child", "other"}},
	{Name: "member_type_member", Kind: FeatureKindBool, Description: "One-hot member_type=member", UsedByGLM: true},
	{Name: "member_type_spouse", Kind: FeatureKindBool, Description: "One-hot member_type=spouse", UsedByGLM: true},
	{Name: "member_type_child", Kind: FeatureKindBool, Description: "One-hot member_type=child", UsedByGLM: true},
	{Name: "relationship_to_member", Kind: FeatureKindString, Description: "Claimant's relationship to member", UsedByRules: true},
	{Name: "relationship_self", Kind: FeatureKindBool, Description: "Claimant is the member (self)", UsedByGLM: true, UsedByRules: true},
	{Name: "cause_type", Kind: FeatureKindString, Description: "Cause of claim", UsedByRules: true},
	{Name: "cause_natural", Kind: FeatureKindBool, Description: "Cause is natural causes", UsedByGLM: true, UsedByRules: true},
	{Name: "priority", Kind: FeatureKindString, Description: "Claim priority", UsedByRules: true},
	{Name: "priority_high", Kind: FeatureKindBool, Description: "Priority is high/urgent", UsedByGLM: true, UsedByRules: true},
}

// glmFeatureOrder is the canonical column order for the GLM feature vector.
// Coefficients are stored by name (JSONMapFloat) so order is only a fitting
// convention; we keep one slice to ensure refit/score agree.
var glmFeatureOrder = func() []string {
	out := make([]string, 0, len(FraudFeatureCatalogue))
	for _, f := range FraudFeatureCatalogue {
		if f.UsedByGLM {
			out = append(out, f.Name)
		}
	}
	return out
}()

// ExtractRawFeatures returns a feature map suitable for rule evaluation. Includes
// both raw categoricals (strings) and numerics. Numeric values are float64.
func ExtractRawFeatures(claim models.GroupSchemeClaim) map[string]any {
	out := map[string]any{}

	out["claim_amount"] = claim.ClaimAmount
	out["claim_amount_log"] = math.Log1p(math.Max(0, claim.ClaimAmount))

	dEvent := parseClaimDate(claim.DateOfEvent)
	dNotif := parseClaimDate(claim.DateNotified)
	dReg := parseClaimDate(claim.DateRegistered)
	out["days_event_to_notification"] = daysBetween(dEvent, dNotif)
	out["days_notification_to_registration"] = daysBetween(dNotif, dReg)

	out["attachment_count"] = float64(len(claim.Attachments))

	memberType := strings.ToLower(strings.TrimSpace(claim.MemberType))
	out["member_type"] = memberType
	out["member_type_member"] = boolFloat(memberType == "member")
	out["member_type_spouse"] = boolFloat(memberType == "spouse")
	out["member_type_child"] = boolFloat(memberType == "child")

	rel := strings.ToLower(strings.TrimSpace(claim.RelationshipToMember))
	out["relationship_to_member"] = rel
	out["relationship_self"] = boolFloat(rel == "self" || rel == "member")

	cause := strings.ToLower(strings.TrimSpace(claim.CauseType))
	out["cause_type"] = cause
	out["cause_natural"] = boolFloat(strings.Contains(cause, "natural"))

	priority := strings.ToLower(strings.TrimSpace(claim.Priority))
	out["priority"] = priority
	out["priority_high"] = boolFloat(priority == "high" || priority == "urgent" || priority == "critical")

	return out
}

// GLMFeatureVector projects raw features into the ordered numeric vector the
// GLM consumes. Missing keys default to 0.
func GLMFeatureVector(raw map[string]any) []float64 {
	vec := make([]float64, len(glmFeatureOrder))
	for i, name := range glmFeatureOrder {
		if v, ok := raw[name]; ok {
			if f, ok := v.(float64); ok {
				vec[i] = f
			}
		}
	}
	return vec
}

// fraudSigmoid is the logistic link (local copy; services/win_probability.go
// declares its own sigmoid for a different use case).
func fraudSigmoid(z float64) float64 {
	if z >= 0 {
		ez := math.Exp(-z)
		return 1.0 / (1.0 + ez)
	}
	ez := math.Exp(z)
	return ez / (1.0 + ez)
}

// bandScore translates a [0,1] probability into a Low/Medium/High/Critical
// band. Thresholds are fixed for this iteration.
func bandScore(p float64) string {
	switch {
	case p < 0.25:
		return RiskLevelLow
	case p < 0.5:
		return RiskLevelMedium
	case p < 0.75:
		return RiskLevelHigh
	default:
		return RiskLevelCritical
	}
}

// LoadFraudRiskModel returns the singleton model row, creating it if missing.
func LoadFraudRiskModel() (models.FraudRiskModel, error) {
	var m models.FraudRiskModel
	err := DB.First(&m, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		m = models.FraudRiskModel{ID: 1}
		if cerr := DB.Create(&m).Error; cerr != nil {
			return m, cerr
		}
		return m, nil
	}
	return m, err
}

// ScoreClaimResult is the payload returned to the frontend after a fraud check.
type ScoreClaimResult struct {
	GLMScore       float64        `json:"glm_score"`
	GLMBand        string         `json:"glm_band"`
	MatchedRule    *MatchedRule   `json:"matched_rule"`
	FinalRiskLevel string         `json:"final_risk_level"`
	Rationale      string         `json:"rationale"`
	Features       map[string]any `json:"features"`
	AssessmentID   int            `json:"assessment_id"`
}

// MatchedRule is the trimmed rule shape included in a score result.
type MatchedRule struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	RiskLevel string `json:"risk_level"`
}

// ScoreClaim runs the GLM + rule engine for a single claim and persists an
// audit row.
func ScoreClaim(claimID int, user models.AppUser) (ScoreClaimResult, error) {
	var result ScoreClaimResult

	var claim models.GroupSchemeClaim
	if err := DB.Preload("Attachments").First(&claim, claimID).Error; err != nil {
		return result, err
	}

	model, err := LoadFraudRiskModel()
	if err != nil {
		return result, err
	}

	raw := ExtractRawFeatures(claim)
	vec := GLMFeatureVector(raw)

	z := model.Intercept
	for i, name := range glmFeatureOrder {
		z += model.Coefficients[name] * vec[i]
	}
	score := fraudSigmoid(z)
	glmBand := bandScore(score)

	matched, err := evaluateRules(raw)
	if err != nil {
		return result, err
	}

	result.GLMScore = score
	result.GLMBand = glmBand
	result.Features = raw

	if matched != nil {
		result.MatchedRule = &MatchedRule{ID: matched.ID, Name: matched.Name, RiskLevel: matched.RiskLevel}
		result.FinalRiskLevel = matched.RiskLevel
		result.Rationale = fmt.Sprintf("Rule %q matched (risk level: %s)", matched.Name, matched.RiskLevel)
	} else {
		result.FinalRiskLevel = glmBand
		if model.TrainedAt == nil {
			result.Rationale = fmt.Sprintf("GLM untrained; defaulted to %s. Train the model from Administration → Fraud Risk Config.", glmBand)
		} else {
			result.Rationale = fmt.Sprintf("GLM score %.3f → %s", score, glmBand)
		}
	}

	featuresJSON, _ := json.Marshal(raw)
	var matchedID *int
	if matched != nil {
		mid := matched.ID
		matchedID = &mid
	}
	audit := models.FraudRiskAssessment{
		ClaimID:        claim.ID,
		GLMScore:       score,
		GLMBand:        glmBand,
		MatchedRuleID:  matchedID,
		FinalRiskLevel: result.FinalRiskLevel,
		Features:       models.JSON(featuresJSON),
		Rationale:      result.Rationale,
		ComputedBy:     user.UserName,
	}
	if err := DB.Create(&audit).Error; err != nil {
		return result, err
	}
	result.AssessmentID = audit.ID

	return result, nil
}

// evaluateRules iterates enabled rules in priority desc and returns the first match.
func evaluateRules(raw map[string]any) (*models.FraudRiskRule, error) {
	var rules []models.FraudRiskRule
	if err := DB.Where("enabled = ?", true).Order("priority DESC, id ASC").Find(&rules).Error; err != nil {
		return nil, err
	}
	for i := range rules {
		ok, err := matchConditions(rules[i].Conditions, raw)
		if err != nil {
			return nil, fmt.Errorf("rule %d (%s): %w", rules[i].ID, rules[i].Name, err)
		}
		if ok {
			return &rules[i], nil
		}
	}
	return nil, nil
}

// matchConditions recursively evaluates a conditions JSON tree against the raw
// feature map. Supports {"all":[...]}, {"any":[...]} and leaf
// {"field":..., "op":..., "value":...}.
func matchConditions(raw models.JSON, features map[string]any) (bool, error) {
	if len(raw) == 0 {
		return false, nil
	}
	var node map[string]any
	if err := json.Unmarshal([]byte(raw), &node); err != nil {
		return false, fmt.Errorf("invalid conditions JSON: %w", err)
	}
	return evalNode(node, features)
}

func evalNode(node map[string]any, features map[string]any) (bool, error) {
	if children, ok := node["all"]; ok {
		arr, ok := children.([]any)
		if !ok {
			return false, errors.New(`"all" must be an array`)
		}
		for _, c := range arr {
			cm, ok := c.(map[string]any)
			if !ok {
				return false, errors.New(`"all" entries must be objects`)
			}
			ok2, err := evalNode(cm, features)
			if err != nil {
				return false, err
			}
			if !ok2 {
				return false, nil
			}
		}
		return true, nil
	}
	if children, ok := node["any"]; ok {
		arr, ok := children.([]any)
		if !ok {
			return false, errors.New(`"any" must be an array`)
		}
		for _, c := range arr {
			cm, ok := c.(map[string]any)
			if !ok {
				return false, errors.New(`"any" entries must be objects`)
			}
			ok2, err := evalNode(cm, features)
			if err != nil {
				return false, err
			}
			if ok2 {
				return true, nil
			}
		}
		return false, nil
	}

	field, _ := node["field"].(string)
	op, _ := node["op"].(string)
	value := node["value"]
	if field == "" || op == "" {
		return false, errors.New("leaf condition requires field and op")
	}
	fv, ok := features[field]
	if !ok {
		return false, nil
	}
	return compare(fv, op, value)
}

func compare(left any, op string, right any) (bool, error) {
	switch op {
	case "eq":
		return equalsAny(left, right), nil
	case "ne":
		return !equalsAny(left, right), nil
	case "gt", "gte", "lt", "lte":
		l, okL := toFloat(left)
		r, okR := toFloat(right)
		if !okL || !okR {
			return false, nil
		}
		switch op {
		case "gt":
			return l > r, nil
		case "gte":
			return l >= r, nil
		case "lt":
			return l < r, nil
		case "lte":
			return l <= r, nil
		}
	case "in":
		arr, ok := right.([]any)
		if !ok {
			return false, errors.New(`"in" value must be an array`)
		}
		for _, v := range arr {
			if equalsAny(left, v) {
				return true, nil
			}
		}
		return false, nil
	}
	return false, fmt.Errorf("unsupported op: %s", op)
}

func equalsAny(a, b any) bool {
	if la, okA := toFloat(a); okA {
		if lb, okB := toFloat(b); okB {
			return la == lb
		}
	}
	return strings.EqualFold(fmt.Sprint(a), fmt.Sprint(b))
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int32:
		return float64(t), true
	case int64:
		return float64(t), true
	case json.Number:
		f, err := t.Float64()
		return f, err == nil
	}
	return 0, false
}

// ---------- GLM fitting (logistic regression via L-BFGS) ----------

// RefitResult summarises a refit run for the admin UI.
type RefitResult struct {
	SampleSize    int                `json:"sample_size"`
	PositiveCount int                `json:"positive_count"`
	AUC           float64            `json:"auc"`
	Intercept     float64            `json:"intercept"`
	Coefficients  map[string]float64 `json:"coefficients"`
}

// Minimum sample requirements — below these the model would be useless and we
// reject the refit so the admin knows to gather more labelled data.
const (
	minRefitSamples  = 100
	minRefitPositive = 30
)

// RefitFraudRiskModel pulls labelled historical claims, fits a logistic
// regression with small L2 regularisation, and upserts the singleton model
// row. The label is 1 if the assessment marked the claim as high/critical
// fraud risk OR flagged for investigation.
func RefitFraudRiskModel(user models.AppUser) (RefitResult, error) {
	var res RefitResult

	type sample struct {
		x []float64
		y float64
	}

	// Load all closed assessments and join to their claims.
	var assessments []models.GroupSchemeClaimAssessment
	if err := DB.Where("assessment_outcome <> ''").Find(&assessments).Error; err != nil {
		return res, err
	}
	if len(assessments) == 0 {
		return res, errors.New("no closed assessments available for training")
	}
	claimIDs := make([]int, 0, len(assessments))
	for _, a := range assessments {
		claimIDs = append(claimIDs, a.ClaimID)
	}
	var claims []models.GroupSchemeClaim
	if err := DB.Preload("Attachments").Where("id IN ?", claimIDs).Find(&claims).Error; err != nil {
		return res, err
	}
	claimByID := make(map[int]models.GroupSchemeClaim, len(claims))
	for _, c := range claims {
		claimByID[c.ID] = c
	}

	// Train on the assessor's judgement (assessor_risk_level), not the system's
	// own past output — using fraud_risk_level here would be circular and
	// would teach the model to reinforce its prior decisions.
	samples := make([]sample, 0, len(assessments))
	positives := 0
	for _, a := range assessments {
		c, ok := claimByID[a.ClaimID]
		if !ok {
			continue
		}
		level := strings.ToLower(strings.TrimSpace(a.AssessorRiskLevel))
		if level == "" {
			continue // unlabelled — skip
		}
		raw := ExtractRawFeatures(c)
		x := GLMFeatureVector(raw)
		y := 0.0
		if level == RiskLevelHigh || level == RiskLevelCritical {
			y = 1.0
			positives++
		}
		samples = append(samples, sample{x: x, y: y})
	}

	if len(samples) < minRefitSamples {
		return res, fmt.Errorf("not enough labelled samples (need %d closed assessments with an assessor risk level, have %d)", minRefitSamples, len(samples))
	}
	if positives < minRefitPositive {
		return res, fmt.Errorf("not enough fraud-positive samples (need %d assessor-flagged high/critical claims, have %d)", minRefitPositive, positives)
	}

	// Optimise negative log-likelihood + L2 reg with L-BFGS.
	// Params layout: [intercept, β1, β2, ...]
	n := len(glmFeatureOrder)
	const lambda = 0.01 // small L2 to keep things stable on modest n

	problem := optimize.Problem{
		Func: func(theta []float64) float64 {
			nll := 0.0
			for _, s := range samples {
				z := theta[0]
				for i := 0; i < n; i++ {
					z += theta[1+i] * s.x[i]
				}
				p := fraudSigmoid(z)
				eps := 1e-12
				p = math.Min(math.Max(p, eps), 1-eps)
				nll -= s.y*math.Log(p) + (1-s.y)*math.Log(1-p)
			}
			// L2 reg (not on intercept).
			for i := 0; i < n; i++ {
				nll += 0.5 * lambda * theta[1+i] * theta[1+i]
			}
			return nll
		},
		Grad: func(grad, theta []float64) {
			for i := range grad {
				grad[i] = 0
			}
			for _, s := range samples {
				z := theta[0]
				for i := 0; i < n; i++ {
					z += theta[1+i] * s.x[i]
				}
				err := fraudSigmoid(z) - s.y
				grad[0] += err
				for i := 0; i < n; i++ {
					grad[1+i] += err * s.x[i]
				}
			}
			for i := 0; i < n; i++ {
				grad[1+i] += lambda * theta[1+i]
			}
		},
	}

	init := make([]float64, n+1)
	result, err := optimize.Minimize(problem, init, nil, &optimize.LBFGS{})
	if err != nil {
		return res, fmt.Errorf("optimiser failed: %w", err)
	}
	theta := result.X

	// Compute training AUC.
	ys := make([]float64, len(samples))
	ps := make([]float64, len(samples))
	for i, s := range samples {
		z := theta[0]
		for j := 0; j < n; j++ {
			z += theta[1+j] * s.x[j]
		}
		ys[i] = s.y
		ps[i] = fraudSigmoid(z)
	}
	auc := rocAUC(ys, ps)

	coefs := make(models.JSONMapFloat, n)
	for i, name := range glmFeatureOrder {
		coefs[name] = theta[1+i]
	}

	// Upsert singleton row.
	now := time.Now()
	updates := models.FraudRiskModel{
		ID:            1,
		Intercept:     theta[0],
		Coefficients:  coefs,
		TrainedAt:     &now,
		TrainedBy:     user.UserName,
		SampleSize:    len(samples),
		PositiveCount: positives,
		AUC:           auc,
	}
	// Use Save so existing row is overwritten.
	if err := DB.Save(&updates).Error; err != nil {
		return res, err
	}

	res = RefitResult{
		SampleSize:    len(samples),
		PositiveCount: positives,
		AUC:           auc,
		Intercept:     theta[0],
		Coefficients:  map[string]float64(coefs),
	}
	return res, nil
}

// rocAUC computes the area under the ROC curve using the Mann-Whitney U
// formula. Robust to ties.
func rocAUC(ys, ps []float64) float64 {
	n := len(ys)
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool { return ps[idx[a]] < ps[idx[b]] })

	// Assign average ranks (1-indexed), handling ties.
	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i
		for j+1 < n && ps[idx[j+1]] == ps[idx[i]] {
			j++
		}
		avg := float64(i+j+2) / 2.0
		for k := i; k <= j; k++ {
			ranks[k] = avg
		}
		i = j + 1
	}

	rankSumPos := 0.0
	pos := 0
	for k, orig := range idx {
		if ys[orig] >= 0.5 {
			rankSumPos += ranks[k]
			pos++
		}
	}
	neg := n - pos
	if pos == 0 || neg == 0 {
		return 0
	}
	return (rankSumPos - float64(pos*(pos+1))/2.0) / float64(pos*neg)
}

// ---------- Rule CRUD ----------

// ValidateFraudRiskRule checks an incoming rule from the admin UI.
func ValidateFraudRiskRule(r *models.FraudRiskRule) error {
	r.Name = strings.TrimSpace(r.Name)
	r.RiskLevel = strings.ToLower(strings.TrimSpace(r.RiskLevel))
	if r.Name == "" {
		return errors.New("rule name is required")
	}
	if !validRiskLevels[r.RiskLevel] {
		return errors.New("risk_level must be low|medium|high|critical")
	}
	if len(r.Conditions) == 0 {
		return errors.New("conditions are required")
	}
	// Sanity-check the conditions tree by attempting an evaluation against an
	// empty feature map — any structural error will surface here.
	if _, err := matchConditions(r.Conditions, map[string]any{}); err != nil {
		return fmt.Errorf("invalid conditions: %w", err)
	}
	if r.Priority == 0 {
		r.Priority = 50
	}
	return nil
}

func ListFraudRiskRules() ([]models.FraudRiskRule, error) {
	var rules []models.FraudRiskRule
	err := DB.Order("priority DESC, id ASC").Find(&rules).Error
	return rules, err
}

func CreateFraudRiskRule(r models.FraudRiskRule, user models.AppUser) (models.FraudRiskRule, error) {
	if err := ValidateFraudRiskRule(&r); err != nil {
		return r, err
	}
	r.ID = 0
	r.UpdatedBy = user.UserName
	if err := DB.Create(&r).Error; err != nil {
		return r, err
	}
	return r, nil
}

func UpdateFraudRiskRule(id int, r models.FraudRiskRule, user models.AppUser) (models.FraudRiskRule, error) {
	var existing models.FraudRiskRule
	if err := DB.First(&existing, id).Error; err != nil {
		return existing, err
	}
	r.ID = existing.ID
	r.CreatedAt = existing.CreatedAt
	if err := ValidateFraudRiskRule(&r); err != nil {
		return existing, err
	}
	r.UpdatedBy = user.UserName
	if err := DB.Save(&r).Error; err != nil {
		return existing, err
	}
	return r, nil
}

func DeleteFraudRiskRule(id int) error {
	res := DB.Delete(&models.FraudRiskRule{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ---------- helpers ----------

func boolFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

var claimDateLayouts = []string{
	"2006-01-02",
	time.RFC3339,
	"2006-01-02T15:04:05",
	"02/01/2006",
	"2006/01/02",
}

func parseClaimDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	for _, layout := range claimDateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func daysBetween(a, b time.Time) float64 {
	if a.IsZero() || b.IsZero() {
		return 0
	}
	return b.Sub(a).Hours() / 24.0
}

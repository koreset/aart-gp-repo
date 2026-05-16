package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// JSONMapFloat is a JSON-serializable map[string]float64 that satisfies
// driver.Valuer and sql.Scanner so GORM can persist it. Mirrors JSONMapBool
// in group_pricing.go.
type JSONMapFloat map[string]float64

func (m JSONMapFloat) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	b, err := json.Marshal(map[string]float64(m))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (m *JSONMapFloat) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONMapFloat Scan: %T", value)
	}
	if len(data) == 0 {
		*m = nil
		return nil
	}
	var tmp map[string]float64
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*m = JSONMapFloat(tmp)
	return nil
}

// FraudRiskModel stores the active logistic-regression coefficients used to
// score a claim's fraud probability. Singleton row (ID=1) mirrors the
// GroupPricingSetting pattern.
type FraudRiskModel struct {
	ID            int          `json:"id" gorm:"primaryKey"`
	Intercept     float64      `json:"intercept"`
	Coefficients  JSONMapFloat `json:"coefficients" gorm:"type:json"`
	TrainedAt     *time.Time   `json:"trained_at"`
	TrainedBy     string       `json:"trained_by"`
	SampleSize    int          `json:"sample_size"`
	PositiveCount int          `json:"positive_count"`
	AUC           float64      `json:"auc"`
	UpdatedAt     time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// FraudRiskRule is a company-configurable rule that asserts a fraud risk level
// when its Conditions match against a claim's features. Matching rules
// override the GLM result.
type FraudRiskRule struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Conditions  JSON      `json:"conditions" gorm:"type:json"`
	RiskLevel   string    `json:"risk_level"`
	Priority    int       `json:"priority" gorm:"default:50"`
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// FraudRiskAssessment is an append-only audit row for each fraud-check run.
type FraudRiskAssessment struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	ClaimID        int       `json:"claim_id" gorm:"index"`
	GLMScore       float64   `json:"glm_score"`
	GLMBand        string    `json:"glm_band"`
	MatchedRuleID  *int      `json:"matched_rule_id"`
	FinalRiskLevel string    `json:"final_risk_level"`
	Features       JSON      `json:"features" gorm:"type:json"`
	Rationale      string    `json:"rationale"`
	ComputedAt     time.Time `json:"computed_at" gorm:"autoCreateTime"`
	ComputedBy     string    `json:"computed_by"`
}

package services

import (
	"fmt"
	"math"

	"api/models"
)

// GetMethodAssignments returns all per-product method assignments for a portfolio.
func GetMethodAssignments(portfolioID int) ([]models.LicIbnrMethodAssignment, error) {
	var assignments []models.LicIbnrMethodAssignment
	err := DB.Where("portfolio_id = ?", portfolioID).Find(&assignments).Error
	return assignments, err
}

// SaveMethodAssignment upserts a single assignment.
// The unique key is (portfolio_id, product_code, accident_year_from, accident_year_to):
//   - Both year fields = 0 → product-level assignment (one per product)
//   - Non-zero years → AY-range assignment (multiple can exist per product, one per range)
func SaveMethodAssignment(a models.LicIbnrMethodAssignment) (models.LicIbnrMethodAssignment, error) {
	var existing models.LicIbnrMethodAssignment
	result := DB.Where(
		"portfolio_id = ? AND product_code = ? AND accident_year_from = ? AND accident_year_to = ?",
		a.PortfolioID, a.ProductCode, a.AccidentYearFrom, a.AccidentYearTo,
	).First(&existing)
	if result.Error == nil {
		// Update in place — preserve original creator
		a.ID = existing.ID
		if a.CreatedBy == "" {
			a.CreatedBy = existing.CreatedBy
		}
		if a.UserName == "" {
			a.UserName = existing.UserName
		}
	}
	err := DB.Save(&a).Error
	return a, err
}

// DeleteMethodAssignment removes an assignment by primary key.
func DeleteMethodAssignment(id int) error {
	return DB.Delete(&models.LicIbnrMethodAssignment{}, id).Error
}

// AutoAssignMethods evaluates two data-quality signals for every product in
// the portfolio and writes rule-type assignments.  Manual assignments are
// left untouched.
//
// Signals:
//   - ActualDataYears = COUNT(DISTINCT damage_year) per product_code
//   - ActualCV       = STDDEV(claim_amount) / AVG(claim_amount) per product_code
//
// Decision:
//   Both pass  →  chain-ladder
//   Either fails → fallbackMethod  (default: "bornhuetter-ferguson")
func AutoAssignMethods(
	portfolioID int,
	portfolioName string,
	claimsYear string,
	version string,
	minDataYears int,
	maxCV float64,
	fallbackMethod string,
	createdBy string,
	userName string,
) ([]models.LicIbnrMethodAssignment, error) {

	if fallbackMethod == "" {
		fallbackMethod = "bornhuetter-ferguson"
	}

	// Fetch existing manual assignments so we don't overwrite them
	var existing []models.LicIbnrMethodAssignment
	DB.Where("portfolio_id = ? AND assignment_type = 'manual'", portfolioID).Find(&existing)
	manualProducts := make(map[string]bool)
	for _, m := range existing {
		manualProducts[m.ProductCode] = true
	}

	// Compute data-quality signals per product from LICClaimsInput
	type productSignal struct {
		ProductCode string
		DataYears   int
		AvgClaim    float64
		StdClaim    float64
	}

	// We'll do this in application code to stay database-agnostic
	var claims []models.LICClaimsInput
	err := DB.Where(
		"portfolio_name = ? AND year = ? AND version_name = ?",
		portfolioName, claimsYear, version,
	).Find(&claims).Error
	if err != nil {
		return nil, fmt.Errorf("auto-assign: querying claims: %w", err)
	}

	// Group by product code
	type accum struct {
		years  map[int]bool
		values []float64
	}
	grouped := make(map[string]*accum)
	for _, c := range claims {
		if _, ok := grouped[c.ProductCode]; !ok {
			grouped[c.ProductCode] = &accum{years: make(map[int]bool)}
		}
		grouped[c.ProductCode].years[c.DamageYear] = true
		grouped[c.ProductCode].values = append(grouped[c.ProductCode].values, c.ClaimAmount)
	}

	var results []models.LicIbnrMethodAssignment

	for productCode, acc := range grouped {
		if manualProducts[productCode] {
			continue // don't touch manual assignments
		}

		dataYears := len(acc.years)
		cv := computeCV(acc.values)

		// Determine method
		method := "chain-ladder"
		reason := fmt.Sprintf("data years %d ≥ %d and CV %.3f ≤ %.3f → Chain Ladder", dataYears, minDataYears, cv, maxCV)
		if dataYears < minDataYears {
			method = fallbackMethod
			reason = fmt.Sprintf("data years %d < threshold %d → %s", dataYears, minDataYears, fallbackMethod)
		} else if cv > maxCV {
			method = fallbackMethod
			reason = fmt.Sprintf("CV %.3f > threshold %.3f → %s", cv, maxCV, fallbackMethod)
		}

		assignment := models.LicIbnrMethodAssignment{
			PortfolioID:     portfolioID,
			PortfolioName:   portfolioName,
			ProductCode:     productCode,
			Method:          method,
			AssignmentType:  "rule",
			MinDataYears:    minDataYears,
			MaxCVThreshold:  maxCV,
			ActualDataYears: dataYears,
			ActualCV:        cv,
			RuleReason:      reason,
			CreatedBy:       createdBy,
			UserName:        userName,
			IsActive:        true,
		}
		saved, err := SaveMethodAssignment(assignment)
		if err != nil {
			return nil, fmt.Errorf("auto-assign: saving %s: %w", productCode, err)
		}
		results = append(results, saved)
	}

	return results, nil
}

// BuildProductMethodMap loads all assignments for a portfolio and returns a
// map[productCode]method.  The caller falls back to the run-level method when
// a product has no entry — preserving the existing single-method behaviour.
// Kept for backward compatibility; prefer BuildProductAssignmentList for new code.
func BuildProductMethodMap(portfolioID int) map[string]string {
	assignments, err := GetMethodAssignments(portfolioID)
	m := make(map[string]string)
	if err != nil {
		return m
	}
	for _, a := range assignments {
		m[a.ProductCode] = a.Method
	}
	return m
}

// BuildProductAssignmentList returns all assignments for a portfolio as a slice.
// The engine passes this slice to ResolveAYMethod for per-accident-year method resolution.
// Returns an empty (non-nil) slice when there are no assignments — safe to pass
// to ResolveAYMethod, which then falls back to the run-level method for every row.
func BuildProductAssignmentList(portfolioID int) []models.LicIbnrMethodAssignment {
	assignments, err := GetMethodAssignments(portfolioID)
	if err != nil {
		return []models.LicIbnrMethodAssignment{}
	}
	return assignments
}

// ResolveAYMethod returns the most specific IBNR method for a given product code
// and accident year, using the following precedence:
//
//  1. AY-range assignment where AccidentYearFrom ≤ accidentYear ≤ AccidentYearTo.
//     When multiple ranges overlap the given year, the narrowest range wins.
//  2. Product-level assignment (AccidentYearFrom = 0 and AccidentYearTo = 0).
//  3. fallbackMethod (the run-level IBNRMethod) — preserves single-method behaviour.
func ResolveAYMethod(
	assignments []models.LicIbnrMethodAssignment,
	productCode string,
	accidentYear int,
	fallbackMethod string,
) string {
	var productLevel *models.LicIbnrMethodAssignment
	var bestRange *models.LicIbnrMethodAssignment
	bestWidth := math.MaxInt32

	for i := range assignments {
		a := &assignments[i]
		if !a.IsActive {
			continue
		}
		if a.ProductCode != productCode {
			continue
		}
		// Product-level assignment: both year fields are zero.
		if a.AccidentYearFrom == 0 && a.AccidentYearTo == 0 {
			if productLevel == nil {
				productLevel = a
			}
			continue
		}
		// AY-range assignment: check if this accident year falls inside the range.
		if accidentYear >= a.AccidentYearFrom && accidentYear <= a.AccidentYearTo {
			width := a.AccidentYearTo - a.AccidentYearFrom
			if width < bestWidth {
				bestWidth = width
				bestRange = a
			}
		}
	}

	// Precedence: narrowest AY-range > product-level > run-level fallback
	if bestRange != nil {
		return bestRange.Method
	}
	if productLevel != nil {
		return productLevel.Method
	}
	return fallbackMethod
}

// computeCV returns coefficient of variation (σ/μ) for a slice of values.
// Returns 0 when the slice is empty or the mean is zero.
func computeCV(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(n)
	if mean == 0 {
		return 0
	}
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(n)
	return math.Sqrt(variance) / mean
}

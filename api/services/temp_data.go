package services

// AggregationData holds the sums needed for calculation for each duration
type AggregationData struct {
	// Common
	Exposure float64

	// Actuals
	ActualLapses float64

	// Expected
	ExpectedLapses float64
}

// OutputRow defines the structure for each row in the AG-Grid
type OutputRow struct {
	Duration string `json:"Duration"`

	ActualLapses   float64 `json:"ActualLapses"`
	ExpectedLapses float64 `json:"ExpectedLapses"`
	Exposure       float64 `json:"Exposure"` // Common exposure
	ActualUX       float64 `json:"ActualUX"`
	ExpectedUX     float64 `json:"ExpectedUX"`

	ActualAnnualRate    float64 `json:"ActualAnnualRate"`
	ExpectedAnnualRate  float64 `json:"ExpectedAnnualRate"`
	ActualMonthlyRate   float64 `json:"ActualMonthlyRate"`
	ExpectedMonthlyRate float64 `json:"ExpectedMonthlyRate"`
}

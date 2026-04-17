package services

import (
	"api/models"
	"math"
	"testing"
)

func ptrFloat(v float64) *float64 { return &v }

func approxEqual(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

func TestComputeProgressiveCommissionFromBands(t *testing.T) {
	// FSCA-style sliding scale for brokers.
	scale := []models.CommissionStructure{
		{Channel: "broker", LowerBound: 0, UpperBound: ptrFloat(200000), MaximumCommission: 0.075, ApplicableRate: 0.075},
		{Channel: "broker", LowerBound: 200000, UpperBound: ptrFloat(300000), MaximumCommission: 0.05, ApplicableRate: 0.05},
		{Channel: "broker", LowerBound: 300000, UpperBound: ptrFloat(600000), MaximumCommission: 0.03, ApplicableRate: 0.03},
		{Channel: "broker", LowerBound: 600000, UpperBound: ptrFloat(2000000), MaximumCommission: 0.02, ApplicableRate: 0.02},
		{Channel: "broker", LowerBound: 2000000, UpperBound: nil, MaximumCommission: 0.01, ApplicableRate: 0.01},
	}

	tests := []struct {
		name    string
		bands   []models.CommissionStructure
		premium float64
		want    float64
	}{
		{
			name:    "no bands returns 0",
			bands:   nil,
			premium: 500000,
			want:    0,
		},
		{
			name:    "zero premium returns 0",
			bands:   scale,
			premium: 0,
			want:    0,
		},
		{
			name:    "negative premium returns 0",
			bands:   scale,
			premium: -1,
			want:    0,
		},
		{
			name:    "premium inside first band pays first rate",
			bands:   scale,
			premium: 100000,
			want:    0.075,
		},
		{
			name:    "premium exactly at first band boundary",
			bands:   scale,
			premium: 200000,
			// 200000 * 0.075 = 15000, divided by 200000 = 0.075
			want: 0.075,
		},
		{
			name:    "premium spans two bands blends rates",
			bands:   scale,
			premium: 250000,
			// 200000 * 0.075 + 50000 * 0.05 = 15000 + 2500 = 17500; / 250000 = 0.07
			want: 0.07,
		},
		{
			name:    "premium 700000 example from the plan",
			bands:   scale,
			premium: 700000,
			// 200000*0.075 + 100000*0.05 + 300000*0.03 + 100000*0.02
			// = 15000 + 5000 + 9000 + 2000 = 31000
			// / 700000 = 0.044285714...
			want: 31000.0 / 700000.0,
		},
		{
			name:    "premium lands in unbounded last band",
			bands:   scale,
			premium: 5000000,
			// 200000*0.075 + 100000*0.05 + 300000*0.03 + 1400000*0.02 + 3000000*0.01
			// = 15000 + 5000 + 9000 + 28000 + 30000 = 87000
			// / 5000000 = 0.0174
			want: 87000.0 / 5000000.0,
		},
		{
			name: "single band covers everything",
			bands: []models.CommissionStructure{
				{Channel: "broker", LowerBound: 0, UpperBound: nil, MaximumCommission: 0.1, ApplicableRate: 0.1},
			},
			premium: 123456,
			want:    0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeProgressiveCommissionFromBands(tt.bands, tt.premium)
			if !approxEqual(got, tt.want, 1e-9) {
				t.Errorf("got %.10f, want %.10f", got, tt.want)
			}
		})
	}
}

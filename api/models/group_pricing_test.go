package models

import (
	"math"
	"testing"
)

func TestSchemeTotalLoading(t *testing.T) {
	tests := []struct {
		name string
		in   MemberRatingResultSummary
		want float64
	}{
		{
			name: "zero loading",
			in:   MemberRatingResultSummary{},
			want: 0,
		},
		{
			// SchemeTotalLoading sums expense + profit + admin + other +
			// binder + outsource. CommissionLoading is excluded because
			// commission is added on top of the pre-comm office premium via
			// the progressive allocation, not baked into the gross-up
			// denominator.
			name: "typical loading: expense + profit",
			in: MemberRatingResultSummary{
				ExpenseLoading:    0.05,
				CommissionLoading: 0.10, // ignored
				ProfitLoading:     0.05,
			},
			want: 0.10,
		},
		{
			name: "all six fields sum",
			in: MemberRatingResultSummary{
				ExpenseLoading:    0.05,
				CommissionLoading: 0.10, // ignored
				ProfitLoading:     0.05,
				AdminLoading:      0.01,
				OtherLoading:      0.01,
				BinderFeeRate:     0.02,
				OutsourceFeeRate:  0.01,
			},
			want: 0.15,
		},
		{
			name: "non-binder quote (binder/outsource zero)",
			in: MemberRatingResultSummary{
				ExpenseLoading: 0.10,
				ProfitLoading:  0.15,
				AdminLoading:   0.02,
				OtherLoading:   0.03,
			},
			want: 0.30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.SchemeTotalLoading()
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("SchemeTotalLoading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeOfficePremium(t *testing.T) {
	tests := []struct {
		name        string
		riskPremium float64
		summary     *MemberRatingResultSummary
		want        float64
	}{
		{
			name:        "nil summary returns zero",
			riskPremium: 100,
			summary:     nil,
			want:        0,
		},
		{
			name:        "zero loading gives risk premium unchanged",
			riskPremium: 100,
			summary:     &MemberRatingResultSummary{},
			want:        100,
		},
		{
			// CommissionLoading is intentionally excluded from the gross-up;
			// see SchemeTotalLoading() doc.
			name:        "typical 0.1 loading (expense + profit)",
			riskPremium: 90,
			summary: &MemberRatingResultSummary{
				ExpenseLoading:    0.05,
				CommissionLoading: 0.10, // ignored
				ProfitLoading:     0.05,
			},
			want: 100, // 90 / (1 - 0.10)
		},
		{
			name:        "binder + outsource included in denominator",
			riskPremium: 85,
			summary: &MemberRatingResultSummary{
				ExpenseLoading:   0.05,
				ProfitLoading:    0.05,
				AdminLoading:     0.01,
				OtherLoading:     0.01,
				BinderFeeRate:    0.02,
				OutsourceFeeRate: 0.01,
			},
			want: 100, // 85 / (1 - 0.15)
		},
		{
			name:        "denom guard at 1.0 loading",
			riskPremium: 100,
			summary: &MemberRatingResultSummary{
				ExpenseLoading:    0.50,
				CommissionLoading: 0.30, // ignored
				ProfitLoading:     0.50,
			},
			want: 0,
		},
		{
			name:        "denom guard above 1.0 loading",
			riskPremium: 100,
			summary: &MemberRatingResultSummary{
				ExpenseLoading:    0.60,
				CommissionLoading: 0.30, // ignored
				ProfitLoading:     0.50,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeOfficePremium(tt.riskPremium, tt.summary)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("ComputeOfficePremium(%v, %+v) = %v, want %v",
					tt.riskPremium, tt.summary, got, tt.want)
			}
		})
	}
}

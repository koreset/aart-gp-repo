package models

type System struct {
	TimingLapseZero        float64 `json:"timing_lapse"`
	TimingLapseHalf        float64 `json:"timing_lapse1"`
	TimingLapseOne         float64 `json:"timing_lapse2"`
	TimingDisabilityZero   float64 `json:"timing_disability"`
	TimingDisabilityHalf   float64 `json:"timing_disability1"`
	TimingDisabilityOne    float64 `json:"timing_disability2"`
	TimingMortalityZero    float64 `json:"timing_mortality"`
	TimingMortalityHalf    float64 `json:"timing_mortality1"`
	TimingMortalityOne     float64 `json:"timing_mortality2"`
	TimingRetrenchmentZero float64 `json:"timing_retrenchment"`
	TimingRetrenchmentHalf float64 `json:"timing_retrenchment1"`
	TimingRetrenchmentOne  float64 `json:"timing_retrenchment2"`
}

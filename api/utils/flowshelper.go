package utils

import (
	"api/models"
)

// AggregatedProjectionsContain tells whether a contains x.
func AggregatedProjectionsContain(a *[]models.AggregatedProjection, p models.AggregatedProjection) (int, bool) {
	for i, n := range *a {
		if p.ProjectionMonth == n.ProjectionMonth && p.ProductCode == n.ProductCode && p.JobProductID == n.JobProductID && p.SpCode == n.SpCode {
			return i, true
		}
	}
	return -1, false
}

// ScopedAggregatedProjectionsContain tells whether a contains x.
func ScopedAggregatedProjectionsContain(a []models.ScopedAggregatedProjection, p models.ScopedAggregatedProjection) (int, bool) {
	for i, n := range a {
		if p.RunDate == n.RunDate && p.ProjectionMonth == n.ProjectionMonth && p.IFRS17Group == n.IFRS17Group {
			return i, true
		}
	}
	return -1, false
}

// LICAggregatedProjectionsContain tells whether a contains x.
func LICAggregatedProjectionsContain(a []models.LICAggregatedProjections, p models.LICAggregatedProjections) (int, bool) {
	for i, n := range a {
		if p.RunDate == n.RunDate && p.ProjectionMonth == n.ProjectionMonth && p.IFRS17Group == n.IFRS17Group {
			return i, true
		}
	}
	return -1, false
}

func ModifiedGMMScopedAggregatedProjectionsContain(a []models.ModifiedGMMScopedAggregation, p models.ModifiedGMMScopedAggregation) (int, bool) {
	for i, n := range a {
		if p.ProjectionMonth == n.ProjectionMonth && p.IFRS17Group == n.IFRS17Group {
			return i, true
		}
	}
	return -1, false
}

func AggregatedModifiedGMMProjectionsContain(a []models.AggregatedModifiedGMMProjection, p models.AggregatedModifiedGMMProjection) (int, bool) {
	for i, n := range a {
		if p.ProjectionMonth == n.ProjectionMonth && p.ProductCode == n.ProductCode {
			return i, true
		}
	}
	return -1, false
}

func TriangulationsContain(a []models.LicTriangulation, p models.LicModelPoint) (int, bool) {
	for i, n := range a {
		if p.AccidentYear == n.AccidentYear && p.AccidentMonth == n.AccidentMonth {
			return i, true
		}
	}
	return -1, false
}

func Unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func StatesContains(a *[]models.ProductTransitionState, state string) bool {
	for _, n := range *a {
		if n.State == state {
			return true
		}
	}
	return false
}

func StringArrayContains(strArray *[]string, item string) bool {
	for _, v := range *strArray {
		if v == item {
			return true
		}
	}
	return false
}

func FactorsContains(a *[]models.Fds, state string) bool {
	for _, n := range *a {
		if n.Factor == state {
			return true
		}
	}
	return false
}

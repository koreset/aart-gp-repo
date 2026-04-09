package utils

import (
	"github.com/montanaflynn/stats"
)

func GetMax(input []float64) (float64, error) {
	return stats.Max(input)
}

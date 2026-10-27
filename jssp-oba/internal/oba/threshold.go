package oba

import "math"

// clampThreshold evita que le humbral de OBA2-N sea negativo
func clampThreshold(value float64) float64 {
	return math.Max(0, value)
}

package utilities

import "math"

func RoundTo2Decimals(v float64) float64 {
	return math.Round(v*100) / 100
}

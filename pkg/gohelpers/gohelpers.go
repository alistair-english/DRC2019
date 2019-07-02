package gohelpers

import "math"

// MaxFloat64 returns the index and value of the largest Float64 in a slice
func MaxFloat64(slice []float64) (i int, val float64) {
	var maxI int = 0
	var maxV float64 = 0

	for i, val := range slice {
		if val > maxV {
			maxV = val
			maxI = i
		}
	}

	return maxI, maxV
}

// B2i converts a boolean to an integer
func B2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// DegToRad converts degrees to radians
func DegToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// RadToDeg converts radians to degrees
func RadToDeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}

// IntMin finds the minimum of two integers
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IntMax finds the minimum of two integers
func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

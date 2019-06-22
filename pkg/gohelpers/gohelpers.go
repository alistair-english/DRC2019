package gohelpers

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

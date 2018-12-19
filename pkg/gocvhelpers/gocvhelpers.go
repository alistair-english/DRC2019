package gocvhelpers

import "gocv.io/x/gocv"

// HSVMask creates a HSV mask that can be used in gocv.InRange
func HSVMask(in gocv.Scalar, dst *gocv.Mat, channels int, rows int, cols int) {

	source := gocv.NewMatFromScalar(in, gocv.MatTypeCV8UC3)

	// input is a 1x1x3 with the 3 HSV values we need
	inputChannels := gocv.Split(source)

	// the mask we want is a copy of the input channels but at the size defined by the caller
	mask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	maskChannels := gocv.Split(mask)

	// copy HSV values to the mask
	for c := 0; c < channels; c++ {
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				maskChannels[c].SetUCharAt(row, col, inputChannels[c].GetUCharAt(0, 0))
			}
		}
	}

	gocv.Merge(maskChannels, dst)
}

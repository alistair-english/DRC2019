package cvhelpers

import (
	"image"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/gohelpers"
)

// ReadHSV will read in an image from the supplied VideoCapture and output a HSV file to dst
func ReadHSV(cam *gocv.VideoCapture, dst *gocv.Mat) {
	tempMat := gocv.NewMat()
	defer tempMat.Close()

	cam.Read(&tempMat)
	gocv.CvtColor(tempMat, dst, gocv.ColorBGRToHSV)
}

// HSVMask creates a HSV mask that can be used in gocv.InRange
func HSVMask(in gocv.Scalar, dst *gocv.Mat, channels int, rows int, cols int) {

	source := gocv.NewMatFromScalar(in, gocv.MatTypeCV8UC3)
	defer source.Close()

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

// FindLargestContour finds the largest contour in a binary image using some default settings and returns the countour
func FindLargestContour(in gocv.Mat) []image.Point {
	contours := gocv.FindContours(in, gocv.RetrievalTree, gocv.ChainApproxNone)

	if len(contours) == 0 {
		return nil
	}

	var areas []float64
	for _, contour := range contours {
		areas = append(areas, gocv.ContourArea(contour))
	}

	i, _ := gohelpers.MaxFloat64(areas)

	return contours[i]
}

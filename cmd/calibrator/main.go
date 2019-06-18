package main

import (
	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
)

func b2i(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func hsvScalarFromSliders(hTb *gocv.Trackbar, sTb *gocv.Trackbar, vTb *gocv.Trackbar) gocv.Scalar {
	return gocv.NewScalar(
		float64(hTb.GetPos()),
		float64(sTb.GetPos()),
		float64(vTb.GetPos()),
		0.0)
}

func main() {
	cam, _ := gocv.VideoCaptureDevice(0)
	defer cam.Close()

	displayWindow := gocv.NewWindow("HSV Calibrator")
	defer displayWindow.Close()

	sourceWindow := gocv.NewWindow("Source Image")
	defer displayWindow.Close()

	// Make the sliders
	var (
		upperH = displayWindow.CreateTrackbar("Upper H", 255)
		lowerH = displayWindow.CreateTrackbar("Lower H", 255)

		upperS = displayWindow.CreateTrackbar("Upper S", 255)
		lowerS = displayWindow.CreateTrackbar("Lower S", 255)

		upperV = displayWindow.CreateTrackbar("Upper V", 255)
		lowerV = displayWindow.CreateTrackbar("Lower V", 255)
	)

	var (
		sourceImg = gocv.NewMat()
		hsvImg    = gocv.NewMat()
		threshImg = gocv.NewMat()
	)

	// Image closes
	defer hsvImg.Close()
	defer threshImg.Close()

	// Window Slider Bars
	cvhelpers.ReadHSV(cam, &hsvImg)

	// Calculate our HSV masks
	channels, rows, cols := hsvImg.Channels(), hsvImg.Rows(), hsvImg.Cols()

	var (
		lowerMask = cvhelpers.NewHSVMask(
			gocv.NewScalar(
				0,
				0,
				0,
				0.0),
			channels,
			rows,
			cols)

		upperMask = cvhelpers.NewHSVMask(
			gocv.NewScalar(
				0,
				0,
				0,
				0.0),
			channels,
			rows,
			cols)
	)

	var (
		lowerHSV     = gocv.NewScalar(0, 0, 0, 0)
		upperHSV     = gocv.NewScalar(0, 0, 0, 0)
		prevLowerHSV = gocv.NewScalar(0, 0, 0, 0)
		prevUpperHSV = gocv.NewScalar(0, 0, 0, 0)
	)

	for { // foreva

		lowerHSV = hsvScalarFromSliders(
			lowerH,
			lowerS,
			lowerV)

		upperHSV = hsvScalarFromSliders(
			upperH,
			upperS,
			upperV)

		if lowerHSV != prevLowerHSV || upperHSV != prevUpperHSV {
			lowerMask = cvhelpers.NewHSVMask(
				lowerHSV,
				channels,
				rows,
				cols)

			upperMask = cvhelpers.NewHSVMask(
				upperHSV,
				channels,
				rows,
				cols)
		}

		// Read Image
		cam.Read(&sourceImg)

		// convert to HSV
		gocv.CvtColor(sourceImg, &hsvImg, gocv.ColorBGRToHSV)

		// Calculate threshold
		gocv.InRange(hsvImg, lowerMask, upperMask, &threshImg)

		prevLowerHSV = lowerHSV
		prevUpperHSV = upperHSV

		// Display Images
		displayWindow.IMShow(threshImg)
		sourceWindow.IMShow(sourceImg)

		displayWindow.WaitKey(1)
		sourceWindow.WaitKey(1)
	}

}

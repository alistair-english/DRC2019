package main

import (
	"image"
	"image/color"
	"os"

	"github.com/alistair-english/DRC2019/internal/pkg/config"
	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
)

func main() {

	// Load Configurations
	cvConfig := config.GetCVConfig()

	cam, _ := gocv.VideoCaptureDevice(0)
	defer cam.Close()

	displayWindow := gocv.NewWindow("Display")
	defer displayWindow.Close()

	// Image Mats
	var (
		sourceImg = gocv.NewMat()
		hsvImg    = gocv.NewMat()
		mask      = gocv.NewMat()
		maskBlur  = gocv.NewMat()
	)

	// Image closes
	defer sourceImg.Close()
	defer hsvImg.Close()
	defer mask.Close()
	defer maskBlur.Close()

	cvhelpers.ReadHSV(cam, &hsvImg)

	// Calculate our HSV masks
	channels, rows, cols := hsvImg.Channels(), hsvImg.Rows(), hsvImg.Cols()

	lowerMask := gocv.NewMat()
	upperMask := gocv.NewMat()

	defer lowerMask.Close()
	defer upperMask.Close()

	cvhelpers.HSVMask(gocv.NewScalar(cvConfig.LeftLower.H, cvConfig.LeftLower.S, cvConfig.LeftLower.V, 0.0), &lowerMask, channels, rows, cols)
	cvhelpers.HSVMask(gocv.NewScalar(cvConfig.LeftUpper.H, cvConfig.LeftUpper.S, cvConfig.LeftUpper.V, 0.0), &upperMask, channels, rows, cols)

	for { // inifinte loop

		// get the image
		cam.Read(&sourceImg)

		// blur the image
		gocv.GaussianBlur(sourceImg, &hsvImg, image.Point{5, 5}, 0, 0, gocv.BorderReflect101)

		// convert to HSV
		gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorBGRToHSV)

		// Apply the threshold
		gocv.InRange(hsvImg, lowerMask, upperMask, &mask)

		// Blur the mask
		gocv.GaussianBlur(mask, &maskBlur, image.Point{5, 5}, 0, 0, gocv.BorderReflect101)

		// Find largest contour
		contour := cvhelpers.FindLargestContour(mask)

		// Draw contour on source img
		if contour != nil {
			rect := gocv.BoundingRect(contour)
			gocv.Rectangle(&sourceImg, rect, color.RGBA{255, 0, 0, 0}, 3)
		}

		os.Open()

		// Display source img
		displayWindow.IMShow(mask)
		displayWindow.WaitKey(1)
	}

}

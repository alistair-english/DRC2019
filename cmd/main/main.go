package main

import (
	"fmt"
	"image"

	"github.com/alistair-english/DRC2019/internal/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"gocv.io/x/gocv"
)

func main() {

	const (
		LEFT_LINE  = "LEFT_LINE"
		RIGHT_LINE = "RIGHT_LINE"
	)

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
	)

	// Masks
	var (
		leftMask  = gocv.NewMat()
		rightMask = gocv.NewMat()
	)

	// Image closes
	defer sourceImg.Close()
	defer hsvImg.Close()
	defer leftMask.Close()
	defer rightMask.Close()

	cvhelpers.ReadHSV(cam, &hsvImg)

	// Calculate our HSV masks
	channels, rows, cols := hsvImg.Channels(), hsvImg.Rows(), hsvImg.Cols()

	leftLower := cvhelpers.NewHSVMask(
		gocv.NewScalar(
			cvConfig.LeftLower.H,
			cvConfig.LeftLower.S,
			cvConfig.LeftLower.V,
			0.0),
		channels,
		rows,
		cols)

	leftUpper := cvhelpers.NewHSVMask(
		gocv.NewScalar(
			cvConfig.LeftUpper.H,
			cvConfig.LeftUpper.S,
			cvConfig.LeftUpper.V,
			0.0),
		channels,
		rows,
		cols)

	rightLower := cvhelpers.NewHSVMask(
		gocv.NewScalar(
			cvConfig.RightLower.H,
			cvConfig.RightLower.S,
			cvConfig.RightLower.V,
			0.0),
		channels,
		rows,
		cols)

	rightUpper := cvhelpers.NewHSVMask(
		gocv.NewScalar(
			cvConfig.RightUpper.H,
			cvConfig.RightUpper.S,
			cvConfig.RightUpper.V,
			0.0),
		channels,
		rows,
		cols)

	for { // inifinte loop

		// get the image
		cam.Read(&sourceImg)

		// blur the image
		gocv.GaussianBlur(sourceImg, &hsvImg, image.Point{11, 11}, 0, 0, gocv.BorderReflect101)

		// convert to HSV
		gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorBGRToHSV)

		// Left mask
		gocv.InRange(hsvImg, leftLower, leftUpper, &leftMask)

		// Right mask
		gocv.InRange(hsvImg, rightLower, rightUpper, &rightMask)

		fmt.Println(leftMask.GetIntAt(0, 0))

		displayWindow.IMShow(leftMask)
		displayWindow.WaitKey(1)
	}
}

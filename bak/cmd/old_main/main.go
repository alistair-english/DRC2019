package main

import (
	"image"
	"image/color"

	"github.com/alistair-english/DRC2019/internal/pkg/config"
	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
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

	// Image closes
	defer sourceImg.Close()
	defer hsvImg.Close()

	cvhelpers.ReadHSV(cam, &hsvImg)

	// Calculate our HSV masks
	channels, rows, cols := hsvImg.Channels(), hsvImg.Rows(), hsvImg.Cols()

	var processMask cvhelpers.ImageMod = func(src gocv.Mat, dst *gocv.Mat) {
		// Blur the mask
		// gocv.GaussianBlur(mask, &maskBlur, image.Point{5, 5}, 0, 0, gocv.BorderReflect101)
		return
	}

	// Create objects slice
	objects := []cvhelpers.HSVObject{
		cvhelpers.NewHSVObject(
			LEFT_LINE,

			cvhelpers.NewHSVMask(
				gocv.NewScalar(
					cvConfig.LeftLower.H,
					cvConfig.LeftLower.S,
					cvConfig.LeftLower.V,
					0.0),
				channels,
				rows,
				cols),

			cvhelpers.NewHSVMask(
				gocv.NewScalar(
					cvConfig.LeftUpper.H,
					cvConfig.LeftUpper.S,
					cvConfig.LeftUpper.V,
					0.0),
				channels,
				rows,
				cols),
		),

		cvhelpers.NewHSVObject(
			RIGHT_LINE,

			cvhelpers.NewHSVMask(
				gocv.NewScalar(
					cvConfig.RightLower.H,
					cvConfig.RightLower.S,
					cvConfig.RightLower.V,
					0.0),
				channels,
				rows,
				cols),

			cvhelpers.NewHSVMask(
				gocv.NewScalar(
					cvConfig.RightUpper.H,
					cvConfig.RightUpper.S,
					cvConfig.RightUpper.V,
					0.0),
				channels,
				rows,
				cols),
		),
	}

	// Cleanup masks afer use
	defer func() {
		for i := range objects {
			objects[i].Masks.Lower.Close()
			objects[i].Masks.Upper.Close()
		}
	}()

	for { // inifinte loop

		// get the image
		cam.Read(&sourceImg)

		// blur the image
		gocv.GaussianBlur(sourceImg, &hsvImg, image.Point{11, 11}, 0, 0, gocv.BorderReflect101)

		// convert to HSV
		gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorBGRToHSV)

		// Find the HSV objects in the image
		result := cvhelpers.FindHSVObjects(hsvImg, objects, processMask)

		for _, obj := range result {
			var colour color.RGBA
			if obj.Name == LEFT_LINE {
				colour = color.RGBA{255, 0, 0, 0}
			} else if obj.Name == RIGHT_LINE {
				colour = color.RGBA{0, 255, 0, 0}
			}

			if obj.Area > 1000 {
				gocv.Rectangle(&sourceImg, obj.BoundingBox, colour, 3)
			}
		}
		// fmt.Println()

		// Display source img
		displayWindow.IMShow(sourceImg)
		displayWindow.WaitKey(1)
	}

}

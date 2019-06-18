package main

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/internal/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
)

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {

	const (
		LEFT_LINE  = "LEFT_LINE"
		RIGHT_LINE = "RIGHT_LINE"
	)

	// Load Configurations
	cvConfig := config.GetCVConfig()

	cam, _ := gocv.VideoCaptureDevice(0)
	defer cam.Close()

	cam.Set(gocv.VideoCaptureBrightness, 120)

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
					0.0,
				),
				channels,
				rows,
				cols,
			),

			cvhelpers.NewHSVMask(
				gocv.NewScalar(
					cvConfig.LeftUpper.H,
					cvConfig.LeftUpper.S,
					cvConfig.LeftUpper.V,
					0.0,
				),
				channels,
				rows,
				cols,
			),
		),

		cvhelpers.NewHSVObject(
			RIGHT_LINE,

			cvhelpers.NewHSVMask(
				gocv.NewScalar(
					cvConfig.RightLower.H,
					cvConfig.RightLower.S,
					cvConfig.RightLower.V,
					0.0,
				),
				channels,
				rows,
				cols,
			),

			cvhelpers.NewHSVMask(
				gocv.NewScalar(
					cvConfig.RightUpper.H,
					cvConfig.RightUpper.S,
					cvConfig.RightUpper.V,
					0.0,
				),
				channels,
				rows,
				cols,
			),
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

		found := make(map[string]cvhelpers.HSVObjectResult)

		for _, obj := range result {
			if obj.Area > 50 {
				found[obj.Name] = obj
			}
		}

		// fmt.Printf("%v %v\n", found[RIGHT_LINE].BoundingBox.Min.X, found[LEFT_LINE].BoundingBox.Max.X)

		leftLine, lExists := found[LEFT_LINE]
		rightLine, rExists := found[RIGHT_LINE]

		leftX := b2i(lExists)*leftLine.BoundingBox.Max.X + b2i(!lExists)*0
		rightX := b2i(rExists)*rightLine.BoundingBox.Min.X + b2i(!rExists)*cols

		horDiff := rightX - leftX
		horCoord := leftX + horDiff/2

		vertCoord := int(math.Min(
			float64(b2i(lExists)*leftLine.BoundingBox.Min.Y+b2i(!lExists)*rows),
			float64(b2i(rExists)*rightLine.BoundingBox.Min.Y+b2i(!rExists)*rows),
		))

		gocv.Rectangle(&sourceImg, found[RIGHT_LINE].BoundingBox, color.RGBA{255, 255, 0, 0}, 3)
		gocv.Rectangle(&sourceImg, found[LEFT_LINE].BoundingBox, color.RGBA{0, 0, 255, 0}, 3)

		gocv.Line(&sourceImg, image.Point{cols / 2, rows}, image.Point{horCoord, vertCoord}, color.RGBA{0, 255, 0, 0}, 3)
		gocv.Circle(&sourceImg, image.Point{horCoord, vertCoord}, 5, color.RGBA{255, 0, 0, 0}, 3)

		// Display source img
		displayWindow.IMShow(sourceImg)
		displayWindow.WaitKey(1)
	}

}

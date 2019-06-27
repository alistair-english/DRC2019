package cvservice

import (
	"reflect"

	"github.com/alistair-english/DRC2019/pkg/cvhelpers"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"gocv.io/x/gocv"
)

// CalibratorService provides calibration service
type CalibratorService struct {
	actionRequestChannel chan<- arch.ActionRequest
}

// GetActionRequestType from Service interface
func (c *CalibratorService) GetActionRequestType() reflect.Type {
	// Does not fulfill requests, only creates
	return nil
}

// SetActionRequestChannel from Service interface
func (c *CalibratorService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	c.actionRequestChannel = channel
}

// FulfullActionRequest from Service interface
func (c *CalibratorService) FulfullActionRequest(request arch.ActionRequest) {
	// Does not fulfill requests, only creates
}

func hsvScalarFromSliders(hTb *gocv.Trackbar, sTb *gocv.Trackbar, vTb *gocv.Trackbar) gocv.Scalar {
	return gocv.NewScalar(
		float64(hTb.GetPos()),
		float64(sTb.GetPos()),
		float64(vTb.GetPos()),
		0.0)
}

// Start from Service interface - provides main functionality
func (c *CalibratorService) Start() {
	go func() {
		displayWindow := gocv.NewWindow("HSV Calibrator")
		defer displayWindow.Close()

		sourceWindow := gocv.NewWindow("Source Image")
		defer displayWindow.Close()

		sliderWindow := gocv.NewWindow("Sliders")
		defer sliderWindow.Close()

		// Make the sliders
		var (
			lowerH = sliderWindow.CreateTrackbar("Lower H", 255)
			lowerS = sliderWindow.CreateTrackbar("Lower S", 255)
			lowerV = sliderWindow.CreateTrackbar("Lower V", 255)

			upperH = sliderWindow.CreateTrackbar("Upper H", 255)
			upperS = sliderWindow.CreateTrackbar("Upper S", 255)
			upperV = sliderWindow.CreateTrackbar("Upper V", 255)
		)

		var (
			sourceImg = gocv.NewMat()
			hsvImg    = gocv.NewMat()
			threshImg = gocv.NewMat()
		)

		// Image closes
		defer sourceImg.Close()
		defer hsvImg.Close()
		defer threshImg.Close()

		// Img request setup
		imgReadChannel := make(chan bool, 1)

		// Get an image
		getImgBlocking(c.actionRequestChannel, &sourceImg, imgReadChannel)

		// Convert to HSV
		gocv.CvtColor(sourceImg, &hsvImg, gocv.ColorBGRToHSV)

		// Calculate our HSV masks
		channels, rows, cols := hsvImg.Channels(), hsvImg.Rows(), hsvImg.Cols()

		var (
			lowerMask = cvhelpers.NewHSVMask(
				gocv.NewScalar(
					0,
					0,
					0,
					0.0,
				),
				channels,
				rows,
				cols,
			)

			upperMask = cvhelpers.NewHSVMask(
				gocv.NewScalar(
					0,
					0,
					0,
					0.0,
				),
				channels,
				rows,
				cols,
			)
		)

		var (
			lowerHSV     = gocv.NewScalar(0, 0, 0, 0)
			upperHSV     = gocv.NewScalar(0, 0, 0, 0)
			prevLowerHSV = gocv.NewScalar(0, 0, 0, 0)
			prevUpperHSV = gocv.NewScalar(0, 0, 0, 0)
		)

		// Read Image
		getImgBlocking(c.actionRequestChannel, &sourceImg, imgReadChannel)

		for { // foreva

			lowerHSV = hsvScalarFromSliders(
				lowerH,
				lowerS,
				lowerV,
			)

			upperHSV = hsvScalarFromSliders(
				upperH,
				upperS,
				upperV,
			)

			// Wait for space
			key := displayWindow.WaitKey(500)
			if key == 13 {
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

					prevLowerHSV = lowerHSV
					prevUpperHSV = upperHSV
				}
			}

			// convert to HSV
			gocv.CvtColor(sourceImg, &hsvImg, gocv.ColorBGRToHSV)

			// Calculate threshold
			gocv.InRange(hsvImg, lowerMask, upperMask, &threshImg)

			// Display Images
			displayWindow.IMShow(threshImg)
			sourceWindow.IMShow(sourceImg)

			// Wait for space
			key = sourceWindow.WaitKey(500)
			if key == 32 {
				// Read Image
				getImgBlocking(c.actionRequestChannel, &sourceImg, imgReadChannel)
			}
		}
	}()
}

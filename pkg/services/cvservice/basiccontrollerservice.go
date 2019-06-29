package cvservice

import (
	"fmt"
	"image"
	"math"
	"reflect"

	"github.com/alistair-english/DRC2019/pkg/services/serialservice"

	"github.com/alistair-english/DRC2019/pkg/gohelpers"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"gocv.io/x/gocv"

	"github.com/felixge/pidctrl"
)

// BasicControllerService provides recording service
type BasicControllerService struct {
	actionRequestChannel chan<- arch.ActionRequest
}

// GetActionRequestType from Service interface
func (c *BasicControllerService) GetActionRequestType() reflect.Type {
	// Does not fulfill requests, only creates
	return nil
}

// SetActionRequestChannel from Service interface
func (c *BasicControllerService) SetActionRequestChannel(channel chan<- arch.ActionRequest) {
	c.actionRequestChannel = channel
}

// FulfullActionRequest from Service interface
func (c *BasicControllerService) FulfullActionRequest(request arch.ActionRequest) {
	// Does not fulfill requests, only creates
}

const (
	LEFT_LINE  = "LEFT_LINE"
	RIGHT_LINE = "RIGHT_LINE"
)

// Start from Service interface - provides main functionality
func (c *BasicControllerService) Start() {
	go func() {
		// Load Configurations
		cvConfig := config.GetCVConfig()

		// displayWindow := gocv.NewWindow("Display")
		// defer displayWindow.Close()

		// Image Mats
		var (
			sourceImg = gocv.NewMat()
			hsvImg    = gocv.NewMat()
		)

		// Image closes
		defer sourceImg.Close()
		defer hsvImg.Close()

		// Img Read Channel
		imgReadChannel := make(chan bool, 1)

		getImgBlocking(c.actionRequestChannel, &sourceImg, imgReadChannel)
		gocv.CvtColor(sourceImg, &hsvImg, gocv.ColorBGRToHSV)

		// Calculate our HSV masks
		channels, rows, cols := hsvImg.Channels(), hsvImg.Rows(), hsvImg.Cols()

		// diagonalLen := math.Sqrt(math.Pow(float64(rows), 2) + math.Pow(float64(cols)/2, 2))

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
			getImgBlocking(c.actionRequestChannel, &sourceImg, imgReadChannel)

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

			leftX := gohelpers.B2i(lExists)*leftLine.BoundingBox.Max.X + gohelpers.B2i(!lExists)*0
			rightX := gohelpers.B2i(rExists)*rightLine.BoundingBox.Min.X + gohelpers.B2i(!rExists)*cols

			horDiff := rightX - leftX
			horCoord := leftX + horDiff/2

			vertCoord := int(math.Min(
				float64(gohelpers.B2i(lExists)*leftLine.BoundingBox.Min.Y+gohelpers.B2i(!lExists)*rows),
				float64(gohelpers.B2i(rExists)*rightLine.BoundingBox.Min.Y+gohelpers.B2i(!rExists)*rows),
			))

			// gocv.Rectangle(&sourceImg, found[RIGHT_LINE].BoundingBox, color.RGBA{255, 255, 0, 0}, 3)
			// gocv.Rectangle(&sourceImg, found[LEFT_LINE].BoundingBox, color.RGBA{0, 0, 255, 0}, 3)

			// gocv.Line(&sourceImg, image.Point{cols / 2, rows}, image.Point{horCoord, vertCoord}, color.RGBA{0, 255, 0, 0}, 3)
			// gocv.Circle(&sourceImg, image.Point{horCoord, vertCoord}, 5, color.RGBA{255, 0, 0, 0}, 3)

			cartX := horCoord - (cols / 2)
			cartY := rows - vertCoord

			cartAngle := gohelpers.RadToDeg(math.Atan2(float64(cartY), float64(cartX)))
			// cartLen := math.Sqrt(math.Pow(float64(cartY), 2) + math.Pow(float64(cartX), 2))

			driveAngle := CartesianToDriveAngle(cartAngle)
			// driveSpeed := int8((cartLen / diagonalLen) * 100)

			c.actionRequestChannel <- serialservice.SerialSendActionReq{
				serialservice.Control{
					Dir: -driveAngle,
					Spd: 100,
				},
			}

			fmt.Println(driveAngle, -driveAngle)

			// Display source img
			// displayWindow.IMShow(sourceImg)
			// displayWindow.WaitKey(0)
		}
	}()
}

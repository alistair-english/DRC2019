package cvservice

import (
	"fmt"
	"image"
	"reflect"

	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"gocv.io/x/gocv"
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

// Start from Service interface - provides main functionality
func (c *BasicControllerService) Start() {
	go func() {
		// Image Mats
		var sourceImg = gocv.NewMat()
		defer sourceImg.Close()
		var hsvImg = gocv.NewMat()
		defer hsvImg.Close()

		// Img Read Channel
		imgReadChannel := make(chan bool, 1)

		// Create objects
		objects := getObjectsFromConfig()

		// displayWindow := gocv.NewWindow("Display")
		// defer displayWindow.Close()

		controller := newBasicDriveController()

		for { // inifinte loop

			// get the image
			getImgBlocking(c.actionRequestChannel, &sourceImg, imgReadChannel)

			// blur the image
			gocv.GaussianBlur(sourceImg, &hsvImg, image.Point{11, 11}, 0, 0, gocv.BorderReflect101)

			// convert to HSV
			gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorBGRToHSV)

			// Find the HSV objects in the image
			result := cvhelpers.FindHSVObjects(hsvImg, objects)

			found := make(map[string]cvhelpers.HSVObjectGroupResult)

			for _, obj := range result {
				found[obj.Name] = obj
			}

			control := controller.update(found)

			fmt.Println(control)

			c.actionRequestChannel <- serialservice.SerialSendActionReq{control}

			// fmt.Printf("%v %v\n", found[RIGHT_LINE].BoundingBox.Min.X, found[LEFT_LINE].BoundingBox.Max.X)

			// leftLine, lExists := found[LEFT_LINE]
			// rightLine, rExists := found[RIGHT_LINE]

			// leftX := gohelpers.B2i(lExists)*leftLine.BoundingBox.Max.X + gohelpers.B2i(!lExists)*0
			// rightX := gohelpers.B2i(rExists)*rightLine.BoundingBox.Min.X + gohelpers.B2i(!rExists)*cols

			// horDiff := rightX - leftX
			// horCoord := leftX + horDiff/2

			// vertCoord := int(math.Min(
			// 	float64(gohelpers.B2i(lExists)*leftLine.BoundingBox.Min.Y+gohelpers.B2i(!lExists)*rows),
			// 	float64(gohelpers.B2i(rExists)*rightLine.BoundingBox.Min.Y+gohelpers.B2i(!rExists)*rows),
			// ))

			// gocv.Rectangle(&sourceImg, found[RIGHT_LINE].BoundingBox, color.RGBA{255, 255, 0, 0}, 3)
			// gocv.Rectangle(&sourceImg, found[LEFT_LINE].BoundingBox, color.RGBA{0, 0, 255, 0}, 3)

			// gocv.Line(&sourceImg, image.Point{cols / 2, rows}, image.Point{horCoord, vertCoord}, color.RGBA{0, 255, 0, 0}, 3)
			// gocv.Circle(&sourceImg, image.Point{horCoord, vertCoord}, 5, color.RGBA{255, 0, 0, 0}, 3)

			// cartX := horCoord - (cols / 2)
			// cartY := rows - vertCoord

			// cartAngle := gohelpers.RadToDeg(math.Atan2(float64(cartY), float64(cartX)))
			// // cartLen := math.Sqrt(math.Pow(float64(cartY), 2) + math.Pow(float64(cartX), 2))

			// driveAngle := CartesianToDriveAngle(cartAngle)
			// // driveSpeed := int8((cartLen / diagonalLen) * 100)

			// tunedAngle := controlPID.Update(float64(driveAngle))
			// fmt.Printf("Corrected Angle: %v \n", tunedAngle)

			// c.actionRequestChannel <- serialservice.SerialSendActionReq{
			// 	serialservice.Control{
			// 		Dir: -driveAngle,
			// 		Spd: 100,
			// 	},
			// }

			// fmt.Println(driveAngle, -driveAngle)

			// Display source img
			// displayWindow.IMShow(sourceImg)
			// displayWindow.WaitKey(0)
		}
	}()
}

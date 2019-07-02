package cvservice

import (
	"fmt"
	"image"
	"math"

	"github.com/alistair-english/DRC2019/pkg/logging"

	"github.com/alistair-english/DRC2019/pkg/services/serialservice"

	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"github.com/alistair-english/DRC2019/pkg/gohelpers"
	"github.com/felixge/pidctrl"
)

type basicDriveController struct {
	controlPID *pidctrl.PIDController
	width      int
	height     int
}

func newBasicDriveController() *basicDriveController {
	controller := basicDriveController{}

	controlConfig := config.GetControlPIDConfig()
	controller.controlPID = pidctrl.NewPIDController(controlConfig.Pid.P, controlConfig.Pid.I, controlConfig.Pid.D)
	controller.controlPID.SetOutputLimits(-90.0, 90.0)
	controller.controlPID.Set(0.00)

	cvConfig := config.GetCVConfig()
	controller.width = cvConfig.ImgWidth
	controller.height = cvConfig.ImgHeight

	return &controller
}

func (c *basicDriveController) update(objs []cvhelpers.HSVObjectGroupResult) serialservice.Control {
	var (
		leftLineGroup  cvhelpers.HSVObjectGroupResult
		rightLineGroup cvhelpers.HSVObjectGroupResult
	)

	for _, obj := range objs {
		switch obj.Name {
		case LEFT_LINE:
			leftLineGroup = obj
		case RIGHT_LINE:
			rightLineGroup = obj
		default:
			logging.L().Logln(TAG, logging.All, "Unknown obj detected: %v", obj)
		}
	}

	ang, spd := c.getTrackAngleAndDriveSpeed(leftLineGroup, rightLineGroup)
	return serialservice.Control{
		Dir: ang,
		Spd: spd,
	}
}

func (c *basicDriveController) getTrackAngleAndDriveSpeed(leftLineGroup, rightLineGroup cvhelpers.HSVObjectGroupResult) (int8, int8) {
	var leftLine cvhelpers.HSVObject
	var rightLine cvhelpers.HSVObject

	// Extract the points out of the detected objects
	if len(leftLineGroup.Objects) > 0 {
		leftLine.BoundingBox = leftLineGroup.Objects[0].BoundingBox
	} else {
		// no line found -> create a line out to the left
		leftLine.BoundingBox = image.Rect(0, c.height/2, 0, c.height/2)
	}

	if len(rightLineGroup.Objects) > 0 {
		rightLine.BoundingBox = rightLineGroup.Objects[0].BoundingBox
	} else {
		// no line found -> create a line out to the right
		rightLine.BoundingBox = image.Rect(c.width, c.height/2, c.width, c.height/2)
	}

	fmt.Println("Left X:", leftLine.BoundingBox.Max.X, "Y: ", leftLine.BoundingBox.Min.Y)
	fmt.Println("Right X:", rightLine.BoundingBox.Min.X, "Y: ", rightLine.BoundingBox.Min.Y)

	horDiff := rightLine.BoundingBox.Min.X - leftLine.BoundingBox.Max.X
	horX := leftLine.BoundingBox.Max.X + horDiff/2

	cartX := horX - (c.width / 2)
	cartY := gohelpers.IntMax(c.height-leftLine.BoundingBox.Min.Y, c.height-rightLine.BoundingBox.Min.Y)
	fmt.Println(cartY)

	cartAngle := gohelpers.RadToDeg(math.Atan2(float64(cartY), float64(cartX)))

	trackAngle := CartesianToDriveAngle(cartAngle)
	driveSpeed := int8((float64(cartY) / float64(c.height)) * 100)
	fmt.Println(float64(cartY) / float64(c.height))

	return trackAngle, driveSpeed
}

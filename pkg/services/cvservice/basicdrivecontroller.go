package cvservice

import (
	"fmt"
	"image"
	"math"

	"github.com/alistair-english/DRC2019/pkg/services/serialservice"

	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"github.com/alistair-english/DRC2019/pkg/gohelpers"
	"github.com/felixge/pidctrl"
)

type basicDriveController struct {
	currentObjects map[string]cvhelpers.HSVObjectGroupResult
	controlPID     *pidctrl.PIDController
	width          int
	height         int
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

func (c *basicDriveController) update(objs map[string]cvhelpers.HSVObjectGroupResult) serialservice.Control {
	c.currentObjects = objs
	ang, spd := c.getTrackAngleAndDriveSpeed()
	return serialservice.Control{
		Dir: ang,
		Spd: spd,
	}
}

func (c *basicDriveController) getTrackAngleAndDriveSpeed() (int8, int8) {
	var leftLine cvhelpers.HSVObject
	var rightLine cvhelpers.HSVObject

	// Extract the points out of the detected objects
	if len(c.currentObjects[LEFT_LINE].Objects) > 0 {
		leftLine.BoundingBox = c.currentObjects[LEFT_LINE].Objects[0].BoundingBox
	} else {
		// no line found -> create a line out to the left
		leftLine.BoundingBox = image.Rect(0, c.height, 0, c.height)
	}

	if len(c.currentObjects[RIGHT_LINE].Objects) > 0 {
		rightLine.BoundingBox = c.currentObjects[RIGHT_LINE].Objects[0].BoundingBox
	} else {
		// no line found -> create a line out to the right
		rightLine.BoundingBox = image.Rect(c.width, c.height, c.width, c.height)
	}

	fmt.Println("Left X:", leftLine.BoundingBox.Max.X)
	fmt.Println("Right X:", rightLine.BoundingBox.Min.X)
	fmt.Println()

	horDiff := rightLine.BoundingBox.Min.X - leftLine.BoundingBox.Max.X
	horX := leftLine.BoundingBox.Max.X + horDiff/2

	cartX := horX - (c.width / 2)
	cartY := c.height

	cartAngle := gohelpers.RadToDeg(math.Atan2(float64(c.height-10), float64(cartX)))

	trackAngle := CartesianToDriveAngle(cartAngle)
	driveSpeed := int8((cartY / c.height) * 100)

	return trackAngle, driveSpeed
}

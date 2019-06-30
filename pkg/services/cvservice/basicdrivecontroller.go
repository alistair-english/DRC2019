package cvservice

import (
	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"github.com/felixge/pidctrl"
)

type basicDriveController struct {
	currentObjects map[string]cvhelpers.HSVObjectGroupResult
	controlPID     *pidctrl.PIDController
}

func newBasicDriveController() *basicDriveController {
	controller := basicDriveController{}

	controlConfig := config.GetControlPIDConfig()
	controller.controlPID = pidctrl.NewPIDController(controlConfig.Pid.P, controlConfig.Pid.I, controlConfig.Pid.D)
	controller.controlPID.SetOutputLimits(-90.0, 90.0)
	controller.controlPID.Set(0.00)

	return &controller
}

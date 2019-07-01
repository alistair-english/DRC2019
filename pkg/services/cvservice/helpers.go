package cvservice

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/config"
	"github.com/alistair-english/DRC2019/pkg/cvhelpers"
	"github.com/alistair-english/DRC2019/pkg/logging"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"gocv.io/x/gocv"
)

func getImgBlocking(actionReqChannel chan<- arch.ActionRequest, img *gocv.Mat) {
	timingChannel := make(chan bool, 1)
	actionReqChannel <- cameraservice.GetImageActionReq{
		Img:             img,
		ResponseChannel: timingChannel,
	}
	<-timingChannel
}

func CartesianToDriveAngle(cartAngle float64) int8 {
	if cartAngle <= 90 {
		return int8(-(cartAngle - 90))
	} else if cartAngle <= 180 {
		return int8(90 - cartAngle)
	} else {
		return 0
	}
}

func getObjectsFromConfig() []cvhelpers.HSVObjectGroup {
	cvConfig := config.GetCVConfig()

	objs := make([]cvhelpers.HSVObjectGroup, len(cvConfig.Objects))

	for i, obj := range cvConfig.Objects {
		objs[i] = cvhelpers.NewHSVObjectGroup(
			obj.Name,
			newMaskFromConfig(obj.LowerMask, obj.Name, cvConfig.ImgHeight, cvConfig.ImgWidth, cvConfig.ImgChannels),
			newMaskFromConfig(obj.UpperMask, obj.Name, cvConfig.ImgHeight, cvConfig.ImgWidth, cvConfig.ImgChannels),
			obj.NumToFind,
			obj.MinArea,
		)
	}
	return objs
}

func newMaskFromConfig(cfg config.Mask, objName string, rows, cols, channels int) cvhelpers.HSVMask {
	if cfg.Path != "" {
		mask, err := cvhelpers.NewHSVMaskFromFile(cfg.Path, channels, rows, cols)
		if err == nil {
			return mask
		}
		logging.L().Logln(TAG, logging.All,
			"newMaskFromConfig - WARNING: error creating mask from file for %v (%v) from file. Defaulting to calcs.\n"+err.Error(), cfg, objName)
	}

	return cvhelpers.NewHSVMask(
		cfg.H, cfg.S, cfg.V,
		rows, cols,
	)
}

package cvservice

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"gocv.io/x/gocv"
)

func getImgBlocking(actionReqChannel chan<- arch.ActionRequest, img *gocv.Mat, timingChannel chan bool) {
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

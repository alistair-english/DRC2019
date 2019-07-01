package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
)

func main() {
	router := arch.NewRouter()

	spdTstService := &cameraservice.CameraSpeedTestService{}
	camService, _ := cameraservice.NewPiCamera()

	router.Register(spdTstService)
	router.Register(camService)

	spdTstService.Start()
	camService.Start()

	// router is blocking
	router.Start()
}

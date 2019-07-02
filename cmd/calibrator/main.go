package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
)

func main() {
	router := arch.NewRouter()

	calService := &cvservice.CalibratorService{}
	camService, _ := cameraservice.NewFileReaderCamera("../recorder/recording_07-02-2019_12:50:29.avi")

	router.Register(calService)
	router.Register(camService)

	calService.Start()
	camService.Start()

	// router is blocking
	router.Start()
}

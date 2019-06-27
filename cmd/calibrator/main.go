package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
)

func main() {
	router := arch.NewRouter()

	calService := &cvservice.CalibratorService{}
	camService, _ := cameraservice.NewFileReaderCamera("/home/alistair/Dev/go/src/github.com/alistair-english/DRC2019/cmd/recorder/recording_06-26-2019_48:34:48.avi")

	router.Register(calService)
	router.Register(camService)

	calService.Start()
	camService.Start()

	// router is blocking
	router.Start()
}

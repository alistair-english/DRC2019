package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
)

func main() {
	router := arch.NewRouter()

	calService := &cvservice.CalibratorService{}
	camService, _ := cameraservice.NewFileReaderCamera("/home/alistair/Dev/go/src/github.com/alistair-english/DRC2019/cmd/recordtodisk/recording_06-25-2019_16:10:16.avi")

	router.Register(calService)
	router.Register(camService)

	calService.Start()
	camService.Start()

	// router is blocking
	router.Start()
}

package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"
)

func main() {
	router := arch.NewRouter()

	serService, _ := serialservice.NewFakeSerial()
	camService, _ := cameraservice.NewFileReaderCamera("/home/alistair/Dev/go/src/github.com/alistair-english/DRC2019/cmd/recordtodisk/recording_06-25-2019_16:10:16.avi")

	router.Register(serService)
	router.Register(camService)

	serService.Start()
	camService.Start()

	// router is blocking
	router.Start()
}

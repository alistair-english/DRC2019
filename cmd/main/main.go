package main

import (
	"fmt"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"
)

func main() {
	router := arch.NewRouter()

	serService, err := serialservice.NewFakeSerial()
	if err != nil {
		fmt.Println(err)
		return
	}

	camService, err := cameraservice.NewFileReaderCamera("/home/alistair/Dev/go/src/github.com/alistair-english/DRC2019/cmd/recorder/recording_06-26-2019_48:34:48.avi")
	if err != nil {
		fmt.Println(err)
		return
	}

	ctrlService := &cvservice.BasicControllerService{}

	router.Register(serService)
	router.Register(camService)
	router.Register(ctrlService)

	serService.Start()
	camService.Start()
	ctrlService.Start()

	// router is blocking
	router.Start()
}

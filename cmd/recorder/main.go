package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
)

func main() {
	router := arch.NewRouter()

	recService := cvservice.NewRecorderService("")
	camService, _ := cameraservice.NewPiCamera()

	router.Register(recService)
	router.Register(camService)

	recService.Start()
	camService.Start()

	// router is blocking
	router.Start()
}

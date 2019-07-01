package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
)

func main() {
	router := arch.NewRouter()

	viewService := &cvservice.ViewerService{}
	camService, _ := cameraservice.NewPiCamera()

	router.Register(viewService)
	router.Register(camService)

	viewService.Start()
	camService.Start()

	// router is blocking
	router.Start()
}

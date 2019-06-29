package main

import (
	"fmt"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
	"github.com/alistair-english/DRC2019/pkg/services/seriallogservice"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"

	"github.com/alistair-english/DRC2019/pkg/logging"
)

func main() {
	router := arch.NewRouter()

	fmt.Print("Getting serial... ")
	serService, err := serialservice.NewPiSerial()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Done.")

	fmt.Print("Getting camera... ")
	camService, err := cameraservice.NewPiCamera()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Done.")

	fmt.Print("Getting controller... ")
	ctrlService := &cvservice.BasicControllerService{}
	fmt.Println("Done.")

	fmt.Print("Getting serial log...")
	serLogService := seriallogservice.NewSerialLogService()
	fmt.Println("Done.")

	fmt.Print("Registering services... ")
	router.Register(serService)
	router.Register(camService)
	router.Register(ctrlService)
	router.Register(serLogService)
	fmt.Println("Done.")

	fmt.Print("Starting Serial... ")
	serService.Start()
	fmt.Println("Done.")

	fmt.Print("Staring Camera... ")
	camService.Start()
	fmt.Println("Done.")

	fmt.Print("Starting Controller... ")
	ctrlService.Start()
	fmt.Println("Done.")

	fmt.Print("Starting SerLog... ")
	serLogService.Start()
	fmt.Println("Done.")

	fmt.Print("Starting Logger... ")
	log := logging.Logger()
	log.Init(serLogService)
	fmt.Println("Done.")

	fmt.Println("Starting Router (blocking)")
	// router is blocking
	router.Start()
}

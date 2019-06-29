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

	fmt.Print("Registering services... ")
	router.Register(serService)
	router.Register(camService)
	router.Register(ctrlService)
	fmt.Println("Done.")

	fmt.Print("Starting Serial... ")
	serService.Start()
	fmt.Println("Done.")

	fmt.Print("Staring Camera... ")
	camService.Start()
	fmt.Println("Done.")

	fmt.Println("Starting Controller... ")
	ctrlService.Start()
	fmt.Println("Done.")

	fmt.Println("Starting Router (blocking)")
	// router is blocking
	router.Start()
}

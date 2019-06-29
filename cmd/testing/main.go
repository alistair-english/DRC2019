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

const TAG = "MAIN"

func main() {
	router := arch.NewRouter()

	fmt.Print("Getting Logger... ")
	log := logging.Logger()
	log.Init()
	fmt.Println("Done.")

	log.Print(TAG, logging.All, "Getting serial... ")
	serService, err := serialservice.NewFakeSerial()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Getting camera... ")
	camService, err := cameraservice.NewFileReaderCamera("")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Getting controller... ")
	ctrlService := &cvservice.BasicControllerService{}
	log.Println(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Getting serial log...")
	serLogService := seriallogservice.NewSerialLogService()
	log.Print(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Registering services... ")
	router.Register(serService)
	router.Register(camService)
	router.Register(ctrlService)
	router.Register(serLogService)
	log.Println(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Starting Serial... ")
	serService.Start()
	log.Println(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Staring Camera... ")
	camService.Start()
	log.Println(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Starting Controller... ")
	ctrlService.Start()
	log.Println(TAG, logging.All, "Done.")

	log.Print(TAG, logging.All, "Starting SerLog... ")
	serLogService.Start()
	log.Println(TAG, logging.All, "Done.")

	log.AddSerialLogService(&serLogService)

	log.Println(TAG, logging.All, "Starting Router (blocking)")
	// router is blocking
	router.Start()
}

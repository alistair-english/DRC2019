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

	fmt.Println("Getting Logger... ")
	log := logging.Logger()
	log.Init(true)
	fmt.Println("Done.")

	log.Log(TAG, logging.All, "Getting serial... ")
	serService, err := serialservice.NewFakeSerial()
	if err != nil {
		log.Logln("", logging.All, err.Error())
		return
	}
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Getting camera... ")
	camService, err := cameraservice.NewFileReaderCamera("../recorder/m5-sport.avi")
	if err != nil {
		log.Logln("", logging.All, err.Error())
		return
	}
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Getting controller... ")
	ctrlService := &cvservice.BasicControllerService{}
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Getting serial log...")
	serLogService := seriallogservice.NewSerialLogService()
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Registering services... ")
	router.Register(serService)
	router.Register(camService)
	router.Register(ctrlService)
	router.Register(serLogService)
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Starting Serial... ")
	serService.Start()
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Staring Camera... ")
	camService.Start()
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Starting Controller... ")
	ctrlService.Start()
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Starting SerLog... ")
	serLogService.Start()
	log.Logln("", logging.All, "Done.")

	// log.AddSerialLogService(serLogService)

	log.Logln(TAG, logging.All, "Starting Router (blocking)")
	// router is blocking
	router.Start()
}

package main

import (
	"fmt"

	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/services/cameraservice"
	"github.com/alistair-english/DRC2019/pkg/services/cvservice"
	"github.com/alistair-english/DRC2019/pkg/services/seriallogservice"
	"github.com/alistair-english/DRC2019/pkg/services/serialservice"

	"github.com/alistair-english/DRC2019/pkg/logging"

	"flag"
	"github.com/fatih/color"
)

const TAG = "MAIN"

func main() {
	// Get flags from CLI inputs (Use -log to enable logging)
	var logEnabled bool
	flag.BoolVar(&logEnabled, "log", false, "Add this flag to enable logging")
	flag.Parse()

	if logEnabled {
		color.Green("Logging Enabled")
	} else {
		color.Red("Logging Disabled")
	}

	router := arch.NewRouter()

	fmt.Println("Getting Logger... ")
	log := logging.Logger()
	log.Init()
	fmt.Println("Done.")

	log.Log(TAG, logging.All, "Getting serial... ")
	serService, err := serialservice.NewPiSerial()
	if err != nil {
		log.Logln("", logging.All, err.Error())
		return
	}
	log.Logln("", logging.All, "Done.")

	log.Log(TAG, logging.All, "Getting camera... ")
	camService, err := cameraservice.NewPiCamera()
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
	log.Log("", logging.All, "Done.")

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

	log.AddSerialLogService(serLogService)

	log.Logln(TAG, logging.All, "Starting Router (blocking)")
	// router is blocking
	router.Start()
}

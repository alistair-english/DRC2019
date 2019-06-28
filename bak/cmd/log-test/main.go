package main

import (
	"github.com/alistair-english/DRC2019/pkg/logging"
)

func main() {
	log := logging.Logger()
	log.Init()
	log.AddStream("rpi-serial.log", "SERIAL_LOG")
	log.Logln("Hello World! (in default)")
	log.ChangeStream("SERIAL_LOG")
	log.Log("From the Logger! %v \n", 10)

	log.LogOutln("This should log to terminal: %v", 69)
}

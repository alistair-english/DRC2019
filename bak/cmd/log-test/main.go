package main

import (
	"os"
	"time"

	"github.com/alistair-english/DRC2019/pkg/logging"
)

const TAG string = "LOG TEST"

func main() {
	// Create the logger, not matter how much you do this and where, it will be the same logger instance
	log := logging.Logger()
	// Create a default stream and set it
	log.Init()
	// Add the rpi-serial.log file stream with a title SERIAL_LOG
	log.AddStream("rpi-serial.log", "SERIAL_LOG")
	// Log to the default stream
	log.Logln("Hello World! (in default)")
	// Change stream
	log.ChangeStream("SERIAL_LOG")
	// Log to the new stream
	log.Log("From the Logger! %v \n", 10)
	// Log while outputting to stdio
	log.LogOutln("This should log to terminal: %v", 69)
	// List the current streams available
	log.ListStreams()
}

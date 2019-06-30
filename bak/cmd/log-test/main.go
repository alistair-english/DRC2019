package main

import (
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
	log.Logln(TAG, logging.LogStd, "Hello World! (in default)")
	// Change stream
	log.ChangeStream("SERIAL_LOG")
	// Log to the new stream
	log.Log(TAG, logging.LogStd, "From the Logger! %v \n", 10)
	// List the current streams available
	log.ListStreams()
}

package main

import (
	"github.com/alistair-english/DRC2019/pkg/logging"
	"os"
)

func main() {
	log := logging.Logger()
	log.Init()
	f, err := os.OpenFile(os.Getenv("GOPATH")+"/src/github.com/alistair-english/DRC2019/logs/rpi-serial.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetStream(f, "SERIAL_LOG")
	log.Logln("Hello World!")
	log.Log("From the Logger! %v", 10)
}

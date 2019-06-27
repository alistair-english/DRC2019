package main

import (
	"os"
	"time"

	"github.com/alistair-english/DRC2019/pkg/logging"
)

const TAG string = "LOG TEST"

func main() {
	log := logging.Logger()
	log.Init()
	f, err := os.OpenFile(os.Getenv("GOPATH")+"/src/github.com/alistair-english/DRC2019/logs/rpi-serial.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	log.SetStream(f, "SERIAL_LOG")

	for {
		log.Logln(TAG, "Hello World!")
		log.Log(TAG, "From the Logger! %v", 10)
		time.Sleep(1 * time.Second)
	}

}

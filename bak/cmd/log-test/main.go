package main

import "github.com/alistair-english/DRC2019/pkg/logging"

func main() {
	log := logging.Logger()
	log.Init()
	log.Logln("Hello World!")
	log.Log("From the Logger! %v", 10)
}

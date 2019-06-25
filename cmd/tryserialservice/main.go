package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/dummyservice"
	"github.com/alistair-english/DRC2019/pkg/serialservice"
)

func main() {
	router := arch.NewRouter()

	serviceA := dummyservice.NewDummyServiceA()
	serService, _ := serialservice.NewFakeSerial()

	router.Register(serviceA)
	router.Register(serService)

	serviceA.Start()
	serService.Start()

	// router is blocking
	router.Start()
}

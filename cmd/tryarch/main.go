package main

import (
	"github.com/alistair-english/DRC2019/pkg/arch"
	"github.com/alistair-english/DRC2019/pkg/dummyservice"
)

func main() {
	router := arch.NewRouter()

	serviceA := dummyservice.NewDummyServiceA()
	serviceB := dummyservice.NewDummyServiceB()

	router.Register(serviceA)
	router.Register(serviceB)

	serviceA.Start()
	serviceB.Start()

	// router is blocking
	router.Start()
}

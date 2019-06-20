package main

import (
	"fmt"
	"time"
)

func cameraWorker(imageRequests <-chan bool, imageResults chan<- int) {
	for range imageRequests {
		// Wait until channel is empty (don't want to start taking the pic until we are ready to recieve it)
		// for {
		// 	if len(imageResults) == cap(imageResults) {
		// 		break
		// 	}
		// }
		time.Sleep(2 * time.Second)
		imageResults <- 69
	}
}

func initaliseSerial(inputs <-chan int, outputs chan<- int) {
	// Probably spawns some other routines like:
	go func(inputs <-chan int) {
		for i := range inputs {
			// Do something intelligent with the serial
			fmt.Println(i)
		}
	}(inputs)

	go func(outputs chan<- int) {
		for {
			time.Sleep(5 * time.Second)
			outputs <- 420
		}
	}(outputs)
}

func main() {
	// Setup
	// Required Channels
	imageRequests := make(chan bool, 1)
	imageResults := make(chan int, 1)

	serialInputs := make(chan int, 100)
	serialOutputs := make(chan int, 100)

	// Initalise the camera (has an image request channel and image result channel)
	go cameraWorker(imageRequests, imageResults)

	// Initalise the serial (has an input and output channel)
	initaliseSerial(serialInputs, serialOutputs)

	// Put an inital image request on the channel
	imageRequests <- true

	for {
		select {
		case img := <-imageResults:
			// Code to handle an image is ready
			// e.g. process img - this realllyyy needs to be a goroutine otherwise huge block
			x := img / 13

			// Then wack an output on the serial channel
			serialOutputs <- x

			// Then request another image
			imageRequests <- true

		case serMsg := <-serialOutputs:
			// Act on the message
			fmt.Println(serMsg)
		}
	}
}

package main

import (
	"gocv.io/x/gocv"
)

func main() {

	webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for { // foreva
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}

	// - get hsv image
	// - get 'top' and 'bottom' ROI (dunno why -> probs just use one thresholded img)
	// - threshold the image based on colours
	// - blur the thresholds (they used medianBlur)
	// - could do some morphological transformations to make the binary img better
	// - then on these binary imgs find the contours - maybe just find the biggest one or something to use for the line
	// - now they do some 'get centered coord' bs that i cant be bothered to decipher...
	//   so here ^^ once we have the contours, get their bounding boxes, get the center top corners and center bottom corners
	//   of each line, and find a center points (center_top and center_bottom)
	//   - these center_top and center_bottom will give us an idea of the direction of the line and also an idea of how offset
	//     the car is from the center of the track
	// - from this info, we can act upon it and generate some steering and stuff

	// some other notes:
	// - they used a kalman filter to filter their error (their 'offset' from the center)
	// - then they used PID with that filtered val

	// and thats about it. only other thing is need some sort of way to setup the colour thresholds

	// thatll do for now, can go back and look at their object aviodance once we get to that
}

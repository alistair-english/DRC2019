package main

import (
	"image/color"

	"gocv.io/x/gocv"
)

func main() {

	// THIS PROGRAM THRESHOLDS FOR BLUE HSV AND DRAWS BIG OBJECTS ON THE SCREEN (area > 1000)

	webcam, _ := gocv.VideoCaptureDevice(0)
	defer webcam.Close()

	// display windows
	sourceWindow := gocv.NewWindow("Source")
	defer sourceWindow.Close()

	maskWindow := gocv.NewWindow("Mask")
	defer maskWindow.Close()

	// image matricies
	sourceImg := gocv.NewMat()
	defer sourceImg.Close()

	blurredImg := gocv.NewMat()
	defer blurredImg.Close()

	hsvImg := gocv.NewMat()
	defer hsvImg.Close()

	mask := gocv.NewMat()
	defer mask.Close()

	blurMask := gocv.NewMat()
	defer blurMask.Close()

	finalMask := gocv.NewMat()
	defer finalMask.Close()

	webcam.Read(&sourceImg)

	gocv.CvtColor(sourceImg, &hsvImg, gocv.ColorBGRToHSV)

	channels, rows, cols := hsvImg.Channels(), hsvImg.Rows(), hsvImg.Cols()

	// define HSV color upper and lower bound ranges
	lower := gocv.NewMatFromScalar(gocv.NewScalar(90.0, 100.0, 150.0, 0.0), gocv.MatTypeCV8UC3)
	upper := gocv.NewMatFromScalar(gocv.NewScalar(130.0, 255.0, 255.0, 0.0), gocv.MatTypeCV8UC3)

	// split HSV lower bounds into H, S, V channels
	lowerChans := gocv.Split(lower)
	lowerMask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	lowerMaskChans := gocv.Split(lowerMask)

	// split HSV lower bounds into H, S, V channels
	upperChans := gocv.Split(upper)
	upperMask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	upperMaskChans := gocv.Split(upperMask)

	// copy HSV values to upper and lower masks
	for c := 0; c < channels; c++ {
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				lowerMaskChans[c].SetUCharAt(i, j, lowerChans[c].GetUCharAt(0, 0))
				upperMaskChans[c].SetUCharAt(i, j, upperChans[c].GetUCharAt(0, 0))
			}
		}
	}

	gocv.Merge(lowerMaskChans, &lowerMask)
	gocv.Merge(upperMaskChans, &upperMask)

	for { // for3vA

		webcam.Read(&sourceImg)

		// gocv.GaussianBlur(sourceImg, &blurredImg, image.Point{11, 11}, 0, 0, gocv.BorderReflect101)

		gocv.CvtColor(sourceImg, &hsvImg, gocv.ColorBGRToHSV)
		gocv.InRange(hsvImg, lowerMask, upperMask, &mask)

		// gocv.GaussianBlur(mask, &blurMask, image.Point{11, 11}, 0, 0, gocv.BorderReflect101)

		gocv.Merge([]gocv.Mat{mask, mask, mask}, &finalMask)

		// gocv.BitwiseAnd(hsvImg, finalMask, &hsvImg)

		// gocv.CvtColor(hsvImg, &hsvImg, gocv.ColorHSVToBGR)

		contours := gocv.FindContours(mask, gocv.RetrievalTree, gocv.ChainApproxNone)

		for _, contour := range contours {
			if gocv.ContourArea(contour) > 1000 {
				// gocv.DrawContours(&sourceImg, contours, i, color.RGBA{255, 0, 0, 0}, 3)
				rect := gocv.BoundingRect(contour)
				gocv.Rectangle(&sourceImg, rect, color.RGBA{255, 0, 0, 0}, 3)
			}
		}

		maskWindow.IMShow(mask)

		sourceWindow.IMShow(sourceImg)

		sourceWindow.WaitKey(1)
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

package cvhelpers

//TODO: Move this package to internal

import (
	"image"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/gohelpers"
)

const TAG = "CV_HELPERS"

// ReadHSV will read in an image from the supplied VideoCapture and output a HSV file to dst
func ReadHSV(cam *gocv.VideoCapture, dst *gocv.Mat) {
	tempMat := gocv.NewMat()
	defer tempMat.Close()

	cam.Read(&tempMat)
	gocv.CvtColor(tempMat, dst, gocv.ColorBGRToHSV)
}

// NewHSVMask creates a HSV mask that can be used in gocv.InRange
func NewHSVMask(in gocv.Scalar, channels int, rows int, cols int) gocv.Mat {

	source := gocv.NewMatFromScalar(in, gocv.MatTypeCV8UC3)
	defer source.Close()

	output := gocv.NewMat()

	// input is a 1x1x3 with the 3 HSV values we need
	inputChannels := gocv.Split(source)

	// the mask we want is a copy of the input channels but at the size defined by the caller
	mask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	maskChannels := gocv.Split(mask)

	// copy HSV values to the mask
	for c := 0; c < channels; c++ {
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				maskChannels[c].SetUCharAt(row, col, inputChannels[c].GetUCharAt(0, 0))
			}
		}
	}

	gocv.Merge(maskChannels, &output)

	return output
}

// FindLargestContour finds the largest contour in a binary image using some default settings and returns the countour
func FindLargestContour(in gocv.Mat) []image.Point {
	contours := gocv.FindContours(in, gocv.RetrievalTree, gocv.ChainApproxNone)

	if len(contours) == 0 {
		return nil
	}

	var areas []float64
	for _, contour := range contours {
		areas = append(areas, gocv.ContourArea(contour))
	}

	i, _ := gohelpers.MaxFloat64(areas)

	return contours[i]
}

// Thresholds is a struct that contains upper and lower colour bounds in the form of gocv.Scalar
type Thresholds struct {
	Lower gocv.Scalar
	Upper gocv.Scalar
}

// HSVMasks is a struct that contains upper and lower HSV masks that can be generated with cvhelpers.HSVMask
type HSVMasks struct {
	Lower gocv.Mat
	Upper gocv.Mat
}

// HSVObjectGroup describes an object with a name and a HSV masks
type HSVObjectGroup struct {
	Name      string
	Masks     HSVMasks
	NumToFind int
}

// HSVObject describes a found HSV object
type HSVObject struct {
	Countour    []image.Point
	BoundingBox image.Rectangle
	Area        float64
}

// HSVObjectGroupResult describes a group of objects found with cvhelpers.FindHSVObjects
type HSVObjectGroupResult struct {
	Name    string
	Objects []HSVObject
}

// NewHSVObjectGroup creates a new cvhelpers.HSVObjectGroup
func NewHSVObjectGroup(name string, lowerMask gocv.Mat, upperMask gocv.Mat, numToFind int) HSVObjectGroup {
	return HSVObjectGroup{
		name,
		HSVMasks{
			lowerMask,
			upperMask,
		},
		numToFind,
	}
}

// FindHSVObjects finds all HSVObjects from a []cvhelpers.HSVObject in a given image
func FindHSVObjects(img gocv.Mat, objectGroups []HSVObjectGroup) []HSVObjectGroupResult {

	tempMask := gocv.NewMat()
	defer tempMask.Close()

	results := make([]HSVObjectGroupResult, len(objectGroups))

	resultChan := make(chan HSVObjectGroupResult, len(objectGroups))

	// Spawn goroutines to find all the objects
	for _, group := range objectGroups {
		go findHSVObjectGroup(img, group, resultChan)
	}

	// Now wait for the results and pull them out of the channel
	for i := 0; i < len(objectGroups); i++ {
		results[i] = <-resultChan
	}

	return results
}

func findHSVObjectGroup(img gocv.Mat, objectGroup HSVObjectGroup, resultChan chan<- HSVObjectGroupResult) {
	/// Find the object
	mask := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	InRangeBySegments(img, objectGroup.Masks.Lower, objectGroup.Masks.Upper, 2, 2, &mask)

	// Apply any user processing
	// processMask(tempMask, &tempMask)

	// // Find the largest contour
	// contour := FindLargestContour(tempMask)
	// results[i].Countour = contour

	// // Find the bounding box
	// if contour != nil {
	// 	results[i].BoundingBox = gocv.BoundingRect(contour)
	// 	results[i].Area = gocv.ContourArea(contour)
	// }

	// // Copy the name
	// results[i].Name = obj.Name
}

// InRangeBySegments runs InRange gocv function by splitting the image into segments and calculating concurrently
func InRangeBySegments(img, lowerMask, upperMask gocv.Mat, numSegHor, numSegVert int, dst *gocv.Mat) {
	if img.Rows() != dst.Rows() || img.Cols() != dst.Cols() {
		// this should throw error
	}

	// dst must be same size as img and of type MatTypeCV8U
	segWidth := img.Cols() / numSegHor
	segHeight := img.Rows() / numSegVert

	doneChan := make(chan bool, segWidth*segHeight)

	// Spawn a bunch of inrange go routines acting on segments on the image
	for r := 0; r < numSegVert; r++ {
		for c := 0; c < numSegHor; c++ {
			seg := image.Rect(
				c*segWidth,      // xMin
				r*segHeight,     // yMin
				(c+1)*segWidth,  // xMax
				(r+1)*segHeight, // yMax
			)

			sourceSeg := img.Region(seg)
			destSeg := dst.Region(seg)
			go func() {
				gocv.InRange(sourceSeg, lowerMask.Region(seg), upperMask.Region(seg), &destSeg)
				doneChan <- true
			}()
		}
	}

	for r := 0; r < numSegVert; r++ {
		for c := 0; c < numSegHor; c++ {
			<-doneChan
		}
	}

}

package cvhelpers

//TODO: Move this package to internal

import (
	"image"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/pkg/gohelpers"
)

// ReadHSV will read in an image from the supplied VideoCapture and output a HSV file to dst
func ReadHSV(cam *gocv.VideoCapture, dst *gocv.Mat) {
	tempMat := gocv.NewMat()
	defer tempMat.Close()

	cam.Read(&tempMat)
	gocv.CvtColor(tempMat, dst, gocv.ColorBGRToHSV)
}

// HSVMask creates a HSV mask that can be used in gocv.InRange
func HSVMask(in gocv.Scalar, dst *gocv.Mat, channels int, rows int, cols int) {

	source := gocv.NewMatFromScalar(in, gocv.MatTypeCV8UC3)
	defer source.Close()

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

	gocv.Merge(maskChannels, dst)
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
	lower gocv.Scalar
	upper gocv.Scalar
}

// HSVMasks is a struct that contains upper and lower HSV masks that can be generated with cvhelpers.HSVMask
type HSVMasks struct {
	lower gocv.Mat
	upper gocv.Mat
}

// HSVObject describes an object with a name and a HSV masks
type HSVObject struct {
	name  string
	masks HSVMasks
}

type HSVObjectResult struct {
	name        string
	countour    []image.Point
	boundingBox image.Rectangle
}

func NewHSVObject(name string, lowerMask gocv.Mat, upperMask gocv.Mat) HSVObject {
	return HSVObject{
		name,
		HSVMasks{
			lowerMask,
			upperMask,
		},
	}
}

// FindHSVObjects finds all HSVObjects from a []cvhelpers.HSVObject in a given image
func FindHSVObjects(img gocv.Mat, objects []HSVObject) []HSVObjectResult {

	tempMask := gocv.NewMat()
	defer tempMask.Close()

	results := make([]HSVObjectResult, len(objects))

	for i, obj := range objects {
		// Find the object
		gocv.InRange(img, obj.masks.lower, obj.masks.upper, &tempMask)

		// Find the largest contour
		contour := FindLargestContour(tempMask)
		results[i].countour = contour

		// Find the bounding box
		if contour != nil {
			rect := gocv.BoundingRect(contour)
			results[i].boundingBox = rect
		}

		// Copy the name
		results[i].name = obj.name
	}

	return results
}

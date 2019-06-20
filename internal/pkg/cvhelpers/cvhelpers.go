package cvhelpers

//TODO: Move this package to internal

import (
	"image"

	"gocv.io/x/gocv"

	"github.com/alistair-english/DRC2019/internal/pkg/gohelpers"
)

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

// HSVObject describes an object with a name and a HSV masks
type HSVObject struct {
	Name  string
	Masks HSVMasks
}

// HSVObjectResult describes an object found with cvhelpers.FindHSVObjects
type HSVObjectResult struct {
	Name        string
	Countour    []image.Point
	BoundingBox image.Rectangle
	Area        float64
}

// ImageMod is a function to be defined by the user to do some specific processing to a given src image gocv.Mat and output it to dst
type ImageMod func(src gocv.Mat, dst *gocv.Mat)

// NewHSVObject creates a new cvhelpers.HSVObject
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
func FindHSVObjects(img gocv.Mat, objects []HSVObject, processMask ImageMod) []HSVObjectResult {

	tempMask := gocv.NewMat()
	defer tempMask.Close()

	results := make([]HSVObjectResult, len(objects))

	for i, obj := range objects {
		// Find the object
		gocv.InRange(img, obj.Masks.Lower, obj.Masks.Upper, &tempMask)

		// Apply any user processing
		processMask(tempMask, &tempMask)

		// Find the largest contour
		contour := FindLargestContour(tempMask)
		results[i].Countour = contour

		// Find the bounding box
		if contour != nil {
			results[i].BoundingBox = gocv.BoundingRect(contour)
			results[i].Area = gocv.ContourArea(contour)
		}

		// Copy the name
		results[i].Name = obj.Name
	}

	return results
}

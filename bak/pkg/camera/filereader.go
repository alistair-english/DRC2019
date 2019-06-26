package camera

import "gocv.io/x/gocv"

// FileReaderCamera is a camera implementation that serves up images from a video
type FileReaderCamera struct {
	Path    string
	Capture *gocv.VideoCapture
}

// NewFileReaderCamera creates a new camera that will serve up from a video
func NewFileReaderCamera(path string) *FileReaderCamera {
	vid, _ := gocv.VideoCaptureFile(path)
	return &FileReaderCamera{path, vid}
}

// Init from the camera Implementation
func (cam FileReaderCamera) Init() {
}

// RunImagePoller from the camera Implementation
func (cam FileReaderCamera) RunImagePoller(imageRequest <-chan bool, imageResult chan<- bool, outputImg *gocv.Mat) {
	for range imageRequest {
		cam.Capture.Read(outputImg)
		imageResult <- true
	}
}

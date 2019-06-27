package cameraservice

import "gocv.io/x/gocv"

// FileReaderCamera is a camera implementation that serves up images from a video
type FileReaderCamera struct {
	Path    string
	Capture *gocv.VideoCapture
}

// NewFileReaderCameraImplementation creates a new camera implementation that will serve up from a video
func newFileReaderCameraImplementation(path string) (*FileReaderCamera, error) {
	vid, err := gocv.VideoCaptureFile(path)
	if err != nil {
		return nil, err
	}
	return &FileReaderCamera{path, vid}, nil
}

// RunCameraConnection from camera Implementation
func (cam FileReaderCamera) RunCameraConnection(imgRequests <-chan GetImageActionReq) {
	for req := range imgRequests {
		cam.Capture.Read(req.Img)
		select {
		case req.ResponseChannel <- true:
		default:
		}
	}
}

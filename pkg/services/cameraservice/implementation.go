package cameraservice

// Implementation is the interface for a camera
type Implementation interface {
	RunCameraConnection(imgRequests <-chan GetImageActionReq)
}

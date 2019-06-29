package config

import (
	"encoding/json"
	"os"
)

const (
	PINS_CONF_FILE      = "/src/github.com/alistair-english/DRC2019/config/pins.json"
	CV_CONF_FILE        = "/src/github.com/alistair-english/DRC2019/config/cv.json"
	SERIAL_CONF_FILE    = "/src/github.com/alistair-english/DRC2019/config/serial.json"
	PI_CAMERA_CONF_FILE = "/src/github.com/alistair-english/DRC2019/config/picamera.json"
)

// PinConfig is the datatype for control pin information
type PinConfig struct {
	SteeringPin int `json:"steeringPin"`
	DrivePin    int `json:"drivePin"`
	MaxSpeed    int `json:"maxSpeed"`
}

// HSV is the datatype for a hsv value
type HSV struct {
	H float64 `json:"H"`
	S float64 `json:"S"`
	V float64 `json:"V"`
}

//CVConfig is the datatype for the CV information
type CVConfig struct {
	LeftLower  HSV `json:"leftLower"`
	LeftUpper  HSV `json:"leftUpper"`
	RightLower HSV `json:"rightLower"`
	RightUpper HSV `json:"rightUpper"`
}

// SerialConfig is the datatype for the Serial information
type SerialConfig struct {
	Port      string `json:"port"`
	Baud      int    `json:"baud"`
	TimeoutMs int    `json:"timeoutMs"`
}

// PiCameraConfig is the datatype for the Pi Camera configuration
type PiCameraConfig struct {
	Height int `json:"height"`
	Width  int `json:"width"`

	AWB        string `json:"awb"`
	Bitrate    int    `json:"bitrate"`
	Brightness int    `json:"brightness"`
	Contrast   int    `json:"contrast"`
	Exposure   string `json:"exposure"`
	FPS        int    `json:"fps"` // Must be between 2 and 30
	Mode       int    `json:"mode"`
	ROI        string `json:"roi"`
	Sharpness  int    `json:"sharpness"`
	Saturation int    `json:"saturation"`
}

// Possible White Balances (AWB)
//     off           Turn off white balance calculation
//     auto          Automatic mode (default)
//     sun           Sunny mode
//     cloud         Cloudy mode
//     shade         Shaded mode
//     tungsten      Tungsten lighting mode
//     fluorescent   Fluorescent lighting mode
//     incandescent  Incandescent lighting mode
//     flash         Flash mode
//     horizon       Horizon mode

// Possible Exposures:
// 	   off
//     auto          Use automatic exposure mode
//     night         Select setting for night shooting
//     nightpreview
//     backlight     Select setting for back lit subject
//     spotlight
//     sports        Select setting for sports (fast shutter etc)
//     snow          Select setting optimised for snowy scenery
//     beach         Select setting optimised for beach
//     verylong      Select setting for long exposures
//     fixedfps      Constrain fps to a fixed value
//     antishake     Antishake mode
//     fireworks     Select settings

// GetPinConfig gets the pin configuration information from a json file
func GetPinConfig() PinConfig {
	var pins PinConfig
	pinsConfigFile, err := os.Open(os.Getenv("GOPATH") + PINS_CONF_FILE)
	defer pinsConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(pinsConfigFile)
	jsonParser.Decode(&pins)
	return pins
}

// GetCVConfig gets the CV configuration information from a json file
func GetCVConfig() CVConfig {
	var cv CVConfig
	cvConfigFile, err := os.Open(os.Getenv("GOPATH") + CV_CONF_FILE)
	defer cvConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(cvConfigFile)
	jsonParser.Decode(&cv)
	return cv
}

// GetSerialConfig gets the serial configuration information from a json file
func GetSerialConfig() SerialConfig {
	var scf SerialConfig
	serialConfigFile, err := os.Open(os.Getenv("GOPATH") + SERIAL_CONF_FILE)
	defer serialConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(serialConfigFile)
	jsonParser.Decode(&scf)
	return scf
}

// GetPiCameraConfig get the pi camera configuration information from a json file
func GetPiCameraConfig() PiCameraConfig {
	var picf PiCameraConfig
	piConfigFile, err := os.Open(os.Getenv("GOPATH") + PI_CAMERA_CONF_FILE)
	defer piConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(piConfigFile)
	jsonParser.Decode(&picf)
	return picf
}

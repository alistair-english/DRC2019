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
	CONTROL_CONF_FILE   = "/src/github.com/alistair-english/DRC2019/config/control.json"
)

// PiCameraConfig is the datatype for the Pi Camera configuration
type PiCameraConfig struct {
	AWB        string `json:"awb"`
	AWBGains   string `json:"awb_gains"`
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

package config

import (
	"encoding/json"
	"os"
)

// HSV is the datatype for a hsv value
type HSV struct {
	H float64 `json:"H"`
	S float64 `json:"S"`
	V float64 `json:"V"`
}

// Object is the HSV config for an object
type Object struct {
	Name      string  `json:"name"`
	NumToFind int     `json:"numToFind"`
	MinArea   float64 `json:"minArea"`
	UpperMask HSV     `json:"upperMask"`
	LowerMask HSV     `json:"lowerMask"`
}

//CVConfig is the datatype for the CV information
type CVConfig struct {
	ImgWidth    int `json:"imgWidth"`
	ImgHeight   int `json:"imgHeight"`
	ImgChannels int `json:"imgChannels"`

	Objects []Object `json:"objects"`
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

package config

import (
	"encoding/json"
	"os"
)

const (
	PINS_CONF_FILE = "/src/github.com/alistair-english/DRC2019/config/pins.json"
	CV_CONF_FILE   = "/src/github.com/alistair-english/DRC2019/config/cv.json"
)

type PinConfig struct {
	SteeringPin int `json:"steeringPin"`
	DrivePin    int `json:"drivePin"`
	MaxSpeed    int `json:"maxSpeed"`
}

type HSV struct {
	H float64 `json:"H"`
	S float64 `json:"S"`
	V float64 `json:"V"`
}

type CVConfig struct {
	LeftLower  HSV `json:"leftLower"`
	LeftUpper  HSV `json:"leftUpper"`
	RightLower HSV `json:"rightLower"`
	RightUpper HSV `json:"rightUpper"`
}

func GetPinsConfig() PinConfig {
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

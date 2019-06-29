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

//CVConfig is the datatype for the CV information
type CVConfig struct {
	LeftLower  HSV `json:"leftLower"`
	LeftUpper  HSV `json:"leftUpper"`
	RightLower HSV `json:"rightLower"`
	RightUpper HSV `json:"rightUpper"`
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

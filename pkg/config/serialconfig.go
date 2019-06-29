package config

import (
	"encoding/json"
	"os"
)

// SerialConfig is the datatype for the Serial information
type SerialConfig struct {
	Port      string `json:"port"`
	Baud      int    `json:"baud"`
	TimeoutMs int    `json:"timeoutMs"`
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

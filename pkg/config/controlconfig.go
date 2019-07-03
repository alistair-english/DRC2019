package config

import (
	"encoding/json"
	"os"
)

// PID is a pid struct for pulling out of the JSON file
type PID struct {
	P float64 `json:"P"`
	I float64 `json:"I"`
	D float64 `json:"D"`
}

// ControlPID contains the PID config for control
type ControlPID struct {
	Pid PID `json:"controlPid"`
}

// BoundaryPID contains the PID config for Boundary avoid
type BoundaryPID struct {
	Pid PID `json:"boundaryPid"`
}

// OpponentPID contains the PID config for dodging other boyos
type OpponentPID struct {
	Pid PID `json:"opponentPid"`
}

type ObstacleXExclusion struct {
	ObstacleXExclusion int `json:"obstacleXExclusion"`
}

// GetControlPIDConfig returns a control pid struct populated from the config file
func GetControlPIDConfig() ControlPID {
	var ctrlPid ControlPID
	controlConfigFile, err := os.Open(os.Getenv("GOPATH") + CONTROL_CONF_FILE)
	defer controlConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(controlConfigFile)
	jsonParser.Decode(&ctrlPid)
	return ctrlPid
}

// GetBoundaryPIDConfig returns a boundary pid struct populated from the config file
func GetBoundaryPIDConfig() BoundaryPID {
	var boundPid BoundaryPID
	controlConfigFile, err := os.Open(os.Getenv("GOPATH") + CONTROL_CONF_FILE)
	defer controlConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(controlConfigFile)
	jsonParser.Decode(&boundPid)
	return boundPid
}

// GetOpponentPIDConfig returns a opponent pid struct populated from the config file
func GetOpponentPIDConfig() OpponentPID {
	var oppPid OpponentPID
	controlConfigFile, err := os.Open(os.Getenv("GOPATH") + CONTROL_CONF_FILE)
	defer controlConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(controlConfigFile)
	jsonParser.Decode(&oppPid)
	return oppPid
}

func GetObstacleXExclusion() ObstacleXExclusion {
	var obstExc ObstacleXExclusion
	controlConfigFile, err := os.Open(os.Getenv("GOPATH") + CONTROL_CONF_FILE)
	defer controlConfigFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(controlConfigFile)
	jsonParser.Decode(&obstExc)
	return obstExc
}

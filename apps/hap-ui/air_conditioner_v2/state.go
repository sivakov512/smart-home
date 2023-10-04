package airconditionerv2

import (
	"encoding/json"
)

type Mode string

const (
	Cool Mode = "cool"
	Heat      = "heat"
)

type State struct {
	IsActive           bool    `json:"is_active"`
	Mode               Mode    `json:"mode"`
	TargetTemperature  float64 `json:"target_temperature"`
	CurrentTemperature float64 `json:"current_temperature"`
}

func NewState() *State {
	return &State{
		Mode: Cool,
	}
}

func (s *State) Serialize() []byte {
	c, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return c
}

func DeserializeState(i []byte, s *State) {
	json.Unmarshal(i, &s)
}
